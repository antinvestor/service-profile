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
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"testing"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/v2/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
)

// secp256k1Generator is the curve's base point in SEC1 uncompressed form — a
// known good public key to validate against.
const secp256k1Generator = "0479be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798" +
	"483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}

func newP256Key(t *testing.T) []byte {
	t.Helper()
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	require.NoError(t, err)
	// Bytes() is the SEC1 uncompressed form; the devices service stores the
	// raw coordinates the GFOS verifier reads.
	return priv.PublicKey().Bytes()[1:]
}

func webAuthnExtra() data.JSONMap {
	credentialID := base64.RawURLEncoding.EncodeToString([]byte("credential-id-of-enough-length"))
	return data.JSONMap{
		"credential_id": credentialID,
		"rp_id":         "stawi.app",
		"cose_alg":      float64(-7),
	}
}

func TestValidateSecp256k1KeyMaterial(t *testing.T) {
	valid := mustHex(t, secp256k1Generator)

	t.Run("accepts an uncompressed point and derives the key id", func(t *testing.T) {
		out, err := business.ValidateKeyMaterial(
			business.KeyTypeSecp256k1PublicKey, valid, data.JSONMap{"label": "primary"},
		)
		require.NoError(t, err)

		keyID, ok := out["key_id"].(string)
		require.True(t, ok, "key_id must be derived")
		assert.Len(t, keyID, 42, "key id is 0x plus 40 hex characters")
		assert.Equal(t, "secp256k1", out["curve"])
		assert.Equal(t, "sec1-uncompressed", out["encoding"])
		assert.Equal(t, "primary", out["label"])
		assert.NotEmpty(t, out["key_sha256"])
	})

	t.Run("rejects a declared key id that does not match the key", func(t *testing.T) {
		_, err := business.ValidateKeyMaterial(
			business.KeyTypeSecp256k1PublicKey, valid,
			data.JSONMap{"key_id": "0x0000000000000000000000000000000000000000"},
		)
		require.Error(t, err)
	})

	testCases := []struct {
		name string
		key  []byte
	}{
		{name: "too short", key: valid[:64]},
		{name: "compressed form", key: append([]byte{0x02}, valid[1:33]...)},
		{name: "wrong prefix", key: append([]byte{0x05}, valid[1:]...)},
		{name: "point at infinity", key: make([]byte, 65)},
		{name: "not on the curve", key: append(append([]byte{0x04}, valid[1:33]...), make([]byte, 32)...)},
	}
	for _, tc := range testCases {
		t.Run("rejects "+tc.name, func(t *testing.T) {
			_, err := business.ValidateKeyMaterial(business.KeyTypeSecp256k1PublicKey, tc.key, data.JSONMap{})
			require.Error(t, err)
		})
	}

	t.Run("rejects an unknown extra field", func(t *testing.T) {
		_, err := business.ValidateKeyMaterial(
			business.KeyTypeSecp256k1PublicKey, valid, data.JSONMap{"private_key": "never"},
		)
		require.Error(t, err)
	})
}

func TestValidateP256WebAuthnKeyMaterial(t *testing.T) {
	valid := newP256Key(t)

	t.Run("accepts raw coordinates with the required webauthn fields", func(t *testing.T) {
		extra := webAuthnExtra()
		extra["transports"] = []any{"internal", "hybrid"}
		extra["sign_count"] = float64(0)
		extra["user_verification"] = "required"

		out, err := business.ValidateKeyMaterial(business.KeyTypeP256WebauthnPublicKey, valid, extra)
		require.NoError(t, err)

		keyID, ok := out["key_id"].(string)
		require.True(t, ok)
		assert.Len(t, keyID, 42)
		assert.Equal(t, "p-256", out["curve"])
		assert.Equal(t, "raw-xy", out["encoding"])
		assert.Equal(t, "stawi.app", out["rp_id"])
	})

	t.Run("rejects the uncompressed 65 byte form", func(t *testing.T) {
		withPrefix := append([]byte{0x04}, valid...)
		_, err := business.ValidateKeyMaterial(
			business.KeyTypeP256WebauthnPublicKey, withPrefix, webAuthnExtra(),
		)
		require.Error(t, err)
	})

	t.Run("rejects coordinates that are not on P-256", func(t *testing.T) {
		offCurve := make([]byte, len(valid))
		copy(offCurve, valid)
		offCurve[63] ^= 0xff
		_, err := business.ValidateKeyMaterial(
			business.KeyTypeP256WebauthnPublicKey, offCurve, webAuthnExtra(),
		)
		require.Error(t, err)
	})

	missingCases := []struct {
		name   string
		mutate func(data.JSONMap)
	}{
		{name: "credential_id", mutate: func(m data.JSONMap) { delete(m, "credential_id") }},
		{name: "rp_id", mutate: func(m data.JSONMap) { delete(m, "rp_id") }},
		{name: "cose_alg", mutate: func(m data.JSONMap) { delete(m, "cose_alg") }},
		{name: "a short credential_id", mutate: func(m data.JSONMap) {
			m["credential_id"] = base64.RawURLEncoding.EncodeToString([]byte("short"))
		}},
		{name: "a credential_id that is not base64url", mutate: func(m data.JSONMap) {
			m["credential_id"] = "not base64url!!"
		}},
		{name: "an algorithm other than ES256", mutate: func(m data.JSONMap) { m["cose_alg"] = float64(-257) }},
		{name: "an unknown transport", mutate: func(m data.JSONMap) { m["transports"] = []any{"telepathy"} }},
		{name: "backup_state without backup_eligible", mutate: func(m data.JSONMap) { m["backup_state"] = true }},
		{name: "an unknown field", mutate: func(m data.JSONMap) { m["private_key"] = "never" }},
	}
	for _, tc := range missingCases {
		t.Run("rejects "+tc.name, func(t *testing.T) {
			extra := webAuthnExtra()
			tc.mutate(extra)
			_, err := business.ValidateKeyMaterial(business.KeyTypeP256WebauthnPublicKey, valid, extra)
			require.Error(t, err)
		})
	}
}

func TestKeyTypeNaming(t *testing.T) {
	assert.Equal(t, "secp256k1", business.KeyTypeName(business.KeyTypeSecp256k1PublicKey))
	assert.Equal(t, "p256-webauthn", business.KeyTypeName(business.KeyTypeP256WebauthnPublicKey))
	assert.True(t, business.IsPublicKeyType(business.KeyTypeSecp256k1PublicKey))
	assert.True(t, business.IsPublicKeyType(business.KeyTypeP256WebauthnPublicKey))
	assert.False(t, business.IsPublicKeyType(devicev1.KeyType_FCM_TOKEN),
		"a token may be secret and must never be published")
	assert.False(t, business.IsPublicKeyType(devicev1.KeyType_PICKLE_KEY))
}
