import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'device_transport_provider.dart';

/// Notifier for presence updates.
class PresenceNotifier extends StateNotifier<AsyncValue<void>> {
  PresenceNotifier(this._client) : super(const AsyncValue.data(null));
  final DeviceServiceClient _client;

  Future<PresenceObject> updatePresence(
      UpdatePresenceRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.updatePresence(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final presenceNotifierProvider =
    StateNotifierProvider<PresenceNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(deviceServiceClientProvider);
  return PresenceNotifier(client);
});
