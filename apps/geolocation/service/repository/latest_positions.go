package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type latestPositionRepository struct {
	datastore.BaseRepository[*models.LatestPosition]
}

func NewLatestPositionRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) LatestPositionRepository {
	return &latestPositionRepository{
		BaseRepository: datastore.NewBaseRepository[*models.LatestPosition](
			ctx,
			dbPool,
			workMan,
			func() *models.LatestPosition { return &models.LatestPosition{} },
		),
	}
}

func (r *latestPositionRepository) Upsert(ctx context.Context, pos *models.LatestPosition) error {
	pos.GenID(ctx)
	now := time.Now()

	result := r.Pool().DB(ctx, false).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "partition_id"},
				{Name: "subject_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"device_id":   pos.DeviceID,
				"latitude":    pos.Latitude,
				"longitude":   pos.Longitude,
				"accuracy":    pos.Accuracy,
				"ts":          pos.TS,
				"modified_at": now,
				"version":     clause.Expr{SQL: "latest_positions.version + 1"},
			}),
			Where: clause.Where{Exprs: []clause.Expression{
				clause.Expr{SQL: "EXCLUDED.ts >= latest_positions.ts"},
			}},
		}).
		Create(pos)
	if result.Error != nil {
		return fmt.Errorf("upsert latest position for subject %s: %w", pos.SubjectID, result.Error)
	}

	return nil
}

func (r *latestPositionRepository) Get(ctx context.Context, subjectID string) (*models.LatestPosition, error) {
	var pos models.LatestPosition
	if err := r.Pool().DB(ctx, true).Where("subject_id = ?", subjectID).Take(&pos).Error; err != nil {
		return nil, fmt.Errorf("get latest position for subject %s: %w", subjectID, err)
	}
	return &pos, nil
}

func (r *latestPositionRepository) GetNearbySubjects(
	ctx context.Context,
	lat, lon, radiusMeters float64,
	excludeSubjectID string,
	staleHours int,
	limit int,
) ([]*SubjectWithDistance, error) {
	var results []*SubjectWithDistance

	if staleHours <= 0 {
		staleHours = 1
	}
	staleThreshold := time.Now().Add(-time.Duration(staleHours) * time.Hour)

	db := r.Pool().DB(ctx, true).
		Model(&models.LatestPosition{}).
		Select(
			"subject_id, latitude, longitude, ts AS last_seen, "+
				"ST_Distance(geom::geography, ST_SetSRID(ST_Point(?, ?), 4326)::geography) AS distance_meters",
			lon,
			lat,
		).
		Where("subject_id <> ?", excludeSubjectID).
		Where("ts > ?", staleThreshold).
		Where("geom IS NOT NULL").
		Where(
			"ST_DWithin(geom::geography, ST_SetSRID(ST_Point(?, ?), 4326)::geography, ?)",
			lon,
			lat,
			radiusMeters,
		).
		Order("distance_meters ASC")
	if limit > 0 {
		db = db.Limit(limit)
	}

	if err := db.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("get nearby subjects at (%f, %f) radius %f: %w", lat, lon, radiusMeters, err)
	}

	return results, nil
}
