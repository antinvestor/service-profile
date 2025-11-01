package repository

import (
	"context"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/pitabwire/frame"
)

type deviceSessionRepository struct {
	service *frame.Service
}

func NewDeviceSessionRepository(service *frame.Service) DeviceSessionRepository {
	return &deviceSessionRepository{service: service}
}

func (r *deviceSessionRepository) Save(ctx context.Context, session *models.DeviceSession) error {
	return r.service.DB(ctx, false).Save(session).Error
}

func (r *deviceSessionRepository) GetByID(ctx context.Context, id string) (*models.DeviceSession, error) {
	var session models.DeviceSession
	if err := r.service.DB(ctx, true).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *deviceSessionRepository) GetLastByDeviceID(
	ctx context.Context,
	deviceID string,
) (*models.DeviceSession, error) {
	var session models.DeviceSession
	if err := r.service.DB(ctx, true).Where("device_id = ?", deviceID).Order("created_at DESC").First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
