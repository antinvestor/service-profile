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
	"encoding/hex"
	"fmt"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/antinvestor/service-profile/apps/default/service/accounts"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

func randomContact() string {
	return util.RandomAlphaNumericString(12) + "@accounts.testing.com"
}

// TestCreatePersonDerivesPrimaryAccount: a PERSON profile owns its account
// from creation, the account is the one stawi's derive computes, and the
// profile.account_created fact is staged with the profile.
func (pts *ProfileTestSuite) TestCreatePersonDerivesPrimaryAccount() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)
		deriver := tests.NewTestDeriver(t)

		created, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_PERSON, Contact: randomContact(),
		})
		require.NoError(t, err)
		require.Len(t, created.GetAccounts(), 1)

		want, err := deriver.Derive(ctx, created.GetId())
		require.NoError(t, err)
		acct := created.GetAccounts()[0]
		assert.Equal(t, want.AddressHex(), acct.GetAddress())
		assert.Equal(t, accounts.FamilyEVM, acct.GetFamily())
		assert.EqualValues(t, 1, acct.GetVersion())
		assert.Equal(t, want.IdentitySaltHash[:], acct.GetIdentitySaltHash())
		assert.True(t, acct.GetPrimary())

		read, err := pb.GetByID(ctx, created.GetId())
		require.NoError(t, err)
		require.Len(t, read.GetAccounts(), 1, "every profile read carries the accounts")
		assert.Equal(t, acct.GetAddress(), read.GetAccounts()[0].GetAddress())

		staged, err := tests.StagedFacts(ctx, svc, business.FactAccountCreated)
		require.NoError(t, err)
		require.Len(t, staged, 1)
		payload := staged[0].Payload
		assert.Equal(t, created.GetId(), payload["profile_id"])
		assert.Equal(t, want.AddressHex(), payload["address"])
		assert.Equal(t, "EVM", payload["family"])
		assert.Equal(t, "1", fmt.Sprint(payload["version"]))
		assert.Equal(t, "0x"+hex.EncodeToString(want.IdentitySaltHash[:]), payload["identity_salt_hash"])
		assert.NotContains(t, payload, "identity_salt")

		// Creating again with the same contact returns the same profile and
		// derives nothing new.
		again, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_PERSON, Contact: created.GetContacts()[0].GetDetail(),
		})
		require.NoError(t, err)
		assert.Equal(t, created.GetId(), again.GetId())
		staged, err = tests.StagedFacts(ctx, svc, business.FactAccountCreated)
		require.NoError(t, err)
		assert.Len(t, staged, 1)
	})
}

func (pts *ProfileTestSuite) TestBotsAndInstitutionsGetNoAccount() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		for _, pt := range []profilev1.ProfileType{profilev1.ProfileType_BOT, profilev1.ProfileType_INSTITUTION} {
			created, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{Type: pt, Contact: randomContact()})
			require.NoError(t, err, pt.String())
			assert.Empty(t, created.GetAccounts(), pt.String())
		}
		staged, err := tests.StagedFacts(ctx, svc, business.FactAccountCreated)
		require.NoError(t, err)
		assert.Empty(t, staged)
	})
}

// TestCreateWithoutSalterDerivesNothing is the development posture: no salter
// configured, profiles are still created, without accounts.
func (pts *ProfileTestSuite) TestCreateWithoutSalterDerivesNothing() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusinessWith(ctx, svc, nil)

		created, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_PERSON, Contact: randomContact(),
		})
		require.NoError(t, err)
		assert.Empty(t, created.GetAccounts())
	})
}

