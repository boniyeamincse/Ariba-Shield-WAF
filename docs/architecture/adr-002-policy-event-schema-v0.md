# ADR-002 — Policy and Event Schema v0

- **Status:** Proposed (for review)
- **Date:** 2026-08-28
- **Owner:** Platform Architect
- **References:** `../enterprise_waf_development_master_plan.md` (§8, §11), `../phase0_srs.md` (Data model v0, API v0), `adr-001-initial-architecture.md`

---

## Context

The data plane and control plane must share versioned, machine-readable contracts from the first release (§11). Two contracts are foundational and every later feature builds on them:

1. **Policy/config schema** — the declarative configuration the control plane compiles and gateways apply.
2. **Event schema** — the structured security/access events gateways emit and the control plane ingests.

These must be forward-compatible with Phase 2+ (WAF rules, anomaly scoring, rate limits, bot defense) without breaking Release 0.1 consumers.

---

## Decision

Adopt a **versioned, additive-only** contract pair. Major version in the bundle/event header; fields are only ever added, never removed or renamed; deprecated fields carry explicit `deprecated` markers and a removal horizon.

### D1 — Common conventions

- Identifiers: ULID strings. `schema_version` present at the top level of every document.
- Timestamps: RFC 3339 with timezone offset (`2026-08-28T10:00:00Z`).
- Byte counts: integers; durations: milliseconds.
- All enums are **stable codes** (e.g., `block`, `log`, `pass`). Human labels live in the console localization layer only (master plan §6.10).
- Unknown fields are preserved and ignored by older consumers; consumers must never fail on unknown fields (additive forward-compat).

### D2 — Policy/config schema v0 (declarative, versioned, hashable)

Top-level structure (JSON; also stored as PostgreSQL `config_versions.blob`):

```jsonc
{
  "schema_version": "0.1",
  "config_id": "01J...",            // ULID, unique per document
  "created_at": "2026-08-28T10:00:00Z",
  "created_by": "01J...",           // user ULID
  "replaces": "01J...",             // previous active config_id, null if none
  "gateway_targets": ["gateway-a"], // Release 0.1: single-node "local"
  "settings": {
    "log_level": "info",
    "event_retention_days": 30
  },
  "virtual_servers": [
    {
      "id": "01J...",
      "name": "app-prod",
      "listen_addr": "0.0.0.0",
      "listen_port": 443,
      "tls": {
        "enabled": true,
        "certificate_ref": "01J...",   // reference, never the key material
        "min_version": "1.2",
        "protocols": ["http/1.1", "h2"]
      },
      "default_backend_pool_id": "01J...",
      "routes": [
        {
          "id": "01J...",
          "path": "/api/",
          "match": "prefix",
          "backend_pool_id": "01J..."
        }
      ],
      "limits": {
        "max_request_line": 8192,
        "max_header_size": 32768,
        "max_body_size": 10485760
      }
    }
  ],
  "backend_pools": [
    {
      "id": "01J...",
      "name": "api-nodes",
      "application_id": "01J...",
      "lb_algorithm": "round_robin",   // round_robin | weighted | least_conn (later)
      "health_monitor": {
        "type": "http",                // tcp | http
        "interval_ms": 5000,
        "timeout_ms": 2000,
        "fail_threshold": 3,
        "pass_threshold": 2,
        "http_path": "/healthz",
        "http_expected_status": [200]
      },
      "nodes": [
        { "id": "01J...", "host": "10.0.0.11", "port": 8080, "weight": 1, "active": true }
      ]
    }
  ],
  "headers": {
    "trusted_proxy_headers": ["x-forwarded-for", "x-forwarded-proto", "x-forwarded-host"],
    "request_id_header": "x-request-id"
  },
  "signature": { "algorithm": "ed25519", "key_id": "01J...", "value": "base64..." }
}
```

**Rules of the policy schema:**

- Every deployment document is **immutable**; changes create a new `config_id` with `replaces` pointing at the prior one.
- The bundle **hash** (SHA-256 of canonical JSON) is computed by the control plane and verified by the gateway on receipt.
- The control plane **signs** the bundle (`signature` above); gateways verify before applying (mTLS + signature; master plan §7.2).
- Extension points reserved (empty or omitted in v0, defined in later ADRs): `waf_policy` (Phase 2), `rate_limits` (Phase 3), `ip_lists` (Phase 3), `bot_policy` (Phase 7), `api_schemas` (Phase 5).

