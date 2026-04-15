package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"net/http"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	profilepb "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	apis "github.com/antinvestor/common"
	"github.com/antinvestor/common/audit"
	"github.com/antinvestor/common/connection"
	"github.com/antinvestor/common/permissions"
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

//go:embed spec/profile.openapi.yaml
var profileAPISpecFile []byte

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

	dek, dekErr := decodeDEK(cfg)
	if dekErr != nil {
		log.WithError(dekErr).Fatal("main -- Could not decode DEK encryption keys")
	}

	// Handle database migration if requested.
	// The seed MUST run after migration to create encrypted contacts for
	// bootstrap profiles. Without contacts, login creates duplicate profiles.
	// The migration job must have DEK env vars configured:
	//   DEK_ACTIVE_KEY_ID, DEK_AES256GCM_KEY, DEK_LOOKUP_TOKEN_HMACSHA256_KEY
	if handleDatabaseMigration(ctx, dbManager, cfg) {
		seedDefaultData(ctx, svc, dek)
		return
	}

	// Setup clients and services
	notificationCli, nErr := setupNotificationClient(ctx, cfg)
	if nErr != nil {
		log.WithError(nErr).Fatal("main -- Could not setup notification svc")
	}

	// Setup Connect server
	connectHandler := setupConnectServer(ctx, svc, dek, notificationCli)

	// Register permission manifest for the profile service namespace.
	profileSD := profilepb.File_profile_v1_profile_proto.Services().ByName("ProfileService")

	// Setup HTTP handlers
	serviceOptions := []frame.Option{
		frame.WithHTTPHandler(connectHandler),
		frame.WithPermissionRegistration(profileSD),
	}

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
	return connection.NewServiceClient(ctx, &cfg, apis.ServiceTarget{
		Endpoint:              cfg.NotificationSvcURI,
		WorkloadAPITargetPath: cfg.NotificationServiceWorkloadAPITargetPath,
		Audiences:             []string{"service_notification"},
	}, notificationv1connect.NewNotificationServiceClient)
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

// seedDefaultData ensures bootstrap profiles have their encrypted contacts.
// Profile rows are created by SQL migration; this function adds contacts
// since they require application-level encryption (DEK).
//
// This runs as part of the migration job and is critical — if it fails,
// bootstrap profiles will have no contacts and user logins will create
// duplicate profiles. The migration job MUST have these env vars:
//
//	DEK_ACTIVE_KEY_ID              — identifies the active encryption key
//	DEK_AES256GCM_KEY              — AES-256-GCM key for encrypting contact details
//	DEK_LOOKUP_TOKEN_HMACSHA256_KEY — HMAC key for deterministic contact lookup tokens
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

	if err := business.SeedBootstrapContacts(ctx, profileBiz, contactBiz); err != nil {
		log.WithError(err).Fatal("failed to seed bootstrap contacts — migration incomplete, check DEK env vars")
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

	functionChecker := authorizer.NewFunctionChecker(auth, permissions.ForService(sd).Namespace)
	functionAccessInterceptor := connectInterceptors.NewFunctionAccessInterceptor(functionChecker, procMap)

	// Audit interceptor — sends entries to audit service if configured.
	profileCfg, _ := svc.Config().(*aconfig.ProfileConfig)
	var auditInterceptor connect.Interceptor
	if profileCfg != nil && profileCfg.AuditServiceURI != "" {
		auditCli, auditErr := connection.NewServiceClient(ctx, profileCfg, apis.ServiceTarget{
			Endpoint:  profileCfg.AuditServiceURI,
			Audiences: []string{"service_audit"},
		}, audit.NewConnectClient)
		if auditErr != nil {
			util.Log(ctx).WithError(auditErr).Warn("audit client not available — audit entries will only be logged")
			auditInterceptor = audit.NewInterceptor("service_profile", nil)
		} else {
			auditInterceptor = audit.NewInterceptor("service_profile", auditCli)
		}
	} else {
		auditInterceptor = audit.NewInterceptor("service_profile", nil)
	}

	defaultInterceptorList, err := connectInterceptors.DefaultList(
		ctx,
		authenticator,
		tenancyAccessInterceptor,
		functionAccessInterceptor,
		auditInterceptor,
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
	mux.Handle("/openapi.yaml", apis.NewOpenAPIHandler(profileAPISpecFile, nil))

	return mux
}
