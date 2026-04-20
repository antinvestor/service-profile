package models

import (
	"fmt"
	"time"

	"github.com/pitabwire/frame/data"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// LocationSource indicates how a location was determined.
type LocationSource int32

const (
	LocationSourceGPS     LocationSource = 0
	LocationSourceNetwork LocationSource = 1
	LocationSourceIP      LocationSource = 2
	LocationSourceManual  LocationSource = 3
)

func (ls LocationSource) String() string {
	switch ls {
	case LocationSourceGPS:
		return "GPS"
	case LocationSourceNetwork:
		return "NETWORK"
	case LocationSourceIP:
		return "IP"
	case LocationSourceManual:
		return "MANUAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", ls)
	}
}

// AreaType classifies the kind of area geometry.
type AreaType int32

const (
	AreaTypeLand     AreaType = 0
	AreaTypeBuilding AreaType = 1
	AreaTypeZone     AreaType = 2
	AreaTypeFence    AreaType = 3
	AreaTypeCustom   AreaType = 4
)

func (at AreaType) String() string {
	switch at {
	case AreaTypeLand:
		return "LAND"
	case AreaTypeBuilding:
		return "BUILDING"
	case AreaTypeZone:
		return "ZONE"
	case AreaTypeFence:
		return "FENCE"
	case AreaTypeCustom:
		return "CUSTOM"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", at)
	}
}

// GeoEventType identifies the type of spatial event.
type GeoEventType int32

const (
	GeoEventTypeEnter GeoEventType = 0
	GeoEventTypeExit  GeoEventType = 1
	GeoEventTypeDwell GeoEventType = 2
)

func (et GeoEventType) String() string {
	switch et {
	case GeoEventTypeEnter:
		return "ENTER"
	case GeoEventTypeExit:
		return "EXIT"
	case GeoEventTypeDwell:
		return "DWELL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", et)
	}
}

type LocationPointProcessingState int32

const (
	LocationPointProcessingStatePending LocationPointProcessingState = iota
	LocationPointProcessingStateProcessed
	LocationPointProcessingStateFailed
)

// LocationPoint represents a single location observation from a device or subject.
// The geom column is computed by a database trigger from latitude/longitude.
type LocationPoint struct {
	data.BaseModel

	SubjectID       string                       `gorm:"type:varchar(40);not null;index:idx_lp_subject_true_created_at"`
	DeviceID        string                       `gorm:"type:varchar(80);not null;index:idx_lp_device_true_created_at"`
	TrueCreatedAt   time.Time                    `gorm:"type:timestamptz;not null;index:idx_lp_subject_true_created_at,sort:desc;column:true_created_at"`
	IngestedAt      time.Time                    `gorm:"type:timestamptz;not null;default:now()"`
	Latitude        float64                      `gorm:"type:double precision;not null"`
	Longitude       float64                      `gorm:"type:double precision;not null"`
	Altitude        *float64                     `gorm:"type:double precision"`
	Accuracy        float64                      `gorm:"type:double precision;not null;default:0"`
	Speed           *float64                     `gorm:"type:double precision"`
	Bearing         *float64                     `gorm:"type:double precision"`
	Source          LocationSource               `gorm:"type:smallint;not null;default:0"`
	Extras          data.JSONMap                 `gorm:"serializer:json;type:jsonb;default:'{}'"`
	ProcessingState LocationPointProcessingState `gorm:"type:smallint;not null;default:0;index:idx_lp_processing_state"`
	ProcessedAt     *time.Time                   `gorm:"type:timestamptz"`
	ProcessingError string                       `gorm:"type:text;not null;default:''"`
	Geom            *string                      `gorm:"type:geometry(Point,4326);column:geom"                                                           json:"-"`
}

func (*LocationPoint) TableName() string {
	return "location_points"
}

// ToAPI converts a LocationPoint to an API-compatible map structure.
func (lp *LocationPoint) ToAPI() *LocationPointAPI {
	api := &LocationPointAPI{
		Id:        lp.GetID(),
		SubjectId: lp.SubjectID,
		DeviceId:  lp.DeviceID,
		Timestamp: timestamppb.New(lp.TrueCreatedAt),
		Latitude:  lp.Latitude,
		Longitude: lp.Longitude,
		Accuracy:  lp.Accuracy,
		Source:    ToProtoLocationSource(lp.Source),
		Extra:     jsonMapToStruct(lp.Extras),
		CreatedAt: timestamppb.New(lp.CreatedAt),
	}
	if lp.Altitude != nil {
		api.Altitude = lp.Altitude
	}
	if lp.Speed != nil {
		api.Speed = lp.Speed
	}
	if lp.Bearing != nil {
		api.Bearing = lp.Bearing
	}
	return api
}

