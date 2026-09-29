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

package repository_test

import (
	"testing"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/tests"
)

// secp256k1KeyType and p256KeyType mirror device.v1.KeyType.SECP256K1_PUBLIC_KEY
// and P256_WEBAUTHN_PUBLIC_KEY.
const (
	secp256k1KeyType devicev1.KeyType = 6
	p256KeyType      devicev1.KeyType = 7
)

type DeviceKeyRepositoryTestSuite struct {
	tests.DeviceBaseTestSuite
}

func TestDeviceKeyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceKeyRepositoryTestSuite))
}

// TestKeyMaterialUniquenessIsEnforcedByTheDatabase proves the constraint is in
// the schema and not only in the business layer: two writes of the same
// (key_type, key) cannot both land, whichever code path made them.
func (suite *DeviceKeyRepositoryTestSuite) TestKeyMaterialUniquenessIsEnforcedByTheDatabase() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		firstDevice := &models.Device{Name: "first", OS: "Linux"}
		require.NoError(t, deps.DeviceRepo.Create(ctx, firstDevice))

		secondDevice := &models.Device{Name: "second", OS: "Linux"}
		require.NoError(t, deps.DeviceRepo.Create(ctx, secondDevice))

		material := []byte(util.RandomAlphaNumericString(48))

		original := &models.DeviceKey{
			DeviceID: firstDevice.GetID(),
			KeyType:  secp256k1KeyType,
			Key:      material,
		}
		require.NoError(t, deps.KeyRepo.Create(ctx, original))

		t.Run("the same material cannot be registered on another device", func(t *testing.T) {
			duplicate := &models.DeviceKey{
				DeviceID: secondDevice.GetID(),
				KeyType:  secp256k1KeyType,
				Key:      material,
			}
			err := deps.KeyRepo.Create(ctx, duplicate)
			require.Error(t, err, "the unique index must reject the duplicate")
			assert.Contains(t, err.Error(), "idx_device_keys_key_type_key")
		})

		t.Run("the same material cannot be registered twice on one device", func(t *testing.T) {
			duplicate := &models.DeviceKey{
				DeviceID: firstDevice.GetID(),
				KeyType:  secp256k1KeyType,
				Key:      material,
			}
			require.Error(t, deps.KeyRepo.Create(ctx, duplicate))
		})

		t.Run("the same bytes under another key type are a different key", func(t *testing.T) {
			other := &models.DeviceKey{
				DeviceID: secondDevice.GetID(),
				KeyType:  p256KeyType,
				Key:      material,
			}
			require.NoError(t, deps.KeyRepo.Create(ctx, other))
		})

		t.Run("large key material is unique without hitting the index row limit", func(t *testing.T) {
			// Pickled Olm sessions and Matrix keys run to several KB; a btree
			// entry over ~2.7KB cannot be stored, so the index must not hold
			// the raw bytes.
			large := []byte(util.RandomAlphaNumericString(8192))
			first := &models.DeviceKey{DeviceID: firstDevice.GetID(), KeyType: devicev1.KeyType_PICKLE_KEY, Key: large}
			require.NoError(t, deps.KeyRepo.Create(ctx, first))

			duplicate := &models.DeviceKey{
				DeviceID: secondDevice.GetID(),
				KeyType:  devicev1.KeyType_PICKLE_KEY,
				Key:      large,
			}
			err := deps.KeyRepo.Create(ctx, duplicate)
			require.Error(t, err, "the unique index must reject the duplicate")
			assert.Contains(t, err.Error(), "idx_device_keys_key_type_key")
		})

		t.Run("material may be registered again once the key is withdrawn", func(t *testing.T) {
			removed, err := deps.KeyRepo.RemoveByID(ctx, original.GetID())
			require.NoError(t, err)
			require.NotNil(t, removed)

			reRegistered := &models.DeviceKey{
				DeviceID: firstDevice.GetID(),
				KeyType:  secp256k1KeyType,
				Key:      material,
			}
			require.NoError(t, deps.KeyRepo.Create(ctx, reRegistered))
		})
	})
}
