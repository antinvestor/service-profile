package repository

import (
	"context"
	"fmt"
	"github.com/antinvestor/service-profile/apps/default/service/models"

	"gorm.io/gorm/clause"

	"github.com/pitabwire/frame"
)

type rosterRepository struct {
	service *frame.Service
}

func (cr *rosterRepository) Search(
	ctx context.Context,
	query *SearchQuery,
) (frame.JobResultPipe[[]*models.Roster], error) {
	service := cr.service
	job := frame.NewJob[[]*models.Roster](
		func(ctx context.Context, jobResult frame.JobResultPipe[[]*models.Roster]) error {
			paginator := query.Pagination
			for paginator.canLoad() {
				var rosterList []*models.Roster

				db := service.DB(ctx, true).
					Joins("LEFT JOIN contacts ON rosters.contact_id = contacts.id").
					Preload("Contact").
					Limit(paginator.limit).Offset(paginator.offset)

				if query.ProfileID != "" {
					db = db.Where("rosters.profile_id = ?", query.ProfileID)
				}

				if query.StartAt != nil && query.EndAt != nil {
					startDate := query.StartAt.Format("2020-01-31T00:00:00Z")
					endDate := query.EndAt.Format("2020-01-31T00:00:00Z")
					db = db.Where("rosters.created_at @@@ '[ ? TO ?]'", startDate, endDate)
				}

				if query.Query != "" {
					whereConditionParams := []any{query.Query}
					whereQueryStr := " contacts.detail  @@@ ? "

					for _, property := range query.PropertiesToSearchOn {
						whereConditionParams = append(whereConditionParams, query.Query)
						searchTerm := fmt.Sprintf(" OR rosters.id  @@@ paradedb.match( 'properties.%s', ?) ", property)
						whereQueryStr += searchTerm
					}

					db = db.Where(whereQueryStr, whereConditionParams...)
				}

				err := db.Find(&rosterList).Error
				if err != nil {
					return jobResult.WriteError(ctx, err)
				}

				err = jobResult.WriteResult(ctx, rosterList)
				if err != nil {
					return err
				}

				if paginator.stop(len(rosterList)) {
					break
				}
			}
			return nil
		},
	)

	err := frame.SubmitJob(ctx, service, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (cr *rosterRepository) GetByID(ctx context.Context, id string) (*models.Roster, error) {
	roster := &models.Roster{}
	err := cr.service.DB(ctx, true).Preload(clause.Associations).First(roster, "id = ?", id).Error
	return roster, err
}

func (cr *rosterRepository) GetByContactAndProfileID(
	ctx context.Context,
	profileID, contactID string,
) (*models.Roster, error) {
	roster := &models.Roster{}
	err := cr.service.DB(ctx, true).
		Preload(clause.Associations).
		Where("profile_id = ? AND contact_id = ?", profileID, contactID).
		First(roster).
		Error
	return roster, err
}

func (cr *rosterRepository) GetByContactID(ctx context.Context, contactID string) ([]*models.Roster, error) {
	rosterList := make([]*models.Roster, 0)
	err := cr.service.DB(ctx, true).
		Preload(clause.Associations).
		Where("contact_id = ?", contactID).
		Find(&rosterList).
		Error
	return rosterList, err
}

func (cr *rosterRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Roster, error) {
	rosterList := make([]*models.Roster, 0)
	err := cr.service.DB(ctx, true).
		Preload(clause.Associations).
		Where("profile_id = ?", profileID).
		Find(&rosterList).
		Error
	return rosterList, err
}

func (cr *rosterRepository) Save(ctx context.Context, roster *models.Roster) (*models.Roster, error) {
	if roster.ID == "" {
		roster.GenID(ctx)
	}

	err := cr.service.DB(ctx, false).Save(roster).Error
	return roster, err
}

func (cr *rosterRepository) Delete(ctx context.Context, id string) error {
	roster, err := cr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return cr.service.DB(ctx, false).Delete(roster).Error
}

func NewRosterRepository(service *frame.Service) RosterRepository {
	rosterRepo := rosterRepository{
		service: service,
	}
	return &rosterRepo
}
