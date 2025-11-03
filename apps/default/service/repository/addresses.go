package repository

import (
	"context"
	"strings"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type addressRepository struct {
	datastore.BaseRepository[*models.Address]
}

func NewAddressRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) AddressRepository {
	repo := addressRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Address](
			ctx, dbPool, workMan, func() *models.Address { return &models.Address{} },
		),
	}
	return &repo
}

func (ar *addressRepository) SaveLink(ctx context.Context, profileAddress *models.ProfileAddress) error {
	return ar.Pool().DB(ctx, false).Save(profileAddress).Error
}

func (ar *addressRepository) DeleteLink(ctx context.Context, id string) error {
	pAddress := &models.ProfileAddress{}
	err := ar.Pool().DB(ctx, true).First(pAddress, "id = ?", id).Error
	if err != nil {
		return err
	}
	return ar.Pool().DB(ctx, false).Delete(pAddress).Error
}

func (ar *addressRepository) GetByNameAdminUnitAndCountry(
	ctx context.Context,
	name string,
	adminUnit string,
	countryID string,
) (*models.Address, error) {
	address := &models.Address{}
	err := ar.Pool().DB(ctx, true).
		First(address, "name ilike ? AND admin_unit ilike ? AND country_id ilike ?", name, adminUnit, countryID).
		Error
	return address, err
}

func (ar *addressRepository) GetByProfileID(ctx context.Context, id string) ([]*models.ProfileAddress, error) {
	var addressList []*models.ProfileAddress
	err := ar.Pool().DB(ctx, true).Preload(clause.Associations).Where("profile_id = ?", id).Find(&addressList).Error
	return addressList, err
}

func (ar *addressRepository) CountryGetByISO3(ctx context.Context, countryISO3 string) (*models.Country, error) {
	country := &models.Country{}
	err := ar.Pool().DB(ctx, true).Where("ISO3 = ?", countryISO3).First(country).Error
	return country, err
}

func (ar *addressRepository) CountryGetByAny(ctx context.Context, c string) (*models.Country, error) {
	if c == "" {
		return nil, service.ErrCountryDoesNotExist
	}

	country := &models.Country{}
	upperC := strings.ToUpper(c)

	err := ar.Pool().DB(ctx, true).
		Where("ISO3 ilike ? OR ISO2 ilike ? OR Name ilike ?", upperC, upperC, upperC).
		First(country).
		Error
	return country, err
}

func (ar *addressRepository) CountryGetByName(ctx context.Context, name string) (*models.Country, error) {
	country := &models.Country{}
	err := ar.Pool().DB(ctx, true).Where("name = ilike", strings.ToUpper(name)).First(country).Error
	return country, err
}
