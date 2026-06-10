package observability //nolint:testpackage // tests access unexported metrics internals

import (
	"context"
	"testing"
	"time"

	"github.com/pitabwire/frame/security"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMetricsRecorders(t *testing.T) {
	t.Parallel()

	metrics := NewMetrics()
	ctx, span := metrics.StartSpan(context.Background(), "test-span")
	require.NotNil(t, span)
	metrics.EndSpan(ctx, span, nil)

	metrics.RecordIngestBatch(ctx, 1, 2, time.Second)
	metrics.RecordGeofenceEval(ctx, time.Second)
	metrics.RecordGeofenceTransition(ctx, "enter")
	metrics.RecordRouteDeviationEval(ctx, time.Second)
	metrics.RecordRouteDeviationTransition(ctx, "deviated")
	metrics.RecordProximityQuery(ctx, time.Second, 3)
}

// findMetric returns the metric with the given name, or nil when absent.
func findMetric(rm metricdata.ResourceMetrics, name string) *metricdata.Metrics {
	for _, sm := range rm.ScopeMetrics {
		for i := range sm.Metrics {
			if sm.Metrics[i].Name == name {
				return &sm.Metrics[i]
			}
		}
	}
	return nil
}

// TestMetricsTenantAttribution proves, via a ManualReader, that instruments
// created through the BusinessMetrics factory keep their exact metric names
// and automatically attach tenant_id/partition_id from the context claims
// alongside the explicit non-tenant attributes.
func TestMetricsTenantAttribution(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	otel.SetMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)))

	metrics := NewMetrics()

	claims := &security.AuthenticationClaims{
		TenantID:    "tenant-geo-attr",
		PartitionID: "partition-geo-attr",
	}
	claims.Subject = "subject-geo-attr"
	ctx := claims.ClaimsToContext(context.Background())

	metrics.RecordGeofenceTransition(ctx, "enter")
	metrics.RecordIngestBatch(ctx, 3, 1, 250*time.Millisecond)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	// Counter: exact name preserved, tenant + event_type attributes present.
	transitions := findMetric(rm, "service_geolocation/geofence/transitions")
	require.NotNil(t, transitions, "geofence transitions counter must keep its metric name")
	sum, ok := transitions.Data.(metricdata.Sum[int64])
	require.True(t, ok, "transitions must be an int64 sum")

	var matched bool
	for _, dp := range sum.DataPoints {
		tenant, hasTenant := dp.Attributes.Value("tenant_id")
		if !hasTenant || tenant.AsString() != "tenant-geo-attr" {
			continue
		}
		matched = true
		partition, hasPartition := dp.Attributes.Value("partition_id")
		require.True(t, hasPartition, "partition_id must accompany tenant_id")
		require.Equal(t, "partition-geo-attr", partition.AsString())
		eventType, hasEventType := dp.Attributes.Value("event_type")
		require.True(t, hasEventType, "non-tenant attributes must be preserved")
		require.Equal(t, "enter", eventType.AsString())
		require.Equal(t, int64(1), dp.Value)
	}
	require.True(t, matched, "expected a datapoint attributed to tenant-geo-attr")

	// Histogram: exact latency metric name preserved with tenant attribution.
	latency := findMetric(rm, "service_geolocation/ingestion/latency")
	require.NotNil(t, latency, "ingestion latency histogram must keep its metric name")
	hist, ok := latency.Data.(metricdata.Histogram[float64])
	require.True(t, ok, "latency must be a float64 histogram")

	var histMatched bool
	for _, dp := range hist.DataPoints {
		tenant, hasTenant := dp.Attributes.Value("tenant_id")
		if !hasTenant || tenant.AsString() != "tenant-geo-attr" {
			continue
		}
		histMatched = true
		require.Equal(t, uint64(1), dp.Count)
		require.InEpsilon(t, 250.0, dp.Sum, 0.001)
	}
	require.True(t, histMatched, "expected a latency datapoint attributed to tenant-geo-attr")

	// Companion counters from the same batch keep their names too.
	require.NotNil(t, findMetric(rm, "service_geolocation/ingestion/accepted"))
	require.NotNil(t, findMetric(rm, "service_geolocation/ingestion/rejected"))
}
