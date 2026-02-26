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

type geoEventRepository struct {
	datastore.BaseRepository[*models.GeoEvent]
}

// NewGeoEventRepository creates a new repository for geo events.
func NewGeoEventRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) GeoEventRepository {
	return &geoEventRepository{
		BaseRepository: datastore.NewBaseRepository[*models.GeoEvent](
			ctx, dbPool, workMan, func() *models.GeoEvent { return &models.GeoEvent{} },
		),
	}
}

// CreateTx creates a geo event within an existing transaction.
func (r *geoEventRepository) CreateTx(tx *gorm.DB, event *models.GeoEvent) error {
	if err := tx.Create(event).Error; err != nil {
		return fmt.Errorf("create geo event in tx: %w", err)
	}
	return nil
}

// GetBySubject returns geo events for a subject within an optional time range.
func (r *geoEventRepository) GetBySubject(
	ctx context.Context,
	subjectID string,
	from, to *time.Time,
	limit, offset int,
) ([]*models.GeoEvent, error) {
	return r.queryEvents(ctx, "subject_id", subjectID, from, to, limit, offset)
}

// GetByArea returns geo events for an area within an optional time range.
func (r *geoEventRepository) GetByArea(
	ctx context.Context,
	areaID string,
	from, to *time.Time,
	limit, offset int,
) ([]*models.GeoEvent, error) {
	return r.queryEvents(ctx, "area_id", areaID, from, to, limit, offset)
}

// HasDwellEvent checks if a DWELL event exists for the given subject/area after the given timestamp.
func (r *geoEventRepository) HasDwellEvent(
	ctx context.Context,
	subjectID, areaID string,
	from *time.Time,
) (bool, error) {
	db := r.Pool().DB(ctx, true)
	return hasDwellEventQuery(db, subjectID, areaID, from)
}

// HasDwellEventTx checks for a DWELL event within an existing transaction.
// Must be used inside the geofence SELECT FOR UPDATE transaction to prevent duplicate dwell events.
func (r *geoEventRepository) HasDwellEventTx(
	tx *gorm.DB,
	subjectID, areaID string,
	from *time.Time,
) (bool, error) {
	return hasDwellEventQuery(tx, subjectID, areaID, from)
}

// hasDwellEventQuery is the shared implementation for dwell event existence checks.
// Uses EXISTS for O(1) short-circuit instead of COUNT(*).
func hasDwellEventQuery(
	db *gorm.DB,
	subjectID, areaID string,
	from *time.Time,
) (bool, error) {
	var exists bool
	query := db.Raw(
		buildDwellExistsQuery(from),
		buildDwellExistsArgs(subjectID, areaID, from)...,
	).Scan(&exists)

	if query.Error != nil {
		return false, fmt.Errorf(
			"check dwell event for subject %s area %s: %w",
			subjectID, areaID, query.Error,
		)
	}
	return exists, nil
}

func buildDwellExistsQuery(from *time.Time) string {
	base := "SELECT EXISTS(SELECT 1 FROM geo_events WHERE subject_id = ? AND area_id = ? AND event_type = ?"
	if from != nil {
		return base + " AND ts >= ? LIMIT 1)"
	}
	return base + " LIMIT 1)"
}

func buildDwellExistsArgs(subjectID, areaID string, from *time.Time) []any {
	args := []any{subjectID, areaID, models.GeoEventTypeDwell}
	if from != nil {
		args = append(args, *from)
	}
	return args
}

// queryEvents is a shared implementation for querying geo events by a single filter column.
func (r *geoEventRepository) queryEvents(
	ctx context.Context,
	filterCol, filterVal string,
	from, to *time.Time,
	limit, offset int,
) ([]*models.GeoEvent, error) {
	var events []*models.GeoEvent

	db := r.Pool().DB(ctx, true)
	query := db.Where(fmt.Sprintf("%s = ?", filterCol), filterVal)
	query = applyTimeRange(query, from, to)
	query = query.Order("ts DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("get geo events by %s=%s: %w", filterCol, filterVal, result.Error)
	}
	return events, nil
}

// applyTimeRange adds optional from/to time filters to a GORM query.
func applyTimeRange(query *gorm.DB, from, to *time.Time) *gorm.DB {
	if from != nil {
		query = query.Where("ts >= ?", *from)
	}
	if to != nil {
		query = query.Where("ts <= ?", *to)
	}
	return query
}
