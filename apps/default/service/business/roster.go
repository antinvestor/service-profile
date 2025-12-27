package business

import (
	"context"
	"errors"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

type RosterBusiness interface {
	Search(
		ctx context.Context,
		request *profilev1.SearchRosterRequest,
	) (workerpool.JobResultPipe[[]*models.Roster], error)
	GetByID(ctx context.Context, rosterID string) (*models.Roster, error)
	CreateRoster(ctx context.Context, request *profilev1.AddRosterRequest) ([]*profilev1.RosterObject, error)
	RemoveRoster(ctx context.Context, rosterID string) (*profilev1.RosterObject, error)
}

func NewRosterBusiness(
	_ context.Context,
	cfg *config.ProfileConfig, dek *config.DEK,
	contactBusiness ContactBusiness,
	rosterRepo repository.RosterRepository,
) RosterBusiness {
	return &rosterBusiness{
		cfg:              cfg,
		dek:              dek,
		rosterRepository: rosterRepo,
		contactBusiness:  contactBusiness,
	}
}

type rosterBusiness struct {
	cfg              *config.ProfileConfig
	dek              *config.DEK
	rosterRepository repository.RosterRepository
	contactBusiness  ContactBusiness
}

func (rb *rosterBusiness) GetByID(ctx context.Context, rosterID string) (*models.Roster, error) {
	roster, err := rb.rosterRepository.GetByID(ctx, rosterID)
	if err != nil {
		return nil, err
	}
	return roster, nil
}

func (rb *rosterBusiness) Search(ctx context.Context,
	request *profilev1.SearchRosterRequest) (workerpool.JobResultPipe[[]*models.Roster], error) {
	profileID := request.GetProfileId()
	claims := security.ClaimsFromContext(ctx)
	if claims != nil {
		if claims.GetServiceName() == "" || profileID == "" {
			profileID, _ = claims.GetSubject()
		}
	}

	var orSearchFilter = make(map[string]any)
	if request.GetQuery() != "" {
		orSearchFilter["contacts.look_up_token = ?"] = util.ComputeLookupToken(
			rb.dek.LookUpKey,
			Normalize(ctx, request.GetQuery()),
		)
		orSearchFilter["rosters.searchable  @@ websearch_to_tsquery( 'english', ?) "] = request.GetQuery()
	}

	query := data.NewSearchQuery(

		data.WithSearchLimit(int(request.GetCount())),
		data.WithSearchOffset(int(request.GetPage())),
		data.WithSearchFiltersAndByValue(map[string]any{"rosters.profile_id": profileID}),
		data.WithSearchFiltersOrByValue(orSearchFilter))

	result, err := rb.rosterRepository.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rb *rosterBusiness) CreateRoster(
	ctx context.Context,
	request *profilev1.AddRosterRequest,
) ([]*profilev1.RosterObject, error) {
	claims := security.ClaimsFromContext(ctx)

	if claims == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no claims found in context"))
	}

	profileID, err := claims.GetSubject()
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}

	newRosterList := request.GetData()
	if len(newRosterList) == 0 {
		return []*profilev1.RosterObject{}, nil
	}

	// Pre-allocate result slice for better memory efficiency
	rosterObjectList := make([]*profilev1.RosterObject, 0, len(newRosterList))

	// Batch size for processing - optimized for database performance
	const batchSize = 50

	// Process rosters in batches to optimize database operations
	for i := 0; i < len(newRosterList); i += batchSize {
		end := i + batchSize
		if end > len(newRosterList) {
			end = len(newRosterList)
		}

		batch := newRosterList[i:end]
		batchResults, batchErr := rb.processRosterBatch(ctx, profileID, batch)
		if batchErr != nil {
			return nil, batchErr
		}

		rosterObjectList = append(rosterObjectList, batchResults...)
	}

	return rosterObjectList, nil
}

// processRosterBatch processes a batch of roster items efficiently.
func (rb *rosterBusiness) processRosterBatch(
	ctx context.Context,
	profileID string,
	batch []*profilev1.RawContact,
) ([]*profilev1.RosterObject, error) {
	// Step 1: Extract all contact details from the batch
	detailList := make([]string, 0, len(batch))
	for _, newRoster := range batch {
		detailList = append(detailList, newRoster.GetContact())
	}

	// Step 2: Bulk check existing contacts using GetByDetailMap
	existingContactMap, err := rb.contactBusiness.GetByDetailMap(ctx, detailList...)
	if err != nil {
		return nil, err
	}

	// Step 3: Create only contacts that don't exist
	contacts := make([]*models.Contact, 0, len(batch))
	contactDetails := make([]string, 0, len(batch))

	for _, newRoster := range batch {
		detail := newRoster.GetContact()

		// Check if contact already exists
		if existingContact, exists := existingContactMap[detail]; exists {
			// Use existing contact
			contacts = append(contacts, existingContact)
			contactDetails = append(contactDetails, existingContact.GetID())
		} else {
			// Create new contact
			requestExtras := data.JSONMap{}
			contact, createErr := rb.contactBusiness.CreateContact(
				ctx,
				newRoster.GetContact(),
				requestExtras.FromProtoStruct(newRoster.GetExtras()),
			)
			if createErr != nil {
				return nil, createErr
			}
			contacts = append(contacts, contact)
			contactDetails = append(contactDetails, contact.GetID())
		}
	}

	// Step 4: Batch lookup existing rosters for all contact IDs
	existingRosters, err := rb.rosterRepository.GetByContactIDsAndProfileID(ctx, contactDetails, profileID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}

	// Create map for quick lookup of existing rosters
	existingRosterMap := make(map[string]*models.Roster)
	for _, roster := range existingRosters {
		existingRosterMap[roster.ContactID] = roster
	}

	// Step 5: Process each contact and create roster if needed
	rosterObjectList := make([]*profilev1.RosterObject, 0, len(batch))
	for i, contact := range contacts {
		newRoster := batch[i]

		roster, exists := existingRosterMap[contact.GetID()]
		if !exists {
			// Create new roster
			rosterExtras := data.JSONMap{}
			roster = &models.Roster{
				ProfileID:  profileID,
				ContactID:  contact.GetID(),
				Contact:    contact,
				Properties: rosterExtras.FromProtoStruct(newRoster.GetExtras()),
			}

			var createRosterErr = rb.rosterRepository.Create(ctx, roster)
			if createRosterErr != nil {
				return nil, data.ErrorConvertToAPI(createRosterErr)
			}
		}

		rosterObj, apiErr := roster.ToAPI(rb.dek)
		if apiErr != nil {
			return nil, apiErr
		}
		rosterObjectList = append(rosterObjectList, rosterObj)
	}

	return rosterObjectList, nil
}

func (rb *rosterBusiness) RemoveRoster(ctx context.Context, rosterID string) (*profilev1.RosterObject, error) {
	roster, err := rb.rosterRepository.GetByID(ctx, rosterID)
	if err != nil {
		return nil, err
	}

	err = rb.rosterRepository.Delete(ctx, roster.GetID())
	if err != nil {
		return nil, err
	}

	rosterObj, err := roster.ToAPI(rb.dek)
	if err != nil {
		return nil, err
	}
	return rosterObj, nil
}
