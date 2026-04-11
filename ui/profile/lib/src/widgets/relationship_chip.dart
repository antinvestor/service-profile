import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter/material.dart';

import 'profile_card.dart';

/// A chip for displaying a [RelationshipObject].
class RelationshipChip extends StatelessWidget {
  const RelationshipChip({
    super.key,
    required this.relationship,
    this.onTap,
    this.onDelete,
  });

  final RelationshipObject relationship;
  final VoidCallback? onTap;
  final VoidCallback? onDelete;

  (String, Color, IconData) _typeInfo(RelationshipType type) {
    return switch (type) {
      RelationshipType.MEMBER => ('Member', Colors.blue, Icons.group),
      RelationshipType.AFFILIATED =>
        ('Affiliated', Colors.teal, Icons.handshake),
      RelationshipType.BLACK_LISTED =>
        ('Blocked', Colors.red, Icons.block),
      _ => ('Unknown', Colors.grey, Icons.help_outline),
    };
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final (label, color, icon) = _typeInfo(relationship.type);

    final peerName = relationship.hasPeerProfile()
        ? profileName(relationship.peerProfile)
        : relationship.hasChildEntry()
            ? relationship.childEntry.objectName
            : relationship.id.length > 12
                ? '${relationship.id.substring(0, 12)}...'
                : relationship.id;

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: color.withAlpha(60)),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: color.withAlpha(25),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(icon, size: 16, color: color),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      peerName,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    Text(
                      label,
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: color,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
              if (onDelete != null)
                IconButton(
                  icon: Icon(Icons.close, size: 16),
                  onPressed: onDelete,
                  visualDensity: VisualDensity.compact,
                  tooltip: 'Remove',
                ),
            ],
          ),
        ),
      ),
    );
  }
}
