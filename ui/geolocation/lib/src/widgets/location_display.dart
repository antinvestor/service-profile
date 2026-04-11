import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/// Displays a latitude/longitude pair in a formatted, human-readable way.
///
/// Optionally shows a copy button and a label.
///
/// ```dart
/// LocationDisplay(latitude: -1.2921, longitude: 36.8219)
/// LocationDisplay(latitude: lat, longitude: lon, label: 'Home')
/// ```
class LocationDisplay extends StatelessWidget {
  const LocationDisplay({
    super.key,
    required this.latitude,
    required this.longitude,
    this.label,
    this.style,
    this.copiable = false,
    this.compact = false,
    this.icon = true,
  });

  final double latitude;
  final double longitude;
  final String? label;
  final TextStyle? style;
  final bool copiable;
  final bool compact;

  /// Whether to show the location pin icon.
  final bool icon;

  String _formatCoord(double value, bool isLatitude) {
    final dir = isLatitude
        ? (value >= 0 ? 'N' : 'S')
        : (value >= 0 ? 'E' : 'W');
    return '${value.abs().toStringAsFixed(6)}$dir';
  }

  String get _formatted =>
      '${_formatCoord(latitude, true)}, ${_formatCoord(longitude, false)}';

  String get _copyText => '$latitude, $longitude';

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final textStyle = style ?? theme.textTheme.bodyMedium;

    if (compact) {
      return Text(_formatted, style: textStyle);
    }

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        if (icon) ...[
          Icon(
            Icons.location_on_outlined,
            size: 16,
            color: theme.colorScheme.onSurfaceVariant,
          ),
          const SizedBox(width: 6),
        ],
        Flexible(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              if (label != null && label!.isNotEmpty)
                Text(
                  label!,
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              Text(_formatted, style: textStyle),
            ],
          ),
        ),
        if (copiable) ...[
          const SizedBox(width: 4),
          Tooltip(
            message: 'Copy coordinates',
            child: InkWell(
              borderRadius: BorderRadius.circular(12),
              onTap: () {
                Clipboard.setData(ClipboardData(text: _copyText));
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('Coordinates copied'),
                    duration: Duration(seconds: 2),
                    behavior: SnackBarBehavior.floating,
                    width: 200,
                  ),
                );
              },
              child: Padding(
                padding: const EdgeInsets.all(4),
                child: Icon(
                  Icons.copy_rounded,
                  size: 14,
                  color: theme.colorScheme.onSurfaceVariant,
                ),
              ),
            ),
          ),
        ],
      ],
    );
  }
}
