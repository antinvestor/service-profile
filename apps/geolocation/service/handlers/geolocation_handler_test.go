package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/security/authorizer"
	"github.com/pitabwire/frame/workerpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	aconfig "github.com/antinvestor/service-profile/apps/geolocation/config"
	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	geoevents "github.com/antinvestor/service-profile/apps/geolocation/service/events"
	"github.com/antinvestor/service-profile/apps/geolocation/service/handlers"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
	geotests "github.com/antinvestor/service-profile/apps/geolocation/tests"
	geolocationv1 "github.com/antinvestor/service-profile/proto/geolocation/geolocation/v1"
)

const (
	handlerTenantID    = "tenant-handlers"
	handlerPartitionID = "partition-handlers"
)

//nolint:gochecknoglobals // test fixture data
var handlerAreaGeometry = `{"type":"Polygon","coordinates":[[[32.58,0.34],[32.60,0.34],[32.60,0.36],[32.58,0.36],[32.58,0.34]]]}`

type HandlerSuite struct {
	geotests.GeolocationBaseTestSuite
}

type handlerStack struct {
	GeoServer *handlers.GeolocationServer
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) scopedContext(ctx context.Context, subjectID string) context.Context {
	return s.WithAuthClaims(ctx, handlerTenantID, handlerPartitionID, subjectID)
}

func newHandlerStack(ctx context.Context, svc *frame.Service) *handlerStack {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	var workMan workerpool.Manager
	cfg := svc.Config().(*aconfig.GeolocationConfig)
	evtsMan := svc.EventsManager()
	metrics := observability.NewMetrics()

	pointRepo := repository.NewLocationPointRepository(ctx, dbPool, workMan)
	areaRepo := repository.NewAreaRepository(ctx, dbPool, workMan)
	geoEventRepo := repository.NewGeoEventRepository(ctx, dbPool, workMan)
	stateRepo := repository.NewGeofenceStateRepository(ctx, dbPool, workMan)
	latestPosRepo := repository.NewLatestPositionRepository(ctx, dbPool, workMan)
	routeRepo := repository.NewRouteRepository(ctx, dbPool, workMan)
	routeAssignmentRepo := repository.NewRouteAssignmentRepository(ctx, dbPool, workMan)
	routeDeviationStateRepo := repository.NewRouteDeviationStateRepository(ctx, dbPool, workMan)
	ingestionBiz := business.NewIngestionBusiness(evtsMan, pointRepo, cfg.IngestionBusinessConfig())
	areaBiz := business.NewAreaBusiness(evtsMan, areaRepo, stateRepo)
	routeBiz := business.NewRouteBusiness(evtsMan, routeRepo, routeAssignmentRepo, routeDeviationStateRepo)
	proximityBiz := business.NewProximityBusiness(latestPosRepo, areaRepo, cfg.ProximityBusinessConfig())
	trackBiz := business.NewTrackBusiness(pointRepo, geoEventRepo, stateRepo, cfg.TrackBusinessConfig())
	geofenceBiz := business.NewGeofenceBusiness(
		evtsMan,
		areaRepo,
		stateRepo,
		geoEventRepo,
		metrics,
		cfg.GeofenceBusinessConfig(),
	)
	routeDeviationBiz := business.NewRouteDeviationBusiness(
		evtsMan,
		routeRepo,
		routeDeviationStateRepo,
		metrics,
		cfg.RouteDeviationBusinessConfig(),
	)

	evtsMan.Add(geoevents.NewLocationPointConsumer(
		pointRepo,
		proximityBiz,
		geofenceBiz,
		routeDeviationBiz,
		metrics,
	))
	evtsMan.Add(geoevents.NewAreaChangeConsumer(areaBiz, stateRepo))
	evtsMan.Add(geoevents.NewRouteChangeConsumer(routeAssignmentRepo, routeDeviationStateRepo))
	evtsMan.Add(geoevents.NewGeoEventConsumer())
	evtsMan.Add(geoevents.NewRouteDeviationConsumer())

	return &handlerStack{
		GeoServer: handlers.NewGeolocationServer(
			svc,
			authorizer.NewFunctionChecker(svc.SecurityManager().GetAuthorizer(ctx), "service_profile"),
			ingestionBiz,
			areaBiz,
			routeBiz,
			proximityBiz,
			trackBiz,
			metrics,
			cfg.MaxRequestBodyBytes,
		),
	}
}

