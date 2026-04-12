import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/error_helpers.dart';
import 'package:antinvestor_ui_core/widgets/form_field_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/area_providers.dart';

/// Screen for editing an existing area: name, description, type, and geometry.
class AreaEditScreen extends ConsumerStatefulWidget {
  const AreaEditScreen({super.key, required this.areaId});

  final String areaId;

  @override
  ConsumerState<AreaEditScreen> createState() => _AreaEditScreenState();
}

class _AreaEditScreenState extends ConsumerState<AreaEditScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _geometryController = TextEditingController();
  AreaType _areaType = AreaType.AREA_TYPE_LAND;
  bool _submitting = false;
  bool _initialized = false;

  void _initFromArea(AreaObject area) {
    if (_initialized) return;
    _initialized = true;
    _nameController.text = area.name;
    _descriptionController.text = area.description;
    _geometryController.text = area.geometry;
    _areaType = area.areaType;
  }

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    _geometryController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final asyncArea = ref.watch(getAreaProvider(widget.areaId));

    return asyncArea.when(
      loading: () => const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      ),
      error: (error, _) => Scaffold(
        appBar: AppBar(title: const Text('Edit Area')),
        body: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.error_outline,
                  size: 48,
                  color: Theme.of(context).colorScheme.error),
              const SizedBox(height: 16),
              Text(friendlyError(error)),
              const SizedBox(height: 16),
              FilledButton.tonal(
                onPressed: () =>
                    ref.invalidate(getAreaProvider(widget.areaId)),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
      data: (area) {
        _initFromArea(area);
        return _buildForm(area);
      },
    );
  }

  Widget _buildForm(AreaObject area) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Edit ${area.name.isNotEmpty ? area.name : 'Area'}'),
        actions: [
          FilledButton.icon(
            onPressed: _submitting ? null : () => _submit(area),
            icon: _submitting
                ? const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.check, size: 18),
            label: const Text('Save'),
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
                decoration:
                    const InputDecoration(hintText: 'Enter area name'),
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
              description: 'Paste GeoJSON geometry for the area boundary.',
              child: TextFormField(
                controller: _geometryController,
                decoration: const InputDecoration(
                    hintText: '{"type":"Polygon",...}'),
                maxLines: 4,
                style:
                    const TextStyle(fontFamily: 'monospace', fontSize: 12),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _submit(AreaObject original) async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _submitting = true);

    final request = UpdateAreaRequest()
      ..id = widget.areaId
      ..name = _nameController.text.trim()
      ..description = _descriptionController.text.trim()
      ..areaType = _areaType
      ..geometry = _geometryController.text.trim();

    try {
      await ref.read(areaNotifierProvider.notifier).updateArea(request);
      if (mounted) {
        ref.invalidate(getAreaProvider(widget.areaId));
        context.go('/geo/areas/${widget.areaId}');
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
