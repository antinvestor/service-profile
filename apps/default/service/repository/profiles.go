package repository

import (
	"context"
	"fmt"
	"strings"

	profilev1 "github.com/antinvestor/apis/go/profile/v1"
	"github.com/pitabwire/frame"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/dbutil"
)

type profileRepository struct {
	service *frame.Service
}

func (pr *profileRepository) Search(
	ctx context.Context,
	query *dbutil.SearchQuery,
) (frame.JobResultPipe[[]*models.Profile], error) {
	service := pr.service
	job := frame.NewJob(func(ctx context.Context, jobResult frame.JobResultPipe[[]*models.Profile]) error {
		paginator := query.Pagination
		for paginator.CanLoad() {
			profileList, err := pr.searchWithLimits(ctx, query)
			if err != nil {
				return jobResult.WriteError(ctx, err)
			}

			err = jobResult.WriteResult(ctx, profileList)
			if err != nil {
				return err
			}

			if paginator.Stop(len(profileList)) {
				break
			}
		}
		return nil
	})

	err := frame.SubmitJob(ctx, service, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (pr *profileRepository) searchWithLimits(
	ctx context.Context,
	query *dbutil.SearchQuery,
) ([]*models.Profile, error) {
	var profileList []*models.Profile

	paginator := query.Pagination

	db := pr.service.DB(ctx, true).
		Limit(paginator.Limit).Offset(paginator.Offset)

	if query.StartAt != nil && query.EndAt != nil {
		startDate := query.StartAt.Format("2020-01-31T00:00:00Z")
		endDate := query.EndAt.Format("2020-01-31T00:00:00Z")
		db = db.Where("created_at @@@ '[ ? TO ?]'", startDate, endDate)
	}

	if query.Query != "" {
		var whereConditionParams []any
		var whereQueryStrings []string

		for _, property := range query.PropertiesToSearchOn {
			whereConditionParams = append(whereConditionParams, query.Query)
			searchTerm := fmt.Sprintf(
				" id  @@@ paradedb.match( field => 'properties.%s', value => ?, distance => 0) ",
				property,
			)
			whereQueryStrings = append(whereQueryStrings, searchTerm)
		}

		if len(whereQueryStrings) > 0 {
			whereQueryStr := strings.Join(whereQueryStrings, " OR ")
			db = db.Where(whereQueryStr, whereConditionParams...)
		}
	}

	err := db.Find(&profileList).Error
	if err != nil {
		return nil, err
	}

	return profileList, nil
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
