package events

import (
	"context"
	"errors"

	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const ContactKeyRotationEventHandlerName = "contact.key.rotation.queue"

type ContactKeyRotationQueue struct {
	cfg         *config.ProfileConfig
	dek         *config.DEK
	contactRepo repository.ContactRepository
}

func NewContactKeyRotationQueue(
	cfg *config.ProfileConfig, dek *config.DEK, contactRepo repository.ContactRepository,
) *ContactKeyRotationQueue {
	return &ContactKeyRotationQueue{
		cfg:         cfg,
		dek:         dek,
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

	if vq.dek.OldKeyID != contact.EncryptionKeyID {
		return nil
	}

	contactDetail, err := contact.DecryptDetail(vq.dek.OldKeyID, vq.dek.OldKey)
	if err != nil {
		return err
	}

	contact.EncryptedDetail, err = util.EncryptValue(vq.dek.Key, []byte(contactDetail))
	if err != nil {
		return err
	}

	contact.EncryptionKeyID = vq.dek.KeyID

	_, err = vq.contactRepo.Update(ctx, contact, "encrypted_detail", "encryption_key_id")
	if err != nil {
		return err
	}

	logger.WithField("contact", contact.ID).
		Debug("successfully processed key rotation for contact")

	return nil
}
