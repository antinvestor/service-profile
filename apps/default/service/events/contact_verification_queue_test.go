package events_test

import (
	"testing"
	"time"

	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type ContactVerificationQueueTestSuite struct {
	tests.BaseTestSuite
}

func TestContactVerificationQueueSuite(t *testing.T) {
	suite.Run(t, new(ContactVerificationQueueTestSuite))
}

func (cvqts *ContactVerificationQueueTestSuite) TestContactVerificationQueue_Name() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)
		require.Equal(t, events.VerificationEventHandlerName, queue.Name())
	})
}

func (cvqts *ContactVerificationQueueTestSuite) TestContactVerificationQueue_PayloadType() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)
		payloadType := queue.PayloadType()

		// Should return a pointer to models.Verification
		_, ok := payloadType.(*models.Verification)
		require.True(t, ok, "PayloadType should return *models.Verification")
	})
}

func (cvqts *ContactVerificationQueueTestSuite) TestContactVerificationQueue_Validate() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)

		// Test valid payload
		validPayload := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		validPayload.GenID(ctx)

		err := queue.Validate(ctx, validPayload)
		require.NoError(t, err)

		// Test invalid payload type
		invalidPayload := "invalid"
		err = queue.Validate(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *models.Verification")
	})
}

func (cvqts *ContactVerificationQueueTestSuite) TestContactVerificationQueue_Execute_Success() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		// Create test contact first
		contactRepo := repository.NewContactRepository(svc)
		contact := &models.Contact{
			Detail:      "test@example.com",
			ContactType: "EMAIL",
		}
		contact.GenID(ctx)
		savedContact, err := contactRepo.Save(ctx, contact)
		require.NoError(t, err)

		// Create verification
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: savedContact.GetID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		verification.GenID(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)

		// Execute the queue handler
		err = queue.Execute(ctx, verification)
		require.NoError(t, err)
	})
}

func (cvqts *ContactVerificationQueueTestSuite) TestContactVerificationQueue_Execute_InvalidPayload() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)

		// Test with invalid payload type
		invalidPayload := "invalid"
		err := queue.Execute(ctx, invalidPayload)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid payload type, expected *models.Verification")
	})
}

func (cvqts *ContactVerificationQueueTestSuite) TestContactVerificationQueue_Execute_NonExistentContact() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		// Create verification with non-existent contact
		verification := &models.Verification{
			ProfileID: util.IDString(),
			ContactID: util.IDString(), // Non-existent contact ID
			Code:      "123456",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		verification.GenID(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)

		// Execute should handle non-existent contact gracefully
		err := queue.Execute(ctx, verification)
		require.Error(t, err) // Should return error when contact doesn't exist
	})
}

func (cvqts *ContactVerificationQueueTestSuite) TestNewContactVerificationQueue() {
	t := cvqts.T()

	cvqts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		svc, ctx := cvqts.CreateService(t, dep)
		notificationCli := cvqts.GetNotificationCli(ctx)

		queue := events.NewContactVerificationQueue(svc, notificationCli)
		require.NotNil(t, queue)
		require.Equal(t, svc, queue.Service)
		require.Equal(t, notificationCli, queue.NotificationCli)
	})
}
