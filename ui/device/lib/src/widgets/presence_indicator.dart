import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter/material.dart';

/// A small coloured dot indicating device presence status.
///
/// - ONLINE: green
/// - OFFLINE: grey
/// - AWAY: amber
/// - BUSY: red
/// - INVISIBLE: purple
class PresenceIndicator extends StatelessWidget {
  const PresenceIndicator({
    super.key,
    required this.status,
    this.size = 12,
    this.showLabel = false,
  });

  final PresenceStatus status;
  final double size;
  final bool showLabel;

  Color _color() {
    switch (status) {
      case PresenceStatus.ONLINE:
        return Colors.green;
      case PresenceStatus.OFFLINE:
        return Colors.grey;
      case PresenceStatus.AWAY:
        return Colors.amber;
      case PresenceStatus.BUSY:
        return Colors.red;
      case PresenceStatus.INVISIBLE:
        return Colors.purple;
      default:
        return Colors.grey;
    }
  }

  String _label() {
    switch (status) {
      case PresenceStatus.ONLINE:
        return 'Online';
      case PresenceStatus.OFFLINE:
        return 'Offline';
      case PresenceStatus.AWAY:
        return 'Away';
      case PresenceStatus.BUSY:
        return 'Busy';
      case PresenceStatus.INVISIBLE:
        return 'Invisible';
      default:
        return 'Unknown';
    }
  }

  @override
  Widget build(BuildContext context) {
    final dot = Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: _color(),
        shape: BoxShape.circle,
        border: Border.all(
          color: Theme.of(context).colorScheme.surface,
          width: size > 10 ? 2 : 1.5,
        ),
      ),
    );

    if (!showLabel) return dot;

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        dot,
        const SizedBox(width: 6),
        Text(
          _label(),
          style: Theme.of(context).textTheme.labelSmall?.copyWith(
                color: _color(),
                fontWeight: FontWeight.w600,
              ),
        ),
      ],
    );
  }
}
