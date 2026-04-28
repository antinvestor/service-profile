package business_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/workerpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	geobusiness "github.com/antinvestor/service-profile/apps/geolocation/service/business"
	geoevents "github.com/antinvestor/service-profile/apps/geolocation/service/events"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
	geotests "github.com/antinvestor/service-profile/apps/geolocation/tests"
	geolocationv1 "github.com/antinvestor/service-profile/geolocation/v1"
)

const (
	businessTenantID    = "tenant-business"
	businessPartitionID = "partition-business"
)

var (
	//nolint:gochecknoglobals // test fixture data
	businessAreaGeometry = `{"type":"Polygon","coordinates":[[[32.58,0.34],[32.60,0.34],[32.60,0.36],[32.58,0.36],[32.58,0.34]]]}`
	//nolint:gochecknoglobals // test fixture data
	businessRouteGeometry = `{"type":"LineString","coordinates":[[32.58,0.35],[32.60,0.35]]}`
)

type BusinessSuite struct {
	geotests.GeolocationBaseTestSuite
}

type businessStack struct {
	PointRepo               repository.LocationPointRepository
	AreaRepo                repository.AreaRepository
	GeoEventRepo            repository.GeoEventRepository
	GeofenceStateRepo       repository.GeofenceStateRepository
	LatestPositionRepo      repository.LatestPositionRepository
	RouteRepo               repository.RouteRepository
	RouteAssignmentRepo     repository.RouteAssignmentRepository
	RouteDeviationStateRepo repository.RouteDeviationStateRepository
	IngestionBiz            geobusiness.IngestionBusiness
	AreaBiz                 geobusiness.AreaBusiness
	GeofenceBiz             geobusiness.GeofenceBusiness
	ProximityBiz            geobusiness.ProximityBusiness
	TrackBiz                geobusiness.TrackBusiness
	RouteDeviationBiz       geobusiness.RouteDeviationBusiness
	RouteBiz                geobusiness.RouteBusiness
	CatchupBiz              geobusiness.CatchupBusiness
	LocationPointConsumer   *geoevents.LocationPointConsumer
}

func TestBusinessSuite(t *testing.T) {
	suite.Run(t, new(BusinessSuite))
}

func (s *BusinessSuite) scopedContext(ctx context.Context, subjectID string) context.Context {
	return s.WithAuthClaims(ctx, businessTenantID, businessPartitionID, subjectID)
}

func newBusinessStack(ctx context.Context, svc *frame.Service) *businessStack {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	var workMan workerpool.Manager
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

	ingestionBiz := geobusiness.NewIngestionBusiness(evtsMan, pointRepo, geobusiness.IngestionConfig{})
	areaBiz := geobusiness.NewAreaBusiness(evtsMan, areaRepo, stateRepo)
	geofenceBiz := geobusiness.NewGeofenceBusiness(
		evtsMan,
		areaRepo,
		stateRepo,
		geoEventRepo,
		metrics,
		geobusiness.GeofenceConfig{},
	)
	proximityBiz := geobusiness.NewProximityBusiness(latestPosRepo, areaRepo, geobusiness.ProximityConfig{})
	trackBiz := geobusiness.NewTrackBusiness(pointRepo, geoEventRepo, stateRepo, geobusiness.TrackConfig{})
	routeDeviationBiz := geobusiness.NewRouteDeviationBusiness(
		evtsMan,
		routeRepo,
		routeDeviationStateRepo,
		metrics,
		geobusiness.RouteDeviationConfig{},
	)
	routeBiz := geobusiness.NewRouteBusiness(evtsMan, routeRepo, routeAssignmentRepo, routeDeviationStateRepo)
	catchupBiz := geobusiness.NewCatchupBusiness(pointRepo, evtsMan, geobusiness.CatchupConfig{})

	evtsMan.Add(geoevents.NewAreaChangeConsumer(areaBiz, stateRepo))
	evtsMan.Add(geoevents.NewRouteChangeConsumer(routeAssignmentRepo, routeDeviationStateRepo))
	evtsMan.Add(geoevents.NewGeoEventConsumer())
	evtsMan.Add(geoevents.NewRouteDeviationConsumer())

	return &businessStack{
		PointRepo:               pointRepo,
		AreaRepo:                areaRepo,
		GeoEventRepo:            geoEventRepo,
		GeofenceStateRepo:       stateRepo,
		LatestPositionRepo:      latestPosRepo,
		RouteRepo:               routeRepo,
		RouteAssignmentRepo:     routeAssignmentRepo,
		RouteDeviationStateRepo: routeDeviationStateRepo,
		IngestionBiz:            ingestionBiz,
		AreaBiz:                 areaBiz,
		GeofenceBiz:             geofenceBiz,
		ProximityBiz:            proximityBiz,
		TrackBiz:                trackBiz,
		RouteDeviationBiz:       routeDeviationBiz,
		RouteBiz:                routeBiz,
		CatchupBiz:              catchupBiz,
		LocationPointConsumer: geoevents.NewLocationPointConsumer(
			pointRepo,
			proximityBiz,
			geofenceBiz,
			routeDeviationBiz,
			metrics,
		),
	}
}

