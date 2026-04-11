import 'package:antinvestor_api_common/antinvestor_api_common.dart' as common;
import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:antinvestor_ui_core/widgets/state_badge.dart';
import 'package:csv/csv.dart' show CsvEncoder;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/profile_providers.dart';
import '../widgets/profile_card.dart';
import '../widgets/profile_type_badge.dart';

/// Admin-grade profile search screen with DataTable, CSV export, search,
/// pagination, row selection with detail panel, and create/edit actions.
class ProfileSearchScreen extends ConsumerStatefulWidget {
  const ProfileSearchScreen({
    super.key,
    this.onProfileSelected,
    this.onNavigateToDetail,
  });

  /// Optional callback when a profile is selected in the table.
  final void Function(ProfileObject profile)? onProfileSelected;

  /// Optional callback to navigate to a profile detail page.
  /// If null, uses `context.go('/profiles/$id')`.
  final void Function(ProfileObject profile)? onNavigateToDetail;

  @override
  ConsumerState<ProfileSearchScreen> createState() =>
      _ProfileSearchScreenState();
}

class _ProfileSearchScreenState extends ConsumerState<ProfileSearchScreen> {
  String _query = '';
  ProfileType? _selectedType;
  int? _selectedIndex;
  int _currentPage = 0;
  int _pageSize = 25;
  final _searchController = TextEditingController();

  static const _pageSizeOptions = [10, 25, 50, 100];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _exportCsv(List<ProfileObject> profiles) {
    final headers = ['Name', 'Type', 'Contacts', 'State', 'ID'];
    final rows = profiles
        .map((p) => [
              profileName(p),
              p.type.name,
              '${p.contacts.length}',
              p.state.name,
              p.id,
            ])
        .toList();
    final csv = const CsvEncoder().convert([headers, ...rows]);
    Clipboard.setData(ClipboardData(text: csv));
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content:
              Text('Profiles CSV copied to clipboard (${rows.length} rows)'),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final asyncProfiles = ref.watch(profileSearchProvider(_query));

    return asyncProfiles.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text('Failed to load profiles',
                style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(friendlyError(error),
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: () => ref.invalidate(profileSearchProvider(_query)),
              icon: const Icon(Icons.refresh, size: 18),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (profiles) {
        final filtered = _selectedType != null
            ? profiles.where((p) => p.type == _selectedType).toList()
            : profiles;
        return _buildContent(context, theme, filtered);
      },
    );
  }

  Widget _buildContent(
      BuildContext context, ThemeData theme, List<ProfileObject> profiles) {
    // Reset stale index
    if (_selectedIndex != null && _selectedIndex! >= profiles.length) {
      _selectedIndex = null;
    }

    final showDetail = _selectedIndex != null;
    final totalPages = (profiles.length / _pageSize).ceil();
    final pageStart = _currentPage * _pageSize;
    final pageEnd = (pageStart + _pageSize).clamp(0, profiles.length);
    final pageItems =
        profiles.isNotEmpty ? profiles.sublist(pageStart, pageEnd) : [];

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
                    Icon(Icons.people_outlined,
                        size: 28, color: theme.colorScheme.primary),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Profiles',
                              style: theme.textTheme.headlineSmall?.copyWith(
                                  fontWeight: FontWeight.w600)),
                          if (profiles.isNotEmpty)
                            Text('${profiles.length} profiles',
                                style: theme.textTheme.bodySmall?.copyWith(
                                    color:
                                        theme.colorScheme.onSurfaceVariant)),
                        ],
                      ),
                    ),
                    FilledButton.icon(
                      onPressed: () => context.go('/profiles/new'),
                      icon: const Icon(Icons.add, size: 18),
                      label: const Text('New Profile'),
                    ),
                  ],
                ),
                const SizedBox(height: 20),

                // Search + type filter + export
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _searchController,
                        onChanged: (value) {
                          setState(() {
                            _query = value.trim();
                            _currentPage = 0;
                            _selectedIndex = null;
                          });
                        },
                        decoration: const InputDecoration(
                          hintText: 'Search profiles by name or contact...',
                          prefixIcon: Icon(Icons.search, size: 20),
                          isDense: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    _buildTypeFilter(theme),
                    const SizedBox(width: 12),
                    OutlinedButton.icon(
                      onPressed: () => _exportCsv(profiles),
                      icon: const Icon(Icons.download, size: 18),
                      label: const Text('Export'),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // Data table
                Container(
                  width: double.infinity,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                        color: theme.colorScheme.outlineVariant),
                  ),
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: SingleChildScrollView(
                      scrollDirection: Axis.horizontal,
                      child: DataTable(
                        showCheckboxColumn: false,
                        columns: const [
                          DataColumn(label: Text('NAME')),
                          DataColumn(label: Text('TYPE')),
                          DataColumn(
                              label: Text('CONTACTS'), numeric: true),
                          DataColumn(label: Text('STATE')),
                        ],
                        rows: List.generate(pageItems.length, (i) {
                          final profile = pageItems[i];
                          final globalIndex = pageStart + i;
                          final name = profileName(profile);
                          return DataRow(
                            selected: _selectedIndex == globalIndex,
                            onSelectChanged: (_) {
                              setState(
                                  () => _selectedIndex = globalIndex);
                              widget.onProfileSelected?.call(profile);
                            },
                            color: WidgetStateProperty.resolveWith(
                                (states) {
                              if (states
                                  .contains(WidgetState.selected)) {
                                return theme.colorScheme.primaryContainer
                                    .withAlpha(40);
                              }
                              return null;
                            }),
                            cells: [
                              DataCell(
                                Row(
                                  mainAxisSize: MainAxisSize.min,
                                  children: [
                                    ProfileAvatar(
                                      profileId: profile.id,
                                      name: name,
                                      size: 28,
                                    ),
                                    const SizedBox(width: 10),
                                    Text(name),
                                  ],
                                ),
                                onTap: () => _navigateToProfile(profile),
                              ),
                              DataCell(
                                  ProfileTypeBadge(type: profile.type)),
                              DataCell(
                                  Text('${profile.contacts.length}')),
                              DataCell(
                                  StateBadge(state: common.STATE.valueOf(profile.state.value)!)),
                            ],
                          );
                        }),
                      ),
                    ),
                  ),
                ),

                // Pagination footer
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
                          profiles.isEmpty
                              ? '0 of 0'
                              : '${pageStart + 1}-$pageEnd of ${profiles.length}',
                          style: theme.textTheme.bodySmall,
                        ),
                        IconButton(
                          icon: const Icon(Icons.chevron_left, size: 20),
                          onPressed: _currentPage > 0
                              ? () =>
                                  setState(() => _currentPage--)
                              : null,
                        ),
                        IconButton(
                          icon:
                              const Icon(Icons.chevron_right, size: 20),
                          onPressed:
                              _currentPage < totalPages - 1
                                  ? () => setState(
                                      () => _currentPage++)
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
        if (showDetail && _selectedIndex! < profiles.length)
          SizedBox(
            width: 380,
            child: Container(
              margin:
                  const EdgeInsets.only(top: 24, right: 24, bottom: 24),
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
                          icon: const Icon(Icons.open_in_new, size: 20),
                          tooltip: 'Open detail page',
                          onPressed: () => _navigateToProfile(
                              profiles[_selectedIndex!]),
                        ),
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
                      child: _ProfileDetailPanel(
                          profile: profiles[_selectedIndex!]),
                    ),
                  ),
                ],
              ),
            ),
          ),
      ],
    );
  }

  void _navigateToProfile(ProfileObject profile) {
    if (widget.onNavigateToDetail != null) {
      widget.onNavigateToDetail!(profile);
    } else {
      context.go('/profiles/${profile.id}');
    }
  }

  Widget _buildTypeFilter(ThemeData theme) {
    return PopupMenuButton<ProfileType?>(
      icon: Icon(
        _selectedType != null ? Icons.filter_alt : Icons.filter_alt_outlined,
        size: 20,
      ),
      tooltip: 'Filter by type',
      onSelected: (type) => setState(() {
        _selectedType = type;
        _currentPage = 0;
      }),
      itemBuilder: (context) => [
        const PopupMenuItem(value: null, child: Text('All Types')),
        const PopupMenuDivider(),
        const PopupMenuItem(
            value: ProfileType.PERSON, child: Text('Person')),
        const PopupMenuItem(
            value: ProfileType.INSTITUTION, child: Text('Institution')),
        const PopupMenuItem(value: ProfileType.BOT, child: Text('Bot')),
      ],
    );
  }
}

