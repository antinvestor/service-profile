package repository

import (
	"context"
	"fmt"

	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/workerpool"

	"github.com/antinvestor/service-profile/apps/chatagent/service/models"
)

// ContextRepository manages versioned intake context definitions.
type ContextRepository interface {
	datastore.BaseRepository[*models.IntakeContext]
	GetLatest(ctx context.Context, contextKey string) (*models.IntakeContext, error)
	GetVersion(ctx context.Context, contextKey string, version int) (*models.IntakeContext, error)
	ListLatest(ctx context.Context) ([]*models.IntakeContext, error)
	NextVersion(ctx context.Context, contextKey string) (int, error)
}

type contextRepository struct {
	datastore.BaseRepository[*models.IntakeContext]
}

func NewContextRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) ContextRepository {
	return &contextRepository{
		BaseRepository: datastore.NewBaseRepository[*models.IntakeContext](
			ctx, dbPool, workMan, func() *models.IntakeContext { return &models.IntakeContext{} },
		),
	}
}

func (r *contextRepository) GetLatest(ctx context.Context, contextKey string) (*models.IntakeContext, error) {
	var row models.IntakeContext
	err := r.Pool().DB(ctx, true).
		Where("context_key = ? AND active = ?", contextKey, true).
		Order("version DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contextRepository) GetVersion(ctx context.Context, contextKey string, version int) (*models.IntakeContext, error) {
	var row models.IntakeContext
	err := r.Pool().DB(ctx, true).
		Where("context_key = ? AND version = ?", contextKey, version).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contextRepository) ListLatest(ctx context.Context) ([]*models.IntakeContext, error) {
	// Distinct on context_key keeping highest version per key.
	var rows []*models.IntakeContext
	sub := r.Pool().DB(ctx, true).Model(&models.IntakeContext{}).
		Select("context_key, MAX(version) AS version").
		Where("active = ?", true).
		Group("context_key")
	err := r.Pool().DB(ctx, true).
		Joins("JOIN (?) AS latest ON chat_contexts.context_key = latest.context_key AND chat_contexts.version = latest.version", sub).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list latest contexts: %w", err)
	}
	return rows, nil
}

func (r *contextRepository) NextVersion(ctx context.Context, contextKey string) (int, error) {
	var maxVer *int
	err := r.Pool().DB(ctx, true).Model(&models.IntakeContext{}).
		Select("MAX(version)").
		Where("context_key = ?", contextKey).
		Scan(&maxVer).Error
	if err != nil {
		return 0, err
	}
	if maxVer == nil {
		return 1, nil
	}
	return *maxVer + 1, nil
}
