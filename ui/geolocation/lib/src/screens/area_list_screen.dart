import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/area_providers.dart';
import '../widgets/area_type_badge.dart';

/// Screen for listing and searching areas with a paginated DataTable
/// and CSV export.
class AreaListScreen extends ConsumerStatefulWidget {
  const AreaListScreen({super.key});

  @override
  ConsumerState<AreaListScreen> createState() => _AreaListScreenState();
}

class _AreaListScreenState extends ConsumerState<AreaListScreen> {
  String _query = '';
  AreaType? _selectedType;

  String _areaTypeName(AreaType type) {
    return switch (type) {
      AreaType.AREA_TYPE_LAND => 'Land',
      AreaType.AREA_TYPE_BUILDING => 'Building',
      AreaType.AREA_TYPE_ZONE => 'Zone',
      AreaType.AREA_TYPE_FENCE => 'Fence',
      AreaType.AREA_TYPE_CUSTOM => 'Custom',
      _ => 'Unknown',
    };
  }

  String _stateLabel(int state) {
    return switch (state) {
      0 => 'Inactive',
      1 => 'Active',
      _ => 'State $state',
    };
  }

  @override
  Widget build(BuildContext context) {
    final asyncAreas = ref.watch(searchAreasProvider(_query));

    return asyncAreas.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(friendlyError(error)),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () => ref.invalidate(searchAreasProvider(_query)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (areas) {
        final filtered = _selectedType != null
            ? areas.where((a) => a.areaType == _selectedType).toList()
            : areas;
        return _buildTable(filtered);
      },
    );
  }

  Widget _buildTable(List<AreaObject> items) {
    return AdminEntityListPage<AreaObject>(
      title: 'Areas',
      breadcrumbs: const ['Home', 'Geolocation', 'Areas'],
      columns: const [
        DataColumn(label: Text('Name')),
        DataColumn(label: Text('Type')),
        DataColumn(label: Text('Description')),
        DataColumn(label: Text('State')),
      ],
      items: items,
      searchHint: 'Search areas by name...',
      onSearch: (query) {
        setState(() => _query = query.trim());
      },
      onAdd: () => context.go('/geo/areas/new'),
      addLabel: 'New Area',
      onRowNavigate: (area) => context.go('/geo/areas/${area.id}'),
      actions: [
        _buildTypeFilter(),
      ],
      rowBuilder: (area, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              area.name.isNotEmpty ? area.name : 'Unnamed Area',
              style: const TextStyle(fontWeight: FontWeight.w600),
            )),
            DataCell(AreaTypeBadge(areaType: area.areaType)),
            DataCell(ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 200),
              child: Text(
                area.description.isNotEmpty ? area.description : '-',
                overflow: TextOverflow.ellipsis,
                maxLines: 1,
              ),
            )),
            DataCell(Text(_stateLabel(area.state))),
          ],
        );
      },
      exportRow: (area) => [
        area.name,
        _areaTypeName(area.areaType),
        area.description,
        _stateLabel(area.state),
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Areas rows as $format'),
    );
  }

  Widget _buildTypeFilter() {
    return PopupMenuButton<AreaType?>(
      icon: Icon(
        _selectedType != null ? Icons.filter_alt : Icons.filter_alt_outlined,
        size: 20,
      ),
      tooltip: 'Filter by type',
      onSelected: (type) => setState(() => _selectedType = type),
      itemBuilder: (context) => [
        const PopupMenuItem(value: null, child: Text('All Types')),
        const PopupMenuDivider(),
        const PopupMenuItem(
          value: AreaType.AREA_TYPE_LAND,
          child: Text('Land'),
        ),
        const PopupMenuItem(
          value: AreaType.AREA_TYPE_BUILDING,
          child: Text('Building'),
        ),
        const PopupMenuItem(
          value: AreaType.AREA_TYPE_ZONE,
          child: Text('Zone'),
        ),
        const PopupMenuItem(
          value: AreaType.AREA_TYPE_FENCE,
          child: Text('Fence'),
        ),
        const PopupMenuItem(
          value: AreaType.AREA_TYPE_CUSTOM,
          child: Text('Custom'),
        ),
      ],
    );
  }
}
