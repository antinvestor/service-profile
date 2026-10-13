// This is a generated file - do not edit.
//
// Generated from geolocation/v1/geolocation.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class LocationSource extends $pb.ProtobufEnum {
  static const LocationSource LOCATION_SOURCE_UNSPECIFIED =
      LocationSource._(0, _omitEnumNames ? '' : 'LOCATION_SOURCE_UNSPECIFIED');
  static const LocationSource LOCATION_SOURCE_GPS =
      LocationSource._(1, _omitEnumNames ? '' : 'LOCATION_SOURCE_GPS');
  static const LocationSource LOCATION_SOURCE_NETWORK =
      LocationSource._(2, _omitEnumNames ? '' : 'LOCATION_SOURCE_NETWORK');
  static const LocationSource LOCATION_SOURCE_IP =
      LocationSource._(3, _omitEnumNames ? '' : 'LOCATION_SOURCE_IP');
  static const LocationSource LOCATION_SOURCE_MANUAL =
      LocationSource._(4, _omitEnumNames ? '' : 'LOCATION_SOURCE_MANUAL');

  static const $core.List<LocationSource> values = <LocationSource>[
    LOCATION_SOURCE_UNSPECIFIED,
    LOCATION_SOURCE_GPS,
    LOCATION_SOURCE_NETWORK,
    LOCATION_SOURCE_IP,
    LOCATION_SOURCE_MANUAL,
  ];

  static final $core.List<LocationSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static LocationSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LocationSource._(super.value, super.name);
}

class AreaType extends $pb.ProtobufEnum {
  static const AreaType AREA_TYPE_UNSPECIFIED =
      AreaType._(0, _omitEnumNames ? '' : 'AREA_TYPE_UNSPECIFIED');
  static const AreaType AREA_TYPE_LAND =
      AreaType._(1, _omitEnumNames ? '' : 'AREA_TYPE_LAND');
  static const AreaType AREA_TYPE_BUILDING =
      AreaType._(2, _omitEnumNames ? '' : 'AREA_TYPE_BUILDING');
  static const AreaType AREA_TYPE_ZONE =
      AreaType._(3, _omitEnumNames ? '' : 'AREA_TYPE_ZONE');
  static const AreaType AREA_TYPE_FENCE =
      AreaType._(4, _omitEnumNames ? '' : 'AREA_TYPE_FENCE');
  static const AreaType AREA_TYPE_CUSTOM =
      AreaType._(5, _omitEnumNames ? '' : 'AREA_TYPE_CUSTOM');

  static const $core.List<AreaType> values = <AreaType>[
    AREA_TYPE_UNSPECIFIED,
    AREA_TYPE_LAND,
    AREA_TYPE_BUILDING,
    AREA_TYPE_ZONE,
    AREA_TYPE_FENCE,
    AREA_TYPE_CUSTOM,
  ];

  static final $core.List<AreaType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static AreaType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AreaType._(super.value, super.name);
}

class GeoEventType extends $pb.ProtobufEnum {
  static const GeoEventType GEO_EVENT_TYPE_UNSPECIFIED =
      GeoEventType._(0, _omitEnumNames ? '' : 'GEO_EVENT_TYPE_UNSPECIFIED');
  static const GeoEventType GEO_EVENT_TYPE_ENTER =
      GeoEventType._(1, _omitEnumNames ? '' : 'GEO_EVENT_TYPE_ENTER');
  static const GeoEventType GEO_EVENT_TYPE_EXIT =
      GeoEventType._(2, _omitEnumNames ? '' : 'GEO_EVENT_TYPE_EXIT');
  static const GeoEventType GEO_EVENT_TYPE_DWELL =
      GeoEventType._(3, _omitEnumNames ? '' : 'GEO_EVENT_TYPE_DWELL');

  static const $core.List<GeoEventType> values = <GeoEventType>[
    GEO_EVENT_TYPE_UNSPECIFIED,
    GEO_EVENT_TYPE_ENTER,
    GEO_EVENT_TYPE_EXIT,
    GEO_EVENT_TYPE_DWELL,
  ];

  static final $core.List<GeoEventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static GeoEventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GeoEventType._(super.value, super.name);
}

class RouteDeviationEventType extends $pb.ProtobufEnum {
  static const RouteDeviationEventType ROUTE_DEVIATION_EVENT_TYPE_UNSPECIFIED =
      RouteDeviationEventType._(
          0, _omitEnumNames ? '' : 'ROUTE_DEVIATION_EVENT_TYPE_UNSPECIFIED');
  static const RouteDeviationEventType ROUTE_DEVIATION_EVENT_TYPE_DEVIATED =
      RouteDeviationEventType._(
          1, _omitEnumNames ? '' : 'ROUTE_DEVIATION_EVENT_TYPE_DEVIATED');
  static const RouteDeviationEventType
      ROUTE_DEVIATION_EVENT_TYPE_BACK_ON_ROUTE = RouteDeviationEventType._(
          2, _omitEnumNames ? '' : 'ROUTE_DEVIATION_EVENT_TYPE_BACK_ON_ROUTE');

  static const $core.List<RouteDeviationEventType> values =
      <RouteDeviationEventType>[
    ROUTE_DEVIATION_EVENT_TYPE_UNSPECIFIED,
    ROUTE_DEVIATION_EVENT_TYPE_DEVIATED,
    ROUTE_DEVIATION_EVENT_TYPE_BACK_ON_ROUTE,
  ];

  static final $core.List<RouteDeviationEventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static RouteDeviationEventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RouteDeviationEventType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
