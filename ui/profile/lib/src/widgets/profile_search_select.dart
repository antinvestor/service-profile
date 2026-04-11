import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/profile_badge.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/profile_providers.dart';

/// An embeddable search-and-select widget for profiles.
///
/// Shows a text field that searches profiles as the user types, displays
/// results in a dropdown overlay, and calls [onSelected] when a profile
/// is picked.
///
/// ```dart
/// ProfileSearchSelect(
///   onSelected: (profile) => setState(() => _selectedId = profile.id),
/// )
/// ```
class ProfileSearchSelect extends ConsumerStatefulWidget {
  const ProfileSearchSelect({
    super.key,
    required this.onSelected,
    this.label = 'Search profiles',
    this.initialQuery = '',
    this.autofocus = false,
  });

  final ValueChanged<ProfileObject> onSelected;
  final String label;
  final String initialQuery;
  final bool autofocus;

  @override
  ConsumerState<ProfileSearchSelect> createState() =>
      _ProfileSearchSelectState();
}

class _ProfileSearchSelectState extends ConsumerState<ProfileSearchSelect> {
  late final TextEditingController _controller;
  final _focusNode = FocusNode();
  final _layerLink = LayerLink();
  OverlayEntry? _overlay;
  String _query = '';

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: widget.initialQuery);
    _query = widget.initialQuery;
    _focusNode.addListener(_onFocusChanged);
  }

  @override
  void dispose() {
    _removeOverlay();
    _focusNode.removeListener(_onFocusChanged);
    _focusNode.dispose();
    _controller.dispose();
    super.dispose();
  }

  void _onFocusChanged() {
    if (!_focusNode.hasFocus) {
      // Delay removal so tap on result can register.
      Future.delayed(const Duration(milliseconds: 200), _removeOverlay);
    }
  }

  void _onQueryChanged(String value) {
    setState(() => _query = value.trim());
    if (_query.length >= 2) {
      _showOverlay();
    } else {
      _removeOverlay();
    }
  }

  void _showOverlay() {
    _removeOverlay();
    final renderBox = context.findRenderObject() as RenderBox?;
    if (renderBox == null) return;

    final width = renderBox.size.width;

    _overlay = OverlayEntry(
      builder: (_) => Positioned(
        width: width,
        child: CompositedTransformFollower(
          link: _layerLink,
          showWhenUnlinked: false,
          offset: Offset(0, renderBox.size.height + 4),
          child: Material(
            elevation: 4,
            borderRadius: BorderRadius.circular(12),
            child: _ResultsList(
              query: _query,
              onSelected: _onProfileSelected,
            ),
          ),
        ),
      ),
    );
    Overlay.of(context).insert(_overlay!);
  }

  void _removeOverlay() {
    _overlay?.remove();
    _overlay = null;
  }

  void _onProfileSelected(ProfileObject profile) {
    _removeOverlay();
    // Show selected name in the field.
    final name = _extractName(profile);
    _controller.text = name;
    _query = '';
    widget.onSelected(profile);
  }

  String _extractName(ProfileObject profile) {
    try {
      final props = profile.properties;
      if (props.fields.containsKey('name')) {
        final name = props.fields['name']!.stringValue;
        if (name.isNotEmpty) return name;
      }
    } catch (_) {}
    try {
      if (profile.contacts.isNotEmpty) return profile.contacts.first.detail;
    } catch (_) {}
    return profile.id.isNotEmpty
        ? (profile.id.length > 12
            ? '${profile.id.substring(0, 12)}...'
            : profile.id)
        : 'Profile';
  }

  @override
  Widget build(BuildContext context) {
    return CompositedTransformTarget(
      link: _layerLink,
      child: TextField(
        controller: _controller,
        focusNode: _focusNode,
        autofocus: widget.autofocus,
        onChanged: _onQueryChanged,
        decoration: InputDecoration(
          labelText: widget.label,
          prefixIcon: const Icon(Icons.search, size: 20),
          border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
          contentPadding:
              const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          isDense: true,
        ),
      ),
    );
  }
}

/// Internal results list rendered inside the overlay.
class _ResultsList extends ConsumerWidget {
  const _ResultsList({required this.query, required this.onSelected});

  final String query;
  final ValueChanged<ProfileObject> onSelected;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final results = ref.watch(profileSearchProvider(query));

    return ConstrainedBox(
      constraints: const BoxConstraints(maxHeight: 260),
      child: results.when(
        loading: () => const Padding(
          padding: EdgeInsets.all(16),
          child: Center(child: CircularProgressIndicator(strokeWidth: 2)),
        ),
        error: (e, _) => Padding(
          padding: const EdgeInsets.all(16),
          child: Text(
            'Search failed',
            style: TextStyle(color: Theme.of(context).colorScheme.error),
          ),
        ),
        data: (profiles) {
          if (profiles.isEmpty) {
            return const Padding(
              padding: EdgeInsets.all(16),
              child: Text('No profiles found'),
            );
          }
          return ListView.builder(
            shrinkWrap: true,
            padding: const EdgeInsets.symmetric(vertical: 4),
            itemCount: profiles.length,
            itemBuilder: (context, index) {
              final profile = profiles[index];
              final name = _nameOf(profile);
              return ListTile(
                dense: true,
                leading: ProfileAvatar(
                  profileId: profile.id,
                  name: name,
                  size: 32,
                ),
                title: Text(name),
                onTap: () => onSelected(profile),
              );
            },
          );
        },
      ),
    );
  }

  String _nameOf(ProfileObject profile) {
    try {
      final props = profile.properties;
      if (props.fields.containsKey('name')) {
        final n = props.fields['name']!.stringValue;
        if (n.isNotEmpty) return n;
      }
    } catch (_) {}
    try {
      if (profile.contacts.isNotEmpty) return profile.contacts.first.detail;
    } catch (_) {}
    return profile.id.length > 12
        ? '${profile.id.substring(0, 12)}...'
        : profile.id;
  }
}
