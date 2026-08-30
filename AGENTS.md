# AGENTS.md

Guidance for AI agents and humans working in this repository. Read the
[master plan](docs/enterprise_waf_development_master_plan.md) and
[Phase 0 SRS](docs/phase0_srs.md) before changing architecture or contracts.

## Ground rules

- **Contracts are versioned:** policy/event JSON Schema in `packages/*/schema` is
  the source of truth. SDK types are generated, never hand-edited.
- **Architecture changes need an ADR:** anything touching a trust boundary, data
  flow, store, or protocol requires a new ADR in `docs/architecture/` using the
  template in `docs/architecture/repository-ci-coding-standards.md`.
- **No secrets:** no default passwords, no keys in code, no bodies/cookies/
  credentials in logs. Never commit `.env` or certificate material.
- **Traffic path invariants (master plan §22):** never block/crash the live
  traffic path on a log write; never deploy partially compiled policy; never
  add a protocol without conformance/fuzz/evasion/resource-limit tests.
- **Do not claim production or F5-class capability** in code, docs, or commits
  without evidence from the QA program.

## Commands

- `make lint` / `make test` / `make build` / `make gen` — see Makefile.
- `make gen-check` / `make test-replay` / `make test-failover` / `make schema-check` — additional CI targets.
- `make check-i18n` — verify en/bn message catalogs match.
- Always run lint + tests before finishing a task.
- **Commit tasks:** Every finished task must be committed to git with a clear, descriptive commit message relating to the task. This ensures any AI or human can understand the repository history and context.

## Docs to keep updated

- `docs/phase0_srs.md` — the day-to-day contract (FR/NFR IDs, use cases).
- `docs/api/endpoint.md` — The complete API Master Plan and Architecture roadmap (Phases 1-5). **Must be checked before adding any new API endpoint.**
- `docs/api/API Specification Document.md` — Detailed JSON schemas and payloads for the enterprise API.
- `docs/api/openapi-v0.yaml` — OpenAPI spec; must be kept in sync with `apps/control-api/internal/api/router.go`.
- `docs/architecture/` — ADRs + diagrams (update diagrams with any boundary/data-flow change).
- `docs/AI_SKILL.md` — AI-oriented project summary; keep in sync with the source docs (it must not override them).

