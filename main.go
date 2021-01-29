package main

import (
	"context"
	"fmt"
	"github.com/antinvestor/apis"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/handlers"
	"github.com/antinvestor/service-profile/service/models"
	grpc_log "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"

	"os"
	"strconv"

	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/pitabwire/frame"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
)


func main() {


	serviceName := "Profile"

	ctx := context.Background()

	logrusEntry := log.NewEntry(log.New())
	grpc_log.ReplaceGrpcLogger(logrusEntry)

	var err error
	var serviceOptions []frame.Option

	datasource := frame.GetEnv(config.EnvDatabaseUrl, "postgres://ant:@nt@localhost/service_profile")
	mainDb := frame.Datastore(ctx, datasource, false)
	serviceOptions = append(serviceOptions, mainDb)

	readOnlydatasource := frame.GetEnv(config.EnvReplicaDatabaseUrl, datasource)
	readDb := frame.Datastore(ctx, readOnlydatasource, true)
	serviceOptions = append(serviceOptions, readDb)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_log.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	implementation := &handlers.ProfileServer{}

	papi.RegisterProfileServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	implementation.Service = frame.NewService(serviceName, serviceOptions...)

	notificationServiceUrl := frame.GetEnv(config.EnvNotificationServiceUri, "127.0.0.1:7020")
	implementation.NotificationCli, err = napi.NewNotificationClient(ctx, apis.WithEndpoint(notificationServiceUrl))
	if err != nil {
		log.Printf("main -- Could not setup notification service : %v", err)
	}


	isMigration, err := strconv.ParseBool(frame.GetEnv(config.EnvMigrate, "false"))
	if err != nil {
		isMigration = false
	}

	stdArgs := os.Args[1:]
	if (len(stdArgs) > 0 && stdArgs[0] == "migrate") || isMigration {

		migrationPath := frame.GetEnv(config.EnvMigrationPath, "./migrations/0001")
		err := implementation.Service.MigrateDatastore(ctx, migrationPath,
			models.ProfileType{},
			models.Profile{}, models.ContactType{}, models.CommunicationLevel{},
			models.Contact{}, models.Country{}, &models.Address{}, models.ProfileAddress{},
			models.Verification{}, models.VerificationAttempt{})

		if err != nil {
			log.Printf("main -- Could not migrate successfully because : %v", err)
		}

	} else {

		serverPort := frame.GetEnv(config.EnvServerPort, "7005")

		log.Printf(" main -- Initiating server operations on : %s", serverPort)
		err := implementation.Service.Run(ctx, fmt.Sprintf(":%v", serverPort))
		if err != nil {
			log.Printf("main -- Could not run Server : %v", err)
		}

	}

}
