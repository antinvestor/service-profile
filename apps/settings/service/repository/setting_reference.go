package repository

import (
	"context"
	"fmt"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
)

type ReferenceRepository interface {
	GetByID(ctx context.Context, id string) (*models.SettingRef, error)
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
	Search(
		ctx context.Context,
		module string,
		query string,
		object string,
		objectID string,
		language string,
	) ([]*models.SettingRef, error)
	Save(ctx context.Context, settingRef *models.SettingRef) error
}

type referenceRepository struct {
	service *frame.Service
}

func NewReferenceRepository(_ context.Context, service *frame.Service) ReferenceRepository {
	return &referenceRepository{service: service}
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

	refQuery := repo.service.DB(ctx, true).Where("module = ? AND name = ?", module, name)
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

func (repo *referenceRepository) Search(
	ctx context.Context,
	module string,
	query string,
	object string,
	objectID string,
	language string,
) ([]*models.SettingRef, error) {
	var settingRefs []*models.SettingRef

	queryStr := fmt.Sprintf("%%%s%%", query)

	refQuery := repo.service.DB(ctx, true).Where(" name iLike ?", queryStr)

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

func (repo *referenceRepository) GetByID(ctx context.Context, id string) (*models.SettingRef, error) {
	ref := models.SettingRef{}
	err := repo.service.DB(ctx, true).First(&ref, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (repo *referenceRepository) Save(ctx context.Context, settingRef *models.SettingRef) error {
	return repo.service.DB(ctx, false).Save(settingRef).Error
}
