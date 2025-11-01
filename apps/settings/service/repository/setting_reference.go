package repository

import (
	"context"
	"fmt"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
)

type ReferenceRepository interface {
	datastore.BaseRepository[*models.SettingRef]
	GetByName(ctx context.Context, module string, name string) (*models.SettingRef, error)
	GetByNameAndLanguage(ctx context.Context, module string, name string, language string) (*models.SettingRef, error)
	GetByNameAndObject(
		ctx context.Context,
		module string,
		name string,
		object string,
		objectID string,
	) (*models.SettingRef, error)
	GetByNameAndObjectAndLanguage(
		ctx context.Context,
		module string,
		name string,
		object string,
		objectID string,
		language string,
	) (*models.SettingRef, error)
	SearchRef(
		ctx context.Context,
		module string,
		query string,
		object string,
		objectID string,
		language string,
	) ([]*models.SettingRef, error)
}

type referenceRepository struct {
	datastore.BaseRepository[*models.SettingRef]
}

func NewReferenceRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) ReferenceRepository {
	return &referenceRepository{
		BaseRepository: datastore.NewBaseRepository[*models.SettingRef](
			ctx, dbPool, workMan, func() *models.SettingRef { return &models.SettingRef{} },
		),
	}
}

func (repo *referenceRepository) GetByName(
	ctx context.Context,
	module string,
	name string,
) (*models.SettingRef, error) {
	return repo.GetByNameAndObjectAndLanguage(ctx, module, name, "", "", "")
}

func (repo *referenceRepository) GetByNameAndLanguage(ctx context.Context, module string,
	name string, language string) (*models.SettingRef, error) {
	return repo.GetByNameAndObjectAndLanguage(ctx, module, name, "", "", language)
}

func (repo *referenceRepository) GetByNameAndObject(
	ctx context.Context,
	module string,
	name string,
	object string,
	objectID string,
) (*models.SettingRef, error) {
	return repo.GetByNameAndObjectAndLanguage(ctx, module, name, object, objectID, "")
}

func (repo *referenceRepository) GetByNameAndObjectAndLanguage(
	ctx context.Context,
	module string,
	name string,
	object string,
	objectID string,
	language string,
) (*models.SettingRef, error) {
	var settingRef models.SettingRef

	refQuery := repo.Pool().DB(ctx, true).Where("module = ? AND name = ?", module, name)
	if objectID != "" && object != "" {
		refQuery = refQuery.Where(" object = ? AND object_id = ? ", object, objectID)
	}

	if language != "" {
		refQuery = refQuery.Where(" language = ?", language)
	}

	err := refQuery.First(&settingRef).Error
	if err != nil {
		return nil, err
	}
	return &settingRef, nil
}

func (repo *referenceRepository) SearchRef(
	ctx context.Context,
	module string,
	query string,
	object string,
	objectID string,
	language string,
) ([]*models.SettingRef, error) {
	var settingRefs []*models.SettingRef

	queryStr := fmt.Sprintf("%%%s%%", query)

	refQuery := repo.Pool().DB(ctx, true).Where(" name iLike ?", queryStr)

	if module != "" {
		refQuery = refQuery.Where(" module = ? ", module)
	}

	if objectID != "" && object != "" {
		refQuery = refQuery.Where(" object = ? AND object_id = ? ", object, objectID)
	}

	if language != "" {
		refQuery = refQuery.Where(" language = ?", language)
	}

	err := refQuery.Find(&settingRefs).Error
	if err != nil {
		return nil, err
	}
	return settingRefs, nil
}
