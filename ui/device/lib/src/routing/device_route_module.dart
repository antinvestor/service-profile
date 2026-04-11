import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/navigation/nav_items.dart';
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
        children: [
          NavItem(
            id: 'devices-list',
            label: 'All Devices',
            icon: Icons.list,
            route: '/devices',
          ),
          NavItem(
            id: 'devices-link',
            label: 'Register Device',
            icon: Icons.add_circle_outline,
            route: '/devices/link',
          ),
        ],
      ),
    ];
  }

  @override
  Map<String, Set<String>> get routePermissions => {
        '/devices': {'device:read', 'admin'},
        '/devices/detail': {'device:read', 'admin'},
        '/devices/link': {'device:write', 'admin'},
      };
}
