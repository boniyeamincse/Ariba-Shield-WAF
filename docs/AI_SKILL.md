# AI Skill / Project Context

## Project name
**Ariba Shield WAF**

## Project type
Centralized enterprise **Web Application Firewall (WAF)** platform for Linux.

## Purpose
This project is building a centralized platform that protects multiple web applications and APIs through managed gateway instances. The architecture separates:

- a **data plane** that handles live traffic
- a **control plane** that manages configuration, policies, applications, users, logs, and lifecycle operations

The project is designed to grow incrementally from a safe reverse-proxy platform into a full enterprise WAF.

---

## Current development phase
The project documentation currently defines:

1. A **master roadmap** in `enterprise_waf_development_master_plan.md`
2. A **day-to-day implementation contract for Release 0.1 (Lab)** in `phase0_srs.md`

### Current practical scope
The active near-term build target is **Release 0.1 (Lab)**.

Release 0.1 is a **working WAF** with:

- a single-node Linux lab deployment
- a centralized management UI/API (300+ API routes, 23 console routes)
- Coraza WAF engine (detection-only, baseline CRS rules for SQLi, XSS, cmd, traversal, LFI)
- security events pipeline (engine → ingestor → PostgreSQL)
- sensitive-field masking (rule 5)
- an HTTPS reverse proxy gateway
- application and backend pool management
- active/passive health checks
- safe config validation and atomic activation
- last-known-good rollback/reload behavior
- **rate limiting** and **IP lists** (Phase 3)
- **RBAC** (7 roles, session-cookie auth, MFA)
- **incident management** (Phase 4)
- **dashboard** (8 widgets: overview, traffic, security, attacks, top-IPs, top-rules, applications, gateways)
- **learning engine** (baseline sessions and suggestions)
- structured request logging
- English/Bangla UI shell

---

## Architecture summary

### 1. Data plane
Responsible for:

- TLS termination
- HTTP parsing and normalization
- reverse proxying
- routing to backend pools
- health checks
- later: inspection, WAF enforcement, rate limiting, bot defense, API protection

### 2. Control plane
Responsible for:

- applications
- virtual servers
- backend pools and nodes
- policies and versions
- certificates and secret references
- users, auth, RBAC
- audit events
- telemetry and reporting
- update and lifecycle operations

### Critical architecture rule
The **control plane must never be in the live traffic path**.

Gateways must continue serving traffic with a **signed, versioned, locally cached last-known-good configuration** if the control plane becomes unavailable.

---

## Confirmed technology direction

### Release 0.1 / early platform
- **Gateway/Data plane:** OpenResty / Nginx
- **Frontend/UI:** Next.js
- **Backend API:** Go
- **Database:** PostgreSQL
- **Metrics/observability:** Prometheus + structured JSON logs + OpenTelemetry-compatible pipeline

### Future direction
- A custom Rust data plane may be considered later, but only after semantics, tests, baselines, and operations are mature.

---

## Important scope rules

### In scope now
- Reverse proxy behavior
- Virtual server and backend pool configuration
- Health monitors
- Basic authentication
- Role-aware authorization foundation
- Config generation/validation/activation/rollback
- Auditability
- Bilingual UI foundation

### Explicitly out of scope for Release 0.1
- Full OWASP CRS integration (baseline rules only)
- Blocking WAF enforcement (detection-only in 0.1)
- Multi-gateway HA control
- Bot defense (policies scaffolded)
- API schema protection (policies scaffolded)
- Full learning engine (baseline scaffolded)
- Production blocking posture
- Custom TLS implementation or cryptography — always use maintained system libraries
- L3/L4 volumetric DDoS mitigation

---

## Product goals from the master plan
The long-term platform should support:

- centralized protection for multiple applications and APIs
- Linux, VM, container, and Kubernetes deployment models
- HTTP/1.1, HTTP/2, WebSocket, REST, GraphQL, JSON, XML, multipart, and file uploads
- negative and positive security models
- anomaly scoring
- API protection
- rate limiting
- bot and abuse defense
- centralized multi-gateway management
- SIEM/SOAR/Wazuh/syslog/webhook integrations
- English and Bangla UI with extensible i18n

