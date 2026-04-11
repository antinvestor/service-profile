import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/roster_providers.dart';
import '../widgets/profile_badge_by_id.dart';

/// Admin-grade roster screen with search-by-profile, DataTable results,
/// and remove actions.
class RosterScreen extends ConsumerStatefulWidget {
  const RosterScreen({super.key, required this.profileId});

  final String profileId;

  @override
  ConsumerState<RosterScreen> createState() => _RosterScreenState();
}

class _RosterScreenState extends ConsumerState<RosterScreen> {
  final _profileIdCtl = TextEditingController();
  final _queryCtl = TextEditingController();
  String _activeProfileId = '';

  @override
  void initState() {
    super.initState();
    _profileIdCtl.text = widget.profileId;
    _activeProfileId = widget.profileId;
  }

  @override
  void dispose() {
    _profileIdCtl.dispose();
    _queryCtl.dispose();
    super.dispose();
  }

  void _search() {
    final profileId = _profileIdCtl.text.trim();
    if (profileId.isEmpty) return;
    setState(() {
      _activeProfileId = profileId;
    });
  }

  Future<void> _removeEntry(RosterObject entry) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Remove Roster Entry'),
        content: Text('Remove "${entry.contact.detail}" from roster?'),
        actions: [
          TextButton(
              onPressed: () => Navigator.of(ctx).pop(false),
              child: const Text('Cancel')),
          FilledButton(
              onPressed: () => Navigator.of(ctx).pop(true),
              child: const Text('Remove')),
        ],
      ),
    );
    if (confirmed != true || !mounted) return;

    try {
      await ref
          .read(rosterNotifierProvider.notifier)
          .remove(RemoveRosterRequest(id: entry.id));
      ref.invalidate(rosterSearchProvider(_activeProfileId));
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Roster entry removed')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Row(
            children: [
              Icon(Icons.contacts_outlined,
                  size: 28, color: theme.colorScheme.primary),
              const SizedBox(width: 12),
              Expanded(
                child: Text('Roster',
                    style: theme.textTheme.headlineSmall
                        ?.copyWith(fontWeight: FontWeight.w600)),
              ),
            ],
          ),
          const SizedBox(height: 20),

          // Search bar
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: theme.colorScheme.outlineVariant),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Search Roster',
                    style: theme.textTheme.titleSmall
                        ?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _profileIdCtl,
                        decoration: const InputDecoration(
                          labelText: 'Profile ID',
                          hintText: 'Enter profile ID...',
                          prefixIcon: Icon(Icons.person_outlined, size: 20),
                          isDense: true,
                        ),
                        onSubmitted: (_) => _search(),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: TextField(
                        controller: _queryCtl,
                        decoration: const InputDecoration(
                          labelText: 'Search Query (optional)',
                          hintText: 'Filter by name or contact...',
                          prefixIcon: Icon(Icons.search, size: 20),
                          isDense: true,
                        ),
                        onSubmitted: (_) => _search(),
                      ),
                    ),
                    const SizedBox(width: 12),
                    ElevatedButton.icon(
                      onPressed: _search,
                      icon: const Icon(Icons.search, size: 18),
                      label: const Text('Search'),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // Results
          if (_activeProfileId.isNotEmpty)
            _RosterResults(
              profileId: _activeProfileId,
              onRemove: _removeEntry,
            ),
        ],
      ),
    );
  }
}

class _RosterResults extends ConsumerWidget {
  const _RosterResults({
    required this.profileId,
    required this.onRemove,
  });

  final String profileId;
  final void Function(RosterObject entry) onRemove;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final asyncRoster = ref.watch(rosterSearchProvider(profileId));

    return asyncRoster.when(
      loading: () => const Padding(
        padding: EdgeInsets.all(48),
        child: Center(child: CircularProgressIndicator()),
      ),
      error: (error, _) => Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: theme.colorScheme.errorContainer.withAlpha(30),
          borderRadius: BorderRadius.circular(10),
        ),
        child: Row(
          children: [
            Icon(Icons.error_outline,
                size: 20, color: theme.colorScheme.error),
            const SizedBox(width: 8),
            Expanded(child: Text(friendlyError(error))),
          ],
        ),
      ),
      data: (entries) {
        if (entries.isEmpty) {
          return Center(
            child: Padding(
              padding: const EdgeInsets.all(48),
              child: Column(
                children: [
                  Icon(Icons.contacts_outlined,
                      size: 48,
                      color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(height: 12),
                  Text('No roster entries found',
                      style: TextStyle(
                          color: theme.colorScheme.onSurfaceVariant)),
                ],
              ),
            ),
          );
        }

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('${entries.length} roster entries found',
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 8),
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(12),
                border:
                    Border.all(color: theme.colorScheme.outlineVariant),
              ),
              child: Column(
                children: [
                  for (int i = 0; i < entries.length; i++) ...[
                    if (i > 0) const Divider(height: 1),
                    _RosterEntryTile(
                      entry: entries[i],
                      onRemove: () => onRemove(entries[i]),
                    ),
                  ],
                ],
              ),
            ),
          ],
        );
      },
    );
  }
}

class _RosterEntryTile extends StatelessWidget {
  const _RosterEntryTile({required this.entry, required this.onRemove});

  final RosterObject entry;
  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final contact = entry.contact;
    final isEmail = contact.type == ContactType.EMAIL;

    return ListTile(
      leading: Icon(
        isEmail ? Icons.email_outlined : Icons.phone_outlined,
        size: 20,
        color: theme.colorScheme.primary,
      ),
      title: Text(contact.detail,
          style: const TextStyle(fontWeight: FontWeight.w500)),
      subtitle: ProfileBadgeById(
          profileId: entry.profileId, avatarSize: 20),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (contact.verified)
            const Icon(Icons.verified, size: 16, color: Colors.green)
          else
            Icon(Icons.pending_outlined,
                size: 16, color: theme.colorScheme.onSurfaceVariant),
          const SizedBox(width: 8),
          IconButton(
            icon: Icon(Icons.delete_outline,
                size: 18, color: theme.colorScheme.error),
            tooltip: 'Remove',
            onPressed: onRemove,
          ),
        ],
      ),
    );
  }
}
