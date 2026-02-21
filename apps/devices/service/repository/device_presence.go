package repository

import (
	"context"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type devicePresenceRepository struct {
	datastore.BaseRepository[*models.DevicePresence]
}

func NewDevicePresenceRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) DevicePresenceRepository {
	return &devicePresenceRepository{
		BaseRepository: datastore.NewBaseRepository[*models.DevicePresence](
			ctx, dbPool, workMan, func() *models.DevicePresence { return &models.DevicePresence{} },
		),
	}
}

func (dlr *devicePresenceRepository) GetLatestByDeviceID(
	ctx context.Context, deviceID string,
) (*models.DevicePresence, error) {
	var presence models.DevicePresence
	if err := dlr.Pool().
		DB(ctx, true).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&presence).
		Error; err != nil {
		return nil, err
	}
	return &presence, nil
}

func (dlr *devicePresenceRepository) GetByDeviceID(
	ctx context.Context,
	deviceID string,
) (workerpool.JobResultPipe[[]*models.DevicePresence], error) {
	query := data.NewSearchQuery(data.WithSearchFiltersAndByValue(map[string]any{
		"device_id": deviceID,
	}))

	return dlr.Search(ctx, query)
}

// Upsert creates or updates the presence record for a device.
// Uses device_id as the conflict key to ensure one row per device.
func (dlr *devicePresenceRepository) Upsert(ctx context.Context, presence *models.DevicePresence) error {
	db := dlr.Pool().DB(ctx, false)
	return db.
		Where("device_id = ?", presence.DeviceID).
		Assign(map[string]any{
			"profile_id":     presence.ProfileID,
			"status":         presence.Status,
			"status_message": presence.StatusMessage,
			"expiry_time":    presence.ExpiryTime,
			"data":           presence.Data,
		}).
		FirstOrCreate(presence).Error
}
