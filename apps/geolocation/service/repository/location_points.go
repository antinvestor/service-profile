package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type locationPointRepository struct {
	datastore.BaseRepository[*models.LocationPoint]
}

// NewLocationPointRepository creates a new repository for location points.
func NewLocationPointRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) LocationPointRepository {
	return &locationPointRepository{
		BaseRepository: datastore.NewBaseRepository[*models.LocationPoint](
			ctx, dbPool, workMan, func() *models.LocationPoint { return &models.LocationPoint{} },
		),
	}
}

// GetTrack retrieves location history for a subject within a time range, ordered by timestamp descending.
func (r *locationPointRepository) GetTrack(
	ctx context.Context,
	subjectID string,
	from, to time.Time,
	limit, offset int,
) ([]*models.LocationPoint, error) {
	var points []*models.LocationPoint

	db := r.Pool().DB(ctx, true)
	query := db.Where("subject_id = ? AND ts >= ? AND ts <= ?", subjectID, from, to).
		Order("ts DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&points)
	if result.Error != nil {
		return nil, fmt.Errorf("get track for subject %s: %w", subjectID, result.Error)
	}

	return points, nil
}

func (r *locationPointRepository) GetPendingForProcessing(
	ctx context.Context,
	limit int,
) ([]*models.LocationPoint, error) {
	var points []*models.LocationPoint

	db := r.Pool().DB(ctx, true)
	query := db.Where(
		"processing_state IN ? AND deleted_at IS NULL",
		[]models.LocationPointProcessingState{
			models.LocationPointProcessingStatePending,
			models.LocationPointProcessingStateFailed,
		},
	).Order("ingested_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&points).Error; err != nil {
		return nil, fmt.Errorf("get pending location points: %w", err)
	}

	return points, nil
}

func (r *locationPointRepository) MarkProcessed(
	ctx context.Context,
	pointID string,
) error {
	return r.updateProcessingState(ctx, pointID, models.LocationPointProcessingStateProcessed, "")
}

func (r *locationPointRepository) MarkFailed(
	ctx context.Context,
	pointID string,
	processingErr error,
) error {
	if processingErr == nil {
		processingErr = gorm.ErrInvalidData
	}
	return r.updateProcessingState(
		ctx,
		pointID,
		models.LocationPointProcessingStateFailed,
		processingErr.Error(),
	)
}

func (r *locationPointRepository) updateProcessingState(
	ctx context.Context,
	pointID string,
	state models.LocationPointProcessingState,
	processingErr string,
) error {
	tableName := (&models.LocationPoint{}).TableName()
	updates := map[string]any{
		"processing_state": state,
		"processing_error": processingErr,
	}
	if state == models.LocationPointProcessingStateProcessed {
		updates["processed_at"] = time.Now()
	} else {
		updates["processed_at"] = nil
	}

	result := r.Pool().DB(ctx, false).
		Table(tableName).
		Where("id = ?", pointID)
	result = result.
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update processing state for point %s: %w", pointID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("location point %s not found or deleted", pointID)
	}

	return nil
}
