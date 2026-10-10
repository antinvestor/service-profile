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
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	vault "github.com/hashicorp/vault/api"
	k8sauth "github.com/hashicorp/vault/api/auth/kubernetes"
)

// tokenRenewMargin re-authenticates this far before the Vault token expires.
const tokenRenewMargin = 30 * time.Second

// ErrTransitConfig reports an incomplete Transit configuration.
var ErrTransitConfig = errors.New("accounts: incomplete vault transit configuration")

// TransitConfig locates the Transit key and how to authenticate to Vault.
type TransitConfig struct {
	// Address is the Vault address (VAULT_ADDR).
	Address string
	// AuthRole is the Kubernetes auth role. Empty means the client token
	// already in the environment (VAULT_TOKEN) is used, as on a dev server.
	AuthRole string
	// AuthMount is the Kubernetes auth mount path (default "kubernetes").
	AuthMount string
	// ServiceAccountTokenPath overrides the projected service-account token
	// path; empty uses the in-pod default.
	ServiceAccountTokenPath string
	// TransitMount is the Transit secrets engine mount (default "transit").
	TransitMount string
	// Key is the Transit key name (default "stawi-identity").
	Key string
	// KeyVersion pins the key version so a rotation, which must never
	// happen to this key, cannot silently move every address.
	KeyVersion int
}

// TransitSalter computes identity salts with Vault Transit
// (transit/hmac/<key>/sha2-256), so K_identity never leaves Vault.
type TransitSalter struct {
	client *vault.Client
	cfg    TransitConfig

	mu        sync.Mutex
	expiresAt time.Time
}

// NewTransitSalter builds a salter; it authenticates lazily on first use.
func NewTransitSalter(cfg TransitConfig) (*TransitSalter, error) {
	if cfg.Address == "" || cfg.Key == "" || cfg.TransitMount == "" {
		return nil, fmt.Errorf("%w: address, transit mount and key are required", ErrTransitConfig)
	}
	if cfg.KeyVersion < 1 {
		return nil, fmt.Errorf("%w: key version must be pinned (>= 1)", ErrTransitConfig)
	}
	if cfg.AuthMount == "" {
		cfg.AuthMount = "kubernetes"
	}

	vcfg := vault.DefaultConfig()
	if vcfg.Error != nil {
		return nil, vcfg.Error
	}
	vcfg.Address = cfg.Address
	client, err := vault.NewClient(vcfg)
	if err != nil {
		return nil, err
	}
	return &TransitSalter{client: client, cfg: cfg}, nil
}

// login authenticates with Kubernetes auth when a role is configured and the
// current token is missing or close to expiry.
func (t *TransitSalter) login(ctx context.Context) error {
	if t.cfg.AuthRole == "" {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.expiresAt.IsZero() && time.Now().Before(t.expiresAt) {
		return nil
	}

	opts := []k8sauth.LoginOption{k8sauth.WithMountPath(t.cfg.AuthMount)}
	if t.cfg.ServiceAccountTokenPath != "" {
		opts = append(opts, k8sauth.WithServiceAccountTokenPath(t.cfg.ServiceAccountTokenPath))
	}
	auth, err := k8sauth.NewKubernetesAuth(t.cfg.AuthRole, opts...)
	if err != nil {
		return fmt.Errorf("accounts: vault kubernetes auth: %w", err)
	}
	secret, err := t.client.Auth().Login(ctx, auth)
	if err != nil {
		return fmt.Errorf("accounts: vault kubernetes login: %w", err)
	}
	if secret == nil || secret.Auth == nil {
		return errors.New("accounts: vault kubernetes login returned no token")
	}
	ttl := time.Duration(secret.Auth.LeaseDuration) * time.Second
	t.expiresAt = time.Now().Add(max(ttl-tokenRenewMargin, ttl/2))
	return nil
}

// Salt implements Salter.
func (t *TransitSalter) Salt(ctx context.Context, profileID string) ([32]byte, error) {
	salts, err := t.SaltBatch(ctx, []string{profileID})
	if err != nil {
		return [32]byte{}, err
	}
	return salts[0], nil
}

// SaltBatch implements Salter with one Transit batch request.
func (t *TransitSalter) SaltBatch(ctx context.Context, profileIDs []string) ([][32]byte, error) {
	if len(profileIDs) == 0 {
		return nil, nil
	}
	if err := t.login(ctx); err != nil {
		return nil, err
	}

	batch := make([]map[string]any, len(profileIDs))
	for i, id := range profileIDs {
		batch[i] = map[string]any{"input": base64.StdEncoding.EncodeToString(IdentityMessage(id))}
	}
	path := fmt.Sprintf("%s/hmac/%s/sha2-256", t.cfg.TransitMount, t.cfg.Key)
	secret, err := t.client.Logical().WriteWithContext(ctx, path, map[string]any{
		"batch_input": batch,
		"key_version": t.cfg.KeyVersion,
	})
	if err != nil {
		t.forgetToken()
		return nil, fmt.Errorf("accounts: vault transit hmac: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("accounts: vault transit hmac returned no data")
	}
	results, ok := secret.Data["batch_results"].([]any)
	if !ok || len(results) != len(profileIDs) {
		return nil, fmt.Errorf("accounts: vault transit hmac returned %d results for %d inputs",
			len(results), len(profileIDs))
	}

	out := make([][32]byte, len(profileIDs))
	for i, r := range results {
		salt, parseErr := parseTransitHMAC(r, t.cfg.KeyVersion)
		if parseErr != nil {
			return nil, fmt.Errorf("accounts: transit hmac for input %d: %w", i, parseErr)
		}
		out[i] = salt
	}
	return out, nil
}

func (t *TransitSalter) forgetToken() {
	t.mu.Lock()
	t.expiresAt = time.Time{}
	t.mu.Unlock()
}

// parseTransitHMAC reads one batch result: "vault:v<N>:<base64 mac>".
func parseTransitHMAC(result any, wantVersion int) ([32]byte, error) {
	entry, ok := result.(map[string]any)
	if !ok {
		return [32]byte{}, errors.New("malformed batch result")
	}
	if msg, _ := entry["error"].(string); msg != "" {
		return [32]byte{}, errors.New(msg)
	}
	value, _ := entry["hmac"].(string)
	prefix := fmt.Sprintf("vault:v%d:", wantVersion)
	if !strings.HasPrefix(value, prefix) {
		return [32]byte{}, fmt.Errorf("hmac not produced by key version %d", wantVersion)
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, prefix))
	if err != nil {
		return [32]byte{}, err
	}
	if len(raw) != sha256Size {
		return [32]byte{}, fmt.Errorf("hmac is %d bytes, want %d", len(raw), sha256Size)
	}
	var out [32]byte
	copy(out[:], raw)
	return out, nil
}

const sha256Size = 32
