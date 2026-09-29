# Device key types and published domain facts

This describes two platform changes the Stawi Group Financial OS (GFOS) requires
of this service — K4 (cryptographic device key types) and K5 (durable domain
facts) — as they are implemented here.

## K4 — cryptographic device key types

`device.v1.KeyType` gains two values:

| Value | Number | Published name | Key bytes |
|-------|--------|----------------|-----------|
| `SECP256K1_PUBLIC_KEY` | 6 | `secp256k1` | 65 bytes, SEC1 uncompressed: `0x04 ‖ x ‖ y` |
| `P256_WEBAUTHN_PUBLIC_KEY` | 7 | `p256-webauthn` | 64 bytes, raw coordinates: `x ‖ y`, no prefix |

Both encodings are the ones the GFOS verifier reads. A compressed secp256k1 key
and a `0x04`-prefixed P-256 key are both rejected, because accepting a second
encoding would mean the same key could be registered twice under two
representations and the uniqueness constraint would not hold.

Validation, in `apps/devices/service/business/keymaterial.go`:

* the length and the prefix byte;
* that the coordinates are a point on the curve — P-256 through
  `crypto/ecdh`, secp256k1 through the curve equation directly, as the standard
  library has no secp256k1;
* the `extra` document against the closed schema below. A field the schema does
  not define is refused, which is what keeps anything unreviewed — a private
  key above all — from riding along on a registration.

### `extra` schema — `SECP256K1_PUBLIC_KEY`

Client supplied:

| Field | Type | Required | Rule |
|-------|------|----------|------|
| `key_id` | string | no | `0x` + 40 hex; must equal the derived id |
| `label` | string | no | ≤ 64 characters |
| `session_proof_id` | string | no | ≤ 128 characters; the authentication-service session proof (GFOS K9) the registration happened under |

Written by the service, always:

| Field | Value |
|-------|-------|
| `key_id` | `0x` + hex(`keccak256(x ‖ y)[12:]`) — the EVM address, which GFOS resolves as `device_key_id` |
| `curve` | `secp256k1` |
| `encoding` | `sec1-uncompressed` |
| `key_sha256` | hex sha256 of the stored key bytes |

### `extra` schema — `P256_WEBAUTHN_PUBLIC_KEY`

Client supplied:

| Field | Type | Required | Rule |
|-------|------|----------|------|
| `credential_id` | string | **yes** | base64url; 16..1023 decoded bytes |
| `rp_id` | string | **yes** | relying party id, ≤ 253 characters |
| `cose_alg` | number | **yes** | `-7` (ES256) — the only algorithm the verifier implements |
| `aaguid` | string | no | 36-character UUID |
| `sign_count` | number | no | ≥ 0 |
| `transports` | string[] | no | subset of `usb`, `nfc`, `ble`, `smart-card`, `hybrid`, `internal`, `cable` |
| `user_verification` | string | no | `required`, `preferred` or `discouraged` |
| `backup_eligible` | bool | no | |
| `backup_state` | bool | no | true only when `backup_eligible` is true |
| `origin` | string | no | ≤ 255 characters |
| `key_id` | string | no | must equal the derived id |
| `label` | string | no | ≤ 64 characters |
| `session_proof_id` | string | no | ≤ 128 characters |

Written by the service, always:

| Field | Value |
|-------|-------|
| `key_id` | `0x` + hex(`keccak256(x ‖ y)[12:]`) — the key id the Stawi account contract derives |
| `curve` | `p-256` |
| `encoding` | `raw-xy` |
| `key_sha256` | hex sha256 of the stored key bytes |

The assertion fields a WebAuthn verification needs at signing time
(`authenticatorData`, `clientDataJSON`, the signature) belong to the
authorization that carries them, not to the registered credential, so they are
not part of this schema.

### Uniqueness

`apps/devices/migrations/0001/20260916_device_key_uniqueness_and_outbox.sql`
creates

