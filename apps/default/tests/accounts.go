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

package tests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pitabwire/frame/v2"
	"github.com/pitabwire/frame/v2/datastore"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/service/accounts"
	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

// TestIdentityKey keys the static salter in tests.
const TestIdentityKey = "0123456789abcdef0123456789abcdef"

// TestAccountParams are the factory, creation code and version of stawi's
// contracts derivation vector (accounts/testdata/account_derivation.json).
func TestAccountParams(t *testing.T) accounts.Params {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	path := filepath.Join(filepath.Dir(file), "..", "service", "accounts", "testdata", "account_derivation.json")
	raw, err := os.ReadFile(path)
	require.NoError(t, err)

	var v struct {
		Factory        string `json:"factory"`
		AccountVersion uint32 `json:"account_version"`
		CreationCode   string `json:"creation_code"`
	}
	require.NoError(t, json.Unmarshal(raw, &v))
	factory, err := accounts.ParseAddress(v.Factory)
	require.NoError(t, err)
	code, err := hex.DecodeString(strings.TrimPrefix(v.CreationCode, "0x"))
	require.NoError(t, err)
	return accounts.Params{Factory: factory, CreationCode: code, Version: v.AccountVersion}
}

// NewTestDeriver derives accounts with the static test key.
func NewTestDeriver(t *testing.T) *accounts.Deriver {
	t.Helper()
	salter, err := accounts.NewStaticSalter([]byte(TestIdentityKey))
	require.NoError(t, err)
	deriver, err := accounts.NewDeriver(salter, TestAccountParams(t))
	require.NoError(t, err)
	return deriver
}

// NewAccountBusiness builds the account business over the service's pool. A
// nil deriver disables derivation.
func NewAccountBusiness(
	ctx context.Context,
	svc *frame.Service,
	deriver *accounts.Deriver,
	batchSize int,
) business.AccountBusiness {
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)
	workMan := svc.WorkManager()
	return business.NewAccountBusiness(
		deriver,
		repository.NewProfileAccountRepository(ctx, dbPool, workMan),
		repository.NewProfileRepository(ctx, dbPool, workMan),
		batchSize,
	)
}
