import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/location_providers.dart';

/// Screen for listing geo events with filtering by event type, displayed
/// in a paginated DataTable with CSV export.
class GeoEventsScreen extends ConsumerStatefulWidget {
  const GeoEventsScreen({super.key});

  @override
  ConsumerState<GeoEventsScreen> createState() => _GeoEventsScreenState();
}

class _GeoEventsScreenState extends ConsumerState<GeoEventsScreen> {
  String _subjectId = '';
  GeoEventType? _selectedType;

  String _eventTypeName(GeoEventType type) {
    return switch (type) {
      GeoEventType.GEO_EVENT_TYPE_ENTER => 'Enter',
      GeoEventType.GEO_EVENT_TYPE_EXIT => 'Exit',
      GeoEventType.GEO_EVENT_TYPE_DWELL => 'Dwell',
      _ => 'Unknown',
    };
  }

  String _formatTimestamp(GeoEventObject event) {
    if (!event.hasTimestamp()) return '-';
    final dt = DateTime.fromMillisecondsSinceEpoch(
      event.timestamp.seconds.toInt() * 1000,
    );
    final now = DateTime.now();
    final diff = now.difference(dt);
    if (diff.inMinutes < 1) return 'Just now';
    if (diff.inMinutes < 60) return '${diff.inMinutes}m ago';
    if (diff.inHours < 24) return '${diff.inHours}h ago';
    if (diff.inDays < 7) return '${diff.inDays}d ago';
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-'
        '${dt.day.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    if (_subjectId.isEmpty) {
      return _buildSubjectInput();
    }

    final asyncEvents = ref.watch(getSubjectEventsProvider(_subjectId));

    return asyncEvents.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(friendlyError(error)),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () =>
                  ref.invalidate(getSubjectEventsProvider(_subjectId)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (events) {
        final filtered = _selectedType != null
            ? events.where((e) => e.eventType == _selectedType).toList()
            : events;
        return _buildTable(filtered);
      },
    );
  }

  Widget _buildSubjectInput() {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Geo Events')),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.event_note, size: 48,
                  color: theme.colorScheme.primary),
              const SizedBox(height: 16),
              Text(
                'Enter a Subject ID to view events',
                style: theme.textTheme.titleMedium,
              ),
              const SizedBox(height: 16),
              SizedBox(
                width: 400,
                child: TextField(
                  decoration: const InputDecoration(
                    hintText: 'Subject ID',
                    prefixIcon: Icon(Icons.person_outline),
                  ),
                  onSubmitted: (value) {
                    if (value.trim().isNotEmpty) {
                      setState(() => _subjectId = value.trim());
                    }
                  },
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTable(List<GeoEventObject> items) {
    return AdminEntityListPage<GeoEventObject>(
      title: 'Geo Events',
      breadcrumbs: const ['Home', 'Geolocation', 'Events'],
      columns: const [
        DataColumn(label: Text('Subject')),
        DataColumn(label: Text('Area')),
        DataColumn(label: Text('Event Type')),
        DataColumn(label: Text('Timestamp')),
      ],
      items: items,
      searchHint: 'Subject: $_subjectId',
      onSearch: (query) {
        if (query.trim().isNotEmpty) {
          setState(() => _subjectId = query.trim());
        }
      },
      actions: [
        _buildEventTypeFilter(),
      ],
      rowBuilder: (event, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              event.subjectId.isNotEmpty
                  ? event.subjectId.length > 16
                      ? '${event.subjectId.substring(0, 13)}...'
                      : event.subjectId
                  : '-',
            )),
            DataCell(Text(
              event.areaId.isNotEmpty
                  ? event.areaId.length > 16
                      ? '${event.areaId.substring(0, 13)}...'
                      : event.areaId
                  : '-',
            )),
            DataCell(Text(_eventTypeName(event.eventType))),
            DataCell(Text(_formatTimestamp(event))),
          ],
        );
      },
      exportRow: (event) => [
        event.subjectId,
        event.areaId,
        _eventTypeName(event.eventType),
        _formatTimestamp(event),
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Geo Events rows as $format'),
    );
  }

  Widget _buildEventTypeFilter() {
    return PopupMenuButton<GeoEventType?>(
      icon: Icon(
        _selectedType != null ? Icons.filter_alt : Icons.filter_alt_outlined,
        size: 20,
      ),
      tooltip: 'Filter by event type',
      onSelected: (type) => setState(() => _selectedType = type),
      itemBuilder: (context) => [
        const PopupMenuItem(value: null, child: Text('All Events')),
        const PopupMenuDivider(),
        const PopupMenuItem(
          value: GeoEventType.GEO_EVENT_TYPE_ENTER,
          child: Text('Enter'),
        ),
        const PopupMenuItem(
          value: GeoEventType.GEO_EVENT_TYPE_EXIT,
          child: Text('Exit'),
        ),
        const PopupMenuItem(
          value: GeoEventType.GEO_EVENT_TYPE_DWELL,
          child: Text('Dwell'),
        ),
      ],
    );
  }
}
