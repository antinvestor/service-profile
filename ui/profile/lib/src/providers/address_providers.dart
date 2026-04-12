import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// Notifier for address operations.
class AddressNotifier extends Notifier<AsyncValue<void>> {
  @override
  AsyncValue<void> build() => const AsyncValue.data(null);

  ProfileServiceClient get _client =>
      ref.read(profileServiceClientProvider);

  Future<ProfileObject> addAddress(AddAddressRequest request) async {
    state = const AsyncValue.loading();
    try {
      final response = await _client.addAddress(request);
      state = const AsyncValue.data(null);
      return response.data;
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }
}

final addressNotifierProvider =
    NotifierProvider<AddressNotifier, AsyncValue<void>>(AddressNotifier.new);
