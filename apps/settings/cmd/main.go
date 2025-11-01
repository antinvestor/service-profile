package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/antinvestor/apis/go/settings/v1/settingsv1connect"
	aconfig "github.com/antinvestor/service-profile/apps/settings/config"
	"github.com/antinvestor/service-profile/apps/settings/service/handlers"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	securityconnect "github.com/pitabwire/frame/security/interceptors/connect"
	"github.com/pitabwire/util"
)

func main() {
	serviceName := "service_settings"
	ctx := context.Background()

	cfg, err := config.LoadWithOIDC[aconfig.SettingsConfig](ctx)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not process configs")
		return
	}

	ctx, svc := frame.NewServiceWithContext(ctx, serviceName, frame.WithConfig(&cfg), frame.WithRegisterServerOauth2Client())
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	// Handle database migration if requested
	if handleDatabaseMigration(ctx, svc, &cfg, log) {
		return
	}

	// Setup Connect server
	connectHandler := setupConnectServer(ctx, svc)

	// Setup HTTP handlers
	// Start with datastore option
	serviceOptions := []frame.Option{frame.WithDatastore(), frame.WithHTTPHandler(connectHandler)}

	svc.Init(ctx, serviceOptions...)

	log.WithField("server port", cfg.HTTPPort()).
		Info(" Initiating server operations")

	defer svc.Stop(ctx)
	err = svc.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}

// handleDatabaseMigration performs database migration if configured to do so.
func handleDatabaseMigration(
	ctx context.Context,
	svc *frame.Service,
	cfg *aconfig.SettingsConfig,
	log *util.LogEntry,
) bool {
	serviceOptions := []frame.Option{frame.WithDatastore()}

	if cfg.DoDatabaseMigrate() {
		svc.Init(ctx, serviceOptions...)

		dbMan := svc.DatastoreManager()
		dbPool := dbMan.GetPool(ctx, datastore.DefaultPoolName)
		err := repository.Migrate(ctx, dbMan, dbPool, cfg.GetDatabaseMigrationPath())
		if err != nil {
			log.WithError(err).Fatal("main -- Could not migrate successfully")
		}
		return true
	}
	return false
}

// setupConnectServer initializes and configures the gRPC server.
func setupConnectServer(ctx context.Context, svc *frame.Service) http.Handler {

	securityMan := svc.SecurityManager()

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not configure open telemetry")
	}

	validateInterceptor, err := securityconnect.NewValidationInterceptor()
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not configure validation interceptor")
	}

	authenticator := securityMan.GetAuthenticator(ctx)
	authInterceptor := securityconnect.NewAuthInterceptor(authenticator)

	implementation := handlers.NewSettingsServer(ctx, svc)

	_, serverHandler := settingsv1connect.NewSettingsServiceHandler(
		implementation, connect.WithInterceptors(authInterceptor, otelInterceptor, validateInterceptor))

	return serverHandler
}
