package business_test

import (
	"context"
	"testing"
	"time"

	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	geobusiness "github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	geolocationv1 "github.com/antinvestor/service-profile/proto/geolocation/geolocation/v1"
)

func (s *BusinessSuite) TestLookupsSchedulersAndEdgeCases() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newBusinessStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-coverage")

		area, err := stack.AreaBiz.CreateArea(ctx, &models.CreateAreaRequest{
			Data: &geolocationv1.AreaObject{
				OwnerId:  "owner-coverage",
				Name:     "Coverage Area",
				Geometry: businessAreaGeometry,
			},
		})
		require.NoError(t, err)

		gotArea, err := stack.AreaBiz.GetArea(ctx, area.GetId())
		require.NoError(t, err)
		require.Equal(t, area.GetId(), gotArea.GetId())

		threshold := 15.0
		route, err := stack.RouteBiz.CreateRoute(ctx, &models.CreateRouteRequest{
			Data: &geolocationv1.RouteObject{
				OwnerId:             "owner-coverage",
				Name:                "Coverage Route",
				Geometry:            businessRouteGeometry,
				DeviationThresholdM: &threshold,
			},
		})
		require.NoError(t, err)

		gotRoute, err := stack.RouteBiz.GetRoute(ctx, route.GetId())
		require.NoError(t, err)
		require.Equal(t, route.GetId(), gotRoute.GetId())

		routes, err := stack.RouteBiz.SearchRoutes(ctx, "owner-coverage", 0)
		require.NoError(t, err)
		require.Len(t, routes, 1)

		validFrom := time.Now().UTC()
		validUntil := validFrom.Add(time.Hour)
		assignment, err := stack.RouteBiz.AssignRoute(ctx, &models.AssignRouteRequest{
			SubjectId:  "subject-coverage",
			RouteId:    route.GetId(),
			ValidFrom:  timestamppb.New(validFrom),
			ValidUntil: timestamppb.New(validUntil),
		})
		require.NoError(t, err)
		require.NotEmpty(t, assignment.GetId())
		require.NoError(t, stack.RouteBiz.UnassignRoute(ctx, assignment.GetId()))

		validationCases := []struct {
			name string
			run  func() error
		}{
			{
				name: "get missing area",
				run:  func() error { _, e := stack.AreaBiz.GetArea(ctx, "missing"); return e },
			},
			{
				name: "get missing route",
				run:  func() error { _, e := stack.RouteBiz.GetRoute(ctx, "missing"); return e },
			},
			{
				name: "assign nil request",
				run:  func() error { _, e := stack.RouteBiz.AssignRoute(ctx, nil); return e },
			},
			{name: "assign missing subject", run: func() error {
				_, e := stack.RouteBiz.AssignRoute(ctx, &models.AssignRouteRequest{RouteId: route.GetId()})
				return e
			}},
			{name: "assign missing route", run: func() error {
				_, e := stack.RouteBiz.AssignRoute(ctx, &models.AssignRouteRequest{SubjectId: "subject-coverage"})
				return e
			}},
			{name: "assign missing route entity", run: func() error {
				_, e := stack.RouteBiz.AssignRoute(
					ctx,
					&models.AssignRouteRequest{SubjectId: "subject-coverage", RouteId: "missing"},
				)
				return e
			}},
			{
				name: "unassign missing assignment",
				run:  func() error { return stack.RouteBiz.UnassignRoute(ctx, "missing") },
			},
			{name: "geofence nil event", run: func() error { return stack.GeofenceBiz.EvaluatePoint(ctx, nil) }},
			{
				name: "route deviation nil event",
				run:  func() error { return stack.RouteDeviationBiz.EvaluatePoint(ctx, nil) },
			},
		}
		for _, tc := range validationCases {
			t.Run(tc.name, func(t *testing.T) {
				require.Error(t, tc.run())
			})
		}

		require.NoError(t, stack.GeofenceBiz.EvaluatePoint(ctx, &models.LocationPointIngestedEvent{
			EventTenancy: models.EventTenancy{
				TenantID:    businessTenantID,
				PartitionID: businessPartitionID,
				AccessID:    "access-coverage",
			},
			PointID:   "skip-accuracy",
			SubjectID: "subject-coverage",
			DeviceID:  "device-coverage",
			Latitude:  0.35,
			Longitude: 32.59,
			Accuracy:  0,
			Timestamp: time.Now().UTC().UnixMilli(),
		}))
		require.NoError(t, stack.RouteDeviationBiz.EvaluatePoint(ctx, &models.LocationPointIngestedEvent{
			EventTenancy: models.EventTenancy{
				TenantID:    businessTenantID,
				PartitionID: businessPartitionID,
				AccessID:    "access-coverage",
			},
			PointID:   "skip-accuracy-route",
			SubjectID: "subject-coverage",
			DeviceID:  "device-coverage",
			Latitude:  0.35,
			Longitude: 32.59,
			Accuracy:  1000,
			Timestamp: time.Now().UTC().UnixMilli(),
		}))

		retention := geobusiness.NewRetentionBusiness(stack.PointRepo.Pool(), geobusiness.RetentionConfig{
			LocationPointRetentionDays: 1,
			GeoEventRetentionDays:      1,
			GeofenceStateStaleDays:     1,
			RetentionBatchSize:         10,
			PartitionMaintenanceMonths: 1,
			RetentionInterval:          time.Millisecond,
		})
		retentionCtx, cancelRetention := context.WithCancel(ctx)
		retentionDone := make(chan struct{})
		go func() {
			defer close(retentionDone)
			retention.StartScheduler(retentionCtx)
		}()
		time.Sleep(5 * time.Millisecond)
		cancelRetention()
		<-retentionDone

		catchup := geobusiness.NewCatchupBusiness(stack.PointRepo, svc.EventsManager(), geobusiness.CatchupConfig{
			BatchSize: 1,
			Interval:  time.Millisecond,
		})
		catchupCtx, cancelCatchup := context.WithCancel(ctx)
		catchupDone := make(chan struct{})
		go func() {
			defer close(catchupDone)
			catchup.StartScheduler(catchupCtx)
		}()
		time.Sleep(5 * time.Millisecond)
		cancelCatchup()
		<-catchupDone
	})
}

