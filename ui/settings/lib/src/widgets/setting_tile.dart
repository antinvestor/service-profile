import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:flutter/material.dart';

/// A list tile for displaying a single setting with its key, a preview
/// of the value, and the last-updated timestamp.
class SettingTile extends StatelessWidget {
  const SettingTile({
    super.key,
    required this.setting,
    this.onTap,
  });

  final SettingObject setting;
  final VoidCallback? onTap;

  String _valuePreview(String value) {
    if (value.length <= 80) return value;
    return '${value.substring(0, 77)}...';
  }

  String _formatUpdated(SettingObject s) {
    if (!s.hasUpdated()) return '';
    final ts = s.updated.toDateTime();
    final now = DateTime.now();
    final diff = now.difference(ts);
    if (diff.inMinutes < 1) return 'just now';
    if (diff.inHours < 1) return '${diff.inMinutes}m ago';
    if (diff.inDays < 1) return '${diff.inHours}h ago';
    if (diff.inDays < 30) return '${diff.inDays}d ago';
    return '${ts.year}-${ts.month.toString().padLeft(2, '0')}-${ts.day.toString().padLeft(2, '0')}';
  }

  IconData _inferIcon(String value) {
    final lower = value.toLowerCase().trim();
    if (lower == 'true' || lower == 'false') return Icons.toggle_on_outlined;
    if (double.tryParse(lower) != null) return Icons.pin_outlined;
    if (lower.startsWith('{') || lower.startsWith('[')) {
      return Icons.data_object;
    }
    return Icons.text_fields;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final updated = _formatUpdated(setting);

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: theme.colorScheme.outlineVariant),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(
                  _inferIcon(setting.value),
                  size: 20,
                  color: theme.colorScheme.onPrimaryContainer,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      setting.key.name,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      _valuePreview(setting.value),
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              if (updated.isNotEmpty) ...[
                const SizedBox(width: 8),
                Text(
                  updated,
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
              const SizedBox(width: 4),
              Icon(
                Icons.chevron_right,
                size: 20,
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
