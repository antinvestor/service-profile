import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter/material.dart';

/// A list tile for displaying a [ContactObject].
class ContactListTile extends StatelessWidget {
  const ContactListTile({
    super.key,
    required this.contact,
    this.trailing,
    this.onTap,
  });

  final ContactObject contact;
  final Widget? trailing;
  final VoidCallback? onTap;

  IconData _typeIcon(ContactType type) {
    return switch (type) {
      ContactType.EMAIL => Icons.email_outlined,
      ContactType.MSISDN => Icons.phone_outlined,
      _ => Icons.contact_mail_outlined,
    };
  }

  String _communicationLabel(CommunicationLevel level) {
    return switch (level) {
      CommunicationLevel.ALL => 'All communications',
      CommunicationLevel.NO_CONTACT => 'No communications',
      _ => level.name,
    };
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ListTile(
      leading: Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: theme.colorScheme.primaryContainer,
          borderRadius: BorderRadius.circular(10),
        ),
        child: Icon(
          _typeIcon(contact.type),
          size: 20,
          color: theme.colorScheme.onPrimaryContainer,
        ),
      ),
      title: Row(
        children: [
          Expanded(
            child: Text(
              contact.detail,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ),
          if (contact.verified)
            Padding(
              padding: const EdgeInsets.only(left: 6),
              child: Icon(
                Icons.verified,
                size: 16,
                color: Colors.green.shade600,
              ),
            ),
        ],
      ),
      subtitle: Text(
        _communicationLabel(contact.communicationLevel),
        style: theme.textTheme.bodySmall?.copyWith(
          color: theme.colorScheme.onSurfaceVariant,
        ),
      ),
      trailing: trailing,
      onTap: onTap,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
    );
  }
}
