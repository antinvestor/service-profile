//
//  Generated code. Do not modify.
//  source: profile/v1/profile.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import '../../google/protobuf/struct.pbjson.dart' as $6;
import '../../google/protobuf/timestamp.pbjson.dart' as $2;

@$core.Deprecated('Use contactTypeDescriptor instead')
const ContactType$json = {
  '1': 'ContactType',
  '2': [
    {'1': 'EMAIL', '2': 0},
    {'1': 'MSISDN', '2': 1},
  ],
};

/// Descriptor for `ContactType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contactTypeDescriptor = $convert.base64Decode(
    'CgtDb250YWN0VHlwZRIJCgVFTUFJTBAAEgoKBk1TSVNEThAB');

@$core.Deprecated('Use communicationLevelDescriptor instead')
const CommunicationLevel$json = {
  '1': 'CommunicationLevel',
  '2': [
    {'1': 'ALL', '2': 0},
    {'1': 'INTERNAL_MARKETING', '2': 1},
    {'1': 'IMPORTANT_ALERTS', '2': 2},
    {'1': 'SYSTEM_ALERTS', '2': 3},
    {'1': 'NO_CONTACT', '2': 4},
  ],
};

/// Descriptor for `CommunicationLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List communicationLevelDescriptor = $convert.base64Decode(
    'ChJDb21tdW5pY2F0aW9uTGV2ZWwSBwoDQUxMEAASFgoSSU5URVJOQUxfTUFSS0VUSU5HEAESFA'
    'oQSU1QT1JUQU5UX0FMRVJUUxACEhEKDVNZU1RFTV9BTEVSVFMQAxIOCgpOT19DT05UQUNUEAQ=');

@$core.Deprecated('Use profileTypeDescriptor instead')
const ProfileType$json = {
  '1': 'ProfileType',
  '2': [
    {'1': 'PERSON', '2': 0},
    {'1': 'INSTITUTION', '2': 1},
    {'1': 'BOT', '2': 2},
  ],
};

/// Descriptor for `ProfileType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List profileTypeDescriptor = $convert.base64Decode(
    'CgtQcm9maWxlVHlwZRIKCgZQRVJTT04QABIPCgtJTlNUSVRVVElPThABEgcKA0JPVBAC');

@$core.Deprecated('Use relationshipTypeDescriptor instead')
const RelationshipType$json = {
  '1': 'RelationshipType',
  '2': [
    {'1': 'MEMBER', '2': 0},
    {'1': 'AFFILIATED', '2': 1},
    {'1': 'BLACK_LISTED', '2': 2},
  ],
};

/// Descriptor for `RelationshipType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List relationshipTypeDescriptor = $convert.base64Decode(
    'ChBSZWxhdGlvbnNoaXBUeXBlEgoKBk1FTUJFUhAAEg4KCkFGRklMSUFURUQQARIQCgxCTEFDS1'
    '9MSVNURUQQAg==');

@$core.Deprecated('Use contactObjectDescriptor instead')
const ContactObject$json = {
  '1': 'ContactObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'type', '3': 2, '4': 1, '5': 14, '6': '.profile.v1.ContactType', '10': 'type'},
    {'1': 'detail', '3': 3, '4': 1, '5': 9, '10': 'detail'},
    {'1': 'verified', '3': 4, '4': 1, '5': 8, '10': 'verified'},
    {'1': 'communication_level', '3': 5, '4': 1, '5': 14, '6': '.profile.v1.CommunicationLevel', '10': 'communicationLevel'},
    {'1': 'state', '3': 6, '4': 1, '5': 14, '6': '.common.v1.STATE', '10': 'state'},
    {'1': 'extra', '3': 7, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
  ],
};

/// Descriptor for `ContactObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contactObjectDescriptor = $convert.base64Decode(
    'Cg1Db250YWN0T2JqZWN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLTlhLXpfLV17Myw0MH'
    '1SAmlkEisKBHR5cGUYAiABKA4yFy5wcm9maWxlLnYxLkNvbnRhY3RUeXBlUgR0eXBlEhYKBmRl'
    'dGFpbBgDIAEoCVIGZGV0YWlsEhoKCHZlcmlmaWVkGAQgASgIUgh2ZXJpZmllZBJPChNjb21tdW'
    '5pY2F0aW9uX2xldmVsGAUgASgOMh4ucHJvZmlsZS52MS5Db21tdW5pY2F0aW9uTGV2ZWxSEmNv'
    'bW11bmljYXRpb25MZXZlbBImCgVzdGF0ZRgGIAEoDjIQLmNvbW1vbi52MS5TVEFURVIFc3RhdG'
    'USLQoFZXh0cmEYByABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0UgVleHRyYQ==');

@$core.Deprecated('Use rosterObjectDescriptor instead')
const RosterObject$json = {
  '1': 'RosterObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'profile_id', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'profileId'},
    {'1': 'contact', '3': 3, '4': 1, '5': 11, '6': '.profile.v1.ContactObject', '10': 'contact'},
    {'1': 'extra', '3': 4, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extra'},
    {'1': 'name', '3': 5, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `RosterObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rosterObjectDescriptor = $convert.base64Decode(
    'CgxSb3N0ZXJPYmplY3QSKwoCaWQYASABKAlCG7pIGHIWEAMYKDIQWzAtOWEtel8tXXszLDQwfV'
    'ICaWQSOgoKcHJvZmlsZV9pZBgCIAEoCUIbukgYchYQAxgoMhBbMC05YS16Xy1dezMsNDB9Uglw'
    'cm9maWxlSWQSMwoHY29udGFjdBgDIAEoCzIZLnByb2ZpbGUudjEuQ29udGFjdE9iamVjdFIHY2'
    '9udGFjdBItCgVleHRyYRgEIAEoCzIXLmdvb2dsZS5wcm90b2J1Zi5TdHJ1Y3RSBWV4dHJhEhIK'
    'BG5hbWUYBSABKAlSBG5hbWU=');

@$core.Deprecated('Use addressObjectDescriptor instead')
const AddressObject$json = {
  '1': 'AddressObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'name'},
    {'1': 'country', '3': 3, '4': 1, '5': 9, '10': 'country'},
    {'1': 'city', '3': 4, '4': 1, '5': 9, '10': 'city'},
    {'1': 'area', '3': 5, '4': 1, '5': 9, '10': 'area'},
    {'1': 'street', '3': 6, '4': 1, '5': 9, '10': 'street'},
    {'1': 'house', '3': 7, '4': 1, '5': 9, '10': 'house'},
    {'1': 'postcode', '3': 8, '4': 1, '5': 9, '10': 'postcode'},
    {'1': 'latitude', '3': 9, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 10, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'extra', '3': 11, '4': 1, '5': 9, '8': {}, '10': 'extra'},
  ],
};

/// Descriptor for `AddressObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addressObjectDescriptor = $convert.base64Decode(
    'Cg1BZGRyZXNzT2JqZWN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLTlhLXpfLV17Myw0MH'
    '1SAmlkEh0KBG5hbWUYAiABKAlCCbpIBnIEEAMYZFIEbmFtZRIYCgdjb3VudHJ5GAMgASgJUgdj'
    'b3VudHJ5EhIKBGNpdHkYBCABKAlSBGNpdHkSEgoEYXJlYRgFIAEoCVIEYXJlYRIWCgZzdHJlZX'
    'QYBiABKAlSBnN0cmVldBIUCgVob3VzZRgHIAEoCVIFaG91c2USGgoIcG9zdGNvZGUYCCABKAlS'
    'CHBvc3Rjb2RlEhoKCGxhdGl0dWRlGAkgASgBUghsYXRpdHVkZRIcCglsb25naXR1ZGUYCiABKA'
    'FSCWxvbmdpdHVkZRIgCgVleHRyYRgLIAEoCUIKukgHcgUQChj0A1IFZXh0cmE=');

@$core.Deprecated('Use profileObjectDescriptor instead')
const ProfileObject$json = {
  '1': 'ProfileObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'type', '3': 2, '4': 1, '5': 14, '6': '.profile.v1.ProfileType', '10': 'type'},
    {'1': 'properties', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'properties'},
    {'1': 'contacts', '3': 4, '4': 3, '5': 11, '6': '.profile.v1.ContactObject', '10': 'contacts'},
    {'1': 'addresses', '3': 5, '4': 3, '5': 11, '6': '.profile.v1.AddressObject', '10': 'addresses'},
    {'1': 'state', '3': 6, '4': 1, '5': 14, '6': '.common.v1.STATE', '10': 'state'},
    {'1': 'accounts', '3': 7, '4': 3, '5': 11, '6': '.profile.v1.ProfileAccount', '10': 'accounts'},
  ],
};

/// Descriptor for `ProfileObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileObjectDescriptor = $convert.base64Decode(
    'Cg1Qcm9maWxlT2JqZWN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLTlhLXpfLV17Myw0MH'
    '1SAmlkEisKBHR5cGUYAiABKA4yFy5wcm9maWxlLnYxLlByb2ZpbGVUeXBlUgR0eXBlEjcKCnBy'
    'b3BlcnRpZXMYAyABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0Ugpwcm9wZXJ0aWVzEjUKCG'
    'NvbnRhY3RzGAQgAygLMhkucHJvZmlsZS52MS5Db250YWN0T2JqZWN0Ughjb250YWN0cxI3Cglh'
    'ZGRyZXNzZXMYBSADKAsyGS5wcm9maWxlLnYxLkFkZHJlc3NPYmplY3RSCWFkZHJlc3NlcxImCg'
    'VzdGF0ZRgGIAEoDjIQLmNvbW1vbi52MS5TVEFURVIFc3RhdGUSNgoIYWNjb3VudHMYByADKAsy'
    'Gi5wcm9maWxlLnYxLlByb2ZpbGVBY2NvdW50UghhY2NvdW50cw==');

@$core.Deprecated('Use profileAccountDescriptor instead')
const ProfileAccount$json = {
  '1': 'ProfileAccount',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    {'1': 'family', '3': 2, '4': 1, '5': 9, '10': 'family'},
    {'1': 'version', '3': 3, '4': 1, '5': 13, '10': 'version'},
    {'1': 'identity_salt_hash', '3': 4, '4': 1, '5': 12, '10': 'identitySaltHash'},
    {'1': 'primary', '3': 5, '4': 1, '5': 8, '10': 'primary'},
  ],
};

