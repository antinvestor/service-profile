package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame"
)

type deviceKeyRepository struct {
	service *frame.Service
}

func NewDeviceKeyRepository(service *frame.Service) DeviceKeyRepository {
	return &deviceKeyRepository{service: service}
}

func (r *deviceKeyRepository) Save(ctx context.Context, key *models.DeviceKey) error {
	return r.service.DB(ctx, false).Save(key).Error
}

func (r *deviceKeyRepository) GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceKey, error) {
	var keys []*models.DeviceKey
	err := r.service.DB(ctx, true).Where("device_id = ?", deviceID).Find(&keys).Error
	return keys, err
}

func (r *deviceKeyRepository) RemoveByID(ctx context.Context, id string) (*models.DeviceKey, error) {
	var key models.DeviceKey
	if err := r.service.DB(ctx, true).First(&key, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.service.DB(ctx, false).Delete(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}
