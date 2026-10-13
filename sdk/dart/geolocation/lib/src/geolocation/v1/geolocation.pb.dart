// This is a generated file - do not edit.
//
// Generated from geolocation/v1/geolocation.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;
import 'package:protobuf/well_known_types/google/protobuf/empty.pb.dart' as $2;
import 'package:protobuf/well_known_types/google/protobuf/struct.pb.dart' as $1;
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $0;

import 'geolocation.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'geolocation.pbenum.dart';

class LocationPointInput extends $pb.GeneratedMessage {
  factory LocationPointInput({
    $0.Timestamp? timestamp,
    $core.double? latitude,
    $core.double? longitude,
    $core.double? altitude,
    $core.double? accuracy,
    $core.double? speed,
    $core.double? bearing,
    LocationSource? source,
    $1.Struct? extra,
    $core.String? deviceId,
  }) {
    final result = LocationPointInput._();
    if (timestamp != null) result.timestamp = timestamp;
    if (latitude != null) result.latitude = latitude;
    if (longitude != null) result.longitude = longitude;
    if (altitude != null) result.altitude = altitude;
    if (accuracy != null) result.accuracy = accuracy;
    if (speed != null) result.speed = speed;
    if (bearing != null) result.bearing = bearing;
    if (source != null) result.source = source;
    if (extra != null) result.extra = extra;
    if (deviceId != null) result.deviceId = deviceId;
    return result;
  }

  LocationPointInput._();

