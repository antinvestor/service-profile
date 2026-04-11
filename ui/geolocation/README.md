# antinvestor_ui_geolocation

Embeddable geolocation UI for Antinvestor applications. Provides screens and widgets for managing areas, routes, location tracking, and geo-events.

## Installation

```yaml
dependencies:
  antinvestor_ui_geolocation: ^0.1.0
```

## Features

- **Area Management**: List, create, edit, and view areas with type badges
- **Route Management**: Define and manage routes with assignment tracking
- **Location Tracking**: Real-time location point display
- **Geo-Events**: Browse and filter geolocation events
- **Embeddable Widgets**: `AreaTypeBadge`, `AreaBadge`, `GeoEventTile`, `LocationPointTile`, `RouteAssignmentChip`, `NearbyCard`, `LocationDisplay`
- **Routing**: `GeolocationRouteModule` with GoRouter integration

## Usage

```dart
import 'package:antinvestor_ui_geolocation/antinvestor_ui_geolocation.dart';

// Display a location inline
LocationDisplay(latitude: -1.286, longitude: 36.817)

// Area type indicator
AreaTypeBadge(type: areaType)

// Register routes in your host app
final module = GeolocationRouteModule();
ShellRoute(
  routes: [...ownRoutes, ...module.buildRoutes()],
);
```

## Routes

| Path | Screen |
|------|--------|
| `/geo/areas` | Area list |
| `/geo/areas/new` | Create area |
| `/geo/areas/:areaId` | Area detail |
| `/geo/areas/:areaId/edit` | Edit area |
| `/geo/routes` | Route list |
| `/geo/routes/new` | Create route |
| `/geo/routes/:routeId` | Route detail |
| `/geo/tracking` | Location tracking |
| `/geo/events` | Geo-events log |

## Embedding Widgets

```dart
// Compact area badge for references
AreaBadge(area: areaObject)

// Nearby location card
NearbyCard(location: locationObject)

// Route assignment indicator
RouteAssignmentChip(route: routeObject)
```
