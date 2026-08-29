# Ariba Shield WAF — Task List

Generated: 2026-08-29  
Source: Full project analysis (data plane, control plane, docs, tests, CI)  
Priorities: **P0** (critical bug) → **P1** (high) → **P2** (medium) → **P3** (low)

---

## P0 — Critical bugs (breaks functionality or security)

- [x] **P0.1** `apply-config.sh:28` validates the bare generated config fragment, not the wrapper main config. `nginx -t -c "$staging_dir/generated.conf"` always fails because the fragment has no `http {}` wrapper. Every config apply is rejected. Fix: write a wrapper main config (as `entrypoint.sh` does), then `nginx -t -c` on the wrapper.
- [x] **P0.2** `engine.go` detect-only path: `handleInterruption` returns an empty 200 **without calling `proxy.ServeHTTP`**. The backend never receives matched requests. In "transparent" mode, the WAF silently drops traffic instead of forwarding it. Fix: log the event, then call `proxy.ServeHTTP` to forward to the backend.
- [x] **P0.3** `docker-compose.yml:88` mounts the gateway config volume at `/etc/shield-waf/config`, but the Dockerfile sets `CONFIG_STORE=/var/lib/shield-waf/config`. The dev stack lands at the wrong path → gateway always boots into safe-mode 503.
- [x] **P0.4** `middleware/auth.go` trusts `X-User-Email` and `X-User-Role` headers verbatim. A caller can send `X-User-Role: Super Admin` and get full access. Anonymous (no headers) gets Read Only. Fix: session-cookie auth; mock gated behind `AUTH_MOCK_ENABLED=true`.
- [x] **P0.5** `store.go:66-68` stores the admin password as plaintext in `password_hash`. No bcrypt/argon2 in the dependency tree. Fix: hash with bcrypt before storing.
- [x] **P0.6** `middleware/audit.go` is placed **before** Auth and RequestID in the middleware stack (`router.go:53-59`). The goroutine reads `AuditFromContext(r.Context())` and `RequestIDFromContext(r.Context())` which are always empty because they haven't been set yet. The audit trail is untraceable (empty actor, empty request ID, path-as-resource). Fix: reorder so Auth+RequestID run before Audit.
- [x] **P0.7** `engine.go:209` `io.ReadAll(r.Body)` buffers the entire request body before Coraza's `SecRequestBodyLimit` is checked. A 1 GB upload is held in RAM. Fix: wrap with `io.LimitReader` (13 MB cap) + 413 rejection.
- [x] **P0.8** `engine.go` `clientIP(r.RemoteAddr)` ignores `X-Forwarded-For`/`X-Real-IP`. Behind the gateway, every client is the gateway's IP → IP blocklist and rate limiter apply to the proxy, not the client. Also splits on last colon, breaking IPv6 (`[::1]:8080` → `[::1]` which `netip.ParseAddr` rejects). Fix: `net.SplitHostPort` + optional trusted-proxy XFF handling.

---

## P1 — High (blocking Phase 4 deliverables or major features)

