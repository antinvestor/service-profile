package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"buf.build/go/protovalidate"
	apis "github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/handlers"
	"github.com/antinvestor/service-profile/apps/default/service/queue"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

func main() {
	serviceName := "service_profile"

	// Initialize configuration
	cfg, err := frame.ConfigFromEnv[config.ProfileConfig]()
	if err != nil {
		slog.With("err", err).Error("could not process configs")
		return
	}

	// Create service
	ctx, svc := frame.NewService(serviceName, frame.WithConfig(&cfg))
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	// Handle database migration if requested
	if handleDatabaseMigration(ctx, svc, cfg, log) {
		return
	}

	// Register for JWT
	err = svc.RegisterForJwt(ctx)
	if err != nil {
		log.WithError(err).Fatal("main -- could not register for jwt")
	}

	// Setup clients and services
	notificationCli, nErr := setupNotificationClient(ctx, svc, cfg)
	if err != nil {
		log.WithError(nErr).Fatal("main -- Could not setup notification svc")
	}

	// Setup GRPC server
	grpcServer, implementation := setupGRPCServer(ctx, svc, notificationCli, cfg, serviceName, log)

	// Setup HTTP handlers and proxy
	serviceOptions, httpErr := setupHTTPHandlers(ctx, svc, implementation, cfg, grpcServer)
	if err != nil {
		log.WithError(httpErr).Fatal("could not setup HTTP handlers")
	}

	// Register queue handlers
	serviceOptions = registerQueueHandlers(svc, serviceOptions, cfg, notificationCli)

	// Initialize the service with all options
	svc.Init(ctx, serviceOptions...)

	// Start the service
	log.WithField("server http port", cfg.HTTPServerPort).
		WithField("server grpc port", cfg.GrpcServerPort).
		Info(" Initiating server operations")

	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server")
	}
}

// handleDatabaseMigration performs database migration if configured to do so.
func handleDatabaseMigration(
	ctx context.Context,
	svc *frame.Service,
	cfg config.ProfileConfig,
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

// setupNotificationClient creates and configures the notification client.
func setupNotificationClient(
	ctx context.Context,
	svc *frame.Service,
	cfg config.ProfileConfig,
) (*notificationv1.NotificationClient, error) {
	oauth2ServiceHost := cfg.GetOauth2ServiceURI()
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)

	audienceList := make([]string, 0)
	oauth2ServiceAudience := cfg.Oauth2ServiceAudience
	if oauth2ServiceAudience != "" {
		audienceList = strings.Split(oauth2ServiceAudience, ",")
	}

	return notificationv1.NewNotificationClient(ctx,
		apis.WithEndpoint(cfg.NotificationServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(svc.JwtClientID()),
		apis.WithTokenPassword(cfg.Oauth2ServiceClientSecret),
		apis.WithAudiences(audienceList...))
}

// setupGRPCServer initializes and configures the gRPC server.
func setupGRPCServer(_ context.Context, svc *frame.Service,
	notificationCli *notificationv1.NotificationClient,
	cfg config.ProfileConfig,
	serviceName string,
	log *util.LogEntry) (*grpc.Server, *handlers.ProfileServer) {
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

	implementation := &handlers.ProfileServer{
		Service:         svc,
		NotificationCli: notificationCli,
	}
	profilev1.RegisterProfileServiceServer(grpcServer, implementation)

	return grpcServer, implementation
}

// setupHTTPHandlers configures HTTP handlers and proxy.
func setupHTTPHandlers(
	ctx context.Context,
	svc *frame.Service,
	implementation *handlers.ProfileServer,
	cfg config.ProfileConfig,
	grpcServer *grpc.Server,
) ([]frame.Option, error) {
	// Start with datastore option
	serviceOptions := []frame.Option{frame.WithDatastore()}

	// Add GRPC server option
	grpcServerOpt := frame.WithGRPCServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	// Setup proxy
	proxyOptions := apis.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost:%s", cfg.GrpcServerPort),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := profilev1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		return nil, err
	}

	// Setup REST handlers
	jwtAudience := cfg.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = "service_profile"
	}

	profileServiceRestHandlers := svc.AuthenticationMiddleware(
		implementation.NewSecureRouterV1(), jwtAudience, cfg.Oauth2JwtVerifyIssuer)

	proxyMux.Handle("/public/", http.StripPrefix("/public", profileServiceRestHandlers))
	serviceOptions = append(serviceOptions, frame.WithHTTPHandler(proxyMux))

	return serviceOptions, nil
}

// registerQueueHandlers sets up all message queue handlers and publishers.
func registerQueueHandlers(
	svc *frame.Service,
	serviceOptions []frame.Option,
	cfg config.ProfileConfig,
	notificationCli *notificationv1.NotificationClient,
) []frame.Option {
	// Setup verification queue
	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         svc,
		ContactRepo:     repository.NewContactRepository(svc),
		NotificationCli: notificationCli,
	}

	verificationQueue := frame.WithRegisterSubscriber(
		cfg.QueueVerificationName,
		cfg.QueueVerification,
		&verificationQueueHandler,
	)
	verificationQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueVerificationName,
		cfg.QueueVerification,
	)

	// Setup relationship queues
	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipConnectName,
		cfg.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipDisConnectName,
		cfg.QueueRelationshipDisConnectURI,
	)

	// Combine all options
	return append(serviceOptions,
		verificationQueue, verificationQueuePublisher,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			&events.ClientConnectedSetupQueue{Service: svc},
		))
}
