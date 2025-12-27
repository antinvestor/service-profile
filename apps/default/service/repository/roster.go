package repository

import (
	"context"

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
	baseRepo := datastore.NewBaseRepository[*models.Roster](
		ctx, dbPool, workMan, func() *models.Roster { return &models.Roster{} },
	)

	baseRepo.ExtendFieldsAllowed("rosters.profile_id", "rosters.searchable",
		"contacts.look_up_token")

	rosterRepo := rosterRepository{
		BaseRepository: baseRepo,
	}
	return &rosterRepo
}

func (rr *rosterRepository) Search(
	ctx context.Context,
	query *data.SearchQuery,
) (workerpool.JobResultPipe[[]*models.Roster], error) {
	rr.Pool()

	return data.StableSearch[*models.Roster](ctx, rr.WorkManager(), query, func(
		ctx context.Context,
		sq *data.SearchQuery,
	) ([]*models.Roster, error) {
		sq.OrderBy = "rosters.created_at DESC"

		db := rr.Pool().DB(ctx, true).Table("rosters")
		db = db.Select(`rosters.*`)

		db = db.Joins("LEFT JOIN contacts ON rosters.contact_id = contacts.id").
			Preload("Contact")

		// Query optimization: ensure proper index usage through query structure
		// The JOIN and WHERE clauses are structured to utilize existing indexes

		result, err := datastore.SearchFunc[*models.Roster](ctx, db, sq, rr.IsFieldAllowed)
		if err != nil {
			return nil, err
		}

		return result, nil
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

func (rr *rosterRepository) GetByContactIDsAndProfileID(
	ctx context.Context,
	contactIDs []string,
	profileID string,
) ([]*models.Roster, error) {
	rosterList := make([]*models.Roster, 0, len(contactIDs))
	err := rr.Pool().DB(ctx, true).
		Preload(clause.Associations).
		Where("profile_id = ? AND contact_id IN ?", profileID, contactIDs).
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
