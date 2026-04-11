/// Geolocation UI library for Antinvestor. Location tracking, areas, routes,
/// events.
library antinvestor_ui_geolocation;

// Providers
export 'src/providers/geolocation_transport_provider.dart';
export 'src/providers/area_providers.dart';
export 'src/providers/route_providers.dart';
export 'src/providers/location_providers.dart';

// Widgets
export 'src/widgets/area_type_badge.dart';
export 'src/widgets/geo_event_tile.dart';
export 'src/widgets/location_point_tile.dart';
export 'src/widgets/route_assignment_chip.dart';
export 'src/widgets/nearby_card.dart';
export 'src/widgets/location_display.dart';
export 'src/widgets/area_badge.dart';

// Screens
export 'src/screens/area_list_screen.dart';
export 'src/screens/area_detail_screen.dart';
export 'src/screens/area_create_screen.dart';
export 'src/screens/area_edit_screen.dart';
export 'src/screens/route_list_screen.dart';
export 'src/screens/route_detail_screen.dart';
export 'src/screens/route_create_screen.dart';
export 'src/screens/location_track_screen.dart';
export 'src/screens/geo_events_screen.dart';

// Routing
export 'src/routing/geolocation_route_module.dart';
