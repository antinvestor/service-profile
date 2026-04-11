import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/settings_providers.dart';
import '../widgets/setting_value_editor.dart';

/// Screen for viewing and editing a single setting's value.
class SettingDetailScreen extends ConsumerStatefulWidget {
  const SettingDetailScreen({
    super.key,
    required this.settingKey,
    this.initialSetting,
  });

  /// The setting key (used in the URL).
  final String settingKey;

  /// Optional pre-loaded setting passed via route extra.
  final SettingObject? initialSetting;

  @override
  ConsumerState<SettingDetailScreen> createState() =>
      _SettingDetailScreenState();
}

class _SettingDetailScreenState extends ConsumerState<SettingDetailScreen> {
  bool _editing = false;
  late String _editedValue;
  bool _saving = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _editedValue = widget.initialSetting?.value ?? '';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final setting = widget.initialSetting;

    if (setting == null) {
      return _buildNotFound(theme);
    }

    final mutationState = ref.watch(settingsNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.canPop() ? context.pop() : context.go('/settings'),
        ),
        title: Text(
          setting.key.name,
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
                _editedValue = setting.value;
              }),
            ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Metadata section
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
                    _metadataRow(theme, 'Key', setting.key.name),
                    _metadataRow(theme, 'ID', setting.id),
                    if (setting.hasUpdated())
                      _metadataRow(
                        theme,
                        'Last Updated',
                        _formatTimestamp(setting),
                      ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            // Value editor section
            FormFieldCard(
              label: 'Value',
              description: _editing
                  ? 'Edit the setting value below. Choose the appropriate type.'
                  : 'Current setting value.',
              child: SettingValueEditor(
                value: _editing ? _editedValue : setting.value,
                readOnly: !_editing,
                onChanged: (value) => _editedValue = value,
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
                              _editedValue = setting.value;
                            }),
                    child: const Text('Cancel'),
                  ),
                  const SizedBox(width: 12),
                  FilledButton.icon(
                    onPressed:
                        _saving || mutationState.isLoading ? null : _save,
                    icon: _saving
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child:
                                CircularProgressIndicator(strokeWidth: 2),
                          )
                        : const Icon(Icons.save, size: 18),
                    label: Text(_saving ? 'Saving...' : 'Save'),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _metadataRow(ThemeData theme, String label, String value) {
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

  String _formatTimestamp(SettingObject s) {
    if (!s.hasUpdated()) return 'N/A';
    final ts = DateTime.tryParse(s.updated);
    if (ts == null) return s.updated;
    return '${ts.year}-${ts.month.toString().padLeft(2, '0')}-'
        '${ts.day.toString().padLeft(2, '0')} '
        '${ts.hour.toString().padLeft(2, '0')}:'
        '${ts.minute.toString().padLeft(2, '0')}';
  }

  Future<void> _save() async {
    setState(() {
      _saving = true;
      _error = null;
    });

    try {
      final setting = widget.initialSetting!;
      final originalKey = setting.key;
      final request = SetRequest()
        ..key = (Setting()
          ..name = originalKey.name
          ..object = originalKey.object
          ..objectId = originalKey.objectId
          ..lang = originalKey.lang
          ..module = originalKey.module)
        ..value = _editedValue;

      final notifier = ref.read(settingsNotifierProvider.notifier);
      await notifier.set(request);

      if (mounted) {
        setState(() {
          _editing = false;
          _saving = false;
        });
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Setting saved successfully'),
            behavior: SnackBarBehavior.floating,
          ),
        );
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

  Widget _buildNotFound(ThemeData theme) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.canPop() ? context.pop() : context.go('/settings'),
        ),
        title: const Text('Setting Not Found'),
      ),
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.search_off,
              size: 48,
              color: theme.colorScheme.onSurfaceVariant,
            ),
            const SizedBox(height: 16),
            Text(
              'Setting "${widget.settingKey}" was not found.',
              style: theme.textTheme.titleMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () => context.go('/settings'),
              child: const Text('Back to Settings'),
            ),
          ],
        ),
      ),
    );
  }
}
