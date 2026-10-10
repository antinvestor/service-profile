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

// Package accounts derives the chain account a person's profile owns
// (chain-first rooms design §4, profile accounts design §2.1).
//
//	identity_salt      = HMAC-SHA256(K_identity, "stawi/identity/v1" ‖ profile_id)
//	identity_salt_hash = keccak256(identity_salt)
//	address            = derive.AccountAddress(FamilyEVM, version, identity_salt_hash,
//	                       factory, keccak256(creation_code ‖ identity_salt_hash))
//
// K_identity never leaves Vault Transit in production; the salt itself never
// leaves this package. Only the salt hash and the address are stored or
// published.
package accounts

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
)

// IdentityDomain separates identity salts from every other use of K_identity.
const IdentityDomain = "stawi/identity/v1"

// minStaticKeyLen is the shortest static key accepted: a full SHA-256 block of
// entropy is 32 bytes.
const minStaticKeyLen = 32

// ErrStaticKeyTooShort is returned for a static identity key under 32 bytes.
var ErrStaticKeyTooShort = errors.New("accounts: static identity key must be at least 32 bytes")

// Salter computes the identity salt of a profile. Implementations are
// deterministic: the same profile id always yields the same salt, because the
// account address is derived from it and cannot move.
type Salter interface {
	Salt(ctx context.Context, profileID string) ([32]byte, error)
	// SaltBatch returns one salt per profile id, in the order given.
	SaltBatch(ctx context.Context, profileIDs []string) ([][32]byte, error)
}

// IdentityMessage is the HMAC input for a profile: the domain tag followed by
// the profile id, with no separator (the tag has a fixed length).
func IdentityMessage(profileID string) []byte {
	msg := make([]byte, 0, len(IdentityDomain)+len(profileID))
	msg = append(msg, IdentityDomain...)
	return append(msg, profileID...)
}

// StaticSalter keys the HMAC with a key held in configuration. It exists for
// tests and local stacks; production uses TransitSalter so the key never
// reaches the process.
type StaticSalter struct {
	key []byte
}

// NewStaticSalter returns a salter over the given key.
func NewStaticSalter(key []byte) (*StaticSalter, error) {
	if len(key) < minStaticKeyLen {
		return nil, ErrStaticKeyTooShort
	}
	return &StaticSalter{key: append([]byte(nil), key...)}, nil
}

// Salt implements Salter.
func (s *StaticSalter) Salt(_ context.Context, profileID string) ([32]byte, error) {
	mac := hmac.New(sha256.New, s.key)
	mac.Write(IdentityMessage(profileID))
	var out [32]byte
	copy(out[:], mac.Sum(nil))
	return out, nil
}

// SaltBatch implements Salter.
func (s *StaticSalter) SaltBatch(ctx context.Context, profileIDs []string) ([][32]byte, error) {
	out := make([][32]byte, len(profileIDs))
	for i, id := range profileIDs {
		salt, err := s.Salt(ctx, id)
		if err != nil {
			return nil, err
		}
		out[i] = salt
	}
	return out, nil
}
