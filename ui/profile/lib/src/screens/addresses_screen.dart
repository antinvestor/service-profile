import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:csv/csv.dart' show CsvEncoder;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/profile_providers.dart';
import '../widgets/profile_badge_by_id.dart';
import '../widgets/profile_card.dart';

/// An address entry paired with the owning profile's ID and name.
typedef ProfileAddress = ({
  String profileId,
  String profileName,
  AddressObject address,
});

/// Provider that extracts all addresses from loaded profiles into a flat list.
final allAddressesProvider =
    FutureProvider<List<ProfileAddress>>((ref) async {
  final profiles = await ref.watch(profileSearchProvider('').future);
  final addresses = <ProfileAddress>[];
  for (final profile in profiles) {
    final name = profileName(profile);
    for (final address in profile.addresses) {
      addresses.add((
        profileId: profile.id,
        profileName: name,
        address: address,
      ));
    }
  }
  return addresses;
});

/// Admin-grade addresses screen with DataTable, CSV export, search,
/// pagination, and detail panel.
class AddressesScreen extends ConsumerStatefulWidget {
  const AddressesScreen({super.key, this.profileId});

  /// When provided, only shows addresses for this profile.
  final String? profileId;

  @override
  ConsumerState<AddressesScreen> createState() => _AddressesScreenState();
}

class _AddressesScreenState extends ConsumerState<AddressesScreen> {
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

