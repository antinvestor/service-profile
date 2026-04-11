import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/api/api_base.dart';
import 'package:connectrpc/connect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

const _profileUrl = String.fromEnvironment(
  'PROFILE_URL',
  defaultValue: 'https://api.antinvestor.com/profile',
);

final profileTransportProvider = Provider<Transport>((ref) {
  final tokenProvider = ref.watch(authTokenProviderProvider);
  return createTransport(tokenProvider, baseUrl: _profileUrl);
});

final profileServiceClientProvider = Provider<ProfileServiceClient>((ref) {
  final transport = ref.watch(profileTransportProvider);
  return ProfileServiceClient(transport);
});
