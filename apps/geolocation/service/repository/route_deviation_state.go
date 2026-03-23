package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type routeDeviationStateRepository struct {
	datastore.BaseRepository[*models.RouteDeviationState]
}

func NewRouteDeviationStateRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) RouteDeviationStateRepository {
	return &routeDeviationStateRepository{
		BaseRepository: datastore.NewBaseRepository[*models.RouteDeviationState](
			ctx,
			dbPool,
			workMan,
			func() *models.RouteDeviationState { return &models.RouteDeviationState{} },
		),
	}
}

func (r *routeDeviationStateRepository) UpsertTx( //nolint:dupl // similar upsert pattern, different models
	tx *gorm.DB,
	state *models.RouteDeviationState,
) error {
	state.GenID(tx.Statement.Context)
	now := time.Now()

	result := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "tenant_id"},
			{Name: "partition_id"},
			{Name: "subject_id"},
			{Name: "route_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"deviated":                state.Deviated,
			"consecutive_off_route":   state.ConsecutiveOffRoute,
			"last_deviation_event_at": state.LastDeviationEventAt,
			"last_point_ts":           state.LastPointTS,
			"last_lat":                state.LastLat,
			"last_lon":                state.LastLon,
			"modified_at":             now,
			"version":                 clause.Expr{SQL: "route_deviation_states.version + 1"},
		}),
	}).Create(state)
	if result.Error != nil {
		return fmt.Errorf("upsert route deviation state (%s, %s): %w", state.SubjectID, state.RouteID, result.Error)
	}
	return nil
}

func (r *routeDeviationStateRepository) GetForUpdate(
	tx *gorm.DB,
	subjectID, routeID string,
) (*models.RouteDeviationState, error) {
	var state models.RouteDeviationState
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("subject_id = ? AND route_id = ?", subjectID, routeID).
		Take(&state).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil //nolint:nilnil // not found is a valid empty result
	}
	if err != nil {
		return nil, fmt.Errorf("get route deviation state for update (%s, %s): %w", subjectID, routeID, err)
	}
	return &state, nil
}

func (r *routeDeviationStateRepository) DeleteByRoute(
	ctx context.Context,
	routeID string,
) error {
	if err := r.Pool().DB(ctx, false).
		Where("route_id = ?", routeID).
		Delete(&models.RouteDeviationState{}).Error; err != nil {
		return fmt.Errorf("delete route deviation states for route %s: %w", routeID, err)
	}
	return nil
}
