# Enterprise WAF - Complete API Architecture & Master Plan

This document outlines the full 37-module REST API structure required to elevate the **Ariba Shield WAF** into a world-class Control Plane + Data Plane + Security Analytics Platform.

## Phase 1 — Core WAF (Foundation)
1. **Applications** (`/api/v1/applications`)
2. **Domains** (`/api/v1/applications/{id}/domains`)
3. **Origins/Upstreams** (`/api/v1/applications/{id}/origins`)
4. **Security Policies** (`/api/v1/security-policies`)
5. **Managed Rules** (`/api/v1/managed-rules`)
6. **Custom Rules** (`/api/v1/custom-rules`)
7. **Traffic Control** (`/api/v1/traffic-control` - Rate Limiting, Geo, IP Lists)
8. **Security Events & Investigation** (`/api/v1/security-events`)

## Phase 2 — Production Ready (Operations & Safety)
9. **Gateways & Nodes** (`/api/v1/gateways`)
10. **Deployments** (`/api/v1/deployments`)
11. **Policy Versioning** (`/api/v1/policy-versions`)
12. **Rollback Management** (`/api/v1/security-policies/{id}/rollback`)
13. **Configuration Validation & Dry-Run** (`/api/v1/config-validation`)
14. **Exceptions & False Positives** (`/api/v1/exceptions`)
15. **Audit Logs & Compliance** (`/api/v1/audit-logs`)
16. **Certificates & TLS** (`/api/v1/certificates`)

## Phase 3 — Enterprise Security (WAAP Capabilities)
17. **Threat Intelligence** (`/api/v1/threat-intelligence` - IOCs & Feeds)
18. **API Security** (`/api/v1/api-security` - Schema Validation, Discovery)
19. **Bot Management** (`/api/v1/bot-management` - Captcha, JS Challenges)
20. **Data Loss Prevention (DLP)** (`/api/v1/dlp` - Masking SSN/CC)
21. **SIEM & Log Forwarding** (`/api/v1/integrations`)
22. **IAM, RBAC & SSO** (`/api/v1/iam`)
23. **Service Accounts & API Keys** (`/api/v1/service-accounts`, `/api/v1/api-keys`)
24. **Secrets Management** (`/api/v1/secrets`)

## Phase 4 — Advanced Platform (Scale & Analytics)
25. **Security & Threat Analytics** (`/api/v1/analytics`)
26. **Rule Hit Analytics** (`/api/v1/rule-analytics`)
27. **Attack Response Automation** (`/api/v1/automation`)
28. **Incident Response & Webhooks** (`/api/v1/incidents`)
29. **Backup & Disaster Recovery** (`/api/v1/backups`)
30. **Multi-Tenancy** (`/api/v1/organizations`, `/api/v1/workspaces`)
31. **HA & Cluster Management** (`/api/v1/clusters`)
32. **Edge Caching & CDN** (`/api/v1/caching`)

## Phase 5 — Next-Gen WAAP (Newly Added)
33. **GraphQL & gRPC Security** (`/api/v1/graphql-security`)
    - Deep query inspection, query depth limiting, and introspection control.
34. **Client-Side Protection (Magecart/CSP)** (`/api/v1/client-side-protection`)
    - Injecting and managing Content Security Policies (CSP) to stop in-browser JavaScript tampering.
35. **Advanced API Quotas & JWT Claims** (`/api/v1/api-quotas`)
    - Rate limiting not just by IP, but by JWT consumer ID or API token tiers.
36. **Machine Learning / Behavioral Baselines** (`/api/v1/ml-baselines`)
    - Train AI on normal traffic patterns to automatically detect zero-day anomalies without signatures.
37. **L3/L4 DDoS Mitigation** (`/api/v1/network-protection`)
    - BGP route injection configs and Anycast routing status for volumetric attacks.

## Phase 6 — Learning and policy builder (Months 18–22)

38. **Learning sessions** (`/api/v1/learning/sessions`)
    - POST — Create a new learning session from a trusted source.
    - GET — List sessions with filters (source, status).
    - GET /{id} — Retrieve a single session.

39. **Learning suggestions** (`/api/v1/learning/suggestions`)
    - GET — List suggestions filtered by status, rule_id, session_id.
    - GET /{id} — Retrieve a single suggestion.
    - POST /{id}/accept — Accept a suggestion (gates policy update via approval workflow; exit criteria: learning cannot directly weaken policy without configured approval).
    - POST /{id}/reject — Reject a suggestion.


41. **Roles** (`/api/v1/roles`)
    - GET — List all roles with their permission sets.
    - GET /{id} — Retrieve a single role by ID.

42. **Permissions** (`/api/v1/permissions`)
    - GET — Derive permissions from the role system. Each role carries a permissions[] array; this endpoint lists all roles and their associated permissions.

## Phase 4 — Reports (Compliance & Analytics)

43. **Reports** (`/api/v1/reports`)
    - GET — List all generated reports (filter by kind, status).
    - POST — Create a report from a kind + params in the body.
    - GET /{id} — Retrieve a single report with its summary data.
    - DELETE /{id} — Delete a report.
    - POST /security — Generate a security events report (severity distribution).
    - POST /traffic — Generate a traffic summary report (method, count, avg latency).
    - POST /incidents — Generate an incidents report (severity distribution).
    - POST /compliance — Generate a compliance summary report (audit event volume, range).
    - GET /{id}/download — Download the full report as JSON attachment.

## Phase 4 — Dashboard (Overview Widgets)

44. **Dashboard** (`/api/v1/dashboard`)
    - GET /overview — High-level counts (events, blocked, requests, applications, gateways, active incidents).
    - GET /traffic — Request volume, avg/p99 latency, status-code distribution.
    - GET /security — Security event volume, blocked count, unique IPs, severity distribution.
    - GET /attacks — Top attack types by reason in the period.
    - GET /top-ips — Top client IPs by event volume and blocked count.
    - GET /top-rules — Top rules by hit count.
    - GET /applications — Per-application request/event/blocked counts.
    - GET /gateways — Gateway fleet status (total, active, offline, detail list).

    All widgets accept an optional `days` query param (default 7, max 90).

## Phase 8 — System Settings & License

45. **System Settings** (`/api/v1/settings`)
    - GET — Retrieve all settings grouped by category.
    - PATCH — Update general settings (upsert key/value pairs).
    - GET/PATCH /security — Security settings (session timeout, lockout, etc.).
    - GET/PATCH /localization — Localization settings (default language, locale).
    - GET/PATCH /retention — Retention settings (log retention days, event TTL).

46. **License / Entitlements** (`/api/v1/license`)
    - GET — Retrieve the current active license.
    - POST /activate — Activate a license key (sets edition, seats, limits, expiry).
    - POST /deactivate — Deactivate the current license.
    - GET /usage — Usage against license limits (gateways, applications).
    - GET /entitlements — Feature set granted by the license edition.
