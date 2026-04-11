import 'package:antinvestor_ui_core/navigation/nav_items.dart';
import 'package:antinvestor_ui_core/routing/route_module.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../screens/profile_search_screen.dart';
import '../screens/profile_detail_screen.dart';
import '../screens/profile_edit_screen.dart';
import '../screens/profile_create_screen.dart';
import '../screens/contacts_screen.dart';
import '../screens/addresses_screen.dart';
import '../screens/relationships_screen.dart';
import '../screens/roster_screen.dart';
import '../screens/profile_merge_screen.dart';
import '../screens/profile_analytics_screen.dart';

class ProfileRouteModule extends RouteModule {
  @override
  String get moduleId => 'profile';

  @override
  List<RouteBase> buildRoutes() => [
        GoRoute(
          path: '/profiles',
          builder: (context, state) => const ProfileSearchScreen(),
          routes: [
            GoRoute(
              path: 'analytics',
              builder: (context, state) =>
                  const ProfileAnalyticsScreen(),
            ),
            GoRoute(
              path: 'new',
              builder: (context, state) => const ProfileCreateScreen(),
            ),
            GoRoute(
              path: 'merge',
              builder: (context, state) => const ProfileMergeScreen(),
            ),
            GoRoute(
              path: ':profileId',
              builder: (context, state) => ProfileDetailScreen(
                profileId: state.pathParameters['profileId']!,
              ),
              routes: [
                GoRoute(
                  path: 'edit',
                  builder: (context, state) => ProfileEditScreen(
                    profileId: state.pathParameters['profileId']!,
                  ),
                ),
                GoRoute(
                  path: 'contacts',
                  builder: (context, state) => ContactsScreen(
                    profileId: state.pathParameters['profileId']!,
                  ),
                ),
                GoRoute(
                  path: 'addresses',
                  builder: (context, state) => AddressesScreen(
                    profileId: state.pathParameters['profileId']!,
                  ),
                ),
                GoRoute(
                  path: 'relationships',
                  builder: (context, state) => RelationshipsScreen(
                    profileId: state.pathParameters['profileId']!,
                  ),
                ),
                GoRoute(
                  path: 'roster',
                  builder: (context, state) => RosterScreen(
                    profileId: state.pathParameters['profileId']!,
                  ),
                ),
              ],
            ),
          ],
        ),
      ];

  @override
  List<NavItem> buildNavItems() => [
        const NavItem(
          id: 'profiles',
          label: 'Profiles',
          icon: Icons.people_outline,
          activeIcon: Icons.people,
          route: '/profiles',
        ),
      ];

  @override
  Map<String, Set<String>> get routePermissions => {
        '/profiles': {},
      };
}