```sql
CREATE UNIQUE INDEX idx_device_keys_key_type_key
    ON device_keys (key_type, sha256(key))
    WHERE deleted_at IS NULL;
```

so one piece of key material is registered once, enforced by the database and
not by a read before the write. The index is partial on `deleted_at` so a key
that has been withdrawn may be registered again — that is what rotation looks
like. Duplicates that predate the index are retired by the same migration,
keeping the most recent registration. The index stores a SHA-256 of the
material rather than the bytes, because pickled sessions and Matrix keys can
exceed the btree row size limit.

Consequences at the API:

* the same device re-sending the same material gets its existing key back; this
  is a retry and publishes no second fact;
* another device sending material that is already registered is refused with
  `ALREADY_EXISTS`.

## K5 — durable domain facts

Facts are staged in the same database transaction as the state change they
describe (`internal/outbox`), and one relay publishes them to the service's
registered frame queue publisher. There is no second publishing path.

| Fact | Topic | Aggregate |
|------|-------|-----------|
| `profile.created` | `profile.events` (`QUEUE_PROFILE_EVENTS_NAME`) | profile |
| `device.linked` | `device.events` (`QUEUE_DEVICE_EVENTS_NAME`) | device |
| `device.key.added` | `device.events` | device |
| `device.key.removed` | `device.events` | device |

Every message carries the headers `frame._internal.event.header` (frame's own
`events.EventHeaderName`, so a consumer can register these facts as frame
events instead of writing a queue worker), `event_name` — the same value in
plain text — plus `event_id` and `aggregate_id`. Every body carries `event_id`,
`name` and `occurred_at`. Delivery is
at-least-once: consumers deduplicate on `event_id`.

### Egress must be configured

The relay runs only when the queue URI is set — `QUEUE_PROFILE_EVENTS_URI`
(profile) and `QUEUE_DEVICE_EVENTS_URI` (devices), e.g. a
`gcppubsub://{project}/{topic}` URL. Unset (the default), no publisher or relay
is registered, a warning is logged at startup, and facts stay in
`outbox_events` with `published_at` NULL. Once an egress is configured, the
relay publishes the whole backlog in order. A default in-process queue would
instead mark every fact published while delivering it nowhere.

### Payloads

`profile.created`

```json
{
  "event_id": "…", "name": "profile.created", "occurred_at": "RFC3339Nano",
  "profile_id": "…", "profile_type": "PERSON", "profile_type_uid": 1,
  "profile_type_id": "…", "created_at": "RFC3339Nano"
}
```

`device.linked`

```json
{
  "event_id": "…", "name": "device.linked", "occurred_at": "RFC3339Nano",
  "profile_id": "…", "device_id": "…", "session_id": "…",
  "device_name": "…", "os": "…", "linked_at": "RFC3339Nano"
}
```

`device.key.added`

```json
{
  "event_id": "…", "name": "device.key.added", "occurred_at": "RFC3339Nano",
  "profile_id": "…", "device_id": "…",
  "key_id": "…", "derived_key_id": "0x…",
  "key_type": "p256-webauthn", "key_type_code": 7,
  "public_key": "base64 of the stored key bytes",
  "key_sha256": "…",
  "extra": { "…": "the stored extra document" },
  "created_at": "RFC3339Nano", "expires_at": "RFC3339Nano (when set)"
}
```

`device.key.removed` carries the same identification plus `removed_at`.

`key_id` is this service's row id; `derived_key_id` is the 20-byte id GFOS
resolves as `device_key_id`.

### What is never published

`public_key` and `extra` are present only for key types that hold public
material (`secp256k1`, `p256-webauthn`, `ed25519`, `curve25519`). For any other
type — an FCM token, a pickle key, a Matrix session key — the fact states that a
key of that type was added or removed and identifies it by `key_sha256`, and the
material itself never leaves the database. Private key material is never held by
this service at all.
