package observability //nolint:testpackage // tests access unexported metrics internals

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
