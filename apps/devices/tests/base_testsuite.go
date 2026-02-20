package tests

import (
	"context"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"

	aconfig "github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/caching"
	devQueue "github.com/antinvestor/service-profile/apps/devices/service/queue"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

const (
	DefaultRandomStringLength = 8
)

type DeviceBaseTestSuite struct {
	frametests.FrameBaseTestSuite
}

type DepsBuilder struct {
	DeviceRepo    repository.DeviceRepository
	DeviceLogRepo repository.DeviceLogRepository
	SessionRepo   repository.DeviceSessionRepository
	KeyRepo       repository.DeviceKeyRepository
	PresenceRepo  repository.DevicePresenceRepository

	DeviceBusiness business.DeviceBusiness
	KeyBusiness    business.KeysBusiness

	AnalysisQueueHandler *devQueue.DeviceAnalysisQueueHandler
}

func BuildRepos(ctx context.Context, svc *frame.Service) *DepsBuilder {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	workMan := svc.WorkManager()
	qMan := svc.QueueManager()

	cfg, _ := svc.Config().(*aconfig.DevicesConfig)

	// Initialize cache service (may be nil in test environments without cache).
	var cacheSvc *caching.DeviceCacheService
	if cacheMan := svc.CacheManager(); cacheMan != nil {
		cacheSvc = caching.NewDeviceCacheService(cacheMan)
	}

	deviceRepo := repository.NewDeviceRepository(ctx, dbPool, workMan)
	deviceLogRepo := repository.NewDeviceLogRepository(ctx, dbPool, workMan)
	sessionRepo := repository.NewDeviceSessionRepository(ctx, dbPool, workMan)
	keyRepo := repository.NewDeviceKeyRepository(ctx, dbPool, workMan)
	presenceRepo := repository.NewDevicePresenceRepository(ctx, dbPool, workMan)

	deviceBusiness := business.NewDeviceBusiness(
		ctx,
		cfg,
		qMan,
		workMan,
		deviceRepo,
		deviceLogRepo,
		sessionRepo,
		cacheSvc,
	)
	keyBusiness := business.NewKeysBusiness(ctx, cfg, qMan, workMan, deviceRepo, keyRepo, cacheSvc)

	return &DepsBuilder{
		DeviceRepo:    deviceRepo,
		SessionRepo:   sessionRepo,
		DeviceLogRepo: deviceLogRepo,

		KeyRepo:      keyRepo,
		PresenceRepo: presenceRepo,

		DeviceBusiness: deviceBusiness,
		KeyBusiness:    keyBusiness,

		AnalysisQueueHandler: devQueue.NewDeviceAnalysisQueueHandler(
			svc.HTTPClientManager(),
			deviceRepo,
			deviceLogRepo,
			sessionRepo,
			cacheSvc,
		),
	}
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
) (context.Context, *frame.Service, *DepsBuilder) {
	ctx := t.Context()
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	cfg, err := config.FromEnv[aconfig.DevicesConfig]()
	require.NoError(t, err)

	cfg.LogLevel = "debug"
	cfg.RunServiceSecurely = false
	cfg.DatabaseMigrate = true
	cfg.DatabaseTraceQueries = true
	cfg.ServerPort = ""

	res := depOpts.ByIsDatabase(ctx)
	testDS, cleanup, err0 := res.GetRandomisedDS(ctx, depOpts.Prefix())
	require.NoError(t, err0)

	t.Cleanup(func() {
		cleanup(ctx)
	})

	cfg.DatabasePrimaryURL = []string{testDS.String()}
	cfg.DatabaseReplicaURL = []string{testDS.String()}

	ctx, svc := frame.NewServiceWithContext(ctx, frame.WithName("device tests"),
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
		frame.WithCacheManager(),
		frame.WithInMemoryCache(aconfig.CacheNameDevices),
		frame.WithInMemoryCache(aconfig.CacheNamePresence),
		frame.WithInMemoryCache(aconfig.CacheNameGeoIP),
		frame.WithInMemoryCache(aconfig.CacheNameRate),
		frametests.WithNoopDriver())

	depsBuilder := BuildRepos(ctx, svc)

	analysisQueueTopic := frame.WithRegisterPublisher(
		cfg.QueueDeviceAnalysisName,
		cfg.QueueDeviceAnalysis,
	)

	analysisQueue := frame.WithRegisterSubscriber(
		cfg.QueueDeviceAnalysisName,
		cfg.QueueDeviceAnalysis,
		depsBuilder.AnalysisQueueHandler,
	)

	svc.Init(ctx, analysisQueueTopic, analysisQueue)

	err = repository.Migrate(ctx, svc.DatastoreManager(), "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return security.SkipTenancyChecksOnClaims(ctx), svc, depsBuilder
}

func (bs *DeviceBaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependencies Creates subtests with each known DependancyOption.
func (bs *DeviceBaseTestSuite) WithTestDependencies(
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
