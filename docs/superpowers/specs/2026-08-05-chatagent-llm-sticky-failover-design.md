# Chat Agent LLM Sticky Multi-Key Failover

**Status:** Approved for implementation  
**Date:** 2026-08-05  
**Scope:** `apps/chatagent` inference client only (engine, business, and RPC contracts unchanged)

## Problem

ChatAgent uses a single OpenAI-compatible inference endpoint (`INFERENCE_BASE_URL` + `INFERENCE_API_KEY` + `INFERENCE_MODEL`). Production needs:

1. **Multiple API keys** for the preferred model (quota / rate-limit resilience).
2. **Optional secondary provider** (OpenAI ↔ Google Gemini OpenAI-compat) after the primary key pool is exhausted.
3. **Sticky selection** — use the highest-priority healthy candidate until it degrades; do **not** rotate keys every request for load balancing.
4. **Cooldown recovery** — after a key is marked degraded, prefer the next key; when cooldown expires, prefer the primary again.

## Goals

| Goal | Success criteria |
|------|------------------|
| Sticky primary | Successive successful calls always use the highest-priority non-degraded candidate |
| Same-request failover | If candidate N fails with a degradable error, try N+1 in the same `Complete` call |
| Cooldown | Degraded candidates are skipped until `now + cooldown`; then eligible again |
| Keys-first chain | All primary keys (same model) before any secondary provider keys |
| Provider support | `openai`, `google` (Gemini OpenAI-compat), and `custom` (explicit base URL) |
| Backward compatible | Existing `INFERENCE_BASE_URL` / `INFERENCE_API_KEY` / `INFERENCE_MODEL` still work |
| Safe logging | Never log full API keys; fingerprint last 4 chars only |
| No engine changes | `engine.Completer` interface remains `Complete(ctx, prompt) (string, error)` |

## Non-goals

- Round-robin or weighted load balancing across keys
- Cross-replica shared degraded state (each process keeps its own map)
- Anthropic / native Gemini REST (non–OpenAI-compat) APIs
- Per-tenant or per-session key selection
- Changing extract prompts, readiness rules, or Notification delivery
- Persistent metrics store (optional counters/logs only)

## Decisions (from product discussion)

1. **Keys first, same model** — primary provider+model with ordered keys; secondary provider only after all primary keys are degraded.
2. **Cooldown then try primary again** — not stick-until-restart; not fail-over-for-this-request-only without state.
3. **Primary + optional secondary provider** — not single-provider-only and not an unbounded multi-provider list beyond two slots.

## Architecture

```
engine.Agent
    │ Complete(ctx, prompt)
    ▼
llm.FailoverCompleter          ← sticky selector + same-request retry
    │ picks highest-priority healthy Candidate
    ▼
llm.Completer (per call)       ← existing HTTP OpenAI-compat client
    │ POST {base}/…/chat/completions
    ▼
OpenAI API  |  Gemini OpenAI-compat  |  custom gateway
```

### Candidate chain (priority order)

```
index 0: primary  provider + model + key[0]
index 1: primary  provider + model + key[1]
…
index P: secondary provider + model + key[0]   (if secondary configured)
…
```

Built once at process start from env config. Empty keys are skipped. If no candidates, inference stays disabled (evidence-only mode), same as today when base URL is empty.

### Sticky selection algorithm

State (in-memory, mutex-protected):

```text
degradedUntil[i] time.Time   // zero ⇒ not degraded
cooldown         time.Duration
```

On each `Complete(ctx, prompt)`:

1. Let `start` be the lowest index `i` where `degradedUntil[i].IsZero() || !time.Now().Before(degradedUntil[i])`.
2. If no such index, return error (all degraded and cooldown not yet elapsed — still try all indices once if all are "degraded" but we need a response: prefer trying every candidate whose cooldown has expired; if none, try the soonest-to-recover / all keys anyway as last resort — **policy: if every candidate is still within cooldown, still attempt all of them once in priority order rather than hard-failing without I/O**. Rationale: cooldowns are hints to avoid thundering the primary; total outage should still probe).
3. For each eligible candidate in priority order (skip only those still cooling down **when at least one is healthy**; if none healthy, probe all):
   - Call HTTP completer with that candidate’s base URL, model, key, path style.
   - On **success**: return content. Do **not** mark others healthy/unhealthy. Next request still re-selects from index 0’s health rules (so primary returns after cooldown without “locking” on secondary).
   - On **degradable error**: set `degradedUntil[i] = now + cooldown`, log failover, continue to next eligible.
   - On **hard error** (client/request fault, e.g. HTTP 400/422): return immediately; do **not** try other keys (same bad prompt would fail everywhere and burn quota).
