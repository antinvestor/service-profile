import 'package:antinvestor_ui_core/navigation/nav_items.dart';
import 'package:antinvestor_ui_core/permissions/permission_manifest.dart';
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
          requiredPermissions: {'profile_view'},
        ),
      ];

  @override
  Map<String, Set<String>> get routePermissions => {
        '/profiles': {'profile_view'},
        '/profiles/new': {'profile_create'},
        '/profiles/merge': {'profile_merge'},
      };

  @override
  PermissionManifest get permissionManifest => const PermissionManifest(
        namespace: 'service_profile',
        permissions: [
          PermissionEntry(
            key: 'profile_view',
            label: 'View Profiles',
            scope: PermissionScope.service,
          ),
          PermissionEntry(
            key: 'profile_create',
            label: 'Create Profiles',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'profile_update',
            label: 'Update Profiles',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'profile_merge',
            label: 'Merge Profiles',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'contact_manage',
            label: 'Manage Contacts',
            scope: PermissionScope.feature,
          ),
          PermissionEntry(
            key: 'roster_view',
            label: 'View Roster',
            scope: PermissionScope.feature,
          ),
          PermissionEntry(
            key: 'roster_manage',
            label: 'Manage Roster',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'address_manage',
            label: 'Manage Addresses',
            scope: PermissionScope.feature,
          ),
          PermissionEntry(
            key: 'relationship_view',
            label: 'View Relationships',
            scope: PermissionScope.feature,
          ),
          PermissionEntry(
            key: 'relationship_manage',
            label: 'Manage Relationships',
            scope: PermissionScope.action,
          ),
        ],
      );
}
