package repository

import (
	"context"
	"sort"
	"time"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/security"
	"github.com/pitabwire/frame/v2/workerpool"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
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
	// Profile types are global seed data with NULL partition/tenant columns.
	// Skip tenancy checks so the query isn't scoped to a specific partition.
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	err := pr.Pool().DB(unscopedCtx, true).First(profileType, "id = ?", profileTypeID).Error
	return profileType, err
}

func (pr *profileRepository) GetTypeByUID(
	ctx context.Context,
	profileType profilev1.ProfileType,
) (*models.ProfileType, error) {
	profileTypeUID := models.ProfileTypeIDMap[profileType]
	profileTypeM := &models.ProfileType{}
	// Profile types are global seed data with NULL partition/tenant columns.
	// Skip tenancy checks so the query isn't scoped to a specific partition.
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	err := pr.Pool().DB(unscopedCtx, true).First(profileTypeM, "uid = ?", profileTypeUID).Error
	return profileTypeM, err
}

func (pr *profileRepository) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	// Profiles are accessed cross-tenant (e.g. reading back a just-created profile,
	// or looking up profiles by contact across tenants). Skip tenancy scoping.
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	profile := &models.Profile{}
	err := pr.Pool().DB(unscopedCtx, true).Preload(clause.Associations).First(profile, "id = ?", id).Error
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

// CreateWithFact persists a profile and stages the profile.created fact in the
// same transaction, so the fact exists exactly when the profile does
// (GFOS K5).
func (pr *profileRepository) CreateWithFact(
	ctx context.Context,
	profile *models.Profile,
	fact func(*models.Profile) *outbox.Event,
) error {
	return pr.CreateWithAccount(ctx, profile, nil, func(p *models.Profile) []*outbox.Event {
		return []*outbox.Event{fact(p)}
	})
}

// CreateWithAccount persists a profile, its primary account (when given) and
// the facts describing them in one transaction: a person's account exists
// exactly when the profile does. The account takes the profile's id and
// tenancy.
func (pr *profileRepository) CreateWithAccount(
	ctx context.Context,
	profile *models.Profile,
	account *models.ProfileAccount,
	facts func(*models.Profile) []*outbox.Event,
) error {
	profile.GenID(ctx)

	return pr.Pool().DB(ctx, false).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		if account != nil {
			account.ProfileID = profile.GetID()
			account.TenantID = profile.TenantID
			account.PartitionID = profile.PartitionID
			account.AccessID = profile.AccessID
			account.GenID(ctx)
			if err := tx.Create(account).Error; err != nil {
				return err
			}
		}
		for _, evt := range facts(profile) {
			if err := outbox.Enqueue(tx, evt); err != nil {
				return err
			}
		}
		return nil
	})
}

// Merge folds merging into target in one transaction (see the interface).
func (pr *profileRepository) Merge(
	ctx context.Context,
	target, merging *models.Profile,
	inTx func(tx *gorm.DB) ([]*outbox.Event, error),
) error {
	return pr.Pool().DB(ctx, false).Transaction(func(tx *gorm.DB) error {
		// Lock both profiles, in id order so concurrent merges of the same
		// pair cannot deadlock, and require both to be live: a survivor
		// deleted (or merged away) since it was read must not receive
		// accounts.
		ids := []string{target.GetID(), merging.GetID()}
		sort.Strings(ids)
		var live []string
		if err := tx.Table("profiles").Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ? AND deleted_at IS NULL", ids).
			Order("id asc").
			Pluck("id", &live).Error; err != nil {
			return err
		}
		if len(live) != len(ids) {
			return gorm.ErrRecordNotFound
		}

		facts, err := inTx(tx)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		res := tx.Table("profiles").Where("id = ? AND deleted_at IS NULL", target.GetID()).
			Updates(map[string]any{"properties": target.Properties, "modified_at": now})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		res = tx.Table("profiles").Where("id = ? AND deleted_at IS NULL", merging.GetID()).
			Update("deleted_at", now)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		for _, evt := range facts {
			if err = outbox.Enqueue(tx, evt); err != nil {
				return err
			}
		}
		return nil
	})
}
