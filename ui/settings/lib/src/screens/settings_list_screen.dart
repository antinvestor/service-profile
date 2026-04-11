import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/widgets/admin_entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/settings_providers.dart';
import '../widgets/settings_scope_selector.dart';

/// Screen that displays settings in a paginated DataTable with search,
/// scope selector, and CSV export.
class SettingsListScreen extends ConsumerStatefulWidget {
  const SettingsListScreen({super.key});

  @override
  ConsumerState<SettingsListScreen> createState() => _SettingsListScreenState();
}

class _SettingsListScreenState extends ConsumerState<SettingsListScreen> {
  SettingsScope _scope = const SettingsScope(
    object: 'application',
    objectId: 'default',
  );
  String _searchQuery = '';
  bool _scopeExpanded = true;

  SettingListParams get _listParams => SettingListParams(
        object: _scope.object,
        objectId: _scope.objectId,
      );

  /// Extracts a module/group name from the setting key.
  String _extractModule(String key) {
    if (key.contains('.')) return key.split('.').first;
    if (key.contains('_')) return key.split('_').first;
    return 'general';
  }

  String _valuePreview(String value) {
    if (value.length <= 60) return value;
    return '${value.substring(0, 57)}...';
  }

  String _formatUpdated(SettingObject s) {
    if (!s.hasUpdated()) return '-';
    return s.updated;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    // Choose between search results and full list.
    final asyncSettings = _searchQuery.isNotEmpty
        ? ref.watch(settingsSearchProvider(_searchQuery))
        : ref.watch(settingsListProvider(_listParams));

    return asyncSettings.when(
      loading: () => _buildShell(theme, isLoading: true, items: const []),
      error: (error, _) =>
          _buildShell(theme, error: friendlyError(error), items: const []),
      data: (settings) => _buildShell(theme, items: settings),
    );
  }

  Widget _buildShell(
    ThemeData theme, {
    required List<SettingObject> items,
    bool isLoading = false,
    String? error,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // Scope selector
        Padding(
          padding: const EdgeInsets.fromLTRB(24, 16, 24, 0),
          child: Column(
            children: [
              InkWell(
                onTap: () => setState(() => _scopeExpanded = !_scopeExpanded),
                borderRadius: BorderRadius.circular(8),
                child: Row(
                  children: [
                    Icon(
                      _scopeExpanded
                          ? Icons.expand_less
                          : Icons.expand_more,
                      size: 20,
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      'Scope Filter',
                      style: theme.textTheme.labelMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ),
              ),
              if (_scopeExpanded) ...[
                const SizedBox(height: 8),
                SettingsScopeSelector(
                  scope: _scope,
                  onScopeChanged: (scope) {
                    setState(() => _scope = scope);
                  },
                ),
              ],
            ],
          ),
        ),

        // Main table
        Expanded(
          child: _buildContent(items, isLoading, error),
        ),
      ],
    );
  }

  Widget _buildContent(
    List<SettingObject> items,
    bool isLoading,
    String? error,
  ) {
    if (isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(error),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: _refresh,
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    return AdminEntityListPage<SettingObject>(
      title: 'Settings',
      breadcrumbs: const ['Home', 'Settings'],
      columns: const [
        DataColumn(label: Text('Key')),
        DataColumn(label: Text('Value')),
        DataColumn(label: Text('Module')),
        DataColumn(label: Text('Updated')),
      ],
      items: items,
      searchHint: 'Search settings by key...',
      onSearch: (query) {
        setState(() => _searchQuery = query.trim());
      },
      onAdd: () => context.go('/settings/bulk-edit'),
      addLabel: 'Bulk Edit',
      onRowNavigate: (setting) {
        context.go(
          '/settings/detail/${Uri.encodeComponent(setting.key.name)}',
          extra: setting,
        );
      },
      rowBuilder: (setting, selected, onSelect) {
        return DataRow(
          selected: selected,
          onSelectChanged: (_) => onSelect(),
          cells: [
            DataCell(Text(
              setting.key.name,
              style: const TextStyle(fontWeight: FontWeight.w600),
            )),
            DataCell(ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 250),
              child: Text(
                _valuePreview(setting.value),
                overflow: TextOverflow.ellipsis,
                maxLines: 1,
              ),
            )),
            DataCell(Text(_extractModule(setting.key.name))),
            DataCell(Text(_formatUpdated(setting))),
          ],
        );
      },
      exportRow: (setting) => [
        setting.key.name,
        setting.value,
        _extractModule(setting.key.name),
        _formatUpdated(setting),
      ],
      onExport: (format, count) =>
          debugPrint('Exported $count Settings rows as $format'),
    );
  }

  void _refresh() {
    if (_searchQuery.isNotEmpty) {
      ref.invalidate(settingsSearchProvider(_searchQuery));
    } else {
      ref.invalidate(settingsListProvider(_listParams));
    }
  }
}
