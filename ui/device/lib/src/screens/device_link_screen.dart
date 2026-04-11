import 'package:antinvestor_api_device/antinvestor_api_device.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/device_providers.dart';

/// Screen for registering a new device or linking an existing device
/// to a user profile.
class DeviceLinkScreen extends ConsumerStatefulWidget {
  const DeviceLinkScreen({super.key});

  @override
  ConsumerState<DeviceLinkScreen> createState() => _DeviceLinkScreenState();
}

class _DeviceLinkScreenState extends ConsumerState<DeviceLinkScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _profileIdController = TextEditingController();
  final _deviceIdController = TextEditingController();
  bool _isLinking = false; // true = link existing, false = create new
  bool _saving = false;
  String? _error;

  @override
  void dispose() {
    _nameController.dispose();
    _profileIdController.dispose();
    _deviceIdController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () =>
              context.canPop() ? context.pop() : context.go('/devices'),
        ),
        title: Text(
          _isLinking ? 'Link Device' : 'Register Device',
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Mode toggle
              SegmentedButton<bool>(
                segments: const [
                  ButtonSegment(
                    value: false,
                    label: Text('Register New'),
                    icon: Icon(Icons.add_circle_outline, size: 18),
                  ),
                  ButtonSegment(
                    value: true,
                    label: Text('Link Existing'),
                    icon: Icon(Icons.link, size: 18),
                  ),
                ],
                selected: {_isLinking},
                onSelectionChanged: (set) {
                  setState(() {
                    _isLinking = set.first;
                    _error = null;
                  });
                },
              ),
              const SizedBox(height: 24),

              if (!_isLinking) ...[
                // Register new device
                FormFieldCard(
                  label: 'Device Name',
                  description:
                      'A friendly name to identify this device.',
                  child: TextFormField(
                    controller: _nameController,
                    decoration: InputDecoration(
                      hintText: 'e.g., My Laptop',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    validator: (value) {
                      if (value == null || value.trim().isEmpty) {
                        return 'Device name is required';
                      }
                      return null;
                    },
                  ),
                ),
              ] else ...[
                // Link existing device
                FormFieldCard(
                  label: 'Device ID',
                  description:
                      'The ID of the device to link to a profile.',
                  child: TextFormField(
                    controller: _deviceIdController,
                    decoration: InputDecoration(
                      hintText: 'Device ID',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    validator: (value) {
                      if (value == null || value.trim().isEmpty) {
                        return 'Device ID is required';
                      }
                      return null;
                    },
                  ),
                ),
                const SizedBox(height: 16),
                FormFieldCard(
                  label: 'Profile ID',
                  description:
                      'The user profile to link this device to.',
                  child: TextFormField(
                    controller: _profileIdController,
                    decoration: InputDecoration(
                      hintText: 'Profile ID',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    validator: (value) {
                      if (value == null || value.trim().isEmpty) {
                        return 'Profile ID is required';
                      }
                      return null;
                    },
                  ),
                ),
              ],

              // Error
              if (_error != null) ...[
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.errorContainer,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        Icons.error_outline,
                        size: 20,
                        color: theme.colorScheme.onErrorContainer,
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          _error!,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onErrorContainer,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],

              // Submit
              const SizedBox(height: 24),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  OutlinedButton(
                    onPressed: _saving ? null : () => context.go('/devices'),
                    child: const Text('Cancel'),
                  ),
                  const SizedBox(width: 12),
                  FilledButton.icon(
                    onPressed: _saving ? null : _submit,
                    icon: _saving
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child:
                                CircularProgressIndicator(strokeWidth: 2),
                          )
                        : Icon(
                            _isLinking ? Icons.link : Icons.add,
                            size: 18,
                          ),
                    label: Text(
                      _saving
                          ? 'Processing...'
                          : _isLinking
                              ? 'Link Device'
                              : 'Register Device',
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _saving = true;
      _error = null;
    });

    try {
      final notifier = ref.read(deviceNotifierProvider.notifier);

      if (_isLinking) {
        await notifier.link(
          LinkRequest()
            ..id = _deviceIdController.text.trim()
            ..profileId = _profileIdController.text.trim(),
        );
      } else {
        await notifier.create(
          CreateRequest()..name = _nameController.text.trim(),
        );
      }

      if (mounted) {
        context.go('/devices');
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              _isLinking
                  ? 'Device linked successfully'
                  : 'Device registered successfully',
            ),
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _saving = false;
          _error = friendlyError(e);
        });
      }
    }
  }
}
