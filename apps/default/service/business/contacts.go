package business

import (
	"context"
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/ttacon/libphonenumber"

	"github.com/pitabwire/frame"
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
	GetVerification(ctx context.Context, contactID string) (*models.Verification, error)

	ToAPI(ctx context.Context, contact *models.Contact, partial bool) (*profilev1.ContactObject, error)
}

func NewContactBusiness(_ context.Context, service *frame.Service) ContactBusiness {
	contactRepo := repository.NewContactRepository(service)
	return &contactBusiness{
		service:           service,
		contactRepository: contactRepo,
	}
}

type contactBusiness struct {
	service           *frame.Service
	contactRepository repository.ContactRepository
}

func (cb *contactBusiness) ToAPI(
	ctx context.Context,
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
		verification, _ := cb.GetVerification(ctx, contact.GetID())
		contactObject.Verified = verification != nil
	}

	return &contactObject, nil
}

func (cb *contactBusiness) ContactTypeFromDetail(_ context.Context, detail string) (string, error) {
	if EmailPattern.MatchString(detail) {
		return profilev1.ContactType_EMAIL.String(), nil
	} else {
		possibleNumber, err := libphonenumber.Parse(detail, "")

		if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
			return profilev1.ContactType_MSISDN.String(), err
		}
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

func (cb *contactBusiness) GetVerification(ctx context.Context, contactID string) (*models.Verification, error) {
	return cb.contactRepository.GetVerificationByContactID(ctx, contactID)
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

	err = cb.VerifyContact(ctx, contact)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (cb *contactBusiness) RemoveContact(ctx context.Context, contactID, profileID string) (*models.Contact, error) {
	return cb.contactRepository.DelinkFromProfile(ctx, contactID, profileID)
}

func (cb *contactBusiness) VerifyContact(ctx context.Context, contact *models.Contact) error {
	if contact.ProfileID == "" {
		return nil
	}

	cfg := cb.service.Config().(*config.ProfileConfig)
	expiryTime := time.Now().Add(time.Duration(cfg.VerificationPinExpiryTimeInSec))

	verification := &models.Verification{
		ProfileID: contact.ProfileID,
		ContactID: contact.ID,
		Pin:       GeneratePin(cfg.LengthOfVerificationPin),
		LinkHash:  GeneratePin(cfg.LengthOfVerificationLinkHash),
		ExpiresAt: &expiryTime,
	}

	verification.GenID(ctx)

	return cb.service.Publish(ctx, cfg.QueueVerificationName, verification)
}

// GeneratePin returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GeneratePin(n int) string {
	if n <= 0 {
		return ""
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
