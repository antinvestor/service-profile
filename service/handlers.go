package service

import (
	"bitbucket.org/antinvestor/service-profile/profile"
	"context"
	"fmt"
	"strings"
)

type ProfileServer struct {
	Env *Env
}

func (server *ProfileServer) getProfileByID(ctx context.Context, profileID string, ) (*profile.ProfileObject, error) {
	p := Profile{}
	p.ProfileID = profileID

	server.Env.GetRDb(ctx).First(&p)

	return p.ToObject(server.Env.GetRDb(ctx))
}

func (server *ProfileServer) GetByHash(ctx context.Context,
	request *profile.ProfileIDRequest, ) (*profile.ProfileObject, error) {
	profileID := strings.TrimSpace(request.GetID())
	return server.getProfileByID(ctx, profileID)
}

func (server *ProfileServer) Search(request *profile.ProfileSearchRequest,
	stream profile.ProfileService_SearchServer, ) error {

	profiles := []Profile{}

	// creating WHERE clause to query by properties JSONB
	scope := server.Env.GetRDb(ctx).New()
	for _, property := range request.GetProperties() {
		column := fmt.Sprintf("properties->>'%s'", EscapeColumnName(property), )
		scope = scope.Or(column+" LIKE ?", "%"+request.GetQuery()+"%")
	}

	scope.Find(&profiles)

	for _, p := range profiles {
		profileObject, err := p.ToObject(server.Env.GetRDb(ctx))
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

func (server *ProfileServer) Merge(ctx context.Context, request *profile.ProfileMergeRequest, ) (
	*profile.ProfileObject, error) {

	var target Profile
	var merging Profile

	target.ProfileID = request.GetID()

	if err := target.GetByID(server.Env.GetRDb(ctx)); err != nil {
		return nil, err
	}

	merging.ProfileID = request.GetMergeID()

	if err := merging.GetByID(server.Env.GetRDb(ctx)); err != nil {
		return nil, err
	}

	for key, value := range merging.Properties {

		existingValue := target.Properties[key]

		if existingValue == value {
			continue
		}

		target.Properties[key] = value
	}

	target.UpdateProperties(server.Env.GeWtDb(ctx), merging.Properties)

	return target.ToObject(server.Env.GetRDb(ctx))
}

func (server *ProfileServer) Create(ctx context.Context, request *profile.ProfileCreateRequest, ) (
	*profile.ProfileObject, error) {

	properties := make(map[string]interface{})

	for key, value := range request.GetProperties() {
		properties[key] = value
	}

	p := Profile{}

	err := p.Create(
		server.Env.GeWtDb(ctx),
		request.GetType(),
		strings.TrimSpace(request.GetContact()),
		properties,
	)
	if err != nil {
		return nil, err
	}

	return server.getProfileByID(ctx, p.ProfileID)
}

func (server *ProfileServer) Update(
	ctx context.Context,
	request *profile.ProfileUpdateRequest,
) (*profile.ProfileObject, error) {
	p := Profile{
		ProfileID: strings.TrimSpace(request.GetID()),
	}

	err := p.GetByID(server.Env.GetRDb(ctx))
	if err != nil {
		return nil, profile.ErrorProfileDoesNotExist
	}

	properties := map[string]interface{}{}
	for key, value := range request.GetProperties() {
		properties[key] = value
	}

	p.UpdateProperties(server.Env.GeWtDb(ctx), properties)

	return server.getProfileByID(ctx, p.ProfileID)
}

func EscapeColumnName(name string) string {
	return strings.NewReplacer(
		"'", "",
		"\\", "",
		"\r", "",
		"\n", "",
	).Replace(name)
}

func (server *ProfileServer) AddContact(ctx context.Context, request *profile.ProfileAddContactRequest,
) (*profile.ProfileObject, error) {

	p := Profile{}
	p.ProfileID = request.GetID()
	if err := server.Env.GetRDb(ctx).Find(&p).Error; err != nil {
		return nil, err
	}

	contact := Contact{}
	if err := contact.Create(server.Env.GeWtDb(ctx), p.ProfileID, request.GetContact()); err != nil {
		return nil, err
	}

	return p.ToObject(server.Env.GetRDb(ctx))
}

// Adds a new address based on the request.
func (server *ProfileServer) AddAddress(ctx context.Context, request *profile.ProfileAddAddressRequest) (*profile.ProfileObject, error) {
	p := Profile{}
	p.ProfileID = request.GetID()
	if err := server.Env.GetRDb(ctx).Find(&p).Error; err != nil {
		return nil, err
	}

	obj := request.GetAddress()

	address := Address{}

	if err := address.CreateFull(server.Env.GeWtDb(ctx), obj.GetCountry(), obj.GetTown(), obj.GetLocation(), obj.GetArea(), obj.GetStreet(),
		obj.GetHouse(), obj.GetPostcode(), obj.GetLatitude(), obj.GetLongitude(), ); err != nil {
		return nil, err
	}

	profileAddress := ProfileAddress{}
	profileAddress.Create(server.Env.GeWtDb(ctx), p.ProfileID, address.AddressID, obj.GetName())

	return p.ToObject(server.Env.GetRDb(ctx))
}