/// Descriptor for `ProfileAccount`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List profileAccountDescriptor = $convert.base64Decode(
    'Cg5Qcm9maWxlQWNjb3VudBIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNzEhYKBmZhbWlseRgCIA'
    'EoCVIGZmFtaWx5EhgKB3ZlcnNpb24YAyABKA1SB3ZlcnNpb24SLAoSaWRlbnRpdHlfc2FsdF9o'
    'YXNoGAQgASgMUhBpZGVudGl0eVNhbHRIYXNoEhgKB3ByaW1hcnkYBSABKAhSB3ByaW1hcnk=');

@$core.Deprecated('Use entryItemDescriptor instead')
const EntryItem$json = {
  '1': 'EntryItem',
  '2': [
    {'1': 'object_name', '3': 1, '4': 1, '5': 9, '10': 'objectName'},
    {'1': 'object_id', '3': 2, '4': 1, '5': 9, '10': 'objectId'},
  ],
};

/// Descriptor for `EntryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entryItemDescriptor = $convert.base64Decode(
    'CglFbnRyeUl0ZW0SHwoLb2JqZWN0X25hbWUYASABKAlSCm9iamVjdE5hbWUSGwoJb2JqZWN0X2'
    'lkGAIgASgJUghvYmplY3RJZA==');

@$core.Deprecated('Use relationshipObjectDescriptor instead')
const RelationshipObject$json = {
  '1': 'RelationshipObject',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'type', '3': 2, '4': 1, '5': 14, '6': '.profile.v1.RelationshipType', '10': 'type'},
    {'1': 'properties', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'properties'},
    {'1': 'child_entry', '3': 4, '4': 1, '5': 11, '6': '.profile.v1.EntryItem', '10': 'childEntry'},
    {'1': 'parent_entry', '3': 5, '4': 1, '5': 11, '6': '.profile.v1.EntryItem', '10': 'parentEntry'},
    {'1': 'peer_profile', '3': 6, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'peerProfile'},
  ],
};

/// Descriptor for `RelationshipObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List relationshipObjectDescriptor = $convert.base64Decode(
    'ChJSZWxhdGlvbnNoaXBPYmplY3QSKwoCaWQYASABKAlCG7pIGHIWEAMYKDIQWzAtOWEtel8tXX'
    'szLDQwfVICaWQSMAoEdHlwZRgCIAEoDjIcLnByb2ZpbGUudjEuUmVsYXRpb25zaGlwVHlwZVIE'
    'dHlwZRI3Cgpwcm9wZXJ0aWVzGAMgASgLMhcuZ29vZ2xlLnByb3RvYnVmLlN0cnVjdFIKcHJvcG'
    'VydGllcxI2CgtjaGlsZF9lbnRyeRgEIAEoCzIVLnByb2ZpbGUudjEuRW50cnlJdGVtUgpjaGls'
    'ZEVudHJ5EjgKDHBhcmVudF9lbnRyeRgFIAEoCzIVLnByb2ZpbGUudjEuRW50cnlJdGVtUgtwYX'
    'JlbnRFbnRyeRI8CgxwZWVyX3Byb2ZpbGUYBiABKAsyGS5wcm9maWxlLnYxLlByb2ZpbGVPYmpl'
    'Y3RSC3BlZXJQcm9maWxl');

@$core.Deprecated('Use getByIdRequestDescriptor instead')
const GetByIdRequest$json = {
  '1': 'GetByIdRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
  ],
};

/// Descriptor for `GetByIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByIdRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRCeUlkUmVxdWVzdBIrCgJpZBgBIAEoCUIbukgYchYQAxgoMhBbMC05YS16Xy1dezMsND'
    'B9UgJpZA==');

@$core.Deprecated('Use getByIdResponseDescriptor instead')
const GetByIdResponse$json = {
  '1': 'GetByIdResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `GetByIdResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByIdResponseDescriptor = $convert.base64Decode(
    'Cg9HZXRCeUlkUmVzcG9uc2USLQoEZGF0YRgBIAEoCzIZLnByb2ZpbGUudjEuUHJvZmlsZU9iam'
    'VjdFIEZGF0YQ==');

@$core.Deprecated('Use searchRequestDescriptor instead')
const SearchRequest$json = {
  '1': 'SearchRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'page', '3': 2, '4': 1, '5': 3, '10': 'page'},
    {'1': 'count', '3': 3, '4': 1, '5': 5, '10': 'count'},
    {'1': 'start_date', '3': 4, '4': 1, '5': 9, '10': 'startDate'},
    {'1': 'end_date', '3': 5, '4': 1, '5': 9, '10': 'endDate'},
    {'1': 'properties', '3': 6, '4': 3, '5': 9, '10': 'properties'},
    {'1': 'extras', '3': 7, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extras'},
  ],
};

/// Descriptor for `SearchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchRequestDescriptor = $convert.base64Decode(
    'Cg1TZWFyY2hSZXF1ZXN0EhQKBXF1ZXJ5GAEgASgJUgVxdWVyeRISCgRwYWdlGAIgASgDUgRwYW'
    'dlEhQKBWNvdW50GAMgASgFUgVjb3VudBIdCgpzdGFydF9kYXRlGAQgASgJUglzdGFydERhdGUS'
    'GQoIZW5kX2RhdGUYBSABKAlSB2VuZERhdGUSHgoKcHJvcGVydGllcxgGIAMoCVIKcHJvcGVydG'
    'llcxIvCgZleHRyYXMYByABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0UgZleHRyYXM=');

@$core.Deprecated('Use searchResponseDescriptor instead')
const SearchResponse$json = {
  '1': 'SearchResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `SearchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchResponseDescriptor = $convert.base64Decode(
    'Cg5TZWFyY2hSZXNwb25zZRItCgRkYXRhGAEgAygLMhkucHJvZmlsZS52MS5Qcm9maWxlT2JqZW'
    'N0UgRkYXRh');

@$core.Deprecated('Use mergeRequestDescriptor instead')
const MergeRequest$json = {
  '1': 'MergeRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'mergeid', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'mergeid'},
  ],
};

/// Descriptor for `MergeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mergeRequestDescriptor = $convert.base64Decode(
    'CgxNZXJnZVJlcXVlc3QSKwoCaWQYASABKAlCG7pIGHIWEAMYKDIQWzAtOWEtel8tXXszLDQwfV'
    'ICaWQSNQoHbWVyZ2VpZBgCIAEoCUIbukgYchYQAxgoMhBbMC05YS16Xy1dezMsNDB9UgdtZXJn'
    'ZWlk');

@$core.Deprecated('Use mergeResponseDescriptor instead')
const MergeResponse$json = {
  '1': 'MergeResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `MergeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mergeResponseDescriptor = $convert.base64Decode(
    'Cg1NZXJnZVJlc3BvbnNlEi0KBGRhdGEYASABKAsyGS5wcm9maWxlLnYxLlByb2ZpbGVPYmplY3'
    'RSBGRhdGE=');

@$core.Deprecated('Use createRequestDescriptor instead')
const CreateRequest$json = {
  '1': 'CreateRequest',
  '2': [
    {'1': 'type', '3': 1, '4': 1, '5': 14, '6': '.profile.v1.ProfileType', '8': {}, '10': 'type'},
    {'1': 'contact', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'contact'},
    {'1': 'properties', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'properties'},
  ],
};

/// Descriptor for `CreateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRequestDescriptor = $convert.base64Decode(
    'Cg1DcmVhdGVSZXF1ZXN0EjUKBHR5cGUYASABKA4yFy5wcm9maWxlLnYxLlByb2ZpbGVUeXBlQg'
    'i6SAWCAQIQAVIEdHlwZRIkCgdjb250YWN0GAIgASgJQgq6SAdyBRADGP8BUgdjb250YWN0EjcK'
    'CnByb3BlcnRpZXMYAyABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0Ugpwcm9wZXJ0aWVz');

@$core.Deprecated('Use createResponseDescriptor instead')
const CreateResponse$json = {
  '1': 'CreateResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `CreateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResponseDescriptor = $convert.base64Decode(
    'Cg5DcmVhdGVSZXNwb25zZRItCgRkYXRhGAEgASgLMhkucHJvZmlsZS52MS5Qcm9maWxlT2JqZW'
    'N0UgRkYXRh');

@$core.Deprecated('Use updateRequestDescriptor instead')
const UpdateRequest$json = {
  '1': 'UpdateRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'properties', '3': 2, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'properties'},
    {'1': 'state', '3': 3, '4': 1, '5': 14, '6': '.common.v1.STATE', '10': 'state'},
    {'1': 'scoped', '3': 4, '4': 1, '5': 8, '10': 'scoped'},
  ],
};

/// Descriptor for `UpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRequestDescriptor = $convert.base64Decode(
    'Cg1VcGRhdGVSZXF1ZXN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLTlhLXpfLV17Myw0MH'
    '1SAmlkEjcKCnByb3BlcnRpZXMYAiABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0Ugpwcm9w'
    'ZXJ0aWVzEiYKBXN0YXRlGAMgASgOMhAuY29tbW9uLnYxLlNUQVRFUgVzdGF0ZRIWCgZzY29wZW'
    'QYBCABKAhSBnNjb3BlZA==');

@$core.Deprecated('Use updateResponseDescriptor instead')
const UpdateResponse$json = {
  '1': 'UpdateResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `UpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResponseDescriptor = $convert.base64Decode(
    'Cg5VcGRhdGVSZXNwb25zZRItCgRkYXRhGAEgASgLMhkucHJvZmlsZS52MS5Qcm9maWxlT2JqZW'
    'N0UgRkYXRh');

