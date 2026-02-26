package models

import (
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// API types mirror the proto message definitions that will be published at
// buf.build/antinvestor/geolocation. Once the proto package exists, these
// can be replaced with the generated types. Until then they serve as the
// canonical API contract.

// LocationPointAPI is the API representation of a location point.
type LocationPointAPI struct {
	ID        string                 `json:"id"`
	SubjectID string                 `json:"subject_id"`
	Timestamp *timestamppb.Timestamp `json:"timestamp"`
	Latitude  float64                `json:"latitude"`
	Longitude float64                `json:"longitude"`
	Altitude  float64                `json:"altitude,omitempty"`
	Accuracy  float64                `json:"accuracy"`
	Speed     float64                `json:"speed,omitempty"`
	Bearing   float64                `json:"bearing,omitempty"`
	Source    LocationSource         `json:"source"`
	Extras    *structpb.Struct       `json:"extras,omitempty"`
	CreatedAt *timestamppb.Timestamp `json:"created_at,omitempty"`
}

// AreaAPI is the API representation of a geographic area.
type AreaAPI struct {
	ID           string                 `json:"id"`
	OwnerID      string                 `json:"owner_id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	AreaType     AreaType               `json:"area_type"`
	GeometryJSON string                 `json:"geometry"`
	AreaM2       float64                `json:"area_m2,omitempty"`
	PerimeterM   float64                `json:"perimeter_m,omitempty"`
	State        int32                  `json:"state"`
	Extras       *structpb.Struct       `json:"extras,omitempty"`
	CreatedAt    *timestamppb.Timestamp `json:"created_at,omitempty"`
}

// GeoEventAPI is the API representation of a spatial event.
type GeoEventAPI struct {
	ID         string                 `json:"id"`
	SubjectID  string                 `json:"subject_id"`
	AreaID     string                 `json:"area_id"`
	EventType  GeoEventType           `json:"event_type"`
	Timestamp  *timestamppb.Timestamp `json:"timestamp"`
	Confidence float64                `json:"confidence"`
	PointID    string                 `json:"point_id,omitempty"`
	Extras     *structpb.Struct       `json:"extras,omitempty"`
}

// NearbySubjectAPI represents a subject found within a proximity query.
type NearbySubjectAPI struct {
	SubjectID      string                 `json:"subject_id"`
	DistanceMeters float64                `json:"distance_meters"`
	LastSeen       *timestamppb.Timestamp `json:"last_seen"`
}

// NearbyAreaAPI represents an area found within a proximity query.
type NearbyAreaAPI struct {
	AreaID         string   `json:"area_id"`
	Name           string   `json:"name"`
	AreaType       AreaType `json:"area_type"`
	DistanceMeters float64  `json:"distance_meters"`
}

// IngestLocationsRequest is the batch ingestion request.
type IngestLocationsRequest struct {
	SubjectID string              `json:"subject_id"`
	Points    []*LocationPointAPI `json:"points"`
}

// IngestLocationsResponse is the batch ingestion response.
type IngestLocationsResponse struct {
	Accepted int32 `json:"accepted"`
	Rejected int32 `json:"rejected"`
}

// CreateAreaRequest is the request to create a new area.
type CreateAreaRequest struct {
	Data *AreaAPI `json:"data"`
}

// CreateAreaResponse is the response after creating an area.
type CreateAreaResponse struct {
	Data *AreaAPI `json:"data"`
}

// UpdateAreaRequest is the request to update an existing area.
type UpdateAreaRequest struct {
	ID          string           `json:"id"`
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Geometry    string           `json:"geometry,omitempty"`
	AreaType    *AreaType        `json:"area_type,omitempty"`
	Extras      *structpb.Struct `json:"extras,omitempty"`
}

// UpdateAreaResponse is the response after updating an area.
type UpdateAreaResponse struct {
	Data *AreaAPI `json:"data"`
}

// GetTrackRequest is the request for location history.
type GetTrackRequest struct {
	SubjectID string                 `json:"subject_id"`
	From      *timestamppb.Timestamp `json:"from"`
	To        *timestamppb.Timestamp `json:"to"`
	Limit     int32                  `json:"limit,omitempty"`
	Offset    int32                  `json:"offset,omitempty"`
}

// GetSubjectEventsRequest is the request for a subject's geo events.
type GetSubjectEventsRequest struct {
	SubjectID string                 `json:"subject_id"`
	From      *timestamppb.Timestamp `json:"from,omitempty"`
	To        *timestamppb.Timestamp `json:"to,omitempty"`
	Limit     int32                  `json:"limit,omitempty"`
	Offset    int32                  `json:"offset,omitempty"`
}

// GetAreaSubjectsRequest is the request to find subjects inside an area.
type GetAreaSubjectsRequest struct {
	AreaID string `json:"area_id"`
}

// AreaSubjectAPI represents a subject currently inside an area.
type AreaSubjectAPI struct {
	SubjectID      string                 `json:"subject_id"`
	EnterTimestamp *timestamppb.Timestamp `json:"enter_timestamp,omitempty"`
}

// GetNearbySubjectsRequest is a proximity query for nearby subjects.
type GetNearbySubjectsRequest struct {
	SubjectID    string  `json:"subject_id"`
	RadiusMeters float64 `json:"radius_meters"`
	Limit        int32   `json:"limit,omitempty"`
}

// GetNearbyAreasRequest is a proximity query for nearby areas.
type GetNearbyAreasRequest struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters float64 `json:"radius_meters"`
	Limit        int32   `json:"limit,omitempty"`
}

// LocationPointIngestedEvent is the NATS event payload when a point is ingested.
type LocationPointIngestedEvent struct {
	PointID   string  `json:"point_id"`
	SubjectID string  `json:"subject_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp int64   `json:"timestamp"` // Unix millis
}

// AreaChangedEvent is the NATS event payload when an area is created/updated/deleted.
type AreaChangedEvent struct {
	AreaID  string `json:"area_id"`
	Action  string `json:"action"` // "created", "updated", "deleted"
	OwnerID string `json:"owner_id"`
}

// RouteAPI is the API representation of a route.
type RouteAPI struct {
	ID                        string                 `json:"id"`
	OwnerID                   string                 `json:"owner_id"`
	Name                      string                 `json:"name"`
	Description               string                 `json:"description,omitempty"`
	GeometryJSON              string                 `json:"geometry"`
	LengthM                   float64                `json:"length_m,omitempty"`
	State                     int32                  `json:"state"`
	DeviationThresholdM       *float64               `json:"deviation_threshold_m,omitempty"`
	DeviationConsecutiveCount *int                   `json:"deviation_consecutive_count,omitempty"`
	DeviationCooldownSec      *int                   `json:"deviation_cooldown_sec,omitempty"`
	Extras                    *structpb.Struct       `json:"extras,omitempty"`
	CreatedAt                 *timestamppb.Timestamp `json:"created_at,omitempty"`
}

// RouteAssignmentAPI is the API representation of a route assignment.
type RouteAssignmentAPI struct {
	ID         string                 `json:"id"`
	SubjectID  string                 `json:"subject_id"`
	RouteID    string                 `json:"route_id"`
	ValidFrom  *timestamppb.Timestamp `json:"valid_from,omitempty"`
	ValidUntil *timestamppb.Timestamp `json:"valid_until,omitempty"`
	State      int32                  `json:"state"`
	Extras     *structpb.Struct       `json:"extras,omitempty"`
	CreatedAt  *timestamppb.Timestamp `json:"created_at,omitempty"`
}

// CreateRouteRequest is the request to create a new route.
type CreateRouteRequest struct {
	Data *RouteAPI `json:"data"`
}

// CreateRouteResponse is the response after creating a route.
type CreateRouteResponse struct {
	Data *RouteAPI `json:"data"`
}

// UpdateRouteRequest is the request to update an existing route.
type UpdateRouteRequest struct {
	ID                        string           `json:"id"`
	Name                      string           `json:"name,omitempty"`
	Description               string           `json:"description,omitempty"`
	Geometry                  string           `json:"geometry,omitempty"`
	DeviationThresholdM       *float64         `json:"deviation_threshold_m,omitempty"`
	DeviationConsecutiveCount *int             `json:"deviation_consecutive_count,omitempty"`
	DeviationCooldownSec      *int             `json:"deviation_cooldown_sec,omitempty"`
	Extras                    *structpb.Struct `json:"extras,omitempty"`
}

// UpdateRouteResponse is the response after updating a route.
type UpdateRouteResponse struct {
	Data *RouteAPI `json:"data"`
}

// AssignRouteRequest is the request to assign a subject to a route.
type AssignRouteRequest struct {
	SubjectID  string                 `json:"subject_id"`
	RouteID    string                 `json:"route_id"`
	ValidFrom  *timestamppb.Timestamp `json:"valid_from,omitempty"`
	ValidUntil *timestamppb.Timestamp `json:"valid_until,omitempty"`
}

// AssignRouteResponse is the response after assigning a subject to a route.
type AssignRouteResponse struct {
	Data *RouteAssignmentAPI `json:"data"`
}

// UnassignRouteRequest is the request to remove a route assignment.
type UnassignRouteRequest struct {
	AssignmentID string `json:"assignment_id"`
}

// GetSubjectRouteAssignmentsRequest is the request to get a subject's route assignments.
type GetSubjectRouteAssignmentsRequest struct {
	SubjectID string `json:"subject_id"`
}

// RouteDeviationDetectedEvent is the event payload when route deviation is detected.
type RouteDeviationDetectedEvent struct {
	SubjectID      string  `json:"subject_id"`
	RouteID        string  `json:"route_id"`
	EventType      string  `json:"event_type"` // "deviated" or "back_on_route"
	DistanceMeters float64 `json:"distance_meters"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Timestamp      int64   `json:"timestamp"` // Unix millis
}

// RouteChangedEvent is the event payload when a route is created/updated/deleted.
type RouteChangedEvent struct {
	RouteID string `json:"route_id"`
	Action  string `json:"action"` // "created", "updated", "deleted"
	OwnerID string `json:"owner_id"`
}

// GeoEventEmitted is the NATS event payload when a geofence event occurs.
type GeoEventEmitted struct {
	EventID    string       `json:"event_id"`
	SubjectID  string       `json:"subject_id"`
	AreaID     string       `json:"area_id"`
	EventType  GeoEventType `json:"event_type"`
	Timestamp  int64        `json:"timestamp"` // Unix millis
	Confidence float64      `json:"confidence"`
}
