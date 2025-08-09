package events_test

import (
	"testing"

	"github.com/pitabwire/frame/tests/testdef"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type ContactVerificationAttemptQueueTestSuite struct {
	tests.BaseTestSuite
}

func TestContactVerificationAttemptQueueSuite(t *testing.T) {
	suite.Run(t, new(ContactVerificationAttemptQueueTestSuite))
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_Name() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, _ := cvaqts.CreateService(t, dep)

		queue := events.NewContactVerificationAttemptedQueue(svc)
		require.Equal(t, events.VerificationAttemptEventHandlerName, queue.Name())
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_PayloadType() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, _ := cvaqts.CreateService(t, dep)

		queue := events.NewContactVerificationAttemptedQueue(svc)
		payloadType := queue.PayloadType()

		// Should return a pointer to models.VerificationAttempt
		_, ok := payloadType.(*models.VerificationAttempt)
		require.True(t, ok, "PayloadType should return *models.VerificationAttempt")
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestContactVerificationAttemptedQueue_Validate() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue := events.NewContactVerificationAttemptedQueue(svc)

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

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		// Create test verification first
		verificationRepo := repository.NewVerificationRepository(svc)
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
		}
		verification.GenID(ctx)
		err := verificationRepo.Save(ctx, verification)
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

		queue := events.NewContactVerificationAttemptedQueue(svc)

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

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := cvaqts.CreateService(t, dep)

		queue := events.NewContactVerificationAttemptedQueue(svc)

		// Test with invalid payload type
		invalidPayload := "invalid"
		err := queue.Execute(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *models.VerificationAttempt")
	})
}

func (cvaqts *ContactVerificationAttemptQueueTestSuite) TestNewContactVerificationAttemptedQueue() {
	t := cvaqts.T()

	cvaqts.WithTestDependancies(t, func(t *testing.T, dep *testdef.DependancyOption) {
		svc, _ := cvaqts.CreateService(t, dep)

		queue := events.NewContactVerificationAttemptedQueue(svc)
		require.NotNil(t, queue)
		require.Equal(t, svc, queue.Service)
		require.NotNil(t, queue.ContactRepo)
		require.NotNil(t, queue.VerificationRepo)
	})
}
