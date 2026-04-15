import 'dart:async';

import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/area_providers.dart';
import 'area_badge.dart';
import 'area_type_badge.dart';

/// A form field for picking a geographic area from the geolocation service.
///
/// Shows the currently selected area (resolved by ID) with a button to open
/// a search dialog. The dialog provides debounced search with area type badges.
///
/// ```dart
/// AreaPickerField(
///   selectedAreaId: _geoId,
///   label: 'Coverage Area',
///   onSelected: (area) => setState(() => _geoId = area?.id ?? ''),
/// )
/// ```
class AreaPickerField extends ConsumerWidget {
  const AreaPickerField({
    super.key,
    required this.selectedAreaId,
    required this.onSelected,
    this.label = 'Coverage Area',
    this.description,
    this.errorText,
    this.enabled = true,
    this.isRequired = false,
  });

  /// Currently selected area ID (may be empty for no selection).
  final String selectedAreaId;

  /// Called when the user picks an area (or clears the selection).
  final ValueChanged<AreaObject?> onSelected;

  /// Label displayed above the field.
  final String label;

  /// Optional description text below the label.
  final String? description;

  /// Error text shown below the field (e.g. for validation).
  final String? errorText;

  /// Whether the field is interactive.
  final bool enabled;

  /// Whether this field is required (shows asterisk).
  final bool isRequired;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Label
        Padding(
          padding: const EdgeInsets.only(bottom: 4),
          child: Row(
            children: [
              Text(
                label,
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              if (isRequired)
                Text(' *',
                    style: TextStyle(
                        color: theme.colorScheme.error,
                        fontWeight: FontWeight.w600)),
            ],
          ),
        ),
        if (description != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Text(
              description!,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ),

        // Selected area display + pick button
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            border: Border.all(
              color: errorText != null
                  ? theme.colorScheme.error
                  : theme.colorScheme.outline.withAlpha(100),
            ),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              Icon(
                Icons.place_outlined,
                size: 20,
                color: selectedAreaId.isNotEmpty
                    ? theme.colorScheme.primary
                    : theme.colorScheme.onSurfaceVariant,
              ),
              const SizedBox(width: 8),
              Expanded(
                child: selectedAreaId.isNotEmpty
                    ? AreaBadge(areaId: selectedAreaId)
                    : Text(
                        'No area selected',
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
              ),
              if (selectedAreaId.isNotEmpty && enabled)
                IconButton(
                  icon: const Icon(Icons.close, size: 18),
                  tooltip: 'Clear selection',
                  onPressed: () => onSelected(null),
                  visualDensity: VisualDensity.compact,
                ),
              if (enabled)
                FilledButton.tonal(
                  onPressed: () => _openPicker(context, ref),
                  child: const Text('Select'),
                ),
            ],
          ),
        ),

        // Error text
        if (errorText != null)
          Padding(
            padding: const EdgeInsets.only(top: 4, left: 12),
            child: Text(
              errorText!,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.error,
              ),
            ),
          ),

        const SizedBox(height: 12),
      ],
    );
  }

  Future<void> _openPicker(BuildContext context, WidgetRef ref) async {
    final result = await showDialog<AreaObject>(
      context: context,
      builder: (context) => _AreaSearchDialog(
        initialSelectedId: selectedAreaId,
      ),
    );
    if (result != null) {
      onSelected(result);
    }
  }
}

// ---------------------------------------------------------------------------
// Area search dialog
// ---------------------------------------------------------------------------

class _AreaSearchDialog extends ConsumerStatefulWidget {
  const _AreaSearchDialog({this.initialSelectedId = ''});
  final String initialSelectedId;

  @override
  ConsumerState<_AreaSearchDialog> createState() => _AreaSearchDialogState();
}

class _AreaSearchDialogState extends ConsumerState<_AreaSearchDialog> {
  final _searchCtrl = TextEditingController();
  Timer? _debounce;
  String _query = '';

  @override
  void dispose() {
    _searchCtrl.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (mounted) setState(() => _query = value.trim());
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final areasAsync =
        _query.length >= 2 ? ref.watch(searchAreasProvider(_query)) : null;

    return Dialog(
      clipBehavior: Clip.antiAlias,
      insetPadding: const EdgeInsets.symmetric(horizontal: 24, vertical: 24),
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 640, maxHeight: 520),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Header
            Container(
              padding: const EdgeInsets.fromLTRB(24, 20, 24, 12),
              decoration: BoxDecoration(
                border: Border(
                  bottom: BorderSide(
                      color: theme.colorScheme.outlineVariant.withAlpha(60)),
                ),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(Icons.place, color: theme.colorScheme.primary),
                      const SizedBox(width: 8),
                      Text(
                        'Select Area',
                        style: theme.textTheme.titleLarge
                            ?.copyWith(fontWeight: FontWeight.w700),
                      ),
                      const Spacer(),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(context),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  TextField(
                    controller: _searchCtrl,
                    autofocus: true,
                    decoration: InputDecoration(
                      hintText: 'Search areas by name...',
                      prefixIcon: const Icon(Icons.search, size: 20),
                      isDense: true,
                      border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8)),
                    ),
                    onChanged: _onSearchChanged,
                  ),
                ],
              ),
            ),

            // Results
            Expanded(
              child: _buildResults(theme, areasAsync),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildResults(
      ThemeData theme, AsyncValue<List<AreaObject>>? areasAsync) {
    if (_query.length < 2) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.search,
                  size: 48,
                  color: theme.colorScheme.onSurface.withAlpha(60)),
              const SizedBox(height: 12),
              Text(
                'Type at least 2 characters to search',
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                ),
              ),
            ],
          ),
        ),
      );
    }

    if (areasAsync == null) {
      return const Center(child: CircularProgressIndicator());
    }

    return areasAsync.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 12),
            Text('Failed to search: $e',
                style: theme.textTheme.bodyMedium),
          ],
        ),
      ),
      data: (areas) {
        if (areas.isEmpty) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.location_off,
                    size: 48,
                    color: theme.colorScheme.onSurface.withAlpha(60)),
                const SizedBox(height: 12),
                Text('No areas found for "$_query"',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    )),
              ],
            ),
          );
        }

        return ListView.builder(
          padding: const EdgeInsets.symmetric(vertical: 8),
          itemCount: areas.length,
          itemBuilder: (context, index) {
            final area = areas[index];
            final isSelected = area.id == widget.initialSelectedId;

            return ListTile(
              leading: Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: isSelected
                      ? theme.colorScheme.primaryContainer
                      : theme.colorScheme.surfaceContainerHighest,
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(
                  Icons.place_outlined,
                  size: 20,
                  color: isSelected
                      ? theme.colorScheme.primary
                      : theme.colorScheme.onSurfaceVariant,
                ),
              ),
              title: Text(
                area.name,
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              subtitle: area.description.isNotEmpty
                  ? Text(
                      area.description,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    )
                  : null,
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  AreaTypeBadge(areaType: area.areaType),
                  if (isSelected) ...[
                    const SizedBox(width: 8),
                    Icon(Icons.check_circle,
                        size: 20, color: theme.colorScheme.primary),
                  ],
                ],
              ),
              selected: isSelected,
              selectedTileColor:
                  theme.colorScheme.primaryContainer.withAlpha(40),
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8)),
              onTap: () => Navigator.pop(context, area),
            );
          },
        );
      },
    );
  }
}
