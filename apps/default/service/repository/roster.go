package repository

import (
	"context"
	"strings"
	"time"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/datastore"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type rosterRepository struct {
	service *frame.Service
}

func (cr *rosterRepository) Search(
	ctx context.Context,
	query *datastore.SearchQuery,
) (frame.JobResultPipe[[]*models.Roster], error) {
	return datastore.StableSearch[models.Roster](ctx, cr.service, query, func(
		ctx context.Context,
		query *datastore.SearchQuery,
	) ([]*models.Roster, error) {
		var rosterList []*models.Roster

		paginator := query.Pagination

		db := cr.service.DB(ctx, true).
			Joins("LEFT JOIN contacts ON rosters.contact_id = contacts.id").
			Preload("Contact").
			Limit(paginator.Limit).Offset(paginator.Offset)

		if query.Fields != nil {
			startAt, sok := query.Fields["start_date"]
			stopAt, stok := query.Fields["end_date"]
			if sok && startAt != nil && stok && stopAt != nil {
				startDate, ok1 := startAt.(*time.Time)
				endDate, ok2 := stopAt.(*time.Time)
				if ok1 && ok2 {
					db = db.Where(
						"rosters.created_at BETWEEN ? AND ? ",
						startDate.Format("2020-01-31T00:00:00Z"),
						endDate.Format("2020-01-31T00:00:00Z"),
					)
				}
			}

			profileID, pok := query.Fields["profile_id"]
			if pok {
				db = db.Where(" rosters.profile_id = ?", profileID)
			}
		}

		if query.Query != "" {
			// Use TSVector with prefix matching for partial searches
			// Handle multi-word queries by replacing spaces with & (AND operator)
			searchQuery := strings.ReplaceAll(query.Query, " ", " & ") + ":*"

			// Hybrid approach: Use indexed rosters.search_properties for roster properties
			// and LIKE search for contact details (emails/phones) since TSVector doesn't
			// support partial matching within email tokens
			searchTerm := "%" + query.Query + "%"
			db = db.Where(
				"rosters.search_properties @@ to_tsquery('simple', ?) OR "+
					"contacts.detail ILIKE ?",
				searchQuery, searchTerm,
			)
		}

		err := db.Find(&rosterList).Error
		if err != nil {
			return nil, err
		}

		return rosterList, nil
	})
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
