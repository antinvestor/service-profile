import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/api/api_base.dart';
import 'package:connectrpc/connect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

const _geolocationUrl = String.fromEnvironment(
  'GEOLOCATION_URL',
  defaultValue: 'https://api.antinvestor.com/geolocation',
);

final geolocationTransportProvider = Provider<Transport>((ref) {
  final tokenProvider = ref.watch(authTokenProviderProvider);
  return createTransport(tokenProvider, baseUrl: _geolocationUrl);
});

final geolocationServiceClientProvider =
    Provider<GeolocationServiceClient>((ref) {
  final transport = ref.watch(geolocationTransportProvider);
  return GeolocationServiceClient(transport);
});
