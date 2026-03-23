package events //nolint:testpackage // tests access unexported event internals

import (
	"context"
	"testing"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/workerpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	aconfig "github.com/antinvestor/service-profile/apps/geolocation/config"
	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
	geotests "github.com/antinvestor/service-profile/apps/geolocation/tests"
	geolocationv1 "github.com/antinvestor/service-profile/proto/geolocation/v1"
)

const (
	eventsTenantID    = "tenant-events"
	eventsPartitionID = "partition-events"
)

//nolint:gochecknoglobals // test fixture data
var eventsAreaGeometry = `{"type":"Polygon","coordinates":[[[32.58,0.34],[32.60,0.34],[32.60,0.36],[32.58,0.36],[32.58,0.34]]]}`

type EventsSuite struct {
	geotests.GeolocationBaseTestSuite
}

type eventsStack struct {
	PointRepo               repository.LocationPointRepository
	AreaRepo                repository.AreaRepository
	GeoEventRepo            repository.GeoEventRepository
	GeofenceStateRepo       repository.GeofenceStateRepository
	LatestPositionRepo      repository.LatestPositionRepository
	RouteRepo               repository.RouteRepository
	RouteAssignmentRepo     repository.RouteAssignmentRepository
	RouteDeviationStateRepo repository.RouteDeviationStateRepository
	AreaBiz                 business.AreaBusiness
	LocationPointConsumer   *LocationPointConsumer
	AreaChangeConsumer      *AreaChangeConsumer
	RouteChangeConsumer     *RouteChangeConsumer
	GeoEventConsumer        *GeoEventConsumer
	RouteDeviationConsumer  *RouteDeviationConsumer
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsSuite))
}

func (s *EventsSuite) scopedContext(ctx context.Context, subjectID string) context.Context {
	return s.WithAuthClaims(ctx, eventsTenantID, eventsPartitionID, subjectID)
}

func newEventsStack(ctx context.Context, svc *frame.Service) *eventsStack {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	var workMan workerpool.Manager
	cfg := svc.Config().(*aconfig.GeolocationConfig)
	metrics := observability.NewMetrics()
	evtsMan := svc.EventsManager()

	pointRepo := repository.NewLocationPointRepository(ctx, dbPool, workMan)
	areaRepo := repository.NewAreaRepository(ctx, dbPool, workMan)
	geoEventRepo := repository.NewGeoEventRepository(ctx, dbPool, workMan)
	stateRepo := repository.NewGeofenceStateRepository(ctx, dbPool, workMan)
	latestPosRepo := repository.NewLatestPositionRepository(ctx, dbPool, workMan)
	routeRepo := repository.NewRouteRepository(ctx, dbPool, workMan)
	routeAssignmentRepo := repository.NewRouteAssignmentRepository(ctx, dbPool, workMan)
	routeDeviationStateRepo := repository.NewRouteDeviationStateRepository(ctx, dbPool, workMan)

	areaBiz := business.NewAreaBusiness(evtsMan, areaRepo, stateRepo)
	proximityBiz := business.NewProximityBusiness(latestPosRepo, areaRepo, cfg.ProximityBusinessConfig())
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

	return &eventsStack{
		PointRepo:               pointRepo,
		AreaRepo:                areaRepo,
		GeoEventRepo:            geoEventRepo,
		GeofenceStateRepo:       stateRepo,
		LatestPositionRepo:      latestPosRepo,
		RouteRepo:               routeRepo,
		RouteAssignmentRepo:     routeAssignmentRepo,
		RouteDeviationStateRepo: routeDeviationStateRepo,
		AreaBiz:                 areaBiz,
		LocationPointConsumer: NewLocationPointConsumer(
			pointRepo,
			proximityBiz,
			geofenceBiz,
			routeDeviationBiz,
			metrics,
		),
		AreaChangeConsumer:     NewAreaChangeConsumer(areaBiz, stateRepo),
		RouteChangeConsumer:    NewRouteChangeConsumer(routeAssignmentRepo, routeDeviationStateRepo),
		GeoEventConsumer:       NewGeoEventConsumer(),
		RouteDeviationConsumer: NewRouteDeviationConsumer(),
	}
}

