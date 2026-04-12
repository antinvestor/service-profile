/// Device management UI library for Antinvestor.
///
/// Provides embeddable screens, widgets, and Riverpod providers for device
/// registration, key management, session logs, and presence tracking.
library;

// Providers
export 'src/providers/device_transport_provider.dart';
export 'src/providers/device_providers.dart';
export 'src/providers/device_key_providers.dart';
export 'src/providers/device_log_providers.dart';
export 'src/providers/presence_providers.dart';
export 'src/providers/presence_by_profile_provider.dart';

// Widgets
export 'src/widgets/device_card.dart';
export 'src/widgets/presence_indicator.dart';
export 'src/widgets/device_key_tile.dart';
export 'src/widgets/session_log_entry.dart';
export 'src/widgets/device_presence_dot.dart';

// Screens
export 'src/screens/device_list_screen.dart';
export 'src/screens/device_detail_screen.dart';
export 'src/screens/device_keys_screen.dart';
export 'src/screens/device_logs_screen.dart';
export 'src/screens/device_link_screen.dart';

// Routing
export 'src/routing/device_route_module.dart';
