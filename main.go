package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/antinvestor/apis"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/handlers"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/queue"
	"github.com/antinvestor/service-profile/service/repository"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"log"
	"strings"

	"os"
	"strconv"

	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/pitabwire/frame"

	"github.com/grpc-ecosystem/go-grpc-middleware"
)

func main() {

	serviceName := "service_profile"
	service := frame.NewService(serviceName)
	ctx := context.Background()

	var err error
	var serviceOptions []frame.Option

	datasource := frame.GetEnv(config.EnvDatabaseURL, "postgres://ant:@nt@localhost/service_profile")
	mainDb := frame.Datastore(ctx, datasource, false)
	serviceOptions = append(serviceOptions, mainDb)

	readOnlydatasource := frame.GetEnv(config.EnvReplicaDatabaseURL, datasource)
	readDb := frame.Datastore(ctx, readOnlydatasource, true)
	serviceOptions = append(serviceOptions, readDb)

	isMigration, err := strconv.ParseBool(frame.GetEnv(config.EnvMigrate, "false"))
	if err != nil {
		isMigration = false
	}

	stdArgs := os.Args[1:]
	if (len(stdArgs) > 0 && stdArgs[0] == "migrate") || isMigration {

		service.Init(serviceOptions...)

		migrationPath := frame.GetEnv(config.EnvMigrationPath, "./migrations/0001")
		err := service.MigrateDatastore(ctx, migrationPath,
			models.ProfileType{}, models.Profile{}, models.ContactType{},
			models.CommunicationLevel{}, models.Contact{}, models.Country{},
			&models.Address{}, models.ProfileAddress{}, models.Verification{},
			models.VerificationAttempt{})

		if err != nil {
			log.Fatalf("main -- Could not migrate successfully because : %+v", err)
		}

		return

	}

	oauth2ServiceHost := frame.GetEnv(config.EnvOauth2ServiceURI, "")
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)
	oauth2ServiceSecret := frame.GetEnv(config.EnvOauth2ServiceClientSecret, "")

	audienceList := make([]string, 0)
	oauth2ServiceAudience := frame.GetEnv(config.EnvOauth2ServiceAudience, "")
	if oauth2ServiceAudience != "" {
		audienceList = strings.Split(oauth2ServiceAudience, ",")
	}
	notificationServiceURL := frame.GetEnv(config.EnvNotificationServiceURI, "127.0.0.1:7020")
	notificationCli, err := napi.NewNotificationClient(ctx,
		apis.WithEndpoint(notificationServiceURL),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(serviceName),
		apis.WithTokenPassword(oauth2ServiceSecret),
		apis.WithAudiences(audienceList...))
	if err != nil {
		log.Fatalf("main -- Could not setup notification service : %+v", err)
	}

	jwtAudience := frame.GetEnv(config.EnvOauth2JwtVerifyAudience, serviceName)
	jwtIssuer := frame.GetEnv(config.EnvOauth2JwtVerifyIssuer, "")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			frame.UnaryAuthInterceptor(jwtAudience, jwtIssuer),
		)),
		grpc.StreamInterceptor(frame.StreamAuthInterceptor(jwtAudience, jwtIssuer)),
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
		SystemAccessID:  frame.GetEnv(config.EnvSystemAccessID, "c8cf0ldstmdlinc3eva0"),
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: notificationCli,
	}
	verificationQueueURL := frame.GetEnv(config.EnvQueueVerification, fmt.Sprintf("mem://%s", config.QueueVerificationName))
	verificationQueue := frame.RegisterSubscriber(config.QueueVerificationName, verificationQueueURL, 2, &verificationQueueHandler)
	verificationQueuePublisher := frame.RegisterPublisher(config.QueueVerificationName, verificationQueueURL)

	serviceOptions = append(serviceOptions, verificationQueue, verificationQueuePublisher)

	service.Init(serviceOptions...)

	encryptionKey := frame.GetEnv(config.EnvContactEncryptionKey, "")
	if encryptionKey == "" {
		err := errors.New("an encryption key has to be specified")
		log.Fatalf("main -- Could not start service because : %+v", err)
	}

	encryptionSalt := frame.GetEnv(config.EnvContactEncryptionSalt, "")
	if encryptionSalt == "" {
		err := errors.New("an encryption salt has to be specified")
		log.Fatalf("main -- Could not start service because : %+v", err)
	}

	contactEncryptionKey := pbkdf2.Key([]byte(encryptionKey), []byte(encryptionSalt), 4096, 32, sha256.New)
	implementation.EncryptionKey = contactEncryptionKey

	serverPort := frame.GetEnv(config.EnvServerPort, "7005")

	log.Printf(" main -- Initiating server operations on : %s", serverPort)
	err = implementation.Service.Run(ctx, fmt.Sprintf(":%v", serverPort))
	if err != nil {
		log.Fatalf("main -- Could not run Server : %+v", err)
	}

}
