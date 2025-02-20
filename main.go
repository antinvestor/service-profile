package main

import (
	"fmt"
	apis "github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/events"
	"github.com/antinvestor/service-profile/service/handlers"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/queue"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/bufbuild/protovalidate-go"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strings"
)

func main() {

	serviceName := "service_profile"

	profileConfig, err := frame.ConfigFromEnv[config.ProfileConfig]()
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewService(serviceName, frame.Config(&profileConfig))
	defer service.Stop(ctx)
	log := service.L(ctx)

	serviceOptions := []frame.Option{frame.Datastore(ctx)}

	if profileConfig.DoDatabaseMigrate() {

		service.Init(serviceOptions...)

		err = service.DB(ctx, false).Exec(`
			CREATE EXTENSION IF NOT EXISTS pg_search;
			CREATE EXTENSION IF NOT EXISTS pg_analytics;
			CREATE EXTENSION IF NOT EXISTS pg_ivm;
			CREATE EXTENSION IF NOT EXISTS vector;
			CREATE EXTENSION IF NOT EXISTS postgis;
			CREATE EXTENSION IF NOT EXISTS postgis_topology;
			CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
			CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;
		`).Error
		if err != nil {
			log.Fatalf("main -- Failed to create extensions: %v", err)
		}

		err := service.MigrateDatastore(ctx, profileConfig.GetDatabaseMigrationPath(),
			&models.ProfileType{}, &models.Profile{}, &models.Contact{}, &models.Country{},
			&models.Address{}, &models.ProfileAddress{}, &models.Verification{},
			&models.VerificationAttempt{}, &models.RelationshipType{}, &models.Relationship{},
			&models.Device{}, &models.DeviceLog{}, &models.Roster{})

		if err != nil {
			log.Fatalf("main -- Could not migrate successfully because : %+v", err)
		}

		return
	}

	err = service.RegisterForJwt(ctx)
	if err != nil {
		log.WithError(err).Fatal("main -- could not register fo jwt")
	}

	oauth2ServiceHost := profileConfig.GetOauth2ServiceURI()
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)

	audienceList := make([]string, 0)
	oauth2ServiceAudience := profileConfig.Oauth2ServiceAudience
	if oauth2ServiceAudience != "" {
		audienceList = strings.Split(oauth2ServiceAudience, ",")
	}

	notificationCli, err := notificationv1.NewNotificationClient(ctx,
		apis.WithEndpoint(profileConfig.NotificationServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(service.JwtClientID()),
		apis.WithTokenPassword(profileConfig.Oauth2ServiceClientSecret),
		apis.WithAudiences(audienceList...))
	if err != nil {
		log.Fatalf("main -- Could not setup notification service : %+v", err)
	}

	jwtAudience := profileConfig.Oauth2JwtVerifyAudience
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
			service.UnaryAuthInterceptor(jwtAudience, profileConfig.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.UnaryServerInterceptor(validator)),

		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(frame.RecoveryHandlerFun)),
			service.StreamAuthInterceptor(jwtAudience, profileConfig.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.StreamServerInterceptor(validator),
		),
	)

	implementation := &handlers.ProfileServer{
		Service:         service,
		NotificationCli: notificationCli,
	}
	profilev1.RegisterProfileServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	proxyOptions := apis.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost:%s", profileConfig.GrpcServerPort),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := profilev1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		log.WithError(err).Fatal("could not create the proxy handler")
		return
	}

	profileServiceRestHandlers := service.AuthenticationMiddleware(
		implementation.NewSecureRouterV1(), jwtAudience, profileConfig.Oauth2JwtVerifyIssuer)

	proxyMux.Handle("/public/", http.StripPrefix("/public", profileServiceRestHandlers))
	proxyMux.Handle("/_public/", http.StripPrefix("/_public", implementation.NewInSecureRouterV1()))

	serviceOptions = append(serviceOptions, frame.HttpHandler(proxyMux))

	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         service,
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: notificationCli,
	}

	verificationQueue := frame.RegisterSubscriber(profileConfig.QueueVerificationName, profileConfig.QueueVerification, 2, &verificationQueueHandler)
	verificationQueuePublisher := frame.RegisterPublisher(profileConfig.QueueVerificationName, profileConfig.QueueVerification)

	deviceAnalysisQueueHandler := queue.DeviceAnalysisQueueHandler{
		Service:             service,
		DeviceRepository:    repository.NewDeviceRepository(service),
		DeviceLogRepository: repository.NewDeviceLogRepository(service),
	}

	deviceAnalysisQueue := frame.RegisterSubscriber(profileConfig.QueueDeviceAnalysisName, profileConfig.QueueDeviceAnalysis, 2, &deviceAnalysisQueueHandler)
	deviceAnalysisQueuePublisher := frame.RegisterPublisher(profileConfig.QueueDeviceAnalysisName, profileConfig.QueueDeviceAnalysis)

	relationshipConnectQueuePublisher := frame.RegisterPublisher(profileConfig.QueueRelationshipConnectName, profileConfig.QueueRelationshipConnectURI)
	relationshipDisConnectQueuePublisher := frame.RegisterPublisher(profileConfig.QueueRelationshipDisConnectName, profileConfig.QueueRelationshipDisConnectURI)

	serviceOptions = append(serviceOptions,
		verificationQueue, verificationQueuePublisher,
		deviceAnalysisQueue, deviceAnalysisQueuePublisher,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.RegisterEvents(
			&events.ClientConnectedSetupQueue{Service: service},
		))
	service.Init(serviceOptions...)

	log.WithField("server http port", profileConfig.HttpServerPort).
		WithField("server grpc port", profileConfig.GrpcServerPort).
		Info(" Initiating server operations")
	defer implementation.Service.Stop(ctx)
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}

}
