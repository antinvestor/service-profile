import 'package:antinvestor_api_geolocation/antinvestor_api_geolocation.dart';
import 'package:antinvestor_ui_core/widgets/status_badge.dart';
import 'package:flutter/material.dart';

/// Badge widget that displays an area type with a color-coded label.
class AreaTypeBadge extends StatelessWidget {
  const AreaTypeBadge({super.key, required this.areaType});

  final AreaType areaType;

  @override
  Widget build(BuildContext context) {
    return StatusBadge.fromEnum(
      value: areaType,
      mapper: (type) => switch (type) {
        AreaType.AREA_TYPE_LAND => ('Land', Colors.green, Icons.terrain),
        AreaType.AREA_TYPE_BUILDING =>
          ('Building', Colors.blue, Icons.apartment),
        AreaType.AREA_TYPE_ZONE =>
          ('Zone', Colors.purple, Icons.layers_outlined),
        AreaType.AREA_TYPE_FENCE =>
          ('Fence', Colors.orange, Icons.fence_outlined),
        AreaType.AREA_TYPE_CUSTOM =>
          ('Custom', Colors.grey, Icons.category_outlined),
        _ => ('Unknown', Colors.grey, null),
      },
    );
  }
}
