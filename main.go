package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/antinvestor/apis"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/handlers"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/queue"
	"github.com/antinvestor/service-profile/service/repository"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"strings"

	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/pitabwire/frame"

	"github.com/grpc-ecosystem/go-grpc-middleware"
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
			models.ProfileType{}, models.Profile{}, models.ContactType{},
			models.CommunicationLevel{}, models.Contact{}, models.Country{},
			&models.Address{}, models.ProfileAddress{}, models.Verification{},
			models.VerificationAttempt{})

		if err != nil {
			log.Fatalf("main -- Could not migrate successfully because : %+v", err)
		}

		return
	}

	oauth2ServiceHost := profileConfig.GetOauth2ServiceURI()
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)

	audienceList := make([]string, 0)
	oauth2ServiceAudience := profileConfig.Oauth2ServiceAudience
	if oauth2ServiceAudience != "" {
		audienceList = strings.Split(oauth2ServiceAudience, ",")
	}

	notificationCli, err := napi.NewNotificationClient(ctx,
		apis.WithEndpoint(profileConfig.NotificationServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(serviceName),
		apis.WithTokenPassword(profileConfig.Oauth2ServiceClientSecret),
		apis.WithAudiences(audienceList...))
	if err != nil {
		log.Fatalf("main -- Could not setup notification service : %+v", err)
	}

	jwtAudience := profileConfig.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			service.UnaryAuthInterceptor(jwtAudience, profileConfig.Oauth2JwtVerifyIssuer),
		)),
		grpc.StreamInterceptor(service.StreamAuthInterceptor(jwtAudience, profileConfig.Oauth2JwtVerifyIssuer)),
	)

	implementation := &handlers.ProfileServer{
		Service:         service,
		NotificationCli: notificationCli,
	}
	papi.RegisterProfileServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         service,
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: notificationCli,
	}

	verificationQueue := frame.RegisterSubscriber(profileConfig.QueueVerificationName, profileConfig.QueueVerification, 2, &verificationQueueHandler)
	verificationQueuePublisher := frame.RegisterPublisher(profileConfig.QueueVerificationName, profileConfig.QueueVerification)

	serviceOptions = append(serviceOptions, verificationQueue, verificationQueuePublisher)

	service.Init(serviceOptions...)

	contactEncryptionKey := pbkdf2.Key([]byte(profileConfig.ContactEncryptionKey),
		[]byte(profileConfig.ContactEncryptionSalt), 4096, 32, sha256.New)
	implementation.EncryptionKey = contactEncryptionKey

	log.WithField("server http port", profileConfig.HttpServerPort).
		WithField("server grpc port", profileConfig.GrpcServerPort).
		Info(" Initiating server operations")
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}

}
