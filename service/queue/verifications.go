package queue

import (
	"context"
	"encoding/json"
	napi "github.com/antinvestor/service-notification-api"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
)

type VerificationsQueueHandler struct {
	Service         *frame.Service
	ContactRepo     repository.ContactRepository
	NotificationCli *napi.NotificationClient
}

func (vq *VerificationsQueueHandler) Handle(ctx context.Context, payload []byte) error {

	verification := &models.Verification{}
	err := json.Unmarshal(payload, verification)
	if err != nil {
		return err
	}

	contact, err := vq.ContactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		return err
	}

	err = vq.ContactRepo.VerificationSave(ctx, verification)
	if err != nil {
		return err
	}

	variables := make(map[string]string)
	variables["pin"] = verification.Pin
	variables["linkHash"] = verification.LinkHash
	variables["expiryDate"] = verification.ExpiresAt.String()

	_, err = vq.NotificationCli.Send(ctx, contact.Profile.ID, contact.ID, contact.Language, config.MessageTemplateContactVerification, variables)
	return err

}
