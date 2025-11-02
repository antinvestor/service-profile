package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
)

type deviceSessionRepository struct {
	datastore.BaseRepository[*models.DeviceSession]
}

func NewDeviceSessionRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) DeviceSessionRepository {
	return &deviceSessionRepository{

		BaseRepository: datastore.NewBaseRepository[*models.DeviceSession](
			ctx, dbPool, workMan, func() *models.DeviceSession { return &models.DeviceSession{} },
		),
	}
}

func (r *deviceSessionRepository) Save(ctx context.Context, session *models.DeviceSession) error {
	return r.Pool().DB(ctx, false).Save(session).Error
}

func (r *deviceSessionRepository) GetByID(ctx context.Context, id string) (*models.DeviceSession, error) {
	var session models.DeviceSession
	if err := r.Pool().DB(ctx, true).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *deviceSessionRepository) GetLastByDeviceID(
	ctx context.Context,
	deviceID string,
) (*models.DeviceSession, error) {
	var session models.DeviceSession
	if err := r.Pool().DB(ctx, true).Where("device_id = ?", deviceID).Order("created_at DESC").First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
