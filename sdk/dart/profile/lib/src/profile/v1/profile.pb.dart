// This is a generated file - do not edit.
//
// Generated from profile/v1/profile.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;
import 'package:protobuf/well_known_types/google/protobuf/struct.pb.dart' as $0;
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

import '../../common/v1/common.pbenum.dart' as $2;
import 'profile.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'profile.pbenum.dart';

/// ContactObject represents a contact method (email or phone).
class ContactObject extends $pb.GeneratedMessage {
  factory ContactObject({
    $core.String? id,
    ContactType? type,
    $core.String? detail,
    $core.bool? verified,
    CommunicationLevel? communicationLevel,
    $2.STATE? state,
    $0.Struct? extra,
  }) {
    final result = ContactObject._();
    if (id != null) result.id = id;
    if (type != null) result.type = type;
    if (detail != null) result.detail = detail;
    if (verified != null) result.verified = verified;
    if (communicationLevel != null)
      result.communicationLevel = communicationLevel;
    if (state != null) result.state = state;
    if (extra != null) result.extra = extra;
    return result;
  }

  ContactObject._();

  factory ContactObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ContactObject()..mergeFromBuffer(data, registry);
  factory ContactObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ContactObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContactObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ContactObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aE<ContactType>(2, _omitFieldNames ? '' : 'type',
        enumValues: ContactType.values)
    ..aOS(3, _omitFieldNames ? '' : 'detail')
    ..aOB(4, _omitFieldNames ? '' : 'verified')
    ..aE<CommunicationLevel>(5, _omitFieldNames ? '' : 'communicationLevel',
        enumValues: CommunicationLevel.values)
    ..aE<$2.STATE>(6, _omitFieldNames ? '' : 'state',
        enumValues: $2.STATE.values)
    ..aOM<$0.Struct>(7, _omitFieldNames ? '' : 'extra',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContactObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContactObject copyWith(void Function(ContactObject) updates) =>
      super.copyWith((message) => updates(message as ContactObject))
          as ContactObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ContactObject() / ContactObject.new instead')
  static ContactObject create() => ContactObject._();
  static $pb.GeneratedMessage $_createMessage() => ContactObject._();
  @$core.override
  ContactObject createEmptyInstance() => ContactObject._();
  @$core.pragma('dart2js:noInline')
  static ContactObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ContactObject>(
          ContactObject.$_createMessage);
  static ContactObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  ContactType get type => $_getN(1);
  @$pb.TagNumber(2)
  set type(ContactType value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(2)
  void clearType() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get detail => $_getSZ(2);
  @$pb.TagNumber(3)
  set detail($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasDetail() => $_has(2);
  @$pb.TagNumber(3)
  void clearDetail() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get verified => $_getBF(3);
  @$pb.TagNumber(4)
  set verified($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasVerified() => $_has(3);
  @$pb.TagNumber(4)
  void clearVerified() => $_clearField(4);

  @$pb.TagNumber(5)
  CommunicationLevel get communicationLevel => $_getN(4);
  @$pb.TagNumber(5)
  set communicationLevel(CommunicationLevel value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasCommunicationLevel() => $_has(4);
  @$pb.TagNumber(5)
  void clearCommunicationLevel() => $_clearField(5);

  @$pb.TagNumber(6)
  $2.STATE get state => $_getN(5);
  @$pb.TagNumber(6)
  set state($2.STATE value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(6)
  void clearState() => $_clearField(6);

  @$pb.TagNumber(7)
  $0.Struct get extra => $_getN(6);
  @$pb.TagNumber(7)
  set extra($0.Struct value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasExtra() => $_has(6);
  @$pb.TagNumber(7)
  void clearExtra() => $_clearField(7);
  @$pb.TagNumber(7)
  $0.Struct ensureExtra() => $_ensure(6);
}

/// RosterObject represents a contact in a user's roster/contact list.
class RosterObject extends $pb.GeneratedMessage {
  factory RosterObject({
    $core.String? id,
    $core.String? profileId,
    ContactObject? contact,
    $0.Struct? extra,
    $core.String? name,
  }) {
    final result = RosterObject._();
    if (id != null) result.id = id;
    if (profileId != null) result.profileId = profileId;
    if (contact != null) result.contact = contact;
    if (extra != null) result.extra = extra;
    if (name != null) result.name = name;
    return result;
  }

  RosterObject._();

  factory RosterObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RosterObject()..mergeFromBuffer(data, registry);
  factory RosterObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RosterObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RosterObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RosterObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'profileId')
    ..aOM<ContactObject>(3, _omitFieldNames ? '' : 'contact',
        subBuilder: ContactObject.$_createMessage)
    ..aOM<$0.Struct>(4, _omitFieldNames ? '' : 'extra',
        subBuilder: $0.Struct.$_createMessage)
    ..aOS(5, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RosterObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RosterObject copyWith(void Function(RosterObject) updates) =>
      super.copyWith((message) => updates(message as RosterObject))
          as RosterObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use RosterObject() / RosterObject.new instead')
  static RosterObject create() => RosterObject._();
  static $pb.GeneratedMessage $_createMessage() => RosterObject._();
  @$core.override
  RosterObject createEmptyInstance() => RosterObject._();
  @$core.pragma('dart2js:noInline')
  static RosterObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RosterObject>(
          RosterObject.$_createMessage);
  static RosterObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get profileId => $_getSZ(1);
  @$pb.TagNumber(2)
  set profileId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasProfileId() => $_has(1);
  @$pb.TagNumber(2)
  void clearProfileId() => $_clearField(2);

  @$pb.TagNumber(3)
  ContactObject get contact => $_getN(2);
  @$pb.TagNumber(3)
  set contact(ContactObject value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasContact() => $_has(2);
  @$pb.TagNumber(3)
  void clearContact() => $_clearField(3);
  @$pb.TagNumber(3)
  ContactObject ensureContact() => $_ensure(2);

  @$pb.TagNumber(4)
  $0.Struct get extra => $_getN(3);
  @$pb.TagNumber(4)
  set extra($0.Struct value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasExtra() => $_has(3);
  @$pb.TagNumber(4)
  void clearExtra() => $_clearField(4);
  @$pb.TagNumber(4)
  $0.Struct ensureExtra() => $_ensure(3);

  @$pb.TagNumber(5)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(5)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(5)
  void clearName() => $_clearField(5);
}

/// AddressObject represents a physical address with geocoding.
class AddressObject extends $pb.GeneratedMessage {
  factory AddressObject({
    $core.String? id,
    $core.String? name,
    $core.String? country,
    $core.String? city,
    $core.String? area,
    $core.String? street,
    $core.String? house,
    $core.String? postcode,
    $core.double? latitude,
    $core.double? longitude,
    $core.String? extra,
  }) {
    final result = AddressObject._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (country != null) result.country = country;
    if (city != null) result.city = city;
    if (area != null) result.area = area;
    if (street != null) result.street = street;
    if (house != null) result.house = house;
    if (postcode != null) result.postcode = postcode;
    if (latitude != null) result.latitude = latitude;
    if (longitude != null) result.longitude = longitude;
    if (extra != null) result.extra = extra;
    return result;
  }

  AddressObject._();

  factory AddressObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddressObject()..mergeFromBuffer(data, registry);
  factory AddressObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddressObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddressObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddressObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'country')
    ..aOS(4, _omitFieldNames ? '' : 'city')
    ..aOS(5, _omitFieldNames ? '' : 'area')
    ..aOS(6, _omitFieldNames ? '' : 'street')
    ..aOS(7, _omitFieldNames ? '' : 'house')
    ..aOS(8, _omitFieldNames ? '' : 'postcode')
    ..aD(9, _omitFieldNames ? '' : 'latitude')
    ..aD(10, _omitFieldNames ? '' : 'longitude')
    ..aOS(11, _omitFieldNames ? '' : 'extra')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddressObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddressObject copyWith(void Function(AddressObject) updates) =>
      super.copyWith((message) => updates(message as AddressObject))
          as AddressObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddressObject() / AddressObject.new instead')
  static AddressObject create() => AddressObject._();
  static $pb.GeneratedMessage $_createMessage() => AddressObject._();
  @$core.override
  AddressObject createEmptyInstance() => AddressObject._();
  @$core.pragma('dart2js:noInline')
  static AddressObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddressObject>(
          AddressObject.$_createMessage);
  static AddressObject? _defaultInstance;

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
  $core.String get country => $_getSZ(2);
  @$pb.TagNumber(3)
  set country($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCountry() => $_has(2);
  @$pb.TagNumber(3)
  void clearCountry() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get city => $_getSZ(3);
  @$pb.TagNumber(4)
  set city($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCity() => $_has(3);
  @$pb.TagNumber(4)
  void clearCity() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get area => $_getSZ(4);
  @$pb.TagNumber(5)
  set area($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasArea() => $_has(4);
  @$pb.TagNumber(5)
  void clearArea() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get street => $_getSZ(5);
  @$pb.TagNumber(6)
  set street($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasStreet() => $_has(5);
  @$pb.TagNumber(6)
  void clearStreet() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get house => $_getSZ(6);
  @$pb.TagNumber(7)
  set house($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasHouse() => $_has(6);
  @$pb.TagNumber(7)
  void clearHouse() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get postcode => $_getSZ(7);
  @$pb.TagNumber(8)
  set postcode($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPostcode() => $_has(7);
  @$pb.TagNumber(8)
  void clearPostcode() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.double get latitude => $_getN(8);
  @$pb.TagNumber(9)
  set latitude($core.double value) => $_setDouble(8, value);
  @$pb.TagNumber(9)
  $core.bool hasLatitude() => $_has(8);
  @$pb.TagNumber(9)
  void clearLatitude() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.double get longitude => $_getN(9);
  @$pb.TagNumber(10)
  set longitude($core.double value) => $_setDouble(9, value);
  @$pb.TagNumber(10)
  $core.bool hasLongitude() => $_has(9);
  @$pb.TagNumber(10)
  void clearLongitude() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.String get extra => $_getSZ(10);
  @$pb.TagNumber(11)
  set extra($core.String value) => $_setString(10, value);
  @$pb.TagNumber(11)
  $core.bool hasExtra() => $_has(10);
  @$pb.TagNumber(11)
  void clearExtra() => $_clearField(11);
}

/// ProfileObject represents a complete user or entity profile.
class ProfileObject extends $pb.GeneratedMessage {
  factory ProfileObject({
    $core.String? id,
    ProfileType? type,
    $0.Struct? properties,
    $core.Iterable<ContactObject>? contacts,
    $core.Iterable<AddressObject>? addresses,
    $2.STATE? state,
    $core.Iterable<ProfileAccount>? accounts,
  }) {
    final result = ProfileObject._();
    if (id != null) result.id = id;
    if (type != null) result.type = type;
    if (properties != null) result.properties = properties;
    if (contacts != null) result.contacts.addAll(contacts);
    if (addresses != null) result.addresses.addAll(addresses);
    if (state != null) result.state = state;
    if (accounts != null) result.accounts.addAll(accounts);
    return result;
  }

  ProfileObject._();

  factory ProfileObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ProfileObject()..mergeFromBuffer(data, registry);
  factory ProfileObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ProfileObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ProfileObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aE<ProfileType>(2, _omitFieldNames ? '' : 'type',
        enumValues: ProfileType.values)
    ..aOM<$0.Struct>(3, _omitFieldNames ? '' : 'properties',
        subBuilder: $0.Struct.$_createMessage)
    ..pPM<ContactObject>(4, _omitFieldNames ? '' : 'contacts',
        subBuilder: ContactObject.$_createMessage)
    ..pPM<AddressObject>(5, _omitFieldNames ? '' : 'addresses',
        subBuilder: AddressObject.$_createMessage)
    ..aE<$2.STATE>(6, _omitFieldNames ? '' : 'state',
        enumValues: $2.STATE.values)
    ..pPM<ProfileAccount>(7, _omitFieldNames ? '' : 'accounts',
        subBuilder: ProfileAccount.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileObject copyWith(void Function(ProfileObject) updates) =>
      super.copyWith((message) => updates(message as ProfileObject))
          as ProfileObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ProfileObject() / ProfileObject.new instead')
  static ProfileObject create() => ProfileObject._();
  static $pb.GeneratedMessage $_createMessage() => ProfileObject._();
  @$core.override
  ProfileObject createEmptyInstance() => ProfileObject._();
  @$core.pragma('dart2js:noInline')
  static ProfileObject getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileObject>(
          ProfileObject.$_createMessage);
  static ProfileObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  ProfileType get type => $_getN(1);
  @$pb.TagNumber(2)
  set type(ProfileType value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(2)
  void clearType() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Struct get properties => $_getN(2);
  @$pb.TagNumber(3)
  set properties($0.Struct value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasProperties() => $_has(2);
  @$pb.TagNumber(3)
  void clearProperties() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Struct ensureProperties() => $_ensure(2);

  @$pb.TagNumber(4)
  $pb.PbList<ContactObject> get contacts => $_getList(3);

  @$pb.TagNumber(5)
  $pb.PbList<AddressObject> get addresses => $_getList(4);

  @$pb.TagNumber(6)
  $2.STATE get state => $_getN(5);
  @$pb.TagNumber(6)
  set state($2.STATE value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(6)
  void clearState() => $_clearField(6);

  /// Chain accounts the profile owns, primary first. Only PERSON profiles
  /// own accounts; a merge adds the merged profile's accounts as secondary.
  @$pb.TagNumber(7)
  $pb.PbList<ProfileAccount> get accounts => $_getList(6);
}

/// ProfileAccount is a chain account derived for a profile:
/// address = CREATE2(factory, ACCOUNT_SALT(family, version, identity_salt_hash),
/// keccak256(creation_code || identity_salt_hash)). The identity salt itself
/// never leaves the profile service.
class ProfileAccount extends $pb.GeneratedMessage {
  factory ProfileAccount({
    $core.String? address,
    $core.String? family,
    $core.int? version,
    $core.List<$core.int>? identitySaltHash,
    $core.bool? primary,
  }) {
    final result = ProfileAccount._();
    if (address != null) result.address = address;
    if (family != null) result.family = family;
    if (version != null) result.version = version;
    if (identitySaltHash != null) result.identitySaltHash = identitySaltHash;
    if (primary != null) result.primary = primary;
    return result;
  }

  ProfileAccount._();

  factory ProfileAccount.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ProfileAccount()..mergeFromBuffer(data, registry);
  factory ProfileAccount.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ProfileAccount()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProfileAccount',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ProfileAccount.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOS(2, _omitFieldNames ? '' : 'family')
    ..aI(3, _omitFieldNames ? '' : 'version', fieldType: $pb.PbFieldType.OU3)
    ..a<$core.List<$core.int>>(
        4, _omitFieldNames ? '' : 'identitySaltHash', $pb.PbFieldType.OY)
    ..aOB(5, _omitFieldNames ? '' : 'primary')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileAccount clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProfileAccount copyWith(void Function(ProfileAccount) updates) =>
      super.copyWith((message) => updates(message as ProfileAccount))
          as ProfileAccount;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ProfileAccount() / ProfileAccount.new instead')
  static ProfileAccount create() => ProfileAccount._();
  static $pb.GeneratedMessage $_createMessage() => ProfileAccount._();
  @$core.override
  ProfileAccount createEmptyInstance() => ProfileAccount._();
  @$core.pragma('dart2js:noInline')
  static ProfileAccount getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProfileAccount>(
          ProfileAccount.$_createMessage);
  static ProfileAccount? _defaultInstance;

  /// Account address, lowercase 0x-prefixed hex (20 bytes).
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// Account family; "EVM" in protocol v1.
  @$pb.TagNumber(2)
  $core.String get family => $_getSZ(1);
  @$pb.TagNumber(2)
  set family($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasFamily() => $_has(1);
  @$pb.TagNumber(2)
  void clearFamily() => $_clearField(2);

  /// Protocol account version the address was derived for.
  @$pb.TagNumber(3)
  $core.int get version => $_getIZ(2);
  @$pb.TagNumber(3)
  set version($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasVersion() => $_has(2);
  @$pb.TagNumber(3)
  void clearVersion() => $_clearField(3);

  /// keccak256 of the profile's identity salt (32 bytes); the StawiAccount
  /// constructor argument.
  @$pb.TagNumber(4)
  $core.List<$core.int> get identitySaltHash => $_getN(3);
  @$pb.TagNumber(4)
  set identitySaltHash($core.List<$core.int> value) => $_setBytes(3, value);
  @$pb.TagNumber(4)
  $core.bool hasIdentitySaltHash() => $_has(3);
  @$pb.TagNumber(4)
  void clearIdentitySaltHash() => $_clearField(4);

  /// True for the account derived for this profile; false for accounts that
  /// arrived through a merge.
  @$pb.TagNumber(5)
  $core.bool get primary => $_getBF(4);
  @$pb.TagNumber(5)
  set primary($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPrimary() => $_has(4);
  @$pb.TagNumber(5)
  void clearPrimary() => $_clearField(5);
}

/// EntryItem identifies an object in a relationship.
class EntryItem extends $pb.GeneratedMessage {
  factory EntryItem({
    $core.String? objectName,
    $core.String? objectId,
  }) {
    final result = EntryItem._();
    if (objectName != null) result.objectName = objectName;
    if (objectId != null) result.objectId = objectId;
    return result;
  }

  EntryItem._();

  factory EntryItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      EntryItem()..mergeFromBuffer(data, registry);
  factory EntryItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      EntryItem()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EntryItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: EntryItem.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'objectName')
    ..aOS(2, _omitFieldNames ? '' : 'objectId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EntryItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EntryItem copyWith(void Function(EntryItem) updates) =>
      super.copyWith((message) => updates(message as EntryItem)) as EntryItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use EntryItem() / EntryItem.new instead')
  static EntryItem create() => EntryItem._();
  static $pb.GeneratedMessage $_createMessage() => EntryItem._();
  @$core.override
  EntryItem createEmptyInstance() => EntryItem._();
  @$core.pragma('dart2js:noInline')
  static EntryItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EntryItem>(EntryItem.$_createMessage);
  static EntryItem? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get objectName => $_getSZ(0);
  @$pb.TagNumber(1)
  set objectName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasObjectName() => $_has(0);
  @$pb.TagNumber(1)
  void clearObjectName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get objectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set objectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasObjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearObjectId() => $_clearField(2);
}

/// RelationshipObject represents a relationship between two profiles.
class RelationshipObject extends $pb.GeneratedMessage {
  factory RelationshipObject({
    $core.String? id,
    RelationshipType? type,
    $0.Struct? properties,
    EntryItem? childEntry,
    EntryItem? parentEntry,
    ProfileObject? peerProfile,
  }) {
    final result = RelationshipObject._();
    if (id != null) result.id = id;
    if (type != null) result.type = type;
    if (properties != null) result.properties = properties;
    if (childEntry != null) result.childEntry = childEntry;
    if (parentEntry != null) result.parentEntry = parentEntry;
    if (peerProfile != null) result.peerProfile = peerProfile;
    return result;
  }

  RelationshipObject._();

  factory RelationshipObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RelationshipObject()..mergeFromBuffer(data, registry);
  factory RelationshipObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RelationshipObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RelationshipObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RelationshipObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aE<RelationshipType>(2, _omitFieldNames ? '' : 'type',
        enumValues: RelationshipType.values)
    ..aOM<$0.Struct>(3, _omitFieldNames ? '' : 'properties',
        subBuilder: $0.Struct.$_createMessage)
    ..aOM<EntryItem>(4, _omitFieldNames ? '' : 'childEntry',
        subBuilder: EntryItem.$_createMessage)
    ..aOM<EntryItem>(5, _omitFieldNames ? '' : 'parentEntry',
        subBuilder: EntryItem.$_createMessage)
    ..aOM<ProfileObject>(6, _omitFieldNames ? '' : 'peerProfile',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RelationshipObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RelationshipObject copyWith(void Function(RelationshipObject) updates) =>
      super.copyWith((message) => updates(message as RelationshipObject))
          as RelationshipObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use RelationshipObject() / RelationshipObject.new instead')
  static RelationshipObject create() => RelationshipObject._();
  static $pb.GeneratedMessage $_createMessage() => RelationshipObject._();
  @$core.override
  RelationshipObject createEmptyInstance() => RelationshipObject._();
  @$core.pragma('dart2js:noInline')
  static RelationshipObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RelationshipObject>(
          RelationshipObject.$_createMessage);
  static RelationshipObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  RelationshipType get type => $_getN(1);
  @$pb.TagNumber(2)
  set type(RelationshipType value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(2)
  void clearType() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Struct get properties => $_getN(2);
  @$pb.TagNumber(3)
  set properties($0.Struct value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasProperties() => $_has(2);
  @$pb.TagNumber(3)
  void clearProperties() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Struct ensureProperties() => $_ensure(2);

  @$pb.TagNumber(4)
  EntryItem get childEntry => $_getN(3);
  @$pb.TagNumber(4)
  set childEntry(EntryItem value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasChildEntry() => $_has(3);
  @$pb.TagNumber(4)
  void clearChildEntry() => $_clearField(4);
  @$pb.TagNumber(4)
  EntryItem ensureChildEntry() => $_ensure(3);

  @$pb.TagNumber(5)
  EntryItem get parentEntry => $_getN(4);
  @$pb.TagNumber(5)
  set parentEntry(EntryItem value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasParentEntry() => $_has(4);
  @$pb.TagNumber(5)
  void clearParentEntry() => $_clearField(5);
  @$pb.TagNumber(5)
  EntryItem ensureParentEntry() => $_ensure(4);

  @$pb.TagNumber(6)
  ProfileObject get peerProfile => $_getN(5);
  @$pb.TagNumber(6)
  set peerProfile(ProfileObject value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasPeerProfile() => $_has(5);
  @$pb.TagNumber(6)
  void clearPeerProfile() => $_clearField(6);
  @$pb.TagNumber(6)
  ProfileObject ensurePeerProfile() => $_ensure(5);
}

/// The request message containing the profile's hash
class GetByIdRequest extends $pb.GeneratedMessage {
  factory GetByIdRequest({
    $core.String? id,
  }) {
    final result = GetByIdRequest._();
    if (id != null) result.id = id;
    return result;
  }

  GetByIdRequest._();

  factory GetByIdRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIdRequest()..mergeFromBuffer(data, registry);
  factory GetByIdRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIdRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByIdRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetByIdRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIdRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIdRequest copyWith(void Function(GetByIdRequest) updates) =>
      super.copyWith((message) => updates(message as GetByIdRequest))
          as GetByIdRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetByIdRequest() / GetByIdRequest.new instead')
  static GetByIdRequest create() => GetByIdRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetByIdRequest._();
  @$core.override
  GetByIdRequest createEmptyInstance() => GetByIdRequest._();
  @$core.pragma('dart2js:noInline')
  static GetByIdRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetByIdRequest>(
          GetByIdRequest.$_createMessage);
  static GetByIdRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class GetByIdResponse extends $pb.GeneratedMessage {
  factory GetByIdResponse({
    ProfileObject? data,
  }) {
    final result = GetByIdResponse._();
    if (data != null) result.data = data;
    return result;
  }

  GetByIdResponse._();

  factory GetByIdResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIdResponse()..mergeFromBuffer(data, registry);
  factory GetByIdResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIdResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByIdResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetByIdResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIdResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIdResponse copyWith(void Function(GetByIdResponse) updates) =>
      super.copyWith((message) => updates(message as GetByIdResponse))
          as GetByIdResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetByIdResponse() / GetByIdResponse.new instead')
  static GetByIdResponse create() => GetByIdResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetByIdResponse._();
  @$core.override
  GetByIdResponse createEmptyInstance() => GetByIdResponse._();
  @$core.pragma('dart2js:noInline')
  static GetByIdResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetByIdResponse>(
          GetByIdResponse.$_createMessage);
  static GetByIdResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

class SearchRequest extends $pb.GeneratedMessage {
  factory SearchRequest({
    $core.String? query,
    $fixnum.Int64? page,
    $core.int? count,
    $core.String? startDate,
    $core.String? endDate,
    $core.Iterable<$core.String>? properties,
    $0.Struct? extras,
  }) {
    final result = SearchRequest._();
    if (query != null) result.query = query;
    if (page != null) result.page = page;
    if (count != null) result.count = count;
    if (startDate != null) result.startDate = startDate;
    if (endDate != null) result.endDate = endDate;
    if (properties != null) result.properties.addAll(properties);
    if (extras != null) result.extras = extras;
    return result;
  }

  SearchRequest._();

  factory SearchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRequest()..mergeFromBuffer(data, registry);
  factory SearchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: SearchRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..aInt64(2, _omitFieldNames ? '' : 'page')
    ..aI(3, _omitFieldNames ? '' : 'count')
    ..aOS(4, _omitFieldNames ? '' : 'startDate')
    ..aOS(5, _omitFieldNames ? '' : 'endDate')
    ..pPS(6, _omitFieldNames ? '' : 'properties')
    ..aOM<$0.Struct>(7, _omitFieldNames ? '' : 'extras',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRequest copyWith(void Function(SearchRequest) updates) =>
      super.copyWith((message) => updates(message as SearchRequest))
          as SearchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use SearchRequest() / SearchRequest.new instead')
  static SearchRequest create() => SearchRequest._();
  static $pb.GeneratedMessage $_createMessage() => SearchRequest._();
  @$core.override
  SearchRequest createEmptyInstance() => SearchRequest._();
  @$core.pragma('dart2js:noInline')
  static SearchRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SearchRequest>(
          SearchRequest.$_createMessage);
  static SearchRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get page => $_getI64(1);
  @$pb.TagNumber(2)
  set page($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPage() => $_has(1);
  @$pb.TagNumber(2)
  void clearPage() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get count => $_getIZ(2);
  @$pb.TagNumber(3)
  set count($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCount() => $_has(2);
  @$pb.TagNumber(3)
  void clearCount() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get startDate => $_getSZ(3);
  @$pb.TagNumber(4)
  set startDate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasStartDate() => $_has(3);
  @$pb.TagNumber(4)
  void clearStartDate() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get endDate => $_getSZ(4);
  @$pb.TagNumber(5)
  set endDate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasEndDate() => $_has(4);
  @$pb.TagNumber(5)
  void clearEndDate() => $_clearField(5);

  @$pb.TagNumber(6)
  $pb.PbList<$core.String> get properties => $_getList(5);

  @$pb.TagNumber(7)
  $0.Struct get extras => $_getN(6);
  @$pb.TagNumber(7)
  set extras($0.Struct value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasExtras() => $_has(6);
  @$pb.TagNumber(7)
  void clearExtras() => $_clearField(7);
  @$pb.TagNumber(7)
  $0.Struct ensureExtras() => $_ensure(6);
}

class SearchResponse extends $pb.GeneratedMessage {
  factory SearchResponse({
    $core.Iterable<ProfileObject>? data,
  }) {
    final result = SearchResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  SearchResponse._();

  factory SearchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchResponse()..mergeFromBuffer(data, registry);
  factory SearchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: SearchResponse.$_createMessage)
    ..pPM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchResponse copyWith(void Function(SearchResponse) updates) =>
      super.copyWith((message) => updates(message as SearchResponse))
          as SearchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use SearchResponse() / SearchResponse.new instead')
  static SearchResponse create() => SearchResponse._();
  static $pb.GeneratedMessage $_createMessage() => SearchResponse._();
  @$core.override
  SearchResponse createEmptyInstance() => SearchResponse._();
  @$core.pragma('dart2js:noInline')
  static SearchResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SearchResponse>(
          SearchResponse.$_createMessage);
  static SearchResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<ProfileObject> get data => $_getList(0);
}

/// The request message containing the profile's hash
class MergeRequest extends $pb.GeneratedMessage {
  factory MergeRequest({
    $core.String? id,
    $core.String? mergeid,
  }) {
    final result = MergeRequest._();
    if (id != null) result.id = id;
    if (mergeid != null) result.mergeid = mergeid;
    return result;
  }

  MergeRequest._();

  factory MergeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MergeRequest()..mergeFromBuffer(data, registry);
  factory MergeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MergeRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MergeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: MergeRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'mergeid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeRequest copyWith(void Function(MergeRequest) updates) =>
      super.copyWith((message) => updates(message as MergeRequest))
          as MergeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use MergeRequest() / MergeRequest.new instead')
  static MergeRequest create() => MergeRequest._();
  static $pb.GeneratedMessage $_createMessage() => MergeRequest._();
  @$core.override
  MergeRequest createEmptyInstance() => MergeRequest._();
  @$core.pragma('dart2js:noInline')
  static MergeRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MergeRequest>(
          MergeRequest.$_createMessage);
  static MergeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get mergeid => $_getSZ(1);
  @$pb.TagNumber(2)
  set mergeid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMergeid() => $_has(1);
  @$pb.TagNumber(2)
  void clearMergeid() => $_clearField(2);
}

class MergeResponse extends $pb.GeneratedMessage {
  factory MergeResponse({
    ProfileObject? data,
  }) {
    final result = MergeResponse._();
    if (data != null) result.data = data;
    return result;
  }

  MergeResponse._();

  factory MergeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MergeResponse()..mergeFromBuffer(data, registry);
  factory MergeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MergeResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MergeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: MergeResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeResponse copyWith(void Function(MergeResponse) updates) =>
      super.copyWith((message) => updates(message as MergeResponse))
          as MergeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use MergeResponse() / MergeResponse.new instead')
  static MergeResponse create() => MergeResponse._();
  static $pb.GeneratedMessage $_createMessage() => MergeResponse._();
  @$core.override
  MergeResponse createEmptyInstance() => MergeResponse._();
  @$core.pragma('dart2js:noInline')
  static MergeResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MergeResponse>(
          MergeResponse.$_createMessage);
  static MergeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

/// The request message containing the information necessary to create a profile
class CreateRequest extends $pb.GeneratedMessage {
  factory CreateRequest({
    ProfileType? type,
    $core.String? contact,
    $0.Struct? properties,
  }) {
    final result = CreateRequest._();
    if (type != null) result.type = type;
    if (contact != null) result.contact = contact;
    if (properties != null) result.properties = properties;
    return result;
  }

  CreateRequest._();

  factory CreateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateRequest()..mergeFromBuffer(data, registry);
  factory CreateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CreateRequest.$_createMessage)
    ..aE<ProfileType>(1, _omitFieldNames ? '' : 'type',
        enumValues: ProfileType.values)
    ..aOS(2, _omitFieldNames ? '' : 'contact')
    ..aOM<$0.Struct>(3, _omitFieldNames ? '' : 'properties',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRequest copyWith(void Function(CreateRequest) updates) =>
      super.copyWith((message) => updates(message as CreateRequest))
          as CreateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CreateRequest() / CreateRequest.new instead')
  static CreateRequest create() => CreateRequest._();
  static $pb.GeneratedMessage $_createMessage() => CreateRequest._();
  @$core.override
  CreateRequest createEmptyInstance() => CreateRequest._();
  @$core.pragma('dart2js:noInline')
  static CreateRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateRequest>(
          CreateRequest.$_createMessage);
  static CreateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileType get type => $_getN(0);
  @$pb.TagNumber(1)
  set type(ProfileType value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get contact => $_getSZ(1);
  @$pb.TagNumber(2)
  set contact($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasContact() => $_has(1);
  @$pb.TagNumber(2)
  void clearContact() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Struct get properties => $_getN(2);
  @$pb.TagNumber(3)
  set properties($0.Struct value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasProperties() => $_has(2);
  @$pb.TagNumber(3)
  void clearProperties() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Struct ensureProperties() => $_ensure(2);
}

class CreateResponse extends $pb.GeneratedMessage {
  factory CreateResponse({
    ProfileObject? data,
  }) {
    final result = CreateResponse._();
    if (data != null) result.data = data;
    return result;
  }

  CreateResponse._();

  factory CreateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateResponse()..mergeFromBuffer(data, registry);
  factory CreateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CreateResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateResponse copyWith(void Function(CreateResponse) updates) =>
      super.copyWith((message) => updates(message as CreateResponse))
          as CreateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CreateResponse() / CreateResponse.new instead')
  static CreateResponse create() => CreateResponse._();
  static $pb.GeneratedMessage $_createMessage() => CreateResponse._();
  @$core.override
  CreateResponse createEmptyInstance() => CreateResponse._();
  @$core.pragma('dart2js:noInline')
  static CreateResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateResponse>(
          CreateResponse.$_createMessage);
  static CreateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

/// The request message containing the information necessary to create a profile
class UpdateRequest extends $pb.GeneratedMessage {
  factory UpdateRequest({
    $core.String? id,
    $0.Struct? properties,
    $2.STATE? state,
    $core.bool? scoped,
  }) {
    final result = UpdateRequest._();
    if (id != null) result.id = id;
    if (properties != null) result.properties = properties;
    if (state != null) result.state = state;
    if (scoped != null) result.scoped = scoped;
    return result;
  }

  UpdateRequest._();

  factory UpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRequest()..mergeFromBuffer(data, registry);
  factory UpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: UpdateRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOM<$0.Struct>(2, _omitFieldNames ? '' : 'properties',
        subBuilder: $0.Struct.$_createMessage)
    ..aE<$2.STATE>(3, _omitFieldNames ? '' : 'state',
        enumValues: $2.STATE.values)
    ..aOB(4, _omitFieldNames ? '' : 'scoped')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRequest copyWith(void Function(UpdateRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateRequest))
          as UpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UpdateRequest() / UpdateRequest.new instead')
  static UpdateRequest create() => UpdateRequest._();
  static $pb.GeneratedMessage $_createMessage() => UpdateRequest._();
  @$core.override
  UpdateRequest createEmptyInstance() => UpdateRequest._();
  @$core.pragma('dart2js:noInline')
  static UpdateRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateRequest>(
          UpdateRequest.$_createMessage);
  static UpdateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $0.Struct get properties => $_getN(1);
  @$pb.TagNumber(2)
  set properties($0.Struct value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasProperties() => $_has(1);
  @$pb.TagNumber(2)
  void clearProperties() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.Struct ensureProperties() => $_ensure(1);

  @$pb.TagNumber(3)
  $2.STATE get state => $_getN(2);
  @$pb.TagNumber(3)
  set state($2.STATE value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasState() => $_has(2);
  @$pb.TagNumber(3)
  void clearState() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get scoped => $_getBF(3);
  @$pb.TagNumber(4)
  set scoped($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasScoped() => $_has(3);
  @$pb.TagNumber(4)
  void clearScoped() => $_clearField(4);
}

class UpdateResponse extends $pb.GeneratedMessage {
  factory UpdateResponse({
    ProfileObject? data,
  }) {
    final result = UpdateResponse._();
    if (data != null) result.data = data;
    return result;
  }

  UpdateResponse._();

  factory UpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateResponse()..mergeFromBuffer(data, registry);
  factory UpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: UpdateResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateResponse copyWith(void Function(UpdateResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateResponse))
          as UpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UpdateResponse() / UpdateResponse.new instead')
  static UpdateResponse create() => UpdateResponse._();
  static $pb.GeneratedMessage $_createMessage() => UpdateResponse._();
  @$core.override
  UpdateResponse createEmptyInstance() => UpdateResponse._();
  @$core.pragma('dart2js:noInline')
  static UpdateResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateResponse>(
          UpdateResponse.$_createMessage);
  static UpdateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

/// The request message containing the information necessary to create a profile
class AddContactRequest extends $pb.GeneratedMessage {
  factory AddContactRequest({
    $core.String? id,
    $core.String? contact,
    $0.Struct? extras,
  }) {
    final result = AddContactRequest._();
    if (id != null) result.id = id;
    if (contact != null) result.contact = contact;
    if (extras != null) result.extras = extras;
    return result;
  }

  AddContactRequest._();

  factory AddContactRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddContactRequest()..mergeFromBuffer(data, registry);
  factory AddContactRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddContactRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddContactRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddContactRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'contact')
    ..aOM<$0.Struct>(3, _omitFieldNames ? '' : 'extras',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddContactRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddContactRequest copyWith(void Function(AddContactRequest) updates) =>
      super.copyWith((message) => updates(message as AddContactRequest))
          as AddContactRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddContactRequest() / AddContactRequest.new instead')
  static AddContactRequest create() => AddContactRequest._();
  static $pb.GeneratedMessage $_createMessage() => AddContactRequest._();
  @$core.override
  AddContactRequest createEmptyInstance() => AddContactRequest._();
  @$core.pragma('dart2js:noInline')
  static AddContactRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddContactRequest>(
          AddContactRequest.$_createMessage);
  static AddContactRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get contact => $_getSZ(1);
  @$pb.TagNumber(2)
  set contact($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasContact() => $_has(1);
  @$pb.TagNumber(2)
  void clearContact() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Struct get extras => $_getN(2);
  @$pb.TagNumber(3)
  set extras($0.Struct value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasExtras() => $_has(2);
  @$pb.TagNumber(3)
  void clearExtras() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Struct ensureExtras() => $_ensure(2);
}

class AddContactResponse extends $pb.GeneratedMessage {
  factory AddContactResponse({
    ProfileObject? data,
    $core.String? verificationId,
  }) {
    final result = AddContactResponse._();
    if (data != null) result.data = data;
    if (verificationId != null) result.verificationId = verificationId;
    return result;
  }

  AddContactResponse._();

  factory AddContactResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddContactResponse()..mergeFromBuffer(data, registry);
  factory AddContactResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddContactResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddContactResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddContactResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'verificationId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddContactResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddContactResponse copyWith(void Function(AddContactResponse) updates) =>
      super.copyWith((message) => updates(message as AddContactResponse))
          as AddContactResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddContactResponse() / AddContactResponse.new instead')
  static AddContactResponse create() => AddContactResponse._();
  static $pb.GeneratedMessage $_createMessage() => AddContactResponse._();
  @$core.override
  AddContactResponse createEmptyInstance() => AddContactResponse._();
  @$core.pragma('dart2js:noInline')
  static AddContactResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddContactResponse>(
          AddContactResponse.$_createMessage);
  static AddContactResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.String get verificationId => $_getSZ(1);
  @$pb.TagNumber(2)
  set verificationId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasVerificationId() => $_has(1);
  @$pb.TagNumber(2)
  void clearVerificationId() => $_clearField(2);
}

/// The request message containing the information necessary to create a profile
class CreateContactRequest extends $pb.GeneratedMessage {
  factory CreateContactRequest({
    $core.String? id,
    $core.String? contact,
    $0.Struct? extras,
  }) {
    final result = CreateContactRequest._();
    if (id != null) result.id = id;
    if (contact != null) result.contact = contact;
    if (extras != null) result.extras = extras;
    return result;
  }

  CreateContactRequest._();

  factory CreateContactRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactRequest()..mergeFromBuffer(data, registry);
  factory CreateContactRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateContactRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CreateContactRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'contact')
    ..aOM<$0.Struct>(3, _omitFieldNames ? '' : 'extras',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactRequest copyWith(void Function(CreateContactRequest) updates) =>
      super.copyWith((message) => updates(message as CreateContactRequest))
          as CreateContactRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreateContactRequest() / CreateContactRequest.new instead')
  static CreateContactRequest create() => CreateContactRequest._();
  static $pb.GeneratedMessage $_createMessage() => CreateContactRequest._();
  @$core.override
  CreateContactRequest createEmptyInstance() => CreateContactRequest._();
  @$core.pragma('dart2js:noInline')
  static CreateContactRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateContactRequest>(
          CreateContactRequest.$_createMessage);
  static CreateContactRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get contact => $_getSZ(1);
  @$pb.TagNumber(2)
  set contact($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasContact() => $_has(1);
  @$pb.TagNumber(2)
  void clearContact() => $_clearField(2);

  @$pb.TagNumber(3)
  $0.Struct get extras => $_getN(2);
  @$pb.TagNumber(3)
  set extras($0.Struct value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasExtras() => $_has(2);
  @$pb.TagNumber(3)
  void clearExtras() => $_clearField(3);
  @$pb.TagNumber(3)
  $0.Struct ensureExtras() => $_ensure(2);
}

class CreateContactResponse extends $pb.GeneratedMessage {
  factory CreateContactResponse({
    ContactObject? data,
  }) {
    final result = CreateContactResponse._();
    if (data != null) result.data = data;
    return result;
  }

  CreateContactResponse._();

  factory CreateContactResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactResponse()..mergeFromBuffer(data, registry);
  factory CreateContactResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateContactResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CreateContactResponse.$_createMessage)
    ..aOM<ContactObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ContactObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactResponse copyWith(
          void Function(CreateContactResponse) updates) =>
      super.copyWith((message) => updates(message as CreateContactResponse))
          as CreateContactResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreateContactResponse() / CreateContactResponse.new instead')
  static CreateContactResponse create() => CreateContactResponse._();
  static $pb.GeneratedMessage $_createMessage() => CreateContactResponse._();
  @$core.override
  CreateContactResponse createEmptyInstance() => CreateContactResponse._();
  @$core.pragma('dart2js:noInline')
  static CreateContactResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateContactResponse>(
          CreateContactResponse.$_createMessage);
  static CreateContactResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ContactObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ContactObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ContactObject ensureData() => $_ensure(0);
}

/// GetContactsRequest resolves contacts by id (one or many).
class GetContactsRequest extends $pb.GeneratedMessage {
  factory GetContactsRequest({
    $core.Iterable<$core.String>? ids,
  }) {
    final result = GetContactsRequest._();
    if (ids != null) result.ids.addAll(ids);
    return result;
  }

  GetContactsRequest._();

  factory GetContactsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContactsRequest()..mergeFromBuffer(data, registry);
  factory GetContactsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContactsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetContactsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetContactsRequest.$_createMessage)
    ..pPS(1, _omitFieldNames ? '' : 'ids')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContactsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContactsRequest copyWith(void Function(GetContactsRequest) updates) =>
      super.copyWith((message) => updates(message as GetContactsRequest))
          as GetContactsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GetContactsRequest() / GetContactsRequest.new instead')
  static GetContactsRequest create() => GetContactsRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetContactsRequest._();
  @$core.override
  GetContactsRequest createEmptyInstance() => GetContactsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetContactsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetContactsRequest>(
          GetContactsRequest.$_createMessage);
  static GetContactsRequest? _defaultInstance;

  /// Contact ids to resolve. Pass one id or a batch (max 100).
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get ids => $_getList(0);
}

class GetContactsResponse extends $pb.GeneratedMessage {
  factory GetContactsResponse({
    $core.Iterable<ContactObject>? data,
    $core.Iterable<$core.String>? missingIds,
  }) {
    final result = GetContactsResponse._();
    if (data != null) result.data.addAll(data);
    if (missingIds != null) result.missingIds.addAll(missingIds);
    return result;
  }

  GetContactsResponse._();

  factory GetContactsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContactsResponse()..mergeFromBuffer(data, registry);
  factory GetContactsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContactsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetContactsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetContactsResponse.$_createMessage)
    ..pPM<ContactObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ContactObject.$_createMessage)
    ..pPS(2, _omitFieldNames ? '' : 'missingIds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContactsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContactsResponse copyWith(void Function(GetContactsResponse) updates) =>
      super.copyWith((message) => updates(message as GetContactsResponse))
          as GetContactsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use GetContactsResponse() / GetContactsResponse.new instead')
  static GetContactsResponse create() => GetContactsResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetContactsResponse._();
  @$core.override
  GetContactsResponse createEmptyInstance() => GetContactsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetContactsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetContactsResponse>(
          GetContactsResponse.$_createMessage);
  static GetContactsResponse? _defaultInstance;

  /// Found contacts (standalone or profile-attached).
  @$pb.TagNumber(1)
  $pb.PbList<ContactObject> get data => $_getList(0);

  /// Requested ids that were not found or not readable.
  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get missingIds => $_getList(1);
}

/// ResolveAccountsRequest maps chain account addresses to profiles.
class ResolveAccountsRequest extends $pb.GeneratedMessage {
  factory ResolveAccountsRequest({
    $core.Iterable<$core.String>? addresses,
  }) {
    final result = ResolveAccountsRequest._();
    if (addresses != null) result.addresses.addAll(addresses);
    return result;
  }

  ResolveAccountsRequest._();

  factory ResolveAccountsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResolveAccountsRequest()..mergeFromBuffer(data, registry);
  factory ResolveAccountsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResolveAccountsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResolveAccountsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ResolveAccountsRequest.$_createMessage)
    ..pPS(1, _omitFieldNames ? '' : 'addresses')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResolveAccountsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResolveAccountsRequest copyWith(
          void Function(ResolveAccountsRequest) updates) =>
      super.copyWith((message) => updates(message as ResolveAccountsRequest))
          as ResolveAccountsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ResolveAccountsRequest() / ResolveAccountsRequest.new instead')
  static ResolveAccountsRequest create() => ResolveAccountsRequest._();
  static $pb.GeneratedMessage $_createMessage() => ResolveAccountsRequest._();
  @$core.override
  ResolveAccountsRequest createEmptyInstance() => ResolveAccountsRequest._();
  @$core.pragma('dart2js:noInline')
  static ResolveAccountsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResolveAccountsRequest>(
          ResolveAccountsRequest.$_createMessage);
  static ResolveAccountsRequest? _defaultInstance;

  /// 0x-prefixed hex addresses, any case (max 500).
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get addresses => $_getList(0);
}

/// AccountOwner is the profile that owns an address.
class AccountOwner extends $pb.GeneratedMessage {
  factory AccountOwner({
    $core.String? address,
    $core.String? profileId,
  }) {
    final result = AccountOwner._();
    if (address != null) result.address = address;
    if (profileId != null) result.profileId = profileId;
    return result;
  }

  AccountOwner._();

  factory AccountOwner.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountOwner()..mergeFromBuffer(data, registry);
  factory AccountOwner.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AccountOwner()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountOwner',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AccountOwner.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOS(2, _omitFieldNames ? '' : 'profileId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountOwner clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountOwner copyWith(void Function(AccountOwner) updates) =>
      super.copyWith((message) => updates(message as AccountOwner))
          as AccountOwner;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AccountOwner() / AccountOwner.new instead')
  static AccountOwner create() => AccountOwner._();
  static $pb.GeneratedMessage $_createMessage() => AccountOwner._();
  @$core.override
  AccountOwner createEmptyInstance() => AccountOwner._();
  @$core.pragma('dart2js:noInline')
  static AccountOwner getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AccountOwner>(
          AccountOwner.$_createMessage);
  static AccountOwner? _defaultInstance;

  /// Account address, lowercase 0x-prefixed hex.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// Owning profile id.
  @$pb.TagNumber(2)
  $core.String get profileId => $_getSZ(1);
  @$pb.TagNumber(2)
  set profileId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasProfileId() => $_has(1);
  @$pb.TagNumber(2)
  void clearProfileId() => $_clearField(2);
}

/// ResolveAccountsResponse lists the owners of the known addresses. Addresses
/// no profile owns are omitted.
class ResolveAccountsResponse extends $pb.GeneratedMessage {
  factory ResolveAccountsResponse({
    $core.Iterable<AccountOwner>? data,
  }) {
    final result = ResolveAccountsResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  ResolveAccountsResponse._();

  factory ResolveAccountsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResolveAccountsResponse()..mergeFromBuffer(data, registry);
  factory ResolveAccountsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResolveAccountsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResolveAccountsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ResolveAccountsResponse.$_createMessage)
    ..pPM<AccountOwner>(1, _omitFieldNames ? '' : 'data',
        subBuilder: AccountOwner.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResolveAccountsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResolveAccountsResponse copyWith(
          void Function(ResolveAccountsResponse) updates) =>
      super.copyWith((message) => updates(message as ResolveAccountsResponse))
          as ResolveAccountsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ResolveAccountsResponse() / ResolveAccountsResponse.new instead')
  static ResolveAccountsResponse create() => ResolveAccountsResponse._();
  static $pb.GeneratedMessage $_createMessage() => ResolveAccountsResponse._();
  @$core.override
  ResolveAccountsResponse createEmptyInstance() => ResolveAccountsResponse._();
  @$core.pragma('dart2js:noInline')
  static ResolveAccountsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResolveAccountsResponse>(
          ResolveAccountsResponse.$_createMessage);
  static ResolveAccountsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<AccountOwner> get data => $_getList(0);
}

/// The request message containing the information necessary to verify a contact
class CreateContactVerificationRequest extends $pb.GeneratedMessage {
  factory CreateContactVerificationRequest({
    $core.String? id,
    $core.String? contactId,
    $core.String? code,
    $core.String? durationToExpire,
  }) {
    final result = CreateContactVerificationRequest._();
    if (id != null) result.id = id;
    if (contactId != null) result.contactId = contactId;
    if (code != null) result.code = code;
    if (durationToExpire != null) result.durationToExpire = durationToExpire;
    return result;
  }

  CreateContactVerificationRequest._();

  factory CreateContactVerificationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactVerificationRequest()..mergeFromBuffer(data, registry);
  factory CreateContactVerificationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactVerificationRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateContactVerificationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CreateContactVerificationRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'contactId')
    ..aOS(3, _omitFieldNames ? '' : 'code')
    ..aOS(4, _omitFieldNames ? '' : 'durationToExpire')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactVerificationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactVerificationRequest copyWith(
          void Function(CreateContactVerificationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateContactVerificationRequest))
          as CreateContactVerificationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreateContactVerificationRequest() / CreateContactVerificationRequest.new instead')
  static CreateContactVerificationRequest create() =>
      CreateContactVerificationRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      CreateContactVerificationRequest._();
  @$core.override
  CreateContactVerificationRequest createEmptyInstance() =>
      CreateContactVerificationRequest._();
  @$core.pragma('dart2js:noInline')
  static CreateContactVerificationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateContactVerificationRequest>(
          CreateContactVerificationRequest.$_createMessage);
  static CreateContactVerificationRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get contactId => $_getSZ(1);
  @$pb.TagNumber(2)
  set contactId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasContactId() => $_has(1);
  @$pb.TagNumber(2)
  void clearContactId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get code => $_getSZ(2);
  @$pb.TagNumber(3)
  set code($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCode() => $_has(2);
  @$pb.TagNumber(3)
  void clearCode() => $_clearField(3);

  /// This is the string amount of time before code expires e.g. 3m or 500s
  @$pb.TagNumber(4)
  $core.String get durationToExpire => $_getSZ(3);
  @$pb.TagNumber(4)
  set durationToExpire($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDurationToExpire() => $_has(3);
  @$pb.TagNumber(4)
  void clearDurationToExpire() => $_clearField(4);
}

class CreateContactVerificationResponse extends $pb.GeneratedMessage {
  factory CreateContactVerificationResponse({
    $core.String? id,
    $core.bool? success,
  }) {
    final result = CreateContactVerificationResponse._();
    if (id != null) result.id = id;
    if (success != null) result.success = success;
    return result;
  }

  CreateContactVerificationResponse._();

  factory CreateContactVerificationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactVerificationResponse()..mergeFromBuffer(data, registry);
  factory CreateContactVerificationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateContactVerificationResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateContactVerificationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CreateContactVerificationResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOB(2, _omitFieldNames ? '' : 'success')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactVerificationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateContactVerificationResponse copyWith(
          void Function(CreateContactVerificationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateContactVerificationResponse))
          as CreateContactVerificationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreateContactVerificationResponse() / CreateContactVerificationResponse.new instead')
  static CreateContactVerificationResponse create() =>
      CreateContactVerificationResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      CreateContactVerificationResponse._();
  @$core.override
  CreateContactVerificationResponse createEmptyInstance() =>
      CreateContactVerificationResponse._();
  @$core.pragma('dart2js:noInline')
  static CreateContactVerificationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateContactVerificationResponse>(
          CreateContactVerificationResponse.$_createMessage);
  static CreateContactVerificationResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get success => $_getBF(1);
  @$pb.TagNumber(2)
  set success($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSuccess() => $_has(1);
  @$pb.TagNumber(2)
  void clearSuccess() => $_clearField(2);
}

class CheckVerificationRequest extends $pb.GeneratedMessage {
  factory CheckVerificationRequest({
    $core.String? id,
    $core.String? code,
  }) {
    final result = CheckVerificationRequest._();
    if (id != null) result.id = id;
    if (code != null) result.code = code;
    return result;
  }

  CheckVerificationRequest._();

  factory CheckVerificationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckVerificationRequest()..mergeFromBuffer(data, registry);
  factory CheckVerificationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckVerificationRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckVerificationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CheckVerificationRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'code')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckVerificationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckVerificationRequest copyWith(
          void Function(CheckVerificationRequest) updates) =>
      super.copyWith((message) => updates(message as CheckVerificationRequest))
          as CheckVerificationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CheckVerificationRequest() / CheckVerificationRequest.new instead')
  static CheckVerificationRequest create() => CheckVerificationRequest._();
  static $pb.GeneratedMessage $_createMessage() => CheckVerificationRequest._();
  @$core.override
  CheckVerificationRequest createEmptyInstance() =>
      CheckVerificationRequest._();
  @$core.pragma('dart2js:noInline')
  static CheckVerificationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckVerificationRequest>(
          CheckVerificationRequest.$_createMessage);
  static CheckVerificationRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get code => $_getSZ(1);
  @$pb.TagNumber(2)
  set code($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCode() => $_has(1);
  @$pb.TagNumber(2)
  void clearCode() => $_clearField(2);
}

class CheckVerificationResponse extends $pb.GeneratedMessage {
  factory CheckVerificationResponse({
    $core.String? id,
    $core.int? checkAttempts,
    $core.bool? success,
  }) {
    final result = CheckVerificationResponse._();
    if (id != null) result.id = id;
    if (checkAttempts != null) result.checkAttempts = checkAttempts;
    if (success != null) result.success = success;
    return result;
  }

  CheckVerificationResponse._();

  factory CheckVerificationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckVerificationResponse()..mergeFromBuffer(data, registry);
  factory CheckVerificationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckVerificationResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckVerificationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: CheckVerificationResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aI(2, _omitFieldNames ? '' : 'checkAttempts')
    ..aOB(3, _omitFieldNames ? '' : 'success')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckVerificationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckVerificationResponse copyWith(
          void Function(CheckVerificationResponse) updates) =>
      super.copyWith((message) => updates(message as CheckVerificationResponse))
          as CheckVerificationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CheckVerificationResponse() / CheckVerificationResponse.new instead')
  static CheckVerificationResponse create() => CheckVerificationResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      CheckVerificationResponse._();
  @$core.override
  CheckVerificationResponse createEmptyInstance() =>
      CheckVerificationResponse._();
  @$core.pragma('dart2js:noInline')
  static CheckVerificationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckVerificationResponse>(
          CheckVerificationResponse.$_createMessage);
  static CheckVerificationResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get checkAttempts => $_getIZ(1);
  @$pb.TagNumber(2)
  set checkAttempts($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCheckAttempts() => $_has(1);
  @$pb.TagNumber(2)
  void clearCheckAttempts() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.bool get success => $_getBF(2);
  @$pb.TagNumber(3)
  set success($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSuccess() => $_has(2);
  @$pb.TagNumber(3)
  void clearSuccess() => $_clearField(3);
}

/// The request message containing the information necessary to remove a contact
class RemoveContactRequest extends $pb.GeneratedMessage {
  factory RemoveContactRequest({
    $core.String? id,
  }) {
    final result = RemoveContactRequest._();
    if (id != null) result.id = id;
    return result;
  }

  RemoveContactRequest._();

  factory RemoveContactRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveContactRequest()..mergeFromBuffer(data, registry);
  factory RemoveContactRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveContactRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveContactRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RemoveContactRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveContactRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveContactRequest copyWith(void Function(RemoveContactRequest) updates) =>
      super.copyWith((message) => updates(message as RemoveContactRequest))
          as RemoveContactRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RemoveContactRequest() / RemoveContactRequest.new instead')
  static RemoveContactRequest create() => RemoveContactRequest._();
  static $pb.GeneratedMessage $_createMessage() => RemoveContactRequest._();
  @$core.override
  RemoveContactRequest createEmptyInstance() => RemoveContactRequest._();
  @$core.pragma('dart2js:noInline')
  static RemoveContactRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveContactRequest>(
          RemoveContactRequest.$_createMessage);
  static RemoveContactRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class RemoveContactResponse extends $pb.GeneratedMessage {
  factory RemoveContactResponse({
    ProfileObject? data,
  }) {
    final result = RemoveContactResponse._();
    if (data != null) result.data = data;
    return result;
  }

  RemoveContactResponse._();

  factory RemoveContactResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveContactResponse()..mergeFromBuffer(data, registry);
  factory RemoveContactResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveContactResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveContactResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RemoveContactResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveContactResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveContactResponse copyWith(
          void Function(RemoveContactResponse) updates) =>
      super.copyWith((message) => updates(message as RemoveContactResponse))
          as RemoveContactResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RemoveContactResponse() / RemoveContactResponse.new instead')
  static RemoveContactResponse create() => RemoveContactResponse._();
  static $pb.GeneratedMessage $_createMessage() => RemoveContactResponse._();
  @$core.override
  RemoveContactResponse createEmptyInstance() => RemoveContactResponse._();
  @$core.pragma('dart2js:noInline')
  static RemoveContactResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveContactResponse>(
          RemoveContactResponse.$_createMessage);
  static RemoveContactResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

class SearchRosterRequest extends $pb.GeneratedMessage {
  factory SearchRosterRequest({
    $core.String? query,
    $fixnum.Int64? page,
    $core.int? count,
    $core.String? startDate,
    $core.String? endDate,
    $core.Iterable<$core.String>? properties,
    $0.Struct? extras,
    $core.String? profileId,
    $core.String? name,
  }) {
    final result = SearchRosterRequest._();
    if (query != null) result.query = query;
    if (page != null) result.page = page;
    if (count != null) result.count = count;
    if (startDate != null) result.startDate = startDate;
    if (endDate != null) result.endDate = endDate;
    if (properties != null) result.properties.addAll(properties);
    if (extras != null) result.extras = extras;
    if (profileId != null) result.profileId = profileId;
    if (name != null) result.name = name;
    return result;
  }

  SearchRosterRequest._();

  factory SearchRosterRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRosterRequest()..mergeFromBuffer(data, registry);
  factory SearchRosterRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRosterRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchRosterRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: SearchRosterRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'query')
    ..aInt64(2, _omitFieldNames ? '' : 'page')
    ..aI(3, _omitFieldNames ? '' : 'count')
    ..aOS(4, _omitFieldNames ? '' : 'startDate')
    ..aOS(5, _omitFieldNames ? '' : 'endDate')
    ..pPS(6, _omitFieldNames ? '' : 'properties')
    ..aOM<$0.Struct>(7, _omitFieldNames ? '' : 'extras',
        subBuilder: $0.Struct.$_createMessage)
    ..aOS(8, _omitFieldNames ? '' : 'profileId')
    ..aOS(9, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRosterRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRosterRequest copyWith(void Function(SearchRosterRequest) updates) =>
      super.copyWith((message) => updates(message as SearchRosterRequest))
          as SearchRosterRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use SearchRosterRequest() / SearchRosterRequest.new instead')
  static SearchRosterRequest create() => SearchRosterRequest._();
  static $pb.GeneratedMessage $_createMessage() => SearchRosterRequest._();
  @$core.override
  SearchRosterRequest createEmptyInstance() => SearchRosterRequest._();
  @$core.pragma('dart2js:noInline')
  static SearchRosterRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchRosterRequest>(
          SearchRosterRequest.$_createMessage);
  static SearchRosterRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get query => $_getSZ(0);
  @$pb.TagNumber(1)
  set query($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasQuery() => $_has(0);
  @$pb.TagNumber(1)
  void clearQuery() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get page => $_getI64(1);
  @$pb.TagNumber(2)
  set page($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPage() => $_has(1);
  @$pb.TagNumber(2)
  void clearPage() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get count => $_getIZ(2);
  @$pb.TagNumber(3)
  set count($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCount() => $_has(2);
  @$pb.TagNumber(3)
  void clearCount() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get startDate => $_getSZ(3);
  @$pb.TagNumber(4)
  set startDate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasStartDate() => $_has(3);
  @$pb.TagNumber(4)
  void clearStartDate() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get endDate => $_getSZ(4);
  @$pb.TagNumber(5)
  set endDate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasEndDate() => $_has(4);
  @$pb.TagNumber(5)
  void clearEndDate() => $_clearField(5);

  @$pb.TagNumber(6)
  $pb.PbList<$core.String> get properties => $_getList(5);

  @$pb.TagNumber(7)
  $0.Struct get extras => $_getN(6);
  @$pb.TagNumber(7)
  set extras($0.Struct value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasExtras() => $_has(6);
  @$pb.TagNumber(7)
  void clearExtras() => $_clearField(7);
  @$pb.TagNumber(7)
  $0.Struct ensureExtras() => $_ensure(6);

  @$pb.TagNumber(8)
  $core.String get profileId => $_getSZ(7);
  @$pb.TagNumber(8)
  set profileId($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasProfileId() => $_has(7);
  @$pb.TagNumber(8)
  void clearProfileId() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get name => $_getSZ(8);
  @$pb.TagNumber(9)
  set name($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasName() => $_has(8);
  @$pb.TagNumber(9)
  void clearName() => $_clearField(9);
}

class SearchRosterResponse extends $pb.GeneratedMessage {
  factory SearchRosterResponse({
    $core.Iterable<RosterObject>? data,
  }) {
    final result = SearchRosterResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  SearchRosterResponse._();

  factory SearchRosterResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRosterResponse()..mergeFromBuffer(data, registry);
  factory SearchRosterResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SearchRosterResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchRosterResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: SearchRosterResponse.$_createMessage)
    ..pPM<RosterObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RosterObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRosterResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchRosterResponse copyWith(void Function(SearchRosterResponse) updates) =>
      super.copyWith((message) => updates(message as SearchRosterResponse))
          as SearchRosterResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SearchRosterResponse() / SearchRosterResponse.new instead')
  static SearchRosterResponse create() => SearchRosterResponse._();
  static $pb.GeneratedMessage $_createMessage() => SearchRosterResponse._();
  @$core.override
  SearchRosterResponse createEmptyInstance() => SearchRosterResponse._();
  @$core.pragma('dart2js:noInline')
  static SearchRosterResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchRosterResponse>(
          SearchRosterResponse.$_createMessage);
  static SearchRosterResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<RosterObject> get data => $_getList(0);
}

class RawContact extends $pb.GeneratedMessage {
  factory RawContact({
    $core.String? contact,
    $0.Struct? extras,
  }) {
    final result = RawContact._();
    if (contact != null) result.contact = contact;
    if (extras != null) result.extras = extras;
    return result;
  }

  RawContact._();

  factory RawContact.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RawContact()..mergeFromBuffer(data, registry);
  factory RawContact.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RawContact()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RawContact',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RawContact.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'contact')
    ..aOM<$0.Struct>(2, _omitFieldNames ? '' : 'extras',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RawContact clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RawContact copyWith(void Function(RawContact) updates) =>
      super.copyWith((message) => updates(message as RawContact)) as RawContact;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use RawContact() / RawContact.new instead')
  static RawContact create() => RawContact._();
  static $pb.GeneratedMessage $_createMessage() => RawContact._();
  @$core.override
  RawContact createEmptyInstance() => RawContact._();
  @$core.pragma('dart2js:noInline')
  static RawContact getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RawContact>(RawContact.$_createMessage);
  static RawContact? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get contact => $_getSZ(0);
  @$pb.TagNumber(1)
  set contact($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContact() => $_has(0);
  @$pb.TagNumber(1)
  void clearContact() => $_clearField(1);

  @$pb.TagNumber(2)
  $0.Struct get extras => $_getN(1);
  @$pb.TagNumber(2)
  set extras($0.Struct value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasExtras() => $_has(1);
  @$pb.TagNumber(2)
  void clearExtras() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.Struct ensureExtras() => $_ensure(1);
}

class AddRosterRequest extends $pb.GeneratedMessage {
  factory AddRosterRequest({
    $core.Iterable<RawContact>? data,
    $core.String? name,
  }) {
    final result = AddRosterRequest._();
    if (data != null) result.data.addAll(data);
    if (name != null) result.name = name;
    return result;
  }

  AddRosterRequest._();

  factory AddRosterRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRosterRequest()..mergeFromBuffer(data, registry);
  factory AddRosterRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRosterRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddRosterRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddRosterRequest.$_createMessage)
    ..pPM<RawContact>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RawContact.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRosterRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRosterRequest copyWith(void Function(AddRosterRequest) updates) =>
      super.copyWith((message) => updates(message as AddRosterRequest))
          as AddRosterRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddRosterRequest() / AddRosterRequest.new instead')
  static AddRosterRequest create() => AddRosterRequest._();
  static $pb.GeneratedMessage $_createMessage() => AddRosterRequest._();
  @$core.override
  AddRosterRequest createEmptyInstance() => AddRosterRequest._();
  @$core.pragma('dart2js:noInline')
  static AddRosterRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddRosterRequest>(
          AddRosterRequest.$_createMessage);
  static AddRosterRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<RawContact> get data => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);
}

class AddRosterResponse extends $pb.GeneratedMessage {
  factory AddRosterResponse({
    $core.Iterable<RosterObject>? data,
  }) {
    final result = AddRosterResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  AddRosterResponse._();

  factory AddRosterResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRosterResponse()..mergeFromBuffer(data, registry);
  factory AddRosterResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRosterResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddRosterResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddRosterResponse.$_createMessage)
    ..pPM<RosterObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RosterObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRosterResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRosterResponse copyWith(void Function(AddRosterResponse) updates) =>
      super.copyWith((message) => updates(message as AddRosterResponse))
          as AddRosterResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddRosterResponse() / AddRosterResponse.new instead')
  static AddRosterResponse create() => AddRosterResponse._();
  static $pb.GeneratedMessage $_createMessage() => AddRosterResponse._();
  @$core.override
  AddRosterResponse createEmptyInstance() => AddRosterResponse._();
  @$core.pragma('dart2js:noInline')
  static AddRosterResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddRosterResponse>(
          AddRosterResponse.$_createMessage);
  static AddRosterResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<RosterObject> get data => $_getList(0);
}

class RemoveRosterRequest extends $pb.GeneratedMessage {
  factory RemoveRosterRequest({
    $core.String? id,
    $core.String? name,
  }) {
    final result = RemoveRosterRequest._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    return result;
  }

  RemoveRosterRequest._();

  factory RemoveRosterRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveRosterRequest()..mergeFromBuffer(data, registry);
  factory RemoveRosterRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveRosterRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveRosterRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RemoveRosterRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRosterRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRosterRequest copyWith(void Function(RemoveRosterRequest) updates) =>
      super.copyWith((message) => updates(message as RemoveRosterRequest))
          as RemoveRosterRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use RemoveRosterRequest() / RemoveRosterRequest.new instead')
  static RemoveRosterRequest create() => RemoveRosterRequest._();
  static $pb.GeneratedMessage $_createMessage() => RemoveRosterRequest._();
  @$core.override
  RemoveRosterRequest createEmptyInstance() => RemoveRosterRequest._();
  @$core.pragma('dart2js:noInline')
  static RemoveRosterRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveRosterRequest>(
          RemoveRosterRequest.$_createMessage);
  static RemoveRosterRequest? _defaultInstance;

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
}

class RemoveRosterResponse extends $pb.GeneratedMessage {
  factory RemoveRosterResponse({
    RosterObject? roster,
  }) {
    final result = RemoveRosterResponse._();
    if (roster != null) result.roster = roster;
    return result;
  }

  RemoveRosterResponse._();

  factory RemoveRosterResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveRosterResponse()..mergeFromBuffer(data, registry);
  factory RemoveRosterResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveRosterResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveRosterResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: RemoveRosterResponse.$_createMessage)
    ..aOM<RosterObject>(1, _omitFieldNames ? '' : 'roster',
        subBuilder: RosterObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRosterResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRosterResponse copyWith(void Function(RemoveRosterResponse) updates) =>
      super.copyWith((message) => updates(message as RemoveRosterResponse))
          as RemoveRosterResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RemoveRosterResponse() / RemoveRosterResponse.new instead')
  static RemoveRosterResponse create() => RemoveRosterResponse._();
  static $pb.GeneratedMessage $_createMessage() => RemoveRosterResponse._();
  @$core.override
  RemoveRosterResponse createEmptyInstance() => RemoveRosterResponse._();
  @$core.pragma('dart2js:noInline')
  static RemoveRosterResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveRosterResponse>(
          RemoveRosterResponse.$_createMessage);
  static RemoveRosterResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RosterObject get roster => $_getN(0);
  @$pb.TagNumber(1)
  set roster(RosterObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasRoster() => $_has(0);
  @$pb.TagNumber(1)
  void clearRoster() => $_clearField(1);
  @$pb.TagNumber(1)
  RosterObject ensureRoster() => $_ensure(0);
}

/// The request message containing the information necessary to create a profile
class AddAddressRequest extends $pb.GeneratedMessage {
  factory AddAddressRequest({
    $core.String? id,
    AddressObject? address,
  }) {
    final result = AddAddressRequest._();
    if (id != null) result.id = id;
    if (address != null) result.address = address;
    return result;
  }

  AddAddressRequest._();

  factory AddAddressRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddAddressRequest()..mergeFromBuffer(data, registry);
  factory AddAddressRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddAddressRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddAddressRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddAddressRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOM<AddressObject>(2, _omitFieldNames ? '' : 'address',
        subBuilder: AddressObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddAddressRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddAddressRequest copyWith(void Function(AddAddressRequest) updates) =>
      super.copyWith((message) => updates(message as AddAddressRequest))
          as AddAddressRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddAddressRequest() / AddAddressRequest.new instead')
  static AddAddressRequest create() => AddAddressRequest._();
  static $pb.GeneratedMessage $_createMessage() => AddAddressRequest._();
  @$core.override
  AddAddressRequest createEmptyInstance() => AddAddressRequest._();
  @$core.pragma('dart2js:noInline')
  static AddAddressRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddAddressRequest>(
          AddAddressRequest.$_createMessage);
  static AddAddressRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  AddressObject get address => $_getN(1);
  @$pb.TagNumber(2)
  set address(AddressObject value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);
  @$pb.TagNumber(2)
  AddressObject ensureAddress() => $_ensure(1);
}

class AddAddressResponse extends $pb.GeneratedMessage {
  factory AddAddressResponse({
    ProfileObject? data,
  }) {
    final result = AddAddressResponse._();
    if (data != null) result.data = data;
    return result;
  }

  AddAddressResponse._();

  factory AddAddressResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddAddressResponse()..mergeFromBuffer(data, registry);
  factory AddAddressResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddAddressResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddAddressResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddAddressResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddAddressResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddAddressResponse copyWith(void Function(AddAddressResponse) updates) =>
      super.copyWith((message) => updates(message as AddAddressResponse))
          as AddAddressResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddAddressResponse() / AddAddressResponse.new instead')
  static AddAddressResponse create() => AddAddressResponse._();
  static $pb.GeneratedMessage $_createMessage() => AddAddressResponse._();
  @$core.override
  AddAddressResponse createEmptyInstance() => AddAddressResponse._();
  @$core.pragma('dart2js:noInline')
  static AddAddressResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddAddressResponse>(
          AddAddressResponse.$_createMessage);
  static AddAddressResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

/// Request message containing a contact necessary to link to a profile
class GetByContactRequest extends $pb.GeneratedMessage {
  factory GetByContactRequest({
    $core.String? contact,
  }) {
    final result = GetByContactRequest._();
    if (contact != null) result.contact = contact;
    return result;
  }

  GetByContactRequest._();

  factory GetByContactRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByContactRequest()..mergeFromBuffer(data, registry);
  factory GetByContactRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByContactRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByContactRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetByContactRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'contact')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByContactRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByContactRequest copyWith(void Function(GetByContactRequest) updates) =>
      super.copyWith((message) => updates(message as GetByContactRequest))
          as GetByContactRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use GetByContactRequest() / GetByContactRequest.new instead')
  static GetByContactRequest create() => GetByContactRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetByContactRequest._();
  @$core.override
  GetByContactRequest createEmptyInstance() => GetByContactRequest._();
  @$core.pragma('dart2js:noInline')
  static GetByContactRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetByContactRequest>(
          GetByContactRequest.$_createMessage);
  static GetByContactRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get contact => $_getSZ(0);
  @$pb.TagNumber(1)
  set contact($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContact() => $_has(0);
  @$pb.TagNumber(1)
  void clearContact() => $_clearField(1);
}

class GetByContactResponse extends $pb.GeneratedMessage {
  factory GetByContactResponse({
    ProfileObject? data,
  }) {
    final result = GetByContactResponse._();
    if (data != null) result.data = data;
    return result;
  }

  GetByContactResponse._();

  factory GetByContactResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByContactResponse()..mergeFromBuffer(data, registry);
  factory GetByContactResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByContactResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByContactResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetByContactResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByContactResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByContactResponse copyWith(void Function(GetByContactResponse) updates) =>
      super.copyWith((message) => updates(message as GetByContactResponse))
          as GetByContactResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetByContactResponse() / GetByContactResponse.new instead')
  static GetByContactResponse create() => GetByContactResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetByContactResponse._();
  @$core.override
  GetByContactResponse createEmptyInstance() => GetByContactResponse._();
  @$core.pragma('dart2js:noInline')
  static GetByContactResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetByContactResponse>(
          GetByContactResponse.$_createMessage);
  static GetByContactResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

class GetByIDAndPartitionRequest extends $pb.GeneratedMessage {
  factory GetByIDAndPartitionRequest({
    $core.String? id,
    $core.String? partitionId,
  }) {
    final result = GetByIDAndPartitionRequest._();
    if (id != null) result.id = id;
    if (partitionId != null) result.partitionId = partitionId;
    return result;
  }

  GetByIDAndPartitionRequest._();

  factory GetByIDAndPartitionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIDAndPartitionRequest()..mergeFromBuffer(data, registry);
  factory GetByIDAndPartitionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIDAndPartitionRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByIDAndPartitionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetByIDAndPartitionRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'partitionId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIDAndPartitionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIDAndPartitionRequest copyWith(
          void Function(GetByIDAndPartitionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetByIDAndPartitionRequest))
          as GetByIDAndPartitionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetByIDAndPartitionRequest() / GetByIDAndPartitionRequest.new instead')
  static GetByIDAndPartitionRequest create() => GetByIDAndPartitionRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetByIDAndPartitionRequest._();
  @$core.override
  GetByIDAndPartitionRequest createEmptyInstance() =>
      GetByIDAndPartitionRequest._();
  @$core.pragma('dart2js:noInline')
  static GetByIDAndPartitionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetByIDAndPartitionRequest>(
          GetByIDAndPartitionRequest.$_createMessage);
  static GetByIDAndPartitionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get partitionId => $_getSZ(1);
  @$pb.TagNumber(2)
  set partitionId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPartitionId() => $_has(1);
  @$pb.TagNumber(2)
  void clearPartitionId() => $_clearField(2);
}

class GetByIDAndPartitionResponse extends $pb.GeneratedMessage {
  factory GetByIDAndPartitionResponse({
    ProfileObject? data,
  }) {
    final result = GetByIDAndPartitionResponse._();
    if (data != null) result.data = data;
    return result;
  }

  GetByIDAndPartitionResponse._();

  factory GetByIDAndPartitionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIDAndPartitionResponse()..mergeFromBuffer(data, registry);
  factory GetByIDAndPartitionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetByIDAndPartitionResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByIDAndPartitionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: GetByIDAndPartitionResponse.$_createMessage)
    ..aOM<ProfileObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: ProfileObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIDAndPartitionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByIDAndPartitionResponse copyWith(
          void Function(GetByIDAndPartitionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetByIDAndPartitionResponse))
          as GetByIDAndPartitionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetByIDAndPartitionResponse() / GetByIDAndPartitionResponse.new instead')
  static GetByIDAndPartitionResponse create() =>
      GetByIDAndPartitionResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetByIDAndPartitionResponse._();
  @$core.override
  GetByIDAndPartitionResponse createEmptyInstance() =>
      GetByIDAndPartitionResponse._();
  @$core.pragma('dart2js:noInline')
  static GetByIDAndPartitionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetByIDAndPartitionResponse>(
          GetByIDAndPartitionResponse.$_createMessage);
  static GetByIDAndPartitionResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProfileObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(ProfileObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  ProfileObject ensureData() => $_ensure(0);
}

class PropertyHistoryRequest extends $pb.GeneratedMessage {
  factory PropertyHistoryRequest({
    $core.String? id,
    $core.String? key,
  }) {
    final result = PropertyHistoryRequest._();
    if (id != null) result.id = id;
    if (key != null) result.key = key;
    return result;
  }

  PropertyHistoryRequest._();

  factory PropertyHistoryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PropertyHistoryRequest()..mergeFromBuffer(data, registry);
  factory PropertyHistoryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PropertyHistoryRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PropertyHistoryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: PropertyHistoryRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'key')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PropertyHistoryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PropertyHistoryRequest copyWith(
          void Function(PropertyHistoryRequest) updates) =>
      super.copyWith((message) => updates(message as PropertyHistoryRequest))
          as PropertyHistoryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use PropertyHistoryRequest() / PropertyHistoryRequest.new instead')
  static PropertyHistoryRequest create() => PropertyHistoryRequest._();
  static $pb.GeneratedMessage $_createMessage() => PropertyHistoryRequest._();
  @$core.override
  PropertyHistoryRequest createEmptyInstance() => PropertyHistoryRequest._();
  @$core.pragma('dart2js:noInline')
  static PropertyHistoryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PropertyHistoryRequest>(
          PropertyHistoryRequest.$_createMessage);
  static PropertyHistoryRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get key => $_getSZ(1);
  @$pb.TagNumber(2)
  set key($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasKey() => $_has(1);
  @$pb.TagNumber(2)
  void clearKey() => $_clearField(2);
}

class PropertyEntryObject extends $pb.GeneratedMessage {
  factory PropertyEntryObject({
    $core.String? key,
    $core.String? value,
    $core.String? tenantId,
    $core.String? createdBy,
    $1.Timestamp? createdAt,
    $core.bool? scoped,
  }) {
    final result = PropertyEntryObject._();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    if (tenantId != null) result.tenantId = tenantId;
    if (createdBy != null) result.createdBy = createdBy;
    if (createdAt != null) result.createdAt = createdAt;
    if (scoped != null) result.scoped = scoped;
    return result;
  }

  PropertyEntryObject._();

  factory PropertyEntryObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PropertyEntryObject()..mergeFromBuffer(data, registry);
  factory PropertyEntryObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PropertyEntryObject()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PropertyEntryObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: PropertyEntryObject.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'key')
    ..aOS(2, _omitFieldNames ? '' : 'value')
    ..aOS(3, _omitFieldNames ? '' : 'tenantId')
    ..aOS(4, _omitFieldNames ? '' : 'createdBy')
    ..aOM<$1.Timestamp>(5, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOB(6, _omitFieldNames ? '' : 'scoped')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PropertyEntryObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PropertyEntryObject copyWith(void Function(PropertyEntryObject) updates) =>
      super.copyWith((message) => updates(message as PropertyEntryObject))
          as PropertyEntryObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use PropertyEntryObject() / PropertyEntryObject.new instead')
  static PropertyEntryObject create() => PropertyEntryObject._();
  static $pb.GeneratedMessage $_createMessage() => PropertyEntryObject._();
  @$core.override
  PropertyEntryObject createEmptyInstance() => PropertyEntryObject._();
  @$core.pragma('dart2js:noInline')
  static PropertyEntryObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PropertyEntryObject>(
          PropertyEntryObject.$_createMessage);
  static PropertyEntryObject? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearKey() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(2)
  void clearValue() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get tenantId => $_getSZ(2);
  @$pb.TagNumber(3)
  set tenantId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTenantId() => $_has(2);
  @$pb.TagNumber(3)
  void clearTenantId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get createdBy => $_getSZ(3);
  @$pb.TagNumber(4)
  set createdBy($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCreatedBy() => $_has(3);
  @$pb.TagNumber(4)
  void clearCreatedBy() => $_clearField(4);

  @$pb.TagNumber(5)
  $1.Timestamp get createdAt => $_getN(4);
  @$pb.TagNumber(5)
  set createdAt($1.Timestamp value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasCreatedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearCreatedAt() => $_clearField(5);
  @$pb.TagNumber(5)
  $1.Timestamp ensureCreatedAt() => $_ensure(4);

  @$pb.TagNumber(6)
  $core.bool get scoped => $_getBF(5);
  @$pb.TagNumber(6)
  set scoped($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasScoped() => $_has(5);
  @$pb.TagNumber(6)
  void clearScoped() => $_clearField(6);
}

class PropertyHistoryResponse extends $pb.GeneratedMessage {
  factory PropertyHistoryResponse({
    $core.Iterable<PropertyEntryObject>? entries,
  }) {
    final result = PropertyHistoryResponse._();
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  PropertyHistoryResponse._();

  factory PropertyHistoryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PropertyHistoryResponse()..mergeFromBuffer(data, registry);
  factory PropertyHistoryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PropertyHistoryResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PropertyHistoryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: PropertyHistoryResponse.$_createMessage)
    ..pPM<PropertyEntryObject>(1, _omitFieldNames ? '' : 'entries',
        subBuilder: PropertyEntryObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PropertyHistoryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PropertyHistoryResponse copyWith(
          void Function(PropertyHistoryResponse) updates) =>
      super.copyWith((message) => updates(message as PropertyHistoryResponse))
          as PropertyHistoryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use PropertyHistoryResponse() / PropertyHistoryResponse.new instead')
  static PropertyHistoryResponse create() => PropertyHistoryResponse._();
  static $pb.GeneratedMessage $_createMessage() => PropertyHistoryResponse._();
  @$core.override
  PropertyHistoryResponse createEmptyInstance() => PropertyHistoryResponse._();
  @$core.pragma('dart2js:noInline')
  static PropertyHistoryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PropertyHistoryResponse>(
          PropertyHistoryResponse.$_createMessage);
  static PropertyHistoryResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<PropertyEntryObject> get entries => $_getList(0);
}

/// Request message to list all profiles related to the supplied profile id
class ListRelationshipRequest extends $pb.GeneratedMessage {
  factory ListRelationshipRequest({
    $core.String? peerName,
    $core.String? peerId,
    $core.String? lastRelationshipId,
    $core.Iterable<$core.String>? relatedChildrenId,
    $core.int? count,
    $core.bool? invertRelation,
  }) {
    final result = ListRelationshipRequest._();
    if (peerName != null) result.peerName = peerName;
    if (peerId != null) result.peerId = peerId;
    if (lastRelationshipId != null)
      result.lastRelationshipId = lastRelationshipId;
    if (relatedChildrenId != null)
      result.relatedChildrenId.addAll(relatedChildrenId);
    if (count != null) result.count = count;
    if (invertRelation != null) result.invertRelation = invertRelation;
    return result;
  }

  ListRelationshipRequest._();

  factory ListRelationshipRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListRelationshipRequest()..mergeFromBuffer(data, registry);
  factory ListRelationshipRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListRelationshipRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRelationshipRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ListRelationshipRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'peerName')
    ..aOS(2, _omitFieldNames ? '' : 'peerId')
    ..aOS(3, _omitFieldNames ? '' : 'lastRelationshipId')
    ..pPS(4, _omitFieldNames ? '' : 'relatedChildrenId')
    ..aI(5, _omitFieldNames ? '' : 'count')
    ..aOB(6, _omitFieldNames ? '' : 'invertRelation')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRelationshipRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRelationshipRequest copyWith(
          void Function(ListRelationshipRequest) updates) =>
      super.copyWith((message) => updates(message as ListRelationshipRequest))
          as ListRelationshipRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListRelationshipRequest() / ListRelationshipRequest.new instead')
  static ListRelationshipRequest create() => ListRelationshipRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListRelationshipRequest._();
  @$core.override
  ListRelationshipRequest createEmptyInstance() => ListRelationshipRequest._();
  @$core.pragma('dart2js:noInline')
  static ListRelationshipRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRelationshipRequest>(
          ListRelationshipRequest.$_createMessage);
  static ListRelationshipRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get peerName => $_getSZ(0);
  @$pb.TagNumber(1)
  set peerName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPeerName() => $_has(0);
  @$pb.TagNumber(1)
  void clearPeerName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get peerId => $_getSZ(1);
  @$pb.TagNumber(2)
  set peerId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPeerId() => $_has(1);
  @$pb.TagNumber(2)
  void clearPeerId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get lastRelationshipId => $_getSZ(2);
  @$pb.TagNumber(3)
  set lastRelationshipId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLastRelationshipId() => $_has(2);
  @$pb.TagNumber(3)
  void clearLastRelationshipId() => $_clearField(3);

  @$pb.TagNumber(4)
  $pb.PbList<$core.String> get relatedChildrenId => $_getList(3);

  @$pb.TagNumber(5)
  $core.int get count => $_getIZ(4);
  @$pb.TagNumber(5)
  set count($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasCount() => $_has(4);
  @$pb.TagNumber(5)
  void clearCount() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.bool get invertRelation => $_getBF(5);
  @$pb.TagNumber(6)
  set invertRelation($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasInvertRelation() => $_has(5);
  @$pb.TagNumber(6)
  void clearInvertRelation() => $_clearField(6);
}

class ListRelationshipResponse extends $pb.GeneratedMessage {
  factory ListRelationshipResponse({
    $core.Iterable<RelationshipObject>? data,
  }) {
    final result = ListRelationshipResponse._();
    if (data != null) result.data.addAll(data);
    return result;
  }

  ListRelationshipResponse._();

  factory ListRelationshipResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListRelationshipResponse()..mergeFromBuffer(data, registry);
  factory ListRelationshipResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListRelationshipResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRelationshipResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: ListRelationshipResponse.$_createMessage)
    ..pPM<RelationshipObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RelationshipObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRelationshipResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRelationshipResponse copyWith(
          void Function(ListRelationshipResponse) updates) =>
      super.copyWith((message) => updates(message as ListRelationshipResponse))
          as ListRelationshipResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListRelationshipResponse() / ListRelationshipResponse.new instead')
  static ListRelationshipResponse create() => ListRelationshipResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListRelationshipResponse._();
  @$core.override
  ListRelationshipResponse createEmptyInstance() =>
      ListRelationshipResponse._();
  @$core.pragma('dart2js:noInline')
  static ListRelationshipResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRelationshipResponse>(
          ListRelationshipResponse.$_createMessage);
  static ListRelationshipResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<RelationshipObject> get data => $_getList(0);
}

/// The request message containing the information necessary to create a profile
class AddRelationshipRequest extends $pb.GeneratedMessage {
  factory AddRelationshipRequest({
    $core.String? id,
    $core.String? parent,
    $core.String? parentId,
    $core.String? child,
    $core.String? childId,
    RelationshipType? type,
    $0.Struct? properties,
  }) {
    final result = AddRelationshipRequest._();
    if (id != null) result.id = id;
    if (parent != null) result.parent = parent;
    if (parentId != null) result.parentId = parentId;
    if (child != null) result.child = child;
    if (childId != null) result.childId = childId;
    if (type != null) result.type = type;
    if (properties != null) result.properties = properties;
    return result;
  }

  AddRelationshipRequest._();

  factory AddRelationshipRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRelationshipRequest()..mergeFromBuffer(data, registry);
  factory AddRelationshipRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRelationshipRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddRelationshipRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddRelationshipRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'parent')
    ..aOS(3, _omitFieldNames ? '' : 'parentId')
    ..aOS(4, _omitFieldNames ? '' : 'child')
    ..aOS(5, _omitFieldNames ? '' : 'childId')
    ..aE<RelationshipType>(6, _omitFieldNames ? '' : 'type',
        enumValues: RelationshipType.values)
    ..aOM<$0.Struct>(7, _omitFieldNames ? '' : 'properties',
        subBuilder: $0.Struct.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRelationshipRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRelationshipRequest copyWith(
          void Function(AddRelationshipRequest) updates) =>
      super.copyWith((message) => updates(message as AddRelationshipRequest))
          as AddRelationshipRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AddRelationshipRequest() / AddRelationshipRequest.new instead')
  static AddRelationshipRequest create() => AddRelationshipRequest._();
  static $pb.GeneratedMessage $_createMessage() => AddRelationshipRequest._();
  @$core.override
  AddRelationshipRequest createEmptyInstance() => AddRelationshipRequest._();
  @$core.pragma('dart2js:noInline')
  static AddRelationshipRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddRelationshipRequest>(
          AddRelationshipRequest.$_createMessage);
  static AddRelationshipRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get parent => $_getSZ(1);
  @$pb.TagNumber(2)
  set parent($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasParent() => $_has(1);
  @$pb.TagNumber(2)
  void clearParent() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get parentId => $_getSZ(2);
  @$pb.TagNumber(3)
  set parentId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasParentId() => $_has(2);
  @$pb.TagNumber(3)
  void clearParentId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get child => $_getSZ(3);
  @$pb.TagNumber(4)
  set child($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasChild() => $_has(3);
  @$pb.TagNumber(4)
  void clearChild() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get childId => $_getSZ(4);
  @$pb.TagNumber(5)
  set childId($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasChildId() => $_has(4);
  @$pb.TagNumber(5)
  void clearChildId() => $_clearField(5);

  @$pb.TagNumber(6)
  RelationshipType get type => $_getN(5);
  @$pb.TagNumber(6)
  set type(RelationshipType value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasType() => $_has(5);
  @$pb.TagNumber(6)
  void clearType() => $_clearField(6);

  @$pb.TagNumber(7)
  $0.Struct get properties => $_getN(6);
  @$pb.TagNumber(7)
  set properties($0.Struct value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasProperties() => $_has(6);
  @$pb.TagNumber(7)
  void clearProperties() => $_clearField(7);
  @$pb.TagNumber(7)
  $0.Struct ensureProperties() => $_ensure(6);
}

class AddRelationshipResponse extends $pb.GeneratedMessage {
  factory AddRelationshipResponse({
    RelationshipObject? data,
  }) {
    final result = AddRelationshipResponse._();
    if (data != null) result.data = data;
    return result;
  }

  AddRelationshipResponse._();

  factory AddRelationshipResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRelationshipResponse()..mergeFromBuffer(data, registry);
  factory AddRelationshipResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddRelationshipResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddRelationshipResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: AddRelationshipResponse.$_createMessage)
    ..aOM<RelationshipObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RelationshipObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRelationshipResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddRelationshipResponse copyWith(
          void Function(AddRelationshipResponse) updates) =>
      super.copyWith((message) => updates(message as AddRelationshipResponse))
          as AddRelationshipResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use AddRelationshipResponse() / AddRelationshipResponse.new instead')
  static AddRelationshipResponse create() => AddRelationshipResponse._();
  static $pb.GeneratedMessage $_createMessage() => AddRelationshipResponse._();
  @$core.override
  AddRelationshipResponse createEmptyInstance() => AddRelationshipResponse._();
  @$core.pragma('dart2js:noInline')
  static AddRelationshipResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddRelationshipResponse>(
          AddRelationshipResponse.$_createMessage);
  static AddRelationshipResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RelationshipObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RelationshipObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RelationshipObject ensureData() => $_ensure(0);
}

/// The request message containing the information necessary to delete relationship to a profile
class DeleteRelationshipRequest extends $pb.GeneratedMessage {
  factory DeleteRelationshipRequest({
    $core.String? id,
    $core.String? parentId,
  }) {
    final result = DeleteRelationshipRequest._();
    if (id != null) result.id = id;
    if (parentId != null) result.parentId = parentId;
    return result;
  }

  DeleteRelationshipRequest._();

  factory DeleteRelationshipRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteRelationshipRequest()..mergeFromBuffer(data, registry);
  factory DeleteRelationshipRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteRelationshipRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRelationshipRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: DeleteRelationshipRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'parentId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRelationshipRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRelationshipRequest copyWith(
          void Function(DeleteRelationshipRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteRelationshipRequest))
          as DeleteRelationshipRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DeleteRelationshipRequest() / DeleteRelationshipRequest.new instead')
  static DeleteRelationshipRequest create() => DeleteRelationshipRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      DeleteRelationshipRequest._();
  @$core.override
  DeleteRelationshipRequest createEmptyInstance() =>
      DeleteRelationshipRequest._();
  @$core.pragma('dart2js:noInline')
  static DeleteRelationshipRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRelationshipRequest>(
          DeleteRelationshipRequest.$_createMessage);
  static DeleteRelationshipRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get parentId => $_getSZ(1);
  @$pb.TagNumber(2)
  set parentId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasParentId() => $_has(1);
  @$pb.TagNumber(2)
  void clearParentId() => $_clearField(2);
}

class DeleteRelationshipResponse extends $pb.GeneratedMessage {
  factory DeleteRelationshipResponse({
    RelationshipObject? data,
  }) {
    final result = DeleteRelationshipResponse._();
    if (data != null) result.data = data;
    return result;
  }

  DeleteRelationshipResponse._();

  factory DeleteRelationshipResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteRelationshipResponse()..mergeFromBuffer(data, registry);
  factory DeleteRelationshipResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteRelationshipResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRelationshipResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'profile.v1'),
      createEmptyInstance: DeleteRelationshipResponse.$_createMessage)
    ..aOM<RelationshipObject>(1, _omitFieldNames ? '' : 'data',
        subBuilder: RelationshipObject.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRelationshipResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRelationshipResponse copyWith(
          void Function(DeleteRelationshipResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRelationshipResponse))
          as DeleteRelationshipResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DeleteRelationshipResponse() / DeleteRelationshipResponse.new instead')
  static DeleteRelationshipResponse create() => DeleteRelationshipResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      DeleteRelationshipResponse._();
  @$core.override
  DeleteRelationshipResponse createEmptyInstance() =>
      DeleteRelationshipResponse._();
  @$core.pragma('dart2js:noInline')
  static DeleteRelationshipResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRelationshipResponse>(
          DeleteRelationshipResponse.$_createMessage);
  static DeleteRelationshipResponse? _defaultInstance;

  @$pb.TagNumber(1)
  RelationshipObject get data => $_getN(0);
  @$pb.TagNumber(1)
  set data(RelationshipObject value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => $_clearField(1);
  @$pb.TagNumber(1)
  RelationshipObject ensureData() => $_ensure(0);
}

/// ProfileService manages user and entity profiles.
/// All RPCs require authentication via Bearer token.
class ProfileServiceApi {
  final $pb.RpcClient _client;

  ProfileServiceApi(this._client);

  /// GetById retrieves a profile by its unique ID.
  $async.Future<GetByIdResponse> getById(
          $pb.ClientContext? ctx, GetByIdRequest request) =>
      _client.invoke<GetByIdResponse>(
          ctx, 'ProfileService', 'GetById', request, GetByIdResponse());

  /// GetByContact retrieves a profile by contact information.
  $async.Future<GetByContactResponse> getByContact(
          $pb.ClientContext? ctx, GetByContactRequest request) =>
      _client.invoke<GetByContactResponse>(ctx, 'ProfileService',
          'GetByContact', request, GetByContactResponse());

  /// Search finds profiles matching specified criteria.
  $async.Future<SearchResponse> search(
          $pb.ClientContext? ctx, SearchRequest request) =>
      _client.invoke<SearchResponse>(
          ctx, 'ProfileService', 'Search', request, SearchResponse());

  /// Merge combines two profiles into one.
  $async.Future<MergeResponse> merge(
          $pb.ClientContext? ctx, MergeRequest request) =>
      _client.invoke<MergeResponse>(
          ctx, 'ProfileService', 'Merge', request, MergeResponse());

  /// Create creates a new profile.
  $async.Future<CreateResponse> create_(
          $pb.ClientContext? ctx, CreateRequest request) =>
      _client.invoke<CreateResponse>(
          ctx, 'ProfileService', 'Create', request, CreateResponse());

  /// Update updates an existing profile's properties.
  $async.Future<UpdateResponse> update(
          $pb.ClientContext? ctx, UpdateRequest request) =>
      _client.invoke<UpdateResponse>(
          ctx, 'ProfileService', 'Update', request, UpdateResponse());

  /// AddContact adds a new contact to a profile with automatic verification.
  $async.Future<AddContactResponse> addContact(
          $pb.ClientContext? ctx, AddContactRequest request) =>
      _client.invoke<AddContactResponse>(
          ctx, 'ProfileService', 'AddContact', request, AddContactResponse());

  /// CreateContact creates a standalone contact not linked to a profile.
  $async.Future<CreateContactResponse> createContact(
          $pb.ClientContext? ctx, CreateContactRequest request) =>
      _client.invoke<CreateContactResponse>(ctx, 'ProfileService',
          'CreateContact', request, CreateContactResponse());

  /// GetContacts resolves one or many contacts by id, whether standalone or
  /// attached to a profile. Missing ids are listed in missing_ids; found
  /// contacts are returned in data (order not guaranteed).
  $async.Future<GetContactsResponse> getContacts(
          $pb.ClientContext? ctx, GetContactsRequest request) =>
      _client.invoke<GetContactsResponse>(
          ctx, 'ProfileService', 'GetContacts', request, GetContactsResponse());

  /// CreateContactVerification initiates contact verification.
  $async.Future<CreateContactVerificationResponse> createContactVerification(
          $pb.ClientContext? ctx, CreateContactVerificationRequest request) =>
      _client.invoke<CreateContactVerificationResponse>(
          ctx,
          'ProfileService',
          'CreateContactVerification',
          request,
          CreateContactVerificationResponse());

  /// CheckVerification verifies a contact using the provided code.
  $async.Future<CheckVerificationResponse> checkVerification(
          $pb.ClientContext? ctx, CheckVerificationRequest request) =>
      _client.invoke<CheckVerificationResponse>(ctx, 'ProfileService',
          'CheckVerification', request, CheckVerificationResponse());

  /// RemoveContact removes a contact from a profile.
  $async.Future<RemoveContactResponse> removeContact(
          $pb.ClientContext? ctx, RemoveContactRequest request) =>
      _client.invoke<RemoveContactResponse>(ctx, 'ProfileService',
          'RemoveContact', request, RemoveContactResponse());

  /// SearchRoster searches a user's contact roster.
  $async.Future<SearchRosterResponse> searchRoster(
          $pb.ClientContext? ctx, SearchRosterRequest request) =>
      _client.invoke<SearchRosterResponse>(ctx, 'ProfileService',
          'SearchRoster', request, SearchRosterResponse());

  /// AddRoster adds multiple contacts to a user's roster.
  $async.Future<AddRosterResponse> addRoster(
          $pb.ClientContext? ctx, AddRosterRequest request) =>
      _client.invoke<AddRosterResponse>(
          ctx, 'ProfileService', 'AddRoster', request, AddRosterResponse());

  /// RemoveRoster removes a contact from a user's roster.
  $async.Future<RemoveRosterResponse> removeRoster(
          $pb.ClientContext? ctx, RemoveRosterRequest request) =>
      _client.invoke<RemoveRosterResponse>(ctx, 'ProfileService',
          'RemoveRoster', request, RemoveRosterResponse());

  /// AddAddress adds a new address to a profile.
  $async.Future<AddAddressResponse> addAddress(
          $pb.ClientContext? ctx, AddAddressRequest request) =>
      _client.invoke<AddAddressResponse>(
          ctx, 'ProfileService', 'AddAddress', request, AddAddressResponse());

  /// AddRelationship creates a relationship between profiles.
  $async.Future<AddRelationshipResponse> addRelationship(
          $pb.ClientContext? ctx, AddRelationshipRequest request) =>
      _client.invoke<AddRelationshipResponse>(ctx, 'ProfileService',
          'AddRelationship', request, AddRelationshipResponse());

  /// DeleteRelationship removes a relationship between profiles.
  $async.Future<DeleteRelationshipResponse> deleteRelationship(
          $pb.ClientContext? ctx, DeleteRelationshipRequest request) =>
      _client.invoke<DeleteRelationshipResponse>(ctx, 'ProfileService',
          'DeleteRelationship', request, DeleteRelationshipResponse());

  /// ListRelationship lists all relationships for a profile.
  $async.Future<ListRelationshipResponse> listRelationship(
          $pb.ClientContext? ctx, ListRelationshipRequest request) =>
      _client.invoke<ListRelationshipResponse>(ctx, 'ProfileService',
          'ListRelationship', request, ListRelationshipResponse());

  /// GetByIDAndPartition retrieves a profile with partition-scoped properties merged in.
  $async.Future<GetByIDAndPartitionResponse> getByIDAndPartition(
          $pb.ClientContext? ctx, GetByIDAndPartitionRequest request) =>
      _client.invoke<GetByIDAndPartitionResponse>(ctx, 'ProfileService',
          'GetByIDAndPartition', request, GetByIDAndPartitionResponse());

  /// PropertyHistory returns the change history for a specific property key on a profile.
  $async.Future<PropertyHistoryResponse> propertyHistory(
          $pb.ClientContext? ctx, PropertyHistoryRequest request) =>
      _client.invoke<PropertyHistoryResponse>(ctx, 'ProfileService',
          'PropertyHistory', request, PropertyHistoryResponse());

  /// ResolveAccounts maps chain account addresses to the profiles that own
  /// them. Service-to-service only (account_resolve is bound to ROLE_SERVICE
  /// alone); unknown addresses are omitted from the response.
  $async.Future<ResolveAccountsResponse> resolveAccounts(
          $pb.ClientContext? ctx, ResolveAccountsRequest request) =>
      _client.invoke<ResolveAccountsResponse>(ctx, 'ProfileService',
          'ResolveAccounts', request, ResolveAccountsResponse());
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
