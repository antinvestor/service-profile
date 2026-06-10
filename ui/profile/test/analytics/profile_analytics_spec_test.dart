import 'dart:convert';

import 'package:antinvestor_ui_core/antinvestor_ui_core.dart';
import 'package:antinvestor_ui_profile/antinvestor_ui_profile.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;

/// A transport call recorded with its decoded JSON body.
class RecordedRequest {
  RecordedRequest(this.path, this.body);
  final String path;
  final Map<String, dynamic> body;
}

/// Records every request and answers via a per-test handler.
class MockTransport {
  MockTransport(this.handler);

  final http.Response Function(String path, Map<String, dynamic> body) handler;
  final List<RecordedRequest> requests = [];

  Future<http.Response> call(String path, {Object? body}) async {
    final decoded = json.decode(body! as String) as Map<String, dynamic>;
    requests.add(RecordedRequest(path, decoded));
    return handler(path, decoded);
  }
}

http.Response ok(Object payload) => http.Response(
  json.encode(payload),
  200,
  headers: {'content-type': 'application/json'},
);

final fixedRange = AnalyticsTimeRange(
  start: DateTime.utc(2026, 6, 1),
  end: DateTime.utc(2026, 6, 8),
  granularity: TimeGranularity.day,
);

const fixedRangeJson = {
  'start': '2026-06-01T00:00:00.000Z',
  'end': '2026-06-08T00:00:00.000Z',
};

ThesaAnalyticsDataSource sourceWith(MockTransport transport) =>
    ThesaAnalyticsDataSource(
      transport.call,
      specs: const [profileAnalyticsSpec],
    );

void main() {
  group('device cache KPI contracts', () {
    test('getMetrics posts one exact scalar body per cache counter', () async {
      final transport = MockTransport((_, body) {
        return body['metric'] == deviceCacheHitsMetric
            ? ok({'value': 75})
            : ok({'value': 25});
      });
      final source = sourceWith(transport);

      final metrics = await source.getMetrics('profile', timeRange: fixedRange);

      expect(transport.requests, hasLength(2));
      expect(transport.requests.map((r) => r.path).toSet(), {
        '/api/analytics/query/scalar',
      });
      expect(transport.requests[0].body, {
        'metric': 'devices/caching/cache_hits',
        'aggregation': 'sum',
        'time_range': fixedRangeJson,
      });
      expect(transport.requests[1].body, {
        'metric': 'devices/caching/cache_misses',
        'aggregation': 'sum',
        'time_range': fixedRangeJson,
      });

      expect(metrics.map((m) => m.key), ['cache_hits', 'cache_misses']);
      expect(metrics[0].value, 75.0);
      expect(metrics[1].value, 25.0);
    });
  });

  group('geolocation ingestion contracts', () {
    const expectedLabels = {
      'service_geolocation/ingestion/accepted': 'Accepted Points',
      'service_geolocation/ingestion/rejected': 'Rejected Points',
    };

    for (final entry in expectedLabels.entries) {
      test('${entry.key} posts an exact sum time-series body', () async {
        final transport = MockTransport(
          (_, _) => ok({
            'points': [
              {'timestamp': '2026-06-01T00:00:00Z', 'value': 11},
            ],
          }),
        );
        final source = sourceWith(transport);

        final series = await source.getTimeSeries(
          'profile',
          entry.key,
          timeRange: fixedRange,
        );

        expect(
          transport.requests.single.path,
          '/api/analytics/query/timeseries',
        );
        expect(transport.requests.single.body, {
          'metric': entry.key,
          'aggregation': 'sum',
          'time_range': fixedRangeJson,
          'step': 'day',
        });
        expect(series.single.label, entry.value);
        expect(series.single.points.single.value, 11.0);
      });
    }
  });

  group('tenancy', () {
    test('never sends tenant_id or partition_id', () async {
      final transport = MockTransport((_, _) => ok({'value': 1}));
      final source = sourceWith(transport);

      await source.getMetrics('profile', timeRange: fixedRange);
      for (final metric in geoIngestionMetrics) {
        await source.getTimeSeries('profile', metric, timeRange: fixedRange);
      }

      for (final request in transport.requests) {
        final filters =
            request.body['filters'] as Map<String, dynamic>? ?? const {};
        expect(filters.keys, isNot(contains('tenant_id')));
        expect(filters.keys, isNot(contains('partition_id')));
      }
    });

    test('spec declares no tenancy filters anywhere', () {
      expect(profileAnalyticsSpec.baseFilters, isEmpty);
      for (final kpi in profileAnalyticsSpec.kpis) {
        expect(kpi.filters, isNull);
      }
      for (final chart in profileAnalyticsSpec.charts) {
        expect(chart.filters, isNull);
      }
    });
  });

  group('spec declaration', () {
    test('covers cache counters and ingestion outcomes', () {
      expect(profileAnalyticsSpec.service, 'profile');
      expect(profileAnalyticsSpec.metricKeys, ['cache_hits', 'cache_misses']);
      for (final metric in geoIngestionMetrics) {
        final chart = profileAnalyticsSpec.chartFor(metric);
        expect(chart, isNotNull, reason: metric);
        expect(chart!.type, ChartType.timeSeries, reason: metric);
        expect(chart.aggregation, AnalyticsAggregation.sum, reason: metric);
      }
    });
  });
}
