import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/status_badge.dart';
import 'package:flutter/material.dart';

/// Displays a colored badge indicating the profile type.
class ProfileTypeBadge extends StatelessWidget {
  const ProfileTypeBadge({super.key, required this.type});

  final ProfileType type;

  @override
  Widget build(BuildContext context) {
    return StatusBadge.fromEnum(
      value: type,
      mapper: (t) => switch (t) {
        ProfileType.PERSON => ('Person', Colors.blue, Icons.person),
        ProfileType.INSTITUTION =>
          ('Institution', Colors.teal, Icons.business),
        ProfileType.BOT => ('Bot', Colors.purple, Icons.smart_toy),
        _ => ('Unknown', Colors.grey, Icons.help_outline),
      },
    );
  }
}
