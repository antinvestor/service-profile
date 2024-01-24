package business

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/antinvestor/service-profile/utils"
	"github.com/pitabwire/frame"
	"github.com/ttacon/libphonenumber"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	EmailPattern = regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type ContactBusiness interface {
	GetByID(ctx context.Context, contactID string) (*models.Contact, error)
	GetByDetail(ctx context.Context, detail string) (*models.Contact, error)
	GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, error)
	CreateContact(ctx context.Context, key []byte, profileID string, detail string) error

	GetVerification(ctx context.Context, contactID string) (*models.Verification, error)

	ToAPI(ctx context.Context, contact *models.Contact, key []byte) (*profilev1.ContactObject, error)
}

func NewContactBusiness(ctx context.Context, service *frame.Service) ContactBusiness {
	contactRepo := repository.NewContactRepository(service)
	return &contactBusiness{
		service:    service,
		contactRep: contactRepo,
	}
}

type contactBusiness struct {
	service    *frame.Service
	contactRep repository.ContactRepository
}

func (cb *contactBusiness) ToAPI(ctx context.Context, contact *models.Contact, key []byte) (*profilev1.ContactObject, error) {

	detail, err := utils.AesDecrypt(key, contact.Nonce, contact.Detail)
	if err != nil {
		return nil, err
	}

	contactType, err := cb.contactRep.ContactTypeByID(ctx, contact.ContactTypeID)
	if err != nil {
		return nil, err
	}

	communicationLevel, err := cb.contactRep.CommunicationLevelByID(ctx, contact.CommunicationLevelID)
	if err != nil {
		return nil, err
	}

	isVerified := false
	verification, _ := cb.GetVerification(ctx, contact.GetID())

	isVerified = verification != nil

	contactObject := profilev1.ContactObject{
		Id:                 contact.ID,
		Detail:             detail,
		Verified:           isVerified,
		Type:               profilev1.ContactType(contactType.UID),
		CommunicationLevel: profilev1.CommunicationLevel(communicationLevel.UID),
	}

	return &contactObject, err

}

func (cb *contactBusiness) FromDetail(ctx context.Context, detail string) (*models.ContactType, error) {

	if EmailPattern.MatchString(detail) {
		ct, err := cb.contactRep.ContactType(ctx, profilev1.ContactType_EMAIL)
		return ct, err

	} else {

		possibleNumber, err := libphonenumber.Parse(detail, "")

		if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
			ct, err := cb.contactRep.ContactType(ctx, profilev1.ContactType_PHONE)
			return ct, err

		}
	}

	return nil, service.ErrorContactDetailsNotValid
}

func (cb *contactBusiness) GetByID(ctx context.Context, contactID string) (*models.Contact, error) {
	return cb.contactRep.GetByID(ctx, contactID)
}

func (cb *contactBusiness) GetByDetail(ctx context.Context, detail string) (*models.Contact, error) {
	return cb.contactRep.GetByDetail(ctx, detail)
}

func (cb *contactBusiness) GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, error) {
	return cb.contactRep.GetByProfileID(ctx, profileID)
}

func (cb *contactBusiness) GetVerification(ctx context.Context, contactID string) (*models.Verification, error) {
	return cb.contactRep.GetVerificationByContactID(ctx, contactID)
}
func (cb *contactBusiness) CreateContact(ctx context.Context, key []byte, profileID string, detail string) error {

	contact := &models.Contact{}

	detail = strings.ToLower(strings.TrimSpace(detail))

	ct, err := cb.FromDetail(ctx, detail)
	if err != nil {
		return err
	}
	contact.ContactTypeID = ct.ID
	contact.ContactType = ct

	cl, err := cb.contactRep.CommunicationLevel(ctx, profilev1.CommunicationLevel_ALL)
	if err != nil {
		return err
	}
	contact.CommunicationLevelID = cl.ID
	contact.CommunicationLevel = cl

	contact.ProfileID = profileID

	contact.Detail, contact.Nonce, err = utils.AesEncrypt(key, detail)
	if err != nil {
		return err
	}
	contact.Tokens = detail

	contact, err = cb.contactRep.Save(ctx, contact)
	if err != nil {
		return err
	}

	return cb.VerifyContact(ctx, contact)

}

func (cb *contactBusiness) VerifyContact(ctx context.Context, contact *models.Contact) error {

	profileConfig := cb.service.Config().(*config.ProfileConfig)
	expiryTime := time.Now().Add(time.Duration(profileConfig.VerificationPinExpiryTimeInSec))

	verification := &models.Verification{
		ContactID: contact.ID,
		Pin:       GeneratePin(profileConfig.LengthOfVerificationPin),
		LinkHash:  GeneratePin(profileConfig.LengthOfVerificationLinkHash),
		ExpiresAt: &expiryTime,
	}

	verification.GenID(ctx)

	return cb.service.Publish(ctx, profileConfig.QueueVerificationName, verification)
}

// GeneratePin returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GeneratePin(n int) string {

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
