package events

import (
	"context"
	"errors"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const VerificationEventHandlerName = "contact.verification.queue"

type ContactVerificationQueue struct {
	cfg              *config.ProfileConfig
	contactRepo      repository.ContactRepository
	verificationRepo repository.VerificationRepository
	notificationCli  notificationv1connect.NotificationServiceClient
}

func NewContactVerificationQueue(
	cfg *config.ProfileConfig, contactRepo repository.ContactRepository,
	verificationRepo repository.VerificationRepository, notificationCli notificationv1connect.NotificationServiceClient,
) *ContactVerificationQueue {
	return &ContactVerificationQueue{
		cfg:              cfg,
		contactRepo:      contactRepo,
		verificationRepo: verificationRepo,
		notificationCli:  notificationCli,
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
		return errors.New(" invalid payload type, expected *models.Verification")
	}

	if notification.GetID() == "" {
		return errors.New(" invalid payload type, expected Id on *models.Verification")
	}
	return nil
}

func (vq *ContactVerificationQueue) Execute(ctx context.Context, payload any) error {
	verification, ok := payload.(*models.Verification)
	if !ok {
		return errors.New(" invalid payload type, expected *models.Verification")
	}

	logger := util.Log(ctx).WithFields(map[string]any{
		"verification_id": verification.GetID(),
		"type":            vq.Name(),
	})

	contact, err := vq.contactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		return err
	}

	err = vq.verificationRepo.Create(ctx, verification)
	if err != nil {
		if data.ErrorIsDuplicateKey(err) {
			logger.Debug("verification already exists, skipping duplicate")
			return nil
		}
		logger.WithError(err).Error("failed to save verification")
		return err
	}

	variables := make(data.JSONMap)
	variables["verification_id"] = verification.GetID()
	variables["code"] = verification.Code
	variables["expiryDate"] = verification.ExpiresAt.String()

	variablePayload, _ := structpb.NewStruct(variables)

	recipient := &commonv1.ContactLink{
		ProfileType: "Profile",
		ProfileId:   verification.ProfileID,
		ContactId:   contact.ID,
	}

	nMessages := &notificationv1.Notification{
		Recipient:   recipient,
		Payload:     variablePayload,
		Language:    contact.Language,
		Template:    vq.cfg.MessageTemplateContactVerification,
		OutBound:    true,
		AutoRelease: true,
	}

	req := connect.NewRequest(&notificationv1.SendRequest{
		Data: []*notificationv1.Notification{nMessages},
	})

	resp, err := vq.notificationCli.Send(ctx, req)
	if err != nil {
		logger.WithField("contact_id", contact.ID).WithError(err).Error("failed to send verification notification")
		return err
	}

	if resp == nil {
		logger.Debug("notification response is nil, likely in test mode")
		return nil
	}

	for resp.Receive() {
		err = resp.Err()
		if err != nil {
			logger.WithField("contact_id", contact.ID).
				WithError(err).Error("failed submitting verification for contact")
			return err
		}

		logger.WithField("contact_id", contact.ID).Debug("verification notification submitted")
	}
	return nil
}
