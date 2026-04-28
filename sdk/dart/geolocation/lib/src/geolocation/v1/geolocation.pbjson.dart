//
//  Generated code. Do not modify.
//  source: geolocation/v1/geolocation.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import '../../google/protobuf/empty.pbjson.dart' as $2;
import '../../google/protobuf/struct.pbjson.dart' as $1;
import '../../google/protobuf/timestamp.pbjson.dart' as $0;

@$core.Deprecated('Use locationSourceDescriptor instead')
const LocationSource$json = {
  '1': 'LocationSource',
  '2': [
    {'1': 'LOCATION_SOURCE_UNSPECIFIED', '2': 0},
    {'1': 'LOCATION_SOURCE_GPS', '2': 1},
    {'1': 'LOCATION_SOURCE_NETWORK', '2': 2},
    {'1': 'LOCATION_SOURCE_IP', '2': 3},
    {'1': 'LOCATION_SOURCE_MANUAL', '2': 4},
  ],
};

/// Descriptor for `LocationSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List locationSourceDescriptor = $convert.base64Decode(
    'Cg5Mb2NhdGlvblNvdXJjZRIfChtMT0NBVElPTl9TT1VSQ0VfVU5TUEVDSUZJRUQQABIXChNMT0'
    'NBVElPTl9TT1VSQ0VfR1BTEAESGwoXTE9DQVRJT05fU09VUkNFX05FVFdPUksQAhIWChJMT0NB'
    'VElPTl9TT1VSQ0VfSVAQAxIaChZMT0NBVElPTl9TT1VSQ0VfTUFOVUFMEAQ=');

@$core.Deprecated('Use areaTypeDescriptor instead')
const AreaType$json = {
  '1': 'AreaType',
  '2': [
    {'1': 'AREA_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'AREA_TYPE_LAND', '2': 1},
    {'1': 'AREA_TYPE_BUILDING', '2': 2},
    {'1': 'AREA_TYPE_ZONE', '2': 3},
    {'1': 'AREA_TYPE_FENCE', '2': 4},
    {'1': 'AREA_TYPE_CUSTOM', '2': 5},
  ],
};

/// Descriptor for `AreaType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List areaTypeDescriptor = $convert.base64Decode(
    'CghBcmVhVHlwZRIZChVBUkVBX1RZUEVfVU5TUEVDSUZJRUQQABISCg5BUkVBX1RZUEVfTEFORB'
    'ABEhYKEkFSRUFfVFlQRV9CVUlMRElORxACEhIKDkFSRUFfVFlQRV9aT05FEAMSEwoPQVJFQV9U'
    'WVBFX0ZFTkNFEAQSFAoQQVJFQV9UWVBFX0NVU1RPTRAF');

@$core.Deprecated('Use geoEventTypeDescriptor instead')
const GeoEventType$json = {
  '1': 'GeoEventType',
  '2': [
    {'1': 'GEO_EVENT_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'GEO_EVENT_TYPE_ENTER', '2': 1},
    {'1': 'GEO_EVENT_TYPE_EXIT', '2': 2},
    {'1': 'GEO_EVENT_TYPE_DWELL', '2': 3},
  ],
};

/// Descriptor for `GeoEventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List geoEventTypeDescriptor = $convert.base64Decode(
    'CgxHZW9FdmVudFR5cGUSHgoaR0VPX0VWRU5UX1RZUEVfVU5TUEVDSUZJRUQQABIYChRHRU9fRV'
    'ZFTlRfVFlQRV9FTlRFUhABEhcKE0dFT19FVkVOVF9UWVBFX0VYSVQQAhIYChRHRU9fRVZFTlRf'
    'VFlQRV9EV0VMTBAD');

@$core.Deprecated('Use routeDeviationEventTypeDescriptor instead')
const RouteDeviationEventType$json = {
  '1': 'RouteDeviationEventType',
  '2': [
    {'1': 'ROUTE_DEVIATION_EVENT_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'ROUTE_DEVIATION_EVENT_TYPE_DEVIATED', '2': 1},
    {'1': 'ROUTE_DEVIATION_EVENT_TYPE_BACK_ON_ROUTE', '2': 2},
  ],
};

/// Descriptor for `RouteDeviationEventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List routeDeviationEventTypeDescriptor = $convert.base64Decode(
    'ChdSb3V0ZURldmlhdGlvbkV2ZW50VHlwZRIqCiZST1VURV9ERVZJQVRJT05fRVZFTlRfVFlQRV'
    '9VTlNQRUNJRklFRBAAEicKI1JPVVRFX0RFVklBVElPTl9FVkVOVF9UWVBFX0RFVklBVEVEEAES'
    'LAooUk9VVEVfREVWSUFUSU9OX0VWRU5UX1RZUEVfQkFDS19PTl9ST1VURRAC');

@$core.Deprecated('Use locationPointInputDescriptor instead')
const LocationPointInput$json = {
  '1': 'LocationPointInput',
  '2': [
    {'1': 'timestamp', '3': 1, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'timestamp'},
    {'1': 'latitude', '3': 2, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 3, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'altitude', '3': 4, '4': 1, '5': 1, '9': 0, '10': 'altitude', '17': true},
    {'1': 'accuracy', '3': 5, '4': 1, '5': 1, '10': 'accuracy'},
    {'1': 'speed', '3': 6, '4': 1, '5': 1, '9': 1, '10': 'speed', '17': true},
    {'1': 'bearing', '3': 7, '4': 1, '5': 1, '9': 2, '10': 'bearing', '17': true},
    {'1': 'source', '3': 8, '4': 1, '5': 14, '6': '.geolocation.v1.LocationSource', '10': 'source'},
    {'1': 'extra', '3': 9, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
    {'1': 'device_id', '3': 10, '4': 1, '5': 9, '10': 'deviceId'},
  ],
  '8': [
    {'1': '_altitude'},
    {'1': '_speed'},
    {'1': '_bearing'},
  ],
};

/// Descriptor for `LocationPointInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locationPointInputDescriptor = $convert.base64Decode(
    'ChJMb2NhdGlvblBvaW50SW5wdXQSOAoJdGltZXN0YW1wGAEgASgLMhouZ29vZ2xlLnByb3RvYn'
    'VmLlRpbWVzdGFtcFIJdGltZXN0YW1wEhoKCGxhdGl0dWRlGAIgASgBUghsYXRpdHVkZRIcCgls'
    'b25naXR1ZGUYAyABKAFSCWxvbmdpdHVkZRIfCghhbHRpdHVkZRgEIAEoAUgAUghhbHRpdHVkZY'
    'gBARIaCghhY2N1cmFjeRgFIAEoAVIIYWNjdXJhY3kSGQoFc3BlZWQYBiABKAFIAVIFc3BlZWSI'
    'AQESHQoHYmVhcmluZxgHIAEoAUgCUgdiZWFyaW5niAEBEjYKBnNvdXJjZRgIIAEoDjIeLmdlb2'
    'xvY2F0aW9uLnYxLkxvY2F0aW9uU291cmNlUgZzb3VyY2USLQoFZXh0cmEYCSABKAsyFy5nb29n'
    'bGUucHJvdG9idWYuU3RydWN0UgVleHRyYRIbCglkZXZpY2VfaWQYCiABKAlSCGRldmljZUlkQg'
    'sKCV9hbHRpdHVkZUIICgZfc3BlZWRCCgoIX2JlYXJpbmc=');