func (s *HandlerSuite) TestGeolocationServerFlow() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newHandlerStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "profile-1")
		s.SeedTenantAccess(ctx, svc, handlerTenantID, handlerPartitionID, "profile-1")
		s.SeedTenantRole(ctx, svc, handlerTenantID, handlerPartitionID, "profile-1", "owner")

		server := stack.GeoServer

		ingestResp, err := server.IngestLocations(ctx, connect.NewRequest(&geolocationv1.IngestLocationsRequest{
			SubjectId: "profile-1",
			Points: []*geolocationv1.LocationPointInput{{
				DeviceId: "device-profile-1",
				Latitude: 0.3500, Longitude: 32.5900, Accuracy: 4, Timestamp: timestamppb.New(time.Now().UTC()),
			}},
		}))
		require.NoError(t, err)
		require.EqualValues(t, 1, ingestResp.Msg.GetAccepted())

		areaResp, err := server.CreateArea(ctx, connect.NewRequest(&geolocationv1.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{OwnerId: "owner-1", Name: "Handler Area", Geometry: handlerAreaGeometry},
		}))
		require.NoError(t, err)
		areaID := areaResp.Msg.GetData().GetId()

		_, err = server.GetArea(ctx, connect.NewRequest(&geolocationv1.GetAreaRequest{Id: areaID}))
		require.NoError(t, err)

		updatedArea, err := server.UpdateArea(
			ctx,
			connect.NewRequest(&geolocationv1.UpdateAreaRequest{Id: areaID, Name: ptr("Handler Area Updated")}),
		)
		require.NoError(t, err)
		require.Equal(t, "Handler Area Updated", updatedArea.Msg.GetData().GetName())

		searchedAreas, err := server.SearchAreas(
			ctx,
			connect.NewRequest(&geolocationv1.SearchAreasRequest{OwnerId: "owner-1", Limit: 10}),
		)
		require.NoError(t, err)
		require.Len(t, searchedAreas.Msg.GetData(), 1)

		routeResp, err := server.CreateRoute(ctx, connect.NewRequest(&geolocationv1.CreateRouteRequest{
			Data: &geolocationv1.RouteObject{
				OwnerId: "owner-1", Name: "Handler Route",
				Geometry:            `{"type":"LineString","coordinates":[[32.58,0.35],[32.60,0.35]]}`,
				DeviationThresholdM: float64Ptr(20),
			},
		}))
		require.NoError(t, err)
		routeID := routeResp.Msg.GetData().GetId()

		_, err = server.GetRoute(ctx, connect.NewRequest(&geolocationv1.GetRouteRequest{Id: routeID}))
		require.NoError(t, err)

		updatedRoute, err := server.UpdateRoute(
			ctx,
			connect.NewRequest(&geolocationv1.UpdateRouteRequest{Id: routeID, Name: ptr("Handler Route Updated")}),
		)
		require.NoError(t, err)
		require.Equal(t, "Handler Route Updated", updatedRoute.Msg.GetData().GetName())

		searchedRoutes, err := server.SearchRoutes(
			ctx,
			connect.NewRequest(&geolocationv1.SearchRoutesRequest{OwnerId: "owner-1", Limit: 10}),
		)
		require.NoError(t, err)
		require.Len(t, searchedRoutes.Msg.GetData(), 1)

		assignResp, err := server.AssignRoute(
			ctx,
			connect.NewRequest(&geolocationv1.AssignRouteRequest{SubjectId: "profile-1", RouteId: routeID}),
		)
		require.NoError(t, err)
		assignmentID := assignResp.Msg.GetData().GetId()

		assignments, err := server.GetSubjectRouteAssignments(
			ctx,
			connect.NewRequest(&geolocationv1.GetSubjectRouteAssignmentsRequest{SubjectId: "profile-1"}),
		)
		require.NoError(t, err)
		require.Len(t, assignments.Msg.GetData(), 1)

		trackResp, err := server.GetTrack(
			ctx,
			connect.NewRequest(&geolocationv1.GetTrackRequest{SubjectId: "profile-1", Limit: 10}),
		)
		require.NoError(t, err)
		require.Len(t, trackResp.Msg.GetData(), 1)

		nearbyAreas, err := server.GetNearbyAreas(
			ctx,
			connect.NewRequest(
				&geolocationv1.GetNearbyAreasRequest{Latitude: 0.35, Longitude: 32.59, RadiusMeters: 1000},
			),
		)
		require.NoError(t, err)
		require.Len(t, nearbyAreas.Msg.GetData(), 1)

		eventsResp, err := server.GetSubjectEvents(
			ctx,
			connect.NewRequest(&geolocationv1.GetSubjectEventsRequest{SubjectId: "profile-1", Limit: 10}),
		)
		require.NoError(t, err)
		require.NotNil(t, eventsResp.Msg)

		subjectsResp, err := server.GetAreaSubjects(
			ctx,
			connect.NewRequest(&geolocationv1.GetAreaSubjectsRequest{AreaId: areaID}),
		)
		require.NoError(t, err)
		require.NotNil(t, subjectsResp.Msg)

		var nearbySubjects *connect.Response[geolocationv1.GetNearbySubjectsResponse]
		require.Eventually(t, func() bool {
			var nearErr error
			nearbySubjects, nearErr = server.GetNearbySubjects(
				ctx,
				connect.NewRequest(&geolocationv1.GetNearbySubjectsRequest{
					SubjectId:    "profile-1",
					RadiusMeters: 1000,
					Limit:        10,
				}),
			)
			return nearErr == nil
		}, 5*time.Second, 100*time.Millisecond)
		require.Empty(t, nearbySubjects.Msg.GetData())

		unassignResp, err := server.UnassignRoute(
			ctx,
			connect.NewRequest(&geolocationv1.UnassignRouteRequest{Id: assignmentID}),
		)
		require.NoError(t, err)
		require.IsType(t, &emptypb.Empty{}, unassignResp.Msg)

		deleteRouteResp, err := server.DeleteRoute(
			ctx,
			connect.NewRequest(&geolocationv1.DeleteRouteRequest{Id: routeID}),
		)
		require.NoError(t, err)
		require.IsType(t, &emptypb.Empty{}, deleteRouteResp.Msg)

		deleteAreaResp, err := server.DeleteArea(ctx, connect.NewRequest(&geolocationv1.DeleteAreaRequest{Id: areaID}))
		require.NoError(t, err)
		require.IsType(t, &emptypb.Empty{}, deleteAreaResp.Msg)
	})
}

