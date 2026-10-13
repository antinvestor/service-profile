package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect/v2"
	"github.com/pitabwire/frame/v3"
	"github.com/pitabwire/frame/v3/datastore"
	"github.com/pitabwire/frame/v3/security"
	"github.com/pitabwire/frame/v3/security/authorizer"
	"github.com/pitabwire/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"buf.build/gen/go/antinvestor/geolocation/connectrpc/go/v2/geolocation/v1/geolocationv1connect"
	geolocationv1 "buf.build/gen/go/antinvestor/geolocation/protocolbuffers/go/geolocation/v1"

	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
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
	req *geolocationv1.IngestLocationsRequest,
) (*geolocationv1.IngestLocationsResponse, error) {
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != req.GetSubjectId() {
		if err := s.checker.Check(ctx, "location_ingest"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	resp, err := s.ingestionBiz.IngestBatch(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return resp, nil
}

func (s *GeolocationServer) CreateArea(
	ctx context.Context,
	req *geolocationv1.CreateAreaRequest,
) (*geolocationv1.CreateAreaResponse, error) {
	if err := s.checker.Check(ctx, "geolocation_manage"); err != nil {
		return nil, authorizer.ToConnectError(err)
	}

	area, err := s.areaBiz.CreateArea(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.CreateAreaResponse{Data: area}, nil
}

func (s *GeolocationServer) GetArea(
	ctx context.Context,
	req *geolocationv1.GetAreaRequest,
) (*geolocationv1.GetAreaResponse, error) {
	area, err := s.areaBiz.GetArea(ctx, req.GetId())
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetAreaResponse{Data: area}, nil
}

func (s *GeolocationServer) UpdateArea(
	ctx context.Context,
	req *geolocationv1.UpdateAreaRequest,
) (*geolocationv1.UpdateAreaResponse, error) {
	area, err := s.areaBiz.UpdateArea(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.UpdateAreaResponse{Data: area}, nil
}

func (s *GeolocationServer) DeleteArea(
	ctx context.Context,
	req *geolocationv1.DeleteAreaRequest,
) (*emptypb.Empty, error) {
	if err := s.areaBiz.DeleteArea(ctx, req.GetId()); err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *GeolocationServer) SearchAreas(
	ctx context.Context,
	req *geolocationv1.SearchAreasRequest,
) (*geolocationv1.SearchAreasResponse, error) {
	areas, err := s.areaBiz.SearchAreas(ctx, req.GetQuery(), req.GetOwnerId(), int(req.GetLimit()))
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.SearchAreasResponse{Data: areas}, nil
}

func (s *GeolocationServer) CreateRoute(
	ctx context.Context,
	req *geolocationv1.CreateRouteRequest,
) (*geolocationv1.CreateRouteResponse, error) {
	route, err := s.routeBiz.CreateRoute(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.CreateRouteResponse{Data: route}, nil
}

func (s *GeolocationServer) GetRoute(
	ctx context.Context,
	req *geolocationv1.GetRouteRequest,
) (*geolocationv1.GetRouteResponse, error) {
	route, err := s.routeBiz.GetRoute(ctx, req.GetId())
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetRouteResponse{Data: route}, nil
}

func (s *GeolocationServer) UpdateRoute(
	ctx context.Context,
	req *geolocationv1.UpdateRouteRequest,
) (*geolocationv1.UpdateRouteResponse, error) {
	route, err := s.routeBiz.UpdateRoute(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.UpdateRouteResponse{Data: route}, nil
}

func (s *GeolocationServer) DeleteRoute(
	ctx context.Context,
	req *geolocationv1.DeleteRouteRequest,
) (*emptypb.Empty, error) {
	if err := s.routeBiz.DeleteRoute(ctx, req.GetId()); err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *GeolocationServer) SearchRoutes(
	ctx context.Context,
	req *geolocationv1.SearchRoutesRequest,
) (*geolocationv1.SearchRoutesResponse, error) {
	routes, err := s.routeBiz.SearchRoutes(ctx, req.GetOwnerId(), int(req.GetLimit()))
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.SearchRoutesResponse{Data: routes}, nil
}

func (s *GeolocationServer) AssignRoute(
	ctx context.Context,
	req *geolocationv1.AssignRouteRequest,
) (*geolocationv1.AssignRouteResponse, error) {
	assignment, err := s.routeBiz.AssignRoute(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.AssignRouteResponse{Data: assignment}, nil
}

func (s *GeolocationServer) UnassignRoute(
	ctx context.Context,
	req *geolocationv1.UnassignRouteRequest,
) (*emptypb.Empty, error) {
	if err := s.routeBiz.UnassignRoute(ctx, req.GetId()); err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *GeolocationServer) GetSubjectRouteAssignments(
	ctx context.Context,
	req *geolocationv1.GetSubjectRouteAssignmentsRequest,
) (*geolocationv1.GetSubjectRouteAssignmentsResponse, error) {
	subjectID := req.GetSubjectId()
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
	return &geolocationv1.GetSubjectRouteAssignmentsResponse{Data: assignments}, nil
}

func (s *GeolocationServer) GetTrack(
	ctx context.Context,
	req *geolocationv1.GetTrackRequest,
) (*geolocationv1.GetTrackResponse, error) {
	subjectID := req.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	points, err := s.trackBiz.GetTrack(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetTrackResponse{Data: points}, nil
}

func (s *GeolocationServer) GetSubjectEvents(
	ctx context.Context,
	req *geolocationv1.GetSubjectEventsRequest,
) (*geolocationv1.GetSubjectEventsResponse, error) {
	subjectID := req.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	events, err := s.trackBiz.GetSubjectEvents(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetSubjectEventsResponse{Data: events}, nil
}

func (s *GeolocationServer) GetAreaSubjects(
	ctx context.Context,
	req *geolocationv1.GetAreaSubjectsRequest,
) (*geolocationv1.GetAreaSubjectsResponse, error) {
	subjects, err := s.trackBiz.GetAreaSubjects(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetAreaSubjectsResponse{Data: subjects}, nil
}

func (s *GeolocationServer) GetNearbySubjects(
	ctx context.Context,
	req *geolocationv1.GetNearbySubjectsRequest,
) (*geolocationv1.GetNearbySubjectsResponse, error) {
	subjectID := req.GetSubjectId()
	claims := security.ClaimsFromContext(ctx)
	if sub, _ := claims.GetSubject(); sub != subjectID {
		if err := s.checker.Check(ctx, "geolocation_view"); err != nil {
			return nil, authorizer.ToConnectError(err)
		}
	}

	subjects, err := s.proximityBiz.GetNearbySubjects(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetNearbySubjectsResponse{Data: subjects}, nil
}

func (s *GeolocationServer) GetNearbyAreas(
	ctx context.Context,
	req *geolocationv1.GetNearbyAreasRequest,
) (*geolocationv1.GetNearbyAreasResponse, error) {
	areas, err := s.proximityBiz.GetNearbyAreas(ctx, req)
	if err != nil {
		return nil, s.cleanErr(ctx, err)
	}
	return &geolocationv1.GetNearbyAreasResponse{Data: areas}, nil
}

func (s *GeolocationServer) cleanErr(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), strings.Contains(err.Error(), "not found"):
		return connect.NewError(connect.CodeNotFound, "resource not found")
	case isValidationError(err):
		return connect.NewError(connect.CodeInvalidArgument, err.Error()).WithCause(err)
	default:
		util.Log(ctx).WithError(err).Error("internal error processing geolocation request")
		return connect.NewError(connect.CodeInternal, "internal server error")
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
