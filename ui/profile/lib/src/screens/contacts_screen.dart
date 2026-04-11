import 'package:antinvestor_api_common/antinvestor_api_common.dart' as common;
import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/state_badge.dart';
import 'package:csv/csv.dart' show CsvEncoder;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/profile_providers.dart';
import '../widgets/profile_badge_by_id.dart';

/// A contact entry paired with the owning profile's ID.
typedef ProfileContact = ({String profileId, ContactObject contact});

/// Provider that extracts all contacts from loaded profiles into a flat list.
final allContactsProvider =
    FutureProvider<List<ProfileContact>>((ref) async {
  final profiles = await ref.watch(profileSearchProvider('').future);
  final contacts = <ProfileContact>[];
  for (final profile in profiles) {
    for (final contact in profile.contacts) {
      contacts.add((profileId: profile.id, contact: contact));
    }
  }
  return contacts;
});

/// Admin-grade contacts screen with DataTable, verification actions,
/// CSV export, search, pagination, and detail panel.
class ContactsScreen extends ConsumerStatefulWidget {
  const ContactsScreen({super.key, this.profileId});

  /// When provided, only shows contacts for this profile.
  final String? profileId;

  @override
  ConsumerState<ContactsScreen> createState() => _ContactsScreenState();
}

class _ContactsScreenState extends ConsumerState<ContactsScreen> {
  int? _selectedIndex;
  int _currentPage = 0;
  int _pageSize = 25;
  String _searchQuery = '';
  final _searchController = TextEditingController();

  static const _pageSizeOptions = [10, 25, 50, 100];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _exportCsv(List<ProfileContact> contacts) {
    final headers = ['Detail', 'Type', 'Verified', 'Profile ID', 'State'];
    final rows = contacts
        .map((e) => [
              e.contact.detail,
              e.contact.type.name,
              e.contact.verified ? 'Yes' : 'No',
              e.profileId,
              e.contact.state.name,
            ])
        .toList();
    final csv = const CsvEncoder().convert([headers, ...rows]);
    Clipboard.setData(ClipboardData(text: csv));
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
              'Contacts CSV copied to clipboard (${rows.length} rows)'),
        ),
      );
    }
  }

  List<ProfileContact> _filter(List<ProfileContact> items) {
    if (_searchQuery.isEmpty) return items;
    final q = _searchQuery.toLowerCase();
    return items.where((e) {
      return e.contact.detail.toLowerCase().contains(q) ||
          e.contact.type.name.toLowerCase().contains(q) ||
          e.profileId.toLowerCase().contains(q);
    }).toList();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final asyncContacts = widget.profileId != null
        ? ref.watch(profileByIdProvider(widget.profileId!)).whenData((p) =>
            p.contacts
                .map((c) =>
                    (profileId: widget.profileId!, contact: c))
                .toList())
        : ref.watch(allContactsProvider);

    return asyncContacts.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text('Failed to load contacts',
                style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(friendlyError(error),
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: () {
                if (widget.profileId != null) {
                  ref.invalidate(profileByIdProvider(widget.profileId!));
                } else {
                  ref.invalidate(allContactsProvider);
                }
              },
              icon: const Icon(Icons.refresh, size: 18),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (allContacts) {
        final contacts = _filter(allContacts);
        return _buildContent(context, theme, contacts, allContacts.length);
      },
    );
  }

  Widget _buildContent(BuildContext context, ThemeData theme,
      List<ProfileContact> contacts, int totalUnfiltered) {
    if (_selectedIndex != null && _selectedIndex! >= contacts.length) {
      _selectedIndex = null;
    }

    final showDetail = _selectedIndex != null;
    final totalPages = (contacts.length / _pageSize).ceil();
    final pageStart = _currentPage * _pageSize;
    final pageEnd = (pageStart + _pageSize).clamp(0, contacts.length);
    final pageItems =
        contacts.isNotEmpty ? contacts.sublist(pageStart, pageEnd) : [];

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          flex: showDetail ? 3 : 1,
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Header
                Row(
                  children: [
                    Icon(Icons.contact_phone_outlined,
                        size: 28, color: theme.colorScheme.primary),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Contacts',
                              style: theme.textTheme.headlineSmall
                                  ?.copyWith(fontWeight: FontWeight.w600)),
                          Text('$totalUnfiltered contacts',
                              style: theme.textTheme.bodySmall?.copyWith(
                                  color:
                                      theme.colorScheme.onSurfaceVariant)),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),

                // Search + export
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _searchController,
                        onChanged: (v) => setState(() {
                          _searchQuery = v.trim();
                          _currentPage = 0;
                          _selectedIndex = null;
                        }),
                        decoration: const InputDecoration(
                          hintText: 'Search contacts...',
                          prefixIcon: Icon(Icons.search, size: 20),
                          isDense: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    OutlinedButton.icon(
                      onPressed: () => _exportCsv(contacts),
                      icon: const Icon(Icons.download, size: 18),
                      label: const Text('Export'),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // DataTable
                Container(
                  width: double.infinity,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(12),
                    border:
                        Border.all(color: theme.colorScheme.outlineVariant),
                  ),
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: SingleChildScrollView(
                      scrollDirection: Axis.horizontal,
                      child: DataTable(
                        showCheckboxColumn: false,
                        columns: const [
                          DataColumn(label: Text('DETAIL')),
                          DataColumn(label: Text('TYPE')),
                          DataColumn(label: Text('VERIFIED')),
                          DataColumn(label: Text('PROFILE')),
                          DataColumn(label: Text('STATE')),
                        ],
                        rows: List.generate(pageItems.length, (i) {
                          final entry = pageItems[i];
                          final contact = entry.contact;
                          final globalIndex = pageStart + i;
                          final typeColor =
                              contact.type == ContactType.EMAIL
                                  ? theme.colorScheme.primary
                                  : Colors.green;

                          return DataRow(
                            selected: _selectedIndex == globalIndex,
                            onSelectChanged: (_) =>
                                setState(() => _selectedIndex = globalIndex),
                            color: WidgetStateProperty.resolveWith(
                                (states) {
                              if (states.contains(WidgetState.selected)) {
                                return theme.colorScheme.primaryContainer
                                    .withAlpha(40);
                              }
                              return null;
                            }),
                            cells: [
                              DataCell(Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Icon(
                                    contact.type == ContactType.EMAIL
                                        ? Icons.email_outlined
                                        : Icons.phone_outlined,
                                    size: 16,
                                    color: typeColor,
                                  ),
                                  const SizedBox(width: 8),
                                  Text(contact.detail),
                                ],
                              )),
                              DataCell(_ColorBadge(
                                  contact.type.name, typeColor)),
                              DataCell(Icon(
                                contact.verified
                                    ? Icons.check_circle
                                    : Icons.cancel_outlined,
                                size: 18,
                                color: contact.verified
                                    ? Colors.green
                                    : theme.colorScheme.onSurfaceVariant,
                              )),
                              DataCell(ProfileBadgeById(
                                profileId: entry.profileId,
                                avatarSize: 24,
                              )),
                              DataCell(
                                  StateBadge(state: common.STATE.valueOf(contact.state.value)!)),
                            ],
                          );
                        }),
                      ),
                    ),
                  ),
                ),

                // Pagination
                const SizedBox(height: 8),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        Text('Rows per page:',
                            style: theme.textTheme.bodySmall),
                        const SizedBox(width: 8),
                        DropdownButton<int>(
                          value: _pageSize,
                          underline: const SizedBox(),
                          items: _pageSizeOptions
                              .map((s) => DropdownMenuItem(
                                  value: s, child: Text('$s')))
                              .toList(),
                          onChanged: (v) {
                            if (v != null) {
                              setState(() {
                                _pageSize = v;
                                _currentPage = 0;
                              });
                            }
                          },
                        ),
                      ],
                    ),
                    Row(
                      children: [
                        Text(
                          contacts.isEmpty
                              ? '0 of 0'
                              : '${pageStart + 1}-$pageEnd of ${contacts.length}',
                          style: theme.textTheme.bodySmall,
                        ),
                        IconButton(
                          icon: const Icon(Icons.chevron_left, size: 20),
                          onPressed: _currentPage > 0
                              ? () => setState(() => _currentPage--)
                              : null,
                        ),
                        IconButton(
                          icon: const Icon(Icons.chevron_right, size: 20),
                          onPressed: _currentPage < totalPages - 1
                              ? () => setState(() => _currentPage++)
                              : null,
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),

        // Detail panel
        if (showDetail && _selectedIndex! < contacts.length)
          SizedBox(
            width: 380,
            child: Container(
              margin: const EdgeInsets.only(
                  top: 24, right: 24, bottom: 24),
              decoration: BoxDecoration(
                color: theme.colorScheme.surface,
                borderRadius: BorderRadius.circular(12),
                border:
                    Border.all(color: theme.colorScheme.outlineVariant),
              ),
              child: Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.all(12),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        IconButton(
                          icon: const Icon(Icons.close, size: 20),
                          onPressed: () =>
                              setState(() => _selectedIndex = null),
                        ),
                      ],
                    ),
                  ),
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.fromLTRB(20, 0, 20, 20),
                      child: _ContactDetail(
                          entry: contacts[_selectedIndex!]),
                    ),
                  ),
                ],
              ),
            ),
          ),
      ],
    );
  }
}

