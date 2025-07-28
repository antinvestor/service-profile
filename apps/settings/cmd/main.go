package main

import (
	"context"

	"buf.build/go/protovalidate"
	settingsV1 "github.com/antinvestor/apis/go/settings/v1"
	"github.com/antinvestor/service-profile/apps/settings/config"
	"github.com/antinvestor/service-profile/apps/settings/service/handlers"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
	"google.golang.org/grpc"
)

func main() {
	serviceName := "service_settings"
	ctx := context.Background()

	cfg, err := frame.ConfigFromEnv[config.SettingsConfig]()
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not process configs")
		return
	}

	ctx, svc := frame.NewServiceWithContext(ctx, serviceName, frame.WithConfig(&cfg))
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	serviceOptions := []frame.Option{frame.WithDatastore()}

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

	jwtAudience := cfg.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	validator, err := protovalidate.New()
	if err != nil {
		log.WithError(err).Fatal("could not load validator for proto messages")
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		protovalidateinterceptor.UnaryServerInterceptor(validator),
		recovery.UnaryServerInterceptor(),
	}

	if cfg.SecurelyRunService {
		unaryInterceptors = append(
			[]grpc.UnaryServerInterceptor{svc.UnaryAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer)},
			unaryInterceptors...)
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		protovalidateinterceptor.StreamServerInterceptor(validator),
		recovery.StreamServerInterceptor(),
	}

	if cfg.SecurelyRunService {
		streamInterceptors = append(
			[]grpc.StreamServerInterceptor{svc.StreamAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer)},
			streamInterceptors...)
	} else {
		log.Warn("svc is running insecurely: secure by setting SECURELY_RUN_SERVICE=True")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	implementation := &handlers.SettingsServer{
		Service: svc,
	}

	settingsV1.RegisterSettingsServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.WithGRPCServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	svc.Init(ctx, serviceOptions...)

	log.WithField("server http port", cfg.HTTPServerPort).
		WithField("server grpc port", cfg.GrpcServerPort).
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