---

## Non-goals / safety rules
These are important constraints from the docs:

- Do **not** build a custom TLS implementation.
- Do **not** invent cryptography.
- Do **not** claim F5-equivalent capability based only on feature count.
- Do **not** use generative AI as the blocking authority.
- Do **not** send full sensitive request bodies to cloud AI services.
- Do **not** depend on successful log writing to make traffic decisions.

---

## Security expectations
The management platform itself must be secure.

Expected controls include:

- no default passwords
- RBAC from the beginning
- MFA support
- audit logging for all config mutations
- HTTPS-only management API
- encrypted secrets at rest
- signed config/rule bundles
- protected software supply chain
- branch/review/scanning/release controls

---

## Data model concepts
Core entities mentioned in the docs include:

- organizations
- users
- groups
- roles
- sessions
- audit events
- applications
- virtual servers
- routes
- backend pools
- backend nodes
- health monitors
- config versions
- config deployments

### Identifier conventions
- Public IDs should use **ULID/UUID-style identifiers**
- Write models should support **optimistic concurrency**
- Config deployments are versioned and hash-based

---

## API conventions
Documented API conventions for Release 0.1:

- Base path: `/api/v1`
- HTTPS only
- authenticated access
- OpenAPI-documented
- write endpoints accept `Idempotency-Key`
- responses use version/ETag semantics for optimistic concurrency

Main endpoint groups include:

- auth
- current user/profile
- applications
- virtual servers
- backend pools
- backend nodes
- health monitors
- config versions / activate / rollback
- traffic request queries
- gateways
- audit events
- metrics
- health

---

## Deployment model for Release 0.1
Single-node lab environment:

- control plane and gateway on one Linux host
- management traffic on management network
- gateway exposed for application/service traffic
- gateway keeps last-known-good config locally
- logs and metrics exported asynchronously

---

## Development priorities
Based on the SRS sprint plan, the current work should generally align to:

1. repository/CI/foundation
2. Docker Compose development environment
3. PostgreSQL schema and migrations
4. Go API skeleton + OpenAPI
5. Next.js login/layout/i18n shell
6. OpenResty gateway image and config generation
7. health checks and safe reload/rollback
8. management UI for apps/backends
9. observability, audit logs, idempotency, concurrency
10. installer, upgrade, rollback, and release validation

---

## What an AI coding agent should understand before changing code

### 1. This project is phase-driven
Do not implement future-phase enterprise WAF features inside Release 0.1 work unless explicitly requested.

### 2. Safety is a first-class requirement
Config validation, atomic activation, rollback behavior, and operational safety matter as much as features.

### 3. The control plane is not the traffic path
Do not redesign features in a way that makes management-plane availability required for request handling.

### 4. Follow the documented stack
Only use technologies already aligned with the docs unless the repository clearly shows a different implemented decision.

### 5. Internationalization matters from the beginning
UI strings should not be hard-coded if the app is following the documented bilingual architecture.

### 6. Auditability matters
Config-changing actions should be designed with audit logging and traceability.

### 7. Performance claims must be measured
Do not state throughput/latency claims without context or benchmark basis.

---

## Recommended contributor rules for AI agents
- Read `enterprise_waf_development_master_plan.md` before major architecture changes.
- Read `phase0_srs.md` before implementing Release 0.1 features.
- Prefer incremental, testable changes.
- Preserve last-known-good and rollback safety patterns.
- Avoid introducing hidden defaults that weaken security.
- Keep machine-readable codes stable even if user-facing labels are translated.
- Match naming and API contracts to the docs unless the codebase has an approved newer pattern.

---

## Source documents
- `/home/boni/Desktop/WAF/docs/enterprise_waf_development_master_plan.md`
- `/home/boni/Desktop/WAF/docs/phase0_srs.md`

This file is an AI-oriented summary and must not override the source documents.