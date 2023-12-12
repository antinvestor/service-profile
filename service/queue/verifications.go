package queue

import (
	"context"
	"encoding/json"
	notificationv1 "github.com/antinvestor/apis/notification/v1"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
)

type VerificationsQueueHandler struct {
	Service         *frame.Service
	ContactRepo     repository.ContactRepository
	NotificationCli *notificationv1.NotificationClient
}

func (vq *VerificationsQueueHandler) Handle(ctx context.Context, _ map[string]string, payload []byte) error {
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

	profileConfig := vq.Service.Config().(*config.ProfileConfig)

	variables := make(map[string]string)
	variables["pin"] = verification.Pin
	variables["linkHash"] = verification.LinkHash
	variables["expiryDate"] = verification.ExpiresAt.String()
	_, err = vq.NotificationCli.Send(ctx, profileConfig.SystemAccessID,
		contact.ID, "", contact.Language,
		profileConfig.MessageTemplateContactVerification, variables)
	return err
}
