package business

import (
	"context"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

type AddressBusiness interface {
	GetByProfile(ctx context.Context, profileID string) ([]*models.ProfileAddress, error)
	CreateAddress(ctx context.Context, request *profilev1.AddressObject) (*profilev1.AddressObject, error)
	LinkAddressToProfile(ctx context.Context, profile string, name string, address *profilev1.AddressObject) error

	ToAPI(address *models.Address) *profilev1.AddressObject
}

func NewAddressBusiness(_ context.Context, addressRepo repository.AddressRepository) AddressBusiness {
	return &addressBusiness{
		addressRepo: addressRepo,
	}
}

type addressBusiness struct {
	addressRepo repository.AddressRepository
}

func (aB *addressBusiness) ToAPI(address *models.Address) *profilev1.AddressObject {
	countryName := ""
	if address.Country != nil {
		countryName = address.Country.Name
	}

	addressObj := &profilev1.AddressObject{
		Id:      address.GetID(),
		Name:    address.Name,
		Area:    address.AdminUnit,
		Country: countryName,
	}

	return addressObj
}

func (aB *addressBusiness) GetByProfile(ctx context.Context, profileID string) ([]*models.ProfileAddress, error) {
	return aB.addressRepo.GetByProfileID(ctx, profileID)
}

func (aB *addressBusiness) CreateAddress(
	ctx context.Context,
	request *profilev1.AddressObject,
) (*profilev1.AddressObject, error) {
	logger := util.Log(ctx).WithField("request", request)

	country, err := aB.addressRepo.CountryGetByAny(ctx, request.GetCountry())
	if err != nil {
		logger.WithError(err).Warn("get country error")
		return nil, err
	}

	address, err := aB.addressRepo.GetByNameAdminUnitAndCountry(ctx, request.GetName(), request.GetArea(), country.ISO3)
	if err != nil {
		logger.WithError(err).Warn("get address error")

		if !data.ErrorIsNoRows(err) {
			return nil, err
		}

		a := models.Address{
			Name:      request.GetName(),
			AdminUnit: request.GetArea(),
			CountryID: country.ISO3,
			Country:   country,
		}

		saveErr := aB.addressRepo.Create(ctx, &a)
		if saveErr != nil {
			return nil, saveErr
		}
		address = &a
	}

	return aB.ToAPI(address), nil
}

func (aB *addressBusiness) LinkAddressToProfile(
	ctx context.Context,
	profileID string,
	name string,
	address *profilev1.AddressObject,
) error {
	profileAddresses, err := aB.addressRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	for _, pAddress := range profileAddresses {
		if address.GetId() == pAddress.AddressID {
			return nil
		}
	}

	profileAddress := models.ProfileAddress{
		Name:      name,
		AddressID: address.GetId(),
		ProfileID: profileID,
	}
	return aB.addressRepo.SaveLink(ctx, &profileAddress)
}
