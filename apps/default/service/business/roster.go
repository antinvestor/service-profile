package business

import (
	"context"
	"errors"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"

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
	contactBusiness ContactBusiness,
	rosterRepo repository.RosterRepository,
) RosterBusiness {
	return &rosterBusiness{
		rosterRepository: rosterRepo,
		contactBusiness:  contactBusiness,
	}
}

type rosterBusiness struct {
	rosterRepository repository.RosterRepository
	contactBusiness  ContactBusiness
}

func (rb *rosterBusiness) GetByID(ctx context.Context, rosterID string) (*models.Roster, error) {
	return rb.rosterRepository.GetByID(ctx, rosterID)
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
		orSearchFilter["SIMILARITY(contacts.detail,?) > 0"] = request.GetQuery()
		orSearchFilter["rosters.searchable  @@ websearch_to_tsquery( 'english', ?) "] = request.GetQuery()
	}

	query := data.NewSearchQuery(

		data.WithSearchLimit(int(request.GetCount())),
		data.WithSearchOffset(int(request.GetPage())),
		data.WithSearchFiltersAndByValue(map[string]any{"rosters.profile_id": profileID}),
		data.WithSearchFiltersOrByValue(orSearchFilter))

	return rb.rosterRepository.Search(ctx, query)
}

func (rb *rosterBusiness) CreateRoster(
	ctx context.Context,
	request *profilev1.AddRosterRequest,
) ([]*profilev1.RosterObject, error) {
	claims := security.ClaimsFromContext(ctx)

	if claims == nil {
		return nil, errors.New("no claims found in context")
	}

	profileID, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	var rosterObjectList []*profilev1.RosterObject
	newRosterList := request.GetData()
	for _, newRoster := range newRosterList {
		var contact *models.Contact
		var roster *models.Roster

		requestExtras := data.JSONMap{}

		contact, err = rb.contactBusiness.CreateContact(
			ctx,
			newRoster.GetContact(),
			requestExtras.FromProtoStruct(newRoster.GetExtras()),
		)
		if err != nil {
			return nil, err
		}

		roster, err = rb.rosterRepository.GetByContactAndProfileID(ctx, profileID, contact.GetID())
		if err != nil {
			if !data.ErrorIsNoRows(err) {
				return nil, err
			}

			rosterExtras := data.JSONMap{}

			roster = &models.Roster{
				ProfileID:  profileID,
				ContactID:  contact.GetID(),
				Contact:    contact,
				Properties: rosterExtras.FromProtoStruct(newRoster.GetExtras()),
			}

			err = rb.rosterRepository.Create(ctx, roster)
			if err != nil {
				return nil, err
			}
		}

		rosterObjectList = append(rosterObjectList, roster.ToAPI())
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

	return roster.ToAPI(), nil
}
