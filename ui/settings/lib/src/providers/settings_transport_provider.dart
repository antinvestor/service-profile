import 'package:antinvestor_api_settings/antinvestor_api_settings.dart';
import 'package:antinvestor_ui_core/api/api_base.dart';
import 'package:connectrpc/connect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

const _settingsUrl = String.fromEnvironment(
  'SETTINGS_URL',
  defaultValue: 'https://api.antinvestor.com/settings',
);

final settingsTransportProvider = Provider<Transport>((ref) {
  final tokenProvider = ref.watch(authTokenProviderProvider);
  return createTransport(tokenProvider, baseUrl: _settingsUrl);
});

final settingsServiceClientProvider = Provider<SettingsServiceClient>((ref) {
  final transport = ref.watch(settingsTransportProvider);
  return SettingsServiceClient(transport);
});
