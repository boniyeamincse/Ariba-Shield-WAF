# ADR-005 — Coraza WAF Engine Integration

- **Status:** Accepted
- **Date:** 2026-08-29
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§4, §6.2, §6.3, §23), `../phase0_srs.md` (Phase 2 deliverables), `adr-001-initial-architecture.md`, `adr-002-policy-event-schema-v0.md`, `adr-003-config-generation-atomic-reload.md`

---

## Context

Phase 2 introduces the WAF inspection engine into the data plane. The master plan (§4) requires a proven rule engine first — "Coraza/ModSecurity-compatible engine + OWASP CRS" — with a custom compiled engine deferred until semantics and baselines are stable.

The engine must:
1. Sit between the gateway and the backend without becoming a second parser (rule 1).
2. Evaluate request headers, args, and bodies against rules.
3. Support **transparent (detection-only)** and **blocking** modes.
4. **Never fail closed on engine error** — the traffic path must not be blocked by an engine fault (rule 2 analog).
5. Emit structured security events (ADR-002) with masked payloads (rule 5).
6. Support anomaly scoring and a configurable blocking threshold (Phase 3).

---

## Decision

Integrate **Coraza v3** (Go) as a standalone sidecar (`services/waf-engine`) that receives requests from the OpenResty gateway, inspects them, and forwards to the backend. Rules live in `rules/core/` as a **minimal CRS-style baseline**; full OWASP CRS is loaded at Phase 2+ when corpus baselines stabilize.

### D1 — Deployment topology

```text
Client → OpenResty gateway → waf-engine (Coraza) → backend
                              └─ security events → stdout → event-ingestor
```

- The engine is a **separate Go process** (not in-process Lua/nginx), giving it memory safety, its own failure domain, and a clear boundary for the future Rust rewrite.
- `httputil.ReverseProxy` forwards the request after inspection; the body is buffered (bounded, `io.LimitReader` cap) so Coraza can inspect it and the proxy can replay it.

### D2 — Transaction lifecycle

1. Resolve client IP (trusted-proxy aware when behind the gateway).
2. IP reputation pre-check (block-list, allow-list bypass).
3. Rate limit (per-IP sliding window).
4. Create Coraza transaction; seed connection, URI, headers.
5. `ProcessRequestHeaders` (phase 1) → `ProcessRequestBody`/`ReadRequestBodyFrom` (phase 2).
6. On interruption: log event; in **detect-only** forward to backend (transparent), in **blocking** serve the block page.
7. `proxy.ServeHTTP` only when no block was served.

### D3 — Failure semantics

- Any engine/parse/read error logs a security event and **proceeds to proxy the request** (fail open). Blocking decisions are made only by deterministic rule evaluation, never by an error path.
- Body exceeding the limit → 413 in blocking mode; in detect-only, log and forward.
- Coraza rule-load failure at startup is fatal (do not start a gateway without a valid rule set) — this is a **boot-time** gate, not a per-request failure.

### D4 — Events and masking

- On rule match, emit a `security` event (ADR-002) with `rule_ids`, `match_details`, severity, and reason.
- Sensitive args (password/token/api_key/authorization/session/cookie) are **masked** in match data before emission (rule 5; `mask.go`).
- `virtual_server_id` / `application_id` are stamped from trusted gateway headers (`X-Shield-VS-ID` / `X-Shield-App-ID`, set by the compiler template).

### D5 — Rules and anomaly scoring

- Baseline rules in `rules/core/baseline.conf`: SQLi, XSS, command injection, path traversal, LFI, sensitive-file read — each `pass` + `setvar:tx.anomaly_score=+1` (CRS-style scoring).
- A single phase-2 gate rule (`949110`) blocks when `tx.anomaly_score >= tx.blocking_anomaly_score`.
- Threshold is configurable (`--anomaly-threshold`); default baseline sets it to 1. Detect-only ignores the gate (logs only).
- Full OWASP CRS replaces `baseline.conf` when replay corpus + FP baseline pass (Phase 2 exit criteria).

---

## Alternatives considered

### A1 — Coraza in-process with the gateway (OpenResty Lua)
**Rejected.** Ties the WAF to the proxy's failure domain and Lua's limitations (master plan §4.1: avoid Lua as long-term business logic). A separate Go process isolates faults and matches the future Rust gateway boundary.

### A2 — ModSecurity (C) directly
**Rejected.** Coraza is the maintained, pure-Go, OWASP CRS-compatible engine aligned with the Go control plane and the eventual Rust rewrite. ModSecurity's C integration adds native-build and packaging burden without a semantic advantage for this product.

### A3 — Fail-closed on engine error
**Rejected.** A WAF that takes down an application on an internal error violates the availability goals (§1, §12). Fail-open is correct for detection; blocking must be driven by rules, not faults.

---

## Consequences

### Positive
- Isolation of the WAF fault domain; engine errors never drop traffic.
- Deterministic, corpus-verifiable detection/blocking (replay harness).
- Clean migration path to a compiled internal engine (Phase 9) behind the same request-processing contract.

### Negative
- Extra hop adds latency (mitigated by same-host sidecar + keepalive).
- Buffering request bodies costs memory (bounded to 13 MB).
- Baseline rules are intentionally small; enterprise coverage requires CRS + tuning.

### Risks and mitigations
| Risk | Mitigation |
|---|---|
| Rule false positives in blocking mode | Default threshold 1 is aggressive — gate on it; use detect-only first, then raise threshold after FP baseline |
| Body buffering DoS | Bounded `io.LimitReader` (13 MB) + 413 |
| Detect-only "transparent" dropping traffic | P0.2 fix: forward to backend after logging (verified by TestDetectOnlyForwardsToBackend) |
| Rule-set drift vs events | Baseline rules carry stable IDs; tests assert rule IDs in emitted events |

---

## Acceptance for this ADR

- Replay corpus 35/35 in both detect and block modes.
- Detect-only forwards matched requests to the backend (tested).
- Engine error (malformed body/read) logs and forwards, never blocks.
- Sensitive args masked in emitted events (mask_test.go).
- Anomaly threshold blocks at the configured score; detect-only ignores it.

---

## Follow-up ADRs expected

- **ADR-007:** High-volume event transport (NATS/Kafka) and ClickHouse analytics.
- Phase 9: compiled rule matching engine replaces Coraza behind this contract.
