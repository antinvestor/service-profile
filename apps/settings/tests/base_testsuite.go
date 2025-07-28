package tests

import (
	"context"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests"
	"github.com/pitabwire/frame/tests/deps/testpostgres"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/settings/config"
	"github.com/antinvestor/service-profile/apps/settings/service/events"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

const (
	PostgresqlDBImage = "paradedb/paradedb:latest"

	DefaultRandomStringLength = 8
)

type SettingsBaseTestSuite struct {
	tests.FrameBaseTestSuite
}

func initResources(_ context.Context) []testdef.TestResource {
	pg := testpostgres.NewPGDepWithCred(PostgresqlDBImage, "ant", "s3cr3t", "service_settings")
	resources := []testdef.TestResource{pg}
	return resources
}

func (bs *SettingsBaseTestSuite) SetupSuite() {
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()
}

func (bs *SettingsBaseTestSuite) CreateService(
	t *testing.T,
	depOpts *testdef.DependancyOption,
) (*frame.Service, context.Context) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	cfg, err := frame.ConfigFromEnv[config.SettingsConfig]()
	require.NoError(t, err)

	cfg.LogLevel = "debug"
	cfg.RunServiceSecurely = false
	cfg.ServerPort = ""

	for _, res := range depOpts.Database() {
		testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
		require.NoError(t, err0)

		t.Cleanup(func() {
			cleanup(t.Context())
		})

		cfg.DatabasePrimaryURL = []string{testDS.String()}
		cfg.DatabaseReplicaURL = []string{testDS.String()}
	}

	ctx, svc := frame.NewServiceWithContext(t.Context(), "settings tests",
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
		frame.WithNoopDriver())

	eventList := frame.WithRegisterEvents(
		&events.SettingsAuditor{Service: svc})

	svc.Init(ctx, eventList)

	err = repository.Migrate(ctx, svc, "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return svc, ctx
}

func (bs *SettingsBaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (bs *SettingsBaseTestSuite) WithTestDependancies(
	t *testing.T,
	testFn func(t *testing.T, dep *testdef.DependancyOption),
) {
	options := []*testdef.DependancyOption{
		testdef.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), bs.Resources()),
	}

	tests.WithTestDependancies(t, options, testFn)
}
