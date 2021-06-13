package handlers

import (
	"context"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/service/business"
)

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *papi.ProfileContactRequest, ) (*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.GetByContact(ctx, ps.EncryptionKey, request.GetContact())
}

func (ps *ProfileServer) AddContact(ctx context.Context, request *papi.ProfileAddContactRequest,
) (*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.AddContact(ctx, ps.EncryptionKey, request)
}
