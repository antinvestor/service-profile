/// Dart client library for Ant Investor Geolocation Service.
///
/// Provides Geolocation service functionality using Connect RPC protocol.
library;

export 'src/client.dart';
export 'src/geolocation/v1/geolocation.connect.client.dart';
export 'src/geolocation/v1/geolocation.connect.spec.dart';
export 'src/geolocation/v1/geolocation.pb.dart';
export 'src/geolocation/v1/geolocation.pbenum.dart';
export 'src/geolocation/v1/geolocation.pbjson.dart';
export 'src/geolocation/v1/geolocation.pbserver.dart';

// Well-known types
export 'package:protobuf/well_known_types/google/protobuf/struct.pb.dart';
export 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart';
