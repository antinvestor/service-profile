package repository

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type contactRepository struct {
	datastore.BaseRepository[*models.Contact]
}

func NewContactRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) ContactRepository {
	repo := contactRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Contact](
			ctx, dbPool, workMan, func() *models.Contact { return &models.Contact{} },
		),
	}
	return &repo
}

func (cr *contactRepository) GetVerificationByID(
	ctx context.Context,
	verificationID string,
) (*models.Verification, error) {
	verification := &models.Verification{}
	err := cr.Pool().DB(ctx, false).First(verification, "id = ?", verificationID).Error
	return verification, err
}

func (cr *contactRepository) VerificationSave(ctx context.Context, verification *models.Verification) error {
	return cr.Pool().DB(ctx, false).FirstOrCreate(verification).Error
}

func (cr *contactRepository) VerificationAttemptSave(ctx context.Context, attempt *models.VerificationAttempt) error {
	return cr.Pool().DB(ctx, false).Save(attempt).Error
}

func (cr *contactRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error) {
	contactList := make([]*models.Contact, 0)
	err := cr.Pool().DB(ctx, true).Where("profile_id = ?", profileID).Find(&contactList).Error
	return contactList, err
}

func (cr *contactRepository) GetByLookupToken(ctx context.Context, lookupTokenList ...[]byte) ([]*models.Contact, error) {
	var contactList []*models.Contact

	if err := cr.Pool().DB(ctx, true).Where(" look_up_token IN ?", lookupTokenList).Find(&contactList).Error; err != nil {
		return nil, err
	}

	return contactList, nil
}

func (cr *contactRepository) DelinkFromProfile(ctx context.Context, id, profileID string) (*models.Contact, error) {
	contact, err := cr.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if profileID != contact.ProfileID {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("contact profile is invalid"))
	}

	contact.ProfileID = ""

	err = cr.Pool().DB(ctx, false).Save(contact).Error
	return contact, err
}
