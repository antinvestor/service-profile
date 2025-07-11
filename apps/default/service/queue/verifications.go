package queue

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	commonv1 "github.com/antinvestor/apis/go/common/v1"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

type VerificationsQueueHandler struct {
	Service         *frame.Service
	ContactRepo     repository.ContactRepository
	NotificationCli *notificationv1.NotificationClient
}

func (vq *VerificationsQueueHandler) Handle(ctx context.Context, _ map[string]string, payload []byte) error {
	var verification models.Verification
	err := json.Unmarshal(payload, &verification) // Added & to pass pointer correctly
	if err != nil {
		slog.Error("Failed to unmarshal verification payload", "error", err)
		return err
	}

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	contact, err := vq.ContactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		return err
	}

	err = vq.ContactRepo.VerificationSave(ctx, &verification)
	if err != nil {
		return err
	}

	profileConfig, ok := vq.Service.Config().(*config.ProfileConfig)
	if !ok {
		return errors.New("invalid service configuration")
	}

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

	_, err = vq.NotificationCli.Send(ctx, []*notificationv1.Notification{nMessages})
	return err
}
