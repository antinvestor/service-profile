package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore/pool"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type latestPositionRepository struct {
	dbPool pool.Pool
}

// NewLatestPositionRepository creates a new repository for latest positions.
// LatestPosition has a single-column PK (subject_id) with no BaseModel, so uses raw pool.
func NewLatestPositionRepository(dbPool pool.Pool) LatestPositionRepository {
	return &latestPositionRepository{dbPool: dbPool}
}

// Upsert creates or updates the latest known position for a subject.
// The geom column is computed by a database trigger from latitude/longitude.
// The freshness guard (WHERE EXCLUDED.ts >= latest_positions.ts) ensures we never
// overwrite a newer position with an older one from out-of-order events.
func (r *latestPositionRepository) Upsert(ctx context.Context, pos *models.LatestPosition) error {
	db := r.dbPool.DB(ctx, false)

	result := db.Exec(
		`INSERT INTO latest_positions (subject_id, latitude, longitude, accuracy, ts, updated_at)
		 VALUES (?, ?, ?, ?, ?, NOW())
		 ON CONFLICT (subject_id)
		 DO UPDATE SET
		     latitude = EXCLUDED.latitude,
		     longitude = EXCLUDED.longitude,
		     accuracy = EXCLUDED.accuracy,
		     ts = EXCLUDED.ts,
		     updated_at = NOW()
		 WHERE EXCLUDED.ts >= latest_positions.ts`,
		pos.SubjectID, pos.Latitude, pos.Longitude, pos.Accuracy, pos.TS,
	)

	if result.Error != nil {
		return fmt.Errorf("upsert latest position for subject %s: %w", pos.SubjectID, result.Error)
	}
	return nil
}

// Get retrieves the latest known position for a subject.
func (r *latestPositionRepository) Get(ctx context.Context, subjectID string) (*models.LatestPosition, error) {
	db := r.dbPool.DB(ctx, true)

	var pos models.LatestPosition
	result := db.Where("subject_id = ?", subjectID).First(&pos)
	if result.Error != nil {
		return nil, fmt.Errorf("get latest position for subject %s: %w", subjectID, result.Error)
	}
	return &pos, nil
}

// GetNearbySubjects finds subjects within radiusMeters of the given point using PostGIS ST_DWithin.
// The geom column on latest_positions is maintained by a database trigger.
// The excludeSubjectID parameter allows excluding the querying subject from results.
// staleHours controls how old a position can be before it's considered stale and excluded.
func (r *latestPositionRepository) GetNearbySubjects(
	ctx context.Context,
	lat, lon, radiusMeters float64,
	excludeSubjectID string,
	staleHours int,
	limit int,
) ([]*SubjectWithDistance, error) {
	db := r.dbPool.DB(ctx, true)

	if staleHours <= 0 {
		staleHours = 1
	}
	staleThreshold := time.Now().Add(-time.Duration(staleHours) * time.Hour)

	const pgSQL = `
		SELECT
			subject_id,
			latitude,
			longitude,
			ts AS last_seen,
			ST_Distance(
				geom::geography,
				ST_SetSRID(ST_Point($1, $2), 4326)::geography
			) AS distance_meters
		FROM latest_positions
		WHERE subject_id != $3
		  AND ts > $4
		  AND geom IS NOT NULL
		  AND ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_Point($1, $2), 4326)::geography,
			$5
		  )
		ORDER BY distance_meters ASC
		LIMIT $6`

	var results []*SubjectWithDistance
	err := db.Raw(pgSQL, lon, lat, excludeSubjectID, staleThreshold, radiusMeters, limit).
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get nearby subjects at (%f, %f) radius %f: %w", lat, lon, radiusMeters, err)
	}

	return results, nil
}
