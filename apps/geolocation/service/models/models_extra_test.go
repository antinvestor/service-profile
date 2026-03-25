package models //nolint:testpackage // tests access unexported model helpers

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/security"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	geolocationv1 "github.com/antinvestor/service-profile/proto/geolocation/geolocation/v1"
)

func TestModelStringsAndTableNames(t *testing.T) {
	t.Parallel()

	require.Equal(t, "GPS", LocationSourceGPS.String())
	require.Equal(t, "NETWORK", LocationSourceNetwork.String())
	require.Equal(t, "IP", LocationSourceIP.String())
	require.Equal(t, "MANUAL", LocationSourceManual.String())
	require.Contains(t, LocationSource(99).String(), "UNKNOWN")
	require.Equal(t, "LAND", AreaTypeLand.String())
	require.Equal(t, "BUILDING", AreaTypeBuilding.String())
	require.Equal(t, "ZONE", AreaTypeZone.String())
	require.Equal(t, "FENCE", AreaTypeFence.String())
	require.Equal(t, "CUSTOM", AreaTypeCustom.String())
	require.Contains(t, AreaType(99).String(), "UNKNOWN")
	require.Equal(t, "ENTER", GeoEventTypeEnter.String())
	require.Equal(t, "EXIT", GeoEventTypeExit.String())
	require.Equal(t, "DWELL", GeoEventTypeDwell.String())
	require.Contains(t, GeoEventType(99).String(), "UNKNOWN")

	require.Equal(t, "location_points", (&LocationPoint{}).TableName())
	require.Equal(t, "areas", (&Area{}).TableName())
	require.Equal(t, "geo_events", (&GeoEvent{}).TableName())
	require.Equal(t, "geofence_states", (&GeofenceState{}).TableName())
	require.Equal(t, "latest_positions", (&LatestPosition{}).TableName())
	require.Equal(t, "routes", (&Route{}).TableName())
	require.Equal(t, "route_assignments", (&RouteAssignment{}).TableName())
	require.Equal(t, "route_deviation_states", (&RouteDeviationState{}).TableName())
}

