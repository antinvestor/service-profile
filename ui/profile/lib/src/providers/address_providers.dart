import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'profile_transport_provider.dart';

/// Notifier for address operations.
class AddressNotifier extends StateNotifier<AsyncValue<void>> {
  AddressNotifier(this._client) : super(const AsyncValue.data(null));
  final ProfileServiceClient _client;

  Future<AddressObject> addAddress(AddAddressRequest request) async {
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
    StateNotifierProvider<AddressNotifier, AsyncValue<void>>((ref) {
  final client = ref.watch(profileServiceClientProvider);
  return AddressNotifier(client);
});
