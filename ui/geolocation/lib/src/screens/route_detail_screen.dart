import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/route_providers.dart';
import '../widgets/route_assignment_chip.dart';

/// Detail screen for a single route showing waypoints, assignments,
/// and description.
class RouteDetailScreen extends ConsumerWidget {
  const RouteDetailScreen({super.key, required this.routeId});

  final String routeId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncRoute = ref.watch(getRouteProvider(routeId));

    return asyncRoute.when(
      loading: () => const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      ),
      error: (error, _) => Scaffold(
        appBar: AppBar(title: const Text('Route')),
        body: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.error_outline,
                  size: 48,
                  color: Theme.of(context).colorScheme.error),
              const SizedBox(height: 16),
              Text(friendlyError(error)),
              const SizedBox(height: 16),
              FilledButton.tonal(
                onPressed: () =>
                    ref.invalidate(getRouteProvider(routeId)),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
      data: (route) =>
          _RouteDetailBody(route: route, routeId: routeId),
    );
  }
}

class _RouteDetailBody extends ConsumerStatefulWidget {
  const _RouteDetailBody({required this.route, required this.routeId});

  final RouteObject route;
  final String routeId;

  @override
  ConsumerState<_RouteDetailBody> createState() =>
      _RouteDetailBodyState();
}

class _RouteDetailBodyState extends ConsumerState<_RouteDetailBody> {
  final _subjectIdController = TextEditingController();

  @override
  void dispose() {
    _subjectIdController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final route = widget.route;

    return Scaffold(
      appBar: AppBar(
        title: Text(route.name.isNotEmpty ? route.name : 'Route'),
        actions: [
          IconButton(
            icon: const Icon(Icons.delete_outline),
            tooltip: 'Delete Route',
            onPressed: () => _confirmDelete(context),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
          // Header
          Row(
            children: [
              Container(
                width: 56,
                height: 56,
                decoration: BoxDecoration(
                  color: theme.colorScheme.secondaryContainer,
                  borderRadius: BorderRadius.circular(14),
                ),
                child: Icon(
                  Icons.route_outlined,
                  color: theme.colorScheme.onSecondaryContainer,
                  size: 28,
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      route.name.isNotEmpty ? route.name : 'Unnamed Route',
                      style: theme.textTheme.headlineSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    if (route.lengthM > 0)
                      Text(
                        _formatLength(route.lengthM),
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Description
          if (route.description.isNotEmpty) ...[
            Text(
              'Description',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.primary,
              ),
            ),
            const SizedBox(height: 8),
            Text(route.description, style: theme.textTheme.bodyMedium),
            const SizedBox(height: 24),
          ],

          // Properties
          Text(
            'Properties',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w700,
              color: theme.colorScheme.primary,
            ),
          ),
          const SizedBox(height: 12),
          _PropertyRow(label: 'ID', value: route.id),
          _PropertyRow(label: 'Owner ID', value: route.ownerId),
          _PropertyRow(label: 'State', value: route.state.toString()),
          if (route.deviationThresholdM > 0)
            _PropertyRow(
              label: 'Dev. Threshold',
              value: '${route.deviationThresholdM.toStringAsFixed(0)} m',
            ),
          if (route.deviationConsecutiveCount > 0)
            _PropertyRow(
              label: 'Dev. Count',
              value: route.deviationConsecutiveCount.toString(),
            ),
          if (route.deviationCooldownSec > 0)
            _PropertyRow(
              label: 'Dev. Cooldown',
              value: '${route.deviationCooldownSec} s',
            ),
          const SizedBox(height: 24),

          // Geometry / Waypoints
          if (route.geometry.isNotEmpty) ...[
            Text(
              'Waypoints / Geometry',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.primary,
              ),
            ),
            const SizedBox(height: 8),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: theme.colorScheme.surfaceContainerLow,
                borderRadius: BorderRadius.circular(8),
              ),
              child: SelectableText(
                route.geometry,
                style: theme.textTheme.bodySmall?.copyWith(
                  fontFamily: 'monospace',
                ),
              ),
            ),
            const SizedBox(height: 24),
          ],

          // Assignments
          Text(
            'Assignments',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w700,
              color: theme.colorScheme.primary,
            ),
          ),
          const SizedBox(height: 12),
          _AssignmentsSection(routeId: widget.routeId),

          // Assign new subject
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _subjectIdController,
                  decoration: const InputDecoration(
                    hintText: 'Subject ID to assign',
                    isDense: true,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              FilledButton.tonal(
                onPressed: _assignSubject,
                child: const Text('Assign'),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Future<void> _assignSubject() async {
    final subjectId = _subjectIdController.text.trim();
    if (subjectId.isEmpty) return;
    try {
      await ref.read(routeNotifierProvider.notifier).assignRoute(
            subjectId: subjectId,
            routeId: widget.routeId,
          );
      _subjectIdController.clear();
      ref.invalidate(subjectRouteAssignmentsProvider);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Subject assigned successfully.')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed: ${friendlyError(e)}')),
        );
      }
    }
  }

  Future<void> _confirmDelete(BuildContext context) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Delete Route'),
        content: Text(
            'Are you sure you want to delete "${widget.route.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
    if (confirmed == true && context.mounted) {
      try {
        await ref
            .read(routeNotifierProvider.notifier)
            .deleteRoute(widget.routeId);
        if (context.mounted) {
          context.go('/geo/routes');
        }
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed: ${friendlyError(e)}')),
          );
        }
      }
    }
  }

  String _formatLength(double meters) {
    if (meters < 1000) return '${meters.toStringAsFixed(0)} m';
    return '${(meters / 1000).toStringAsFixed(1)} km';
  }
}

class _AssignmentsSection extends ConsumerWidget {
  const _AssignmentsSection({required this.routeId});

  final String routeId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Fetch assignments for the route by using an empty-string subject query.
    // The route_detail shows assignments related to this route. Since the API
    // fetches by subject, we display any cached assignments or a placeholder.
    final theme = Theme.of(context);
    final asyncAssignments =
        ref.watch(subjectRouteAssignmentsProvider(routeId));

    return asyncAssignments.when(
      loading: () => const Padding(
        padding: EdgeInsets.all(12),
        child: Center(child: CircularProgressIndicator()),
      ),
      error: (error, _) => Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: theme.colorScheme.surfaceContainerLow,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          'Failed to load assignments: ${friendlyError(error)}',
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.error,
          ),
        ),
      ),
      data: (assignments) {
        if (assignments.isEmpty) {
          return Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: theme.colorScheme.surfaceContainerLow,
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              'No subjects assigned to this route yet.',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          );
        }
        return Column(
          children: assignments.map((assignment) {
            return RouteAssignmentChip(
              assignment: assignment,
              onUnassign: () async {
                await ref
                    .read(routeNotifierProvider.notifier)
                    .unassignRoute(assignment.id);
                ref.invalidate(
                    subjectRouteAssignmentsProvider(routeId));
              },
            );
          }).toList(),
        );
      },
    );
  }
}

class _PropertyRow extends StatelessWidget {
  const _PropertyRow({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: theme.textTheme.labelMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          Expanded(
            child: SelectableText(
              value.isNotEmpty ? value : '-',
              style: theme.textTheme.bodyMedium,
            ),
          ),
        ],
      ),
    );
  }
}
