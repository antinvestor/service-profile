package events

import (
	"context"
	"errors"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/pitabwire/util"
)

const ContactKeyRotationEventHandlerName = "contact.key.rotation.queue"

type ContactKeyRotationQueue struct {
	cfg         *config.ProfileConfig
	contactRepo repository.ContactRepository
}

func NewContactKeyRotationQueue(
	cfg *config.ProfileConfig, contactRepo repository.ContactRepository,
) *ContactVerificationQueue {
	return &ContactVerificationQueue{
		cfg:         cfg,
		contactRepo: contactRepo,
	}
}

func (vq *ContactKeyRotationQueue) Name() string {
	return ContactKeyRotationEventHandlerName
}

func (vq *ContactKeyRotationQueue) PayloadType() any {
	return &models.Contact{}
}

func (vq *ContactKeyRotationQueue) Validate(_ context.Context, payload any) error {
	_, ok := payload.(*string)
	if !ok {
		return errors.New("invalid payload type, expected *string")
	}

	return nil
}

func (vq *ContactKeyRotationQueue) Execute(ctx context.Context, payload any) error {
	contactIDPtr, ok := payload.(*string)
	if !ok {
		return errors.New("invalid payload type, expected *string")
	}

	contactID := *contactIDPtr

	logger := util.Log(ctx).WithField("payload", contactID).WithField("type", vq.Name())

	contact, err := vq.contactRepo.GetByID(ctx, contactID)
	if err != nil {
		return err
	}

	logger.WithField("contact", contact.Detail).
		WithField("resp", resp.Msg()).
		Info("successfully submitted verification for contact")

	return nil
}
