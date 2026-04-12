import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'geolocation_transport_provider.dart';

/// Search routes by owner ID.
final searchRoutesProvider =
    FutureProvider.family<List<RouteObject>, String>((ref, ownerId) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = SearchRoutesRequest()
    ..ownerId = ownerId
    ..limit = 50;
  final response = await client.searchRoutes(request);
  return response.data;
});

/// Get a single route by ID.
final getRouteProvider =
    FutureProvider.family<RouteObject, String>((ref, id) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetRouteRequest()..id = id;
  final response = await client.getRoute(request);
  return response.data;
});

/// Get route assignments for a subject.
final subjectRouteAssignmentsProvider =
    FutureProvider.family<List<RouteAssignmentObject>, String>(
        (ref, subjectId) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetSubjectRouteAssignmentsRequest()
    ..subjectId = subjectId;
  final response = await client.getSubjectRouteAssignments(request);
  return response.data;
});

/// Notifier for route mutations.
class RouteNotifier extends Notifier<AsyncValue<void>> {
  @override
  AsyncValue<void> build() => const AsyncValue.data(null);

  GeolocationServiceClient get _client =>
      ref.read(geolocationServiceClientProvider);

  Future<RouteObject> createRoute(RouteObject data) async {
    state = const AsyncValue.loading();
    try {
      final request = CreateRouteRequest()..data = data;
      final response = await _client.createRoute(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<RouteObject> updateRoute(UpdateRouteRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.updateRoute(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> deleteRoute(String id) async {
    state = const AsyncValue.loading();
    try {
      final request = DeleteRouteRequest()..id = id;
      await _client.deleteRoute(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<RouteAssignmentObject> assignRoute({
    required String subjectId,
    required String routeId,
    Timestamp? validFrom,
    Timestamp? validUntil,
  }) async {
    state = const AsyncValue.loading();
    try {
      final request = AssignRouteRequest()
        ..subjectId = subjectId
        ..routeId = routeId;
      if (validFrom != null) request.validFrom = validFrom;
      if (validUntil != null) request.validUntil = validUntil;
      final response = await _client.assignRoute(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> unassignRoute(String assignmentId) async {
    state = const AsyncValue.loading();
    try {
      final request = UnassignRouteRequest()..id = assignmentId;
      await _client.unassignRoute(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final routeNotifierProvider =
    NotifierProvider<RouteNotifier, AsyncValue<void>>(RouteNotifier.new);
