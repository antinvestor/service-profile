# Chat Agent (`apps/chatagent`)

Product-agnostic **conversational data collection** tool. Products register a
**context** (required fields + purpose/prompt). The agent evaluates evidence
already available (seed fields, CV/documents, prior conversation) and only
asks for what is still missing.

## Why it exists

Rigid forms force users to re-enter data the system already has. Chat agent:

1. Accepts a **context definition** (schema + guidance) — the only product-specific input
2. Opens a **session** with seed fields + documents (e.g. uploaded CV text)
3. Runs **turns** that extract from all evidence before asking the next gap
4. Returns `ready` + structured `fields` when required data is complete
5. Optionally delivers assistant replies via **Notification service** so the same
   conversation works on web, SMS, WhatsApp, email, push, in-app, or USSD

Different products (placement, KYC, support intake, …) reuse the same service
and only change the context (and channel binding for omnichannel).

## Architecture (channel-agnostic core)

```
                    ┌─────────────────────────┐
  Web / product ──► │ ChatAgentService.Turn   │
  SMS/WA adapter ─► │ IngestChannelMessage    │──► evidence-first engine
                    └───────────┬─────────────┘
                                │ assistant reply
                    ┌───────────▼─────────────┐
                    │ notify.Deliverer        │  (no-op when channel=web)
                    │ → Notification.Send     │──► sms | email | push | wa | …
                    └─────────────────────────┘
```

- **Engine** never knows about SMS/WhatsApp — only text + fields.
- **ChannelBinding** on the session is the delivery address book entry.
- **Notification** owns routing, templates, and channel integrations.

## API (Connect)

| RPC | Purpose |
|-----|---------|
| `UpsertContext` | Register/version a context definition |
| `GetContext` / `ListContexts` | Read registry |
| `CreateSession` | Start session; optional `evaluate_evidence` + `channel` binding |
| `GetSession` | Resume transcript + field state |
| `Turn` | One message (and/or new documents/structured inputs); delivers reply when channel bound |
| `IngestChannelMessage` | Inbound SMS/WhatsApp/email path: resolve/create session → Turn → deliver reply |
| `EndSession` | Close session |

Audience path (catalog): `/chat-agent` (`servicecatalog.ServiceChatAgent`).

## Omnichannel usage

### 1. Web (default)

Omit `channel` or set `CHANNEL_WEB`. Replies return only on the RPC (`Turn.reply`).

### 2. Outbound-capable session (SMS / WhatsApp / …)

```json
// CreateSession
{
  "subject_id": "profile-abc",
  "context_key": "stawi.placement.intake",
  "evaluate_evidence": true,
  "channel": {
    "channel": "CHANNEL_SMS",
    "contact_id": "contact-phone-id",
    "profile_id": "profile-abc",
    "language": "en"
  }
}
```

Every subsequent `Turn` (and the optional evaluate reply) is also queued on
Notification with `type=sms`, `out_bound=true`, `auto_release=true`, and
`data` = assistant text (or a template when `channel.template` is set).

### 3. Inbound channel adapter

Channel integrations (or a thin adapter on Notification inbound routes) call:

```json
// IngestChannelMessage
{
  "subject_id": "profile-abc",
  "context_key": "stawi.placement.intake",
  "create_if_missing": true,
  "message": "Backend engineer, full-time, Kenya",
  "channel": {
    "channel": "CHANNEL_WHATSAPP",
    "contact_id": "contact-wa-id",
    "profile_id": "profile-abc"
  }
}
```

This resolves the active session for that contact/channel (or creates one),
runs a turn, and delivers the reply on the same channel.

### Delivery rules

| Condition | Behavior |
|-----------|----------|
| `channel` empty / `WEB` / `UNSPECIFIED` | RPC only |
| `skip_delivery=true` | RPC only even if SMS/WA bound |
| Non-web + contact or profile | `Notification.Send` after reply |
| Delivery error | Logged; **Turn still succeeds** (soft-fail) |
| `NOTIFICATION_SERVICE_URI` empty | Delivery skipped; RPC replies still work |

Supported channel enums: `WEB`, `SMS`, `EMAIL`, `PUSH`, `IN_APP`, `WHATSAPP`, `USSD`.

## Config

| Env | Purpose |
|-----|---------|
| `DATABASE_PRIMARY_URL` | PostgreSQL |
| `DATABASE_MIGRATION_PATH` | e.g. `file://migrations` |
| `INFERENCE_BASE_URL` | OpenAI-compatible root (optional) |
| `INFERENCE_API_KEY` | Bearer token |
| `INFERENCE_MODEL` | Default `meta/llama-3.1-8b-instruct` |
| `NOTIFICATION_SERVICE_URI` | Notification Connect endpoint (optional; enables omnichannel) |
| `NOTIFICATION_SERVICE_WORKLOAD_API_TARGET_PATH` | SPIFFE path for S2S to notification (default `/ns/notifications/sa/service-notification`) |
| OIDC / Frame security | Same as other profile apps |

Without inference, the agent still merges seed/documents and returns guided
follow-ups (evidence-only mode). Without Notification URI, channel bindings are
stored but replies are not delivered off-web.

## Local generate protos

```bash
cd proto
buf generate --template buf.gen.chatagent.yaml chatagent
# Ensure gen imports use buf.build packages for common + gnostic (see gen/go/chatagent).
```

Generated Go lives under `gen/go/chatagent/v1` (local until BSR module is published).

## Example: placement-style context

```json
{
  "context_key": "stawi.placement.intake",
  "purpose": "Collect qualifications and preferences for opportunity matching.",
  "fields": [
    {"name": "target_job_title", "type": "FIELD_TYPE_STRING", "required": true, "priority": 1, "ask": "What role are you targeting?"},
    {"name": "capabilities", "type": "FIELD_TYPE_STRING", "required": true, "priority": 2, "min_length": 80, "evidence_hints": ["document"], "ask": "Share your CV or work history."},
    {"name": "job_types", "type": "FIELD_TYPE_STRING_LIST", "required": true, "priority": 3, "enum_values": ["Full-time", "Part-time", "Contract"]},
    {"name": "salary_min", "type": "FIELD_TYPE_NUMBER", "required": true, "priority": 4},
    {"name": "preferred_countries", "type": "FIELD_TYPE_STRING_LIST", "required": true, "priority": 5},
    {"name": "experience_level", "type": "FIELD_TYPE_STRING", "required": true, "priority": 6, "enum_values": ["entry","junior","mid","senior","lead","executive"]}
  ]
}
```

CreateSession with CV document + `evaluate_evidence=true` will mark fields
satisfied by the CV before the user types anything.

## Layout

```
apps/chatagent/
  cmd/                 entrypoint (+ Notification client wiring)
  config/              env config
  service/
    engine/            pure evidence-first turn engine (unit tested)
    business/          session/context orchestration + channel ingest
    handlers/          Connect RPC
    notify/            Notification.Send deliverer (omnichannel)
    repository/        PostgreSQL
    models/
    llm/               OpenAI-compatible client
    authz/
  migrations/
```
