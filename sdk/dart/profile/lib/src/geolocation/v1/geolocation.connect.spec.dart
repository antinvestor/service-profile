//
//  Generated code. Do not modify.
//  source: geolocation/v1/geolocation.proto
//

import "package:connectrpc/connect.dart" as connect;
import "geolocation.pb.dart" as geolocationv1geolocation;
import "../../google/protobuf/empty.pb.dart" as googleprotobufempty;

abstract final class GeolocationService {
  /// Fully-qualified name of the GeolocationService service.
  static const name = 'geolocation.v1.GeolocationService';

  static const ingestLocations = connect.Spec(
    '/$name/IngestLocations',
    connect.StreamType.unary,
    geolocationv1geolocation.IngestLocationsRequest.new,
    geolocationv1geolocation.IngestLocationsResponse.new,
  );

  static const createArea = connect.Spec(
    '/$name/CreateArea',
    connect.StreamType.unary,
    geolocationv1geolocation.CreateAreaRequest.new,
    geolocationv1geolocation.CreateAreaResponse.new,
  );

  static const getArea = connect.Spec(
    '/$name/GetArea',
    connect.StreamType.unary,
    geolocationv1geolocation.GetAreaRequest.new,
    geolocationv1geolocation.GetAreaResponse.new,
  );

  static const updateArea = connect.Spec(
    '/$name/UpdateArea',
    connect.StreamType.unary,
    geolocationv1geolocation.UpdateAreaRequest.new,
    geolocationv1geolocation.UpdateAreaResponse.new,
  );

  static const deleteArea = connect.Spec(
    '/$name/DeleteArea',
    connect.StreamType.unary,
    geolocationv1geolocation.DeleteAreaRequest.new,
    googleprotobufempty.Empty.new,
  );

  static const searchAreas = connect.Spec(
    '/$name/SearchAreas',
    connect.StreamType.unary,
    geolocationv1geolocation.SearchAreasRequest.new,
    geolocationv1geolocation.SearchAreasResponse.new,
  );

  static const createRoute = connect.Spec(
    '/$name/CreateRoute',
    connect.StreamType.unary,
    geolocationv1geolocation.CreateRouteRequest.new,
    geolocationv1geolocation.CreateRouteResponse.new,
  );

  static const getRoute = connect.Spec(
    '/$name/GetRoute',
    connect.StreamType.unary,
    geolocationv1geolocation.GetRouteRequest.new,
    geolocationv1geolocation.GetRouteResponse.new,
  );

  static const updateRoute = connect.Spec(
    '/$name/UpdateRoute',
    connect.StreamType.unary,
    geolocationv1geolocation.UpdateRouteRequest.new,
    geolocationv1geolocation.UpdateRouteResponse.new,
  );

  static const deleteRoute = connect.Spec(
    '/$name/DeleteRoute',
    connect.StreamType.unary,
    geolocationv1geolocation.DeleteRouteRequest.new,
    googleprotobufempty.Empty.new,
  );

  static const searchRoutes = connect.Spec(
    '/$name/SearchRoutes',
    connect.StreamType.unary,
    geolocationv1geolocation.SearchRoutesRequest.new,
    geolocationv1geolocation.SearchRoutesResponse.new,
  );

  static const assignRoute = connect.Spec(
    '/$name/AssignRoute',
    connect.StreamType.unary,
    geolocationv1geolocation.AssignRouteRequest.new,
    geolocationv1geolocation.AssignRouteResponse.new,
  );

  static const unassignRoute = connect.Spec(
    '/$name/UnassignRoute',
    connect.StreamType.unary,
    geolocationv1geolocation.UnassignRouteRequest.new,
    googleprotobufempty.Empty.new,
  );

  static const getSubjectRouteAssignments = connect.Spec(
    '/$name/GetSubjectRouteAssignments',
    connect.StreamType.unary,
    geolocationv1geolocation.GetSubjectRouteAssignmentsRequest.new,
    geolocationv1geolocation.GetSubjectRouteAssignmentsResponse.new,
  );

  static const getTrack = connect.Spec(
    '/$name/GetTrack',
    connect.StreamType.unary,
    geolocationv1geolocation.GetTrackRequest.new,
    geolocationv1geolocation.GetTrackResponse.new,
  );

  static const getSubjectEvents = connect.Spec(
    '/$name/GetSubjectEvents',
    connect.StreamType.unary,
    geolocationv1geolocation.GetSubjectEventsRequest.new,
    geolocationv1geolocation.GetSubjectEventsResponse.new,
  );

  static const getAreaSubjects = connect.Spec(
    '/$name/GetAreaSubjects',
    connect.StreamType.unary,
    geolocationv1geolocation.GetAreaSubjectsRequest.new,
    geolocationv1geolocation.GetAreaSubjectsResponse.new,
  );

  static const getNearbySubjects = connect.Spec(
    '/$name/GetNearbySubjects',
    connect.StreamType.unary,
    geolocationv1geolocation.GetNearbySubjectsRequest.new,
    geolocationv1geolocation.GetNearbySubjectsResponse.new,
  );

  static const getNearbyAreas = connect.Spec(
    '/$name/GetNearbyAreas',
    connect.StreamType.unary,
    geolocationv1geolocation.GetNearbyAreasRequest.new,
    geolocationv1geolocation.GetNearbyAreasResponse.new,
  );
}
