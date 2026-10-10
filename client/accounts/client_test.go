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

package accounts_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/antinvestor/service-profile/client/accounts"
)

type fakeProfiles struct {
	profilev1connect.UnimplementedProfileServiceHandler

	owners   map[string]string
	profiles map[string]*profilev1.ProfileObject
	batches  []int
}

func (f *fakeProfiles) ResolveAccounts(
	_ context.Context,
	req *connect.Request[profilev1.ResolveAccountsRequest],
) (*connect.Response[profilev1.ResolveAccountsResponse], error) {
	f.batches = append(f.batches, len(req.Msg.GetAddresses()))
	out := &profilev1.ResolveAccountsResponse{}
	for _, a := range req.Msg.GetAddresses() {
		if p, ok := f.owners[a]; ok {
			out.Data = append(out.Data, &profilev1.AccountOwner{Address: a, ProfileId: p})
		}
	}
	return connect.NewResponse(out), nil
}

//nolint:revive,staticcheck // generated RPC name
func (f *fakeProfiles) GetById(
	_ context.Context,
	req *connect.Request[profilev1.GetByIdRequest],
) (*connect.Response[profilev1.GetByIdResponse], error) {
	p, ok := f.profiles[req.Msg.GetId()]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no profile"))
	}
	return connect.NewResponse(&profilev1.GetByIdResponse{Data: p}), nil
}

const (
	addrA = "0x00000000000000000000000000000000000000a1"
	addrB = "0x00000000000000000000000000000000000000b2"
)

func TestProfileOfResolvesOnlyOwnedAddresses(t *testing.T) {
	c := accounts.New(&fakeProfiles{owners: map[string]string{addrA: "p1"}})
	_, err := c.ProfileOf(context.Background(), strings.ToUpper(addrA[2:]))
	require.ErrorIs(t, err, accounts.ErrNotOwned, "an address without 0x is not the stored form")
	p, err := c.ProfileOf(context.Background(), "0x"+strings.ToUpper(addrA[2:]))
	require.NoError(t, err)
	require.Equal(t, "p1", p)
	_, err = c.ProfileOf(context.Background(), addrB)
	require.ErrorIs(t, err, accounts.ErrNotOwned)
}

func TestResolveBatchesAtTheServiceLimit(t *testing.T) {
	f := &fakeProfiles{owners: map[string]string{}}
	addrs := make([]string, 1201)
	for i := range addrs {
		addrs[i] = "0x" + strings.Repeat("0", 36) + strings.Repeat("1", 4)
	}
	_, err := accounts.New(f).Resolve(context.Background(), addrs)
	require.NoError(t, err)
	require.Equal(t, []int{500, 500, 201}, f.batches)
}

func TestIdentitySaltHashAndJurisdictionComeFromTheProfile(t *testing.T) {
	salt := make([]byte, 32)
	salt[0] = 0x5a
	props, err := structpb.NewStruct(map[string]any{"jurisdiction": "ke"})
	require.NoError(t, err)
	f := &fakeProfiles{profiles: map[string]*profilev1.ProfileObject{"p1": {
		Id: "p1", Properties: props,
		Accounts: []*profilev1.ProfileAccount{
			{Address: addrA, Family: "EVM", Version: 1, IdentitySaltHash: salt, Primary: true},
		},
	}}}
	c := accounts.New(f)
	h, err := c.IdentitySaltHash(context.Background(), "p1", addrA)
	require.NoError(t, err)
	require.Equal(t, byte(0x5a), h[0])
	_, err = c.IdentitySaltHash(context.Background(), "p1", addrB)
	require.ErrorIs(t, err, accounts.ErrNotOwned, "an address the profile does not own has no salt here")

	accts, err := c.Accounts(context.Background(), "p1")
	require.NoError(t, err)
	require.Equal(t, []accounts.Account{{Address: addrA, Family: "EVM", Version: 1, Primary: true}}, accts)
	j, err := c.Jurisdiction(context.Background(), "p1")
	require.NoError(t, err)
	require.Equal(t, "KE", j)

	_, err = c.Accounts(context.Background(), "missing")
	require.Error(t, err)
}

func TestDialWithoutEndpointIsNil(t *testing.T) {
	c, err := accounts.Dial(context.Background(), nil, accounts.Target{})
	require.NoError(t, err)
	require.Nil(t, c)
}
