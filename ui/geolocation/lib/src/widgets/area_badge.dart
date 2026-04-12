import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/area_providers.dart';
import 'area_type_badge.dart';

/// Resolves an area by ID and shows the area name with its type badge.
///
/// Drop this into ANY screen to display an area from just its ID.
///
/// ```dart
/// AreaBadge(areaId: route.areaId)
/// AreaBadge(areaId: event.areaId, showType: false)
/// ```
class AreaBadge extends ConsumerWidget {
  const AreaBadge({
    super.key,
    required this.areaId,
    this.showType = true,
    this.style,
  });

  final String areaId;
  final bool showType;
  final TextStyle? style;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (areaId.isEmpty) return const SizedBox.shrink();

    final theme = Theme.of(context);
    final areaAsync = ref.watch(getAreaProvider(areaId));
    final fallback = areaId.length > 12
        ? '${areaId.substring(0, 12)}...'
        : areaId;

    return areaAsync.when(
      data: (area) {
        final name = area.name.isNotEmpty ? area.name : fallback;
        return Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.place_outlined,
              size: 16,
              color: theme.colorScheme.onSurfaceVariant,
            ),
            const SizedBox(width: 6),
            Flexible(
              child: Text(
                name,
                style: style ??
                    theme.textTheme.bodyMedium?.copyWith(
                      fontWeight: FontWeight.w500,
                    ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            if (showType) ...[
              const SizedBox(width: 8),
              AreaTypeBadge(areaType: area.areaType),
            ],
          ],
        );
      },
      loading: () => Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
            width: 14,
            height: 14,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
          const SizedBox(width: 6),
          Text(fallback, style: style ?? theme.textTheme.bodyMedium),
        ],
      ),
      error: (_, _) => Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.place_outlined,
            size: 16,
            color: theme.colorScheme.onSurfaceVariant,
          ),
          const SizedBox(width: 6),
          Text(fallback, style: style ?? theme.textTheme.bodyMedium),
        ],
      ),
    );
  }
}
