package main

import (
	"context"
	"net/http"

	"buf.build/gen/go/antinvestor/device/connectrpc/go/device/v1/devicev1connect"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	connectInterceptors "github.com/pitabwire/frame/security/interceptors/connect"
	securityhttp "github.com/pitabwire/frame/security/interceptors/httptor"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
	"github.com/antinvestor/service-profile/apps/devices/service/handlers"
	"github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

func main() {
	serviceName := "service_devices"

	ctx := context.Background()
	cfg, err := config.LoadWithOIDC[aconfig.DevicesConfig](ctx)
	if err != nil {
		util.Log(ctx).With("err", err).Error("could not process configs")
		return
	}

	if cfg.Name() == "" {
		cfg.ServiceName = serviceName
	}

	ctx, svc := frame.NewServiceWithContext(
		ctx,
		frame.WithConfig(&cfg),
		frame.WithRegisterServerOauth2Client(),
		frame.WithDatastore(),
		frame.WithCacheManager(),
		frame.WithInMemoryCache(aconfig.CacheNameDevices),
		frame.WithInMemoryCache(aconfig.CacheNamePresence),
		frame.WithInMemoryCache(aconfig.CacheNameGeoIP),
		frame.WithInMemoryCache(aconfig.CacheNameRate),
	)
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	if cfg.DoDatabaseMigrate() {
		err = repository.Migrate(ctx, svc.DatastoreManager(), cfg.GetDatabaseMigrationPath())
		if err != nil {
			log.WithError(err).Fatal("main -- Could not migrate successfully because : %+v", err)
		}

		return
	}

	serviceOptions := initServiceComponents(ctx, svc, &cfg)
	svc.Init(ctx, serviceOptions...)

	// Start service.
	err = svc.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}

// initServiceComponents initializes repositories, business layer, handlers, and queue subscriptions.
func initServiceComponents(
	ctx context.Context,
	svc *frame.Service,
	cfg *aconfig.DevicesConfig,
) []frame.Option {
	securityMan := svc.SecurityManager()
	queueMan := svc.QueueManager()
	workMan := svc.WorkManager()
	httpClientMan := svc.HTTPClientManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	// Initialize cache service.
	cacheSvc := caching.NewDeviceCacheService(svc.CacheManager())

	// Initialize repositories.
	deviceLogRepo := repository.NewDeviceLogRepository(ctx, dbPool, workMan)
	deviceSessionRepo := repository.NewDeviceSessionRepository(ctx, dbPool, workMan)
	deviceRepo := repository.NewDeviceRepository(ctx, dbPool, workMan)
	deviceKeyRepo := repository.NewDeviceKeyRepository(ctx, dbPool, workMan)
	devicePresenceRepo := repository.NewDevicePresenceRepository(ctx, dbPool, workMan)

	// Initialize business layer with cache.
	deviceBusiness := business.NewDeviceBusiness(
		ctx, cfg, queueMan, workMan, deviceRepo, deviceLogRepo, deviceSessionRepo, cacheSvc,
	)
	keyBusiness := business.NewKeysBusiness(ctx, cfg, queueMan, workMan, deviceRepo, deviceKeyRepo, cacheSvc)
	presenceBusiness := business.NewPresenceBusiness(
		ctx, cfg, queueMan, workMan, deviceRepo, devicePresenceRepo, cacheSvc,
	)
	notifyBusiness, err := business.NewNotifyBusiness(ctx, cfg, queueMan, workMan, keyBusiness, deviceRepo)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not configure device server")
	}

	// TURN provider is optional — nil turnBusiness is handled gracefully by the handler.
	var turnBiz business.TURNBusiness
	turnBiz, err = business.NewTURNBusiness(cfg, httpClientMan)
	if err != nil {
		util.Log(ctx).WithError(err).Warn("TURN credentials provider not configured, endpoint will return 503")
	}

	implementation := handlers.NewDeviceServer(ctx, deviceBusiness, presenceBusiness, keyBusiness, notifyBusiness)
	connectHandler := setupConnectServer(ctx, securityMan, implementation, turnBiz, cacheSvc, cfg)

	analysisHandler := queue.NewDeviceAnalysisQueueHandler(
		httpClientMan, deviceRepo, deviceLogRepo, deviceSessionRepo, cacheSvc,
	)

	return []frame.Option{
		frame.WithHTTPHandler(connectHandler),
		frame.WithRegisterSubscriber(cfg.QueueDeviceAnalysisName, cfg.QueueDeviceAnalysis, analysisHandler),
		frame.WithRegisterPublisher(cfg.QueueDeviceAnalysisName, cfg.QueueDeviceAnalysis),
	}
}

// setupConnectServer initializes and configures the connect server.
func setupConnectServer(
	ctx context.Context,
	securityMan security.Manager,
	implementation *handlers.DevicesServer,
	turnBusiness business.TURNBusiness,
	cacheSvc *caching.DeviceCacheService,
	cfg *aconfig.DevicesConfig,
) http.Handler {
	authenticator := securityMan.GetAuthenticator(ctx)

	defaultInterceptorList, err := connectInterceptors.DefaultList(ctx, authenticator)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("main -- Could not create default interceptors")
	}

	_, serverHandler := devicev1connect.NewDeviceServiceHandler(
		implementation, connect.WithInterceptors(defaultInterceptorList...))

	turnHandler := handlers.NewTURNHandler(turnBusiness)
	turnRouter := http.NewServeMux()
	turnRouter.HandleFunc("POST /v1/turn/credentials",
		handlers.RateLimitTURN(turnHandler.GetTurnCredentials, cacheSvc, cfg.RateLimitTURNPerMinute))
	secureTurnRouter := securityhttp.AuthenticationMiddleware(turnRouter, authenticator)

	mux := http.NewServeMux()
	mux.Handle("/", serverHandler)
	mux.Handle("/v1/turn/", secureTurnRouter)

	return mux
}