func TestModelConversionsAndHelpers(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	altitude := 123.4
	speed := 5.6
	bearing := 7.8
	length := 12.5
	threshold := 55.0
	consecutive := 3
	cooldown := 120
	validFrom := now.Add(-time.Hour)
	validUntil := now.Add(time.Hour)

	route := &Route{
		OwnerID:                   "owner-1",
		Name:                      "route",
		Description:               "desc",
		GeometryJSON:              "{}",
		LengthM:                   &length,
		State:                     2,
		Extras:                    data.JSONMap{"key": "value"},
		DeviationThresholdM:       &threshold,
		DeviationConsecutiveCount: &consecutive,
		DeviationCooldownSec:      &cooldown,
	}
	route.ID = "route-1"
	route.CreatedAt = now
	require.True(t, route.HasDeviationConfig())
	routeAPI := route.ToAPI()
	require.Equal(t, "route-1", routeAPI.GetId())
	require.Equal(t, int32(consecutive), routeAPI.GetDeviationConsecutiveCount())
	require.Equal(t, int32(cooldown), routeAPI.GetDeviationCooldownSec())
	require.Equal(t, "value", routeAPI.GetExtra().GetFields()["key"].GetStringValue())

	assignment := &RouteAssignment{
		SubjectID:  "subject-1",
		RouteID:    "route-1",
		ValidFrom:  &validFrom,
		ValidUntil: &validUntil,
		State:      2,
		Extras:     data.JSONMap{"mode": "test"},
	}
	assignment.ID = "assignment-1"
	assignment.CreatedAt = now
	assignmentAPI := assignment.ToAPI()
	require.Equal(t, "assignment-1", assignmentAPI.GetId())
	require.Equal(t, "test", assignmentAPI.GetExtra().GetFields()["mode"].GetStringValue())
	require.Equal(t, validFrom.Unix(), assignmentAPI.GetValidFrom().AsTime().Unix())
	require.Equal(t, validUntil.Unix(), assignmentAPI.GetValidUntil().AsTime().Unix())

	require.Nil(t, jsonMapToStruct(nil))
	require.Nil(t, StructToJSONMap(nil))

	extraStruct, err := structpb.NewStruct(map[string]any{"a": "b"})
	require.NoError(t, err)
	require.Equal(t, "b", StructToJSONMap(extraStruct)["a"])
	require.Equal(t, "b", jsonMapToStruct(data.JSONMap{"a": "b"}).GetFields()["a"].GetStringValue())

	point := &LocationPoint{
		SubjectID: "subject-1",
		DeviceID:  "device-1",
		TS:        now,
		Latitude:  1.1,
		Longitude: 2.2,
		Altitude:  &altitude,
		Accuracy:  3.3,
		Speed:     &speed,
		Bearing:   &bearing,
		Source:    LocationSourceManual,
		Extras:    data.JSONMap{"mode": "manual"},
	}
	point.ID = "point-1"
	point.CreatedAt = now
	pointAPI := point.ToAPI()
	require.Equal(t, "point-1", pointAPI.GetId())
	require.Equal(t, "device-1", pointAPI.GetDeviceId())
	require.InDelta(t, altitude, pointAPI.GetAltitude(), 0.001)
	require.InDelta(t, speed, pointAPI.GetSpeed(), 0.001)
	require.InDelta(t, bearing, pointAPI.GetBearing(), 0.001)
	require.Equal(t, "manual", pointAPI.GetExtra().GetFields()["mode"].GetStringValue())

	areaM2 := 44.0
	perimeter := 22.0
	area := &Area{
		OwnerID:      "owner-1",
		Name:         "area",
		Description:  "desc",
		AreaType:     AreaTypeFence,
		GeometryJSON: "{}",
		AreaM2:       &areaM2,
		PerimeterM:   &perimeter,
		State:        2,
		Extras:       data.JSONMap{"scope": "yard"},
	}
	area.ID = "area-1"
	area.CreatedAt = now
	areaAPI := area.ToAPI()
	require.InDelta(t, areaM2, areaAPI.GetAreaM2(), 0.001)
	require.InDelta(t, perimeter, areaAPI.GetPerimeterM(), 0.001)
	require.Equal(t, "yard", areaAPI.GetExtra().GetFields()["scope"].GetStringValue())
}

