# antinvestor_ui_settings

Embeddable settings management UI for Antinvestor applications. Provides screens and widgets for viewing, searching, and editing application settings with scope-based organization.

## Installation

```yaml
dependencies:
  antinvestor_ui_settings: ^0.1.0
```

## Features

- **Settings List**: Grouped by module with search and filtering
- **Setting Detail**: View and edit individual settings with type-aware editors
- **Bulk Edit**: Edit multiple settings at once
- **Scope Selection**: Filter settings by scope (global, tenant, partition)
- **Embeddable Widgets**: `SettingTile`, `SettingValueEditor`, `SettingsScopeSelector`, `SettingValueWidget`
- **Routing**: `SettingsRouteModule` with GoRouter integration

## Usage

```dart
import 'package:antinvestor_ui_settings/antinvestor_ui_settings.dart';

// Display a setting value inline
SettingValueWidget(settingKey: 'app.theme.mode')

// Scope selector for filtering
SettingsScopeSelector(
  onScopeChanged: (scope) => print(scope),
)

// Register routes in your host app
final module = SettingsRouteModule();
ShellRoute(
  routes: [...ownRoutes, ...module.buildRoutes()],
);
```

## Routes

| Path | Screen |
|------|--------|
| `/settings` | Settings list (grouped, searchable) |
| `/settings/detail/:key` | View/edit a single setting |
| `/settings/bulk-edit` | Bulk editor for multiple settings |

## Embedding Widgets

```dart
// Setting row with value preview
SettingTile(setting: settingObject)

// Type-aware value editor (string, number, boolean, JSON)
SettingValueEditor(setting: settingObject, onSaved: (v) {})
```
