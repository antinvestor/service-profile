package business

import (
	"context"
	"errors"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
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
	) (workerpool.JobResultPipe[[]*models.Roster], *connect.Error)
	GetByID(ctx context.Context, rosterID string) (*models.Roster, *connect.Error)
	CreateRoster(ctx context.Context, request *profilev1.AddRosterRequest) ([]*profilev1.RosterObject, *connect.Error)
	RemoveRoster(ctx context.Context, rosterID string) (*profilev1.RosterObject, *connect.Error)
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

func (rb *rosterBusiness) GetByID(ctx context.Context, rosterID string) (*models.Roster, *connect.Error) {
	roster, err := rb.rosterRepository.GetByID(ctx, rosterID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return roster, nil
}

func (rb *rosterBusiness) Search(ctx context.Context,
	request *profilev1.SearchRosterRequest) (workerpool.JobResultPipe[[]*models.Roster], *connect.Error) {
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

	result, err := rb.rosterRepository.Search(ctx, query)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return result, nil
}

func (rb *rosterBusiness) CreateRoster(
	ctx context.Context,
	request *profilev1.AddRosterRequest,
) ([]*profilev1.RosterObject, *connect.Error) {
	claims := security.ClaimsFromContext(ctx)

	if claims == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no claims found in context"))
	}

	profileID, subjectErr := claims.GetSubject()
	if subjectErr != nil {
		return nil, data.ErrorConvertToAPI(subjectErr)
	}

	var rosterObjectList []*profilev1.RosterObject
	newRosterList := request.GetData()
	for _, newRoster := range newRosterList {
		var contact *models.Contact
		var roster *models.Roster

		requestExtras := data.JSONMap{}

		var createContactErr *connect.Error
		contact, createContactErr = rb.contactBusiness.CreateContact(
			ctx,
			newRoster.GetContact(),
			requestExtras.FromProtoStruct(newRoster.GetExtras()),
		)
		if createContactErr != nil {
			return nil, createContactErr
		}

		var getRosterErr error
		roster, getRosterErr = rb.rosterRepository.GetByContactAndProfileID(ctx, profileID, contact.GetID())
		if getRosterErr != nil {
			if !data.ErrorIsNoRows(getRosterErr) {
				return nil, data.ErrorConvertToAPI(getRosterErr)
			}

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

		rosterObjectList = append(rosterObjectList, roster.ToAPI())
	}

	return rosterObjectList, nil
}

func (rb *rosterBusiness) RemoveRoster(ctx context.Context, rosterID string) (*profilev1.RosterObject, *connect.Error) {
	roster, err := rb.rosterRepository.GetByID(ctx, rosterID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}

	err = rb.rosterRepository.Delete(ctx, roster.GetID())
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}

	return roster.ToAPI(), nil
}
