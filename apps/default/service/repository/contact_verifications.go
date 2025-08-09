package repository

import (
	"context"
	"strings"

	"github.com/pitabwire/frame"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type verificationRepository struct {
	service *frame.Service
}

func NewVerificationRepository(service *frame.Service) VerificationRepository {
	repo := verificationRepository{
		service: service,
	}
	return &repo
}

func (vr *verificationRepository) GetByID(
	ctx context.Context,
	verificationID string,
) (*models.Verification, error) {
	verification := &models.Verification{}
	err := vr.service.DB(ctx, false).First(verification, "id = ?", verificationID).Error
	return verification, err
}

func (vr *verificationRepository) Save(ctx context.Context, verification *models.Verification) error {
	err := vr.service.DB(ctx, false).Create(verification).Error
	if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return err
	}
	return nil
}

func (vr *verificationRepository) GetAttempts(
	ctx context.Context,
	verificationID string,
) ([]*models.VerificationAttempt, error) {
	verificationAttemptList := make([]*models.VerificationAttempt, 0)
	err := vr.service.DB(ctx, true).
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
	err := vr.service.DB(ctx, false).Create(verificationAttempt).Error
	if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return err
	}
	return nil
}
