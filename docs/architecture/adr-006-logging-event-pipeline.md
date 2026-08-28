# ADR-006 — Logging and Event Pipeline Design

- **Status:** Proposed (for review)
- **Date:** 2026-08-28
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§22 rules 2 & 5, §6.8, §12.3), `../phase0_srs.md` (FR-0.1-030..034, NFR-0.1-005, NFR-0.1-013), `adr-002-policy-event-schema-v0.md`

---

## Context

Two invariants shape the logging design (master plan §22):

> **Rule 2:** Never let logging failure block or crash the live traffic path.
> **Rule 5:** Never store unmasked secrets merely for analyst convenience.

Release 0.1 must emit structured JSON access events (event schema v0 from ADR-002), export baseline Prometheus metrics, and keep zero audit-event loss on the control plane — all without ever coupling a traffic decision to a log write. The design must also be forward-compatible with the Phase 2 security-event pipeline and Wazuh/SIEM output (§6.8).

---

## Decision

Adopt a **decoupled, asynchronous, fail-open** event pipeline. The gateway writes structured events to stdout via a **non-blocking bounded buffer**; a separate sidecar/collector consumes, normalizes, and forwards them. Traffic paths never write to a blocking destination.

### D1 — Gateway emission (data plane)

- Access events are produced per request using **event schema v0** (ADR-002 D3) and written to **stdout** as newline-delimited JSON.
- Emission is **non-blocking**: the worker writes to a bounded in-process buffer; if the buffer is full, events are **dropped with a counter metric** (never block, never crash, never retry synchronously) — master plan rule 2.
- **No payloads**: events never contain request bodies, cookies, credentials, or full headers (ADR-002 D3; SRS FR-0.1-032). Only metadata and status fields are emitted in 0.1.
- Sensitive field masking (Phase 2+) is implemented as an upstream masking layer, not by adding payloads to events.

### D2 — Collector/forwarder (observability plane)

- A lightweight sidecar (Go, part of the control-plane monorepo) reads the gateway stdout stream.
- Collector responsibilities:
  1. Normalize (already structured), timestamp, and add collector metadata.
  2. Forward to the configured sinks **asynchronously** with bounded retries + backoff.
  3. Never backpressure the gateway (collector failure only loses its own buffer, recorded as a drop counter).
- Sinks for Release 0.1: stdout/stderr capture, optional file (JSON-lines), OpenTelemetry Collector, Prometheus (metrics only).
- **Forward-compatible**: the same collector gains syslog/CEF/LEEF/JSON and Wazuh output adapters in Phase 2 without changing the event shape.

### D3 — Control-plane event persistence

- Control plane ingests normalized events into **PostgreSQL** (append-mostly table with indexes on time, application_id, virtual_server_id, request_id).
- **Zero audit-event loss** is a Release 0.1 NFR (NFR-0.1-005) — applies to *audit* events (config mutations) and *health* events, which the control plane writes itself and are loss-tolerant-tolerant by design: they are written through a durable path, while high-volume *access* events are **best-effort** with drop counters and are explicitly allowed to be lost under overload (traffic must keep flowing).
- Events are immutable after insert; retention policy (default 30 days) is configurable; archival to object storage is a later-phase concern.

### D4 — Metrics

- Prometheus baseline (FR-0.1-034): `shield_http_requests_total{vs,status,method}`, `shield_http_request_duration_ms_histogram`, `shield_upstream_status_total{pool,node}`, `shield_event_drops_total`, `shield_active_connections`.
- All traffic-path metrics are updated **synchronously in-memory** (cheap counters/histograms) — never gated on I/O.

### D5 — Correlation and audit

- Every request gets a `request_id` (ULID) injected into upstream headers and returned to the client (FR-0.1-030); it appears in access events, metrics labels, and audit events, enabling end-to-end correlation.
- Audit events (config mutations, auth changes) are **separate, durable, append-only** — never mixed into the lossy access-event stream (SRS §3.5; master plan §7.4).
- Access events are correlated to security events later by `request_id`; event/incident system (§6.8) reuses this correlation key.

---

## Alternatives considered

### A1 — Synchronous blocking log writes in the traffic path
**Rejected.** A slow/blocked log destination would stall or crash the proxy — directly violates rule 2.

### A2 — Ship directly to Wazuh/SIEM from the gateway
**Rejected for 0.1.** Couples the data plane to external sink availability and makes the traffic path depend on a network destination. The collector decouples this (Phase 2 adapter pattern).

### A3 — Store full request bodies in events for analyst convenience
**Rejected.** Violates rule 5 and ADR-002 D3. Masking is the sanctioned mechanism, added at the right phase.

### A4 — Single shared queue for both access and audit events
**Rejected.** Audit events must be durable and lossless (NFR-0.1-005, §7.4); access events are high-volume and loss-tolerant. Mixing them forces one policy on both, either losing audit events or blocking on access volume.

---

## Consequences

### Positive

- Traffic path is decoupled from any I/O failure (rule 2 satisfied by construction).
- Forward-compatible event shape: Phase 2 security events + Wazuh/SIEM adapters add no breaking change.
- `request_id` correlation key is established now, reused by events/incidents (§6.8).
- Loss metrics (`shield_event_drops_total`) make overload visible rather than silent.

### Negative

- High-volume access events are explicitly best-effort under overload — acceptable for 0.1, must be a documented, measured property (backlog/backpressure metrics in Phase 4 with Kafka/NATS, master plan §4).
- Collector is another moving part to operate (mitigated: it is a thin Go sidecar with no state of its own).
- Two event streams (access vs audit) require clear documentation to avoid misuse.

### Risks and mitigations

| Risk | Mitigation |
|---|---|
| Collector silently losing events | Drop counters + health metric per sink; alerting later |
| Access event loss mistaken for security issue | Documented best-effort semantics + metrics (loss is measurable) |
| Audit event loss | Durable write path, never mixed with lossy stream |
| Log injection | Structured JSON, no raw multi-line user strings in freeform fields (SRS T5) |

---

## Acceptance for this ADR

- `kill` the collector and flood traffic: gateway continues serving; event drop counter increments; no traffic impact (chaos test).
- Audit events survive collector outage (durable path).
- A live request produces an access event in PostgreSQL with a valid `request_id` and no payload fields (UC-2 + FR-0.1-032).
- Prometheus metrics scrapeable; request histogram matches p50/p99 budget (NFR-0.1-001/002).
- Wazuh/syslog adapter in Phase 2 consumes the same event shape without schema change.

---

## Follow-up ADRs expected

- **ADR-004:** Gateway registration and heartbeat protocol (health/status events feed this pipeline).
- **ADR-005:** Coraza/ModSecurity integration and security-event extension of schema v0.
- **ADR-007 (later):** High-volume event analytics migration to ClickHouse (§8) and Kafka/NATS transport (§4).
