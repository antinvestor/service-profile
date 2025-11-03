package events_test

import (
	"context"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type ContactVerificationAttemptQueueTestSuite struct {
	tests.ProfileBaseTestSuite
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) getVerificationAttemptEvtQ(
	ctx context.Context,
	svc *frame.Service,
) (*events.ContactVerificationAttemptedQueue, repository.VerificationRepository) {
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	contactRepo := repository.NewContactRepository(ctx, dbPool, workMan)
	verificationRepo := repository.NewVerificationRepository(ctx, dbPool, workMan)

	return events.NewContactVerificationAttemptedQueue(contactRepo, verificationRepo), verificationRepo
}

func TestContactVerificationAttemptQueueSuite(t *testing.T) {
	suite.Run(t, new(ContactVerificationAttemptQueueTestSuite))
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_Name() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue, _ := cvaqts.getVerificationAttemptEvtQ(ctx, svc)
		require.Equal(t, events.VerificationAttemptEventHandlerName, queue.Name())
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_PayloadType() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue, _ := cvaqts.getVerificationAttemptEvtQ(ctx, svc)
		payloadType := queue.PayloadType()

		// Should return a pointer to models.VerificationAttempt
		_, ok := payloadType.(*models.VerificationAttempt)
		require.True(t, ok, "PayloadType should return *models.VerificationAttempt")
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_Validate() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue, _ := cvaqts.getVerificationAttemptEvtQ(ctx, svc)

		// Test valid payload
		validPayload := &models.VerificationAttempt{
			VerificationID: util.IDString(),
			Data:           "123456",
			State:          "Success",
			DeviceID:       util.IDString(),
			IPAddress:      "192.168.1.1",
			RequestID:      util.IDString(),
		}
		validPayload.GenID(ctx)

		err := queue.Validate(ctx, validPayload)
		require.NoError(t, err)

		// Test invalid payload type
		invalidPayload := "invalid"
		err = queue.Validate(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *models.VerificationAttempt")
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_Execute_Success() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue, verificationRepo := cvaqts.getVerificationAttemptEvtQ(ctx, svc)

		// Create test verification first
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
		}
		verification.GenID(ctx)
		err := verificationRepo.Create(ctx, verification)
		require.NoError(t, err)

		// Create verification attempt
		attempt := &models.VerificationAttempt{
			VerificationID: verification.GetID(),
			Data:           "123456",
			State:          "Success",
			DeviceID:       util.IDString(),
			IPAddress:      "192.168.1.1",
			RequestID:      util.IDString(),
		}
		attempt.GenID(ctx)

		// Execute the queue handler
		err = queue.Execute(ctx, attempt)
		require.NoError(t, err)

		// Verify the final state
		attemptList, err := verificationRepo.GetAttempts(ctx, verification.GetID())
		require.NoError(t, err)
		require.Len(t, attemptList, 1)
		require.Equal(t, attempt.GetID(), attemptList[0].GetID())
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_Execute_InvalidPayload() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue, _ := cvaqts.getVerificationAttemptEvtQ(ctx, svc)

		// Test with invalid payload type
		invalidPayload := "invalid"
		err := queue.Execute(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *models.VerificationAttempt")
	})
}
