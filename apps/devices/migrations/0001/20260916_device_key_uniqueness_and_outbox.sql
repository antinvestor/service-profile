--- Copyright 2023-2026 Ant Investor Ltd
---
--- Licensed under the Apache License, Version 2.0 (the "License");
--- you may not use this file except in compliance with the License.
--- You may obtain a copy of the License at
---
---      http://www.apache.org/licenses/LICENSE-2.0
---
--- Unless required by applicable law or agreed to in writing, software
--- distributed under the License is distributed on an "AS IS" BASIS,
--- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
--- See the License for the specific language governing permissions and
--- limitations under the License.

-- GFOS platform change K4: a piece of device key material is registered once.
--
-- Without this the same public key could be bound to two devices, and a
-- consumer resolving a key id to a device would have no single answer. The
-- index is partial on deleted_at so a key that was withdrawn can be registered
-- again, which is what key rotation looks like.
--
-- Historic duplicates are retired before the index is built, keeping the most
-- recent registration. They are possible for the pre-existing token types (an
-- FCM token that moved between device rows on reinstall); the cryptographic
-- types are new here and cannot have any.
UPDATE device_keys AS dk
SET deleted_at = NOW()
WHERE dk.deleted_at IS NULL
  AND EXISTS (SELECT 1
              FROM device_keys AS newer
              WHERE newer.deleted_at IS NULL
                AND newer.key_type = dk.key_type
                AND newer.key = dk.key
                AND (newer.created_at, newer.id) > (dk.created_at, dk.id));

CREATE UNIQUE INDEX IF NOT EXISTS idx_device_keys_key_type_key
    ON device_keys (key_type, key)
    WHERE deleted_at IS NULL;

-- GFOS platform change K5: the outbox relay claims the oldest staged facts
-- that have not been published yet, so that is the access path to index.
CREATE INDEX IF NOT EXISTS idx_outbox_events_unpublished
    ON outbox_events (created_at)
    WHERE published_at IS NULL;