func TestContextAndProtoConversions(t *testing.T) {
	t.Parallel()

	ctx := ContextWithEventTenancy(context.Background(), EventTenancy{
		TenantID:    "tenant-1",
		PartitionID: "partition-1",
		AccessID:    "access-1",
	}, "subject-1")

	claims := security.ClaimsFromContext(ctx)
	require.NotNil(t, claims)
	require.Equal(t, "tenant-1", claims.GetTenantID())
	subject, err := claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, "subject-1", subject)

	require.Equal(
		t,
		geolocationv1.LocationSource_LOCATION_SOURCE_UNSPECIFIED,
		ToProtoLocationSource(LocationSource(99)),
	)
	require.Equal(t, geolocationv1.LocationSource_LOCATION_SOURCE_GPS, ToProtoLocationSource(LocationSourceGPS))
	require.Equal(t, geolocationv1.LocationSource_LOCATION_SOURCE_NETWORK, ToProtoLocationSource(LocationSourceNetwork))
	require.Equal(t, geolocationv1.LocationSource_LOCATION_SOURCE_IP, ToProtoLocationSource(LocationSourceIP))
	require.Equal(
		t,
		LocationSourceNetwork,
		LocationSourceFromProto(geolocationv1.LocationSource_LOCATION_SOURCE_NETWORK),
	)
	require.Equal(t, LocationSourceIP, LocationSourceFromProto(geolocationv1.LocationSource_LOCATION_SOURCE_IP))
	require.Equal(t, LocationSourceManual, LocationSourceFromProto(geolocationv1.LocationSource_LOCATION_SOURCE_MANUAL))
	require.Equal(t, LocationSourceGPS, LocationSourceFromProto(geolocationv1.LocationSource(99)))
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_UNSPECIFIED, ToProtoAreaType(AreaType(99)))
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_LAND, ToProtoAreaType(AreaTypeLand))
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_BUILDING, ToProtoAreaType(AreaTypeBuilding))
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_FENCE, ToProtoAreaType(AreaTypeFence))
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_CUSTOM, ToProtoAreaType(AreaTypeCustom))
	require.Equal(t, AreaTypeLand, AreaTypeFromProto(geolocationv1.AreaType(99)))
	require.Equal(t, AreaTypeBuilding, AreaTypeFromProto(geolocationv1.AreaType_AREA_TYPE_BUILDING))
	require.Equal(t, AreaTypeZone, AreaTypeFromProto(geolocationv1.AreaType_AREA_TYPE_ZONE))
	require.Equal(t, AreaTypeCustom, AreaTypeFromProto(geolocationv1.AreaType_AREA_TYPE_CUSTOM))
	require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_UNSPECIFIED, ToProtoGeoEventType(GeoEventType(99)))
	require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_ENTER, ToProtoGeoEventType(GeoEventTypeEnter))
	require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_EXIT, ToProtoGeoEventType(GeoEventTypeExit))
}

func TestValidationExtraBranches(t *testing.T) {
	t.Parallel()

	require.NoError(t, ValidateRouteName("route"))
	require.Error(t, ValidateRouteName(""))
	require.Error(t, ValidateRouteName(strings.Repeat("r", MaxRouteNameLength+1)))

	require.NoError(t, ValidateAreaName("area"))
	require.Error(t, ValidateAreaName(""))
	require.Error(t, ValidateAreaName(strings.Repeat("a", MaxAreaNameLength+1)))

	require.Error(t, ValidateGeoJSON(`{"coordinates":[]}`))
	require.Error(t, ValidateGeoJSON(`{"type":"Triangle","coordinates":[]}`))
	require.Error(t, ValidateGeoJSON(`{"type":"Polygon"}`))
	require.Error(t, ValidateGeoJSON(""))
	require.Error(t, ValidateGeoJSON("not-json"))

	require.Error(t, ValidateRouteGeoJSON(`{"coordinates":[]}`))
	require.Error(t, ValidateRouteGeoJSON(`{"type":"Point","coordinates":[]}`))
	require.Error(t, ValidateRouteGeoJSON(`{"type":"LineString","coordinates":[[32,0]]}`))
	require.Error(t, ValidateRouteGeoJSON(""))
	require.Error(t, ValidateRouteGeoJSON("not-json"))

	require.Equal(t, 0, countVerticesRecursive("nope"))
	require.Equal(t, 1, countGeoJSONVertices([]byte(`[1,2]`)))
	require.Error(t, ValidateLatLon(math.NaN(), 1))
	require.Error(t, ValidateLatLon(1, math.Inf(1)))
	require.Error(t, ValidateLocationPoint(nil))
	require.Error(t, ValidateLocationPoint(&LocationPointInput{
		DeviceId:  "device-1",
		Latitude:  1,
		Longitude: 1,
		Accuracy:  1,
		Timestamp: nil,
		Source:    geolocationv1.LocationSource(99),
	}))
}

func TestContextWithEventTenancyNoop(t *testing.T) {
	t.Parallel()

	base := (&security.AuthenticationClaims{
		TenantID:    "tenant-0",
		PartitionID: "partition-0",
		AccessID:    "access-0",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "subject-0",
		},
	}).ClaimsToContext(context.Background())

	require.Same(t, base, ContextWithEventTenancy(base, EventTenancy{}, "subject-1"))
}
