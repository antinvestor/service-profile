package caching

import (
	"context"

	"github.com/pitabwire/frame/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "devices/caching"

// Tracer provides OpenTelemetry tracing for the devices cache layer.
//
//nolint:gochecknoglobals // OpenTelemetry instruments must be global for reuse
var Tracer = telemetry.NewTracer(instrumentationName)

// Metric instruments for cache observability. Created through the
// BusinessMetrics factory so every measurement is transparently
// tenant-scoped (tenant_id/partition_id from the context's claims).
//
//nolint:gochecknoglobals // OpenTelemetry instruments must be global for reuse
var (
	businessMetrics  = telemetry.NewBusinessMetrics(instrumentationName)
	cacheHitCounter  = businessMetrics.Counter("devices/caching/cache_hits", "Number of cache hits")
	cacheMissCounter = businessMetrics.Counter("devices/caching/cache_misses", "Number of cache misses")
	rateLimitCounter = businessMetrics.Counter(
		"devices/caching/rate_limited",
		"Number of rate-limited requests",
	)
)

// RecordCacheHit records a cache hit for the given cache type.
func RecordCacheHit(ctx context.Context, cacheType string) {
	cacheHitCounter.Add(ctx, 1, attribute.String("cache_type", cacheType))
}

// RecordCacheMiss records a cache miss for the given cache type.
func RecordCacheMiss(ctx context.Context, cacheType string) {
	cacheMissCounter.Add(ctx, 1, attribute.String("cache_type", cacheType))
}

// RecordRateLimited records a rate-limited request for the given operation.
func RecordRateLimited(ctx context.Context, operation string) {
	rateLimitCounter.Add(ctx, 1, attribute.String("operation", operation))
}

// StartSpan creates a new trace span for a device service operation.
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// EndSpan completes a trace span, recording any error.
func EndSpan(ctx context.Context, span trace.Span, err error) {
	Tracer.End(ctx, span, err)
}
