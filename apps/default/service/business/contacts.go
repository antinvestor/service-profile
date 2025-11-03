package business

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame/data"
	frevents "github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"
	"github.com/ttacon/libphonenumber"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service"
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
	GetByID(ctx context.Context, contactID string) (*models.Contact, error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, error)
	GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, error)
	CreateContact(ctx context.Context, detail string, extra data.JSONMap) (*models.Contact, error)
	UpdateContact(
		ctx context.Context,
		contactID string,
		profileID string,
		extra data.JSONMap,
	) (*models.Contact, error)
	RemoveContact(ctx context.Context, contactID, profileID string) (*models.Contact, error)
	VerifyContact(
		ctx context.Context,
		contact *models.Contact,
		verificationID string,
		code string,
		duration time.Duration,
	) (*models.Verification, error)
	GetVerification(ctx context.Context, verificationID string) (*models.Verification, error)
	GetVerificationAttempts(ctx context.Context, verificationID string) ([]*models.VerificationAttempt, error)
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

func ContactTypeFromDetail(_ context.Context, detail string) (string, error) {
	if EmailPattern.MatchString(detail) {
		return profilev1.ContactType_EMAIL.String(), nil
	}

	possibleNumber, err := libphonenumber.Parse(detail, "")
	if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
		return profilev1.ContactType_MSISDN.String(), nil
	}

	return "", service.ErrContactDetailsNotValid
}

func (cb *contactBusiness) GetByID(ctx context.Context, contactID string) (*models.Contact, error) {
	contact, err := cb.contactRepository.GetByID(ctx, contactID)
	if err != nil {
		if data.ErrorIsNoRows(err) {
			return nil, service.ErrContactDoesNotExist
		}
	}
	return contact, err
}

func (cb *contactBusiness) GetByDetail(ctx context.Context, detail string) (*models.Contact, error) {
	contact, err := cb.contactRepository.GetByDetail(ctx, detail)
	if err != nil {
		if data.ErrorIsNoRows(err) {
			return nil, service.ErrContactDoesNotExist
		}
	}
	return contact, err
}

func (cb *contactBusiness) GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, error) {
	if profileID == "" {
		return nil, errors.New("profile ID is empty")
	}

	return cb.contactRepository.GetByProfileID(ctx, profileID)
}

func (cb *contactBusiness) UpdateContact(
	ctx context.Context,
	contactID string,
	profileID string,
	extra data.JSONMap,
) (*models.Contact, error) {
	contact, err := cb.contactRepository.GetByID(ctx, contactID)
	if err != nil {
		return nil, err
	}
	if contact.ProfileID == "" {
		contact.ProfileID = profileID
	}

	contact.Properties = contact.Properties.Update(extra)

	_, err = cb.contactRepository.Update(ctx, contact, "profile_id", "properties")
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (cb *contactBusiness) CreateContact(
	ctx context.Context,
	detail string,
	extra data.JSONMap,
) (*models.Contact, error) {
	detail = strings.ToLower(strings.TrimSpace(detail))

	contactType, err := ContactTypeFromDetail(ctx, detail)
	if err != nil {
		return nil, err
	}

	contact, err := cb.GetByDetail(ctx, detail)
	if err == nil {
		return contact, nil
	}

	if !errors.Is(err, service.ErrContactDoesNotExist) {
		return nil, err
	}

	contact = &models.Contact{
		Detail:      detail,
		ContactType: contactType,
	}

	contact.Properties = contact.Properties.Update(extra)

	err = cb.contactRepository.Create(ctx, contact)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (cb *contactBusiness) RemoveContact(ctx context.Context, contactID, profileID string) (*models.Contact, error) {
	return cb.contactRepository.DelinkFromProfile(ctx, contactID, profileID)
}

func (cb *contactBusiness) GetVerification(ctx context.Context, verificationID string) (*models.Verification, error) {
	return cb.verificationRepository.GetByID(ctx, verificationID)
}

func (cb *contactBusiness) VerifyContact(
	ctx context.Context,
	contact *models.Contact,
	verificationID string,
	code string,
	durationToExpiry time.Duration,
) (*models.Verification, error) {
	logger := util.Log(ctx).WithField("contact", contact)

	if contact == nil {
		return nil, errors.New("no contact specified")
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
) ([]*models.VerificationAttempt, error) {
	return cb.verificationRepository.GetAttempts(ctx, verificationID)
}
