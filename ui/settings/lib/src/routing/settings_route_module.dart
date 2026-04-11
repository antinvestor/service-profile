import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/navigation/nav_items.dart';
import 'package:antinvestor_ui_core/permissions/permission_manifest.dart';
import 'package:antinvestor_ui_core/routing/route_module.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../screens/setting_detail_screen.dart';
import '../screens/settings_bulk_edit_screen.dart';
import '../screens/settings_list_screen.dart';

/// Route module for settings management.
///
/// Registers the following routes:
/// - `/settings` - settings list (grouped by module, with search)
/// - `/settings/detail/:key` - view / edit a single setting
/// - `/settings/bulk-edit` - bulk editor for multiple settings
class SettingsRouteModule extends RouteModule {
  @override
  String get moduleId => 'settings';

  @override
  List<RouteBase> buildRoutes() {
    return [
      GoRoute(
        path: '/settings',
        builder: (context, state) => const SettingsListScreen(),
        routes: [
          GoRoute(
            path: 'detail/:key',
            builder: (context, state) {
              final key = Uri.decodeComponent(
                state.pathParameters['key'] ?? '',
              );
              final extra = state.extra;
              final setting =
                  extra is SettingObject ? extra : null;
              return SettingDetailScreen(
                settingKey: key,
                initialSetting: setting,
              );
            },
          ),
          GoRoute(
            path: 'bulk-edit',
            builder: (context, state) => const SettingsBulkEditScreen(),
          ),
        ],
      ),
    ];
  }

  @override
  List<NavItem> buildNavItems() {
    return [
      const NavItem(
        id: 'settings',
        label: 'Settings',
        icon: Icons.settings_outlined,
        activeIcon: Icons.settings,
        route: '/settings',
        requiredPermissions: {'setting_view'},
        children: [
          NavItem(
            id: 'settings-list',
            label: 'All Settings',
            icon: Icons.list,
            route: '/settings',
            requiredPermissions: {'setting_view'},
          ),
          NavItem(
            id: 'settings-bulk-edit',
            label: 'Bulk Edit',
            icon: Icons.edit_note,
            route: '/settings/bulk-edit',
            requiredPermissions: {'setting_update'},
          ),
        ],
      ),
    ];
  }

  @override
  Map<String, Set<String>> get routePermissions => {
        '/settings': {'setting_view'},
        '/settings/detail': {'setting_view'},
        '/settings/bulk-edit': {'setting_update'},
      };

  @override
  PermissionManifest get permissionManifest => const PermissionManifest(
        namespace: 'service_setting',
        permissions: [
          PermissionEntry(
            key: 'setting_view',
            label: 'View Settings',
            scope: PermissionScope.service,
          ),
          PermissionEntry(
            key: 'setting_update',
            label: 'Update Settings',
            scope: PermissionScope.action,
          ),
        ],
      );
}
