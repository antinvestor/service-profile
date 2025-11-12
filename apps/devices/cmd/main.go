package main

import (
	"net/http"

	"buf.build/gen/go/antinvestor/device/connectrpc/go/device/v1/devicev1connect"
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	securityconnect "github.com/pitabwire/frame/security/interceptors/connect"
	"github.com/pitabwire/util"
	"golang.org/x/net/context"

	aconfig "github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
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
	)
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	serviceOptions := []frame.Option{frame.WithDatastore()}

	if cfg.DoDatabaseMigrate() {
		svc.Init(ctx, serviceOptions...)

		err = repository.Migrate(ctx, svc.DatastoreManager(), cfg.GetDatabaseMigrationPath())
		if err != nil {
			log.WithError(err).Fatal("main -- Could not migrate successfully because : %+v", err)
		}

		return
	}

	securityMan := svc.SecurityManager()
	queueMan := svc.QueueManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	deviceLogRepo := repository.NewDeviceLogRepository(ctx, dbPool, workMan)
	deviceSessionRepo := repository.NewDeviceSessionRepository(ctx, dbPool, workMan)
	deviceRepo := repository.NewDeviceRepository(ctx, dbPool, workMan)
	deviceKeyRepo := repository.NewDeviceKeyRepository(ctx, dbPool, workMan)
	devicePresenceRepo := repository.NewDevicePresenceRepository(ctx, dbPool, workMan)

	deviceBusiness := business.NewDeviceBusiness(
		ctx,
		&cfg,
		queueMan,
		workMan,
		deviceRepo,
		deviceLogRepo,
		deviceSessionRepo,
	)
	keyBusiness := business.NewKeysBusiness(ctx, &cfg, queueMan, workMan, deviceRepo, deviceKeyRepo)
	presenceBusiness := business.NewPresenceBusiness(ctx, &cfg, queueMan, workMan, deviceRepo, devicePresenceRepo)
	notifyBusiness, err := business.NewNotifyBusiness(ctx, &cfg, queueMan, workMan, keyBusiness, deviceRepo)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not configure device server")
	}

	implementation := handlers.NewDeviceServer(ctx, deviceBusiness, presenceBusiness, keyBusiness, notifyBusiness)

	// Setup Connect server
	connectHandler := setupConnectServer(ctx, securityMan, implementation)

	// Setup HTTP handlers
	// Start with datastore option
	serviceOptions = []frame.Option{frame.WithDatastore(), frame.WithHTTPHandler(connectHandler)}

	deviceAnalysisQueue := frame.WithRegisterSubscriber(
		cfg.QueueDeviceAnalysisName,
		cfg.QueueDeviceAnalysis,
		queue.NewDeviceAnalysisQueueHandler(svc.HTTPClientManager(), deviceRepo, deviceLogRepo, deviceSessionRepo),
	)
	deviceAnalysisQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueDeviceAnalysisName,
		cfg.QueueDeviceAnalysis,
	)

	serviceOptions = append(serviceOptions,
		deviceAnalysisQueue, deviceAnalysisQueuePublisher,
	)
	svc.Init(ctx, serviceOptions...)

	log.
		WithField("server port", cfg.HTTPPort()).
		Info(" Initiating server operations")
	defer svc.Stop(ctx)
	err = svc.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}

// setupConnectServer initializes and configures the connect server.
func setupConnectServer(
	ctx context.Context,
	securityMan security.Manager,
	implementation *handlers.DevicesServer,
) http.Handler {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not configure open telemetry")
	}

	validateInterceptor, err := securityconnect.NewValidationInterceptor()
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not configure validation interceptor")
	}

	authenticator := securityMan.GetAuthenticator(ctx)
	authInterceptor := securityconnect.NewAuthInterceptor(authenticator)

	_, serverHandler := devicev1connect.NewDeviceServiceHandler(
		implementation, connect.WithInterceptors(authInterceptor, otelInterceptor, validateInterceptor))

	return serverHandler
}
