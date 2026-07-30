package main

import (
	"context"
	"net/http"
	"strings"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"connectrpc.com/connect"
	apis "github.com/antinvestor/common/v2"
	"github.com/antinvestor/common/v2/connection"
	"github.com/antinvestor/common/v2/permissions"
	"github.com/antinvestor/common/v2/servicecatalog"
	"github.com/pitabwire/frame/v2"
	"github.com/pitabwire/frame/v2/config"
	"github.com/pitabwire/frame/v2/security/authorizer"
	connectInterceptors "github.com/pitabwire/frame/v2/security/interceptors/connect"
	"github.com/pitabwire/frame/v2/setup"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/chatagent/config"
	"github.com/antinvestor/service-profile/apps/chatagent/service/authz"
	"github.com/antinvestor/service-profile/apps/chatagent/service/handlers"
	"github.com/antinvestor/service-profile/apps/chatagent/service/repository"
	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
	"github.com/antinvestor/service-profile/gen/go/chatagent/v1/chatagentv1connect"
)

func main() {
	serviceName := "service_chat_agent"
	ctx := context.Background()

	cfg, err := config.LoadWithOIDC[aconfig.ChatAgentConfig](ctx)
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

	chatSD := chatagentv1.File_chatagent_v1_chatagent_proto.Services().ByName("ChatAgentService")
	svc.Setup().RegisterFunc(setup.NameMigrate, func(ctx context.Context) error {
		return repository.Migrate(ctx, svc.DatastoreManager(), cfg.GetDatabaseMigrationPath())
	})

	if frame.ShouldRunSetup(&cfg) {
		svc.Init(ctx, frame.WithPermissionRegistration(chatSD))
		if setupErr := svc.RunSetupForProcess(ctx, &cfg); setupErr != nil {
			log.WithError(setupErr).Fatal("setup plan failed")
		}
		log.Info("setup plan complete — exiting")
		return
	}

	connectHandler := setupConnectServer(ctx, svc, &cfg)

	serviceOptions := []frame.Option{
		frame.WithHTTPHandler(connectHandler),
		frame.WithPermissionRegistration(chatSD),
	}

	svc.Init(ctx, serviceOptions...)

	if runErr := svc.Run(ctx, ""); runErr != nil {
		log.WithError(runErr).Fatal("could not run Server")
	}
}

func setupConnectServer(ctx context.Context, svc *frame.Service, cfg *aconfig.ChatAgentConfig) http.Handler {
	securityMan := svc.SecurityManager()
	authenticator := securityMan.GetAuthenticator(ctx)
	auth := securityMan.GetAuthorizer(ctx)

	tenancyAccessChecker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)
	tenancyAccessInterceptor := connectInterceptors.NewTenancyAccessInterceptor(tenancyAccessChecker)

	sd := chatagentv1.File_chatagent_v1_chatagent_proto.Services().ByName("ChatAgentService")
	procMap := permissions.BuildProcedureMap(sd)

	functionChecker := authorizer.NewFunctionChecker(auth, permissions.ForService(sd).Namespace)
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

	notificationCli := setupNotificationClient(ctx, cfg)

	implementation := handlers.NewChatAgentServer(ctx, svc, handlers.ServerDeps{
		LLM: handlers.LLMConfig{
			BaseURL: cfg.InferenceBaseURL,
			APIKey:  cfg.InferenceAPIKey,
			Model:   cfg.InferenceModel,
		},
		NotificationClient: notificationCli,
	})

	_, serverHandler := chatagentv1connect.NewChatAgentServiceHandler(
		implementation, connect.WithInterceptors(defaultInterceptorList...))

	return serverHandler
}

// setupNotificationClient creates the Notification service client for omnichannel reply delivery.
// Returns nil when NOTIFICATION_SERVICE_URI is empty or client setup fails (web-only mode).
func setupNotificationClient(
	ctx context.Context,
	cfg *aconfig.ChatAgentConfig,
) notificationv1connect.NotificationServiceClient {
	if strings.TrimSpace(cfg.NotificationSvcURI) == "" {
		return nil
	}
	cli, err := connection.NewServiceClient(ctx, cfg, apis.ServiceTarget{
		Endpoint:              cfg.NotificationSvcURI,
		WorkloadAPITargetPath: cfg.NotificationServiceWorkloadAPITargetPath,
		ServiceID:             servicecatalog.ServiceNotification,
	}, notificationv1connect.NewNotificationServiceClient)
	if err != nil {
		util.Log(ctx).WithError(err).Warn("chatagent: notification client unavailable; omnichannel delivery disabled")
		return nil
	}
	return cli
}
