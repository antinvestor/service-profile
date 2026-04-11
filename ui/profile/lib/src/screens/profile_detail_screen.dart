import 'package:antinvestor_api_common/antinvestor_api_common.dart' as common;
import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:antinvestor_ui_core/widgets/state_badge.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/contact_providers.dart';
import '../providers/address_providers.dart';
import '../providers/profile_providers.dart';
import '../providers/relationship_providers.dart';
import '../providers/roster_providers.dart';
import '../widgets/profile_card.dart';
import '../widgets/profile_type_badge.dart';

/// Rich admin detail page for a single profile with tabs:
/// Overview | Contacts | Addresses | Relationships | Roster
class ProfileDetailScreen extends ConsumerWidget {
  const ProfileDetailScreen({super.key, required this.profileId});

  final String profileId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final asyncProfile = ref.watch(profileByIdProvider(profileId));

    return asyncProfile.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text('Failed to load profile',
                style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(friendlyError(error),
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: () =>
                  ref.invalidate(profileByIdProvider(profileId)),
              icon: const Icon(Icons.refresh, size: 18),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (profile) =>
          _ProfileDetailContent(profile: profile, profileId: profileId),
    );
  }
}

class _ProfileDetailContent extends ConsumerWidget {
  const _ProfileDetailContent({
    required this.profile,
    required this.profileId,
  });

  final ProfileObject profile;
  final String profileId;

  String get _displayName => profileName(profile);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final name = _displayName;

