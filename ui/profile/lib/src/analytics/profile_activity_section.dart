import 'package:antinvestor_ui_core/antinvestor_ui_core.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'analytics_states.dart';
import 'profile_analytics_spec.dart';

/// Devices & geolocation activity from the Thesa analytics gate.
///
/// Complements the entity counts on the profile analytics screen with
/// tenant-scoped activity signals: device cache hits vs misses (with a
/// client-computed hit ratio) and geolocation ingestion accepted vs
/// rejected time series.
class ProfileActivitySection extends ConsumerStatefulWidget {
  const ProfileActivitySection({super.key});

  @override
  ConsumerState<ProfileActivitySection> createState() =>
      _ProfileActivitySectionState();
}

class _ProfileActivitySectionState
    extends ConsumerState<ProfileActivitySection> {
  AnalyticsTimeRange _range = AnalyticsTimeRange.last30Days();

  void _refresh() {
    ref.invalidate(serviceMetricsProvider);
    ref.invalidate(serviceTimeSeriesProvider);
  }

  /// Appends the client-computed cache hit ratio to the gate-provided
  /// hit/miss counters. The gate is queried per metric; the ratio is
  /// derived locally so no extra endpoint is needed.
  List<MetricValue> _withHitRatio(List<MetricValue> metrics) {
    double? hits;
    double? misses;
    for (final metric in metrics) {
      if (metric.key == 'cache_hits') hits = metric.value;
      if (metric.key == 'cache_misses') misses = metric.value;
    }
    final total = (hits ?? 0) + (misses ?? 0);
    if (hits == null || misses == null || total <= 0) return metrics;
    return [
      ...metrics,
      MetricValue(
        key: 'cache_hit_ratio',
        label: 'Cache Hit Ratio',
        value: hits / total * 100,
        unit: 'percent',
        icon: Icons.speed_outlined,
      ),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final metricsAsync = ref.watch(
      serviceMetricsProvider(
        ServiceMetricsParams(profileAnalyticsSpec.service, timeRange: _range),
      ),
    );
    final ingestionAsync = [
      for (final metric in geoIngestionMetrics)
        ref.watch(
          serviceTimeSeriesProvider(
            ServiceTimeSeriesParams(
              profileAnalyticsSpec.service,
              metric,
              timeRange: _range,
            ),
          ),
        ),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Expanded(
              child: Text(
                'Devices & Geolocation Activity',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            IconButton(
              icon: const Icon(Icons.refresh),
              tooltip: 'Refresh activity',
              onPressed: _refresh,
            ),
          ],
        ),
        const SizedBox(height: 8),
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: TimeRangeSelector(
            value: _range,
            onChanged: (range) => setState(() => _range = range),
          ),
        ),
        const SizedBox(height: 16),
        metricsAsync.when(
          data: (metrics) => MetricsRow(metrics: _withHitRatio(metrics)),
          loading: () =>
              const MetricsRow(metrics: [], isLoading: true, skeletonCount: 3),
          error: (error, _) =>
              AnalyticsErrorCard(error: error, onRetry: _refresh),
        ),
        const SizedBox(height: 16),
        _ingestionChart(theme, ingestionAsync),
      ],
    );
  }

  /// Geolocation ingestion accepted vs rejected as one combined chart.
  Widget _ingestionChart(
    ThemeData theme,
    List<AsyncValue<List<TimeSeries>>> asyncSeries,
  ) {
    final cs = theme.colorScheme;

    Widget body;
    final failed = asyncSeries.where((a) => a.hasError).toList();
    if (failed.isNotEmpty) {
      body = AnalyticsErrorCard(error: failed.first.error!, onRetry: _refresh);
    } else if (asyncSeries.any((a) => a.isLoading)) {
      body = const SizedBox(
        height: 240,
        child: Center(child: CircularProgressIndicator()),
      );
    } else {
      body = TimeSeriesChart(
        series: [for (final a in asyncSeries) ...a.requireValue],
        granularity: _range.granularity,
      );
    }

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: cs.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: cs.outlineVariant),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Geolocation Ingestion (accepted vs rejected)',
            style: theme.textTheme.titleSmall?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          body,
        ],
      ),
    );
  }
}
