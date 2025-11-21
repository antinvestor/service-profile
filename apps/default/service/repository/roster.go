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

type RosterWithSim struct {
	*models.Roster
	Sim float64
}

func NewRosterRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) RosterRepository {
	baseRepo := datastore.NewBaseRepository[*models.Roster](
		ctx, dbPool, workMan, func() *models.Roster { return &models.Roster{} },
	)

	fieldMap := baseRepo.FieldsAllowed()
	fieldMap["rosters.profile_id"] = struct{}{}
	fieldMap["rosters.searchable"] = struct{}{}
	fieldMap["contacts.detail"] = struct{}{}
	fieldMap["sim"] = struct{}{}
	fieldMap["similarity(contacts.detail,?)"] = struct{}{}

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
		sq.OrderBy = "sim DESC"

		db := rr.Pool().DB(ctx, true).Table("rosters")

		searchQuery, ok := sq.FiltersOrByValue["SIMILARITY(contacts.detail,?) > 0"]
		if !ok {
			db = db.Select(`rosters.*, SIMILARITY(contacts.detail,?) AS sim`, searchQuery)
		} else {
			db = db.Select(`rosters.*, 1 AS sim`)
		}

		db = db.Joins("LEFT JOIN contacts ON rosters.contact_id = contacts.id").
			Preload("Contact")

		result, err := datastore.SearchFunc[*RosterWithSim](ctx, db, sq, rr.IsFieldAllowed)
		if err != nil {
			return nil, err
		}

		rosters := make([]*models.Roster, len(result))
		for i, r := range result {
			rosters[i] = r.Roster
		}
		return rosters, nil
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
