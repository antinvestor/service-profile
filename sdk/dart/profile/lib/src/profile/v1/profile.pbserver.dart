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

import 'package:protobuf/protobuf.dart' as $pb;

import 'profile.pb.dart' as $3;
import 'profile.pbjson.dart';

export 'profile.pb.dart';

abstract class ProfileServiceBase extends $pb.GeneratedService {
  $async.Future<$3.GetByIdResponse> getById(
      $pb.ServerContext ctx, $3.GetByIdRequest request);
  $async.Future<$3.GetByContactResponse> getByContact(
      $pb.ServerContext ctx, $3.GetByContactRequest request);
  $async.Future<$3.SearchResponse> search(
      $pb.ServerContext ctx, $3.SearchRequest request);
  $async.Future<$3.MergeResponse> merge(
      $pb.ServerContext ctx, $3.MergeRequest request);
  $async.Future<$3.CreateResponse> create(
      $pb.ServerContext ctx, $3.CreateRequest request);
  $async.Future<$3.UpdateResponse> update(
      $pb.ServerContext ctx, $3.UpdateRequest request);
  $async.Future<$3.AddContactResponse> addContact(
      $pb.ServerContext ctx, $3.AddContactRequest request);
  $async.Future<$3.CreateContactResponse> createContact(
      $pb.ServerContext ctx, $3.CreateContactRequest request);
  $async.Future<$3.GetContactsResponse> getContacts(
      $pb.ServerContext ctx, $3.GetContactsRequest request);
  $async.Future<$3.CreateContactVerificationResponse> createContactVerification(
      $pb.ServerContext ctx, $3.CreateContactVerificationRequest request);
  $async.Future<$3.CheckVerificationResponse> checkVerification(
      $pb.ServerContext ctx, $3.CheckVerificationRequest request);
  $async.Future<$3.RemoveContactResponse> removeContact(
      $pb.ServerContext ctx, $3.RemoveContactRequest request);
  $async.Future<$3.SearchRosterResponse> searchRoster(
      $pb.ServerContext ctx, $3.SearchRosterRequest request);
  $async.Future<$3.AddRosterResponse> addRoster(
      $pb.ServerContext ctx, $3.AddRosterRequest request);
  $async.Future<$3.RemoveRosterResponse> removeRoster(
      $pb.ServerContext ctx, $3.RemoveRosterRequest request);
  $async.Future<$3.AddAddressResponse> addAddress(
      $pb.ServerContext ctx, $3.AddAddressRequest request);
  $async.Future<$3.AddRelationshipResponse> addRelationship(
      $pb.ServerContext ctx, $3.AddRelationshipRequest request);
  $async.Future<$3.DeleteRelationshipResponse> deleteRelationship(
      $pb.ServerContext ctx, $3.DeleteRelationshipRequest request);
  $async.Future<$3.ListRelationshipResponse> listRelationship(
      $pb.ServerContext ctx, $3.ListRelationshipRequest request);
  $async.Future<$3.GetByIDAndPartitionResponse> getByIDAndPartition(
      $pb.ServerContext ctx, $3.GetByIDAndPartitionRequest request);
  $async.Future<$3.PropertyHistoryResponse> propertyHistory(
      $pb.ServerContext ctx, $3.PropertyHistoryRequest request);
  $async.Future<$3.ResolveAccountsResponse> resolveAccounts(
      $pb.ServerContext ctx, $3.ResolveAccountsRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'GetById':
        return $3.GetByIdRequest();
      case 'GetByContact':
        return $3.GetByContactRequest();
      case 'Search':
        return $3.SearchRequest();
      case 'Merge':
        return $3.MergeRequest();
      case 'Create':
        return $3.CreateRequest();
      case 'Update':
        return $3.UpdateRequest();
      case 'AddContact':
        return $3.AddContactRequest();
      case 'CreateContact':
        return $3.CreateContactRequest();
      case 'GetContacts':
        return $3.GetContactsRequest();
      case 'CreateContactVerification':
        return $3.CreateContactVerificationRequest();
      case 'CheckVerification':
        return $3.CheckVerificationRequest();
      case 'RemoveContact':
        return $3.RemoveContactRequest();
      case 'SearchRoster':
        return $3.SearchRosterRequest();
      case 'AddRoster':
        return $3.AddRosterRequest();
      case 'RemoveRoster':
        return $3.RemoveRosterRequest();
      case 'AddAddress':
        return $3.AddAddressRequest();
      case 'AddRelationship':
        return $3.AddRelationshipRequest();
      case 'DeleteRelationship':
        return $3.DeleteRelationshipRequest();
      case 'ListRelationship':
        return $3.ListRelationshipRequest();
      case 'GetByIDAndPartition':
        return $3.GetByIDAndPartitionRequest();
      case 'PropertyHistory':
        return $3.PropertyHistoryRequest();
      case 'ResolveAccounts':
        return $3.ResolveAccountsRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'GetById':
        return getById(ctx, request as $3.GetByIdRequest);
      case 'GetByContact':
        return getByContact(ctx, request as $3.GetByContactRequest);
      case 'Search':
        return search(ctx, request as $3.SearchRequest);
      case 'Merge':
        return merge(ctx, request as $3.MergeRequest);
      case 'Create':
        return create(ctx, request as $3.CreateRequest);
      case 'Update':
        return update(ctx, request as $3.UpdateRequest);
      case 'AddContact':
        return addContact(ctx, request as $3.AddContactRequest);
      case 'CreateContact':
        return createContact(ctx, request as $3.CreateContactRequest);
      case 'GetContacts':
        return getContacts(ctx, request as $3.GetContactsRequest);
      case 'CreateContactVerification':
        return createContactVerification(
            ctx, request as $3.CreateContactVerificationRequest);
      case 'CheckVerification':
        return checkVerification(ctx, request as $3.CheckVerificationRequest);
      case 'RemoveContact':
        return removeContact(ctx, request as $3.RemoveContactRequest);
      case 'SearchRoster':
        return searchRoster(ctx, request as $3.SearchRosterRequest);
      case 'AddRoster':
        return addRoster(ctx, request as $3.AddRosterRequest);
      case 'RemoveRoster':
        return removeRoster(ctx, request as $3.RemoveRosterRequest);
      case 'AddAddress':
        return addAddress(ctx, request as $3.AddAddressRequest);
      case 'AddRelationship':
        return addRelationship(ctx, request as $3.AddRelationshipRequest);
      case 'DeleteRelationship':
        return deleteRelationship(ctx, request as $3.DeleteRelationshipRequest);
      case 'ListRelationship':
        return listRelationship(ctx, request as $3.ListRelationshipRequest);
      case 'GetByIDAndPartition':
        return getByIDAndPartition(
            ctx, request as $3.GetByIDAndPartitionRequest);
      case 'PropertyHistory':
        return propertyHistory(ctx, request as $3.PropertyHistoryRequest);
      case 'ResolveAccounts':
        return resolveAccounts(ctx, request as $3.ResolveAccountsRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => ProfileServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => ProfileServiceBase$messageJson;
}
