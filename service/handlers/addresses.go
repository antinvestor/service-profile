package handlers

import (
	"github.com/antinvestor/service-profile/grpc/profile"
	"github.com/antinvestor/service-profile/models"
	"context"
)

// Adds a new address based on the request.
func (server *ProfileServer) AddAddress(ctx context.Context, request *profile.ProfileAddAddressRequest) (*profile.ProfileObject, error) {
	p := models.Profile{}
	p.ProfileID = request.GetID()
	if err := server.Env.GetRDb(ctx).Find(&p).Error; err != nil {
		return nil, err
	}

	obj := request.GetAddress()

	address := models.Address{}

	if err := address.CreateFull(server.Env.GeWtDb(ctx), obj.GetCountry(), obj.GetArea(), obj.GetStreet(),
		obj.GetHouse(), obj.GetPostcode(), obj.GetLatitude(), obj.GetLongitude(), ); err != nil {
		return nil, err
	}

	profileAddress := models.ProfileAddress{}
	profileAddress.Create(server.Env.GeWtDb(ctx), p.ProfileID, address.AddressID, obj.GetName())

	return p.ToObject(server.Env.GetRDb(ctx))
}