@$core.Deprecated('Use locationPointObjectDescriptor instead')
const LocationPointObject$json = {
  '1': 'LocationPointObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'subject_id', '3': 2, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'timestamp', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'timestamp'},
    {'1': 'latitude', '3': 4, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 5, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'altitude', '3': 6, '4': 1, '5': 1, '9': 0, '10': 'altitude', '17': true},
    {'1': 'accuracy', '3': 7, '4': 1, '5': 1, '10': 'accuracy'},
    {'1': 'speed', '3': 8, '4': 1, '5': 1, '9': 1, '10': 'speed', '17': true},
    {'1': 'bearing', '3': 9, '4': 1, '5': 1, '9': 2, '10': 'bearing', '17': true},
    {'1': 'source', '3': 10, '4': 1, '5': 14, '6': '.geolocation.v1.LocationSource', '10': 'source'},
    {'1': 'extra', '3': 11, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
    {'1': 'created_at', '3': 12, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'createdAt'},
    {'1': 'device_id', '3': 13, '4': 1, '5': 9, '10': 'deviceId'},
  ],
  '8': [
    {'1': '_altitude'},
    {'1': '_speed'},
    {'1': '_bearing'},
  ],
};

/// Descriptor for `LocationPointObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locationPointObjectDescriptor = $convert.base64Decode(
    'ChNMb2NhdGlvblBvaW50T2JqZWN0Eg4KAmlkGAEgASgJUgJpZBIdCgpzdWJqZWN0X2lkGAIgAS'
    'gJUglzdWJqZWN0SWQSOAoJdGltZXN0YW1wGAMgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVz'
    'dGFtcFIJdGltZXN0YW1wEhoKCGxhdGl0dWRlGAQgASgBUghsYXRpdHVkZRIcCglsb25naXR1ZG'
    'UYBSABKAFSCWxvbmdpdHVkZRIfCghhbHRpdHVkZRgGIAEoAUgAUghhbHRpdHVkZYgBARIaCghh'
    'Y2N1cmFjeRgHIAEoAVIIYWNjdXJhY3kSGQoFc3BlZWQYCCABKAFIAVIFc3BlZWSIAQESHQoHYm'
    'VhcmluZxgJIAEoAUgCUgdiZWFyaW5niAEBEjYKBnNvdXJjZRgKIAEoDjIeLmdlb2xvY2F0aW9u'
    'LnYxLkxvY2F0aW9uU291cmNlUgZzb3VyY2USLQoFZXh0cmEYCyABKAsyFy5nb29nbGUucHJvdG'
    '9idWYuU3RydWN0UgVleHRyYRI5CgpjcmVhdGVkX2F0GAwgASgLMhouZ29vZ2xlLnByb3RvYnVm'
    'LlRpbWVzdGFtcFIJY3JlYXRlZEF0EhsKCWRldmljZV9pZBgNIAEoCVIIZGV2aWNlSWRCCwoJX2'
    'FsdGl0dWRlQggKBl9zcGVlZEIKCghfYmVhcmluZw==');

@$core.Deprecated('Use areaObjectDescriptor instead')
const AreaObject$json = {
  '1': 'AreaObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'owner_id', '3': 2, '4': 1, '5': 9, '10': 'ownerId'},
    {'1': 'name', '3': 3, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'area_type', '3': 5, '4': 1, '5': 14, '6': '.geolocation.v1.AreaType', '10': 'areaType'},
    {'1': 'geometry', '3': 6, '4': 1, '5': 9, '10': 'geometry'},
    {'1': 'area_m2', '3': 7, '4': 1, '5': 1, '10': 'areaM2'},
    {'1': 'perimeter_m', '3': 8, '4': 1, '5': 1, '10': 'perimeterM'},
    {'1': 'state', '3': 9, '4': 1, '5': 5, '10': 'state'},
    {'1': 'extra', '3': 10, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
    {'1': 'created_at', '3': 11, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'createdAt'},
  ],
};

/// Descriptor for `AreaObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List areaObjectDescriptor = $convert.base64Decode(
    'CgpBcmVhT2JqZWN0Eg4KAmlkGAEgASgJUgJpZBIZCghvd25lcl9pZBgCIAEoCVIHb3duZXJJZB'
    'ISCgRuYW1lGAMgASgJUgRuYW1lEiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcmlwdGlvbhI1'
    'CglhcmVhX3R5cGUYBSABKA4yGC5nZW9sb2NhdGlvbi52MS5BcmVhVHlwZVIIYXJlYVR5cGUSGg'
    'oIZ2VvbWV0cnkYBiABKAlSCGdlb21ldHJ5EhcKB2FyZWFfbTIYByABKAFSBmFyZWFNMhIfCgtw'
    'ZXJpbWV0ZXJfbRgIIAEoAVIKcGVyaW1ldGVyTRIUCgVzdGF0ZRgJIAEoBVIFc3RhdGUSLQoFZX'
    'h0cmEYCiABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0UgVleHRyYRI5CgpjcmVhdGVkX2F0'
    'GAsgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJY3JlYXRlZEF0');

@$core.Deprecated('Use geoEventObjectDescriptor instead')
const GeoEventObject$json = {
  '1': 'GeoEventObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'subject_id', '3': 2, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'area_id', '3': 3, '4': 1, '5': 9, '10': 'areaId'},
    {'1': 'event_type', '3': 4, '4': 1, '5': 14, '6': '.geolocation.v1.GeoEventType', '10': 'eventType'},
    {'1': 'timestamp', '3': 5, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'timestamp'},
    {'1': 'confidence', '3': 6, '4': 1, '5': 1, '10': 'confidence'},
    {'1': 'point_id', '3': 7, '4': 1, '5': 9, '10': 'pointId'},
    {'1': 'extra', '3': 8, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
  ],
};

/// Descriptor for `GeoEventObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List geoEventObjectDescriptor = $convert.base64Decode(
    'Cg5HZW9FdmVudE9iamVjdBIOCgJpZBgBIAEoCVICaWQSHQoKc3ViamVjdF9pZBgCIAEoCVIJc3'
    'ViamVjdElkEhcKB2FyZWFfaWQYAyABKAlSBmFyZWFJZBI7CgpldmVudF90eXBlGAQgASgOMhwu'
    'Z2VvbG9jYXRpb24udjEuR2VvRXZlbnRUeXBlUglldmVudFR5cGUSOAoJdGltZXN0YW1wGAUgAS'
    'gLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJdGltZXN0YW1wEh4KCmNvbmZpZGVuY2UY'
    'BiABKAFSCmNvbmZpZGVuY2USGQoIcG9pbnRfaWQYByABKAlSB3BvaW50SWQSLQoFZXh0cmEYCC'
    'ABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0UgVleHRyYQ==');

@$core.Deprecated('Use areaSubjectObjectDescriptor instead')
const AreaSubjectObject$json = {
  '1': 'AreaSubjectObject',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'enter_timestamp', '3': 2, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'enterTimestamp'},
  ],
};

/// Descriptor for `AreaSubjectObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List areaSubjectObjectDescriptor = $convert.base64Decode(
    'ChFBcmVhU3ViamVjdE9iamVjdBIdCgpzdWJqZWN0X2lkGAEgASgJUglzdWJqZWN0SWQSQwoPZW'
    '50ZXJfdGltZXN0YW1wGAIgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIOZW50ZXJU'
    'aW1lc3RhbXA=');

@$core.Deprecated('Use nearbySubjectObjectDescriptor instead')
const NearbySubjectObject$json = {
  '1': 'NearbySubjectObject',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'distance_meters', '3': 2, '4': 1, '5': 1, '10': 'distanceMeters'},
    {'1': 'last_seen', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'lastSeen'},
  ],
};

/// Descriptor for `NearbySubjectObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nearbySubjectObjectDescriptor = $convert.base64Decode(
    'ChNOZWFyYnlTdWJqZWN0T2JqZWN0Eh0KCnN1YmplY3RfaWQYASABKAlSCXN1YmplY3RJZBInCg'
    '9kaXN0YW5jZV9tZXRlcnMYAiABKAFSDmRpc3RhbmNlTWV0ZXJzEjcKCWxhc3Rfc2VlbhgDIAEo'
    'CzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCGxhc3RTZWVu');

