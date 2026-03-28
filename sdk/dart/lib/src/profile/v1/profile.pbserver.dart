//
//  Generated code. Do not modify.
//  source: profile/v1/profile.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'profile.pb.dart' as $12;
import 'profile.pbjson.dart';

export 'profile.pb.dart';

abstract class ProfileServiceBase extends $pb.GeneratedService {
  $async.Future<$12.GetByIdResponse> getById($pb.ServerContext ctx, $12.GetByIdRequest request);
  $async.Future<$12.GetByContactResponse> getByContact($pb.ServerContext ctx, $12.GetByContactRequest request);
  $async.Future<$12.SearchResponse> search($pb.ServerContext ctx, $12.SearchRequest request);
  $async.Future<$12.MergeResponse> merge($pb.ServerContext ctx, $12.MergeRequest request);
  $async.Future<$12.CreateResponse> create($pb.ServerContext ctx, $12.CreateRequest request);
  $async.Future<$12.UpdateResponse> update($pb.ServerContext ctx, $12.UpdateRequest request);
  $async.Future<$12.AddContactResponse> addContact($pb.ServerContext ctx, $12.AddContactRequest request);
  $async.Future<$12.CreateContactResponse> createContact($pb.ServerContext ctx, $12.CreateContactRequest request);
  $async.Future<$12.CreateContactVerificationResponse> createContactVerification($pb.ServerContext ctx, $12.CreateContactVerificationRequest request);
  $async.Future<$12.CheckVerificationResponse> checkVerification($pb.ServerContext ctx, $12.CheckVerificationRequest request);
  $async.Future<$12.RemoveContactResponse> removeContact($pb.ServerContext ctx, $12.RemoveContactRequest request);
  $async.Future<$12.SearchRosterResponse> searchRoster($pb.ServerContext ctx, $12.SearchRosterRequest request);
  $async.Future<$12.AddRosterResponse> addRoster($pb.ServerContext ctx, $12.AddRosterRequest request);
  $async.Future<$12.RemoveRosterResponse> removeRoster($pb.ServerContext ctx, $12.RemoveRosterRequest request);
  $async.Future<$12.AddAddressResponse> addAddress($pb.ServerContext ctx, $12.AddAddressRequest request);
  $async.Future<$12.AddRelationshipResponse> addRelationship($pb.ServerContext ctx, $12.AddRelationshipRequest request);
  $async.Future<$12.DeleteRelationshipResponse> deleteRelationship($pb.ServerContext ctx, $12.DeleteRelationshipRequest request);
  $async.Future<$12.ListRelationshipResponse> listRelationship($pb.ServerContext ctx, $12.ListRelationshipRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'GetById': return $12.GetByIdRequest();
      case 'GetByContact': return $12.GetByContactRequest();
      case 'Search': return $12.SearchRequest();
      case 'Merge': return $12.MergeRequest();
      case 'Create': return $12.CreateRequest();
      case 'Update': return $12.UpdateRequest();
      case 'AddContact': return $12.AddContactRequest();
      case 'CreateContact': return $12.CreateContactRequest();
      case 'CreateContactVerification': return $12.CreateContactVerificationRequest();
      case 'CheckVerification': return $12.CheckVerificationRequest();
      case 'RemoveContact': return $12.RemoveContactRequest();
      case 'SearchRoster': return $12.SearchRosterRequest();
      case 'AddRoster': return $12.AddRosterRequest();
      case 'RemoveRoster': return $12.RemoveRosterRequest();
      case 'AddAddress': return $12.AddAddressRequest();
      case 'AddRelationship': return $12.AddRelationshipRequest();
      case 'DeleteRelationship': return $12.DeleteRelationshipRequest();
      case 'ListRelationship': return $12.ListRelationshipRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'GetById': return this.getById(ctx, request as $12.GetByIdRequest);
      case 'GetByContact': return this.getByContact(ctx, request as $12.GetByContactRequest);
      case 'Search': return this.search(ctx, request as $12.SearchRequest);
      case 'Merge': return this.merge(ctx, request as $12.MergeRequest);
      case 'Create': return this.create(ctx, request as $12.CreateRequest);
      case 'Update': return this.update(ctx, request as $12.UpdateRequest);
      case 'AddContact': return this.addContact(ctx, request as $12.AddContactRequest);
      case 'CreateContact': return this.createContact(ctx, request as $12.CreateContactRequest);
      case 'CreateContactVerification': return this.createContactVerification(ctx, request as $12.CreateContactVerificationRequest);
      case 'CheckVerification': return this.checkVerification(ctx, request as $12.CheckVerificationRequest);
      case 'RemoveContact': return this.removeContact(ctx, request as $12.RemoveContactRequest);
      case 'SearchRoster': return this.searchRoster(ctx, request as $12.SearchRosterRequest);
      case 'AddRoster': return this.addRoster(ctx, request as $12.AddRosterRequest);
      case 'RemoveRoster': return this.removeRoster(ctx, request as $12.RemoveRosterRequest);
      case 'AddAddress': return this.addAddress(ctx, request as $12.AddAddressRequest);
      case 'AddRelationship': return this.addRelationship(ctx, request as $12.AddRelationshipRequest);
      case 'DeleteRelationship': return this.deleteRelationship(ctx, request as $12.DeleteRelationshipRequest);
      case 'ListRelationship': return this.listRelationship(ctx, request as $12.ListRelationshipRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => ProfileServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => ProfileServiceBase$messageJson;
}

