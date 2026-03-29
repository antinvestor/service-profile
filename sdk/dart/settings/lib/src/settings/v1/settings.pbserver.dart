//
//  Generated code. Do not modify.
//  source: settings/v1/settings.proto
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

import '../../common/v1/common.pb.dart' as $10;
import 'settings.pb.dart' as $13;
import 'settings.pbjson.dart';

export 'settings.pb.dart';

abstract class SettingsServiceBase extends $pb.GeneratedService {
  $async.Future<$13.GetResponse> get($pb.ServerContext ctx, $13.GetRequest request);
  $async.Future<$13.ListResponse> list($pb.ServerContext ctx, $13.ListRequest request);
  $async.Future<$13.SearchResponse> search($pb.ServerContext ctx, $10.SearchRequest request);
  $async.Future<$13.SetResponse> set($pb.ServerContext ctx, $13.SetRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'Get': return $13.GetRequest();
      case 'List': return $13.ListRequest();
      case 'Search': return $10.SearchRequest();
      case 'Set': return $13.SetRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'Get': return this.get(ctx, request as $13.GetRequest);
      case 'List': return this.list(ctx, request as $13.ListRequest);
      case 'Search': return this.search(ctx, request as $10.SearchRequest);
      case 'Set': return this.set(ctx, request as $13.SetRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => SettingsServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => SettingsServiceBase$messageJson;
}

