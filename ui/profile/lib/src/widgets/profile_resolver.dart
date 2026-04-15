import 'dart:typed_data';

import 'package:antinvestor_api_profile/antinvestor_api_profile.dart' as profile;
import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

import '../providers/profile_providers.dart';
import '../providers/profile_transport_provider.dart';

/// Reusable profile resolution widget.
///
/// Implements the contact search → find/create profile flow that is
/// common across organization, org unit, investor, client, and workforce
/// member creation.
///
/// Usage:
/// ```dart
/// ProfileResolver(
///   profileType: ProfileType.INSTITUTION, // or PERSON
///   onProfileResolved: (profile) {
///     setState(() => _profile = profile);
///   },
///   onPickAvatar: (bytes, filename) => uploadPublicImage(ref, bytes, filename),
/// )
/// ```
class ProfileResolver extends ConsumerStatefulWidget {
  const ProfileResolver({
    super.key,
    required this.onProfileResolved,
    this.profileType = profile.ProfileType.INSTITUTION,
    this.onPickAvatar,
    this.resolvedProfile,
    this.contactLabel = 'Enter contact (email or phone)',
    this.nameLabel = 'Name',
    this.nameHint = 'e.g. Stawi Capital Limited',
  });

  /// Called when a profile is found or created.
  final ValueChanged<profile.ProfileObject> onProfileResolved;

  /// Type of profile to create if not found.
  final profile.ProfileType profileType;

  /// Callback to upload avatar. Returns content URI.
  final Future<String> Function(Uint8List bytes, String filename)? onPickAvatar;

  /// Already resolved profile (for edit mode).
  final profile.ProfileObject? resolvedProfile;

  /// Label for the contact search field.
  final String contactLabel;

  /// Label for the name field when creating.
  final String nameLabel;

  /// Hint for the name field when creating.
  final String nameHint;

  @override
  ConsumerState<ProfileResolver> createState() => _ProfileResolverState();
}

class _ProfileResolverState extends ConsumerState<ProfileResolver> {
  final _contactCtrl = TextEditingController();
  final _nameCtrl = TextEditingController();
  final _descriptionCtrl = TextEditingController();
  Uint8List? _avatarBytes;
  String? _avatarFileName;

  profile.ProfileObject? _foundProfile;
  bool _searched = false;
  bool _searching = false;
  bool _creating = false;

  @override
  void initState() {
    super.initState();
    if (widget.resolvedProfile != null) {
      _foundProfile = widget.resolvedProfile;
      _searched = true;
      _nameCtrl.text = _extractName(widget.resolvedProfile!);
      _descriptionCtrl.text = _extractProp(widget.resolvedProfile!, 'description');
    }
  }

  @override
  void dispose() {
    _contactCtrl.dispose();
    _nameCtrl.dispose();
    _descriptionCtrl.dispose();
    super.dispose();
  }

  String _extractName(profile.ProfileObject p) {
    try {
      if (p.properties.fields.containsKey('name')) {
        return p.properties.fields['name']!.stringValue;
      }
    } catch (_) {}
    if (p.contacts.isNotEmpty) return p.contacts.first.detail;
    return '';
  }

  String _extractProp(profile.ProfileObject p, String key) {
    try {
      if (p.properties.fields.containsKey(key)) {
        return p.properties.fields[key]!.stringValue;
      }
    } catch (_) {}
    return '';
  }

