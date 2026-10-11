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
// K_identity is read once from a mounted file (ACCOUNTS_SALT_FILE) and kept
// only inside the salter. The salt itself never leaves this package. Only the
// salt hash and the address are stored or published.
package accounts

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

// IdentityDomain separates identity salts from every other use of K_identity.
const IdentityDomain = "stawi/identity/v1"

const (
	// MinKeyLen is the shortest identity key accepted: 32 bytes, the
	// SHA-256 output size.
	MinKeyLen = 32
	// maxKeyLen bounds what is read from the key file, so a mount pointing
	// at the wrong file fails instead of loading it.
	maxKeyLen = 1024
)

// ErrKeyTooShort is returned for an identity key under MinKeyLen bytes.
var ErrKeyTooShort = fmt.Errorf("accounts: identity key must be at least %d bytes", MinKeyLen)

// ErrKeyTooLong is returned for an identity key file over maxKeyLen bytes.
var ErrKeyTooLong = fmt.Errorf("accounts: identity key file is over %d bytes", maxKeyLen)

// ErrKeyFile reports an identity key file that cannot be used.
var ErrKeyFile = errors.New("accounts: identity key file")

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

// HMACSalter computes identity salts in process with HMAC-SHA256 under
// K_identity.
type HMACSalter struct {
	key []byte
}

// NewHMACSalter returns a salter over a copy of key.
func NewHMACSalter(key []byte) (*HMACSalter, error) {
	if len(key) < MinKeyLen {
		return nil, ErrKeyTooShort
	}
	return &HMACSalter{key: append([]byte(nil), key...)}, nil
}

// NewHMACSalterFromFile reads K_identity from path. The file's bytes are the
// key exactly as stored: nothing is trimmed or decoded. The read buffer is
// zeroed once the salter holds its copy. Errors name the path, never the
// contents.
func NewHMACSalterFromFile(path string) (*HMACSalter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrKeyFile, err)
	}
	defer func() { _ = f.Close() }()

	buf := make([]byte, maxKeyLen+1)
	defer clear(buf)
	n, err := io.ReadFull(f, buf)
	switch {
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
	case err != nil:
		return nil, fmt.Errorf("%w: %s: %w", ErrKeyFile, path, err)
	default:
		return nil, fmt.Errorf("%w: %s", ErrKeyTooLong, path)
	}
	salter, err := NewHMACSalter(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("%w: %s holds %d bytes", err, path, n)
	}
	return salter, nil
}

// String keeps the key out of formatted output.
func (s *HMACSalter) String() string { return "accounts.HMACSalter{key: redacted}" }

// GoString keeps the key out of %#v output.
func (s *HMACSalter) GoString() string { return s.String() }

// Salt implements Salter.
func (s *HMACSalter) Salt(_ context.Context, profileID string) ([32]byte, error) {
	mac := hmac.New(sha256.New, s.key)
	mac.Write(IdentityMessage(profileID))
	var out [32]byte
	copy(out[:], mac.Sum(nil))
	return out, nil
}

// SaltBatch implements Salter.
func (s *HMACSalter) SaltBatch(ctx context.Context, profileIDs []string) ([][32]byte, error) {
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
