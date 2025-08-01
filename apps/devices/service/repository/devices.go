package repository

import (
	"context"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type deviceRepository struct {
	service *frame.Service
}

func NewDeviceRepository(service *frame.Service) DeviceRepository {
	return &deviceRepository{service: service}
}

func (dr *deviceRepository) Save(ctx context.Context, device *models.Device) error {
	return dr.service.DB(ctx, false).Save(device).Error
}

func (dr *deviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	if err := dr.service.DB(ctx, true).First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (dr *deviceRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Device, error) {
	var devices []*models.Device
	db := dr.service.DB(ctx, true)
	if err := db.Where("profile_id = ?", profileID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (dr *deviceRepository) RemoveByID(ctx context.Context, id string) (*models.Device, error) {
	var device models.Device
	if err := dr.service.DB(ctx, true).First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := dr.service.DB(ctx, false).Delete(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

type deviceLogRepository struct {
	service *frame.Service
}

func NewDeviceLogRepository(service *frame.Service) DeviceLogRepository {
	return &deviceLogRepository{service: service}
}

func (r *deviceLogRepository) Save(ctx context.Context, log *models.DeviceLog) error {
	return r.service.DB(ctx, false).Save(log).Error
}

func (r *deviceLogRepository) GetByID(ctx context.Context, id string) (*models.DeviceLog, error) {
	var log models.DeviceLog
	if err := r.service.DB(ctx, true).First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *deviceLogRepository) GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceLog, error) {
	var logs []*models.DeviceLog
	db := r.service.DB(ctx, true)
	if err := db.Where("device_id = ?", deviceID).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
