import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/location_providers.dart';

/// Screen showing location history as a paginated DataTable with CSV export.
class LocationTrackScreen extends ConsumerWidget {
  const LocationTrackScreen({super.key, required this.subjectId});

  final String subjectId;

  String _sourceName(LocationSource source) {
    return switch (source) {
      LocationSource.LOCATION_SOURCE_GPS => 'GPS',
      LocationSource.LOCATION_SOURCE_NETWORK => 'Network',
      LocationSource.LOCATION_SOURCE_IP => 'IP',
      LocationSource.LOCATION_SOURCE_MANUAL => 'Manual',
      _ => 'Unknown',
    };
  }

  String _formatTimestamp(LocationPointObject point) {
    if (!point.hasTimestamp()) return '-';
    final dt = DateTime.fromMillisecondsSinceEpoch(
      point.timestamp.seconds.toInt() * 1000,
    );
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-'
        '${dt.day.toString().padLeft(2, '0')} '
        '${dt.hour.toString().padLeft(2, '0')}:'
        '${dt.minute.toString().padLeft(2, '0')}:'
        '${dt.second.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncTrack = ref.watch(getTrackProvider(subjectId));

    return asyncTrack.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline, size: 48,
                color: Theme.of(context).colorScheme.error),
            const SizedBox(height: 16),
            Text(friendlyError(error)),
            const SizedBox(height: 16),
            FilledButton.tonal(
              onPressed: () => ref.invalidate(getTrackProvider(subjectId)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (points) => _buildTable(context, ref, points),
    );
  }

  Widget _buildTable(
    BuildContext context,
    WidgetRef ref,
    List<LocationPointObject> points,
  ) {
    return AdminEntityListPage<LocationPointObject>(
      title: 'Location Track',
      breadcrumbs: const ['Home', 'Geolocation', 'Track'],
      columns: const [
        DataColumn(label: Text('Latitude')),
        DataColumn(label: Text('Longitude')),
        DataColumn(label: Text('Accuracy')),
        DataColumn(label: Text('Speed')),
        DataColumn(label: Text('Source')),
        DataColumn(label: Text('Timestamp')),
      ],
      items: points,
      searchHint: 'Subject: $subjectId',
      rowBuilder: (point, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              point.latitude.toStringAsFixed(6),
              style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
            )),
            DataCell(Text(
              point.longitude.toStringAsFixed(6),
              style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
            )),
            DataCell(Text(
              point.accuracy > 0
                  ? '${point.accuracy.toStringAsFixed(0)}m'
                  : '-',
            )),
            DataCell(Text(
              point.speed > 0
                  ? '${point.speed.toStringAsFixed(1)} m/s'
                  : '-',
            )),
            DataCell(Text(_sourceName(point.source))),
            DataCell(Text(_formatTimestamp(point))),
          ],
        );
      },
      exportRow: (point) => [
        point.latitude.toStringAsFixed(6),
        point.longitude.toStringAsFixed(6),
        point.accuracy > 0 ? point.accuracy.toStringAsFixed(0) : '',
        point.speed > 0 ? point.speed.toStringAsFixed(1) : '',
        _sourceName(point.source),
        _formatTimestamp(point),
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Location Track rows as $format'),
      actions: [
        IconButton(
          icon: const Icon(Icons.refresh, size: 20),
          tooltip: 'Refresh',
          onPressed: () => ref.invalidate(getTrackProvider(subjectId)),
        ),
      ],
    );
  }
}
