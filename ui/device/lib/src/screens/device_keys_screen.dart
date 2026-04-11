import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/device_key_providers.dart';
import '../widgets/device_key_tile.dart';

/// Screen displaying keys associated with a device.
/// Supports adding and removing keys.
class DeviceKeysScreen extends ConsumerStatefulWidget {
  const DeviceKeysScreen({
    super.key,
    required this.deviceId,
  });

  final String deviceId;

  @override
  ConsumerState<DeviceKeysScreen> createState() => _DeviceKeysScreenState();
}

class _DeviceKeysScreenState extends ConsumerState<DeviceKeysScreen> {
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final asyncKeys = ref.watch(deviceKeysProvider(widget.deviceId));

    return Column(
      children: [
        // Header with add button
        Padding(
          padding: const EdgeInsets.fromLTRB(24, 16, 24, 8),
          child: Row(
            children: [
              Text(
                'Device Keys',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const Spacer(),
              FilledButton.tonalIcon(
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add Key'),
                onPressed: () => _showAddKeyDialog(context),
              ),
            ],
          ),
        ),

        // Keys list
        Expanded(
          child: asyncKeys.when(
            loading: () =>
                const Center(child: CircularProgressIndicator()),
            error: (error, _) => Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 48,
                    color: theme.colorScheme.error,
                  ),
                  const SizedBox(height: 16),
                  Text(friendlyError(error)),
                  const SizedBox(height: 12),
                  FilledButton.tonal(
                    onPressed: () => ref.invalidate(
                      deviceKeysProvider(widget.deviceId),
                    ),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            ),
            data: (keys) {
              if (keys.isEmpty) {
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
                          Icons.vpn_key_off,
                          size: 28,
                          color: theme.colorScheme.primary.withAlpha(160),
                        ),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'No keys registered for this device.',
                        style: theme.textTheme.titleMedium?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ],
                  ),
                );
              }

              return ListView.separated(
                padding: const EdgeInsets.fromLTRB(24, 8, 24, 24),
                itemCount: keys.length,
                separatorBuilder: (_, __) => const SizedBox(height: 8),
                itemBuilder: (context, index) {
                  final keyObj = keys[index];
                  return DeviceKeyTile(
                    keyObject: keyObj,
                    onRemove: () => _confirmRemoveKey(context, keyObj),
                  );
                },
              );
            },
          ),
        ),
      ],
    );
  }

  Future<void> _showAddKeyDialog(BuildContext context) async {
    KeyType selectedType = KeyType.FCM_TOKEN;
    final keyValueController = TextEditingController();

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: const Text('Add Key'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              DropdownButtonFormField<KeyType>(
                value: selectedType,
                decoration: InputDecoration(
                  labelText: 'Key Type',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                items: KeyType.values.map((type) {
                  return DropdownMenuItem(
                    value: type,
                    child: Text(type.name),
                  );
                }).toList(),
                onChanged: (value) {
                  if (value != null) {
                    setDialogState(() => selectedType = value);
                  }
                },
              ),
              const SizedBox(height: 16),
              TextField(
                controller: keyValueController,
                decoration: InputDecoration(
                  labelText: 'Key Value',
                  hintText: 'Enter the key value (e.g. FCM token)',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                maxLines: 3,
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text('Cancel'),
            ),
            FilledButton(
              onPressed: () => Navigator.pop(context, true),
              child: const Text('Add'),
            ),
          ],
        ),
      ),
    );

    if (confirmed == true && mounted) {
      try {
        final notifier = ref.read(deviceKeyNotifierProvider.notifier);
        await notifier.addKey(
          AddKeyRequest()
            ..deviceId = widget.deviceId
            ..keyType = selectedType
            ..data = keyValueController.text.trim().codeUnits,
        );
        ref.invalidate(deviceKeysProvider(widget.deviceId));
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(friendlyError(e)),
              behavior: SnackBarBehavior.floating,
            ),
          );
        }
      }
    }
    keyValueController.dispose();
  }

  Future<void> _confirmRemoveKey(
    BuildContext context,
    KeyObject keyObj,
  ) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Remove Key'),
        content: Text(
          'Remove ${keyObj.keyType.name} key "${keyObj.id}"? '
          'This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(context).colorScheme.error,
            ),
            child: const Text('Remove'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      try {
        final notifier = ref.read(deviceKeyNotifierProvider.notifier);
        await notifier.removeKey(
          RemoveKeyRequest()..id.add(keyObj.id),
        );
        ref.invalidate(deviceKeysProvider(widget.deviceId));
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(friendlyError(e)),
              behavior: SnackBarBehavior.floating,
            ),
          );
        }
      }
    }
  }
}
