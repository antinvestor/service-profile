import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/api/stream_helpers.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'device_transport_provider.dart';

/// Search devices by query string.
final deviceSearchProvider =
    FutureProvider.family<List<DeviceObject>, String>((ref, query) async {
  final client = ref.watch(deviceServiceClientProvider);
  final request = SearchRequest()
    ..query = query
    ..count = 50;
  final stream = client.search(request);
  return collectStream<SearchResponse, DeviceObject>(
    stream,
    extract: (r) => r.data,
  );
});

/// Get device(s) by ID.
final deviceByIdProvider =
    FutureProvider.family<DeviceObject, String>((ref, id) async {
  final client = ref.watch(deviceServiceClientProvider);
  final request = GetByIdRequest()..id.add(id);
  final response = await client.getById(request);
  if (response.data.isEmpty) {
    throw Exception('Device not found');
  }
  return response.data.first;
});

/// Notifier for device mutations (create, update, remove, link).
class DeviceNotifier extends StateNotifier<AsyncValue<void>> {
  DeviceNotifier(this._client) : super(const AsyncValue.data(null));
  final DeviceServiceClient _client;

  Future<DeviceObject> create(CreateRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.create(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<DeviceObject> update(UpdateRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.update(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> remove(RemoveRequest request) async {
    state = const AsyncValue.loading();
    try {
      await _client.remove(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<DeviceObject> link(LinkRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.link(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final deviceNotifierProvider =
    StateNotifierProvider<DeviceNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(deviceServiceClientProvider);
  return DeviceNotifier(client);
});
