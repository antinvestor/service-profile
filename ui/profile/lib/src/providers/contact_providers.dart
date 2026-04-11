import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// Notifier for contact management operations.
class ContactNotifier extends StateNotifier<AsyncValue<void>> {
  ContactNotifier(this._client) : super(const AsyncValue.data(null));
  final ProfileServiceClient _client;

  Future<ContactObject> addContact(AddContactRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.addContact(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<ContactObject> createContact(CreateContactRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.createContact(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> removeContact(RemoveContactRequest request) async {
    state = const AsyncValue.loading();
    try {
      await _client.removeContact(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> createVerification(
    CreateContactVerificationRequest request,
  ) async {
    state = const AsyncValue.loading();
    try {
      await _client.createContactVerification(request);
      state = const AsyncValue.data(null);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<bool> checkVerification(CheckVerificationRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.checkVerification(request);
      state = const AsyncValue.data(null);
      return response.passed;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final contactNotifierProvider =
    StateNotifierProvider<ContactNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(profileServiceClientProvider);
  return ContactNotifier(client);
});
