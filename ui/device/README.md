# antinvestor_ui_device

Embeddable device management UI for Antinvestor applications. Provides screens and widgets for device registration, key management, session logs, and presence tracking.

## Installation

```yaml
dependencies:
  antinvestor_ui_device: ^0.1.0
```

## Features

- **Device List**: Searchable device list with status indicators
- **Device Detail**: Tabbed view with device info, keys, and session logs
- **Key Management**: View and manage device cryptographic keys
- **Session Logs**: Browse device session history
- **Device Linking**: Register or link new devices
- **Presence Tracking**: Real-time online/offline indicators
- **Embeddable Widgets**: `DeviceCard`, `PresenceIndicator`, `DevicePresenceDot`, `DeviceKeyTile`, `SessionLogEntry`
- **Routing**: `DeviceRouteModule` with GoRouter integration

## Usage

```dart
import 'package:antinvestor_ui_device/antinvestor_ui_device.dart';

// Show presence status for a device
DevicePresenceDot(deviceId: 'device-xyz')

// Presence indicator with label
PresenceIndicator(deviceId: 'device-xyz')

// Register routes in your host app
final module = DeviceRouteModule();
ShellRoute(
  routes: [...ownRoutes, ...module.buildRoutes()],
);
```

## Routes

| Path | Screen |
|------|--------|
| `/devices` | Device list with search |
| `/devices/detail/:id` | Device detail (Info, Keys, Sessions) |
| `/devices/link` | Register or link a device |

## Embedding Widgets

```dart
// Device summary card
DeviceCard(device: deviceObject)

// Cryptographic key tile
DeviceKeyTile(key: keyObject)

// Session log row
SessionLogEntry(log: logObject)
```
