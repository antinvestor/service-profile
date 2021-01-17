package handlers

import (
	"context"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/models"
)

// Adds a new address based on the request.
func (ps *ProfileServer) AddAddress(ctx context.Context, request *papi.ProfileAddAddressRequest) (*papi.ProfileObject, error) {
	p := models.Profile{}
	p.ProfileID = request.GetID()
	if err := ps.Service.DB(ctx, true).Find(&p).Error; err != nil {
		return nil, err
	}

	obj := request.GetAddress()

	address := models.Address{}

	if err := address.CreateFull(ps.Service.DB(ctx, false), obj.GetCountry(), obj.GetArea(), obj.GetStreet(),
		obj.GetHouse(), obj.GetPostcode(), obj.GetLatitude(), obj.GetLongitude(), ); err != nil {
		return nil, err
	}

	profileAddress := models.ProfileAddress{}
	profileAddress.Create(ps.Service.DB(ctx, false), p.ProfileID, address.AddressID, obj.GetName())

	return p.ToObject(ps.Service.DB(ctx, true))
}


