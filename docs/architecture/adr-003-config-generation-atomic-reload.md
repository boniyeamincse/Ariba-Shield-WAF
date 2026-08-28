# ADR-003 — Configuration Generation and Atomic Reload Mechanism

- **Status:** Proposed (for review)
- **Date:** 2026-08-28
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§11, §22 rules 3 & 9, §15 Phase 1 exit criteria), `../phase0_srs.md` (FR-0.1-051..056, UC-3, UC-4), `adr-002-policy-event-schema-v0.md`

---

## Context

Release 0.1 must satisfy three hard requirements:

1. **Invalid config cannot break the active listener** (Phase 1 exit criterion).
2. **Last-known-good config restores automatically** after restart or reload failure.
3. **No partial policy activation** — validate, stage, atomically switch (§11).

OpenResty/Nginx loads config at process start; naive reload (`nginx -s reload`) is not safe against semantic errors, and a half-written file can take down the gateway. This ADR defines the mechanism that makes configuration safe on the data plane.

---

## Decision

Adopt a **validate → stage → switch** pipeline with a **write-ahead config store** and a **versioned active pointer**, using Nginx's process signaling only after validation, and falling back to the last-known-good on any failure.

### D1 — Config generation

- The control plane's **policy compiler** (Go) renders the declarative policy document (ADR-002) into a gateway-native configuration (nginx `http`/`server`/`upstream` blocks for Release 0.1).
- Generation is **deterministic**: identical policy in → identical bytes out. Golden tests lock this.
- Compiler output is validated in two layers:
  1. **Schema/semantic validation** (control plane): duplicate listeners, unknown refs, port conflicts, empty pools, path conflicts — rejected before anything reaches the gateway.
  2. **Gateway-side validation** (`nginx -t` / `openresty -t`): syntax and directive validity on the target image.

### D2 — Write-ahead config store (on gateway)

The gateway keeps a small local store of **immutable config versions**:

```text
/var/lib/shield-waf/config/
  <bundle_hash>/
    config.json        # policy document (signed)
    generated.conf     # compiled nginx config
    meta.json          # applied_at, applied_by, status, error
  active                 # file containing the bundle_hash of the ACTIVE config
  last-known-good        # file containing the bundle_hash of the last fully-valid active config
```

Rules:
- Versions are **append-only**; a version never overwrites another.
- `active` is updated **only after** a successful switch + verify.
- `last-known-good` is updated **only after** the active config has served traffic successfully for a short soak period (default 60 s; configurable) with the gateway reporting healthy.

### D3 — Atomic activation

On receipt of a valid, correctly signed/hashed bundle:

1. Write new version directory (`config.json`, `generated.conf`) with `status: staged`.
2. Run `nginx -t` against the staged file with a temporary config path.
3. If validation fails → mark `status: rejected`, record `error`, **never touch `active`**, return failure to control plane (audit event).
4. If validation passes → atomically update `active` pointer (single `rename()`), then send `SIGHUP` to the running worker for a **graceful reload**.
5. Wait for reload completion + verify the new config is serving (health probe). On success mark `status: active` and update `last-known-good`.
6. On any failure after switch → immediately re-point `active` to the last-known-good bundle hash and reload again (rollback path, §D4).

### D4 — Rollback and restart recovery

- **Reload failure:** if the reload signal fails or the gateway fails to come up healthy, the wrapper reverts `active` to `last-known-good` and reloads. Logs and audit events document the rollback with the rejected hash.
- **Process restart:** on gateway start, the bootstrap reads `active`; if the active config fails to validate/start, it falls back to `last-known-good` (which by construction always validated). If even last-known-good is absent, the gateway starts in a **safe-mode** that serves a controlled maintenance/503 response rather than crashing or serving with no config.
- **Ordering guarantee:** the write-ahead store is on local disk; `active`/`last-known-good` pointer updates are atomic `rename()` operations, so a crash mid-write can never leave a torn pointer.

### D5 — Single parse path invariant (master plan rule 1)

Release 0.1 performs **no security parsing**; when Phase 2 adds the Coraza path, the *same normalized request representation* produced by the gateway's parse/normalization layer must be the only representation inspected by both security and proxy logic. This ADR records the invariant now so the config generation and reload mechanism never introduces a second, divergent parse path.

---

## Alternatives considered

### A1 — Direct edit + `nginx -s reload` on live config file
**Rejected.** Editing the live file non-atomically can leave a torn config; `reload` does not guard semantic validity; a bad reload can take listeners down. Violates Phase 1 exit criteria.

### A2 — Control plane pushes config and gateway applies blindly
**Rejected.** Gateway has no local authority to self-heal; control-plane outage or a bad push would require manual intervention. Violates last-known-good and §12.3 failure matrix (control plane unreachable).

### A3 — Config in a KV store read at request time
**Rejected.** Adds latency to the traffic path and a runtime dependency; conflicts with "traffic decisions must never depend on a log write" and deterministic offline operation. Also incompatible with Nginx's static-config model for 0.1.

### A4 — Just store a single config file, no versioning
**Rejected.** Breaks rollback and canary (Phase 3) and makes restart recovery ambiguous.

---

## Consequences

### Positive

- Invalid config can never break the active listener (Release 0.1 exit criterion met).
- Restart recovery is deterministic and automatic via last-known-good.
- Rollback is cheap and versioned; same mechanism scales to Phase 3 canary/rollback.
- Immutable versioned store gives an audit trail of every applied/rejected config.
- Single parse-path invariant is recorded before Phase 2, preventing divergence.

### Negative

- Disk usage for version history on the gateway (bounded by retention/compaction policy).
- Reload orchestration complexity (wrapper process managing signal + health verification).
- `nginx -t` validation only catches directive-level errors; semantic errors must be caught at the control-plane compiler layer (requires strong compiler tests).

### Risks and mitigations

| Risk | Mitigation |
|---|---|
| Torn pointer on crash | Atomic `rename()` for pointer files |
| Reload succeeds but config is wrong semantically | Soak period before promoting to last-known-good + gateway health probes |
| `nginx -t` passes but generated config has runtime pitfalls | Golden generation tests + a reference test battery that runs generated configs under load (Sprint 4) |
| Version store grows unbounded | Retention policy (keep last N active + all active-origin ancestors; configurable) |

---

## Acceptance for this ADR

- Restart test: gateway restarts and resumes with active config (or falls back to last-known-good) without traffic loss on the test app.
- Invalid-config test: a config that fails `nginx -t` is rejected; active listener unaffected (UC-3).
- Reload-failure test: induced reload failure reverts to last-known-good (UC-4).
- No torn-pointer state survives `kill -9` at arbitrary points (chaos test in Sprint 4).
- Golden compiler tests: same policy → same bytes; hash matches ADR-002 canonical serialization.

---

## Follow-up ADRs expected

- **ADR-004:** Gateway registration and heartbeat protocol (reports applied hash, status, validation errors).
- **ADR-005:** Coraza/ModSecurity integration (extends generated config with the WAF path).
- **ADR-006:** Logging and event pipeline design (non-blocking).
