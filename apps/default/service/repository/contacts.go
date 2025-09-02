package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/default/service"
	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type contactRepository struct {
	service *frame.Service
}

func (cr *contactRepository) GetVerificationByID(
	ctx context.Context,
	verificationID string,
) (*models.Verification, error) {
	verification := &models.Verification{}
	err := cr.service.DB(ctx, false).First(verification, "id = ?", verificationID).Error
	return verification, err
}

func (cr *contactRepository) VerificationSave(ctx context.Context, verification *models.Verification) error {
	return cr.service.DB(ctx, false).FirstOrCreate(verification).Error
}

func (cr *contactRepository) VerificationAttemptSave(ctx context.Context, attempt *models.VerificationAttempt) error {
	return cr.service.DB(ctx, false).Save(attempt).Error
}

func (cr *contactRepository) GetByID(ctx context.Context, id string) (*models.Contact, error) {
	var contact models.Contact
	err := cr.service.DB(ctx, true).First(&contact, "id = ?", id).Error
	return &contact, err
}

func (cr *contactRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error) {
	contactList := make([]*models.Contact, 0)
	err := cr.service.DB(ctx, true).Where("profile_id = ?", profileID).Find(&contactList).Error
	return contactList, err
}

func (cr *contactRepository) GetByDetail(ctx context.Context, detail string) (*models.Contact, error) {
	contact := &models.Contact{}

	detail = strings.ToLower(strings.TrimSpace(detail))
	if err := cr.service.DB(ctx, true).First(contact, " detail = ?", detail).Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (cr *contactRepository) Save(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	db := cr.service.DB(ctx, false)

	if contact.ID == "" {
		contact.GenID(ctx)
	}

	if contact.Version == 0 {
		// Create new contact
		if err := db.Create(contact).Error; err != nil {
			return nil, fmt.Errorf("failed to create contact: %w", err)
		}
	} else {
		// Update existing contact
		if err := db.Save(contact).Error; err != nil {
			return nil, fmt.Errorf("failed to update contact: %w", err)
		}
	}

	return contact, nil
}

func (cr *contactRepository) DelinkFromProfile(ctx context.Context, id, profileID string) (*models.Contact, error) {
	contact, err := cr.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if profileID != contact.ProfileID {
		return nil, service.ErrContactProfileNotValid
	}

	contact.ProfileID = ""

	err = cr.service.DB(ctx, false).Save(contact).Error
	return contact, err
}

func (cr *contactRepository) Delete(ctx context.Context, id string) error {
	contact, err := cr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return cr.service.DB(ctx, false).Delete(contact).Error
}

func NewContactRepository(service *frame.Service) ContactRepository {
	repo := contactRepository{
		service: service,
	}
	return &repo
}
