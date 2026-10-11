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
	"fmt"
	"os"
	"strings"

	"github.com/pitabwire/util"
	"golang.org/x/crypto/sha3"

	"github.com/antinvestor/service-profile/apps/default/config"
)

// FromConfig builds the deriver the configuration describes.
//
// With nothing configured (neither ACCOUNTS_SALT_FILE nor STAWI_ACCOUNT_*) it
// returns (nil, nil) and logs a warning: profiles are created without
// accounts, which is the development posture. A configuration that is present
// but wrong, including a key file that is missing, empty or under 32 bytes, is
// an error, and the service refuses to start.
func FromConfig(ctx context.Context, cfg *config.ProfileConfig) (*Deriver, error) {
	accountConfigured := cfg.AccountFactory != "" ||
		cfg.AccountCreationCode != "" || cfg.AccountCreationCodeFile != ""

	if cfg.AccountsSaltFile == "" {
		if accountConfigured {
			return nil, fmt.Errorf("%w: STAWI_ACCOUNT_* is set but ACCOUNTS_SALT_FILE is not", ErrAccountConfig)
		}
		util.Log(ctx).Warn("ACCOUNTS_SALT_FILE is not set: profile accounts are not derived")
		return nil, nil //nolint:nilnil // no deriver is a valid, logged state
	}

	salter, err := NewHMACSalterFromFile(cfg.AccountsSaltFile)
	if err != nil {
		return nil, fmt.Errorf("ACCOUNTS_SALT_FILE: %w", err)
	}
	params, err := ParamsFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return NewDeriver(salter, params)
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

func decodeHex(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	return hex.DecodeString(s)
}
