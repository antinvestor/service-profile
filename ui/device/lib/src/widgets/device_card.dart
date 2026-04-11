import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter/material.dart';

import 'presence_indicator.dart';

/// A card widget displaying a device summary with presence status,
/// name, OS, IP, and last-seen time.
class DeviceCard extends StatelessWidget {
  const DeviceCard({
    super.key,
    required this.device,
    this.onTap,
  });

  final DeviceObject device;
  final VoidCallback? onTap;

  IconData _osIcon(String os) {
    final lower = os.toLowerCase();
    if (lower.contains('android')) return Icons.phone_android;
    if (lower.contains('ios') || lower.contains('iphone')) {
      return Icons.phone_iphone;
    }
    if (lower.contains('windows')) return Icons.desktop_windows;
    if (lower.contains('mac') || lower.contains('darwin')) {
      return Icons.laptop_mac;
    }
    if (lower.contains('linux')) return Icons.computer;
    return Icons.devices;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: theme.colorScheme.outlineVariant),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          child: Row(
            children: [
              // Device icon with presence badge
              Stack(
                children: [
                  Container(
                    width: 44,
                    height: 44,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      _osIcon(device.os),
                      size: 22,
                      color: theme.colorScheme.onPrimaryContainer,
                    ),
                  ),
                  Positioned(
                    right: 0,
                    bottom: 0,
                    child: PresenceIndicator(status: device.presence),
                  ),
                ],
              ),
              const SizedBox(width: 12),

              // Device info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      device.name.isNotEmpty ? device.name : 'Unnamed Device',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      [
                        if (device.os.isNotEmpty) device.os,
                        if (device.ip.isNotEmpty) device.ip,
                      ].join(' - '),
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),

              // Last seen
              if (device.lastSeen.isNotEmpty) ...[
                const SizedBox(width: 8),
                Text(
                  device.lastSeen,
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
              const SizedBox(width: 4),
              Icon(
                Icons.chevron_right,
                size: 20,
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
