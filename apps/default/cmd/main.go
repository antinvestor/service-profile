package main

import (
	"fmt"
	"net/http"
	"strings"

	"buf.build/go/protovalidate"
	apis "github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/handlers"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/queue"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pitabwire/frame"
)

func main() {
	serviceName := "service_profile"

	profileConfig, err := frame.ConfigFromEnv[config.ProfileConfig]()
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewService(serviceName, frame.WithConfig(&profileConfig))
	defer service.Stop(ctx)
	log := service.Log(ctx)

	serviceOptions := []frame.Option{frame.WithDatastore()}

	if profileConfig.DoDatabaseMigrate() {
		service.Init(ctx, serviceOptions...)

		err = service.MigrateDatastore(ctx, profileConfig.GetDatabaseMigrationPath(),
			&models.ProfileType{}, &models.Profile{}, &models.Contact{}, &models.Country{},
			&models.Address{}, &models.ProfileAddress{}, &models.Verification{},
			&models.VerificationAttempt{}, &models.RelationshipType{}, &models.Relationship{},
			&models.Device{}, &models.DeviceLog{}, &models.Roster{})

		if err != nil {
			log.WithError(err).Fatal("main -- Could not migrate successfully because : %+v", err)
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
		log.WithError(err).Fatal("main -- Could not setup notification service : %+v", err)
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

	grpcServerOpt := frame.WithGRPCServer(grpcServer)
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

	serviceOptions = append(serviceOptions, frame.WithHTTPHandler(proxyMux))

	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         service,
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: notificationCli,
	}

	verificationQueue := frame.WithRegisterSubscriber(
		profileConfig.QueueVerificationName,
		profileConfig.QueueVerification,
		&verificationQueueHandler,
	)
	verificationQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueVerificationName,
		profileConfig.QueueVerification,
	)

	deviceAnalysisQueueHandler := queue.DeviceAnalysisQueueHandler{
		Service:             service,
		DeviceRepository:    repository.NewDeviceRepository(service),
		DeviceLogRepository: repository.NewDeviceLogRepository(service),
	}

	deviceAnalysisQueue := frame.WithRegisterSubscriber(
		profileConfig.QueueDeviceAnalysisName,
		profileConfig.QueueDeviceAnalysis,
		&deviceAnalysisQueueHandler,
	)
	deviceAnalysisQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueDeviceAnalysisName,
		profileConfig.QueueDeviceAnalysis,
	)

	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueRelationshipConnectName,
		profileConfig.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueRelationshipDisConnectName,
		profileConfig.QueueRelationshipDisConnectURI,
	)

	serviceOptions = append(serviceOptions,
		verificationQueue, verificationQueuePublisher,
		deviceAnalysisQueue, deviceAnalysisQueuePublisher,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			&events.ClientConnectedSetupQueue{Service: service},
		))
	service.Init(ctx, serviceOptions...)

	log.WithField("server http port", profileConfig.HTTPServerPort).
		WithField("server grpc port", profileConfig.GrpcServerPort).
		Info(" Initiating server operations")
	defer implementation.Service.Stop(ctx)
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}
