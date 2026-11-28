package repository

import (
	"context"
	"errors"

	"github.com/pitabwire/frame/v3/datastore"
	"github.com/pitabwire/frame/v3/datastore/pool"
	"github.com/pitabwire/frame/v3/workerpool"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

type SettingValRepository interface {
	datastore.BaseRepository[*models.SettingVal]
	GetByRef(ctx context.Context, id ...string) ([]*models.SettingVal, error)
	// UpdateDetail writes val.Detail under optimistic locking and bumps the
	// version. It returns ErrSettingValVersionConflict when the row changed
	// since val was read.
	UpdateDetail(ctx context.Context, val *models.SettingVal) error
}

// ErrSettingValVersionConflict reports that a setting value was modified
// between being read and being updated.
var ErrSettingValVersionConflict = errors.New("setting value version conflict")

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

func (repo *settingValRepository) UpdateDetail(ctx context.Context, val *models.SettingVal) error {
	// BaseRepository.Update matches on the version that was read; BeforeUpdate
	// bumps it, so the version and audit columns must be part of the write.
	rows, err := repo.Update(ctx, val, "detail", "version", "modified_at", "modified_by")
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSettingValVersionConflict
	}
	return nil
}
