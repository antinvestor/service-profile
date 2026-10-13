package tests

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/v2/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connectinprocess"
	"github.com/pitabwire/frame/v3"
	"github.com/pitabwire/frame/v3/config"
	"github.com/pitabwire/frame/v3/datastore"
	"github.com/pitabwire/frame/v3/datastore/pool"
	"github.com/pitabwire/frame/v3/frametests"
	"github.com/pitabwire/frame/v3/frametests/definition"
	"github.com/pitabwire/frame/v3/frametests/deps/testpostgres"
	"github.com/pitabwire/frame/v3/frametests/rlstest"
	"github.com/pitabwire/frame/v3/security"
	"github.com/pitabwire/frame/v3/security/authorizer"
	"github.com/pitabwire/frame/v3/tenancy"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	aconfig "github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests/testketo"
	"github.com/antinvestor/service-profile/internal/outbox"
	"github.com/antinvestor/service-profile/internal/outboxtest"
	"github.com/antinvestor/service-profile/internal/rlsadmin"
)

const PostgresqlDBImage = "postgres:latest"

const (
	DefaultRandomStringLength = 8
)

type ProfileBaseTestSuite struct {
	frametests.FrameBaseTestSuite

	FunctionChecker *authorizer.FunctionChecker
	ketoReadURI     string
	ketoWriteURI    string

	// Facts records the domain facts the outbox relay published, and
	// FactRelay drains the outbox on demand rather than on the background
	// ticker.
	Facts     *outboxtest.Recorder
	FactRelay *outbox.Relay

	ContactRepo      repository.ContactRepository
	VerificationRepo repository.VerificationRepository
	AddressRepo      repository.AddressRepository
	ProfileRepo      repository.ProfileRepository
	RosterRepo       repository.RosterRepository
	RelationshipRepo repository.RelationshipRepository
}

func initResources(_ context.Context) []definition.TestResource {
	pg := testpostgres.NewWithOpts("service_profile", definition.WithUserName("ant"))

	keto := testketo.NewWithOpts(
		definition.WithDependancies(pg),
		definition.WithEnableLogging(true),
	)

	return []definition.TestResource{pg, keto}
}

func (bs *ProfileBaseTestSuite) SetupSuite() {
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()

	ctx := bs.T().Context()

	// Find Keto dependency and extract read/write URIs
	var ketoDep definition.DependancyConn
	for _, res := range bs.Resources() {
		if res.Name() == testketo.ImageName {
			ketoDep = res
			break
		}
	}
	bs.Require().NotNil(ketoDep, "keto dependency should be available")

	// Write API: default port (4467/tcp, first in port list)
	writeURL, err := url.Parse(string(ketoDep.GetDS(ctx)))
	bs.Require().NoError(err)
	bs.ketoWriteURI = writeURL.Host

	// Read API: port 4466/tcp (second in port list)
	readPort, err := ketoDep.PortMapping(ctx, "4466/tcp")
	bs.Require().NoError(err)
	bs.ketoReadURI = fmt.Sprintf("%s:%s", writeURL.Hostname(), readPort)
}

func (bs *ProfileBaseTestSuite) CreateService(
	t *testing.T,
	depOpts *definition.DependencyOption,
) (context.Context, *frame.Service) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")

	ctx := t.Context()
	cfg, err := config.FromEnv[aconfig.ProfileConfig]()
	require.NoError(t, err)

	cfg.LogLevel = "debug"
	cfg.RunServiceSecurely = false
	cfg.ServerPort = ""
	cfg.DatabaseMigrate = true
	cfg.DatabaseTraceQueries = true
	cfg.QueueProfileEventsURI = "mem://profile.events"

	res := depOpts.ByIsDatabase(ctx)
	testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
	require.NoError(t, err0)

	t.Cleanup(func() {
		cleanup(t.Context())
	})

	cfg.DatabasePrimaryURL = []string{testDS.String()}
	cfg.DatabaseReplicaURL = []string{testDS.String()}

	// Configure real Keto authoriser URIs
	cfg.AuthorizationServiceReadURI = bs.ketoReadURI
	cfg.AuthorizationServiceWriteURI = bs.ketoWriteURI

	// The postgres testcontainer user is a superuser which bypasses RLS,
	// so application queries are dropped to an unprivileged role to
	// actually exercise the tenancy-isolation policies installed during
	// migration.
	require.NoError(t, rlstest.CreateRole(ctx, testDS.String()))
	rlsProv := rlstest.New()

	ctx, svc := frame.NewServiceWithContext(t.Context(), frame.WithName("profile tests"),
		frame.WithConfig(&cfg),
		frame.WithTenancyProvider(rlsProv),
		frame.WithDatastore(pool.WithTraceConfig(&cfg)),
		frametests.WithNoopDriver())
	t.Cleanup(func() { svc.Stop(ctx) })

	// Wire real Keto authoriser via SecurityManager
	sm := svc.SecurityManager()
	bs.FunctionChecker = authorizer.NewFunctionChecker(sm.GetAuthorizer(ctx), "service_profile")

	bs.Facts = &outboxtest.Recorder{}

	profileFactsPublisher := frame.WithRegisterPublisher(
		cfg.QueueProfileEventsName,
		cfg.QueueProfileEventsURI,
	)
	profileFactsSubscriber := frame.WithRegisterSubscriber(
		cfg.QueueProfileEventsName,
		cfg.QueueProfileEventsURI,
		bs.Facts,
	)

	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipConnectName,
		cfg.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipDisConnectName,
		cfg.QueueRelationshipDisConnectURI,
	)

	evtsMan := svc.EventsManager()
	qMan := svc.QueueManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)
	relationshipRepo := repository.NewRelationshipRepository(ctx, dbPool, workMan)

	bs.FactRelay = outbox.NewRelay(dbPool, qMan, cfg.QueueProfileEventsName)

	svc.Init(ctx,
		profileFactsPublisher, profileFactsSubscriber,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			events.NewClientConnectedSetupQueue(ctx, &cfg, qMan, evtsMan, relationshipRepo),
			events.NewContactVerificationQueue(&cfg, contactRepo, verificationRepo, bs.GetNotificationCli(t)),
			events.NewContactVerificationAttemptedQueue(contactRepo, verificationRepo),
		),
	)

	err = repository.Migrate(ctx, svc.DatastoreManager(), "../../migrations/0001")
	require.NoError(t, err)

	require.NoError(t, rlstest.GrantAll(ctx, testDS.String()))
	require.NoError(t, rlsadmin.GrantOwnership(ctx, testDS.String()))
	rlsProv.Enable()

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return ctx, svc
}

