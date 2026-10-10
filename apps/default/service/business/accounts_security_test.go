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

package business_test

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2"
	"github.com/pitabwire/frame/v2/data"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

var errInjected = errors.New("injected failure after the account update")

// failingAfterMove moves the accounts for real, then fails, so the merge
// transaction must undo the move.
type failingAfterMove struct {
	business.AccountBusiness
}

func (f failingAfterMove) MoveOnMerge(tx *gorm.DB, survivor, merged *models.Profile) ([]*models.ProfileAccount, error) {
	moved, err := f.AccountBusiness.MoveOnMerge(tx, survivor, merged)
	if err != nil {
		return nil, err
	}
	if len(moved) == 0 {
		return nil, errors.New("test setup: nothing was moved")
	}
	return nil, errInjected
}

type tenancyRef struct{ tenant, partition string }

func newTenancy() tenancyRef { return tenancyRef{tenant: util.IDString(), partition: util.IDString()} }

func (pts *ProfileTestSuite) as(ctx context.Context, tn tenancyRef, subject string) context.Context {
	return pts.WithAuthClaims(ctx, tn.tenant, tn.partition, subject)
}

func (pts *ProfileTestSuite) createPersonIn(
	ctx context.Context, t *testing.T, pb business.ProfileBusiness, tn tenancyRef,
) *profilev1.ProfileObject {
	t.Helper()
	p, err := pb.CreateProfile(pts.as(ctx, tn, "creator-"+util.RandomAlphaNumericString(6)),
		&profilev1.CreateRequest{Type: profilev1.ProfileType_PERSON, Contact: randomContact()})
	require.NoError(t, err)
	return p
}

// requireUntouched asserts the profile still owns exactly its own primary
// account.
func requireUntouched(ctx context.Context, t *testing.T, svc *frame.Service, profileID string) {
	t.Helper()
	rows, err := tests.AccountsOf(ctx, svc, profileID)
	require.NoError(t, err)
	require.Len(t, rows, 1, "the account still belongs to %s", profileID)
	assert.True(t, rows[0].Primary)
}

// TestMergeCannotPullAnotherTenantsAccounts: a caller authorised to merge in
// its own tenancy cannot name a profile of another tenant (or partition) and
// take its account.
func (pts *ProfileTestSuite) TestMergeCannotPullAnotherTenantsAccounts() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		mine, theirs := newTenancy(), newTenancy()
		survivor := pts.createPersonIn(ctx, t, pb, mine)
		victim := pts.createPersonIn(ctx, t, pb, theirs)
		otherPartition := pts.createPersonIn(ctx, t, pb, tenancyRef{tenant: mine.tenant, partition: util.IDString()})

		caller := pts.as(ctx, mine, "admin-"+util.RandomAlphaNumericString(6))
		for _, mergeID := range []string{victim.GetId(), otherPartition.GetId()} {
			_, err := pb.MergeProfile(caller, &profilev1.MergeRequest{Id: survivor.GetId(), Mergeid: mergeID})
			require.Error(t, err)
			assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
			requireUntouched(ctx, t, svc, mergeID)
		}
		// Nor the other way round: the survivor named from another tenancy.
		_, err := pb.MergeProfile(pts.as(ctx, theirs, "admin-x"),
			&profilev1.MergeRequest{Id: victim.GetId(), Mergeid: survivor.GetId()})
		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
		requireUntouched(ctx, t, svc, survivor.GetId())

		// An internal caller without claims still cannot mix tenancies.
		_, err = pb.MergeProfile(ctx, &profilev1.MergeRequest{Id: survivor.GetId(), Mergeid: victim.GetId()})
		require.Error(t, err)
		assert.Equal(t, connect.CodeFailedPrecondition, connect.CodeOf(err))
		requireUntouched(ctx, t, svc, victim.GetId())

		staged, err := tests.StagedFacts(ctx, svc, business.FactAccountsMerged)
		require.NoError(t, err)
		assert.Empty(t, staged)

		// Inside the caller's own tenancy the merge goes through.
		sibling := pts.createPersonIn(ctx, t, pb, mine)
		merged, err := pb.MergeProfile(caller, &profilev1.MergeRequest{Id: survivor.GetId(), Mergeid: sibling.GetId()})
		require.NoError(t, err)
		assert.Equal(t, survivor.GetId(), merged.GetId())
		rows, err := tests.AccountsOf(ctx, svc, survivor.GetId())
		require.NoError(t, err)
		assert.Len(t, rows, 2)
	})
}

