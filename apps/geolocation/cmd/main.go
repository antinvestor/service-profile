package main

import (
	"context"
	"net/http"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	securityhttp "github.com/pitabwire/frame/security/interceptors/httptor"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/geolocation/config"
	"github.com/antinvestor/service-profile/apps/geolocation/service/authz"
	"github.com/antinvestor/service-profile/apps/geolocation/service/business"
	"github.com/antinvestor/service-profile/apps/geolocation/service/events"
	"github.com/antinvestor/service-profile/apps/geolocation/service/handlers"
	"github.com/antinvestor/service-profile/apps/geolocation/service/observability"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
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
		frame.WithRegisterServerOauth2Client(),
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
	stateRepo := repository.NewGeofenceStateRepository(dbPool)
	latestPosRepo := repository.NewLatestPositionRepository(dbPool)
	routeRepo := repository.NewRouteRepository(ctx, dbPool, workMan)
	routeAssignmentRepo := repository.NewRouteAssignmentRepository(ctx, dbPool, workMan)
	routeDeviationStateRepo := repository.NewRouteDeviationStateRepository(dbPool)

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
	catchupBiz := business.NewCatchupBusiness(dbPool, evtsMan, business.CatchupConfig{})

	// Start the retention scheduler as a background goroutine.
	// It runs immediately (including EnsurePartitions) and then every RetentionInterval.
	retentionCtx, cancelRetention := context.WithCancel(ctx)
	defer cancelRetention()
	go retentionBiz.StartScheduler(retentionCtx)

	// Run catch-up at startup to re-emit events for any location points that were
	// persisted but never processed (e.g., due to a crash between INSERT and event emission).
	if catchupErr := catchupBiz.RunCatchup(ctx); catchupErr != nil {
		log.WithError(catchupErr).Warn("startup catchup failed")
	}

	// Setup HTTP handler with authentication and rate limiting.
	sm := svc.SecurityManager()
	authzMiddleware := authz.NewMiddleware(sm.GetAuthorizer(ctx))

	geoServer := handlers.NewGeolocationServer(
		svc, authzMiddleware, ingestionBiz, areaBiz, routeBiz, proximityBiz, trackBiz, metrics,
		cfg.MaxRequestBodyBytes,
	)

	// Rate limiter middleware.
	rl := handlers.NewRateLimitMiddleware(cfg.RateLimitConfig())
	defer rl.Stop()

	// Health check and OpenAPI are unauthenticated; all other routes require authentication.
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("GET /healthz", geoServer.HealthCheck)
	healthMux.HandleFunc("GET /openapi.yaml", handlers.OpenAPIHandler())

	authenticatedRouter := authenticateRouter(ctx, sm, geoServer.NewRouter())

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthMux)
	mux.Handle("/openapi.yaml", healthMux)
	mux.Handle("/", rl.Middleware(authenticatedRouter))

	// Register event consumers and start service.
	svc.Init(ctx,
		frame.WithHTTPHandler(mux),
		frame.WithRegisterEvents(
			events.NewLocationPointConsumer(
				proximityBiz, geofenceBiz, routeDeviationBiz, metrics,
			),
			events.NewAreaChangeConsumer(areaBiz, stateRepo),
			events.NewGeoEventConsumer(),
			events.NewRouteDeviationConsumer(),
		),
	)

	if runErr := svc.Run(ctx, ""); runErr != nil {
		log.WithError(runErr).Fatal("could not run server")
	}
}

// authenticateRouter wraps the given handler with OAuth2 authentication middleware.
func authenticateRouter(
	ctx context.Context,
	sm security.Manager,
	handler http.Handler,
) http.Handler {
	authenticator := sm.GetAuthenticator(ctx)
	return securityhttp.AuthenticationMiddleware(handler, authenticator)
}
