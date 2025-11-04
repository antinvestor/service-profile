package repository

import (
	"context"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type profileRepository struct {
	datastore.BaseRepository[*models.Profile]
}

func NewProfileRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) ProfileRepository {
	repo := profileRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Profile](
			ctx, dbPool, workMan, func() *models.Profile { return &models.Profile{} },
		),
	}
	return &repo
}

func (pr *profileRepository) GetTypeByID(ctx context.Context, profileTypeID string) (*models.ProfileType, error) {
	profileType := &models.ProfileType{}
	err := pr.Pool().DB(ctx, true).First(profileType, "id = ?", profileTypeID).Error
	return profileType, err
}

func (pr *profileRepository) GetTypeByUID(
	ctx context.Context,
	profileType profilev1.ProfileType,
) (*models.ProfileType, error) {
	profileTypeUID := models.ProfileTypeIDMap[profileType]
	profileTypeM := &models.ProfileType{}
	err := pr.Pool().DB(ctx, true).First(profileTypeM, "uid = ?", profileTypeUID).Error
	return profileTypeM, err
}

func (pr *profileRepository) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	emptyClaims := &security.AuthenticationClaims{}
	emptyCtx := emptyClaims.ClaimsToContext(ctx)

	profile := &models.Profile{}
	err := pr.Pool().DB(emptyCtx, true).Preload(clause.Associations).First(profile, "id = ?", id).Error
	return profile, err
}

func (pr *profileRepository) Save(ctx context.Context, tenant *models.Profile) error {
	return pr.Pool().DB(ctx, false).Save(tenant).Error
}

func (pr *profileRepository) Delete(ctx context.Context, id string) error {
	profile, err := pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return pr.Pool().DB(ctx, false).Delete(profile).Error
}