func (s *EventsSuite) pointEvent(point *models.LocationPoint) *models.LocationPointIngestedEvent {
	return &models.LocationPointIngestedEvent{
		EventTenancy: models.EventTenancy{
			TenantID:    eventsTenantID,
			PartitionID: eventsPartitionID,
			AccessID:    "access-events",
		},
		PointID:   point.GetID(),
		SubjectID: point.SubjectID,
		DeviceID:  point.DeviceID,
		Latitude:  point.Latitude,
		Longitude: point.Longitude,
		Accuracy:  point.Accuracy,
		Timestamp: point.TS.UnixMilli(),
	}
}

func (s *EventsSuite) TestLocationPointConsumerSuccessAndFailure() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newEventsStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-1")

		_, err := stack.AreaBiz.CreateArea(ctx, &models.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{
				OwnerId:  "owner-1",
				Name:     "Events Area",
				Geometry: eventsAreaGeometry,
			},
		})
		require.NoError(t, err)

		point := &models.LocationPoint{
			SubjectID:       "subject-1",
			DeviceID:        "device-subject-1",
			TS:              time.Now().UTC(),
			IngestedAt:      time.Now().UTC(),
			Latitude:        0.35,
			Longitude:       32.59,
			Accuracy:        4,
			Source:          models.LocationSourceGPS,
			ProcessingState: models.LocationPointProcessingStatePending,
		}
		point.GenID(ctx)
		require.NoError(t, stack.PointRepo.Create(ctx, point))

		require.Equal(t, business.LocationPointIngestedEventName, stack.LocationPointConsumer.Name())
		require.IsType(t, &models.LocationPointIngestedEvent{}, stack.LocationPointConsumer.PayloadType())
		require.NoError(t, stack.LocationPointConsumer.Validate(ctx, s.pointEvent(point)))
		require.NoError(t, stack.LocationPointConsumer.Execute(ctx, s.pointEvent(point)))

		processed, err := stack.PointRepo.GetByID(ctx, point.GetID())
		require.NoError(t, err)
		require.Equal(t, models.LocationPointProcessingStateProcessed, processed.ProcessingState)

		latest, err := stack.LatestPositionRepo.Get(ctx, "subject-1")
		require.NoError(t, err)
		require.InDelta(t, 0.35, latest.Latitude, 0.00001)
		require.Equal(t, "device-subject-1", latest.DeviceID)

		require.Error(t, stack.LocationPointConsumer.Validate(ctx, &models.LocationPointIngestedEvent{}))
		require.Error(
			t,
			stack.LocationPointConsumer.Validate(ctx, &models.LocationPointIngestedEvent{SubjectID: "subject-1"}),
		)
		require.Error(
			t,
			stack.LocationPointConsumer.Validate(
				ctx,
				&models.LocationPointIngestedEvent{SubjectID: "subject-1", DeviceID: "device-1"},
			),
		)
		require.Error(t, stack.LocationPointConsumer.Execute(ctx, "wrong-type"))

		_, err = stack.AreaBiz.CreateArea(ctx, &models.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{
				OwnerId:  "owner-bad",
				Name:     "Broken Area",
				Geometry: eventsAreaGeometry,
			},
		})
		require.NoError(t, err)

		badPoint := &models.LocationPoint{
			SubjectID:       "subject-bad",
			DeviceID:        "device-subject-bad",
			TS:              time.Now().UTC(),
			IngestedAt:      time.Now().UTC(),
			Latitude:        0.35,
			Longitude:       32.59,
			Accuracy:        4,
			Source:          models.LocationSourceGPS,
			ProcessingState: models.LocationPointProcessingStatePending,
		}
		badPoint.GenID(ctx)
		require.NoError(t, stack.PointRepo.Create(ctx, badPoint))

		require.NoError(t, stack.AreaRepo.Pool().DB(ctx, false).
			Exec("ALTER TABLE areas DROP COLUMN bbox").Error)

		err = stack.LocationPointConsumer.Execute(ctx, s.pointEvent(badPoint))
		require.Error(t, err)

		failed, err := stack.PointRepo.GetByID(ctx, badPoint.GetID())
		require.NoError(t, err)
		require.Equal(t, models.LocationPointProcessingStateFailed, failed.ProcessingState)
		require.NotEmpty(t, failed.ProcessingError)
	})
}

