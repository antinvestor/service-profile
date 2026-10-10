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

package accounts

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/stawilabs/stawi/pkg/protocol/derive"
)

// FamilyEVM is the account family name stored and published; its on-chain
// identifier is derive.FamilyEVM (keccak256("EVM")).
const FamilyEVM = "EVM"

// ErrAccountConfig reports an unusable account derivation configuration.
var ErrAccountConfig = errors.New("accounts: invalid account configuration")

// Params are the per-deployment inputs of the derivation, from the protocol
// manifest: the canonical factory, the StawiAccount creation code it deploys,
// and the account version it registers.
type Params struct {
	Factory      [20]byte
	CreationCode []byte
	Version      uint32
}

// Account is a derived account. The identity salt is deliberately absent.
type Account struct {
	ProfileID        string
	Family           string
	Version          uint32
	Address          [20]byte
	IdentitySaltHash [32]byte
	Factory          [20]byte
}

// AddressHex is the address as lowercase 0x-prefixed hex, the form stored in
// facts and returned by the API.
func (a Account) AddressHex() string { return HexAddress(a.Address[:]) }

// HexAddress formats a 20-byte address as lowercase 0x-prefixed hex.
func HexAddress(addr []byte) string { return "0x" + hex.EncodeToString(addr) }

// ParseAddress parses a 0x-prefixed (or bare) 40-hex-digit address in any case.
func ParseAddress(s string) ([20]byte, error) {
	var out [20]byte
	raw, err := hex.DecodeString(strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(s), "0x"), "0X"))
	if err != nil || len(raw) != len(out) {
		return out, fmt.Errorf("accounts: %q is not a 20-byte hex address", s)
	}
	copy(out[:], raw)
	return out, nil
}

// FromSaltHash derives the account address for an identity salt hash with the
// exact code the contracts are vector-tested against.
func FromSaltHash(p Params, identitySaltHash [32]byte) ([20]byte, error) {
	initCodeHash := derive.InitCodeHash(p.CreationCode, identitySaltHash)
	return derive.AccountAddress(derive.FamilyEVM, p.Version, identitySaltHash, p.Factory, initCodeHash)
}

// Deriver turns profile ids into accounts.
type Deriver struct {
	salter Salter
	params Params
}

// NewDeriver validates the parameters and returns a deriver.
func NewDeriver(salter Salter, p Params) (*Deriver, error) {
	if salter == nil {
		return nil, fmt.Errorf("%w: no salter", ErrAccountConfig)
	}
	if p.Factory == ([20]byte{}) {
		return nil, fmt.Errorf("%w: factory address is required", ErrAccountConfig)
	}
	if len(p.CreationCode) == 0 {
		return nil, fmt.Errorf("%w: account creation code is required", ErrAccountConfig)
	}
	if p.Version == 0 {
		return nil, fmt.Errorf("%w: account version must be >= 1", ErrAccountConfig)
	}
	return &Deriver{salter: salter, params: p}, nil
}

// Version is the account version this deriver produces.
func (d *Deriver) Version() uint32 { return d.params.Version }

// Derive computes the primary account of one profile.
func (d *Deriver) Derive(ctx context.Context, profileID string) (Account, error) {
	accounts, err := d.DeriveBatch(ctx, []string{profileID})
	if err != nil {
		return Account{}, err
	}
	return accounts[0], nil
}

// DeriveBatch computes accounts for several profiles with one salter call.
func (d *Deriver) DeriveBatch(ctx context.Context, profileIDs []string) ([]Account, error) {
	for _, id := range profileIDs {
		if strings.TrimSpace(id) == "" {
			return nil, errors.New("accounts: empty profile id")
		}
	}
	salts, err := d.salter.SaltBatch(ctx, profileIDs)
	if err != nil {
		return nil, err
	}
	if len(salts) != len(profileIDs) {
		return nil, fmt.Errorf("accounts: salter returned %d salts for %d profiles", len(salts), len(profileIDs))
	}
	out := make([]Account, len(profileIDs))
	for i, profileID := range profileIDs {
		salt := salts[i]
		saltHash, hashErr := derive.IdentitySaltHash(salt)
		if hashErr != nil {
			return nil, hashErr
		}
		addr, addrErr := FromSaltHash(d.params, saltHash)
		if addrErr != nil {
			return nil, addrErr
		}
		out[i] = Account{
			ProfileID:        profileID,
			Family:           FamilyEVM,
			Version:          d.params.Version,
			Address:          addr,
			IdentitySaltHash: saltHash,
			Factory:          d.params.Factory,
		}
	}
	return out, nil
}
