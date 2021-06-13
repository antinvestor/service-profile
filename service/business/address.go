package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
)

type AddressBusiness interface {
	GetByProfile(ctx context.Context, profileID string) ([]*models.ProfileAddress, error)
	CreateAddress(ctx context.Context, request *profileV1.AddressObject) (*profileV1.AddressObject, error)
	LinkAddressToProfile(ctx context.Context, profile string, name string, address *profileV1.AddressObject)  error

	ToApi(address *models.Address) *profileV1.AddressObject
}

func NewAddressBusiness(ctx context.Context, service *frame.Service) AddressBusiness {
	addressRepo := repository.NewAddressRepository(service)
	return &addressBusiness{
		service:    service,
		addressRepo: addressRepo,
	}
}

type addressBusiness struct {
	service     *frame.Service
	addressRepo repository.AddressRepository
}

func (aB *addressBusiness) ToApi(address *models.Address) *profileV1.AddressObject {

	addressObj := &profileV1.AddressObject{
		ID: address.GetID(),
		Name: address.Name,
		Area: address.AdminUnit,
		Country: address.Country.Name,

	}

	return addressObj

}

func (aB *addressBusiness) GetByProfile(ctx context.Context, profileID string) ([]*models.ProfileAddress, error) {
	return aB.addressRepo.GetByProfileID(ctx, profileID)
}


func (aB *addressBusiness) CreateAddress(ctx context.Context, request *profileV1.AddressObject) (*profileV1.AddressObject, error){


	country, err := aB.addressRepo.CountryGetByAny(ctx, request.GetCountry())
	if err != nil {
		return nil, err
	}

	address, err := aB.addressRepo.GetByNameAdminUnitAndCountry(ctx, request.GetName(), request.GetArea(), country.ISO3)
	if err != nil {

		if !frame.DBErrorIsRecordNotFound(err) {
			return nil, err
		}

		a := models.Address{
			Name: request.GetName(),
			AdminUnit: request.GetArea(),
			CountryID: country.ISO3,
			Country: country,
		}

		err := aB.addressRepo.Save(ctx, &a)
		if err != nil {
			return nil, err
		}
		address = &a
	}

	return aB.ToApi(address), nil

}
func (aB *addressBusiness) LinkAddressToProfile(ctx context.Context, profileID string, name string, address *profileV1.AddressObject)  error{


	profileAddresses, err := aB.addressRepo.GetByProfileID(ctx, profileID)
	if err != nil{
		return err
	}

	for _, pAddress := range  profileAddresses{
		if  address.GetID() == pAddress.AddressID {
			return nil
		}
	}

	profileAddress := models.ProfileAddress{
		Name: name,
		AddressID: address.GetID(),
		ProfileID: profileID,
	}
	return aB.addressRepo.SaveLink(ctx, &profileAddress)

}