@$core.Deprecated('Use addContactRequestDescriptor instead')
const AddContactRequest$json = {
  '1': 'AddContactRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'contact', '3': 2, '4': 1, '5': 9, '10': 'contact'},
    {'1': 'extras', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extras'},
  ],
};

/// Descriptor for `AddContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addContactRequestDescriptor = $convert.base64Decode(
    'ChFBZGRDb250YWN0UmVxdWVzdBIuCgJpZBgBIAEoCUIeukgb2AEBchYQAxgoMhBbMC05YS16Xy'
    '1dezMsNDB9UgJpZBIYCgdjb250YWN0GAIgASgJUgdjb250YWN0Ei8KBmV4dHJhcxgDIAEoCzIX'
    'Lmdvb2dsZS5wcm90b2J1Zi5TdHJ1Y3RSBmV4dHJhcw==');

@$core.Deprecated('Use addContactResponseDescriptor instead')
const AddContactResponse$json = {
  '1': 'AddContactResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
    {'1': 'verification_id', '3': 2, '4': 1, '5': 9, '10': 'verificationId'},
  ],
};

/// Descriptor for `AddContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addContactResponseDescriptor = $convert.base64Decode(
    'ChJBZGRDb250YWN0UmVzcG9uc2USLQoEZGF0YRgBIAEoCzIZLnByb2ZpbGUudjEuUHJvZmlsZU'
    '9iamVjdFIEZGF0YRInCg92ZXJpZmljYXRpb25faWQYAiABKAlSDnZlcmlmaWNhdGlvbklk');

@$core.Deprecated('Use createContactRequestDescriptor instead')
const CreateContactRequest$json = {
  '1': 'CreateContactRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'contact', '3': 2, '4': 1, '5': 9, '10': 'contact'},
    {'1': 'extras', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extras'},
  ],
};

/// Descriptor for `CreateContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVDb250YWN0UmVxdWVzdBIuCgJpZBgBIAEoCUIeukgb2AEBchYQAxgoMhBbMC05YS'
    '16Xy1dezMsNDB9UgJpZBIYCgdjb250YWN0GAIgASgJUgdjb250YWN0Ei8KBmV4dHJhcxgDIAEo'
    'CzIXLmdvb2dsZS5wcm90b2J1Zi5TdHJ1Y3RSBmV4dHJhcw==');

@$core.Deprecated('Use createContactResponseDescriptor instead')
const CreateContactResponse$json = {
  '1': 'CreateContactResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ContactObject', '10': 'data'},
  ],
};

/// Descriptor for `CreateContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactResponseDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVDb250YWN0UmVzcG9uc2USLQoEZGF0YRgBIAEoCzIZLnByb2ZpbGUudjEuQ29udG'
    'FjdE9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use getContactsRequestDescriptor instead')
const GetContactsRequest$json = {
  '1': 'GetContactsRequest',
  '2': [
    {'1': 'ids', '3': 1, '4': 3, '5': 9, '8': {}, '10': 'ids'},
  ],
};

/// Descriptor for `GetContactsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContactsRequestDescriptor = $convert.base64Decode(
    'ChJHZXRDb250YWN0c1JlcXVlc3QSNgoDaWRzGAEgAygJQiS6SCGSAR4IARBkIhhyFhADGCgyEF'
    'swLTlhLXpfLV17Myw0MH1SA2lkcw==');

@$core.Deprecated('Use getContactsResponseDescriptor instead')
const GetContactsResponse$json = {
  '1': 'GetContactsResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.ContactObject', '10': 'data'},
    {'1': 'missing_ids', '3': 2, '4': 3, '5': 9, '10': 'missingIds'},
  ],
};

/// Descriptor for `GetContactsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContactsResponseDescriptor = $convert.base64Decode(
    'ChNHZXRDb250YWN0c1Jlc3BvbnNlEi0KBGRhdGEYASADKAsyGS5wcm9maWxlLnYxLkNvbnRhY3'
    'RPYmplY3RSBGRhdGESHwoLbWlzc2luZ19pZHMYAiADKAlSCm1pc3NpbmdJZHM=');

@$core.Deprecated('Use resolveAccountsRequestDescriptor instead')
const ResolveAccountsRequest$json = {
  '1': 'ResolveAccountsRequest',
  '2': [
    {'1': 'addresses', '3': 1, '4': 3, '5': 9, '8': {}, '10': 'addresses'},
  ],
};

/// Descriptor for `ResolveAccountsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resolveAccountsRequestDescriptor = $convert.base64Decode(
    'ChZSZXNvbHZlQWNjb3VudHNSZXF1ZXN0EkUKCWFkZHJlc3NlcxgBIAMoCUInukgkkgEhCAEQ9A'
    'MiGnIYMhZeMFt4WF1bMC05YS1mQS1GXXs0MH0kUglhZGRyZXNzZXM=');

@$core.Deprecated('Use accountOwnerDescriptor instead')
const AccountOwner$json = {
  '1': 'AccountOwner',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    {'1': 'profile_id', '3': 2, '4': 1, '5': 9, '10': 'profileId'},
  ],
};

/// Descriptor for `AccountOwner`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountOwnerDescriptor = $convert.base64Decode(
    'CgxBY2NvdW50T3duZXISGAoHYWRkcmVzcxgBIAEoCVIHYWRkcmVzcxIdCgpwcm9maWxlX2lkGA'
    'IgASgJUglwcm9maWxlSWQ=');

@$core.Deprecated('Use resolveAccountsResponseDescriptor instead')
const ResolveAccountsResponse$json = {
  '1': 'ResolveAccountsResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.AccountOwner', '10': 'data'},
  ],
};

/// Descriptor for `ResolveAccountsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resolveAccountsResponseDescriptor = $convert.base64Decode(
    'ChdSZXNvbHZlQWNjb3VudHNSZXNwb25zZRIsCgRkYXRhGAEgAygLMhgucHJvZmlsZS52MS5BY2'
    'NvdW50T3duZXJSBGRhdGE=');

@$core.Deprecated('Use createContactVerificationRequestDescriptor instead')
const CreateContactVerificationRequest$json = {
  '1': 'CreateContactVerificationRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'contact_id', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'contactId'},
    {'1': 'code', '3': 3, '4': 1, '5': 9, '10': 'code'},
    {'1': 'duration_to_expire', '3': 4, '4': 1, '5': 9, '10': 'durationToExpire'},
  ],
};

/// Descriptor for `CreateContactVerificationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactVerificationRequestDescriptor = $convert.base64Decode(
    'CiBDcmVhdGVDb250YWN0VmVyaWZpY2F0aW9uUmVxdWVzdBIrCgJpZBgBIAEoCUIbukgYchYQAx'
    'goMhBbMC05YS16Xy1dezMsNDB9UgJpZBI6Cgpjb250YWN0X2lkGAIgASgJQhu6SBhyFhADGCgy'
    'EFswLTlhLXpfLV17Myw0MH1SCWNvbnRhY3RJZBISCgRjb2RlGAMgASgJUgRjb2RlEiwKEmR1cm'
    'F0aW9uX3RvX2V4cGlyZRgEIAEoCVIQZHVyYXRpb25Ub0V4cGlyZQ==');

