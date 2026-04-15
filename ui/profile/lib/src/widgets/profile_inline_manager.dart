import 'package:antinvestor_api_profile/antinvestor_api_profile.dart' as profile;
import 'package:antinvestor_ui_profile/antinvestor_ui_profile.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Inline contact and address management for a profile.
///
/// Shows existing contacts and addresses from the profile service,
/// with buttons to add new ones directly. Drop this into any detail
/// page that has a profileId.
///
/// ```dart
/// ProfileInlineManager(profileId: organization.profileId)
/// ```
class ProfileInlineManager extends ConsumerStatefulWidget {
  const ProfileInlineManager({
    super.key,
    required this.profileId,
  });

  final String profileId;

  @override
  ConsumerState<ProfileInlineManager> createState() =>
      _ProfileInlineManagerState();
}

class _ProfileInlineManagerState extends ConsumerState<ProfileInlineManager> {
  final _contactCtrl = TextEditingController();
  bool _contactIsPhone = false;
  bool _addingContact = false;

  final _addressNameCtrl = TextEditingController();
  final _streetCtrl = TextEditingController();
  final _cityCtrl = TextEditingController();
  final _countryCtrl = TextEditingController();
  final _postalCodeCtrl = TextEditingController();
  bool _addingAddress = false;
  bool _showAddressForm = false;

  @override
  void dispose() {
    _contactCtrl.dispose();
    _addressNameCtrl.dispose();
    _streetCtrl.dispose();
    _cityCtrl.dispose();
    _countryCtrl.dispose();
    _postalCodeCtrl.dispose();
    super.dispose();
  }

