import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/route_providers.dart';

/// Screen for listing and searching routes with a paginated DataTable
/// and CSV export.
class RouteListScreen extends ConsumerStatefulWidget {
  const RouteListScreen({super.key});

  @override
  ConsumerState<RouteListScreen> createState() => _RouteListScreenState();
}

class _RouteListScreenState extends ConsumerState<RouteListScreen> {
  String _ownerId = '';

  String _stateLabel(int state) {
    return switch (state) {
      0 => 'Inactive',
      1 => 'Active',
      _ => 'State $state',
    };
  }

  String _formatLength(double meters) {
    if (meters <= 0) return '-';
    if (meters < 1000) return '${meters.toStringAsFixed(0)} m';
    return '${(meters / 1000).toStringAsFixed(1)} km';
  }

  int _waypointCount(RouteObject route) {
    // The geometry field holds GeoJSON; count coordinates as a proxy
    // for waypoints. If geometry is empty, show 0.
    if (route.geometry.isEmpty) return 0;
    // Count coordinate pairs by occurrences of ','
    return RegExp(r'\[[-\d]').allMatches(route.geometry).length;
  }

  @override
  Widget build(BuildContext context) {
    final asyncRoutes = ref.watch(searchRoutesProvider(_ownerId));

    return asyncRoutes.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(friendlyError(error)),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () =>
                  ref.invalidate(searchRoutesProvider(_ownerId)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (routes) => _buildTable(routes),
    );
  }

  Widget _buildTable(List<RouteObject> items) {
    return AdminEntityListPage<RouteObject>(
      title: 'Routes',
      breadcrumbs: const ['Home', 'Geolocation', 'Routes'],
      columns: const [
        DataColumn(label: Text('Name')),
        DataColumn(label: Text('Description')),
        DataColumn(label: Text('Length')),
        DataColumn(label: Text('State')),
      ],
      items: items,
      searchHint: 'Filter by owner ID...',
      onSearch: (query) {
        setState(() => _ownerId = query.trim());
      },
      onAdd: () => context.go('/geo/routes/new'),
      addLabel: 'New Route',
      onRowNavigate: (route) => context.go('/geo/routes/${route.id}'),
      rowBuilder: (route, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              route.name.isNotEmpty ? route.name : 'Unnamed Route',
              style: const TextStyle(fontWeight: FontWeight.w600),
            )),
            DataCell(ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 200),
              child: Text(
                route.description.isNotEmpty ? route.description : '-',
                overflow: TextOverflow.ellipsis,
                maxLines: 1,
              ),
            )),
            DataCell(Text(_formatLength(route.lengthM))),
            DataCell(Text(_stateLabel(route.state))),
          ],
        );
      },
      exportRow: (route) => [
        route.name,
        route.description,
        _formatLength(route.lengthM),
        _stateLabel(route.state),
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Routes rows as $format'),
    );
  }
}
