import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// Notifier for contact management operations.
class ContactNotifier extends Notifier<AsyncValue<void>> {
  @override
  AsyncValue<void> build() => const AsyncValue.data(null);

  ProfileServiceClient get _client =>
      ref.read(profileServiceClientProvider);

  Future<ProfileObject> addContact(AddContactRequest request) async {
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
      return response.success;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final contactNotifierProvider =
    NotifierProvider<ContactNotifier, AsyncValue<void>>(ContactNotifier.new);
