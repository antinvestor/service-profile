package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type DeviceRepository interface {
	Save(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id string) (*models.Device, error)
	GetByLinkID(ctx context.Context, linkId string) (*models.Device, error)
	List(ctx context.Context, profileID string) ([]*models.Device, error)
}

type DeviceLogRepository interface {
	Save(ctx context.Context, deviceLog *models.DeviceLog) error
	GetByID(ctx context.Context, id string) (*models.DeviceLog, error)
	GetByLinkID(ctx context.Context, linkID string) (*models.DeviceLog, error)
	ListByDeviceID(ctx context.Context, deviceLogID string) ([]*models.DeviceLog, error)
}
