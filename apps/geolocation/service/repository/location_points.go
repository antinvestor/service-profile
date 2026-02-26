package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

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
