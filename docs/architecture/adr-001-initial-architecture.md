# ADR-001 — Initial System Architecture

- **Status:** Accepted
- **Date:** 2026-08-28
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§4, §10, §23, §24), `../phase0_srs.md` (Data model v0, Deployment topology)

---

## Context

Ariba Shield WAF must evolve from a single-node lab proxy to an F5-class centralized enterprise WAF over a 30–48 month horizon. Early architectural choices must:

1. Deliver a usable **Release 0.1 (Lab)** within 12 weeks.
2. Preserve a credible path to a **custom high-performance Rust data plane** later without reworking the control plane.
3. Keep the data plane and control plane decoupled with versioned contracts from day one.
4. Never let the control plane sit in the traffic path.
5. Keep the traffic path safe even when the control plane or event pipeline is unavailable.

The master plan (§23) recommends a specific starting stack. This ADR records why, what alternatives were rejected, and the consequences the team must accept.

---

## Decision

Adopt a **two-plane architecture** with the following decisions, in force for Release 0.1 and expected to persist with incremental replacement:

### D1 — Two planes, one monorepo

- **Data plane** handles live client traffic: TLS, HTTP parsing, normalization, inspection, enforcement, proxying, load balancing, health checks.
- **Control plane** is the centralized web software: apps, policies, certs, rules, users, analytics, updates, licensing, cluster management, audit.
- Single **monorepo** (`shield-waf/`, structure per master plan §10) for atomic schema and policy changes. Split into services only when scaling, isolation, or team ownership justifies it.

### D2 — Control plane is out-of-band

- Control plane distributes policies and receives telemetry **asynchronously**.
- It is never in the client→gateway→backend traffic path (master plan §3.2).
- Gateways keep a **signed, versioned, locally cached last-known-good configuration**; if the control plane fails, protected applications keep serving (§3.1).

### D3 — Proven components first, custom later

| Area | Release 0.1 decision | Later trajectory |
|---|---|---|
| Traffic proxy | OpenResty/Nginx | Rust gateway where justified (Phase 9) |
| Rule engine | None in 0.1; Coraza/ModSecurity-compatible + OWASP CRS in Phase 2 | Optimized compiled internal engine + compatibility layer |
| Control API | Go (modular monolith) | Go microservices/modules |
| High-perf components | Go initially | Rust for parser/normalization/matching/proxy hot paths |
| Console | Next.js + TypeScript | Same, modular enterprise console |
| Database | PostgreSQL | PostgreSQL HA + partitioning |
| Config/cache distribution | Redis Streams (introduced when needed) | NATS JetStream or Kafka per measured scale |
| Event analytics | PostgreSQL (via log pipeline) | ClickHouse for high-volume events |
| Object storage | MinIO/S3-compatible (later phases) | Distributed object storage |
| Telemetry | OpenTelemetry + Prometheus | Same with scalable collectors |
| Packaging | Docker Compose | OCI images, VM appliance, Kubernetes/Helm |
| Automation | Ansible | Ansible + Kubernetes Operator |

### D4 — Language boundaries

- **Rust:** latency-sensitive, memory-safe traffic components only.
- **Go:** control plane, agents, orchestration, policy compiler, update service.
- **TypeScript:** web console only.
- **Python:** offline analytics, training, QA tooling, rule research — never the core synchronous traffic path.
- **Lua:** temporary OpenResty extensions only; never the long-term business-logic layer.

Each service must have an explicit reason to exist. Start as a **modular monolith** for the control plane.

### D5 — Versioned contracts between planes

- REST/JSON management API with OpenAPI specification (§11 of plan).
- Every deployment produces an immutable **bundle hash**; gateways report applied hash, status, and validation errors.
- **No partial policy activation**: validate, stage, atomically switch.
- Backward compatibility matrix between control-plane and gateway versions from the first release.

### D6 — Hard engineering invariants (from master plan §22)

1. Never parse the same request differently in security and proxy layers.
2. Never let logging failure block or crash the live traffic path.
3. Never deploy partially compiled policy.
4. Never publish performance without exact test conditions.
5. Never add a protocol until it has conformance, fuzz, evasion, and resource-limit tests.
6. Never call the product enterprise/F5-grade until independent testing supports the claim.