func (s *HandlerSuite) TestGeolocationServerAuthorization() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newHandlerStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "profile-2")
		s.SeedTenantAccess(ctx, svc, handlerTenantID, handlerPartitionID, "profile-2")

		server := stack.GeoServer

		resp, err := server.IngestLocations(ctx, connect.NewRequest(&geolocationv1.IngestLocationsRequest{
			SubjectId: "profile-2",
			Points: []*geolocationv1.LocationPointInput{
				{DeviceId: "device-profile-2", Latitude: 0.35, Longitude: 32.59, Accuracy: 4},
			},
		}))
		require.NoError(t, err)
		require.EqualValues(t, 1, resp.Msg.GetAccepted())

		_, err = server.CreateArea(ctx, connect.NewRequest(&geolocationv1.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{OwnerId: "owner", Name: "x", Geometry: handlerAreaGeometry},
		}))
		require.Error(t, err)
		require.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))

		_, err = server.GetNearbySubjects(
			ctx,
			connect.NewRequest(&geolocationv1.GetNearbySubjectsRequest{SubjectId: "someone-else"}),
		)
		require.Error(t, err)
		require.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})
}

func (s *HandlerSuite) TestGeolocationServerErrorPaths() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newHandlerStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "profile-errors")
		s.SeedTenantAccess(ctx, svc, handlerTenantID, handlerPartitionID, "profile-errors")
		s.SeedTenantRole(ctx, svc, handlerTenantID, handlerPartitionID, "profile-errors", "owner")

		server := stack.GeoServer

		cases := []struct {
			name string
			run  func() error
			code connect.Code
		}{
			{
				name: "create invalid area",
				run: func() error {
					_, err := server.CreateArea(ctx, connect.NewRequest(&geolocationv1.CreateAreaRequest{
						Data: &geolocationv1.AreaObject{
							OwnerId:  "owner",
							Name:     "x",
							Geometry: `{"type":"Point","coordinates":[32.58,0.35]}`,
						},
					}))
					return err
				},
				code: connect.CodeInvalidArgument,
			},
			{
				name: "get missing area",
				run: func() error {
					_, err := server.GetArea(ctx, connect.NewRequest(&geolocationv1.GetAreaRequest{Id: "missing"}))
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "update missing area",
				run: func() error {
					_, err := server.UpdateArea(
						ctx,
						connect.NewRequest(&geolocationv1.UpdateAreaRequest{Id: "missing", Name: ptr("x")}),
					)
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "delete missing area",
				run: func() error {
					_, err := server.DeleteArea(
						ctx,
						connect.NewRequest(&geolocationv1.DeleteAreaRequest{Id: "missing"}),
					)
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "search areas invalid",
				run: func() error {
					_, err := server.SearchAreas(ctx, connect.NewRequest(&geolocationv1.SearchAreasRequest{}))
					return err
				},
				code: connect.CodeInvalidArgument,
			},
			{
				name: "create invalid route",
				run: func() error {
					_, err := server.CreateRoute(ctx, connect.NewRequest(&geolocationv1.CreateRouteRequest{
						Data: &geolocationv1.RouteObject{
							OwnerId:  "owner",
							Name:     "x",
							Geometry: `{"type":"Polygon","coordinates":[]}`,
						},
					}))
					return err
				},
				code: connect.CodeInvalidArgument,
			},
			{
				name: "get missing route",
				run: func() error {
					_, err := server.GetRoute(ctx, connect.NewRequest(&geolocationv1.GetRouteRequest{Id: "missing"}))
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "update missing route",
				run: func() error {
					_, err := server.UpdateRoute(
						ctx,
						connect.NewRequest(&geolocationv1.UpdateRouteRequest{Id: "missing", Name: ptr("x")}),
					)
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "delete missing route",
				run: func() error {
					_, err := server.DeleteRoute(
						ctx,
						connect.NewRequest(&geolocationv1.DeleteRouteRequest{Id: "missing"}),
					)
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "search routes invalid",
				run: func() error {
					_, err := server.SearchRoutes(ctx, connect.NewRequest(&geolocationv1.SearchRoutesRequest{}))
					return err
				},
				code: connect.CodeInvalidArgument,
			},
			{
				name: "assign route missing ids",
				run: func() error {
					_, err := server.AssignRoute(ctx, connect.NewRequest(&geolocationv1.AssignRouteRequest{}))
					return err
				},
				code: connect.CodeInvalidArgument,
			},
			{
				name: "unassign missing assignment",
				run: func() error {
					_, err := server.UnassignRoute(
						ctx,
						connect.NewRequest(&geolocationv1.UnassignRouteRequest{Id: "missing"}),
					)
					return err
				},
				code: connect.CodeNotFound,
			},
			{
				name: "nearby areas invalid radius",
				run: func() error {
					_, err := server.GetNearbyAreas(ctx, connect.NewRequest(&geolocationv1.GetNearbyAreasRequest{
						Latitude: 0.35, Longitude: 32.59, RadiusMeters: models.MaxProximityRadiusMeters + 1,
					}))
					return err
				},
				code: connect.CodeInvalidArgument,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.run()
				require.Error(t, err)
				require.Equal(t, tc.code, connect.CodeOf(err))
			})
		}
	})
}

func (s *HandlerSuite) TestHealthCheck() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newHandlerStack(baseCtx, svc)

		rec := httptest.NewRecorder()
		stack.GeoServer.HealthCheck(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		require.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Parallel()

	result := handlers.NewRateLimitMiddleware(&handlers.RateLimiterConfig{
		RequestsPerWindow: 1,
		WindowDuration:    time.Hour,
		CleanupInterval:   time.Millisecond,
	})
	defer result.Stop()

	calls := 0
	handler := result.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.1:1234"

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, req)
	require.Equal(t, http.StatusOK, first.Code)

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, req)
	require.Equal(t, http.StatusTooManyRequests, second.Code)
	require.Equal(t, "60", second.Header().Get("Retry-After"))
	require.Equal(t, 1, calls)
}

func TestRateLimiterHelpersAndOpenAPI(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	handlers.OpenAPIHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))
	require.NotEmpty(t, rec.Body.Bytes())
}

func ptr[T any](value T) *T             { return &value }
func float64Ptr(value float64) *float64 { return &value }
