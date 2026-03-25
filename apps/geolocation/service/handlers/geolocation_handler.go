package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
	"github.com/pitabwire/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	geolocationv1 "github.com/antinvestor/service-profile/proto/geolocation/v1"
	"github.com/antinvestor/service-profile/proto/geolocation/v1/geolocationv1connect"
)

type GeolocationServer struct {
	geolocationv1connect.UnimplementedGeolocationServiceHandler

	Service      *frame.Service
	checker      *authorizer.FunctionChecker
	ingestionBiz business.IngestionBusiness
	areaBiz      business.AreaBusiness
	routeBiz     business.RouteBusiness
	proximityBiz business.ProximityBusiness
	trackBiz     business.TrackBusiness
	metrics      *observability.Metrics
}

func NewGeolocationServer(
	svc *frame.Service,
	checker *authorizer.FunctionChecker,
	ingestionBiz business.IngestionBusiness,
	areaBiz business.AreaBusiness,
	routeBiz business.RouteBusiness,
	proximityBiz business.ProximityBusiness,
	trackBiz business.TrackBusiness,
	metrics *observability.Metrics,
	_ int64,
) *GeolocationServer {
	return &GeolocationServer{
		Service:      svc,
		checker:      checker,
		ingestionBiz: ingestionBiz,
		areaBiz:      areaBiz,
		routeBiz:     routeBiz,
		proximityBiz: proximityBiz,
		trackBiz:     trackBiz,
		metrics:      metrics,
	}
}

func (s *GeolocationServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	db := s.Service.DatastoreManager().GetPool(r.Context(), datastore.DefaultPoolName).DB(r.Context(), true)
	if err := db.Exec("SELECT 1").Error; err != nil {
		http.Error(w, "unhealthy", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *GeolocationServer) IngestLocations(
	ctx context.Context,
	req *connect.Request[geolocationv1.IngestLocationsRequest],
) (*connect.Response[geolocationv1.IngestLocationsResponse], error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != req.Msg.GetSubjectId() {
		if err := s.checker.Check(ctx, "location_ingest"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	resp, err := s.ingestionBiz.IngestBatch(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(resp), nil
}

func (s *GeolocationServer) CreateArea(
	ctx context.Context,
	req *connect.Request[geolocationv1.CreateAreaRequest],
) (*connect.Response[geolocationv1.CreateAreaResponse], error) {
	area, err := s.areaBiz.CreateArea(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.CreateAreaResponse{Data: area}), nil
}

func (s *GeolocationServer) GetArea(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetAreaRequest],
) (*connect.Response[geolocationv1.GetAreaResponse], error) {
	area, err := s.areaBiz.GetArea(ctx, req.Msg.GetId())
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetAreaResponse{Data: area}), nil
}

func (s *GeolocationServer) UpdateArea(
	ctx context.Context,
	req *connect.Request[geolocationv1.UpdateAreaRequest],
) (*connect.Response[geolocationv1.UpdateAreaResponse], error) {
	area, err := s.areaBiz.UpdateArea(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.UpdateAreaResponse{Data: area}), nil
}

func (s *GeolocationServer) DeleteArea(
	ctx context.Context,
	req *connect.Request[geolocationv1.DeleteAreaRequest],
) (*connect.Response[emptypb.Empty], error) {
	if err := s.areaBiz.DeleteArea(ctx, req.Msg.GetId()); err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *GeolocationServer) SearchAreas(
	ctx context.Context,
	req *connect.Request[geolocationv1.SearchAreasRequest],
) (*connect.Response[geolocationv1.SearchAreasResponse], error) {
	areas, err := s.areaBiz.SearchAreas(ctx, req.Msg.GetQuery(), req.Msg.GetOwnerId(), int(req.Msg.GetLimit()))
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.SearchAreasResponse{Data: areas}), nil
}

func (s *GeolocationServer) CreateRoute(
	ctx context.Context,
	req *connect.Request[geolocationv1.CreateRouteRequest],
) (*connect.Response[geolocationv1.CreateRouteResponse], error) {
	route, err := s.routeBiz.CreateRoute(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.CreateRouteResponse{Data: route}), nil
}

func (s *GeolocationServer) GetRoute(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetRouteRequest],
) (*connect.Response[geolocationv1.GetRouteResponse], error) {
	route, err := s.routeBiz.GetRoute(ctx, req.Msg.GetId())
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetRouteResponse{Data: route}), nil
}

