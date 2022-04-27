package handlers

import (
	"context"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/service/business"
)

// AddAddress Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context, request *papi.ProfileAddAddressRequest) (*papi.ProfileObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	profileBusiness := business.NewProfileBusiness(ctx, ps.Service)
	return profileBusiness.AddAddress(ctx, ps.EncryptionKey, request)

}
