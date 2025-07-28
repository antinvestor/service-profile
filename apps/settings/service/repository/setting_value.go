package repository

import (
	"context"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

type SettingValRepository interface {
	GetByID(ctx context.Context, id string) (*models.SettingVal, error)
	GetByRef(ctx context.Context, id string) (*models.SettingVal, error)
	Save(ctx context.Context, settingVal *models.SettingVal) error
}

type settingValRepository struct {
	service *frame.Service
}

func NewSettingValRepository(_ context.Context, service *frame.Service) SettingValRepository {
	return &settingValRepository{service: service}
}

func (repo *settingValRepository) GetByID(ctx context.Context, id string) (*models.SettingVal, error) {
	settingVal := models.SettingVal{}
	err := repo.service.DB(ctx, true).First(&settingVal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &settingVal, nil
}

func (repo *settingValRepository) GetByRef(ctx context.Context, id string) (*models.SettingVal, error) {
	settingVal := models.SettingVal{}
	err := repo.service.DB(ctx, true).First(&settingVal, "ref = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &settingVal, nil
}

func (repo *settingValRepository) Save(ctx context.Context, sVal *models.SettingVal) error {
	return repo.service.DB(ctx, false).Save(sVal).Error
}
