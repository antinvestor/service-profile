import 'dart:convert';

import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/antinvestor_ui_core.dart';
import 'package:antinvestor_ui_profile/antinvestor_ui_profile.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;

void main() {
  testWidgets('keeps entity KPI cards and contains no placeholder event data', (
    tester,
  ) async {
    Future<http.Response> transport(String path, {Object? body}) async {
      final payload = path.endsWith('/scalar')
          ? {'value': 0}
          : {'points': <Object>[]};
      return http.Response(
        json.encode(payload),
        200,
        headers: {'content-type': 'application/json'},
      );
    }

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          profileSearchProvider.overrideWith(
            (ref, query) async => <ProfileObject>[],
          ),
          analyticsDataSourceProvider.overrideWithValue(
            ThesaAnalyticsDataSource(
              transport,
              specs: const [profileAnalyticsSpec],
            ),
          ),
        ],
        child: const MaterialApp(
          home: Scaffold(body: ProfileAnalyticsScreen()),
        ),
      ),
    );
    await tester.pumpAndSettle();

    // Entity-API counts stay.
    expect(find.text('Total Profiles'), findsOneWidget);
    expect(find.text('Person Profiles'), findsOneWidget);
    expect(find.text('Verified Contacts'), findsOneWidget);

    // The gate-backed activity section replaces the mock event feed.
    expect(find.text('Devices & Geolocation Activity'), findsOneWidget);
    expect(find.text('Recent Events'), findsNothing);
    expect(find.text('New profile registered'), findsNothing);
    expect(find.text('Contact verification completed'), findsNothing);
    expect(find.text('Profile merge executed'), findsNothing);
  });
}