class _ContactDetail extends StatelessWidget {
  const _ContactDetail({required this.entry});

  final ProfileContact entry;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final contact = entry.contact;
    final typeColor = contact.type == ContactType.EMAIL
        ? theme.colorScheme.primary
        : Colors.green;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            CircleAvatar(
              radius: 20,
              backgroundColor: typeColor.withAlpha(30),
              child: Icon(
                contact.type == ContactType.EMAIL
                    ? Icons.email_outlined
                    : Icons.phone_outlined,
                color: typeColor,
                size: 20,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(contact.detail,
                      style: theme.textTheme.titleMedium
                          ?.copyWith(fontWeight: FontWeight.w600)),
                  Text(contact.id,
                      style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                          fontFamily: 'monospace')),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 20),
        _DetailRow(label: 'Type', value: contact.type.name),
        _DetailRow(
            label: 'Verified', value: contact.verified ? 'Yes' : 'No'),
        _DetailRow(label: 'State', value: contact.state.name),
        const SizedBox(height: 8),
        Text('Profile',
            style: theme.textTheme.labelMedium
                ?.copyWith(fontWeight: FontWeight.w600)),
        const SizedBox(height: 4),
        ProfileBadgeById(profileId: entry.profileId),
      ],
    );
  }
}

class _ColorBadge extends StatelessWidget {
  const _ColorBadge(this.label, this.color);
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withAlpha(25),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(label,
          style: TextStyle(
              fontSize: 11, fontWeight: FontWeight.w600, color: color)),
    );
  }
}

class _DetailRow extends StatelessWidget {
  const _DetailRow({required this.label, required this.value});
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(label,
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ),
          Expanded(
            child: Text(value,
                style: theme.textTheme.bodySmall
                    ?.copyWith(fontWeight: FontWeight.w500)),
          ),
        ],
      ),
    );
  }
}
