import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:antinvestor_ui_core/widgets/gradient_button.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/profile_providers.dart';
import '../widgets/profile_card.dart';

/// Screen for editing an existing profile's properties.
class ProfileEditScreen extends ConsumerStatefulWidget {
  const ProfileEditScreen({super.key, required this.profileId});

  final String profileId;

  @override
  ConsumerState<ProfileEditScreen> createState() => _ProfileEditScreenState();
}

class _ProfileEditScreenState extends ConsumerState<ProfileEditScreen> {
  final _formKey = GlobalKey<FormState>();
  final _controllers = <String, TextEditingController>{};
  bool _isLoading = false;
  String? _error;
  bool _initialized = false;

  @override
  void dispose() {
    for (final c in _controllers.values) {
      c.dispose();
    }
    super.dispose();
  }

  void _initControllers(ProfileObject profile) {
    if (_initialized) return;
    _initialized = true;

    final fields = profile.properties.fields;
    for (final entry in fields.entries) {
      _controllers[entry.key] = TextEditingController(
        text: _extractValue(entry.value),
      );
    }

    // Ensure common fields exist.
    for (final key in ['name', 'email', 'phone', 'description']) {
      _controllers.putIfAbsent(key, () => TextEditingController());
    }
  }

  String _extractValue(Value value) {
    if (value.hasStringValue()) return value.stringValue;
    if (value.hasNumberValue()) return value.numberValue.toString();
    if (value.hasBoolValue()) return value.boolValue.toString();
    return value.toString();
  }

  Future<void> _save(ProfileObject profile) async {
    if (!_formKey.currentState!.validate()) return;
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final props = Struct();
      for (final entry in _controllers.entries) {
        final text = entry.value.text.trim();
        if (text.isNotEmpty) {
          props.fields[entry.key] = Value()..stringValue = text;
        }
      }

      final request = UpdateRequest()
        ..id = widget.profileId
        ..properties = props
        ..state = profile.state;

      await ref.read(profileNotifierProvider.notifier).update(request);
      ref.invalidate(profileByIdProvider(widget.profileId));

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Profile updated'),
            behavior: SnackBarBehavior.floating,
            width: 220,
          ),
        );
        context.pop();
      }
    } catch (e) {
      setState(() {
        _error = friendlyError(e);
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final asyncProfile = ref.watch(profileByIdProvider(widget.profileId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Edit Profile')),
      body: asyncProfile.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text(friendlyError(e))),
        data: (profile) {
          _initControllers(profile);
          final name = profileName(profile);

          return Form(
            key: _formKey,
            child: ListView(
              padding: const EdgeInsets.all(24),
              children: [
                Text(
                  'Editing: $name',
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 24),

                ..._controllers.entries.map((entry) {
                  return FormFieldCard(
                    label: _formatLabel(entry.key),
                    child: TextFormField(
                      controller: entry.value,
                      decoration: InputDecoration(
                        hintText: 'Enter ${_formatLabel(entry.key).toLowerCase()}...',
                      ),
                    ),
                  );
                }),

                if (_error != null) ...[
                  const SizedBox(height: 12),
                  Text(
                    _error!,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.error,
                    ),
                  ),
                ],

                const SizedBox(height: 24),
                Align(
                  alignment: Alignment.centerLeft,
                  child: GradientButton(
                    onPressed: _isLoading ? null : () => _save(profile),
                    label: 'Save Changes',
                    icon: Icons.save_outlined,
                    isLoading: _isLoading,
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  String _formatLabel(String key) {
    return key
        .replaceAll('_', ' ')
        .replaceAllMapped(
          RegExp(r'(^|_)([a-z])'),
          (m) => '${m[1] ?? ''}${m[2]!.toUpperCase()}',
        )
        .replaceRange(0, 1, key[0].toUpperCase());
  }
}
