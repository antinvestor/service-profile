package events

import (
	"context"
	"errors"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/security"
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

	logger := util.Log(ctx).WithField("payload", verification.GetID()).WithField("type", vq.Name())

	ctx = security.SkipTenancyChecksOnClaims(ctx)

	contact, err := vq.contactRepo.GetByID(ctx, verification.ContactID)
	if err != nil {
		return err
	}

	err = vq.verificationRepo.Create(ctx, verification)
	if err != nil {
		logger.WithError(err).Error("Failed to save verification")
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
		Recipient: recipient,
		Payload:   variablePayload,
		Language:  contact.Language,
		Template:  vq.cfg.MessageTemplateContactVerification,
		OutBound:  true,
	}

	req := connect.NewRequest(&notificationv1.SendRequest{
		Data: []*notificationv1.Notification{nMessages},
	})

	_, err = vq.notificationCli.Send(ctx, req)
	if err != nil {
		logger.WithError(err).Error("Failed to send out verification")
		return err
	}

	return nil
}
