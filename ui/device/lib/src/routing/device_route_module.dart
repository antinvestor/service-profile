import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/navigation/nav_items.dart';
import 'package:antinvestor_ui_core/permissions/permission_manifest.dart';
import 'package:antinvestor_ui_core/routing/route_module.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../screens/device_detail_screen.dart';
import '../screens/device_link_screen.dart';
import '../screens/device_list_screen.dart';

/// Route module for device management.
///
/// Registers the following routes:
/// - `/devices` - device list with search
/// - `/devices/detail/:id` - device detail with tabs (Info, Keys, Sessions)
/// - `/devices/link` - register or link a device
class DeviceRouteModule extends RouteModule {
  @override
  String get moduleId => 'devices';

  @override
  List<RouteBase> buildRoutes() {
    return [
      GoRoute(
        path: '/devices',
        builder: (context, state) => const DeviceListScreen(),
        routes: [
          GoRoute(
            path: 'detail/:id',
            builder: (context, state) {
              final id = Uri.decodeComponent(
                state.pathParameters['id'] ?? '',
              );
              final extra = state.extra;
              final device = extra is DeviceObject ? extra : null;
              return DeviceDetailScreen(
                deviceId: id,
                initialDevice: device,
              );
            },
          ),
          GoRoute(
            path: 'link',
            builder: (context, state) => const DeviceLinkScreen(),
          ),
        ],
      ),
    ];
  }

  @override
  List<NavItem> buildNavItems() {
    return [
      const NavItem(
        id: 'devices',
        label: 'Devices',
        icon: Icons.devices_outlined,
        activeIcon: Icons.devices,
        route: '/devices',
        requiredPermissions: {'device_view'},
        children: [
          NavItem(
            id: 'devices-list',
            label: 'All Devices',
            icon: Icons.list,
            route: '/devices',
            requiredPermissions: {'device_view'},
          ),
          NavItem(
            id: 'devices-link',
            label: 'Register Device',
            icon: Icons.add_circle_outline,
            route: '/devices/link',
            requiredPermissions: {'device_create'},
          ),
        ],
      ),
    ];
  }

  @override
  Map<String, Set<String>> get routePermissions => {
        '/devices': {'device_view'},
        '/devices/detail': {'device_view'},
        '/devices/link': {'device_create'},
      };

  @override
  PermissionManifest get permissionManifest => const PermissionManifest(
        namespace: 'service_device',
        permissions: [
          PermissionEntry(
            key: 'device_view',
            label: 'View Devices',
            scope: PermissionScope.service,
          ),
          PermissionEntry(
            key: 'device_create',
            label: 'Create Devices',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'device_update',
            label: 'Update Devices',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'device_remove',
            label: 'Remove Devices',
            scope: PermissionScope.action,
          ),
          PermissionEntry(
            key: 'device_key_manage',
            label: 'Manage Device Keys',
            scope: PermissionScope.feature,
          ),
          PermissionEntry(
            key: 'device_log_view',
            label: 'View Device Logs',
            scope: PermissionScope.feature,
          ),
        ],
      );
}