// noopNotificationHandler accepts every Send and streams no responses. Every
// other RPC answers CodeUnimplemented.
type noopNotificationHandler struct {
	notificationv1connect.UnimplementedNotificationServiceHandler
}

func (noopNotificationHandler) Send(
	_ context.Context,
	_ *notificationv1.SendRequest,
	_ notificationv1connect.NotificationServiceSendServerStream,
) error {
	return nil
}

// GetNotificationCli returns a notification client wired in-process to a
// handler that accepts every Send.
func (bs *ProfileBaseTestSuite) GetNotificationCli(_ *testing.T) notificationv1connect.NotificationServiceClient {
	server := connect.NewServer()
	notificationv1connect.RegisterNotificationServiceHandler(server, noopNotificationHandler{})
	return notificationv1connect.NewNotificationServiceClient(connect.NewClient(connectinprocess.New(server)))
}

func (bs *ProfileBaseTestSuite) CreateTestProfiles(
	ctx context.Context,
	profileBiz business.ProfileBusiness,
	contacts []string,
) ([]*profilev1.ProfileObject, error) {
	var profileSlice []*profilev1.ProfileObject

	for _, contact := range contacts {
		prof := &profilev1.CreateRequest{
			Contact: contact,
		}
		profile, err := profileBiz.CreateProfile(ctx, prof)
		if err != nil {
			return nil, err
		}

		profileSlice = append(profileSlice, profile)
	}

	return profileSlice, nil
}

// WithAuthClaims adds authentication claims to a context for testing.
func (bs *ProfileBaseTestSuite) WithAuthClaims(
	ctx context.Context,
	tenantID, partitionID, profileID string,
) context.Context {
	claims := &security.AuthenticationClaims{
		TenantID:    tenantID,
		PartitionID: partitionID,
		AccessID:    util.IDString(),
		ContactID:   profileID,
		SessionID:   util.IDString(),
		DeviceID:    "test-device",
	}
	claims.Subject = profileID
	return claims.ClaimsToContext(ctx)
}

// SeedTenantAccess writes a tenancy_access member tuple so the profile can pass
// the TenancyAccessChecker (data access layer).
func (bs *ProfileBaseTestSuite) SeedTenantAccess(
	ctx context.Context,
	svc *frame.Service,
	tenantID, partitionID, profileID string,
) {
	auth := svc.SecurityManager().GetAuthorizer(ctx)
	tenancyPath := fmt.Sprintf("%s/%s", tenantID, partitionID)
	err := auth.WriteTuple(ctx, authz.BuildAccessTuple(tenancyPath, profileID))
	bs.Require().NoError(err, "failed to seed tenant access")
}

// SeedTenantRole writes functional permission tuples in the service_profile
// namespace for the given role.
func (bs *ProfileBaseTestSuite) SeedTenantRole(
	ctx context.Context,
	svc *frame.Service,
	tenantID, partitionID, profileID, role string,
) {
	auth := svc.SecurityManager().GetAuthorizer(ctx)
	tenancyPath := fmt.Sprintf("%s/%s", tenantID, partitionID)

	permissions := authz.RolePermissions()[role]
	tuples := make([]security.RelationTuple, 0, 1+len(permissions))

	tuples = append(tuples, security.RelationTuple{
		Object:   security.ObjectRef{Namespace: authz.NamespaceProfile, ID: tenancyPath},
		Relation: role,
		Subject:  security.SubjectRef{Namespace: authz.NamespaceProfileUser, ID: profileID},
	})
	for _, perm := range permissions {
		tuples = append(tuples, security.RelationTuple{
			Object:   security.ObjectRef{Namespace: authz.NamespaceProfile, ID: tenancyPath},
			Relation: authz.GrantedRelation(perm),
			Subject:  security.SubjectRef{Namespace: authz.NamespaceProfileUser, ID: profileID},
		})
	}

	err := auth.WriteTuples(ctx, tuples)
	bs.Require().NoError(err, "failed to seed tenant role")
}

func (bs *ProfileBaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (bs *ProfileBaseTestSuite) WithTestDependancies(
	t *testing.T,
	testFn func(t *testing.T, dep *definition.DependencyOption),
) {
	options := []*definition.DependencyOption{
		definition.NewDependancyOption(
			"default",
			util.RandomAlphaNumericString(DefaultRandomStringLength),
			bs.Resources(),
		),
	}

	frametests.WithTestDependencies(t, options, testFn)
}

// StagedFacts reads the outbox rows with the given fact name, whatever tenant
// staged them.
func StagedFacts(ctx context.Context, svc *frame.Service, name string) ([]*outbox.Event, error) {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	var staged []*outbox.Event
	err := dbPool.DB(tenancy.WithSkipEnforcement(ctx), true).
		Where("name = ?", name).
		Order("created_at asc").
		Find(&staged).Error
	return staged, err
}
