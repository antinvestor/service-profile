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

package handlers_test

import (
	"encoding/base64"
	"strings"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/service/handlers"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

type AccountsHandlerSuite struct {
	tests.ProfileBaseTestSuite
}

func TestAccountsHandlerSuite(t *testing.T) {
	suite.Run(t, new(AccountsHandlerSuite))
}

func testDEK(t *testing.T, cfg *config.ProfileConfig) *config.DEK {
	t.Helper()
	key, err := base64.StdEncoding.DecodeString(cfg.DEKActiveAES256GCMKey)
	require.NoError(t, err)
	lookup, err := base64.StdEncoding.DecodeString(cfg.DEKLookupTokenHMACSHA256Key)
	require.NoError(t, err)
	return &config.DEK{KeyID: cfg.DEKActiveKeyID, Key: key, LookUpKey: lookup}
}

// TestResolveAccountsIsServiceOnly: a service principal resolves addresses to
// profiles (unknown addresses omitted); the profile's owner, holding every
// owner permission, is refused.
func (s *AccountsHandlerSuite) TestResolveAccountsIsServiceOnly() {
	s.WithTestDependancies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := s.CreateService(t, dep)
		cfg, ok := svc.Config().(*config.ProfileConfig)
		require.True(t, ok)

		server := handlers.NewProfileServer(ctx, svc, testDEK(t, cfg), s.GetNotificationCli(t),
			s.FunctionChecker, tests.NewTestDeriver(t))

		tenantID, partitionID := util.IDString(), util.IDString()
		serviceID := "svc-" + util.RandomAlphaNumericString(8)
		serviceCtx := s.WithAuthClaims(ctx, tenantID, partitionID, serviceID)
		s.SeedTenantRole(serviceCtx, svc, tenantID, partitionID, serviceID, authz.RoleService)

		// A service (e.g. authentication) creates the person in its tenancy.
		created, err := server.Create(serviceCtx, connect.NewRequest(&profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: util.RandomAlphaNumericString(10) + "@resolve.testing.com",
		}))
		require.NoError(t, err)
		profile := created.Msg.GetData()
		require.Len(t, profile.GetAccounts(), 1)
		address := profile.GetAccounts()[0].GetAddress()
		unknown := "0x" + strings.Repeat("ab", 20)

		resp, err := server.ResolveAccounts(serviceCtx, connect.NewRequest(&profilev1.ResolveAccountsRequest{
			Addresses: []string{strings.ToUpper(address[:2]) + strings.ToUpper(address[2:]), unknown},
		}))
		require.NoError(t, err)
		require.Len(t, resp.Msg.GetData(), 1, "unknown addresses are omitted")
		assert.Equal(t, address, resp.Msg.GetData()[0].GetAddress())
		assert.Equal(t, profile.GetId(), resp.Msg.GetData()[0].GetProfileId())

		// A service of another tenancy resolves nothing.
		otherTenant, otherPartition := util.IDString(), util.IDString()
		foreignCtx := s.WithAuthClaims(ctx, otherTenant, otherPartition, serviceID)
		s.SeedTenantRole(foreignCtx, svc, otherTenant, otherPartition, serviceID, authz.RoleService)
		foreign, err := server.ResolveAccounts(foreignCtx, connect.NewRequest(&profilev1.ResolveAccountsRequest{
			Addresses: []string{address},
		}))
		require.NoError(t, err)
		assert.Empty(t, foreign.Msg.GetData())

		ownerCtx := s.WithAuthClaims(ctx, tenantID, partitionID, profile.GetId())
		s.SeedTenantRole(ownerCtx, svc, tenantID, partitionID, profile.GetId(), authz.RoleOwner)
		s.SeedTenantRole(ownerCtx, svc, tenantID, partitionID, profile.GetId(), authz.RoleAdmin)
		_, err = server.ResolveAccounts(ownerCtx, connect.NewRequest(&profilev1.ResolveAccountsRequest{
			Addresses: []string{address},
		}))
		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))

		// The owner still sees the account on their own profile.
		own, err := server.GetById(ownerCtx, connect.NewRequest(&profilev1.GetByIdRequest{Id: profile.GetId()}))
		require.NoError(t, err)
		require.Len(t, own.Msg.GetData().GetAccounts(), 1)
		assert.Equal(t, address, own.Msg.GetData().GetAccounts()[0].GetAddress())

		// Another user of the tenancy holding profile_view (a viewer) reads
		// the profile without its accounts.
		viewerID := "viewer-" + util.RandomAlphaNumericString(6)
		viewerCtx := s.WithAuthClaims(ctx, tenantID, partitionID, viewerID)
		s.SeedTenantRole(viewerCtx, svc, tenantID, partitionID, viewerID, authz.RoleViewer)
		viewed, err := server.GetById(viewerCtx, connect.NewRequest(&profilev1.GetByIdRequest{Id: profile.GetId()}))
		require.NoError(t, err)
		assert.Empty(t, viewed.Msg.GetData().GetAccounts())

		_, err = server.ResolveAccounts(serviceCtx, connect.NewRequest(&profilev1.ResolveAccountsRequest{
			Addresses: []string{"not-an-address"},
		}))
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}
