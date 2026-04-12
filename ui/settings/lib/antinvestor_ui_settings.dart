/// Settings management UI library for Antinvestor.
///
/// Provides embeddable screens, widgets, and Riverpod providers for viewing,
/// searching, and editing application settings.
library;

// Providers
export 'src/providers/settings_transport_provider.dart';
export 'src/providers/settings_providers.dart';

// Widgets
export 'src/widgets/setting_tile.dart';
export 'src/widgets/setting_value_editor.dart';
export 'src/widgets/settings_scope_selector.dart';
export 'src/widgets/setting_value_widget.dart';

// Screens
export 'src/screens/settings_list_screen.dart';
export 'src/screens/setting_detail_screen.dart';
export 'src/screens/settings_bulk_edit_screen.dart';

// Routing
export 'src/routing/settings_route_module.dart';
