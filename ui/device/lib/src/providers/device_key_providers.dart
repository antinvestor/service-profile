import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'device_transport_provider.dart';

/// Search keys for a device by device ID.
final deviceKeysProvider =
    FutureProvider.family<List<KeyObject>, String>((ref, deviceId) async {
  final client = ref.watch(deviceServiceClientProvider);
  final request = SearchKeyRequest()..deviceId = deviceId;
  final response = await client.searchKey(request);
  return response.data;
});

/// Notifier for key mutations (add, remove).
class DeviceKeyNotifier extends StateNotifier<AsyncValue<void>> {
  DeviceKeyNotifier(this._client) : super(const AsyncValue.data(null));
  final DeviceServiceClient _client;

  Future<KeyObject> addKey(AddKeyRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.addKey(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> removeKey(RemoveKeyRequest request) async {
    state = const AsyncValue.loading();
    try {
      await _client.removeKey(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final deviceKeyNotifierProvider =
    StateNotifierProvider<DeviceKeyNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(deviceServiceClientProvider);
  return DeviceKeyNotifier(client);
});
