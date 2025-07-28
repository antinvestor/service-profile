package queue

import (
	"context"
	"encoding/json"
	"errors"

	commonv1 "github.com/antinvestor/apis/go/common/v1"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
)

type VerificationsQueueHandler struct {
	Service         *frame.Service
	ContactRepo     repository.ContactRepository
	NotificationCli *notificationv1.NotificationClient
}

func (vq *VerificationsQueueHandler) Handle(
	ctx context.Context,
	_ map[string]string,
	payload []byte,
) error {
	// Create a new verification object with the proper JSON struct tags
	var verification models.Verification // This struct has JSON tags defined in models.go
	//nolint:musttag // The struct has proper JSON tags in models.go
	err := json.Unmarshal(
		payload,
		&verification,
	)
	if err != nil {
		util.Log(ctx).WithError(err).Error("Failed to unmarshal verification payload")
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