  void _exportCsv(List<ProfileAddress> addresses) {
    final headers = ['Name', 'City', 'Country', 'Profile', 'Profile ID'];
    final rows = addresses
        .map((pa) => [
              pa.address.name,
              pa.address.city,
              pa.address.country,
              pa.profileName,
              pa.profileId,
            ])
        .toList();
    final csv = const CsvEncoder().convert([headers, ...rows]);
    Clipboard.setData(ClipboardData(text: csv));

    // Audit: log the export event
    debugPrint('[AUDIT] Exported ${rows.length} Addresses as csv');

    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
              'Addresses CSV copied to clipboard (${rows.length} rows)'),
        ),
      );
    }
  }

  List<ProfileAddress> _filter(List<ProfileAddress> items) {
    if (_searchQuery.isEmpty) return items;
    final q = _searchQuery.toLowerCase();
    return items.where((pa) {
      return pa.address.name.toLowerCase().contains(q) ||
          pa.address.city.toLowerCase().contains(q) ||
          pa.address.country.toLowerCase().contains(q) ||
          pa.profileName.toLowerCase().contains(q);
    }).toList();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    final asyncAddresses = widget.profileId != null
        ? ref
            .watch(profileByIdProvider(widget.profileId!))
            .whenData((p) {
            final name = profileName(p);
            return p.addresses
                .map((a) => (
                      profileId: widget.profileId!,
                      profileName: name,
                      address: a,
                    ))
                .toList();
          })
        : ref.watch(allAddressesProvider);

    return asyncAddresses.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, _) => Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline,
                size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text('Failed to load addresses',
                style: theme.textTheme.titleMedium),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: () {
                if (widget.profileId != null) {
                  ref.invalidate(profileByIdProvider(widget.profileId!));
                } else {
                  ref.invalidate(allAddressesProvider);
                }
              },
              icon: const Icon(Icons.refresh, size: 18),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
      data: (allAddresses) {
        final addresses = _filter(allAddresses);
        return _buildContent(
            context, theme, addresses, allAddresses.length);
      },
    );
  }

  Widget _buildContent(BuildContext context, ThemeData theme,
      List<ProfileAddress> addresses, int totalUnfiltered) {
    if (_selectedIndex != null && _selectedIndex! >= addresses.length) {
      _selectedIndex = null;
    }

    final showDetail = _selectedIndex != null;
    final totalPages = (addresses.length / _pageSize).ceil();
    final pageStart = _currentPage * _pageSize;
    final pageEnd = (pageStart + _pageSize).clamp(0, addresses.length);
    final pageItems =
        addresses.isNotEmpty ? addresses.sublist(pageStart, pageEnd) : [];

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
                    Icon(Icons.location_on_outlined,
                        size: 28, color: theme.colorScheme.primary),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Addresses',
                              style: theme.textTheme.headlineSmall
                                  ?.copyWith(fontWeight: FontWeight.w600)),
                          Text('$totalUnfiltered addresses',
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
                          hintText: 'Search addresses...',
                          prefixIcon: Icon(Icons.search, size: 20),
                          isDense: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    OutlinedButton.icon(
                      onPressed: () => _exportCsv(addresses),
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
                          DataColumn(label: Text('LOCATION')),
                          DataColumn(label: Text('PROFILE')),
                        ],
                        rows: List.generate(pageItems.length, (i) {
                          final pa = pageItems[i];
                          final globalIndex = pageStart + i;
                          final parts = [
                            pa.address.city,
                            pa.address.country,
                          ].where((s) => s.isNotEmpty).join(', ');

                          return DataRow(
                            selected: _selectedIndex == globalIndex,
                            onSelectChanged: (_) => setState(
                                () => _selectedIndex = globalIndex),
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
                                  Icon(Icons.location_on_outlined,
                                      size: 16,
                                      color: theme.colorScheme.primary),
                                  const SizedBox(width: 8),
                                  Text(
                                    pa.address.name.isNotEmpty
                                        ? pa.address.name
                                        : 'Address',
                                    style: const TextStyle(
                                        fontWeight: FontWeight.w500),
                                  ),
                                ],
                              )),
                              DataCell(Text(
                                  parts.isNotEmpty ? parts : '\u2014')),
                              DataCell(ProfileBadgeById(
                                profileId: pa.profileId,
                                avatarSize: 24,
                              )),
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
                          addresses.isEmpty
                              ? '0 of 0'
                              : '${pageStart + 1}-$pageEnd of ${addresses.length}',
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
        if (showDetail && _selectedIndex! < addresses.length)
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
                      child: _AddressDetail(
                          pa: addresses[_selectedIndex!]),
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

class _AddressDetail extends StatelessWidget {
  const _AddressDetail({required this.pa});
  final ProfileAddress pa;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final addr = pa.address;
    final fullAddress = [
      addr.house,
      addr.street,
      addr.area,
      addr.city,
      addr.country,
      addr.postcode,
    ].where((s) => s.isNotEmpty).join(', ');

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: theme.colorScheme.primary.withAlpha(25),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(Icons.location_on_outlined,
                  size: 24, color: theme.colorScheme.primary),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(addr.name.isNotEmpty ? addr.name : 'Address',
                      style: theme.textTheme.titleMedium
                          ?.copyWith(fontWeight: FontWeight.w600)),
                  ProfileBadgeById(
                      profileId: pa.profileId, avatarSize: 20),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 20),
        _Row('Address', fullAddress.isNotEmpty ? fullAddress : '\u2014'),
        if (addr.street.isNotEmpty) _Row('Street', addr.street),
        if (addr.area.isNotEmpty) _Row('Area', addr.area),
        if (addr.city.isNotEmpty) _Row('City', addr.city),
        if (addr.country.isNotEmpty) _Row('Country', addr.country),
        if (addr.postcode.isNotEmpty) _Row('Postcode', addr.postcode),
        if (addr.latitude != 0 || addr.longitude != 0)
          _Row('Coordinates',
              '${addr.latitude.toStringAsFixed(5)}, ${addr.longitude.toStringAsFixed(5)}'),
        const SizedBox(height: 4),
        ProfileBadgeById(profileId: pa.profileId),
      ],
    );
  }
}

class _Row extends StatelessWidget {
  const _Row(this.label, this.value);
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
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
            child: Text(value,
                style: theme.textTheme.bodySmall
                    ?.copyWith(fontWeight: FontWeight.w500)),
          ),
        ],
      ),
    );
  }
}