  Future<void> _addContact() async {
    final value = _contactCtrl.text.trim();
    if (value.isEmpty || widget.profileId.isEmpty) return;
    setState(() => _addingContact = true);
    try {
      await ref.read(contactNotifierProvider.notifier).addContact(
            profile.AddContactRequest(
              id: widget.profileId,
              contact: value,
            ),
          );
      if (mounted) {
        _contactCtrl.clear();
        ref.invalidate(profileByIdProvider(widget.profileId));
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Contact added')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to add contact: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _addingContact = false);
    }
  }

  Future<void> _addAddress() async {
    if (widget.profileId.isEmpty) return;
    setState(() => _addingAddress = true);
    try {
      await ref.read(addressNotifierProvider.notifier).addAddress(
            profile.AddAddressRequest(
              id: widget.profileId,
              address: profile.AddressObject(
                name: _addressNameCtrl.text.trim(),
                street: _streetCtrl.text.trim(),
                city: _cityCtrl.text.trim(),
                country: _countryCtrl.text.trim(),
                postcode: _postalCodeCtrl.text.trim(),
              ),
            ),
          );
      if (mounted) {
        _addressNameCtrl.clear();
        _streetCtrl.clear();
        _cityCtrl.clear();
        _countryCtrl.clear();
        _postalCodeCtrl.clear();
        setState(() => _showAddressForm = false);
        ref.invalidate(profileByIdProvider(widget.profileId));
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Address added')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to add address: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _addingAddress = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (widget.profileId.isEmpty) return const SizedBox.shrink();

    final theme = Theme.of(context);
    final profileAsync = ref.watch(profileByIdProvider(widget.profileId));

    return profileAsync.when(
      loading: () => const Padding(
        padding: EdgeInsets.all(16),
        child: Center(child: CircularProgressIndicator()),
      ),
      error: (e, _) => Padding(
        padding: const EdgeInsets.all(16),
        child: Text('Could not load profile: $e',
            style: TextStyle(color: theme.colorScheme.error)),
      ),
      data: (p) => Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildContactsSection(theme, p.contacts),
          const SizedBox(height: 16),
          _buildAddressesSection(theme, p.addresses),
        ],
      ),
    );
  }

  Widget _buildContactsSection(
      ThemeData theme, List<profile.ContactObject> contacts) {
    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
        side: BorderSide(
            color: theme.colorScheme.outlineVariant.withAlpha(38)),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.contacts_outlined,
                    size: 20, color: theme.colorScheme.primary),
                const SizedBox(width: 8),
                Text('Contacts',
                    style: theme.textTheme.titleMedium
                        ?.copyWith(fontWeight: FontWeight.w600)),
                const Spacer(),
                Text('${contacts.length}',
                    style: theme.textTheme.labelMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant)),
              ],
            ),
            if (contacts.isNotEmpty) ...[
              const SizedBox(height: 8),
              ...contacts.map((c) => ContactListTile(contact: c)),
            ],
            const SizedBox(height: 12),
            // Inline add contact
            Row(
              children: [
                ChoiceChip(
                  label: const Text('Email'),
                  selected: !_contactIsPhone,
                  onSelected: (_) =>
                      setState(() => _contactIsPhone = false),
                  visualDensity: VisualDensity.compact,
                ),
                const SizedBox(width: 6),
                ChoiceChip(
                  label: const Text('Phone'),
                  selected: _contactIsPhone,
                  onSelected: (_) =>
                      setState(() => _contactIsPhone = true),
                  visualDensity: VisualDensity.compact,
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _contactCtrl,
                    decoration: InputDecoration(
                      hintText: _contactIsPhone
                          ? '+254 700 000 000'
                          : 'info@org.com',
                      prefixIcon: Icon(
                        _contactIsPhone
                            ? Icons.phone_outlined
                            : Icons.email_outlined,
                        size: 18,
                      ),
                      isDense: true,
                      border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8)),
                    ),
                    keyboardType: _contactIsPhone
                        ? TextInputType.phone
                        : TextInputType.emailAddress,
                    onSubmitted: (_) => _addContact(),
                  ),
                ),
                const SizedBox(width: 8),
                FilledButton.tonal(
                  onPressed: _addingContact ? null : _addContact,
                  child: _addingContact
                      ? const SizedBox(
                          width: 16,
                          height: 16,
                          child: CircularProgressIndicator(
                              strokeWidth: 2),
                        )
                      : const Text('Add'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAddressesSection(
      ThemeData theme, List<profile.AddressObject> addresses) {
    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
        side: BorderSide(
            color: theme.colorScheme.outlineVariant.withAlpha(38)),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.location_on_outlined,
                    size: 20, color: theme.colorScheme.primary),
                const SizedBox(width: 8),
                Text('Addresses',
                    style: theme.textTheme.titleMedium
                        ?.copyWith(fontWeight: FontWeight.w600)),
                const Spacer(),
                if (!_showAddressForm)
                  FilledButton.tonal(
                    onPressed: () =>
                        setState(() => _showAddressForm = true),
                    child: const Text('Add Address'),
                  ),
              ],
            ),
            if (addresses.isNotEmpty) ...[
              const SizedBox(height: 8),
              ...addresses.map((a) => Padding(
                    padding: const EdgeInsets.only(bottom: 8),
                    child: AddressTile(address: a),
                  )),
            ] else if (!_showAddressForm)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 12),
                child: Text('No addresses yet',
                    style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant)),
              ),
            if (_showAddressForm) ...[
              const SizedBox(height: 12),
              TextField(
                controller: _addressNameCtrl,
                decoration: InputDecoration(
                  labelText: 'Label',
                  hintText: 'e.g. Head Office',
                  isDense: true,
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8)),
                ),
              ),
              const SizedBox(height: 8),
              TextField(
                controller: _streetCtrl,
                decoration: InputDecoration(
                  labelText: 'Street',
                  hintText: 'e.g. 123 Kenyatta Avenue',
                  isDense: true,
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8)),
                ),
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _cityCtrl,
                      decoration: InputDecoration(
                        labelText: 'City',
                        hintText: 'e.g. Nairobi',
                        isDense: true,
                        border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(8)),
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: TextField(
                      controller: _countryCtrl,
                      decoration: InputDecoration(
                        labelText: 'Country',
                        hintText: 'e.g. Kenya',
                        isDense: true,
                        border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(8)),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              TextField(
                controller: _postalCodeCtrl,
                decoration: InputDecoration(
                  labelText: 'Postal Code',
                  hintText: 'e.g. 00100',
                  isDense: true,
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8)),
                ),
              ),
              const SizedBox(height: 12),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  TextButton(
                    onPressed: () =>
                        setState(() => _showAddressForm = false),
                    child: const Text('Cancel'),
                  ),
                  const SizedBox(width: 8),
                  FilledButton(
                    onPressed: _addingAddress ? null : _addAddress,
                    child: _addingAddress
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                                strokeWidth: 2, color: Colors.white),
                          )
                        : const Text('Save Address'),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}
