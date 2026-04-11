import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'geolocation_transport_provider.dart';

/// Get location track for a subject.
final getTrackProvider =
    FutureProvider.family<List<LocationPointObject>, String>(
        (ref, subjectId) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetTrackRequest()
    ..subjectId = subjectId
    ..limit = 100;
  final response = await client.getTrack(request);
  return response.data;
});

/// Get geo events for a subject.
final getSubjectEventsProvider =
    FutureProvider.family<List<GeoEventObject>, String>(
        (ref, subjectId) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetSubjectEventsRequest()
    ..subjectId = subjectId
    ..limit = 100;
  final response = await client.getSubjectEvents(request);
  return response.data;
});

/// Parameters for nearby subjects query.
class NearbySubjectsParams {
  const NearbySubjectsParams({
    required this.subjectId,
    this.radiusMeters = 1000,
    this.limit = 20,
  });

  final String subjectId;
  final double radiusMeters;
  final int limit;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is NearbySubjectsParams &&
          other.subjectId == subjectId &&
          other.radiusMeters == radiusMeters &&
          other.limit == limit;

  @override
  int get hashCode =>
      Object.hash(subjectId, radiusMeters, limit);
}

/// Get nearby subjects.
final getNearbySubjectsProvider =
    FutureProvider.family<List<NearbySubjectObject>, NearbySubjectsParams>(
        (ref, params) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetNearbySubjectsRequest()
    ..subjectId = params.subjectId
    ..radiusMeters = params.radiusMeters
    ..limit = params.limit;
  final response = await client.getNearbySubjects(request);
  return response.data;
});

/// Parameters for nearby areas query.
class NearbyAreasParams {
  const NearbyAreasParams({
    required this.latitude,
    required this.longitude,
    this.radiusMeters = 1000,
    this.limit = 20,
  });

  final double latitude;
  final double longitude;
  final double radiusMeters;
  final int limit;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is NearbyAreasParams &&
          other.latitude == latitude &&
          other.longitude == longitude &&
          other.radiusMeters == radiusMeters &&
          other.limit == limit;

  @override
  int get hashCode =>
      Object.hash(latitude, longitude, radiusMeters, limit);
}

/// Get nearby areas.
final getNearbyAreasProvider =
    FutureProvider.family<List<NearbyAreaObject>, NearbyAreasParams>(
        (ref, params) async {
  final client = ref.watch(geolocationServiceClientProvider);
  final request = GetNearbyAreasRequest()
    ..latitude = params.latitude
    ..longitude = params.longitude
    ..radiusMeters = params.radiusMeters
    ..limit = params.limit;
  final response = await client.getNearbyAreas(request);
  return response.data;
});
