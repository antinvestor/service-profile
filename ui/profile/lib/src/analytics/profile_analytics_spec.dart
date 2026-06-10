import 'package:antinvestor_ui_core/antinvestor_ui_core.dart';
import 'package:flutter/material.dart';

/// Metric names emitted by the devices caching layer
/// (`apps/devices/service/caching`) and the geolocation ingestion pipeline
/// (`apps/geolocation/service/observability`).
const deviceCacheHitsMetric = 'devices/caching/cache_hits';
const deviceCacheMissesMetric = 'devices/caching/cache_misses';
const geoIngestionAcceptedMetric = 'service_geolocation/ingestion/accepted';
const geoIngestionRejectedMetric = 'service_geolocation/ingestion/rejected';

/// The geolocation ingestion outcomes charted as activity time series.
const geoIngestionMetrics = [
  geoIngestionAcceptedMetric,
  geoIngestionRejectedMetric,
];

/// Analytics catalog for the profile service family, served by the Thesa
/// analytics gate.
///
/// Profile entity counts stay on the profile API; this spec covers the
/// devices/geolocation activity signals. Tenant scoping is injected
/// server-side from the caller's JWT; no tenant filters are (or may be)
/// declared here.
const profileAnalyticsSpec = ServiceAnalyticsSpec(
  service: 'profile',
  kpis: [
    KpiSpec(
      'cache_hits',
      label: 'Device Cache Hits',
      metric: deviceCacheHitsMetric,
      unit: 'count',
      icon: Icons.memory_outlined,
    ),
    KpiSpec(
      'cache_misses',
      label: 'Device Cache Misses',
      metric: deviceCacheMissesMetric,
      unit: 'count',
      icon: Icons.sd_card_alert_outlined,
    ),
  ],
  charts: [
    ChartConfig.timeSeries(
      geoIngestionAcceptedMetric,
      label: 'Accepted Points',
    ),
    ChartConfig.timeSeries(
      geoIngestionRejectedMetric,
      label: 'Rejected Points',
    ),
  ],
);
