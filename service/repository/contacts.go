package repository

import (
	"context"
	"errors"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
	"strings"
)

type contactRepository struct {
	service *frame.Service
}

func (cr *contactRepository) GetVerificationByContactID(ctx context.Context, contactID string) (*models.Verification, error) {
	verification := &models.Verification{}
	err := cr.service.DB(ctx, false).Last(verification, "contact_id = ?", contactID).Error
	return verification, err
}

func (cr *contactRepository) VerificationSave(ctx context.Context, verification *models.Verification) error {
	return cr.service.DB(ctx, false).FirstOrCreate(verification).Error
}

func (cr *contactRepository) VerificationAttemptSave(ctx context.Context, attempt *models.VerificationAttempt) error {
	return cr.service.DB(ctx, false).Save(attempt).Error
}

func (cr *contactRepository) GetByID(ctx context.Context, id string) (*models.Contact, error) {
	contact := &models.Contact{}
	err := cr.service.DB(ctx, true).First(contact, "id = ?", id).Error
	return contact, err
}

func (cr *contactRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error) {
	contactList := make([]*models.Contact, 0)
	err := cr.service.DB(ctx, true).Find(&contactList, "profile_id = ?", profileID).Error
	return contactList, err
}

func (cr *contactRepository) GetByDetail(ctx context.Context, detail string) (*models.Contact, error) {
	contact := &models.Contact{}

	detail = strings.ToLower(strings.TrimSpace(detail))
	if err := cr.service.DB(ctx, true).First(contact, " tokens @@ to_tsquery(?)", detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrorContactDoesNotExist
		}
		return nil, err
	}

	return contact, nil
}

func (cr *contactRepository) Save(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	if contact.ID == "" {
		contact.GenID(ctx)
		err := cr.service.DB(ctx, false).Model(contact).Create(map[string]interface{}{
			"ID":                   contact.ID,
			"ContactTypeID":        contact.ContactTypeID,
			"CommunicationLevelID": contact.CommunicationLevelID,
			"ProfileID":            contact.ProfileID,
			"Detail":               contact.Detail,
			"Nonce":                contact.Nonce,
			"Tokens":               gorm.Expr("to_tsvector(?)", contact.Tokens),
		}).Error
		return contact, err
	}

	err := cr.service.DB(ctx, false).Save(contact).Error
	return contact, err
}

func (cr *contactRepository) Delete(ctx context.Context, id string) error {
	contact, err := cr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return cr.service.DB(ctx, false).Delete(contact).Error
}

func (cr *contactRepository) ContactType(ctx context.Context,
	contactType profilev1.ContactType) (*models.ContactType, error) {
	uid := models.ContactTypeUIDMap[contactType]
	ct := &models.ContactType{}
	err := cr.service.DB(ctx, true).First(ct, "uid = ?", uid).Error
	return ct, err

}

func (cr *contactRepository) ContactTypeByID(ctx context.Context, contactTypeID string) (*models.ContactType, error) {
	ct := &models.ContactType{}
	err := cr.service.DB(ctx, true).First(ct, "id = ?", contactTypeID).Error
	return ct, err

}

func (cr *contactRepository) CommunicationLevel(ctx context.Context,
	communicationLevel profilev1.CommunicationLevel) (*models.CommunicationLevel, error) {

	uid := models.CommunicationLevelUIDMap[communicationLevel]
	cl := &models.CommunicationLevel{}
	err := cr.service.DB(ctx, true).First(cl, "uid = ?", uid).Error
	return cl, err

}

func (cr *contactRepository) CommunicationLevelByID(ctx context.Context,
	communicationLevelID string) (*models.CommunicationLevel, error) {

	cl := &models.CommunicationLevel{}
	err := cr.service.DB(ctx, true).First(cl, "id = ?", communicationLevelID).Error
	return cl, err

}

func NewContactRepository(service *frame.Service) ContactRepository {
	contactRepository := contactRepository{
		service: service,
	}
	return &contactRepository
}
