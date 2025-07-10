package repository

import (
	"context"
	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"strings"

	"gorm.io/gorm/clause"

	"github.com/pitabwire/frame"
)

type addressRepository struct {
	service *frame.Service
}

func (ar *addressRepository) SaveLink(ctx context.Context, profileAddress *models.ProfileAddress) error {
	return ar.service.DB(ctx, false).Save(profileAddress).Error
}

func (ar *addressRepository) DeleteLink(ctx context.Context, id string) error {
	pAddress := &models.ProfileAddress{}
	err := ar.service.DB(ctx, true).First(pAddress, "id = ?", id).Error
	if err != nil {
		return err
	}
	return ar.service.DB(ctx, false).Delete(pAddress).Error
}

func (ar *addressRepository) GetByID(ctx context.Context, id string) (*models.Address, error) {
	address := &models.Address{}
	err := ar.service.DB(ctx, true).First(address, "id = ?", id).Error
	return address, err
}

func (ar *addressRepository) GetByNameAdminUnitAndCountry(
	ctx context.Context,
	name string,
	adminUnit string,
	countryID string,
) (*models.Address, error) {
	address := &models.Address{}
	err := ar.service.DB(ctx, true).
		First(address, "name ilike ? AND admin_unit ilike ? AND country_id ilike ?", name, adminUnit, countryID).
		Error
	return address, err
}

func (ar *addressRepository) GetByProfileID(ctx context.Context, id string) ([]*models.ProfileAddress, error) {
	var addressList []*models.ProfileAddress
	err := ar.service.DB(ctx, true).Preload(clause.Associations).Where("profile_id = ?", id).Find(&addressList).Error
	return addressList, err
}

func (ar *addressRepository) Save(ctx context.Context, address *models.Address) error {
	return ar.service.DB(ctx, false).Save(address).Error
}

func (ar *addressRepository) Delete(ctx context.Context, id string) error {
	address, err := ar.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return ar.service.DB(ctx, false).Delete(address).Error
}

func (ar *addressRepository) CountryGetByISO3(ctx context.Context, countryISO3 string) (*models.Country, error) {
	country := &models.Country{}
	err := ar.service.DB(ctx, true).Where("ISO3 = ?", countryISO3).First(country).Error
	return country, err
}

func (ar *addressRepository) CountryGetByAny(ctx context.Context, c string) (*models.Country, error) {
	if c == "" {
		return nil, service.ErrorCountryDoesNotExist
	}

	country := &models.Country{}
	upperC := strings.ToUpper(c)

	err := ar.service.DB(ctx, true).
		Where("ISO3 ilike ? OR ISO2 ilike ? OR Name ilike ?", upperC, upperC, upperC).
		First(country).
		Error
	return country, err
}

func (ar *addressRepository) CountryGetByName(ctx context.Context, name string) (*models.Country, error) {
	country := &models.Country{}
	err := ar.service.DB(ctx, true).Where("name = ilike", strings.ToUpper(name)).First(country).Error
	return country, err
}

func NewAddressRepository(service *frame.Service) AddressRepository {
	addressRepository := addressRepository{
		service: service,
	}
	return &addressRepository
}
