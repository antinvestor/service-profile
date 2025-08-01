package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

// DeviceRepository defines the operations for managing devices in storage.
type DeviceRepository interface {
	Save(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id string) (*models.Device, error)
	GetByProfileID(ctx context.Context, profileID string) ([]*models.Device, error)
	RemoveByID(ctx context.Context, id string) (*models.Device, error)
}

// DeviceSessionRepository defines the operations for managing device sessions.
type DeviceSessionRepository interface {
	Save(ctx context.Context, session *models.DeviceSession) error
	GetByID(ctx context.Context, id string) (*models.DeviceSession, error)
	GetLastByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error)
}

// DeviceLogRepository defines the operations for managing device logs.
type DeviceLogRepository interface {
	Save(ctx context.Context, deviceLog *models.DeviceLog) error
	GetByID(ctx context.Context, id string) (*models.DeviceLog, error)
	GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceLog, error)
}

// DeviceKeyRepository defines the operations for managing matrix keys.
type DeviceKeyRepository interface {
	Save(ctx context.Context, key *models.DeviceKey) error
	GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceKey, error)
	RemoveByID(ctx context.Context, id string) (*models.DeviceKey, error)
}