func (s *GeolocationServer) UpdateRoute(
	ctx context.Context,
	req *connect.Request[geolocationv1.UpdateRouteRequest],
) (*connect.Response[geolocationv1.UpdateRouteResponse], error) {
	route, err := s.routeBiz.UpdateRoute(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.UpdateRouteResponse{Data: route}), nil
}

func (s *GeolocationServer) DeleteRoute(
	ctx context.Context,
	req *connect.Request[geolocationv1.DeleteRouteRequest],
) (*connect.Response[emptypb.Empty], error) {
	if err := s.routeBiz.DeleteRoute(ctx, req.Msg.GetId()); err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *GeolocationServer) SearchRoutes(
	ctx context.Context,
	req *connect.Request[geolocationv1.SearchRoutesRequest],
) (*connect.Response[geolocationv1.SearchRoutesResponse], error) {
	routes, err := s.routeBiz.SearchRoutes(ctx, req.Msg.GetOwnerId(), int(req.Msg.GetLimit()))
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.SearchRoutesResponse{Data: routes}), nil
}

func (s *GeolocationServer) AssignRoute(
	ctx context.Context,
	req *connect.Request[geolocationv1.AssignRouteRequest],
) (*connect.Response[geolocationv1.AssignRouteResponse], error) {
	assignment, err := s.routeBiz.AssignRoute(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.AssignRouteResponse{Data: assignment}), nil
}

func (s *GeolocationServer) UnassignRoute(
	ctx context.Context,
	req *connect.Request[geolocationv1.UnassignRouteRequest],
) (*connect.Response[emptypb.Empty], error) {
	if err := s.routeBiz.UnassignRoute(ctx, req.Msg.GetId()); err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *GeolocationServer) GetSubjectRouteAssignments(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetSubjectRouteAssignmentsRequest],
) (*connect.Response[geolocationv1.GetSubjectRouteAssignmentsResponse], error) {
	subjectID := req.Msg.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	assignments, err := s.routeBiz.GetSubjectAssignments(ctx, subjectID)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetSubjectRouteAssignmentsResponse{Data: assignments}), nil
}

func (s *GeolocationServer) GetTrack(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetTrackRequest],
) (*connect.Response[geolocationv1.GetTrackResponse], error) {
	subjectID := req.Msg.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	points, err := s.trackBiz.GetTrack(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetTrackResponse{Data: points}), nil
}

func (s *GeolocationServer) GetSubjectEvents(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetSubjectEventsRequest],
) (*connect.Response[geolocationv1.GetSubjectEventsResponse], error) {
	subjectID := req.Msg.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	events, err := s.trackBiz.GetSubjectEvents(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetSubjectEventsResponse{Data: events}), nil
}

func (s *GeolocationServer) GetAreaSubjects(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetAreaSubjectsRequest],
) (*connect.Response[geolocationv1.GetAreaSubjectsResponse], error) {
	subjects, err := s.trackBiz.GetAreaSubjects(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetAreaSubjectsResponse{Data: subjects}), nil
}

func (s *GeolocationServer) GetNearbySubjects(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetNearbySubjectsRequest],
) (*connect.Response[geolocationv1.GetNearbySubjectsResponse], error) {
	subjectID := req.Msg.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	subjects, err := s.proximityBiz.GetNearbySubjects(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetNearbySubjectsResponse{Data: subjects}), nil
}

func (s *GeolocationServer) GetNearbyAreas(
	ctx context.Context,
	req *connect.Request[geolocationv1.GetNearbyAreasRequest],
) (*connect.Response[geolocationv1.GetNearbyAreasResponse], error) {
	areas, err := s.proximityBiz.GetNearbyAreas(ctx, req.Msg)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return connect.NewResponse(&geolocationv1.GetNearbyAreasResponse{Data: areas}), nil
}

func (s *GeolocationServer) cleanErr(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), strings.Contains(err.Error(), "not found"):
		return connect.NewError(connect.CodeNotFound, errors.New("resource not found"))
	case isValidationError(err):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		util.Log(ctx).WithError(err).Error("internal error processing geolocation request")
		return connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	for _, prefix := range []string{
		"invalid",
		"required",
		"must be",
		"batch size",
		"either query or owner_id",
		"radius_meters",
		"subject_id",
		"owner_id",
		"route_id",
		"area_id",
	} {
		if strings.HasPrefix(msg, prefix) {
			return true
		}
	}

	return false
}

var _ geolocationv1connect.GeolocationServiceHandler = (*GeolocationServer)(nil)
