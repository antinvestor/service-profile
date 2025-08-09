package events

import (
	"context"
	"errors"

	commonv1 "github.com/antinvestor/apis/go/common/v1"
	notificationv1 "github.com/antinvestor/apis/go/notification/v1"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const VerificationEventHandlerName = "contact.verification.queue"

type ContactVerificationQueue struct {
	Service          *frame.Service
	ContactRepo      repository.ContactRepository
	VerificationRepo repository.VerificationRepository
	NotificationCli  *notificationv1.NotificationClient
}

func NewContactVerificationQueue(
	service *frame.Service,
	notificationCli *notificationv1.NotificationClient,
) *ContactVerificationQueue {
	return &ContactVerificationQueue{
		Service:          service,
		ContactRepo:      repository.NewContactRepository(service),
		VerificationRepo: repository.NewVerificationRepository(service),
		NotificationCli:  notificationCli,
	}
}

func (vq *ContactVerificationQueue) Name() string {
	return VerificationEventHandlerName
}

func (vq *ContactVerificationQueue) PayloadType() any {
	return &models.Verification{}
}

func (vq *ContactVerificationQueue) Validate(_ context.Context, payload any) error {
	notification, ok := payload.(*models.Verification)
	if !ok {
		return errors.New(" payload is not of type models.Verification")
	}

	if notification.GetID() == "" {
		return errors.New(" verification Id should already have been set ")
	}
	return nil
}

func (vq *ContactVerificationQueue) Execute(ctx context.Context, payload any) error {
	verification, ok := payload.(*models.Verification)
	if !ok {
		return errors.New(" payload is not of type models.Verification")
	}

	logger := vq.Service.Log(ctx).WithField("payload", verification.GetID()).WithField("type", vq.Name())

	ctx = frame.SkipTenancyChecksOnClaims(ctx)

	contact, err := vq.ContactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		return err
	}

	err = vq.VerificationRepo.Save(ctx, verification)
	if err != nil {
		logger.WithError(err).Error("Failed to save verification")
		return err
	}

	profileConfig, ok := vq.Service.Config().(*config.ProfileConfig)
	if !ok {
		return errors.New("invalid service configuration")
	}

	variables := make(map[string]string)
	variables["code"] = verification.Code
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
	if err != nil {
		logger.WithError(err).Error("Failed to send out verification")
		return err
	}

	return nil
}
