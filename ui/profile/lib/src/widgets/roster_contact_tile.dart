import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter/material.dart';

/// A list tile for displaying a [RosterObject].
class RosterContactTile extends StatelessWidget {
  const RosterContactTile({
    super.key,
    required this.roster,
    this.onRemove,
  });

  final RosterObject roster;
  final VoidCallback? onRemove;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final contact = roster.hasContact() ? roster.contact : null;
    final displayText =
        contact != null ? contact.detail : roster.profileId;

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: theme.colorScheme.outlineVariant),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
        child: Row(
          children: [
            Container(
              width: 36,
              height: 36,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Icon(
                contact?.type == ContactType.EMAIL
                    ? Icons.email_outlined
                    : Icons.phone_outlined,
                size: 18,
                color: theme.colorScheme.onPrimaryContainer,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    displayText,
                    style: theme.textTheme.bodyMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  if (roster.profileId.isNotEmpty)
                    Text(
                      'Profile: ${roster.profileId.length > 16 ? '${roster.profileId.substring(0, 16)}...' : roster.profileId}',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                ],
              ),
            ),
            if (onRemove != null)
              IconButton(
                icon: Icon(
                  Icons.remove_circle_outline,
                  color: theme.colorScheme.error,
                  size: 20,
                ),
                onPressed: onRemove,
                tooltip: 'Remove',
              ),
          ],
        ),
      ),
    );
  }
}
