package events

import (
	"context"
	"errors"
	"time"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const VerificationAttemptEventHandlerName = "contact.verification.attempt.queue"

type ContactVerificationAttemptedQueue struct {
	ContactRepo      repository.ContactRepository
	VerificationRepo repository.VerificationRepository
}

func NewContactVerificationAttemptedQueue(
	contactRepo repository.ContactRepository,
	verificationRepo repository.VerificationRepository,
) *ContactVerificationAttemptedQueue {
	return &ContactVerificationAttemptedQueue{
		ContactRepo:      contactRepo,
		VerificationRepo: verificationRepo,
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

	logger := util.Log(ctx).WithFields(map[string]any{
		"attempt_id": attempt.GetID(),
		"type":       vaq.Name(),
	})

	err := vaq.VerificationRepo.SaveAttempt(ctx, attempt)
	if err != nil {
		if data.ErrorIsDuplicateKey(err) {
			logger.Debug("verification attempt already exists, skipping duplicate")
			return nil
		}
		logger.WithError(err).Error("failed to save verification attempt")
		return err
	}

	verification, err := vaq.VerificationRepo.GetByID(ctx, attempt.VerificationID)
	if err != nil {
		logger.WithError(err).Error("failed to get verification")
		return nil
	}

	if verification.Code != attempt.Data || verification.ExpiresAt.Before(time.Now()) {
		return nil
	}

	contact, err := vaq.ContactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		logger.WithField("contact_id", verification.ContactID).WithError(err).Error("failed to get contact")
		return nil
	}

	verification.VerifiedAt = time.Now()
	_, err = vaq.VerificationRepo.Update(ctx, verification, "verified_at")
	if err != nil {
		logger.WithError(err).Error("failed to save verification")
		return err
	}

	contact.VerificationID = verification.ID

	_, err = vaq.ContactRepo.Update(ctx, contact, "verification_id")
	if err != nil {
		logger.WithField("contact_id", contact.GetID()).WithError(err).Error("failed to save contact")
		return nil
	}

	return nil
}
