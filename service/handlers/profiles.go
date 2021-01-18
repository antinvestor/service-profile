package handlers

import (
	"context"
	"errors"
	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/models"
	"github.com/antinvestor/service-profile/service"
	"github.com/pitabwire/frame"

	"strings"
)

type ProfileServer struct {
	Service *frame.Service
	NotificationCli *napi.NotificationClient

	papi.ProfileServiceServer
}

func (ps *ProfileServer) getProfileByID(ctx context.Context, profileID string, ) (*papi.ProfileObject, error) {
	p := models.Profile{}
	p.ID = profileID

	ps.Service.DB(ctx, true).First(&p)

	return p.ToObject(ps.Service.DB(ctx, true))
}

func (ps *ProfileServer) GetByID(ctx context.Context,
	request *papi.ProfileIDRequest, ) (*papi.ProfileObject, error) {
	profileID := strings.TrimSpace(request.GetID())
	return ps.getProfileByID(ctx, profileID)
}

func (ps *ProfileServer) Search(request *papi.ProfileSearchRequest,
	stream papi.ProfileService_SearchServer, ) error {

	var profiles []models.Profile

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//// creating WHERE clause to query by properties JSONB
	//scope := ps.Service.DB(ctx, true).New()
	//for _, property := range request.GetProperties() {
	//	column := fmt.Sprintf("properties->>'%s'", EscapeColumnName(property), )
	//	scope = scope.Or(column+" LIKE ?", "%"+request.GetQuery()+"%")
	//}
	//
	//scope.Find(&profiles)

	for _, p := range profiles {
		profileObject, err := p.ToObject(ps.Service.DB(ctx, true))
		if err != nil {
			return err
		}

		err = stream.Send(profileObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *ProfileServer) Merge(ctx context.Context, request *papi.ProfileMergeRequest, ) (
	*papi.ProfileObject, error) {

	var target models.Profile
	var merging models.Profile

	target.ID = request.GetID()

	if err := target.GetByID(ps.Service.DB(ctx, true)); err != nil {
		return nil, err
	}

	merging.ID = request.GetMergeID()

	if err := merging.GetByID(ps.Service.DB(ctx, true)); err != nil {
		return nil, err
	}

	for key, value := range merging.Properties {

		existingValue := target.Properties[key]

		if existingValue == value {
			continue
		}

		target.Properties[key] = value
	}

	target.UpdateProperties(ps.Service.DB(ctx, false), merging.Properties)

	return target.ToObject(ps.Service.DB(ctx, true))
}

func (ps *ProfileServer) Create(ctx context.Context, request *papi.ProfileCreateRequest, ) (
	*papi.ProfileObject, error) {

	properties := make(map[string]interface{})

	for key, value := range request.GetProperties() {
		properties[key] = value
	}

	contactDetail := strings.TrimSpace(request.GetContact())

	if contactDetail == "" {
		return nil, service.ErrorContactDetailsNotValid
	}

	p := models.Profile{}

	contact := models.Contact{Detail: contactDetail}

	err := contact.GetByDetail(ps.Service.DB(ctx, true))

	if err != nil {
		if  !errors.Is(service.ErrorContactDoesNotExist, err) {
			return nil, err
		}


		err := p.Create(ps.Service.DB(ctx, false), request.GetType(), properties)
		if err != nil {
			return nil, err
		}

		contact, err := createContact( ctx, ps.Service, ps.NotificationCli, p.ID, contactDetail)
		if err != nil && contact == nil{
			return nil, err
		}

	}else{

		p.ID = contact.ProfileID

		err = ps.Service.DB(ctx, true).First(p).Error
		if err != nil {
			return nil, err
		}

		err = p.UpdateProperties(ps.Service.DB(ctx, false), properties)
		if err != nil {
			return nil, err
		}
	}

	return ps.getProfileByID(ctx, p.ID)
}

func (ps *ProfileServer) Update(
	ctx context.Context,
	request *papi.ProfileUpdateRequest,
) (*papi.ProfileObject, error) {
	p := models.Profile{}
	p.ID= strings.TrimSpace(request.GetID())

	err := p.GetByID(ps.Service.DB(ctx, true))
	if err != nil {
		return nil, service.ErrorProfileDoesNotExist
	}

	properties := map[string]interface{}{}
	for key, value := range request.GetProperties() {
		properties[key] = value
	}

	err = p.UpdateProperties(ps.Service.DB(ctx, false), properties)
	if err != nil {
		return nil, err
	}

	return ps.getProfileByID(ctx, p.ID)
}

func EscapeColumnName(name string) string {
	return strings.NewReplacer(
		"'", "",
		"\\", "",
		"\r", "",
		"\n", "",
	).Replace(name)
}
