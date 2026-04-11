import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/area_providers.dart';
import '../widgets/area_type_badge.dart';

/// Screen for listing and searching areas with type filtering.
class AreaListScreen extends ConsumerStatefulWidget {
  const AreaListScreen({super.key});

  @override
  ConsumerState<AreaListScreen> createState() => _AreaListScreenState();
}

class _AreaListScreenState extends ConsumerState<AreaListScreen> {
  String _query = '';
  AreaType? _selectedType;

  @override
  Widget build(BuildContext context) {
    final asyncAreas = ref.watch(searchAreasProvider(_query));

    return asyncAreas.when(
      loading: () => _buildShell(isLoading: true, items: const []),
      error: (error, _) =>
          _buildShell(error: friendlyError(error), items: const []),
      data: (areas) {
        final filtered = _selectedType != null
            ? areas.where((a) => a.areaType == _selectedType).toList()
            : areas;
        return _buildShell(items: filtered);
      },
    );
  }

  Widget _buildShell({
    required List<AreaObject> items,
    bool isLoading = false,
    String? error,
  }) {
    return EntityListPage<AreaObject>(
      title: 'Areas',
      icon: Icons.map_outlined,
      items: items,
      isLoading: isLoading,
      error: error,
      onRetry: () => ref.invalidate(searchAreasProvider(_query)),
      searchHint: 'Search areas by name...',
      onSearchChanged: (query) {
        setState(() => _query = query.trim());
      },
      actionLabel: 'New Area',
      onAction: () => context.go('/geo/areas/new'),
      filterWidget: _buildTypeFilter(),
      itemBuilder: (context, area) {
        return _AreaCard(
          area: area,
          onTap: () => context.go('/geo/areas/${area.id}'),
        );
      },
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

class _AreaCard extends StatelessWidget {
  const _AreaCard({required this.area, this.onTap});

  final AreaObject area;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 44,
                height: 44,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  Icons.map_outlined,
                  color: theme.colorScheme.onPrimaryContainer,
                  size: 22,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      area.name.isNotEmpty ? area.name : 'Unnamed Area',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        AreaTypeBadge(areaType: area.areaType),
                        if (area.areaM2 > 0) ...[
                          const SizedBox(width: 8),
                          Text(
                            '${area.areaM2.toStringAsFixed(0)} m\u00B2',
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                        ],
                      ],
                    ),
                    if (area.description.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(
                        area.description,
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ],
                  ],
                ),
              ),
              const Icon(Icons.chevron_right, size: 20),
            ],
          ),
        ),
      ),
    );
  }
}
