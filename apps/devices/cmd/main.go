package main

import (
	"fmt"
	"net/http"

	"buf.build/go/protovalidate"
	apis "github.com/antinvestor/apis/go/common"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/handlers"
	"github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

func main() {
	serviceName := "service_devices"

	ctx := context.Background()
	cfg, err := frame.ConfigLoadWithOIDC[config.DevicesConfig](ctx)
	if err != nil {
		util.Log(ctx).With("err", err).Error("could not process configs")
		return
	}

	ctx, svc := frame.NewServiceWithContext(ctx, serviceName, frame.WithConfig(&cfg))
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

	// Setup GRPC server
	grpcServer, implementation := setupGRPCServer(ctx, svc, cfg, serviceName, log)

	// Setup HTTP handlers and proxy
	serviceOptions, httpErr := setupHTTPHandlers(ctx, implementation, cfg, grpcServer)
	if err != nil {
		log.WithError(httpErr).Fatal("could not setup HTTP handlers")
	}

	deviceAnalysisQueueHandler := queue.DeviceAnalysisQueueHandler{
		Service:             svc,
		DeviceRepository:    repository.NewDeviceRepository(svc),
		DeviceLogRepository: repository.NewDeviceLogRepository(svc),
		SessionRepository:   repository.NewDeviceSessionRepository(svc),
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

// setupGRPCServer initializes and configures the gRPC server.
func setupGRPCServer(ctx context.Context, svc *frame.Service,
	cfg config.DevicesConfig,
	serviceName string,
	log *util.LogEntry) (*grpc.Server, *handlers.DevicesServer) {
	jwtAudience := cfg.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	validator, err := protovalidate.New()
	if err != nil {
		log.WithError(err).Fatal("could not load validator for proto messages")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(frame.RecoveryHandlerFun)),
			svc.UnaryAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.UnaryServerInterceptor(validator)),

		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(frame.RecoveryHandlerFun)),
			svc.StreamAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.StreamServerInterceptor(validator),
		),
	)

	implementation := handlers.NewDeviceServer(ctx, svc)
	devicev1.RegisterDeviceServiceServer(grpcServer, implementation)

	return grpcServer, implementation
}

// setupHTTPHandlers configures HTTP handlers and proxy.
func setupHTTPHandlers(
	ctx context.Context,
	implementation *handlers.DevicesServer,
	cfg config.DevicesConfig, grpcServer *grpc.Server,
) ([]frame.Option, error) {
	// Start with datastore option
	serviceOptions := []frame.Option{frame.WithDatastore()}

	// Add GRPC server option
	grpcServerOpt := frame.WithGRPCServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	// Setup proxy
	proxyOptions := apis.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost:%s", cfg.GrpcPort()),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := devicev1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		return nil, err
	}

	// Setup REST handlers
	proxyMux.Handle("/_public/", http.StripPrefix("/_public", implementation.NewInSecureRouterV1()))
	serviceOptions = append(serviceOptions, frame.WithHTTPHandler(proxyMux))

	return serviceOptions, nil
}
