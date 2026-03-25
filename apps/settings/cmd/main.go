package main

import (
	"context"
	"net/http"

	"buf.build/gen/go/antinvestor/settingz/connectrpc/go/settings/v1/settingsv1connect"
	settingspb "buf.build/gen/go/antinvestor/settingz/protocolbuffers/go/settings/v1"
	"connectrpc.com/connect"
	"github.com/antinvestor/common/permissions"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security/authorizer"
	connectInterceptors "github.com/pitabwire/frame/security/interceptors/connect"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/settings/config"
	"github.com/antinvestor/service-profile/apps/settings/service/authz"
	"github.com/antinvestor/service-profile/apps/settings/service/handlers"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

func main() {
	serviceName := "service_setting"
	ctx := context.Background()

	cfg, err := config.LoadWithOIDC[aconfig.SettingsConfig](ctx)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not process configs")
		return
	}

	if cfg.Name() == "" {
		cfg.ServiceName = serviceName
	}

	ctx, svc := frame.NewServiceWithContext(
		ctx,
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
	)
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	// Handle database migration if requested
	if handleDatabaseMigration(ctx, svc.DatastoreManager(), &cfg) {
		return
	}

	// Setup Connect server
	connectHandler := setupConnectServer(ctx, svc)

	// Setup HTTP handlers
	// Start with datastore option
	serviceOptions := []frame.Option{frame.WithHTTPHandler(connectHandler)}

	svc.Init(ctx, serviceOptions...)

	// Start service
	err = svc.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server ")
	}
}

// handleDatabaseMigration performs database migration if configured to do so.
func handleDatabaseMigration(
	ctx context.Context,
	dbManager datastore.Manager,
	cfg *aconfig.SettingsConfig,
) bool {
	if cfg.DoDatabaseMigrate() {
		err := repository.Migrate(ctx, dbManager, cfg.GetDatabaseMigrationPath())
		if err != nil {
			util.Log(ctx).WithError(err).Fatal("main -- Could not migrate successfully")
		}
		return true
	}
	return false
}

// setupConnectServer initializes and configures the gRPC server.
func setupConnectServer(ctx context.Context, svc *frame.Service) http.Handler {
	securityMan := svc.SecurityManager()

	authenticator := securityMan.GetAuthenticator(ctx)

	auth := securityMan.GetAuthorizer(ctx)
	tenancyAccessChecker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)
	tenancyAccessInterceptor := connectInterceptors.NewTenancyAccessInterceptor(tenancyAccessChecker)

	// Build procedure map from proto annotations — no self-bypass RPCs in settings.
	sd := settingspb.File_settings_v1_settings_proto.Services().ByName("SettingsService")
	procMap := permissions.BuildProcedureMap(sd)

	functionChecker := authorizer.NewFunctionChecker(auth, "service_profile")
	functionAccessInterceptor := connectInterceptors.NewFunctionAccessInterceptor(functionChecker, procMap)

	defaultInterceptorList, err := connectInterceptors.DefaultList(
		ctx,
		authenticator,
		tenancyAccessInterceptor,
		functionAccessInterceptor,
	)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("main -- Could not create default interceptors")
	}

	implementation := handlers.NewSettingsServer(ctx, svc)

	_, serverHandler := settingsv1connect.NewSettingsServiceHandler(
		implementation, connect.WithInterceptors(defaultInterceptorList...))

	return serverHandler
}
