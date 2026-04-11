import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter/material.dart';

/// A tile displaying a device key with type badge, ID preview,
/// active status, and expiry info.
class DeviceKeyTile extends StatelessWidget {
  const DeviceKeyTile({
    super.key,
    required this.keyObject,
    this.onRemove,
  });

  final KeyObject keyObject;
  final VoidCallback? onRemove;

  String _keyTypeLabel(KeyType type) {
    switch (type) {
      case KeyType.MATRIX_KEY:
        return 'Matrix';
      case KeyType.NOTIFICATION_KEY:
        return 'Notification';
      case KeyType.FCM_TOKEN:
        return 'FCM';
      case KeyType.CURVE25519_KEY:
        return 'Curve25519';
      case KeyType.ED25519_KEY:
        return 'Ed25519';
      case KeyType.PICKLE_KEY:
        return 'Pickle';
      default:
        return 'Unknown';
    }
  }

  IconData _keyTypeIcon(KeyType type) {
    switch (type) {
      case KeyType.FCM_TOKEN:
        return Icons.notifications_active;
      case KeyType.NOTIFICATION_KEY:
        return Icons.notifications;
      case KeyType.MATRIX_KEY:
        return Icons.lock;
      case KeyType.CURVE25519_KEY:
      case KeyType.ED25519_KEY:
        return Icons.vpn_key;
      case KeyType.PICKLE_KEY:
        return Icons.enhanced_encryption;
      default:
        return Icons.key;
    }
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
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: keyObject.isActive
                    ? theme.colorScheme.primaryContainer
                    : theme.colorScheme.surfaceContainerHighest,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                _keyTypeIcon(keyObject.keyType),
                size: 20,
                color: keyObject.isActive
                    ? theme.colorScheme.onPrimaryContainer
                    : theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: theme.colorScheme.secondaryContainer,
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: Text(
                          _keyTypeLabel(keyObject.keyType),
                          style: theme.textTheme.labelSmall?.copyWith(
                            fontWeight: FontWeight.w600,
                            color: theme.colorScheme.onSecondaryContainer,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      if (!keyObject.isActive)
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 6,
                            vertical: 2,
                          ),
                          decoration: BoxDecoration(
                            color: theme.colorScheme.errorContainer,
                            borderRadius: BorderRadius.circular(6),
                          ),
                          child: Text(
                            'Inactive',
                            style: theme.textTheme.labelSmall?.copyWith(
                              color: theme.colorScheme.onErrorContainer,
                            ),
                          ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'ID: ${keyObject.id}',
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                      fontFamily: 'monospace',
                      fontSize: 12,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  if (keyObject.expiresAt.isNotEmpty) ...[
                    const SizedBox(height: 2),
                    Text(
                      'Expires: ${keyObject.expiresAt}',
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ],
              ),
            ),
            if (onRemove != null)
              IconButton(
                icon: Icon(
                  Icons.delete_outline,
                  size: 20,
                  color: theme.colorScheme.error,
                ),
                tooltip: 'Remove key',
                onPressed: onRemove,
              ),
          ],
        ),
      ),
    );
  }
}
