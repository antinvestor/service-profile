package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"
)

const (
	// MaxLatitude is the maximum valid latitude (90 degrees).
	MaxLatitude = 90.0
	// MinLatitude is the minimum valid latitude (-90 degrees).
	MinLatitude = -90.0
	// MaxLongitude is the maximum valid longitude (180 degrees).
	MaxLongitude = 180.0
	// MinLongitude is the minimum valid longitude (-180 degrees).
	MinLongitude = -180.0

	// MaxAccuracyMeters is the maximum acceptable accuracy value.
	MaxAccuracyMeters = 10000.0

	// MaxClockSkew is the maximum clock drift tolerated for future timestamps.
	MaxClockSkew = 5 * time.Minute

	// MaxBatchSize is the maximum number of points in a single ingestion batch.
	MaxBatchSize = 1000

	// MaxGeoJSONVertices is the maximum number of vertices allowed in an area geometry.
	MaxGeoJSONVertices = 10000

	// MaxRouteGeoJSONVertices is the maximum number of vertices allowed in a route geometry.
	MaxRouteGeoJSONVertices = 50000

	// MinRouteVertices is the minimum number of vertices required in a route geometry.
	MinRouteVertices = 2

	// MaxRouteNameLength is the maximum length of a route name.
	MaxRouteNameLength = 250

	// MaxAreaSqMeters is the maximum area allowed (roughly 10,000 sq km).
	MaxAreaSqMeters = 10_000_000_000.0

	// DefaultProximityLimit is the default limit for proximity queries.
	DefaultProximityLimit = 100

	// MaxProximityRadiusMeters is the maximum proximity search radius.
	MaxProximityRadiusMeters = 50_000.0

	// DefaultTrackLimit is the default limit for track queries.
	DefaultTrackLimit = 1000

	// MaxTrackLimit is the maximum limit for track queries.
	MaxTrackLimit = 10000

	// MinSubjectIDLength is the minimum length of a subject ID.
	MinSubjectIDLength = 3

	// MaxSubjectIDLength is the maximum length of a subject ID.
	MaxSubjectIDLength = 40

	// MaxAreaNameLength is the maximum length of an area name.
	MaxAreaNameLength = 250
)

// ValidateLatLon checks that latitude and longitude are within valid ranges.
func ValidateLatLon(lat, lon float64) error {
	if lat < MinLatitude || lat > MaxLatitude {
		return fmt.Errorf("latitude %f out of range [%f, %f]", lat, MinLatitude, MaxLatitude)
	}
	if lon < MinLongitude || lon > MaxLongitude {
		return fmt.Errorf("longitude %f out of range [%f, %f]", lon, MinLongitude, MaxLongitude)
	}
	if math.IsNaN(lat) || math.IsNaN(lon) || math.IsInf(lat, 0) || math.IsInf(lon, 0) {
		return errors.New("latitude and longitude must be finite numbers")
	}
	return nil
}

// ValidateAccuracy checks that accuracy is a positive, reasonable value.
func ValidateAccuracy(accuracy float64) error {
	if accuracy < 0 {
		return fmt.Errorf("accuracy must be non-negative, got %f", accuracy)
	}
	if accuracy > MaxAccuracyMeters {
		return fmt.Errorf("accuracy %f exceeds maximum %f meters", accuracy, MaxAccuracyMeters)
	}
	return nil
}

// ValidateTimestamp checks that a point's timestamp is not too far in the future.
func ValidateTimestamp(ts time.Time) error {
	if ts.IsZero() {
		return errors.New("timestamp must not be zero")
	}
	if ts.After(time.Now().Add(MaxClockSkew)) {
		return fmt.Errorf(
			"timestamp %s is too far in the future (max clock skew: %s)",
			ts,
			MaxClockSkew,
		)
	}
	return nil
}

// ValidateSubjectID checks that a subject ID is non-empty and within bounds.
func ValidateSubjectID(id string) error {
	if len(id) < MinSubjectIDLength {
		return errors.New("subject_id must be at least 3 characters")
	}
	if len(id) > MaxSubjectIDLength {
		return errors.New("subject_id must be at most 40 characters")
	}
	return nil
}

// ValidateLocationPoint performs full validation on an incoming location point.
func ValidateLocationPoint(pt *LocationPointAPI) error {
	if pt == nil {
		return errors.New("location point is nil")
	}
	if err := ValidateSubjectID(pt.SubjectID); err != nil {
		return fmt.Errorf("invalid subject_id: %w", err)
	}
	if err := ValidateLatLon(pt.Latitude, pt.Longitude); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}
	if err := ValidateAccuracy(pt.Accuracy); err != nil {
		return fmt.Errorf("invalid accuracy: %w", err)
	}
	if pt.Timestamp != nil {
		if err := ValidateTimestamp(pt.Timestamp.AsTime()); err != nil {
			return fmt.Errorf("invalid timestamp: %w", err)
		}
	}
	if pt.Source < LocationSourceGPS || pt.Source > LocationSourceManual {
		return fmt.Errorf("invalid source: %d", pt.Source)
	}
	return nil
}

