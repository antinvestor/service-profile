package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
)

type deviceRepository struct {
	datastore.BaseRepository[*models.Device]
}

func NewDeviceRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) DeviceRepository {
	return &deviceRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Device](
			ctx, dbPool, workMan, func() *models.Device { return &models.Device{} },
		),
	}
}

func (dr *deviceRepository) RemoveByID(ctx context.Context, id string) (*models.Device, error) {
	device, err := dr.BaseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := dr.Pool().DB(ctx, false).Delete(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}
