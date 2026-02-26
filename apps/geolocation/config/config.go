package config

import (
	"strings"
	"time"

	"github.com/pitabwire/frame/config"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/handlers"
)

// GeolocationConfig holds configuration for the geolocation service.
type GeolocationConfig struct {
	config.ConfigurationDefault

	// Geofence detection tuning.
	GeofenceHysteresisBufferM float64 `envDefault:"30.0"  env:"GEOFENCE_HYSTERESIS_BUFFER_M"`
	GeofenceDwellThresholdSec int     `envDefault:"120"   env:"GEOFENCE_DWELL_THRESHOLD_SEC"`
	GeofenceMaxAccuracyM      float64 `envDefault:"500.0" env:"GEOFENCE_MAX_ACCURACY_M"`
	GeofenceMaxCandidateAreas int     `envDefault:"100"   env:"GEOFENCE_MAX_CANDIDATE_AREAS"`

	// Route deviation detection tuning (global GPS accuracy filter only).
	RouteDeviationMaxAccuracyM float64 `envDefault:"500.0" env:"ROUTE_DEVIATION_MAX_ACCURACY_M"`

	// Proximity query defaults.
	ProximityDefaultRadiusM float64 `envDefault:"1000.0"  env:"PROXIMITY_DEFAULT_RADIUS_M"`
	ProximityMaxRadiusM     float64 `envDefault:"50000.0" env:"PROXIMITY_MAX_RADIUS_M"`
	ProximityStaleHours     int     `envDefault:"1"       env:"PROXIMITY_STALE_HOURS"`
	ProximityDefaultLimit   int     `envDefault:"50"      env:"PROXIMITY_DEFAULT_LIMIT"`
	ProximityMaxLimit       int     `envDefault:"200"     env:"PROXIMITY_MAX_LIMIT"`

	// Track/history query limits.
	TrackDefaultLimit int `envDefault:"100"  env:"TRACK_DEFAULT_LIMIT"`
	TrackMaxLimit     int `envDefault:"1000" env:"TRACK_MAX_LIMIT"`
	EventDefaultLimit int `envDefault:"100"  env:"EVENT_DEFAULT_LIMIT"`
	EventMaxLimit     int `envDefault:"500"  env:"EVENT_MAX_LIMIT"`
	MaxAreaSubjects   int `envDefault:"1000" env:"MAX_AREA_SUBJECTS"`

	// Ingestion limits.
	IngestionMaxBatchSize int `envDefault:"1000" env:"INGESTION_MAX_BATCH_SIZE"`

	// Search limits.
	SearchDefaultLimit int `envDefault:"50" env:"SEARCH_DEFAULT_LIMIT"`

	// Request body size limit in bytes (default 2MB).
	MaxRequestBodyBytes int64 `envDefault:"2097152" env:"MAX_REQUEST_BODY_BYTES"`

	// Data retention policy.
	RetentionLocationPointDays int `envDefault:"90"    env:"RETENTION_LOCATION_POINT_DAYS"`
	RetentionGeoEventDays      int `envDefault:"365"   env:"RETENTION_GEO_EVENT_DAYS"`
	RetentionGeofenceStateDays int `envDefault:"30"    env:"RETENTION_GEOFENCE_STATE_DAYS"`
	RetentionPartitionMonths   int `envDefault:"3"     env:"RETENTION_PARTITION_MONTHS"`
	RetentionBatchSize         int `envDefault:"10000" env:"RETENTION_BATCH_SIZE"`

	// Rate limiting.
	RateLimitRequestsPerMinute int    `envDefault:"600" env:"RATE_LIMIT_REQUESTS_PER_MINUTE"`
	RateLimitTrustedProxies    string `envDefault:""    env:"RATE_LIMIT_TRUSTED_PROXIES"`
}

// GeofenceBusinessConfig returns the GeofenceConfig derived from this configuration.
func (c *GeolocationConfig) GeofenceBusinessConfig() business.GeofenceConfig {
	return business.GeofenceConfig{
		HysteresisBufferM:  c.GeofenceHysteresisBufferM,
		DwellThreshold:     time.Duration(c.GeofenceDwellThresholdSec) * time.Second,
		MaxCandidateAreas:  c.GeofenceMaxCandidateAreas,
		MaxAccuracyForEval: c.GeofenceMaxAccuracyM,
	}
}

// RouteDeviationBusinessConfig returns the RouteDeviationConfig derived from this configuration.
func (c *GeolocationConfig) RouteDeviationBusinessConfig() business.RouteDeviationConfig {
	return business.RouteDeviationConfig{
		MaxAccuracyForEval: c.RouteDeviationMaxAccuracyM,
	}
}

// ProximityBusinessConfig returns the ProximityConfig derived from this configuration.
func (c *GeolocationConfig) ProximityBusinessConfig() business.ProximityConfig {
	return business.ProximityConfig{
		DefaultRadiusM: c.ProximityDefaultRadiusM,
		MaxRadiusM:     c.ProximityMaxRadiusM,
		StaleHours:     c.ProximityStaleHours,
		DefaultLimit:   c.ProximityDefaultLimit,
		MaxLimit:       c.ProximityMaxLimit,
	}
}

// TrackBusinessConfig returns the TrackConfig derived from this configuration.
func (c *GeolocationConfig) TrackBusinessConfig() business.TrackConfig {
	return business.TrackConfig{
		DefaultTrackLimit: c.TrackDefaultLimit,
		MaxTrackLimit:     c.TrackMaxLimit,
		DefaultEventLimit: c.EventDefaultLimit,
		MaxEventLimit:     c.EventMaxLimit,
		MaxAreaSubjects:   c.MaxAreaSubjects,
	}
}

// IngestionBusinessConfig returns the IngestionConfig derived from this configuration.
func (c *GeolocationConfig) IngestionBusinessConfig() business.IngestionConfig {
	return business.IngestionConfig{
		MaxBatchSize: c.IngestionMaxBatchSize,
	}
}

// RetentionBusinessConfig returns the RetentionConfig derived from this configuration.
func (c *GeolocationConfig) RetentionBusinessConfig() business.RetentionConfig {
	return business.RetentionConfig{
		LocationPointRetentionDays: c.RetentionLocationPointDays,
		GeoEventRetentionDays:      c.RetentionGeoEventDays,
		GeofenceStateStaleDays:     c.RetentionGeofenceStateDays,
		PartitionMaintenanceMonths: c.RetentionPartitionMonths,
		RetentionBatchSize:         c.RetentionBatchSize,
	}
}

// RateLimitConfig returns the rate limiter configuration.
func (c *GeolocationConfig) RateLimitConfig() *handlers.RateLimiterConfig {
	var proxies []string
	if c.RateLimitTrustedProxies != "" {
		proxies = strings.Split(c.RateLimitTrustedProxies, ",")
	}
	return &handlers.RateLimiterConfig{
		RequestsPerWindow: c.RateLimitRequestsPerMinute,
		WindowDuration:    time.Minute,
		TrustedProxies:    proxies,
	}
}