    return DefaultTabController(
      length: 5,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header with actions
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 24, 24, 0),
            child: Row(
              children: [
                ProfileAvatar(profileId: profile.id, name: name, size: 48),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(name,
                          style: theme.textTheme.headlineSmall
                              ?.copyWith(fontWeight: FontWeight.w600)),
                      const SizedBox(height: 4),
                      SelectableText('ID: ${profile.id}',
                          style: theme.textTheme.bodySmall?.copyWith(
                              fontFamily: 'monospace',
                              color:
                                  theme.colorScheme.onSurfaceVariant)),
                    ],
                  ),
                ),
                OutlinedButton.icon(
                  onPressed: () => context.go('/profiles'),
                  icon: const Icon(Icons.arrow_back, size: 18),
                  label: const Text('Back'),
                ),
                const SizedBox(width: 8),
                OutlinedButton.icon(
                  onPressed: () =>
                      context.go('/profiles/$profileId/edit'),
                  icon: const Icon(Icons.edit_outlined, size: 18),
                  label: const Text('Edit'),
                ),
              ],
            ),
          ),

          // Status badges
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 12, 24, 0),
            child: Wrap(
              spacing: 12,
              runSpacing: 8,
              crossAxisAlignment: WrapCrossAlignment.center,
              children: [
                StateBadge(state: common.STATE.valueOf(profile.state.value)!),
                ProfileTypeBadge(type: profile.type),
                if (profile.contacts.isNotEmpty)
                  Text(
                    '${profile.contacts.length} contact${profile.contacts.length == 1 ? '' : 's'}',
                    style: theme.textTheme.bodySmall
                        ?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                  ),
                if (profile.addresses.isNotEmpty)
                  Text(
                    '${profile.addresses.length} address${profile.addresses.length == 1 ? '' : 'es'}',
                    style: theme.textTheme.bodySmall
                        ?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                  ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // Tabs
          const TabBar(
            isScrollable: true,
            tabAlignment: TabAlignment.start,
            tabs: [
              Tab(text: 'Overview'),
              Tab(text: 'Contacts'),
              Tab(text: 'Addresses'),
              Tab(text: 'Relationships'),
              Tab(text: 'Roster'),
            ],
          ),
          const Divider(height: 1),

          // Tab views
          Expanded(
            child: TabBarView(
              children: [
                _OverviewTab(profile: profile),
                _ContactsTab(profile: profile, profileId: profileId),
                _AddressesTab(profile: profile, profileId: profileId),
                _RelationshipsTab(profileId: profileId),
                _RosterTab(profileId: profileId),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// -- Overview Tab --

class _OverviewTab extends StatelessWidget {
  const _OverviewTab({required this.profile});

  final ProfileObject profile;

  String get _displayName => profileName(profile);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          LayoutBuilder(builder: (context, constraints) {
            final wide = constraints.maxWidth > 700;
            final detailCard = _buildCard(
              context,
              title: 'Profile Details',
              icon: Icons.person_outlined,
              child: Column(children: [
                _OvRow('Name', _displayName),
                _OvRow('Type', profile.type.name),
                _OvRow('State', profile.state.name),
                _OvRow('ID', profile.id),
              ]),
            );

            final propCard = (profile.hasProperties() &&
                    profile.properties.fields.isNotEmpty)
                ? _buildCard(
                    context,
                    title: 'Properties',
                    icon: Icons.data_object,
                    child: Column(
                      children: [
                        for (final e in profile.properties.fields.entries)
                          _OvRow(e.key, _fmtValue(e.value)),
                      ],
                    ),
                  )
                : null;

            if (wide && propCard != null) {
              return Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(child: detailCard),
                  const SizedBox(width: 16),
                  Expanded(child: propCard),
                ],
              );
            }
            return Column(children: [
              detailCard,
              if (propCard != null) ...[const SizedBox(height: 16), propCard],
            ]);
          }),
          const SizedBox(height: 16),

          // Contacts summary
          if (profile.contacts.isNotEmpty)
            _buildCard(
              context,
              title: 'Contacts (${profile.contacts.length})',
              icon: Icons.contact_phone_outlined,
              child: Column(
                children: [
                  for (final c in profile.contacts)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 6),
                      child: Row(
                        children: [
                          Icon(
                            c.type == ContactType.EMAIL
                                ? Icons.email_outlined
                                : Icons.phone_outlined,
                            size: 16,
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(c.detail,
                                style: theme.textTheme.bodySmall
                                    ?.copyWith(fontWeight: FontWeight.w500)),
                          ),
                          if (c.verified)
                            const Icon(Icons.verified,
                                size: 14, color: Colors.green)
                          else
                            Icon(Icons.pending_outlined,
                                size: 14,
                                color: theme.colorScheme.onSurfaceVariant),
                        ],
                      ),
                    ),
                ],
              ),
            ),
          if (profile.contacts.isNotEmpty) const SizedBox(height: 16),

          // Addresses summary
          if (profile.addresses.isNotEmpty)
            _buildCard(
              context,
              title: 'Addresses (${profile.addresses.length})',
              icon: Icons.location_on_outlined,
              child: Column(
                children: [
                  for (final a in profile.addresses)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 6),
                      child: Row(
                        children: [
                          Icon(Icons.location_on_outlined,
                              size: 16,
                              color: theme.colorScheme.onSurfaceVariant),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              [
                                if (a.name.isNotEmpty) a.name,
                                if (a.city.isNotEmpty) a.city,
                                if (a.country.isNotEmpty) a.country,
                              ].join(', '),
                              style: theme.textTheme.bodySmall
                                  ?.copyWith(fontWeight: FontWeight.w500),
                            ),
                          ),
                        ],
                      ),
                    ),
                ],
              ),
            ),
        ],
      ),
    );
  }

  String _fmtValue(Value v) {
    if (v.hasStringValue()) return v.stringValue;
    if (v.hasBoolValue()) return v.boolValue ? 'true' : 'false';
    if (v.hasNumberValue()) return v.numberValue.toString();
    if (v.hasStructValue()) {
      return v.structValue.fields.entries
          .map((e) => '${e.key}: ${_fmtValue(e.value)}')
          .join(', ');
    }
    return '\u2014';
  }

  Widget _buildCard(BuildContext context,
      {required String title, required IconData icon, required Widget child}) {
    final theme = Theme.of(context);
    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(10),
        side: BorderSide(color: theme.colorScheme.outlineVariant),
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(icon, size: 18, color: theme.colorScheme.primary),
              const SizedBox(width: 8),
              Text(title,
                  style: theme.textTheme.titleSmall
                      ?.copyWith(fontWeight: FontWeight.w600)),
            ]),
            const SizedBox(height: 12),
            child,
          ],
        ),
      ),
    );
  }
}

