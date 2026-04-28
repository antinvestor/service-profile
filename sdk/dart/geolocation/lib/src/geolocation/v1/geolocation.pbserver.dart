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

import '../../google/protobuf/empty.pb.dart' as $2;
import 'geolocation.pb.dart' as $3;
import 'geolocation.pbjson.dart';

export 'geolocation.pb.dart';

abstract class GeolocationServiceBase extends $pb.GeneratedService {
  $async.Future<$3.IngestLocationsResponse> ingestLocations($pb.ServerContext ctx, $3.IngestLocationsRequest request);
  $async.Future<$3.CreateAreaResponse> createArea($pb.ServerContext ctx, $3.CreateAreaRequest request);
  $async.Future<$3.GetAreaResponse> getArea($pb.ServerContext ctx, $3.GetAreaRequest request);
  $async.Future<$3.UpdateAreaResponse> updateArea($pb.ServerContext ctx, $3.UpdateAreaRequest request);
  $async.Future<$2.Empty> deleteArea($pb.ServerContext ctx, $3.DeleteAreaRequest request);
  $async.Future<$3.SearchAreasResponse> searchAreas($pb.ServerContext ctx, $3.SearchAreasRequest request);
  $async.Future<$3.CreateRouteResponse> createRoute($pb.ServerContext ctx, $3.CreateRouteRequest request);
  $async.Future<$3.GetRouteResponse> getRoute($pb.ServerContext ctx, $3.GetRouteRequest request);
  $async.Future<$3.UpdateRouteResponse> updateRoute($pb.ServerContext ctx, $3.UpdateRouteRequest request);
  $async.Future<$2.Empty> deleteRoute($pb.ServerContext ctx, $3.DeleteRouteRequest request);
  $async.Future<$3.SearchRoutesResponse> searchRoutes($pb.ServerContext ctx, $3.SearchRoutesRequest request);
  $async.Future<$3.AssignRouteResponse> assignRoute($pb.ServerContext ctx, $3.AssignRouteRequest request);
  $async.Future<$2.Empty> unassignRoute($pb.ServerContext ctx, $3.UnassignRouteRequest request);
  $async.Future<$3.GetSubjectRouteAssignmentsResponse> getSubjectRouteAssignments($pb.ServerContext ctx, $3.GetSubjectRouteAssignmentsRequest request);
  $async.Future<$3.GetTrackResponse> getTrack($pb.ServerContext ctx, $3.GetTrackRequest request);
  $async.Future<$3.GetSubjectEventsResponse> getSubjectEvents($pb.ServerContext ctx, $3.GetSubjectEventsRequest request);
  $async.Future<$3.GetAreaSubjectsResponse> getAreaSubjects($pb.ServerContext ctx, $3.GetAreaSubjectsRequest request);
  $async.Future<$3.GetNearbySubjectsResponse> getNearbySubjects($pb.ServerContext ctx, $3.GetNearbySubjectsRequest request);
  $async.Future<$3.GetNearbyAreasResponse> getNearbyAreas($pb.ServerContext ctx, $3.GetNearbyAreasRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'IngestLocations': return $3.IngestLocationsRequest();
      case 'CreateArea': return $3.CreateAreaRequest();
      case 'GetArea': return $3.GetAreaRequest();
      case 'UpdateArea': return $3.UpdateAreaRequest();
      case 'DeleteArea': return $3.DeleteAreaRequest();
      case 'SearchAreas': return $3.SearchAreasRequest();
      case 'CreateRoute': return $3.CreateRouteRequest();
      case 'GetRoute': return $3.GetRouteRequest();
      case 'UpdateRoute': return $3.UpdateRouteRequest();
      case 'DeleteRoute': return $3.DeleteRouteRequest();
      case 'SearchRoutes': return $3.SearchRoutesRequest();
      case 'AssignRoute': return $3.AssignRouteRequest();
      case 'UnassignRoute': return $3.UnassignRouteRequest();
      case 'GetSubjectRouteAssignments': return $3.GetSubjectRouteAssignmentsRequest();
      case 'GetTrack': return $3.GetTrackRequest();
      case 'GetSubjectEvents': return $3.GetSubjectEventsRequest();
      case 'GetAreaSubjects': return $3.GetAreaSubjectsRequest();
      case 'GetNearbySubjects': return $3.GetNearbySubjectsRequest();
      case 'GetNearbyAreas': return $3.GetNearbyAreasRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'IngestLocations': return this.ingestLocations(ctx, request as $3.IngestLocationsRequest);
      case 'CreateArea': return this.createArea(ctx, request as $3.CreateAreaRequest);
      case 'GetArea': return this.getArea(ctx, request as $3.GetAreaRequest);
      case 'UpdateArea': return this.updateArea(ctx, request as $3.UpdateAreaRequest);
      case 'DeleteArea': return this.deleteArea(ctx, request as $3.DeleteAreaRequest);
      case 'SearchAreas': return this.searchAreas(ctx, request as $3.SearchAreasRequest);
      case 'CreateRoute': return this.createRoute(ctx, request as $3.CreateRouteRequest);
      case 'GetRoute': return this.getRoute(ctx, request as $3.GetRouteRequest);
      case 'UpdateRoute': return this.updateRoute(ctx, request as $3.UpdateRouteRequest);
      case 'DeleteRoute': return this.deleteRoute(ctx, request as $3.DeleteRouteRequest);
      case 'SearchRoutes': return this.searchRoutes(ctx, request as $3.SearchRoutesRequest);
      case 'AssignRoute': return this.assignRoute(ctx, request as $3.AssignRouteRequest);
      case 'UnassignRoute': return this.unassignRoute(ctx, request as $3.UnassignRouteRequest);
      case 'GetSubjectRouteAssignments': return this.getSubjectRouteAssignments(ctx, request as $3.GetSubjectRouteAssignmentsRequest);
      case 'GetTrack': return this.getTrack(ctx, request as $3.GetTrackRequest);
      case 'GetSubjectEvents': return this.getSubjectEvents(ctx, request as $3.GetSubjectEventsRequest);
      case 'GetAreaSubjects': return this.getAreaSubjects(ctx, request as $3.GetAreaSubjectsRequest);
      case 'GetNearbySubjects': return this.getNearbySubjects(ctx, request as $3.GetNearbySubjectsRequest);
      case 'GetNearbyAreas': return this.getNearbyAreas(ctx, request as $3.GetNearbyAreasRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => GeolocationServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => GeolocationServiceBase$messageJson;
}

