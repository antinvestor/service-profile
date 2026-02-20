package repository

import (
	"context"
	"fmt"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
)

type deviceSessionRepository struct {
	datastore.BaseRepository[*models.DeviceSession]
}

func NewDeviceSessionRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) DeviceSessionRepository {
	return &deviceSessionRepository{

		BaseRepository: datastore.NewBaseRepository[*models.DeviceSession](
			ctx, dbPool, workMan, func() *models.DeviceSession { return &models.DeviceSession{} },
		),
	}
}

func (r *deviceSessionRepository) GetLastByDeviceID(
	ctx context.Context,
	deviceID string,
) (*models.DeviceSession, error) {
	var session models.DeviceSession
	if err := r.Pool().
		DB(ctx, true).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&session).
		Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetLatestByDeviceIDs retrieves the most recent session for each of the given device IDs
// in a single query using PostgreSQL DISTINCT ON, eliminating N+1 query patterns.
func (r *deviceSessionRepository) GetLatestByDeviceIDs(
	ctx context.Context,
	deviceIDs []string,
) (map[string]*models.DeviceSession, error) {
	if len(deviceIDs) == 0 {
		return map[string]*models.DeviceSession{}, nil
	}

	var sessions []*models.DeviceSession
	err := r.Pool().
		DB(ctx, true).
		Raw(`SELECT DISTINCT ON (device_id) *
			 FROM device_sessions
			 WHERE device_id IN (?)
			 ORDER BY device_id, created_at DESC`, deviceIDs).
		Scan(&sessions).
		Error
	if err != nil {
		return nil, fmt.Errorf("get latest sessions by device ids: %w", err)
	}

	result := make(map[string]*models.DeviceSession, len(sessions))
	for _, s := range sessions {
		result[s.DeviceID] = s
	}
	return result, nil
}