@$core.Deprecated('Use nearbyAreaObjectDescriptor instead')
const NearbyAreaObject$json = {
  '1': 'NearbyAreaObject',
  '2': [
    {'1': 'area_id', '3': 1, '4': 1, '5': 9, '10': 'areaId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'area_type', '3': 3, '4': 1, '5': 14, '6': '.geolocation.v1.AreaType', '10': 'areaType'},
    {'1': 'distance_meters', '3': 4, '4': 1, '5': 1, '10': 'distanceMeters'},
  ],
};

/// Descriptor for `NearbyAreaObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nearbyAreaObjectDescriptor = $convert.base64Decode(
    'ChBOZWFyYnlBcmVhT2JqZWN0EhcKB2FyZWFfaWQYASABKAlSBmFyZWFJZBISCgRuYW1lGAIgAS'
    'gJUgRuYW1lEjUKCWFyZWFfdHlwZRgDIAEoDjIYLmdlb2xvY2F0aW9uLnYxLkFyZWFUeXBlUghh'
    'cmVhVHlwZRInCg9kaXN0YW5jZV9tZXRlcnMYBCABKAFSDmRpc3RhbmNlTWV0ZXJz');

@$core.Deprecated('Use routeObjectDescriptor instead')
const RouteObject$json = {
  '1': 'RouteObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'owner_id', '3': 2, '4': 1, '5': 9, '10': 'ownerId'},
    {'1': 'name', '3': 3, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'geometry', '3': 5, '4': 1, '5': 9, '10': 'geometry'},
    {'1': 'length_m', '3': 6, '4': 1, '5': 1, '10': 'lengthM'},
    {'1': 'state', '3': 7, '4': 1, '5': 5, '10': 'state'},
    {'1': 'deviation_threshold_m', '3': 8, '4': 1, '5': 1, '9': 0, '10': 'deviationThresholdM', '17': true},
    {'1': 'deviation_consecutive_count', '3': 9, '4': 1, '5': 5, '9': 1, '10': 'deviationConsecutiveCount', '17': true},
    {'1': 'deviation_cooldown_sec', '3': 10, '4': 1, '5': 5, '9': 2, '10': 'deviationCooldownSec', '17': true},
    {'1': 'extra', '3': 11, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
    {'1': 'created_at', '3': 12, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'createdAt'},
  ],
  '8': [
    {'1': '_deviation_threshold_m'},
    {'1': '_deviation_consecutive_count'},
    {'1': '_deviation_cooldown_sec'},
  ],
};

/// Descriptor for `RouteObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routeObjectDescriptor = $convert.base64Decode(
    'CgtSb3V0ZU9iamVjdBIOCgJpZBgBIAEoCVICaWQSGQoIb3duZXJfaWQYAiABKAlSB293bmVySW'
    'QSEgoEbmFtZRgDIAEoCVIEbmFtZRIgCgtkZXNjcmlwdGlvbhgEIAEoCVILZGVzY3JpcHRpb24S'
    'GgoIZ2VvbWV0cnkYBSABKAlSCGdlb21ldHJ5EhkKCGxlbmd0aF9tGAYgASgBUgdsZW5ndGhNEh'
    'QKBXN0YXRlGAcgASgFUgVzdGF0ZRI3ChVkZXZpYXRpb25fdGhyZXNob2xkX20YCCABKAFIAFIT'
    'ZGV2aWF0aW9uVGhyZXNob2xkTYgBARJDChtkZXZpYXRpb25fY29uc2VjdXRpdmVfY291bnQYCS'
    'ABKAVIAVIZZGV2aWF0aW9uQ29uc2VjdXRpdmVDb3VudIgBARI5ChZkZXZpYXRpb25fY29vbGRv'
    'd25fc2VjGAogASgFSAJSFGRldmlhdGlvbkNvb2xkb3duU2VjiAEBEi0KBWV4dHJhGAsgASgLMh'
    'cuZ29vZ2xlLnByb3RvYnVmLlN0cnVjdFIFZXh0cmESOQoKY3JlYXRlZF9hdBgMIAEoCzIaLmdv'
    'b2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdEIYChZfZGV2aWF0aW9uX3RocmVzaG'
    '9sZF9tQh4KHF9kZXZpYXRpb25fY29uc2VjdXRpdmVfY291bnRCGQoXX2RldmlhdGlvbl9jb29s'
    'ZG93bl9zZWM=');

@$core.Deprecated('Use routeAssignmentObjectDescriptor instead')
const RouteAssignmentObject$json = {
  '1': 'RouteAssignmentObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'subject_id', '3': 2, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'route_id', '3': 3, '4': 1, '5': 9, '10': 'routeId'},
    {'1': 'valid_from', '3': 4, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'validFrom'},
    {'1': 'valid_until', '3': 5, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'validUntil'},
    {'1': 'state', '3': 6, '4': 1, '5': 5, '10': 'state'},
    {'1': 'extra', '3': 7, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
    {'1': 'created_at', '3': 8, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'createdAt'},
  ],
};

/// Descriptor for `RouteAssignmentObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routeAssignmentObjectDescriptor = $convert.base64Decode(
    'ChVSb3V0ZUFzc2lnbm1lbnRPYmplY3QSDgoCaWQYASABKAlSAmlkEh0KCnN1YmplY3RfaWQYAi'
    'ABKAlSCXN1YmplY3RJZBIZCghyb3V0ZV9pZBgDIAEoCVIHcm91dGVJZBI5Cgp2YWxpZF9mcm9t'
    'GAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJdmFsaWRGcm9tEjsKC3ZhbGlkX3'
    'VudGlsGAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIKdmFsaWRVbnRpbBIUCgVz'
    'dGF0ZRgGIAEoBVIFc3RhdGUSLQoFZXh0cmEYByABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydW'
    'N0UgVleHRyYRI5CgpjcmVhdGVkX2F0GAggASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFt'
    'cFIJY3JlYXRlZEF0');

@$core.Deprecated('Use routeDeviationEventObjectDescriptor instead')
const RouteDeviationEventObject$json = {
  '1': 'RouteDeviationEventObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'subject_id', '3': 2, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'route_id', '3': 3, '4': 1, '5': 9, '10': 'routeId'},
    {'1': 'event_type', '3': 4, '4': 1, '5': 14, '6': '.geolocation.v1.RouteDeviationEventType', '10': 'eventType'},
    {'1': 'distance_meters', '3': 5, '4': 1, '5': 1, '10': 'distanceMeters'},
    {'1': 'latitude', '3': 6, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 7, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'timestamp', '3': 8, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'timestamp'},
    {'1': 'extra', '3': 9, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
  ],
};

/// Descriptor for `RouteDeviationEventObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routeDeviationEventObjectDescriptor = $convert.base64Decode(
    'ChlSb3V0ZURldmlhdGlvbkV2ZW50T2JqZWN0Eg4KAmlkGAEgASgJUgJpZBIdCgpzdWJqZWN0X2'
    'lkGAIgASgJUglzdWJqZWN0SWQSGQoIcm91dGVfaWQYAyABKAlSB3JvdXRlSWQSRgoKZXZlbnRf'
    'dHlwZRgEIAEoDjInLmdlb2xvY2F0aW9uLnYxLlJvdXRlRGV2aWF0aW9uRXZlbnRUeXBlUglldm'
    'VudFR5cGUSJwoPZGlzdGFuY2VfbWV0ZXJzGAUgASgBUg5kaXN0YW5jZU1ldGVycxIaCghsYXRp'
    'dHVkZRgGIAEoAVIIbGF0aXR1ZGUSHAoJbG9uZ2l0dWRlGAcgASgBUglsb25naXR1ZGUSOAoJdG'
    'ltZXN0YW1wGAggASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJdGltZXN0YW1wEi0K'
    'BWV4dHJhGAkgASgLMhcuZ29vZ2xlLnByb3RvYnVmLlN0cnVjdFIFZXh0cmE=');

@$core.Deprecated('Use ingestLocationsRequestDescriptor instead')
const IngestLocationsRequest$json = {
  '1': 'IngestLocationsRequest',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'points', '3': 2, '4': 3, '5': 11, '6': '.geolocation.v1.LocationPointInput', '10': 'points'},
  ],
};