// ValidateGeoJSON performs basic structural validation on a GeoJSON string.
// It checks that the JSON is valid, has a recognized type, and coordinates are present.
// Full geometry validation (closed rings, no self-intersection, etc.) is done by PostGIS.
func ValidateGeoJSON(geoJSON string) error {
	if geoJSON == "" {
		return errors.New("geometry is empty")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(geoJSON), &parsed); err != nil {
		return fmt.Errorf("geometry is not valid JSON: %w", err)
	}

	geoType, ok := parsed["type"].(string)
	if !ok {
		return errors.New("geometry missing 'type' field")
	}

	switch geoType {
	case "Polygon", "MultiPolygon":
		// valid types for areas
	case "Point":
		return errors.New("area geometry must be Polygon or MultiPolygon, not Point")
	default:
		return fmt.Errorf(
			"unsupported geometry type: %s (expected Polygon or MultiPolygon)",
			geoType,
		)
	}

	coords, ok := parsed["coordinates"]
	if !ok || coords == nil {
		return errors.New("geometry missing 'coordinates' field")
	}

	// Count approximate vertices by marshaling coordinates back and counting arrays.
	coordBytes, err := json.Marshal(coords)
	if err != nil {
		return fmt.Errorf("cannot re-marshal coordinates: %w", err)
	}
	vertexCount := countGeoJSONVertices(coordBytes)
	if vertexCount > MaxGeoJSONVertices {
		return fmt.Errorf(
			"geometry has %d vertices, maximum allowed is %d",
			vertexCount,
			MaxGeoJSONVertices,
		)
	}

	return nil
}

// countGeoJSONVertices gives an approximate count of coordinate pairs in GeoJSON coordinates.
// It counts the innermost arrays (pairs of numbers).
func countGeoJSONVertices(coordBytes []byte) int {
	var raw any
	if err := json.Unmarshal(coordBytes, &raw); err != nil {
		return 0
	}
	return countVerticesRecursive(raw)
}

func countVerticesRecursive(v any) int {
	arr, ok := v.([]any)
	if !ok {
		return 0
	}
	if len(arr) == 0 {
		return 0
	}

	// If the first element is a number, this is a coordinate pair/triple.
	if _, isNum := arr[0].(float64); isNum {
		return 1
	}

	count := 0
	for _, elem := range arr {
		count += countVerticesRecursive(elem)
	}
	return count
}

// ValidateRouteName checks that a route name meets requirements.
func ValidateRouteName(name string) error {
	if len(name) < 1 {
		return errors.New("route name must not be empty")
	}
	if len(name) > MaxRouteNameLength {
		return errors.New("route name must be at most 250 characters")
	}
	return nil
}

// ValidateRouteGeoJSON validates a GeoJSON string for route use (LineString).
func ValidateRouteGeoJSON(geoJSON string) error {
	if geoJSON == "" {
		return errors.New("route geometry is empty")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(geoJSON), &parsed); err != nil {
		return fmt.Errorf("route geometry is not valid JSON: %w", err)
	}

	geoType, ok := parsed["type"].(string)
	if !ok {
		return errors.New("route geometry missing 'type' field")
	}

	if geoType != "LineString" {
		return fmt.Errorf("route geometry must be LineString, got %s", geoType)
	}

	coords, ok := parsed["coordinates"]
	if !ok || coords == nil {
		return errors.New("route geometry missing 'coordinates' field")
	}

	coordBytes, err := json.Marshal(coords)
	if err != nil {
		return fmt.Errorf("cannot re-marshal coordinates: %w", err)
	}
	vertexCount := countGeoJSONVertices(coordBytes)
	if vertexCount < MinRouteVertices {
		return fmt.Errorf(
			"route geometry must have at least %d vertices, got %d",
			MinRouteVertices, vertexCount,
		)
	}
	if vertexCount > MaxRouteGeoJSONVertices {
		return fmt.Errorf(
			"route geometry has %d vertices, maximum allowed is %d",
			vertexCount, MaxRouteGeoJSONVertices,
		)
	}

	return nil
}

// ValidateAreaName checks that an area name meets requirements.
func ValidateAreaName(name string) error {
	if len(name) < 1 {
		return errors.New("area name must not be empty")
	}
	if len(name) > MaxAreaNameLength {
		return errors.New("area name must be at most 250 characters")
	}
	return nil
}
