import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:antinvestor_ui_core/widgets/state_badge.dart';
import 'package:flutter/material.dart';

import 'profile_type_badge.dart';

/// Extracts a display name from a [ProfileObject].
String profileName(ProfileObject profile) {
  final props = profile.properties;
  if (props.fields.containsKey('name')) {
    return props.fields['name']!.stringValue;
  }
  if (profile.contacts.isNotEmpty) {
    return profile.contacts.first.detail;
  }
  return profile.id.length > 12
      ? '${profile.id.substring(0, 12)}...'
      : profile.id;
}

/// A card widget for displaying a [ProfileObject] in list views.
class ProfileCard extends StatelessWidget {
  const ProfileCard({
    super.key,
    required this.profile,
    this.onTap,
  });

  final ProfileObject profile;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = profileName(profile);

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
              ProfileAvatar(
                profileId: profile.id,
                name: name,
                size: 44,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      name,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        ProfileTypeBadge(type: profile.type),
                        const SizedBox(width: 8),
                        StateBadge(state: profile.state),
                      ],
                    ),
                  ],
                ),
              ),
              if (profile.contacts.isNotEmpty) ...[
                const SizedBox(width: 8),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.surfaceContainerLow,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        Icons.contacts_outlined,
                        size: 14,
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '${profile.contacts.length}',
                        style: theme.textTheme.labelSmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ],
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