@$core.Deprecated('Use createContactVerificationResponseDescriptor instead')
const CreateContactVerificationResponse$json = {
  '1': 'CreateContactVerificationResponse',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'success', '3': 2, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `CreateContactVerificationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactVerificationResponseDescriptor = $convert.base64Decode(
    'CiFDcmVhdGVDb250YWN0VmVyaWZpY2F0aW9uUmVzcG9uc2USKwoCaWQYASABKAlCG7pIGHIWEA'
    'MYKDIQWzAtOWEtel8tXXszLDQwfVICaWQSGAoHc3VjY2VzcxgCIAEoCFIHc3VjY2Vzcw==');

@$core.Deprecated('Use checkVerificationRequestDescriptor instead')
const CheckVerificationRequest$json = {
  '1': 'CheckVerificationRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'code', '3': 2, '4': 1, '5': 9, '10': 'code'},
  ],
};

/// Descriptor for `CheckVerificationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkVerificationRequestDescriptor = $convert.base64Decode(
    'ChhDaGVja1ZlcmlmaWNhdGlvblJlcXVlc3QSKwoCaWQYASABKAlCG7pIGHIWEAMYKDIQWzAtOW'
    'Etel8tXXszLDQwfVICaWQSEgoEY29kZRgCIAEoCVIEY29kZQ==');

@$core.Deprecated('Use checkVerificationResponseDescriptor instead')
const CheckVerificationResponse$json = {
  '1': 'CheckVerificationResponse',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'check_attempts', '3': 2, '4': 1, '5': 5, '10': 'checkAttempts'},
    {'1': 'success', '3': 3, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `CheckVerificationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkVerificationResponseDescriptor = $convert.base64Decode(
    'ChlDaGVja1ZlcmlmaWNhdGlvblJlc3BvbnNlEg4KAmlkGAEgASgJUgJpZBIlCg5jaGVja19hdH'
    'RlbXB0cxgCIAEoBVINY2hlY2tBdHRlbXB0cxIYCgdzdWNjZXNzGAMgASgIUgdzdWNjZXNz');

@$core.Deprecated('Use removeContactRequestDescriptor instead')
const RemoveContactRequest$json = {
  '1': 'RemoveContactRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
  ],
};

/// Descriptor for `RemoveContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeContactRequestDescriptor = $convert.base64Decode(
    'ChRSZW1vdmVDb250YWN0UmVxdWVzdBIrCgJpZBgBIAEoCUIbukgYchYQAxgoMhBbMC05YS16Xy'
    '1dezMsNDB9UgJpZA==');

@$core.Deprecated('Use removeContactResponseDescriptor instead')
const RemoveContactResponse$json = {
  '1': 'RemoveContactResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `RemoveContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeContactResponseDescriptor = $convert.base64Decode(
    'ChVSZW1vdmVDb250YWN0UmVzcG9uc2USLQoEZGF0YRgBIAEoCzIZLnByb2ZpbGUudjEuUHJvZm'
    'lsZU9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use searchRosterRequestDescriptor instead')
const SearchRosterRequest$json = {
  '1': 'SearchRosterRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'page', '3': 2, '4': 1, '5': 3, '10': 'page'},
    {'1': 'count', '3': 3, '4': 1, '5': 5, '10': 'count'},
    {'1': 'start_date', '3': 4, '4': 1, '5': 9, '10': 'startDate'},
    {'1': 'end_date', '3': 5, '4': 1, '5': 9, '10': 'endDate'},
    {'1': 'properties', '3': 6, '4': 3, '5': 9, '10': 'properties'},
    {'1': 'extras', '3': 7, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extras'},
    {'1': 'profile_id', '3': 8, '4': 1, '5': 9, '8': {}, '10': 'profileId'},
    {'1': 'name', '3': 9, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `SearchRosterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchRosterRequestDescriptor = $convert.base64Decode(
    'ChNTZWFyY2hSb3N0ZXJSZXF1ZXN0EhQKBXF1ZXJ5GAEgASgJUgVxdWVyeRISCgRwYWdlGAIgAS'
    'gDUgRwYWdlEhQKBWNvdW50GAMgASgFUgVjb3VudBIdCgpzdGFydF9kYXRlGAQgASgJUglzdGFy'
    'dERhdGUSGQoIZW5kX2RhdGUYBSABKAlSB2VuZERhdGUSHgoKcHJvcGVydGllcxgGIAMoCVIKcH'
    'JvcGVydGllcxIvCgZleHRyYXMYByABKAsyFy5nb29nbGUucHJvdG9idWYuU3RydWN0UgZleHRy'
    'YXMSPwoKcHJvZmlsZV9pZBgIIAEoCUIgukgd2AEBchgQAxj6ATIRWzAtOWEtel8tXXszLDI1MH'
    '1SCXByb2ZpbGVJZBISCgRuYW1lGAkgASgJUgRuYW1l');

@$core.Deprecated('Use searchRosterResponseDescriptor instead')
const SearchRosterResponse$json = {
  '1': 'SearchRosterResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.RosterObject', '10': 'data'},
  ],
};

/// Descriptor for `SearchRosterResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchRosterResponseDescriptor = $convert.base64Decode(
    'ChRTZWFyY2hSb3N0ZXJSZXNwb25zZRIsCgRkYXRhGAEgAygLMhgucHJvZmlsZS52MS5Sb3N0ZX'
    'JPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use rawContactDescriptor instead')
const RawContact$json = {
  '1': 'RawContact',
  '2': [
    {'1': 'contact', '3': 1, '4': 1, '5': 9, '10': 'contact'},
    {'1': 'extras', '3': 2, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'extras'},
  ],
};

/// Descriptor for `RawContact`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rawContactDescriptor = $convert.base64Decode(
    'CgpSYXdDb250YWN0EhgKB2NvbnRhY3QYASABKAlSB2NvbnRhY3QSLwoGZXh0cmFzGAIgASgLMh'
    'cuZ29vZ2xlLnByb3RvYnVmLlN0cnVjdFIGZXh0cmFz');

@$core.Deprecated('Use addRosterRequestDescriptor instead')
const AddRosterRequest$json = {
  '1': 'AddRosterRequest',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.RawContact', '10': 'data'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `AddRosterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addRosterRequestDescriptor = $convert.base64Decode(
    'ChBBZGRSb3N0ZXJSZXF1ZXN0EioKBGRhdGEYASADKAsyFi5wcm9maWxlLnYxLlJhd0NvbnRhY3'
    'RSBGRhdGESEgoEbmFtZRgCIAEoCVIEbmFtZQ==');

@$core.Deprecated('Use addRosterResponseDescriptor instead')
const AddRosterResponse$json = {
  '1': 'AddRosterResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.RosterObject', '10': 'data'},
  ],
};

/// Descriptor for `AddRosterResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addRosterResponseDescriptor = $convert.base64Decode(
    'ChFBZGRSb3N0ZXJSZXNwb25zZRIsCgRkYXRhGAEgAygLMhgucHJvZmlsZS52MS5Sb3N0ZXJPYm'
    'plY3RSBGRhdGE=');

@$core.Deprecated('Use removeRosterRequestDescriptor instead')
const RemoveRosterRequest$json = {
  '1': 'RemoveRosterRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `RemoveRosterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeRosterRequestDescriptor = $convert.base64Decode(
    'ChNSZW1vdmVSb3N0ZXJSZXF1ZXN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLTlhLXpfLV'
    '17Myw0MH1SAmlkEhIKBG5hbWUYAiABKAlSBG5hbWU=');

@$core.Deprecated('Use removeRosterResponseDescriptor instead')
const RemoveRosterResponse$json = {
  '1': 'RemoveRosterResponse',
  '2': [
    {'1': 'roster', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.RosterObject', '10': 'roster'},
  ],
};

/// Descriptor for `RemoveRosterResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeRosterResponseDescriptor = $convert.base64Decode(
    'ChRSZW1vdmVSb3N0ZXJSZXNwb25zZRIwCgZyb3N0ZXIYASABKAsyGC5wcm9maWxlLnYxLlJvc3'
    'Rlck9iamVjdFIGcm9zdGVy');

@$core.Deprecated('Use addAddressRequestDescriptor instead')
const AddAddressRequest$json = {
  '1': 'AddAddressRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'address', '3': 2, '4': 1, '5': 11, '6': '.profile.v1.AddressObject', '10': 'address'},
  ],
};

/// Descriptor for `AddAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addAddressRequestDescriptor = $convert.base64Decode(
    'ChFBZGRBZGRyZXNzUmVxdWVzdBIrCgJpZBgBIAEoCUIbukgYchYQAxgoMhBbMC05YS16Xy1dez'
    'MsNDB9UgJpZBIzCgdhZGRyZXNzGAIgASgLMhkucHJvZmlsZS52MS5BZGRyZXNzT2JqZWN0Ugdh'
    'ZGRyZXNz');

@$core.Deprecated('Use addAddressResponseDescriptor instead')
const AddAddressResponse$json = {
  '1': 'AddAddressResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `AddAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addAddressResponseDescriptor = $convert.base64Decode(
    'ChJBZGRBZGRyZXNzUmVzcG9uc2USLQoEZGF0YRgBIAEoCzIZLnByb2ZpbGUudjEuUHJvZmlsZU'
    '9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use getByContactRequestDescriptor instead')
const GetByContactRequest$json = {
  '1': 'GetByContactRequest',
  '2': [
    {'1': 'contact', '3': 1, '4': 1, '5': 9, '10': 'contact'},
  ],
};

/// Descriptor for `GetByContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByContactRequestDescriptor = $convert.base64Decode(
    'ChNHZXRCeUNvbnRhY3RSZXF1ZXN0EhgKB2NvbnRhY3QYASABKAlSB2NvbnRhY3Q=');

@$core.Deprecated('Use getByContactResponseDescriptor instead')
const GetByContactResponse$json = {
  '1': 'GetByContactResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `GetByContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByContactResponseDescriptor = $convert.base64Decode(
    'ChRHZXRCeUNvbnRhY3RSZXNwb25zZRItCgRkYXRhGAEgASgLMhkucHJvZmlsZS52MS5Qcm9maW'
    'xlT2JqZWN0UgRkYXRh');

@$core.Deprecated('Use getByIDAndPartitionRequestDescriptor instead')
const GetByIDAndPartitionRequest$json = {
  '1': 'GetByIDAndPartitionRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'partition_id', '3': 2, '4': 1, '5': 9, '10': 'partitionId'},
  ],
};

/// Descriptor for `GetByIDAndPartitionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByIDAndPartitionRequestDescriptor = $convert.base64Decode(
    'ChpHZXRCeUlEQW5kUGFydGl0aW9uUmVxdWVzdBIOCgJpZBgBIAEoCVICaWQSIQoMcGFydGl0aW'
    '9uX2lkGAIgASgJUgtwYXJ0aXRpb25JZA==');

@$core.Deprecated('Use getByIDAndPartitionResponseDescriptor instead')
const GetByIDAndPartitionResponse$json = {
  '1': 'GetByIDAndPartitionResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.ProfileObject', '10': 'data'},
  ],
};

/// Descriptor for `GetByIDAndPartitionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByIDAndPartitionResponseDescriptor = $convert.base64Decode(
    'ChtHZXRCeUlEQW5kUGFydGl0aW9uUmVzcG9uc2USLQoEZGF0YRgBIAEoCzIZLnByb2ZpbGUudj'
    'EuUHJvZmlsZU9iamVjdFIEZGF0YQ==');

@$core.Deprecated('Use propertyHistoryRequestDescriptor instead')
const PropertyHistoryRequest$json = {
  '1': 'PropertyHistoryRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'key', '3': 2, '4': 1, '5': 9, '10': 'key'},
  ],
};

/// Descriptor for `PropertyHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List propertyHistoryRequestDescriptor = $convert.base64Decode(
    'ChZQcm9wZXJ0eUhpc3RvcnlSZXF1ZXN0Eg4KAmlkGAEgASgJUgJpZBIQCgNrZXkYAiABKAlSA2'
    'tleQ==');

@$core.Deprecated('Use propertyEntryObjectDescriptor instead')
const PropertyEntryObject$json = {
  '1': 'PropertyEntryObject',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
    {'1': 'tenant_id', '3': 3, '4': 1, '5': 9, '10': 'tenantId'},
    {'1': 'created_by', '3': 4, '4': 1, '5': 9, '10': 'createdBy'},
    {'1': 'created_at', '3': 5, '4': 1, '5': 11, '6': '.google.protobuf.Timestamp', '10': 'createdAt'},
    {'1': 'scoped', '3': 6, '4': 1, '5': 8, '10': 'scoped'},
  ],
};

/// Descriptor for `PropertyEntryObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List propertyEntryObjectDescriptor = $convert.base64Decode(
    'ChNQcm9wZXJ0eUVudHJ5T2JqZWN0EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZRIbCgl0ZW5hbnRfaWQYAyABKAlSCHRlbmFudElkEh0KCmNyZWF0ZWRfYnkYBCABKAlS'
    'CWNyZWF0ZWRCeRI5CgpjcmVhdGVkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdG'
    'FtcFIJY3JlYXRlZEF0EhYKBnNjb3BlZBgGIAEoCFIGc2NvcGVk');

