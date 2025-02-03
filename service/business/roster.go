package business

import (
	"context"
	"fmt"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
)

type RosterBusiness interface {
	Search(ctx context.Context, request *profilev1.SearchRosterRequest) (frame.JobResultPipe[[]*models.Roster], error)
	GetByID(ctx context.Context, rosterID string) (*models.Roster, error)
	CreateRoster(ctx context.Context, request *profilev1.AddRosterRequest) ([]*profilev1.RosterObject, error)
	RemoveRoster(ctx context.Context, rosterID string) (*profilev1.RosterObject, error)

	ToApi(ctx context.Context, roster *models.Roster) (*profilev1.RosterObject, error)
}

func NewRosterBusiness(ctx context.Context, service *frame.Service) RosterBusiness {
	rosterRepo := repository.NewRosterRepository(service)
	return &rosterBusiness{
		service:          service,
		rosterRepository: rosterRepo,
		contactBusiness:  NewContactBusiness(ctx, service),
	}
}

type rosterBusiness struct {
	service          *frame.Service
	rosterRepository repository.RosterRepository
	contactBusiness  ContactBusiness
}

func (rb *rosterBusiness) ToApi(ctx context.Context, roster *models.Roster) (*profilev1.RosterObject, error) {

	contact := roster.Contact
	contactObject, err := rb.contactBusiness.ToAPI(ctx, contact, true)
	if err != nil {
		return nil, err
	}

	return &profilev1.RosterObject{
		Id:        roster.GetID(),
		ProfileId: roster.ProfileID,
		Contact:   contactObject,
		Extra:     frame.DBPropertiesToMap(roster.Properties),
	}, nil
}

func (rb *rosterBusiness) GetByID(ctx context.Context, rosterID string) (*models.Roster, error) {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	return rb.rosterRepository.GetByID(ctx, rosterID)
}

func (rb *rosterBusiness) Search(ctx context.Context,
	request *profilev1.SearchRosterRequest) (frame.JobResultPipe[[]*models.Roster], error) {

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	profileID := request.GetProfileId()
	claims := frame.ClaimsFromContext(ctx)
	if claims != nil {
		if claims.ServiceName() == "" || profileID == "" {
			profileID, _ = claims.GetSubject()
		}
	}

	query, err := repository.NewSearchQuery(ctx, profileID, request.GetQuery(), request.GetProperties(), request.GetStartDate(), request.GetEndDate(), int(request.GetCount()), int(request.GetPage()))
	if err != nil {
		return nil, err
	}

	return rb.rosterRepository.Search(ctx, query)
}

func (rb *rosterBusiness) CreateRoster(ctx context.Context, request *profilev1.AddRosterRequest) ([]*profilev1.RosterObject, error) {

	claims := frame.ClaimsFromContext(ctx)

	if claims == nil {
		return nil, fmt.Errorf("no claims found in context")
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

		contact, err = rb.contactBusiness.CreateContact(ctx, newRoster.GetContact(), newRoster.GetExtras())
		if err != nil {
			return nil, err
		}

		roster, err = rb.rosterRepository.GetByContactAndProfileID(ctx, profileID, contact.GetID())
		if err != nil {
			if !frame.DBErrorIsRecordNotFound(err) {
				return nil, err
			}

			roster = &models.Roster{
				ProfileID:  profileID,
				ContactID:  contact.GetID(),
				Contact:    contact,
				Properties: frame.DBPropertiesFromMap(newRoster.GetExtras()),
			}

			roster, err = rb.rosterRepository.Save(ctx, roster)
			if err != nil {
				return nil, err
			}
		}
		rosterObject, err := rb.ToApi(ctx, roster)
		if err != nil {
			return nil, err
		}
		rosterObjectList = append(rosterObjectList, rosterObject)
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

	return rb.ToApi(ctx, roster)
}
