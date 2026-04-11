import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/api/api_base.dart';
import 'package:connectrpc/connect.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

const _deviceUrl = String.fromEnvironment(
  'DEVICE_URL',
  defaultValue: 'https://api.antinvestor.com/device',
);

final deviceTransportProvider = Provider<Transport>((ref) {
  final tokenProvider = ref.watch(authTokenProviderProvider);
  return createTransport(tokenProvider, baseUrl: _deviceUrl);
});

final deviceServiceClientProvider = Provider<DeviceServiceClient>((ref) {
  final transport = ref.watch(deviceTransportProvider);
  return DeviceServiceClient(transport);
});