func (s *EventsSuite) TestOtherConsumers() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newEventsStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-1")

		area, err := stack.AreaBiz.CreateArea(ctx, &models.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{
				OwnerId:  "owner-1",
				Name:     "Delete Area",
				Geometry: eventsAreaGeometry,
			},
		})
		require.NoError(t, err)

		state := &models.GeofenceState{SubjectID: "subject-1", AreaID: area.GetId(), Inside: true}
		state.TenantID = eventsTenantID
		state.PartitionID = eventsPartitionID
		state.AccessID = "access-events"
		tx := stack.GeofenceStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.GeofenceStateRepo.UpsertTx(tx, state))
		require.NoError(t, tx.Commit().Error)

		areaConsumer := stack.AreaChangeConsumer
		require.Equal(t, business.AreaChangedEventName, areaConsumer.Name())
		require.IsType(t, &models.AreaChangedEvent{}, areaConsumer.PayloadType())
		require.NoError(
			t,
			areaConsumer.Validate(ctx, &models.AreaChangedEvent{AreaID: area.GetId(), Action: "deleted"}),
		)
		require.NoError(t, areaConsumer.Execute(ctx, &models.AreaChangedEvent{
			EventTenancy: models.EventTenancy{
				TenantID:    eventsTenantID,
				PartitionID: eventsPartitionID,
				AccessID:    "access-events",
			},
			AreaID:  area.GetId(),
			Action:  "deleted",
			OwnerID: "owner-1",
		}))
		states, err := stack.GeofenceStateRepo.GetInsideByArea(ctx, area.GetId(), 10)
		require.NoError(t, err)
		require.Empty(t, states)

		geoConsumer := stack.GeoEventConsumer
		require.Equal(t, business.GeoEventEmittedEventName, geoConsumer.Name())
		require.IsType(t, &models.GeoEventEmitted{}, geoConsumer.PayloadType())
		require.NoError(t, geoConsumer.Validate(ctx, &models.GeoEventEmitted{
			EventID:   "event-1",
			SubjectID: "subject-1",
			AreaID:    area.GetId(),
		}))
		require.NoError(t, geoConsumer.Execute(ctx, &models.GeoEventEmitted{
			EventTenancy: models.EventTenancy{
				TenantID:    eventsTenantID,
				PartitionID: eventsPartitionID,
				AccessID:    "access-events",
			},
			EventID:    "event-1",
			SubjectID:  "subject-1",
			AreaID:     area.GetId(),
			EventType:  models.GeoEventTypeEnter,
			Timestamp:  time.Now().UTC().UnixMilli(),
			Confidence: 0.9,
		}))

		routeConsumer := stack.RouteDeviationConsumer
		require.Equal(t, business.RouteDeviationDetectedEventName, routeConsumer.Name())
		require.IsType(t, &models.RouteDeviationDetectedEvent{}, routeConsumer.PayloadType())
		require.NoError(t, routeConsumer.Validate(ctx, &models.RouteDeviationDetectedEvent{
			SubjectID: "subject-1",
			RouteID:   "route-1",
			EventType: "deviated",
		}))
		require.Error(
			t,
			routeConsumer.Validate(ctx, &models.RouteDeviationDetectedEvent{RouteID: "route-1", EventType: "deviated"}),
		)
		require.Error(
			t,
			routeConsumer.Validate(
				ctx,
				&models.RouteDeviationDetectedEvent{SubjectID: "subject-1", EventType: "deviated"},
			),
		)
		require.Error(
			t,
			routeConsumer.Validate(
				ctx,
				&models.RouteDeviationDetectedEvent{SubjectID: "subject-1", RouteID: "route-1"},
			),
		)
		require.NoError(t, routeConsumer.Execute(ctx, &models.RouteDeviationDetectedEvent{
			EventTenancy: models.EventTenancy{
				TenantID:    eventsTenantID,
				PartitionID: eventsPartitionID,
				AccessID:    "access-events",
			},
			SubjectID:      "subject-1",
			RouteID:        "route-1",
			EventType:      "deviated",
			DistanceMeters: 120,
			Timestamp:      time.Now().UTC().UnixMilli(),
		}))

		threshold := 20.0
		route := &models.Route{
			OwnerID:             "owner-1",
			Name:                "Cleanup Route",
			GeometryJSON:        `{"type":"LineString","coordinates":[[32.58,0.35],[32.60,0.35]]}`,
			State:               2,
			DeviationThresholdM: &threshold,
		}
		route.GenID(ctx)
		require.NoError(t, stack.RouteRepo.Create(ctx, route))
		require.NoError(t, stack.RouteRepo.UpdateGeometry(ctx, route.GetID(), route.GeometryJSON))

		assignment := &models.RouteAssignment{SubjectID: "subject-1", RouteID: route.GetID(), State: 2}
		assignment.GenID(ctx)
		require.NoError(t, stack.RouteAssignmentRepo.Create(ctx, assignment))

		deviationState := &models.RouteDeviationState{SubjectID: "subject-1", RouteID: route.GetID(), Deviated: true}
		deviationState.TenantID = eventsTenantID
		deviationState.PartitionID = eventsPartitionID
		deviationState.AccessID = "access-events"
		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.RouteDeviationStateRepo.UpsertTx(tx, deviationState))
		require.NoError(t, tx.Commit().Error)

		routeChangeConsumer := stack.RouteChangeConsumer
		require.Equal(t, business.RouteChangedEventName, routeChangeConsumer.Name())
		require.IsType(t, &models.RouteChangedEvent{}, routeChangeConsumer.PayloadType())
		require.NoError(
			t,
			routeChangeConsumer.Validate(ctx, &models.RouteChangedEvent{RouteID: route.GetID(), Action: "deleted"}),
		)
		require.NoError(t, routeChangeConsumer.Execute(ctx, &models.RouteChangedEvent{
			EventTenancy: models.EventTenancy{
				TenantID:    eventsTenantID,
				PartitionID: eventsPartitionID,
				AccessID:    "access-events",
			},
			RouteID: route.GetID(),
			Action:  "deleted",
			OwnerID: "owner-1",
		}))
		assignments, err := stack.RouteAssignmentRepo.GetBySubject(ctx, "subject-1")
		require.NoError(t, err)
		require.Empty(t, assignments)
		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		lockedDeviationState, err := stack.RouteDeviationStateRepo.GetForUpdate(tx, "subject-1", route.GetID())
		require.NoError(t, err)
		require.Nil(t, lockedDeviationState)
		require.NoError(t, tx.Rollback().Error)

		require.Error(t, geoConsumer.Validate(ctx, "wrong"))
		require.Error(
			t,
			geoConsumer.Validate(ctx, &models.GeoEventEmitted{SubjectID: "subject-1", AreaID: area.GetId()}),
		)
		require.Error(t, geoConsumer.Validate(ctx, &models.GeoEventEmitted{EventID: "event-1", AreaID: area.GetId()}))
		require.Error(t, geoConsumer.Validate(ctx, &models.GeoEventEmitted{EventID: "event-1", SubjectID: "subject-1"}))
		require.Error(t, areaConsumer.Validate(ctx, "wrong"))
		require.Error(t, areaConsumer.Validate(ctx, &models.AreaChangedEvent{Action: "deleted"}))
		require.Error(t, areaConsumer.Validate(ctx, &models.AreaChangedEvent{AreaID: area.GetId()}))
		require.Error(t, routeChangeConsumer.Validate(ctx, "wrong"))
		require.Error(t, routeChangeConsumer.Validate(ctx, &models.RouteChangedEvent{Action: "deleted"}))
		require.Error(t, routeChangeConsumer.Validate(ctx, &models.RouteChangedEvent{RouteID: route.GetID()}))
		require.Error(t, routeConsumer.Execute(ctx, "wrong"))
		require.Error(t, routeChangeConsumer.Execute(ctx, "wrong"))
	})
}
