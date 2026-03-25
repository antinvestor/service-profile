package main

import (
	"context"
	"encoding/base64"
	"net/http"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	profilepb "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	apis "github.com/antinvestor/apis/go/common"
	"github.com/antinvestor/apis/go/notification"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/common/permissions"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security/authorizer"
	connectInterceptors "github.com/pitabwire/frame/security/interceptors/connect"
	securityhttp "github.com/pitabwire/frame/security/interceptors/httptor"
	"github.com/pitabwire/util"

	aconfig "github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/business"
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
		frame.WithDatastore(),
	)
	defer svc.Stop(ctx)
	log := svc.Log(ctx)

	dbManager := svc.DatastoreManager()

	// Handle database migration if requested
	if handleDatabaseMigration(ctx, dbManager, cfg) {
		return
	}

	// Setup clients and services
	notificationCli, nErr := setupNotificationClient(ctx, cfg)
	if nErr != nil {
		log.WithError(nErr).Fatal("main -- Could not setup notification svc")
	}

	dek, dekErr := decodeDEK(cfg)
	if dekErr != nil {
		log.WithError(dekErr).Fatal("main -- Could not decode DEK encryption keys")
	}

	// Seed default data (system bot contact) after migration creates the profile row
	seedDefaultData(ctx, svc, dek)

	// Setup Connect server
	connectHandler := setupConnectServer(ctx, svc, dek, notificationCli)

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
	cfg aconfig.ProfileConfig) (notificationv1connect.NotificationServiceClient, error) {
	return notification.NewClient(ctx, &cfg, apis.ServiceTarget{
		Endpoint:              cfg.NotificationSvcURI,
		WorkloadAPITargetPath: cfg.NotificationServiceWorkloadAPITargetPath,
		Audiences:             []string{"service_notification"},
	})
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

// seedDefaultData ensures seed data exists (e.g. the system bot contact).
// The bot profile row is created by SQL migration; this adds the encrypted contact.
func seedDefaultData(ctx context.Context, svc *frame.Service, dek *aconfig.DEK) {
	log := util.Log(ctx)
	cfg, _ := svc.Config().(*aconfig.ProfileConfig)
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	evtsMan := svc.EventsManager()

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)
	contactBiz := business.NewContactBusiness(ctx, cfg, dek, evtsMan, contactRepo, verificationRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBiz := business.NewAddressBusiness(ctx, addressRepo)
	profileBiz := business.NewProfileBusiness(ctx, cfg, dek, evtsMan, contactBiz, addressBiz, profileRepo)

	if err := business.SeedSystemBotContact(ctx, profileBiz, contactBiz); err != nil {
		log.WithError(err).Warn("failed to seed system bot contact — will retry on next startup")
	}
}

// setupConnectServer initializes and configures the gRPC server.
func setupConnectServer(ctx context.Context, svc *frame.Service, dek *aconfig.DEK,
	notificationCli notificationv1connect.NotificationServiceClient) http.Handler {
	securityMan := svc.SecurityManager()

	authenticator := securityMan.GetAuthenticator(ctx)

	auth := securityMan.GetAuthorizer(ctx)
	tenancyAccessChecker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)
	tenancyAccessInterceptor := connectInterceptors.NewTenancyAccessInterceptor(tenancyAccessChecker)

	// Build procedure map from proto annotations and exclude self-bypass RPCs.
	sd := profilepb.File_profile_v1_profile_proto.Services().ByName("ProfileService")
	procMap := permissions.BuildProcedureMap(sd)

	// Exclude self-bypass RPCs from auto-enforcement.
	// These are checked manually in handlers with self-bypass logic.
	delete(procMap, "/profile.v1.ProfileService/GetById")
	delete(procMap, "/profile.v1.ProfileService/Update")
	delete(procMap, "/profile.v1.ProfileService/AddAddress")
	delete(procMap, "/profile.v1.ProfileService/AddContact")
	delete(procMap, "/profile.v1.ProfileService/RemoveContact")
	delete(procMap, "/profile.v1.ProfileService/SearchRoster")
	delete(procMap, "/profile.v1.ProfileService/AddRelationship")
	delete(procMap, "/profile.v1.ProfileService/ListRelationships")

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

	implementation := handlers.NewProfileServer(ctx, svc, dek, notificationCli, functionChecker)

	_, serverHandler := profilev1connect.NewProfileServiceHandler(
		implementation, connect.WithInterceptors(defaultInterceptorList...))

	publicRestHandler := securityhttp.AuthenticationMiddleware(implementation.NewSecureRouterV1(), authenticator)

	mux := http.NewServeMux()
	mux.Handle("/", serverHandler)
	mux.Handle("/public/", http.StripPrefix("/public", publicRestHandler))
	mux.Handle("/openapi.yaml", apis.NewOpenAPIHandler(profilev1.ApiSpecFile, nil))

	return mux
}
