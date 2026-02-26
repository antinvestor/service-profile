package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/pitabwire/frame/datastore/pool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type routeDeviationStateRepository struct {
	dbPool pool.Pool
}

// NewRouteDeviationStateRepository creates a new repository for route deviation state.
func NewRouteDeviationStateRepository(dbPool pool.Pool) RouteDeviationStateRepository {
	return &routeDeviationStateRepository{dbPool: dbPool}
}

// Pool returns the underlying database pool for transaction management.
func (r *routeDeviationStateRepository) Pool() pool.Pool {
	return r.dbPool
}

// UpsertTx creates or updates the route deviation state within an existing transaction.
func (r *routeDeviationStateRepository) UpsertTx(
	tx *gorm.DB,
	state *models.RouteDeviationState,
) error {
	result := tx.Exec(
		`INSERT INTO route_deviation_states
		     (subject_id, route_id, deviated, consecutive_off_route,
		      last_deviation_event_at, last_point_ts, last_lat, last_lon, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
		 ON CONFLICT (subject_id, route_id)
		 DO UPDATE SET
		     deviated = EXCLUDED.deviated,
		     consecutive_off_route = EXCLUDED.consecutive_off_route,
		     last_deviation_event_at = EXCLUDED.last_deviation_event_at,
		     last_point_ts = EXCLUDED.last_point_ts,
		     last_lat = EXCLUDED.last_lat,
		     last_lon = EXCLUDED.last_lon,
		     updated_at = NOW()`,
		state.SubjectID,
		state.RouteID,
		state.Deviated,
		state.ConsecutiveOffRoute,
		state.LastDeviationEventAt,
		state.LastPointTS,
		state.LastLat,
		state.LastLon,
	)
	if result.Error != nil {
		return fmt.Errorf(
			"upsert route deviation state (%s, %s): %w",
			state.SubjectID, state.RouteID, result.Error,
		)
	}
	return nil
}

// GetForUpdate retrieves the route deviation state with a row-level lock.
// Returns (nil, nil) if no state exists for the given (subject, route) pair.
func (r *routeDeviationStateRepository) GetForUpdate(
	tx *gorm.DB,
	subjectID, routeID string,
) (*models.RouteDeviationState, error) {
	var state models.RouteDeviationState
	result := tx.Raw(
		`SELECT * FROM route_deviation_states
		 WHERE subject_id = ? AND route_id = ? FOR UPDATE`,
		subjectID, routeID,
	).Scan(&state)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // intentional: nil state means "no prior state exists"
		}
		return nil, fmt.Errorf(
			"get route deviation state for update (%s, %s): %w",
			subjectID, routeID, result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return nil, nil //nolint:nilnil // intentional: nil state means "no prior state exists"
	}

	return &state, nil
}

// DeleteByRoute removes all route deviation state entries for a given route.
func (r *routeDeviationStateRepository) DeleteByRoute(
	ctx context.Context,
	routeID string,
) error {
	db := r.dbPool.DB(ctx, false)
	result := db.Where("route_id = ?", routeID).Delete(&models.RouteDeviationState{})
	if result.Error != nil {
		return fmt.Errorf(
			"delete route deviation states for route %s: %w",
			routeID, result.Error,
		)
	}
	return nil
}