// -- Detail Panel --

class _ProfileDetailPanel extends StatelessWidget {
  const _ProfileDetailPanel({required this.profile});

  final ProfileObject profile;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = profileName(profile);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            ProfileAvatar(profileId: profile.id, name: name, size: 40),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(name,
                      style: theme.textTheme.titleMedium
                          ?.copyWith(fontWeight: FontWeight.w600)),
                  Text(profile.id,
                      style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                          fontFamily: 'monospace')),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 20),
        _DetailRow(label: 'Type', value: profile.type.name),
        _DetailRow(label: 'State', value: profile.state.name),
        const SizedBox(height: 16),
        // Contacts
        if (profile.contacts.isNotEmpty) ...[
          Text('Contacts',
              style: theme.textTheme.labelMedium
                  ?.copyWith(fontWeight: FontWeight.w500)),
          const SizedBox(height: 8),
          for (final contact in profile.contacts)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: Row(
                children: [
                  Icon(
                    contact.type == ContactType.EMAIL
                        ? Icons.email_outlined
                        : Icons.phone_outlined,
                    size: 16,
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(contact.detail,
                        style: theme.textTheme.bodySmall),
                  ),
                  if (contact.verified)
                    Icon(Icons.verified, size: 14, color: Colors.green),
                ],
              ),
            ),
          const SizedBox(height: 16),
        ],
        // Addresses
        if (profile.addresses.isNotEmpty) ...[
          Text('Addresses',
              style: theme.textTheme.labelMedium
                  ?.copyWith(fontWeight: FontWeight.w500)),
          const SizedBox(height: 8),
          for (final address in profile.addresses)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Icon(Icons.location_on_outlined,
                      size: 16,
                      color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      _formatAddress(address),
                      style: theme.textTheme.bodySmall,
                    ),
                  ),
                ],
              ),
            ),
        ],
      ],
    );
  }

  String _formatAddress(AddressObject address) {
    final parts = <String>[
      if (address.name.isNotEmpty) address.name,
      if (address.street.isNotEmpty) address.street,
      if (address.area.isNotEmpty) address.area,
      if (address.city.isNotEmpty) address.city,
      if (address.country.isNotEmpty) address.country,
      if (address.postcode.isNotEmpty) address.postcode,
    ];
    return parts.isNotEmpty ? parts.join(', ') : 'N/A';
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
            width: 110,
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
