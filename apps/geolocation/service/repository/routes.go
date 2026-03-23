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

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

type routeRepository struct {
	datastore.BaseRepository[*models.Route]
}

// NewRouteRepository creates a new repository for routes.
func NewRouteRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) RouteRepository {
	return &routeRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Route](
			ctx, dbPool, workMan, func() *models.Route { return &models.Route{} },
		),
	}
}

// UpdateGeometry sets the PostGIS geom column from GeoJSON for a route.
func (r *routeRepository) UpdateGeometry(
	ctx context.Context,
	routeID string,
	geoJSON string,
) error {
	db := r.Pool().DB(ctx, false)
	return executeUpdateRouteGeometry(db, routeID, geoJSON)
}

// UpdateGeometryTx sets the PostGIS geom column within an existing transaction.
func (r *routeRepository) UpdateGeometryTx(
	tx *gorm.DB,
	routeID string,
	geoJSON string,
) error {
	return executeUpdateRouteGeometry(tx, routeID, geoJSON)
}

func executeUpdateRouteGeometry(db *gorm.DB, routeID string, geoJSON string) error {
	result := db.Table((&models.Route{}).TableName()).
		Where("id = ?", routeID).
		Updates(map[string]any{
			"geom":          gorm.Expr("ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)", geoJSON),
			"geometry_json": geoJSON,
		})
	if result.Error != nil {
		return fmt.Errorf("update geometry for route %s: %w", routeID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("route %s not found or deleted", routeID)
	}
	return nil
}

// GetActiveAssignmentsForSubject returns active route assignments for a subject,
// joined with routes that have non-null deviation config and valid geometry.
func (r *routeRepository) GetActiveAssignmentsForSubject(
	ctx context.Context,
	subjectID string,
	at time.Time,
) ([]*RouteWithAssignment, error) {
	db := r.Pool().DB(ctx, true)

	type joinResult struct {
		models.Route
		AssignmentID string     `gorm:"column:assignment_id"`
		SubjectID    string     `gorm:"column:subject_id"`
		ValidFrom    *time.Time `gorm:"column:valid_from"`
		ValidUntil   *time.Time `gorm:"column:valid_until"`
	}

	var results []joinResult

	query := db.Table((&models.RouteAssignment{}).TableName()).
		Select(
			"routes.*, route_assignments.id AS assignment_id, route_assignments.subject_id, "+
				"route_assignments.valid_from, route_assignments.valid_until",
		).
		Joins("JOIN routes ON routes.id = route_assignments.route_id").
		Where("route_assignments.subject_id = ?", subjectID).
		Where("route_assignments.state = ?", 2). //nolint:mnd // 2 = active state
		Where("routes.state = ?", 2).            //nolint:mnd // 2 = active state
		Where("routes.geom IS NOT NULL").
		Where("routes.deviation_threshold_m IS NOT NULL").
		Where("(route_assignments.valid_from IS NULL OR route_assignments.valid_from <= ?)", at).
		Where("(route_assignments.valid_until IS NULL OR route_assignments.valid_until >= ?)", at)

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf(
			"get active route assignments for subject %s: %w",
			subjectID, err,
		)
	}

	out := make([]*RouteWithAssignment, len(results))
	for i := range results {
		route := results[i].Route
		out[i] = &RouteWithAssignment{
			Route: &route,
			Assignment: &models.RouteAssignment{
				SubjectID:  results[i].SubjectID,
				ValidFrom:  results[i].ValidFrom,
				ValidUntil: results[i].ValidUntil,
			},
		}
		out[i].Assignment.ID = results[i].AssignmentID
	}
	return out, nil
}

// DistanceToRouteMeters computes the distance in meters from a point to a route.
func (r *routeRepository) DistanceToRouteMeters(
	ctx context.Context,
	routeID string,
	lat, lon float64,
) (float64, error) {
	var distance float64

	db := r.Pool().DB(ctx, true)
	result := db.Table((&models.Route{}).TableName()).
		Select(
			"ST_Distance(geom::geography, ST_SetSRID(ST_Point(?, ?), 4326)::geography)",
			lon,
			lat,
		).
		Where("id = ?", routeID).
		Scan(&distance)

	if result.Error != nil {
		return 0, fmt.Errorf(
			"distance to route %s from (%f, %f): %w",
			routeID, lat, lon, result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return 0, errors.New("route not found")
	}
	return distance, nil
}

// SearchByOwner returns routes owned by the given owner, with a limit.
func (r *routeRepository) SearchByOwner(
	ctx context.Context,
	ownerID string,
	limit int,
) ([]*models.Route, error) {
	var routes []*models.Route
	db := r.Pool().DB(ctx, true)
	query := db.Where("owner_id = ? AND deleted_at IS NULL", ownerID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	result := query.Find(&routes)
	if result.Error != nil {
		return nil, fmt.Errorf("search routes by owner %s: %w", ownerID, result.Error)
	}
	return routes, nil
}

// --- Route Assignment Repository ---

type routeAssignmentRepository struct {
	datastore.BaseRepository[*models.RouteAssignment]
}

// NewRouteAssignmentRepository creates a new repository for route assignments.
func NewRouteAssignmentRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) RouteAssignmentRepository {
	return &routeAssignmentRepository{
		BaseRepository: datastore.NewBaseRepository[*models.RouteAssignment](
			ctx, dbPool, workMan,
			func() *models.RouteAssignment { return &models.RouteAssignment{} },
		),
	}
}

// GetBySubject returns active assignments for a subject.
func (r *routeAssignmentRepository) GetBySubject(
	ctx context.Context,
	subjectID string,
) ([]*models.RouteAssignment, error) {
	var assignments []*models.RouteAssignment
	db := r.Pool().DB(ctx, true)
	result := db.Where(
		"subject_id = ? AND state = 2 AND deleted_at IS NULL",
		subjectID,
	).Find(&assignments)
	if result.Error != nil {
		return nil, fmt.Errorf(
			"get route assignments for subject %s: %w",
			subjectID, result.Error,
		)
	}
	return assignments, nil
}

// DeleteByRoute removes all assignments for a given route.
func (r *routeAssignmentRepository) DeleteByRoute(
	ctx context.Context,
	routeID string,
) error {
	db := r.Pool().DB(ctx, false)
	result := db.Where("route_id = ?", routeID).Delete(&models.RouteAssignment{})
	if result.Error != nil {
		return fmt.Errorf(
			"delete route assignments for route %s: %w",
			routeID, result.Error,
		)
	}
	return nil
}