4. If all attempts fail: return a multierror summarizing attempts (status + key fingerprint + provider).

**Sticky vs sticky-active-index:** There is no “current index that advances on success.” Stickiness means “always the best healthy key,” not “stay on whatever worked last.” After key0 fails and key1 succeeds, subsequent calls use key1 only while key0 is degraded; when key0’s cooldown ends, key0 is preferred again. That is intentional and matches “most reliable until it degrades” + “cooldown then try primary again.”

### Error classification

| Condition | Class | Behavior |
|-----------|--------|----------|
| Context canceled / deadline from **caller** | hard (propagate) | Stop; do not mark degraded solely for caller cancel mid-flight if request never started; if HTTP returned cancel because **client timeout** while provider hung, treat as degradable |
| Dial / connection / TLS failure | degradable | Mark key, try next |
| HTTP timeout (client) | degradable | Mark key, try next |
| HTTP 429 | degradable | Mark key, try next |
| HTTP 500, 502, 503, 504 | degradable | Mark key, try next |
| HTTP 401, 403 | degradable | Mark key (invalid/revoked/forbidden for that key), try next |
| HTTP 400, 404, 413, 422 | hard | Fail request; do not try next key |
| Other 4xx | hard | Fail request |
| Empty choices / response decode error after 200 | degradable | Provider glitch; try next |
| Empty base URL / misconfiguration at call site | hard | Fail |

Caller `ctx` cancellation that occurs before or during the call: if `errors.Is(err, context.Canceled)` and `ctx.Err() != nil`, treat as **hard** (do not degrade and do not burn other keys for a cancelled request). If the transport times out independently, degradable.

### Provider URL resolution

| Provider | Default base URL | Completions path |
|----------|------------------|------------------|
| `openai` | `https://api.openai.com` | `/v1/chat/completions` |
| `google` | `https://generativelanguage.googleapis.com/v1beta/openai` | `/chat/completions` |
| `custom` | **required** `INFERENCE_BASE_URL` | `/v1/chat/completions` (default) |

Rules:

