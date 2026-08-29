# ADR-004 — Gateway Registration and Heartbeat Protocol

- **Status:** Proposed (for review)
- **Date:** 2026-08-28
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§3.1, §11, §12, §7.2), `../phase0_srs.md` (FR-0.1-070, §8, §9 T6/T9), `adr-001-initial-architecture.md`, `adr-002-policy-event-schema-v0.md`, `adr-003-config-generation-atomic-reload.md`

---

## Context

The control plane must centrally manage one or more gateways and know their state:

- Which gateways exist, where they are, and what version of software/config they run.
- Whether a gateway applied a config bundle (and if validation failed, why) — per §11, gateways report applied hash, status, and validation errors.
- Whether a gateway is alive and healthy, so the console can show fleet status (FR-0.1-070).
- **Config reach:** the control plane must be able to deliver signed, versioned bundles (§3.1) and the gateway must acknowledge applied or rejected state.

Release 0.1 runs a single node, but the protocol must be designed from day one for the Phase 4 fleet (two or more gateways, mTLS, central management). The failure matrix (§12.3) requires defined behavior when the control plane is unreachable — the gateway must keep serving with its last-known-good config.

---

## Decision

Adopt a **registration + periodic heartbeat + config pull** protocol over HTTPS, with **mutual TLS (mTLS)** as the gateway identity and the only cross-trust-boundary transport. The control plane is the **authority**; the gateway is a **client** that pulls, never pushes.

### D1 — Registration

- On first boot (or when its identity is re-provisioned), a gateway **registers** with the control plane:
  - `POST /api/v1/gateways/register` (HTTPS + mTLS or a one-time registration token in 0.1).
  - Payload: `gateway_id` (ULID, provisioned), `hostname`, `ip`, `version`, `capabilities[]` (e.g., `["http/1.1", "h2"]`), `site`/`cluster` (Phase 4).
- Control plane responds with the **gateway's configuration bundle URL** and a **pull interval**.
- Registration is **idempotent**: re-registering an existing gateway returns its current state (server-side upsert keyed on `gateway_id`).

### D2 — Heartbeat

- Gateway sends a heartbeat every `heartbeat_interval` (default 30 s):
  - `POST /api/v1/gateways/{id}/heartbeat`
  - Payload: `status` (`starting | active | degraded | stopping`), `applied_hash` (current active config bundle hash), `version`, `health` (gateway + backend pool summaries), `last_error` (optional).
- Control plane stores the heartbeat in `gateway_heartbeats` and updates the gateway's `last_seen_at`.
- A gateway with `last_seen_at` older than `offline_threshold` (default 3× interval) shows as **offline** in the console and is excluded from fleet metrics. **No traffic behavior changes** — the gateway keeps serving from its cache (§3.1).

### D3 — Config distribution (pull model)

- The gateway **pulls** config at its pull interval and on heartbeat acknowledgment:
  1. `GET /api/v1/gateways/{id}/config/current` → returns the latest **signed bundle** (ADR-002) with a `bundle_hash` and `next_pull_after`.
  2. Gateway verifies **signature + hash** (must match `signature.value` over canonical JSON).
  3. Gateway runs the ADR-003 pipeline: validate → stage → switch → verify, then reports back via heartbeat: `applied_hash` = new hash, or `last_error` = rejection reason.
  4. If the bundle hash is unchanged, the control plane returns `304 Not Modified` to save bandwidth.
- **Pull, not push:** the gateway controls when it applies config, which naturally provides backpressure and allows offline operation. This is deliberate — see alternatives.
- Config bundles are **signed with an ed25519 control-plane key** (`key_id` in the signature, ADR-002). In Release 0.1, gateway trusts the control-plane CA over mTLS and verifies bundle signatures; key rotation is a Phase 4 concern.

### D4 — Identity and secrets

- Each gateway gets a **provisioned identity**: a unique `gateway_id` (ULID) and a client certificate signed by the control plane's CA (mTLS).
- Private keys are generated **on the gateway**; only the CSR/public part is ever sent to the control plane (master plan §7.2: never return private keys from APIs after import).
- Gateway identity rotation is **automatic** and scheduled (Phase 4); in 0.1, manual re-issue is documented.
- The control plane stores `certificate_metadata` and `secret_references`, never the key material (§7.2).

