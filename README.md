# Ariba Shield WAF

Centralized, Linux-hosted, web-managed enterprise Web Application Firewall.
See the [master plan](docs/enterprise_waf_development_master_plan.md) for scope, the
[Phase 0 SRS](docs/phase0_srs.md) as the day-to-day development contract, and
[ADR-001](docs/architecture/adr-001-initial-architecture.md) for the architecture.

**Status:** Phase 0 documentation complete; Sprint 1 (repo + CI) in progress. This is a monorepo.

## Repository layout

| Path | Purpose |
|---|---|
| `apps/console-web/` | Next.js + TypeScript management console (control plane UI) |
| `apps/control-api/` | Go control-plane API (management REST/JSON, OpenAPI) |
| `gateways/openresty-gateway/` | Release 0.1 data plane (OpenResty/Nginx) |
| `gateways/rust-gateway/` | Phase 9 custom high-performance data plane (placeholder) |
| `services/` | control-plane services (policy-compiler, event-ingestor, ...) |
| `packages/` | shared schema, localization, generated SDKs |
| `rules/` | signature source (Phase 2+, placeholder) |
| `deployments/` | compose / ansible / helm / appliance |
| `tests/` | conformance, evasions, performance, failover |
| `docs/` | architecture ADRs + diagrams, operations, security, API |

## Core invariants (master plan §22)

1. Never parse the same request differently in security and proxy layers.
2. Never let logging failure block or crash the live traffic path.
3. Never deploy partially compiled policy.
4. Never trust automatically learned traffic without poisoning controls.
5. Never store unmasked secrets merely for analyst convenience.
6. Never make AI output a direct production block rule without deterministic validation and human approval.
7. Never publish performance without exact test conditions.
8. Never add a protocol until it has conformance, fuzz, evasion, and resource-limit tests.
9. Never update all gateways simultaneously; use canary and rollback.
10. Never call the product enterprise/F5-grade until independent testing supports the claim.

## Getting started

See `deployments/compose/` (Sprint 2) for the development environment.

## Development

```sh
make lint      # lint all languages
make test      # run all tests
make build     # build all artifacts
make gen       # regenerate SDK types from schema
```

See [repository-ci-coding-standards.md](docs/architecture/repository-ci-coding-standards.md)
for coding standards and the ADR template.

## RBAC Test Users (For Local UI Testing)

Use the following credentials in the local development environment (`http://localhost:3002/login`) to test the Role-Based Access Control (RBAC) UI features:

| Role | Email | Password | Allowed Access |
|---|---|---|---|
| **Super Admin** | `superadmin@aribashield.local` | `admin` | Full access to all settings, global policies, and tenant management. |
| **Platform Admin** | `platform@aribashield.local` | `admin` | Gateway nodes, load balancers, deployments, and cluster health. |
| **Security Admin** | `security@aribashield.local` | `admin` | Policy creation, custom rules, WAF configurations, and exceptions. |
| **App Owner** | `appowner@aribashield.local` | `admin` | Only their assigned Applications and related traffic logs. |
| **SOC Analyst** | `soc@aribashield.local` | `admin` | Security Events, Incident Playbooks, Analytics, and Webhooks. |
| **Auditor** | `auditor@aribashield.local` | `admin` | Read-only access to Audit Logs, Reports, and System Configurations. |
| **Read Only** | `readonly@aribashield.local` | `admin` | Read-only view of the main dashboard and live traffic. |

*(Note: In production, these mock passwords must be replaced, and MFA must be enabled according to Phase 3 requirements.)*
