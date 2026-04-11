import 'dart:convert';

import 'package:flutter/material.dart';

/// The detected type of a setting value.
enum SettingValueType { text, number, boolean, json }

/// A type-aware editor widget for setting values.
///
/// Automatically detects the value type and shows the appropriate editor:
/// - Boolean: toggle switch
/// - Number: numeric text field
/// - JSON: multi-line code editor with validation
/// - Text: single-line text field (default)
class SettingValueEditor extends StatefulWidget {
  const SettingValueEditor({
    super.key,
    required this.value,
    required this.onChanged,
    this.readOnly = false,
  });

  final String value;
  final ValueChanged<String> onChanged;
  final bool readOnly;

  @override
  State<SettingValueEditor> createState() => _SettingValueEditorState();
}

class _SettingValueEditorState extends State<SettingValueEditor> {
  late TextEditingController _controller;
  late SettingValueType _type;
  String? _jsonError;

  @override
  void initState() {
    super.initState();
    _type = _detectType(widget.value);
    _controller = TextEditingController(text: widget.value);
  }

  @override
  void didUpdateWidget(SettingValueEditor oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.value != widget.value) {
      _type = _detectType(widget.value);
      _controller.text = widget.value;
      _jsonError = null;
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  static SettingValueType _detectType(String value) {
    final trimmed = value.trim().toLowerCase();
    if (trimmed == 'true' || trimmed == 'false') return SettingValueType.boolean;
    if (double.tryParse(trimmed) != null) return SettingValueType.number;
    if (_isJson(value)) return SettingValueType.json;
    return SettingValueType.text;
  }

  static bool _isJson(String value) {
    final trimmed = value.trim();
    if ((!trimmed.startsWith('{') || !trimmed.endsWith('}')) &&
        (!trimmed.startsWith('[') || !trimmed.endsWith(']'))) {
      return false;
    }
    try {
      json.decode(trimmed);
      return true;
    } catch (_) {
      return false;
    }
  }

  void _validateJson(String text) {
    if (_type != SettingValueType.json) {
      setState(() => _jsonError = null);
      return;
    }
    try {
      json.decode(text.trim());
      setState(() => _jsonError = null);
    } on FormatException catch (e) {
      setState(() => _jsonError = 'Invalid JSON: ${e.message}');
    }
  }

  void _onTypeChanged(SettingValueType? newType) {
    if (newType == null || newType == _type) return;
    setState(() {
      _type = newType;
      _jsonError = null;
      // Reset value to a sensible default for the new type.
      switch (newType) {
        case SettingValueType.boolean:
          _controller.text = 'false';
        case SettingValueType.number:
          _controller.text = '0';
        case SettingValueType.json:
          _controller.text = '{}';
        case SettingValueType.text:
          // Keep current value.
          break;
      }
      widget.onChanged(_controller.text);
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Type selector
        Row(
          children: [
            Text(
              'Type:',
              style: theme.textTheme.labelMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(width: 8),
            SegmentedButton<SettingValueType>(
              segments: const [
                ButtonSegment(
                  value: SettingValueType.text,
                  label: Text('Text'),
                  icon: Icon(Icons.text_fields, size: 16),
                ),
                ButtonSegment(
                  value: SettingValueType.number,
                  label: Text('Number'),
                  icon: Icon(Icons.pin_outlined, size: 16),
                ),
                ButtonSegment(
                  value: SettingValueType.boolean,
                  label: Text('Bool'),
                  icon: Icon(Icons.toggle_on_outlined, size: 16),
                ),
                ButtonSegment(
                  value: SettingValueType.json,
                  label: Text('JSON'),
                  icon: Icon(Icons.data_object, size: 16),
                ),
              ],
              selected: {_type},
              onSelectionChanged: widget.readOnly
                  ? null
                  : (set) => _onTypeChanged(set.first),
              showSelectedIcon: false,
            ),
          ],
        ),
        const SizedBox(height: 16),

        // Editor
        _buildEditor(theme),
      ],
    );
  }

  Widget _buildEditor(ThemeData theme) {
    switch (_type) {
      case SettingValueType.boolean:
        return _buildBooleanEditor(theme);
      case SettingValueType.number:
        return _buildNumberEditor(theme);
      case SettingValueType.json:
        return _buildJsonEditor(theme);
      case SettingValueType.text:
        return _buildTextEditor(theme);
    }
  }

  Widget _buildBooleanEditor(ThemeData theme) {
    final isTrue = _controller.text.trim().toLowerCase() == 'true';
    return Row(
      children: [
        Switch(
          value: isTrue,
          onChanged: widget.readOnly
              ? null
              : (val) {
                  final text = val.toString();
                  _controller.text = text;
                  widget.onChanged(text);
                },
        ),
        const SizedBox(width: 8),
        Text(
          isTrue ? 'true' : 'false',
          style: theme.textTheme.bodyLarge?.copyWith(
            fontWeight: FontWeight.w600,
            color: isTrue
                ? theme.colorScheme.primary
                : theme.colorScheme.onSurfaceVariant,
          ),
        ),
      ],
    );
  }

  Widget _buildNumberEditor(ThemeData theme) {
    return TextField(
      controller: _controller,
      readOnly: widget.readOnly,
      keyboardType: const TextInputType.numberWithOptions(decimal: true),
      decoration: InputDecoration(
        labelText: 'Value',
        prefixIcon: const Icon(Icons.pin_outlined),
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
      ),
      onChanged: widget.onChanged,
    );
  }

  Widget _buildJsonEditor(ThemeData theme) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextField(
          controller: _controller,
          readOnly: widget.readOnly,
          maxLines: 12,
          minLines: 4,
          style: theme.textTheme.bodyMedium?.copyWith(
            fontFamily: 'monospace',
            fontSize: 13,
          ),
          decoration: InputDecoration(
            labelText: 'JSON Value',
            alignLabelWithHint: true,
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
            errorText: _jsonError,
            suffixIcon: widget.readOnly
                ? null
                : IconButton(
                    icon: const Icon(Icons.auto_fix_high, size: 20),
                    tooltip: 'Format JSON',
                    onPressed: () {
                      try {
                        final parsed = json.decode(_controller.text.trim());
                        final formatted =
                            const JsonEncoder.withIndent('  ').convert(parsed);
                        _controller.text = formatted;
                        widget.onChanged(formatted);
                        setState(() => _jsonError = null);
                      } on FormatException catch (e) {
                        setState(
                            () => _jsonError = 'Invalid JSON: ${e.message}');
                      }
                    },
                  ),
          ),
          onChanged: (text) {
            _validateJson(text);
            widget.onChanged(text);
          },
        ),
      ],
    );
  }

  Widget _buildTextEditor(ThemeData theme) {
    return TextField(
      controller: _controller,
      readOnly: widget.readOnly,
      maxLines: 3,
      minLines: 1,
      decoration: InputDecoration(
        labelText: 'Value',
        prefixIcon: const Icon(Icons.text_fields),
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
      ),
      onChanged: widget.onChanged,
    );
  }
}
