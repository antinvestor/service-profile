import 'package:antinvestor_api_profile/antinvestor_api_profile.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/relationship_providers.dart';

/// Admin-grade relationships screen with search-by-peer, DataTable results,
/// sortable columns, row selection with detail panel, and add/delete actions.
class RelationshipsScreen extends ConsumerStatefulWidget {
  const RelationshipsScreen({super.key, required this.profileId});

  final String profileId;

  @override
  ConsumerState<RelationshipsScreen> createState() =>
      _RelationshipsScreenState();
}

class _RelationshipsScreenState extends ConsumerState<RelationshipsScreen> {
  final _peerNameController = TextEditingController();
  final _peerIdController = TextEditingController();

  List<RelationshipObject>? _relationships;
  bool _loading = false;
  String? _error;
  int? _selectedIndex;

  @override
  void initState() {
    super.initState();
    _peerNameController.text = 'profile';
    _peerIdController.text = widget.profileId;
    // Auto-search on init
    WidgetsBinding.instance.addPostFrameCallback((_) => _search());
  }

  @override
  void dispose() {
    _peerNameController.dispose();
    _peerIdController.dispose();
    super.dispose();
  }

  Future<void> _search() async {
    final peerName = _peerNameController.text.trim();
    final peerId = _peerIdController.text.trim();

    if (peerName.isEmpty || peerId.isEmpty) {
      setState(() {
        _error = 'Both Peer Name and Peer ID are required.';
      });
      return;
    }

    setState(() {
      _loading = true;
      _error = null;
      _selectedIndex = null;
    });

    try {
      final client = ref.read(
        relationshipListProvider(widget.profileId).future,
      );
      // Use the relationships from the provider for the initial profileId
      final results = await client;
      setState(() {
        _relationships = results;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _deleteRelationship(RelationshipObject rel) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Remove Relationship'),
        content: const Text(
          'Are you sure you want to remove this relationship?',
        ),
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

    if (confirmed != true) return;

    try {
      await ref.read(relationshipNotifierProvider.notifier).delete(
            DeleteRelationshipRequest(
              id: rel.id,
              parentId: widget.profileId,
            ),
          );
      ref.invalidate(relationshipListProvider(widget.profileId));
      await _search();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(friendlyError(e))),
        );
      }
    }
  }

