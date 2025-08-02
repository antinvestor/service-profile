package repository

import (
	"context"
	"time"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type profileRepository struct {
	service *frame.Service
}

func (pr *profileRepository) Search(
	ctx context.Context,
	query *datastore.SearchQuery,
) (frame.JobResultPipe[[]*models.Profile], error) {

	return datastore.StableSearch[models.Profile](ctx, pr.service, query, func(
		ctx context.Context,
		query *datastore.SearchQuery,
	) ([]*models.Profile, error) {
		var profileList []*models.Profile

		paginator := query.Pagination

		db := pr.service.DB(ctx, true).
			Limit(paginator.Limit).Offset(paginator.Offset)

		if query.Fields != nil {

			startAt, sok := query.Fields["start_date"]
			stopAt, stok := query.Fields["end_date"]
			if sok && startAt != nil && stok && stopAt != nil {
				startDate := startAt.(*time.Time).Format("2020-01-31T00:00:00Z")
				endDate := stopAt.(*time.Time).Format("2020-01-31T00:00:00Z")
				db = db.Where("created_at BETWEEN ? AND ? ", startDate, endDate)
			}

			profileID, pok := query.Fields["profile_id"]
			if pok {
				db = db.Where("id = ?", profileID)
			}
		}

		if query.Query != "" {
			db = db.Where(" search_column @@ plainto_tsquery(?) ", query.Query)
		}

		err := db.Find(&profileList).Error
		if err != nil {
			return nil, err
		}

		return profileList, nil

	})
}

func (pr *profileRepository) GetTypeByID(ctx context.Context, profileTypeID string) (*models.ProfileType, error) {
	profileType := &models.ProfileType{}
	err := pr.service.DB(ctx, true).First(profileType, "id = ?", profileTypeID).Error
	return profileType, err
}

func (pr *profileRepository) GetTypeByUID(
	ctx context.Context,
	profileType profilev1.ProfileType,
) (*models.ProfileType, error) {
	profileTypeUID := models.ProfileTypeIDMap[profileType]
	profileTypeM := &models.ProfileType{}
	err := pr.service.DB(ctx, true).First(profileTypeM, "uid = ?", profileTypeUID).Error
	return profileTypeM, err
}

func (pr *profileRepository) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	emptyClaims := &frame.AuthenticationClaims{}
	emptyCtx := emptyClaims.ClaimsToContext(ctx)

	profile := &models.Profile{}
	err := pr.service.DB(emptyCtx, true).Preload(clause.Associations).First(profile, "id = ?", id).Error
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
	repo := profileRepository{
		service: service,
	}
	return &repo
}
