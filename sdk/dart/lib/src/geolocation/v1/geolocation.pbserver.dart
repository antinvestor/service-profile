//
//  Generated code. Do not modify.
//  source: geolocation/v1/geolocation.proto
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

import '../../google/protobuf/empty.pb.dart' as $8;
import 'geolocation.pb.dart' as $9;
import 'geolocation.pbjson.dart';

export 'geolocation.pb.dart';

abstract class GeolocationServiceBase extends $pb.GeneratedService {
  $async.Future<$9.IngestLocationsResponse> ingestLocations($pb.ServerContext ctx, $9.IngestLocationsRequest request);
  $async.Future<$9.CreateAreaResponse> createArea($pb.ServerContext ctx, $9.CreateAreaRequest request);
  $async.Future<$9.GetAreaResponse> getArea($pb.ServerContext ctx, $9.GetAreaRequest request);
  $async.Future<$9.UpdateAreaResponse> updateArea($pb.ServerContext ctx, $9.UpdateAreaRequest request);
  $async.Future<$8.Empty> deleteArea($pb.ServerContext ctx, $9.DeleteAreaRequest request);
  $async.Future<$9.SearchAreasResponse> searchAreas($pb.ServerContext ctx, $9.SearchAreasRequest request);
  $async.Future<$9.CreateRouteResponse> createRoute($pb.ServerContext ctx, $9.CreateRouteRequest request);
  $async.Future<$9.GetRouteResponse> getRoute($pb.ServerContext ctx, $9.GetRouteRequest request);
  $async.Future<$9.UpdateRouteResponse> updateRoute($pb.ServerContext ctx, $9.UpdateRouteRequest request);
  $async.Future<$8.Empty> deleteRoute($pb.ServerContext ctx, $9.DeleteRouteRequest request);
  $async.Future<$9.SearchRoutesResponse> searchRoutes($pb.ServerContext ctx, $9.SearchRoutesRequest request);
  $async.Future<$9.AssignRouteResponse> assignRoute($pb.ServerContext ctx, $9.AssignRouteRequest request);
  $async.Future<$8.Empty> unassignRoute($pb.ServerContext ctx, $9.UnassignRouteRequest request);
  $async.Future<$9.GetSubjectRouteAssignmentsResponse> getSubjectRouteAssignments($pb.ServerContext ctx, $9.GetSubjectRouteAssignmentsRequest request);
  $async.Future<$9.GetTrackResponse> getTrack($pb.ServerContext ctx, $9.GetTrackRequest request);
  $async.Future<$9.GetSubjectEventsResponse> getSubjectEvents($pb.ServerContext ctx, $9.GetSubjectEventsRequest request);
  $async.Future<$9.GetAreaSubjectsResponse> getAreaSubjects($pb.ServerContext ctx, $9.GetAreaSubjectsRequest request);
  $async.Future<$9.GetNearbySubjectsResponse> getNearbySubjects($pb.ServerContext ctx, $9.GetNearbySubjectsRequest request);
  $async.Future<$9.GetNearbyAreasResponse> getNearbyAreas($pb.ServerContext ctx, $9.GetNearbyAreasRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'IngestLocations': return $9.IngestLocationsRequest();
      case 'CreateArea': return $9.CreateAreaRequest();
      case 'GetArea': return $9.GetAreaRequest();
      case 'UpdateArea': return $9.UpdateAreaRequest();
      case 'DeleteArea': return $9.DeleteAreaRequest();
      case 'SearchAreas': return $9.SearchAreasRequest();
      case 'CreateRoute': return $9.CreateRouteRequest();
      case 'GetRoute': return $9.GetRouteRequest();
      case 'UpdateRoute': return $9.UpdateRouteRequest();
      case 'DeleteRoute': return $9.DeleteRouteRequest();
      case 'SearchRoutes': return $9.SearchRoutesRequest();
      case 'AssignRoute': return $9.AssignRouteRequest();
      case 'UnassignRoute': return $9.UnassignRouteRequest();
      case 'GetSubjectRouteAssignments': return $9.GetSubjectRouteAssignmentsRequest();
      case 'GetTrack': return $9.GetTrackRequest();
      case 'GetSubjectEvents': return $9.GetSubjectEventsRequest();
      case 'GetAreaSubjects': return $9.GetAreaSubjectsRequest();
      case 'GetNearbySubjects': return $9.GetNearbySubjectsRequest();
      case 'GetNearbyAreas': return $9.GetNearbyAreasRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'IngestLocations': return this.ingestLocations(ctx, request as $9.IngestLocationsRequest);
      case 'CreateArea': return this.createArea(ctx, request as $9.CreateAreaRequest);
      case 'GetArea': return this.getArea(ctx, request as $9.GetAreaRequest);
      case 'UpdateArea': return this.updateArea(ctx, request as $9.UpdateAreaRequest);
      case 'DeleteArea': return this.deleteArea(ctx, request as $9.DeleteAreaRequest);
      case 'SearchAreas': return this.searchAreas(ctx, request as $9.SearchAreasRequest);
      case 'CreateRoute': return this.createRoute(ctx, request as $9.CreateRouteRequest);
      case 'GetRoute': return this.getRoute(ctx, request as $9.GetRouteRequest);
      case 'UpdateRoute': return this.updateRoute(ctx, request as $9.UpdateRouteRequest);
      case 'DeleteRoute': return this.deleteRoute(ctx, request as $9.DeleteRouteRequest);
      case 'SearchRoutes': return this.searchRoutes(ctx, request as $9.SearchRoutesRequest);
      case 'AssignRoute': return this.assignRoute(ctx, request as $9.AssignRouteRequest);
      case 'UnassignRoute': return this.unassignRoute(ctx, request as $9.UnassignRouteRequest);
      case 'GetSubjectRouteAssignments': return this.getSubjectRouteAssignments(ctx, request as $9.GetSubjectRouteAssignmentsRequest);
      case 'GetTrack': return this.getTrack(ctx, request as $9.GetTrackRequest);
      case 'GetSubjectEvents': return this.getSubjectEvents(ctx, request as $9.GetSubjectEventsRequest);
      case 'GetAreaSubjects': return this.getAreaSubjects(ctx, request as $9.GetAreaSubjectsRequest);
      case 'GetNearbySubjects': return this.getNearbySubjects(ctx, request as $9.GetNearbySubjectsRequest);
      case 'GetNearbyAreas': return this.getNearbyAreas(ctx, request as $9.GetNearbyAreasRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => GeolocationServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => GeolocationServiceBase$messageJson;
}

