## 2.0.0

- Regenerated with protoc_plugin 25.1 / connect-dart 2.0; requires protobuf ^6.1.0, connectrpc ^2.0.0, Dart ^3.7. WKTs now come from package:protobuf.
- Requires `antinvestor_api_common` ^2.0.0.
- Removed the vendored `google/protobuf` well-known-type copies and the empty `.pbserver.dart` files that protoc_plugin 25 no longer emits.

## 1.55.0

- `ProfileObject.accounts` (`ProfileAccount`: address, family, version, identity_salt_hash, primary).
- `ResolveAccounts` RPC (`account_resolve`, service role only).

## 1.53.0

- Initial tracked release
