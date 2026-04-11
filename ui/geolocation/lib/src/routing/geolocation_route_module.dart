import 'package:antinvestor_ui_core/navigation/nav_items.dart';
import 'package:antinvestor_ui_core/routing/route_module.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../screens/area_list_screen.dart';
import '../screens/area_detail_screen.dart';
import '../screens/area_create_screen.dart';
import '../screens/area_edit_screen.dart';
import '../screens/route_list_screen.dart';
import '../screens/route_detail_screen.dart';
import '../screens/route_create_screen.dart';
import '../screens/location_track_screen.dart';
import '../screens/geo_events_screen.dart';

class GeolocationRouteModule extends RouteModule {
  @override
  String get moduleId => 'geolocation';

  @override
  List<RouteBase> buildRoutes() => [
        GoRoute(
          path: '/geo/areas',
          builder: (context, state) => const AreaListScreen(),
          routes: [
            GoRoute(
              path: 'new',
              builder: (context, state) => const AreaCreateScreen(),
            ),
            GoRoute(
              path: ':areaId',
              builder: (context, state) => AreaDetailScreen(
                areaId: state.pathParameters['areaId']!,
              ),
              routes: [
                GoRoute(
                  path: 'edit',
                  builder: (context, state) => AreaEditScreen(
                    areaId: state.pathParameters['areaId']!,
                  ),
                ),
              ],
            ),
          ],
        ),
        GoRoute(
          path: '/geo/routes',
          builder: (context, state) => const RouteListScreen(),
          routes: [
            GoRoute(
              path: 'new',
              builder: (context, state) => const RouteCreateScreen(),
            ),
            GoRoute(
              path: ':routeId',
              builder: (context, state) => RouteDetailScreen(
                routeId: state.pathParameters['routeId']!,
              ),
            ),
          ],
        ),
        GoRoute(
          path: '/geo/track/:subjectId',
          builder: (context, state) => LocationTrackScreen(
            subjectId: state.pathParameters['subjectId']!,
          ),
        ),
        GoRoute(
          path: '/geo/events',
          builder: (context, state) => const GeoEventsScreen(),
        ),
      ];

  @override
  List<NavItem> buildNavItems() => [
        NavItem(
          id: 'geolocation',
          label: 'Geolocation',
          icon: Icons.map_outlined,
          activeIcon: Icons.map,
          route: '/geo/areas',
          children: [
            const NavItem(
              id: 'geo-areas',
              label: 'Areas',
              icon: Icons.layers_outlined,
              activeIcon: Icons.layers,
              route: '/geo/areas',
            ),
            const NavItem(
              id: 'geo-routes',
              label: 'Routes',
              icon: Icons.route_outlined,
              activeIcon: Icons.route,
              route: '/geo/routes',
            ),
            const NavItem(
              id: 'geo-events',
              label: 'Events',
              icon: Icons.event_note_outlined,
              activeIcon: Icons.event_note,
              route: '/geo/events',
            ),
          ],
        ),
      ];

  @override
  Map<String, Set<String>> get routePermissions => {
        '/geo/areas': {},
        '/geo/routes': {},
        '/geo/track': {},
        '/geo/events': {},
      };
}
