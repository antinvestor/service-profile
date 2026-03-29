/// Dart client library for Ant Investor OCR Service.
///
/// Provides OCR service functionality using Connect RPC protocol.
library;

export 'src/ocr/v1/ocr.connect.client.dart';
export 'src/ocr/v1/ocr.connect.spec.dart';
export 'src/ocr/v1/ocr.pb.dart';
export 'src/ocr/v1/ocr.pbenum.dart';
export 'src/ocr/v1/ocr.pbjson.dart';
export 'src/ocr/v1/ocr.pbserver.dart';

// Common types
export 'src/common/v1/common.pb.dart' hide StatusResponse;
export 'src/common/v1/common.pbenum.dart';
export 'src/google/protobuf/struct.pb.dart';
export 'src/google/protobuf/timestamp.pb.dart';
