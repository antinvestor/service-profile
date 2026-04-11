import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/route_providers.dart';

/// Screen for listing and searching routes.
class RouteListScreen extends ConsumerStatefulWidget {
  const RouteListScreen({super.key});

  @override
  ConsumerState<RouteListScreen> createState() => _RouteListScreenState();
}

class _RouteListScreenState extends ConsumerState<RouteListScreen> {
  String _ownerId = '';

  @override
  Widget build(BuildContext context) {
    final asyncRoutes = ref.watch(searchRoutesProvider(_ownerId));

    return asyncRoutes.when(
      loading: () => _buildShell(isLoading: true, items: const []),
      error: (error, _) =>
          _buildShell(error: friendlyError(error), items: const []),
      data: (routes) => _buildShell(items: routes),
    );
  }

  Widget _buildShell({
    required List<RouteObject> items,
    bool isLoading = false,
    String? error,
  }) {
    return EntityListPage<RouteObject>(
      title: 'Routes',
      icon: Icons.route_outlined,
      items: items,
      isLoading: isLoading,
      error: error,
      onRetry: () => ref.invalidate(searchRoutesProvider(_ownerId)),
      searchHint: 'Filter by owner ID...',
      onSearchChanged: (query) {
        setState(() => _ownerId = query.trim());
      },
      actionLabel: 'New Route',
      onAction: () => context.go('/geo/routes/new'),
      itemBuilder: (context, route) {
        return _RouteCard(
          route: route,
          onTap: () => context.go('/geo/routes/${route.id}'),
        );
      },
    );
  }
}

class _RouteCard extends StatelessWidget {
  const _RouteCard({required this.route, this.onTap});

  final RouteObject route;
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
                  color: theme.colorScheme.secondaryContainer,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  Icons.route_outlined,
                  color: theme.colorScheme.onSecondaryContainer,
                  size: 22,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      route.name.isNotEmpty ? route.name : 'Unnamed Route',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        if (route.lengthM > 0) ...[
                          Icon(Icons.straighten, size: 14,
                              color: theme.colorScheme.onSurfaceVariant),
                          const SizedBox(width: 4),
                          Text(
                            _formatLength(route.lengthM),
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                          const SizedBox(width: 12),
                        ],
                        if (route.deviationThresholdM > 0) ...[
                          Icon(Icons.warning_amber, size: 14,
                              color: theme.colorScheme.onSurfaceVariant),
                          const SizedBox(width: 4),
                          Text(
                            '${route.deviationThresholdM.toStringAsFixed(0)}m threshold',
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                        ],
                      ],
                    ),
                    if (route.description.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(
                        route.description,
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

  String _formatLength(double meters) {
    if (meters < 1000) return '${meters.toStringAsFixed(0)} m';
    return '${(meters / 1000).toStringAsFixed(1)} km';
  }
}
