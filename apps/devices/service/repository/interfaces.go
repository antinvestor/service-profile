package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

// DeviceRepository defines the operations for managing devices in storage.
type DeviceRepository interface {
	datastore.BaseRepository[*models.Device]
	RemoveByID(ctx context.Context, id string) (*models.Device, error)
}

// DeviceSessionRepository defines the operations for managing device sessions.
type DeviceSessionRepository interface {
	datastore.BaseRepository[*models.DeviceSession]
	GetLastByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error)
}

// DeviceLogRepository defines the operations for managing device logs.
type DeviceLogRepository interface {
	datastore.BaseRepository[*models.DeviceLog]
	GetByDeviceID(ctx context.Context, deviceID string) (workerpool.JobResultPipe[[]*models.DeviceLog], error)
}

// DevicePresenceRepository defines the operations for managing device presence.
type DevicePresenceRepository interface {
	datastore.BaseRepository[*models.DevicePresence]
	GetByDeviceID(ctx context.Context, deviceID string) (workerpool.JobResultPipe[[]*models.DevicePresence], error)
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.DevicePresence, error)
}

// DeviceKeyRepository defines the operations for managing matrix keys.
type DeviceKeyRepository interface {
	datastore.BaseRepository[*models.DeviceKey]
	GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceKey, error)
	RemoveByID(ctx context.Context, id string) (*models.DeviceKey, error)
}
