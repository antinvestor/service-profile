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
5. Optionally delivers assistant replies via the **existing Notification service**
   (`NotificationService.Send`) — ChatAgent does **not** invent channels

Different products reuse the same service and only change the context (and
optional `NotificationTarget` when they want non-web delivery).

## Omnichannel = reuse Notification service

ChatAgent stays **channel-agnostic**. Channels (SMS, email, push, in-app, …)
already live in **service-notification** (routes, templates, integrations).

```
  Web / product RPC ──► Turn / CreateSession
  Notification adapter ─► IngestMessage
                │
                ▼ evidence-first engine (text + fields only)
                │
                ▼ assistant reply
  notificationv1connect.NotificationServiceClient.Send
                │  (same client/pattern as profile verification)
                ▼
         Notification service (type + recipient + template/data)
                │
                ▼ sms | email | push | …  (routes in Notification)
```

### How delivery is configured

`NotificationTarget` is a **thin subset of `notification.v1.Notification`**:

| Field | Same meaning as |
|-------|-----------------|
| `type` | `Notification.type` (`"sms"`, `"email"`, `"push"`, `"in-app"`, …) |
| `recipient` | `Notification.recipient` (`common.v1.ContactLink`) |
| `source` | `Notification.source` |
| `language` | `Notification.language` |
| `template` | `Notification.template` (empty → raw body in `Notification.data`) |
| `payload` | Merged into `Notification.payload` (`reply` + `session_id` always set) |
| `route_id` | `Notification.route_id` |
| `skip` | Opt out of Send even if type is set |

Empty `type` (or no target) → **RPC-only** (web). No parallel channel enum.

### Client wiring (same as profile app)

```go
// apps/chatagent/cmd — identical pattern to apps/default setupNotificationClient
notificationCli := setupNotificationClient(ctx, cfg) // connection.NewServiceClient → NotificationServiceClient

handlers.NewChatAgentServer(ctx, svc, handlers.ServerDeps{
    NotificationClient: notificationCli,
})
```

Outbound Send follows the same shape as contact verification:

```go
n := &notificationv1.Notification{
    Recipient:   target.Recipient,
    Type:        target.Type,
    Template:    target.Template,
    Payload:     payload, // includes reply, session_id
    Data:        reply,   // when template empty
    Language:    target.Language,
    OutBound:    true,
    AutoRelease: true,
}
cli.Send(ctx, connect.NewRequest(&notificationv1.SendRequest{Data: []*notificationv1.Notification{n}}))
```

## API (Connect)

| RPC | Purpose |
|-----|---------|
| `UpsertContext` | Register/version a context definition |
| `GetContext` / `ListContexts` | Read registry |
| `CreateSession` | Start session; optional `evaluate_evidence` + `notification` target |
| `GetSession` | Resume transcript + field state |
| `Turn` | One message; Send via Notification when session has a target |
| `IngestMessage` | Inbound message (from Notification adapters) → Turn → Notification.Send |
| `EndSession` | Close session |

Audience path (catalog): `/chat-agent` (`servicecatalog.ServiceChatAgent`).

## Examples

### Web (default)

Omit `notification`. Replies return only on the RPC (`Turn.reply`).

### SMS via Notification service

```json
// CreateSession
{
  "subject_id": "profile-abc",
  "context_key": "stawi.placement.intake",
  "evaluate_evidence": true,
  "notification": {
    "type": "sms",
    "recipient": {
      "contact_id": "contact-phone-id",
      "profile_id": "profile-abc",
      "profile_type": "Profile"
    },
    "language": "en"
  }
}
```

Every subsequent `Turn` (and the optional evaluate reply) is also queued with
`NotificationService.Send` (`type=sms`, `out_bound=true`, `auto_release=true`).

### Inbound adapter (Notification integration path)

```json
// IngestMessage
{
  "subject_id": "profile-abc",
  "context_key": "stawi.placement.intake",
  "create_if_missing": true,
  "message": "Backend engineer, full-time, Kenya",
  "notification": {
    "type": "whatsapp",
    "recipient": {
      "contact_id": "contact-wa-id",
      "profile_id": "profile-abc"
    }
  }
}
```

Channel **routing** (which gateway, which route) stays in Notification —
ChatAgent only supplies the same `type` / `recipient` / `template` fields
Notification already uses.

### Delivery rules