// Area represents a geographic boundary (polygon/multipolygon) stored via PostGIS.
// The geom and bbox columns are managed by database triggers that compute from GeometryJSON.
type Area struct {
	data.BaseModel

	OwnerID      string   `gorm:"type:varchar(40);not null;index:idx_areas_owner"`
	Name         string   `gorm:"type:varchar(250);not null"`
	Description  string   `gorm:"type:text"`
	AreaType     AreaType `gorm:"type:smallint;not null;default:0"`
	GeometryJSON string   `gorm:"type:text;column:geometry_json"`

	// PostGIS spatial columns — computed by DB triggers, queried by raw SQL.
	Geom *string `gorm:"type:geometry(Geometry,4326);column:geom" json:"-"`
	Bbox *string `gorm:"type:geometry(Polygon,4326);column:bbox"  json:"-"`

	// Computed fields (populated by DB trigger on geom column write).
	AreaM2     *float64     `gorm:"type:double precision;column:area_m2"`
	PerimeterM *float64     `gorm:"type:double precision;column:perimeter_m"`
	State      int32        `gorm:"type:smallint;not null;default:0"`
	Extras     data.JSONMap `gorm:"serializer:json;type:jsonb;default:'{}'"`
}

func (*Area) TableName() string {
	return "areas"
}

func (a *Area) ToAPI() *AreaAPI {
	api := &AreaAPI{
		Id:          a.GetID(),
		OwnerId:     a.OwnerID,
		Name:        a.Name,
		Description: a.Description,
		AreaType:    ToProtoAreaType(a.AreaType),
		Geometry:    a.GeometryJSON,
		State:       a.State,
		Extra:       jsonMapToStruct(a.Extras),
		CreatedAt:   timestamppb.New(a.CreatedAt),
	}
	if a.AreaM2 != nil {
		api.AreaM2 = *a.AreaM2
	}
	if a.PerimeterM != nil {
		api.PerimeterM = *a.PerimeterM
	}
	return api
}

// GeoEvent represents a spatial transition event (enter/exit/dwell).
// Append-only — never updated after creation.
type GeoEvent struct {
	data.BaseModel

	SubjectID     string       `gorm:"type:varchar(40);not null;index:idx_ge_subject_true_created_at"`
	AreaID        string       `gorm:"type:varchar(40);not null;index:idx_ge_area_true_created_at"`
	EventType     GeoEventType `gorm:"type:smallint;not null"`
	TrueCreatedAt time.Time    `gorm:"type:timestamptz;not null;index:idx_ge_subject_true_created_at,sort:desc;index:idx_ge_area_true_created_at,sort:desc;column:true_created_at"`
	Confidence    float64      `gorm:"type:double precision;not null;default:1.0"`
	PointID       string       `gorm:"type:varchar(40)"`
	Extras        data.JSONMap `gorm:"serializer:json;type:jsonb;default:'{}'"`
}

func (*GeoEvent) TableName() string {
	return "geo_events"
}

func (ge *GeoEvent) ToAPI() *GeoEventAPI {
	return &GeoEventAPI{
		Id:         ge.GetID(),
		SubjectId:  ge.SubjectID,
		AreaId:     ge.AreaID,
		EventType:  ToProtoGeoEventType(ge.EventType),
		Timestamp:  timestamppb.New(ge.TrueCreatedAt),
		Confidence: ge.Confidence,
		PointId:    ge.PointID,
		Extra:      jsonMapToStruct(ge.Extras),
	}
}

// GeofenceState tracks the current inside/outside state for a (subject, area) pair.
// This is the mutable state that drives enter/exit detection.
type GeofenceState struct {
	data.BaseModel

	SubjectID      string     `gorm:"type:varchar(40);not null;index:idx_geofence_state_subject_area,unique"`
	AreaID         string     `gorm:"type:varchar(40);not null;index:idx_geofence_state_subject_area,unique"`
	Inside         bool       `gorm:"not null;default:false"`
	LastTransition *time.Time `gorm:"type:timestamptz"`
	LastPointTS    *time.Time `gorm:"type:timestamptz;column:last_point_ts"`
	EnterTS        *time.Time `gorm:"type:timestamptz;column:enter_ts"`

	// Last evaluated position (lat/lon, not PostGIS column — spatial queries use raw SQL).
	LastLat float64 `gorm:"type:double precision"`
	LastLon float64 `gorm:"type:double precision"`
}

func (*GeofenceState) TableName() string {
	return "geofence_states"
}

// LatestPosition is a materialized "most recent position" per subject.
// Used for proximity queries. Maintained by the detection engine.
// The geom column is computed by a database trigger from latitude/longitude.
type LatestPosition struct {
	data.BaseModel

	SubjectID string    `gorm:"type:varchar(40);not null;index:idx_latest_position_subject_tenant,unique"`
	DeviceID  string    `gorm:"type:varchar(80);not null"`
	Latitude  float64   `gorm:"type:double precision;not null"`
	Longitude float64   `gorm:"type:double precision;not null"`
	Accuracy  float64   `gorm:"type:double precision;not null;default:0"`
	TS        time.Time `gorm:"type:timestamptz;not null;column:ts"`
	Geom      *string   `gorm:"type:geometry(Point,4326);column:geom"                                     json:"-"`
}

func (*LatestPosition) TableName() string {
	return "latest_positions"
}

