import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/widgets/entity_list_page.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/settings_providers.dart';
import '../widgets/setting_tile.dart';
import '../widgets/settings_scope_selector.dart';

/// Screen that displays settings grouped by module with a search bar
/// and scope selector.
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

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    // Choose between search results and full list.
    final asyncSettings = _searchQuery.isNotEmpty
        ? ref.watch(settingsSearchProvider(_searchQuery))
        : ref.watch(settingsListProvider(_listParams));

    return asyncSettings.when(
      loading: () => _buildShell(
        theme,
        isLoading: true,
        items: const [],
      ),
      error: (error, _) => _buildShell(
        theme,
        error: friendlyError(error),
        items: const [],
      ),
      data: (settings) => _buildShell(
        theme,
        items: settings,
      ),
    );
  }

  Widget _buildShell(
    ThemeData theme, {
    required List<SettingObject> items,
    bool isLoading = false,
    String? error,
  }) {
    // Group settings by module for display.
    final grouped = <String, List<SettingObject>>{};
    for (final s in items) {
      // Use the setting's key prefix or a default group label.
      final module = _extractModule(s.key);
      grouped.putIfAbsent(module, () => []).add(s);
    }
    final sortedModules = grouped.keys.toList()..sort();

    // Flatten into a list with section headers for EntityListPage.
    final flatItems = <_ListItem>[];
    for (final module in sortedModules) {
      flatItems.add(_ListItem.header(module));
      for (final s in grouped[module]!) {
        flatItems.add(_ListItem.setting(s));
      }
    }

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

        // Main list
        Expanded(
          child: EntityListPage<_ListItem>(
            title: 'Settings',
            icon: Icons.settings,
            items: flatItems,
            isLoading: isLoading,
            error: error,
            onRetry: () => _refresh(),
            searchHint: 'Search settings by key...',
            onSearchChanged: (query) {
              setState(() => _searchQuery = query.trim());
            },
            actionLabel: 'Bulk Edit',
            onAction: () => context.go('/settings/bulk-edit'),
            itemBuilder: (context, item) {
              if (item.isHeader) {
                return Padding(
                  padding: const EdgeInsets.only(top: 16, bottom: 4),
                  child: Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 4,
                        ),
                        decoration: BoxDecoration(
                          color: theme.colorScheme.secondaryContainer,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          item.module,
                          style: theme.textTheme.labelMedium?.copyWith(
                            fontWeight: FontWeight.w600,
                            color: theme.colorScheme.onSecondaryContainer,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Divider(
                          color: theme.colorScheme.outlineVariant,
                        ),
                      ),
                    ],
                  ),
                );
              }
              return SettingTile(
                setting: item.setting!,
                onTap: () {
                  context.go(
                    '/settings/detail/${Uri.encodeComponent(item.setting!.key)}',
                    extra: item.setting,
                  );
                },
              );
            },
          ),
        ),
      ],
    );
  }

  void _refresh() {
    if (_searchQuery.isNotEmpty) {
      ref.invalidate(settingsSearchProvider(_searchQuery));
    } else {
      ref.invalidate(settingsListProvider(_listParams));
    }
  }

  /// Extracts a module/group name from the setting key.
  /// E.g., "auth.session.timeout" -> "auth",
  ///       "notifications_email_enabled" -> "notifications".
  String _extractModule(String key) {
    if (key.contains('.')) return key.split('.').first;
    if (key.contains('_')) return key.split('_').first;
    return 'general';
  }
}

/// Internal helper to represent either a section header or a setting row.
class _ListItem {
  _ListItem.header(this.module) : setting = null;
  _ListItem.setting(SettingObject s)
      : setting = s,
        module = '';

  final SettingObject? setting;
  final String module;

  bool get isHeader => setting == null;
}