func (s *BusinessSuite) TestBusinessValidationAndLookupErrors() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		stack := newBusinessStack(baseCtx, svc)
		ctx := s.scopedContext(baseCtx, "subject-validation")

		cases := []struct {
			name string
			run  func() error
		}{
			{name: "create nil area", run: func() error { _, err := stack.AreaBiz.CreateArea(ctx, nil); return err }},
			{name: "update nil area", run: func() error { _, err := stack.AreaBiz.UpdateArea(ctx, nil); return err }},
			{name: "delete missing area", run: func() error { return stack.AreaBiz.DeleteArea(ctx, "missing") }},
			{
				name: "get missing area",
				run:  func() error { _, err := stack.AreaBiz.GetArea(ctx, "missing"); return err },
			},
			{
				name: "create nil route",
				run:  func() error { _, err := stack.RouteBiz.CreateRoute(ctx, nil); return err },
			},
			{
				name: "update nil route",
				run:  func() error { _, err := stack.RouteBiz.UpdateRoute(ctx, nil); return err },
			},
			{name: "delete missing route", run: func() error { return stack.RouteBiz.DeleteRoute(ctx, "missing") }},
			{
				name: "get missing route",
				run:  func() error { _, err := stack.RouteBiz.GetRoute(ctx, "missing"); return err },
			},
			{
				name: "subject events nil request",
				run:  func() error { _, err := stack.TrackBiz.GetSubjectEvents(ctx, nil); return err },
			},
			{
				name: "area subjects nil request",
				run:  func() error { _, err := stack.TrackBiz.GetAreaSubjects(ctx, nil); return err },
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.Error(t, tc.run())
			})
		}
	})
}
