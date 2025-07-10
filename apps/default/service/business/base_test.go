package business_test

import (
	"context"
	"fmt"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/queue"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"testing"

	"github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/mock/gomock"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests"
	"github.com/pitabwire/frame/tests/deps/testpostgres"
	"github.com/pitabwire/frame/tests/testdef"
)

const PostgresqlDbImage = "paradedb/paradedb:latest"

// StdoutLogConsumer is a LogConsumer that prints the log to stdout.
type StdoutLogConsumer struct{}

// Accept prints the log to stdout.
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}

type BaseTestSuite struct {
	tests.FrameBaseTestSuite
}

func initResources(_ context.Context) []testdef.TestResource {
	pg := testpostgres.NewPGDepWithCred(PostgresqlDbImage, "ant", "s3cr3t", "service_profile")
	resources := []testdef.TestResource{pg}
	return resources
}

func (bs *BaseTestSuite) SetupSuite() {
	bs.MigrationImageContext = "../../"
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()
}

func (bs *BaseTestSuite) CreateService(
	t *testing.T,
	depOpts *testdef.DependancyOption,
) (*frame.Service, context.Context) {
	profileConfig, err := frame.ConfigFromEnv[config.ProfileConfig]()
	require.NoError(t, err)

	profileConfig.LogLevel = "debug"
	profileConfig.RunServiceSecurely = false
	profileConfig.ServerPort = ""

	for _, res := range depOpts.Database() {
		testDS, cleanup, err0 := res.GetPrefixedDS(t.Context(), depOpts.Prefix())
		require.NoError(t, err0)

		t.Cleanup(func() {
			cleanup(t.Context())
		})

		profileConfig.DatabasePrimaryURL = []string{testDS.String()}
		profileConfig.DatabaseReplicaURL = []string{testDS.String()}

		err0 = bs.Migrate(t.Context(), testDS)
		require.NoError(t, err0)
	}

	ctx, service := frame.NewServiceWithContext(t.Context(), "profile tests",
		frame.WithConfig(&profileConfig),
		frame.WithDatastore(),
		frame.WithNoopDriver())

	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         service,
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: bs.GetNotificationCli(ctx),
	}

	verificationQueue := frame.WithRegisterSubscriber(
		profileConfig.QueueVerificationName,
		profileConfig.QueueVerification,
		&verificationQueueHandler,
	)
	verificationQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueVerificationName,
		profileConfig.QueueVerification,
	)
	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueRelationshipConnectName,
		profileConfig.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		profileConfig.QueueRelationshipDisConnectName,
		profileConfig.QueueRelationshipDisConnectURI,
	)

	service.Init(ctx, verificationQueue, verificationQueuePublisher,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(&events.ClientConnectedSetupQueue{Service: service}),
	)

	err = service.Run(ctx, "")
	require.NoError(t, err)

	return service, ctx
}

func (bs *BaseTestSuite) GetNotificationCli(_ context.Context) *notificationv1.NotificationClient {
	t := bs.T()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNotificationService := notificationv1.NewMockNotificationServiceClient(ctrl)
	mockNotificationService.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()
	notificationCli := notificationv1.Init(&common.GrpcClientBase{}, mockNotificationService)

	return notificationCli
}

func (bs *BaseTestSuite) CreateTestProfiles(
	ctx context.Context,
	svc *frame.Service,
	contacts []string,
) ([]*profilev1.ProfileObject, error) {
	profBuss := business.NewProfileBusiness(ctx, svc)

	var profileSlice []*profilev1.ProfileObject

	for _, contact := range contacts {
		prof := &profilev1.CreateRequest{
			Contact: contact,
		}
		profile, err := profBuss.CreateProfile(ctx, prof)
		if err != nil {
			return nil, err
		}

		profileSlice = append(profileSlice, profile)
	}

	return profileSlice, nil
}

func (bs *BaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (bs *BaseTestSuite) WithTestDependancies(t *testing.T, testFn func(t *testing.T, dep *testdef.DependancyOption)) {
	options := []*testdef.DependancyOption{
		testdef.NewDependancyOption("default", util.RandomString(8), bs.Resources()),
	}

	tests.WithTestDependancies(t, options, testFn)
}
