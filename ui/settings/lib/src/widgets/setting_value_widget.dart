import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/settings_providers.dart';

/// Displays a setting value inline by looking it up with [settingByKeyProvider].
///
/// Pass the setting name, object type, and object ID. The widget resolves the
/// value and renders it as text, falling back to [placeholder] if not found.
///
/// ```dart
/// SettingValueWidget(name: 'theme', object: 'profile', objectId: profileId)
/// SettingValueWidget(name: 'currency', object: 'tenant', objectId: tenantId,
///     prefix: 'Currency: ')
/// ```
class SettingValueWidget extends ConsumerWidget {
  const SettingValueWidget({
    super.key,
    required this.name,
    required this.object,
    required this.objectId,
    this.lang = '',
    this.module = '',
    this.style,
    this.prefix = '',
    this.placeholder = '\u2014',
    this.icon,
  });

  final String name;
  final String object;
  final String objectId;
  final String lang;
  final String module;
  final TextStyle? style;
  final String prefix;
  final String placeholder;
  final IconData? icon;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final key = SettingKey(
      name: name,
      object: object,
      objectId: objectId,
      lang: lang,
      module: module,
    );

    final settingAsync = ref.watch(settingByKeyProvider(key));

    final display = settingAsync.when(
      data: (setting) {
        final value = setting.value;
        if (value.isEmpty) return placeholder;
        return '$prefix$value';
      },
      loading: () => '$prefix...',
      error: (_, _) => placeholder,
    );

    final textStyle = style ?? theme.textTheme.bodyMedium;

    if (icon != null) {
      return Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 16, color: theme.colorScheme.onSurfaceVariant),
          const SizedBox(width: 6),
          Flexible(
            child: Text(
              display,
              style: textStyle,
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ],
      );
    }

    return Text(display, style: textStyle, overflow: TextOverflow.ellipsis);
  }
}