class _OvRow extends StatelessWidget {
  const _OvRow(this.label, this.value);
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(label,
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ),
          Expanded(
            child: SelectableText(value,
                style: theme.textTheme.bodySmall
                    ?.copyWith(fontWeight: FontWeight.w500)),
          ),
        ],
      ),
    );
  }
}

// -- Contacts Tab --

class _ContactsTab extends ConsumerWidget {
  const _ContactsTab({required this.profile, required this.profileId});

  final ProfileObject profile;
  final String profileId;

  Future<void> _addContact(BuildContext context, WidgetRef ref) async {
    final contactController = TextEditingController();
    final result = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Add Contact'),
        content: SizedBox(
          width: 340,
          child: TextField(
            controller: contactController,
            decoration: const InputDecoration(
              labelText: 'Contact (email or phone)',
              hintText: 'e.g. user@example.com or +254712345678',
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () =>
                Navigator.of(ctx).pop(contactController.text.trim()),
            child: const Text('Add'),
          ),
        ],
      ),
    );
    contactController.dispose();
    if (result == null || result.isEmpty || !context.mounted) return;

    try {
      final notifier = ref.read(contactNotifierProvider.notifier);
      await notifier.addContact(
        AddContactRequest(id: profileId, contact: result),
      );
      ref.invalidate(profileByIdProvider(profileId));
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('Contact added')));
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  Future<void> _verifyContact(
      BuildContext context, WidgetRef ref, ContactObject contact) async {
    try {
      final notifier = ref.read(contactNotifierProvider.notifier);
      await notifier.createVerification(
        CreateContactVerificationRequest(
          id: profileId,
          contactId: contact.id,
        ),
      );
      if (!context.mounted) return;

      final codeController = TextEditingController();
      final code = await showDialog<String>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: Text('Verify ${contact.detail}'),
          content: TextField(
            controller: codeController,
            decoration: InputDecoration(
              labelText: 'Verification Code',
              hintText: 'Enter the code sent to ${contact.detail}',
            ),
          ),
          actions: [
            TextButton(
                onPressed: () => Navigator.of(ctx).pop(),
                child: const Text('Cancel')),
            FilledButton(
              onPressed: () =>
                  Navigator.of(ctx).pop(codeController.text.trim()),
              child: const Text('Verify'),
            ),
          ],
        ),
      );
      codeController.dispose();
      if (code == null || code.isEmpty || !context.mounted) return;

      final success = await notifier.checkVerification(
        CheckVerificationRequest(code: code),
      );
      ref.invalidate(profileByIdProvider(profileId));
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(
              success ? 'Contact verified successfully' : 'Verification failed'),
        ));
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  Future<void> _removeContact(
      BuildContext context, WidgetRef ref, ContactObject contact) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Remove Contact'),
        content: Text('Remove "${contact.detail}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('Remove'),
          ),
        ],
      ),
    );
    if (confirmed != true || !context.mounted) return;

    try {
      final notifier = ref.read(contactNotifierProvider.notifier);
      await notifier.removeContact(RemoveContactRequest(id: contact.id));
      ref.invalidate(profileByIdProvider(profileId));
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('Contact removed')));
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final contacts = profile.contacts;

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              ElevatedButton.icon(
                onPressed: () => _addContact(context, ref),
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add Contact'),
              ),
            ],
          ),
        ),
        if (contacts.isEmpty)
          const Expanded(
            child: Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.contact_phone_outlined,
                      size: 48, color: Colors.grey),
                  SizedBox(height: 12),
                  Text('No contacts'),
                ],
              ),
            ),
          )
        else
          Expanded(
            child: ListView.separated(
              padding: const EdgeInsets.all(16),
              itemCount: contacts.length,
              separatorBuilder: (_, __) => const Divider(height: 1),
              itemBuilder: (context, index) {
                final contact = contacts[index];
                return ExpansionTile(
                  leading: Icon(
                    contact.type == ContactType.EMAIL
                        ? Icons.email_outlined
                        : Icons.phone_outlined,
                    size: 20,
                    color: theme.colorScheme.primary,
                  ),
                  title: Text(contact.detail,
                      style: const TextStyle(fontWeight: FontWeight.w500)),
                  subtitle: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(contact.type.name,
                          style: TextStyle(
                              fontSize: 11,
                              color: theme.colorScheme.onSurfaceVariant)),
                      const SizedBox(width: 8),
                      if (contact.verified)
                        Row(mainAxisSize: MainAxisSize.min, children: [
                          const Icon(Icons.verified,
                              size: 14, color: Colors.green),
                          const SizedBox(width: 4),
                          Text('Verified',
                              style: TextStyle(
                                  fontSize: 11, color: Colors.green)),
                        ])
                      else
                        Text('Unverified',
                            style: TextStyle(
                                fontSize: 11,
                                color: theme.colorScheme.onSurfaceVariant)),
                    ],
                  ),
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (!contact.verified)
                        TextButton.icon(
                          onPressed: () =>
                              _verifyContact(context, ref, contact),
                          icon:
                              const Icon(Icons.verified_outlined, size: 16),
                          label: const Text('Verify',
                              style: TextStyle(fontSize: 12)),
                        ),
                      IconButton(
                        icon: Icon(Icons.delete_outline,
                            size: 18, color: theme.colorScheme.error),
                        tooltip: 'Remove',
                        onPressed: () =>
                            _removeContact(context, ref, contact),
                      ),
                    ],
                  ),
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          _OvRow('Contact ID', contact.id),
                          _OvRow('Type', contact.type.name),
                          _OvRow('Detail', contact.detail),
                          _OvRow('Verified',
                              contact.verified ? 'Yes' : 'No'),
                          _OvRow('State', contact.state.name),
                        ],
                      ),
                    ),
                  ],
                );
              },
            ),
          ),
      ],
    );
  }
}