/// Descriptor for `IngestLocationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ingestLocationsRequestDescriptor = $convert.base64Decode(
    'ChZJbmdlc3RMb2NhdGlvbnNSZXF1ZXN0Eh0KCnN1YmplY3RfaWQYASABKAlSCXN1YmplY3RJZB'
    'I6CgZwb2ludHMYAiADKAsyIi5nZW9sb2NhdGlvbi52MS5Mb2NhdGlvblBvaW50SW5wdXRSBnBv'
    'aW50cw==');

@$core.Deprecated('Use ingestLocationsResponseDescriptor instead')
const IngestLocationsResponse$json = {
  '1': 'IngestLocationsResponse',
  '2': [
    {'1': 'accepted', '3': 1, '4': 1, '5': 5, '10': 'accepted'},
    {'1': 'rejected', '3': 2, '4': 1, '5': 5, '10': 'rejected'},
  ],
};

/// Descriptor for `IngestLocationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ingestLocationsResponseDescriptor = $convert.base64Decode(
    'ChdJbmdlc3RMb2NhdGlvbnNSZXNwb25zZRIaCghhY2NlcHRlZBgBIAEoBVIIYWNjZXB0ZWQSGg'
    'oIcmVqZWN0ZWQYAiABKAVSCHJlamVjdGVk');

@$core.Deprecated('Use createAreaRequestDescriptor instead')
const CreateAreaRequest$json = {
  '1': 'CreateAreaRequest',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.AreaObject', '10': 'data'},
  ],
};

/// Descriptor for `CreateAreaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAreaRequestDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVBcmVhUmVxdWVzdBIuCgRkYXRhGAEgASgLMhouZ2VvbG9jYXRpb24udjEuQXJlYU'
    '9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use createAreaResponseDescriptor instead')
const CreateAreaResponse$json = {
  '1': 'CreateAreaResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.AreaObject', '10': 'data'},
  ],
};

/// Descriptor for `CreateAreaResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAreaResponseDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVBcmVhUmVzcG9uc2USLgoEZGF0YRgBIAEoCzIaLmdlb2xvY2F0aW9uLnYxLkFyZW'
    'FPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use getAreaRequestDescriptor instead')
const GetAreaRequest$json = {
  '1': 'GetAreaRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetAreaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAreaRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRBcmVhUmVxdWVzdBIOCgJpZBgBIAEoCVICaWQ=');

@$core.Deprecated('Use getAreaResponseDescriptor instead')
const GetAreaResponse$json = {
  '1': 'GetAreaResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.AreaObject', '10': 'data'},
  ],
};

/// Descriptor for `GetAreaResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAreaResponseDescriptor = $convert.base64Decode(
    'Cg9HZXRBcmVhUmVzcG9uc2USLgoEZGF0YRgBIAEoCzIaLmdlb2xvY2F0aW9uLnYxLkFyZWFPYm'
    'plY3RSBGRhdGE=');

@$core.Deprecated('Use updateAreaRequestDescriptor instead')
const UpdateAreaRequest$json = {
  '1': 'UpdateAreaRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '9': 0, '10': 'name', '17': true},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '9': 1, '10': 'description', '17': true},
    {'1': 'area_type', '3': 4, '4': 1, '5': 14, '6': '.geolocation.v1.AreaType', '9': 2, '10': 'areaType', '17': true},
    {'1': 'geometry', '3': 5, '4': 1, '5': 9, '9': 3, '10': 'geometry', '17': true},
    {'1': 'extra', '3': 6, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
  ],
  '8': [
    {'1': '_name'},
    {'1': '_description'},
    {'1': '_area_type'},
    {'1': '_geometry'},
  ],
};

/// Descriptor for `UpdateAreaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAreaRequestDescriptor = $convert.base64Decode(
    'ChFVcGRhdGVBcmVhUmVxdWVzdBIOCgJpZBgBIAEoCVICaWQSFwoEbmFtZRgCIAEoCUgAUgRuYW'
    '1liAEBEiUKC2Rlc2NyaXB0aW9uGAMgASgJSAFSC2Rlc2NyaXB0aW9uiAEBEjoKCWFyZWFfdHlw'
    'ZRgEIAEoDjIYLmdlb2xvY2F0aW9uLnYxLkFyZWFUeXBlSAJSCGFyZWFUeXBliAEBEh8KCGdlb2'
    '1ldHJ5GAUgASgJSANSCGdlb21ldHJ5iAEBEi0KBWV4dHJhGAYgASgLMhcuZ29vZ2xlLnByb3Rv'
    'YnVmLlN0cnVjdFIFZXh0cmFCBwoFX25hbWVCDgoMX2Rlc2NyaXB0aW9uQgwKCl9hcmVhX3R5cG'
    'VCCwoJX2dlb21ldHJ5');

@$core.Deprecated('Use updateAreaResponseDescriptor instead')
const UpdateAreaResponse$json = {
  '1': 'UpdateAreaResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.AreaObject', '10': 'data'},
  ],
};

/// Descriptor for `UpdateAreaResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAreaResponseDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVBcmVhUmVzcG9uc2USLgoEZGF0YRgBIAEoCzIaLmdlb2xvY2F0aW9uLnYxLkFyZW'
    'FPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use deleteAreaRequestDescriptor instead')
const DeleteAreaRequest$json = {
  '1': 'DeleteAreaRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteAreaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAreaRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVBcmVhUmVxdWVzdBIOCgJpZBgBIAEoCVICaWQ=');

@$core.Deprecated('Use searchAreasRequestDescriptor instead')
const SearchAreasRequest$json = {
  '1': 'SearchAreasRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'owner_id', '3': 2, '4': 1, '5': 9, '10': 'ownerId'},
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `SearchAreasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchAreasRequestDescriptor = $convert.base64Decode(
    'ChJTZWFyY2hBcmVhc1JlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5EhkKCG93bmVyX2lkGA'
    'IgASgJUgdvd25lcklkEhQKBWxpbWl0GAMgASgFUgVsaW1pdA==');

@$core.Deprecated('Use searchAreasResponseDescriptor instead')
const SearchAreasResponse$json = {
  '1': 'SearchAreasResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.AreaObject', '10': 'data'},
  ],
};

/// Descriptor for `SearchAreasResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchAreasResponseDescriptor = $convert.base64Decode(
    'ChNTZWFyY2hBcmVhc1Jlc3BvbnNlEi4KBGRhdGEYASADKAsyGi5nZW9sb2NhdGlvbi52MS5Bcm'
    'VhT2JqZWN0UgRkYXRh');

@$core.Deprecated('Use createRouteRequestDescriptor instead')
const CreateRouteRequest$json = {
  '1': 'CreateRouteRequest',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.RouteObject', '10': 'data'},
  ],
};

/// Descriptor for `CreateRouteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRouteRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVSb3V0ZVJlcXVlc3QSLwoEZGF0YRgBIAEoCzIbLmdlb2xvY2F0aW9uLnYxLlJvdX'
    'RlT2JqZWN0UgRkYXRh');

@$core.Deprecated('Use createRouteResponseDescriptor instead')
const CreateRouteResponse$json = {
  '1': 'CreateRouteResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.RouteObject', '10': 'data'},
  ],
};

/// Descriptor for `CreateRouteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRouteResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVSb3V0ZVJlc3BvbnNlEi8KBGRhdGEYASABKAsyGy5nZW9sb2NhdGlvbi52MS5Sb3'
    'V0ZU9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use getRouteRequestDescriptor instead')
const GetRouteRequest$json = {
  '1': 'GetRouteRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetRouteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRouteRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRSb3V0ZVJlcXVlc3QSDgoCaWQYASABKAlSAmlk');

@$core.Deprecated('Use getRouteResponseDescriptor instead')
const GetRouteResponse$json = {
  '1': 'GetRouteResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.RouteObject', '10': 'data'},
  ],
};

