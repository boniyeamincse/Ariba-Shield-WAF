# Ariba Shield WAF — Release 0.1 (Lab)

**Status:** Lab release (Sprint 6)
**Production position:** Lab / internal evaluation only. NOT production blocking.
**References:** `enterprise_waf_development_master_plan.md` (§15 Phase 1, §16 release editions), `phase0_srs.md`

---

## 1. What Release 0.1 provides

A single-node lab platform on Ubuntu Server LTS:

| Capability | Status |
|---|---|
| Centralized management console (English/Bangla) | ✅ |
| Go control-plane API (`/api/v1`, OpenAPI documented) | ✅ |
| OpenResty/nginx reverse-proxy gateway | ✅ |
| Virtual servers, TLS termination, backend pools | ✅ |
| Active/passive health checks | ✅ |
| Safe config generation, validation, atomic reload, last-known-good | ✅ |
| Gateway registration + heartbeat (ADR-004 contract) | ✅ |
| Structured JSON access logs (request ID) | ✅ |
| Prometheus metrics endpoint | ✅ |
| Coraza WAF engine — **detection-only** (transparent) | ✅ |
| CRS-style baseline rules (SQLi, XSS, cmd, traversal, LFI) | ✅ |
| Rate limiting + IP lists (Phase 3 safe blocking) | ✅ |
| Security events pipeline (engine → ingestor → PostgreSQL) | ✅ |
| Sensitive-field masking (rule 5) | ✅ |
| Wazuh/syslog forwarder (adapter, not yet wired to a live Wazuh) | ⚠️ adapter ready |
| Audit log (append-only, immutable) | ✅ |
| Real authentication/RBAC | ✅ session-cookie auth + 7 roles (mock gated behind `AUTH_MOCK_ENABLED`) |
| TOTP MFA | ✅ |
| Incidents (assign/escalate/close/reopen) | ✅ |
| Dashboard (8 widgets) | ✅ |
| Learning sessions/suggestions | ✅ |
| Reports, integrations, notification channels, license, settings | ✅ |
| Full OWASP CRS set | ❌ Phase 2+ |
| Blocking enforcement posture | ❌ detection-only in 0.1 |
| Bot defense / API schema protection | ❌ Phase 3+ (CRUD scaffolded) |
| HA, mTLS fleet, multi-tenancy, SSO | ❌ Phase 4+ |

**Explicitly NOT included:** blocking enforcement, full OWASP CRS set,
bot defense, API schema protection, HA, mTLS fleet, multi-tenancy, SSO.

---

## 2. Quick start

### 2.1 Prerequisites

- Ubuntu Server LTS (24.04+), Docker + Docker Compose plugin.
- DNS name for the gateway or `localhost` for a pure lab.

### 2.2 Run the stack

```sh
cd infra/compose
cp .env.example .env          # EDIT the passwords first — no defaults
docker compose up -d --build
```

Services (all management ports bound to 127.0.0.1 by default):

| Service | Bind | Notes |
|---|---|---|
| console-web | http://127.0.0.1:3000 | bilingual UI |
| control-api | https://127.0.0.1:8443 | management API |
| gateway | :443 | reverse proxy |
| waf-engine | 127.0.0.1:8082 | detection-only Coraza |
| postgres / redis | 127.0.0.1 | state |

### 2.3 First-run admin

The API creates the initial admin from `ADMIN_INITIAL_EMAIL` / `ADMIN_INITIAL_PASSWORD`
on first boot (no default credentials; `config.go` fails fast if unset).

### 2.4 Register a gateway and proxy a test app

```sh
# Register gateway (idempotent)
curl -k -X POST https://127.0.0.1:8443/api/v1/gateways/register \
  -H 'Content-Type: application/json' \
  -d '{"gateway_id":"gw-01","hostname":"lab-gw","capabilities":["http/1.1"]}'

# Create an application
curl -k -X POST https://127.0.0.1:8443/api/v1/applications \
  -H 'Content-Type: application/json' -d '{"name":"demo-app"}'

# Create a security policy bound to it (transparent)
curl -k -X POST https://127.0.0.1:8443/api/v1/security-policies \
  -H 'Content-Type: application/json' -d '{"name":"demo","enforcement_mode":"transparent"}'
```

---

## 3. Pipeline overview (Release 0.1)

```text
Client ──HTTPS──► OpenResty gateway ──► waf-engine (Coraza, detect-only) ──► backend-a
                        │                       │
                        │  JSON-lines stdout    │  security events (masked)
                        ▼                       ▼
                   access log            event-ingestor ──► PostgreSQL security_events
                                                      └─► wazuh-forward (adapter)
```

**Safety invariant (master plan §22):** the WAF engine fails **open** in detection
mode — an engine error never blocks or crashes traffic.

---

## 4. Testing

```sh
# Engine unit tests (detect/block/body/masking/rule matching)
cd services/waf-engine && go test ./...

# Replay harness against a live engine (detect or block)
bash tests/replay/replay.sh --host http://127.0.0.1:8082 --mode detect
bash tests/replay/replay.sh --host http://127.0.0.1:8082 --mode block

# Gateway failure drills (restart, invalid config, last-known-good, safe-mode)
bash tests/failover/gateway-failure-tests.sh

# API drills (registration, heartbeat, policies, invalid-mode rejection)
bash tests/failover/api-failure-tests.sh   # requires compose stack up
```

Reference result on the lab bench (Oct 2026):
- Replay detect: **6/6** corpus cases (cmdi, legitimate, lfi, sqli, traversal, xss)
- Replay block: **6/6** corpus cases
- Gateway failure drills: **4/4**
- Go tests: **control-api 18, waf-engine 27, policy-compiler 15** (event-ingestor pending)
- Console-web Vitest: **10** component/permission tests

---

## 5. Known limitations (Release 0.1)

1. **No SSO** — session-cookie auth + role-based access + TOTP MFA is implemented,
   but OIDC/SAML SSO and break-glass are Phase 8.
2. **No full OWASP CRS** — baseline rules only; CRS integration is Phase 2+.
3. **Wazuh forwarder** is an adapter that reads JSON-lines stdin; a live Wazuh
   agent + output transport is not yet wired.
4. **Single node** — no HA, no mTLS fleet, no canary distribution (Phase 4).
5. **TLS termination uses self-signed/generated certs in lab**; certificate
   management is Phase 4.
6. Engine emits status=200 for detection-only blocks by design; the console
   event stream shows the detection (decision_action=log).

---

## 6. Upgrade / rollback note

`docker compose up -d --build` rebuilds in place. For the gateway, the
last-known-good mechanism (ADR-003) restores the previous active config on
restart if a new config fails validation — this is the rollback path for 0.1.
