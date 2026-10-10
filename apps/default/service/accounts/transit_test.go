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
	"encoding/base64"
	"fmt"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/antinvestor/service-profile/apps/default/service/accounts"
)

const (
	vaultImage     = "hashicorp/vault:1.20"
	vaultRootToken = "root-token"
	transitKey     = "stawi-identity"
)

// startVault runs a Vault dev server with Transit enabled and an exportable
// key, so the test can compare Transit's HMAC with a local one.
func startVault(t *testing.T) (string, []byte) {
	t.Helper()
	ctx := t.Context()

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        vaultImage,
			ExposedPorts: []string{"8200/tcp"},
			Env: map[string]string{
				"VAULT_DEV_ROOT_TOKEN_ID":  vaultRootToken,
				"VAULT_DEV_LISTEN_ADDRESS": "0.0.0.0:8200",
				"SKIP_SETCAP":              "true",
			},
			WaitingFor: wait.ForHTTP("/v1/sys/health").WithPort("8200/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })

	endpoint, err := ctr.PortEndpoint(ctx, "8200/tcp", "http")
	require.NoError(t, err)

	cfg := vault.DefaultConfig()
	cfg.Address = endpoint
	client, err := vault.NewClient(cfg)
	require.NoError(t, err)
	client.SetToken(vaultRootToken)

	require.NoError(t, client.Sys().MountWithContext(ctx, "transit", &vault.MountInput{Type: "transit"}))
	_, err = client.Logical().WriteWithContext(ctx, "transit/keys/"+transitKey, map[string]any{
		"type":       "aes256-gcm96",
		"exportable": true,
	})
	require.NoError(t, err)

	exported, err := client.Logical().ReadWithContext(ctx, "transit/export/hmac-key/"+transitKey+"/1")
	require.NoError(t, err)
	keys, ok := exported.Data["keys"].(map[string]any)
	require.True(t, ok)
	encoded, ok := keys["1"].(string)
	require.True(t, ok)
	key, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)

	return endpoint, key
}

func localSalt(key []byte, profileID string) [32]byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("stawi/identity/v1" + profileID))
	var out [32]byte
	copy(out[:], mac.Sum(nil))
	return out
}

func TestTransitSalterAgainstVault(t *testing.T) {
	endpoint, key := startVault(t)
	ctx := t.Context()

	salter, err := accounts.NewTransitSalter(accounts.TransitConfig{
		Address: endpoint, AuthMethod: accounts.AuthToken, Token: vaultRootToken,
		TransitMount: "transit", Key: transitKey, KeyVersion: 1,
	})
	require.NoError(t, err)

	salt, err := salter.Salt(ctx, "profile-one")
	require.NoError(t, err)
	assert.Equal(t, localSalt(key, "profile-one"), salt, "transit HMAC equals HMAC-SHA256 under K_identity")

	ids := make([]string, 25)
	for i := range ids {
		ids[i] = fmt.Sprintf("profile-%02d", i)
	}
	batch, err := salter.SaltBatch(ctx, ids)
	require.NoError(t, err)
	require.Len(t, batch, len(ids))
	for i, id := range ids {
		assert.Equal(t, localSalt(key, id), batch[i], id)
	}

	// The static salter keyed with the same material agrees, so local stacks
	// and production derive identically from the same key.
	static, err := accounts.NewStaticSalter(key)
	require.NoError(t, err)
	staticSalt, err := static.Salt(ctx, "profile-one")
	require.NoError(t, err)
	assert.Equal(t, salt, staticSalt)

	// A rotation must not move addresses: the pinned version keeps answering.
	admin, err := vault.NewClient(&vault.Config{Address: endpoint})
	require.NoError(t, err)
	admin.SetToken(vaultRootToken)
	_, err = admin.Logical().WriteWithContext(ctx, "transit/keys/"+transitKey+"/rotate", nil)
	require.NoError(t, err)
	afterRotate, err := salter.Salt(ctx, "profile-one")
	require.NoError(t, err)
	assert.Equal(t, salt, afterRotate)

	// A version that does not exist is an error, never a different salt.
	wrong, err := accounts.NewTransitSalter(accounts.TransitConfig{
		Address: endpoint, AuthMethod: accounts.AuthToken, Token: vaultRootToken,
		TransitMount: "transit", Key: transitKey, KeyVersion: 9,
	})
	require.NoError(t, err)
	_, err = wrong.Salt(ctx, "profile-one")
	require.Error(t, err)
}
