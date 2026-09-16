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
	"context"
	"encoding/base64"
	"testing"
	"time"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"github.com/pitabwire/frame/v2/data"
	frevents "github.com/pitabwire/frame/v2/events"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/apps/devices/tests"
)

const factWaitTimeout = 10 * time.Second

type DeviceFactsTestSuite struct {
	tests.DeviceBaseTestSuite
}

func TestDeviceFactsTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceFactsTestSuite))
}

// addTestDevice creates a device row directly so a fact test starts from a
// known state without going through the asynchronous device creation queue.
func (suite *DeviceFactsTestSuite) addTestDevice(
	ctx context.Context,
	deps *tests.DepsBuilder,
	profileID string,
) *models.Device {
	device := &models.Device{Name: "Fact Device", OS: "Linux", ProfileID: profileID}
	suite.Require().NoError(deps.DeviceRepo.Create(ctx, device))
	return device
}

// waitForFacts blocks until at least one fact with the given name has reached
// the topic, and returns everything recorded under that name.
func waitForFacts(t *testing.T, recorder *tests.FactRecorder, name string) []tests.RecordedFact {
	t.Helper()
	var facts []tests.RecordedFact
	require.Eventually(t, func() bool {
		facts = recorder.Named(name)
		return len(facts) > 0
	}, factWaitTimeout, 25*time.Millisecond, "no %s fact reached the topic", name)
	return facts
}

// assertFactHeaders checks the routing headers every published fact carries.
// frevents.EventHeaderName is the one frame's own events.Manager routes on, so
// a consumer can register these facts as frame events rather than writing a
// queue worker of its own; event_name is the plain-text alias and event_id is
// the deduplication key.
func assertFactHeaders(t *testing.T, fact tests.RecordedFact, name string) {
	t.Helper()
	assert.Equal(t, name, fact.Metadata[frevents.EventHeaderName],
		"the frame event header must carry the fact name")
	assert.Equal(t, name, fact.Metadata["event_name"])
	assert.Equal(t, fact.EventID, fact.Metadata["event_id"])
}

// TestKeyAddedFactIsPublishedExactlyOnce is the K5 guarantee end to end: the
// fact is staged in the transaction that registers the key, the relay publishes
// it once, and a second relay pass — the shape a redelivery or a restart takes
// — publishes nothing further.
func (suite *DeviceFactsTestSuite) TestKeyAddedFactIsPublishedExactlyOnce() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc, deps := suite.CreateService(t, dep)

		profileID := util.IDString()
		device := suite.addTestDevice(ctx, deps, profileID)

		key := newP256Key(t)
		extra := webAuthnExtra()

		stored, err := deps.KeyBusiness.AddKey(ctx, device.GetID(), business.KeyTypeP256WebauthnPublicKey, key, extra)
		require.NoError(t, err)

		staged, err := tests.StagedFacts(ctx, svc, business.FactDeviceKeyAdded)
		require.NoError(t, err)
		require.Len(t, staged, 1, "the key write must stage exactly one fact")
		require.Nil(t, staged[0].PublishedAt, "a staged fact is not published yet")

		published, err := deps.FactRelay.Drain(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, published)

		facts := waitForFacts(t, deps.Facts, business.FactDeviceKeyAdded)
		require.Len(t, facts, 1)

		fact := facts[0]
		assert.Equal(t, profileID, fact.Payload["profile_id"])
		assert.Equal(t, device.GetID(), fact.Payload["device_id"])
		assert.Equal(t, stored.GetId(), fact.Payload["key_id"])
		assert.Equal(t, "p256-webauthn", fact.Payload["key_type"])
		assert.Equal(t, base64.StdEncoding.EncodeToString(key), fact.Payload["public_key"])
		assert.NotEmpty(t, fact.Payload["derived_key_id"])
		assert.NotEmpty(t, fact.Payload["created_at"])
		assert.NotEmpty(t, fact.EventID)
		assertFactHeaders(t, fact, business.FactDeviceKeyAdded)

		republished, err := deps.FactRelay.Drain(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, republished, "a second relay pass must not republish a delivered fact")

		assert.Len(t, deps.Facts.Named(business.FactDeviceKeyAdded), 1)
	})
}

// TestDuplicateKeyRegistrationStagesNoSecondFact covers both duplicate cases:
// the client retrying on the same device is idempotent, and another device
// claiming the same material is refused with nothing written at all.
func (suite *DeviceFactsTestSuite) TestDuplicateKeyRegistrationStagesNoSecondFact() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc, deps := suite.CreateService(t, dep)

		device := suite.addTestDevice(ctx, deps, util.IDString())
		other := suite.addTestDevice(ctx, deps, util.IDString())

		key := newP256Key(t)
		extra := webAuthnExtra()

		first, err := deps.KeyBusiness.AddKey(
			ctx, device.GetID(), business.KeyTypeP256WebauthnPublicKey, key, extra,
		)
		require.NoError(t, err)

		again, err := deps.KeyBusiness.AddKey(
			ctx, device.GetID(), business.KeyTypeP256WebauthnPublicKey, key, webAuthnExtra(),
		)
		require.NoError(t, err, "re-registering the same key on the same device is a retry")
		assert.Equal(t, first.GetId(), again.GetId())

		_, err = deps.KeyBusiness.AddKey(
			ctx, other.GetID(), business.KeyTypeP256WebauthnPublicKey, key, webAuthnExtra(),
		)
		require.Error(t, err, "one public key belongs to one device")

		staged, err := tests.StagedFacts(ctx, svc, business.FactDeviceKeyAdded)
		require.NoError(t, err)
		assert.Len(t, staged, 1, "only the first registration is a fact")
	})
}

