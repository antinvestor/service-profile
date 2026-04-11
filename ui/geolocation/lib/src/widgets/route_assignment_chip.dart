import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter/material.dart';

/// Chip showing a subject assigned to a route with optional remove action.
class RouteAssignmentChip extends StatelessWidget {
  const RouteAssignmentChip({
    super.key,
    required this.assignment,
    this.onRemove,
  });

  final RouteAssignmentObject assignment;
  final VoidCallback? onRemove;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final validFrom = assignment.hasValidFrom()
        ? DateTime.fromMillisecondsSinceEpoch(
            assignment.validFrom.seconds.toInt() * 1000,
          )
        : null;
    final validUntil = assignment.hasValidUntil()
        ? DateTime.fromMillisecondsSinceEpoch(
            assignment.validUntil.seconds.toInt() * 1000,
          )
        : null;

    return Chip(
      avatar: CircleAvatar(
        backgroundColor: theme.colorScheme.primaryContainer,
        child: Icon(
          Icons.person_pin_circle,
          size: 16,
          color: theme.colorScheme.onPrimaryContainer,
        ),
      ),
      label: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            assignment.subjectId,
            style: theme.textTheme.labelMedium?.copyWith(
              fontWeight: FontWeight.w600,
            ),
            overflow: TextOverflow.ellipsis,
          ),
          if (validFrom != null || validUntil != null)
            Text(
              [
                if (validFrom != null) 'From: ${_formatDate(validFrom)}',
                if (validUntil != null) 'Until: ${_formatDate(validUntil)}',
              ].join(' | '),
              style: theme.textTheme.labelSmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
        ],
      ),
      deleteIcon: onRemove != null
          ? const Icon(Icons.close, size: 16)
          : null,
      onDeleted: onRemove,
    );
  }

  String _formatDate(DateTime dt) {
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-'
        '${dt.day.toString().padLeft(2, '0')}';
  }
}
