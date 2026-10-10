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

// Package accounts is the client other services use to work with profile
// accounts: who owns an address, which accounts a profile has, the identity
// salt hash an address was derived from, and a profile's jurisdiction.
//
// The profile service derives and stores every person's account, and only it
// maps an address to a profile (ResolveAccounts, service principals only).
// Consumers keep no copy of that link; they ask through this package.
//
// Addresses are plain strings: 0x-prefixed hex, accepted in any case and
// returned lowercase.
package accounts

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect"
	common "github.com/antinvestor/common/v2"
	"github.com/antinvestor/common/v2/connection"
	"github.com/antinvestor/common/v2/servicecatalog"
)

// ErrNotOwned is returned when no profile owns the address, or the named
// profile does not own it.
var ErrNotOwned = errors.New("accounts: the address is not owned by that profile")

// ResolveLimit is the most addresses one ResolveAccounts request may carry.
const ResolveLimit = 500

const saltHashLen = 32

// API is the part of profile.v1 ProfileService this client uses. The
// generated client satisfies it, and so does a handler (for tests).
type API interface {
	ResolveAccounts(
		ctx context.Context,
		req *connect.Request[profilev1.ResolveAccountsRequest],
	) (*connect.Response[profilev1.ResolveAccountsResponse], error)
	GetById(
		ctx context.Context,
		req *connect.Request[profilev1.GetByIdRequest],
	) (*connect.Response[profilev1.GetByIdResponse], error)
}

// Account is one of a profile's accounts.
type Account struct {
	Address string // lowercase 0x hex
	Family  string
	Version uint32
	Primary bool
}

// Client answers account questions through the profile service.
type Client struct{ api API }

// New wraps a profile API.
func New(api API) *Client { return &Client{api: api} }

// Target is where the profile service is reached.
type Target struct {
	Endpoint              string
	WorkloadAPITargetPath string
}

// Dial builds the profile client the house way: workload identity, OAuth2
// and retries come from cfg. Resolving addresses needs account_resolve
// (ROLE_SERVICE). An empty endpoint returns (nil, nil): running without a
// profile service is a development posture.
func Dial(ctx context.Context, cfg any, t Target) (*Client, error) {
	if t.Endpoint == "" {
		return nil, nil //nolint:nilnil // no endpoint is a development posture, not an error
	}
	cli, err := connection.NewServiceClient(ctx, cfg, common.ServiceTarget{
		Endpoint:              t.Endpoint,
		WorkloadAPITargetPath: t.WorkloadAPITargetPath,
		ServiceID:             servicecatalog.ServiceProfile,
	}, profilev1connect.NewProfileServiceClient)
	if err != nil {
		return nil, err
	}
	return New(cli), nil
}

// Resolve maps addresses (any case) to their owners' profile ids, in batches
// of ResolveLimit. Unknown addresses are absent; keys are lowercase 0x hex.
func (c *Client) Resolve(ctx context.Context, addresses []string) (map[string]string, error) {
	out := make(map[string]string, len(addresses))
	for start := 0; start < len(addresses); start += ResolveLimit {
		end := min(start+ResolveLimit, len(addresses))
		batch := make([]string, 0, end-start)
		for _, a := range addresses[start:end] {
			batch = append(batch, strings.ToLower(a))
		}
		res, err := c.api.ResolveAccounts(ctx,
			connect.NewRequest(&profilev1.ResolveAccountsRequest{Addresses: batch}))
		if err != nil {
			return nil, fmt.Errorf("accounts: resolve: %w", err)
		}
		for _, o := range res.Msg.GetData() {
			out[strings.ToLower(o.GetAddress())] = o.GetProfileId()
		}
	}
	return out, nil
}

// ProfileOf is the id of the profile owning address, or ErrNotOwned.
func (c *Client) ProfileOf(ctx context.Context, address string) (string, error) {
	key := strings.ToLower(address)
	m, err := c.Resolve(ctx, []string{key})
	if err != nil {
		return "", err
	}
	p, ok := m[key]
	if !ok {
		return "", ErrNotOwned
	}
	return p, nil
}

func (c *Client) profile(ctx context.Context, profileID string) (*profilev1.ProfileObject, error) {
	res, err := c.api.GetById(ctx, connect.NewRequest(&profilev1.GetByIdRequest{Id: profileID}))
	if err != nil {
		return nil, fmt.Errorf("accounts: profile %s: %w", profileID, err)
	}
	return res.Msg.GetData(), nil
}

// Accounts lists the profile's accounts, primary first.
func (c *Client) Accounts(ctx context.Context, profileID string) ([]Account, error) {
	p, err := c.profile(ctx, profileID)
	if err != nil {
		return nil, err
	}
	out := make([]Account, 0, len(p.GetAccounts()))
	for _, a := range p.GetAccounts() {
		out = append(out, Account{
			Address: strings.ToLower(a.GetAddress()),
			Family:  a.GetFamily(),
			Version: a.GetVersion(),
			Primary: a.GetPrimary(),
		})
	}
	return out, nil
}

// IdentitySaltHash is keccak256 of the identity salt the profile's address
// was derived from; only the hash ever leaves the profile service.
func (c *Client) IdentitySaltHash(ctx context.Context, profileID, address string) ([32]byte, error) {
	var out [32]byte
	p, err := c.profile(ctx, profileID)
	if err != nil {
		return out, err
	}
	want := strings.ToLower(address)
	for _, a := range p.GetAccounts() {
		if strings.ToLower(a.GetAddress()) != want {
			continue
		}
		h := a.GetIdentitySaltHash()
		if len(h) != saltHashLen || bytes.Equal(h, out[:]) {
			return out, fmt.Errorf("accounts: account %s has no identity salt hash", want)
		}
		copy(out[:], h)
		return out, nil
	}
	return out, ErrNotOwned
}

// Jurisdiction is the profile's ISO 3166-1 alpha-2 country, or "".
func (c *Client) Jurisdiction(ctx context.Context, profileID string) (string, error) {
	p, err := c.profile(ctx, profileID)
	if err != nil {
		return "", err
	}
	j, _ := p.GetProperties().AsMap()["jurisdiction"].(string)
	return strings.ToUpper(j), nil
}