  Future<void> _searchProfile() async {
    final contact = _contactCtrl.text.trim();
    if (contact.isEmpty) return;
    setState(() {
      _searching = true;
      _searched = false;
      _foundProfile = null;
    });
    try {
      // Force a fresh lookup (don't use cached provider state).
      final client = ref.read(profileServiceClientProvider);
      final request = profile.GetByContactRequest()..contact = contact;
      final response = await client.getByContact(request);
      final result = response.data;
      if (mounted) {
        setState(() {
          _foundProfile = result;
          _searched = true;
          _searching = false;
          _nameCtrl.text = _extractName(result);
          _descriptionCtrl.text = _extractProp(result, 'description');
        });
        widget.onProfileResolved(result);
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _foundProfile = null;
          _searched = true;
          _searching = false;
        });
      }
    }
  }

  Future<void> _createProfile() async {
    final name = _nameCtrl.text.trim();
    if (name.isEmpty) return;
    setState(() => _creating = true);
    try {
      String? avatarUri;
      if (_avatarBytes != null &&
          _avatarFileName != null &&
          widget.onPickAvatar != null) {
        avatarUri =
            await widget.onPickAvatar!(_avatarBytes!, _avatarFileName!);
      }

      final props = <String, profile.Value>{};
      props['name'] = profile.Value(stringValue: name);
      if (_descriptionCtrl.text.trim().isNotEmpty) {
        props['description'] =
            profile.Value(stringValue: _descriptionCtrl.text.trim());
      }
      if (avatarUri != null) {
        props['avatar'] = profile.Value(stringValue: avatarUri);
      }

      final created =
          await ref.read(profileNotifierProvider.notifier).create(
                profile.CreateRequest(
                  type: widget.profileType,
                  contact: _contactCtrl.text.trim(),
                  properties: profile.Struct(fields: props),
                ),
              );

      if (mounted) {
        setState(() {
          _foundProfile = created;
          _creating = false;
        });
        widget.onProfileResolved(created);
      }
    } catch (e) {
      if (mounted) {
        setState(() => _creating = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to create profile: $e')),
        );
      }
    }
  }

  Future<void> _pickAvatar() async {
    final picker = ImagePicker();
    final picked = await picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 512,
      maxHeight: 512,
      imageQuality: 85,
    );
    if (picked == null) return;
    final bytes = await picked.readAsBytes();
    setState(() {
      _avatarBytes = bytes;
      _avatarFileName = picked.name;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Search field
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _contactCtrl,
                decoration: InputDecoration(
                  hintText: widget.contactLabel,
                  prefixIcon: const Icon(Icons.contact_mail_outlined),
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8)),
                ),
                onSubmitted: (_) => _searchProfile(),
              ),
            ),
            const SizedBox(width: 8),
            FilledButton.icon(
              onPressed: _searching ? null : _searchProfile,
              icon: _searching
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(
                          strokeWidth: 2, color: Colors.white),
                    )
                  : const Icon(Icons.search, size: 18),
              label: const Text('Search'),
            ),
          ],
        ),

        const SizedBox(height: 16),

        // Profile found
        if (_searched && _foundProfile != null)
          Card(
            elevation: 0,
            color: theme.colorScheme.primaryContainer.withAlpha(40),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
              side: BorderSide(
                  color: theme.colorScheme.primary.withAlpha(60)),
            ),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Icon(Icons.check_circle,
                      color: theme.colorScheme.primary, size: 24),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('Profile found',
                            style: theme.textTheme.titleSmall?.copyWith(
                                fontWeight: FontWeight.w600,
                                color: theme.colorScheme.primary)),
                        ProfileBadge(
                          profileId: _foundProfile!.id,
                          name: _extractName(_foundProfile!),
                          avatarSize: 32,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),

        // Profile not found — create form
        if (_searched && _foundProfile == null)
          Card(
            elevation: 0,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
              side: BorderSide(
                  color: theme.colorScheme.outlineVariant.withAlpha(60)),
            ),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(Icons.info_outline,
                          color: theme.colorScheme.secondary, size: 20),
                      const SizedBox(width: 8),
                      Text('No profile found. Create one below.',
                          style: theme.textTheme.bodyMedium?.copyWith(
                              color: theme.colorScheme.secondary)),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Center(
                    child: GestureDetector(
                      onTap: _pickAvatar,
                      child: CircleAvatar(
                        radius: 40,
                        backgroundColor:
                            theme.colorScheme.primaryContainer,
                        backgroundImage: _avatarBytes != null
                            ? MemoryImage(_avatarBytes!)
                            : null,
                        child: _avatarBytes == null
                            ? Icon(Icons.add_a_photo,
                                color: theme.colorScheme.primary,
                                size: 24)
                            : null,
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextField(
                    controller: _nameCtrl,
                    decoration: InputDecoration(
                      labelText: widget.nameLabel,
                      hintText: widget.nameHint,
                      prefixIcon: const Icon(Icons.business),
                      border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8)),
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextField(
                    controller: _descriptionCtrl,
                    decoration: InputDecoration(
                      labelText: 'Description',
                      hintText: 'Brief description',
                      prefixIcon:
                          const Icon(Icons.description_outlined),
                      border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8)),
                    ),
                    maxLines: 2,
                  ),
                  const SizedBox(height: 12),
                  FilledButton.icon(
                    onPressed: _creating ? null : _createProfile,
                    icon: _creating
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                                strokeWidth: 2, color: Colors.white),
                          )
                        : const Icon(Icons.person_add, size: 18),
                    label: const Text('Create Profile'),
                  ),
                ],
              ),
            ),
          ),
      ],
    );
  }
}
