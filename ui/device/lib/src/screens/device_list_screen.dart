import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/device_providers.dart';
import '../widgets/device_card.dart';

/// Screen displaying a searchable list of registered devices.
class DeviceListScreen extends ConsumerStatefulWidget {
  const DeviceListScreen({super.key});

  @override
  ConsumerState<DeviceListScreen> createState() => _DeviceListScreenState();
}

class _DeviceListScreenState extends ConsumerState<DeviceListScreen> {
  String _searchQuery = '';

  @override
  Widget build(BuildContext context) {
    final asyncDevices = ref.watch(
      deviceSearchProvider(_searchQuery),
    );

    return asyncDevices.when(
      loading: () => _buildShell(isLoading: true, items: const []),
      error: (error, _) => _buildShell(
        error: friendlyError(error),
        items: const [],
      ),
      data: (devices) => _buildShell(items: devices),
    );
  }

  Widget _buildShell({
    required List<DeviceObject> items,
    bool isLoading = false,
    String? error,
  }) {
    return EntityListPage<DeviceObject>(
      title: 'Devices',
      icon: Icons.devices,
      items: items,
      isLoading: isLoading,
      error: error,
      onRetry: () => _refresh(),
      searchHint: 'Search devices by name, OS, or IP...',
      onSearchChanged: (query) {
        setState(() => _searchQuery = query.trim());
      },
      actionLabel: 'Register Device',
      onAction: () => context.go('/devices/link'),
      itemBuilder: (context, device) {
        return DeviceCard(
          device: device,
          onTap: () {
            context.go(
              '/devices/detail/${Uri.encodeComponent(device.id)}',
              extra: device,
            );
          },
        );
      },
    );
  }

  void _refresh() {
    ref.invalidate(
      deviceSearchProvider(_searchQuery),
    );
  }
}
