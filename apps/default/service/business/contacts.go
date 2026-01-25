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
	GetByID(ctx context.Context, contactID string) (*models.Contact, error)
	GetByDetail(ctx context.Context, detailList ...string) ([]*models.Contact, error)
	GetByDetailMap(ctx context.Context, detailList ...string) (map[string]*models.Contact, error)
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

func NewContactBusiness(_ context.Context, cfg *config.ProfileConfig, dek *config.DEK,
	evtMan frevents.Manager, contactRepository repository.ContactRepository,
	verificationRepository repository.VerificationRepository) ContactBusiness {
	return &contactBusiness{
		cfg:                    cfg,
		dek:                    dek,
		eventsMan:              evtMan,
		contactRepository:      contactRepository,
		verificationRepository: verificationRepository,
	}
}

type contactBusiness struct {
	cfg                    *config.ProfileConfig
	dek                    *config.DEK
	eventsMan              frevents.Manager
	contactRepository      repository.ContactRepository
	verificationRepository repository.VerificationRepository
}

func ContactTypeFromDetail(ctx context.Context, detail string) (string, error) {
	normalizedDetail := Normalize(ctx, detail)

	if EmailPattern.MatchString(normalizedDetail) {
		return profilev1.ContactType_EMAIL.String(), nil
	}

	possibleNumber, err := libphonenumber.Parse(normalizedDetail, "")
	if err == nil && libphonenumber.IsValidNumber(possibleNumber) {
		return profilev1.ContactType_MSISDN.String(), nil
	}

	return "", connect.NewError(connect.CodeInvalidArgument, errors.New("contact details are invalid"))
}

func Normalize(_ context.Context, detail string) string {
	return strings.ToLower(strings.TrimSpace(detail))
}

func (cb *contactBusiness) GetByID(ctx context.Context, contactID string) (*models.Contact, error) {
	contact, err := cb.contactRepository.GetByID(ctx, contactID)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (cb *contactBusiness) GetByDetail(ctx context.Context, detailList ...string) ([]*models.Contact, error) {
	var lookUpTokenList [][]byte

	for _, detail := range detailList {
		normalizedDetail := Normalize(ctx, detail)

		token := util.ComputeLookupToken(cb.dek.LookUpKey, normalizedDetail)
		lookUpTokenList = append(lookUpTokenList, token)
	}
	contact, err := cb.contactRepository.GetByLookupToken(ctx, lookUpTokenList...)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

// GetByDetailMap returns a map of detail to contact for efficient bulk lookups.
func (cb *contactBusiness) GetByDetailMap(
	ctx context.Context,
	detailList ...string,
) (map[string]*models.Contact, error) {
	contacts, err := cb.GetByDetail(ctx, detailList...)
	if err != nil {
		return nil, err
	}

	// Create map for efficient O(1) lookups
	contactMap := make(map[string]*models.Contact)
	for _, contact := range contacts {
		// Decrypt the detail to use as map key
		detail, decryptErr := contact.DecryptDetail(cb.dek.KeyID, cb.dek.Key)
		if decryptErr != nil {
			// Skip contacts that can't be decrypted
			continue
		}
		contactMap[detail] = contact
	}

	return contactMap, nil
}

func (cb *contactBusiness) GetByProfile(ctx context.Context, profileID string) ([]*models.Contact, error) {
	if profileID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("profile ID is empty"))
	}

	contacts, err := cb.contactRepository.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return contacts, nil
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
	normalizedDetail := Normalize(ctx, detail)

	contactType, err := ContactTypeFromDetail(ctx, normalizedDetail)
	if err != nil {
		return nil, err
	}

	lookupToken := util.ComputeLookupToken(cb.dek.LookUpKey, normalizedDetail)

	encryptedDetail, err := util.EncryptValue(cb.dek.Key, []byte(normalizedDetail))
	if err != nil {
		return nil, err
	}

	contact := &models.Contact{
		EncryptedDetail: encryptedDetail,
		EncryptionKeyID: cb.cfg.DEKActiveKeyID,
		LookUpToken:     lookupToken,
		ContactType:     contactType,
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
) (*models.Contact, error) {
	contact, err := cb.contactRepository.DelinkFromProfile(ctx, contactID, profileID)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (cb *contactBusiness) GetVerification(
	ctx context.Context,
	verificationID string,
) (*models.Verification, error) {
	verification, err := cb.verificationRepository.GetByID(ctx, verificationID)
	if err != nil {
		return nil, err
	}
	return verification, nil
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
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no contact specified"))
	}

	if durationToExpiry == 0 {
		durationToExpiry = time.Duration(cb.cfg.VerificationPinExpiryTimeInSec)
	}

	expiryTime := time.Now().Add(durationToExpiry)

	if code == "" {
		code = util.RandomNumericString(cb.cfg.LengthOfVerificationCode)
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
	attempts, err := cb.verificationRepository.GetAttempts(ctx, verificationID)
	if err != nil {
		return nil, err
	}
	return attempts, nil
}
