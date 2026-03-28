//
//  Generated code. Do not modify.
//  source: geolocation/v1/geolocation.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import '../../google/protobuf/empty.pb.dart' as $8;
import '../../google/protobuf/struct.pb.dart' as $6;
import '../../google/protobuf/timestamp.pb.dart' as $2;
import 'geolocation.pbenum.dart';

export 'geolocation.pbenum.dart';

class LocationPointInput extends $pb.GeneratedMessage {
  factory LocationPointInput({
    $2.Timestamp? timestamp,
    $core.double? latitude,
    $core.double? longitude,
    $core.double? altitude,
    $core.double? accuracy,
    $core.double? speed,
    $core.double? bearing,
    LocationSource? source,
    $6.Struct? extra,
    $core.String? deviceId,
  }) {
    final $result = create();
    if (timestamp != null) {
      $result.timestamp = timestamp;
    }
    if (latitude != null) {
      $result.latitude = latitude;
    }
    if (longitude != null) {
      $result.longitude = longitude;
    }
    if (altitude != null) {
      $result.altitude = altitude;
    }
    if (accuracy != null) {
      $result.accuracy = accuracy;
    }
    if (speed != null) {
      $result.speed = speed;
    }
    if (bearing != null) {
      $result.bearing = bearing;
    }
    if (source != null) {
      $result.source = source;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    if (deviceId != null) {
      $result.deviceId = deviceId;
    }
    return $result;
  }
  LocationPointInput._() : super();
  factory LocationPointInput.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LocationPointInput.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LocationPointInput', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<$2.Timestamp>(1, _omitFieldNames ? '' : 'timestamp', subBuilder: $2.Timestamp.create)
    ..a<$core.double>(2, _omitFieldNames ? '' : 'latitude', $pb.PbFieldType.OD)
    ..a<$core.double>(3, _omitFieldNames ? '' : 'longitude', $pb.PbFieldType.OD)
    ..a<$core.double>(4, _omitFieldNames ? '' : 'altitude', $pb.PbFieldType.OD)
    ..a<$core.double>(5, _omitFieldNames ? '' : 'accuracy', $pb.PbFieldType.OD)
    ..a<$core.double>(6, _omitFieldNames ? '' : 'speed', $pb.PbFieldType.OD)
    ..a<$core.double>(7, _omitFieldNames ? '' : 'bearing', $pb.PbFieldType.OD)
    ..e<LocationSource>(8, _omitFieldNames ? '' : 'source', $pb.PbFieldType.OE, defaultOrMaker: LocationSource.LOCATION_SOURCE_UNSPECIFIED, valueOf: LocationSource.valueOf, enumValues: LocationSource.values)
    ..aOM<$6.Struct>(9, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..aOS(10, _omitFieldNames ? '' : 'deviceId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LocationPointInput clone() => LocationPointInput()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LocationPointInput copyWith(void Function(LocationPointInput) updates) => super.copyWith((message) => updates(message as LocationPointInput)) as LocationPointInput;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LocationPointInput create() => LocationPointInput._();
  LocationPointInput createEmptyInstance() => create();
  static $pb.PbList<LocationPointInput> createRepeated() => $pb.PbList<LocationPointInput>();
  @$core.pragma('dart2js:noInline')
  static LocationPointInput getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LocationPointInput>(create);
  static LocationPointInput? _defaultInstance;

  @$pb.TagNumber(1)
  $2.Timestamp get timestamp => $_getN(0);
  @$pb.TagNumber(1)
  set timestamp($2.Timestamp v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasTimestamp() => $_has(0);
  @$pb.TagNumber(1)
  void clearTimestamp() => clearField(1);
  @$pb.TagNumber(1)
  $2.Timestamp ensureTimestamp() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.double get latitude => $_getN(1);
  @$pb.TagNumber(2)
  set latitude($core.double v) { $_setDouble(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLatitude() => $_has(1);
  @$pb.TagNumber(2)
  void clearLatitude() => clearField(2);

  @$pb.TagNumber(3)
  $core.double get longitude => $_getN(2);
  @$pb.TagNumber(3)
  set longitude($core.double v) { $_setDouble(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLongitude() => $_has(2);
  @$pb.TagNumber(3)
  void clearLongitude() => clearField(3);

  @$pb.TagNumber(4)
  $core.double get altitude => $_getN(3);
  @$pb.TagNumber(4)
  set altitude($core.double v) { $_setDouble(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasAltitude() => $_has(3);
  @$pb.TagNumber(4)
  void clearAltitude() => clearField(4);

  @$pb.TagNumber(5)
  $core.double get accuracy => $_getN(4);
  @$pb.TagNumber(5)
  set accuracy($core.double v) { $_setDouble(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasAccuracy() => $_has(4);
  @$pb.TagNumber(5)
  void clearAccuracy() => clearField(5);

  @$pb.TagNumber(6)
  $core.double get speed => $_getN(5);
  @$pb.TagNumber(6)
  set speed($core.double v) { $_setDouble(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasSpeed() => $_has(5);
  @$pb.TagNumber(6)
  void clearSpeed() => clearField(6);

  @$pb.TagNumber(7)
  $core.double get bearing => $_getN(6);
  @$pb.TagNumber(7)
  set bearing($core.double v) { $_setDouble(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasBearing() => $_has(6);
  @$pb.TagNumber(7)
  void clearBearing() => clearField(7);

  @$pb.TagNumber(8)
  LocationSource get source => $_getN(7);
  @$pb.TagNumber(8)
  set source(LocationSource v) { setField(8, v); }
  @$pb.TagNumber(8)
  $core.bool hasSource() => $_has(7);
  @$pb.TagNumber(8)
  void clearSource() => clearField(8);

  @$pb.TagNumber(9)
  $6.Struct get extra => $_getN(8);
  @$pb.TagNumber(9)
  set extra($6.Struct v) { setField(9, v); }
  @$pb.TagNumber(9)
  $core.bool hasExtra() => $_has(8);
  @$pb.TagNumber(9)
  void clearExtra() => clearField(9);
  @$pb.TagNumber(9)
  $6.Struct ensureExtra() => $_ensure(8);

  @$pb.TagNumber(10)
  $core.String get deviceId => $_getSZ(9);
  @$pb.TagNumber(10)
  set deviceId($core.String v) { $_setString(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasDeviceId() => $_has(9);
  @$pb.TagNumber(10)
  void clearDeviceId() => clearField(10);
}

class LocationPointObject extends $pb.GeneratedMessage {
  factory LocationPointObject({
    $core.String? id,
    $core.String? subjectId,
    $2.Timestamp? timestamp,
    $core.double? latitude,
    $core.double? longitude,
    $core.double? altitude,
    $core.double? accuracy,
    $core.double? speed,
    $core.double? bearing,
    LocationSource? source,
    $6.Struct? extra,
    $2.Timestamp? createdAt,
    $core.String? deviceId,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (timestamp != null) {
      $result.timestamp = timestamp;
    }
    if (latitude != null) {
      $result.latitude = latitude;
    }
    if (longitude != null) {
      $result.longitude = longitude;
    }
    if (altitude != null) {
      $result.altitude = altitude;
    }
    if (accuracy != null) {
      $result.accuracy = accuracy;
    }
    if (speed != null) {
      $result.speed = speed;
    }
    if (bearing != null) {
      $result.bearing = bearing;
    }
    if (source != null) {
      $result.source = source;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    if (deviceId != null) {
      $result.deviceId = deviceId;
    }
    return $result;
  }
  LocationPointObject._() : super();
  factory LocationPointObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LocationPointObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LocationPointObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$2.Timestamp>(3, _omitFieldNames ? '' : 'timestamp', subBuilder: $2.Timestamp.create)
    ..a<$core.double>(4, _omitFieldNames ? '' : 'latitude', $pb.PbFieldType.OD)
    ..a<$core.double>(5, _omitFieldNames ? '' : 'longitude', $pb.PbFieldType.OD)
    ..a<$core.double>(6, _omitFieldNames ? '' : 'altitude', $pb.PbFieldType.OD)
    ..a<$core.double>(7, _omitFieldNames ? '' : 'accuracy', $pb.PbFieldType.OD)
    ..a<$core.double>(8, _omitFieldNames ? '' : 'speed', $pb.PbFieldType.OD)
    ..a<$core.double>(9, _omitFieldNames ? '' : 'bearing', $pb.PbFieldType.OD)
    ..e<LocationSource>(10, _omitFieldNames ? '' : 'source', $pb.PbFieldType.OE, defaultOrMaker: LocationSource.LOCATION_SOURCE_UNSPECIFIED, valueOf: LocationSource.valueOf, enumValues: LocationSource.values)
    ..aOM<$6.Struct>(11, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..aOM<$2.Timestamp>(12, _omitFieldNames ? '' : 'createdAt', subBuilder: $2.Timestamp.create)
    ..aOS(13, _omitFieldNames ? '' : 'deviceId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LocationPointObject clone() => LocationPointObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LocationPointObject copyWith(void Function(LocationPointObject) updates) => super.copyWith((message) => updates(message as LocationPointObject)) as LocationPointObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LocationPointObject create() => LocationPointObject._();
  LocationPointObject createEmptyInstance() => create();
  static $pb.PbList<LocationPointObject> createRepeated() => $pb.PbList<LocationPointObject>();
  @$core.pragma('dart2js:noInline')
  static LocationPointObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LocationPointObject>(create);
  static LocationPointObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => clearField(2);

  @$pb.TagNumber(3)
  $2.Timestamp get timestamp => $_getN(2);
  @$pb.TagNumber(3)
  set timestamp($2.Timestamp v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(3)
  void clearTimestamp() => clearField(3);
  @$pb.TagNumber(3)
  $2.Timestamp ensureTimestamp() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.double get latitude => $_getN(3);
  @$pb.TagNumber(4)
  set latitude($core.double v) { $_setDouble(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasLatitude() => $_has(3);
  @$pb.TagNumber(4)
  void clearLatitude() => clearField(4);

  @$pb.TagNumber(5)
  $core.double get longitude => $_getN(4);
  @$pb.TagNumber(5)
  set longitude($core.double v) { $_setDouble(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasLongitude() => $_has(4);
  @$pb.TagNumber(5)
  void clearLongitude() => clearField(5);

  @$pb.TagNumber(6)
  $core.double get altitude => $_getN(5);
  @$pb.TagNumber(6)
  set altitude($core.double v) { $_setDouble(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasAltitude() => $_has(5);
  @$pb.TagNumber(6)
  void clearAltitude() => clearField(6);

  @$pb.TagNumber(7)
  $core.double get accuracy => $_getN(6);
  @$pb.TagNumber(7)
  set accuracy($core.double v) { $_setDouble(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasAccuracy() => $_has(6);
  @$pb.TagNumber(7)
  void clearAccuracy() => clearField(7);

  @$pb.TagNumber(8)
  $core.double get speed => $_getN(7);
  @$pb.TagNumber(8)
  set speed($core.double v) { $_setDouble(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasSpeed() => $_has(7);
  @$pb.TagNumber(8)
  void clearSpeed() => clearField(8);

  @$pb.TagNumber(9)
  $core.double get bearing => $_getN(8);
  @$pb.TagNumber(9)
  set bearing($core.double v) { $_setDouble(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasBearing() => $_has(8);
  @$pb.TagNumber(9)
  void clearBearing() => clearField(9);

  @$pb.TagNumber(10)
  LocationSource get source => $_getN(9);
  @$pb.TagNumber(10)
  set source(LocationSource v) { setField(10, v); }
  @$pb.TagNumber(10)
  $core.bool hasSource() => $_has(9);
  @$pb.TagNumber(10)
  void clearSource() => clearField(10);

  @$pb.TagNumber(11)
  $6.Struct get extra => $_getN(10);
  @$pb.TagNumber(11)
  set extra($6.Struct v) { setField(11, v); }
  @$pb.TagNumber(11)
  $core.bool hasExtra() => $_has(10);
  @$pb.TagNumber(11)
  void clearExtra() => clearField(11);
  @$pb.TagNumber(11)
  $6.Struct ensureExtra() => $_ensure(10);

  @$pb.TagNumber(12)
  $2.Timestamp get createdAt => $_getN(11);
  @$pb.TagNumber(12)
  set createdAt($2.Timestamp v) { setField(12, v); }
  @$pb.TagNumber(12)
  $core.bool hasCreatedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearCreatedAt() => clearField(12);
  @$pb.TagNumber(12)
  $2.Timestamp ensureCreatedAt() => $_ensure(11);

  @$pb.TagNumber(13)
  $core.String get deviceId => $_getSZ(12);
  @$pb.TagNumber(13)
  set deviceId($core.String v) { $_setString(12, v); }
  @$pb.TagNumber(13)
  $core.bool hasDeviceId() => $_has(12);
  @$pb.TagNumber(13)
  void clearDeviceId() => clearField(13);
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
    $6.Struct? extra,
    $2.Timestamp? createdAt,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (ownerId != null) {
      $result.ownerId = ownerId;
    }
    if (name != null) {
      $result.name = name;
    }
    if (description != null) {
      $result.description = description;
    }
    if (areaType != null) {
      $result.areaType = areaType;
    }
    if (geometry != null) {
      $result.geometry = geometry;
    }
    if (areaM2 != null) {
      $result.areaM2 = areaM2;
    }
    if (perimeterM != null) {
      $result.perimeterM = perimeterM;
    }
    if (state != null) {
      $result.state = state;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    return $result;
  }
  AreaObject._() : super();
  factory AreaObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AreaObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AreaObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'ownerId')
    ..aOS(3, _omitFieldNames ? '' : 'name')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..e<AreaType>(5, _omitFieldNames ? '' : 'areaType', $pb.PbFieldType.OE, defaultOrMaker: AreaType.AREA_TYPE_UNSPECIFIED, valueOf: AreaType.valueOf, enumValues: AreaType.values)
    ..aOS(6, _omitFieldNames ? '' : 'geometry')
    ..a<$core.double>(7, _omitFieldNames ? '' : 'areaM2', $pb.PbFieldType.OD)
    ..a<$core.double>(8, _omitFieldNames ? '' : 'perimeterM', $pb.PbFieldType.OD)
    ..a<$core.int>(9, _omitFieldNames ? '' : 'state', $pb.PbFieldType.O3)
    ..aOM<$6.Struct>(10, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..aOM<$2.Timestamp>(11, _omitFieldNames ? '' : 'createdAt', subBuilder: $2.Timestamp.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AreaObject clone() => AreaObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AreaObject copyWith(void Function(AreaObject) updates) => super.copyWith((message) => updates(message as AreaObject)) as AreaObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AreaObject create() => AreaObject._();
  AreaObject createEmptyInstance() => create();
  static $pb.PbList<AreaObject> createRepeated() => $pb.PbList<AreaObject>();
  @$core.pragma('dart2js:noInline')
  static AreaObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AreaObject>(create);
  static AreaObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get ownerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ownerId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasOwnerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearOwnerId() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(3)
  set name($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(3)
  void clearName() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => clearField(4);

  @$pb.TagNumber(5)
  AreaType get areaType => $_getN(4);
  @$pb.TagNumber(5)
  set areaType(AreaType v) { setField(5, v); }
  @$pb.TagNumber(5)
  $core.bool hasAreaType() => $_has(4);
  @$pb.TagNumber(5)
  void clearAreaType() => clearField(5);

  @$pb.TagNumber(6)
  $core.String get geometry => $_getSZ(5);
  @$pb.TagNumber(6)
  set geometry($core.String v) { $_setString(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasGeometry() => $_has(5);
  @$pb.TagNumber(6)
  void clearGeometry() => clearField(6);

  @$pb.TagNumber(7)
  $core.double get areaM2 => $_getN(6);
  @$pb.TagNumber(7)
  set areaM2($core.double v) { $_setDouble(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasAreaM2() => $_has(6);
  @$pb.TagNumber(7)
  void clearAreaM2() => clearField(7);

  @$pb.TagNumber(8)
  $core.double get perimeterM => $_getN(7);
  @$pb.TagNumber(8)
  set perimeterM($core.double v) { $_setDouble(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasPerimeterM() => $_has(7);
  @$pb.TagNumber(8)
  void clearPerimeterM() => clearField(8);

  @$pb.TagNumber(9)
  $core.int get state => $_getIZ(8);
  @$pb.TagNumber(9)
  set state($core.int v) { $_setSignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasState() => $_has(8);
  @$pb.TagNumber(9)
  void clearState() => clearField(9);

  @$pb.TagNumber(10)
  $6.Struct get extra => $_getN(9);
  @$pb.TagNumber(10)
  set extra($6.Struct v) { setField(10, v); }
  @$pb.TagNumber(10)
  $core.bool hasExtra() => $_has(9);
  @$pb.TagNumber(10)
  void clearExtra() => clearField(10);
  @$pb.TagNumber(10)
  $6.Struct ensureExtra() => $_ensure(9);

  @$pb.TagNumber(11)
  $2.Timestamp get createdAt => $_getN(10);
  @$pb.TagNumber(11)
  set createdAt($2.Timestamp v) { setField(11, v); }
  @$pb.TagNumber(11)
  $core.bool hasCreatedAt() => $_has(10);
  @$pb.TagNumber(11)
  void clearCreatedAt() => clearField(11);
  @$pb.TagNumber(11)
  $2.Timestamp ensureCreatedAt() => $_ensure(10);
}

class GeoEventObject extends $pb.GeneratedMessage {
  factory GeoEventObject({
    $core.String? id,
    $core.String? subjectId,
    $core.String? areaId,
    GeoEventType? eventType,
    $2.Timestamp? timestamp,
    $core.double? confidence,
    $core.String? pointId,
    $6.Struct? extra,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (areaId != null) {
      $result.areaId = areaId;
    }
    if (eventType != null) {
      $result.eventType = eventType;
    }
    if (timestamp != null) {
      $result.timestamp = timestamp;
    }
    if (confidence != null) {
      $result.confidence = confidence;
    }
    if (pointId != null) {
      $result.pointId = pointId;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    return $result;
  }
  GeoEventObject._() : super();
  factory GeoEventObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GeoEventObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GeoEventObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOS(3, _omitFieldNames ? '' : 'areaId')
    ..e<GeoEventType>(4, _omitFieldNames ? '' : 'eventType', $pb.PbFieldType.OE, defaultOrMaker: GeoEventType.GEO_EVENT_TYPE_UNSPECIFIED, valueOf: GeoEventType.valueOf, enumValues: GeoEventType.values)
    ..aOM<$2.Timestamp>(5, _omitFieldNames ? '' : 'timestamp', subBuilder: $2.Timestamp.create)
    ..a<$core.double>(6, _omitFieldNames ? '' : 'confidence', $pb.PbFieldType.OD)
    ..aOS(7, _omitFieldNames ? '' : 'pointId')
    ..aOM<$6.Struct>(8, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GeoEventObject clone() => GeoEventObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GeoEventObject copyWith(void Function(GeoEventObject) updates) => super.copyWith((message) => updates(message as GeoEventObject)) as GeoEventObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoEventObject create() => GeoEventObject._();
  GeoEventObject createEmptyInstance() => create();
  static $pb.PbList<GeoEventObject> createRepeated() => $pb.PbList<GeoEventObject>();
  @$core.pragma('dart2js:noInline')
  static GeoEventObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GeoEventObject>(create);
  static GeoEventObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get areaId => $_getSZ(2);
  @$pb.TagNumber(3)
  set areaId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAreaId() => $_has(2);
  @$pb.TagNumber(3)
  void clearAreaId() => clearField(3);

  @$pb.TagNumber(4)
  GeoEventType get eventType => $_getN(3);
  @$pb.TagNumber(4)
  set eventType(GeoEventType v) { setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasEventType() => $_has(3);
  @$pb.TagNumber(4)
  void clearEventType() => clearField(4);

  @$pb.TagNumber(5)
  $2.Timestamp get timestamp => $_getN(4);
  @$pb.TagNumber(5)
  set timestamp($2.Timestamp v) { setField(5, v); }
  @$pb.TagNumber(5)
  $core.bool hasTimestamp() => $_has(4);
  @$pb.TagNumber(5)
  void clearTimestamp() => clearField(5);
  @$pb.TagNumber(5)
  $2.Timestamp ensureTimestamp() => $_ensure(4);

  @$pb.TagNumber(6)
  $core.double get confidence => $_getN(5);
  @$pb.TagNumber(6)
  set confidence($core.double v) { $_setDouble(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasConfidence() => $_has(5);
  @$pb.TagNumber(6)
  void clearConfidence() => clearField(6);

  @$pb.TagNumber(7)
  $core.String get pointId => $_getSZ(6);
  @$pb.TagNumber(7)
  set pointId($core.String v) { $_setString(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasPointId() => $_has(6);
  @$pb.TagNumber(7)
  void clearPointId() => clearField(7);

  @$pb.TagNumber(8)
  $6.Struct get extra => $_getN(7);
  @$pb.TagNumber(8)
  set extra($6.Struct v) { setField(8, v); }
  @$pb.TagNumber(8)
  $core.bool hasExtra() => $_has(7);
  @$pb.TagNumber(8)
  void clearExtra() => clearField(8);
  @$pb.TagNumber(8)
  $6.Struct ensureExtra() => $_ensure(7);
}

class AreaSubjectObject extends $pb.GeneratedMessage {
  factory AreaSubjectObject({
    $core.String? subjectId,
    $2.Timestamp? enterTimestamp,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (enterTimestamp != null) {
      $result.enterTimestamp = enterTimestamp;
    }
    return $result;
  }
  AreaSubjectObject._() : super();
  factory AreaSubjectObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AreaSubjectObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AreaSubjectObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$2.Timestamp>(2, _omitFieldNames ? '' : 'enterTimestamp', subBuilder: $2.Timestamp.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AreaSubjectObject clone() => AreaSubjectObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AreaSubjectObject copyWith(void Function(AreaSubjectObject) updates) => super.copyWith((message) => updates(message as AreaSubjectObject)) as AreaSubjectObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AreaSubjectObject create() => AreaSubjectObject._();
  AreaSubjectObject createEmptyInstance() => create();
  static $pb.PbList<AreaSubjectObject> createRepeated() => $pb.PbList<AreaSubjectObject>();
  @$core.pragma('dart2js:noInline')
  static AreaSubjectObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AreaSubjectObject>(create);
  static AreaSubjectObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $2.Timestamp get enterTimestamp => $_getN(1);
  @$pb.TagNumber(2)
  set enterTimestamp($2.Timestamp v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasEnterTimestamp() => $_has(1);
  @$pb.TagNumber(2)
  void clearEnterTimestamp() => clearField(2);
  @$pb.TagNumber(2)
  $2.Timestamp ensureEnterTimestamp() => $_ensure(1);
}

class NearbySubjectObject extends $pb.GeneratedMessage {
  factory NearbySubjectObject({
    $core.String? subjectId,
    $core.double? distanceMeters,
    $2.Timestamp? lastSeen,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (distanceMeters != null) {
      $result.distanceMeters = distanceMeters;
    }
    if (lastSeen != null) {
      $result.lastSeen = lastSeen;
    }
    return $result;
  }
  NearbySubjectObject._() : super();
  factory NearbySubjectObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory NearbySubjectObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'NearbySubjectObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..a<$core.double>(2, _omitFieldNames ? '' : 'distanceMeters', $pb.PbFieldType.OD)
    ..aOM<$2.Timestamp>(3, _omitFieldNames ? '' : 'lastSeen', subBuilder: $2.Timestamp.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  NearbySubjectObject clone() => NearbySubjectObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  NearbySubjectObject copyWith(void Function(NearbySubjectObject) updates) => super.copyWith((message) => updates(message as NearbySubjectObject)) as NearbySubjectObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NearbySubjectObject create() => NearbySubjectObject._();
  NearbySubjectObject createEmptyInstance() => create();
  static $pb.PbList<NearbySubjectObject> createRepeated() => $pb.PbList<NearbySubjectObject>();
  @$core.pragma('dart2js:noInline')
  static NearbySubjectObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<NearbySubjectObject>(create);
  static NearbySubjectObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $core.double get distanceMeters => $_getN(1);
  @$pb.TagNumber(2)
  set distanceMeters($core.double v) { $_setDouble(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasDistanceMeters() => $_has(1);
  @$pb.TagNumber(2)
  void clearDistanceMeters() => clearField(2);

  @$pb.TagNumber(3)
  $2.Timestamp get lastSeen => $_getN(2);
  @$pb.TagNumber(3)
  set lastSeen($2.Timestamp v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasLastSeen() => $_has(2);
  @$pb.TagNumber(3)
  void clearLastSeen() => clearField(3);
  @$pb.TagNumber(3)
  $2.Timestamp ensureLastSeen() => $_ensure(2);
}

class NearbyAreaObject extends $pb.GeneratedMessage {
  factory NearbyAreaObject({
    $core.String? areaId,
    $core.String? name,
    AreaType? areaType,
    $core.double? distanceMeters,
  }) {
    final $result = create();
    if (areaId != null) {
      $result.areaId = areaId;
    }
    if (name != null) {
      $result.name = name;
    }
    if (areaType != null) {
      $result.areaType = areaType;
    }
    if (distanceMeters != null) {
      $result.distanceMeters = distanceMeters;
    }
    return $result;
  }
  NearbyAreaObject._() : super();
  factory NearbyAreaObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory NearbyAreaObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'NearbyAreaObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'areaId')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..e<AreaType>(3, _omitFieldNames ? '' : 'areaType', $pb.PbFieldType.OE, defaultOrMaker: AreaType.AREA_TYPE_UNSPECIFIED, valueOf: AreaType.valueOf, enumValues: AreaType.values)
    ..a<$core.double>(4, _omitFieldNames ? '' : 'distanceMeters', $pb.PbFieldType.OD)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  NearbyAreaObject clone() => NearbyAreaObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  NearbyAreaObject copyWith(void Function(NearbyAreaObject) updates) => super.copyWith((message) => updates(message as NearbyAreaObject)) as NearbyAreaObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NearbyAreaObject create() => NearbyAreaObject._();
  NearbyAreaObject createEmptyInstance() => create();
  static $pb.PbList<NearbyAreaObject> createRepeated() => $pb.PbList<NearbyAreaObject>();
  @$core.pragma('dart2js:noInline')
  static NearbyAreaObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<NearbyAreaObject>(create);
  static NearbyAreaObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get areaId => $_getSZ(0);
  @$pb.TagNumber(1)
  set areaId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAreaId() => $_has(0);
  @$pb.TagNumber(1)
  void clearAreaId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => clearField(2);

  @$pb.TagNumber(3)
  AreaType get areaType => $_getN(2);
  @$pb.TagNumber(3)
  set areaType(AreaType v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasAreaType() => $_has(2);
  @$pb.TagNumber(3)
  void clearAreaType() => clearField(3);

  @$pb.TagNumber(4)
  $core.double get distanceMeters => $_getN(3);
  @$pb.TagNumber(4)
  set distanceMeters($core.double v) { $_setDouble(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasDistanceMeters() => $_has(3);
  @$pb.TagNumber(4)
  void clearDistanceMeters() => clearField(4);
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
    $6.Struct? extra,
    $2.Timestamp? createdAt,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (ownerId != null) {
      $result.ownerId = ownerId;
    }
    if (name != null) {
      $result.name = name;
    }
    if (description != null) {
      $result.description = description;
    }
    if (geometry != null) {
      $result.geometry = geometry;
    }
    if (lengthM != null) {
      $result.lengthM = lengthM;
    }
    if (state != null) {
      $result.state = state;
    }
    if (deviationThresholdM != null) {
      $result.deviationThresholdM = deviationThresholdM;
    }
    if (deviationConsecutiveCount != null) {
      $result.deviationConsecutiveCount = deviationConsecutiveCount;
    }
    if (deviationCooldownSec != null) {
      $result.deviationCooldownSec = deviationCooldownSec;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    return $result;
  }
  RouteObject._() : super();
  factory RouteObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RouteObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RouteObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'ownerId')
    ..aOS(3, _omitFieldNames ? '' : 'name')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aOS(5, _omitFieldNames ? '' : 'geometry')
    ..a<$core.double>(6, _omitFieldNames ? '' : 'lengthM', $pb.PbFieldType.OD)
    ..a<$core.int>(7, _omitFieldNames ? '' : 'state', $pb.PbFieldType.O3)
    ..a<$core.double>(8, _omitFieldNames ? '' : 'deviationThresholdM', $pb.PbFieldType.OD)
    ..a<$core.int>(9, _omitFieldNames ? '' : 'deviationConsecutiveCount', $pb.PbFieldType.O3)
    ..a<$core.int>(10, _omitFieldNames ? '' : 'deviationCooldownSec', $pb.PbFieldType.O3)
    ..aOM<$6.Struct>(11, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..aOM<$2.Timestamp>(12, _omitFieldNames ? '' : 'createdAt', subBuilder: $2.Timestamp.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RouteObject clone() => RouteObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RouteObject copyWith(void Function(RouteObject) updates) => super.copyWith((message) => updates(message as RouteObject)) as RouteObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RouteObject create() => RouteObject._();
  RouteObject createEmptyInstance() => create();
  static $pb.PbList<RouteObject> createRepeated() => $pb.PbList<RouteObject>();
  @$core.pragma('dart2js:noInline')
  static RouteObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RouteObject>(create);
  static RouteObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get ownerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ownerId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasOwnerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearOwnerId() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(3)
  set name($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(3)
  void clearName() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get geometry => $_getSZ(4);
  @$pb.TagNumber(5)
  set geometry($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasGeometry() => $_has(4);
  @$pb.TagNumber(5)
  void clearGeometry() => clearField(5);

  @$pb.TagNumber(6)
  $core.double get lengthM => $_getN(5);
  @$pb.TagNumber(6)
  set lengthM($core.double v) { $_setDouble(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasLengthM() => $_has(5);
  @$pb.TagNumber(6)
  void clearLengthM() => clearField(6);

  @$pb.TagNumber(7)
  $core.int get state => $_getIZ(6);
  @$pb.TagNumber(7)
  set state($core.int v) { $_setSignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasState() => $_has(6);
  @$pb.TagNumber(7)
  void clearState() => clearField(7);

  @$pb.TagNumber(8)
  $core.double get deviationThresholdM => $_getN(7);
  @$pb.TagNumber(8)
  set deviationThresholdM($core.double v) { $_setDouble(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasDeviationThresholdM() => $_has(7);
  @$pb.TagNumber(8)
  void clearDeviationThresholdM() => clearField(8);

  @$pb.TagNumber(9)
  $core.int get deviationConsecutiveCount => $_getIZ(8);
  @$pb.TagNumber(9)
  set deviationConsecutiveCount($core.int v) { $_setSignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasDeviationConsecutiveCount() => $_has(8);
  @$pb.TagNumber(9)
  void clearDeviationConsecutiveCount() => clearField(9);

  @$pb.TagNumber(10)
  $core.int get deviationCooldownSec => $_getIZ(9);
  @$pb.TagNumber(10)
  set deviationCooldownSec($core.int v) { $_setSignedInt32(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasDeviationCooldownSec() => $_has(9);
  @$pb.TagNumber(10)
  void clearDeviationCooldownSec() => clearField(10);

  @$pb.TagNumber(11)
  $6.Struct get extra => $_getN(10);
  @$pb.TagNumber(11)
  set extra($6.Struct v) { setField(11, v); }
  @$pb.TagNumber(11)
  $core.bool hasExtra() => $_has(10);
  @$pb.TagNumber(11)
  void clearExtra() => clearField(11);
  @$pb.TagNumber(11)
  $6.Struct ensureExtra() => $_ensure(10);

  @$pb.TagNumber(12)
  $2.Timestamp get createdAt => $_getN(11);
  @$pb.TagNumber(12)
  set createdAt($2.Timestamp v) { setField(12, v); }
  @$pb.TagNumber(12)
  $core.bool hasCreatedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearCreatedAt() => clearField(12);
  @$pb.TagNumber(12)
  $2.Timestamp ensureCreatedAt() => $_ensure(11);
}

class RouteAssignmentObject extends $pb.GeneratedMessage {
  factory RouteAssignmentObject({
    $core.String? id,
    $core.String? subjectId,
    $core.String? routeId,
    $2.Timestamp? validFrom,
    $2.Timestamp? validUntil,
    $core.int? state,
    $6.Struct? extra,
    $2.Timestamp? createdAt,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (routeId != null) {
      $result.routeId = routeId;
    }
    if (validFrom != null) {
      $result.validFrom = validFrom;
    }
    if (validUntil != null) {
      $result.validUntil = validUntil;
    }
    if (state != null) {
      $result.state = state;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    return $result;
  }
  RouteAssignmentObject._() : super();
  factory RouteAssignmentObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RouteAssignmentObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RouteAssignmentObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOS(3, _omitFieldNames ? '' : 'routeId')
    ..aOM<$2.Timestamp>(4, _omitFieldNames ? '' : 'validFrom', subBuilder: $2.Timestamp.create)
    ..aOM<$2.Timestamp>(5, _omitFieldNames ? '' : 'validUntil', subBuilder: $2.Timestamp.create)
    ..a<$core.int>(6, _omitFieldNames ? '' : 'state', $pb.PbFieldType.O3)
    ..aOM<$6.Struct>(7, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..aOM<$2.Timestamp>(8, _omitFieldNames ? '' : 'createdAt', subBuilder: $2.Timestamp.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RouteAssignmentObject clone() => RouteAssignmentObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RouteAssignmentObject copyWith(void Function(RouteAssignmentObject) updates) => super.copyWith((message) => updates(message as RouteAssignmentObject)) as RouteAssignmentObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RouteAssignmentObject create() => RouteAssignmentObject._();
  RouteAssignmentObject createEmptyInstance() => create();
  static $pb.PbList<RouteAssignmentObject> createRepeated() => $pb.PbList<RouteAssignmentObject>();
  @$core.pragma('dart2js:noInline')
  static RouteAssignmentObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RouteAssignmentObject>(create);
  static RouteAssignmentObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get routeId => $_getSZ(2);
  @$pb.TagNumber(3)
  set routeId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasRouteId() => $_has(2);
  @$pb.TagNumber(3)
  void clearRouteId() => clearField(3);

  @$pb.TagNumber(4)
  $2.Timestamp get validFrom => $_getN(3);
  @$pb.TagNumber(4)
  set validFrom($2.Timestamp v) { setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasValidFrom() => $_has(3);
  @$pb.TagNumber(4)
  void clearValidFrom() => clearField(4);
  @$pb.TagNumber(4)
  $2.Timestamp ensureValidFrom() => $_ensure(3);

  @$pb.TagNumber(5)
  $2.Timestamp get validUntil => $_getN(4);
  @$pb.TagNumber(5)
  set validUntil($2.Timestamp v) { setField(5, v); }
  @$pb.TagNumber(5)
  $core.bool hasValidUntil() => $_has(4);
  @$pb.TagNumber(5)
  void clearValidUntil() => clearField(5);
  @$pb.TagNumber(5)
  $2.Timestamp ensureValidUntil() => $_ensure(4);

  @$pb.TagNumber(6)
  $core.int get state => $_getIZ(5);
  @$pb.TagNumber(6)
  set state($core.int v) { $_setSignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(6)
  void clearState() => clearField(6);

  @$pb.TagNumber(7)
  $6.Struct get extra => $_getN(6);
  @$pb.TagNumber(7)
  set extra($6.Struct v) { setField(7, v); }
  @$pb.TagNumber(7)
  $core.bool hasExtra() => $_has(6);
  @$pb.TagNumber(7)
  void clearExtra() => clearField(7);
  @$pb.TagNumber(7)
  $6.Struct ensureExtra() => $_ensure(6);

  @$pb.TagNumber(8)
  $2.Timestamp get createdAt => $_getN(7);
  @$pb.TagNumber(8)
  set createdAt($2.Timestamp v) { setField(8, v); }
  @$pb.TagNumber(8)
  $core.bool hasCreatedAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearCreatedAt() => clearField(8);
  @$pb.TagNumber(8)
  $2.Timestamp ensureCreatedAt() => $_ensure(7);
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
    $2.Timestamp? timestamp,
    $6.Struct? extra,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (routeId != null) {
      $result.routeId = routeId;
    }
    if (eventType != null) {
      $result.eventType = eventType;
    }
    if (distanceMeters != null) {
      $result.distanceMeters = distanceMeters;
    }
    if (latitude != null) {
      $result.latitude = latitude;
    }
    if (longitude != null) {
      $result.longitude = longitude;
    }
    if (timestamp != null) {
      $result.timestamp = timestamp;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    return $result;
  }
  RouteDeviationEventObject._() : super();
  factory RouteDeviationEventObject.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RouteDeviationEventObject.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RouteDeviationEventObject', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'subjectId')
    ..aOS(3, _omitFieldNames ? '' : 'routeId')
    ..e<RouteDeviationEventType>(4, _omitFieldNames ? '' : 'eventType', $pb.PbFieldType.OE, defaultOrMaker: RouteDeviationEventType.ROUTE_DEVIATION_EVENT_TYPE_UNSPECIFIED, valueOf: RouteDeviationEventType.valueOf, enumValues: RouteDeviationEventType.values)
    ..a<$core.double>(5, _omitFieldNames ? '' : 'distanceMeters', $pb.PbFieldType.OD)
    ..a<$core.double>(6, _omitFieldNames ? '' : 'latitude', $pb.PbFieldType.OD)
    ..a<$core.double>(7, _omitFieldNames ? '' : 'longitude', $pb.PbFieldType.OD)
    ..aOM<$2.Timestamp>(8, _omitFieldNames ? '' : 'timestamp', subBuilder: $2.Timestamp.create)
    ..aOM<$6.Struct>(9, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RouteDeviationEventObject clone() => RouteDeviationEventObject()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RouteDeviationEventObject copyWith(void Function(RouteDeviationEventObject) updates) => super.copyWith((message) => updates(message as RouteDeviationEventObject)) as RouteDeviationEventObject;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RouteDeviationEventObject create() => RouteDeviationEventObject._();
  RouteDeviationEventObject createEmptyInstance() => create();
  static $pb.PbList<RouteDeviationEventObject> createRepeated() => $pb.PbList<RouteDeviationEventObject>();
  @$core.pragma('dart2js:noInline')
  static RouteDeviationEventObject getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RouteDeviationEventObject>(create);
  static RouteDeviationEventObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get subjectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set subjectId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSubjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearSubjectId() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get routeId => $_getSZ(2);
  @$pb.TagNumber(3)
  set routeId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasRouteId() => $_has(2);
  @$pb.TagNumber(3)
  void clearRouteId() => clearField(3);

  @$pb.TagNumber(4)
  RouteDeviationEventType get eventType => $_getN(3);
  @$pb.TagNumber(4)
  set eventType(RouteDeviationEventType v) { setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasEventType() => $_has(3);
  @$pb.TagNumber(4)
  void clearEventType() => clearField(4);

  @$pb.TagNumber(5)
  $core.double get distanceMeters => $_getN(4);
  @$pb.TagNumber(5)
  set distanceMeters($core.double v) { $_setDouble(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasDistanceMeters() => $_has(4);
  @$pb.TagNumber(5)
  void clearDistanceMeters() => clearField(5);

  @$pb.TagNumber(6)
  $core.double get latitude => $_getN(5);
  @$pb.TagNumber(6)
  set latitude($core.double v) { $_setDouble(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasLatitude() => $_has(5);
  @$pb.TagNumber(6)
  void clearLatitude() => clearField(6);

  @$pb.TagNumber(7)
  $core.double get longitude => $_getN(6);
  @$pb.TagNumber(7)
  set longitude($core.double v) { $_setDouble(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasLongitude() => $_has(6);
  @$pb.TagNumber(7)
  void clearLongitude() => clearField(7);

  @$pb.TagNumber(8)
  $2.Timestamp get timestamp => $_getN(7);
  @$pb.TagNumber(8)
  set timestamp($2.Timestamp v) { setField(8, v); }
  @$pb.TagNumber(8)
  $core.bool hasTimestamp() => $_has(7);
  @$pb.TagNumber(8)
  void clearTimestamp() => clearField(8);
  @$pb.TagNumber(8)
  $2.Timestamp ensureTimestamp() => $_ensure(7);

  @$pb.TagNumber(9)
  $6.Struct get extra => $_getN(8);
  @$pb.TagNumber(9)
  set extra($6.Struct v) { setField(9, v); }
  @$pb.TagNumber(9)
  $core.bool hasExtra() => $_has(8);
  @$pb.TagNumber(9)
  void clearExtra() => clearField(9);
  @$pb.TagNumber(9)
  $6.Struct ensureExtra() => $_ensure(8);
}

class IngestLocationsRequest extends $pb.GeneratedMessage {
  factory IngestLocationsRequest({
    $core.String? subjectId,
    $core.Iterable<LocationPointInput>? points,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (points != null) {
      $result.points.addAll(points);
    }
    return $result;
  }
  IngestLocationsRequest._() : super();
  factory IngestLocationsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory IngestLocationsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'IngestLocationsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..pc<LocationPointInput>(2, _omitFieldNames ? '' : 'points', $pb.PbFieldType.PM, subBuilder: LocationPointInput.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  IngestLocationsRequest clone() => IngestLocationsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  IngestLocationsRequest copyWith(void Function(IngestLocationsRequest) updates) => super.copyWith((message) => updates(message as IngestLocationsRequest)) as IngestLocationsRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IngestLocationsRequest create() => IngestLocationsRequest._();
  IngestLocationsRequest createEmptyInstance() => create();
  static $pb.PbList<IngestLocationsRequest> createRepeated() => $pb.PbList<IngestLocationsRequest>();
  @$core.pragma('dart2js:noInline')
  static IngestLocationsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<IngestLocationsRequest>(create);
  static IngestLocationsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<LocationPointInput> get points => $_getList(1);
}

class IngestLocationsResponse extends $pb.GeneratedMessage {
  factory IngestLocationsResponse({
    $core.int? accepted,
    $core.int? rejected,
  }) {
    final $result = create();
    if (accepted != null) {
      $result.accepted = accepted;
    }
    if (rejected != null) {
      $result.rejected = rejected;
    }
    return $result;
  }
  IngestLocationsResponse._() : super();
  factory IngestLocationsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory IngestLocationsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'IngestLocationsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'accepted', $pb.PbFieldType.O3)
    ..a<$core.int>(2, _omitFieldNames ? '' : 'rejected', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  IngestLocationsResponse clone() => IngestLocationsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  IngestLocationsResponse copyWith(void Function(IngestLocationsResponse) updates) => super.copyWith((message) => updates(message as IngestLocationsResponse)) as IngestLocationsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IngestLocationsResponse create() => IngestLocationsResponse._();
  IngestLocationsResponse createEmptyInstance() => create();
  static $pb.PbList<IngestLocationsResponse> createRepeated() => $pb.PbList<IngestLocationsResponse>();
  @$core.pragma('dart2js:noInline')
  static IngestLocationsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<IngestLocationsResponse>(create);
  static IngestLocationsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get accepted => $_getIZ(0);
  @$pb.TagNumber(1)
  set accepted($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAccepted() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccepted() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get rejected => $_getIZ(1);
  @$pb.TagNumber(2)
  set rejected($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRejected() => $_has(1);
  @$pb.TagNumber(2)
  void clearRejected() => clearField(2);
}

class CreateAreaRequest extends $pb.GeneratedMessage {
  factory CreateAreaRequest({
    AreaObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  CreateAreaRequest._() : super();
  factory CreateAreaRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateAreaRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateAreaRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data', subBuilder: AreaObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateAreaRequest clone() => CreateAreaRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateAreaRequest copyWith(void Function(CreateAreaRequest) updates) => super.copyWith((message) => updates(message as CreateAreaRequest)) as CreateAreaRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateAreaRequest create() => CreateAreaRequest._();
  CreateAreaRequest createEmptyInstance() => create();
  static $pb.PbList<CreateAreaRequest> createRepeated() => $pb.PbList<CreateAreaRequest>();
  @$core.pragma('dart2js:noInline')
  static CreateAreaRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateAreaRequest>(create);
  static CreateAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class CreateAreaResponse extends $pb.GeneratedMessage {
  factory CreateAreaResponse({
    AreaObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  CreateAreaResponse._() : super();
  factory CreateAreaResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateAreaResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateAreaResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data', subBuilder: AreaObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateAreaResponse clone() => CreateAreaResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateAreaResponse copyWith(void Function(CreateAreaResponse) updates) => super.copyWith((message) => updates(message as CreateAreaResponse)) as CreateAreaResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateAreaResponse create() => CreateAreaResponse._();
  CreateAreaResponse createEmptyInstance() => create();
  static $pb.PbList<CreateAreaResponse> createRepeated() => $pb.PbList<CreateAreaResponse>();
  @$core.pragma('dart2js:noInline')
  static CreateAreaResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateAreaResponse>(create);
  static CreateAreaResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class GetAreaRequest extends $pb.GeneratedMessage {
  factory GetAreaRequest({
    $core.String? id,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    return $result;
  }
  GetAreaRequest._() : super();
  factory GetAreaRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAreaRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAreaRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAreaRequest clone() => GetAreaRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAreaRequest copyWith(void Function(GetAreaRequest) updates) => super.copyWith((message) => updates(message as GetAreaRequest)) as GetAreaRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAreaRequest create() => GetAreaRequest._();
  GetAreaRequest createEmptyInstance() => create();
  static $pb.PbList<GetAreaRequest> createRepeated() => $pb.PbList<GetAreaRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAreaRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAreaRequest>(create);
  static GetAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);
}

class GetAreaResponse extends $pb.GeneratedMessage {
  factory GetAreaResponse({
    AreaObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  GetAreaResponse._() : super();
  factory GetAreaResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAreaResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAreaResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data', subBuilder: AreaObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAreaResponse clone() => GetAreaResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAreaResponse copyWith(void Function(GetAreaResponse) updates) => super.copyWith((message) => updates(message as GetAreaResponse)) as GetAreaResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAreaResponse create() => GetAreaResponse._();
  GetAreaResponse createEmptyInstance() => create();
  static $pb.PbList<GetAreaResponse> createRepeated() => $pb.PbList<GetAreaResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAreaResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAreaResponse>(create);
  static GetAreaResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
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
    $6.Struct? extra,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (name != null) {
      $result.name = name;
    }
    if (description != null) {
      $result.description = description;
    }
    if (areaType != null) {
      $result.areaType = areaType;
    }
    if (geometry != null) {
      $result.geometry = geometry;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    return $result;
  }
  UpdateAreaRequest._() : super();
  factory UpdateAreaRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UpdateAreaRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateAreaRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..e<AreaType>(4, _omitFieldNames ? '' : 'areaType', $pb.PbFieldType.OE, defaultOrMaker: AreaType.AREA_TYPE_UNSPECIFIED, valueOf: AreaType.valueOf, enumValues: AreaType.values)
    ..aOS(5, _omitFieldNames ? '' : 'geometry')
    ..aOM<$6.Struct>(6, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UpdateAreaRequest clone() => UpdateAreaRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UpdateAreaRequest copyWith(void Function(UpdateAreaRequest) updates) => super.copyWith((message) => updates(message as UpdateAreaRequest)) as UpdateAreaRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAreaRequest create() => UpdateAreaRequest._();
  UpdateAreaRequest createEmptyInstance() => create();
  static $pb.PbList<UpdateAreaRequest> createRepeated() => $pb.PbList<UpdateAreaRequest>();
  @$core.pragma('dart2js:noInline')
  static UpdateAreaRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateAreaRequest>(create);
  static UpdateAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => clearField(3);

  @$pb.TagNumber(4)
  AreaType get areaType => $_getN(3);
  @$pb.TagNumber(4)
  set areaType(AreaType v) { setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasAreaType() => $_has(3);
  @$pb.TagNumber(4)
  void clearAreaType() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get geometry => $_getSZ(4);
  @$pb.TagNumber(5)
  set geometry($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasGeometry() => $_has(4);
  @$pb.TagNumber(5)
  void clearGeometry() => clearField(5);

  @$pb.TagNumber(6)
  $6.Struct get extra => $_getN(5);
  @$pb.TagNumber(6)
  set extra($6.Struct v) { setField(6, v); }
  @$pb.TagNumber(6)
  $core.bool hasExtra() => $_has(5);
  @$pb.TagNumber(6)
  void clearExtra() => clearField(6);
  @$pb.TagNumber(6)
  $6.Struct ensureExtra() => $_ensure(5);
}

class UpdateAreaResponse extends $pb.GeneratedMessage {
  factory UpdateAreaResponse({
    AreaObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  UpdateAreaResponse._() : super();
  factory UpdateAreaResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UpdateAreaResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateAreaResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<AreaObject>(1, _omitFieldNames ? '' : 'data', subBuilder: AreaObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UpdateAreaResponse clone() => UpdateAreaResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UpdateAreaResponse copyWith(void Function(UpdateAreaResponse) updates) => super.copyWith((message) => updates(message as UpdateAreaResponse)) as UpdateAreaResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAreaResponse create() => UpdateAreaResponse._();
  UpdateAreaResponse createEmptyInstance() => create();
  static $pb.PbList<UpdateAreaResponse> createRepeated() => $pb.PbList<UpdateAreaResponse>();
  @$core.pragma('dart2js:noInline')
  static UpdateAreaResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateAreaResponse>(create);
  static UpdateAreaResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AreaObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(AreaObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  AreaObject ensureData() => $_ensure(0);
}

class DeleteAreaRequest extends $pb.GeneratedMessage {
  factory DeleteAreaRequest({
    $core.String? id,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    return $result;
  }
  DeleteAreaRequest._() : super();
  factory DeleteAreaRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DeleteAreaRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteAreaRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DeleteAreaRequest clone() => DeleteAreaRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DeleteAreaRequest copyWith(void Function(DeleteAreaRequest) updates) => super.copyWith((message) => updates(message as DeleteAreaRequest)) as DeleteAreaRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAreaRequest create() => DeleteAreaRequest._();
  DeleteAreaRequest createEmptyInstance() => create();
  static $pb.PbList<DeleteAreaRequest> createRepeated() => $pb.PbList<DeleteAreaRequest>();
  @$core.pragma('dart2js:noInline')
  static DeleteAreaRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteAreaRequest>(create);
  static DeleteAreaRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);
}

class SearchAreasRequest extends $pb.GeneratedMessage {
  factory SearchAreasRequest({
    $core.String? query,
    $core.String? ownerId,
    $core.int? limit,
  }) {
    final $result = create();
    if (query != null) {
      $result.query = query;
    }
    if (ownerId != null) {
      $result.ownerId = ownerId;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    return $result;
  }
  SearchAreasRequest._() : super();
  factory SearchAreasRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SearchAreasRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SearchAreasRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..aOS(2, _omitFieldNames ? '' : 'ownerId')
    ..a<$core.int>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SearchAreasRequest clone() => SearchAreasRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SearchAreasRequest copyWith(void Function(SearchAreasRequest) updates) => super.copyWith((message) => updates(message as SearchAreasRequest)) as SearchAreasRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchAreasRequest create() => SearchAreasRequest._();
  SearchAreasRequest createEmptyInstance() => create();
  static $pb.PbList<SearchAreasRequest> createRepeated() => $pb.PbList<SearchAreasRequest>();
  @$core.pragma('dart2js:noInline')
  static SearchAreasRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SearchAreasRequest>(create);
  static SearchAreasRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get ownerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ownerId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasOwnerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearOwnerId() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(3)
  set limit($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => clearField(3);
}

class SearchAreasResponse extends $pb.GeneratedMessage {
  factory SearchAreasResponse({
    $core.Iterable<AreaObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  SearchAreasResponse._() : super();
  factory SearchAreasResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SearchAreasResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SearchAreasResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<AreaObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: AreaObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SearchAreasResponse clone() => SearchAreasResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SearchAreasResponse copyWith(void Function(SearchAreasResponse) updates) => super.copyWith((message) => updates(message as SearchAreasResponse)) as SearchAreasResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchAreasResponse create() => SearchAreasResponse._();
  SearchAreasResponse createEmptyInstance() => create();
  static $pb.PbList<SearchAreasResponse> createRepeated() => $pb.PbList<SearchAreasResponse>();
  @$core.pragma('dart2js:noInline')
  static SearchAreasResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SearchAreasResponse>(create);
  static SearchAreasResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<AreaObject> get data => $_getList(0);
}

class CreateRouteRequest extends $pb.GeneratedMessage {
  factory CreateRouteRequest({
    RouteObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  CreateRouteRequest._() : super();
  factory CreateRouteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateRouteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateRouteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data', subBuilder: RouteObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateRouteRequest clone() => CreateRouteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateRouteRequest copyWith(void Function(CreateRouteRequest) updates) => super.copyWith((message) => updates(message as CreateRouteRequest)) as CreateRouteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRouteRequest create() => CreateRouteRequest._();
  CreateRouteRequest createEmptyInstance() => create();
  static $pb.PbList<CreateRouteRequest> createRepeated() => $pb.PbList<CreateRouteRequest>();
  @$core.pragma('dart2js:noInline')
  static CreateRouteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateRouteRequest>(create);
  static CreateRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class CreateRouteResponse extends $pb.GeneratedMessage {
  factory CreateRouteResponse({
    RouteObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  CreateRouteResponse._() : super();
  factory CreateRouteResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateRouteResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateRouteResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data', subBuilder: RouteObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateRouteResponse clone() => CreateRouteResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateRouteResponse copyWith(void Function(CreateRouteResponse) updates) => super.copyWith((message) => updates(message as CreateRouteResponse)) as CreateRouteResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRouteResponse create() => CreateRouteResponse._();
  CreateRouteResponse createEmptyInstance() => create();
  static $pb.PbList<CreateRouteResponse> createRepeated() => $pb.PbList<CreateRouteResponse>();
  @$core.pragma('dart2js:noInline')
  static CreateRouteResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateRouteResponse>(create);
  static CreateRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class GetRouteRequest extends $pb.GeneratedMessage {
  factory GetRouteRequest({
    $core.String? id,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    return $result;
  }
  GetRouteRequest._() : super();
  factory GetRouteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRouteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetRouteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRouteRequest clone() => GetRouteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRouteRequest copyWith(void Function(GetRouteRequest) updates) => super.copyWith((message) => updates(message as GetRouteRequest)) as GetRouteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRouteRequest create() => GetRouteRequest._();
  GetRouteRequest createEmptyInstance() => create();
  static $pb.PbList<GetRouteRequest> createRepeated() => $pb.PbList<GetRouteRequest>();
  @$core.pragma('dart2js:noInline')
  static GetRouteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRouteRequest>(create);
  static GetRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);
}

class GetRouteResponse extends $pb.GeneratedMessage {
  factory GetRouteResponse({
    RouteObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  GetRouteResponse._() : super();
  factory GetRouteResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRouteResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetRouteResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data', subBuilder: RouteObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRouteResponse clone() => GetRouteResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRouteResponse copyWith(void Function(GetRouteResponse) updates) => super.copyWith((message) => updates(message as GetRouteResponse)) as GetRouteResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRouteResponse create() => GetRouteResponse._();
  GetRouteResponse createEmptyInstance() => create();
  static $pb.PbList<GetRouteResponse> createRepeated() => $pb.PbList<GetRouteResponse>();
  @$core.pragma('dart2js:noInline')
  static GetRouteResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRouteResponse>(create);
  static GetRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
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
    $6.Struct? extra,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    if (name != null) {
      $result.name = name;
    }
    if (description != null) {
      $result.description = description;
    }
    if (geometry != null) {
      $result.geometry = geometry;
    }
    if (deviationThresholdM != null) {
      $result.deviationThresholdM = deviationThresholdM;
    }
    if (deviationConsecutiveCount != null) {
      $result.deviationConsecutiveCount = deviationConsecutiveCount;
    }
    if (deviationCooldownSec != null) {
      $result.deviationCooldownSec = deviationCooldownSec;
    }
    if (extra != null) {
      $result.extra = extra;
    }
    return $result;
  }
  UpdateRouteRequest._() : super();
  factory UpdateRouteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UpdateRouteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateRouteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'description')
    ..aOS(4, _omitFieldNames ? '' : 'geometry')
    ..a<$core.double>(5, _omitFieldNames ? '' : 'deviationThresholdM', $pb.PbFieldType.OD)
    ..a<$core.int>(6, _omitFieldNames ? '' : 'deviationConsecutiveCount', $pb.PbFieldType.O3)
    ..a<$core.int>(7, _omitFieldNames ? '' : 'deviationCooldownSec', $pb.PbFieldType.O3)
    ..aOM<$6.Struct>(8, _omitFieldNames ? '' : 'extra', subBuilder: $6.Struct.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UpdateRouteRequest clone() => UpdateRouteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UpdateRouteRequest copyWith(void Function(UpdateRouteRequest) updates) => super.copyWith((message) => updates(message as UpdateRouteRequest)) as UpdateRouteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRouteRequest create() => UpdateRouteRequest._();
  UpdateRouteRequest createEmptyInstance() => create();
  static $pb.PbList<UpdateRouteRequest> createRepeated() => $pb.PbList<UpdateRouteRequest>();
  @$core.pragma('dart2js:noInline')
  static UpdateRouteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateRouteRequest>(create);
  static UpdateRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(3)
  void clearDescription() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get geometry => $_getSZ(3);
  @$pb.TagNumber(4)
  set geometry($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasGeometry() => $_has(3);
  @$pb.TagNumber(4)
  void clearGeometry() => clearField(4);

  @$pb.TagNumber(5)
  $core.double get deviationThresholdM => $_getN(4);
  @$pb.TagNumber(5)
  set deviationThresholdM($core.double v) { $_setDouble(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasDeviationThresholdM() => $_has(4);
  @$pb.TagNumber(5)
  void clearDeviationThresholdM() => clearField(5);

  @$pb.TagNumber(6)
  $core.int get deviationConsecutiveCount => $_getIZ(5);
  @$pb.TagNumber(6)
  set deviationConsecutiveCount($core.int v) { $_setSignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasDeviationConsecutiveCount() => $_has(5);
  @$pb.TagNumber(6)
  void clearDeviationConsecutiveCount() => clearField(6);

  @$pb.TagNumber(7)
  $core.int get deviationCooldownSec => $_getIZ(6);
  @$pb.TagNumber(7)
  set deviationCooldownSec($core.int v) { $_setSignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasDeviationCooldownSec() => $_has(6);
  @$pb.TagNumber(7)
  void clearDeviationCooldownSec() => clearField(7);

  @$pb.TagNumber(8)
  $6.Struct get extra => $_getN(7);
  @$pb.TagNumber(8)
  set extra($6.Struct v) { setField(8, v); }
  @$pb.TagNumber(8)
  $core.bool hasExtra() => $_has(7);
  @$pb.TagNumber(8)
  void clearExtra() => clearField(8);
  @$pb.TagNumber(8)
  $6.Struct ensureExtra() => $_ensure(7);
}

class UpdateRouteResponse extends $pb.GeneratedMessage {
  factory UpdateRouteResponse({
    RouteObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  UpdateRouteResponse._() : super();
  factory UpdateRouteResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UpdateRouteResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateRouteResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<RouteObject>(1, _omitFieldNames ? '' : 'data', subBuilder: RouteObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UpdateRouteResponse clone() => UpdateRouteResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UpdateRouteResponse copyWith(void Function(UpdateRouteResponse) updates) => super.copyWith((message) => updates(message as UpdateRouteResponse)) as UpdateRouteResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRouteResponse create() => UpdateRouteResponse._();
  UpdateRouteResponse createEmptyInstance() => create();
  static $pb.PbList<UpdateRouteResponse> createRepeated() => $pb.PbList<UpdateRouteResponse>();
  @$core.pragma('dart2js:noInline')
  static UpdateRouteResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateRouteResponse>(create);
  static UpdateRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  RouteObject ensureData() => $_ensure(0);
}

class DeleteRouteRequest extends $pb.GeneratedMessage {
  factory DeleteRouteRequest({
    $core.String? id,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    return $result;
  }
  DeleteRouteRequest._() : super();
  factory DeleteRouteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory DeleteRouteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteRouteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  DeleteRouteRequest clone() => DeleteRouteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  DeleteRouteRequest copyWith(void Function(DeleteRouteRequest) updates) => super.copyWith((message) => updates(message as DeleteRouteRequest)) as DeleteRouteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRouteRequest create() => DeleteRouteRequest._();
  DeleteRouteRequest createEmptyInstance() => create();
  static $pb.PbList<DeleteRouteRequest> createRepeated() => $pb.PbList<DeleteRouteRequest>();
  @$core.pragma('dart2js:noInline')
  static DeleteRouteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteRouteRequest>(create);
  static DeleteRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);
}

class SearchRoutesRequest extends $pb.GeneratedMessage {
  factory SearchRoutesRequest({
    $core.String? ownerId,
    $core.int? limit,
  }) {
    final $result = create();
    if (ownerId != null) {
      $result.ownerId = ownerId;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    return $result;
  }
  SearchRoutesRequest._() : super();
  factory SearchRoutesRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SearchRoutesRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SearchRoutesRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'ownerId')
    ..a<$core.int>(2, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SearchRoutesRequest clone() => SearchRoutesRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SearchRoutesRequest copyWith(void Function(SearchRoutesRequest) updates) => super.copyWith((message) => updates(message as SearchRoutesRequest)) as SearchRoutesRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchRoutesRequest create() => SearchRoutesRequest._();
  SearchRoutesRequest createEmptyInstance() => create();
  static $pb.PbList<SearchRoutesRequest> createRepeated() => $pb.PbList<SearchRoutesRequest>();
  @$core.pragma('dart2js:noInline')
  static SearchRoutesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SearchRoutesRequest>(create);
  static SearchRoutesRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get ownerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set ownerId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasOwnerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearOwnerId() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(2)
  set limit($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimit() => clearField(2);
}

class SearchRoutesResponse extends $pb.GeneratedMessage {
  factory SearchRoutesResponse({
    $core.Iterable<RouteObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  SearchRoutesResponse._() : super();
  factory SearchRoutesResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SearchRoutesResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SearchRoutesResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<RouteObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: RouteObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SearchRoutesResponse clone() => SearchRoutesResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SearchRoutesResponse copyWith(void Function(SearchRoutesResponse) updates) => super.copyWith((message) => updates(message as SearchRoutesResponse)) as SearchRoutesResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchRoutesResponse create() => SearchRoutesResponse._();
  SearchRoutesResponse createEmptyInstance() => create();
  static $pb.PbList<SearchRoutesResponse> createRepeated() => $pb.PbList<SearchRoutesResponse>();
  @$core.pragma('dart2js:noInline')
  static SearchRoutesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SearchRoutesResponse>(create);
  static SearchRoutesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<RouteObject> get data => $_getList(0);
}

class AssignRouteRequest extends $pb.GeneratedMessage {
  factory AssignRouteRequest({
    $core.String? subjectId,
    $core.String? routeId,
    $2.Timestamp? validFrom,
    $2.Timestamp? validUntil,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (routeId != null) {
      $result.routeId = routeId;
    }
    if (validFrom != null) {
      $result.validFrom = validFrom;
    }
    if (validUntil != null) {
      $result.validUntil = validUntil;
    }
    return $result;
  }
  AssignRouteRequest._() : super();
  factory AssignRouteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AssignRouteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AssignRouteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOS(2, _omitFieldNames ? '' : 'routeId')
    ..aOM<$2.Timestamp>(3, _omitFieldNames ? '' : 'validFrom', subBuilder: $2.Timestamp.create)
    ..aOM<$2.Timestamp>(4, _omitFieldNames ? '' : 'validUntil', subBuilder: $2.Timestamp.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AssignRouteRequest clone() => AssignRouteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AssignRouteRequest copyWith(void Function(AssignRouteRequest) updates) => super.copyWith((message) => updates(message as AssignRouteRequest)) as AssignRouteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssignRouteRequest create() => AssignRouteRequest._();
  AssignRouteRequest createEmptyInstance() => create();
  static $pb.PbList<AssignRouteRequest> createRepeated() => $pb.PbList<AssignRouteRequest>();
  @$core.pragma('dart2js:noInline')
  static AssignRouteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AssignRouteRequest>(create);
  static AssignRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get routeId => $_getSZ(1);
  @$pb.TagNumber(2)
  set routeId($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRouteId() => $_has(1);
  @$pb.TagNumber(2)
  void clearRouteId() => clearField(2);

  @$pb.TagNumber(3)
  $2.Timestamp get validFrom => $_getN(2);
  @$pb.TagNumber(3)
  set validFrom($2.Timestamp v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasValidFrom() => $_has(2);
  @$pb.TagNumber(3)
  void clearValidFrom() => clearField(3);
  @$pb.TagNumber(3)
  $2.Timestamp ensureValidFrom() => $_ensure(2);

  @$pb.TagNumber(4)
  $2.Timestamp get validUntil => $_getN(3);
  @$pb.TagNumber(4)
  set validUntil($2.Timestamp v) { setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasValidUntil() => $_has(3);
  @$pb.TagNumber(4)
  void clearValidUntil() => clearField(4);
  @$pb.TagNumber(4)
  $2.Timestamp ensureValidUntil() => $_ensure(3);
}

class AssignRouteResponse extends $pb.GeneratedMessage {
  factory AssignRouteResponse({
    RouteAssignmentObject? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data = data;
    }
    return $result;
  }
  AssignRouteResponse._() : super();
  factory AssignRouteResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AssignRouteResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AssignRouteResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOM<RouteAssignmentObject>(1, _omitFieldNames ? '' : 'data', subBuilder: RouteAssignmentObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AssignRouteResponse clone() => AssignRouteResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AssignRouteResponse copyWith(void Function(AssignRouteResponse) updates) => super.copyWith((message) => updates(message as AssignRouteResponse)) as AssignRouteResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssignRouteResponse create() => AssignRouteResponse._();
  AssignRouteResponse createEmptyInstance() => create();
  static $pb.PbList<AssignRouteResponse> createRepeated() => $pb.PbList<AssignRouteResponse>();
  @$core.pragma('dart2js:noInline')
  static AssignRouteResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AssignRouteResponse>(create);
  static AssignRouteResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RouteAssignmentObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RouteAssignmentObject v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
  @$pb.TagNumber(1)
  RouteAssignmentObject ensureData() => $_ensure(0);
}

class UnassignRouteRequest extends $pb.GeneratedMessage {
  factory UnassignRouteRequest({
    $core.String? id,
  }) {
    final $result = create();
    if (id != null) {
      $result.id = id;
    }
    return $result;
  }
  UnassignRouteRequest._() : super();
  factory UnassignRouteRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnassignRouteRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UnassignRouteRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnassignRouteRequest clone() => UnassignRouteRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnassignRouteRequest copyWith(void Function(UnassignRouteRequest) updates) => super.copyWith((message) => updates(message as UnassignRouteRequest)) as UnassignRouteRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnassignRouteRequest create() => UnassignRouteRequest._();
  UnassignRouteRequest createEmptyInstance() => create();
  static $pb.PbList<UnassignRouteRequest> createRepeated() => $pb.PbList<UnassignRouteRequest>();
  @$core.pragma('dart2js:noInline')
  static UnassignRouteRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnassignRouteRequest>(create);
  static UnassignRouteRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);
}

class GetSubjectRouteAssignmentsRequest extends $pb.GeneratedMessage {
  factory GetSubjectRouteAssignmentsRequest({
    $core.String? subjectId,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    return $result;
  }
  GetSubjectRouteAssignmentsRequest._() : super();
  factory GetSubjectRouteAssignmentsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetSubjectRouteAssignmentsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetSubjectRouteAssignmentsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetSubjectRouteAssignmentsRequest clone() => GetSubjectRouteAssignmentsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetSubjectRouteAssignmentsRequest copyWith(void Function(GetSubjectRouteAssignmentsRequest) updates) => super.copyWith((message) => updates(message as GetSubjectRouteAssignmentsRequest)) as GetSubjectRouteAssignmentsRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSubjectRouteAssignmentsRequest create() => GetSubjectRouteAssignmentsRequest._();
  GetSubjectRouteAssignmentsRequest createEmptyInstance() => create();
  static $pb.PbList<GetSubjectRouteAssignmentsRequest> createRepeated() => $pb.PbList<GetSubjectRouteAssignmentsRequest>();
  @$core.pragma('dart2js:noInline')
  static GetSubjectRouteAssignmentsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetSubjectRouteAssignmentsRequest>(create);
  static GetSubjectRouteAssignmentsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);
}

class GetSubjectRouteAssignmentsResponse extends $pb.GeneratedMessage {
  factory GetSubjectRouteAssignmentsResponse({
    $core.Iterable<RouteAssignmentObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  GetSubjectRouteAssignmentsResponse._() : super();
  factory GetSubjectRouteAssignmentsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetSubjectRouteAssignmentsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetSubjectRouteAssignmentsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<RouteAssignmentObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: RouteAssignmentObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetSubjectRouteAssignmentsResponse clone() => GetSubjectRouteAssignmentsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetSubjectRouteAssignmentsResponse copyWith(void Function(GetSubjectRouteAssignmentsResponse) updates) => super.copyWith((message) => updates(message as GetSubjectRouteAssignmentsResponse)) as GetSubjectRouteAssignmentsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSubjectRouteAssignmentsResponse create() => GetSubjectRouteAssignmentsResponse._();
  GetSubjectRouteAssignmentsResponse createEmptyInstance() => create();
  static $pb.PbList<GetSubjectRouteAssignmentsResponse> createRepeated() => $pb.PbList<GetSubjectRouteAssignmentsResponse>();
  @$core.pragma('dart2js:noInline')
  static GetSubjectRouteAssignmentsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetSubjectRouteAssignmentsResponse>(create);
  static GetSubjectRouteAssignmentsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<RouteAssignmentObject> get data => $_getList(0);
}

class GetTrackRequest extends $pb.GeneratedMessage {
  factory GetTrackRequest({
    $core.String? subjectId,
    $2.Timestamp? from,
    $2.Timestamp? to,
    $core.int? limit,
    $core.int? offset,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (from != null) {
      $result.from = from;
    }
    if (to != null) {
      $result.to = to;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    if (offset != null) {
      $result.offset = offset;
    }
    return $result;
  }
  GetTrackRequest._() : super();
  factory GetTrackRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTrackRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTrackRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$2.Timestamp>(2, _omitFieldNames ? '' : 'from', subBuilder: $2.Timestamp.create)
    ..aOM<$2.Timestamp>(3, _omitFieldNames ? '' : 'to', subBuilder: $2.Timestamp.create)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.O3)
    ..a<$core.int>(5, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTrackRequest clone() => GetTrackRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTrackRequest copyWith(void Function(GetTrackRequest) updates) => super.copyWith((message) => updates(message as GetTrackRequest)) as GetTrackRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrackRequest create() => GetTrackRequest._();
  GetTrackRequest createEmptyInstance() => create();
  static $pb.PbList<GetTrackRequest> createRepeated() => $pb.PbList<GetTrackRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTrackRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTrackRequest>(create);
  static GetTrackRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $2.Timestamp get from => $_getN(1);
  @$pb.TagNumber(2)
  set from($2.Timestamp v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasFrom() => $_has(1);
  @$pb.TagNumber(2)
  void clearFrom() => clearField(2);
  @$pb.TagNumber(2)
  $2.Timestamp ensureFrom() => $_ensure(1);

  @$pb.TagNumber(3)
  $2.Timestamp get to => $_getN(2);
  @$pb.TagNumber(3)
  set to($2.Timestamp v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasTo() => $_has(2);
  @$pb.TagNumber(3)
  void clearTo() => clearField(3);
  @$pb.TagNumber(3)
  $2.Timestamp ensureTo() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(4)
  set limit($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(4)
  void clearLimit() => clearField(4);

  @$pb.TagNumber(5)
  $core.int get offset => $_getIZ(4);
  @$pb.TagNumber(5)
  set offset($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(5)
  void clearOffset() => clearField(5);
}

class GetTrackResponse extends $pb.GeneratedMessage {
  factory GetTrackResponse({
    $core.Iterable<LocationPointObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  GetTrackResponse._() : super();
  factory GetTrackResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTrackResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTrackResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<LocationPointObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: LocationPointObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTrackResponse clone() => GetTrackResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTrackResponse copyWith(void Function(GetTrackResponse) updates) => super.copyWith((message) => updates(message as GetTrackResponse)) as GetTrackResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrackResponse create() => GetTrackResponse._();
  GetTrackResponse createEmptyInstance() => create();
  static $pb.PbList<GetTrackResponse> createRepeated() => $pb.PbList<GetTrackResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTrackResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTrackResponse>(create);
  static GetTrackResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<LocationPointObject> get data => $_getList(0);
}

class GetSubjectEventsRequest extends $pb.GeneratedMessage {
  factory GetSubjectEventsRequest({
    $core.String? subjectId,
    $2.Timestamp? from,
    $2.Timestamp? to,
    $core.int? limit,
    $core.int? offset,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (from != null) {
      $result.from = from;
    }
    if (to != null) {
      $result.to = to;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    if (offset != null) {
      $result.offset = offset;
    }
    return $result;
  }
  GetSubjectEventsRequest._() : super();
  factory GetSubjectEventsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetSubjectEventsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetSubjectEventsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..aOM<$2.Timestamp>(2, _omitFieldNames ? '' : 'from', subBuilder: $2.Timestamp.create)
    ..aOM<$2.Timestamp>(3, _omitFieldNames ? '' : 'to', subBuilder: $2.Timestamp.create)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.O3)
    ..a<$core.int>(5, _omitFieldNames ? '' : 'offset', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetSubjectEventsRequest clone() => GetSubjectEventsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetSubjectEventsRequest copyWith(void Function(GetSubjectEventsRequest) updates) => super.copyWith((message) => updates(message as GetSubjectEventsRequest)) as GetSubjectEventsRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSubjectEventsRequest create() => GetSubjectEventsRequest._();
  GetSubjectEventsRequest createEmptyInstance() => create();
  static $pb.PbList<GetSubjectEventsRequest> createRepeated() => $pb.PbList<GetSubjectEventsRequest>();
  @$core.pragma('dart2js:noInline')
  static GetSubjectEventsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetSubjectEventsRequest>(create);
  static GetSubjectEventsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $2.Timestamp get from => $_getN(1);
  @$pb.TagNumber(2)
  set from($2.Timestamp v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasFrom() => $_has(1);
  @$pb.TagNumber(2)
  void clearFrom() => clearField(2);
  @$pb.TagNumber(2)
  $2.Timestamp ensureFrom() => $_ensure(1);

  @$pb.TagNumber(3)
  $2.Timestamp get to => $_getN(2);
  @$pb.TagNumber(3)
  set to($2.Timestamp v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasTo() => $_has(2);
  @$pb.TagNumber(3)
  void clearTo() => clearField(3);
  @$pb.TagNumber(3)
  $2.Timestamp ensureTo() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(4)
  set limit($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(4)
  void clearLimit() => clearField(4);

  @$pb.TagNumber(5)
  $core.int get offset => $_getIZ(4);
  @$pb.TagNumber(5)
  set offset($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasOffset() => $_has(4);
  @$pb.TagNumber(5)
  void clearOffset() => clearField(5);
}

class GetSubjectEventsResponse extends $pb.GeneratedMessage {
  factory GetSubjectEventsResponse({
    $core.Iterable<GeoEventObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  GetSubjectEventsResponse._() : super();
  factory GetSubjectEventsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetSubjectEventsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetSubjectEventsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<GeoEventObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: GeoEventObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetSubjectEventsResponse clone() => GetSubjectEventsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetSubjectEventsResponse copyWith(void Function(GetSubjectEventsResponse) updates) => super.copyWith((message) => updates(message as GetSubjectEventsResponse)) as GetSubjectEventsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSubjectEventsResponse create() => GetSubjectEventsResponse._();
  GetSubjectEventsResponse createEmptyInstance() => create();
  static $pb.PbList<GetSubjectEventsResponse> createRepeated() => $pb.PbList<GetSubjectEventsResponse>();
  @$core.pragma('dart2js:noInline')
  static GetSubjectEventsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetSubjectEventsResponse>(create);
  static GetSubjectEventsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<GeoEventObject> get data => $_getList(0);
}

class GetAreaSubjectsRequest extends $pb.GeneratedMessage {
  factory GetAreaSubjectsRequest({
    $core.String? areaId,
  }) {
    final $result = create();
    if (areaId != null) {
      $result.areaId = areaId;
    }
    return $result;
  }
  GetAreaSubjectsRequest._() : super();
  factory GetAreaSubjectsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAreaSubjectsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAreaSubjectsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'areaId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAreaSubjectsRequest clone() => GetAreaSubjectsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAreaSubjectsRequest copyWith(void Function(GetAreaSubjectsRequest) updates) => super.copyWith((message) => updates(message as GetAreaSubjectsRequest)) as GetAreaSubjectsRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAreaSubjectsRequest create() => GetAreaSubjectsRequest._();
  GetAreaSubjectsRequest createEmptyInstance() => create();
  static $pb.PbList<GetAreaSubjectsRequest> createRepeated() => $pb.PbList<GetAreaSubjectsRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAreaSubjectsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAreaSubjectsRequest>(create);
  static GetAreaSubjectsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get areaId => $_getSZ(0);
  @$pb.TagNumber(1)
  set areaId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAreaId() => $_has(0);
  @$pb.TagNumber(1)
  void clearAreaId() => clearField(1);
}

class GetAreaSubjectsResponse extends $pb.GeneratedMessage {
  factory GetAreaSubjectsResponse({
    $core.Iterable<AreaSubjectObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  GetAreaSubjectsResponse._() : super();
  factory GetAreaSubjectsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAreaSubjectsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAreaSubjectsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<AreaSubjectObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: AreaSubjectObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAreaSubjectsResponse clone() => GetAreaSubjectsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAreaSubjectsResponse copyWith(void Function(GetAreaSubjectsResponse) updates) => super.copyWith((message) => updates(message as GetAreaSubjectsResponse)) as GetAreaSubjectsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAreaSubjectsResponse create() => GetAreaSubjectsResponse._();
  GetAreaSubjectsResponse createEmptyInstance() => create();
  static $pb.PbList<GetAreaSubjectsResponse> createRepeated() => $pb.PbList<GetAreaSubjectsResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAreaSubjectsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAreaSubjectsResponse>(create);
  static GetAreaSubjectsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<AreaSubjectObject> get data => $_getList(0);
}

class GetNearbySubjectsRequest extends $pb.GeneratedMessage {
  factory GetNearbySubjectsRequest({
    $core.String? subjectId,
    $core.double? radiusMeters,
    $core.int? limit,
  }) {
    final $result = create();
    if (subjectId != null) {
      $result.subjectId = subjectId;
    }
    if (radiusMeters != null) {
      $result.radiusMeters = radiusMeters;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    return $result;
  }
  GetNearbySubjectsRequest._() : super();
  factory GetNearbySubjectsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNearbySubjectsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNearbySubjectsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'subjectId')
    ..a<$core.double>(2, _omitFieldNames ? '' : 'radiusMeters', $pb.PbFieldType.OD)
    ..a<$core.int>(3, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNearbySubjectsRequest clone() => GetNearbySubjectsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNearbySubjectsRequest copyWith(void Function(GetNearbySubjectsRequest) updates) => super.copyWith((message) => updates(message as GetNearbySubjectsRequest)) as GetNearbySubjectsRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNearbySubjectsRequest create() => GetNearbySubjectsRequest._();
  GetNearbySubjectsRequest createEmptyInstance() => create();
  static $pb.PbList<GetNearbySubjectsRequest> createRepeated() => $pb.PbList<GetNearbySubjectsRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNearbySubjectsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNearbySubjectsRequest>(create);
  static GetNearbySubjectsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get subjectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set subjectId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSubjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSubjectId() => clearField(1);

  @$pb.TagNumber(2)
  $core.double get radiusMeters => $_getN(1);
  @$pb.TagNumber(2)
  set radiusMeters($core.double v) { $_setDouble(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRadiusMeters() => $_has(1);
  @$pb.TagNumber(2)
  void clearRadiusMeters() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(3)
  set limit($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => clearField(3);
}

class GetNearbySubjectsResponse extends $pb.GeneratedMessage {
  factory GetNearbySubjectsResponse({
    $core.Iterable<NearbySubjectObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  GetNearbySubjectsResponse._() : super();
  factory GetNearbySubjectsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNearbySubjectsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNearbySubjectsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<NearbySubjectObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: NearbySubjectObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNearbySubjectsResponse clone() => GetNearbySubjectsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNearbySubjectsResponse copyWith(void Function(GetNearbySubjectsResponse) updates) => super.copyWith((message) => updates(message as GetNearbySubjectsResponse)) as GetNearbySubjectsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNearbySubjectsResponse create() => GetNearbySubjectsResponse._();
  GetNearbySubjectsResponse createEmptyInstance() => create();
  static $pb.PbList<GetNearbySubjectsResponse> createRepeated() => $pb.PbList<GetNearbySubjectsResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNearbySubjectsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNearbySubjectsResponse>(create);
  static GetNearbySubjectsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<NearbySubjectObject> get data => $_getList(0);
}

class GetNearbyAreasRequest extends $pb.GeneratedMessage {
  factory GetNearbyAreasRequest({
    $core.double? latitude,
    $core.double? longitude,
    $core.double? radiusMeters,
    $core.int? limit,
  }) {
    final $result = create();
    if (latitude != null) {
      $result.latitude = latitude;
    }
    if (longitude != null) {
      $result.longitude = longitude;
    }
    if (radiusMeters != null) {
      $result.radiusMeters = radiusMeters;
    }
    if (limit != null) {
      $result.limit = limit;
    }
    return $result;
  }
  GetNearbyAreasRequest._() : super();
  factory GetNearbyAreasRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNearbyAreasRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNearbyAreasRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..a<$core.double>(1, _omitFieldNames ? '' : 'latitude', $pb.PbFieldType.OD)
    ..a<$core.double>(2, _omitFieldNames ? '' : 'longitude', $pb.PbFieldType.OD)
    ..a<$core.double>(3, _omitFieldNames ? '' : 'radiusMeters', $pb.PbFieldType.OD)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'limit', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNearbyAreasRequest clone() => GetNearbyAreasRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNearbyAreasRequest copyWith(void Function(GetNearbyAreasRequest) updates) => super.copyWith((message) => updates(message as GetNearbyAreasRequest)) as GetNearbyAreasRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNearbyAreasRequest create() => GetNearbyAreasRequest._();
  GetNearbyAreasRequest createEmptyInstance() => create();
  static $pb.PbList<GetNearbyAreasRequest> createRepeated() => $pb.PbList<GetNearbyAreasRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNearbyAreasRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNearbyAreasRequest>(create);
  static GetNearbyAreasRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.double get latitude => $_getN(0);
  @$pb.TagNumber(1)
  set latitude($core.double v) { $_setDouble(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLatitude() => $_has(0);
  @$pb.TagNumber(1)
  void clearLatitude() => clearField(1);

  @$pb.TagNumber(2)
  $core.double get longitude => $_getN(1);
  @$pb.TagNumber(2)
  set longitude($core.double v) { $_setDouble(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLongitude() => $_has(1);
  @$pb.TagNumber(2)
  void clearLongitude() => clearField(2);

  @$pb.TagNumber(3)
  $core.double get radiusMeters => $_getN(2);
  @$pb.TagNumber(3)
  set radiusMeters($core.double v) { $_setDouble(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasRadiusMeters() => $_has(2);
  @$pb.TagNumber(3)
  void clearRadiusMeters() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(4)
  set limit($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(4)
  void clearLimit() => clearField(4);
}

class GetNearbyAreasResponse extends $pb.GeneratedMessage {
  factory GetNearbyAreasResponse({
    $core.Iterable<NearbyAreaObject>? data,
  }) {
    final $result = create();
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  GetNearbyAreasResponse._() : super();
  factory GetNearbyAreasResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNearbyAreasResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNearbyAreasResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'geolocation.v1'), createEmptyInstance: create)
    ..pc<NearbyAreaObject>(1, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: NearbyAreaObject.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNearbyAreasResponse clone() => GetNearbyAreasResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNearbyAreasResponse copyWith(void Function(GetNearbyAreasResponse) updates) => super.copyWith((message) => updates(message as GetNearbyAreasResponse)) as GetNearbyAreasResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNearbyAreasResponse create() => GetNearbyAreasResponse._();
  GetNearbyAreasResponse createEmptyInstance() => create();
  static $pb.PbList<GetNearbyAreasResponse> createRepeated() => $pb.PbList<GetNearbyAreasResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNearbyAreasResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNearbyAreasResponse>(create);
  static GetNearbyAreasResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<NearbyAreaObject> get data => $_getList(0);
}

class GeolocationServiceApi {
  $pb.RpcClient _client;
  GeolocationServiceApi(this._client);

  $async.Future<IngestLocationsResponse> ingestLocations($pb.ClientContext? ctx, IngestLocationsRequest request) =>
    _client.invoke<IngestLocationsResponse>(ctx, 'GeolocationService', 'IngestLocations', request, IngestLocationsResponse())
  ;
  $async.Future<CreateAreaResponse> createArea($pb.ClientContext? ctx, CreateAreaRequest request) =>
    _client.invoke<CreateAreaResponse>(ctx, 'GeolocationService', 'CreateArea', request, CreateAreaResponse())
  ;
  $async.Future<GetAreaResponse> getArea($pb.ClientContext? ctx, GetAreaRequest request) =>
    _client.invoke<GetAreaResponse>(ctx, 'GeolocationService', 'GetArea', request, GetAreaResponse())
  ;
  $async.Future<UpdateAreaResponse> updateArea($pb.ClientContext? ctx, UpdateAreaRequest request) =>
    _client.invoke<UpdateAreaResponse>(ctx, 'GeolocationService', 'UpdateArea', request, UpdateAreaResponse())
  ;
  $async.Future<$8.Empty> deleteArea($pb.ClientContext? ctx, DeleteAreaRequest request) =>
    _client.invoke<$8.Empty>(ctx, 'GeolocationService', 'DeleteArea', request, $8.Empty())
  ;
  $async.Future<SearchAreasResponse> searchAreas($pb.ClientContext? ctx, SearchAreasRequest request) =>
    _client.invoke<SearchAreasResponse>(ctx, 'GeolocationService', 'SearchAreas', request, SearchAreasResponse())
  ;
  $async.Future<CreateRouteResponse> createRoute($pb.ClientContext? ctx, CreateRouteRequest request) =>
    _client.invoke<CreateRouteResponse>(ctx, 'GeolocationService', 'CreateRoute', request, CreateRouteResponse())
  ;
  $async.Future<GetRouteResponse> getRoute($pb.ClientContext? ctx, GetRouteRequest request) =>
    _client.invoke<GetRouteResponse>(ctx, 'GeolocationService', 'GetRoute', request, GetRouteResponse())
  ;
  $async.Future<UpdateRouteResponse> updateRoute($pb.ClientContext? ctx, UpdateRouteRequest request) =>
    _client.invoke<UpdateRouteResponse>(ctx, 'GeolocationService', 'UpdateRoute', request, UpdateRouteResponse())
  ;
  $async.Future<$8.Empty> deleteRoute($pb.ClientContext? ctx, DeleteRouteRequest request) =>
    _client.invoke<$8.Empty>(ctx, 'GeolocationService', 'DeleteRoute', request, $8.Empty())
  ;
  $async.Future<SearchRoutesResponse> searchRoutes($pb.ClientContext? ctx, SearchRoutesRequest request) =>
    _client.invoke<SearchRoutesResponse>(ctx, 'GeolocationService', 'SearchRoutes', request, SearchRoutesResponse())
  ;
  $async.Future<AssignRouteResponse> assignRoute($pb.ClientContext? ctx, AssignRouteRequest request) =>
    _client.invoke<AssignRouteResponse>(ctx, 'GeolocationService', 'AssignRoute', request, AssignRouteResponse())
  ;
  $async.Future<$8.Empty> unassignRoute($pb.ClientContext? ctx, UnassignRouteRequest request) =>
    _client.invoke<$8.Empty>(ctx, 'GeolocationService', 'UnassignRoute', request, $8.Empty())
  ;
  $async.Future<GetSubjectRouteAssignmentsResponse> getSubjectRouteAssignments($pb.ClientContext? ctx, GetSubjectRouteAssignmentsRequest request) =>
    _client.invoke<GetSubjectRouteAssignmentsResponse>(ctx, 'GeolocationService', 'GetSubjectRouteAssignments', request, GetSubjectRouteAssignmentsResponse())
  ;
  $async.Future<GetTrackResponse> getTrack($pb.ClientContext? ctx, GetTrackRequest request) =>
    _client.invoke<GetTrackResponse>(ctx, 'GeolocationService', 'GetTrack', request, GetTrackResponse())
  ;
  $async.Future<GetSubjectEventsResponse> getSubjectEvents($pb.ClientContext? ctx, GetSubjectEventsRequest request) =>
    _client.invoke<GetSubjectEventsResponse>(ctx, 'GeolocationService', 'GetSubjectEvents', request, GetSubjectEventsResponse())
  ;
  $async.Future<GetAreaSubjectsResponse> getAreaSubjects($pb.ClientContext? ctx, GetAreaSubjectsRequest request) =>
    _client.invoke<GetAreaSubjectsResponse>(ctx, 'GeolocationService', 'GetAreaSubjects', request, GetAreaSubjectsResponse())
  ;
  $async.Future<GetNearbySubjectsResponse> getNearbySubjects($pb.ClientContext? ctx, GetNearbySubjectsRequest request) =>
    _client.invoke<GetNearbySubjectsResponse>(ctx, 'GeolocationService', 'GetNearbySubjects', request, GetNearbySubjectsResponse())
  ;
  $async.Future<GetNearbyAreasResponse> getNearbyAreas($pb.ClientContext? ctx, GetNearbyAreasRequest request) =>
    _client.invoke<GetNearbyAreasResponse>(ctx, 'GeolocationService', 'GetNearbyAreas', request, GetNearbyAreasResponse())
  ;
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
