import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/area_providers.dart';

/// Screen for creating a new area with name, type, description,
/// and boundary points.
class AreaCreateScreen extends ConsumerStatefulWidget {
  const AreaCreateScreen({super.key});

  @override
  ConsumerState<AreaCreateScreen> createState() => _AreaCreateScreenState();
}

class _AreaCreateScreenState extends ConsumerState<AreaCreateScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _geometryController = TextEditingController();
  AreaType _areaType = AreaType.AREA_TYPE_LAND;
  final List<({double lat, double lon})> _boundaryPoints = [];
  final _latController = TextEditingController();
  final _lonController = TextEditingController();
  bool _submitting = false;

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    _geometryController.dispose();
    _latController.dispose();
    _lonController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('New Area'),
        actions: [
          FilledButton.icon(
            onPressed: _submitting ? null : _submit,
            icon: _submitting
                ? const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.check, size: 18),
            label: const Text('Create'),
          ),
          const SizedBox(width: 16),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(24),
          children: [
            FormFieldCard(
              label: 'Name',
              isRequired: true,
              description: 'A descriptive name for this area.',
              child: TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(hintText: 'Enter area name'),
                validator: (v) =>
                    (v == null || v.trim().isEmpty) ? 'Name is required' : null,
              ),
            ),
            FormFieldCard(
              label: 'Area Type',
              isRequired: true,
              description: 'The classification of this area.',
              child: DropdownButtonFormField<AreaType>(
                initialValue: _areaType,
                onChanged: (v) {
                  if (v != null) setState(() => _areaType = v);
                },
                items: const [
                  DropdownMenuItem(
                    value: AreaType.AREA_TYPE_LAND,
                    child: Text('Land'),
                  ),
                  DropdownMenuItem(
                    value: AreaType.AREA_TYPE_BUILDING,
                    child: Text('Building'),
                  ),
                  DropdownMenuItem(
                    value: AreaType.AREA_TYPE_ZONE,
                    child: Text('Zone'),
                  ),
                  DropdownMenuItem(
                    value: AreaType.AREA_TYPE_FENCE,
                    child: Text('Fence'),
                  ),
                  DropdownMenuItem(
                    value: AreaType.AREA_TYPE_CUSTOM,
                    child: Text('Custom'),
                  ),
                ],
              ),
            ),
            FormFieldCard(
              label: 'Description',
              description: 'Optional description of the area.',
              child: TextFormField(
                controller: _descriptionController,
                decoration:
                    const InputDecoration(hintText: 'Enter description'),
                maxLines: 3,
              ),
            ),
            FormFieldCard(
              label: 'Geometry (GeoJSON)',
              description:
                  'Paste GeoJSON geometry or add boundary points below.',
              child: TextFormField(
                controller: _geometryController,
                decoration:
                    const InputDecoration(hintText: '{"type":"Polygon",...}'),
                maxLines: 4,
                style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
              ),
            ),

            // Boundary points section
            FormSection(
              title: 'Boundary Points',
              description: 'Add latitude/longitude pairs for the area boundary.',
              children: [
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _latController,
                        decoration:
                            const InputDecoration(hintText: 'Latitude'),
                        keyboardType: const TextInputType.numberWithOptions(
                          decimal: true,
                          signed: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: TextFormField(
                        controller: _lonController,
                        decoration:
                            const InputDecoration(hintText: 'Longitude'),
                        keyboardType: const TextInputType.numberWithOptions(
                          decimal: true,
                          signed: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    IconButton.filled(
                      onPressed: _addBoundaryPoint,
                      icon: const Icon(Icons.add, size: 20),
                      tooltip: 'Add point',
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                if (_boundaryPoints.isNotEmpty)
                  ...List.generate(_boundaryPoints.length, (i) {
                    final pt = _boundaryPoints[i];
                    return ListTile(
                      dense: true,
                      leading: CircleAvatar(
                        radius: 14,
                        child: Text(
                          '${i + 1}',
                          style: theme.textTheme.labelSmall,
                        ),
                      ),
                      title: Text(
                        '${pt.lat.toStringAsFixed(6)}, '
                        '${pt.lon.toStringAsFixed(6)}',
                        style: const TextStyle(
                          fontFamily: 'monospace',
                          fontSize: 13,
                        ),
                      ),
                      trailing: IconButton(
                        icon: const Icon(Icons.close, size: 18),
                        onPressed: () {
                          setState(() => _boundaryPoints.removeAt(i));
                        },
                      ),
                    );
                  }),
                if (_boundaryPoints.isEmpty)
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    child: Text(
                      'No boundary points added yet.',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _addBoundaryPoint() {
    final lat = double.tryParse(_latController.text.trim());
    final lon = double.tryParse(_lonController.text.trim());
    if (lat == null || lon == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Enter valid latitude and longitude.')),
      );
      return;
    }
    setState(() {
      _boundaryPoints.add((lat: lat, lon: lon));
      _latController.clear();
      _lonController.clear();
    });
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _submitting = true);

    final area = AreaObject(
      name: _nameController.text.trim(),
      description: _descriptionController.text.trim(),
      areaType: _areaType,
      geometry: _geometryController.text.trim(),
    );

    try {
      final created =
          await ref.read(areaNotifierProvider.notifier).createArea(area);
      if (mounted) {
        context.go('/geo/areas/${created.id}');
      }
    } catch (e) {
      if (mounted) {
        setState(() => _submitting = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed: ${friendlyError(e)}')),
        );
      }
    }
  }
}
