package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/workerpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	georepo "github.com/antinvestor/service-profile/apps/geolocation/service/repository"
	geotests "github.com/antinvestor/service-profile/apps/geolocation/tests"
)

const (
	repoTenantID    = "tenant-repository"
	repoPartitionID = "partition-repository"
)

var (
	//nolint:gochecknoglobals // test fixture data
	repoAreaGeometry = `{"type":"Polygon","coordinates":[[[32.58,0.34],[32.60,0.34],[32.60,0.36],[32.58,0.36],[32.58,0.34]]]}`
	//nolint:gochecknoglobals // test fixture data
	repoRouteGeometry = `{"type":"LineString","coordinates":[[32.58,0.35],[32.60,0.35]]}`
)

type RepositorySuite struct {
	geotests.GeolocationBaseTestSuite
}

type repositoryStack struct {
	PointRepo               georepo.LocationPointRepository
	AreaRepo                georepo.AreaRepository
	GeoEventRepo            georepo.GeoEventRepository
	GeofenceStateRepo       georepo.GeofenceStateRepository
	LatestPositionRepo      georepo.LatestPositionRepository
	RouteRepo               georepo.RouteRepository
	RouteAssignmentRepo     georepo.RouteAssignmentRepository
	RouteDeviationStateRepo georepo.RouteDeviationStateRepository
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) scopedContext(ctx context.Context, subjectID string) context.Context {
	return s.WithAuthClaims(ctx, repoTenantID, repoPartitionID, subjectID)
}

func (s *RepositorySuite) otherTenantContext(ctx context.Context, subjectID string) context.Context {
	return s.WithAuthClaims(ctx, "tenant-other", "partition-other", subjectID)
}

func newRepositoryStack(ctx context.Context, svc *frame.Service) *repositoryStack {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	var workMan workerpool.Manager

	return &repositoryStack{
		PointRepo:               georepo.NewLocationPointRepository(ctx, dbPool, workMan),
		AreaRepo:                georepo.NewAreaRepository(ctx, dbPool, workMan),
		GeoEventRepo:            georepo.NewGeoEventRepository(ctx, dbPool, workMan),
		GeofenceStateRepo:       georepo.NewGeofenceStateRepository(ctx, dbPool, workMan),
		LatestPositionRepo:      georepo.NewLatestPositionRepository(ctx, dbPool, workMan),
		RouteRepo:               georepo.NewRouteRepository(ctx, dbPool, workMan),
		RouteAssignmentRepo:     georepo.NewRouteAssignmentRepository(ctx, dbPool, workMan),
		RouteDeviationStateRepo: georepo.NewRouteDeviationStateRepository(ctx, dbPool, workMan),
	}
}

func (s *RepositorySuite) TestRepositoriesUseMigratedPostGISSchema() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := s.CreateService(t, dep)

		err := georepo.Migrate(ctx, svc.DatastoreManager(), "../../migrations/0001")
		require.NoError(t, err)

		var version string
		require.NoError(t, svc.DatastoreManager().
			GetPool(ctx, datastore.DefaultMigrationPoolName).DB(ctx, true).
			Raw("SELECT PostGIS_Version()").
			Scan(&version).Error)
		require.NotEmpty(t, version)
	})
}

