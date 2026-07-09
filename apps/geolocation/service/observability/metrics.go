package observability

import (
	"context"
	"time"

	"github.com/pitabwire/frame/v2/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const pkgName = "service_geolocation"

// Metrics holds pre-allocated OTel instruments for the geolocation service.
// Instruments are created once at startup through the BusinessMetrics
// factory, so every measurement is transparently tenant-scoped
// (tenant_id/partition_id derived from the context's security claims).
type Metrics struct {
	tracer telemetry.Tracer

	// Ingestion metrics.
	ingestBatchLatency telemetry.Histogram
	ingestAccepted     telemetry.Counter
	ingestRejected     telemetry.Counter

	// Geofence metrics.
	geofenceEvalLatency telemetry.Histogram
	geofenceTransitions telemetry.Counter

	// Route deviation metrics.
	routeDeviationEvalLatency telemetry.Histogram
	routeDeviationTransitions telemetry.Counter

	// Proximity metrics.
	proximityQueryLatency telemetry.Histogram
}

// NewMetrics creates and registers all OTel instruments for the geolocation service.
func NewMetrics() *Metrics {
	t := telemetry.NewTracer(pkgName)
	bm := telemetry.NewBusinessMetrics(pkgName)

	return &Metrics{
		tracer: t,
		ingestBatchLatency: bm.Histogram(
			pkgName+"/ingestion/latency",
			"Latency distribution of batch ingestions",
		),
		ingestAccepted: bm.Counter(
			pkgName+"/ingestion/accepted",
			"Number of accepted location points",
		),
		ingestRejected: bm.Counter(
			pkgName+"/ingestion/rejected",
			"Number of rejected location points",
		),
		geofenceEvalLatency: bm.Histogram(
			pkgName+"/geofence/latency",
			"Latency distribution of geofence evaluations",
		),
		geofenceTransitions: bm.Counter(
			pkgName+"/geofence/transitions",
			"Number of geofence state transitions",
		),
		routeDeviationEvalLatency: bm.Histogram(
			pkgName+"/route_deviation/latency",
			"Latency distribution of route deviation evaluations",
		),
		routeDeviationTransitions: bm.Counter(
			pkgName+"/route_deviation/transitions",
			"Number of route deviation state transitions",
		),
		proximityQueryLatency: bm.Histogram(
			pkgName+"/proximity/latency",
			"Latency distribution of proximity queries",
		),
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
	m.geofenceTransitions.Add(ctx, 1, attribute.String("event_type", eventType))
}

// RecordRouteDeviationEval records metrics for a route deviation evaluation.
func (m *Metrics) RecordRouteDeviationEval(ctx context.Context, elapsed time.Duration) {
	m.routeDeviationEvalLatency.Record(ctx, float64(elapsed.Milliseconds()))
}

// RecordRouteDeviationTransition records a route deviation state transition.
func (m *Metrics) RecordRouteDeviationTransition(ctx context.Context, eventType string) {
	m.routeDeviationTransitions.Add(ctx, 1, attribute.String("event_type", eventType))
}

// RecordProximityQuery records metrics for a proximity query.
func (m *Metrics) RecordProximityQuery(
	ctx context.Context,
	elapsed time.Duration,
	resultCount int,
) {
	m.proximityQueryLatency.Record(ctx, float64(elapsed.Milliseconds()),
		attribute.Int("result_count", resultCount),
	)
}