// TestAccountBackfill derives the missing accounts of existing PERSON
// profiles in batches, emits one fact per account, and is idempotent.
func (pts *ProfileTestSuite) TestAccountBackfill() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		withoutAccounts, _ := pts.getProfileBusinessWith(ctx, svc, nil)
		deriver := tests.NewTestDeriver(t)
		backfill := tests.NewAccountBusiness(ctx, svc, deriver, 2)

		// Persons seeded by the migrations are backfilled first, so the
		// counts below are this test's own profiles.
		_, err := backfill.Backfill(ctx)
		require.NoError(t, err)
		seeded, err := tests.StagedFacts(ctx, svc, business.FactAccountCreated)
		require.NoError(t, err)

		var persons []string
		for range 5 {
			p, createErr := withoutAccounts.CreateProfile(ctx, &profilev1.CreateRequest{
				Type: profilev1.ProfileType_PERSON, Contact: randomContact(),
			})
			require.NoError(t, createErr)
			persons = append(persons, p.GetId())
		}
		bot, err := withoutAccounts.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_BOT, Contact: randomContact(),
		})
		require.NoError(t, err)

		n, err := backfill.Backfill(ctx)
		require.NoError(t, err)
		assert.Equal(t, len(persons), n)

		for _, id := range persons {
			list, listErr := backfill.ListByProfile(ctx, id)
			require.NoError(t, listErr)
			require.Len(t, list, 1)
			want, deriveErr := deriver.Derive(ctx, id)
			require.NoError(t, deriveErr)
			assert.Equal(t, want.Address[:], list[0].Address)
			assert.True(t, list[0].Primary)
		}
		botAccounts, err := backfill.ListByProfile(ctx, bot.GetId())
		require.NoError(t, err)
		assert.Empty(t, botAccounts)

		staged, err := tests.StagedFacts(ctx, svc, business.FactAccountCreated)
		require.NoError(t, err)
		assert.Len(t, staged, len(seeded)+len(persons))

		again, err := backfill.Backfill(ctx)
		require.NoError(t, err)
		assert.Zero(t, again, "a second backfill finds nothing to do")
		staged, err = tests.StagedFacts(ctx, svc, business.FactAccountCreated)
		require.NoError(t, err)
		assert.Len(t, staged, len(seeded)+len(persons), "and stages no further facts")

		disabled := tests.NewAccountBusiness(ctx, svc, nil, 2)
		skipped, err := disabled.Backfill(ctx)
		require.NoError(t, err)
		assert.Zero(t, skipped)
	})
}

// TestMergeMovesAccounts: the survivor keeps its primary account and gains
// the merged profile's account as a secondary one.
func (pts *ProfileTestSuite) TestMergeMovesAccounts() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		survivor, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_PERSON, Contact: randomContact(),
		})
		require.NoError(t, err)
		merged, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_PERSON, Contact: randomContact(),
		})
		require.NoError(t, err)
		survivorAddr := survivor.GetAccounts()[0].GetAddress()
		mergedAddr := merged.GetAccounts()[0].GetAddress()

		result, err := pb.MergeProfile(ctx, &profilev1.MergeRequest{Id: survivor.GetId(), Mergeid: merged.GetId()})
		require.NoError(t, err)
		require.Len(t, result.GetAccounts(), 2)
		assert.Equal(t, survivorAddr, result.GetAccounts()[0].GetAddress())
		assert.True(t, result.GetAccounts()[0].GetPrimary())
		assert.Equal(t, mergedAddr, result.GetAccounts()[1].GetAddress())
		assert.False(t, result.GetAccounts()[1].GetPrimary())

		staged, err := tests.StagedFacts(ctx, svc, business.FactAccountsMerged)
		require.NoError(t, err)
		require.Len(t, staged, 1)
		assert.Equal(t, survivor.GetId(), staged[0].Payload["surviving_profile_id"])
		assert.Equal(t, merged.GetId(), staged[0].Payload["merged_profile_id"])
		assert.Equal(t, []any{mergedAddr}, staged[0].Payload["addresses"])

		resolved, err := tests.NewAccountBusiness(ctx, svc, nil, 0).Resolve(ctx, []string{mergedAddr})
		require.NoError(t, err)
		require.Len(t, resolved, 1)
		assert.Equal(t, survivor.GetId(), resolved[0].ProfileID)

		_, err = pb.MergeProfile(ctx, &profilev1.MergeRequest{Id: survivor.GetId(), Mergeid: survivor.GetId()})
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func (pts *ProfileTestSuite) TestJurisdictionProperty() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		props, err := structpb.NewStruct(map[string]any{"jurisdiction": " ke "})
		require.NoError(t, err)
		created, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type: profilev1.ProfileType_PERSON, Contact: randomContact(), Properties: props,
		})
		require.NoError(t, err)
		assert.Equal(t, "KE", created.GetProperties().AsMap()["jurisdiction"])

		for _, bad := range []any{"XX", "Kenya", "KEN", 254.0, ""} {
			badProps, structErr := structpb.NewStruct(map[string]any{"jurisdiction": bad})
			require.NoError(t, structErr)
			_, err = pb.UpdateProfile(ctx, &profilev1.UpdateRequest{Id: created.GetId(), Properties: badProps})
			require.Error(t, err, bad)
			assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err), bad)
		}

		ugProps, err := structpb.NewStruct(map[string]any{"jurisdiction": "ug"})
		require.NoError(t, err)
		updated, err := pb.UpdateProfile(ctx, &profilev1.UpdateRequest{Id: created.GetId(), Properties: ugProps})
		require.NoError(t, err)
		assert.Equal(t, "UG", updated.GetProperties().AsMap()["jurisdiction"])
	})
}
