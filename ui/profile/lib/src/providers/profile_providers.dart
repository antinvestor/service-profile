import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/api/stream_helpers.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// Search profiles by query string.
final profileSearchProvider =
    FutureProvider.family<List<ProfileObject>, String>((ref, query) async {
  final client = ref.watch(profileServiceClientProvider);
  final request = SearchRequest()..query = query;
  final stream = client.search(request);
  return collectStream<SearchResponse, ProfileObject>(
    stream,
    extract: (r) => r.data,
  );
});

/// Get a single profile by ID.
final profileByIdProvider =
    FutureProvider.family<ProfileObject, String>((ref, id) async {
  final client = ref.watch(profileServiceClientProvider);
  final request = GetByIdRequest()..id = id;
  final response = await client.getById(request);
  return response.data;
});

/// Get a profile by contact (email or phone).
final profileByContactProvider =
    FutureProvider.family<ProfileObject, String>((ref, contact) async {
  final client = ref.watch(profileServiceClientProvider);
  final request = GetByContactRequest()..contact = contact;
  final response = await client.getByContact(request);
  return response.data;
});

/// Notifier for profile mutations (create, update, merge).
class ProfileNotifier extends Notifier<AsyncValue<void>> {
  @override
  AsyncValue<void> build() => const AsyncValue.data(null);

  ProfileServiceClient get _client =>
      ref.read(profileServiceClientProvider);

  Future<ProfileObject> create(CreateRequest request) async {
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

  Future<ProfileObject> update(UpdateRequest request) async {
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

  Future<ProfileObject> merge(MergeRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.merge(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final profileNotifierProvider =
    NotifierProvider<ProfileNotifier, AsyncValue<void>>(ProfileNotifier.new);
