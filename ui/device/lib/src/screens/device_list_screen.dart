import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/device_providers.dart';
import '../widgets/device_card.dart';
import '../widgets/presence_indicator.dart';

/// Screen displaying a searchable, paginated DataTable of registered devices.
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
      loading: () => _buildPage(items: const [], isLoading: true),
      error: (error, _) => _buildPage(items: const [], error: friendlyError(error)),
      data: (devices) => _buildPage(items: devices),
    );
  }

  Widget _buildPage({
    required List<DeviceObject> items,
    bool isLoading = false,
    String? error,
  }) {
    if (isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(error),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: _refresh,
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    return AdminEntityListPage<DeviceObject>(
      title: 'Devices',
      breadcrumbs: const ['Home', 'Devices'],
      columns: const [
        DataColumn(label: Text('Name')),
        DataColumn(label: Text('OS')),
        DataColumn(label: Text('IP')),
        DataColumn(label: Text('Last Seen')),
        DataColumn(label: Text('Profile')),
        DataColumn(label: Text('Presence')),
      ],
      items: items,
      searchHint: 'Search devices by name, OS, or IP...',
      onSearch: (query) {
        setState(() => _searchQuery = query.trim());
      },
      onAdd: () => context.go('/devices/link'),
      addLabel: 'Register Device',
      onRowNavigate: (device) {
        context.go(
          '/devices/detail/${Uri.encodeComponent(device.id)}',
          extra: device,
        );
      },
      rowBuilder: (device, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              device.name.isNotEmpty ? device.name : 'Unnamed Device',
            )),
            DataCell(Text(device.os)),
            DataCell(Text(device.ip)),
            DataCell(Text(device.lastSeen)),
            DataCell(Text(
              device.profileId.isNotEmpty
                  ? device.profileId.length > 12
                      ? '${device.profileId.substring(0, 12)}...'
                      : device.profileId
                  : '-',
            )),
            DataCell(PresenceIndicator(
              status: device.presence,
              size: 10,
              showLabel: true,
            )),
          ],
        );
      },
      exportRow: (device) => [
        device.name,
        device.os,
        device.ip,
        device.lastSeen,
        device.profileId,
        device.presence.name,
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Devices rows as $format'),
      detailBuilder: (device) => DeviceCard(
        device: device,
        onTap: () {
          context.go(
            '/devices/detail/${Uri.encodeComponent(device.id)}',
            extra: device,
          );
        },
      ),
    );
  }

  void _refresh() {
    ref.invalidate(
      deviceSearchProvider(_searchQuery),
    );
  }
}