### D5 — Failure behavior (explicit matrix subset)

| Condition | Behavior |
|---|---|
| Control plane unreachable | Gateway keeps serving from last-known-good config; heartbeats retry with exponential backoff; no traffic impact |
| Heartbeat lost | Console shows offline after threshold; no config revoked remotely (config is signed, not a live session) |
| Config pull fails / bundle rejected | Gateway stays on current active config; `last_error` reported on next successful heartbeat |
| Control plane restarts | Nothing to do for gateways; they re-heartbeat on their own schedule |
| Clock skew | Heartbeats and bundle timestamps use absolute UTC; skew > 5 min logged and flagged as degraded (NTP is a deployment requirement) |
| Re-registration race | Idempotent upsert; the newer registration wins; old connection is a client that will re-pull |

---

## Alternatives considered

### A1 — Control plane pushes config to gateways (server-initiated)
**Rejected.** Push couples traffic availability to control-plane reachability and requires the control plane to track every gateway's liveness for delivery. Pull gives the gateway authority over when it applies config (ADR-003 safety), works offline, and scales to many gateways without fan-out state. Push can be layered later (e.g., NATS/Kafka notification to trigger immediate pull, master plan §4) without changing the contract.

### A2 — Gateway pushes config/telemetry without registration
**Rejected.** Without registration there is no inventory, no per-gateway identity/CA binding, and no way to enforce mTLS scoping. Registration is the trust bootstrap.

### A3 — Anonymous (non-mTLS) heartbeat, API-key auth
**Rejected for gateways.** API keys are replayable and weak for a machine-to-machine trust boundary that will later carry config bundles to production. mTLS binds identity to a certificate and enables revocation (master plan §7.2: mTLS between control plane and gateways).

### A4 — WebSocket/streaming keep-alive for heartbeat
**Deferred.** Streaming is lower-latency for status but adds state to the control plane (connection tracking, reconnect storms) with no benefit for a 30 s status cadence. If Phase 4 needs real-time fleet state, add it behind the same contract.

---

## Consequences

### Positive

- Gateway identity is cryptographically bound (mTLS), matching §7.2.
- Pull model keeps gateways serving through control-plane outages by construction (§3.1, §12.3).
- Heartbeat + applied-hash reporting gives the console a truthful, verifiable fleet state (FR-0.1-070).
- Registration/heartbeat/config-pull are the same contract across one or N gateways — Phase 4 fleet is additive.
- Idempotent registration makes provisioning/restores deterministic.

### Negative

- Pull means config propagation latency up to `pull_interval` (default 30 s), acceptable for 0.1; Phase 3 canary uses a shorter interval for the canary gateway.
- mTLS infrastructure (internal CA, cert issuance/rotation) is required operational tooling from Phase 4 onward.
- The gateway must implement signature verification over canonical JSON (shared codec with the control plane — ADR-002 consequence).

### Risks and mitigations

| Risk | Mitigation |
|---|---|
| Gateway impersonation | mTLS with control-plane CA; per-gateway certs; revocation list at Phase 4 |
| Config replay | Signed bundles with `created_at`; gateway rejects bundles older than current active (anti-replay by monotonic time) |
| Heartbeat flooding | Configurable interval + backoff on control-plane error responses |
| Certificate expiry taking down a gateway | Cert expiry monitoring (§7.2) + automatic rotation (Phase 4); NTP discipline (D5) |

---

## Acceptance for this ADR

- A gateway can register, pull the current signed bundle, apply it via ADR-003, and report `applied_hash` in its next heartbeat.
- Killing the control plane for > 1 minute: gateway keeps serving traffic unchanged; on control-plane return, heartbeats resume and applied-hash state reconciles.
- A tampered bundle (wrong signature or hash) is rejected and `last_error` reflects the reason; active config stays serving.
- Idempotent re-registration returns the existing gateway without duplicate rows.
- Clock-skew test: a gateway with > 5 min skew is flagged `degraded` in console.

---

## Follow-up ADRs expected

- **ADR-005:** Coraza/ModSecurity WAF engine integration (extends the config bundle with a `waf_policy` section).
- **ADR-007 (later):** High-volume event transport (NATS JetStream/Kafka) and ClickHouse analytics — reuses the same gateway identity model.
