package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

type SettingValRepository interface {
	datastore.BaseRepository[*models.SettingVal]
	GetByRef(ctx context.Context, id ...string) ([]*models.SettingVal, error)
}

type settingValRepository struct {
	datastore.BaseRepository[*models.SettingVal]
}

func NewSettingValRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) SettingValRepository {
	return &settingValRepository{
		BaseRepository: datastore.NewBaseRepository[*models.SettingVal](
			ctx, dbPool, workMan, func() *models.SettingVal { return &models.SettingVal{} },
		),
	}
}

func (repo *settingValRepository) GetByRef(ctx context.Context, reference ...string) ([]*models.SettingVal, error) {
	var settingVal []*models.SettingVal
	err := repo.Pool().DB(ctx, true).Find(&settingVal, "ref IN ?", reference).Error
	if err != nil {
		return nil, err
	}
	return settingVal, nil
}
