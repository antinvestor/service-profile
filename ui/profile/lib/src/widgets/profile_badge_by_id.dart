import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/profile_providers.dart';

/// Resolves a profile by ID and renders a ProfileBadge.
/// Drop this into ANY screen across ANY service to show a profile.
class ProfileBadgeById extends ConsumerWidget {
  const ProfileBadgeById({
    super.key,
    required this.profileId,
    this.avatarSize = 36,
    this.trailing,
  });

  final String profileId;
  final double avatarSize;
  final Widget? trailing;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (profileId.isEmpty) return const SizedBox.shrink();

    final profileAsync = ref.watch(profileByIdProvider(profileId));
    final fallback = _truncateId(profileId);

    return profileAsync.when(
      data: (profile) {
        final name = _extractName(profile);
        return ProfileBadge(
          profileId: profileId,
          name: name,
          avatarSize: avatarSize,
          trailing: trailing,
        );
      },
      loading: () => ProfileBadge(
        profileId: profileId,
        name: fallback,
        avatarSize: avatarSize,
        trailing: trailing,
      ),
      error: (_, __) => ProfileBadge(
        profileId: profileId,
        name: fallback,
        avatarSize: avatarSize,
        trailing: trailing,
      ),
    );
  }

  String _extractName(dynamic profile) {
    // Try properties.name first.
    try {
      final props = profile.properties;
      if (props.fields.containsKey('name')) {
        final name = props.fields['name']!.stringValue;
        if (name.isNotEmpty) return name;
      }
    } catch (_) {}
    // Fall back to first contact.
    try {
      if (profile.contacts.isNotEmpty) {
        return profile.contacts.first.detail;
      }
    } catch (_) {}
    return _truncateId(profileId);
  }
}

/// Compact version -- just shows the resolved name as text.
class ProfileNameText extends ConsumerWidget {
  const ProfileNameText({
    super.key,
    required this.profileId,
    this.style,
    this.prefix = '',
  });

  final String profileId;
  final TextStyle? style;
  final String prefix;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (profileId.isEmpty) {
      return Text('$prefix\u2014', style: style);
    }

    final profileAsync = ref.watch(profileByIdProvider(profileId));
    final fallback = _truncateId(profileId);

    final name = profileAsync.when(
      data: (p) {
        try {
          final props = p.properties;
          if (props.fields.containsKey('name')) {
            final n = props.fields['name']!.stringValue;
            if (n.isNotEmpty) return n;
          }
        } catch (_) {}
        try {
          if (p.contacts.isNotEmpty) return p.contacts.first.detail;
        } catch (_) {}
        return fallback;
      },
      loading: () => fallback,
      error: (_, __) => fallback,
    );

    return Text('$prefix$name', style: style);
  }
}

String _truncateId(String id) =>
    id.length > 12 ? '${id.substring(0, 12)}...' : id;