// Route represents a predefined path (LineString) stored via PostGIS.
// The geom column is computed by a database trigger from GeometryJSON.
type Route struct {
	data.BaseModel

	OwnerID      string `gorm:"type:varchar(40);not null;index:idx_routes_owner"`
	Name         string `gorm:"type:varchar(250);not null"`
	Description  string `gorm:"type:text"`
	GeometryJSON string `gorm:"type:text;column:geometry_json"`

	// PostGIS spatial column — computed by DB trigger, queried by raw SQL.
	Geom *string `gorm:"type:geometry(LineString,4326);column:geom" json:"-"`

	// Computed field (populated by DB trigger on geom column write).
	LengthM *float64 `gorm:"type:double precision;column:length_m"`

	State  int32        `gorm:"type:smallint;not null;default:0"`
	Extras data.JSONMap `gorm:"serializer:json;type:jsonb;default:'{}'"`

	// Per-route deviation thresholds. NULL means skip during evaluation.
	DeviationThresholdM       *float64 `gorm:"type:double precision;column:deviation_threshold_m"`
	DeviationConsecutiveCount *int     `gorm:"type:integer;column:deviation_consecutive_count"`
	DeviationCooldownSec      *int     `gorm:"type:integer;column:deviation_cooldown_sec"`
}

func (*Route) TableName() string {
	return "routes"
}

// HasDeviationConfig returns true if the route has deviation detection configured.
func (r *Route) HasDeviationConfig() bool {
	return r.DeviationThresholdM != nil
}

func (r *Route) ToAPI() *RouteAPI {
	api := &RouteAPI{
		Id:          r.GetID(),
		OwnerId:     r.OwnerID,
		Name:        r.Name,
		Description: r.Description,
		Geometry:    r.GeometryJSON,
		State:       r.State,
		Extra:       jsonMapToStruct(r.Extras),
		CreatedAt:   timestamppb.New(r.CreatedAt),
	}
	if r.LengthM != nil {
		api.LengthM = *r.LengthM
	}
	api.DeviationThresholdM = r.DeviationThresholdM
	if r.DeviationConsecutiveCount != nil {
		//nolint:gosec // bounded by business validation
		value := int32(*r.DeviationConsecutiveCount)
		api.DeviationConsecutiveCount = &value
	}
	if r.DeviationCooldownSec != nil {
		value := int32(*r.DeviationCooldownSec) //nolint:gosec // bounded by business validation
		api.DeviationCooldownSec = &value
	}
	return api
}

// RouteAssignment maps a subject to a route with an optional time window.
type RouteAssignment struct {
	data.BaseModel

	SubjectID  string       `gorm:"type:varchar(40);not null;index:idx_ra_subject_state"`
	RouteID    string       `gorm:"type:varchar(40);not null;index:idx_ra_route"`
	ValidFrom  *time.Time   `gorm:"type:timestamptz"`
	ValidUntil *time.Time   `gorm:"type:timestamptz"`
	State      int32        `gorm:"type:smallint;not null;default:0"`
	Extras     data.JSONMap `gorm:"serializer:json;type:jsonb;default:'{}'"`
}

func (*RouteAssignment) TableName() string {
	return "route_assignments"
}

func (ra *RouteAssignment) ToAPI() *RouteAssignmentAPI {
	api := &RouteAssignmentAPI{
		Id:        ra.GetID(),
		SubjectId: ra.SubjectID,
		RouteId:   ra.RouteID,
		State:     ra.State,
		Extra:     jsonMapToStruct(ra.Extras),
		CreatedAt: timestamppb.New(ra.CreatedAt),
	}
	if ra.ValidFrom != nil {
		api.ValidFrom = timestamppb.New(*ra.ValidFrom)
	}
	if ra.ValidUntil != nil {
		api.ValidUntil = timestamppb.New(*ra.ValidUntil)
	}
	return api
}

// RouteDeviationState tracks the current on/off-route state for a (subject, route) pair.
// Composite PK (no BaseModel), same pattern as GeofenceState.
type RouteDeviationState struct {
	data.BaseModel

	SubjectID            string     `gorm:"type:varchar(40);not null;index:idx_route_deviation_state_subject_route,unique"`
	RouteID              string     `gorm:"type:varchar(40);not null;index:idx_route_deviation_state_subject_route,unique"`
	Deviated             bool       `gorm:"not null;default:false"`
	ConsecutiveOffRoute  int        `gorm:"not null;default:0"`
	LastDeviationEventAt *time.Time `gorm:"type:timestamptz"`
	LastPointTS          *time.Time `gorm:"type:timestamptz;column:last_point_ts"`
	LastLat              float64    `gorm:"type:double precision"`
	LastLon              float64    `gorm:"type:double precision"`
}

func (*RouteDeviationState) TableName() string {
	return "route_deviation_states"
}

// jsonMapToStruct converts a data.JSONMap to a protobuf Struct.
func jsonMapToStruct(m data.JSONMap) *structpb.Struct {
	if m == nil {
		return nil
	}
	s, _ := structpb.NewStruct(map[string]any(m))
	return s
}

// StructToJSONMap converts a protobuf Struct to data.JSONMap.
func StructToJSONMap(s *structpb.Struct) data.JSONMap {
	if s == nil {
		return nil
	}
	return data.JSONMap(s.AsMap())
}
