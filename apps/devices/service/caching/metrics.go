package caching

import (
	"context"

	"github.com/pitabwire/frame/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "devices/caching"

// Tracer provides OpenTelemetry tracing for the devices cache layer.
//
//nolint:gochecknoglobals // OpenTelemetry instruments must be global for reuse
var Tracer = telemetry.NewTracer(instrumentationName)

// Metric instruments for cache observability.
//
//nolint:gochecknoglobals // OpenTelemetry instruments must be global for reuse
var (
	cacheHitCounter  = telemetry.DimensionlessMeasure(instrumentationName, "cache_hits", "Number of cache hits")
	cacheMissCounter = telemetry.DimensionlessMeasure(instrumentationName, "cache_misses", "Number of cache misses")
	rateLimitCounter = telemetry.DimensionlessMeasure(
		instrumentationName,
		"rate_limited",
		"Number of rate-limited requests",
	)
)

// RecordCacheHit records a cache hit for the given cache type.
func RecordCacheHit(ctx context.Context, cacheType string) {
	cacheHitCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache_type", cacheType),
	))
}

// RecordCacheMiss records a cache miss for the given cache type.
func RecordCacheMiss(ctx context.Context, cacheType string) {
	cacheMissCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache_type", cacheType),
	))
}

// RecordRateLimited records a rate-limited request for the given operation.
func RecordRateLimited(ctx context.Context, operation string) {
	rateLimitCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", operation),
	))
}

// StartSpan creates a new trace span for a device service operation.
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// EndSpan completes a trace span, recording any error.
func EndSpan(ctx context.Context, span trace.Span, err error) {
	Tracer.End(ctx, span, err)
}
