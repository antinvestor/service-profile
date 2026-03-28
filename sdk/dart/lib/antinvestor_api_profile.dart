/// Dart client library for Ant Investor Profile Service.
///
/// Provides Profile service functionality using Connect RPC protocol.
///
/// For other services in this package, import directly:
/// - `package:antinvestor_api_profile/device.dart` for Device service
/// - `package:antinvestor_api_profile/settings.dart` for Settings service
/// - `package:antinvestor_api_profile/ocr.dart` for OCR service
/// - `package:antinvestor_api_profile/geolocation.dart` for Geolocation service
library;

// Export client wrapper
export 'src/client.dart';

// Profile service
export 'src/profile/v1/profile.pb.dart';
export 'src/profile/v1/profile.pbenum.dart';
export 'src/profile/v1/profile.pbjson.dart';
export 'src/profile/v1/profile.connect.client.dart';
export 'src/profile/v1/profile.connect.spec.dart';

// Export common types used in profile API (hide names that conflict with profile)
export 'src/common/v1/common.pb.dart' hide SearchRequest;
export 'src/common/v1/common.pbenum.dart';
export 'src/google/protobuf/struct.pb.dart';
export 'src/google/protobuf/timestamp.pb.dart';
