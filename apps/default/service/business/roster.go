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
	batchSize := len(batch)

	// Pre-allocate all slices with exact capacity to avoid reallocations
	detailSet := make(map[string]struct{}, batchSize)

	for _, newRoster := range batch {
		detail := newRoster.GetContact()
		detailSet[detail] = struct{}{}
	}

	// Create detailList from deduplicated set
	detailList := make([]string, 0, len(detailSet))
	for detail := range detailSet {
		detailList = append(detailList, detail)
	}

	batchSize = len(detailList)

	// Step 1: Bulk check existing contacts using GetByDetailMap
	existingContactMap, err := rb.contactBusiness.GetByDetailMap(ctx, detailList...)
	if err != nil {
		return nil, err
	}

	// Step 2: Separate existing contacts from new ones and prepare batch creation
	missingContacts := make([]string, 0, batchSize)

	for _, contact := range batch {
		if _, exists := existingContactMap[contact.GetContact()]; !exists {
			missingContacts = append(missingContacts, contact.GetContact())
		}
	}

	// Step 3: Batch create new contacts if any
	newContactsMap, createErr := rb.batchCreateContacts(ctx, missingContacts)
	if createErr != nil {
		return nil, createErr
	}

	// Step 4: Create unified contact map for single lookup and build contactDetails in one pass
	allContacts := make(map[string]*models.Contact, len(existingContactMap)+len(newContactsMap))

	// Add existing contacts to unified map
	for detail, contact := range existingContactMap {
		allContacts[detail] = contact
	}

	// Add new contacts to unified map
	for detail, contact := range newContactsMap {
		allContacts[detail] = contact
	}

	// Build contactDetails in single pass while preserving input order
	contactDetails := make([]string, 0, len(batch))

	for _, rawContact := range batch {
		detail := rawContact.GetContact()
		if contact, exists := allContacts[detail]; exists {
			contactDetails = append(contactDetails, contact.GetID())
		}
	}

	// Step 4: Batch lookup existing rosters for all contact IDs
	existingRosters, err := rb.rosterRepository.GetByContactIDsAndProfileID(ctx, contactDetails, profileID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}

	// Create map for quick lookup of existing rosters
	existingRosterMap := make(map[string]*models.Roster, len(existingRosters))
	for _, roster := range existingRosters {
		existingRosterMap[roster.ContactID] = roster
	}

	// Step 5: Batch create rosters that don't exist - iterate through unique contacts only
	rostersToCreate := make([]*models.Roster, 0, batchSize)
	processedContacts := make(map[string]bool, batchSize) // Track which contacts we've processed

	for _, rawCtc := range batch {
		detail := rawCtc.GetContact()
		// Skip if we've already processed this contact
		if _, processed := processedContacts[detail]; processed {
			continue
		}
		processedContacts[detail] = true

		// Single lookup using unified contact map
		contact, exists := allContacts[detail]
		if !exists {
			continue
		}

		if _, exists := existingRosterMap[contact.GetID()]; !exists {
			rosterExtras := data.JSONMap{}
			roster := &models.Roster{
				ProfileID:  profileID,
				ContactID:  contact.GetID(),
				Contact:    contact,
				Properties: rosterExtras.FromProtoStruct(rawCtc.GetExtras()),
			}
			rostersToCreate = append(rostersToCreate, roster)
		}
	}

	// Batch create rosters if any need to be created
	if len(rostersToCreate) > 0 {
		if createErr := rb.batchCreateRosters(ctx, rostersToCreate); createErr != nil {
			return nil, createErr
		}

		// Add newly created rosters to existing map
		for _, roster := range rostersToCreate {
			existingRosterMap[roster.ContactID] = roster
		}
	}

	// Step 6: Build final result preserving input order with optimized API conversion
	rosterObjectList := make([]*profilev1.RosterObject, 0, batchSize)

	for _, rawContact := range batch {
		detail := rawContact.GetContact()
		// Single lookup to get contact
		contact, exists := allContacts[detail]
		if !exists {
			continue
		}

		// Direct lookup for roster using contact ID
		roster, exists := existingRosterMap[contact.GetID()]
		if !exists {
			continue
		}

		// Add a check to ensure roster is not nil before calling ToAPI
		if roster != nil {
			rosterObj, apiErr := roster.ToAPI(rb.dek)
			if apiErr != nil {
				return nil, apiErr
			}
			rosterObjectList = append(rosterObjectList, rosterObj)
		}
	}

	return rosterObjectList, nil
}

// batchCreateContacts creates multiple contacts in a batch for better performance
func (rb *rosterBusiness) batchCreateContacts(ctx context.Context, newContacts []string) (map[string]*models.Contact, error) {
	contacts := make(map[string]*models.Contact, len(newContacts))

	for _, contactDetail := range newContacts {

		contact, createErr := rb.contactBusiness.CreateContact(
			ctx,
			contactDetail,
			data.JSONMap{},
		)
		if createErr != nil {
			util.Log(ctx).WithField("detail", contactDetail).Error("Failed to create contact", "error", createErr)
			continue
		}

		contacts[contactDetail] = contact
	}

	return contacts, nil
}

// batchCreateRosters creates multiple rosters in a batch for better performance
func (rb *rosterBusiness) batchCreateRosters(ctx context.Context, rosters []*models.Roster) error {
	createErr := rb.rosterRepository.BulkCreate(ctx, rosters)
	if createErr != nil {
		return data.ErrorConvertToAPI(createErr)
	}
	return nil
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