// TestMergeIsAtomic: a failure after the account update rolls back the move,
// the demotion, the profile deletion and the fact together.
func (pts *ProfileTestSuite) TestMergeIsAtomic() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)
		tn := newTenancy()
		survivor := pts.createPersonIn(ctx, t, pb, tn)
		merged := pts.createPersonIn(ctx, t, pb, tn)

		failing, _ := pts.getProfileBusinessWithAccounts(ctx, svc, failingAfterMove{
			AccountBusiness: tests.NewAccountBusiness(ctx, svc, tests.NewTestDeriver(t), 0, tests.AsService),
		})
		_, err := failing.MergeProfile(pts.as(ctx, tn, "admin"),
			&profilev1.MergeRequest{Id: survivor.GetId(), Mergeid: merged.GetId()})
		require.Error(t, err)

		requireUntouched(ctx, t, svc, survivor.GetId())
		requireUntouched(ctx, t, svc, merged.GetId())
		stillThere, err := pb.GetByID(ctx, merged.GetId())
		require.NoError(t, err, "the merged profile is not deleted")
		assert.Equal(t, merged.GetId(), stillThere.GetId())
		staged, err := tests.StagedFacts(ctx, svc, business.FactAccountsMerged)
		require.NoError(t, err)
		assert.Empty(t, staged, "no fact is staged for a merge that did not happen")
	})
}

// TestAccountsVisibleOnlyToOwnerAndServices covers every profile read: the
// owner and service principals see accounts; another user with profile_view
// does not; a caller in another tenancy sees none even as a service.
func (pts *ProfileTestSuite) TestAccountsVisibleOnlyToOwnerAndServices() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		tn := newTenancy()

		userView, _ := pts.getProfileBusinessWithAccounts(ctx, svc,
			tests.NewAccountBusiness(ctx, svc, tests.NewTestDeriver(t), 0, nil))
		serviceView, _ := pts.getProfileBusiness(ctx, svc)

		// Created by another principal (a user with profile_create): the
		// response carries no accounts.
		created, err := userView.CreateProfile(pts.as(ctx, tn, "someone-else"),
			&profilev1.CreateRequest{Type: profilev1.ProfileType_PERSON, Contact: randomContact()})
		require.NoError(t, err)
		assert.Empty(t, created.GetAccounts(), "create response for a non-owner")
		id := created.GetId()
		contact := created.GetContacts()[0].GetDetail()

		reads := map[string]func(ctx context.Context, pb business.ProfileBusiness) (*profilev1.ProfileObject, error){
			"GetByID": func(c context.Context, pb business.ProfileBusiness) (*profilev1.ProfileObject, error) {
				return pb.GetByID(c, id)
			},
			"GetByIDAndPartition": func(c context.Context, pb business.ProfileBusiness) (*profilev1.ProfileObject, error) {
				return pb.GetByIDAndPartition(c, id, tn.partition)
			},
			"GetByContact": func(c context.Context, pb business.ProfileBusiness) (*profilev1.ProfileObject, error) {
				return pb.GetByContact(c, contact)
			},
			"Search/ToAPI": func(c context.Context, pb business.ProfileBusiness) (*profilev1.ProfileObject, error) {
				return pb.ToAPI(c, &models.Profile{BaseModel: created2Model(id)})
			},
		}

		for name, read := range reads {
			other, readErr := read(pts.as(ctx, tn, "other-user"), userView)
			require.NoError(t, readErr, name)
			assert.Empty(t, other.GetAccounts(), "%s: a non-owner with profile_view sees no accounts", name)

			own, readErr := read(pts.as(ctx, tn, id), userView)
			require.NoError(t, readErr, name)
			assert.Len(t, own.GetAccounts(), 1, "%s: the owner sees their account", name)

			svcRead, readErr := read(pts.as(ctx, tn, "settlement"), serviceView)
			require.NoError(t, readErr, name)
			assert.Len(t, svcRead.GetAccounts(), 1, "%s: a service in the tenancy sees it", name)

			foreign, readErr := read(pts.as(ctx, newTenancy(), "settlement"), serviceView)
			require.NoError(t, readErr, name)
			assert.Empty(t, foreign.GetAccounts(), "%s: another tenancy sees nothing", name)
		}

		resolved, err := tests.NewAccountBusiness(ctx, svc, nil, 0, tests.AsService).
			Resolve(pts.as(ctx, newTenancy(), "settlement"), []string{"0x" + hexAddr(ctx, t, svc, id)})
		require.NoError(t, err)
		assert.Empty(t, resolved, "ResolveAccounts from another tenancy finds nothing")
	})
}

func created2Model(id string) data.BaseModel { return data.BaseModel{ID: id} }

func hexAddr(ctx context.Context, t *testing.T, svc *frame.Service, profileID string) string {
	t.Helper()
	rows, err := tests.AccountsOf(ctx, svc, profileID)
	require.NoError(t, err)
	require.NotEmpty(t, rows)
	return hex.EncodeToString(rows[0].Address)
}
