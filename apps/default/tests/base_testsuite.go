package tests

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	notificationv1mocks "github.com/antinvestor/apis/go/notification/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	aconfig "github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests/testketo"
)

const PostgresqlDBImage = "postgres:latest"

const (
	DefaultRandomStringLength = 8
)

type ProfileBaseTestSuite struct {
	frametests.FrameBaseTestSuite

	AuthzMiddleware authz.Middleware
	ketoReadURI     string
	ketoWriteURI    string

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

	ctx, svc := frame.NewServiceWithContext(t.Context(), frame.WithName("profile tests"),
		frame.WithConfig(&cfg),
		frame.WithDatastore(pool.WithTraceConfig(&cfg)),
		frametests.WithNoopDriver())

	// Wire real Keto authoriser via SecurityManager
	sm := svc.SecurityManager()
	bs.AuthzMiddleware = authz.NewMiddleware(sm.GetAuthorizer(ctx))

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

	svc.Init(ctx,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			events.NewClientConnectedSetupQueue(ctx, &cfg, qMan, evtsMan, relationshipRepo),
			events.NewContactVerificationQueue(&cfg, contactRepo, verificationRepo, bs.GetNotificationCli(t)),
			events.NewContactVerificationAttemptedQueue(contactRepo, verificationRepo),
		),
	)

	err = repository.Migrate(ctx, svc.DatastoreManager(), "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return ctx, svc
}

func (bs *ProfileBaseTestSuite) GetNotificationCli(t *testing.T) notificationv1connect.NotificationServiceClient {
	mc := minimock.NewController(t)

	notificationCli := notificationv1mocks.NewNotificationServiceClientMock(mc)

	notificationCli.SendMock.Optional().Return(nil, nil)

	return notificationCli
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
func (bs *ProfileBaseTestSuite) WithAuthClaims(ctx context.Context, tenantID, profileID string) context.Context {
	claims := &security.AuthenticationClaims{
		TenantID:  tenantID,
		AccessID:  util.IDString(),
		ContactID: profileID,
		SessionID: util.IDString(),
		DeviceID:  "test-device",
	}
	claims.Subject = profileID
	return claims.ClaimsToContext(ctx)
}

// SeedTenantRole writes a relation tuple granting the given role to a profile on a tenant.
func (bs *ProfileBaseTestSuite) SeedTenantRole(
	ctx context.Context,
	svc *frame.Service,
	tenantID, profileID, role string,
) {
	auth := svc.SecurityManager().GetAuthorizer(ctx)
	err := auth.WriteTuple(ctx, security.RelationTuple{
		Object:   security.ObjectRef{Namespace: authz.NamespaceTenant, ID: tenantID},
		Relation: role,
		Subject:  security.SubjectRef{Namespace: authz.NamespaceProfile, ID: profileID},
	})
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
