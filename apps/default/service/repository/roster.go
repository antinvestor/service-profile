package repository

import (
	"context"
	"strings"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type rosterRepository struct {
	datastore.BaseRepository[*models.Roster]
}

func NewRosterRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) RosterRepository {
	rosterRepo := rosterRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Roster](
			ctx, dbPool, workMan, func() *models.Roster { return &models.Roster{} },
		),
	}
	return &rosterRepo
}

func (rr *rosterRepository) Search(
	ctx context.Context,
	query *data.SearchQuery,
) (workerpool.JobResultPipe[[]*models.Roster], error) {
	return data.StableSearch[*models.Roster](ctx, rr.WorkManager(), query, func(
		ctx context.Context,
		query *data.SearchQuery,
	) ([]*models.Roster, error) {
		var rosterList []*models.Roster

		db := rr.Pool().DB(ctx, true).
			Joins("LEFT JOIN contacts ON rosters.contact_id = contacts.id").
			Preload("Contact")

		rr.DefaultSearchFunction(ctx, db, query)

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

func (rr *rosterRepository) GetByContactAndProfileID(
	ctx context.Context,
	profileID, contactID string,
) (*models.Roster, error) {
	roster := &models.Roster{}
	err := rr.Pool().DB(ctx, true).
		Preload(clause.Associations).
		Where("profile_id = ? AND contact_id = ?", profileID, contactID).
		First(roster).
		Error
	return roster, err
}

func (rr *rosterRepository) GetByContactID(ctx context.Context, contactID string) ([]*models.Roster, error) {
	rosterList := make([]*models.Roster, 0)
	err := rr.Pool().DB(ctx, true).
		Preload(clause.Associations).
		Where("contact_id = ?", contactID).
		Find(&rosterList).
		Error
	return rosterList, err
}

func (rr *rosterRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Roster, error) {
	rosterList := make([]*models.Roster, 0)
	err := rr.Pool().DB(ctx, true).
		Preload(clause.Associations).
		Where("profile_id = ?", profileID).
		Find(&rosterList).
		Error
	return rosterList, err
}