// -- Addresses Tab --

class _AddressesTab extends ConsumerWidget {
  const _AddressesTab({required this.profile, required this.profileId});

  final ProfileObject profile;
  final String profileId;

  Future<void> _addAddress(BuildContext context, WidgetRef ref) async {
    final controllers = <String, TextEditingController>{
      'name': TextEditingController(),
      'country': TextEditingController(),
      'city': TextEditingController(),
      'area': TextEditingController(),
      'street': TextEditingController(),
      'house': TextEditingController(),
      'postcode': TextEditingController(),
    };

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Add Address'),
        content: SizedBox(
          width: 400,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: controllers.entries.map((e) {
                return Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: TextField(
                    controller: e.value,
                    decoration: InputDecoration(
                      labelText: e.key[0].toUpperCase() + e.key.substring(1),
                      isDense: true,
                    ),
                  ),
                );
              }).toList(),
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('Add'),
          ),
        ],
      ),
    );

    if (confirmed != true || !context.mounted) {
      for (final c in controllers.values) {
        c.dispose();
      }
      return;
    }

    try {
      await ref.read(addressNotifierProvider.notifier).addAddress(
            AddAddressRequest(
              id: profileId,
              address: AddressObject(
                name: controllers['name']!.text.trim(),
                country: controllers['country']!.text.trim(),
                city: controllers['city']!.text.trim(),
                area: controllers['area']!.text.trim(),
                street: controllers['street']!.text.trim(),
                house: controllers['house']!.text.trim(),
                postcode: controllers['postcode']!.text.trim(),
              ),
            ),
          );
      ref.invalidate(profileByIdProvider(profileId));
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('Address added')));
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    } finally {
      for (final c in controllers.values) {
        c.dispose();
      }
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final addresses = profile.addresses;

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              ElevatedButton.icon(
                onPressed: () => _addAddress(context, ref),
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add Address'),
              ),
            ],
          ),
        ),
        if (addresses.isEmpty)
          const Expanded(
            child: Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.location_on_outlined,
                      size: 48, color: Colors.grey),
                  SizedBox(height: 12),
                  Text('No addresses'),
                ],
              ),
            ),
          )
        else
          Expanded(
            child: ListView.separated(
              padding: const EdgeInsets.all(16),
              itemCount: addresses.length,
              separatorBuilder: (_, __) => const Divider(height: 1),
              itemBuilder: (context, index) {
                final addr = addresses[index];
                final parts = [
                  addr.city,
                  addr.country,
                ].where((s) => s.isNotEmpty).join(', ');
                return ExpansionTile(
                  leading: Icon(Icons.location_on_outlined,
                      size: 20, color: theme.colorScheme.primary),
                  title: Text(
                      addr.name.isNotEmpty ? addr.name : 'Address',
                      style: const TextStyle(fontWeight: FontWeight.w500)),
                  subtitle: Text(parts.isNotEmpty ? parts : '\u2014',
                      style: TextStyle(
                          fontSize: 12,
                          color: theme.colorScheme.onSurfaceVariant)),
                  trailing: addr.latitude != 0 || addr.longitude != 0
                      ? Tooltip(
                          message:
                              '${addr.latitude.toStringAsFixed(5)}, ${addr.longitude.toStringAsFixed(5)}',
                          child: Icon(Icons.my_location,
                              size: 16,
                              color: theme.colorScheme.primary),
                        )
                      : null,
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          if (addr.name.isNotEmpty)
                            _OvRow('Name', addr.name),
                          if (addr.house.isNotEmpty)
                            _OvRow('House', addr.house),
                          if (addr.street.isNotEmpty)
                            _OvRow('Street', addr.street),
                          if (addr.area.isNotEmpty)
                            _OvRow('Area', addr.area),
                          if (addr.city.isNotEmpty)
                            _OvRow('City', addr.city),
                          if (addr.country.isNotEmpty)
                            _OvRow('Country', addr.country),
                          if (addr.postcode.isNotEmpty)
                            _OvRow('Postcode', addr.postcode),
                          if (addr.latitude != 0 || addr.longitude != 0)
                            _OvRow('Coordinates',
                                '${addr.latitude.toStringAsFixed(5)}, ${addr.longitude.toStringAsFixed(5)}'),
                        ],
                      ),
                    ),
                  ],
                );
              },
            ),
          ),
      ],
    );
  }
}

