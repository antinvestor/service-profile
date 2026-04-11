import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'geolocation_transport_provider.dart';

/// Search areas by query string.
final searchAreasProvider =
    FutureProvider.family<List<AreaObject>, String>((ref, query) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = SearchAreasRequest()
    ..query = query
    ..limit = 50;
  final response = await client.searchAreas(request);
  return response.data;
});

/// Get a single area by ID.
final getAreaProvider =
    FutureProvider.family<AreaObject, String>((ref, id) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetAreaRequest()..id = id;
  final response = await client.getArea(request);
  return response.data;
});

/// Notifier for area mutations (create, update, delete).
class AreaNotifier extends StateNotifier<AsyncValue<void>> {
  AreaNotifier(this._client) : super(const AsyncValue.data(null));
  final GeolocationServiceClient _client;

  Future<AreaObject> createArea(AreaObject data) async {
    state = const AsyncValue.loading();
    try {
      final request = CreateAreaRequest()..data = data;
      final response = await _client.createArea(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<AreaObject> updateArea(UpdateAreaRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.updateArea(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> deleteArea(String id) async {
    state = const AsyncValue.loading();
    try {
      final request = DeleteAreaRequest()..id = id;
      await _client.deleteArea(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final areaNotifierProvider =
    StateNotifierProvider<AreaNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(geolocationServiceClientProvider);
  return AreaNotifier(client);
});
