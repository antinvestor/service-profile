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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/service/accounts/derive"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/accounts"
)

// accountDerivationVector is stawi's packages/contracts/test/vectors/
// account_derivation.json: the address the Solidity AccountFactory computed
// for a fixed identity salt hash from the creation code it deploys.
type accountDerivationVector struct {
	Factory          string `json:"factory"`
	AccountVersion   uint32 `json:"account_version"`
	IdentitySaltHash string `json:"identity_salt_hash"`
	CreationCode     string `json:"creation_code"`
	InitCodeHash     string `json:"init_code_hash"`
	Account          string `json:"account"`
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.TrimPrefix(strings.ToLower(s), "0x"))
	require.NoError(t, err)
	return b
}

func loadVector(t *testing.T) accountDerivationVector {
	t.Helper()
	raw, err := os.ReadFile("testdata/account_derivation.json")
	require.NoError(t, err)
	var v accountDerivationVector
	require.NoError(t, json.Unmarshal(raw, &v))
	return v
}

func vectorParams(t *testing.T, v accountDerivationVector) accounts.Params {
	t.Helper()
	factory, err := accounts.ParseAddress(v.Factory)
	require.NoError(t, err)
	return accounts.Params{Factory: factory, CreationCode: mustHex(t, v.CreationCode), Version: v.AccountVersion}
}

func TestFromSaltHashMatchesContractsVector(t *testing.T) {
	v := loadVector(t)
	var saltHash [32]byte
	copy(saltHash[:], mustHex(t, v.IdentitySaltHash))

	addr, err := accounts.FromSaltHash(vectorParams(t, v), saltHash)
	require.NoError(t, err)
	assert.Equal(t, strings.ToLower(v.Account), accounts.HexAddress(addr[:]),
		"profile must derive the address the factory computes")

	ich := derive.InitCodeHash(mustHex(t, v.CreationCode), saltHash)
	assert.Equal(t, strings.ToLower(v.InitCodeHash), "0x"+hex.EncodeToString(ich[:]))
}

// TestDeriverMatchesStawiDerive checks the whole chain from a profile id: the
// HMAC salt, its keccak hash and the address all equal what stawi's derive
// package computes for the same inputs.
func TestDeriverMatchesStawiDerive(t *testing.T) {
	v := loadVector(t)
	params := vectorParams(t, v)
	key := []byte("0123456789abcdef0123456789abcdef")
	salter, err := accounts.NewHMACSalter(key)
	require.NoError(t, err)
	deriver, err := accounts.NewDeriver(salter, params)
	require.NoError(t, err)

	const profileID = "d1v0profile0000000000"
	got, err := deriver.Derive(context.Background(), profileID)
	require.NoError(t, err)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("stawi/identity/v1" + profileID))
	var salt [32]byte
	copy(salt[:], mac.Sum(nil))
	saltHash, err := derive.IdentitySaltHash(salt)
	require.NoError(t, err)
	want, err := derive.AccountAddress(derive.FamilyEVM, params.Version, saltHash, params.Factory,
		derive.InitCodeHash(params.CreationCode, saltHash))
	require.NoError(t, err)

	assert.Equal(t, want, got.Address)
	assert.Equal(t, saltHash, got.IdentitySaltHash)
	assert.Equal(t, accounts.FamilyEVM, got.Family)
	assert.Equal(t, params.Version, got.Version)
	assert.Equal(t, profileID, got.ProfileID)

	again, err := deriver.Derive(context.Background(), profileID)
	require.NoError(t, err)
	assert.Equal(t, got, again, "derivation is deterministic")

	other, err := deriver.Derive(context.Background(), profileID+"x")
	require.NoError(t, err)
	assert.NotEqual(t, got.Address, other.Address)
}

// TestHMACSalterKnownAnswer pins the salt, salt hash and address for a fixed
// key and profile id. The salt equals
//
//	printf 'stawi/identity/v1d1v0profile0000000000' |
//	  openssl dgst -sha256 -mac HMAC -macopt key:0123456789abcdef0123456789abcdef
//
// and the salt hash equals `cast keccak 0x<salt>`.
func TestHMACSalterKnownAnswer(t *testing.T) {
	v := loadVector(t)
	path := filepath.Join(t.TempDir(), "salt")
	require.NoError(t, os.WriteFile(path, []byte("0123456789abcdef0123456789abcdef"), 0o600))
	salter, err := accounts.NewHMACSalterFromFile(path)
	require.NoError(t, err)

	salt, err := salter.Salt(context.Background(), "d1v0profile0000000000")
	require.NoError(t, err)
	assert.Equal(t, "29af8699c187b4cca3e9dd23d6d4792993b7e97b24b7418a0a761ac640902f35", hex.EncodeToString(salt[:]))

	deriver, err := accounts.NewDeriver(salter, vectorParams(t, v))
	require.NoError(t, err)
	got, err := deriver.Derive(context.Background(), "d1v0profile0000000000")
	require.NoError(t, err)
	assert.Equal(t, "9f43eef116c8ba7ed765c74a3f827c3b88c4c29c54037bf462e4cb3439f5dc30",
		hex.EncodeToString(got.IdentitySaltHash[:]))
	assert.Equal(t, "0x15c79d9fd0a15ed61d130c5310a974606e07c8db", accounts.HexAddress(got.Address[:]))

	batch, err := salter.SaltBatch(context.Background(), []string{"d1v0profile0000000000", "other"})
	require.NoError(t, err)
	require.Len(t, batch, 2)
	assert.Equal(t, salt, batch[0])
	assert.NotEqual(t, salt, batch[1])
}