^- [x] **P1.1** Multi-gateway config distribution: add `GET /api/v1/gateways/{id}/config/current` that returns the signed bundle (with `304 Not Modified`). Add `applied_hash` reconciliation.
^- [x] **P1.2** Load-balancing algorithms: implement `ip_hash` (sticky sessions by client IP) and `consistent_hash` (by URI) in the policy compiler. The `weighted` algorithm is a no-op alias for `round_robin` — fix.
^- [x] **P1.3** TLS re-encryption to upstreams: add `protocol` (`http`/`https`) per backend node, emit `proxy_ssl_*` directives and `proxy_pass https://` when `protocol=https`.
^- [x] **P1.4** Session persistence: add `sticky_cookie` / `sticky_route` fields to the backend pool schema and render `ip_hash` or `sticky` in the nginx config.
^- [x] **P1.5** Connection draining: add `drain` state + `slow_start` to backend nodes. Active=false nodes should drain gracefully, not be hard-dropped from the upstream.
^- [x] **P1.6** Active health checks: `HealthMonitor` fields (interval, thresholds, HTTP path, expected status) are declared in the schema but never rendered in the nginx config. Wire them to `health_check` directives.
^- [x] **P1.7** Backup/restore: add `POST /api/v1/backups` (trigger encrypted DB dump) and `POST /api/v1/backups/restore`. Add `backups` table.
^- [x] **P1.8** `make test` must include `services/waf-engine`, `services/policy-compiler`, and `services/event-ingestor`. Currently only `apps/control-api` is tested.
^- [x] **P1.9** `make gen-check` must run in CI as part of the `schema` job to catch schema-SDK drift.
- [x] **P1.10** Real auth: session-cookie auth replaces `X-User-Email`/`X-User-Role` headers. Login/logout/me endpoints, bcrypt password hashing, HTTP-only/Secure/SameSite cookies, user_roles table + role seeding. ⚠️ Remaining: refresh endpoint, Redis session storage (Postgres today), CSRF protection.
^- [x] **P1.11** Idempotency keys: implement `Idempotency-Key` middleware + table. FR-0.1-043 (M) — documented in OpenAPI but zero code.
^- [x] **P1.12** Optimistic concurrency: implement `If-Match`/ETag + `WHERE version =` checks on all mutation handlers. FR-0.1-044 (M) — version columns exist but no checks.
- [x] **P1.13** `engine.go`/`iplist`: `IsAllowed` is never called in the handler. `--allowed-ips` has zero effect. Fix: check allow list before block list.

---

## P2 — Medium (important, non-blocking, Phase 5+ features)

