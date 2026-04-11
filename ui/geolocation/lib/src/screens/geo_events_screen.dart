import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/location_providers.dart';
import '../widgets/geo_event_tile.dart';

/// Screen for listing geo events with filtering by event type.
class GeoEventsScreen extends ConsumerStatefulWidget {
  const GeoEventsScreen({super.key});

  @override
  ConsumerState<GeoEventsScreen> createState() => _GeoEventsScreenState();
}

class _GeoEventsScreenState extends ConsumerState<GeoEventsScreen> {
  String _subjectId = '';
  GeoEventType? _selectedType;

  @override
  Widget build(BuildContext context) {
    if (_subjectId.isEmpty) {
      return _buildSubjectInput();
    }

    final asyncEvents = ref.watch(getSubjectEventsProvider(_subjectId));

    return asyncEvents.when(
      loading: () => _buildShell(isLoading: true, items: const []),
      error: (error, _) =>
          _buildShell(error: friendlyError(error), items: const []),
      data: (events) {
        final filtered = _selectedType != null
            ? events.where((e) => e.eventType == _selectedType).toList()
            : events;
        return _buildShell(items: filtered);
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

  Widget _buildShell({
    required List<GeoEventObject> items,
    bool isLoading = false,
    String? error,
  }) {
    return EntityListPage<GeoEventObject>(
      title: 'Geo Events',
      icon: Icons.event_note,
      items: items,
      isLoading: isLoading,
      error: error,
      onRetry: () =>
          ref.invalidate(getSubjectEventsProvider(_subjectId)),
      searchHint: 'Subject: $_subjectId',
      onSearchChanged: (query) {
        if (query.trim().isNotEmpty) {
          setState(() => _subjectId = query.trim());
        }
      },
      filterWidget: _buildEventTypeFilter(),
      itemBuilder: (context, event) => GeoEventTile(event: event),
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
