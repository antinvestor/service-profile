package events_test

import (
	"context"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type ClientSetupQueueTestSuite struct {
	tests.ProfileBaseTestSuite
}

func TestClientSetupQueueSuite(t *testing.T) {
	suite.Run(t, new(ClientSetupQueueTestSuite))
}

func (csqts *ClientSetupQueueTestSuite) getConnectedSetupEvtQ(
	ctx context.Context,
	svc *frame.Service,
) (*events.ClientConnectedSetupQueue, business.ProfileBusiness, repository.RelationshipRepository) {
	evtsMan := svc.EventsManager()
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	cfg := svc.Config().(*config.ProfileConfig)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	contactBusiness := business.NewContactBusiness(ctx, cfg, evtsMan, contactRepo, verificationRepo)

	addressRepo := repository.NewAddressRepository(ctx, dbPool, workMan)
	addressBusiness := business.NewAddressBusiness(ctx, addressRepo)

	profileRepo := repository.NewProfileRepository(ctx, dbPool, workMan)
	profileBusiness := business.NewProfileBusiness(ctx, evtsMan, contactBusiness, addressBusiness, profileRepo)

	relationshipRepo := repository.NewRelationshipRepository(ctx, dbPool, workMan)

	return events.NewClientConnectedSetupQueue(
		ctx,
		cfg,
		svc.QueueManager(),
		evtsMan,
		relationshipRepo,
	), profileBusiness, relationshipRepo
}

func (csqts *ClientSetupQueueTestSuite) TestClientConnectedSetupQueue_Name() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, _, _ := csqts.getConnectedSetupEvtQ(ctx, svc)
		require.Equal(t, events.ClientConnectedSetupQueueName, queue.Name())
	})
}

func (csqts *ClientSetupQueueTestSuite) TestClientConnectedSetupQueue_PayloadType() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, _, _ := csqts.getConnectedSetupEvtQ(ctx, svc)
		payloadType := queue.PayloadType()

		// Should return a pointer to string
		_, ok := payloadType.(*string)
		require.True(t, ok, "PayloadType should return *string")
	})
}

func (csqts *ClientSetupQueueTestSuite) TestClientConnectedSetupQueue_Validate() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, _, _ := csqts.getConnectedSetupEvtQ(ctx, svc)

		// Test valid payload
		validPayload := "test-relationship-id"
		err := queue.Validate(ctx, &validPayload)
		require.NoError(t, err)

		// Test invalid payload type
		invalidPayload := 123
		err = queue.Validate(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *string")
	})
}

func (csqts *ClientSetupQueueTestSuite) TestClientConnectedSetupQueue_Execute_Success() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, profileBusiness, relationshipRepo := csqts.getConnectedSetupEvtQ(ctx, svc)

		// Create a test relationship
		relationship := &models.Relationship{
			ParentObject:   "profile",
			ParentObjectID: util.IDString(),
			ChildObject:    "profile",
			ChildObjectID:  util.IDString(),
		}
		relationship.GenID(ctx)

		relationshipType, err := relationshipRepo.RelationshipType(ctx, profilev1.RelationshipType_MEMBER)
		require.NoError(t, err)

		relationship.RelationshipTypeID = relationshipType.GetID()

		err = relationshipRepo.Create(ctx, relationship)
		require.NoError(t, err)

		relationshipID := relationship.GetID()

		// Execute the queue handler
		err = queue.Execute(ctx, &relationshipID)
		require.NoError(t, err)

		_ = profileBusiness
	})
}

func (csqts *ClientSetupQueueTestSuite) TestClientConnectedSetupQueue_Execute_InvalidPayload() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, _, _ := csqts.getConnectedSetupEvtQ(ctx, svc)

		// Test with invalid payload type
		invalidPayload := 123
		err := queue.Execute(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *string")
	})
}

func (csqts *ClientSetupQueueTestSuite) TestClientConnectedSetupQueue_Execute_NonExistentRelationship() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, _, _ := csqts.getConnectedSetupEvtQ(ctx, svc)
		nonExistentID := util.IDString()

		// Execute with non-existent relationship ID - should not return error (logs and continues)
		err := queue.Execute(ctx, &nonExistentID)
		require.NoError(t, err, "Should handle non-existent relationship gracefully")
	})
}

func (csqts *ClientSetupQueueTestSuite) TestNewClientConnectedSetupQueue() {
	t := csqts.T()

	csqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := csqts.CreateService(t, dep)

		queue, _, _ := csqts.getConnectedSetupEvtQ(ctx, svc)
		require.NotNil(t, queue)
	})
}