- [ ] **P2.1** Authentication & business-abuse test corpus: add `tests/corpus/auth-abuse.txt` (credential stuffing, MFA bypass, session fixation, brute-force indicators).
- [ ] **P2.2** Request smuggling test corpus: add `tests/corpus/smuggling.txt` (CL.TE, TE.CL, H2.CL, H2.TE patterns).
- [ ] **P2.3** Parser differentials test corpus: `tests/corpus/parser-diff.txt`.
- [ ] **P2.4** Encoding & double-encoding evasions corpus: `tests/corpus/encoding.txt`.
- [ ] **P2.5** SSTI test corpus: `tests/corpus/ssti.txt`.
- [ ] **P2.6** XXE/XML bomb test corpus: `tests/corpus/xxe.txt`.
- [ ] **P2.7** JSON depth bomb / parser abuse corpus: `tests/corpus/json-abuse.txt`.
- [ ] **P2.8** Multipart/file upload evasions corpus: `tests/corpus/multipart.txt`.
- [ ] **P2.9** Cache poisoning/deception corpus: `tests/corpus/cache-poison.txt`.
- [ ] **P2.10** Host-header & routing attack corpus: `tests/corpus/host-header.txt`.
- [ ] **P2.11** WebSocket bypass corpus: `tests/corpus/websocket.txt`.
- [ ] **P2.12** Bot/challenge bypass corpus: `tests/corpus/bot-bypass.txt`.
- [ ] **P2.13** Resource exhaustion corpus: `tests/corpus/resource-exhaust.txt`.
- [ ] **P2.14** `replay.sh` should verify security events are emitted in detect mode (not just HTTP status codes).
- [ ] **P2.15** `replay.sh` should support POST bodies and header injection, not just GET query strings.
- [ ] **P2.16** `replay.sh` should verify rule IDs and match details in emitted events.
- [ ] **P2.17** Detach the duplicated nginx template (`gateways/openresty-gateway/nginx/templates/shield.conf.tmpl` vs `services/policy-compiler/internal/compiler/templates/shield.conf.tmpl`). Use a symlink or single source of truth.
- [ ] **P2.18** `engine.go` `newID()` is not a ULID — uses `fmt.Sprintf("%d-%x", time.Now().UnixNano(), os.Getpid())`. Replace with ULID generation.
- [ ] **P2.19** `engine.go` `event.go` never populates `virtual_server_id`/`application_id` (only `GATEWAY_ID` env var). These are empty in all events.
- [ ] **P2.20** `ratelimit/limit.go` buckets map is never evicted — unbounded memory growth per unique client IP. Add TTL eviction.
- [ ] **P2.21** `event-ingestor/ingest.go` — unbounded `in.buf` growth on DB outage (OOM risk). Add a cap and backoff.
- [ ] **P2.22** `event-ingestor/ingest.go` — final flush on shutdown uses a cancelled context, dropping the last batch.
- [ ] **P2.23** `handlers/auth.go` (unrouted mock) — delete or route. Contains a hardcoded password check ("admin").
- [ ] **P2.24** `handlers/certificates.go`, `custom_rules.go`, `managed_rules.go`, `deployments.go`, `webhooks.go`, `exceptions.go` — all unrouted mock handlers. Delete or implement.
- [ ] **P2.25** `security_events.go` hardcodes `LIMIT 50` with no pagination (`limit`/`offset`/cursor). Add pagination.
- [ ] **P2.26** `metrics.go` declares `shield_control_request_duration_ms` histogram but never records observations. Process start time reports request start time, not process start.
- [ ] **P2.27** `mask.go` has no test file. Add `mask_test.go`.
- [ ] **P2.28** `iplist/`, `ratelimit/`, `blockpage/` have no direct unit tests. Add tests.
- [ ] **P2.29** `cmd/wazuh-forward/main.go` "syslog" mode writes formatted lines to `stdout` — no TCP/UDP/TLS socket, no Wazuh agent API. Implement actual syslog transport.
- [ ] **P2.30** `wazuh-forward` forwards the entire `raw` event to Wazuh `data` without redaction. Add masking at this hop.
- [ ] **P2.31** `policy-v0.json` uses `additionalProperties: false` at the top level, contradicting ADR-002 D1 (forward-compat: unknown fields must be preserved). Fix to `true`.
- [ ] **P2.32** `event-v0.json` uses nested objects; Go `SecurityEvent` uses flat fields. Reconcile the schema with the code.
- [ ] **P2.33** `tools/schema-gen` only generates `$defs` — no root document types (`Policy`, `Event`). Inline nested objects are mangled to `string`. Required/const/pattern/format are all dropped. Expand coverage.
- [ ] **P2.34** Missing ADR-005: document the Coraza v3 integration decision (architecture, transaction lifecycle, fail-open semantics, rule loading).
- [ ] **P2.35** Missing ADR-007: document the high-volume event transport decision (NATS/Kafka vs Redis Streams, ClickHouse migration).
- [ ] **P2.36** `docs/security/` directory is empty. Add at least a security model overview.
- [ ] **P2.37** `docs/api/openapi-v0.yaml` is significantly out of sync with the actual router. Reconcile: add missing routes, remove non-existent routes, fix path mismatches.
- [x] **P2.38** `release-0.1.md` RBAC claim says "placeholder" — code has working RBAC. Update.
- [ ] **P2.39** More granular test coverage: `tests/conformance/`, `tests/evasions/`, `tests/performance/` are empty. Add at least placeholder readmes or initial test files.
- [ ] **P2.40** `tests/README.md` lists empty dirs but omits dirs with actual content (corpus, replay). Fix.

---

## P3 — Low (cleanup, docs, ops, nice-to-have)

