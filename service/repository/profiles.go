package repository

import (
	"context"
	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/pitabwire/frame"
)

type profileRepository struct {
	service *frame.Service
}

func (pr *profileRepository) GetTypeByID(ctx context.Context, profileTypeId string) (*models.ProfileType, error) {
	profileType := &models.ProfileType{}
	err := pr.service.DB(ctx, true).First(profileType, "id = ?", profileTypeId).Error
	return profileType, err
}

func (pr *profileRepository) GetTypeByUID(ctx context.Context, profileType profilev1.ProfileType) (*models.ProfileType, error) {

	profileTypeUId := models.ProfileTypeIDMap[profileType]
	profileTypeM := &models.ProfileType{}
	err := pr.service.DB(ctx, true).First(profileTypeM, "uid = ?", profileTypeUId).Error
	return profileTypeM, err
}

func (pr *profileRepository) GetByID(ctx context.Context, id string) (*models.Profile, error) {

	emptyClaims := &frame.AuthenticationClaims{}
	emptyCtx := emptyClaims.ClaimsToContext(ctx)

	profile := &models.Profile{}
	err := pr.service.DB(emptyCtx, true).First(profile, "id = ?", id).Error
	return profile, err
}

func (pr *profileRepository) Save(ctx context.Context, tenant *models.Profile) error {
	return pr.service.DB(ctx, false).Save(tenant).Error
}

func (pr *profileRepository) Delete(ctx context.Context, id string) error {

	profile, err := pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return pr.service.DB(ctx, false).Delete(profile).Error
}

func NewProfileRepository(service *frame.Service) ProfileRepository {
	profileRepository := profileRepository{
		service: service,
	}
	return &profileRepository
}
