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

package business

import (
	"context"

	"errors"
	"fmt"
	"strings"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"

	"connectrpc.com/connect"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/service/accounts"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// MaxResolveAddresses bounds one ResolveAccounts request.
const MaxResolveAddresses = 500

const defaultBackfillBatchSize = 200

// AccountBusiness owns the chain accounts of profiles.
type AccountBusiness interface {
	// Enabled reports whether accounts are derived at all (a salter and the
	// manifest parameters are configured).
	Enabled() bool
	// NewPrimary derives the primary account for a profile id that is about
	// to be created. It calls the salter, so it runs before the profile's
	// transaction opens. It returns nil when accounts are disabled.
	NewPrimary(ctx context.Context, profileID string) (*models.ProfileAccount, error)
	// ListByProfile returns a profile's accounts, primary first.
	ListByProfile(ctx context.Context, profileID string) ([]*models.ProfileAccount, error)
	// Resolve maps addresses to their accounts; unknown addresses are omitted.
	Resolve(ctx context.Context, addresses []string) ([]*models.ProfileAccount, error)
	// MoveOnMerge re-homes the merged profile's accounts onto the survivor
	// as secondary accounts and stages profile.accounts_merged.
	MoveOnMerge(ctx context.Context, survivingProfileID, mergedProfileID string) ([]*models.ProfileAccount, error)
	// Backfill derives the primary account of every PERSON profile that has
	// none for the configured version, staging profile.account_created for
	// each. It is idempotent and returns how many accounts it created.
	Backfill(ctx context.Context) (int, error)
}

// NewAccountBusiness returns the account business. A nil deriver disables
// derivation; reads, resolution and merges keep working on stored rows.
func NewAccountBusiness(
	deriver *accounts.Deriver,
	accountRepo repository.ProfileAccountRepository,
	profileRepo repository.ProfileRepository,
	backfillBatchSize int,
) AccountBusiness {
	if backfillBatchSize <= 0 {
		backfillBatchSize = defaultBackfillBatchSize
	}
	return &accountBusiness{
		deriver:     deriver,
		accountRepo: accountRepo,
		profileRepo: profileRepo,
		batchSize:   backfillBatchSize,
	}
}

type accountBusiness struct {
	deriver     *accounts.Deriver
	accountRepo repository.ProfileAccountRepository
	profileRepo repository.ProfileRepository
	batchSize   int
}

func (ab *accountBusiness) Enabled() bool { return ab.deriver != nil }

func toModel(a accounts.Account, primary bool) *models.ProfileAccount {
	return &models.ProfileAccount{
		ProfileID:        a.ProfileID,
		Family:           a.Family,
		AccountVersion:   a.Version,
		Address:          append([]byte(nil), a.Address[:]...),
		IdentitySaltHash: append([]byte(nil), a.IdentitySaltHash[:]...),
		Factory:          append([]byte(nil), a.Factory[:]...),
		Primary:          primary,
	}
}

func (ab *accountBusiness) NewPrimary(ctx context.Context, profileID string) (*models.ProfileAccount, error) {
	if ab.deriver == nil {
		return nil, nil //nolint:nilnil // accounts disabled
	}
	derived, err := ab.deriver.Derive(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("derive profile account: %w", err)
	}
	return toModel(derived, true), nil
}

func (ab *accountBusiness) ListByProfile(ctx context.Context, profileID string) ([]*models.ProfileAccount, error) {
	return ab.accountRepo.ListByProfileID(ctx, profileID)
}

func (ab *accountBusiness) Resolve(ctx context.Context, addresses []string) ([]*models.ProfileAccount, error) {
	addresses = normalizeAddressList(addresses)
	if len(addresses) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("at least one address is required"))
	}
	if len(addresses) > MaxResolveAddresses {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("at most %d addresses may be resolved at once", MaxResolveAddresses))
	}
	seen := make(map[[20]byte]struct{}, len(addresses))
	raw := make([][]byte, 0, len(addresses))
	for _, s := range addresses {
		addr, err := accounts.ParseAddress(s)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		if _, dup := seen[addr]; dup {
			continue
		}
		seen[addr] = struct{}{}
		raw = append(raw, addr[:])
	}
	return ab.accountRepo.ListByAddresses(ctx, raw)
}

func (ab *accountBusiness) MoveOnMerge(
	ctx context.Context,
	survivingProfileID, mergedProfileID string,
) ([]*models.ProfileAccount, error) {
	return ab.accountRepo.MoveToProfile(ctx, mergedProfileID, survivingProfileID,
		func(moved []*models.ProfileAccount) *outbox.Event {
			return AccountsMergedFact(ctx, survivingProfileID, mergedProfileID, moved)
		})
}

func (ab *accountBusiness) Backfill(ctx context.Context) (int, error) {
	log := util.Log(ctx)
	if ab.deriver == nil {
		log.Warn("account backfill skipped: no identity salter configured")
		return 0, nil
	}
	personType, err := ab.profileRepo.GetTypeByUID(ctx, profilev1.ProfileType_PERSON)
	if err != nil {
		return 0, fmt.Errorf("account backfill: person profile type: %w", err)
	}

	created := 0
	afterID := ""
	for {
		page, pageErr := ab.accountRepo.PersonProfilesWithoutAccount(
			ctx, personType.GetID(), accounts.FamilyEVM, ab.deriver.Version(), afterID, ab.batchSize)
		if pageErr != nil {
			return created, fmt.Errorf("account backfill: list profiles: %w", pageErr)
		}
		if len(page) == 0 {
			break
		}
		ids := make([]string, len(page))
		for i, p := range page {
			ids[i] = p.GetID()
		}
		derived, deriveErr := ab.deriver.DeriveBatch(ctx, ids)
		if deriveErr != nil {
			return created, fmt.Errorf("account backfill: derive: %w", deriveErr)
		}
		for i, p := range page {
			account := toModel(derived[i], true)
			account.TenantID = p.TenantID
			account.PartitionID = p.PartitionID
			account.AccessID = p.AccessID
			ok, createErr := ab.accountRepo.CreateWithFact(ctx, account,
				func(a *models.ProfileAccount) *outbox.Event { return AccountCreatedFact(ctx, a) })
			if createErr != nil {
				return created, fmt.Errorf("account backfill: profile %s: %w", p.GetID(), createErr)
			}
			if ok {
				created++
			}
		}
		afterID = page[len(page)-1].GetID()
		log.Info("account backfill progress", "created", created, "after", afterID)
	}
	log.Info("account backfill complete", "created", created, "version", ab.deriver.Version())
	return created, nil
}

// AccountsToAPI converts stored accounts to their API form.
func AccountsToAPI(list []*models.ProfileAccount) []*profilev1.ProfileAccount {
	out := make([]*profilev1.ProfileAccount, 0, len(list))
	for _, a := range list {
		out = append(out, &profilev1.ProfileAccount{
			Address:          accounts.HexAddress(a.Address),
			Family:           a.Family,
			Version:          a.AccountVersion,
			IdentitySaltHash: append([]byte(nil), a.IdentitySaltHash...),
			Primary:          a.Primary,
		})
	}
	return out
}

// normalizeAddressList trims and drops empty entries.
func normalizeAddressList(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
