import 'dart:convert';

import 'package:antinvestor_ui_core/antinvestor_ui_core.dart';
import 'package:antinvestor_ui_profile/antinvestor_ui_profile.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;

http.Response ok(Object payload) => http.Response(
  json.encode(payload),
  200,
  headers: {'content-type': 'application/json'},
);

http.Response apiError(int status, String message) => http.Response(
  json.encode({'error': message}),
  status,
  headers: {'content-type': 'application/json'},
);

/// Wires the section into a ProviderScope with a stubbed gate transport.
Widget harness(
  http.Response Function(String path, Map<String, dynamic> body) handler,
) {
  Future<http.Response> transport(String path, {Object? body}) async {
    final decoded = json.decode(body! as String) as Map<String, dynamic>;
    return handler(path, decoded);
  }

  return ProviderScope(
    overrides: [
      analyticsDataSourceProvider.overrideWithValue(
        ThesaAnalyticsDataSource(
          transport,
          specs: const [profileAnalyticsSpec],
        ),
      ),
    ],
    child: const MaterialApp(
      home: Scaffold(
        body: SingleChildScrollView(child: ProfileActivitySection()),
      ),
    ),
  );
}

void main() {
  testWidgets(
    'renders cache KPIs with computed hit ratio and ingestion chart',
    (tester) async {
      await tester.pumpWidget(
        harness((path, body) {
          if (path.endsWith('/scalar')) {
            return body['metric'] == deviceCacheHitsMetric
                ? ok({'value': 75})
                : ok({'value': 25});
          }
          return ok({
            'points': [
              {'timestamp': '2026-06-01T00:00:00Z', 'value': 5},
              {'timestamp': '2026-06-02T00:00:00Z', 'value': 8},
            ],
          });
        }),
      );
      await tester.pumpAndSettle();

      expect(find.text('Devices & Geolocation Activity'), findsOneWidget);
      expect(find.text('Device Cache Hits'), findsOneWidget);
      expect(find.text('75'), findsOneWidget);
      expect(find.text('Device Cache Misses'), findsOneWidget);
      expect(find.text('25'), findsOneWidget);
      // Ratio is computed client-side: 75 / (75 + 25) = 75%.
      expect(find.text('Cache Hit Ratio'), findsOneWidget);
      expect(find.text('75.0%'), findsOneWidget);
      expect(
        find.text('Geolocation Ingestion (accepted vs rejected)'),
        findsOneWidget,
      );
      expect(find.text('Accepted Points'), findsOneWidget);
      expect(find.text('Rejected Points'), findsOneWidget);
      expect(find.byType(TimeSeriesChart), findsOneWidget);
    },
  );

  testWidgets('omits the hit ratio card when there is no cache traffic', (
    tester,
  ) async {
    await tester.pumpWidget(
      harness((path, _) {
        if (path.endsWith('/scalar')) return ok({'value': 0});
        return ok({'points': []});
      }),
    );
    await tester.pumpAndSettle();

    expect(find.text('Device Cache Hits'), findsOneWidget);
    expect(find.text('Cache Hit Ratio'), findsNothing);
    expect(find.text('No data'), findsOneWidget);
  });

  testWidgets('shows access message on 403 instead of raw error', (
    tester,
  ) async {
    await tester.pumpWidget(
      harness((_, _) => apiError(403, 'tenant scope required')),
    );
    await tester.pumpAndSettle();

    expect(
      find.textContaining('You do not have access to analytics'),
      findsWidgets,
    );
    expect(find.textContaining('AnalyticsQueryException'), findsNothing);
    expect(find.textContaining('tenant scope required'), findsNothing);
  });

  testWidgets('shows temporary outage message on 503', (tester) async {
    await tester.pumpWidget(
      harness((_, _) => apiError(503, 'backend unavailable')),
    );
    await tester.pumpAndSettle();

    expect(
      find.textContaining('Analytics is temporarily unavailable'),
      findsWidgets,
    );
    expect(find.text('Retry'), findsWidgets);
  });

  testWidgets('shows unsupported-metric message on 400', (tester) async {
    await tester.pumpWidget(
      harness((_, _) => apiError(400, 'metric not allowed')),
    );
    await tester.pumpAndSettle();

    expect(
      find.textContaining('not supported by the analytics service'),
      findsWidgets,
    );
    expect(find.textContaining('metric not allowed'), findsNothing);
  });
}
