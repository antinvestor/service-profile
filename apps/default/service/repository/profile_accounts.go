// Copyright 2023-2026 Ant Investor Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import (
	"context"
	"time"

	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/security"
	"github.com/pitabwire/frame/v2/workerpool"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

type profileAccountRepository struct {
	datastore.BaseRepository[*models.ProfileAccount]
}

// NewProfileAccountRepository returns the repository of profile accounts.
func NewProfileAccountRepository(
	ctx context.Context,
	dbPool pool.Pool,
	workMan workerpool.Manager,
) ProfileAccountRepository {
	return &profileAccountRepository{
		BaseRepository: datastore.NewBaseRepository[*models.ProfileAccount](
			ctx, dbPool, workMan, func() *models.ProfileAccount { return &models.ProfileAccount{} },
		),
	}
}

// ListByProfileID returns a profile's accounts, primary first. Accounts are
// identity data read across tenants, like the profile itself.
func (ar *profileAccountRepository) ListByProfileID(
	ctx context.Context,
	profileID string,
) ([]*models.ProfileAccount, error) {
	var out []*models.ProfileAccount
	err := ar.Pool().DB(security.SkipTenancyChecksOnClaims(ctx), true).
		Where("profile_id = ?", profileID).
		Order("is_primary desc, account_version desc, created_at asc").
		Find(&out).Error
	return out, err
}

// ListByAddresses returns the accounts with any of the given addresses.
func (ar *profileAccountRepository) ListByAddresses(
	ctx context.Context,
	addresses [][]byte,
) ([]*models.ProfileAccount, error) {
	if len(addresses) == 0 {
		return nil, nil
	}
	var out []*models.ProfileAccount
	err := ar.Pool().DB(security.SkipTenancyChecksOnClaims(ctx), true).
		Where("address IN ?", addresses).
		Find(&out).Error
	return out, err
}

// CreateWithFact inserts an account and stages its fact in one transaction.
// An account whose address (or primary slot) already exists is left alone and
// no fact is staged, so repeating the call is harmless.
func (ar *profileAccountRepository) CreateWithFact(
	ctx context.Context,
	account *models.ProfileAccount,
	fact func(*models.ProfileAccount) *outbox.Event,
) (bool, error) {
	account.GenID(ctx)
	created := false
	err := ar.Pool().DB(ctx, false).Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(account)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		created = true
		return outbox.Enqueue(tx, fact(account))
	})
	return created, err
}

// MoveToProfile re-homes every account of one profile onto another as
// secondary accounts and stages the fact describing the move, in one
// transaction. Nothing is staged when there is nothing to move.
func (ar *profileAccountRepository) MoveToProfile(
	ctx context.Context,
	fromProfileID, toProfileID string,
	fact func([]*models.ProfileAccount) *outbox.Event,
) ([]*models.ProfileAccount, error) {
	var moved []*models.ProfileAccount
	err := ar.Pool().DB(security.SkipTenancyChecksOnClaims(ctx), false).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("profile_id = ?", fromProfileID).
			Order("created_at asc").
			Find(&moved).Error; err != nil {
			return err
		}
		if len(moved) == 0 {
			return nil
		}
		ids := make([]string, len(moved))
		for i, a := range moved {
			ids[i] = a.GetID()
			a.ProfileID = toProfileID
			a.Primary = false
		}
		// Table, not Model: a model value would run the BaseModel hooks,
		// which give it an id that then narrows the WHERE clause.
		if err := tx.Table("profile_accounts").
			Where("id IN ?", ids).
			Updates(map[string]any{
				"profile_id":  toProfileID,
				"is_primary":  false,
				"modified_at": time.Now().UTC(),
			}).Error; err != nil {
			return err
		}
		return outbox.Enqueue(tx, fact(moved))
	})
	return moved, err
}

// PersonProfilesWithoutAccount pages, by id, through profiles of the given
// type that have no live primary account for the family and version.
func (ar *profileAccountRepository) PersonProfilesWithoutAccount(
	ctx context.Context,
	profileTypeID, family string,
	accountVersion uint32,
	afterID string,
	limit int,
) ([]*models.Profile, error) {
	var out []*models.Profile
	err := ar.Pool().DB(security.SkipTenancyChecksOnClaims(ctx), true).
		Model(&models.Profile{}).
		Where("profile_type_id = ? AND id > ?", profileTypeID, afterID).
		Where(`NOT EXISTS (SELECT 1 FROM profile_accounts a
			WHERE a.profile_id = profiles.id AND a.family = ? AND a.account_version = ?
			AND a.is_primary AND a.deleted_at IS NULL)`, family, accountVersion).
		Order("id asc").
		Limit(limit).
		Find(&out).Error
	return out, err
}
