## 0.3.0

- **Breaking:** moves to connect-dart 2 and the protobuf 6 runtime:
  `antinvestor_ui_core ^0.6.0`, `antinvestor_api_profile ^2.0.0`,
  `antinvestor_api_common ^2.0.0`, `connectrpc ^2.0.0`, `protobuf ^6.1.0`.
  Apps must use the 2.x `antinvestor_api_*` SDKs. No public API changes.
- Struct properties are built from map entries, as protobuf 6 requires.

## 0.2.1

- Add missing image_picker dependency for ProfileResolver

## 0.2.0

- Bump antinvestor_ui_core dependency to ^0.2.0
- Cross-tenant profile reads via SkipTenancyChecksOnClaims
- Append-only property ledger support

## 0.1.1

- Migrate providers from StateNotifier to Notifier (Riverpod 3.x compatibility)
- Fix proto STATE type conflict, SearchRequest import ambiguity
- Fix CommunicationLevel enum value, contact provider return types
- Remove deprecated library name, fix lint warnings

## 0.1.0

- Initial release
- Profile management UI with DataTable, CSV export, contacts, addresses, relationships, roster, merge
