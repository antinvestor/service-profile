import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/api/stream_helpers.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// List relationships for a profile.
final relationshipListProvider =
    FutureProvider.family<List<RelationshipObject>, String>(
        (ref, profileId) async {
  final client = ref.watch(profileServiceClientProvider);
  final request = ListRelationshipRequest()..profileId = profileId;
  final stream = client.listRelationship(request);
  return collectStream<ListRelationshipResponse, RelationshipObject>(
    stream,
    extract: (r) => r.data,
  );
});

/// Notifier for relationship mutations.
class RelationshipNotifier extends StateNotifier<AsyncValue<void>> {
  RelationshipNotifier(this._client) : super(const AsyncValue.data(null));
  final ProfileServiceClient _client;

  Future<RelationshipObject> add(AddRelationshipRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.addRelationship(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> delete(DeleteRelationshipRequest request) async {
    state = const AsyncValue.loading();
    try {
      await _client.deleteRelationship(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final relationshipNotifierProvider =
    StateNotifierProvider<RelationshipNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(profileServiceClientProvider);
  return RelationshipNotifier(client);
});
