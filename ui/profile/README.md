# antinvestor_ui_profile

Embeddable profile management UI for Antinvestor applications. Provides screens and widgets for searching, viewing, editing, and managing profiles, contacts, addresses, relationships, and roster entries.

## Installation

```yaml
dependencies:
  antinvestor_ui_profile: ^0.1.0
```

## Features

- **Profile Search**: Paginated DataTable with filters and CSV export
- **Profile Detail**: Multi-tab view (contacts, addresses, relationships, roster)
- **Profile CRUD**: Create, edit, and merge profiles
- **Analytics**: Profile analytics dashboard with KPI cards
- **Embeddable Widgets**: `ProfileBadgeById`, `ProfileSearchSelect`, `ProfileCard`, `ProfileTypeBadge`
- **Contact Management**: Contact list with verification dialogs
- **Address Management**: Address tiles with inline editing
- **Relationships**: Relationship chips with type indicators
- **Routing**: `ProfileRouteModule` with GoRouter integration

## Usage

```dart
import 'package:antinvestor_ui_profile/antinvestor_ui_profile.dart';

// Embed a profile badge anywhere in your app
ProfileBadgeById(profileId: 'abc123')

// Profile search/select dropdown
ProfileSearchSelect(
  onSelected: (profile) => print(profile.id),
)

// Register routes in your host app
final module = ProfileRouteModule();
ShellRoute(
  routes: [...ownRoutes, ...module.buildRoutes()],
);
```

## Routes

| Path | Screen |
|------|--------|
| `/profiles` | Profile search list |
| `/profiles/new` | Create profile |
| `/profiles/analytics` | Profile analytics |
| `/profiles/merge` | Merge profiles |
| `/profiles/:profileId` | Profile detail |
| `/profiles/:profileId/edit` | Edit profile |
| `/profiles/:profileId/contacts` | Contacts tab |
| `/profiles/:profileId/addresses` | Addresses tab |
| `/profiles/:profileId/relationships` | Relationships tab |
| `/profiles/:profileId/roster` | Roster tab |

## Embedding Widgets

```dart
// Compact profile card for lists
ProfileCard(profile: profileObject)

// Type indicator badge
ProfileTypeBadge(type: ProfileType.INDIVIDUAL)

// Contact row with verification status
ContactListTile(contact: contactObject)
```