@$core.Deprecated('Use propertyHistoryResponseDescriptor instead')
const PropertyHistoryResponse$json = {
  '1': 'PropertyHistoryResponse',
  '2': [
    {'1': 'entries', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.PropertyEntryObject', '10': 'entries'},
  ],
};

/// Descriptor for `PropertyHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List propertyHistoryResponseDescriptor = $convert.base64Decode(
    'ChdQcm9wZXJ0eUhpc3RvcnlSZXNwb25zZRI5CgdlbnRyaWVzGAEgAygLMh8ucHJvZmlsZS52MS'
    '5Qcm9wZXJ0eUVudHJ5T2JqZWN0UgdlbnRyaWVz');

@$core.Deprecated('Use listRelationshipRequestDescriptor instead')
const ListRelationshipRequest$json = {
  '1': 'ListRelationshipRequest',
  '2': [
    {'1': 'peer_name', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'peerName'},
    {'1': 'peer_id', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'peerId'},
    {'1': 'last_relationship_id', '3': 3, '4': 1, '5': 9, '8': {}, '10': 'lastRelationshipId'},
    {'1': 'related_children_id', '3': 4, '4': 3, '5': 9, '10': 'relatedChildrenId'},
    {'1': 'count', '3': 5, '4': 1, '5': 5, '10': 'count'},
    {'1': 'invert_relation', '3': 6, '4': 1, '5': 8, '10': 'invertRelation'},
  ],
};

/// Descriptor for `ListRelationshipRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRelationshipRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0UmVsYXRpb25zaGlwUmVxdWVzdBI/CglwZWVyX25hbWUYASABKAlCIrpIH3IdEAMYKF'
    'IHQ29udGFjdFIHUHJvZmlsZVIFR3JvdXBSCHBlZXJOYW1lEjQKB3BlZXJfaWQYAiABKAlCG7pI'
    'GHIWEAMYKDIQWzAtOWEtel8tXXszLDQwfVIGcGVlcklkElAKFGxhc3RfcmVsYXRpb25zaGlwX2'
    'lkGAMgASgJQh66SBvYAQFyFhADGCgyEFswLTlhLXpfLV17Myw0MH1SEmxhc3RSZWxhdGlvbnNo'
    'aXBJZBIuChNyZWxhdGVkX2NoaWxkcmVuX2lkGAQgAygJUhFyZWxhdGVkQ2hpbGRyZW5JZBIUCg'
    'Vjb3VudBgFIAEoBVIFY291bnQSJwoPaW52ZXJ0X3JlbGF0aW9uGAYgASgIUg5pbnZlcnRSZWxh'
    'dGlvbg==');

@$core.Deprecated('Use listRelationshipResponseDescriptor instead')
const ListRelationshipResponse$json = {
  '1': 'ListRelationshipResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.profile.v1.RelationshipObject', '10': 'data'},
  ],
};

/// Descriptor for `ListRelationshipResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRelationshipResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0UmVsYXRpb25zaGlwUmVzcG9uc2USMgoEZGF0YRgBIAMoCzIeLnByb2ZpbGUudjEuUm'
    'VsYXRpb25zaGlwT2JqZWN0UgRkYXRh');

@$core.Deprecated('Use addRelationshipRequestDescriptor instead')
const AddRelationshipRequest$json = {
  '1': 'AddRelationshipRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'parent', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'parent'},
    {'1': 'parent_id', '3': 3, '4': 1, '5': 9, '8': {}, '10': 'parentId'},
    {'1': 'child', '3': 4, '4': 1, '5': 9, '8': {}, '10': 'child'},
    {'1': 'child_id', '3': 5, '4': 1, '5': 9, '8': {}, '10': 'childId'},
    {'1': 'type', '3': 6, '4': 1, '5': 14, '6': '.profile.v1.RelationshipType', '10': 'type'},
    {'1': 'properties', '3': 7, '4': 1, '5': 11, '6': '.google.protobuf.Struct', '10': 'properties'},
  ],
};

/// Descriptor for `AddRelationshipRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addRelationshipRequestDescriptor = $convert.base64Decode(
    'ChZBZGRSZWxhdGlvbnNoaXBSZXF1ZXN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLTlhLX'
    'pfLV17Myw0MH1SAmlkEjoKBnBhcmVudBgCIAEoCUIiukgfch0QAxgoUgdDb250YWN0UgdQcm9m'
    'aWxlUgVHcm91cFIGcGFyZW50EjgKCXBhcmVudF9pZBgDIAEoCUIbukgYchYQAxgoMhBbMC05YS'
    '16Xy1dezMsNDB9UghwYXJlbnRJZBI4CgVjaGlsZBgEIAEoCUIiukgfch0QAxgoUgdDb250YWN0'
    'UgdQcm9maWxlUgVHcm91cFIFY2hpbGQSNgoIY2hpbGRfaWQYBSABKAlCG7pIGHIWEAMYKDIQWz'
    'AtOWEtel8tXXszLDQwfVIHY2hpbGRJZBIwCgR0eXBlGAYgASgOMhwucHJvZmlsZS52MS5SZWxh'
    'dGlvbnNoaXBUeXBlUgR0eXBlEjcKCnByb3BlcnRpZXMYByABKAsyFy5nb29nbGUucHJvdG9idW'
    'YuU3RydWN0Ugpwcm9wZXJ0aWVz');

@$core.Deprecated('Use addRelationshipResponseDescriptor instead')
const AddRelationshipResponse$json = {
  '1': 'AddRelationshipResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.RelationshipObject', '10': 'data'},
  ],
};

/// Descriptor for `AddRelationshipResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addRelationshipResponseDescriptor = $convert.base64Decode(
    'ChdBZGRSZWxhdGlvbnNoaXBSZXNwb25zZRIyCgRkYXRhGAEgASgLMh4ucHJvZmlsZS52MS5SZW'
    'xhdGlvbnNoaXBPYmplY3RSBGRhdGE=');

@$core.Deprecated('Use deleteRelationshipRequestDescriptor instead')
const DeleteRelationshipRequest$json = {
  '1': 'DeleteRelationshipRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '8': {}, '10': 'id'},
    {'1': 'parent_id', '3': 2, '4': 1, '5': 9, '8': {}, '10': 'parentId'},
  ],
};

/// Descriptor for `DeleteRelationshipRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRelationshipRequestDescriptor = $convert.base64Decode(
    'ChlEZWxldGVSZWxhdGlvbnNoaXBSZXF1ZXN0EisKAmlkGAEgASgJQhu6SBhyFhADGCgyEFswLT'
    'lhLXpfLV17Myw0MH1SAmlkEjsKCXBhcmVudF9pZBgCIAEoCUIeukgb2AEBchYQAxgoMhBbMC05'
    'YS16Xy1dezMsNDB9UghwYXJlbnRJZA==');

@$core.Deprecated('Use deleteRelationshipResponseDescriptor instead')
const DeleteRelationshipResponse$json = {
  '1': 'DeleteRelationshipResponse',
  '2': [
    {'1': 'data', '3': 1, '4': 1, '5': 11, '6': '.profile.v1.RelationshipObject', '10': 'data'},
  ],
};

/// Descriptor for `DeleteRelationshipResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRelationshipResponseDescriptor = $convert.base64Decode(
    'ChpEZWxldGVSZWxhdGlvbnNoaXBSZXNwb25zZRIyCgRkYXRhGAEgASgLMh4ucHJvZmlsZS52MS'
    '5SZWxhdGlvbnNoaXBPYmplY3RSBGRhdGE=');

const $core.Map<$core.String, $core.dynamic> ProfileServiceBase$json = {
  '1': 'ProfileService',
  '2': [
    {
      '1': 'GetById',
      '2': '.profile.v1.GetByIdRequest',
      '3': '.profile.v1.GetByIdResponse',
      '4': {'34': 1},
    },
    {
      '1': 'GetByContact',
      '2': '.profile.v1.GetByContactRequest',
      '3': '.profile.v1.GetByContactResponse',
      '4': {'34': 1},
    },
    {
      '1': 'Search',
      '2': '.profile.v1.SearchRequest',
      '3': '.profile.v1.SearchResponse',
      '4': {'34': 1},
      '6': true,
    },
    {'1': 'Merge', '2': '.profile.v1.MergeRequest', '3': '.profile.v1.MergeResponse', '4': {}},
    {'1': 'Create', '2': '.profile.v1.CreateRequest', '3': '.profile.v1.CreateResponse', '4': {}},
    {'1': 'Update', '2': '.profile.v1.UpdateRequest', '3': '.profile.v1.UpdateResponse', '4': {}},
    {'1': 'AddContact', '2': '.profile.v1.AddContactRequest', '3': '.profile.v1.AddContactResponse', '4': {}},
    {'1': 'CreateContact', '2': '.profile.v1.CreateContactRequest', '3': '.profile.v1.CreateContactResponse', '4': {}},
    {
      '1': 'GetContacts',
      '2': '.profile.v1.GetContactsRequest',
      '3': '.profile.v1.GetContactsResponse',
      '4': {'34': 1},
    },
    {'1': 'CreateContactVerification', '2': '.profile.v1.CreateContactVerificationRequest', '3': '.profile.v1.CreateContactVerificationResponse', '4': {}},
    {'1': 'CheckVerification', '2': '.profile.v1.CheckVerificationRequest', '3': '.profile.v1.CheckVerificationResponse', '4': {}},
    {'1': 'RemoveContact', '2': '.profile.v1.RemoveContactRequest', '3': '.profile.v1.RemoveContactResponse', '4': {}},
    {
      '1': 'SearchRoster',
      '2': '.profile.v1.SearchRosterRequest',
      '3': '.profile.v1.SearchRosterResponse',
      '4': {'34': 1},
      '6': true,
    },
    {'1': 'AddRoster', '2': '.profile.v1.AddRosterRequest', '3': '.profile.v1.AddRosterResponse', '4': {}},
    {'1': 'RemoveRoster', '2': '.profile.v1.RemoveRosterRequest', '3': '.profile.v1.RemoveRosterResponse', '4': {}},
    {'1': 'AddAddress', '2': '.profile.v1.AddAddressRequest', '3': '.profile.v1.AddAddressResponse', '4': {}},
    {'1': 'AddRelationship', '2': '.profile.v1.AddRelationshipRequest', '3': '.profile.v1.AddRelationshipResponse', '4': {}},
    {'1': 'DeleteRelationship', '2': '.profile.v1.DeleteRelationshipRequest', '3': '.profile.v1.DeleteRelationshipResponse', '4': {}},
    {
      '1': 'ListRelationship',
      '2': '.profile.v1.ListRelationshipRequest',
      '3': '.profile.v1.ListRelationshipResponse',
      '4': {'34': 1},
      '6': true,
    },
    {
      '1': 'GetByIDAndPartition',
      '2': '.profile.v1.GetByIDAndPartitionRequest',
      '3': '.profile.v1.GetByIDAndPartitionResponse',
      '4': {'34': 1},
    },
    {
      '1': 'PropertyHistory',
      '2': '.profile.v1.PropertyHistoryRequest',
      '3': '.profile.v1.PropertyHistoryResponse',
      '4': {'34': 1},
    },
    {
      '1': 'ResolveAccounts',
      '2': '.profile.v1.ResolveAccountsRequest',
      '3': '.profile.v1.ResolveAccountsResponse',
      '4': {'34': 1},
    },
  ],
  '3': {},
};

