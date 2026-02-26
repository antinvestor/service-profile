package tests

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	aconfig "github.com/antinvestor/service-profile/apps/geolocation/config"
	"github.com/antinvestor/service-profile/apps/geolocation/service/authz"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
	"github.com/antinvestor/service-profile/apps/geolocation/tests/testketo"
)

const (
	DefaultRandomStringLength = 8
)

type GeolocationBaseTestSuite struct {
	frametests.FrameBaseTestSuite

	AuthzMiddleware authz.Middleware
	ketoReadURI     string
	ketoWriteURI    string
}

func initResources(_ context.Context) []definition.TestResource {
	pg := testpostgres.NewWithOpts("service_geolocation", definition.WithUserName("ant"))

	keto := testketo.NewWithOpts(
		definition.WithDependancies(pg),
		definition.WithEnableLogging(true),
	)

	return []definition.TestResource{pg, keto}
}

func (bs *GeolocationBaseTestSuite) SetupSuite() {
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

func (bs *GeolocationBaseTestSuite) CreateService(
	t *testing.T,
	depOpts *definition.DependencyOption,
) (context.Context, *frame.Service) {
	ctx := t.Context()
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	cfg, err := config.FromEnv[aconfig.GeolocationConfig]()
	require.NoError(t, err)

	cfg.LogLevel = "debug"
	cfg.RunServiceSecurely = false
	cfg.DatabaseMigrate = true
	cfg.ServerPort = ""

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

	ctx, svc := frame.NewServiceWithContext(t.Context(), frame.WithName("geolocation tests"),
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
		frametests.WithNoopDriver())

	// Wire real Keto authoriser via SecurityManager
	sm := svc.SecurityManager()
	bs.AuthzMiddleware = authz.NewMiddleware(sm.GetAuthorizer(ctx))

	svc.Init(ctx)

	err = repository.Migrate(ctx, svc.DatastoreManager(), "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return ctx, svc
}

// WithAuthClaims adds authentication claims to a context for testing.
func (bs *GeolocationBaseTestSuite) WithAuthClaims(ctx context.Context, tenantID, profileID string) context.Context {
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
func (bs *GeolocationBaseTestSuite) SeedTenantRole(
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

func (bs *GeolocationBaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependencies Creates subtests with each known DependancyOption.
func (bs *GeolocationBaseTestSuite) WithTestDependencies(
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