### D3 — Event schema v0 (structured, stable codes)

Gateway emits one JSON event per request to stdout (async, non-blocking). Security events (Phase 2+) extend this base.

```jsonc
{
  "schema_version": "0.1",
  "event_id": "01J...",             // ULID
  "request_id": "01J...",
  "timestamp": "2026-08-28T10:00:00Z",
  "event_type": "access",           // access | security | audit | health
  "gateway_id": "gateway-a",
  "virtual_server_id": "01J...",
  "application_id": "01J...",
  "client": {
    "ip": "203.0.113.9",
    "port": 51234,
    "forwarded_for": "198.51.100.4" // from trusted proxy headers, if present
  },
  "request": {
    "method": "GET",
    "path": "/api/orders",
    "query": "id=123",
    "http_version": "1.1",
    "host": "app.example.com",
    "content_type": "application/json",
    "body_size": 512,
    "headers_count": 9
  },
  "response": {
    "status": 200,
    "bytes": 2048,
    "latency_ms": 3,
    "backend_node": "10.0.0.11:8080"
  },
  "decision": { "action": "pass" },   // pass | log | block | challenge | rate_limit | redirect (later phases fill)
  "matched_rules": []                 // Phase 2+: array of {rule_id, phase, category, score}
}
```

**Rules of the event schema:**

- No request bodies, cookies, or credentials ever in events (§22 rule 5; SRS FR-0.1-032). Sensitive field masking is a Phase 2+ concern at the payload level; v0 excludes payloads entirely.
- `decision.action` codes are stable and never localize in the event stream — labels are translated at display time (§6.10).
- Security events (Phase 2) must not alter the base shape; they add optional fields only.
- Wazuh/SIEM output (Phase 2) maps this schema to syslog/CEF/LEEF/JSON; v0 keeps fields losslessly available for that mapping.

---

## Alternatives considered

### A1 — No shared schema, gateway-specific config
**Rejected.** Violates the versioned-contract invariant (§11). Multi-gateway fleet management (Phase 4) would require a rewrite and per-gateway drift would be unrecoverable.

### A2 — Single unversioned JSON blob
**Rejected.** Breaks rollback, canary, backward-compat matrix, and safe upgrades. Versioning and immutability are release gates (SRS §11).

### A3 — Event schema as free-form key/value
**Rejected.** Destroys the ability to correlate by application, IP, identity, device, session, and attack type (§6.8) and to validate Wazuh/SIEM mappings.

---

## Consequences

### Positive

- Data plane and control plane evolve independently behind a versioned contract.
- Rollback/canary/safe-upgrade all operate on immutable, hashable, signed documents.
- Additive-only rule means older gateways never break on newer configs.
- Events are forward-compatible with the Phase 2 security pipeline.

### Negative

- Additive-only discipline requires review; deprecation must be scheduled, not ad hoc.
- Canonical JSON + hashing requires a single serialization implementation to avoid hash mismatch across languages.
- Schema is intentionally minimal for 0.1; Phase 2/3 additions will need companion ADRs.

### Risks and mitigations

| Risk | Mitigation |
|---|---|
| Hash mismatch between Go control plane and gateway | One canonical serialization lib + golden hash tests across both sides |
| Field drift between schema and code | Schema is source of truth; codegen/golden tests in CI |
| Unknown-field churn from future phases | Additive-only policy; unknown fields preserved, never dropped |

---

## Acceptance for this ADR

- Schema v0 checked in as the single source of truth (JSON Schema files).
- Golden compilation tests: control-plane output hashes match expected fixtures.
- Gateway rejects unsigned or wrong-hash bundles (Release 0.1: hash verify; mTLS signing per ADR-004/Phase 4).
- Events from a live request round-trip into PostgreSQL without loss (NFR-0.1-005).

---

## Follow-up ADRs expected

- **ADR-003:** Configuration generation and atomic reload mechanism.
- **ADR-004:** Gateway registration and heartbeat protocol.
- **ADR-005:** Coraza/ModSecurity integration and CRS baseline (extends policy schema with `waf_policy`).
- **ADR-006:** Logging and event pipeline design.
