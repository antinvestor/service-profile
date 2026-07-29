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

Different products (placement, KYC, support intake, …) reuse the same service
and only change the context.

## API (Connect)

| RPC | Purpose |
|-----|---------|
| `UpsertContext` | Register/version a context definition |
| `GetContext` / `ListContexts` | Read registry |
| `CreateSession` | Start session; optional `evaluate_evidence` to fill from CV/seed immediately |
| `GetSession` | Resume transcript + field state |
| `Turn` | One message (and/or new documents/structured inputs) |
| `EndSession` | Close session |

Audience path (catalog): `/chat-agent` (`servicecatalog.ServiceChatAgent`).

## Config

| Env | Purpose |
|-----|---------|
| `DATABASE_PRIMARY_URL` | PostgreSQL |
| `DATABASE_MIGRATION_PATH` | e.g. `file://migrations` |
| `INFERENCE_BASE_URL` | OpenAI-compatible root (optional) |
| `INFERENCE_API_KEY` | Bearer token |
| `INFERENCE_MODEL` | Default `meta/llama-3.1-8b-instruct` |
| OIDC / Frame security | Same as other profile apps |

Without inference, the agent still merges seed/documents and returns guided
follow-ups (evidence-only mode).

## Local generate protos

```bash
cd proto
buf generate --template buf.gen.chatagent.yaml chatagent
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
  cmd/                 entrypoint
  config/              env config
  service/
    engine/            pure evidence-first turn engine (unit tested)
    business/          session/context orchestration
    handlers/          Connect RPC
    repository/        PostgreSQL
    models/
    llm/               OpenAI-compatible client
    authz/
  migrations/
```
