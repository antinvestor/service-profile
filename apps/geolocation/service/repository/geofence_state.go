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

type geofenceStateRepository struct {
	datastore.BaseRepository[*models.GeofenceState]
}

func NewGeofenceStateRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) GeofenceStateRepository {
	return &geofenceStateRepository{
		BaseRepository: datastore.NewBaseRepository[*models.GeofenceState](
			ctx,
			dbPool,
			workMan,
			func() *models.GeofenceState { return &models.GeofenceState{} },
		),
	}
}

func (r *geofenceStateRepository) UpsertTx( //nolint:dupl // similar upsert pattern, different models
	tx *gorm.DB, state *models.GeofenceState,
) error {
	state.GenID(tx.Statement.Context)
	now := time.Now()

	result := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "tenant_id"},
			{Name: "partition_id"},
			{Name: "subject_id"},
			{Name: "area_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"inside":          state.Inside,
			"last_transition": state.LastTransition,
			"last_point_ts":   state.LastPointTS,
			"enter_ts":        state.EnterTS,
			"last_lat":        state.LastLat,
			"last_lon":        state.LastLon,
			"modified_at":     now,
			"version":         clause.Expr{SQL: "geofence_states.version + 1"},
		}),
	}).Create(state)
	if result.Error != nil {
		return fmt.Errorf("upsert geofence state (%s, %s): %w", state.SubjectID, state.AreaID, result.Error)
	}

	return nil
}

func (r *geofenceStateRepository) GetForUpdate(
	tx *gorm.DB,
	subjectID, areaID string,
) (*models.GeofenceState, error) {
	var state models.GeofenceState
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("subject_id = ? AND area_id = ?", subjectID, areaID).
		Take(&state).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil //nolint:nilnil // not found is a valid empty result
	}
	if err != nil {
		return nil, fmt.Errorf("get geofence state for update (%s, %s): %w", subjectID, areaID, err)
	}

	return &state, nil
}

func (r *geofenceStateRepository) GetInsideByArea(
	ctx context.Context,
	areaID string,
	limit int,
) ([]*models.GeofenceState, error) {
	var states []*models.GeofenceState

	query := r.Pool().DB(ctx, true).
		Where("area_id = ? AND inside = ?", areaID, true)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&states).Error; err != nil {
		return nil, fmt.Errorf("get inside subjects for area %s: %w", areaID, err)
	}
	return states, nil
}

func (r *geofenceStateRepository) GetInsideBySubject(
	ctx context.Context,
	subjectID string,
	limit int,
) ([]*models.GeofenceState, error) {
	var states []*models.GeofenceState

	query := r.Pool().DB(ctx, true).
		Where("subject_id = ? AND inside = ?", subjectID, true).
		Order("modified_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&states).Error; err != nil {
		return nil, fmt.Errorf("get inside areas for subject %s: %w", subjectID, err)
	}
	return states, nil
}

func (r *geofenceStateRepository) DeleteByArea(ctx context.Context, areaID string) error {
	if err := r.Pool().DB(ctx, false).
		Where("area_id = ?", areaID).
		Delete(&models.GeofenceState{}).Error; err != nil {
		return fmt.Errorf("delete geofence states for area %s: %w", areaID, err)
	}
	return nil
}
