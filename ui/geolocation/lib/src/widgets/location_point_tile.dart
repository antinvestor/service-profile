import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter/material.dart';

/// Tile showing a location point with lat/lon, source icon, accuracy,
/// and timestamp.
class LocationPointTile extends StatelessWidget {
  const LocationPointTile({super.key, required this.point});

  final LocationPointObject point;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final (icon, sourceLabel) = _sourceDisplay(point.source);
    final timestamp = point.hasTimestamp()
        ? DateTime.fromMillisecondsSinceEpoch(
            point.timestamp.seconds.toInt() * 1000,
          )
        : null;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                icon,
                color: theme.colorScheme.onPrimaryContainer,
                size: 20,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '${point.latitude.toStringAsFixed(6)}, '
                    '${point.longitude.toStringAsFixed(6)}',
                    style: theme.textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w600,
                      fontFamily: 'monospace',
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      _InfoChip(
                        label: sourceLabel,
                        icon: icon,
                        theme: theme,
                      ),
                      if (point.accuracy > 0) ...[
                        const SizedBox(width: 8),
                        _InfoChip(
                          label:
                              '${point.accuracy.toStringAsFixed(0)}m',
                          icon: Icons.gps_fixed,
                          theme: theme,
                        ),
                      ],
                      if (point.speed > 0) ...[
                        const SizedBox(width: 8),
                        _InfoChip(
                          label:
                              '${point.speed.toStringAsFixed(1)} m/s',
                          icon: Icons.speed,
                          theme: theme,
                        ),
                      ],
                    ],
                  ),
                ],
              ),
            ),
            if (timestamp != null)
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${timestamp.hour.toString().padLeft(2, '0')}:'
                    '${timestamp.minute.toString().padLeft(2, '0')}:'
                    '${timestamp.second.toString().padLeft(2, '0')}',
                    style: theme.textTheme.labelMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    '${timestamp.year}-'
                    '${timestamp.month.toString().padLeft(2, '0')}-'
                    '${timestamp.day.toString().padLeft(2, '0')}',
                    style: theme.textTheme.labelSmall?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                ],
              ),
          ],
        ),
      ),
    );
  }

  (IconData, String) _sourceDisplay(LocationSource source) {
    return switch (source) {
      LocationSource.LOCATION_SOURCE_GPS => (Icons.gps_fixed, 'GPS'),
      LocationSource.LOCATION_SOURCE_NETWORK =>
        (Icons.cell_tower, 'Network'),
      LocationSource.LOCATION_SOURCE_IP => (Icons.language, 'IP'),
      LocationSource.LOCATION_SOURCE_MANUAL =>
        (Icons.edit_location_alt, 'Manual'),
      _ => (Icons.location_on_outlined, 'Unknown'),
    };
  }
}

class _InfoChip extends StatelessWidget {
  const _InfoChip({
    required this.label,
    required this.icon,
    required this.theme,
  });

  final String label;
  final IconData icon;
  final ThemeData theme;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        label,
        style: theme.textTheme.labelSmall?.copyWith(
          color: theme.colorScheme.onSurfaceVariant,
        ),
      ),
    );
  }
}
