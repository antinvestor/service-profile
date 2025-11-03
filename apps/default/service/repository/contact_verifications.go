package repository

import (
	"context"
	"strings"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type verificationRepository struct {
	datastore.BaseRepository[*models.Verification]
}

func NewVerificationRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) VerificationRepository {
	repo := verificationRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Verification](
			ctx, dbPool, workMan, func() *models.Verification { return &models.Verification{} },
		),
	}
	return &repo
}

func (vr *verificationRepository) GetByID(
	ctx context.Context,
	verificationID string,
) (*models.Verification, error) {
	verification := &models.Verification{}
	err := vr.Pool().DB(ctx, false).First(verification, "id = ?", verificationID).Error
	return verification, err
}

func (vr *verificationRepository) Save(ctx context.Context, verification *models.Verification) error {
	err := vr.Pool().DB(ctx, false).Create(verification).Error
	if err != nil {
		if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return err
		}
	}
	return nil
}

func (vr *verificationRepository) GetAttempts(
	ctx context.Context,
	verificationID string,
) ([]*models.VerificationAttempt, error) {
	verificationAttemptList := make([]*models.VerificationAttempt, 0)
	err := vr.Pool().DB(ctx, true).
		Preload(clause.Associations).
		Where("verification_id = ?", verificationID).
		Find(&verificationAttemptList).
		Error
	return verificationAttemptList, err
}

func (vr *verificationRepository) SaveAttempt(
	ctx context.Context,
	verificationAttempt *models.VerificationAttempt,
) error {
	err := vr.Pool().DB(ctx, false).Create(verificationAttempt).Error
	if err != nil {
		if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return err
		}
	}
	return nil
}
