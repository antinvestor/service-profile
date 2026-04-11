import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter/material.dart';

/// Card displaying a geo event with type icon, area ID, and timestamp.
class GeoEventTile extends StatelessWidget {
  const GeoEventTile({super.key, required this.event, this.onTap});

  final GeoEventObject event;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final (icon, color, label) = _eventTypeDisplay(event.eventType);
    final timestamp = event.hasTimestamp()
        ? DateTime.fromMillisecondsSinceEpoch(
            event.timestamp.seconds.toInt() * 1000,
          )
        : null;

    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: color.withAlpha(25),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(icon, color: color, size: 20),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      label,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      'Area: ${event.areaId}',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    if (event.confidence > 0)
                      Text(
                        'Confidence: ${(event.confidence * 100).toStringAsFixed(0)}%',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                  ],
                ),
              ),
              if (timestamp != null)
                Text(
                  _formatTimestamp(timestamp),
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }

  (IconData, Color, String) _eventTypeDisplay(GeoEventType type) {
    return switch (type) {
      GeoEventType.GEO_EVENT_TYPE_ENTER => (
        Icons.login,
        Colors.green,
        'Enter',
      ),
      GeoEventType.GEO_EVENT_TYPE_EXIT => (
        Icons.logout,
        Colors.red,
        'Exit',
      ),
      GeoEventType.GEO_EVENT_TYPE_DWELL => (
        Icons.hourglass_bottom,
        Colors.amber,
        'Dwell',
      ),
      _ => (
        Icons.help_outline,
        Colors.grey,
        'Unknown',
      ),
    };
  }

  String _formatTimestamp(DateTime dt) {
    final now = DateTime.now();
    final diff = now.difference(dt);
    if (diff.inMinutes < 1) return 'Just now';
    if (diff.inMinutes < 60) return '${diff.inMinutes}m ago';
    if (diff.inHours < 24) return '${diff.inHours}h ago';
    if (diff.inDays < 7) return '${diff.inDays}d ago';
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';
  }
}
