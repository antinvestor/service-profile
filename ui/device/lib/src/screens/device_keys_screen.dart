import 'dart:convert';

import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/device_key_providers.dart';

/// Screen displaying keys associated with a device as a paginated DataTable.
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

  String _truncateKey(KeyObject keyObj) {
    try {
      final decoded = utf8.decode(keyObj.key, allowMalformed: true);
      if (decoded.length <= 20) return decoded;
      return '${decoded.substring(0, 17)}...';
    } catch (_) {
      final hex = keyObj.key
          .take(10)
          .map((b) => b.toRadixString(16).padLeft(2, '0'))
          .join();
      return '$hex...';
    }
  }

  @override
  Widget build(BuildContext context) {
    final asyncKeys = ref.watch(deviceKeysProvider(widget.deviceId));

    return asyncKeys.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline, size: 48,
                color: Theme.of(context).colorScheme.error),
            const SizedBox(height: 16),
            Text(friendlyError(error)),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () =>
                  ref.invalidate(deviceKeysProvider(widget.deviceId)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (keys) => _buildTable(keys),
    );
  }

  Widget _buildTable(List<KeyObject> keys) {
    return AdminEntityListPage<KeyObject>(
      title: 'Device Keys',
      breadcrumbs: const ['Home', 'Devices', 'Keys'],
      columns: const [
        DataColumn(label: Text('Type')),
        DataColumn(label: Text('Key')),
        DataColumn(label: Text('Created')),
        DataColumn(label: Text('Expires')),
        DataColumn(label: Text('Active')),
      ],
      items: keys,
      searchHint: 'Search keys...',
      onAdd: () => _showAddKeyDialog(context),
      addLabel: 'Add Key',
      rowBuilder: (keyObj, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(_keyTypeLabel(keyObj.keyType))),
            DataCell(Text(
              _truncateKey(keyObj),
              style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
            )),
            DataCell(Text(keyObj.createdAt)),
            DataCell(Text(
              keyObj.expiresAt.isNotEmpty ? keyObj.expiresAt : '-',
            )),
            DataCell(Icon(
              keyObj.isActive ? Icons.check_circle : Icons.cancel,
              color: keyObj.isActive ? Colors.green : Colors.red,
              size: 18,
            )),
          ],
        );
      },
      exportRow: (keyObj) => [
        _keyTypeLabel(keyObj.keyType),
        keyObj.id,
        keyObj.createdAt,
        keyObj.expiresAt,
        keyObj.isActive ? 'Active' : 'Inactive',
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Device Keys rows as $format'),
      actions: [
        IconButton(
          icon: const Icon(Icons.refresh, size: 20),
          tooltip: 'Refresh',
          onPressed: () =>
              ref.invalidate(deviceKeysProvider(widget.deviceId)),
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
                initialValue: selectedType,
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
        if (context.mounted) {
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
}