// -- Relationships Tab --

class _RelationshipsTab extends ConsumerWidget {
  const _RelationshipsTab({required this.profileId});
  final String profileId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final asyncRels = ref.watch(relationshipListProvider(profileId));

    return asyncRels.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 12),
            Text(friendlyError(e)),
            const SizedBox(height: 12),
            FilledButton.tonal(
              onPressed: () =>
                  ref.invalidate(relationshipListProvider(profileId)),
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (relationships) {
        if (relationships.isEmpty) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.people_outline,
                    size: 48,
                    color: theme.colorScheme.onSurfaceVariant.withAlpha(120)),
                const SizedBox(height: 12),
                Text('No relationships', style: theme.textTheme.bodyLarge),
              ],
            ),
          );
        }
        return SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: DataTable(
              showCheckboxColumn: false,
              columns: const [
                DataColumn(label: Text('TYPE')),
                DataColumn(label: Text('PARENT')),
                DataColumn(label: Text('CHILD')),
                DataColumn(label: Text('PEER PROFILE')),
              ],
              rows: relationships.map((rel) {
                final typeLabel = rel.type.name;
                final typeColor = switch (rel.type) {
                  RelationshipType.MEMBER => theme.colorScheme.primary,
                  RelationshipType.AFFILIATED => Colors.green,
                  RelationshipType.BLACK_LISTED => theme.colorScheme.error,
                  _ => theme.colorScheme.onSurfaceVariant,
                };

                String truncate(String id) =>
                    id.length >= 8 ? id.substring(0, 8) : id;

                final parentDisplay =
                    '${rel.parentEntry.objectName}:${truncate(rel.parentEntry.objectId)}';
                final childDisplay =
                    '${rel.childEntry.objectName}:${truncate(rel.childEntry.objectId)}';
                final peerDisplay = rel.hasPeerProfile()
                    ? truncate(rel.peerProfile.id)
                    : '-';

                return DataRow(cells: [
                  DataCell(Container(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: typeColor.withAlpha(25),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Text(typeLabel,
                        style: TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w600,
                            color: typeColor)),
                  )),
                  DataCell(Text(parentDisplay,
                      style: const TextStyle(
                          fontFamily: 'monospace', fontSize: 12))),
                  DataCell(Text(childDisplay,
                      style: const TextStyle(
                          fontFamily: 'monospace', fontSize: 12))),
                  DataCell(Text(peerDisplay,
                      style: const TextStyle(
                          fontFamily: 'monospace', fontSize: 12))),
                ]);
              }).toList(),
            ),
          ),
        );
      },
    );
  }
}