func (s *RepositorySuite) TestLocationPointAndLatestPositionRepositories() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newRepositoryStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-1")
		now := time.Now().UTC().Truncate(time.Second)

		points := []*models.LocationPoint{
			{
				SubjectID:       "subject-1",
				DeviceID:        "device-subject-1-a",
				TrueCreatedAt:   now.Add(-2 * time.Minute),
				IngestedAt:      now.Add(-2 * time.Minute),
				Latitude:        0.3500,
				Longitude:       32.5800,
				Accuracy:        4,
				Source:          models.LocationSourceGPS,
				ProcessingState: models.LocationPointProcessingStatePending,
			},
			{
				SubjectID:       "subject-1",
				DeviceID:        "device-subject-1-b",
				TrueCreatedAt:   now.Add(-time.Minute),
				IngestedAt:      now.Add(-time.Minute),
				Latitude:        0.3502,
				Longitude:       32.5802,
				Accuracy:        5,
				Source:          models.LocationSourceGPS,
				ProcessingState: models.LocationPointProcessingStatePending,
			},
			{
				SubjectID:       "subject-2",
				DeviceID:        "device-subject-2-a",
				TrueCreatedAt:   now,
				IngestedAt:      now,
				Latitude:        0.3503,
				Longitude:       32.5803,
				Accuracy:        3,
				Source:          models.LocationSourceGPS,
				ProcessingState: models.LocationPointProcessingStatePending,
			},
		}
		for _, point := range points {
			point.TenantID = repoTenantID
			point.PartitionID = repoPartitionID
			point.AccessID = "access-repo"
			point.GenID(ctx)
			require.NoError(t, stack.PointRepo.Create(ctx, point))
		}

		foreignCtx := s.otherTenantContext(baseCtx, "foreign-subject")
		foreignPos := &models.LatestPosition{
			SubjectID: "subject-3",
			DeviceID:  "device-subject-3-a",
			Latitude:  0.3504,
			Longitude: 32.5804,
			Accuracy:  2,
			TS:        now,
		}
		foreignPos.TenantID = "tenant-other"
		foreignPos.PartitionID = "partition-other"
		foreignPos.AccessID = "access-other"
		foreignPos.GenID(foreignCtx)
		require.NoError(t, stack.LatestPositionRepo.Upsert(foreignCtx, foreignPos))

		pending, err := stack.PointRepo.GetPendingForProcessing(ctx, 10)
		require.NoError(t, err)
		require.Len(t, pending, 3)

		require.NoError(t, stack.PointRepo.MarkFailed(ctx, points[0].GetID(), errors.New("boom")))
		failed, err := stack.PointRepo.GetByID(ctx, points[0].GetID())
		require.NoError(t, err)
		require.Equal(t, models.LocationPointProcessingStateFailed, failed.ProcessingState)
		require.Contains(t, failed.ProcessingError, "boom")

		require.NoError(t, stack.PointRepo.MarkFailed(ctx, points[2].GetID(), nil))
		failedDefault, err := stack.PointRepo.GetByID(ctx, points[2].GetID())
		require.NoError(t, err)
		require.NotEmpty(t, failedDefault.ProcessingError)

		require.NoError(t, stack.PointRepo.MarkProcessed(ctx, points[1].GetID()))
		processed, err := stack.PointRepo.GetByID(ctx, points[1].GetID())
		require.NoError(t, err)
		require.Equal(t, models.LocationPointProcessingStateProcessed, processed.ProcessingState)
		require.NotNil(t, processed.ProcessedAt)

		track, err := stack.PointRepo.GetTrack(ctx, "subject-1", now.Add(-time.Hour), now.Add(time.Hour), 10, 0)
		require.NoError(t, err)
		require.Len(t, track, 2)
		require.Equal(t, points[1].GetID(), track[0].GetID())

		for _, pos := range []*models.LatestPosition{
			{SubjectID: "subject-1", DeviceID: "device-subject-1-a", Latitude: 0.3500, Longitude: 32.5800, Accuracy: 4, TS: now.Add(-time.Minute)},
			{SubjectID: "subject-1", DeviceID: "device-subject-1-b", Latitude: 0.3501, Longitude: 32.5801, Accuracy: 3, TS: now},
			{SubjectID: "subject-1", DeviceID: "device-subject-1-c", Latitude: 0.1000, Longitude: 31.0000, Accuracy: 1, TS: now.Add(-2 * time.Minute)},
			{SubjectID: "subject-2", DeviceID: "device-subject-2-a", Latitude: 0.3502, Longitude: 32.5802, Accuracy: 2, TS: now},
		} {
			pos.TenantID = repoTenantID
			pos.PartitionID = repoPartitionID
			pos.AccessID = "access-repo"
			pos.GenID(ctx)
			require.NoError(t, stack.LatestPositionRepo.Upsert(ctx, pos))
		}

		current, err := stack.LatestPositionRepo.Get(ctx, "subject-1")
		require.NoError(t, err)
		require.InDelta(t, 0.3501, current.Latitude, 0.00001)
		require.InDelta(t, 32.5801, current.Longitude, 0.00001)

		nearby, err := stack.LatestPositionRepo.GetNearbySubjects(ctx, 0.3501, 32.5801, 100, "subject-1", 1, 10)
		require.NoError(t, err)
		require.Len(t, nearby, 1)
		require.Equal(t, "subject-2", nearby[0].SubjectID)
	})
}