/// Descriptor for `GetRouteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRouteResponseDescriptor = $convert.base64Decode(
    'ChBHZXRSb3V0ZVJlc3BvbnNlEi8KBGRhdGEYASABKAsyGy5nZW9sb2NhdGlvbi52MS5Sb3V0ZU'
    '9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use updateRouteRequestDescriptor instead')
const UpdateRouteRequest$json = {
  '1': 'UpdateRouteRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '9': 0, '10': 'name', '17': true},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '9': 1, '10': 'description', '17': true},
    {'1': 'geometry', '3': 4, '4': 1, '5': 9, '9': 2, '10': 'geometry', '17': true},
    {'1': 'deviation_threshold_m', '3': 5, '4': 1, '5': 1, '9': 3, '10': 'deviationThresholdM', '17': true},
    {'1': 'deviation_consecutive_count', '3': 6, '4': 1, '5': 5, '9': 4, '10': 'deviationConsecutiveCount', '17': true},
    {'1': 'deviation_cooldown_sec', '3': 7, '4': 1, '5': 5, '9': 5, '10': 'deviationCooldownSec', '17': true},
    {'1': 'extra', '3': 8, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
  ],
  '8': [
    {'1': '_name'},
    {'1': '_description'},
    {'1': '_geometry'},
    {'1': '_deviation_threshold_m'},
    {'1': '_deviation_consecutive_count'},
    {'1': '_deviation_cooldown_sec'},
  ],
};

/// Descriptor for `UpdateRouteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRouteRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVSb3V0ZVJlcXVlc3QSDgoCaWQYASABKAlSAmlkEhcKBG5hbWUYAiABKAlIAFIEbm'
    'FtZYgBARIlCgtkZXNjcmlwdGlvbhgDIAEoCUgBUgtkZXNjcmlwdGlvbogBARIfCghnZW9tZXRy'
    'eRgEIAEoCUgCUghnZW9tZXRyeYgBARI3ChVkZXZpYXRpb25fdGhyZXNob2xkX20YBSABKAFIA1'
    'ITZGV2aWF0aW9uVGhyZXNob2xkTYgBARJDChtkZXZpYXRpb25fY29uc2VjdXRpdmVfY291bnQY'
    'BiABKAVIBFIZZGV2aWF0aW9uQ29uc2VjdXRpdmVDb3VudIgBARI5ChZkZXZpYXRpb25fY29vbG'
    'Rvd25fc2VjGAcgASgFSAVSFGRldmlhdGlvbkNvb2xkb3duU2VjiAEBEi0KBWV4dHJhGAggASgL'
    'MhcuZ29vZ2xlLnByb3RvYnVmLlN0cnVjdFIFZXh0cmFCBwoFX25hbWVCDgoMX2Rlc2NyaXB0aW'
    '9uQgsKCV9nZW9tZXRyeUIYChZfZGV2aWF0aW9uX3RocmVzaG9sZF9tQh4KHF9kZXZpYXRpb25f'
    'Y29uc2VjdXRpdmVfY291bnRCGQoXX2RldmlhdGlvbl9jb29sZG93bl9zZWM=');

@$core.Deprecated('Use updateRouteResponseDescriptor instead')
const UpdateRouteResponse$json = {
  '1': 'UpdateRouteResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.RouteObject', '10': 'data'},
  ],
};

/// Descriptor for `UpdateRouteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRouteResponseDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVSb3V0ZVJlc3BvbnNlEi8KBGRhdGEYASABKAsyGy5nZW9sb2NhdGlvbi52MS5Sb3'
    'V0ZU9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use deleteRouteRequestDescriptor instead')
const DeleteRouteRequest$json = {
  '1': 'DeleteRouteRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteRouteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRouteRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVSb3V0ZVJlcXVlc3QSDgoCaWQYASABKAlSAmlk');

@$core.Deprecated('Use searchRoutesRequestDescriptor instead')
const SearchRoutesRequest$json = {
  '1': 'SearchRoutesRequest',
  '2': [
    {'1': 'owner_id', '3': 1, '4': 1, '5': 9, '10': 'ownerId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `SearchRoutesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchRoutesRequestDescriptor = $convert.base64Decode(
    'ChNTZWFyY2hSb3V0ZXNSZXF1ZXN0EhkKCG93bmVyX2lkGAEgASgJUgdvd25lcklkEhQKBWxpbW'
    'l0GAIgASgFUgVsaW1pdA==');

@$core.Deprecated('Use searchRoutesResponseDescriptor instead')
const SearchRoutesResponse$json = {
  '1': 'SearchRoutesResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.RouteObject', '10': 'data'},
  ],
};

/// Descriptor for `SearchRoutesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchRoutesResponseDescriptor = $convert.base64Decode(
    'ChRTZWFyY2hSb3V0ZXNSZXNwb25zZRIvCgRkYXRhGAEgAygLMhsuZ2VvbG9jYXRpb24udjEuUm'
    '91dGVPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use assignRouteRequestDescriptor instead')
const AssignRouteRequest$json = {
  '1': 'AssignRouteRequest',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'route_id', '3': 2, '4': 1, '5': 9, '10': 'routeId'},
    {'1': 'valid_from', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'validFrom'},
    {'1': 'valid_until', '3': 4, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'validUntil'},
  ],
};

/// Descriptor for `AssignRouteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assignRouteRequestDescriptor = $convert.base64Decode(
    'ChJBc3NpZ25Sb3V0ZVJlcXVlc3QSHQoKc3ViamVjdF9pZBgBIAEoCVIJc3ViamVjdElkEhkKCH'
    'JvdXRlX2lkGAIgASgJUgdyb3V0ZUlkEjkKCnZhbGlkX2Zyb20YAyABKAsyGi5nb29nbGUucHJv'
    'dG9idWYuVGltZXN0YW1wUgl2YWxpZEZyb20SOwoLdmFsaWRfdW50aWwYBCABKAsyGi5nb29nbG'
    'UucHJvdG9idWYuVGltZXN0YW1wUgp2YWxpZFVudGls');

@$core.Deprecated('Use assignRouteResponseDescriptor instead')
const AssignRouteResponse$json = {
  '1': 'AssignRouteResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.geolocation.v1.RouteAssignmentObject', '10': 'data'},
  ],
};

/// Descriptor for `AssignRouteResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assignRouteResponseDescriptor = $convert.base64Decode(
    'ChNBc3NpZ25Sb3V0ZVJlc3BvbnNlEjkKBGRhdGEYASABKAsyJS5nZW9sb2NhdGlvbi52MS5Sb3'
    'V0ZUFzc2lnbm1lbnRPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use unassignRouteRequestDescriptor instead')
const UnassignRouteRequest$json = {
  '1': 'UnassignRouteRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `UnassignRouteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unassignRouteRequestDescriptor = $convert.base64Decode(
    'ChRVbmFzc2lnblJvdXRlUmVxdWVzdBIOCgJpZBgBIAEoCVICaWQ=');

@$core.Deprecated('Use getSubjectRouteAssignmentsRequestDescriptor instead')
const GetSubjectRouteAssignmentsRequest$json = {
  '1': 'GetSubjectRouteAssignmentsRequest',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
  ],
};

/// Descriptor for `GetSubjectRouteAssignmentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubjectRouteAssignmentsRequestDescriptor = $convert.base64Decode(
    'CiFHZXRTdWJqZWN0Um91dGVBc3NpZ25tZW50c1JlcXVlc3QSHQoKc3ViamVjdF9pZBgBIAEoCV'
    'IJc3ViamVjdElk');

@$core.Deprecated('Use getSubjectRouteAssignmentsResponseDescriptor instead')
const GetSubjectRouteAssignmentsResponse$json = {
  '1': 'GetSubjectRouteAssignmentsResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.RouteAssignmentObject', '10': 'data'},
  ],
};

/// Descriptor for `GetSubjectRouteAssignmentsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubjectRouteAssignmentsResponseDescriptor = $convert.base64Decode(
    'CiJHZXRTdWJqZWN0Um91dGVBc3NpZ25tZW50c1Jlc3BvbnNlEjkKBGRhdGEYASADKAsyJS5nZW'
    '9sb2NhdGlvbi52MS5Sb3V0ZUFzc2lnbm1lbnRPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use getTrackRequestDescriptor instead')
