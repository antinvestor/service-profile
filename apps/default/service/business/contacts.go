package business

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/data"
	frevents "github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"
	"github.com/ttacon/libphonenumber"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

var (
	EmailPattern = regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
)

type ContactBusiness interface {
	GetByID(ctx context.Context, contactID string) (*models.Contact, *connect.Error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, *connect.Error)
	GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, *connect.Error)
	CreateContact(ctx context.Context, detail string, extra data.JSONMap) (*models.Contact, *connect.Error)
	UpdateContact(
		ctx context.Context,
		contactID string,
		profileID string,
		extra data.JSONMap,
	) (*models.Contact, *connect.Error)
	RemoveContact(ctx context.Context, contactID, profileID string) (*models.Contact, *connect.Error)
	VerifyContact(
		ctx context.Context,
		contact *models.Contact,
		verificationID string,
		code string,
		duration time.Duration,
	) (*models.Verification, *connect.Error)
	GetVerification(ctx context.Context, verificationID string) (*models.Verification, *connect.Error)
	GetVerificationAttempts(ctx context.Context, verificationID string) ([]*models.VerificationAttempt, *connect.Error)
}

func NewContactBusiness(_ context.Context, cfg *config.ProfileConfig,
	evtMan frevents.Manager, contactRepository repository.ContactRepository,
	verificationRepository repository.VerificationRepository) ContactBusiness {
	return &contactBusiness{
		cfg:                    cfg,
		eventsMan:              evtMan,
		contactRepository:      contactRepository,
		verificationRepository: verificationRepository,
	}
}

type contactBusiness struct {
	cfg                    *config.ProfileConfig
	eventsMan              frevents.Manager
	contactRepository      repository.ContactRepository
	verificationRepository repository.VerificationRepository
}

func ContactTypeFromDetail(_ context.Context, detail string) (string, *connect.Error) {
	if EmailPattern.MatchString(detail) {
		return profilev1.ContactType_EMAIL.String(), nil
	}

	possibleNumber, err := libphonenumber.Parse(detail, "")
	if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
		return profilev1.ContactType_MSISDN.String(), nil
	}

	return "", connect.NewError(connect.CodeInvalidArgument, errors.New("contact details are invalid"))
}

func (cb *contactBusiness) GetByID(ctx context.Context, contactID string) (*models.Contact, *connect.Error) {
	contact, err := cb.contactRepository.GetByID(ctx, contactID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return contact, nil
}

func (cb *contactBusiness) GetByDetail(ctx context.Context, detail string) (*models.Contact, *connect.Error) {
	contact, err := cb.contactRepository.GetByDetail(ctx, detail)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return contact, nil
}

func (cb *contactBusiness) GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, *connect.Error) {
	if profileID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("profile ID is empty"))
	}

	contacts, err := cb.contactRepository.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return contacts, nil
}

func (cb *contactBusiness) UpdateContact(
	ctx context.Context,
	contactID string,
	profileID string,
	extra data.JSONMap,
) (*models.Contact, *connect.Error) {
	contact, err := cb.contactRepository.GetByID(ctx, contactID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	if contact.ProfileID == "" {
		contact.ProfileID = profileID
	}

	contact.Properties = contact.Properties.Update(extra)

	_, err = cb.contactRepository.Update(ctx, contact, "profile_id", "properties")
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return contact, nil
}

func (cb *contactBusiness) CreateContact(
	ctx context.Context,
	detail string,
	extra data.JSONMap,
) (*models.Contact, *connect.Error) {
	detail = strings.ToLower(strings.TrimSpace(detail))

	contactType, err := ContactTypeFromDetail(ctx, detail)
	if err != nil {
		return nil, err
	}

	contact, getDetailErr := cb.GetByDetail(ctx, detail)
	if getDetailErr == nil {
		return contact, nil
	}

	if getDetailErr.Code() != connect.CodeNotFound {
		return nil, getDetailErr
	}

	contact = &models.Contact{
		Detail:      detail,
		ContactType: contactType,
	}

	contact.Properties = contact.Properties.Update(extra)

	createErr := cb.contactRepository.Create(ctx, contact)
	if createErr != nil {
		return nil, data.ErrorConvertToAPI(createErr)
	}

	return contact, nil
}

func (cb *contactBusiness) RemoveContact(
	ctx context.Context,
	contactID, profileID string,
) (*models.Contact, *connect.Error) {
	contact, err := cb.contactRepository.DelinkFromProfile(ctx, contactID, profileID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return contact, nil
}

func (cb *contactBusiness) GetVerification(
	ctx context.Context,
	verificationID string,
) (*models.Verification, *connect.Error) {
	verification, err := cb.verificationRepository.GetByID(ctx, verificationID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return verification, nil
}

func (cb *contactBusiness) VerifyContact(
	ctx context.Context,
	contact *models.Contact,
	verificationID string,
	code string,
	durationToExpiry time.Duration,
) (*models.Verification, *connect.Error) {
	logger := util.Log(ctx).WithField("contact", contact)

	if contact == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no contact specified"))
	}

	if durationToExpiry == 0 {
		durationToExpiry = time.Duration(cb.cfg.VerificationPinExpiryTimeInSec)
	}

	expiryTime := time.Now().Add(durationToExpiry)

	if code == "" {
		code = util.RandomString(cb.cfg.LengthOfVerificationCode)
	}

	verification := &models.Verification{
		ProfileID: contact.ProfileID,
		ContactID: contact.ID,
		Code:      code,
		ExpiresAt: expiryTime,
	}

	verification.GenID(ctx)

	if verificationID != "" {
		verification.ID = verificationID
	}

	err := cb.eventsMan.Emit(ctx, events.VerificationEventHandlerName, verification)
	if err != nil {
		logger.WithError(err).Error("could not emit verification attempt event")
	}

	return verification, nil
}

func (cb *contactBusiness) GetVerificationAttempts(
	ctx context.Context,
	verificationID string,
) ([]*models.VerificationAttempt, *connect.Error) {
	attempts, err := cb.verificationRepository.GetAttempts(ctx, verificationID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}
	return attempts, nil
}
