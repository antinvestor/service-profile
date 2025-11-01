package tests

import (
	"context"
	"testing"

	"github.com/antinvestor/apis/go/common"
	"github.com/antinvestor/apis/go/common/mocks"
	commonv1 "github.com/antinvestor/apis/go/common/v1"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	notificationv1_mocks "github.com/antinvestor/apis/go/notification/v1_mocks"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	aconfig "github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

const PostgresqlDBImage = "postgres:latest"

const (
	DefaultRandomStringLength = 8
)

type BaseTestSuite struct {
	frametests.FrameBaseTestSuite

	ContactRepo      repository.ContactRepository
	VerificationRepo repository.VerificationRepository
	AddressRepo      repository.AddressRepository
	ProfileRepo      repository.ProfileRepository
	RosterRepo       repository.RosterRepository
	RelationshipRepo repository.RelationshipRepository
}

func initResources(_ context.Context) []definition.TestResource {
	pg := testpostgres.NewWithOpts("service_profile", definition.WithUserName("ant"))
	resources := []definition.TestResource{pg}
	return resources
}

func (bs *BaseTestSuite) SetupSuite() {
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()
}

func (bs *BaseTestSuite) CreateService(
	t *testing.T,
	depOpts *definition.DependencyOption,
) (*frame.Service, context.Context) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")

	ctx := t.Context()
	cfg, err := config.FromEnv[aconfig.ProfileConfig]()
	require.NoError(t, err)

	cfg.LogLevel = "debug"
	cfg.RunServiceSecurely = false
	cfg.ServerPort = ""

	res := depOpts.ByIsDatabase(ctx)
	testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
	require.NoError(t, err0)

	t.Cleanup(func() {
		cleanup(t.Context())
	})

	cfg.DatabasePrimaryURL = []string{testDS.String()}
	cfg.DatabaseReplicaURL = []string{testDS.String()}

	ctx, svc := frame.NewServiceWithContext(t.Context(), "profile tests",
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
		frametests.WithNoopDriver())

	relationshipConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipConnectName,
		cfg.QueueRelationshipConnectURI,
	)
	relationshipDisConnectQueuePublisher := frame.WithRegisterPublisher(
		cfg.QueueRelationshipDisConnectName,
		cfg.QueueRelationshipDisConnectURI,
	)

	evtsMan := svc.EventsManager(ctx)
	qMan := svc.QueueManager(ctx)
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)
	relationshipRepo := repository.NewRelationshipRepository(ctx, dbPool, workMan)

	svc.Init(ctx,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.WithRegisterEvents(
			events.NewClientConnectedSetupQueue(ctx, &cfg, qMan, evtsMan, relationshipRepo),
			events.NewContactVerificationQueue(&cfg, contactRepo, verificationRepo, bs.GetNotificationCli(ctx)),
			events.NewContactVerificationAttemptedQueue(contactRepo, verificationRepo),
		),
	)

	err = repository.Migrate(ctx, svc.DatastoreManager(), dbPool, "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return svc, ctx
}

func (bs *BaseTestSuite) GetNotificationCli(_ context.Context) *notificationv1.NotificationClient {
	mockNotificationService := notificationv1_mocks.NewMockNotificationServiceClient(bs.Ctrl)
	mockNotificationService.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, _ *notificationv1.SendRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[notificationv1.SendResponse], error) {
			// Return a successful response with a generated message ID
			const randomIDLength = 6
			resp := &notificationv1.SendResponse{
				Data: []*commonv1.StatusResponse{
					{
						Id:         util.IDString(),
						State:      commonv1.STATE_ACTIVE,
						Status:     commonv1.STATUS_SUCCESSFUL,
						ExternalId: util.RandomString(randomIDLength),
					},
				},
			}

			// Create a custom mock implementation
			mockStream := mocks.NewMockServerStreamingClient[notificationv1.SendResponse](ctx)
			err := mockStream.SendMsg(resp)
			if err != nil {
				return nil, err
			}

			return mockStream, nil
		}).
		AnyTimes()
	notificationCli := notificationv1.Init(&common.GrpcClientBase{}, mockNotificationService)

	return notificationCli
}

func (bs *BaseTestSuite) CreateTestProfiles(
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

func (bs *BaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (bs *BaseTestSuite) WithTestDependancies(
	t *testing.T,
	testFn func(t *testing.T, dep *definition.DependencyOption),
) {
	options := []*definition.DependencyOption{
		definition.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), bs.Resources()),
	}

	frametests.WithTestDependencies(t, options, testFn)
}
