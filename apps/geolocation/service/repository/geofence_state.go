package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/pitabwire/frame/datastore/pool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type geofenceStateRepository struct {
	dbPool pool.Pool
}

// NewGeofenceStateRepository creates a new repository for geofence state.
// GeofenceState has a composite primary key (subject_id, area_id), so it doesn't embed BaseModel
// and uses a raw pool instead of BaseRepository.
func NewGeofenceStateRepository(dbPool pool.Pool) GeofenceStateRepository {
	return &geofenceStateRepository{dbPool: dbPool}
}

// Pool returns the underlying database pool for transaction management.
func (r *geofenceStateRepository) Pool() pool.Pool {
	return r.dbPool
}

// UpsertTx creates or updates the geofence state within an existing transaction.
func (r *geofenceStateRepository) UpsertTx(tx *gorm.DB, state *models.GeofenceState) error {
	result := tx.Exec(
		`INSERT INTO geofence_states (subject_id, area_id, inside, last_transition, last_point_ts, enter_ts, last_lat, last_lon, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
		 ON CONFLICT (subject_id, area_id)
		 DO UPDATE SET
		     inside = EXCLUDED.inside,
		     last_transition = EXCLUDED.last_transition,
		     last_point_ts = EXCLUDED.last_point_ts,
		     enter_ts = EXCLUDED.enter_ts,
		     last_lat = EXCLUDED.last_lat,
		     last_lon = EXCLUDED.last_lon,
		     updated_at = NOW()`,
		state.SubjectID,
		state.AreaID,
		state.Inside,
		state.LastTransition,
		state.LastPointTS,
		state.EnterTS,
		state.LastLat,
		state.LastLon,
	)

	if result.Error != nil {
		return fmt.Errorf(
			"upsert geofence state (%s, %s): %w",
			state.SubjectID,
			state.AreaID,
			result.Error,
		)
	}
	return nil
}

// GetForUpdate retrieves the geofence state with a row-level lock (SELECT FOR UPDATE).
// Returns (nil, nil) if no state exists for the given (subject, area) pair.
// This must be called within a transaction to hold the lock.
func (r *geofenceStateRepository) GetForUpdate(
	tx *gorm.DB,
	subjectID, areaID string,
) (*models.GeofenceState, error) {
	var state models.GeofenceState
	result := tx.Raw(
		"SELECT * FROM geofence_states WHERE subject_id = ? AND area_id = ? FOR UPDATE",
		subjectID, areaID,
	).Scan(&state)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // intentional: nil state means "no prior state exists"
		}
		return nil, fmt.Errorf(
			"get geofence state for update (%s, %s): %w",
			subjectID,
			areaID,
			result.Error,
		)
	}

	// GORM Raw+Scan returns RowsAffected=0 when no rows match, without setting ErrRecordNotFound.
	if result.RowsAffected == 0 {
		return nil, nil //nolint:nilnil // intentional: nil state means "no prior state exists"
	}

	return &state, nil
}

// GetInsideByArea returns all geofence states where subjects are currently inside the given area.
func (r *geofenceStateRepository) GetInsideByArea(
	ctx context.Context,
	areaID string,
	limit int,
) ([]*models.GeofenceState, error) {
	db := r.dbPool.DB(ctx, true)

	var states []*models.GeofenceState
	query := db.Where("area_id = ? AND inside = true", areaID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	result := query.Find(&states)
	if result.Error != nil {
		return nil, fmt.Errorf("get inside subjects for area %s: %w", areaID, result.Error)
	}
	return states, nil
}

// DeleteByArea removes all geofence state entries for a given area.
// Used when an area is deleted to clean up stale state.
func (r *geofenceStateRepository) DeleteByArea(ctx context.Context, areaID string) error {
	db := r.dbPool.DB(ctx, false)

	result := db.Where("area_id = ?", areaID).Delete(&models.GeofenceState{})
	if result.Error != nil {
		return fmt.Errorf("delete geofence states for area %s: %w", areaID, result.Error)
	}
	return nil
}