  factory LocationPointInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocationPointInput()..mergeFromBuffer(data, registry);
  factory LocationPointInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocationPointInput()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocationPointInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: LocationPointInput.$_createMessage)
    ..aOM<$0.Timestamp>(1, _omitFieldNames ? '' : 'timestamp',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aD(2, _omitFieldNames ? '' : 'latitude')
    ..aD(3, _omitFieldNames ? '' : 'longitude')
    ..aD(4, _omitFieldNames ? '' : 'altitude')
    ..aD(5, _omitFieldNames ? '' : 'accuracy')
    ..aD(6, _omitFieldNames ? '' : 'speed')
    ..aD(7, _omitFieldNames ? '' : 'bearing')
    ..aE<LocationSource>(8, _omitFieldNames ? '' : 'source',
        enumValues: LocationSource.values)
    ..aOM<$1.Struct>(9, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..aOS(10, _omitFieldNames ? '' : 'deviceId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocationPointInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocationPointInput copyWith(void Function(LocationPointInput) updates) =>
      super.copyWith((message) => updates(message as LocationPointInput))
          as LocationPointInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use LocationPointInput() / LocationPointInput.new instead')
  static LocationPointInput create() => LocationPointInput._();
  static $pb.GeneratedMessage $_createMessage() => LocationPointInput._();
  @$core.override
  LocationPointInput createEmptyInstance() => LocationPointInput._();
  @$core.pragma('dart2js:noInline')
  static LocationPointInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocationPointInput>(
          LocationPointInput.$_createMessage);
  static LocationPointInput? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Timestamp get timestamp => $_getN(0);
  @$pb.TagNumber(1)
  set timestamp($0.Timestamp value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasTimestamp() => $_has(0);
  @$pb.TagNumber(1)
  void clearTimestamp() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Timestamp ensureTimestamp() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.double get latitude => $_getN(1);
  @$pb.TagNumber(2)
  set latitude($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLatitude() => $_has(1);
  @$pb.TagNumber(2)
  void clearLatitude() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.double get longitude => $_getN(2);
  @$pb.TagNumber(3)
  set longitude($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLongitude() => $_has(2);
  @$pb.TagNumber(3)
  void clearLongitude() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.double get altitude => $_getN(3);
  @$pb.TagNumber(4)
  set altitude($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAltitude() => $_has(3);
  @$pb.TagNumber(4)
  void clearAltitude() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.double get accuracy => $_getN(4);
  @$pb.TagNumber(5)
  set accuracy($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(5)
  $core.bool hasAccuracy() => $_has(4);
  @$pb.TagNumber(5)
  void clearAccuracy() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.double get speed => $_getN(5);
  @$pb.TagNumber(6)
  set speed($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(6)
  $core.bool hasSpeed() => $_has(5);
  @$pb.TagNumber(6)
  void clearSpeed() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.double get bearing => $_getN(6);
  @$pb.TagNumber(7)
  set bearing($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(7)
  $core.bool hasBearing() => $_has(6);
  @$pb.TagNumber(7)
  void clearBearing() => $_clearField(7);

  @$pb.TagNumber(8)
  LocationSource get source => $_getN(7);
  @$pb.TagNumber(8)
  set source(LocationSource value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasSource() => $_has(7);
  @$pb.TagNumber(8)
  void clearSource() => $_clearField(8);

  @$pb.TagNumber(9)
  $1.Struct get extra => $_getN(8);
  @$pb.TagNumber(9)
  set extra($1.Struct value) => $_setField(9, value);
  @$pb.TagNumber(9)
  $core.bool hasExtra() => $_has(8);
  @$pb.TagNumber(9)
  void clearExtra() => $_clearField(9);
  @$pb.TagNumber(9)
  $1.Struct ensureExtra() => $_ensure(8);

  @$pb.TagNumber(10)
  $core.String get deviceId => $_getSZ(9);
  @$pb.TagNumber(10)
  set deviceId($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDeviceId() => $_has(9);
  @$pb.TagNumber(10)
  void clearDeviceId() => $_clearField(10);
}

class LocationPointObject extends $pb.GeneratedMessage {
  factory LocationPointObject({
    $core.String? id,
    $core.String? subjectId,
    $0.Timestamp? timestamp,
    $core.double? latitude,
    $core.double? longitude,
    $core.double? altitude,
    $core.double? accuracy,
    $core.double? speed,
    $core.double? bearing,
    LocationSource? source,
    $1.Struct? extra,
    $0.Timestamp? createdAt,
    $core.String? deviceId,
  }) {
    final result = LocationPointObject._();
    if (id != null) result.id = id;
    if (subjectId != null) result.subjectId = subjectId;
    if (timestamp != null) result.timestamp = timestamp;
    if (latitude != null) result.latitude = latitude;
    if (longitude != null) result.longitude = longitude;
    if (altitude != null) result.altitude = altitude;
    if (accuracy != null) result.accuracy = accuracy;
    if (speed != null) result.speed = speed;
    if (bearing != null) result.bearing = bearing;
    if (source != null) result.source = source;
    if (extra != null) result.extra = extra;
    if (createdAt != null) result.createdAt = createdAt;
    if (deviceId != null) result.deviceId = deviceId;
    return result;
  }

  LocationPointObject._();

  factory LocationPointObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocationPointObject()..mergeFromBuffer(data, registry);
  factory LocationPointObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      LocationPointObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocationPointObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: LocationPointObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$0.Timestamp>(3, _omitFieldNames ? '' : 'timestamp',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aD(4, _omitFieldNames ? '' : 'latitude')
    ..aD(5, _omitFieldNames ? '' : 'longitude')
    ..aD(6, _omitFieldNames ? '' : 'altitude')
    ..aD(7, _omitFieldNames ? '' : 'accuracy')
    ..aD(8, _omitFieldNames ? '' : 'speed')
    ..aD(9, _omitFieldNames ? '' : 'bearing')
    ..aE<LocationSource>(10, _omitFieldNames ? '' : 'source',
        enumValues: LocationSource.values)
    ..aOM<$1.Struct>(11, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..aOM<$0.Timestamp>(12, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aOS(13, _omitFieldNames ? '' : 'deviceId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocationPointObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocationPointObject copyWith(void Function(LocationPointObject) updates) =>
      super.copyWith((message) => updates(message as LocationPointObject))
          as LocationPointObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use LocationPointObject() / LocationPointObject.new instead')
  static LocationPointObject create() => LocationPointObject._();
  static $pb.GeneratedMessage $_createMessage() => LocationPointObject._();
  @$core.override
  LocationPointObject createEmptyInstance() => LocationPointObject._();
  @$core.pragma('dart2js:noInline')
  static LocationPointObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocationPointObject>(
          LocationPointObject.$_createMessage);
  static LocationPointObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Timestamp get timestamp => $_getN(2);
  @$pb.TagNumber(3)
  set timestamp($0.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(3)
  void clearTimestamp() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Timestamp ensureTimestamp() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.double get latitude => $_getN(3);
  @$pb.TagNumber(4)
  set latitude($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLatitude() => $_has(3);
  @$pb.TagNumber(4)
  void clearLatitude() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.double get longitude => $_getN(4);
  @$pb.TagNumber(5)
  set longitude($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(5)
  $core.bool hasLongitude() => $_has(4);
  @$pb.TagNumber(5)
  void clearLongitude() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.double get altitude => $_getN(5);
  @$pb.TagNumber(6)
  set altitude($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(6)
  $core.bool hasAltitude() => $_has(5);
  @$pb.TagNumber(6)
  void clearAltitude() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.double get accuracy => $_getN(6);
  @$pb.TagNumber(7)
  set accuracy($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(7)
  $core.bool hasAccuracy() => $_has(6);
  @$pb.TagNumber(7)
  void clearAccuracy() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.double get speed => $_getN(7);
  @$pb.TagNumber(8)
  set speed($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(8)
  $core.bool hasSpeed() => $_has(7);
  @$pb.TagNumber(8)
  void clearSpeed() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.double get bearing => $_getN(8);
  @$pb.TagNumber(9)
  set bearing($core.double value) => $_setDouble(8, value);
  @$pb.TagNumber(9)
  $core.bool hasBearing() => $_has(8);
  @$pb.TagNumber(9)
  void clearBearing() => $_clearField(9);

  @$pb.TagNumber(10)
  LocationSource get source => $_getN(9);
  @$pb.TagNumber(10)
  set source(LocationSource value) => $_setField(10, value);
  @$pb.TagNumber(10)
  $core.bool hasSource() => $_has(9);
  @$pb.TagNumber(10)
  void clearSource() => $_clearField(10);

  @$pb.TagNumber(11)
  $1.Struct get extra => $_getN(10);
  @$pb.TagNumber(11)
  set extra($1.Struct value) => $_setField(11, value);
  @$pb.TagNumber(11)
  $core.bool hasExtra() => $_has(10);
  @$pb.TagNumber(11)
  void clearExtra() => $_clearField(11);
  @$pb.TagNumber(11)
  $1.Struct ensureExtra() => $_ensure(10);

  @$pb.TagNumber(12)
  $0.Timestamp get createdAt => $_getN(11);
  @$pb.TagNumber(12)
  set createdAt($0.Timestamp value) => $_setField(12, value);
  @$pb.TagNumber(12)
  $core.bool hasCreatedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearCreatedAt() => $_clearField(12);
  @$pb.TagNumber(12)
  $0.Timestamp ensureCreatedAt() => $_ensure(11);

  @$pb.TagNumber(13)
  $core.String get deviceId => $_getSZ(12);
  @$pb.TagNumber(13)
  set deviceId($core.String value) => $_setString(12, value);
  @$pb.TagNumber(13)
  $core.bool hasDeviceId() => $_has(12);
  @$pb.TagNumber(13)
  void clearDeviceId() => $_clearField(13);
}

class AreaObject extends $pb.GeneratedMessage {
  factory AreaObject({
    $core.String? id,
    $core.String? ownerId,
    $core.String? name,
    $core.String? description,
    AreaType? areaType,
    $core.String? geometry,
    $core.double? areaM2,
    $core.double? perimeterM,
    $core.int? state,
    $1.Struct? extra,
    $0.Timestamp? createdAt,
  }) {
    final result = AreaObject._();
    if (id != null) result.id = id;
    if (ownerId != null) result.ownerId = ownerId;
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (areaType != null) result.areaType = areaType;
    if (geometry != null) result.geometry = geometry;
    if (areaM2 != null) result.areaM2 = areaM2;
    if (perimeterM != null) result.perimeterM = perimeterM;
    if (state != null) result.state = state;
    if (extra != null) result.extra = extra;
    if (createdAt != null) result.createdAt = createdAt;
    return result;
  }

  AreaObject._();

  factory AreaObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AreaObject()..mergeFromBuffer(data, registry);
  factory AreaObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AreaObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AreaObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: AreaObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'ownerId')
    ..aOS(3, _omitFieldNames ? '' : 'name')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aE<AreaType>(5, _omitFieldNames ? '' : 'areaType',
        enumValues: AreaType.values)
    ..aOS(6, _omitFieldNames ? '' : 'geometry')
    ..aD(7, _omitFieldNames ? '' : 'areaM2')
    ..aD(8, _omitFieldNames ? '' : 'perimeterM')
    ..aI(9, _omitFieldNames ? '' : 'state')
    ..aOM<$1.Struct>(10, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..aOM<$0.Timestamp>(11, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $0.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AreaObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AreaObject copyWith(void Function(AreaObject) updates) =>
      super.copyWith((message) => updates(message as AreaObject)) as AreaObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AreaObject() / AreaObject.new instead')
  static AreaObject create() => AreaObject._();
  static $pb.GeneratedMessage $_createMessage() => AreaObject._();
  @$core.override
  AreaObject createEmptyInstance() => AreaObject._();
  @$core.pragma('dart2js:noInline')
  static AreaObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AreaObject>(AreaObject.$_createMessage);
  static AreaObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get ownerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ownerId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOwnerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearOwnerId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(3)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(3)
  void clearName() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  @$pb.TagNumber(5)
  AreaType get areaType => $_getN(4);
  @$pb.TagNumber(5)
  set areaType(AreaType value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasAreaType() => $_has(4);
  @$pb.TagNumber(5)
  void clearAreaType() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get geometry => $_getSZ(5);
  @$pb.TagNumber(6)
  set geometry($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasGeometry() => $_has(5);
  @$pb.TagNumber(6)
  void clearGeometry() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.double get areaM2 => $_getN(6);
  @$pb.TagNumber(7)
  set areaM2($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(7)
  $core.bool hasAreaM2() => $_has(6);
  @$pb.TagNumber(7)
  void clearAreaM2() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.double get perimeterM => $_getN(7);
  @$pb.TagNumber(8)
  set perimeterM($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPerimeterM() => $_has(7);
  @$pb.TagNumber(8)
  void clearPerimeterM() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get state => $_getIZ(8);
  @$pb.TagNumber(9)
  set state($core.int value) => $_setSignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasState() => $_has(8);
  @$pb.TagNumber(9)
  void clearState() => $_clearField(9);

  @$pb.TagNumber(10)
  $1.Struct get extra => $_getN(9);
  @$pb.TagNumber(10)
  set extra($1.Struct value) => $_setField(10, value);
  @$pb.TagNumber(10)
  $core.bool hasExtra() => $_has(9);
  @$pb.TagNumber(10)
  void clearExtra() => $_clearField(10);
  @$pb.TagNumber(10)
  $1.Struct ensureExtra() => $_ensure(9);

  @$pb.TagNumber(11)
  $0.Timestamp get createdAt => $_getN(10);
  @$pb.TagNumber(11)
  set createdAt($0.Timestamp value) => $_setField(11, value);
  @$pb.TagNumber(11)
  $core.bool hasCreatedAt() => $_has(10);
  @$pb.TagNumber(11)
  void clearCreatedAt() => $_clearField(11);
  @$pb.TagNumber(11)
  $0.Timestamp ensureCreatedAt() => $_ensure(10);
}

class GeoEventObject extends $pb.GeneratedMessage {
  factory GeoEventObject({
    $core.String? id,
    $core.String? subjectId,
    $core.String? areaId,
    GeoEventType? eventType,
    $0.Timestamp? timestamp,
    $core.double? confidence,
    $core.String? pointId,
    $1.Struct? extra,
  }) {
    final result = GeoEventObject._();
    if (id != null) result.id = id;
    if (subjectId != null) result.subjectId = subjectId;
    if (areaId != null) result.areaId = areaId;
    if (eventType != null) result.eventType = eventType;
    if (timestamp != null) result.timestamp = timestamp;
    if (confidence != null) result.confidence = confidence;
    if (pointId != null) result.pointId = pointId;
    if (extra != null) result.extra = extra;
    return result;
  }

  GeoEventObject._();

  factory GeoEventObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GeoEventObject()..mergeFromBuffer(data, registry);
  factory GeoEventObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GeoEventObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoEventObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GeoEventObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOS(3, _omitFieldNames ? '' : 'areaId')
    ..aE<GeoEventType>(4, _omitFieldNames ? '' : 'eventType',
        enumValues: GeoEventType.values)
    ..aOM<$0.Timestamp>(5, _omitFieldNames ? '' : 'timestamp',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aD(6, _omitFieldNames ? '' : 'confidence')
    ..aOS(7, _omitFieldNames ? '' : 'pointId')
    ..aOM<$1.Struct>(8, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoEventObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoEventObject copyWith(void Function(GeoEventObject) updates) =>
      super.copyWith((message) => updates(message as GeoEventObject))
          as GeoEventObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GeoEventObject() / GeoEventObject.new instead')
  static GeoEventObject create() => GeoEventObject._();
  static $pb.GeneratedMessage $_createMessage() => GeoEventObject._();
  @$core.override
  GeoEventObject createEmptyInstance() => GeoEventObject._();
  @$core.pragma('dart2js:noInline')
  static GeoEventObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GeoEventObject>(
          GeoEventObject.$_createMessage);
  static GeoEventObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get areaId => $_getSZ(2);
  @$pb.TagNumber(3)
  set areaId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAreaId() => $_has(2);
  @$pb.TagNumber(3)
  void clearAreaId() => $_clearField(3);

  @$pb.TagNumber(4)
  GeoEventType get eventType => $_getN(3);
  @$pb.TagNumber(4)
  set eventType(GeoEventType value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasEventType() => $_has(3);
  @$pb.TagNumber(4)
  void clearEventType() => $_clearField(4);

  @$pb.TagNumber(5)
  $0.Timestamp get timestamp => $_getN(4);
  @$pb.TagNumber(5)
  set timestamp($0.Timestamp value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasTimestamp() => $_has(4);
  @$pb.TagNumber(5)
  void clearTimestamp() => $_clearField(5);
  @$pb.TagNumber(5)
  $0.Timestamp ensureTimestamp() => $_ensure(4);

  @$pb.TagNumber(6)
  $core.double get confidence => $_getN(5);
  @$pb.TagNumber(6)
  set confidence($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(6)
  $core.bool hasConfidence() => $_has(5);
  @$pb.TagNumber(6)
  void clearConfidence() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get pointId => $_getSZ(6);
  @$pb.TagNumber(7)
  set pointId($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasPointId() => $_has(6);
  @$pb.TagNumber(7)
  void clearPointId() => $_clearField(7);

  @$pb.TagNumber(8)
  $1.Struct get extra => $_getN(7);
  @$pb.TagNumber(8)
  set extra($1.Struct value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasExtra() => $_has(7);
  @$pb.TagNumber(8)
  void clearExtra() => $_clearField(8);
  @$pb.TagNumber(8)
  $1.Struct ensureExtra() => $_ensure(7);
}

class AreaSubjectObject extends $pb.GeneratedMessage {
  factory AreaSubjectObject({
    $core.String? subjectId,
    $0.Timestamp? enterTimestamp,
  }) {
    final result = AreaSubjectObject._();
    if (subjectId != null) result.subjectId = subjectId;
    if (enterTimestamp != null) result.enterTimestamp = enterTimestamp;
    return result;
  }

  AreaSubjectObject._();

  factory AreaSubjectObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AreaSubjectObject()..mergeFromBuffer(data, registry);
  factory AreaSubjectObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AreaSubjectObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AreaSubjectObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: AreaSubjectObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$0.Timestamp>(2, _omitFieldNames ? '' : 'enterTimestamp',
        subBuilder: $0.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AreaSubjectObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AreaSubjectObject copyWith(void Function(AreaSubjectObject) updates) =>
      super.copyWith((message) => updates(message as AreaSubjectObject))
          as AreaSubjectObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AreaSubjectObject() / AreaSubjectObject.new instead')
  static AreaSubjectObject create() => AreaSubjectObject._();
  static $pb.GeneratedMessage $_createMessage() => AreaSubjectObject._();
  @$core.override
  AreaSubjectObject createEmptyInstance() => AreaSubjectObject._();
  @$core.pragma('dart2js:noInline')
  static AreaSubjectObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AreaSubjectObject>(
          AreaSubjectObject.$_createMessage);
  static AreaSubjectObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $0.Timestamp get enterTimestamp => $_getN(1);
  @$pb.TagNumber(2)
  set enterTimestamp($0.Timestamp value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasEnterTimestamp() => $_has(1);
  @$pb.TagNumber(2)
  void clearEnterTimestamp() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.Timestamp ensureEnterTimestamp() => $_ensure(1);
}

class NearbySubjectObject extends $pb.GeneratedMessage {
  factory NearbySubjectObject({
    $core.String? subjectId,
    $core.double? distanceMeters,
    $0.Timestamp? lastSeen,
  }) {
    final result = NearbySubjectObject._();
    if (subjectId != null) result.subjectId = subjectId;
    if (distanceMeters != null) result.distanceMeters = distanceMeters;
    if (lastSeen != null) result.lastSeen = lastSeen;
    return result;
  }

  NearbySubjectObject._();

  factory NearbySubjectObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NearbySubjectObject()..mergeFromBuffer(data, registry);
  factory NearbySubjectObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NearbySubjectObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NearbySubjectObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: NearbySubjectObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aD(2, _omitFieldNames ? '' : 'distanceMeters')
    ..aOM<$0.Timestamp>(3, _omitFieldNames ? '' : 'lastSeen',
        subBuilder: $0.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NearbySubjectObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NearbySubjectObject copyWith(void Function(NearbySubjectObject) updates) =>
      super.copyWith((message) => updates(message as NearbySubjectObject))
          as NearbySubjectObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use NearbySubjectObject() / NearbySubjectObject.new instead')
  static NearbySubjectObject create() => NearbySubjectObject._();
  static $pb.GeneratedMessage $_createMessage() => NearbySubjectObject._();
  @$core.override
  NearbySubjectObject createEmptyInstance() => NearbySubjectObject._();
  @$core.pragma('dart2js:noInline')
  static NearbySubjectObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NearbySubjectObject>(
          NearbySubjectObject.$_createMessage);
  static NearbySubjectObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get distanceMeters => $_getN(1);
  @$pb.TagNumber(2)
  set distanceMeters($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDistanceMeters() => $_has(1);
  @$pb.TagNumber(2)
  void clearDistanceMeters() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Timestamp get lastSeen => $_getN(2);
  @$pb.TagNumber(3)
  set lastSeen($0.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasLastSeen() => $_has(2);
  @$pb.TagNumber(3)
  void clearLastSeen() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Timestamp ensureLastSeen() => $_ensure(2);
}

class NearbyAreaObject extends $pb.GeneratedMessage {
  factory NearbyAreaObject({
    $core.String? areaId,
    $core.String? name,
    AreaType? areaType,
    $core.double? distanceMeters,
  }) {
    final result = NearbyAreaObject._();
    if (areaId != null) result.areaId = areaId;
    if (name != null) result.name = name;
    if (areaType != null) result.areaType = areaType;
    if (distanceMeters != null) result.distanceMeters = distanceMeters;
    return result;
  }

  NearbyAreaObject._();

  factory NearbyAreaObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NearbyAreaObject()..mergeFromBuffer(data, registry);
  factory NearbyAreaObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NearbyAreaObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NearbyAreaObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: NearbyAreaObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'areaId')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aE<AreaType>(3, _omitFieldNames ? '' : 'areaType',
        enumValues: AreaType.values)
    ..aD(4, _omitFieldNames ? '' : 'distanceMeters')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NearbyAreaObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NearbyAreaObject copyWith(void Function(NearbyAreaObject) updates) =>
      super.copyWith((message) => updates(message as NearbyAreaObject))
          as NearbyAreaObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use NearbyAreaObject() / NearbyAreaObject.new instead')
  static NearbyAreaObject create() => NearbyAreaObject._();
  static $pb.GeneratedMessage $_createMessage() => NearbyAreaObject._();
  @$core.override
  NearbyAreaObject createEmptyInstance() => NearbyAreaObject._();
  @$core.pragma('dart2js:noInline')
  static NearbyAreaObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<NearbyAreaObject>(
          NearbyAreaObject.$_createMessage);
  static NearbyAreaObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get areaId => $_getSZ(0);
  @$pb.TagNumber(1)
  set areaId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAreaId() => $_has(0);
  @$pb.TagNumber(1)
  void clearAreaId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  AreaType get areaType => $_getN(2);
  @$pb.TagNumber(3)
  set areaType(AreaType value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasAreaType() => $_has(2);
  @$pb.TagNumber(3)
  void clearAreaType() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.double get distanceMeters => $_getN(3);
  @$pb.TagNumber(4)
  set distanceMeters($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDistanceMeters() => $_has(3);
  @$pb.TagNumber(4)
  void clearDistanceMeters() => $_clearField(4);
}

class RouteObject extends $pb.GeneratedMessage {
  factory RouteObject({
    $core.String? id,
    $core.String? ownerId,
    $core.String? name,
    $core.String? description,
    $core.String? geometry,
    $core.double? lengthM,
    $core.int? state,
    $core.double? deviationThresholdM,
    $core.int? deviationConsecutiveCount,
    $core.int? deviationCooldownSec,
    $1.Struct? extra,
    $0.Timestamp? createdAt,
  }) {
    final result = RouteObject._();
    if (id != null) result.id = id;
    if (ownerId != null) result.ownerId = ownerId;
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (geometry != null) result.geometry = geometry;
    if (lengthM != null) result.lengthM = lengthM;
    if (state != null) result.state = state;
    if (deviationThresholdM != null)
      result.deviationThresholdM = deviationThresholdM;
    if (deviationConsecutiveCount != null)
      result.deviationConsecutiveCount = deviationConsecutiveCount;
    if (deviationCooldownSec != null)
      result.deviationCooldownSec = deviationCooldownSec;
    if (extra != null) result.extra = extra;
    if (createdAt != null) result.createdAt = createdAt;
    return result;
  }

  RouteObject._();

  factory RouteObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RouteObject()..mergeFromBuffer(data, registry);
  factory RouteObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RouteObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RouteObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: RouteObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'ownerId')
    ..aOS(3, _omitFieldNames ? '' : 'name')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aOS(5, _omitFieldNames ? '' : 'geometry')
    ..aD(6, _omitFieldNames ? '' : 'lengthM')
    ..aI(7, _omitFieldNames ? '' : 'state')
    ..aD(8, _omitFieldNames ? '' : 'deviationThresholdM')
    ..aI(9, _omitFieldNames ? '' : 'deviationConsecutiveCount')
    ..aI(10, _omitFieldNames ? '' : 'deviationCooldownSec')
    ..aOM<$1.Struct>(11, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..aOM<$0.Timestamp>(12, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $0.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RouteObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RouteObject copyWith(void Function(RouteObject) updates) =>
      super.copyWith((message) => updates(message as RouteObject))
          as RouteObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use RouteObject() / RouteObject.new instead')
  static RouteObject create() => RouteObject._();
  static $pb.GeneratedMessage $_createMessage() => RouteObject._();
  @$core.override
  RouteObject createEmptyInstance() => RouteObject._();
  @$core.pragma('dart2js:noInline')
  static RouteObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RouteObject>(
          RouteObject.$_createMessage);
  static RouteObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get ownerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ownerId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOwnerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearOwnerId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(3)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(3)
  void clearName() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get geometry => $_getSZ(4);
  @$pb.TagNumber(5)
  set geometry($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasGeometry() => $_has(4);
  @$pb.TagNumber(5)
  void clearGeometry() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.double get lengthM => $_getN(5);
  @$pb.TagNumber(6)
  set lengthM($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(6)
  $core.bool hasLengthM() => $_has(5);
  @$pb.TagNumber(6)
  void clearLengthM() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.int get state => $_getIZ(6);
  @$pb.TagNumber(7)
  set state($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasState() => $_has(6);
  @$pb.TagNumber(7)
  void clearState() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.double get deviationThresholdM => $_getN(7);
  @$pb.TagNumber(8)
  set deviationThresholdM($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(8)
  $core.bool hasDeviationThresholdM() => $_has(7);
  @$pb.TagNumber(8)
  void clearDeviationThresholdM() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.int get deviationConsecutiveCount => $_getIZ(8);
  @$pb.TagNumber(9)
  set deviationConsecutiveCount($core.int value) => $_setSignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasDeviationConsecutiveCount() => $_has(8);
  @$pb.TagNumber(9)
  void clearDeviationConsecutiveCount() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.int get deviationCooldownSec => $_getIZ(9);
  @$pb.TagNumber(10)
  set deviationCooldownSec($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDeviationCooldownSec() => $_has(9);
  @$pb.TagNumber(10)
  void clearDeviationCooldownSec() => $_clearField(10);

  @$pb.TagNumber(11)
  $1.Struct get extra => $_getN(10);
  @$pb.TagNumber(11)
  set extra($1.Struct value) => $_setField(11, value);
  @$pb.TagNumber(11)
  $core.bool hasExtra() => $_has(10);
  @$pb.TagNumber(11)
  void clearExtra() => $_clearField(11);
  @$pb.TagNumber(11)
  $1.Struct ensureExtra() => $_ensure(10);

  @$pb.TagNumber(12)
  $0.Timestamp get createdAt => $_getN(11);
  @$pb.TagNumber(12)
  set createdAt($0.Timestamp value) => $_setField(12, value);
  @$pb.TagNumber(12)
  $core.bool hasCreatedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearCreatedAt() => $_clearField(12);
  @$pb.TagNumber(12)
  $0.Timestamp ensureCreatedAt() => $_ensure(11);
}

class RouteAssignmentObject extends $pb.GeneratedMessage {
  factory RouteAssignmentObject({
    $core.String? id,
    $core.String? subjectId,
    $core.String? routeId,
    $0.Timestamp? validFrom,
    $0.Timestamp? validUntil,
    $core.int? state,
    $1.Struct? extra,
    $0.Timestamp? createdAt,
  }) {
    final result = RouteAssignmentObject._();
    if (id != null) result.id = id;
    if (subjectId != null) result.subjectId = subjectId;
    if (routeId != null) result.routeId = routeId;
    if (validFrom != null) result.validFrom = validFrom;
    if (validUntil != null) result.validUntil = validUntil;
    if (state != null) result.state = state;
    if (extra != null) result.extra = extra;
    if (createdAt != null) result.createdAt = createdAt;
    return result;
  }

  RouteAssignmentObject._();

  factory RouteAssignmentObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RouteAssignmentObject()..mergeFromBuffer(data, registry);
  factory RouteAssignmentObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RouteAssignmentObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RouteAssignmentObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: RouteAssignmentObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOS(3, _omitFieldNames ? '' : 'routeId')
    ..aOM<$0.Timestamp>(4, _omitFieldNames ? '' : 'validFrom',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aOM<$0.Timestamp>(5, _omitFieldNames ? '' : 'validUntil',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aI(6, _omitFieldNames ? '' : 'state')
    ..aOM<$1.Struct>(7, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..aOM<$0.Timestamp>(8, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $0.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RouteAssignmentObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RouteAssignmentObject copyWith(
          void Function(RouteAssignmentObject) updates) =>
      super.copyWith((message) => updates(message as RouteAssignmentObject))
          as RouteAssignmentObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RouteAssignmentObject() / RouteAssignmentObject.new instead')
  static RouteAssignmentObject create() => RouteAssignmentObject._();
  static $pb.GeneratedMessage $_createMessage() => RouteAssignmentObject._();
  @$core.override
  RouteAssignmentObject createEmptyInstance() => RouteAssignmentObject._();
  @$core.pragma('dart2js:noInline')
  static RouteAssignmentObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RouteAssignmentObject>(
          RouteAssignmentObject.$_createMessage);
  static RouteAssignmentObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get routeId => $_getSZ(2);
  @$pb.TagNumber(3)
  set routeId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasRouteId() => $_has(2);
  @$pb.TagNumber(3)
  void clearRouteId() => $_clearField(3);

  @$pb.TagNumber(4)
  $0.Timestamp get validFrom => $_getN(3);
  @$pb.TagNumber(4)
  set validFrom($0.Timestamp value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasValidFrom() => $_has(3);
  @$pb.TagNumber(4)
  void clearValidFrom() => $_clearField(4);
  @$pb.TagNumber(4)
  $0.Timestamp ensureValidFrom() => $_ensure(3);

  @$pb.TagNumber(5)
  $0.Timestamp get validUntil => $_getN(4);
  @$pb.TagNumber(5)
  set validUntil($0.Timestamp value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasValidUntil() => $_has(4);
  @$pb.TagNumber(5)
  void clearValidUntil() => $_clearField(5);
  @$pb.TagNumber(5)
  $0.Timestamp ensureValidUntil() => $_ensure(4);

  @$pb.TagNumber(6)
  $core.int get state => $_getIZ(5);
  @$pb.TagNumber(6)
  set state($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(6)
  void clearState() => $_clearField(6);

  @$pb.TagNumber(7)
  $1.Struct get extra => $_getN(6);
  @$pb.TagNumber(7)
  set extra($1.Struct value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasExtra() => $_has(6);
  @$pb.TagNumber(7)
  void clearExtra() => $_clearField(7);
  @$pb.TagNumber(7)
  $1.Struct ensureExtra() => $_ensure(6);

  @$pb.TagNumber(8)
  $0.Timestamp get createdAt => $_getN(7);
  @$pb.TagNumber(8)
  set createdAt($0.Timestamp value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasCreatedAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearCreatedAt() => $_clearField(8);
  @$pb.TagNumber(8)
  $0.Timestamp ensureCreatedAt() => $_ensure(7);
}

class RouteDeviationEventObject extends $pb.GeneratedMessage {
  factory RouteDeviationEventObject({
    $core.String? id,
    $core.String? subjectId,
    $core.String? routeId,
    RouteDeviationEventType? eventType,
    $core.double? distanceMeters,
    $core.double? latitude,
    $core.double? longitude,
    $0.Timestamp? timestamp,
    $1.Struct? extra,
  }) {
    final result = RouteDeviationEventObject._();
    if (id != null) result.id = id;
    if (subjectId != null) result.subjectId = subjectId;
    if (routeId != null) result.routeId = routeId;
    if (eventType != null) result.eventType = eventType;
    if (distanceMeters != null) result.distanceMeters = distanceMeters;
    if (latitude != null) result.latitude = latitude;
    if (longitude != null) result.longitude = longitude;
    if (timestamp != null) result.timestamp = timestamp;
    if (extra != null) result.extra = extra;
    return result;
  }

  RouteDeviationEventObject._();

  factory RouteDeviationEventObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RouteDeviationEventObject()..mergeFromBuffer(data, registry);
  factory RouteDeviationEventObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RouteDeviationEventObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RouteDeviationEventObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: RouteDeviationEventObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOS(3, _omitFieldNames ? '' : 'routeId')
    ..aE<RouteDeviationEventType>(4, _omitFieldNames ? '' : 'eventType',
        enumValues: RouteDeviationEventType.values)
    ..aD(5, _omitFieldNames ? '' : 'distanceMeters')
    ..aD(6, _omitFieldNames ? '' : 'latitude')
    ..aD(7, _omitFieldNames ? '' : 'longitude')
    ..aOM<$0.Timestamp>(8, _omitFieldNames ? '' : 'timestamp',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aOM<$1.Struct>(9, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RouteDeviationEventObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RouteDeviationEventObject copyWith(
          void Function(RouteDeviationEventObject) updates) =>
      super.copyWith((message) => updates(message as RouteDeviationEventObject))
          as RouteDeviationEventObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RouteDeviationEventObject() / RouteDeviationEventObject.new instead')
  static RouteDeviationEventObject create() => RouteDeviationEventObject._();
  static $pb.GeneratedMessage $_createMessage() =>
      RouteDeviationEventObject._();
  @$core.override
  RouteDeviationEventObject createEmptyInstance() =>
      RouteDeviationEventObject._();
  @$core.pragma('dart2js:noInline')
  static RouteDeviationEventObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RouteDeviationEventObject>(
          RouteDeviationEventObject.$_createMessage);
  static RouteDeviationEventObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get routeId => $_getSZ(2);
  @$pb.TagNumber(3)
  set routeId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasRouteId() => $_has(2);
  @$pb.TagNumber(3)
  void clearRouteId() => $_clearField(3);

  @$pb.TagNumber(4)
  RouteDeviationEventType get eventType => $_getN(3);
  @$pb.TagNumber(4)
  set eventType(RouteDeviationEventType value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasEventType() => $_has(3);
  @$pb.TagNumber(4)
  void clearEventType() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.double get distanceMeters => $_getN(4);
  @$pb.TagNumber(5)
  set distanceMeters($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDistanceMeters() => $_has(4);
  @$pb.TagNumber(5)
  void clearDistanceMeters() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.double get latitude => $_getN(5);
  @$pb.TagNumber(6)
  set latitude($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(6)
  $core.bool hasLatitude() => $_has(5);
  @$pb.TagNumber(6)
  void clearLatitude() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.double get longitude => $_getN(6);
  @$pb.TagNumber(7)
  set longitude($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(7)
  $core.bool hasLongitude() => $_has(6);
  @$pb.TagNumber(7)
  void clearLongitude() => $_clearField(7);

  @$pb.TagNumber(8)
  $0.Timestamp get timestamp => $_getN(7);
  @$pb.TagNumber(8)
  set timestamp($0.Timestamp value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasTimestamp() => $_has(7);
  @$pb.TagNumber(8)
  void clearTimestamp() => $_clearField(8);
  @$pb.TagNumber(8)
  $0.Timestamp ensureTimestamp() => $_ensure(7);

  @$pb.TagNumber(9)
  $1.Struct get extra => $_getN(8);
  @$pb.TagNumber(9)
  set extra($1.Struct value) => $_setField(9, value);
  @$pb.TagNumber(9)
  $core.bool hasExtra() => $_has(8);
  @$pb.TagNumber(9)
  void clearExtra() => $_clearField(9);
  @$pb.TagNumber(9)
  $1.Struct ensureExtra() => $_ensure(8);
}

class IngestLocationsRequest extends $pb.GeneratedMessage {
  factory IngestLocationsRequest({
    $core.String? subjectId,
    $core.Iterable<LocationPointInput>? points,
  }) {
    final result = IngestLocationsRequest._();
    if (subjectId != null) result.subjectId = subjectId;
    if (points != null) result.points.addAll(points);
    return result;
  }

  IngestLocationsRequest._();

  factory IngestLocationsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      IngestLocationsRequest()..mergeFromBuffer(data, registry);
  factory IngestLocationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      IngestLocationsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IngestLocationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: IngestLocationsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..pPM<LocationPointInput>(2, _omitFieldNames ? '' : 'points',
        subBuilder: LocationPointInput.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IngestLocationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IngestLocationsRequest copyWith(
          void Function(IngestLocationsRequest) updates) =>
      super.copyWith((message) => updates(message as IngestLocationsRequest))
          as IngestLocationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use IngestLocationsRequest() / IngestLocationsRequest.new instead')
  static IngestLocationsRequest create() => IngestLocationsRequest._();
  static $pb.GeneratedMessage $_createMessage() => IngestLocationsRequest._();
  @$core.override
  IngestLocationsRequest createEmptyInstance() => IngestLocationsRequest._();
  @$core.pragma('dart2js:noInline')
  static IngestLocationsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IngestLocationsRequest>(
          IngestLocationsRequest.$_createMessage);
  static IngestLocationsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<LocationPointInput> get points => $_getList(1);
}

class IngestLocationsResponse extends $pb.GeneratedMessage {
  factory IngestLocationsResponse({
    $core.int? accepted,
    $core.int? rejected,
  }) {
    final result = IngestLocationsResponse._();
    if (accepted != null) result.accepted = accepted;
    if (rejected != null) result.rejected = rejected;
    return result;
  }

  IngestLocationsResponse._();

  factory IngestLocationsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      IngestLocationsResponse()..mergeFromBuffer(data, registry);
  factory IngestLocationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      IngestLocationsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IngestLocationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: IngestLocationsResponse.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'accepted')
    ..aI(2, _omitFieldNames ? '' : 'rejected')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IngestLocationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IngestLocationsResponse copyWith(
          void Function(IngestLocationsResponse) updates) =>
      super.copyWith((message) => updates(message as IngestLocationsResponse))
          as IngestLocationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use IngestLocationsResponse() / IngestLocationsResponse.new instead')
  static IngestLocationsResponse create() => IngestLocationsResponse._();
  static $pb.GeneratedMessage $_createMessage() => IngestLocationsResponse._();
  @$core.override
  IngestLocationsResponse createEmptyInstance() => IngestLocationsResponse._();
  @$core.pragma('dart2js:noInline')
  static IngestLocationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IngestLocationsResponse>(
          IngestLocationsResponse.$_createMessage);
  static IngestLocationsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get accepted => $_getIZ(0);
  @$pb.TagNumber(1)
  set accepted($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAccepted() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccepted() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get rejected => $_getIZ(1);
  @$pb.TagNumber(2)
  set rejected($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRejected() => $_has(1);
  @$pb.TagNumber(2)
  void clearRejected() => $_clearField(2);
}

class CreateAreaRequest extends $pb.GeneratedMessage {
  factory CreateAreaRequest({
    AreaObject? data,
  }) {
    final result = CreateAreaRequest._();
    if (data != null) result.data = data;
    return result;
  }

  CreateAreaRequest._();

  factory CreateAreaRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateAreaRequest()..mergeFromBuffer(data, registry);
  factory CreateAreaRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateAreaRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateAreaRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: CreateAreaRequest.$_createMessage)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AreaObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAreaRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAreaRequest copyWith(void Function(CreateAreaRequest) updates) =>
      super.copyWith((message) => updates(message as CreateAreaRequest))
          as CreateAreaRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CreateAreaRequest() / CreateAreaRequest.new instead')
  static CreateAreaRequest create() => CreateAreaRequest._();
  static $pb.GeneratedMessage $_createMessage() => CreateAreaRequest._();
  @$core.override
  CreateAreaRequest createEmptyInstance() => CreateAreaRequest._();
  @$core.pragma('dart2js:noInline')
  static CreateAreaRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateAreaRequest>(
          CreateAreaRequest.$_createMessage);
  static CreateAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class CreateAreaResponse extends $pb.GeneratedMessage {
  factory CreateAreaResponse({
    AreaObject? data,
  }) {
    final result = CreateAreaResponse._();
    if (data != null) result.data = data;
    return result;
  }

  CreateAreaResponse._();

  factory CreateAreaResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateAreaResponse()..mergeFromBuffer(data, registry);
  factory CreateAreaResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateAreaResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateAreaResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: CreateAreaResponse.$_createMessage)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AreaObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAreaResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAreaResponse copyWith(void Function(CreateAreaResponse) updates) =>
      super.copyWith((message) => updates(message as CreateAreaResponse))
          as CreateAreaResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CreateAreaResponse() / CreateAreaResponse.new instead')
  static CreateAreaResponse create() => CreateAreaResponse._();
  static $pb.GeneratedMessage $_createMessage() => CreateAreaResponse._();
  @$core.override
  CreateAreaResponse createEmptyInstance() => CreateAreaResponse._();
  @$core.pragma('dart2js:noInline')
  static CreateAreaResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateAreaResponse>(
          CreateAreaResponse.$_createMessage);
  static CreateAreaResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class GetAreaRequest extends $pb.GeneratedMessage {
  factory GetAreaRequest({
    $core.String? id,
  }) {
    final result = GetAreaRequest._();
    if (id != null) result.id = id;
    return result;
  }

  GetAreaRequest._();

  factory GetAreaRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaRequest()..mergeFromBuffer(data, registry);
  factory GetAreaRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAreaRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetAreaRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaRequest copyWith(void Function(GetAreaRequest) updates) =>
      super.copyWith((message) => updates(message as GetAreaRequest))
          as GetAreaRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetAreaRequest() / GetAreaRequest.new instead')
  static GetAreaRequest create() => GetAreaRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetAreaRequest._();
  @$core.override
  GetAreaRequest createEmptyInstance() => GetAreaRequest._();
  @$core.pragma('dart2js:noInline')
  static GetAreaRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAreaRequest>(
          GetAreaRequest.$_createMessage);
  static GetAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class GetAreaResponse extends $pb.GeneratedMessage {
  factory GetAreaResponse({
    AreaObject? data,
  }) {
    final result = GetAreaResponse._();
    if (data != null) result.data = data;
    return result;
  }

  GetAreaResponse._();

  factory GetAreaResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaResponse()..mergeFromBuffer(data, registry);
  factory GetAreaResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAreaResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetAreaResponse.$_createMessage)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AreaObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaResponse copyWith(void Function(GetAreaResponse) updates) =>
      super.copyWith((message) => updates(message as GetAreaResponse))
          as GetAreaResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetAreaResponse() / GetAreaResponse.new instead')
  static GetAreaResponse create() => GetAreaResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetAreaResponse._();
  @$core.override
  GetAreaResponse createEmptyInstance() => GetAreaResponse._();
  @$core.pragma('dart2js:noInline')
  static GetAreaResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAreaResponse>(
          GetAreaResponse.$_createMessage);
  static GetAreaResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class UpdateAreaRequest extends $pb.GeneratedMessage {
  factory UpdateAreaRequest({
    $core.String? id,
    $core.String? name,
    $core.String? description,
    AreaType? areaType,
    $core.String? geometry,
    $1.Struct? extra,
  }) {
    final result = UpdateAreaRequest._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (areaType != null) result.areaType = areaType;
    if (geometry != null) result.geometry = geometry;
    if (extra != null) result.extra = extra;
    return result;
  }

  UpdateAreaRequest._();

  factory UpdateAreaRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateAreaRequest()..mergeFromBuffer(data, registry);
  factory UpdateAreaRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateAreaRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAreaRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: UpdateAreaRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..aE<AreaType>(4, _omitFieldNames ? '' : 'areaType',
        enumValues: AreaType.values)
    ..aOS(5, _omitFieldNames ? '' : 'geometry')
    ..aOM<$1.Struct>(6, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAreaRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAreaRequest copyWith(void Function(UpdateAreaRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateAreaRequest))
          as UpdateAreaRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UpdateAreaRequest() / UpdateAreaRequest.new instead')
  static UpdateAreaRequest create() => UpdateAreaRequest._();
  static $pb.GeneratedMessage $_createMessage() => UpdateAreaRequest._();
  @$core.override
  UpdateAreaRequest createEmptyInstance() => UpdateAreaRequest._();
  @$core.pragma('dart2js:noInline')
  static UpdateAreaRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateAreaRequest>(
          UpdateAreaRequest.$_createMessage);
  static UpdateAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => $_clearField(3);

  @$pb.TagNumber(4)
  AreaType get areaType => $_getN(3);
  @$pb.TagNumber(4)
  set areaType(AreaType value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasAreaType() => $_has(3);
  @$pb.TagNumber(4)
  void clearAreaType() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get geometry => $_getSZ(4);
  @$pb.TagNumber(5)
  set geometry($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasGeometry() => $_has(4);
  @$pb.TagNumber(5)
  void clearGeometry() => $_clearField(5);

  @$pb.TagNumber(6)
  $1.Struct get extra => $_getN(5);
  @$pb.TagNumber(6)
  set extra($1.Struct value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasExtra() => $_has(5);
  @$pb.TagNumber(6)
  void clearExtra() => $_clearField(6);
  @$pb.TagNumber(6)
  $1.Struct ensureExtra() => $_ensure(5);
}

class UpdateAreaResponse extends $pb.GeneratedMessage {
  factory UpdateAreaResponse({
    AreaObject? data,
  }) {
    final result = UpdateAreaResponse._();
    if (data != null) result.data = data;
    return result;
  }

  UpdateAreaResponse._();

  factory UpdateAreaResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateAreaResponse()..mergeFromBuffer(data, registry);
  factory UpdateAreaResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateAreaResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAreaResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: UpdateAreaResponse.$_createMessage)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AreaObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAreaResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAreaResponse copyWith(void Function(UpdateAreaResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateAreaResponse))
          as UpdateAreaResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UpdateAreaResponse() / UpdateAreaResponse.new instead')
  static UpdateAreaResponse create() => UpdateAreaResponse._();
  static $pb.GeneratedMessage $_createMessage() => UpdateAreaResponse._();
  @$core.override
  UpdateAreaResponse createEmptyInstance() => UpdateAreaResponse._();
  @$core.pragma('dart2js:noInline')
  static UpdateAreaResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAreaResponse>(
          UpdateAreaResponse.$_createMessage);
  static UpdateAreaResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class DeleteAreaRequest extends $pb.GeneratedMessage {
  factory DeleteAreaRequest({
    $core.String? id,
  }) {
    final result = DeleteAreaRequest._();
    if (id != null) result.id = id;
    return result;
  }

  DeleteAreaRequest._();

  factory DeleteAreaRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteAreaRequest()..mergeFromBuffer(data, registry);
  factory DeleteAreaRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteAreaRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAreaRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: DeleteAreaRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAreaRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAreaRequest copyWith(void Function(DeleteAreaRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteAreaRequest))
          as DeleteAreaRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use DeleteAreaRequest() / DeleteAreaRequest.new instead')
  static DeleteAreaRequest create() => DeleteAreaRequest._();
  static $pb.GeneratedMessage $_createMessage() => DeleteAreaRequest._();
  @$core.override
  DeleteAreaRequest createEmptyInstance() => DeleteAreaRequest._();
  @$core.pragma('dart2js:noInline')
  static DeleteAreaRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteAreaRequest>(
          DeleteAreaRequest.$_createMessage);
  static DeleteAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class SearchAreasRequest extends $pb.GeneratedMessage {
  factory SearchAreasRequest({
    $core.String? query,
    $core.String? ownerId,
    $core.int? limit,
  }) {
    final result = SearchAreasRequest._();
    if (query != null) result.query = query;
    if (ownerId != null) result.ownerId = ownerId;
    if (limit != null) result.limit = limit;
    return result;
  }

  SearchAreasRequest._();

  factory SearchAreasRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchAreasRequest()..mergeFromBuffer(data, registry);
  factory SearchAreasRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchAreasRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchAreasRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: SearchAreasRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..aOS(2, _omitFieldNames ? '' : 'ownerId')
    ..aI(3, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchAreasRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchAreasRequest copyWith(void Function(SearchAreasRequest) updates) =>
      super.copyWith((message) => updates(message as SearchAreasRequest))
          as SearchAreasRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use SearchAreasRequest() / SearchAreasRequest.new instead')
  static SearchAreasRequest create() => SearchAreasRequest._();
  static $pb.GeneratedMessage $_createMessage() => SearchAreasRequest._();
  @$core.override
  SearchAreasRequest createEmptyInstance() => SearchAreasRequest._();
  @$core.pragma('dart2js:noInline')
  static SearchAreasRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchAreasRequest>(
          SearchAreasRequest.$_createMessage);
  static SearchAreasRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get ownerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ownerId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOwnerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearOwnerId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(3)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => $_clearField(3);
}

class SearchAreasResponse extends $pb.GeneratedMessage {
  factory SearchAreasResponse({
    $core.Iterable<AreaObject>? data,
  }) {
    final result = SearchAreasResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  SearchAreasResponse._();

  factory SearchAreasResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchAreasResponse()..mergeFromBuffer(data, registry);
  factory SearchAreasResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchAreasResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchAreasResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: SearchAreasResponse.$_createMessage)
    ..pPM<AreaObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AreaObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchAreasResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchAreasResponse copyWith(void Function(SearchAreasResponse) updates) =>
      super.copyWith((message) => updates(message as SearchAreasResponse))
          as SearchAreasResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use SearchAreasResponse() / SearchAreasResponse.new instead')
  static SearchAreasResponse create() => SearchAreasResponse._();
  static $pb.GeneratedMessage $_createMessage() => SearchAreasResponse._();
  @$core.override
  SearchAreasResponse createEmptyInstance() => SearchAreasResponse._();
  @$core.pragma('dart2js:noInline')
  static SearchAreasResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchAreasResponse>(
          SearchAreasResponse.$_createMessage);
  static SearchAreasResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<AreaObject> get data => $_getList(0);
}

class CreateRouteRequest extends $pb.GeneratedMessage {
  factory CreateRouteRequest({
    RouteObject? data,
  }) {
    final result = CreateRouteRequest._();
    if (data != null) result.data = data;
    return result;
  }

  CreateRouteRequest._();

  factory CreateRouteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateRouteRequest()..mergeFromBuffer(data, registry);
  factory CreateRouteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateRouteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRouteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: CreateRouteRequest.$_createMessage)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRouteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRouteRequest copyWith(void Function(CreateRouteRequest) updates) =>
      super.copyWith((message) => updates(message as CreateRouteRequest))
          as CreateRouteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CreateRouteRequest() / CreateRouteRequest.new instead')
  static CreateRouteRequest create() => CreateRouteRequest._();
  static $pb.GeneratedMessage $_createMessage() => CreateRouteRequest._();
  @$core.override
  CreateRouteRequest createEmptyInstance() => CreateRouteRequest._();
  @$core.pragma('dart2js:noInline')
  static CreateRouteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRouteRequest>(
          CreateRouteRequest.$_createMessage);
  static CreateRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class CreateRouteResponse extends $pb.GeneratedMessage {
  factory CreateRouteResponse({
    RouteObject? data,
  }) {
    final result = CreateRouteResponse._();
    if (data != null) result.data = data;
    return result;
  }

  CreateRouteResponse._();

  factory CreateRouteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateRouteResponse()..mergeFromBuffer(data, registry);
  factory CreateRouteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateRouteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRouteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: CreateRouteResponse.$_createMessage)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRouteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRouteResponse copyWith(void Function(CreateRouteResponse) updates) =>
      super.copyWith((message) => updates(message as CreateRouteResponse))
          as CreateRouteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use CreateRouteResponse() / CreateRouteResponse.new instead')
  static CreateRouteResponse create() => CreateRouteResponse._();
  static $pb.GeneratedMessage $_createMessage() => CreateRouteResponse._();
  @$core.override
  CreateRouteResponse createEmptyInstance() => CreateRouteResponse._();
  @$core.pragma('dart2js:noInline')
  static CreateRouteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRouteResponse>(
          CreateRouteResponse.$_createMessage);
  static CreateRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class GetRouteRequest extends $pb.GeneratedMessage {
  factory GetRouteRequest({
    $core.String? id,
  }) {
    final result = GetRouteRequest._();
    if (id != null) result.id = id;
    return result;
  }

  GetRouteRequest._();

  factory GetRouteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetRouteRequest()..mergeFromBuffer(data, registry);
  factory GetRouteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetRouteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRouteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetRouteRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRouteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRouteRequest copyWith(void Function(GetRouteRequest) updates) =>
      super.copyWith((message) => updates(message as GetRouteRequest))
          as GetRouteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetRouteRequest() / GetRouteRequest.new instead')
  static GetRouteRequest create() => GetRouteRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetRouteRequest._();
  @$core.override
  GetRouteRequest createEmptyInstance() => GetRouteRequest._();
  @$core.pragma('dart2js:noInline')
  static GetRouteRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRouteRequest>(
          GetRouteRequest.$_createMessage);
  static GetRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class GetRouteResponse extends $pb.GeneratedMessage {
  factory GetRouteResponse({
    RouteObject? data,
  }) {
    final result = GetRouteResponse._();
    if (data != null) result.data = data;
    return result;
  }

  GetRouteResponse._();

  factory GetRouteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetRouteResponse()..mergeFromBuffer(data, registry);
  factory GetRouteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetRouteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRouteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetRouteResponse.$_createMessage)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRouteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRouteResponse copyWith(void Function(GetRouteResponse) updates) =>
      super.copyWith((message) => updates(message as GetRouteResponse))
          as GetRouteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetRouteResponse() / GetRouteResponse.new instead')
  static GetRouteResponse create() => GetRouteResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetRouteResponse._();
  @$core.override
  GetRouteResponse createEmptyInstance() => GetRouteResponse._();
  @$core.pragma('dart2js:noInline')
  static GetRouteResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRouteResponse>(
          GetRouteResponse.$_createMessage);
  static GetRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class UpdateRouteRequest extends $pb.GeneratedMessage {
  factory UpdateRouteRequest({
    $core.String? id,
    $core.String? name,
    $core.String? description,
    $core.String? geometry,
    $core.double? deviationThresholdM,
    $core.int? deviationConsecutiveCount,
    $core.int? deviationCooldownSec,
    $1.Struct? extra,
  }) {
    final result = UpdateRouteRequest._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (geometry != null) result.geometry = geometry;
    if (deviationThresholdM != null)
      result.deviationThresholdM = deviationThresholdM;
    if (deviationConsecutiveCount != null)
      result.deviationConsecutiveCount = deviationConsecutiveCount;
    if (deviationCooldownSec != null)
      result.deviationCooldownSec = deviationCooldownSec;
    if (extra != null) result.extra = extra;
    return result;
  }

  UpdateRouteRequest._();

  factory UpdateRouteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRouteRequest()..mergeFromBuffer(data, registry);
  factory UpdateRouteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRouteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRouteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: UpdateRouteRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..aOS(4, _omitFieldNames ? '' : 'geometry')
    ..aD(5, _omitFieldNames ? '' : 'deviationThresholdM')
    ..aI(6, _omitFieldNames ? '' : 'deviationConsecutiveCount')
    ..aI(7, _omitFieldNames ? '' : 'deviationCooldownSec')
    ..aOM<$1.Struct>(8, _omitFieldNames ? '' : 'extra',
        subBuilder: $1.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRouteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRouteRequest copyWith(void Function(UpdateRouteRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateRouteRequest))
          as UpdateRouteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UpdateRouteRequest() / UpdateRouteRequest.new instead')
  static UpdateRouteRequest create() => UpdateRouteRequest._();
  static $pb.GeneratedMessage $_createMessage() => UpdateRouteRequest._();
  @$core.override
  UpdateRouteRequest createEmptyInstance() => UpdateRouteRequest._();
  @$core.pragma('dart2js:noInline')
  static UpdateRouteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRouteRequest>(
          UpdateRouteRequest.$_createMessage);
  static UpdateRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get geometry => $_getSZ(3);
  @$pb.TagNumber(4)
  set geometry($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasGeometry() => $_has(3);
  @$pb.TagNumber(4)
  void clearGeometry() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.double get deviationThresholdM => $_getN(4);
  @$pb.TagNumber(5)
  set deviationThresholdM($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDeviationThresholdM() => $_has(4);
  @$pb.TagNumber(5)
  void clearDeviationThresholdM() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get deviationConsecutiveCount => $_getIZ(5);
  @$pb.TagNumber(6)
  set deviationConsecutiveCount($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasDeviationConsecutiveCount() => $_has(5);
  @$pb.TagNumber(6)
  void clearDeviationConsecutiveCount() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.int get deviationCooldownSec => $_getIZ(6);
  @$pb.TagNumber(7)
  set deviationCooldownSec($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasDeviationCooldownSec() => $_has(6);
  @$pb.TagNumber(7)
  void clearDeviationCooldownSec() => $_clearField(7);

  @$pb.TagNumber(8)
  $1.Struct get extra => $_getN(7);
  @$pb.TagNumber(8)
  set extra($1.Struct value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasExtra() => $_has(7);
  @$pb.TagNumber(8)
  void clearExtra() => $_clearField(8);
  @$pb.TagNumber(8)
  $1.Struct ensureExtra() => $_ensure(7);
}

class UpdateRouteResponse extends $pb.GeneratedMessage {
  factory UpdateRouteResponse({
    RouteObject? data,
  }) {
    final result = UpdateRouteResponse._();
    if (data != null) result.data = data;
    return result;
  }

  UpdateRouteResponse._();

  factory UpdateRouteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRouteResponse()..mergeFromBuffer(data, registry);
  factory UpdateRouteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRouteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRouteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: UpdateRouteResponse.$_createMessage)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRouteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRouteResponse copyWith(void Function(UpdateRouteResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateRouteResponse))
          as UpdateRouteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use UpdateRouteResponse() / UpdateRouteResponse.new instead')
  static UpdateRouteResponse create() => UpdateRouteResponse._();
  static $pb.GeneratedMessage $_createMessage() => UpdateRouteResponse._();
  @$core.override
  UpdateRouteResponse createEmptyInstance() => UpdateRouteResponse._();
  @$core.pragma('dart2js:noInline')
  static UpdateRouteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRouteResponse>(
          UpdateRouteResponse.$_createMessage);
  static UpdateRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class DeleteRouteRequest extends $pb.GeneratedMessage {
  factory DeleteRouteRequest({
    $core.String? id,
  }) {
    final result = DeleteRouteRequest._();
    if (id != null) result.id = id;
    return result;
  }

  DeleteRouteRequest._();

  factory DeleteRouteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteRouteRequest()..mergeFromBuffer(data, registry);
  factory DeleteRouteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteRouteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRouteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: DeleteRouteRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRouteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRouteRequest copyWith(void Function(DeleteRouteRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteRouteRequest))
          as DeleteRouteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use DeleteRouteRequest() / DeleteRouteRequest.new instead')
  static DeleteRouteRequest create() => DeleteRouteRequest._();
  static $pb.GeneratedMessage $_createMessage() => DeleteRouteRequest._();
  @$core.override
  DeleteRouteRequest createEmptyInstance() => DeleteRouteRequest._();
  @$core.pragma('dart2js:noInline')
  static DeleteRouteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRouteRequest>(
          DeleteRouteRequest.$_createMessage);
  static DeleteRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class SearchRoutesRequest extends $pb.GeneratedMessage {
  factory SearchRoutesRequest({
    $core.String? ownerId,
    $core.int? limit,
  }) {
    final result = SearchRoutesRequest._();
    if (ownerId != null) result.ownerId = ownerId;
    if (limit != null) result.limit = limit;
    return result;
  }

  SearchRoutesRequest._();

  factory SearchRoutesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRoutesRequest()..mergeFromBuffer(data, registry);
  factory SearchRoutesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRoutesRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchRoutesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: SearchRoutesRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'ownerId')
    ..aI(2, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRoutesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRoutesRequest copyWith(void Function(SearchRoutesRequest) updates) =>
      super.copyWith((message) => updates(message as SearchRoutesRequest))
          as SearchRoutesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use SearchRoutesRequest() / SearchRoutesRequest.new instead')
  static SearchRoutesRequest create() => SearchRoutesRequest._();
  static $pb.GeneratedMessage $_createMessage() => SearchRoutesRequest._();
  @$core.override
  SearchRoutesRequest createEmptyInstance() => SearchRoutesRequest._();
  @$core.pragma('dart2js:noInline')
  static SearchRoutesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchRoutesRequest>(
          SearchRoutesRequest.$_createMessage);
  static SearchRoutesRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get ownerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set ownerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasOwnerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearOwnerId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(2)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimit() => $_clearField(2);
}

class SearchRoutesResponse extends $pb.GeneratedMessage {
  factory SearchRoutesResponse({
    $core.Iterable<RouteObject>? data,
  }) {
    final result = SearchRoutesResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  SearchRoutesResponse._();

  factory SearchRoutesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRoutesResponse()..mergeFromBuffer(data, registry);
  factory SearchRoutesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRoutesResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchRoutesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: SearchRoutesResponse.$_createMessage)
    ..pPM<RouteObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRoutesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRoutesResponse copyWith(void Function(SearchRoutesResponse) updates) =>
      super.copyWith((message) => updates(message as SearchRoutesResponse))
          as SearchRoutesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SearchRoutesResponse() / SearchRoutesResponse.new instead')
  static SearchRoutesResponse create() => SearchRoutesResponse._();
  static $pb.GeneratedMessage $_createMessage() => SearchRoutesResponse._();
  @$core.override
  SearchRoutesResponse createEmptyInstance() => SearchRoutesResponse._();
  @$core.pragma('dart2js:noInline')
  static SearchRoutesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchRoutesResponse>(
          SearchRoutesResponse.$_createMessage);
  static SearchRoutesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<RouteObject> get data => $_getList(0);
}

class AssignRouteRequest extends $pb.GeneratedMessage {
  factory AssignRouteRequest({
    $core.String? subjectId,
    $core.String? routeId,
    $0.Timestamp? validFrom,
    $0.Timestamp? validUntil,
  }) {
    final result = AssignRouteRequest._();
    if (subjectId != null) result.subjectId = subjectId;
    if (routeId != null) result.routeId = routeId;
    if (validFrom != null) result.validFrom = validFrom;
    if (validUntil != null) result.validUntil = validUntil;
    return result;
  }

  AssignRouteRequest._();

  factory AssignRouteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AssignRouteRequest()..mergeFromBuffer(data, registry);
  factory AssignRouteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AssignRouteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssignRouteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: AssignRouteRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOS(2, _omitFieldNames ? '' : 'routeId')
    ..aOM<$0.Timestamp>(3, _omitFieldNames ? '' : 'validFrom',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aOM<$0.Timestamp>(4, _omitFieldNames ? '' : 'validUntil',
        subBuilder: $0.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssignRouteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssignRouteRequest copyWith(void Function(AssignRouteRequest) updates) =>
      super.copyWith((message) => updates(message as AssignRouteRequest))
          as AssignRouteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AssignRouteRequest() / AssignRouteRequest.new instead')
  static AssignRouteRequest create() => AssignRouteRequest._();
  static $pb.GeneratedMessage $_createMessage() => AssignRouteRequest._();
  @$core.override
  AssignRouteRequest createEmptyInstance() => AssignRouteRequest._();
  @$core.pragma('dart2js:noInline')
  static AssignRouteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssignRouteRequest>(
          AssignRouteRequest.$_createMessage);
  static AssignRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get routeId => $_getSZ(1);
  @$pb.TagNumber(2)
  set routeId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRouteId() => $_has(1);
  @$pb.TagNumber(2)
  void clearRouteId() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Timestamp get validFrom => $_getN(2);
  @$pb.TagNumber(3)
  set validFrom($0.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasValidFrom() => $_has(2);
  @$pb.TagNumber(3)
  void clearValidFrom() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Timestamp ensureValidFrom() => $_ensure(2);

  @$pb.TagNumber(4)
  $0.Timestamp get validUntil => $_getN(3);
  @$pb.TagNumber(4)
  set validUntil($0.Timestamp value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasValidUntil() => $_has(3);
  @$pb.TagNumber(4)
  void clearValidUntil() => $_clearField(4);
  @$pb.TagNumber(4)
  $0.Timestamp ensureValidUntil() => $_ensure(3);
}

class AssignRouteResponse extends $pb.GeneratedMessage {
  factory AssignRouteResponse({
    RouteAssignmentObject? data,
  }) {
    final result = AssignRouteResponse._();
    if (data != null) result.data = data;
    return result;
  }

  AssignRouteResponse._();

  factory AssignRouteResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AssignRouteResponse()..mergeFromBuffer(data, registry);
  factory AssignRouteResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AssignRouteResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssignRouteResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: AssignRouteResponse.$_createMessage)
    ..aOM<RouteAssignmentObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteAssignmentObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssignRouteResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssignRouteResponse copyWith(void Function(AssignRouteResponse) updates) =>
      super.copyWith((message) => updates(message as AssignRouteResponse))
          as AssignRouteResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use AssignRouteResponse() / AssignRouteResponse.new instead')
  static AssignRouteResponse create() => AssignRouteResponse._();
  static $pb.GeneratedMessage $_createMessage() => AssignRouteResponse._();
  @$core.override
  AssignRouteResponse createEmptyInstance() => AssignRouteResponse._();
  @$core.pragma('dart2js:noInline')
  static AssignRouteResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssignRouteResponse>(
          AssignRouteResponse.$_createMessage);
  static AssignRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteAssignmentObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteAssignmentObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RouteAssignmentObject ensureData() => $_ensure(0);
}

class UnassignRouteRequest extends $pb.GeneratedMessage {
  factory UnassignRouteRequest({
    $core.String? id,
  }) {
    final result = UnassignRouteRequest._();
    if (id != null) result.id = id;
    return result;
  }

  UnassignRouteRequest._();

  factory UnassignRouteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UnassignRouteRequest()..mergeFromBuffer(data, registry);
  factory UnassignRouteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UnassignRouteRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnassignRouteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: UnassignRouteRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnassignRouteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnassignRouteRequest copyWith(void Function(UnassignRouteRequest) updates) =>
      super.copyWith((message) => updates(message as UnassignRouteRequest))
          as UnassignRouteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use UnassignRouteRequest() / UnassignRouteRequest.new instead')
  static UnassignRouteRequest create() => UnassignRouteRequest._();
  static $pb.GeneratedMessage $_createMessage() => UnassignRouteRequest._();
  @$core.override
  UnassignRouteRequest createEmptyInstance() => UnassignRouteRequest._();
  @$core.pragma('dart2js:noInline')
  static UnassignRouteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnassignRouteRequest>(
          UnassignRouteRequest.$_createMessage);
  static UnassignRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class GetSubjectRouteAssignmentsRequest extends $pb.GeneratedMessage {
  factory GetSubjectRouteAssignmentsRequest({
    $core.String? subjectId,
  }) {
    final result = GetSubjectRouteAssignmentsRequest._();
    if (subjectId != null) result.subjectId = subjectId;
    return result;
  }

  GetSubjectRouteAssignmentsRequest._();

  factory GetSubjectRouteAssignmentsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectRouteAssignmentsRequest()..mergeFromBuffer(data, registry);
  factory GetSubjectRouteAssignmentsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectRouteAssignmentsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSubjectRouteAssignmentsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetSubjectRouteAssignmentsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectRouteAssignmentsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectRouteAssignmentsRequest copyWith(
          void Function(GetSubjectRouteAssignmentsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetSubjectRouteAssignmentsRequest))
          as GetSubjectRouteAssignmentsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetSubjectRouteAssignmentsRequest() / GetSubjectRouteAssignmentsRequest.new instead')
  static GetSubjectRouteAssignmentsRequest create() =>
      GetSubjectRouteAssignmentsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetSubjectRouteAssignmentsRequest._();
  @$core.override
  GetSubjectRouteAssignmentsRequest createEmptyInstance() =>
      GetSubjectRouteAssignmentsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetSubjectRouteAssignmentsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSubjectRouteAssignmentsRequest>(
          GetSubjectRouteAssignmentsRequest.$_createMessage);
  static GetSubjectRouteAssignmentsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);
}

class GetSubjectRouteAssignmentsResponse extends $pb.GeneratedMessage {
  factory GetSubjectRouteAssignmentsResponse({
    $core.Iterable<RouteAssignmentObject>? data,
  }) {
    final result = GetSubjectRouteAssignmentsResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  GetSubjectRouteAssignmentsResponse._();

  factory GetSubjectRouteAssignmentsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectRouteAssignmentsResponse()..mergeFromBuffer(data, registry);
  factory GetSubjectRouteAssignmentsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectRouteAssignmentsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSubjectRouteAssignmentsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetSubjectRouteAssignmentsResponse.$_createMessage)
    ..pPM<RouteAssignmentObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RouteAssignmentObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectRouteAssignmentsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectRouteAssignmentsResponse copyWith(
          void Function(GetSubjectRouteAssignmentsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetSubjectRouteAssignmentsResponse))
          as GetSubjectRouteAssignmentsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetSubjectRouteAssignmentsResponse() / GetSubjectRouteAssignmentsResponse.new instead')
  static GetSubjectRouteAssignmentsResponse create() =>
      GetSubjectRouteAssignmentsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetSubjectRouteAssignmentsResponse._();
  @$core.override
  GetSubjectRouteAssignmentsResponse createEmptyInstance() =>
      GetSubjectRouteAssignmentsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetSubjectRouteAssignmentsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSubjectRouteAssignmentsResponse>(
          GetSubjectRouteAssignmentsResponse.$_createMessage);
  static GetSubjectRouteAssignmentsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<RouteAssignmentObject> get data => $_getList(0);
}

class GetTrackRequest extends $pb.GeneratedMessage {
  factory GetTrackRequest({
    $core.String? subjectId,
    $0.Timestamp? from,
    $0.Timestamp? to,
    $core.int? limit,
    $core.int? offset,
  }) {
    final result = GetTrackRequest._();
    if (subjectId != null) result.subjectId = subjectId;
    if (from != null) result.from = from;
    if (to != null) result.to = to;
    if (limit != null) result.limit = limit;
    if (offset != null) result.offset = offset;
    return result;
  }

  GetTrackRequest._();

  factory GetTrackRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTrackRequest()..mergeFromBuffer(data, registry);
  factory GetTrackRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTrackRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrackRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetTrackRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$0.Timestamp>(2, _omitFieldNames ? '' : 'from',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aOM<$0.Timestamp>(3, _omitFieldNames ? '' : 'to',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aI(4, _omitFieldNames ? '' : 'limit')
    ..aI(5, _omitFieldNames ? '' : 'offset')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrackRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrackRequest copyWith(void Function(GetTrackRequest) updates) =>
      super.copyWith((message) => updates(message as GetTrackRequest))
          as GetTrackRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetTrackRequest() / GetTrackRequest.new instead')
  static GetTrackRequest create() => GetTrackRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetTrackRequest._();
  @$core.override
  GetTrackRequest createEmptyInstance() => GetTrackRequest._();
  @$core.pragma('dart2js:noInline')
  static GetTrackRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTrackRequest>(
          GetTrackRequest.$_createMessage);
  static GetTrackRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $0.Timestamp get from => $_getN(1);
  @$pb.TagNumber(2)
  set from($0.Timestamp value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasFrom() => $_has(1);
  @$pb.TagNumber(2)
  void clearFrom() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.Timestamp ensureFrom() => $_ensure(1);

  @$pb.TagNumber(3)
  $0.Timestamp get to => $_getN(2);
  @$pb.TagNumber(3)
  set to($0.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasTo() => $_has(2);
  @$pb.TagNumber(3)
  void clearTo() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Timestamp ensureTo() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(4)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(4)
  void clearLimit() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get offset => $_getIZ(4);
  @$pb.TagNumber(5)
  set offset($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(5)
  void clearOffset() => $_clearField(5);
}

class GetTrackResponse extends $pb.GeneratedMessage {
  factory GetTrackResponse({
    $core.Iterable<LocationPointObject>? data,
  }) {
    final result = GetTrackResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  GetTrackResponse._();

  factory GetTrackResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTrackResponse()..mergeFromBuffer(data, registry);
  factory GetTrackResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetTrackResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrackResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetTrackResponse.$_createMessage)
    ..pPM<LocationPointObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: LocationPointObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrackResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrackResponse copyWith(void Function(GetTrackResponse) updates) =>
      super.copyWith((message) => updates(message as GetTrackResponse))
          as GetTrackResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetTrackResponse() / GetTrackResponse.new instead')
  static GetTrackResponse create() => GetTrackResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetTrackResponse._();
  @$core.override
  GetTrackResponse createEmptyInstance() => GetTrackResponse._();
  @$core.pragma('dart2js:noInline')
  static GetTrackResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTrackResponse>(
          GetTrackResponse.$_createMessage);
  static GetTrackResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<LocationPointObject> get data => $_getList(0);
}

class GetSubjectEventsRequest extends $pb.GeneratedMessage {
  factory GetSubjectEventsRequest({
    $core.String? subjectId,
    $0.Timestamp? from,
    $0.Timestamp? to,
    $core.int? limit,
    $core.int? offset,
  }) {
    final result = GetSubjectEventsRequest._();
    if (subjectId != null) result.subjectId = subjectId;
    if (from != null) result.from = from;
    if (to != null) result.to = to;
    if (limit != null) result.limit = limit;
    if (offset != null) result.offset = offset;
    return result;
  }

  GetSubjectEventsRequest._();

  factory GetSubjectEventsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectEventsRequest()..mergeFromBuffer(data, registry);
  factory GetSubjectEventsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectEventsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSubjectEventsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetSubjectEventsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$0.Timestamp>(2, _omitFieldNames ? '' : 'from',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aOM<$0.Timestamp>(3, _omitFieldNames ? '' : 'to',
        subBuilder: $0.Timestamp.$_createMessage)
    ..aI(4, _omitFieldNames ? '' : 'limit')
    ..aI(5, _omitFieldNames ? '' : 'offset')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectEventsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectEventsRequest copyWith(
          void Function(GetSubjectEventsRequest) updates) =>
      super.copyWith((message) => updates(message as GetSubjectEventsRequest))
          as GetSubjectEventsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetSubjectEventsRequest() / GetSubjectEventsRequest.new instead')
  static GetSubjectEventsRequest create() => GetSubjectEventsRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetSubjectEventsRequest._();
  @$core.override
  GetSubjectEventsRequest createEmptyInstance() => GetSubjectEventsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetSubjectEventsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSubjectEventsRequest>(
          GetSubjectEventsRequest.$_createMessage);
  static GetSubjectEventsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $0.Timestamp get from => $_getN(1);
  @$pb.TagNumber(2)
  set from($0.Timestamp value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasFrom() => $_has(1);
  @$pb.TagNumber(2)
  void clearFrom() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.Timestamp ensureFrom() => $_ensure(1);

  @$pb.TagNumber(3)
  $0.Timestamp get to => $_getN(2);
  @$pb.TagNumber(3)
  set to($0.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasTo() => $_has(2);
  @$pb.TagNumber(3)
  void clearTo() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Timestamp ensureTo() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(4)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(4)
  void clearLimit() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get offset => $_getIZ(4);
  @$pb.TagNumber(5)
  set offset($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(5)
  void clearOffset() => $_clearField(5);
}

class GetSubjectEventsResponse extends $pb.GeneratedMessage {
  factory GetSubjectEventsResponse({
    $core.Iterable<GeoEventObject>? data,
  }) {
    final result = GetSubjectEventsResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  GetSubjectEventsResponse._();

  factory GetSubjectEventsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectEventsResponse()..mergeFromBuffer(data, registry);
  factory GetSubjectEventsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSubjectEventsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSubjectEventsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetSubjectEventsResponse.$_createMessage)
    ..pPM<GeoEventObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: GeoEventObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectEventsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubjectEventsResponse copyWith(
          void Function(GetSubjectEventsResponse) updates) =>
      super.copyWith((message) => updates(message as GetSubjectEventsResponse))
          as GetSubjectEventsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetSubjectEventsResponse() / GetSubjectEventsResponse.new instead')
  static GetSubjectEventsResponse create() => GetSubjectEventsResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetSubjectEventsResponse._();
  @$core.override
  GetSubjectEventsResponse createEmptyInstance() =>
      GetSubjectEventsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetSubjectEventsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSubjectEventsResponse>(
          GetSubjectEventsResponse.$_createMessage);
  static GetSubjectEventsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<GeoEventObject> get data => $_getList(0);
}

class GetAreaSubjectsRequest extends $pb.GeneratedMessage {
  factory GetAreaSubjectsRequest({
    $core.String? areaId,
  }) {
    final result = GetAreaSubjectsRequest._();
    if (areaId != null) result.areaId = areaId;
    return result;
  }

  GetAreaSubjectsRequest._();

  factory GetAreaSubjectsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaSubjectsRequest()..mergeFromBuffer(data, registry);
  factory GetAreaSubjectsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaSubjectsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAreaSubjectsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetAreaSubjectsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'areaId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaSubjectsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaSubjectsRequest copyWith(
          void Function(GetAreaSubjectsRequest) updates) =>
      super.copyWith((message) => updates(message as GetAreaSubjectsRequest))
          as GetAreaSubjectsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetAreaSubjectsRequest() / GetAreaSubjectsRequest.new instead')
  static GetAreaSubjectsRequest create() => GetAreaSubjectsRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetAreaSubjectsRequest._();
  @$core.override
  GetAreaSubjectsRequest createEmptyInstance() => GetAreaSubjectsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetAreaSubjectsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAreaSubjectsRequest>(
          GetAreaSubjectsRequest.$_createMessage);
  static GetAreaSubjectsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get areaId => $_getSZ(0);
  @$pb.TagNumber(1)
  set areaId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAreaId() => $_has(0);
  @$pb.TagNumber(1)
  void clearAreaId() => $_clearField(1);
}

class GetAreaSubjectsResponse extends $pb.GeneratedMessage {
  factory GetAreaSubjectsResponse({
    $core.Iterable<AreaSubjectObject>? data,
  }) {
    final result = GetAreaSubjectsResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  GetAreaSubjectsResponse._();

  factory GetAreaSubjectsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaSubjectsResponse()..mergeFromBuffer(data, registry);
  factory GetAreaSubjectsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetAreaSubjectsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAreaSubjectsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetAreaSubjectsResponse.$_createMessage)
    ..pPM<AreaSubjectObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AreaSubjectObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaSubjectsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAreaSubjectsResponse copyWith(
          void Function(GetAreaSubjectsResponse) updates) =>
      super.copyWith((message) => updates(message as GetAreaSubjectsResponse))
          as GetAreaSubjectsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetAreaSubjectsResponse() / GetAreaSubjectsResponse.new instead')
  static GetAreaSubjectsResponse create() => GetAreaSubjectsResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetAreaSubjectsResponse._();
  @$core.override
  GetAreaSubjectsResponse createEmptyInstance() => GetAreaSubjectsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetAreaSubjectsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAreaSubjectsResponse>(
          GetAreaSubjectsResponse.$_createMessage);
  static GetAreaSubjectsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<AreaSubjectObject> get data => $_getList(0);
}

class GetNearbySubjectsRequest extends $pb.GeneratedMessage {
  factory GetNearbySubjectsRequest({
    $core.String? subjectId,
    $core.double? radiusMeters,
    $core.int? limit,
  }) {
    final result = GetNearbySubjectsRequest._();
    if (subjectId != null) result.subjectId = subjectId;
    if (radiusMeters != null) result.radiusMeters = radiusMeters;
    if (limit != null) result.limit = limit;
    return result;
  }

  GetNearbySubjectsRequest._();

  factory GetNearbySubjectsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbySubjectsRequest()..mergeFromBuffer(data, registry);
  factory GetNearbySubjectsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbySubjectsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNearbySubjectsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetNearbySubjectsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aD(2, _omitFieldNames ? '' : 'radiusMeters')
    ..aI(3, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbySubjectsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbySubjectsRequest copyWith(
          void Function(GetNearbySubjectsRequest) updates) =>
      super.copyWith((message) => updates(message as GetNearbySubjectsRequest))
          as GetNearbySubjectsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetNearbySubjectsRequest() / GetNearbySubjectsRequest.new instead')
  static GetNearbySubjectsRequest create() => GetNearbySubjectsRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetNearbySubjectsRequest._();
  @$core.override
  GetNearbySubjectsRequest createEmptyInstance() =>
      GetNearbySubjectsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetNearbySubjectsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNearbySubjectsRequest>(
          GetNearbySubjectsRequest.$_createMessage);
  static GetNearbySubjectsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get radiusMeters => $_getN(1);
  @$pb.TagNumber(2)
  set radiusMeters($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRadiusMeters() => $_has(1);
  @$pb.TagNumber(2)
  void clearRadiusMeters() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(3)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => $_clearField(3);
}

class GetNearbySubjectsResponse extends $pb.GeneratedMessage {
  factory GetNearbySubjectsResponse({
    $core.Iterable<NearbySubjectObject>? data,
  }) {
    final result = GetNearbySubjectsResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  GetNearbySubjectsResponse._();

  factory GetNearbySubjectsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbySubjectsResponse()..mergeFromBuffer(data, registry);
  factory GetNearbySubjectsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbySubjectsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNearbySubjectsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetNearbySubjectsResponse.$_createMessage)
    ..pPM<NearbySubjectObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: NearbySubjectObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbySubjectsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbySubjectsResponse copyWith(
          void Function(GetNearbySubjectsResponse) updates) =>
      super.copyWith((message) => updates(message as GetNearbySubjectsResponse))
          as GetNearbySubjectsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetNearbySubjectsResponse() / GetNearbySubjectsResponse.new instead')
  static GetNearbySubjectsResponse create() => GetNearbySubjectsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetNearbySubjectsResponse._();
  @$core.override
  GetNearbySubjectsResponse createEmptyInstance() =>
      GetNearbySubjectsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetNearbySubjectsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNearbySubjectsResponse>(
          GetNearbySubjectsResponse.$_createMessage);
  static GetNearbySubjectsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<NearbySubjectObject> get data => $_getList(0);
}

class GetNearbyAreasRequest extends $pb.GeneratedMessage {
  factory GetNearbyAreasRequest({
    $core.double? latitude,
    $core.double? longitude,
    $core.double? radiusMeters,
    $core.int? limit,
  }) {
    final result = GetNearbyAreasRequest._();
    if (latitude != null) result.latitude = latitude;
    if (longitude != null) result.longitude = longitude;
    if (radiusMeters != null) result.radiusMeters = radiusMeters;
    if (limit != null) result.limit = limit;
    return result;
  }

  GetNearbyAreasRequest._();

  factory GetNearbyAreasRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbyAreasRequest()..mergeFromBuffer(data, registry);
  factory GetNearbyAreasRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbyAreasRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNearbyAreasRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetNearbyAreasRequest.$_createMessage)
    ..aD(1, _omitFieldNames ? '' : 'latitude')
    ..aD(2, _omitFieldNames ? '' : 'longitude')
    ..aD(3, _omitFieldNames ? '' : 'radiusMeters')
    ..aI(4, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbyAreasRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbyAreasRequest copyWith(
          void Function(GetNearbyAreasRequest) updates) =>
      super.copyWith((message) => updates(message as GetNearbyAreasRequest))
          as GetNearbyAreasRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetNearbyAreasRequest() / GetNearbyAreasRequest.new instead')
  static GetNearbyAreasRequest create() => GetNearbyAreasRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetNearbyAreasRequest._();
  @$core.override
  GetNearbyAreasRequest createEmptyInstance() => GetNearbyAreasRequest._();
  @$core.pragma('dart2js:noInline')
  static GetNearbyAreasRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNearbyAreasRequest>(
          GetNearbyAreasRequest.$_createMessage);
  static GetNearbyAreasRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.double get latitude => $_getN(0);
  @$pb.TagNumber(1)
  set latitude($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLatitude() => $_has(0);
  @$pb.TagNumber(1)
  void clearLatitude() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get longitude => $_getN(1);
  @$pb.TagNumber(2)
  set longitude($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLongitude() => $_has(1);
  @$pb.TagNumber(2)
  void clearLongitude() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.double get radiusMeters => $_getN(2);
  @$pb.TagNumber(3)
  set radiusMeters($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(3)
  $core.bool hasRadiusMeters() => $_has(2);
  @$pb.TagNumber(3)
  void clearRadiusMeters() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(4)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(4)
  void clearLimit() => $_clearField(4);
}

class GetNearbyAreasResponse extends $pb.GeneratedMessage {
  factory GetNearbyAreasResponse({
    $core.Iterable<NearbyAreaObject>? data,
  }) {
    final result = GetNearbyAreasResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  GetNearbyAreasResponse._();

  factory GetNearbyAreasResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbyAreasResponse()..mergeFromBuffer(data, registry);
  factory GetNearbyAreasResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetNearbyAreasResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNearbyAreasResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'),
      createEmptyInstance: GetNearbyAreasResponse.$_createMessage)
    ..pPM<NearbyAreaObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: NearbyAreaObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbyAreasResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNearbyAreasResponse copyWith(
          void Function(GetNearbyAreasResponse) updates) =>
      super.copyWith((message) => updates(message as GetNearbyAreasResponse))
          as GetNearbyAreasResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetNearbyAreasResponse() / GetNearbyAreasResponse.new instead')
  static GetNearbyAreasResponse create() => GetNearbyAreasResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetNearbyAreasResponse._();
  @$core.override
  GetNearbyAreasResponse createEmptyInstance() => GetNearbyAreasResponse._();
  @$core.pragma('dart2js:noInline')
  static GetNearbyAreasResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNearbyAreasResponse>(
          GetNearbyAreasResponse.$_createMessage);
  static GetNearbyAreasResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<NearbyAreaObject> get data => $_getList(0);
}

class GeolocationServiceApi {
  final $pb.RpcClient _client;

  GeolocationServiceApi(this._client);

  $async.Future<IngestLocationsResponse> ingestLocations(
          $pb.ClientContext? ctx, IngestLocationsRequest request) =>
      _client.invoke<IngestLocationsResponse>(ctx, 'GeolocationService',
          'IngestLocations', request, IngestLocationsResponse());
  $async.Future<CreateAreaResponse> createArea(
          $pb.ClientContext? ctx, CreateAreaRequest request) =>
      _client.invoke<CreateAreaResponse>(ctx, 'GeolocationService',
          'CreateArea', request, CreateAreaResponse());
  $async.Future<GetAreaResponse> getArea(
          $pb.ClientContext? ctx, GetAreaRequest request) =>
      _client.invoke<GetAreaResponse>(
          ctx, 'GeolocationService', 'GetArea', request, GetAreaResponse());
  $async.Future<UpdateAreaResponse> updateArea(
          $pb.ClientContext? ctx, UpdateAreaRequest request) =>
      _client.invoke<UpdateAreaResponse>(ctx, 'GeolocationService',
          'UpdateArea', request, UpdateAreaResponse());
  $async.Future<$2.Empty> deleteArea(
          $pb.ClientContext? ctx, DeleteAreaRequest request) =>
      _client.invoke<$2.Empty>(
          ctx, 'GeolocationService', 'DeleteArea', request, $2.Empty());
  $async.Future<SearchAreasResponse> searchAreas(
          $pb.ClientContext? ctx, SearchAreasRequest request) =>
      _client.invoke<SearchAreasResponse>(ctx, 'GeolocationService',
          'SearchAreas', request, SearchAreasResponse());
  $async.Future<CreateRouteResponse> createRoute(
          $pb.ClientContext? ctx, CreateRouteRequest request) =>
      _client.invoke<CreateRouteResponse>(ctx, 'GeolocationService',
          'CreateRoute', request, CreateRouteResponse());
  $async.Future<GetRouteResponse> getRoute(
          $pb.ClientContext? ctx, GetRouteRequest request) =>
      _client.invoke<GetRouteResponse>(
          ctx, 'GeolocationService', 'GetRoute', request, GetRouteResponse());
  $async.Future<UpdateRouteResponse> updateRoute(
          $pb.ClientContext? ctx, UpdateRouteRequest request) =>
      _client.invoke<UpdateRouteResponse>(ctx, 'GeolocationService',
          'UpdateRoute', request, UpdateRouteResponse());
  $async.Future<$2.Empty> deleteRoute(
          $pb.ClientContext? ctx, DeleteRouteRequest request) =>
      _client.invoke<$2.Empty>(
          ctx, 'GeolocationService', 'DeleteRoute', request, $2.Empty());
  $async.Future<SearchRoutesResponse> searchRoutes(
          $pb.ClientContext? ctx, SearchRoutesRequest request) =>
      _client.invoke<SearchRoutesResponse>(ctx, 'GeolocationService',
          'SearchRoutes', request, SearchRoutesResponse());
  $async.Future<AssignRouteResponse> assignRoute(
          $pb.ClientContext? ctx, AssignRouteRequest request) =>
      _client.invoke<AssignRouteResponse>(ctx, 'GeolocationService',
          'AssignRoute', request, AssignRouteResponse());
  $async.Future<$2.Empty> unassignRoute(
          $pb.ClientContext? ctx, UnassignRouteRequest request) =>
      _client.invoke<$2.Empty>(
          ctx, 'GeolocationService', 'UnassignRoute', request, $2.Empty());
  $async.Future<GetSubjectRouteAssignmentsResponse> getSubjectRouteAssignments(
          $pb.ClientContext? ctx, GetSubjectRouteAssignmentsRequest request) =>
      _client.invoke<GetSubjectRouteAssignmentsResponse>(
          ctx,
          'GeolocationService',
          'GetSubjectRouteAssignments',
          request,
          GetSubjectRouteAssignmentsResponse());
  $async.Future<GetTrackResponse> getTrack(
          $pb.ClientContext? ctx, GetTrackRequest request) =>
      _client.invoke<GetTrackResponse>(
          ctx, 'GeolocationService', 'GetTrack', request, GetTrackResponse());
  $async.Future<GetSubjectEventsResponse> getSubjectEvents(
          $pb.ClientContext? ctx, GetSubjectEventsRequest request) =>
      _client.invoke<GetSubjectEventsResponse>(ctx, 'GeolocationService',
          'GetSubjectEvents', request, GetSubjectEventsResponse());
  $async.Future<GetAreaSubjectsResponse> getAreaSubjects(
          $pb.ClientContext? ctx, GetAreaSubjectsRequest request) =>
      _client.invoke<GetAreaSubjectsResponse>(ctx, 'GeolocationService',
          'GetAreaSubjects', request, GetAreaSubjectsResponse());
  $async.Future<GetNearbySubjectsResponse> getNearbySubjects(
          $pb.ClientContext? ctx, GetNearbySubjectsRequest request) =>
      _client.invoke<GetNearbySubjectsResponse>(ctx, 'GeolocationService',
          'GetNearbySubjects', request, GetNearbySubjectsResponse());
  $async.Future<GetNearbyAreasResponse> getNearbyAreas(
          $pb.ClientContext? ctx, GetNearbyAreasRequest request) =>
      _client.invoke<GetNearbyAreasResponse>(ctx, 'GeolocationService',
          'GetNearbyAreas', request, GetNearbyAreasResponse());
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
