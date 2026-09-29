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

package business

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/pitabwire/frame/v2/data"

	"github.com/antinvestor/service-profile/apps/devices/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// Domain facts published by the devices application (GFOS §10.1, K5). Each is
// staged in the transaction that performs the state change and relayed once.
const (
	// FactDeviceLinked states that a device now belongs to a profile.
	FactDeviceLinked = "device.linked"
	// FactDeviceKeyAdded states that a key was registered on a device.
	FactDeviceKeyAdded = "device.key.added"
	// FactDeviceKeyRemoved states that a key was withdrawn from a device.
	FactDeviceKeyRemoved = "device.key.removed"

	aggregateDevice = "device"
)

// DeviceLinkedFact builds the device.linked fact.
//
// Payload: profile_id, device_id, session_id, device_name, os, linked_at,
// plus the envelope fields event_id, name and occurred_at.
func DeviceLinkedFact(ctx context.Context, device *models.Device, sessionID string) *outbox.Event {
	return outbox.NewEvent(ctx, FactDeviceLinked, aggregateDevice, device.GetID(), data.JSONMap{
		"profile_id":  device.ProfileID,
		"device_id":   device.GetID(),
		"session_id":  sessionID,
		"device_name": device.Name,
		"os":          device.OS,
		"linked_at":   time.Now().UTC().Format(time.RFC3339Nano),
	})
}

// DeviceKeyAddedFact builds the device.key.added fact.
//
// Payload: profile_id, device_id, key_id (the devices-service row id),
// derived_key_id ("0x" + 40 hex, what GFOS resolves as device_key_id),
// key_type, key_type_code, public_key (standard base64, present only for key
// types that hold public material), key_sha256, extra, created_at and
// expires_at, plus the envelope fields.
//
// Key material that may be secret — an FCM token, a pickle key, a Matrix
// session key — is never carried: the fact then states only that a key of that
// type exists, identified by its sha256.
func DeviceKeyAddedFact(ctx context.Context, key *models.DeviceKey, profileID string) *outbox.Event {
	payload := keyFactPayload(key, profileID)
	payload["created_at"] = key.CreatedAt.UTC().Format(time.RFC3339Nano)
	if key.ExpiresAt != nil && !key.ExpiresAt.IsZero() {
		payload["expires_at"] = key.ExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	return outbox.NewEvent(ctx, FactDeviceKeyAdded, aggregateDevice, key.DeviceID, payload)
}

// DeviceKeyRemovedFact builds the device.key.removed fact. It carries the same
// identification as device.key.added plus removed_at, so a consumer can revoke
// a credential without holding any earlier state.
func DeviceKeyRemovedFact(ctx context.Context, key *models.DeviceKey, profileID string) *outbox.Event {
	payload := keyFactPayload(key, profileID)
	payload["removed_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	return outbox.NewEvent(ctx, FactDeviceKeyRemoved, aggregateDevice, key.DeviceID, payload)
}

func keyFactPayload(key *models.DeviceKey, profileID string) data.JSONMap {
	payload := data.JSONMap{
		"profile_id":    profileID,
		"device_id":     key.DeviceID,
		"key_id":        key.GetID(),
		"key_type":      KeyTypeName(key.KeyType),
		"key_type_code": int32(key.KeyType),
		"key_sha256":    keySHA256(key.Key),
	}

	if IsPublicKeyType(key.KeyType) {
		payload["public_key"] = base64.StdEncoding.EncodeToString(key.Key)
		payload["extra"] = map[string]any(key.Extra)
		if derived, ok := key.Extra["key_id"].(string); ok {
			payload["derived_key_id"] = derived
		}
	}

	return payload
}