// -- Roster Tab --

class _RosterTab extends ConsumerWidget {
  const _RosterTab({required this.profileId});
  final String profileId;

  Future<void> _removeEntry(
      BuildContext context, WidgetRef ref, RosterObject entry) async {
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
    if (confirmed != true || !context.mounted) return;

    try {
      await ref
          .read(rosterNotifierProvider.notifier)
          .remove(RemoveRosterRequest(id: entry.id));
      ref.invalidate(rosterSearchProvider(profileId));
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Roster entry removed')),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final asyncRoster = ref.watch(rosterSearchProvider(profileId));

    return asyncRoster.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 36, color: theme.colorScheme.error),
            const SizedBox(height: 12),
            Text('Failed to load roster'),
            const SizedBox(height: 8),
            Text(error.toString(),
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 12),
            OutlinedButton.icon(
              onPressed: () =>
                  ref.invalidate(rosterSearchProvider(profileId)),
              icon: const Icon(Icons.refresh, size: 16),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (entries) => Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
            child: Row(
              children: [
                Text('${entries.length} roster entries',
                    style: theme.textTheme.bodySmall
                        ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                const Spacer(),
                OutlinedButton.icon(
                  onPressed: () =>
                      ref.invalidate(rosterSearchProvider(profileId)),
                  icon: const Icon(Icons.refresh, size: 16),
                  label: const Text('Refresh'),
                ),
              ],
            ),
          ),
          if (entries.isEmpty)
            const Expanded(
              child: Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.contacts_outlined,
                        size: 48, color: Colors.grey),
                    SizedBox(height: 12),
                    Text('No roster entries'),
                    SizedBox(height: 4),
                    Text(
                      'Roster entries are synced from the user\'s contact book',
                      style: TextStyle(fontSize: 12, color: Colors.grey),
                    ),
                  ],
                ),
              ),
            )
          else
            Expanded(
              child: ListView.separated(
                padding: const EdgeInsets.all(16),
                itemCount: entries.length,
                separatorBuilder: (_, __) => const Divider(height: 1),
                itemBuilder: (context, index) {
                  final entry = entries[index];
                  final contact = entry.contact;
                  return ExpansionTile(
                    leading: Icon(
                      contact.type == ContactType.EMAIL
                          ? Icons.email_outlined
                          : Icons.phone_outlined,
                      size: 20,
                      color: theme.colorScheme.primary,
                    ),
                    title: Text(contact.detail,
                        style:
                            const TextStyle(fontWeight: FontWeight.w500)),
                    subtitle: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(contact.type.name,
                            style: TextStyle(
                                fontSize: 11,
                                color:
                                    theme.colorScheme.onSurfaceVariant)),
                        const SizedBox(width: 8),
                        if (contact.verified)
                          const Icon(Icons.verified,
                              size: 14, color: Colors.green)
                        else
                          Icon(Icons.pending_outlined,
                              size: 14,
                              color: theme.colorScheme.onSurfaceVariant),
                      ],
                    ),
                    trailing: IconButton(
                      icon: Icon(Icons.delete_outline,
                          size: 18, color: theme.colorScheme.error),
                      tooltip: 'Remove',
                      onPressed: () =>
                          _removeEntry(context, ref, entry),
                    ),
                    children: [
                      Padding(
                        padding:
                            const EdgeInsets.fromLTRB(16, 0, 16, 12),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            _OvRow('Roster ID', entry.id),
                            _OvRow('Profile ID', entry.profileId),
                            _OvRow('Contact', contact.detail),
                            _OvRow('Type', contact.type.name),
                            _OvRow('Verified',
                                contact.verified ? 'Yes' : 'No'),
                          ],
                        ),
                      ),
                    ],
                  );
                },
              ),
            ),
        ],
      ),
    );
  }
}
