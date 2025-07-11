package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"buf.build/go/protovalidate"
	apis "github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/handlers"
	"github.com/antinvestor/service-profile/apps/default/service/queue"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pitabwire/frame"
)

func main() {
	serviceName := "service_profile"

	cfg, err := frame.ConfigFromEnv[config.ProfileConfig]()
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

	oauth2ServiceHost := cfg.GetOauth2ServiceURI()
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)

	audienceList := make([]string, 0)
	oauth2ServiceAudience := cfg.Oauth2ServiceAudience
	if oauth2ServiceAudience != "" {
		audienceList = strings.Split(oauth2ServiceAudience, ",")
	}

	notificationCli, err := notificationv1.NewNotificationClient(ctx,
		apis.WithEndpoint(cfg.NotificationServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(svc.JwtClientID()),
		apis.WithTokenPassword(cfg.Oauth2ServiceClientSecret),
		apis.WithAudiences(audienceList...))
	if err != nil {
		log.WithError(err).Fatal("main -- Could not setup notification svc : %+v", err)
	}

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

	grpcServerOpt := frame.WithGRPCServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	proxyOptions := apis.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost:%s", cfg.GrpcServerPort),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := profilev1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		log.WithError(err).Fatal("could not create the proxy handler")
		return
	}

	profileServiceRestHandlers := svc.AuthenticationMiddleware(
		implementation.NewSecureRouterV1(), jwtAudience, cfg.Oauth2JwtVerifyIssuer)

	proxyMux.Handle("/public/", http.StripPrefix("/public", profileServiceRestHandlers))

	serviceOptions = append(serviceOptions, frame.WithHTTPHandler(proxyMux))

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

	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipConnectName,
		cfg.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipDisConnectName,
		cfg.QueueRelationshipDisConnectURI,
	)

	serviceOptions = append(serviceOptions,
		verificationQueue, verificationQueuePublisher,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			&events.ClientConnectedSetupQueue{Service: svc},
		))
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