func (s *BusinessSuite) makePointEvent(
	point *models.LocationPoint,
	accessID string,
) *models.LocationPointIngestedEvent {
	return &models.LocationPointIngestedEvent{
		EventTenancy: models.EventTenancy{
			TenantID:    businessTenantID,
			PartitionID: businessPartitionID,
			AccessID:    accessID,
		},
		PointID:   point.GetID(),
		SubjectID: point.SubjectID,
		DeviceID:  point.DeviceID,
		Latitude:  point.Latitude,
		Longitude: point.Longitude,
		Accuracy:  point.Accuracy,
		Timestamp: point.TrueCreatedAt.UnixMilli(),
	}
}

func (s *BusinessSuite) TestAreaAndRouteLifecycle() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newBusinessStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "owner-1")

		area, err := stack.AreaBiz.CreateArea(ctx, &models.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{
				OwnerId:     "owner-1",
				Name:        "Campus Yard",
				Description: "Main yard",
				AreaType:    geolocationv1.AreaType_AREA_TYPE_LAND,
				Geometry:    businessAreaGeometry,
			},
		})
		require.NoError(t, err)
		require.Greater(t, area.GetAreaM2(), 0.0)

		updatedArea, err := stack.AreaBiz.UpdateArea(ctx, &models.UpdateAreaRequest{
			Id:          area.GetId(),
			Name:        ptr("Campus Yard Updated"),
			Description: ptr("Expanded yard"),
			Geometry:    ptr(businessAreaGeometry),
			Extra:       mustStruct(t, data.JSONMap{"label": "updated"}),
		})
		require.NoError(t, err)
		require.Equal(t, "Campus Yard Updated", updatedArea.GetName())

		searchedAreas, err := stack.AreaBiz.SearchAreas(ctx, "Updated", "", 10)
		require.NoError(t, err)
		require.Len(t, searchedAreas, 1)

		state := &models.GeofenceState{SubjectID: "subject-1", AreaID: area.GetId(), Inside: true}
		state.TenantID = businessTenantID
		state.PartitionID = businessPartitionID
		state.AccessID = "access-1"
		tx := stack.GeofenceStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.GeofenceStateRepo.UpsertTx(tx, state))
		require.NoError(t, tx.Commit().Error)
		require.NoError(t, stack.AreaBiz.DeleteArea(ctx, area.GetId()))

		states, err := stack.GeofenceStateRepo.GetInsideByArea(ctx, area.GetId(), 10)
		require.NoError(t, err)
		require.Empty(t, states)

		threshold := 25.0
		consecutive := int32(2)
		cooldown := int32(30)
		route, err := stack.RouteBiz.CreateRoute(ctx, &models.CreateRouteRequest{
			Data: &geolocationv1.RouteObject{
				OwnerId:                   "owner-1",
				Name:                      "Campus Route",
				Description:               "Main route",
				Geometry:                  businessRouteGeometry,
				DeviationThresholdM:       &threshold,
				DeviationConsecutiveCount: &consecutive,
				DeviationCooldownSec:      &cooldown,
			},
		})
		require.NoError(t, err)
		require.Greater(t, route.GetLengthM(), 0.0)

		updatedRoute, err := stack.RouteBiz.UpdateRoute(ctx, &models.UpdateRouteRequest{
			Id:                  route.GetId(),
			Name:                ptr("Campus Route Updated"),
			Description:         ptr("Route updated"),
			DeviationThresholdM: float64Ptr(30),
			Geometry:            ptr(businessRouteGeometry),
		})
		require.NoError(t, err)
		require.Equal(t, "Campus Route Updated", updatedRoute.GetName())

		assignments, err := stack.RouteBiz.GetSubjectAssignments(ctx, "subject-1")
		require.NoError(t, err)
		require.Empty(t, assignments)

		assignment, err := stack.RouteBiz.AssignRoute(ctx, &models.AssignRouteRequest{
			SubjectId: "subject-1",
			RouteId:   route.GetId(),
		})
		require.NoError(t, err)

		assignments, err = stack.RouteBiz.GetSubjectAssignments(ctx, "subject-1")
		require.NoError(t, err)
		require.Len(t, assignments, 1)

		deviationState := &models.RouteDeviationState{SubjectID: "subject-1", RouteID: route.GetId(), Deviated: true}
		deviationState.TenantID = businessTenantID
		deviationState.PartitionID = businessPartitionID
		deviationState.AccessID = "access-1"
		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.RouteDeviationStateRepo.UpsertTx(tx, deviationState))
		require.NoError(t, tx.Commit().Error)

		require.NoError(t, stack.RouteBiz.UnassignRoute(ctx, assignment.GetId()))
		require.NoError(t, stack.RouteBiz.DeleteRoute(ctx, route.GetId()))

		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		lockedState, err := stack.RouteDeviationStateRepo.GetForUpdate(tx, "subject-1", route.GetId())
		require.NoError(t, err)
		require.Nil(t, lockedState)
		require.NoError(t, tx.Rollback().Error)
	})
}

