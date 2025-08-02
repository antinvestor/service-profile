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

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const (
	DefaultRandomStringLength = 8
)

type DeviceBaseTestSuite struct {
	tests.FrameBaseTestSuite
}

func initResources(_ context.Context) []testdef.TestResource {
	pg := testpostgres.NewPGDepWithCred(testpostgres.PostgresqlDBImage, "ant", "s3cr3t", "service_profile")
	resources := []testdef.TestResource{pg}
	return resources
}

func (bs *DeviceBaseTestSuite) SetupSuite() {
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()
}

func (bs *DeviceBaseTestSuite) CreateService(
	t *testing.T,
	depOpts *testdef.DependancyOption,
) (*frame.Service, context.Context) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	deviceConfig, err := frame.ConfigFromEnv[config.DevicesConfig]()
	require.NoError(t, err)

	deviceConfig.LogLevel = "debug"
	deviceConfig.RunServiceSecurely = false
	deviceConfig.ServerPort = ""

	for _, res := range depOpts.Database() {
		testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
		require.NoError(t, err0)

		t.Cleanup(func() {
			cleanup(t.Context())
		})

		deviceConfig.DatabasePrimaryURL = []string{testDS.String()}
		deviceConfig.DatabaseReplicaURL = []string{testDS.String()}
	}

	ctx, svc := frame.NewServiceWithContext(t.Context(), "device tests",
		frame.WithConfig(&deviceConfig),
		frame.WithDatastore(),
		frame.WithNoopDriver())

	verificationQueueHandler := queue.DeviceAnalysisQueueHandler{
		Service:          svc,
		DeviceRepository: repository.NewDeviceRepository(svc),
	}

	analysisQueueTopic := frame.WithRegisterPublisher(
		deviceConfig.QueueDeviceAnalysisName,
		deviceConfig.QueueDeviceAnalysis,
	)

	analysisQueue := frame.WithRegisterSubscriber(
		deviceConfig.QueueDeviceAnalysisName,
		deviceConfig.QueueDeviceAnalysis,
		&verificationQueueHandler,
	)

	svc.Init(ctx, analysisQueueTopic, analysisQueue)

	err = repository.Migrate(ctx, svc, "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return svc, ctx
}

func (bs *DeviceBaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (bs *DeviceBaseTestSuite) WithTestDependancies(
	t *testing.T,
	testFn func(t *testing.T, dep *testdef.DependancyOption),
) {
	options := []*testdef.DependancyOption{
		testdef.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), bs.Resources()),
	}

	tests.WithTestDependancies(t, options, testFn)
}