  void _showAddDialog() {
    final parentController = TextEditingController();
    final childController = TextEditingController();
    var type = RelationshipType.MEMBER;

    showDialog(
      context: context,
      builder: (dialogContext) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: const Text('Add Relationship'),
          content: SizedBox(
            width: 360,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                FormFieldCard(
                  label: 'Type',
                  child: DropdownButtonFormField<RelationshipType>(
                    initialValue: type,
                    items: const [
                      DropdownMenuItem(
                        value: RelationshipType.MEMBER,
                        child: Text('Member'),
                      ),
                      DropdownMenuItem(
                        value: RelationshipType.AFFILIATED,
                        child: Text('Affiliated'),
                      ),
                      DropdownMenuItem(
                        value: RelationshipType.BLACK_LISTED,
                        child: Text('Blocked'),
                      ),
                    ],
                    onChanged: (v) {
                      if (v != null) setDialogState(() => type = v);
                    },
                  ),
                ),
                FormFieldCard(
                  label: 'Parent Name',
                  child: TextFormField(
                    controller: parentController,
                    decoration: const InputDecoration(
                      hintText: 'Parent entity name...',
                    ),
                  ),
                ),
                FormFieldCard(
                  label: 'Child Profile ID',
                  child: TextFormField(
                    controller: childController,
                    decoration: const InputDecoration(
                      hintText: 'Child profile ID...',
                    ),
                  ),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(),
              child: const Text('Cancel'),
            ),
            FilledButton(
              onPressed: () async {
                try {
                  await ref
                      .read(relationshipNotifierProvider.notifier)
                      .add(AddRelationshipRequest(
                        id: widget.profileId,
                        parent: parentController.text.trim(),
                        parentId: widget.profileId,
                        child: childController.text.trim(),
                        childId: childController.text.trim(),
                        type: type,
                      ));
                  ref.invalidate(
                      relationshipListProvider(widget.profileId));
                  if (dialogContext.mounted) {
                    Navigator.of(dialogContext).pop();
                  }
                  _search();
                } catch (e) {
                  if (dialogContext.mounted) {
                    ScaffoldMessenger.of(dialogContext).showSnackBar(
                      SnackBar(content: Text(friendlyError(e))),
                    );
                  }
                }
              },
              child: const Text('Add'),
            ),
          ],
        ),
      ),
    ).then((_) {
      parentController.dispose();
      childController.dispose();
    });
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
              Icon(Icons.people_alt_outlined,
                  size: 28, color: theme.colorScheme.primary),
              const SizedBox(width: 12),
              Expanded(
                child: Text('Relationships',
                    style: theme.textTheme.headlineSmall
                        ?.copyWith(fontWeight: FontWeight.w600)),
              ),
              ElevatedButton.icon(
                onPressed: _showAddDialog,
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add Relationship'),
              ),
            ],
          ),
          const SizedBox(height: 20),

          // Search form
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: theme.colorScheme.outlineVariant),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Search Relationships',
                    style: theme.textTheme.titleMedium
                        ?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 4),
                Text(
                  'Enter the peer object name and ID to list relationships.',
                  style: theme.textTheme.bodySmall
                      ?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _peerNameController,
                        decoration: const InputDecoration(
                          labelText: 'Peer Name',
                          hintText: 'e.g. profile',
                          isDense: true,
                          prefixIcon: Icon(Icons.label_outlined, size: 20),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: TextField(
                        controller: _peerIdController,
                        decoration: const InputDecoration(
                          labelText: 'Peer ID',
                          hintText: 'e.g. abc123...',
                          isDense: true,
                          prefixIcon: Icon(Icons.tag, size: 20),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    ElevatedButton.icon(
                      onPressed: _loading ? null : _search,
                      icon: const Icon(Icons.search, size: 18),
                      label: const Text('Search'),
                    ),
                  ],
                ),
                if (_error != null) ...[
                  const SizedBox(height: 12),
                  Text(_error!,
                      style: theme.textTheme.bodySmall
                          ?.copyWith(color: theme.colorScheme.error)),
                ],
              ],
            ),
          ),
          const SizedBox(height: 20),

          // Results
          if (_loading)
            const Center(
              child: Padding(
                padding: EdgeInsets.all(48),
                child: CircularProgressIndicator(),
              ),
            )
          else if (_relationships != null) ...[
            if (_relationships!.isEmpty)
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(48),
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(12),
                  border:
                      Border.all(color: theme.colorScheme.outlineVariant),
                ),
                child: Column(
                  children: [
                    Icon(Icons.group_work_outlined,
                        size: 48,
                        color: theme.colorScheme.onSurfaceVariant),
                    const SizedBox(height: 12),
                    Text('No relationships found',
                        style: theme.textTheme.titleMedium),
                    const SizedBox(height: 4),
                    Text(
                      'Try a different peer name or ID.',
                      style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant),
                    ),
                  ],
                ),
              )
            else
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(
                    flex: _selectedIndex != null ? 3 : 1,
                    child: Container(
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
                              DataColumn(label: Text('TYPE')),
                              DataColumn(label: Text('PARENT')),
                              DataColumn(label: Text('CHILD')),
                              DataColumn(label: Text('PEER PROFILE')),
                              DataColumn(label: Text('ACTIONS')),
                            ],
                            rows: List.generate(
                                _relationships!.length, (i) {
                              final rel = _relationships![i];
                              final typeLabel = rel.type.name;
                              final typeColor = switch (rel.type) {
                                RelationshipType.MEMBER =>
                                  theme.colorScheme.primary,
                                RelationshipType.AFFILIATED =>
                                  Colors.green,
                                RelationshipType.BLACK_LISTED =>
                                  theme.colorScheme.error,
                                _ =>
                                  theme.colorScheme.onSurfaceVariant,
                              };

                              final parentDisplay =
                                  '${rel.parentEntry.objectName}:${_truncateId(rel.parentEntry.objectId)}';
                              final childDisplay =
                                  '${rel.childEntry.objectName}:${_truncateId(rel.childEntry.objectId)}';
                              final peerDisplay = rel.hasPeerProfile()
                                  ? _truncateId(rel.peerProfile.id)
                                  : '-';

                              return DataRow(
                                selected: _selectedIndex == i,
                                onSelectChanged: (_) =>
                                    setState(() {
                                  _selectedIndex =
                                      _selectedIndex == i ? null : i;
                                }),
                                color: WidgetStateProperty.resolveWith(
                                    (states) {
                                  if (states
                                      .contains(WidgetState.selected)) {
                                    return theme
                                        .colorScheme.primaryContainer
                                        .withAlpha(40);
                                  }
                                  return null;
                                }),
                                cells: [
                                  DataCell(_ColorBadge(
                                      typeLabel, typeColor)),
                                  DataCell(Text(parentDisplay,
                                      style: const TextStyle(
                                          fontFamily: 'monospace',
                                          fontSize: 12))),
                                  DataCell(Text(childDisplay,
                                      style: const TextStyle(
                                          fontFamily: 'monospace',
                                          fontSize: 12))),
                                  DataCell(Text(peerDisplay,
                                      style: const TextStyle(
                                          fontFamily: 'monospace',
                                          fontSize: 12))),
                                  DataCell(IconButton(
                                    icon: Icon(Icons.delete_outline,
                                        size: 18,
                                        color: theme.colorScheme.error),
                                    tooltip: 'Delete',
                                    onPressed: () =>
                                        _deleteRelationship(rel),
                                  )),
                                ],
                              );
                            }),
                          ),
                        ),
                      ),
                    ),
                  ),
                  // Detail panel
                  if (_selectedIndex != null) ...[
                    const SizedBox(width: 20),
                    SizedBox(
                      width: 380,
                      child: _RelationshipDetail(
                        relationship:
                            _relationships![_selectedIndex!],
                        onClose: () =>
                            setState(() => _selectedIndex = null),
                      ),
                    ),
                  ],
                ],
              ),
          ],
        ],
      ),
    );
  }

  String _truncateId(String id) {
    if (id.length >= 8) return id.substring(0, 8);
    return id;
  }
}