func (s *BusinessSuite) TestIngestionGeofenceTrackAndProximityFlow() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newBusinessStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-1")

		area, err := stack.AreaBiz.CreateArea(ctx, &models.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{OwnerId: "owner-1", Name: "Tracking Yard", Geometry: businessAreaGeometry},
		})
		require.NoError(t, err)

		now := time.Now().UTC().Add(-10 * time.Minute).Truncate(time.Second)
		resp, err := stack.IngestionBiz.IngestBatch(ctx, &models.IngestLocationsRequest{
			SubjectId: "subject-1",
			Points: []*models.LocationPointInput{
				{
					DeviceId:  "device-subject-1-a",
					Latitude:  0.3390,
					Longitude: 32.5790,
					Accuracy:  5,
					Timestamp: timestamppb.New(now),
				},
				{
					DeviceId:  "device-subject-1-a",
					Latitude:  0.3500,
					Longitude: 32.5900,
					Accuracy:  5,
					Timestamp: timestamppb.New(now.Add(3 * time.Minute)),
				},
				{
					DeviceId:  "device-subject-1-b",
					Latitude:  0.3501,
					Longitude: 32.5901,
					Accuracy:  5,
					Timestamp: timestamppb.New(now.Add(6 * time.Minute)),
				},
				{
					DeviceId:  "device-subject-1-b",
					Latitude:  0.3380,
					Longitude: 32.5780,
					Accuracy:  5,
					Timestamp: timestamppb.New(now.Add(7 * time.Minute)),
				},
				{DeviceId: "device-subject-1-a", Latitude: 100, Longitude: 32.58, Accuracy: 5},
			},
		})
		require.NoError(t, err)
		require.EqualValues(t, 4, resp.GetAccepted())
		require.EqualValues(t, 1, resp.GetRejected())

		resp, err = stack.IngestionBiz.IngestBatch(ctx, &models.IngestLocationsRequest{
			SubjectId: "subject-2",
			Points: []*models.LocationPointInput{
				{
					DeviceId:  "device-subject-2-a",
					Latitude:  0.3502,
					Longitude: 32.5902,
					Accuracy:  4,
					Timestamp: timestamppb.New(now.Add(6 * time.Minute)),
				},
			},
		})
		require.NoError(t, err)
		require.EqualValues(t, 1, resp.GetAccepted())

		pending, err := stack.PointRepo.GetPendingForProcessing(ctx, 10)
		require.NoError(t, err)
		require.Len(t, pending, 5)
		slices.SortFunc(pending, func(left, right *models.LocationPoint) int {
			switch {
			case left.TrueCreatedAt.Before(right.TrueCreatedAt):
				return -1
			case left.TrueCreatedAt.After(right.TrueCreatedAt):
				return 1
			default:
				return 0
			}
		})

		for _, point := range pending {
			require.NoError(t, stack.LocationPointConsumer.Execute(ctx, s.makePointEvent(point, "access-main")))
		}

		track, err := stack.TrackBiz.GetTrack(ctx, &models.GetTrackRequest{SubjectId: "subject-1", Limit: 10})
		require.NoError(t, err)
		require.Len(t, track, 4)

		events, err := stack.TrackBiz.GetSubjectEvents(
			ctx,
			&models.GetSubjectEventsRequest{SubjectId: "subject-1", Limit: 10},
		)
		require.NoError(t, err)
		require.Len(t, events, 3)
		require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_EXIT, events[0].GetEventType())
		require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_DWELL, events[1].GetEventType())
		require.Equal(t, geolocationv1.GeoEventType_GEO_EVENT_TYPE_ENTER, events[2].GetEventType())

		subjects, err := stack.TrackBiz.GetAreaSubjects(ctx, &models.GetAreaSubjectsRequest{AreaId: area.GetId()})
		require.NoError(t, err)
		require.Len(t, subjects, 1)
		require.Equal(t, "subject-2", subjects[0].GetSubjectId())

		latest, err := stack.LatestPositionRepo.Get(ctx, "subject-1")
		require.NoError(t, err)
		require.InDelta(t, 0.3380, latest.Latitude, 0.00001)
		require.Equal(t, "device-subject-1-b", latest.DeviceID)

		nearbySubjects, err := stack.ProximityBiz.GetNearbySubjects(ctx, &models.GetNearbySubjectsRequest{
			SubjectId: "subject-1", RadiusMeters: 5000, Limit: 10,
		})
		require.NoError(t, err)
		require.Len(t, nearbySubjects, 1)
		require.Equal(t, "subject-2", nearbySubjects[0].GetSubjectId())

		nearbyAreas, err := stack.ProximityBiz.GetNearbyAreas(ctx, &models.GetNearbyAreasRequest{
			Latitude: 0.3500, Longitude: 32.5900, RadiusMeters: 1000, Limit: 10,
		})
		require.NoError(t, err)
		require.Len(t, nearbyAreas, 1)
		require.Equal(t, area.GetId(), nearbyAreas[0].GetAreaId())
	})
}