@$core.Deprecated('Use profileServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> ProfileServiceBase$messageJson = {
  '.profile.v1.GetByIdRequest': GetByIdRequest$json,
  '.profile.v1.GetByIdResponse': GetByIdResponse$json,
  '.profile.v1.ProfileObject': ProfileObject$json,
  '.google.protobuf.Struct': $6.Struct$json,
  '.google.protobuf.Struct.FieldsEntry': $6.Struct_FieldsEntry$json,
  '.google.protobuf.Value': $6.Value$json,
  '.google.protobuf.ListValue': $6.ListValue$json,
  '.profile.v1.ContactObject': ContactObject$json,
  '.profile.v1.AddressObject': AddressObject$json,
  '.profile.v1.ProfileAccount': ProfileAccount$json,
  '.profile.v1.GetByContactRequest': GetByContactRequest$json,
  '.profile.v1.GetByContactResponse': GetByContactResponse$json,
  '.profile.v1.SearchRequest': SearchRequest$json,
  '.profile.v1.SearchResponse': SearchResponse$json,
  '.profile.v1.MergeRequest': MergeRequest$json,
  '.profile.v1.MergeResponse': MergeResponse$json,
  '.profile.v1.CreateRequest': CreateRequest$json,
  '.profile.v1.CreateResponse': CreateResponse$json,
  '.profile.v1.UpdateRequest': UpdateRequest$json,
  '.profile.v1.UpdateResponse': UpdateResponse$json,
  '.profile.v1.AddContactRequest': AddContactRequest$json,
  '.profile.v1.AddContactResponse': AddContactResponse$json,
  '.profile.v1.CreateContactRequest': CreateContactRequest$json,
  '.profile.v1.CreateContactResponse': CreateContactResponse$json,
  '.profile.v1.GetContactsRequest': GetContactsRequest$json,
  '.profile.v1.GetContactsResponse': GetContactsResponse$json,
  '.profile.v1.CreateContactVerificationRequest': CreateContactVerificationRequest$json,
  '.profile.v1.CreateContactVerificationResponse': CreateContactVerificationResponse$json,
  '.profile.v1.CheckVerificationRequest': CheckVerificationRequest$json,
  '.profile.v1.CheckVerificationResponse': CheckVerificationResponse$json,
  '.profile.v1.RemoveContactRequest': RemoveContactRequest$json,
  '.profile.v1.RemoveContactResponse': RemoveContactResponse$json,
  '.profile.v1.SearchRosterRequest': SearchRosterRequest$json,
  '.profile.v1.SearchRosterResponse': SearchRosterResponse$json,
  '.profile.v1.RosterObject': RosterObject$json,
  '.profile.v1.AddRosterRequest': AddRosterRequest$json,
  '.profile.v1.RawContact': RawContact$json,
  '.profile.v1.AddRosterResponse': AddRosterResponse$json,
  '.profile.v1.RemoveRosterRequest': RemoveRosterRequest$json,
  '.profile.v1.RemoveRosterResponse': RemoveRosterResponse$json,
  '.profile.v1.AddAddressRequest': AddAddressRequest$json,
  '.profile.v1.AddAddressResponse': AddAddressResponse$json,
  '.profile.v1.AddRelationshipRequest': AddRelationshipRequest$json,
  '.profile.v1.AddRelationshipResponse': AddRelationshipResponse$json,
  '.profile.v1.RelationshipObject': RelationshipObject$json,
  '.profile.v1.EntryItem': EntryItem$json,
  '.profile.v1.DeleteRelationshipRequest': DeleteRelationshipRequest$json,
  '.profile.v1.DeleteRelationshipResponse': DeleteRelationshipResponse$json,
  '.profile.v1.ListRelationshipRequest': ListRelationshipRequest$json,
  '.profile.v1.ListRelationshipResponse': ListRelationshipResponse$json,
  '.profile.v1.GetByIDAndPartitionRequest': GetByIDAndPartitionRequest$json,
  '.profile.v1.GetByIDAndPartitionResponse': GetByIDAndPartitionResponse$json,
  '.profile.v1.PropertyHistoryRequest': PropertyHistoryRequest$json,
  '.profile.v1.PropertyHistoryResponse': PropertyHistoryResponse$json,
  '.profile.v1.PropertyEntryObject': PropertyEntryObject$json,
  '.google.protobuf.Timestamp': $2.Timestamp$json,
  '.profile.v1.ResolveAccountsRequest': ResolveAccountsRequest$json,
  '.profile.v1.ResolveAccountsResponse': ResolveAccountsResponse$json,
  '.profile.v1.AccountOwner': AccountOwner$json,
};

