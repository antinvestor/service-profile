package models

import (
	"fmt"
	"time"

	"github.com/pitabwire/frame/data"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
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

// LocationPoint represents a single location observation from a device or subject.
// Stored in the geo.location_points table with PostGIS POINT geometry.
type LocationPoint struct {
	data.BaseModel

	SubjectID  string         `gorm:"type:varchar(40);not null;index:idx_lp_subject_ts"`
	TS         time.Time      `gorm:"type:timestamptz;not null;index:idx_lp_subject_ts,sort:desc;column:ts"`
	IngestedAt time.Time      `gorm:"type:timestamptz;not null;default:now()"`
	Latitude   float64        `gorm:"type:double precision;not null"`
	Longitude  float64        `gorm:"type:double precision;not null"`
	Altitude   *float64       `gorm:"type:double precision"`
	Accuracy   float64        `gorm:"type:double precision;not null;default:0"`
	Speed      *float64       `gorm:"type:double precision"`
	Bearing    *float64       `gorm:"type:double precision"`
	Source     LocationSource `gorm:"type:smallint;not null;default:0"`
	Extras     data.JSONMap   `gorm:"type:jsonb;default:'{}'"`
}

func (*LocationPoint) TableName() string {
	return "location_points"
}

// ToAPI converts a LocationPoint to an API-compatible map structure.
func (lp *LocationPoint) ToAPI() *LocationPointAPI {
	api := &LocationPointAPI{
		ID:        lp.GetID(),
		SubjectID: lp.SubjectID,
		Timestamp: timestamppb.New(lp.TS),
		Latitude:  lp.Latitude,
		Longitude: lp.Longitude,
		Accuracy:  lp.Accuracy,
		Source:    lp.Source,
		Extras:    jsonMapToStruct(lp.Extras),
		CreatedAt: timestamppb.New(lp.CreatedAt),
	}
	if lp.Altitude != nil {
		api.Altitude = *lp.Altitude
	}
	if lp.Speed != nil {
		api.Speed = *lp.Speed
	}
	if lp.Bearing != nil {
		api.Bearing = *lp.Bearing
	}
	return api
}

// Area represents a geographic boundary (polygon/multipolygon) stored via PostGIS.
// The geometry column (geom) is managed via raw SQL migrations, not GORM auto-migrate.
type Area struct {
	data.BaseModel

	OwnerID     string   `gorm:"type:varchar(40);not null;index:idx_areas_owner"`
	Name        string   `gorm:"type:varchar(250);not null"`
	Description string   `gorm:"type:text"`
	AreaType    AreaType `gorm:"type:smallint;not null;default:0"`

	// GeoJSON representation of the geometry; the PostGIS geom column is managed via SQL.
	// This field is used for API input/output; the actual spatial column is "geom".
	GeometryJSON string `gorm:"type:text;column:geometry_json"`

	// Computed fields (populated by DB trigger on geom column write).
	AreaM2     *float64     `gorm:"type:double precision;column:area_m2"`
	PerimeterM *float64     `gorm:"type:double precision;column:perimeter_m"`
	State      int32        `gorm:"type:smallint;not null;default:0"`
	Extras     data.JSONMap `gorm:"type:jsonb;default:'{}'"`
	ModifiedAt time.Time    `gorm:"type:timestamptz;not null;default:now()"`
}

func (*Area) TableName() string {
	return "areas"
}

func (a *Area) BeforeUpdate(_ *gorm.DB) error {
	a.ModifiedAt = time.Now()
	return nil
}

func (a *Area) ToAPI() *AreaAPI {
	api := &AreaAPI{
		ID:           a.GetID(),
		OwnerID:      a.OwnerID,
		Name:         a.Name,
		Description:  a.Description,
		AreaType:     a.AreaType,
		GeometryJSON: a.GeometryJSON,
		State:        a.State,
		Extras:       jsonMapToStruct(a.Extras),
		CreatedAt:    timestamppb.New(a.CreatedAt),
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

	SubjectID  string       `gorm:"type:varchar(40);not null;index:idx_ge_subject_ts"`
	AreaID     string       `gorm:"type:varchar(40);not null;index:idx_ge_area_ts"`
	EventType  GeoEventType `gorm:"type:smallint;not null"`
	TS         time.Time    `gorm:"type:timestamptz;not null;index:idx_ge_subject_ts,sort:desc;index:idx_ge_area_ts,sort:desc;column:ts"`
	Confidence float64      `gorm:"type:double precision;not null;default:1.0"`
	PointID    string       `gorm:"type:varchar(40)"`
	Extras     data.JSONMap `gorm:"type:jsonb;default:'{}'"`
}

func (*GeoEvent) TableName() string {
	return "geo_events"
}

func (ge *GeoEvent) ToAPI() *GeoEventAPI {
	return &GeoEventAPI{
		ID:         ge.GetID(),
		SubjectID:  ge.SubjectID,
		AreaID:     ge.AreaID,
		EventType:  ge.EventType,
		Timestamp:  timestamppb.New(ge.TS),
		Confidence: ge.Confidence,
		PointID:    ge.PointID,
		Extras:     jsonMapToStruct(ge.Extras),
	}
}

// GeofenceState tracks the current inside/outside state for a (subject, area) pair.
// This is the mutable state that drives enter/exit detection.
type GeofenceState struct {
	SubjectID      string     `gorm:"type:varchar(40);not null;primaryKey"`
	AreaID         string     `gorm:"type:varchar(40);not null;primaryKey"`
	Inside         bool       `gorm:"not null;default:false"`
	LastTransition *time.Time `gorm:"type:timestamptz"`
	LastPointTS    *time.Time `gorm:"type:timestamptz;column:last_point_ts"`
	EnterTS        *time.Time `gorm:"type:timestamptz;column:enter_ts"`

	// Last evaluated position (lat/lon, not PostGIS column — spatial queries use raw SQL).
	LastLat float64 `gorm:"type:double precision"`
	LastLon float64 `gorm:"type:double precision"`

	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
}

func (*GeofenceState) TableName() string {
	return "geofence_states"
}

func (gs *GeofenceState) BeforeUpdate(_ *gorm.DB) error {
	gs.UpdatedAt = time.Now()
	return nil
}

// LatestPosition is a materialized "most recent position" per subject.
// Used for proximity queries. Maintained by the detection engine.
type LatestPosition struct {
	SubjectID string    `gorm:"type:varchar(40);primaryKey"`
	Latitude  float64   `gorm:"type:double precision;not null"`
	Longitude float64   `gorm:"type:double precision;not null"`
	Accuracy  float64   `gorm:"type:double precision;not null;default:0"`
	TS        time.Time `gorm:"type:timestamptz;not null;column:ts"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
}

func (*LatestPosition) TableName() string {
	return "latest_positions"
}

func (lp *LatestPosition) BeforeUpdate(_ *gorm.DB) error {
	lp.UpdatedAt = time.Now()
	return nil
}

// Route represents a predefined path (LineString) stored via PostGIS.
// The geometry column (geom) is managed via raw SQL migrations, not GORM auto-migrate.
type Route struct {
	data.BaseModel

	OwnerID     string `gorm:"type:varchar(40);not null;index:idx_routes_owner"`
	Name        string `gorm:"type:varchar(250);not null"`
	Description string `gorm:"type:text"`

	// GeoJSON representation of the geometry; the PostGIS geom column is managed via SQL.
	GeometryJSON string `gorm:"type:text;column:geometry_json"`

	// Computed field (populated by DB trigger on geom column write).
	LengthM *float64 `gorm:"type:double precision;column:length_m"`

	State      int32        `gorm:"type:smallint;not null;default:0"`
	Extras     data.JSONMap `gorm:"type:jsonb;default:'{}'"`
	ModifiedAt time.Time    `gorm:"type:timestamptz;not null;default:now()"`

	// Per-route deviation thresholds. NULL means skip during evaluation.
	DeviationThresholdM       *float64 `gorm:"type:double precision;column:deviation_threshold_m"`
	DeviationConsecutiveCount *int     `gorm:"type:integer;column:deviation_consecutive_count"`
	DeviationCooldownSec      *int     `gorm:"type:integer;column:deviation_cooldown_sec"`
}

func (*Route) TableName() string {
	return "routes"
}

func (r *Route) BeforeUpdate(_ *gorm.DB) error {
	r.ModifiedAt = time.Now()
	return nil
}

// HasDeviationConfig returns true if the route has deviation detection configured.
func (r *Route) HasDeviationConfig() bool {
	return r.DeviationThresholdM != nil
}

func (r *Route) ToAPI() *RouteAPI {
	api := &RouteAPI{
		ID:           r.GetID(),
		OwnerID:      r.OwnerID,
		Name:         r.Name,
		Description:  r.Description,
		GeometryJSON: r.GeometryJSON,
		State:        r.State,
		Extras:       jsonMapToStruct(r.Extras),
		CreatedAt:    timestamppb.New(r.CreatedAt),
	}
	if r.LengthM != nil {
		api.LengthM = *r.LengthM
	}
	api.DeviationThresholdM = r.DeviationThresholdM
	api.DeviationConsecutiveCount = r.DeviationConsecutiveCount
	api.DeviationCooldownSec = r.DeviationCooldownSec
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
	Extras     data.JSONMap `gorm:"type:jsonb;default:'{}'"`
}

func (*RouteAssignment) TableName() string {
	return "route_assignments"
}

func (ra *RouteAssignment) ToAPI() *RouteAssignmentAPI {
	api := &RouteAssignmentAPI{
		ID:        ra.GetID(),
		SubjectID: ra.SubjectID,
		RouteID:   ra.RouteID,
		State:     ra.State,
		Extras:    jsonMapToStruct(ra.Extras),
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
	SubjectID            string     `gorm:"type:varchar(40);not null;primaryKey"`
	RouteID              string     `gorm:"type:varchar(40);not null;primaryKey"`
	Deviated             bool       `gorm:"not null;default:false"`
	ConsecutiveOffRoute  int        `gorm:"not null;default:0"`
	LastDeviationEventAt *time.Time `gorm:"type:timestamptz"`
	LastPointTS          *time.Time `gorm:"type:timestamptz;column:last_point_ts"`
	LastLat              float64    `gorm:"type:double precision"`
	LastLon              float64    `gorm:"type:double precision"`
	UpdatedAt            time.Time  `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
}

func (*RouteDeviationState) TableName() string {
	return "route_deviation_states"
}

func (rds *RouteDeviationState) BeforeUpdate(_ *gorm.DB) error {
	rds.UpdatedAt = time.Now()
	return nil
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
		return make(data.JSONMap)
	}
	return data.JSONMap(s.AsMap())
}
