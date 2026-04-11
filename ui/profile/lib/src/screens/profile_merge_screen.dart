import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:antinvestor_ui_core/widgets/gradient_button.dart';
import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:antinvestor_ui_core/widgets/state_badge.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/profile_providers.dart';
import '../widgets/profile_card.dart';
import '../widgets/profile_type_badge.dart';

/// Screen for merging two profiles into one.
class ProfileMergeScreen extends ConsumerStatefulWidget {
  const ProfileMergeScreen({super.key});

  @override
  ConsumerState<ProfileMergeScreen> createState() =>
      _ProfileMergeScreenState();
}

class _ProfileMergeScreenState extends ConsumerState<ProfileMergeScreen> {
  final _primaryController = TextEditingController();
  final _mergeController = TextEditingController();
  ProfileObject? _primaryProfile;
  ProfileObject? _mergeProfile;
  bool _isSearchingPrimary = false;
  bool _isSearchingMerge = false;
  bool _isMerging = false;
  String? _error;

  @override
  void dispose() {
    _primaryController.dispose();
    _mergeController.dispose();
    super.dispose();
  }

  Future<void> _searchPrimary() async {
    final query = _primaryController.text.trim();
    if (query.isEmpty) return;

    setState(() {
      _isSearchingPrimary = true;
      _error = null;
    });

    try {
      final results =
          await ref.read(profileSearchProvider(query).future);
      if (results.isNotEmpty) {
        setState(() => _primaryProfile = results.first);
      } else {
        setState(() => _error = 'No profile found for "$query"');
      }
    } catch (e) {
      setState(() => _error = friendlyError(e));
    } finally {
      setState(() => _isSearchingPrimary = false);
    }
  }

  Future<void> _searchMerge() async {
    final query = _mergeController.text.trim();
    if (query.isEmpty) return;

    setState(() {
      _isSearchingMerge = true;
      _error = null;
    });

    try {
      final results =
          await ref.read(profileSearchProvider(query).future);
      if (results.isNotEmpty) {
        setState(() => _mergeProfile = results.first);
      } else {
        setState(() => _error = 'No profile found for "$query"');
      }
    } catch (e) {
      setState(() => _error = friendlyError(e));
    } finally {
      setState(() => _isSearchingMerge = false);
    }
  }

  Future<void> _merge() async {
    if (_primaryProfile == null || _mergeProfile == null) return;
    if (_primaryProfile!.id == _mergeProfile!.id) {
      setState(() => _error = 'Cannot merge a profile with itself.');
      return;
    }

    setState(() {
      _isMerging = true;
      _error = null;
    });

    try {
      final request = MergeRequest()
        ..id = _primaryProfile!.id
        ..mergeid = _mergeProfile!.id;

      final result = await ref
          .read(profileNotifierProvider.notifier)
          .merge(request);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Profiles merged successfully'),
            behavior: SnackBarBehavior.floating,
          ),
        );
        context.go('/profiles/${result.id}');
      }
    } catch (e) {
      setState(() {
        _error = friendlyError(e);
        _isMerging = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Merge Profiles')),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
          Text(
            'Merge two profiles into one',
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'The secondary profile will be merged into the primary profile. '
            'All contacts, addresses, and relationships from the secondary '
            'profile will be transferred.',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
          const SizedBox(height: 24),

          // Primary profile search
          FormFieldCard(
            label: 'Primary Profile',
            description: 'Search by name, email, or phone.',
            isRequired: true,
            child: Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _primaryController,
                    decoration: const InputDecoration(
                      hintText: 'Search primary profile...',
                    ),
                    onFieldSubmitted: (_) => _searchPrimary(),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  onPressed:
                      _isSearchingPrimary ? null : _searchPrimary,
                  icon: _isSearchingPrimary
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                              strokeWidth: 2),
                        )
                      : const Icon(Icons.search),
                ),
              ],
            ),
          ),
          if (_primaryProfile != null) _profilePreview(_primaryProfile!),

          const SizedBox(height: 8),

          // Merge profile search
          FormFieldCard(
            label: 'Secondary Profile (to merge)',
            description: 'This profile will be merged into the primary.',
            isRequired: true,
            child: Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _mergeController,
                    decoration: const InputDecoration(
                      hintText: 'Search secondary profile...',
                    ),
                    onFieldSubmitted: (_) => _searchMerge(),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  onPressed: _isSearchingMerge ? null : _searchMerge,
                  icon: _isSearchingMerge
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                              strokeWidth: 2),
                        )
                      : const Icon(Icons.search),
                ),
              ],
            ),
          ),
          if (_mergeProfile != null) _profilePreview(_mergeProfile!),

          if (_error != null) ...[
            const SizedBox(height: 16),
            Text(
              _error!,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.error,
              ),
            ),
          ],

          const SizedBox(height: 32),
          if (_primaryProfile != null && _mergeProfile != null)
            Align(
              alignment: Alignment.centerLeft,
              child: GradientButton(
                onPressed: _isMerging ? null : _merge,
                label: 'Confirm Merge',
                icon: Icons.merge_type,
                isLoading: _isMerging,
              ),
            ),
        ],
      ),
    );
  }

  Widget _profilePreview(ProfileObject profile) {
    final theme = Theme.of(context);
    final name = profileName(profile);

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: theme.colorScheme.outlineVariant),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            ProfileAvatar(profileId: profile.id, name: name, size: 48),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    name,
                    style: theme.textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      ProfileTypeBadge(type: profile.type),
                      const SizedBox(width: 8),
                      StateBadge(state: profile.state),
                    ],
                  ),
                  if (profile.contacts.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    Text(
                      '${profile.contacts.length} contacts, '
                      '${profile.addresses.length} addresses',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
