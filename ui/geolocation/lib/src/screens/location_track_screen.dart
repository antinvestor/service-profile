import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/location_providers.dart';
import '../widgets/location_point_tile.dart';

/// Screen showing location history as a timeline of LocationPointTile widgets.
class LocationTrackScreen extends ConsumerWidget {
  const LocationTrackScreen({super.key, required this.subjectId});

  final String subjectId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncTrack = ref.watch(getTrackProvider(subjectId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Location Track'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            tooltip: 'Refresh',
            onPressed: () => ref.invalidate(getTrackProvider(subjectId)),
          ),
        ],
      ),
      body: asyncTrack.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.error_outline,
                  size: 48, color: theme.colorScheme.error),
              const SizedBox(height: 16),
              Text(friendlyError(error)),
              const SizedBox(height: 16),
              FilledButton.tonal(
                onPressed: () =>
                    ref.invalidate(getTrackProvider(subjectId)),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
        data: (points) {
          if (points.isEmpty) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    width: 64,
                    height: 64,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surfaceContainerLow,
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Icon(
                      Icons.location_off_outlined,
                      size: 28,
                      color: theme.colorScheme.primary.withAlpha(160),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'No location points found',
                    style: theme.textTheme.titleMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'Subject: $subjectId',
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                ],
              ),
            );
          }

          return Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(24, 16, 24, 8),
                child: Row(
                  children: [
                    Icon(Icons.timeline, size: 20,
                        color: theme.colorScheme.primary),
                    const SizedBox(width: 8),
                    Text(
                      '${points.length} location points',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const Spacer(),
                    Text(
                      'Subject: $subjectId',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ),
              ),
              Expanded(
                child: ListView.separated(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 24, vertical: 8),
                  itemCount: points.length,
                  separatorBuilder: (_, __) => const SizedBox(height: 4),
                  itemBuilder: (context, index) {
                    return Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Timeline connector
                        SizedBox(
                          width: 24,
                          child: Column(
                            children: [
                              if (index > 0)
                                Container(
                                  width: 2,
                                  height: 8,
                                  color: theme.colorScheme.outlineVariant,
                                ),
                              Container(
                                width: 10,
                                height: 10,
                                decoration: BoxDecoration(
                                  color: theme.colorScheme.primary,
                                  shape: BoxShape.circle,
                                ),
                              ),
                              if (index < points.length - 1)
                                Container(
                                  width: 2,
                                  height: 60,
                                  color: theme.colorScheme.outlineVariant,
                                ),
                            ],
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: LocationPointTile(point: points[index]),
                        ),
                      ],
                    );
                  },
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}
