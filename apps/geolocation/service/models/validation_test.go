package models //nolint:testpackage // tests access unexported validation helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	geolocationv1 "buf.build/gen/go/antinvestor/geolocation/protocolbuffers/go/geolocation/v1"
)

func TestValidateLocationPoint(t *testing.T) {
	t.Parallel()

	valid := &LocationPointInput{
		Timestamp: timestamppb.New(time.Now().Add(-time.Minute)),
		DeviceId:  "device-1",
		Latitude:  0.3,
		Longitude: 32.5,
		Accuracy:  5,
		Source:    geolocationv1.LocationSource_LOCATION_SOURCE_GPS,
	}
	require.NoError(t, ValidateLocationPoint(valid))

	invalid := &LocationPointInput{
		DeviceId:  "device-1",
		Latitude:  200,
		Longitude: 32.5,
		Accuracy:  5,
		Source:    geolocationv1.LocationSource_LOCATION_SOURCE_GPS,
	}
	require.Error(t, ValidateLocationPoint(invalid))

	invalid.Source = 99
	invalid.Latitude = 0.3
	require.Error(t, ValidateLocationPoint(invalid))
}

func TestValidateGeoJSONAndRouteGeoJSON(t *testing.T) {
	t.Parallel()

	require.NoError(t, ValidateGeoJSON(`{"type":"Polygon","coordinates":[[[32,0],[33,0],[33,1],[32,0]]]}`))
	require.NoError(t, ValidateRouteGeoJSON(`{"type":"LineString","coordinates":[[32,0],[33,1]]}`))

	require.Error(t, ValidateGeoJSON(`{"type":"Point","coordinates":[32,0]}`))
	require.Error(t, ValidateRouteGeoJSON(`{"type":"Polygon","coordinates":[]}`))
}

func TestPrimitiveValidation(t *testing.T) {
	t.Parallel()

	require.NoError(t, ValidateLatLon(10, 20))
	require.Error(t, ValidateLatLon(100, 20))

	require.NoError(t, ValidateAccuracy(0))
	require.Error(t, ValidateAccuracy(-1))
	require.Error(t, ValidateAccuracy(MaxAccuracyMeters+1))

	require.NoError(t, ValidateTimestamp(time.Now()))
	require.Error(t, ValidateTimestamp(time.Time{}))
	require.Error(t, ValidateTimestamp(time.Now().Add(MaxClockSkew+time.Minute)))

	require.NoError(t, ValidateSubjectID("abc"))
	require.Error(t, ValidateSubjectID("ab"))
	require.Error(t, ValidateSubjectID("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
	require.NoError(t, ValidateDeviceID("device-123"))
	require.Error(t, ValidateDeviceID("dv"))
	require.Error(
		t,
		ValidateDeviceID("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghi"),
	)
}

func TestProtoConversionsAndToAPI(t *testing.T) {
	t.Parallel()

	require.Equal(t, geolocationv1.LocationSource_LOCATION_SOURCE_MANUAL, ToProtoLocationSource(LocationSourceManual))
	require.Equal(
		t,
		LocationSourceGPS,
		LocationSourceFromProto(geolocationv1.LocationSource_LOCATION_SOURCE_UNSPECIFIED),
	)
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_ZONE, ToProtoAreaType(AreaTypeZone))
	require.Equal(t, AreaTypeFence, AreaTypeFromProto(geolocationv1.AreaType_AREA_TYPE_FENCE))
	require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_DWELL, ToProtoGeoEventType(GeoEventTypeDwell))

	lp := &LocationPoint{
		SubjectID:     "subject-1",
		DeviceID:      "device-1",
		TrueCreatedAt: time.Now(),
		Latitude:      1.2,
		Longitude:     2.3,
		Accuracy:      4,
		Source:        LocationSourceNetwork,
	}
	lp.ID = "lp-1"
	api := lp.ToAPI()
	require.Equal(t, "lp-1", api.GetId())
	require.Equal(t, "device-1", api.GetDeviceId())
	require.Equal(t, geolocationv1.LocationSource_LOCATION_SOURCE_NETWORK, api.GetSource())

	area := &Area{
		OwnerID:      "owner-1",
		Name:         "area",
		AreaType:     AreaTypeFence,
		GeometryJSON: "{}",
		State:        2,
	}
	area.ID = "area-1"
	areaAPI := area.ToAPI()
	require.Equal(t, "area-1", areaAPI.GetId())
	require.Equal(t, geolocationv1.AreaType_AREA_TYPE_FENCE, areaAPI.GetAreaType())

	event := &GeoEvent{
		SubjectID:     "subject-1",
		AreaID:        "area-1",
		EventType:     GeoEventTypeEnter,
		TrueCreatedAt: time.Now(),
	}
	event.ID = "event-1"
	eventAPI := event.ToAPI()
	require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_ENTER, eventAPI.GetEventType())
}