---

## Alternatives considered

### A1 — Custom Rust gateway from day one
**Rejected.** Would delay the first usable release by months, risk parser/security correctness before semantics are stabilized, and make it impossible to build a complete test corpus and performance baselines first. The master plan explicitly requires a proven proxy + proven rule engine underneath a custom control plane first (§1, §4).

### A2 — Single-plane monolithic appliance (proxy + management in one binary/UI)
**Rejected.** Violates the core control-plane rule (§3.2): management must never sit in the traffic path. Would couple scaling, upgrade, and failure domains and make centralized multi-gateway management (Phase 4) a rewrite.

### A3 — Go gateway instead of OpenResty/Nginx
**Rejected for 0.1.** Nginx/OpenResty provides mature TLS, HTTP/1.1/2, proxying, and health-check capabilities with a huge conformance base. A Go reverse proxy is viable but reinvents parser/proxy maturity for no Release 0.1 benefit. Go remains the language for the control plane and later can be superseded by Rust for hot paths.

### A4 — Envoy instead of OpenResty/Nginx
**Deferred.** Envoy is a strong candidate with first-class HTTP/2 and dynamic config. OpenResty was chosen for 0.1 because its Lua extension model and Nginx semantics fit the Coraza/ModSecurity WAF path (Phase 2) and the team's near-term tooling. Envoy can be re-evaluated if dynamic xDS-style control becomes the dominant need before the custom Rust gateway matures.

### A5 — Kafka from the start for config/event distribution
**Rejected for now.** Operational overhead is unjustified for Release 0.1 volumes. Redis Streams initially, then NATS JetStream or Kafka per measured scale (§4).

### A6 — Separate microservices from day one
**Rejected.** Modular monolith first; split services only when scaling, isolation, or team ownership justifies it (§4.1).

---

## Consequences

### Positive

- Fastest credible path to a working, tested, explainable Release 0.1.
- Control-plane/gateway contract is versioned and stable from day one; later data-plane replacement (Rust) is a drop-in behind the same contract.
- Failure isolation: control plane can be down with traffic unaffected.
- Established performance baselines on the proven gateway become the benchmark the Rust gateway must beat (§15 Phase 9).
- Monorepo keeps schema, policy, and code changes atomic.

### Negative

- OpenResty/Nginx + Lua carry operational constraints (business logic in Lua discouraged; long-term gateways must move off Lua).
- Coraza/CRS engine (Phase 2) is a proven engine but not a fully custom product feature; differentiation must come from the control plane, semantics, and later the Rust engine.
- The eventual custom Rust gateway is a large, long-horizon effort (Phase 9) that must reach parity before replacing the proven gateway.
- Monorepo grows; guardrails (CI, ownership, ADR review) are needed to avoid cross-team coupling.

### Risks and mitigations

| Risk | Mitigation |
|---|---|
| Gateway and control plane drift on contract | Versioned OpenAPI + bundle hash reporting + backward-compat matrix (D5) |
| Lua creep into business logic | Enforce Lua as temporary extension layer only (D4) |
| Control-plane outage perceived as product outage | Out-of-band design + last-known-good cache + communication of architecture (D2) |
| Performance claims without baselines | Reference-hardware performance method from Phase 0; publish only with exact test conditions (D6) |

---

## Acceptance for this ADR

- Architecture review approved (Phase 0 exit criterion).
- Repo and CI created with the monorepo structure per §10.
- Versioned policy and event schema v0 in place.
- Data plane and control plane speak the versioned contract defined in D5.

---

## Follow-up ADRs expected

- **ADR-002:** Policy and event schema v0 (field-by-field contract).
- **ADR-003:** Configuration generation and atomic reload mechanism (last-known-good semantics).
- **ADR-004:** Gateway registration and heartbeat protocol (Phase 2 placeholder).
- **ADR-005:** Coraza/ModSecurity engine integration and CRS baseline (Phase 2).
- **ADR-006:** Logging and event pipeline design (non-blocking, Wazuh/SIEM forward-compatible).
