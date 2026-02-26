package main

import (
	"context"
	"encoding/base64"
	"net/http"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	"connectrpc.com/connect"
	apis "github.com/antinvestor/apis/go/common"
	"github.com/antinvestor/apis/go/notification"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	connectInterceptors "github.com/pitabwire/frame/security/interceptors/connect"
	securityhttp "github.com/pitabwire/frame/security/interceptors/httptor"
	"github.com/pitabwire/frame/security/openid"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/handlers"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

func main() {
	ctx := context.Background()

	// Initialize configuration
	cfg, err := config.LoadWithOIDC[aconfig.ProfileConfig](ctx)
	if err != nil {
		util.Log(ctx).With("err", err).Error("could not process configs")
		return
	}

	if cfg.Name() == "" {
		cfg.ServiceName = "service_profile"
	}

	// Create service
	ctx, svc := frame.NewServiceWithContext(
		ctx,
		frame.WithConfig(&cfg),
		frame.WithRegisterServerOauth2Client(),
		frame.WithDatastore(),
	)
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	sm := svc.SecurityManager()
	dbManager := svc.DatastoreManager()

	// Handle database migration if requested
	if handleDatabaseMigration(ctx, dbManager, cfg) {
		return
	}

	// Setup clients and services
	notificationCli, nErr := setupNotificationClient(ctx, sm, cfg)
	if nErr != nil {
		log.WithError(nErr).Fatal("main -- Could not setup notification svc")
	}

	dek, dekErr := decodeDEK(cfg)
	if dekErr != nil {
		log.WithError(dekErr).Fatal("main -- Could not decode DEK encryption keys")
	}

	// Setup authz middleware
	authzMiddleware := authz.NewMiddleware(sm.GetAuthorizer(ctx))

	// Setup Connect server
	connectHandler := setupConnectServer(ctx, svc, dek, notificationCli, authzMiddleware)

	// Setup HTTP handlers
	// Start with datastore option
	serviceOptions := []frame.Option{frame.WithHTTPHandler(connectHandler)}

	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipConnectName,
		cfg.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipDisConnectName,
		cfg.QueueRelationshipDisConnectURI,
	)

	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	evtsMan := svc.EventsManager()
	qMan := svc.QueueManager()

	contactRepository := repository.NewContactRepository(ctx, dbPool, workMan)

	// Register queue handlers
	serviceOptions = append(serviceOptions,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			events.NewClientConnectedSetupQueue(
				ctx,
				&cfg,
				qMan,
				evtsMan,
				repository.NewRelationshipRepository(ctx, dbPool, workMan),
			),
			events.NewContactVerificationQueue(
				&cfg,
				contactRepository,
				repository.NewVerificationRepository(ctx, dbPool, workMan),
				notificationCli,
			),
			events.NewContactVerificationAttemptedQueue(
				contactRepository,
				repository.NewVerificationRepository(ctx, dbPool, workMan),
			),
			events.NewContactKeyRotationQueue(
				&cfg, dek, contactRepository,
			),
		))

	// Initialize the service with all options
	svc.Init(ctx, serviceOptions...)

	// Start the service
	err = svc.Run(ctx, "")
	if err != nil {
		log.WithError(err).Fatal("could not run Server")
	}
}

// handleDatabaseMigration performs database migration if configured to do so.
func handleDatabaseMigration(
	ctx context.Context,
	dbManager datastore.Manager,
	cfg aconfig.ProfileConfig,
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

// setupNotificationClient creates and configures the notification client.
func setupNotificationClient(
	ctx context.Context,
	clHolder security.InternalOauth2ClientHolder,
	cfg aconfig.ProfileConfig) (notificationv1connect.NotificationServiceClient, error) {
	return notification.NewClient(ctx,
		apis.WithEndpoint(cfg.NotificationServiceURI),
		apis.WithTokenEndpoint(cfg.GetOauth2TokenEndpoint()),
		apis.WithTokenUsername(clHolder.JwtClientID()),
		apis.WithTokenPassword(clHolder.JwtClientSecret()),
		apis.WithScopes(openid.ConstSystemScopeInternal),
		apis.WithAudiences("service_notifications"))
}

// decodeDEK decodes the Base64-encoded encryption keys from config into raw bytes.
func decodeDEK(cfg aconfig.ProfileConfig) (*aconfig.DEK, error) {
	key, err := base64.StdEncoding.DecodeString(cfg.DEKActiveAES256GCMKey)
	if err != nil {
		return nil, err
	}

	lookupKey, err := base64.StdEncoding.DecodeString(cfg.DEKLookupTokenHMACSHA256Key)
	if err != nil {
		return nil, err
	}

	var oldKey []byte
	if cfg.DEKOldAES256GCMKey != "" {
		oldKey, err = base64.StdEncoding.DecodeString(cfg.DEKOldAES256GCMKey)
		if err != nil {
			return nil, err
		}
	}

	return &aconfig.DEK{
		KeyID:     cfg.DEKActiveKeyID,
		Key:       key,
		OldKey:    oldKey,
		LookUpKey: lookupKey,
	}, nil
}

// setupConnectServer initializes and configures the gRPC server.
func setupConnectServer(ctx context.Context, svc *frame.Service, dek *aconfig.DEK,
	notificationCli notificationv1connect.NotificationServiceClient, authzMiddleware authz.Middleware) http.Handler {
	securityMan := svc.SecurityManager()

	authenticator := securityMan.GetAuthenticator(ctx)

	defaultInterceptorList, err := connectInterceptors.DefaultList(ctx, authenticator)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("main -- Could not create default interceptors")
	}

	implementation := handlers.NewProfileServer(ctx, svc, dek, notificationCli, authzMiddleware)

	_, serverHandler := profilev1connect.NewProfileServiceHandler(
		implementation, connect.WithInterceptors(defaultInterceptorList...))

	publicRestHandler := securityhttp.AuthenticationMiddleware(implementation.NewSecureRouterV1(), authenticator)

	mux := http.NewServeMux()
	mux.Handle("/", serverHandler)
	mux.Handle("/public/", http.StripPrefix("/public", publicRestHandler))
	mux.Handle("/openapi.yaml", apis.NewOpenAPIHandler(profilev1.ApiSpecFile, nil))

	return mux
}