/// Descriptor for `ProfileService`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List profileServiceDescriptor = $convert.base64Decode(
    'Cg5Qcm9maWxlU2VydmljZRLxAQoHR2V0QnlJZBIaLnByb2ZpbGUudjEuR2V0QnlJZFJlcXVlc3'
    'QaGy5wcm9maWxlLnYxLkdldEJ5SWRSZXNwb25zZSKsAZACAbpHkwEKCFByb2ZpbGVzEhFHZXQg'
    'cHJvZmlsZSBieSBJRBpkUmV0cmlldmVzIGEgY29tcGxldGUgcHJvZmlsZSBieSBpdHMgdW5pcX'
    'VlIGlkZW50aWZpZXIgaW5jbHVkaW5nIGNvbnRhY3RzLCBhZGRyZXNzZXMsIGFuZCBwcm9wZXJ0'
    'aWVzLioOZ2V0UHJvZmlsZUJ5SWSCtRgOCgxwcm9maWxlX3ZpZXcS9QEKDEdldEJ5Q29udGFjdB'
    'IfLnByb2ZpbGUudjEuR2V0QnlDb250YWN0UmVxdWVzdBogLnByb2ZpbGUudjEuR2V0QnlDb250'
    'YWN0UmVzcG9uc2UioQGQAgG6R4gBCghQcm9maWxlcxIWR2V0IHByb2ZpbGUgYnkgY29udGFjdB'
    'pPUmV0cmlldmVzIGEgcHJvZmlsZSBhc3NvY2lhdGVkIHdpdGggYSBzcGVjaWZpYyBjb250YWN0'
    'IChlbWFpbCBvciBwaG9uZSBudW1iZXIpLioTZ2V0UHJvZmlsZUJ5Q29udGFjdIK1GA4KDHByb2'
    'ZpbGVfdmlldxKfAgoGU2VhcmNoEhkucHJvZmlsZS52MS5TZWFyY2hSZXF1ZXN0GhoucHJvZmls'
    'ZS52MS5TZWFyY2hSZXNwb25zZSLbAZACAbpHwgEKCFByb2ZpbGVzEg9TZWFyY2ggcHJvZmlsZX'
    'MalAFTZWFyY2hlcyBmb3IgcHJvZmlsZXMgbWF0Y2hpbmcgc3BlY2lmaWVkIGNyaXRlcmlhIGlu'
    'Y2x1ZGluZyBuYW1lLCBjb250YWN0LCBkYXRlIHJhbmdlLCBhbmQgY3VzdG9tIHByb3BlcnRpZX'
    'MuIFJldHVybnMgYSBzdHJlYW0gb2YgbWF0Y2hpbmcgcHJvZmlsZXMuKg5zZWFyY2hQcm9maWxl'
    'c4K1GA4KDHByb2ZpbGVfdmlldzABEvQBCgVNZXJnZRIYLnByb2ZpbGUudjEuTWVyZ2VSZXF1ZX'
    'N0GhkucHJvZmlsZS52MS5NZXJnZVJlc3BvbnNlIrUBukeeAQoIUHJvZmlsZXMSDk1lcmdlIHBy'
    'b2ZpbGVzGnNNZXJnZXMgdHdvIHByb2ZpbGVzIGJ5IGNvbWJpbmluZyB0aGVpciBkYXRhLiBUaG'
    'UgbWVyZ2Ugc291cmNlIHByb2ZpbGUgZGF0YSBpcyBpbmNvcnBvcmF0ZWQgaW50byB0aGUgdGFy'
    'Z2V0IHByb2ZpbGUuKg1tZXJnZVByb2ZpbGVzgrUYDwoNcHJvZmlsZV9tZXJnZRLuAQoGQ3JlYX'
    'RlEhkucHJvZmlsZS52MS5DcmVhdGVSZXF1ZXN0GhoucHJvZmlsZS52MS5DcmVhdGVSZXNwb25z'
    'ZSKsAbpHlAEKCFByb2ZpbGVzEg5DcmVhdGUgcHJvZmlsZRppQ3JlYXRlcyBhIG5ldyBwcm9maW'
    'xlIHdpdGggdGhlIHNwZWNpZmllZCB0eXBlIChwZXJzb24sIGluc3RpdHV0aW9uLCBib3QpIGFu'
    'ZCBpbml0aWFsIGNvbnRhY3QgaW5mb3JtYXRpb24uKg1jcmVhdGVQcm9maWxlgrUYEAoOcHJvZm'
    'lsZV9jcmVhdGUS7gEKBlVwZGF0ZRIZLnByb2ZpbGUudjEuVXBkYXRlUmVxdWVzdBoaLnByb2Zp'
    'bGUudjEuVXBkYXRlUmVzcG9uc2UirAG6R5QBCghQcm9maWxlcxIOVXBkYXRlIHByb2ZpbGUaaV'
    'VwZGF0ZXMgYW4gZXhpc3RpbmcgcHJvZmlsZSdzIHByb3BlcnRpZXMgYW5kIHN0YXRlLiBDb250'
    'YWN0cyBhbmQgYWRkcmVzc2VzIGFyZSBtYW5hZ2VkIHZpYSBzZXBhcmF0ZSBSUENzLioNdXBkYX'
    'RlUHJvZmlsZYK1GBAKDnByb2ZpbGVfdXBkYXRlEp4CCgpBZGRDb250YWN0Eh0ucHJvZmlsZS52'
    'MS5BZGRDb250YWN0UmVxdWVzdBoeLnByb2ZpbGUudjEuQWRkQ29udGFjdFJlc3BvbnNlItABuk'
    'e4AQoIQ29udGFjdHMSFkFkZCBjb250YWN0IHRvIHByb2ZpbGUahwFBZGRzIGEgbmV3IGNvbnRh'
    'Y3QgKGVtYWlsIG9yIHBob25lKSB0byBhIHByb2ZpbGUgYW5kIGluaXRpYXRlcyBhdXRvbWF0aW'
    'MgdmVyaWZpY2F0aW9uLiBSZXR1cm5zIHRoZSB1cGRhdGVkIHByb2ZpbGUgYW5kIHZlcmlmaWNh'
    'dGlvbiBJRC4qCmFkZENvbnRhY3SCtRgQCg5jb250YWN0X21hbmFnZRKPAgoNQ3JlYXRlQ29udG'
    'FjdBIgLnByb2ZpbGUudjEuQ3JlYXRlQ29udGFjdFJlcXVlc3QaIS5wcm9maWxlLnYxLkNyZWF0'
    'ZUNvbnRhY3RSZXNwb25zZSK4AbpHoAEKCENvbnRhY3RzEhlDcmVhdGUgc3RhbmRhbG9uZSBjb2'
    '50YWN0GmpDcmVhdGVzIGEgc3RhbmRhbG9uZSBjb250YWN0IHRoYXQgY2FuIGxhdGVyIGJlIGxp'
    'bmtlZCB0byBhIHByb2ZpbGUuIFVzZWZ1bCBmb3IgcHJlLXJlZ2lzdHJhdGlvbiBzY2VuYXJpb3'
    'MuKg1jcmVhdGVDb250YWN0grUYEAoOY29udGFjdF9tYW5hZ2USzwIKC0dldENvbnRhY3RzEh4u'
    'cHJvZmlsZS52MS5HZXRDb250YWN0c1JlcXVlc3QaHy5wcm9maWxlLnYxLkdldENvbnRhY3RzUm'
    'VzcG9uc2Ui/gGQAgG6R+MBCghDb250YWN0cxISR2V0IGNvbnRhY3RzIGJ5IGlkGrUBRmV0Y2hl'
    'cyBDb250YWN0T2JqZWN0cyBmb3Igb25lIG9yIG1vcmUgY29udGFjdCBpZHMgcmVnYXJkbGVzcy'
    'BvZiBwcm9maWxlIGF0dGFjaG1lbnQuIFVzZSBhIHNpbmdsZSBpZCBmb3Igb25lIGNvbnRhY3Qs'
    'IG9yIG1hbnkgZm9yIGJhdGNoIHJlc29sdmUgKGUuZy4gcHJvZHVjdC1zdG9yZWQgY3ZfY29udG'
    'FjdF9pZHMpLioLZ2V0Q29udGFjdHOCtRgQCg5jb250YWN0X21hbmFnZRLVAgoZQ3JlYXRlQ29u'
    'dGFjdFZlcmlmaWNhdGlvbhIsLnByb2ZpbGUudjEuQ3JlYXRlQ29udGFjdFZlcmlmaWNhdGlvbl'
    'JlcXVlc3QaLS5wcm9maWxlLnYxLkNyZWF0ZUNvbnRhY3RWZXJpZmljYXRpb25SZXNwb25zZSLa'
    'AbpHwgEKCENvbnRhY3RzEhtDcmVhdGUgY29udGFjdCB2ZXJpZmljYXRpb24afkluaXRpYXRlcy'
    'Bjb250YWN0IHZlcmlmaWNhdGlvbiBieSBzZW5kaW5nIGEgdmVyaWZpY2F0aW9uIGNvZGUgdmlh'
    'IGVtYWlsIG9yIFNNUy4gVGhlIGNvZGUgZXhwaXJlcyBhZnRlciB0aGUgc3BlY2lmaWVkIGR1cm'
    'F0aW9uLioZY3JlYXRlQ29udGFjdFZlcmlmaWNhdGlvboK1GBAKDmNvbnRhY3RfbWFuYWdlEqoC'
    'ChFDaGVja1ZlcmlmaWNhdGlvbhIkLnByb2ZpbGUudjEuQ2hlY2tWZXJpZmljYXRpb25SZXF1ZX'
    'N0GiUucHJvZmlsZS52MS5DaGVja1ZlcmlmaWNhdGlvblJlc3BvbnNlIscBukevAQoIQ29udGFj'
    'dHMSF0NoZWNrIHZlcmlmaWNhdGlvbiBjb2RlGndWZXJpZmllcyBhIGNvbnRhY3QgYnkgY2hlY2'
    'tpbmcgdGhlIHByb3ZpZGVkIHZlcmlmaWNhdGlvbiBjb2RlLiBUcmFja3MgdmVyaWZpY2F0aW9u'
    'IGF0dGVtcHRzIGFuZCByZXR1cm5zIHN1Y2Nlc3Mgc3RhdHVzLioRY2hlY2tWZXJpZmljYXRpb2'
    '6CtRgQCg5jb250YWN0X21hbmFnZRL2AQoNUmVtb3ZlQ29udGFjdBIgLnByb2ZpbGUudjEuUmVt'
    'b3ZlQ29udGFjdFJlcXVlc3QaIS5wcm9maWxlLnYxLlJlbW92ZUNvbnRhY3RSZXNwb25zZSKfAb'
    'pHhwEKCENvbnRhY3RzEg5SZW1vdmUgY29udGFjdBpcUmVtb3ZlcyBhIGNvbnRhY3QgZnJvbSBh'
    'IHByb2ZpbGUuIFRoZSBjb250YWN0IGlzIGRpc2Fzc29jaWF0ZWQgYnV0IG1heSByZW1haW4gaW'
    '4gdGhlIHN5c3RlbS4qDXJlbW92ZUNvbnRhY3SCtRgQCg5jb250YWN0X21hbmFnZRKoAgoMU2Vh'
    'cmNoUm9zdGVyEh8ucHJvZmlsZS52MS5TZWFyY2hSb3N0ZXJSZXF1ZXN0GiAucHJvZmlsZS52MS'
    '5TZWFyY2hSb3N0ZXJSZXNwb25zZSLSAZACAbpHugEKBlJvc3RlchINU2VhcmNoIHJvc3RlchqS'
    'AVNlYXJjaGVzIGEgdXNlcidzIGNvbnRhY3Qgcm9zdGVyIChjb250YWN0IGxpc3QpIHdpdGggZm'
    'lsdGVyaW5nIGJ5IGRhdGUgcmFuZ2UsIHByb3BlcnRpZXMsIGFuZCBjdXN0b20gY3JpdGVyaWEu'
    'IFJldHVybnMgYSBzdHJlYW0gb2Ygcm9zdGVyIGVudHJpZXMuKgxzZWFyY2hSb3N0ZXKCtRgNCg'
    'tyb3N0ZXJfdmlldzABEuwBCglBZGRSb3N0ZXISHC5wcm9maWxlLnYxLkFkZFJvc3RlclJlcXVl'
    'c3QaHS5wcm9maWxlLnYxLkFkZFJvc3RlclJlc3BvbnNlIqEBukeKAQoGUm9zdGVyEhJBZGQgcm'
    '9zdGVyIGVudHJpZXMaYUFkZHMgbXVsdGlwbGUgY29udGFjdHMgdG8gYSB1c2VyJ3Mgcm9zdGVy'
    'IChjb250YWN0IGxpc3QpLiBFYWNoIGNvbnRhY3QgaXMgdmVyaWZpZWQgYXV0b21hdGljYWxseS'
    '4qCWFkZFJvc3RlcoK1GA8KDXJvc3Rlcl9tYW5hZ2USiwIKDFJlbW92ZVJvc3RlchIfLnByb2Zp'
    'bGUudjEuUmVtb3ZlUm9zdGVyUmVxdWVzdBogLnByb2ZpbGUudjEuUmVtb3ZlUm9zdGVyUmVzcG'
    '9uc2UitwG6R6ABCgZSb3N0ZXISE1JlbW92ZSByb3N0ZXIgZW50cnkac1JlbW92ZXMgYSBjb250'
    'YWN0IGZyb20gYSB1c2VyJ3Mgcm9zdGVyIChjb250YWN0IGxpc3QpLiBUaGUgcHJvZmlsZSByZW'
    '1haW5zIGJ1dCBpcyBubyBsb25nZXIgaW4gdGhlIHVzZXIncyBjb250YWN0cy4qDHJlbW92ZVJv'
    'c3RlcoK1GA8KDXJvc3Rlcl9tYW5hZ2US4QEKCkFkZEFkZHJlc3MSHS5wcm9maWxlLnYxLkFkZE'
    'FkZHJlc3NSZXF1ZXN0Gh4ucHJvZmlsZS52MS5BZGRBZGRyZXNzUmVzcG9uc2UikwG6R3wKCUFk'
    'ZHJlc3NlcxILQWRkIGFkZHJlc3MaVkFkZHMgYSBuZXcgcGh5c2ljYWwgYWRkcmVzcyB0byBhIH'
    'Byb2ZpbGUgd2l0aCBvcHRpb25hbCBnZW9jb2RpbmcgKGxhdGl0dWRlL2xvbmdpdHVkZSkuKgph'
    'ZGRBZGRyZXNzgrUYEAoOYWRkcmVzc19tYW5hZ2USoQIKD0FkZFJlbGF0aW9uc2hpcBIiLnByb2'
    'ZpbGUudjEuQWRkUmVsYXRpb25zaGlwUmVxdWVzdBojLnByb2ZpbGUudjEuQWRkUmVsYXRpb25z'
    'aGlwUmVzcG9uc2UixAG6R6cBCg1SZWxhdGlvbnNoaXBzEhBBZGQgcmVsYXRpb25zaGlwGnNDcm'
    'VhdGVzIGEgcmVsYXRpb25zaGlwIGJldHdlZW4gdHdvIHByb2ZpbGVzIChtZW1iZXIsIGFmZmls'
    'aWF0ZWQsIGJsYWNrbGlzdGVkKS4gU3VwcG9ydHMgaGllcmFyY2hpY2FsIHJlbGF0aW9uc2hpcH'
    'MuKg9hZGRSZWxhdGlvbnNoaXCCtRgVChNyZWxhdGlvbnNoaXBfbWFuYWdlEp0CChJEZWxldGVS'
    'ZWxhdGlvbnNoaXASJS5wcm9maWxlLnYxLkRlbGV0ZVJlbGF0aW9uc2hpcFJlcXVlc3QaJi5wcm'
    '9maWxlLnYxLkRlbGV0ZVJlbGF0aW9uc2hpcFJlc3BvbnNlIrcBukeaAQoNUmVsYXRpb25zaGlw'
    'cxITRGVsZXRlIHJlbGF0aW9uc2hpcBpgUmVtb3ZlcyBhbiBleGlzdGluZyByZWxhdGlvbnNoaX'
    'AgYmV0d2VlbiBwcm9maWxlcy4gVGhlIHByb2ZpbGVzIHJlbWFpbiBidXQgYXJlIG5vIGxvbmdl'
    'ciBsaW5rZWQuKhJkZWxldGVSZWxhdGlvbnNoaXCCtRgVChNyZWxhdGlvbnNoaXBfbWFuYWdlEs'
    'QCChBMaXN0UmVsYXRpb25zaGlwEiMucHJvZmlsZS52MS5MaXN0UmVsYXRpb25zaGlwUmVxdWVz'
    'dBokLnByb2ZpbGUudjEuTGlzdFJlbGF0aW9uc2hpcFJlc3BvbnNlIuIBkAIBukfEAQoNUmVsYX'
    'Rpb25zaGlwcxISTGlzdCByZWxhdGlvbnNoaXBzGosBTGlzdHMgYWxsIHJlbGF0aW9uc2hpcHMg'
    'Zm9yIGEgcHJvZmlsZSB3aXRoIG9wdGlvbmFsIGZpbHRlcmluZyBieSB0eXBlIGFuZCByZWxhdG'
    'VkIHByb2ZpbGVzLiBTdXBwb3J0cyBwYWdpbmF0aW9uIGFuZCByZWxhdGlvbnNoaXAgaW52ZXJz'
    'aW9uLioRbGlzdFJlbGF0aW9uc2hpcHOCtRgTChFyZWxhdGlvbnNoaXBfdmlldzABEqMCChNHZX'
    'RCeUlEQW5kUGFydGl0aW9uEiYucHJvZmlsZS52MS5HZXRCeUlEQW5kUGFydGl0aW9uUmVxdWVz'
    'dBonLnByb2ZpbGUudjEuR2V0QnlJREFuZFBhcnRpdGlvblJlc3BvbnNlIroBkAIBukehAQoIUH'
    'JvZmlsZXMSH0dldCBwcm9maWxlIGJ5IElEIGFuZCBwYXJ0aXRpb24aWFJldHJpZXZlcyBhIHBy'
    'b2ZpbGUgYnkgSUQgd2l0aCB0ZW5hbnQtc2NvcGVkIHByb3BlcnRpZXMgbWVyZ2VkIGludG8gdG'
    'hlIGJhc2UgcHJvcGVydGllcy4qGmdldFByb2ZpbGVCeUlEQW5kUGFydGl0aW9ugrUYDgoMcHJv'
    'ZmlsZV92aWV3EpoCCg9Qcm9wZXJ0eUhpc3RvcnkSIi5wcm9maWxlLnYxLlByb3BlcnR5SGlzdG'
    '9yeVJlcXVlc3QaIy5wcm9maWxlLnYxLlByb3BlcnR5SGlzdG9yeVJlc3BvbnNlIr0BkAIBukek'
    'AQoIUHJvZmlsZXMSG0dldCBwcm9wZXJ0eSBjaGFuZ2UgaGlzdG9yeRpqUmV0dXJucyB0aGUgY2'
    'hhbmdlIGhpc3RvcnkgZm9yIGEgc3BlY2lmaWMgcHJvcGVydHkga2V5IG9uIGEgcHJvZmlsZSwg'
    'ZmlsdGVyZWQgYnkgY2FsbGVyIHRlbmFudCB2aXNpYmlsaXR5LioPcHJvcGVydHlIaXN0b3J5gr'
    'UYDgoMcHJvZmlsZV92aWV3ErICCg9SZXNvbHZlQWNjb3VudHMSIi5wcm9maWxlLnYxLlJlc29s'
    'dmVBY2NvdW50c1JlcXVlc3QaIy5wcm9maWxlLnYxLlJlc29sdmVBY2NvdW50c1Jlc3BvbnNlIt'
    'UBkAIBuke5AQoIQWNjb3VudHMSFlJlc29sdmUgYWNjb3VudCBvd25lcnMagwFNYXBzIHVwIHRv'
    'IDUwMCBjaGFpbiBhY2NvdW50IGFkZHJlc3NlcyB0byB0aGUgcHJvZmlsZXMgdGhhdCBvd24gdG'
    'hlbS4gU2VydmljZSBhY2NvdW50cyBvbmx5OyBhZGRyZXNzZXMgbm8gcHJvZmlsZSBvd25zIGFy'
    'ZSBvbWl0dGVkLioPcmVzb2x2ZUFjY291bnRzgrUYEQoPYWNjb3VudF9yZXNvbHZlGsAHgrUYuw'
    'cKD3NlcnZpY2VfcHJvZmlsZRIMcHJvZmlsZV92aWV3Eg5wcm9maWxlX2NyZWF0ZRIOcHJvZmls'
    'ZV91cGRhdGUSDXByb2ZpbGVfbWVyZ2USDmNvbnRhY3RfbWFuYWdlEgtyb3N0ZXJfdmlldxINcm'
    '9zdGVyX21hbmFnZRIOYWRkcmVzc19tYW5hZ2USEXJlbGF0aW9uc2hpcF92aWV3EhNyZWxhdGlv'
    'bnNoaXBfbWFuYWdlEg9hY2NvdW50X3Jlc29sdmUaowEIARIMcHJvZmlsZV92aWV3Eg5wcm9maW'
    'xlX2NyZWF0ZRIOcHJvZmlsZV91cGRhdGUSDXByb2ZpbGVfbWVyZ2USDmNvbnRhY3RfbWFuYWdl'
    'Egtyb3N0ZXJfdmlldxINcm9zdGVyX21hbmFnZRIOYWRkcmVzc19tYW5hZ2USEXJlbGF0aW9uc2'
    'hpcF92aWV3EhNyZWxhdGlvbnNoaXBfbWFuYWdlGqMBCAISDHByb2ZpbGVfdmlldxIOcHJvZmls'
    'ZV9jcmVhdGUSDnByb2ZpbGVfdXBkYXRlEg1wcm9maWxlX21lcmdlEg5jb250YWN0X21hbmFnZR'
    'ILcm9zdGVyX3ZpZXcSDXJvc3Rlcl9tYW5hZ2USDmFkZHJlc3NfbWFuYWdlEhFyZWxhdGlvbnNo'
    'aXBfdmlldxITcmVsYXRpb25zaGlwX21hbmFnZRpfCAMSDHByb2ZpbGVfdmlldxIOcHJvZmlsZV'
    '91cGRhdGUSDmNvbnRhY3RfbWFuYWdlEgtyb3N0ZXJfdmlldxINcm9zdGVyX21hbmFnZRIRcmVs'
    'YXRpb25zaGlwX3ZpZXcaMAgEEgxwcm9maWxlX3ZpZXcSC3Jvc3Rlcl92aWV3EhFyZWxhdGlvbn'
    'NoaXBfdmlldxpgCAUSDHByb2ZpbGVfdmlldxIOcHJvZmlsZV91cGRhdGUSDmNvbnRhY3RfbWFu'
    'YWdlEgtyb3N0ZXJfdmlldxIOYWRkcmVzc19tYW5hZ2USEXJlbGF0aW9uc2hpcF92aWV3GrQBCA'
    'YSDHByb2ZpbGVfdmlldxIOcHJvZmlsZV9jcmVhdGUSDnByb2ZpbGVfdXBkYXRlEg1wcm9maWxl'
    'X21lcmdlEg5jb250YWN0X21hbmFnZRILcm9zdGVyX3ZpZXcSDXJvc3Rlcl9tYW5hZ2USDmFkZH'
    'Jlc3NfbWFuYWdlEhFyZWxhdGlvbnNoaXBfdmlldxITcmVsYXRpb25zaGlwX21hbmFnZRIPYWNj'
    'b3VudF9yZXNvbHZl');