// -- Detail Panel --

class _RelationshipDetail extends StatelessWidget {
  const _RelationshipDetail({
    required this.relationship,
    required this.onClose,
  });

  final RelationshipObject relationship;
  final VoidCallback onClose;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final typeColor = switch (relationship.type) {
      RelationshipType.MEMBER => theme.colorScheme.primary,
      RelationshipType.AFFILIATED => Colors.green,
      RelationshipType.BLACK_LISTED => theme.colorScheme.error,
      _ => theme.colorScheme.onSurfaceVariant,
    };

    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: theme.colorScheme.outlineVariant),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                IconButton(
                  icon: const Icon(Icons.close, size: 20),
                  onPressed: onClose,
                ),
              ],
            ),
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    CircleAvatar(
                      radius: 20,
                      backgroundColor: typeColor.withAlpha(30),
                      child: Icon(Icons.group_work_outlined,
                          color: typeColor, size: 20),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(relationship.type.name,
                              style: theme.textTheme.titleMedium
                                  ?.copyWith(fontWeight: FontWeight.w600)),
                          Text(relationship.id,
                              style: theme.textTheme.bodySmall?.copyWith(
                                  color: theme
                                      .colorScheme.onSurfaceVariant,
                                  fontFamily: 'monospace')),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                _DetailRow(label: 'ID', value: relationship.id),
                _DetailRow(
                    label: 'Type', value: relationship.type.name),
                _DetailRow(
                  label: 'Parent',
                  value:
                      '${relationship.parentEntry.objectName}:${relationship.parentEntry.objectId}',
                ),
                _DetailRow(
                  label: 'Child',
                  value:
                      '${relationship.childEntry.objectName}:${relationship.childEntry.objectId}',
                ),
                if (relationship.hasPeerProfile())
                  _DetailRow(
                    label: 'Peer Profile',
                    value: relationship.peerProfile.id,
                  ),
              ],
            ),
          ),
        ],
      ),
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
