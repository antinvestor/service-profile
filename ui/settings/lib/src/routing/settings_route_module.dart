import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/navigation/nav_items.dart';
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
        children: [
          NavItem(
            id: 'settings-list',
            label: 'All Settings',
            icon: Icons.list,
            route: '/settings',
          ),
          NavItem(
            id: 'settings-bulk-edit',
            label: 'Bulk Edit',
            icon: Icons.edit_note,
            route: '/settings/bulk-edit',
          ),
        ],
      ),
    ];
  }

  @override
  Map<String, Set<String>> get routePermissions => {
        '/settings': {'settings:read', 'admin'},
        '/settings/detail': {'settings:read', 'admin'},
        '/settings/bulk-edit': {'settings:write', 'admin'},
      };
}
