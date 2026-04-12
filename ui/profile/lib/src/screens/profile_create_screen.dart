import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:antinvestor_ui_core/widgets/gradient_button.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/profile_providers.dart';

/// Screen for creating a new profile.
class ProfileCreateScreen extends ConsumerStatefulWidget {
  const ProfileCreateScreen({super.key});

  @override
  ConsumerState<ProfileCreateScreen> createState() =>
      _ProfileCreateScreenState();
}

class _ProfileCreateScreenState extends ConsumerState<ProfileCreateScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _contactController = TextEditingController();
  final _descriptionController = TextEditingController();
  ProfileType _type = ProfileType.PERSON;
  bool _isLoading = false;
  String? _error;

  @override
  void dispose() {
    _nameController.dispose();
    _contactController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _create() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final props = Struct();
      final name = _nameController.text.trim();
      if (name.isNotEmpty) {
        props.fields['name'] = Value()..stringValue = name;
      }
      final description = _descriptionController.text.trim();
      if (description.isNotEmpty) {
        props.fields['description'] = Value()..stringValue = description;
      }

      final request = CreateRequest()
        ..type = _type
        ..contact = _contactController.text.trim()
        ..properties = props;

      final created =
          await ref.read(profileNotifierProvider.notifier).create(request);

      if (mounted) {
        context.go('/profiles/${created.id}');
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
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('New Profile')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(24),
          children: [
            Text(
              'Create a new profile',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 24),

            FormFieldCard(
              label: 'Profile Type',
              isRequired: true,
              child: DropdownButtonFormField<ProfileType>(
                initialValue: _type,
                decoration: const InputDecoration(
                  hintText: 'Select type...',
                ),
                items: const [
                  DropdownMenuItem(
                    value: ProfileType.PERSON,
                    child: Text('Person'),
                  ),
                  DropdownMenuItem(
                    value: ProfileType.INSTITUTION,
                    child: Text('Institution'),
                  ),
                  DropdownMenuItem(
                    value: ProfileType.BOT,
                    child: Text('Bot'),
                  ),
                ],
                onChanged: (type) {
                  if (type != null) setState(() => _type = type);
                },
              ),
            ),

            FormFieldCard(
              label: 'Name',
              isRequired: true,
              child: TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(
                  hintText: 'Enter profile name...',
                ),
                validator: (v) =>
                    v == null || v.trim().isEmpty ? 'Name is required' : null,
              ),
            ),

            FormFieldCard(
              label: 'Contact',
              description: 'Email or phone number for the profile.',
              isRequired: true,
              child: TextFormField(
                controller: _contactController,
                decoration: const InputDecoration(
                  hintText: 'Enter email or phone...',
                ),
                validator: (v) =>
                    v == null || v.trim().isEmpty
                        ? 'Contact is required'
                        : null,
              ),
            ),

            FormFieldCard(
              label: 'Description',
              child: TextFormField(
                controller: _descriptionController,
                decoration: const InputDecoration(
                  hintText: 'Optional description...',
                ),
                maxLines: 3,
              ),
            ),

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
                onPressed: _isLoading ? null : _create,
                label: 'Create Profile',
                icon: Icons.person_add_outlined,
                isLoading: _isLoading,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
