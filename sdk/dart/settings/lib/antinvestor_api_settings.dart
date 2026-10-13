/// Dart client library for Ant Investor Settings Service.
///
/// Provides Settings service functionality using Connect RPC protocol.
library;

export 'src/client.dart';
export 'src/settings/v1/settings.connect.client.dart';
export 'src/settings/v1/settings.connect.spec.dart';
export 'src/settings/v1/settings.pb.dart';
export 'src/settings/v1/settings.pbenum.dart';
export 'src/settings/v1/settings.pbjson.dart';
export 'src/settings/v1/settings.pbserver.dart';

// Common types
export 'src/common/v1/common.pb.dart';
export 'src/common/v1/common.pbenum.dart';
export 'package:protobuf/well_known_types/google/protobuf/struct.pb.dart';
export 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart';
