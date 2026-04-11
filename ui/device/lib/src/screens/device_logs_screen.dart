import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/device_log_providers.dart';

/// Screen displaying session activity logs for a device as a paginated
/// DataTable with CSV export.
class DeviceLogsScreen extends ConsumerWidget {
  const DeviceLogsScreen({
    super.key,
    required this.deviceId,
  });

  final String deviceId;

  String _locationSummary(DeviceLog log) {
    if (!log.hasLocation()) return '-';
    final fields = log.location.fields;
    if (fields.isEmpty) return '-';
    final parts = <String>[];
    if (fields.containsKey('city')) {
      parts.add(fields['city']!.stringValue);
    }
    if (fields.containsKey('country')) {
      parts.add(fields['country']!.stringValue);
    }
    return parts.isNotEmpty ? parts.join(', ') : '-';
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncLogs = ref.watch(deviceLogsProvider(deviceId));

    return asyncLogs.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline, size: 48,
                color: Theme.of(context).colorScheme.error),
            const SizedBox(height: 16),
            Text(friendlyError(error)),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () => ref.invalidate(deviceLogsProvider(deviceId)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (logs) => _buildTable(context, ref, logs),
    );
  }

  Widget _buildTable(BuildContext context, WidgetRef ref, List<DeviceLog> logs) {
    return AdminEntityListPage<DeviceLog>(
      title: 'Session Logs',
      breadcrumbs: const ['Home', 'Devices', 'Session Logs'],
      columns: const [
        DataColumn(label: Text('Session')),
        DataColumn(label: Text('IP')),
        DataColumn(label: Text('OS')),
        DataColumn(label: Text('User Agent')),
        DataColumn(label: Text('Last Seen')),
        DataColumn(label: Text('Location')),
      ],
      items: logs,
      searchHint: 'Search logs...',
      rowBuilder: (log, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              log.sessionId.isNotEmpty
                  ? log.sessionId.length > 12
                      ? '${log.sessionId.substring(0, 12)}...'
                      : log.sessionId
                  : '-',
              style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
            )),
            DataCell(Text(log.ip)),
            DataCell(Text(log.os)),
            DataCell(ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 200),
              child: Text(
                log.userAgent,
                overflow: TextOverflow.ellipsis,
                maxLines: 1,
              ),
            )),
            DataCell(Text(log.lastSeen)),
            DataCell(Text(_locationSummary(log))),
          ],
        );
      },
      exportRow: (log) => [
        log.sessionId,
        log.ip,
        log.os,
        log.userAgent,
        log.lastSeen,
        _locationSummary(log),
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Session Logs rows as $format'),
      actions: [
        IconButton(
          icon: const Icon(Icons.refresh, size: 20),
          tooltip: 'Refresh',
          onPressed: () => ref.invalidate(deviceLogsProvider(deviceId)),
        ),
      ],
    );
  }
}
