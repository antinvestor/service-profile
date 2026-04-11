import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/api/stream_helpers.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'device_transport_provider.dart';

/// List activity logs for a device.
final deviceLogsProvider =
    FutureProvider.family<List<DeviceLog>, String>((ref, deviceId) async {
  final client = ref.watch(deviceServiceClientProvider);
  final request = ListLogsRequest()..deviceId = deviceId;
  final stream = client.listLogs(request);
  return collectStream<ListLogsResponse, DeviceLog>(
    stream,
    extract: (r) => r.data,
  );
});
