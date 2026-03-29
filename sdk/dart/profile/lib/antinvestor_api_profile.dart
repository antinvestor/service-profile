/// Dart client library for Ant Investor Profile Service.
///
/// Provides Profile service functionality using Connect RPC protocol.
library;

// Export client wrapper
export 'src/client.dart';

// Profile service
export 'src/profile/v1/profile.pb.dart';
export 'src/profile/v1/profile.pbenum.dart';
export 'src/profile/v1/profile.pbjson.dart';
export 'src/profile/v1/profile.connect.client.dart';
export 'src/profile/v1/profile.connect.spec.dart';

// Common types
export 'src/common/v1/common.pb.dart' hide SearchRequest;
export 'src/common/v1/common.pbenum.dart';
export 'src/google/protobuf/struct.pb.dart';
export 'src/google/protobuf/timestamp.pb.dart';
