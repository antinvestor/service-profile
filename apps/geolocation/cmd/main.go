package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/antinvestor/common/permissions"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
	connectInterceptors "github.com/pitabwire/frame/security/interceptors/connect"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/geolocation/config"
	"github.com/antinvestor/service-profile/apps/geolocation/service/authz"
	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/events"
	"github.com/antinvestor/service-profile/apps/geolocation/service/handlers"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
	geolocationv1 "github.com/antinvestor/service-profile/proto/geolocation/geolocation/v1"
	"github.com/antinvestor/service-profile/proto/geolocation/geolocation/v1/geolocationv1connect"
)

func main() { //nolint:funlen // wiring function
	ctx := context.Background()

	cfg, err := config.LoadWithOIDC[aconfig.GeolocationConfig](ctx)
	if err != nil {
		util.Log(ctx).With("err", err).Error("could not process configs")
		return
	}

	if cfg.Name() == "" {
		cfg.ServiceName = "service_geolocation"
	}

	ctx, svc := frame.NewServiceWithContext(
		ctx,
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
	)
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	dbManager := svc.DatastoreManager()

	// Handle database migration if requested.
	if cfg.DoDatabaseMigrate() {
		if mErr := repository.Migrate(ctx, dbManager, cfg.GetDatabaseMigrationPath()); mErr != nil {
			log.WithError(mErr).Fatal("could not migrate database")
		}
		return
	}

	// Initialize repositories.
	dbPool := dbManager.GetPool(ctx, datastore.DefaultPoolName)
	workMan, evtsMan := svc.WorkManager(), svc.EventsManager()

	pointRepo := repository.NewLocationPointRepository(ctx, dbPool, workMan)
	areaRepo := repository.NewAreaRepository(ctx, dbPool, workMan)
	geoEventRepo := repository.NewGeoEventRepository(ctx, dbPool, workMan)
	stateRepo := repository.NewGeofenceStateRepository(ctx, dbPool, workMan)
	latestPosRepo := repository.NewLatestPositionRepository(ctx, dbPool, workMan)
	routeRepo := repository.NewRouteRepository(ctx, dbPool, workMan)
	routeAssignmentRepo := repository.NewRouteAssignmentRepository(ctx, dbPool, workMan)
	routeDeviationStateRepo := repository.NewRouteDeviationStateRepository(ctx, dbPool, workMan)

	// Initialize observability.
	metrics := observability.NewMetrics()

	// Initialize business layer with config-driven parameters.
	ingestionBiz := business.NewIngestionBusiness(evtsMan, pointRepo, cfg.IngestionBusinessConfig())
	areaBiz := business.NewAreaBusiness(evtsMan, areaRepo, stateRepo)
	geofenceBiz := business.NewGeofenceBusiness(
		evtsMan, areaRepo, stateRepo, geoEventRepo, metrics, cfg.GeofenceBusinessConfig(),
	)
	proximityBiz := business.NewProximityBusiness(
		latestPosRepo,
		areaRepo,
		cfg.ProximityBusinessConfig(),
	)
	trackBiz := business.NewTrackBusiness(
		pointRepo,
		geoEventRepo,
		stateRepo,
		cfg.TrackBusinessConfig(),
	)
	routeDeviationBiz := business.NewRouteDeviationBusiness(
		evtsMan, routeRepo, routeDeviationStateRepo, metrics,
		cfg.RouteDeviationBusinessConfig(),
	)
	routeBiz := business.NewRouteBusiness(
		evtsMan, routeRepo, routeAssignmentRepo, routeDeviationStateRepo,
	)
	retentionBiz := business.NewRetentionBusiness(dbPool, cfg.RetentionBusinessConfig())
	catchupBiz := business.NewCatchupBusiness(pointRepo, evtsMan, business.CatchupConfig{})

	// Start the retention scheduler as a background goroutine.
	// It runs immediately (including EnsurePartitions) and then every RetentionInterval.
	retentionCtx, cancelRetention := context.WithCancel(ctx)
	defer cancelRetention()
	go retentionBiz.StartScheduler(retentionCtx)

	catchupCtx, cancelCatchup := context.WithCancel(ctx)
	defer cancelCatchup()
	go catchupBiz.StartScheduler(catchupCtx)

	// Setup HTTP handler with authentication and rate limiting.
	sm := svc.SecurityManager()

	auth := sm.GetAuthorizer(ctx)
	sd := geolocationv1.File_geolocation_v1_geolocation_proto.Services().ByName("GeolocationService")
	functionChecker := authorizer.NewFunctionChecker(auth, permissions.ForService(sd).Namespace)

	geoServer := handlers.NewGeolocationServer(
		svc, functionChecker, ingestionBiz, areaBiz, routeBiz, proximityBiz, trackBiz, metrics,
		cfg.MaxRequestBodyBytes,
	)

	// Rate limiter middleware.
	rl := handlers.NewRateLimitMiddleware(cfg.RateLimitConfig())
	defer rl.Stop()

	connectHandler := setupConnectServer(ctx, sm, functionChecker, geoServer)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", geoServer.HealthCheck)
	mux.Handle("/openapi.yaml", handlers.OpenAPIHandler())
	mux.Handle("/", rl.Middleware(connectHandler))

	// Register event consumers and start service.
	svc.Init(ctx,
		frame.WithHTTPHandler(mux),
		frame.WithRegisterEvents(
			events.NewLocationPointConsumer(
				pointRepo, proximityBiz, geofenceBiz, routeDeviationBiz, metrics,
			),
			events.NewAreaChangeConsumer(areaBiz, stateRepo),
			events.NewRouteChangeConsumer(routeAssignmentRepo, routeDeviationStateRepo),
			events.NewGeoEventConsumer(),
			events.NewRouteDeviationConsumer(),
		),
	)

	if runErr := svc.Run(ctx, ""); runErr != nil {
		log.WithError(runErr).Fatal("could not run server")
	}
}

