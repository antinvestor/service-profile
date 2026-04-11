import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/area_providers.dart';
import '../widgets/area_type_badge.dart';

/// Detail screen for a single area, showing properties, boundary, and events.
class AreaDetailScreen extends ConsumerWidget {
  const AreaDetailScreen({super.key, required this.areaId});

  final String areaId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncArea = ref.watch(getAreaProvider(areaId));

    return asyncArea.when(
      loading: () => const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      ),
      error: (error, _) => Scaffold(
        appBar: AppBar(title: const Text('Area')),
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
                    ref.invalidate(getAreaProvider(areaId)),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
      data: (area) => _AreaDetailBody(area: area, areaId: areaId),
    );
  }
}

class _AreaDetailBody extends ConsumerWidget {
  const _AreaDetailBody({required this.area, required this.areaId});

  final AreaObject area;
  final String areaId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(area.name.isNotEmpty ? area.name : 'Area'),
        actions: [
          IconButton(
            icon: const Icon(Icons.delete_outline),
            tooltip: 'Delete Area',
            onPressed: () => _confirmDelete(context, ref),
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
                  color: theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(14),
                ),
                child: Icon(
                  Icons.map_outlined,
                  color: theme.colorScheme.onPrimaryContainer,
                  size: 28,
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      area.name.isNotEmpty ? area.name : 'Unnamed Area',
                      style: theme.textTheme.headlineSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 6),
                    AreaTypeBadge(areaType: area.areaType),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Description
          if (area.description.isNotEmpty) ...[
            Text(
              'Description',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.primary,
              ),
            ),
            const SizedBox(height: 8),
            Text(area.description, style: theme.textTheme.bodyMedium),
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
          _PropertyRow(label: 'ID', value: area.id),
          _PropertyRow(label: 'Owner ID', value: area.ownerId),
          if (area.areaM2 > 0)
            _PropertyRow(
              label: 'Area',
              value: '${area.areaM2.toStringAsFixed(2)} m\u00B2',
            ),
          if (area.perimeterM > 0)
            _PropertyRow(
              label: 'Perimeter',
              value: '${area.perimeterM.toStringAsFixed(2)} m',
            ),
          _PropertyRow(label: 'State', value: area.state.toString()),
          const SizedBox(height: 24),

          // Geometry / Boundary points
          if (area.geometry.isNotEmpty) ...[
            Text(
              'Geometry',
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
                area.geometry,
                style: theme.textTheme.bodySmall?.copyWith(
                  fontFamily: 'monospace',
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Future<void> _confirmDelete(BuildContext context, WidgetRef ref) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Delete Area'),
        content: Text(
            'Are you sure you want to delete "${area.name}"?'),
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
        await ref.read(areaNotifierProvider.notifier).deleteArea(areaId);
        if (context.mounted) {
          context.go('/geo/areas');
        }
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to delete: ${friendlyError(e)}')),
          );
        }
      }
    }
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
