package config //nolint:testpackage // tests access unexported config fields

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDerivedConfigs(t *testing.T) {
	t.Parallel()

	cfg := GeolocationConfig{
		GeofenceHysteresisBufferM:  15,
		GeofenceDwellThresholdSec:  30,
		GeofenceMaxCandidateAreas:  5,
		GeofenceMaxAccuracyM:       20,
		RouteDeviationMaxAccuracyM: 25,
		ProximityDefaultRadiusM:    100,
		ProximityMaxRadiusM:        200,
		ProximityStaleHours:        2,
		ProximityDefaultLimit:      3,
		ProximityMaxLimit:          4,
		TrackDefaultLimit:          5,
		TrackMaxLimit:              6,
		EventDefaultLimit:          7,
		EventMaxLimit:              8,
		MaxAreaSubjects:            9,
		IngestionMaxBatchSize:      10,
		RetentionLocationPointDays: 11,
		RetentionGeoEventDays:      12,
		RetentionGeofenceStateDays: 13,
		RetentionPartitionMonths:   14,
		RetentionBatchSize:         15,
		RateLimitRequestsPerMinute: 16,
		RateLimitTrustedProxies:    "1.1.1.1,2.2.2.2",
	}

	require.InDelta(t, 15.0, cfg.GeofenceBusinessConfig().HysteresisBufferM, 0.001)
	require.Equal(t, 30*time.Second, cfg.GeofenceBusinessConfig().DwellThreshold)
	require.InDelta(t, 25.0, cfg.RouteDeviationBusinessConfig().MaxAccuracyForEval, 0.001)
	require.InDelta(t, 100.0, cfg.ProximityBusinessConfig().DefaultRadiusM, 0.001)
	require.Equal(t, 5, cfg.TrackBusinessConfig().DefaultTrackLimit)
	require.Equal(t, 10, cfg.IngestionBusinessConfig().MaxBatchSize)
	require.Equal(t, 11, cfg.RetentionBusinessConfig().LocationPointRetentionDays)
	require.Len(t, cfg.RateLimitConfig().TrustedProxies, 2)
}
