package observability

import (
	"context"
	"time"

	"github.com/pitabwire/frame/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const pkgName = "service_geolocation"

// Metrics holds pre-allocated OTel instruments for the geolocation service.
// Instruments are created once at startup and reused for every measurement.
type Metrics struct {
	tracer telemetry.Tracer

	// Ingestion metrics.
	ingestBatchLatency metric.Float64Histogram
	ingestAccepted     metric.Int64Counter
	ingestRejected     metric.Int64Counter

	// Geofence metrics.
	geofenceEvalLatency metric.Float64Histogram
	geofenceTransitions metric.Int64Counter

	// Route deviation metrics.
	routeDeviationEvalLatency metric.Float64Histogram
	routeDeviationTransitions metric.Int64Counter

	// Proximity metrics.
	proximityQueryLatency metric.Float64Histogram
}

// NewMetrics creates and registers all OTel instruments for the geolocation service.
func NewMetrics() *Metrics {
	t := telemetry.NewTracer(pkgName)

	return &Metrics{
		tracer:             t,
		ingestBatchLatency: telemetry.LatencyMeasure(pkgName + "/ingestion"),
		ingestAccepted: telemetry.DimensionlessMeasure(
			pkgName,
			"/ingestion/accepted",
			"Number of accepted location points",
		),
		ingestRejected: telemetry.DimensionlessMeasure(
			pkgName,
			"/ingestion/rejected",
			"Number of rejected location points",
		),
		geofenceEvalLatency: telemetry.LatencyMeasure(pkgName + "/geofence"),
		geofenceTransitions: telemetry.DimensionlessMeasure(
			pkgName,
			"/geofence/transitions",
			"Number of geofence state transitions",
		),
		routeDeviationEvalLatency: telemetry.LatencyMeasure(pkgName + "/route_deviation"),
		routeDeviationTransitions: telemetry.DimensionlessMeasure(
			pkgName,
			"/route_deviation/transitions",
			"Number of route deviation state transitions",
		),
		proximityQueryLatency: telemetry.LatencyMeasure(pkgName + "/proximity"),
	}
}

// StartSpan starts a new traced span and returns the enriched context and span.
func (m *Metrics) StartSpan(
	ctx context.Context,
	name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return m.tracer.Start(ctx, name, opts...)
}

// EndSpan ends a span and records latency.
func (m *Metrics) EndSpan(ctx context.Context, span trace.Span, err error) {
	m.tracer.End(ctx, span, err)
}

// RecordIngestBatch records metrics for a batch ingestion.
func (m *Metrics) RecordIngestBatch(
	ctx context.Context,
	accepted, rejected int32,
	elapsed time.Duration,
) {
	m.ingestBatchLatency.Record(ctx, float64(elapsed.Milliseconds()))
	m.ingestAccepted.Add(ctx, int64(accepted))
	m.ingestRejected.Add(ctx, int64(rejected))
}

// RecordGeofenceEval records metrics for a geofence evaluation.
func (m *Metrics) RecordGeofenceEval(ctx context.Context, elapsed time.Duration) {
	m.geofenceEvalLatency.Record(ctx, float64(elapsed.Milliseconds()))
}

// RecordGeofenceTransition records a geofence state transition (enter/exit/dwell).
func (m *Metrics) RecordGeofenceTransition(ctx context.Context, eventType string) {
	m.geofenceTransitions.Add(ctx, 1,
		metric.WithAttributes(attribute.String("event_type", eventType)),
	)
}

// RecordRouteDeviationEval records metrics for a route deviation evaluation.
func (m *Metrics) RecordRouteDeviationEval(ctx context.Context, elapsed time.Duration) {
	m.routeDeviationEvalLatency.Record(ctx, float64(elapsed.Milliseconds()))
}

// RecordRouteDeviationTransition records a route deviation state transition.
func (m *Metrics) RecordRouteDeviationTransition(ctx context.Context, eventType string) {
	m.routeDeviationTransitions.Add(ctx, 1,
		metric.WithAttributes(attribute.String("event_type", eventType)),
	)
}

// RecordProximityQuery records metrics for a proximity query.
func (m *Metrics) RecordProximityQuery(
	ctx context.Context,
	elapsed time.Duration,
	resultCount int,
) {
	m.proximityQueryLatency.Record(ctx, float64(elapsed.Milliseconds()),
		metric.WithAttributes(attribute.Int("result_count", resultCount)),
	)
}
