# Chat Agent LLM Sticky Failover Implementation Plan

> **For agentic workers:** Implement task-by-task. Steps use checkbox syntax for tracking.

**Goal:** Multi-key sticky failover for OpenAI and Google (OpenAI-compat) models in chatagent, with cooldown recovery and no per-request rotation.

**Architecture:** Ordered candidate list (primary keys then secondary provider keys) behind `FailoverCompleter` implementing `engine.Completer`. In-memory per-key degrade cooldowns; always select highest-priority healthy candidate.

**Tech Stack:** Go, existing `llm.Completer` HTTP client, Frame HTTP client manager, envconfig.

---

### Task 1: Core llm types + Completer path fix
### Task 2: FailoverCompleter + classification
### Task 3: Config build + wiring
### Task 4: Tests
### Task 5: README

Files:

- Create: `apps/chatagent/service/llm/candidate.go`
- Create: `apps/chatagent/service/llm/errors.go`
- Create: `apps/chatagent/service/llm/failover.go`
- Create: `apps/chatagent/service/llm/build.go`
- Create: `apps/chatagent/service/llm/*_test.go`
- Modify: `apps/chatagent/service/llm/openai.go`
- Modify: `apps/chatagent/config/config.go`
- Modify: `apps/chatagent/service/handlers/chatagent.go`
- Modify: `apps/chatagent/cmd/main.go`
- Modify: `apps/chatagent/README.md`
