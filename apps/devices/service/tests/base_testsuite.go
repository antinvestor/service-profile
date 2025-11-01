package tests

import (
	"context"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
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
	frametests.FrameBaseTestSuite
}

func initResources(_ context.Context) []definition.TestResource {
	pg := testpostgres.NewWithOpts("service_devices", definition.WithUserName("ant"))
	resources := []definition.TestResource{pg}
	return resources
}

func (bs *DeviceBaseTestSuite) SetupSuite() {
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()
}

func (bs *DeviceBaseTestSuite) CreateService(
	t *testing.T,
	depOpts *definition.DependencyOption,
) (*frame.Service, context.Context) {
	ctx := t.Context()
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	deviceConfig, err := frame.ConfigFromEnv[config.DevicesConfig]()
	require.NoError(t, err)

	deviceConfig.LogLevel = "debug"
	deviceConfig.RunServiceSecurely = false
	deviceConfig.ServerPort = ""

	res := depOpts.ByIsDatabase(ctx)
	testDS, cleanup, err0 := res.GetRandomisedDS(ctx, depOpts.Prefix())
	require.NoError(t, err0)

	t.Cleanup(func() {
		cleanup(ctx)
	})

	deviceConfig.DatabasePrimaryURL = []string{testDS.String()}
	deviceConfig.DatabaseReplicaURL = []string{testDS.String()}

	ctx, svc := frame.NewServiceWithContext(ctx, "device tests",
		frame.WithConfig(&deviceConfig),
		frame.WithDatastore(),
		frametests.WithNoopDriver())

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
	testFn func(t *testing.T, dep *definition.DependencyOption),
) {
	options := []*definition.DependencyOption{
		definition.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), bs.Resources()),
	}

	frametests.WithTestDependancies(t, options, testFn)
}
