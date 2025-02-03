package business_test

import (
	"context"
	"fmt"
	"github.com/antinvestor/apis/go/common"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/business"
	"github.com/antinvestor/service-profile/service/events"
	"github.com/antinvestor/service-profile/service/queue"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pitabwire/frame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/mock/gomock"
	"net"
	"time"
)

const PostgresqlDbImage = "paradedb/paradedb"

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct{}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}

type BaseTestSuite struct {
	suite.Suite
	service     *frame.Service
	ctx         context.Context
	pgContainer *tcPostgres.PostgresContainer
	networks    []string
	postgresUri string
}

func (bs *BaseTestSuite) SetupSuite() {
	ctx := context.Background()

	postgresContainer, err := bs.setupPostgres(ctx)
	assert.NoError(bs.T(), err)

	port, _ := nat.NewPort("tcp", "5432")
	port, _ = postgresContainer.MappedPort(ctx, port)
	fmt.Println(" successfully setup postgresql port : ", port.Port())

	bs.pgContainer = postgresContainer

	bs.networks, err = bs.pgContainer.Networks(ctx)
	assert.NoError(bs.T(), err)

	postgresqlIp, err := bs.pgContainer.ContainerIP(ctx)
	assert.NoError(bs.T(), err)

	bs.postgresUri = fmt.Sprintf("postgres://ant:secret@%s/service_profile?sslmode=disable", net.JoinHostPort(postgresqlIp, "5432"))

	databaseUriStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(bs.T(), err)

	err = bs.setupMigrations(ctx)
	assert.NoError(bs.T(), err)

	profileConfig := config.ProfileConfig{}

	err = frame.ConfigProcess("", &profileConfig)
	assert.NoError(bs.T(), err)

	profileConfig.LogLevel = "debug"
	profileConfig.RunServiceSecurely = false
	profileConfig.ServerPort = ""
	profileConfig.DatabasePrimaryURL = []string{databaseUriStr}
	profileConfig.DatabaseReplicaURL = []string{databaseUriStr}

	var service *frame.Service
	ctx, service = frame.NewServiceWithContext(ctx, "profile tests",
		frame.Config(&profileConfig),
		frame.Datastore(ctx),
		frame.NoopDriver())

	verificationQueueHandler := queue.VerificationsQueueHandler{
		Service:         service,
		ContactRepo:     repository.NewContactRepository(service),
		NotificationCli: bs.getNotificationCli(ctx),
	}

	verificationQueue := frame.RegisterSubscriber(profileConfig.QueueVerificationName, profileConfig.QueueVerification, 2, &verificationQueueHandler)
	verificationQueuePublisher := frame.RegisterPublisher(profileConfig.QueueVerificationName, profileConfig.QueueVerification)
	relationshipConnectQueuePublisher := frame.RegisterPublisher(profileConfig.QueueRelationshipConnectName, profileConfig.QueueRelationshipConnectURI)
	relationshipDisConnectQueuePublisher := frame.RegisterPublisher(profileConfig.QueueRelationshipDisConnectName, profileConfig.QueueRelationshipDisConnectURI)

	service.Init(
		verificationQueue, verificationQueuePublisher,
		relationshipConnectQueuePublisher, relationshipDisConnectQueuePublisher,
		frame.RegisterEvents(&events.ClientConnectedSetupQueue{Service: service}),
	)

	err = service.Run(ctx, "")
	bs.ctx = ctx
	bs.service = service

	assert.NoError(bs.T(), err)
}

func (bs *BaseTestSuite) getNotificationCli(_ context.Context) *notificationv1.NotificationClient {

	t := bs.T()
	//ctx := bs.ctx

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNotificationService := notificationv1.NewMockNotificationServiceClient(ctrl)
	mockNotificationService.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()
	notificationCli := notificationv1.Init(&common.GrpcClientBase{}, mockNotificationService)

	return notificationCli
}

func (bs *BaseTestSuite) createTestProfiles(contacts []string) ([]*profilev1.ProfileObject, error) {

	ctx := bs.ctx

	profBuss := business.NewProfileBusiness(ctx, bs.service)

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

func (bs *BaseTestSuite) setupPostgres(ctx context.Context) (*tcPostgres.PostgresContainer, error) {

	postgresContainer, err := tcPostgres.Run(ctx, PostgresqlDbImage,
		tcPostgres.WithDatabase("service_profile"),
		tcPostgres.WithUsername("ant"),
		tcPostgres.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(20*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	return postgresContainer, nil
}

func (bs *BaseTestSuite) setupMigrations(ctx context.Context) error {

	g := StdoutLogConsumer{}

	cRequest := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "../../",
		},
		ConfigModifier: func(config *container.Config) {
			config.Env = []string{
				"LOG_LEVEL=debug",
				"DO_MIGRATION=true",
				fmt.Sprintf("DATABASE_URL=%s", bs.postgresUri),
			}
		},
		Networks:   bs.networks,
		WaitingFor: wait.ForExit().WithExitTimeout(10 * time.Second),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts:      []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(2 * time.Second)},
			Consumers: []testcontainers.LogConsumer{&g},
		},
	}

	migrationC, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: cRequest,
			Started:          true,
		})
	if err != nil {
		return err
	}

	return migrationC.Terminate(ctx)
}

func (bs *BaseTestSuite) TearDownSuite() {

	t := bs.T()

	if bs.service != nil {
		bs.service.Stop(bs.ctx)
	}

	if bs.pgContainer != nil {
		if err := bs.pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	t.Cleanup(func() {})
}
