import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/device_log_providers.dart';
import '../widgets/session_log_entry.dart';

/// Screen displaying session activity logs for a device.
class DeviceLogsScreen extends ConsumerWidget {
  const DeviceLogsScreen({
    super.key,
    required this.deviceId,
  });

  final String deviceId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final asyncLogs = ref.watch(deviceLogsProvider(deviceId));

    return Column(
      children: [
        // Header
        Padding(
          padding: const EdgeInsets.fromLTRB(24, 16, 24, 8),
          child: Row(
            children: [
              Text(
                'Session Logs',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const Spacer(),
              IconButton(
                icon: const Icon(Icons.refresh, size: 20),
                tooltip: 'Refresh',
                onPressed: () => ref.invalidate(deviceLogsProvider(deviceId)),
              ),
            ],
          ),
        ),

        // Logs list
        Expanded(
          child: asyncLogs.when(
            loading: () =>
                const Center(child: CircularProgressIndicator()),
            error: (error, _) => Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 48,
                    color: theme.colorScheme.error,
                  ),
                  const SizedBox(height: 16),
                  Text(friendlyError(error)),
                  const SizedBox(height: 12),
                  FilledButton.tonal(
                    onPressed: () =>
                        ref.invalidate(deviceLogsProvider(deviceId)),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            ),
            data: (logs) {
              if (logs.isEmpty) {
                return Center(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Container(
                        width: 64,
                        height: 64,
                        decoration: BoxDecoration(
                          color: theme.colorScheme.surfaceContainerLow,
                          borderRadius: BorderRadius.circular(16),
                        ),
                        child: Icon(
                          Icons.history,
                          size: 28,
                          color: theme.colorScheme.primary.withAlpha(160),
                        ),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'No session logs found.',
                        style: theme.textTheme.titleMedium?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ],
                  ),
                );
              }

              return ListView.separated(
                padding: const EdgeInsets.fromLTRB(24, 8, 24, 24),
                itemCount: logs.length,
                separatorBuilder: (_, __) => const SizedBox(height: 8),
                itemBuilder: (context, index) {
                  return SessionLogEntry(log: logs[index]);
                },
              );
            },
          ),
        ),
      ],
    );
  }
}