func (s *BusinessSuite) TestRouteDeviationCatchupRetentionAndValidation() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newBusinessStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-route")

		threshold := 20.0
		consecutive := int32(2)
		route, err := stack.RouteBiz.CreateRoute(ctx, &models.CreateRouteRequest{
			Data: &geolocationv1.RouteObject{
				OwnerId: "owner-route", Name: "Deviation Route", Geometry: businessRouteGeometry,
				DeviationThresholdM: &threshold, DeviationConsecutiveCount: &consecutive,
			},
		})
		require.NoError(t, err)

		_, err = stack.RouteBiz.AssignRoute(
			ctx,
			&models.AssignRouteRequest{SubjectId: "subject-route", RouteId: route.GetId()},
		)
		require.NoError(t, err)

		eventBase := models.EventTenancy{
			TenantID:    businessTenantID,
			PartitionID: businessPartitionID,
			AccessID:    "access-route",
		}
		cases := []struct {
			name        string
			event       *models.LocationPointIngestedEvent
			deviated    bool
			consecutive int
		}{
			{
				name: "on route",
				event: &models.LocationPointIngestedEvent{
					EventTenancy: eventBase,
					PointID:      "p1",
					SubjectID:    "subject-route",
					DeviceID:     "device-route-1",
					Latitude:     0.35,
					Longitude:    32.59,
					Accuracy:     3,
					Timestamp:    time.Now().UTC().UnixMilli(),
				},
				deviated:    false,
				consecutive: 0,
			},
			{
				name: "first off route",
				event: &models.LocationPointIngestedEvent{
					EventTenancy: eventBase,
					PointID:      "p2",
					SubjectID:    "subject-route",
					DeviceID:     "device-route-1",
					Latitude:     0.40,
					Longitude:    32.59,
					Accuracy:     3,
					Timestamp:    time.Now().UTC().Add(time.Minute).UnixMilli(),
				},
				deviated:    false,
				consecutive: 1,
			},
			{
				name: "deviated",
				event: &models.LocationPointIngestedEvent{
					EventTenancy: eventBase,
					PointID:      "p3",
					SubjectID:    "subject-route",
					DeviceID:     "device-route-1",
					Latitude:     0.41,
					Longitude:    32.59,
					Accuracy:     3,
					Timestamp:    time.Now().UTC().Add(2 * time.Minute).UnixMilli(),
				},
				deviated:    true,
				consecutive: 2,
			},
			{
				name: "back on route",
				event: &models.LocationPointIngestedEvent{
					EventTenancy: eventBase,
					PointID:      "p4",
					SubjectID:    "subject-route",
					DeviceID:     "device-route-1",
					Latitude:     0.35,
					Longitude:    32.59,
					Accuracy:     3,
					Timestamp:    time.Now().UTC().Add(3 * time.Minute).UnixMilli(),
				},
				deviated:    false,
				consecutive: 0,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.NoError(t, stack.RouteDeviationBiz.EvaluatePoint(ctx, tc.event))

				tx := stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
				state, txErr := stack.RouteDeviationStateRepo.GetForUpdate(tx, "subject-route", route.GetId())
				require.NoError(t, txErr)
				require.NotNil(t, state)
				require.Equal(t, tc.deviated, state.Deviated)
				require.Equal(t, tc.consecutive, state.ConsecutiveOffRoute)
				require.NoError(t, tx.Rollback().Error)
			})
		}

		oldTime := time.Now().UTC().AddDate(0, 0, -10)
		oldPoint := &models.LocationPoint{
			SubjectID:       "subject-old",
			DeviceID:        "device-subject-old",
			TrueCreatedAt:   oldTime,
			IngestedAt:      oldTime,
			Latitude:        0.30,
			Longitude:       32.50,
			Accuracy:        5,
			Source:          models.LocationSourceGPS,
			ProcessingState: models.LocationPointProcessingStatePending,
		}
		oldPoint.GenID(ctx)
		require.NoError(t, stack.PointRepo.Create(ctx, oldPoint))

		oldLatest := &models.LatestPosition{
			SubjectID: "subject-old",
			DeviceID:  "device-subject-old",
			Latitude:  0.30,
			Longitude: 32.50,
			Accuracy:  5,
			TS:        oldTime,
		}
		oldLatest.GenID(ctx)
		require.NoError(t, stack.LatestPositionRepo.Upsert(ctx, oldLatest))

		oldEvent := &models.GeoEvent{
			SubjectID:     "subject-old",
			AreaID:        "area-old",
			EventType:     models.GeoEventTypeEnter,
			TrueCreatedAt: oldTime,
		}
		oldEvent.GenID(ctx)
		require.NoError(t, stack.GeoEventRepo.Create(ctx, oldEvent))

		oldState := &models.GeofenceState{SubjectID: "subject-old", AreaID: "area-old", Inside: true}
		oldState.TenantID = businessTenantID
		oldState.PartitionID = businessPartitionID
		oldState.AccessID = "access-route"
		tx := stack.GeofenceStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.GeofenceStateRepo.UpsertTx(tx, oldState))
		require.NoError(t, tx.Commit().Error)
		require.NoError(t, stack.GeofenceStateRepo.Pool().DB(ctx, false).Table((&models.GeofenceState{}).TableName()).
			Where("subject_id = ?", "subject-old").Update("modified_at", oldTime).Error)

		oldDeviationState := &models.RouteDeviationState{
			SubjectID: "subject-old",
			RouteID:   route.GetId(),
			Deviated:  true,
		}
		oldDeviationState.TenantID = businessTenantID
		oldDeviationState.PartitionID = businessPartitionID
		oldDeviationState.AccessID = "access-route"
		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.RouteDeviationStateRepo.UpsertTx(tx, oldDeviationState))
		require.NoError(t, tx.Commit().Error)
		require.NoError(
			t,
			stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Table((&models.RouteDeviationState{}).TableName()).
				Where("subject_id = ?", "subject-old").Update("modified_at", oldTime).Error,
		)

		require.NoError(t, stack.CatchupBiz.RunCatchup(ctx))

		retention := geobusiness.NewRetentionBusiness(stack.PointRepo.Pool(), geobusiness.RetentionConfig{
			LocationPointRetentionDays: 1,
			GeoEventRetentionDays:      1,
			GeofenceStateStaleDays:     1,
			RetentionBatchSize:         10,
			PartitionMaintenanceMonths: 1,
		})
		require.NoError(t, retention.RunRetention(ctx))

		_, err = stack.LatestPositionRepo.Get(ctx, "subject-old")
		require.Error(t, err)

		points, err := stack.PointRepo.GetTrack(ctx, "subject-old", oldTime.Add(-time.Hour), time.Now().UTC(), 10, 0)
		require.NoError(t, err)
		require.Empty(t, points)

		events, err := stack.GeoEventRepo.GetBySubject(ctx, "subject-old", nil, nil, 10, 0)
		require.NoError(t, err)
		require.Empty(t, events)

		states, err := stack.GeofenceStateRepo.GetInsideByArea(ctx, "area-old", 10)
		require.NoError(t, err)
		require.Empty(t, states)

		validationCases := []struct {
			name string
			run  func() error
		}{
			{
				name: "nil ingest request",
				run:  func() error { _, e := stack.IngestionBiz.IngestBatch(ctx, nil); return e },
			},
			{
				name: "invalid area search",
				run:  func() error { _, e := stack.AreaBiz.SearchAreas(ctx, "", "", 0); return e },
			},
			{
				name: "invalid route search",
				run:  func() error { _, e := stack.RouteBiz.SearchRoutes(ctx, "", 0); return e },
			},
			{name: "invalid nearby subject request", run: func() error {
				_, e := stack.ProximityBiz.GetNearbySubjects(ctx, &models.GetNearbySubjectsRequest{})
				return e
			}},
			{
				name: "invalid track request",
				run:  func() error { _, e := stack.TrackBiz.GetTrack(ctx, &models.GetTrackRequest{}); return e },
			},
		}

		for _, tc := range validationCases {
			t.Run(tc.name, func(t *testing.T) {
				require.Error(t, tc.run())
			})
		}
	})
}

func ptr[T any](value T) *T             { return &value }
func float64Ptr(value float64) *float64 { return &value }

func mustStruct(t *testing.T, value data.JSONMap) *structpb.Struct {
	t.Helper()
	result, err := structpb.NewStruct(value)
	require.NoError(t, err)
	return result
}