| Condition | Behavior |
|-----------|----------|
| No `notification` / empty `type` / `type=web` | RPC only |
| `skip=true` | RPC only |
| Type set + recipient contact or profile | `NotificationService.Send` after reply |
| Send error | Logged; **Turn still succeeds** (soft-fail) |
| `NOTIFICATION_SERVICE_URI` empty / client nil | Send skipped; RPC replies still work |

## Config

| Env | Purpose |
|-----|---------|
| `DATABASE_PRIMARY_URL` | PostgreSQL |
| `DATABASE_MIGRATION_PATH` | e.g. `file://migrations` |
| `NOTIFICATION_SERVICE_URI` | Existing Notification Connect endpoint (optional) |
| `NOTIFICATION_SERVICE_WORKLOAD_API_TARGET_PATH` | SPIFFE path for S2S to notification |
| OIDC / Frame security | Same as other profile apps |

### Inference (sticky multi-key failover)

Inference is optional. When no provider/keys/base URL are set, the agent runs in
**evidence-only** mode (guided replies without LLM extract).

| Env | Purpose |
|-----|---------|
| `INFERENCE_PROVIDER` | `openai`, `google`, or `custom` |
| `INFERENCE_BASE_URL` | Override root URL (required for `custom`; optional override for openai/google) |
| `INFERENCE_MODEL` | Primary model (default `meta/llama-3.1-8b-instruct`) |
| `INFERENCE_API_KEY` | Single bearer token (legacy) |
| `INFERENCE_API_KEYS` | Ordered keys (`key1,key2` or whitespace). Preferred over single key |
| `INFERENCE_SECONDARY_PROVIDER` | Optional fallback provider after primary keys degrade |
| `INFERENCE_SECONDARY_BASE_URL` | Override for secondary |
| `INFERENCE_SECONDARY_MODEL` | Secondary model (defaults to primary model if empty) |
| `INFERENCE_SECONDARY_API_KEY` / `INFERENCE_SECONDARY_API_KEYS` | Secondary key pool |
| `INFERENCE_FAILOVER_COOLDOWN` | How long a failed key is skipped (default `2m`) |

**Behavior (not load balancing):**

1. Always use the **highest-priority healthy** candidate (primary key order, then secondary).
2. On degradable errors (429, 5xx, 401/403, timeouts, network), mark that key degraded for the cooldown and try the next key **in the same request**.
3. Do **not** rotate keys every request while the primary is healthy.
4. When the cooldown expires, the primary key is preferred again.
5. Hard client errors (400/422) fail the request without burning other keys.

**Provider defaults:**

| Provider | Default base | Completions path |
|----------|--------------|------------------|
| `openai` | `https://api.openai.com` | `/v1/chat/completions` |
| `google` | `https://generativelanguage.googleapis.com/v1beta/openai` | `/chat/completions` |
| `custom` | (required via `INFERENCE_BASE_URL`) | `/v1/chat/completions` |

**Example — OpenAI primary keys + Gemini secondary:**

```bash
INFERENCE_PROVIDER=openai
INFERENCE_MODEL=gpt-4o-mini
INFERENCE_API_KEYS=sk-proj-primary,sk-proj-backup
INFERENCE_SECONDARY_PROVIDER=google
INFERENCE_SECONDARY_MODEL=gemini-2.0-flash
INFERENCE_SECONDARY_API_KEYS=AIzaSy…primary,AIzaSy…backup
INFERENCE_FAILOVER_COOLDOWN=2m
```

**Legacy (unchanged):**

```bash
INFERENCE_BASE_URL=https://my-gateway.example
INFERENCE_API_KEY=secret
INFERENCE_MODEL=meta/llama-3.1-8b-instruct
```

## Local generate protos

```bash
cd proto
buf generate --template buf.gen.chatagent.yaml chatagent
# Ensure gen imports use buf.build packages for common + gnostic (see gen/go/chatagent).
```

Generated Go lives under `gen/go/chatagent/v1`.

## Layout

```
apps/chatagent/
  cmd/                 entrypoint (Notification client via connection.NewServiceClient)
  config/              env config
  service/
    engine/            pure evidence-first turn engine (unit tested)
    business/          session orchestration + Notification.Send (no channel model)
    handlers/          Connect RPC
    repository/        PostgreSQL
    models/
    llm/               OpenAI-compatible client + sticky multi-key failover
    authz/
  migrations/
```
