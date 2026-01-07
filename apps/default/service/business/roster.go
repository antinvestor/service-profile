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
	// Step 1: Deduplicate contacts
	detailSet := rb.deduplicateContacts(batch)
	detailList := rb.createDetailList(detailSet)

	// Step 2: Get existing contacts
	existingContactMap, err := rb.contactBusiness.GetByDetailMap(ctx, detailList...)
	if err != nil {
		return nil, err
	}

	// Step 3: Find missing contacts and create them
	missingContacts := rb.findMissingContacts(batch, existingContactMap)
	newContactsMap := rb.batchCreateContacts(ctx, missingContacts)

	// Step 4: Create unified contact map
	allContacts := rb.createUnifiedContactMap(existingContactMap, newContactsMap)
	contactDetails := rb.buildContactDetails(batch, allContacts)

	// Step 5: Get existing rosters
	existingRosters, err := rb.rosterRepository.GetByContactIDsAndProfileID(ctx, contactDetails, profileID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	existingRosterMap := rb.createRosterLookupMap(existingRosters)

	// Step 6: Create missing rosters
	rostersToCreate := rb.findRostersToCreate(batch, allContacts, existingRosterMap, profileID)
	if len(rostersToCreate) > 0 {
		if createErr := rb.batchCreateRosters(ctx, rostersToCreate); createErr != nil {
			return nil, createErr
		}
		// Add newly created rosters to existing map
		for _, roster := range rostersToCreate {
			existingRosterMap[roster.ContactID] = roster
		}
	}

	// Step 7: Build final result preserving input order
	return rb.buildRosterObjects(batch, allContacts, existingRosterMap), nil
}

// deduplicateContacts creates a set of unique contact details from the batch.
func (rb *rosterBusiness) deduplicateContacts(batch []*profilev1.RawContact) map[string]struct{} {
	detailSet := make(map[string]struct{}, len(batch))
	for _, newRoster := range batch {
		detail := newRoster.GetContact()
		detailSet[detail] = struct{}{}
	}
	return detailSet
}

// createDetailList converts a set to a slice of contact details.
func (rb *rosterBusiness) createDetailList(detailSet map[string]struct{}) []string {
	detailList := make([]string, 0, len(detailSet))
	for detail := range detailSet {
		detailList = append(detailList, detail)
	}
	return detailList
}

// findMissingContacts identifies contacts that don't exist yet.
func (rb *rosterBusiness) findMissingContacts(
	batch []*profilev1.RawContact,
	existingContactMap map[string]*models.Contact,
) []string {
	missingContacts := make([]string, 0, len(batch))
	for _, contact := range batch {
		if _, exists := existingContactMap[contact.GetContact()]; !exists {
			missingContacts = append(missingContacts, contact.GetContact())
		}
	}
	return missingContacts
}

// createUnifiedContactMap combines existing and new contacts into a single map.
func (rb *rosterBusiness) createUnifiedContactMap(
	existingContactMap map[string]*models.Contact,
	newContactsMap map[string]*models.Contact,
) map[string]*models.Contact {
	allContacts := make(map[string]*models.Contact, len(existingContactMap)+len(newContactsMap))
	// Add existing contacts
	for detail, contact := range existingContactMap {
		allContacts[detail] = contact
	}
	// Add new contacts
	for detail, contact := range newContactsMap {
		allContacts[detail] = contact
	}
	return allContacts
}

// buildContactDetails creates a list of contact IDs preserving input order.
func (rb *rosterBusiness) buildContactDetails(
	batch []*profilev1.RawContact,
	allContacts map[string]*models.Contact,
) []string {
	contactDetails := make([]string, 0, len(batch))
	for _, rawContact := range batch {
		detail := rawContact.GetContact()
		if contact, exists := allContacts[detail]; exists {
			contactDetails = append(contactDetails, contact.GetID())
		}
	}
	return contactDetails
}

// createRosterLookupMap creates a map for quick roster lookup by contact ID.
func (rb *rosterBusiness) createRosterLookupMap(existingRosters []*models.Roster) map[string]*models.Roster {
	existingRosterMap := make(map[string]*models.Roster, len(existingRosters))
	for _, roster := range existingRosters {
		existingRosterMap[roster.ContactID] = roster
	}
	return existingRosterMap
}

// findRostersToCreate identifies rosters that need to be created.
func (rb *rosterBusiness) findRostersToCreate(
	batch []*profilev1.RawContact,
	allContacts map[string]*models.Contact,
	existingRosterMap map[string]*models.Roster,
	profileID string,
) []*models.Roster {
	rostersToCreate := make([]*models.Roster, 0, len(batch))
	processedContacts := make(map[string]bool, len(batch))

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

		if _, rosterExists := existingRosterMap[contact.GetID()]; !rosterExists {
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
	return rostersToCreate
}

// buildRosterObjects creates the final roster objects preserving input order.
func (rb *rosterBusiness) buildRosterObjects(
	batch []*profilev1.RawContact,
	allContacts map[string]*models.Contact,
	existingRosterMap map[string]*models.Roster,
) []*profilev1.RosterObject {
	rosterObjectList := make([]*profilev1.RosterObject, 0, len(batch))

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
				return nil // Will be handled by the caller
			}
			rosterObjectList = append(rosterObjectList, rosterObj)
		}
	}
	return rosterObjectList
}

// batchCreateContacts creates multiple contacts in a batch for better performance.
func (rb *rosterBusiness) batchCreateContacts(
	ctx context.Context,
	newContacts []string,
) map[string]*models.Contact {
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

	return contacts
}

// batchCreateRosters creates multiple rosters in a batch for better performance.
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
