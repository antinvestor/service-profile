package business

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
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
	CreateContact(ctx context.Context, detail string, extra map[string]string) (*models.Contact, error)
	UpdateContact(
		ctx context.Context,
		contactID string,
		profileID string,
		extra map[string]string,
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

	ToAPI(ctx context.Context, contact *models.Contact, partial bool) (*profilev1.ContactObject, error)
}

func NewContactBusiness(_ context.Context, service *frame.Service) ContactBusiness {
	cfg, _ := service.Config().(*config.ProfileConfig)

	return &contactBusiness{
		cfg:               cfg,
		service:           service,
		contactRepository: repository.NewContactRepository(service),
		verificationRepository: repository.NewVerificationRepository(
			service,
		),
	}
}

type contactBusiness struct {
	service                *frame.Service
	cfg                    *config.ProfileConfig
	contactRepository      repository.ContactRepository
	verificationRepository repository.VerificationRepository
}

func (cb *contactBusiness) ToAPI(
	_ context.Context,
	contact *models.Contact,
	partial bool,
) (*profilev1.ContactObject, error) {
	contactObject := profilev1.ContactObject{
		Id:     contact.ID,
		Detail: contact.Detail,
	}

	contactTypeID, ok := profilev1.ContactType_value[contact.ContactType]
	if !ok {
		return nil, service.ErrContactTypeNotValid
	}
	contactObject.Type = profilev1.ContactType(contactTypeID)

	communicationLevel, ok := profilev1.CommunicationLevel_value[contact.CommunicationLevel]
	if !ok {
		communicationLevel = int32(profilev1.CommunicationLevel_ALL)
	}
	contactObject.CommunicationLevel = profilev1.CommunicationLevel(communicationLevel)

	contactObject.Verified = false
	if !partial {
		contactObject.Verified = contact.VerificationID != ""
	}

	return &contactObject, nil
}

func (cb *contactBusiness) ContactTypeFromDetail(_ context.Context, detail string) (string, error) {
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
	return cb.contactRepository.GetByID(ctx, contactID)
}

func (cb *contactBusiness) GetByDetail(ctx context.Context, detail string) (*models.Contact, error) {
	return cb.contactRepository.GetByDetail(ctx, detail)
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
	extra map[string]string,
) (*models.Contact, error) {
	contact, err := cb.contactRepository.GetByID(ctx, contactID)
	if err != nil {
		return nil, err
	}
	if contact.ProfileID == "" {
		contact.ProfileID = profileID
	}

	// Initialize Properties map if nil to prevent panic
	if contact.Properties == nil {
		contact.Properties = make(map[string]interface{})
	}

	properties := frame.DBPropertiesFromMap(extra)
	for key, value := range properties {
		if value != contact.Properties[key] {
			contact.Properties[key] = value
		}
	}

	return cb.contactRepository.Save(ctx, contact)
}

func (cb *contactBusiness) CreateContact(
	ctx context.Context,
	detail string,
	extra map[string]string,
) (*models.Contact, error) {
	contact := &models.Contact{}

	detail = strings.ToLower(strings.TrimSpace(detail))

	var err error
	contact.ContactType, err = cb.ContactTypeFromDetail(ctx, detail)
	if err != nil {
		return nil, err
	}

	contact.Detail = detail
	if extra != nil {
		contact.Properties = frame.DBPropertiesFromMap(extra)
	}
	contact, err = cb.contactRepository.Save(ctx, contact)
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
	logger := cb.service.Log(ctx).WithField("contact", contact)

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

	err := cb.service.Emit(ctx, events.VerificationEventHandlerName, verification)
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
