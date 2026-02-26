package repository

import (
	"context"
	"time"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

// LocationPointRepository manages location point persistence.
type LocationPointRepository interface {
	datastore.BaseRepository[*models.LocationPoint]
	GetTrack(
		ctx context.Context,
		subjectID string,
		from, to time.Time,
		limit, offset int,
	) ([]*models.LocationPoint, error)
}

// AreaRepository manages area geometry persistence and spatial queries.
type AreaRepository interface {
	datastore.BaseRepository[*models.Area]
	// GetActiveByBoundingBox returns areas whose bounding box intersects the given point.
	// Uses PostGIS ST_Intersects on the bbox column for fast pre-filtering.
	GetActiveByBoundingBox(ctx context.Context, lat, lon float64) ([]*models.Area, error)
	// ContainsPoint checks if the area's actual geometry contains the given point.
	// Uses PostGIS ST_Contains on the geom column.
	ContainsPoint(ctx context.Context, areaID string, lat, lon float64) (bool, error)
	// UpdateGeometry sets the PostGIS geom column from GeoJSON and recomputes metrics.
	UpdateGeometry(ctx context.Context, areaID string, geoJSON string) error
	// UpdateGeometryTx sets the PostGIS geom column within an existing transaction.
	UpdateGeometryTx(tx *gorm.DB, areaID string, geoJSON string) error
	// GetNearbyAreas finds areas within radiusMeters of the given point.
	GetNearbyAreas(ctx context.Context, lat, lon, radiusMeters float64, limit int) ([]*AreaWithDistance, error)
	// SearchByOwner returns areas owned by the given owner ID with a limit.
	SearchByOwner(ctx context.Context, ownerID string, limit int) ([]*models.Area, error)
	// SearchByQuery performs text search on area name/description.
	SearchByQuery(ctx context.Context, query string, limit int) ([]*models.Area, error)
}

// GeoEventRepository manages geo event persistence.
type GeoEventRepository interface {
	datastore.BaseRepository[*models.GeoEvent]
	// CreateTx creates a geo event within an existing transaction.
	CreateTx(tx *gorm.DB, event *models.GeoEvent) error
	GetBySubject(
		ctx context.Context,
		subjectID string,
		from, to *time.Time,
		limit, offset int,
	) ([]*models.GeoEvent, error)
	GetByArea(ctx context.Context, areaID string, from, to *time.Time, limit, offset int) ([]*models.GeoEvent, error)
	// HasDwellEvent checks if a DWELL event exists for the given subject/area within a time range.
	// Used by the geofence engine to avoid duplicate dwell events.
	HasDwellEvent(ctx context.Context, subjectID, areaID string, from *time.Time) (bool, error)
	// HasDwellEventTx is the same check but runs within an existing transaction.
	// Must be used inside the geofence SELECT FOR UPDATE transaction to prevent duplicate dwell events.
	HasDwellEventTx(tx *gorm.DB, subjectID, areaID string, from *time.Time) (bool, error)
}

// GeofenceStateRepository manages mutable geofence state.
type GeofenceStateRepository interface {
	// UpsertTx creates or updates the geofence state within an existing transaction.
	UpsertTx(tx *gorm.DB, state *models.GeofenceState) error
	// GetForUpdate retrieves the geofence state with a row-level lock (SELECT FOR UPDATE).
	// Returns (nil, nil) if no state exists for the given (subject, area) pair.
	GetForUpdate(tx *gorm.DB, subjectID, areaID string) (*models.GeofenceState, error)
	// GetInsideByArea returns subjects currently inside the given area, with a limit.
	GetInsideByArea(ctx context.Context, areaID string, limit int) ([]*models.GeofenceState, error)
	// DeleteByArea removes all geofence state entries for a given area.
	DeleteByArea(ctx context.Context, areaID string) error
	// Pool returns the underlying database pool for transaction management.
	Pool() pool.Pool
}

// LatestPositionRepository manages the materialized latest position per subject.
type LatestPositionRepository interface {
	Upsert(ctx context.Context, pos *models.LatestPosition) error
	Get(ctx context.Context, subjectID string) (*models.LatestPosition, error)
	// GetNearbySubjects finds subjects within radiusMeters of the given point.
	GetNearbySubjects(
		ctx context.Context,
		lat, lon, radiusMeters float64,
		excludeSubjectID string,
		staleHours int,
		limit int,
	) ([]*SubjectWithDistance, error)
}

// RouteRepository manages route persistence and spatial queries.
type RouteRepository interface {
	datastore.BaseRepository[*models.Route]
	// UpdateGeometry sets the PostGIS geom column from GeoJSON.
	UpdateGeometry(ctx context.Context, routeID string, geoJSON string) error
	// UpdateGeometryTx sets the PostGIS geom column within an existing transaction.
	UpdateGeometryTx(tx *gorm.DB, routeID string, geoJSON string) error
	// GetActiveAssignmentsForSubject returns active route assignments for a subject,
	// joined with routes that have deviation config and valid geometry.
	GetActiveAssignmentsForSubject(
		ctx context.Context,
		subjectID string,
		at time.Time,
	) ([]*RouteWithAssignment, error)
	// DistanceToRouteMeters computes the distance in meters from a point to a route.
	DistanceToRouteMeters(ctx context.Context, routeID string, lat, lon float64) (float64, error)
	// SearchByOwner returns routes owned by the given owner ID.
	SearchByOwner(ctx context.Context, ownerID string, limit int) ([]*models.Route, error)
}

// RouteAssignmentRepository manages route assignment persistence.
type RouteAssignmentRepository interface {
	datastore.BaseRepository[*models.RouteAssignment]
	// GetBySubject returns active assignments for a subject.
	GetBySubject(ctx context.Context, subjectID string) ([]*models.RouteAssignment, error)
	// DeleteByRoute removes all assignments for a given route.
	DeleteByRoute(ctx context.Context, routeID string) error
}

// RouteDeviationStateRepository manages mutable route deviation state.
type RouteDeviationStateRepository interface {
	// UpsertTx creates or updates the deviation state within an existing transaction.
	UpsertTx(tx *gorm.DB, state *models.RouteDeviationState) error
	// GetForUpdate retrieves the deviation state with a row-level lock.
	// Returns (nil, nil) if no state exists.
	GetForUpdate(tx *gorm.DB, subjectID, routeID string) (*models.RouteDeviationState, error)
	// DeleteByRoute removes all deviation state entries for a given route.
	DeleteByRoute(ctx context.Context, routeID string) error
	// Pool returns the underlying database pool for transaction management.
	Pool() pool.Pool
}

// RouteWithAssignment pairs a route with its assignment for a subject.
type RouteWithAssignment struct {
	Route      *models.Route
	Assignment *models.RouteAssignment
}

// AreaWithDistance pairs an area with its distance from a query point.
type AreaWithDistance struct {
	Area           *models.Area
	DistanceMeters float64
}

// SubjectWithDistance pairs a subject with its distance from a query point.
type SubjectWithDistance struct {
	SubjectID      string
	Latitude       float64
	Longitude      float64
	DistanceMeters float64
	LastSeen       time.Time
}
