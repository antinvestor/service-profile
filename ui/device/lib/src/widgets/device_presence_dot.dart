import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/presence_by_profile_provider.dart';
import 'presence_indicator.dart';

/// Resolves a profile's most recent device presence and renders a
/// [PresenceIndicator] dot.
///
/// Drop this next to any profile avatar to show online/offline status.
///
/// ```dart
/// Stack(
///   children: [
///     ProfileAvatar(profileId: id, name: name),
///     Positioned(bottom: 0, right: 0, child: DevicePresenceDot(profileId: id)),
///   ],
/// )
/// ```
class DevicePresenceDot extends ConsumerWidget {
  const DevicePresenceDot({
    super.key,
    required this.profileId,
    this.size = 12,
    this.showLabel = false,
  });

  final String profileId;
  final double size;
  final bool showLabel;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (profileId.isEmpty) {
      return PresenceIndicator(
        status: PresenceStatus.OFFLINE,
        size: size,
        showLabel: showLabel,
      );
    }

    final presenceAsync = ref.watch(presenceByProfileProvider(profileId));

    final status = presenceAsync.when(
      data: (s) => s,
      loading: () => PresenceStatus.OFFLINE,
      error: (_, _) => PresenceStatus.OFFLINE,
    );

    return PresenceIndicator(
      status: status,
      size: size,
      showLabel: showLabel,
    );
  }
}
