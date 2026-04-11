import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/device_providers.dart';
import '../providers/presence_providers.dart';
import '../widgets/presence_indicator.dart';
import 'device_keys_screen.dart';
import 'device_logs_screen.dart';

/// Detail screen for a single device with tabbed layout:
/// Info, Keys, Sessions.
class DeviceDetailScreen extends ConsumerStatefulWidget {
  const DeviceDetailScreen({
    super.key,
    required this.deviceId,
    this.initialDevice,
  });

  final String deviceId;
  final DeviceObject? initialDevice;

  @override
  ConsumerState<DeviceDetailScreen> createState() => _DeviceDetailScreenState();
}

class _DeviceDetailScreenState extends ConsumerState<DeviceDetailScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  bool _editing = false;
  late TextEditingController _nameController;
  bool _saving = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _nameController = TextEditingController(
      text: widget.initialDevice?.name ?? '',
    );
  }

  @override
  void dispose() {
    _tabController.dispose();
    _nameController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final asyncDevice = widget.initialDevice != null
        ? AsyncValue.data(widget.initialDevice!)
        : ref.watch(deviceByIdProvider(widget.deviceId));

    return asyncDevice.when(
      loading: () => Scaffold(
        appBar: AppBar(title: const Text('Device')),
        body: const Center(child: CircularProgressIndicator()),
      ),
      error: (error, _) => Scaffold(
        appBar: AppBar(
          leading: _backButton(context),
          title: const Text('Device'),
        ),
        body: Center(child: Text(friendlyError(error))),
      ),
      data: (device) => _buildDetail(context, theme, device),
    );
  }

  Widget _buildDetail(
    BuildContext context,
    ThemeData theme,
    DeviceObject device,
  ) {
    return Scaffold(
      appBar: AppBar(
        leading: _backButton(context),
        title: Text(
          device.name.isNotEmpty ? device.name : 'Device',
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        actions: [
          if (!_editing)
            IconButton(
              icon: const Icon(Icons.edit),
              tooltip: 'Edit',
              onPressed: () => setState(() {
                _editing = true;
                _nameController.text = device.name;
              }),
            ),
          PopupMenuButton<String>(
            onSelected: (value) {
              if (value == 'remove') _confirmRemove(context, device);
            },
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'remove',
                child: Row(
                  children: [
                    Icon(Icons.delete_outline, size: 20),
                    SizedBox(width: 8),
                    Text('Remove Device'),
                  ],
                ),
              ),
            ],
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: 'Info'),
            Tab(text: 'Keys'),
            Tab(text: 'Sessions'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildInfoTab(theme, device),
          DeviceKeysScreen(deviceId: device.id),
          DeviceLogsScreen(deviceId: device.id),
        ],
      ),
    );
  }

  Widget _buildInfoTab(ThemeData theme, DeviceObject device) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Presence section
          Card(
            elevation: 0,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
              side: BorderSide(color: theme.colorScheme.outlineVariant),
            ),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  PresenceIndicator(
                    status: device.presence,
                    size: 14,
                    showLabel: true,
                  ),
                  const Spacer(),
                  OutlinedButton.icon(
                    icon: const Icon(Icons.sync, size: 16),
                    label: const Text('Update'),
                    onPressed: () => _showPresenceDialog(context, device),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Device info
          FormFieldCard(
            label: 'Device Name',
            description: _editing
                ? 'Edit the device display name.'
                : 'The display name for this device.',
            child: _editing
                ? TextField(
                    controller: _nameController,
                    decoration: InputDecoration(
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  )
                : Text(
                    device.name.isNotEmpty ? device.name : 'Unnamed Device',
                    style: theme.textTheme.bodyLarge,
                  ),
          ),
          const SizedBox(height: 16),

          // Metadata
          Card(
            elevation: 0,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
              side: BorderSide(color: theme.colorScheme.outlineVariant),
            ),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Details',
                    style: theme.textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.primary,
                    ),
                  ),
                  const SizedBox(height: 12),
                  _metadataRow(theme, 'ID', device.id),
                  _metadataRow(theme, 'OS', device.os),
                  _metadataRow(theme, 'IP Address', device.ip),
                  _metadataRow(theme, 'User Agent', device.userAgent),
                  _metadataRow(theme, 'Session ID', device.sessionId),
                  _metadataRow(theme, 'Profile ID', device.profileId),
                  _metadataRow(theme, 'Last Seen', device.lastSeen),
                  if (device.hasLocale()) ...[
                    _metadataRow(
                      theme,
                      'Timezone',
                      device.locale.timezone,
                    ),
                    _metadataRow(
                      theme,
                      'Currency',
                      device.locale.currency,
                    ),
                  ],
                ],
              ),
            ),
          ),

          // Error display
          if (_error != null) ...[
            const SizedBox(height: 8),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: theme.colorScheme.errorContainer,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Row(
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 20,
                    color: theme.colorScheme.onErrorContainer,
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      _error!,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onErrorContainer,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],

          // Action buttons
          if (_editing) ...[
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                OutlinedButton(
                  onPressed: _saving
                      ? null
                      : () => setState(() {
                            _editing = false;
                            _error = null;
                          }),
                  child: const Text('Cancel'),
                ),
                const SizedBox(width: 12),
                FilledButton.icon(
                  onPressed: _saving ? null : () => _save(device),
                  icon: _saving
                      ? const SizedBox(
                          width: 16,
                          height: 16,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Icon(Icons.save, size: 18),
                  label: Text(_saving ? 'Saving...' : 'Save'),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }

  Widget _metadataRow(ThemeData theme, String label, String value) {
    if (value.isEmpty) return const SizedBox.shrink();
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: theme.textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _backButton(BuildContext context) {
    return IconButton(
      icon: const Icon(Icons.arrow_back),
      onPressed: () =>
          context.canPop() ? context.pop() : context.go('/devices'),
    );
  }

  Future<void> _save(DeviceObject device) async {
    setState(() {
      _saving = true;
      _error = null;
    });

    try {
      final notifier = ref.read(deviceNotifierProvider.notifier);
      await notifier.update(
        UpdateRequest()
          ..id = device.id
          ..name = _nameController.text.trim(),
      );

      if (mounted) {
        setState(() {
          _editing = false;
          _saving = false;
        });
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Device updated successfully'),
            behavior: SnackBarBehavior.floating,
          ),
        );
        ref.invalidate(deviceByIdProvider(widget.deviceId));
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _saving = false;
          _error = friendlyError(e);
        });
      }
    }
  }

  Future<void> _confirmRemove(
    BuildContext context,
    DeviceObject device,
  ) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Remove Device'),
        content: Text(
          'Remove "${device.name.isNotEmpty ? device.name : device.id}"? '
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
        final notifier = ref.read(deviceNotifierProvider.notifier);
        await notifier.remove(RemoveRequest()..id = device.id);
        if (mounted) context.go('/devices');
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

  Future<void> _showPresenceDialog(
    BuildContext context,
    DeviceObject device,
  ) async {
    final selected = await showDialog<PresenceStatus>(
      context: context,
      builder: (context) => SimpleDialog(
        title: const Text('Update Presence'),
        children: PresenceStatus.values.map((status) {
          return SimpleDialogOption(
            onPressed: () => Navigator.pop(context, status),
            child: Row(
              children: [
                PresenceIndicator(status: status, size: 12),
                const SizedBox(width: 12),
                Text(status.name),
              ],
            ),
          );
        }).toList(),
      ),
    );

    if (selected != null && mounted) {
      try {
        final notifier = ref.read(presenceNotifierProvider.notifier);
        await notifier.updatePresence(
          UpdatePresenceRequest()
            ..deviceId = device.id
            ..status = selected,
        );
        ref.invalidate(deviceByIdProvider(widget.deviceId));
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