const GetTrackRequest$json = {
  '1': 'GetTrackRequest',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'from', '3': 2, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'from'},
    {'1': 'to', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'to'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetTrackRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrackRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRUcmFja1JlcXVlc3QSHQoKc3ViamVjdF9pZBgBIAEoCVIJc3ViamVjdElkEi4KBGZyb2'
    '0YAiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgRmcm9tEioKAnRvGAMgASgLMhou'
    'Z29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFICdG8SFAoFbGltaXQYBCABKAVSBWxpbWl0EhYKBm'
    '9mZnNldBgFIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use getTrackResponseDescriptor instead')
const GetTrackResponse$json = {
  '1': 'GetTrackResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.LocationPointObject', '10': 'data'},
  ],
};

/// Descriptor for `GetTrackResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrackResponseDescriptor = $convert.base64Decode(
    'ChBHZXRUcmFja1Jlc3BvbnNlEjcKBGRhdGEYASADKAsyIy5nZW9sb2NhdGlvbi52MS5Mb2NhdG'
    'lvblBvaW50T2JqZWN0UgRkYXRh');

@$core.Deprecated('Use getSubjectEventsRequestDescriptor instead')
const GetSubjectEventsRequest$json = {
  '1': 'GetSubjectEventsRequest',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'from', '3': 2, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'from'},
    {'1': 'to', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'to'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetSubjectEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubjectEventsRequestDescriptor = $convert.base64Decode(
    'ChdHZXRTdWJqZWN0RXZlbnRzUmVxdWVzdBIdCgpzdWJqZWN0X2lkGAEgASgJUglzdWJqZWN0SW'
    'QSLgoEZnJvbRgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSBGZyb20SKgoCdG8Y'
    'AyABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgJ0bxIUCgVsaW1pdBgEIAEoBVIFbG'
    'ltaXQSFgoGb2Zmc2V0GAUgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use getSubjectEventsResponseDescriptor instead')
const GetSubjectEventsResponse$json = {
  '1': 'GetSubjectEventsResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.GeoEventObject', '10': 'data'},
  ],
};

/// Descriptor for `GetSubjectEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubjectEventsResponseDescriptor = $convert.base64Decode(
    'ChhHZXRTdWJqZWN0RXZlbnRzUmVzcG9uc2USMgoEZGF0YRgBIAMoCzIeLmdlb2xvY2F0aW9uLn'
    'YxLkdlb0V2ZW50T2JqZWN0UgRkYXRh');

@$core.Deprecated('Use getAreaSubjectsRequestDescriptor instead')
const GetAreaSubjectsRequest$json = {
  '1': 'GetAreaSubjectsRequest',
  '2': [
    {'1': 'area_id', '3': 1, '4': 1, '5': 9, '10': 'areaId'},
  ],
};

/// Descriptor for `GetAreaSubjectsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAreaSubjectsRequestDescriptor = $convert.base64Decode(
    'ChZHZXRBcmVhU3ViamVjdHNSZXF1ZXN0EhcKB2FyZWFfaWQYASABKAlSBmFyZWFJZA==');

@$core.Deprecated('Use getAreaSubjectsResponseDescriptor instead')
const GetAreaSubjectsResponse$json = {
  '1': 'GetAreaSubjectsResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.AreaSubjectObject', '10': 'data'},
  ],
};

/// Descriptor for `GetAreaSubjectsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAreaSubjectsResponseDescriptor = $convert.base64Decode(
    'ChdHZXRBcmVhU3ViamVjdHNSZXNwb25zZRI1CgRkYXRhGAEgAygLMiEuZ2VvbG9jYXRpb24udj'
    'EuQXJlYVN1YmplY3RPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use getNearbySubjectsRequestDescriptor instead')