func TestHMACSalterKeyFile(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, content []byte) string {
		path := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(path, content, 0o600))
		return path
	}

	_, err := accounts.NewHMACSalter(make([]byte, 16))
	require.ErrorIs(t, err, accounts.ErrKeyTooShort)

	_, err = accounts.NewHMACSalterFromFile(filepath.Join(dir, "missing"))
	require.ErrorIs(t, err, accounts.ErrKeyFile)

	_, err = accounts.NewHMACSalterFromFile(write("empty", nil))
	require.ErrorIs(t, err, accounts.ErrKeyTooShort)

	_, err = accounts.NewHMACSalterFromFile(write("short", make([]byte, 31)))
	require.ErrorIs(t, err, accounts.ErrKeyTooShort)

	_, err = accounts.NewHMACSalterFromFile(write("long", make([]byte, 2048)))
	require.ErrorIs(t, err, accounts.ErrKeyTooLong)

	// The file's bytes are the key as stored: a trailing newline is key
	// material, not whitespace to trim.
	key := []byte("0123456789abcdef0123456789abcdef")
	plain, err := accounts.NewHMACSalterFromFile(write("plain", key))
	require.NoError(t, err)
	newline, err := accounts.NewHMACSalterFromFile(write("newline", append(append([]byte(nil), key...), '\n')))
	require.NoError(t, err)
	a, err := plain.Salt(context.Background(), "p")
	require.NoError(t, err)
	b, err := newline.Salt(context.Background(), "p")
	require.NoError(t, err)
	assert.NotEqual(t, a, b)

	// Errors and formatting never carry the key.
	secret := []byte("this-key-must-never-be-printed-anywhere!")
	s, err := accounts.NewHMACSalterFromFile(write("secret", secret))
	require.NoError(t, err)
	assert.NotContains(t, fmt.Sprintf("%v %+v %#v %s", s, s, s, s), string(secret))
	_, err = accounts.NewHMACSalterFromFile(write("short-secret", []byte("short-secret-value")))
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "short-secret-value")
}

func TestParseAddress(t *testing.T) {
	addr, err := accounts.ParseAddress("0x56989d378452be71D7Ff078D5222a6B150D11f73")
	require.NoError(t, err)
	assert.Equal(t, "0x56989d378452be71d7ff078d5222a6b150d11f73", accounts.HexAddress(addr[:]))
	for _, bad := range []string{"", "0x12", "zz989d378452be71D7Ff078D5222a6B150D11f73", "0x56989d378452be71D7Ff078D5222a6B150D11f7300"} {
		_, err = accounts.ParseAddress(bad)
		assert.Error(t, err, bad)
	}
}

func TestFromConfig(t *testing.T) {
	v := loadVector(t)
	saltFile := filepath.Join(t.TempDir(), "salt")
	require.NoError(t, os.WriteFile(saltFile, []byte(strings.Repeat("k", 32)), 0o600))

	t.Run("nothing configured derives nothing", func(t *testing.T) {
		d, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{AccountVersion: 1})
		require.NoError(t, err)
		assert.Nil(t, d)
	})

	t.Run("salt file with manifest parameters", func(t *testing.T) {
		d, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountsSaltFile:        saltFile,
			AccountFactory:          v.Factory,
			AccountCreationCode:     v.CreationCode,
			AccountCreationCodeHash: "0xecb9c721045bde978e8458bcccc7d4534c375b33dce1d7e44a0c4fc8ad2dc680",
			AccountVersion:          1,
		})
		require.NoError(t, err)
		require.NotNil(t, d)
		assert.EqualValues(t, 1, d.Version())
	})

	t.Run("creation code read from a file", func(t *testing.T) {
		path := t.TempDir() + "/code.hex"
		require.NoError(t, os.WriteFile(path, []byte(v.CreationCode+"\n"), 0o600))
		d, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountsSaltFile: saltFile, AccountFactory: v.Factory,
			AccountCreationCodeFile: path, AccountVersion: 1,
		})
		require.NoError(t, err)
		require.NotNil(t, d)
	})

	t.Run("creation code hash mismatch", func(t *testing.T) {
		_, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountsSaltFile: saltFile, AccountFactory: v.Factory,
			AccountCreationCode: v.CreationCode, AccountCreationCodeHash: "0x" + strings.Repeat("00", 32),
			AccountVersion: 1,
		})
		require.ErrorIs(t, err, accounts.ErrAccountConfig)
	})

	t.Run("salt file set but missing refuses to start", func(t *testing.T) {
		_, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountsSaltFile: filepath.Join(t.TempDir(), "absent"), AccountFactory: v.Factory,
			AccountCreationCode: v.CreationCode, AccountVersion: 1,
		})
		require.ErrorIs(t, err, accounts.ErrKeyFile)
	})

	t.Run("salt file too short refuses to start", func(t *testing.T) {
		short := filepath.Join(t.TempDir(), "short")
		require.NoError(t, os.WriteFile(short, []byte("abc"), 0o600))
		_, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountsSaltFile: short, AccountFactory: v.Factory,
			AccountCreationCode: v.CreationCode, AccountVersion: 1,
		})
		require.ErrorIs(t, err, accounts.ErrKeyTooShort)
	})

	t.Run("salt file without account parameters", func(t *testing.T) {
		_, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountsSaltFile: saltFile, AccountVersion: 1,
		})
		require.ErrorIs(t, err, accounts.ErrAccountConfig)
	})

	t.Run("account parameters without a salt file", func(t *testing.T) {
		_, err := accounts.FromConfig(context.Background(), &config.ProfileConfig{
			AccountFactory: v.Factory, AccountVersion: 1,
		})
		require.ErrorIs(t, err, accounts.ErrAccountConfig)
	})
}
