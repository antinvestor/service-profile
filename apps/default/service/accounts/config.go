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
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pitabwire/util"
	"golang.org/x/crypto/sha3"

	"github.com/antinvestor/service-profile/apps/default/config"
)

// ErrTokenAuthInProduction refuses static Vault tokens in production unless
// VAULT_ALLOW_TOKEN_AUTH=true.
var ErrTokenAuthInProduction = errors.New(
	"accounts: VAULT_AUTH_METHOD=token is refused in production unless VAULT_ALLOW_TOKEN_AUTH=true")

// ErrStaticKeyInProduction refuses the static identity key in production.
var ErrStaticKeyInProduction = errors.New(
	"accounts: STAWI_IDENTITY_STATIC_KEY is refused in production; configure Vault Transit")

// FromConfig builds the deriver the configuration describes.
//
// With no salter configured (neither VAULT_ADDR nor STAWI_IDENTITY_STATIC_KEY)
// it returns (nil, nil) and logs a warning: profiles are created without
// accounts, which is the development posture. A configuration that is present
// but wrong is an error.
func FromConfig(ctx context.Context, cfg *config.ProfileConfig) (*Deriver, error) {
	salter, err := SalterFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	accountConfigured := cfg.AccountFactory != "" ||
		cfg.AccountCreationCode != "" || cfg.AccountCreationCodeFile != ""

	if salter == nil {
		if accountConfigured {
			return nil, fmt.Errorf("%w: STAWI_ACCOUNT_* is set but no identity salter "+
				"(VAULT_ADDR or STAWI_IDENTITY_STATIC_KEY)", ErrAccountConfig)
		}
		util.Log(ctx).Warn("no identity salter configured: profile accounts are not derived")
		return nil, nil //nolint:nilnil // no deriver is a valid, logged state
	}

	params, err := ParamsFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return NewDeriver(salter, params)
}

// SalterFromConfig returns the Transit salter when VAULT_ADDR is set, the
// static salter when STAWI_IDENTITY_STATIC_KEY is set, or nil.
func SalterFromConfig(cfg *config.ProfileConfig) (Salter, error) {
	if cfg.VaultAddress != "" {
		if cfg.IdentityTransitKeyVer < 1 {
			return nil, fmt.Errorf("%w: STAWI_IDENTITY_KEY_VERSION is required with VAULT_ADDR", ErrTransitConfig)
		}
		method := strings.ToLower(strings.TrimSpace(cfg.VaultAuthMethod))
		if method == AuthToken && isProduction(cfg) && !cfg.VaultAllowTokenAuth {
			return nil, ErrTokenAuthInProduction
		}
		return NewTransitSalter(TransitConfig{
			Address:                 cfg.VaultAddress,
			AuthMethod:              method,
			AuthRole:                cfg.VaultK8sAuthRole,
			AuthMount:               cfg.VaultK8sAuthMount,
			ServiceAccountTokenPath: cfg.VaultK8sTokenPath,
			GCPRole:                 cfg.VaultGCPAuthRole,
			GCPMount:                cfg.VaultGCPAuthMount,
			GCPServiceAccount:       cfg.VaultGCPServiceAccount,
			Token:                   cfg.VaultToken,
			TransitMount:            cfg.VaultTransitMount,
			Key:                     cfg.IdentityTransitKey,
			KeyVersion:              cfg.IdentityTransitKeyVer,
		})
	}
	if cfg.IdentityStaticKeyHex == "" {
		return nil, nil //nolint:nilnil // no salter configured
	}
	if isProduction(cfg) {
		return nil, ErrStaticKeyInProduction
	}
	key, err := decodeHex(cfg.IdentityStaticKeyHex)
	if err != nil {
		return nil, fmt.Errorf("accounts: STAWI_IDENTITY_STATIC_KEY: %w", err)
	}
	return NewStaticSalter(key)
}

// ParamsFromConfig reads the factory, creation code and version.
func ParamsFromConfig(cfg *config.ProfileConfig) (Params, error) {
	var p Params
	factory, err := ParseAddress(cfg.AccountFactory)
	if err != nil {
		return p, fmt.Errorf("%w: STAWI_ACCOUNT_FACTORY: %w", ErrAccountConfig, err)
	}
	p.Factory = factory
	p.Version = cfg.AccountVersion

	codeHex := cfg.AccountCreationCode
	if codeHex == "" && cfg.AccountCreationCodeFile != "" {
		raw, readErr := os.ReadFile(cfg.AccountCreationCodeFile)
		if readErr != nil {
			return p, fmt.Errorf("%w: STAWI_ACCOUNT_CREATION_CODE_FILE: %w", ErrAccountConfig, readErr)
		}
		codeHex = string(raw)
	}
	if codeHex == "" {
		return p, fmt.Errorf("%w: STAWI_ACCOUNT_CREATION_CODE or STAWI_ACCOUNT_CREATION_CODE_FILE is required",
			ErrAccountConfig)
	}
	p.CreationCode, err = decodeHex(codeHex)
	if err != nil {
		return p, fmt.Errorf("%w: account creation code: %w", ErrAccountConfig, err)
	}

	if cfg.AccountCreationCodeHash != "" {
		want, hashErr := decodeHex(cfg.AccountCreationCodeHash)
		if hashErr != nil {
			return p, fmt.Errorf("%w: STAWI_ACCOUNT_CREATION_CODE_HASH: %w", ErrAccountConfig, hashErr)
		}
		h := sha3.NewLegacyKeccak256()
		h.Write(p.CreationCode)
		if got := h.Sum(nil); !bytes.Equal(got, want) {
			return p, fmt.Errorf("%w: creation code hash is 0x%x, manifest says 0x%x", ErrAccountConfig, got, want)
		}
	}
	return p, nil
}

func isProduction(cfg *config.ProfileConfig) bool {
	for _, env := range []string{cfg.DeploymentEnvironment, cfg.Environment()} {
		switch strings.ToLower(strings.TrimSpace(env)) {
		case "prod", "production":
			return true
		}
	}
	return false
}

func decodeHex(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	return hex.DecodeString(s)
}
