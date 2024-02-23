package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/handlers"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/queue"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/bufbuild/protovalidate-go"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"strings"

	apis "github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
)

func main() {

	serviceName := "service_profile"

	var profileConfig config.ProfileConfig
	err := frame.ConfigProcess("", &profileConfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewService(serviceName, frame.Config(&profileConfig))
	defer service.Stop(ctx)
	log := service.L()

	serviceOptions := []frame.Option{frame.Datastore(ctx)}

	if profileConfig.DoDatabaseMigrate() {

		service.Init(serviceOptions...)

		err := service.MigrateDatastore(ctx, profileConfig.GetDatabaseMigrationPath(),
			&models.ProfileType{}, &models.Profile{}, &models.ContactType{},
			&models.CommunicationLevel{}, &models.Contact{}, &models.Country{},
			&models.Address{}, &models.ProfileAddress{}, &models.Verification{},
			&models.VerificationAttempt{}, &models.RelationshipType{}, &models.Relationship{})

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
			service.UnaryAuthInterceptor(jwtAudience, profileConfig.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.UnaryServerInterceptor(validator),
		),
		grpc.ChainStreamInterceptor(
			service.StreamAuthInterceptor(jwtAudience, profileConfig.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.StreamServerInterceptor(validator),
		),
	)

	implementation := &handlers.ProfileServer{
		Service:         service,
		NotificationCli: notificationCli,
		EncryptionKeyFunc: func() []byte {
			return pbkdf2.Key([]byte(profileConfig.ContactEncryptionKey),
				[]byte(profileConfig.ContactEncryptionSalt), 4096, 32, sha256.New)
		},
	}
	profilev1.RegisterProfileServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	profileServiceRestHandlers := service.AuthenticationMiddleware(implementation.NewRouterV1(),
		jwtAudience, profileConfig.Oauth2JwtVerifyIssuer)

	serviceOptions = append(serviceOptions, frame.HttpHandler(profileServiceRestHandlers))

	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         service,
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: notificationCli,
	}

	verificationQueue := frame.RegisterSubscriber(profileConfig.QueueVerificationName, profileConfig.QueueVerification, 2, &verificationQueueHandler)
	verificationQueuePublisher := frame.RegisterPublisher(profileConfig.QueueVerificationName, profileConfig.QueueVerification)

	serviceOptions = append(serviceOptions, verificationQueue, verificationQueuePublisher)

	service.Init(serviceOptions...)

	log.WithField("server http port", profileConfig.HttpServerPort).
		WithField("server grpc port", profileConfig.GrpcServerPort).
		Info(" Initiating server operations")
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}

}
