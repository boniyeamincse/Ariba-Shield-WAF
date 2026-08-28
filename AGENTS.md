# AGENTS.md

Guidance for AI agents and humans working in this repository. Read the
[master plan](enterprise_waf_development_master_plan.md) and
[Phase 0 SRS](phase0_srs.md) before changing architecture or contracts.

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
- Always run lint + tests before finishing a task.
- **Commit tasks:** Every finished task must be committed to git with a clear, descriptive commit message relating to the task. This ensures any AI or human can understand the repository history and context.

## Docs to keep updated

- `phase0_srs.md` — the day-to-day contract (FR/NFR IDs, use cases).
- `docs/architecture/` — ADRs + diagrams (update diagrams with any boundary/data-flow change).
