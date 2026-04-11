import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'device_providers.dart';

/// Searches devices for a given profileId and returns the most recent
/// device's presence status.
///
/// Returns [PresenceStatus.OFFLINE] as fallback if no device is found or
/// if the device has no presence information.
final presenceByProfileProvider =
    FutureProvider.family<PresenceStatus, String>((ref, profileId) async {
  if (profileId.isEmpty) return PresenceStatus.OFFLINE;

  // Search for devices associated with this profile.
  final devices = await ref.watch(deviceSearchProvider(profileId).future);

  if (devices.isEmpty) return PresenceStatus.OFFLINE;

  // Return the most recently seen device's presence status.
  // Devices are assumed to be ordered by recency from the API.
  final device = devices.first;
  try {
    return device.presence;
  } catch (_) {
    return PresenceStatus.OFFLINE;
  }
});
