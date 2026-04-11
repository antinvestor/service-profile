import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/api/stream_helpers.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// Search roster contacts for a profile.
final rosterSearchProvider =
    FutureProvider.family<List<RosterObject>, String>((ref, profileId) async {
  final client = ref.watch(profileServiceClientProvider);
  final request = SearchRosterRequest()..profileId = profileId;
  final stream = client.searchRoster(request);
  return collectStream<SearchRosterResponse, RosterObject>(
    stream,
    extract: (r) => r.data,
  );
});

/// Notifier for roster mutations.
class RosterNotifier extends StateNotifier<AsyncValue<void>> {
  RosterNotifier(this._client) : super(const AsyncValue.data(null));
  final ProfileServiceClient _client;

  Future<RosterObject> add(AddRosterRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.addRoster(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> remove(RemoveRosterRequest request) async {
    state = const AsyncValue.loading();
    try {
      await _client.removeRoster(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final rosterNotifierProvider =
    StateNotifierProvider<RosterNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(profileServiceClientProvider);
  return RosterNotifier(client);
});
