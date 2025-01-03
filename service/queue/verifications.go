package queue

import (
	"context"
	"encoding/json"
	commonv1 "github.com/antinvestor/apis/go/common/v1"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
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

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

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

	recipient := &commonv1.ContactLink{
		ProfileType: "Profile",
		ProfileId:   verification.ProfileID,
		ContactId:   contact.ID,
	}

	nMessages := &notificationv1.Notification{
		Recipient: recipient,
		Payload:   variables,
		Language:  contact.Language,
		Template:  profileConfig.MessageTemplateContactVerification,
		OutBound:  true,
	}

	_, err = vq.NotificationCli.Send(ctx, nMessages)
	return err
}