- Explicit `INFERENCE_BASE_URL` / `INFERENCE_SECONDARY_BASE_URL` **overrides** the default base for that slot.
- Paths are joined with careful slash normalization (no double slashes; no forced `/v1` on Google).
- Google uses the [OpenAI compatibility surface](https://ai.google.dev/gemini-api/docs/openai) so request/response JSON stays the same as today’s client (`messages`, `response_format: json_object`, Bearer auth).

### Auth header

All providers use `Authorization: Bearer <api_key>` on the OpenAI-compat surface (including Google’s OpenAI-compat endpoint).

### Config surface

```go
// Primary
INFERENCE_PROVIDER              // openai | google | custom (empty ⇒ custom if BASE_URL set, else disabled)
INFERENCE_BASE_URL              // optional override / required for custom
INFERENCE_MODEL                 // default meta/llama-3.1-8b-instruct (existing)
INFERENCE_API_KEY               // legacy single key
INFERENCE_API_KEYS              // "key1,key2,key3" or whitespace-separated; preferred over single key

// Secondary (optional)
INFERENCE_SECONDARY_PROVIDER
INFERENCE_SECONDARY_BASE_URL
INFERENCE_SECONDARY_MODEL
INFERENCE_SECONDARY_API_KEYS    // or INFERENCE_SECONDARY_API_KEY

// Failover policy
INFERENCE_FAILOVER_COOLDOWN     // Go duration, default "2m"
```

**Key parsing:** split on commas and/or whitespace; trim; drop empties; preserve order; de-dupe consecutive duplicates optional (de-dupe exact duplicates globally to avoid double-billing the same key — **yes, unique while preserving first-seen order**).

**Enablement:** inference is enabled when at least one candidate can be built (provider defaults or base URL + ≥1 key). Model may have a default. If only model is set with no keys and no base for custom, stay disabled.

**Legacy:**  
`INFERENCE_BASE_URL` + `INFERENCE_API_KEY` + model → one `custom` (or inferred) candidate — identical behavior to today, with failover wrapper of length 1.

### Package layout

```
apps/chatagent/service/llm/
  openai.go       # single-shot Completer (path-aware)
  candidate.go    # Candidate, Provider, ResolveBaseURL, ParseKeys
  errors.go       # Classify(err) → degradable | hard
  failover.go     # FailoverCompleter
  build.go        # BuildCompleter(cfg, httpClient) engine.Completer
  *_test.go
```

Handlers:

```go
completer := llm.BuildFromConfig(llm.Config{...}, httpClient)
```

instead of `llm.New(...)`.

### Concurrency

`FailoverCompleter` is safe for concurrent `Complete` calls (shared by the process). Mutex protects `degradedUntil` map/slice only; HTTP calls run outside the lock (copy candidate under RLock, update degrade under Lock after failure).

### Observability

On degradable failover:

```
level=WARN msg="chatagent llm: candidate degraded; trying next"
  provider=openai model=gpt-4o key=…abcd status=429 attempt=1 next=2
```

On final failure after all candidates:

```
level=ERROR msg="chatagent llm: all candidates failed"
  attempts=…
```

Success path: Debug only (provider, model, key fingerprint) to avoid noise.

### Testing strategy

Unit tests with `httptest.Server` (or `RoundTripper` fakes):

1. Sticky success — two calls, only first key’s transport used.
2. Same-request failover — key0 returns 429, key1 200; response from key1; key0 degraded.
3. Post-failover stickiness — next call only hits key1 while key0 cooling down.
4. Cooldown recovery — advance clock; next call hits key0 again.
5. Hard 400 — second key never called.
6. Caller cancel — no degrade of unused keys; error is cancel.
7. All keys fail — error contains both fingerprints/statuses.
8. Config — `ParseKeys`, primary-then-secondary order, legacy single key, Google path not `/v1/chat/completions`.
9. Race — concurrent Complete with `-race`.

Clock: inject `func() time.Time` (or `clock` interface) on FailoverCompleter for deterministic cooldown tests.

### Rollout / ops

Example production:

```bash
INFERENCE_PROVIDER=openai
INFERENCE_MODEL=gpt-4o-mini
INFERENCE_API_KEYS=sk-proj-primary,sk-proj-backup
INFERENCE_SECONDARY_PROVIDER=google
INFERENCE_SECONDARY_MODEL=gemini-2.0-flash
INFERENCE_SECONDARY_API_KEYS=AIza...backup
INFERENCE_FAILOVER_COOLDOWN=2m
```

Deploy notes:

- Each replica has independent degrade state (expected).
- Rotating keys: update secrets and restart (or new keys appear only after redeploy).
- Invalid primary key: 401 degrades it for cooldown; secondary keeps service available.

### Risks and mitigations

| Risk | Mitigation |
|------|------------|
| Wrong Google path (`/v1/chat/completions`) | Provider-specific path table + unit test |
| Burning all keys on bad prompt | Hard errors do not fan-out |
| Thundering herd after cooldown | Cooldown is per-key; optional future jitter not required for v1 |
| Logging secrets | Fingerprint only (`…` + last 4 runes of key) |
| JSON mode unsupported on some models | Pre-existing risk; hard 400 fails closed without multi-key burn |

## Implementation outline

1. Refactor `Completer` to use `CompletionsPath` + shared request body.
2. Add `Candidate`, key parsing, provider defaults.
3. Add error classification + `FailoverCompleter` with injectable clock.
4. Add `BuildFromConfig` + expand `ChatAgentConfig` / `LLMConfig` / `main` wiring.
5. Tests + README config table.
6. Design/plan docs under `docs/superpowers/`.

## Acceptance checklist

- [ ] Multiple keys: primary used until degradable failure, then next key same request
- [ ] No per-request rotation when primary is healthy
- [ ] Cooldown elapses → primary preferred again
- [ ] Secondary provider only after all primary keys degraded (or cooling)
- [ ] Legacy single-key env still works
- [ ] Google completions path correct
- [ ] `go test ./apps/chatagent/...` passes with `-race` on llm package
- [ ] README documents new env vars
