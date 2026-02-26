package business

import (
	"context"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

// IngestionBusiness handles location point ingestion: validation, normalization, persistence, and event emission.
type IngestionBusiness interface {
	// IngestBatch validates, normalizes, persists, and emits events for a batch of location points.
	IngestBatch(ctx context.Context, req *models.IngestLocationsRequest) (*models.IngestLocationsResponse, error)
}

// AreaBusiness handles area CRUD with geometry validation and event emission.
type AreaBusiness interface {
	CreateArea(ctx context.Context, req *models.CreateAreaRequest) (*models.AreaAPI, error)
	UpdateArea(ctx context.Context, req *models.UpdateAreaRequest) (*models.AreaAPI, error)
	DeleteArea(ctx context.Context, areaID string) error
	GetArea(ctx context.Context, areaID string) (*models.AreaAPI, error)
	SearchAreas(ctx context.Context, query string, ownerID string, limit int) ([]*models.AreaAPI, error)
}

// GeofenceBusiness implements the stateful geofence detection engine.
type GeofenceBusiness interface { //nolint:iface // intentionally mirrors RouteDeviationBusiness pattern
	// EvaluatePoint runs the geofence state machine for a single location point.
	// It checks all candidate areas, updates state, and emits enter/exit/dwell events.
	EvaluatePoint(ctx context.Context, event *models.LocationPointIngestedEvent) error
}

// ProximityBusiness handles proximity queries.
type ProximityBusiness interface {
	// GetNearbySubjects finds subjects near the given subject.
	GetNearbySubjects(ctx context.Context, req *models.GetNearbySubjectsRequest) ([]*models.NearbySubjectAPI, error)
	// GetNearbyAreas finds areas near the given point.
	GetNearbyAreas(ctx context.Context, req *models.GetNearbyAreasRequest) ([]*models.NearbyAreaAPI, error)
	// UpdateLatestPosition updates the materialized latest position for a subject.
	UpdateLatestPosition(ctx context.Context, event *models.LocationPointIngestedEvent) error
}

// RouteDeviationBusiness implements the stateful route deviation detection engine.
type RouteDeviationBusiness interface { //nolint:iface // intentionally mirrors GeofenceBusiness pattern
	// EvaluatePoint runs the route deviation state machine for a single location point.
	EvaluatePoint(ctx context.Context, event *models.LocationPointIngestedEvent) error
}

// RouteBusiness handles route CRUD, assignments, and event emission.
type RouteBusiness interface {
	CreateRoute(ctx context.Context, req *models.CreateRouteRequest) (*models.RouteAPI, error)
	UpdateRoute(ctx context.Context, req *models.UpdateRouteRequest) (*models.RouteAPI, error)
	DeleteRoute(ctx context.Context, routeID string) error
	GetRoute(ctx context.Context, routeID string) (*models.RouteAPI, error)
	SearchRoutes(ctx context.Context, ownerID string, limit int) ([]*models.RouteAPI, error)
	AssignRoute(ctx context.Context, req *models.AssignRouteRequest) (*models.RouteAssignmentAPI, error)
	UnassignRoute(ctx context.Context, assignmentID string) error
	GetSubjectAssignments(
		ctx context.Context,
		subjectID string,
	) ([]*models.RouteAssignmentAPI, error)
}

// TrackBusiness handles location history and event queries.
type TrackBusiness interface {
	GetTrack(ctx context.Context, req *models.GetTrackRequest) ([]*models.LocationPointAPI, error)
	GetSubjectEvents(ctx context.Context, req *models.GetSubjectEventsRequest) ([]*models.GeoEventAPI, error)
	GetAreaSubjects(ctx context.Context, req *models.GetAreaSubjectsRequest) ([]*models.AreaSubjectAPI, error)
}
