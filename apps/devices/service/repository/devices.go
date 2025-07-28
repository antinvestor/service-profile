package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame"
)

type deviceRepository struct {
	service *frame.Service
}

func (dr *deviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	device := &models.Device{}
	err := dr.service.DB(ctx, true).First(device, "id = ?", id).Error
	return device, err
}

// GetByLinkID retrieves a device by its link ID.
func (dr *deviceRepository) GetByLinkID(ctx context.Context, linkID string) (*models.Device, error) {
	db := dr.service.DB(ctx, true) // true for read-only

	var device models.Device
	result := db.Where("link_id = ?", linkID).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}

	return &device, nil
}

// List retrieves devices based on the provided filter.
func (dr *deviceRepository) List(ctx context.Context, profileID string) ([]*models.Device, error) {
	db := dr.service.DB(ctx, true) // true for read-only

	var devices []*models.Device
	query := db.Where(" profile_id = ?", profileID)

	result := query.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func (dr *deviceRepository) Save(ctx context.Context, device *models.Device) error {
	return dr.service.DB(ctx, false).Save(device).Error
}

func NewDeviceRepository(service *frame.Service) DeviceRepository {
	repository := deviceRepository{
		service: service,
	}
	return &repository
}

type deviceLogRepository struct {
	service *frame.Service
}

func (dlr *deviceLogRepository) GetByID(ctx context.Context, id string) (*models.DeviceLog, error) {
	deviceLog := &models.DeviceLog{}
	err := dlr.service.DB(ctx, true).First(deviceLog, "id = ?", id).Error
	return deviceLog, err
}

func (dlr *deviceLogRepository) GetByLinkID(ctx context.Context, linkID string) (*models.DeviceLog, error) {
	deviceLog := &models.DeviceLog{}
	err := dlr.service.DB(ctx, true).First(deviceLog, "link_id = ?", linkID).Error
	return deviceLog, err
}

func (dlr *deviceLogRepository) ListByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceLog, error) {
	var deviceLogs []*models.DeviceLog

	err := dlr.service.DB(ctx, true).Where(" device_id = ?", deviceID).Find(&deviceLogs).Error
	return deviceLogs, err
}

func (dlr *deviceLogRepository) Save(ctx context.Context, device *models.DeviceLog) error {
	return dlr.service.DB(ctx, false).Save(device).Error
}

func NewDeviceLogRepository(service *frame.Service) DeviceLogRepository {
	repository := deviceLogRepository{
		service: service,
	}
	return &repository
}
