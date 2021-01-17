package main

import (
	"context"
	"fmt"
	"github.com/antinvestor/apis"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/models"
	"github.com/antinvestor/service-profile/service/handlers"
	"gocloud.dev/server"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"

	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/service"
)

func main() {

	serviceName := "Profile"

	ctx := context.Background()

	var serviceOptions []frame.Option

	datasource := frame.GetEnv(config.EnvDatabaseUrl, "postgres://ant:@nt@localhost/service_profile")
	mainDb := frame.Datastore(ctx, datasource, false)
	serviceOptions = append(serviceOptions, mainDb)

	readOnlydatasource := frame.GetEnv(config.EnvReplicaDatabaseUrl, datasource)
	readDb := frame.Datastore(ctx, readOnlydatasource, true)
	serviceOptions = append(serviceOptions, readDb)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(service.AuthInterceptor),
	)

	implementation := &handlers.ProfileServer{}

	papi.RegisterProfileServiceServer(grpcServer, implementation)

	httpOptions := &server.Options{}

	defaultServer := frame.GrpcServer(grpcServer, httpOptions)
	serviceOptions = append(serviceOptions, defaultServer)

	sysService := frame.NewService(serviceName, serviceOptions...)


	notificationServiceUrl := frame.GetEnv(config.EnvNotificationServiceUri, "127.0.0.1:7020")
	notificationCli, err := napi.NewNotificationClient(ctx, apis.WithEndpoint(notificationServiceUrl))
	if err != nil {
		log.Printf("main -- Could not setup notification service : %v", err)
	}

	implementation.Service = sysService
	implementation.NotificationCli = notificationCli



	isMigration, err := strconv.ParseBool(frame.GetEnv(config.EnvMigrate, "false"))
	if err != nil {
		isMigration = false
	}

	stdArgs := os.Args[1:]
	if (len(stdArgs) > 0 && stdArgs[0] == "migrate") || isMigration {

		migrationPath := frame.GetEnv(config.EnvMigrationPath, "./migrations/0001")
		err := sysService.MigrateDatastore(ctx, migrationPath,
			models.ProfileType{},
			models.Profile{}, models.ContactType{}, models.CommunicationLevel{},
			models.Contact{}, models.Country{}, &models.Address{}, models.ProfileAddress{},
			models.Verification{}, models.VerificationAttempt{})

		if err != nil {
			log.Printf("main -- Could not migrate successfully because : %v", err)
		}

	} else {

		serverPort := frame.GetEnv(config.EnvServerPort, "7020")

		log.Printf(" main -- Initiating server operations on : %s", serverPort)
		err := sysService.Run(ctx, fmt.Sprintf(":%v", serverPort))
		if err != nil {
			log.Printf("main -- Could not run Server : %v", err)
		}

	}

}
