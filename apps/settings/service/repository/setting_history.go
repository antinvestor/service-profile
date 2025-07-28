package repository

import (
	"context"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

type SettingAuditRepository interface {
	GetByRef(ctx context.Context, ref string) ([]models.SettingAudit, error)
	Save(ctx context.Context, sAudit *models.SettingAudit) error
}

type settingAuditRepository struct {
	service *frame.Service
}

func NewSettingAuditRepository(_ context.Context, service *frame.Service) SettingAuditRepository {
	return &settingAuditRepository{service: service}
}

func (repo *settingAuditRepository) GetByRef(ctx context.Context, ref string) ([]models.SettingAudit, error) {
	var settingAudit []models.SettingAudit
	svc := repo.service
	err := svc.DB(ctx, false).Find(&settingAudit, "ref = ?", ref).Error
	if err != nil {
		return nil, err
	}
	return settingAudit, nil
}

func (repo *settingAuditRepository) Save(ctx context.Context, sAudit *models.SettingAudit) error {
	return repo.service.DB(ctx, false).Save(sAudit).Error
}