const GetNearbySubjectsRequest$json = {
  '1': 'GetNearbySubjectsRequest',
  '2': [
    {'1': 'subject_id', '3': 1, '4': 1, '5': 9, '10': 'subjectId'},
    {'1': 'radius_meters', '3': 2, '4': 1, '5': 1, '10': 'radiusMeters'},
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `GetNearbySubjectsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNearbySubjectsRequestDescriptor = $convert.base64Decode(
    'ChhHZXROZWFyYnlTdWJqZWN0c1JlcXVlc3QSHQoKc3ViamVjdF9pZBgBIAEoCVIJc3ViamVjdE'
    'lkEiMKDXJhZGl1c19tZXRlcnMYAiABKAFSDHJhZGl1c01ldGVycxIUCgVsaW1pdBgDIAEoBVIF'
    'bGltaXQ=');

@$core.Deprecated('Use getNearbySubjectsResponseDescriptor instead')
const GetNearbySubjectsResponse$json = {
  '1': 'GetNearbySubjectsResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.NearbySubjectObject', '10': 'data'},
  ],
};

/// Descriptor for `GetNearbySubjectsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNearbySubjectsResponseDescriptor = $convert.base64Decode(
    'ChlHZXROZWFyYnlTdWJqZWN0c1Jlc3BvbnNlEjcKBGRhdGEYASADKAsyIy5nZW9sb2NhdGlvbi'
    '52MS5OZWFyYnlTdWJqZWN0T2JqZWN0UgRkYXRh');

@$core.Deprecated('Use getNearbyAreasRequestDescriptor instead')
const GetNearbyAreasRequest$json = {
  '1': 'GetNearbyAreasRequest',
  '2': [
    {'1': 'latitude', '3': 1, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 2, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'radius_meters', '3': 3, '4': 1, '5': 1, '10': 'radiusMeters'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `GetNearbyAreasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNearbyAreasRequestDescriptor = $convert.base64Decode(
    'ChVHZXROZWFyYnlBcmVhc1JlcXVlc3QSGgoIbGF0aXR1ZGUYASABKAFSCGxhdGl0dWRlEhwKCW'
    'xvbmdpdHVkZRgCIAEoAVIJbG9uZ2l0dWRlEiMKDXJhZGl1c19tZXRlcnMYAyABKAFSDHJhZGl1'
    'c01ldGVycxIUCgVsaW1pdBgEIAEoBVIFbGltaXQ=');

@$core.Deprecated('Use getNearbyAreasResponseDescriptor instead')
const GetNearbyAreasResponse$json = {
  '1': 'GetNearbyAreasResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.geolocation.v1.NearbyAreaObject', '10': 'data'},
  ],
};

/// Descriptor for `GetNearbyAreasResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNearbyAreasResponseDescriptor = $convert.base64Decode(
    'ChZHZXROZWFyYnlBcmVhc1Jlc3BvbnNlEjQKBGRhdGEYASADKAsyIC5nZW9sb2NhdGlvbi52MS'
    '5OZWFyYnlBcmVhT2JqZWN0UgRkYXRh');

const $core.Map<$core.String, $core.dynamic> GeolocationServiceBase$json = {
  '1': 'GeolocationService',
  '2': [
    {'1': 'IngestLocations', '2': '.geolocation.v1.IngestLocationsRequest', '3': '.geolocation.v1.IngestLocationsResponse', '4': {}},
    {'1': 'CreateArea', '2': '.geolocation.v1.CreateAreaRequest', '3': '.geolocation.v1.CreateAreaResponse', '4': {}},
    {'1': 'GetArea', '2': '.geolocation.v1.GetAreaRequest', '3': '.geolocation.v1.GetAreaResponse', '4': {}},
    {'1': 'UpdateArea', '2': '.geolocation.v1.UpdateAreaRequest', '3': '.geolocation.v1.UpdateAreaResponse', '4': {}},
    {'1': 'DeleteArea', '2': '.geolocation.v1.DeleteAreaRequest', '3': '.google.protobuf.Empty', '4': {}},
    {'1': 'SearchAreas', '2': '.geolocation.v1.SearchAreasRequest', '3': '.geolocation.v1.SearchAreasResponse', '4': {}},
    {'1': 'CreateRoute', '2': '.geolocation.v1.CreateRouteRequest', '3': '.geolocation.v1.CreateRouteResponse', '4': {}},
    {'1': 'GetRoute', '2': '.geolocation.v1.GetRouteRequest', '3': '.geolocation.v1.GetRouteResponse', '4': {}},
    {'1': 'UpdateRoute', '2': '.geolocation.v1.UpdateRouteRequest', '3': '.geolocation.v1.UpdateRouteResponse', '4': {}},
    {'1': 'DeleteRoute', '2': '.geolocation.v1.DeleteRouteRequest', '3': '.google.protobuf.Empty', '4': {}},
    {'1': 'SearchRoutes', '2': '.geolocation.v1.SearchRoutesRequest', '3': '.geolocation.v1.SearchRoutesResponse', '4': {}},
    {'1': 'AssignRoute', '2': '.geolocation.v1.AssignRouteRequest', '3': '.geolocation.v1.AssignRouteResponse', '4': {}},
    {'1': 'UnassignRoute', '2': '.geolocation.v1.UnassignRouteRequest', '3': '.google.protobuf.Empty', '4': {}},
    {'1': 'GetSubjectRouteAssignments', '2': '.geolocation.v1.GetSubjectRouteAssignmentsRequest', '3': '.geolocation.v1.GetSubjectRouteAssignmentsResponse', '4': {}},
    {'1': 'GetTrack', '2': '.geolocation.v1.GetTrackRequest', '3': '.geolocation.v1.GetTrackResponse', '4': {}},
    {'1': 'GetSubjectEvents', '2': '.geolocation.v1.GetSubjectEventsRequest', '3': '.geolocation.v1.GetSubjectEventsResponse', '4': {}},
    {'1': 'GetAreaSubjects', '2': '.geolocation.v1.GetAreaSubjectsRequest', '3': '.geolocation.v1.GetAreaSubjectsResponse', '4': {}},
    {'1': 'GetNearbySubjects', '2': '.geolocation.v1.GetNearbySubjectsRequest', '3': '.geolocation.v1.GetNearbySubjectsResponse', '4': {}},
    {'1': 'GetNearbyAreas', '2': '.geolocation.v1.GetNearbyAreasRequest', '3': '.geolocation.v1.GetNearbyAreasResponse', '4': {}},
  ],
  '3': {},
};

@$core.Deprecated('Use geolocationServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> GeolocationServiceBase$messageJson = {
  '.geolocation.v1.IngestLocationsRequest': IngestLocationsRequest$json,
  '.geolocation.v1.LocationPointInput': LocationPointInput$json,
  '.google.protobuf.Timestamp': $0.Timestamp$json,
  '.google.protobuf.Struct': $1.Struct$json,
  '.google.protobuf.Struct.FieldsEntry': $1.Struct_FieldsEntry$json,
  '.google.protobuf.Value': $1.Value$json,
  '.google.protobuf.ListValue': $1.ListValue$json,
  '.geolocation.v1.IngestLocationsResponse': IngestLocationsResponse$json,
  '.geolocation.v1.CreateAreaRequest': CreateAreaRequest$json,
  '.geolocation.v1.AreaObject': AreaObject$json,
  '.geolocation.v1.CreateAreaResponse': CreateAreaResponse$json,
  '.geolocation.v1.GetAreaRequest': GetAreaRequest$json,
  '.geolocation.v1.GetAreaResponse': GetAreaResponse$json,
  '.geolocation.v1.UpdateAreaRequest': UpdateAreaRequest$json,
  '.geolocation.v1.UpdateAreaResponse': UpdateAreaResponse$json,
  '.geolocation.v1.DeleteAreaRequest': DeleteAreaRequest$json,
  '.google.protobuf.Empty': $2.Empty$json,
  '.geolocation.v1.SearchAreasRequest': SearchAreasRequest$json,
  '.geolocation.v1.SearchAreasResponse': SearchAreasResponse$json,
  '.geolocation.v1.CreateRouteRequest': CreateRouteRequest$json,
  '.geolocation.v1.RouteObject': RouteObject$json,
  '.geolocation.v1.CreateRouteResponse': CreateRouteResponse$json,
  '.geolocation.v1.GetRouteRequest': GetRouteRequest$json,
  '.geolocation.v1.GetRouteResponse': GetRouteResponse$json,
  '.geolocation.v1.UpdateRouteRequest': UpdateRouteRequest$json,
  '.geolocation.v1.UpdateRouteResponse': UpdateRouteResponse$json,
  '.geolocation.v1.DeleteRouteRequest': DeleteRouteRequest$json,
  '.geolocation.v1.SearchRoutesRequest': SearchRoutesRequest$json,
  '.geolocation.v1.SearchRoutesResponse': SearchRoutesResponse$json,
  '.geolocation.v1.AssignRouteRequest': AssignRouteRequest$json,
  '.geolocation.v1.AssignRouteResponse': AssignRouteResponse$json,
  '.geolocation.v1.RouteAssignmentObject': RouteAssignmentObject$json,
  '.geolocation.v1.UnassignRouteRequest': UnassignRouteRequest$json,
  '.geolocation.v1.GetSubjectRouteAssignmentsRequest': GetSubjectRouteAssignmentsRequest$json,
  '.geolocation.v1.GetSubjectRouteAssignmentsResponse': GetSubjectRouteAssignmentsResponse$json,
  '.geolocation.v1.GetTrackRequest': GetTrackRequest$json,
  '.geolocation.v1.GetTrackResponse': GetTrackResponse$json,
  '.geolocation.v1.LocationPointObject': LocationPointObject$json,
  '.geolocation.v1.GetSubjectEventsRequest': GetSubjectEventsRequest$json,
  '.geolocation.v1.GetSubjectEventsResponse': GetSubjectEventsResponse$json,
  '.geolocation.v1.GeoEventObject': GeoEventObject$json,
  '.geolocation.v1.GetAreaSubjectsRequest': GetAreaSubjectsRequest$json,
  '.geolocation.v1.GetAreaSubjectsResponse': GetAreaSubjectsResponse$json,
  '.geolocation.v1.AreaSubjectObject': AreaSubjectObject$json,
  '.geolocation.v1.GetNearbySubjectsRequest': GetNearbySubjectsRequest$json,
  '.geolocation.v1.GetNearbySubjectsResponse': GetNearbySubjectsResponse$json,
  '.geolocation.v1.NearbySubjectObject': NearbySubjectObject$json,
  '.geolocation.v1.GetNearbyAreasRequest': GetNearbyAreasRequest$json,
  '.geolocation.v1.GetNearbyAreasResponse': GetNearbyAreasResponse$json,
  '.geolocation.v1.NearbyAreaObject': NearbyAreaObject$json,
};

/// Descriptor for `GeolocationService`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List geolocationServiceDescriptor = $convert.base64Decode(
    'ChJHZW9sb2NhdGlvblNlcnZpY2USeQoPSW5nZXN0TG9jYXRpb25zEiYuZ2VvbG9jYXRpb24udj'
    'EuSW5nZXN0TG9jYXRpb25zUmVxdWVzdBonLmdlb2xvY2F0aW9uLnYxLkluZ2VzdExvY2F0aW9u'
    'c1Jlc3BvbnNlIhWCtRgRCg9sb2NhdGlvbl9pbmdlc3QSZgoKQ3JlYXRlQXJlYRIhLmdlb2xvY2'
    'F0aW9uLnYxLkNyZWF0ZUFyZWFSZXF1ZXN0GiIuZ2VvbG9jYXRpb24udjEuQ3JlYXRlQXJlYVJl'
    'c3BvbnNlIhGCtRgNCgthcmVhX21hbmFnZRJbCgdHZXRBcmVhEh4uZ2VvbG9jYXRpb24udjEuR2'
    'V0QXJlYVJlcXVlc3QaHy5nZW9sb2NhdGlvbi52MS5HZXRBcmVhUmVzcG9uc2UiD4K1GAsKCWFy'
    'ZWFfdmlldxJmCgpVcGRhdGVBcmVhEiEuZ2VvbG9jYXRpb24udjEuVXBkYXRlQXJlYVJlcXVlc3'
    'QaIi5nZW9sb2NhdGlvbi52MS5VcGRhdGVBcmVhUmVzcG9uc2UiEYK1GA0KC2FyZWFfbWFuYWdl'
    'EloKCkRlbGV0ZUFyZWESIS5nZW9sb2NhdGlvbi52MS5EZWxldGVBcmVhUmVxdWVzdBoWLmdvb2'
    'dsZS5wcm90b2J1Zi5FbXB0eSIRgrUYDQoLYXJlYV9tYW5hZ2USZwoLU2VhcmNoQXJlYXMSIi5n'
    'ZW9sb2NhdGlvbi52MS5TZWFyY2hBcmVhc1JlcXVlc3QaIy5nZW9sb2NhdGlvbi52MS5TZWFyY2'
    'hBcmVhc1Jlc3BvbnNlIg+CtRgLCglhcmVhX3ZpZXcSagoLQ3JlYXRlUm91dGUSIi5nZW9sb2Nh'
    'dGlvbi52MS5DcmVhdGVSb3V0ZVJlcXVlc3QaIy5nZW9sb2NhdGlvbi52MS5DcmVhdGVSb3V0ZV'
    'Jlc3BvbnNlIhKCtRgOCgxyb3V0ZV9tYW5hZ2USXwoIR2V0Um91dGUSHy5nZW9sb2NhdGlvbi52'
    'MS5HZXRSb3V0ZVJlcXVlc3QaIC5nZW9sb2NhdGlvbi52MS5HZXRSb3V0ZVJlc3BvbnNlIhCCtR'
    'gMCgpyb3V0ZV92aWV3EmoKC1VwZGF0ZVJvdXRlEiIuZ2VvbG9jYXRpb24udjEuVXBkYXRlUm91'
    'dGVSZXF1ZXN0GiMuZ2VvbG9jYXRpb24udjEuVXBkYXRlUm91dGVSZXNwb25zZSISgrUYDgoMcm'
    '91dGVfbWFuYWdlEl0KC0RlbGV0ZVJvdXRlEiIuZ2VvbG9jYXRpb24udjEuRGVsZXRlUm91dGVS'
    'ZXF1ZXN0GhYuZ29vZ2xlLnByb3RvYnVmLkVtcHR5IhKCtRgOCgxyb3V0ZV9tYW5hZ2USawoMU2'
    'VhcmNoUm91dGVzEiMuZ2VvbG9jYXRpb24udjEuU2VhcmNoUm91dGVzUmVxdWVzdBokLmdlb2xv'
    'Y2F0aW9uLnYxLlNlYXJjaFJvdXRlc1Jlc3BvbnNlIhCCtRgMCgpyb3V0ZV92aWV3EmoKC0Fzc2'
    'lnblJvdXRlEiIuZ2VvbG9jYXRpb24udjEuQXNzaWduUm91dGVSZXF1ZXN0GiMuZ2VvbG9jYXRp'
    'b24udjEuQXNzaWduUm91dGVSZXNwb25zZSISgrUYDgoMcm91dGVfbWFuYWdlEmEKDVVuYXNzaW'
    'duUm91dGUSJC5nZW9sb2NhdGlvbi52MS5VbmFzc2lnblJvdXRlUmVxdWVzdBoWLmdvb2dsZS5w'
    'cm90b2J1Zi5FbXB0eSISgrUYDgoMcm91dGVfbWFuYWdlEpUBChpHZXRTdWJqZWN0Um91dGVBc3'
    'NpZ25tZW50cxIxLmdlb2xvY2F0aW9uLnYxLkdldFN1YmplY3RSb3V0ZUFzc2lnbm1lbnRzUmVx'
    'dWVzdBoyLmdlb2xvY2F0aW9uLnYxLkdldFN1YmplY3RSb3V0ZUFzc2lnbm1lbnRzUmVzcG9uc2'
    'UiEIK1GAwKCnJvdXRlX3ZpZXcSXwoIR2V0VHJhY2sSHy5nZW9sb2NhdGlvbi52MS5HZXRUcmFj'
    'a1JlcXVlc3QaIC5nZW9sb2NhdGlvbi52MS5HZXRUcmFja1Jlc3BvbnNlIhCCtRgMCgp0cmFja1'
    '92aWV3EncKEEdldFN1YmplY3RFdmVudHMSJy5nZW9sb2NhdGlvbi52MS5HZXRTdWJqZWN0RXZl'
    'bnRzUmVxdWVzdBooLmdlb2xvY2F0aW9uLnYxLkdldFN1YmplY3RFdmVudHNSZXNwb25zZSIQgr'
    'UYDAoKdHJhY2tfdmlldxJzCg9HZXRBcmVhU3ViamVjdHMSJi5nZW9sb2NhdGlvbi52MS5HZXRB'
    'cmVhU3ViamVjdHNSZXF1ZXN0GicuZ2VvbG9jYXRpb24udjEuR2V0QXJlYVN1YmplY3RzUmVzcG'
    '9uc2UiD4K1GAsKCWFyZWFfdmlldxJ7ChFHZXROZWFyYnlTdWJqZWN0cxIoLmdlb2xvY2F0aW9u'
    'LnYxLkdldE5lYXJieVN1YmplY3RzUmVxdWVzdBopLmdlb2xvY2F0aW9uLnYxLkdldE5lYXJieV'
    'N1YmplY3RzUmVzcG9uc2UiEYK1GA0KC25lYXJieV92aWV3EnIKDkdldE5lYXJieUFyZWFzEiUu'
    'Z2VvbG9jYXRpb24udjEuR2V0TmVhcmJ5QXJlYXNSZXF1ZXN0GiYuZ2VvbG9jYXRpb24udjEuR2'
    'V0TmVhcmJ5QXJlYXNSZXNwb25zZSIRgrUYDQoLbmVhcmJ5X3ZpZXcavwSCtRi6BAoPc2Vydmlj'
    'ZV9wcm9maWxlEg9sb2NhdGlvbl9pbmdlc3QSCWFyZWFfdmlldxILYXJlYV9tYW5hZ2USCnJvdX'
    'RlX3ZpZXcSDHJvdXRlX21hbmFnZRIKdHJhY2tfdmlldxILbmVhcmJ5X3ZpZXcaXggBEg9sb2Nh'
    'dGlvbl9pbmdlc3QSCWFyZWFfdmlldxILYXJlYV9tYW5hZ2USCnJvdXRlX3ZpZXcSDHJvdXRlX2'
    '1hbmFnZRIKdHJhY2tfdmlldxILbmVhcmJ5X3ZpZXcaXggCEg9sb2NhdGlvbl9pbmdlc3QSCWFy'
    'ZWFfdmlldxILYXJlYV9tYW5hZ2USCnJvdXRlX3ZpZXcSDHJvdXRlX21hbmFnZRIKdHJhY2tfdm'
    'lldxILbmVhcmJ5X3ZpZXcaQwgDEg9sb2NhdGlvbl9pbmdlc3QSCWFyZWFfdmlldxIKcm91dGVf'
    'dmlldxIKdHJhY2tfdmlldxILbmVhcmJ5X3ZpZXcaMggEEglhcmVhX3ZpZXcSCnJvdXRlX3ZpZX'
    'cSCnRyYWNrX3ZpZXcSC25lYXJieV92aWV3GjIIBRIJYXJlYV92aWV3Egpyb3V0ZV92aWV3Egp0'
    'cmFja192aWV3EgtuZWFyYnlfdmlldxpeCAYSD2xvY2F0aW9uX2luZ2VzdBIJYXJlYV92aWV3Eg'
    'thcmVhX21hbmFnZRIKcm91dGVfdmlldxIMcm91dGVfbWFuYWdlEgp0cmFja192aWV3EgtuZWFy'
    'Ynlfdmlldw==');

