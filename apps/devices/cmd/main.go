package main

import (
	"log/slog"
	"net/http"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/handlers"
	"github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	repository2 "github.com/antinvestor/service-profile/apps/devices/service/repository"

	"github.com/pitabwire/frame"
)

func main() {
	serviceName := "service_devices"

	cfg, err := frame.ConfigFromEnv[config.DevicesConfig]()
	if err != nil {
		slog.With("err", err).Error("could not process configs")
		return
	}

	ctx, svc := frame.NewService(serviceName, frame.WithConfig(&cfg))
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	serviceOptions := []frame.Option{frame.WithDatastore()}

	if cfg.DoDatabaseMigrate() {
		svc.Init(ctx, serviceOptions...)

		err = repository.Migrate(ctx, svc, cfg.GetDatabaseMigrationPath())
		if err != nil {
			log.WithError(err).Fatal("main -- Could not migrate successfully because : %+v", err)
		}

		return
	}

	err = svc.RegisterForJwt(ctx)
	if err != nil {
		log.WithError(err).Fatal("main -- could not register fo jwt")
	}

	jwtAudience := cfg.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	implementation := &handlers.DevicesServer{
		Service: svc,
	}

	profileServiceRestHandlers := svc.AuthenticationMiddleware(
		implementation.NewSecureRouterV1(), jwtAudience, cfg.Oauth2JwtVerifyIssuer)

	proxyMux := http.NewServeMux()

	proxyMux.Handle("/public/", http.StripPrefix("/public", profileServiceRestHandlers))
	proxyMux.Handle("/_public/", http.StripPrefix("/_public", implementation.NewInSecureRouterV1()))

	serviceOptions = append(serviceOptions, frame.WithHTTPHandler(proxyMux))

	deviceAnalysisQueueHandler := queue.DeviceAnalysisQueueHandler{
		Service:             svc,
		DeviceRepository:    repository2.NewDeviceRepository(svc),
		DeviceLogRepository: repository2.NewDeviceLogRepository(svc),
	}

	deviceAnalysisQueue := frame.WithRegisterSubscriber(
		cfg.QueueDeviceAnalysisName,
		cfg.QueueDeviceAnalysis,
		&deviceAnalysisQueueHandler,
	)
	deviceAnalysisQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueDeviceAnalysisName,
		cfg.QueueDeviceAnalysis,
	)

	serviceOptions = append(serviceOptions,
		deviceAnalysisQueue, deviceAnalysisQueuePublisher,
	)
	svc.Init(ctx, serviceOptions...)

	log.WithField("server http port", cfg.HTTPServerPort).
		Info(" Initiating server operations")
	defer implementation.Service.Stop(ctx)
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}
