import 'package:flutter/material.dart';

/// Scope parameters for filtering settings by object type, instance, and language.
class SettingsScope {
  const SettingsScope({
    this.object = '',
    this.objectId = '',
    this.lang = '',
  });

  final String object;
  final String objectId;
  final String lang;

  SettingsScope copyWith({String? object, String? objectId, String? lang}) {
    return SettingsScope(
      object: object ?? this.object,
      objectId: objectId ?? this.objectId,
      lang: lang ?? this.lang,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SettingsScope &&
          object == other.object &&
          objectId == other.objectId &&
          lang == other.lang;

  @override
  int get hashCode => Object.hash(object, objectId, lang);
}

/// A scope picker widget that lets users select the object type, object
/// instance ID, and language for settings filtering.
class SettingsScopeSelector extends StatefulWidget {
  const SettingsScopeSelector({
    super.key,
    required this.scope,
    required this.onScopeChanged,
    this.objectTypes = const ['application', 'tenant', 'user', 'device'],
    this.languages = const ['', 'en', 'fr', 'sw', 'es', 'de', 'pt'],
  });

  final SettingsScope scope;
  final ValueChanged<SettingsScope> onScopeChanged;
  final List<String> objectTypes;
  final List<String> languages;

  @override
  State<SettingsScopeSelector> createState() => _SettingsScopeSelectorState();
}

class _SettingsScopeSelectorState extends State<SettingsScopeSelector> {
  late TextEditingController _objectIdController;

  @override
  void initState() {
    super.initState();
    _objectIdController = TextEditingController(text: widget.scope.objectId);
  }

  @override
  void didUpdateWidget(SettingsScopeSelector oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.scope.objectId != widget.scope.objectId) {
      _objectIdController.text = widget.scope.objectId;
    }
  }

  @override
  void dispose() {
    _objectIdController.dispose();
    super.dispose();
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
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Scope',
              style: theme.textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.w600,
                color: theme.colorScheme.primary,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                // Object type dropdown
                Expanded(
                  flex: 2,
                  child: DropdownButtonFormField<String>(
                    value: widget.scope.object.isEmpty
                        ? null
                        : widget.scope.object,
                    decoration: InputDecoration(
                      labelText: 'Object Type',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 10,
                      ),
                    ),
                    items: [
                      const DropdownMenuItem(
                        value: null,
                        child: Text('All'),
                      ),
                      for (final type in widget.objectTypes)
                        DropdownMenuItem(
                          value: type,
                          child: Text(type),
                        ),
                    ],
                    onChanged: (value) {
                      widget.onScopeChanged(
                        widget.scope.copyWith(object: value ?? ''),
                      );
                    },
                  ),
                ),
                const SizedBox(width: 12),

                // Object ID text field
                Expanded(
                  flex: 3,
                  child: TextField(
                    controller: _objectIdController,
                    decoration: InputDecoration(
                      labelText: 'Instance ID',
                      hintText: 'e.g. tenant-123',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 10,
                      ),
                    ),
                    onSubmitted: (value) {
                      widget.onScopeChanged(
                        widget.scope.copyWith(objectId: value.trim()),
                      );
                    },
                  ),
                ),
                const SizedBox(width: 12),

                // Language dropdown
                Expanded(
                  flex: 2,
                  child: DropdownButtonFormField<String>(
                    value: widget.scope.lang,
                    decoration: InputDecoration(
                      labelText: 'Language',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 10,
                      ),
                    ),
                    items: [
                      for (final lang in widget.languages)
                        DropdownMenuItem(
                          value: lang,
                          child: Text(lang.isEmpty ? 'All' : lang),
                        ),
                    ],
                    onChanged: (value) {
                      widget.onScopeChanged(
                        widget.scope.copyWith(lang: value ?? ''),
                      );
                    },
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
