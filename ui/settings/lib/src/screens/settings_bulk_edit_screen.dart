import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/settings_providers.dart';
import '../widgets/settings_scope_selector.dart';

/// Screen for editing multiple settings at once.
///
/// Shows a table-like view with key, current value, and new value columns.
/// Users can modify multiple values and save them all in one batch.
class SettingsBulkEditScreen extends ConsumerStatefulWidget {
  const SettingsBulkEditScreen({super.key});

  @override
  ConsumerState<SettingsBulkEditScreen> createState() =>
      _SettingsBulkEditScreenState();
}

class _SettingsBulkEditScreenState
    extends ConsumerState<SettingsBulkEditScreen> {
  SettingsScope _scope = const SettingsScope(
    object: 'application',
    objectId: 'default',
  );

  /// Map from setting key to the edited value. Only populated when the user
  /// has changed a value.
  final Map<String, String> _edits = {};

  /// Persistent controllers keyed by setting key name, to avoid creating
  /// new controllers on every rebuild which causes cursor jumps and state loss.
  final Map<String, TextEditingController> _controllers = {};

  bool _saving = false;
  String? _error;
  int _savedCount = 0;

  @override
  void dispose() {
    for (final controller in _controllers.values) {
      controller.dispose();
    }
    super.dispose();
  }

  TextEditingController _controllerFor(String key, String initialValue) {
    return _controllers.putIfAbsent(key, () {
      final controller = TextEditingController(text: initialValue);
      return controller;
    });
  }

  void _clearControllers() {
    for (final controller in _controllers.values) {
      controller.dispose();
    }
    _controllers.clear();
  }

  SettingListParams get _listParams => SettingListParams(
        object: _scope.object,
        objectId: _scope.objectId,
      );

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final asyncSettings = ref.watch(settingsListProvider(_listParams));

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () =>
              context.canPop() ? context.pop() : context.go('/settings'),
        ),
        title: Text(
          'Bulk Edit Settings',
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        actions: [
          if (_edits.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(right: 8),
              child: Badge(
                label: Text('${_edits.length}'),
                child: IconButton(
                  icon: const Icon(Icons.undo),
                  tooltip: 'Reset all changes',
                  onPressed: () => setState(() {
                    _edits.clear();
                    _clearControllers();
                  }),
                ),
              ),
            ),
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: FilledButton.icon(
              onPressed: _edits.isEmpty || _saving ? null : _saveAll,
              icon: _saving
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : const Icon(Icons.save, size: 18),
              label: Text(
                _saving
                    ? 'Saving...'
                    : 'Save ${_edits.length} change${_edits.length == 1 ? '' : 's'}',
              ),
            ),
          ),
        ],
      ),
      body: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Scope selector
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 16, 24, 8),
            child: SettingsScopeSelector(
              scope: _scope,
              onScopeChanged: (scope) {
                setState(() {
                  _scope = scope;
                  _edits.clear();
                  _clearControllers();
                  _error = null;
                });
              },
            ),
          ),

          // Error / success banners
          if (_error != null)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 4),
              child: _buildBanner(
                theme,
                icon: Icons.error_outline,
                message: _error!,
                color: theme.colorScheme.errorContainer,
                textColor: theme.colorScheme.onErrorContainer,
              ),
            ),
          if (_savedCount > 0 && _error == null && !_saving)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 4),
              child: _buildBanner(
                theme,
                icon: Icons.check_circle_outline,
                message: '$_savedCount setting${_savedCount == 1 ? '' : 's'} saved successfully.',
                color: theme.colorScheme.primaryContainer,
                textColor: theme.colorScheme.onPrimaryContainer,
              ),
            ),

          // Settings table
          Expanded(
            child: asyncSettings.when(
              loading: () =>
                  const Center(child: CircularProgressIndicator()),
              error: (error, _) => Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.error_outline,
                        size: 48, color: theme.colorScheme.error),
                    const SizedBox(height: 16),
                    Text(friendlyError(error)),
                    const SizedBox(height: 12),
                    FilledButton.tonal(
                      onPressed: () =>
                          ref.invalidate(settingsListProvider(_listParams)),
                      child: const Text('Retry'),
                    ),
                  ],
                ),
              ),
              data: (settings) => _buildTable(theme, settings),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBanner(
    ThemeData theme, {
    required IconData icon,
    required String message,
    required Color color,
    required Color textColor,
  }) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: color,
        borderRadius: BorderRadius.circular(10),
      ),
      child: Row(
        children: [
          Icon(icon, size: 20, color: textColor),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              message,
              style: theme.textTheme.bodySmall?.copyWith(color: textColor),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTable(ThemeData theme, List<SettingObject> settings) {
    if (settings.isEmpty) {
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
                Icons.settings,
                size: 28,
                color: theme.colorScheme.primary.withAlpha(160),
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'No settings found for this scope.',
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
      itemCount: settings.length + 1, // +1 for header
      separatorBuilder: (_, _) => const SizedBox(height: 2),
      itemBuilder: (context, index) {
        if (index == 0) {
          return _buildTableHeader(theme);
        }
        final setting = settings[index - 1];
        final keyName = setting.key.name;
        final isEdited = _edits.containsKey(keyName);
        final controller = _controllerFor(keyName, setting.value);

        // Sync controller text when edits are reset externally
        if (!isEdited && controller.text != setting.value) {
          controller.text = setting.value;
        }

        return Container(
          decoration: BoxDecoration(
            color: isEdited
                ? theme.colorScheme.primaryContainer.withAlpha(40)
                : null,
            borderRadius: BorderRadius.circular(8),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          child: Row(
            children: [
              // Key column
              Expanded(
                flex: 3,
                child: Text(
                  keyName,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w500,
                    fontFamily: 'monospace',
                    fontSize: 13,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const SizedBox(width: 12),

              // Original value column
              Expanded(
                flex: 3,
                child: Text(
                  setting.value.length > 60
                      ? '${setting.value.substring(0, 57)}...'
                      : setting.value,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const SizedBox(width: 12),

              // Editable value column
              Expanded(
                flex: 4,
                child: TextField(
                  controller: controller,
                  style: theme.textTheme.bodySmall?.copyWith(
                    fontFamily: 'monospace',
                    fontSize: 12,
                  ),
                  decoration: InputDecoration(
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 10,
                      vertical: 8,
                    ),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                    suffixIcon: isEdited
                        ? IconButton(
                            icon: Icon(
                              Icons.undo,
                              size: 16,
                              color: theme.colorScheme.primary,
                            ),
                            tooltip: 'Reset to original',
                            onPressed: () {
                              setState(() {
                                _edits.remove(keyName);
                                controller.text = setting.value;
                              });
                            },
                          )
                        : null,
                  ),
                  onChanged: (value) {
                    setState(() {
                      if (value == setting.value) {
                        _edits.remove(keyName);
                      } else {
                        _edits[keyName] = value;
                      }
                    });
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildTableHeader(ThemeData theme) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Expanded(
            flex: 3,
            child: Text(
              'Key',
              style: theme.textTheme.labelMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            flex: 3,
            child: Text(
              'Current Value',
              style: theme.textTheme.labelMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            flex: 4,
            child: Text(
              'New Value',
              style: theme.textTheme.labelMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _saveAll() async {
    if (_edits.isEmpty) return;

    setState(() {
      _saving = true;
      _error = null;
      _savedCount = 0;
    });

    final notifier = ref.read(settingsNotifierProvider.notifier);
    var successCount = 0;
    final errors = <String>[];

    for (final entry in _edits.entries) {
      try {
        final request = SetRequest()
          ..key = (Setting()
            ..name = entry.key
            ..object = _scope.object
            ..objectId = _scope.objectId
            ..lang = _scope.lang)
          ..value = entry.value;
        await notifier.set(request);
        successCount++;
      } catch (e) {
        errors.add('${entry.key}: ${friendlyError(e)}');
      }
    }

    if (mounted) {
      setState(() {
        _saving = false;
        _savedCount = successCount;
        if (errors.isNotEmpty) {
          _error = 'Failed to save ${errors.length} setting(s):\n${errors.join('\n')}';
        } else {
          _edits.clear();
          _clearControllers();
        }
      });

      // Refresh the list to pick up new values.
      ref.invalidate(settingsListProvider(_listParams));
    }
  }
}
