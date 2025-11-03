package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

type SettingAuditRepository interface {
	datastore.BaseRepository[*models.SettingAudit]
	GetByRef(ctx context.Context, ref string) ([]models.SettingAudit, error)
}

type settingAuditRepository struct {
	datastore.BaseRepository[*models.SettingAudit]
}

func NewSettingAuditRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) SettingAuditRepository {
	return &settingAuditRepository{
		BaseRepository: datastore.NewBaseRepository[*models.SettingAudit](
			ctx, dbPool, workMan, func() *models.SettingAudit { return &models.SettingAudit{} },
		),
	}
}

func (repo *settingAuditRepository) GetByRef(ctx context.Context, ref string) ([]models.SettingAudit, error) {
	var settingAudit []models.SettingAudit
	err := repo.Pool().DB(ctx, false).Find(&settingAudit, "ref = ?", ref).Error
	if err != nil {
		return nil, err
	}
	return settingAudit, nil
}