func (s *RepositorySuite) TestAreaGeoEventAndGeofenceStateRepositories() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newRepositoryStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-area")

		area := &models.Area{
			OwnerID:      "owner-1",
			Name:         "Main Yard",
			Description:  "North yard",
			AreaType:     models.AreaTypeLand,
			GeometryJSON: repoAreaGeometry,
			State:        2,
			Extras:       data.JSONMap{"zone": "north"},
		}
		area.TenantID = repoTenantID
		area.PartitionID = repoPartitionID
		area.AccessID = "access-area"
		area.GenID(ctx)
		require.NoError(t, stack.AreaRepo.Create(ctx, area))
		require.NoError(t, stack.AreaRepo.UpdateGeometry(ctx, area.GetID(), repoAreaGeometry))
		tx := stack.AreaRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.AreaRepo.UpdateGeometryTx(tx, area.GetID(), repoAreaGeometry))
		require.NoError(t, tx.Commit().Error)

		persisted, err := stack.AreaRepo.GetByID(ctx, area.GetID())
		require.NoError(t, err)
		require.NotNil(t, persisted.AreaM2)
		require.NotNil(t, persisted.PerimeterM)
		require.NotNil(t, persisted.Bbox)

		candidates, err := stack.AreaRepo.GetActiveByBoundingBox(ctx, 0.35, 32.59)
		require.NoError(t, err)
		require.Len(t, candidates, 1)

		contains, err := stack.AreaRepo.ContainsPoint(ctx, area.GetID(), 0.35, 32.59)
		require.NoError(t, err)
		require.True(t, contains)

		nearby, err := stack.AreaRepo.GetNearbyAreas(ctx, 0.35, 32.59, 1000, 10)
		require.NoError(t, err)
		require.Len(t, nearby, 1)
		require.Equal(t, area.GetID(), nearby[0].Area.GetID())

		byOwner, err := stack.AreaRepo.SearchByOwner(ctx, "owner-1", 10)
		require.NoError(t, err)
		require.Len(t, byOwner, 1)

		byQuery, err := stack.AreaRepo.SearchByQuery(ctx, "yard", 10)
		require.NoError(t, err)
		require.Len(t, byQuery, 1)

		event := &models.GeoEvent{
			SubjectID:     "subject-area",
			AreaID:        area.GetID(),
			EventType:     models.GeoEventTypeDwell,
			TrueCreatedAt: time.Now().UTC(),
			Confidence:    0.95,
			PointID:       "point-1",
		}
		event.GenID(ctx)

		tx = stack.GeoEventRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.GeoEventRepo.CreateTx(tx, event))
		require.NoError(t, tx.Commit().Error)

		eventsBySubject, err := stack.GeoEventRepo.GetBySubject(ctx, "subject-area", nil, nil, 10, 0)
		require.NoError(t, err)
		require.Len(t, eventsBySubject, 1)

		eventsByArea, err := stack.GeoEventRepo.GetByArea(ctx, area.GetID(), nil, nil, 10, 0)
		require.NoError(t, err)
		require.Len(t, eventsByArea, 1)

		from := event.TrueCreatedAt.Add(-time.Minute)
		hasDwell, err := stack.GeoEventRepo.HasDwellEvent(ctx, "subject-area", area.GetID(), &from)
		require.NoError(t, err)
		require.True(t, hasDwell)

		tx = stack.GeoEventRepo.Pool().DB(ctx, false).Begin()
		hasDwell, err = stack.GeoEventRepo.HasDwellEventTx(tx, "subject-area", area.GetID(), &from)
		require.NoError(t, err)
		require.True(t, hasDwell)
		require.NoError(t, tx.Rollback().Error)

		pointTS := time.Now().UTC()
		enterTS := pointTS.Add(-time.Minute)
		state := &models.GeofenceState{
			SubjectID:      "subject-area",
			AreaID:         area.GetID(),
			Inside:         true,
			LastTransition: &pointTS,
			LastPointTS:    &pointTS,
			EnterTS:        &enterTS,
			LastLat:        0.35,
			LastLon:        32.59,
		}
		state.TenantID = repoTenantID
		state.PartitionID = repoPartitionID
		state.AccessID = "access-state"

		tx = stack.GeofenceStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.GeofenceStateRepo.UpsertTx(tx, state))

		lockedState, err := stack.GeofenceStateRepo.GetForUpdate(tx, "subject-area", area.GetID())
		require.NoError(t, err)
		require.NotNil(t, lockedState)
		require.True(t, lockedState.Inside)
		require.NoError(t, tx.Commit().Error)

		inside, err := stack.GeofenceStateRepo.GetInsideByArea(ctx, area.GetID(), 10)
		require.NoError(t, err)
		require.Len(t, inside, 1)
		insideBySubject, err := stack.GeofenceStateRepo.GetInsideBySubject(ctx, "subject-area", 10)
		require.NoError(t, err)
		require.Len(t, insideBySubject, 1)

		require.NoError(t, stack.GeofenceStateRepo.DeleteByArea(ctx, area.GetID()))
		inside, err = stack.GeofenceStateRepo.GetInsideByArea(ctx, area.GetID(), 10)
		require.NoError(t, err)
		require.Empty(t, inside)
	})
}

