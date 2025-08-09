package events

import (
	"context"
	"errors"
	"time"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const VerificationAttemptEventHandlerName = "contact.verification.attempt.queue"

type ContactVerificationAttemptedQueue struct {
	Service          *frame.Service
	ContactRepo      repository.ContactRepository
	VerificationRepo repository.VerificationRepository
}

func NewContactVerificationAttemptedQueue(service *frame.Service) *ContactVerificationAttemptedQueue {
	return &ContactVerificationAttemptedQueue{
		Service:          service,
		ContactRepo:      repository.NewContactRepository(service),
		VerificationRepo: repository.NewVerificationRepository(service),
	}
}
func (vaq *ContactVerificationAttemptedQueue) Name() string {
	return VerificationAttemptEventHandlerName
}

func (vaq *ContactVerificationAttemptedQueue) PayloadType() any {
	return &models.VerificationAttempt{}
}

func (vaq *ContactVerificationAttemptedQueue) Validate(_ context.Context, payload any) error {
	notification, ok := payload.(*models.VerificationAttempt)
	if !ok {
		return errors.New("invalid payload type, expected *models.VerificationAttempt")
	}

	if notification.GetID() == "" {
		return errors.New("invalid payload type, expected Id on *models.VerificationAttempt to have been set ")
	}
	return nil
}

func (vaq *ContactVerificationAttemptedQueue) Execute(ctx context.Context, payload any) error {
	attempt, ok := payload.(*models.VerificationAttempt)
	if !ok {
		return errors.New("invalid payload type, expected *models.VerificationAttempt")
	}

	logger := vaq.Service.Log(ctx).WithField("attempt", attempt.GetID()).WithField("type", vaq.Name())

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	err := vaq.VerificationRepo.SaveAttempt(ctx, attempt)
	if err != nil {
		logger.WithError(err).Error("Failed to save verification attempt")
		return err
	}

	verification, err := vaq.VerificationRepo.GetByID(ctx, attempt.VerificationID)
	if err != nil {
		logger.WithError(err).Error("Failed to get verification attempted")
		return nil
	}

	if verification.Code != attempt.Data || verification.ExpiresAt.Before(time.Now()) {
		return nil
	}

	contact, err := vaq.ContactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		logger.WithError(err).Error("Failed to get contact")
		return nil
	}

	verification.VerifiedAt = time.Now()
	err = vaq.VerificationRepo.Save(ctx, verification)
	if err != nil {
		logger.WithError(err).Error("Failed to save verification")
		return err
	}

	contact.VerificationID = verification.ID

	_, err = vaq.ContactRepo.Save(ctx, contact)
	if err != nil {
		logger.WithError(err).Error("Failed to save contact")
		return nil
	}

	return nil
}