- [ ] **P3.1** `README.md` project status: update from "Sprint 1 in progress" to "Sprint 6 / Phase 3 complete".
- [ ] **P3.2** `AI_SKILL.md` is stale (says Release 0.1 has no blocking, rate limits, IP lists). Update to reflect Phase 3 reality.
- [ ] **P3.3** `phase0_srs.md` is significantly outdated (scope was Release 0.1 proxy-only; code has Phase 2+3 WAF features). Either update or replace with a Phase 3 SRS.
- [ ] **P3.4** `docs/architecture/system-context-dataflow-diagrams.md` — add the waf-engine sidecar to the container diagram and request flow DFD.
- [ ] **P3.5** `.github/workflows/ci.yml` — add `make gen-check` to the `schema` job. Add `make check-i18n` target to the Makefile.
- [ ] **P3.6** `.github/pull_request_template.md` — referenced in `repository-ci-coding-standards.md` but does not exist.
- [ ] **P3.7** `Makefile` `test-failover` target: don't skip API tests — add a note that compose must be up, but don't hide the test.
- [ ] **P3.8** `Makefile` `check-i18n` target is called by CI but doesn't exist. Add a simple i18n key-completeness check.
- [ ] **P3.9** `deployments/ansible/`, `deployments/helm/`, `deployments/appliance/` — empty directories. Add placeholder READMEs or remove.
- [ ] **P3.10** `rules/cve/`, `rules/technology/`, `rules/tests/` — empty directories. Add placeholder READMEs or remove.
- [ ] **P3.11** `docs/ariba-shield-waf-dashboard.html` — standalone UI artifact not referenced anywhere. Move to a `design/` directory or document its purpose.
- [ ] **P3.12** `handlers/policy_versions.go` — `CreatePolicyVersion` returns `{"id":"...","version":N,"status":"..."}` but the `version` field is typed `any` and the `json:"version"` tag is `int`. The response structure is inconsistent.
- [ ] **P3.13** `handlers/policy_versions.go` — `PromotePolicyVersion` transition validation map is hardcoded. Extract to a shared state machine.
- [ ] **P3.14** `handlers/policy_diff.go` — `simpleDiff` does a shallow JSON field comparison. Not suitable for deeply nested policy documents. Consider a structured diff library.
- [ ] **P3.15** `handlers/ip_lists.go` and `rate_limits.go` — hardcode `orgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"`. No multi-tenant scoping.
- [ ] **P3.16** `handlers/security_events.go` — no filtering by application_id, severity, date range, or rule_id. Add filter parameters.
- [ ] **P3.17** `handlers/gateways.go` — `RegisterGateway` uses `ON CONFLICT (id) DO UPDATE` for idempotency, but the `id` is client-supplied. A gateway could register with another's ID. Fix: provision IDs server-side or bind to mTLS cert.
- [ ] **P3.18** `handlers/metrics.go` — `shield_control_request_duration_ms` histogram is declared but never observed. The process start time is the request time, not the actual process start. Fix.
- [ ] **P3.19** `middleware/middleware.go` — request ID falls back to `crypto/rand` hex if no `X-Request-ID` header. Use ULID for consistency.
- [ ] **P3.20** `middleware/rbac.go` — `RolePermissions` is an in-memory map, not DB-backed. Phase 8 requires DB-backed roles with per-tenant scoping.
- [ ] **P3.21** `middleware/auth.go` `DefaultRoutePermissions` — missing entries for `bind`, `versions`, `promote`, `activate`, `rollback`, `diff` routes. They fall through to `system:admin` by default (only Super Admin can use them).
- [ ] **P3.22** `console-web` — login page hardcodes fetch to `http://localhost:8080/api/v1/auth/login` (wrong port, nonexistent endpoint). Fix.
- [ ] **P3.23** `console-web` — Overview page falls back to hardcoded mock data when API is unreachable, showing fabricated apps/gateways as if real. Remove or flag mock data.
- [ ] **P3.24** `console-web` — no policies page, no incidents page, no certificates page, no rules page, no API security page, no bot defense page, no rate limiting page, no IP intelligence page, no exceptions page, no virtual servers page, no backend pools page, no health monitors page, no integrations page, no reports page, no users/RBAC page, no audit log page, no settings page. All 18 items in §9 IA are missing except Overview, Applications, Security Events, and Gateways.
- [ ] **P3.25** `console-web` — hardcoded English strings in `overview/page.tsx` violate §6.10 (no UI text hard-coded). Move to message catalogs.
- [ ] **P3.26** `console-web` — no session check or route protection in `src/middleware.ts` (only locale routing). Any page renders without auth.
- [ ] **P3.27** `console-web` — login page stores token in `localStorage` (violates FR-0.1-042 HTTP-only cookies). Fix.
- [ ] **P3.28** `console-web` — the `nav` message keys reference pages that don't exist.
- [ ] **P3.29** `console-web` — Applications, Gateways, and Security Events pages are unstyled minimal tables. No charts, no real-time, no pagination, no search.
- [ ] **P3.30** `console-web` — `security-events/page.tsx` bypasses the API client (`api.ts`) and does a raw `fetch`. Fix: use the API client.
- [ ] **P3.31** `console-web` — `API_BASE` defaults to `http://127.0.0.1:8443` (HTTP). The API only serves HTTPS. The console cannot reach the API without a proxy.
- [ ] **P3.32** `mask.go` — `MaskSensitiveValue` checks `strings.Contains` substrings. A key named "my_password" is caught, but "password" as a value is never caught. Also, `authorization` is a header, not an ARGS key. Consider expanding.
- [ ] **P3.33** `engine.go` — `TLSProfile.MinVersion` from the policy is ignored; template hardcodes `TLSv1.2 TLSv1.3`. Fix.
- [ ] **P3.34** `engine.go` — `ServerName` in the nginx template is hardcoded to `vs.Name + ".shield.local"`. Should use the actual hostname from the policy.
- [ ] **P3.35** `engine.go` — `Route.Match` semantics: "exact" renders `location /api` which is a prefix match in nginx, not exact. Fix: `location = /api` for exact.
- [ ] **P3.36** `engine.go` — `Route.Priority` is ignored. No path-conflict detection.
- [ ] **P3.37** `engine.go` — `VirtualServerID` and `ApplicationID` in the access log template are hardcoded placeholders (`"vs"`, `"app"`). Should be the real IDs from the policy.
- [ ] **P3.38** `engine.go` — size limits lose precision via integer division (5.5 MB → `5m`). Fix.
- [ ] **P3.39** `engine.go` — `HealthMonitor` is never rendered in the nginx config. Active health checks do not exist despite being declared in the policy schema.
- [ ] **P3.40** CI `i18n` job: `make check-i18n` target doesn't exist. Add it or remove the CI step.
- [ ] **P3.41** CI `schema` job: runs `make schema-check` without setting up Node.js — `npx tsc --noEmit` will fail unless Node is available. Fix.
- [ ] **P3.42** `.github/workflows/ci.yml` — gosec uses `@master` floating tag. Pin to a specific version.
- [ ] **P3.43** `.github/workflows/ci.yml` — no `pull_request_template.md` exists despite being referenced in coding standards.
- [ ] **P3.44** `docs/architecture/reference-hardware-performance-method.md` — OQ-2 and OQ-4 still open. No `sysctl` tuning snapshot tracked in `deployments/`.
- [ ] **P3.45** `docs/architecture/adr-002-policy-event-schema-v0.md` — D2 says `additionalProperties: false` at top level, contradicting forward-compat requirement. Fix.
- [ ] **P3.46** `docs/operations/release-0.1.md` — capability table says "RBAC ❌" but code has working RBAC middleware. Update.
- [ ] **P3.47** `docs/operations/release-0.1.md` — reference result 35/35 is stale. Update to reflect current test count.
- [ ] **P3.48** `docs/README.md` — does not mention `AI_SKILL.md` or `ariba-shield-waf-dashboard.html`.
- [ ] **P3.49** `AGENTS.md` — does not list `docs/api/openapi-v0.yaml` as a doc to keep updated. Only lists `endpoint.md` (roadmap) and `API Specification Document.md` (aspirational).
- [ ] **P3.50** `AGENTS.md` — "Commands" section does not list `make gen-check`, `make test-replay`, `make test-failover`, `make schema-check`.
- [ ] **P3.51** `Makefile` — `test-failover` target uses `echo "Skipping: run manually"` for API tests. Make them opt-in rather than skipped.