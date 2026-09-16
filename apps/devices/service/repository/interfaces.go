package repository

import (
	"context"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// DeviceRepository defines the operations for managing devices in storage.
type DeviceRepository interface {
	datastore.BaseRepository[*models.Device]
	RemoveByID(ctx context.Context, id string) (*models.Device, error)
	LinkProfileWithFact(ctx context.Context, device *models.Device, fact *outbox.Event) error
}

// DeviceSessionRepository defines the operations for managing device sessions.
type DeviceSessionRepository interface {
	datastore.BaseRepository[*models.DeviceSession]
	GetLastByDeviceID(ctx context.Context, deviceID string) (*models.DeviceSession, error)
	GetLatestByDeviceIDs(ctx context.Context, deviceIDs []string) (map[string]*models.DeviceSession, error)
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
	Upsert(ctx context.Context, presence *models.DevicePresence) error
}

// DeviceKeyRepository defines the operations for managing matrix keys.
type DeviceKeyRepository interface {
	datastore.BaseRepository[*models.DeviceKey]
	GetByDeviceID(ctx context.Context, deviceID string) ([]*models.DeviceKey, error)
	GetByTypeAndKey(ctx context.Context, keyType devicev1.KeyType, key []byte) (*models.DeviceKey, error)
	CreateWithFact(ctx context.Context, key *models.DeviceKey, fact func(*models.DeviceKey) *outbox.Event) error
	RemoveByID(ctx context.Context, id string) (*models.DeviceKey, error)
	RemoveByIDWithFact(
		ctx context.Context,
		id string,
		fact func(*models.DeviceKey) *outbox.Event,
	) (*models.DeviceKey, error)
}