func (s *RepositorySuite) TestRouteAndDeviationRepositories() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newRepositoryStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-route")

		threshold := 25.0
		consecutive := 2
		cooldown := 60
		route := &models.Route{
			OwnerID:                   "owner-route",
			Name:                      "Main Route",
			Description:               "Campus path",
			GeometryJSON:              repoRouteGeometry,
			State:                     2,
			DeviationThresholdM:       &threshold,
			DeviationConsecutiveCount: &consecutive,
			DeviationCooldownSec:      &cooldown,
		}
		route.TenantID = repoTenantID
		route.PartitionID = repoPartitionID
		route.AccessID = "access-route"
		route.GenID(ctx)
		require.NoError(t, stack.RouteRepo.Create(ctx, route))
		require.NoError(t, stack.RouteRepo.UpdateGeometry(ctx, route.GetID(), repoRouteGeometry))
		tx := stack.RouteRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.RouteRepo.UpdateGeometryTx(tx, route.GetID(), repoRouteGeometry))
		require.NoError(t, tx.Commit().Error)

		persisted, err := stack.RouteRepo.GetByID(ctx, route.GetID())
		require.NoError(t, err)
		require.NotNil(t, persisted.LengthM)
		require.Greater(t, *persisted.LengthM, 0.0)

		routes, err := stack.RouteRepo.SearchByOwner(ctx, "owner-route", 10)
		require.NoError(t, err)
		require.Len(t, routes, 1)

		assignment := &models.RouteAssignment{SubjectID: "subject-route", RouteID: route.GetID(), State: 2}
		assignment.TenantID = repoTenantID
		assignment.PartitionID = repoPartitionID
		assignment.AccessID = "access-route"
		assignment.GenID(ctx)
		require.NoError(t, stack.RouteAssignmentRepo.Create(ctx, assignment))

		assignments, err := stack.RouteAssignmentRepo.GetBySubject(ctx, "subject-route")
		require.NoError(t, err)
		require.Len(t, assignments, 1)

		active, err := stack.RouteRepo.GetActiveAssignmentsForSubject(ctx, "subject-route", time.Now().UTC())
		require.NoError(t, err)
		require.Len(t, active, 1)
		require.Equal(t, route.GetID(), active[0].Route.GetID())

		distance, err := stack.RouteRepo.DistanceToRouteMeters(ctx, route.GetID(), 0.35, 32.5901)
		require.NoError(t, err)
		require.GreaterOrEqual(t, distance, 0.0)

		pointTS := time.Now().UTC()
		state := &models.RouteDeviationState{
			SubjectID:            "subject-route",
			RouteID:              route.GetID(),
			Deviated:             true,
			ConsecutiveOffRoute:  3,
			LastDeviationEventAt: &pointTS,
			LastPointTS:          &pointTS,
			LastLat:              0.36,
			LastLon:              32.61,
		}
		state.TenantID = repoTenantID
		state.PartitionID = repoPartitionID
		state.AccessID = "access-route"

		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		require.NoError(t, stack.RouteDeviationStateRepo.UpsertTx(tx, state))

		lockedState, err := stack.RouteDeviationStateRepo.GetForUpdate(tx, "subject-route", route.GetID())
		require.NoError(t, err)
		require.NotNil(t, lockedState)
		require.True(t, lockedState.Deviated)
		require.NoError(t, tx.Commit().Error)

		require.NoError(t, stack.RouteDeviationStateRepo.DeleteByRoute(ctx, route.GetID()))
		tx = stack.RouteDeviationStateRepo.Pool().DB(ctx, false).Begin()
		lockedState, err = stack.RouteDeviationStateRepo.GetForUpdate(tx, "subject-route", route.GetID())
		require.NoError(t, err)
		require.Nil(t, lockedState)
		require.NoError(t, tx.Rollback().Error)

		require.NoError(t, stack.RouteAssignmentRepo.DeleteByRoute(ctx, route.GetID()))
		assignments, err = stack.RouteAssignmentRepo.GetBySubject(ctx, "subject-route")
		require.NoError(t, err)
		require.Empty(t, assignments)
	})
}