func setupConnectServer(
	ctx context.Context,
	sm security.Manager,
	functionChecker *authorizer.FunctionChecker,
	implementation *handlers.GeolocationServer,
) http.Handler {
	authenticator := sm.GetAuthenticator(ctx)
	auth := sm.GetAuthorizer(ctx)
	tenancyAccessChecker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)
	tenancyAccessInterceptor := connectInterceptors.NewTenancyAccessInterceptor(tenancyAccessChecker)

	// Build procedure map from proto annotations and exclude self-bypass RPCs.
	sd := geolocationv1.File_geolocation_v1_geolocation_proto.Services().ByName("GeolocationService")
	procMap := permissions.BuildProcedureMap(sd)

	// Exclude self-bypass RPCs from auto-enforcement.
	// These are checked manually in handlers with self-bypass logic.
	delete(procMap, "/geolocation.v1.GeolocationService/IngestLocations")
	delete(procMap, "/geolocation.v1.GeolocationService/GetSubjectRouteAssignments")
	delete(procMap, "/geolocation.v1.GeolocationService/GetTrack")
	delete(procMap, "/geolocation.v1.GeolocationService/GetSubjectEvents")
	delete(procMap, "/geolocation.v1.GeolocationService/GetNearbySubjects")

	functionAccessInterceptor := connectInterceptors.NewFunctionAccessInterceptor(functionChecker, procMap)

	defaultInterceptorList, err := connectInterceptors.DefaultList(
		ctx,
		authenticator,
		tenancyAccessInterceptor,
		functionAccessInterceptor,
	)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("main -- Could not create default interceptors")
	}

	_, handler := geolocationv1connect.NewGeolocationServiceHandler(
		implementation,
		connect.WithInterceptors(defaultInterceptorList...),
	)

	return handler
}