// TestKeyRemovedFactIsPublishedWithTheWithdrawal checks the removal fact
// carries what a consumer needs to revoke a credential.
func (suite *DeviceFactsTestSuite) TestKeyRemovedFactIsPublishedWithTheWithdrawal() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc, deps := suite.CreateService(t, dep)

		profileID := util.IDString()
		device := suite.addTestDevice(ctx, deps, profileID)
		key := newP256Key(t)

		stored, err := deps.KeyBusiness.AddKey(
			ctx, device.GetID(), business.KeyTypeP256WebauthnPublicKey, key, webAuthnExtra(),
		)
		require.NoError(t, err)

		removeResult, err := deps.KeyBusiness.RemoveKeys(ctx, stored.GetId())
		require.NoError(t, err)
		for res := range removeResult {
			require.NoError(t, res.Error())
		}

		staged, err := tests.StagedFacts(ctx, svc, business.FactDeviceKeyRemoved)
		require.NoError(t, err)
		require.Len(t, staged, 1)

		_, err = deps.FactRelay.Drain(ctx)
		require.NoError(t, err)

		facts := waitForFacts(t, deps.Facts, business.FactDeviceKeyRemoved)
		fact := facts[0]
		assert.Equal(t, profileID, fact.Payload["profile_id"])
		assert.Equal(t, device.GetID(), fact.Payload["device_id"])
		assert.Equal(t, stored.GetId(), fact.Payload["key_id"])
		assert.Equal(t, "p256-webauthn", fact.Payload["key_type"])
		assert.NotEmpty(t, fact.Payload["removed_at"])
		assertFactHeaders(t, fact, business.FactDeviceKeyRemoved)
	})
}

// TestSecretKeyTypesArePublishedWithoutTheirMaterial guards the rule that a
// fact never carries anything that might be a secret.
func (suite *DeviceFactsTestSuite) TestSecretKeyTypesArePublishedWithoutTheirMaterial() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, _, deps := suite.CreateService(t, dep)

		device := suite.addTestDevice(ctx, deps, util.IDString())
		token := []byte(util.RandomAlphaNumericString(64))

		_, err := deps.KeyBusiness.AddKey(
			ctx, device.GetID(), devicev1.KeyType_FCM_TOKEN, token, data.JSONMap{"app_id": "stawi"},
		)
		require.NoError(t, err)

		_, err = deps.FactRelay.Drain(ctx)
		require.NoError(t, err)

		facts := waitForFacts(t, deps.Facts, business.FactDeviceKeyAdded)
		fact := facts[0]
		assert.NotContains(t, fact.Payload, "public_key", "a token is never published")
		assert.NotContains(t, fact.Payload, "extra")
		assert.NotEmpty(t, fact.Payload["key_sha256"], "the fact still identifies the key")
		assertFactHeaders(t, fact, business.FactDeviceKeyAdded)
	})
}

// TestDeviceLinkedFactIsPublishedWithTheLink checks the device.linked fact.
func (suite *DeviceFactsTestSuite) TestDeviceLinkedFactIsPublishedWithTheLink() {
	suite.WithTestDependencies(suite.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc, deps := suite.CreateService(t, dep)

		device := suite.addTestDevice(ctx, deps, "")
		session := &models.DeviceSession{
			DeviceID:  device.GetID(),
			UserAgent: "Test Agent",
			IP:        "127.0.0.1",
			LastSeen:  time.Now(),
		}
		require.NoError(t, deps.SessionRepo.Create(ctx, session))

		profileID := util.IDString()
		linked, err := deps.DeviceBusiness.LinkDeviceToProfile(ctx, session.GetID(), profileID, data.JSONMap{})
		require.NoError(t, err)
		require.NotNil(t, linked)

		staged, err := tests.StagedFacts(ctx, svc, business.FactDeviceLinked)
		require.NoError(t, err)
		require.Len(t, staged, 1)

		published, err := deps.FactRelay.Drain(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, published)

		facts := waitForFacts(t, deps.Facts, business.FactDeviceLinked)
		fact := facts[0]
		assert.Equal(t, profileID, fact.Payload["profile_id"])
		assert.Equal(t, device.GetID(), fact.Payload["device_id"])
		assert.Equal(t, session.GetID(), fact.Payload["session_id"])
		assertFactHeaders(t, fact, business.FactDeviceLinked)

		// Linking again is a no-op, so there is no second fact.
		_, err = deps.DeviceBusiness.LinkDeviceToProfile(ctx, session.GetID(), profileID, data.JSONMap{})
		require.NoError(t, err)

		staged, err = tests.StagedFacts(ctx, svc, business.FactDeviceLinked)
		require.NoError(t, err)
		assert.Len(t, staged, 1)
	})
}
