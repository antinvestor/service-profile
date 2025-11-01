package main

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	apis "github.com/antinvestor/apis/go/common"
	settingsV1 "github.com/antinvestor/apis/go/settings/v1"
	"github.com/antinvestor/service-profile/apps/settings/config"
	"github.com/antinvestor/service-profile/apps/settings/service/handlers"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serviceName := "service_settings"
	ctx := context.Background()

	cfg, err := frame.ConfigLoadWithOIDC[config.SettingsConfig](ctx)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not process configs")
		return
	}

	ctx, svc := frame.NewServiceWithContext(ctx, serviceName, frame.WithConfig(&cfg))
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	// Handle database migration if requested
	if handleDatabaseMigration(ctx, svc, cfg, log) {
		return
	}

	if cfg.SecurelyRunService {
		err = svc.RegisterForJwt(ctx)
		if err != nil {
			log.WithError(err).Fatal("main -- could not register fo jwt")
		}
	}

	// Setup GRPC server
	grpcServer, implementation := setupGRPCServer(ctx, svc, cfg, log)

	// Setup HTTP handlers and proxy
	serviceOptions, httpErr := setupHTTPHandlers(ctx, cfg, grpcServer)
	if err != nil {
		log.WithError(httpErr).Fatal("could not setup HTTP handlers")
	}

	svc.Init(ctx, serviceOptions...)

	log.WithField("server http port", cfg.HTTPPort()).
		WithField("server grpc port", cfg.GrpcPort()).
		Info(" Initiating server operations")

	defer implementation.Service.Stop(ctx)
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}

// handleDatabaseMigration performs database migration if configured to do so.
func handleDatabaseMigration(
	ctx context.Context,
	svc *frame.Service,
	cfg config.SettingsConfig,
	log *util.LogEntry,
) bool {
	serviceOptions := []frame.Option{frame.WithDatastore()}

	if cfg.DoDatabaseMigrate() {
		svc.Init(ctx, serviceOptions...)

		err := repository.Migrate(ctx, svc, cfg.GetDatabaseMigrationPath())
		if err != nil {
			log.WithError(err).Fatal("main -- Could not migrate successfully")
		}
		return true
	}
	return false
}

// setupGRPCServer initializes and configures the gRPC server.
func setupGRPCServer(_ context.Context, svc *frame.Service,
	cfg config.SettingsConfig, log *util.LogEntry) (*grpc.Server, *handlers.SettingsServer) {
	jwtAudience := cfg.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = svc.Name()
	}

	validator, err := protovalidate.New()
	if err != nil {
		log.WithError(err).Fatal("could not load validator for proto messages")
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			svc.UnaryAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.UnaryServerInterceptor(validator)),

		grpc.ChainStreamInterceptor(
			svc.StreamAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.StreamServerInterceptor(validator),
		),
	)

	implementation := &handlers.SettingsServer{
		Service: svc,
	}

	settingsV1.RegisterSettingsServiceServer(grpcServer, implementation)

	return grpcServer, implementation
}

// setupHTTPHandlers configures HTTP handlers and proxy.
func setupHTTPHandlers(
	ctx context.Context,
	cfg config.SettingsConfig,
	grpcServer *grpc.Server,
) ([]frame.Option, error) {
	// Start with datastore option
	serviceOptions := []frame.Option{frame.WithDatastore()}

	// Add GRPC server option
	grpcServerOpt := frame.WithGRPCServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	// Setup proxy
	proxyOptions := apis.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost%s", cfg.GrpcPort()),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := settingsV1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		return nil, err
	}
	serviceOptions = append(serviceOptions, frame.WithHTTPHandler(proxyMux))

	return serviceOptions, nil
}
