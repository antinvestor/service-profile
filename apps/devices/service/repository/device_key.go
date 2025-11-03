package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type deviceKeyRepository struct {
	datastore.BaseRepository[*models.DeviceKey]
}

func NewDeviceKeyRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) DeviceKeyRepository {
	return &deviceKeyRepository{
		BaseRepository: datastore.NewBaseRepository[*models.DeviceKey](
			ctx, dbPool, workMan, func() *models.DeviceKey { return &models.DeviceKey{} },
		),
	}
}

func (r *deviceKeyRepository) Save(ctx context.Context, key *models.DeviceKey) error {
	return r.Pool().DB(ctx, false).Save(key).Error
}

func (r *deviceKeyRepository) GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceKey, error) {
	var keys []*models.DeviceKey
	err := r.Pool().DB(ctx, true).Where("device_id = ?", deviceID).Find(&keys).Error
	return keys, err
}

func (r *deviceKeyRepository) RemoveByID(ctx context.Context, id string) (*models.DeviceKey, error) {
	var key models.DeviceKey
	if err := r.Pool().DB(ctx, true).First(&key, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.Pool().DB(ctx, false).Delete(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}
