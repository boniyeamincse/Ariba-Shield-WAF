<div align="center">
  <h1>🛡️ Ariba Shield WAF</h1>
  <p><b>Centralized, high-performance, enterprise-grade Web Application Firewall (WAF)</b></p>

  <a href="https://github.com/your-org/ariba-shield-waf/releases"><img src="https://img.shields.io/github/v/release/your-org/ariba-shield-waf?style=for-the-badge&color=blue" alt="Release"></a>
  <a href="https://github.com/your-org/ariba-shield-waf/blob/main/LICENSE"><img src="https://img.shields.io/github/license/your-org/ariba-shield-waf?style=for-the-badge" alt="License"></a>
  <a href="CONTRIBUTING.md"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge" alt="PRs Welcome"></a>
</div>

<br>

Ariba Shield is an open-source, next-generation Web Application Firewall designed for modern enterprise infrastructure. It isolates the control plane from the data plane, ensuring that your traffic routing is never slowed down by management APIs or logging bottlenecks.

Powered by a robust **Go** control plane, an **OpenResty (Coraza)** data plane, and a sleek **Next.js** glassmorphism dashboard.

## ✨ Key Features

- **Advanced WAF Engine:** Powered by Coraza and the OWASP Core Rule Set (CRS) to mitigate SQLi, XSS, and LFI.
- **Zero-Downtime Deployments:** Atomic configuration updates with bundle hashing and automatic rollback on failure.
- **Out-of-Path Control Plane:** Traffic routing is 100% physically isolated from the management APIs.
- **Granular RBAC:** 7-tier role system (Super Admin to Read Only) natively built-in.
- **Compliance Ready:** Immutable audit trails and automatic sensitive-field masking in security events.
- **Premium Dashboard:** Bilingual (English/Bengali) Next.js interface with real-time metric visualization.

---

## 🚀 Quick Start

The fastest way to test Ariba Shield is by running the local development stack via Docker Compose.

```bash
# 1. Clone the repository
git clone https://github.com/your-org/ariba-shield-waf.git
cd ariba-shield-waf

# 2. Spin up the infrastructure
docker compose -f deployments/compose/docker-compose.yml up -d --build
```

### Accessing the Dashboard

Once the containers are running, open your browser:
- **URL:** `http://localhost:3000` (or `http://<your-local-ip>:3000`)
- **Login Email:** `superadmin@aribashield.local`
- **Password:** `admin`

*Note: In production, these mock passwords must be replaced and MFA must be enabled.*

---

## 📁 Repository Architecture

This is a monorepo containing both the control plane and data plane.

| Path | Purpose |
|---|---|
| `/apps/console-web/` | Next.js + TypeScript management console (UI) |
| `/apps/control-api/` | Go control-plane API (management REST/JSON, OpenAPI) |
| `/gateways/openresty-gateway/` | OpenResty/Nginx data plane with Lua routing |
| `/packages/` | Shared schemas, localization, and generated SDKs |
| `/deployments/` | Docker Compose, Ansible, and Kubernetes manifests |
| `/docs/` | Architecture ADRs, SRS, API schemas, and User Guide |

---

## 🔒 10 Core Invariants

To maintain enterprise-grade resilience, this project strictly adheres to 10 golden rules (defined in our Master Plan):
1. **Never** parse the same request differently in security and proxy layers.
2. **Never** let logging failure block or crash the live traffic path.
3. **Never** deploy partially compiled policy.
4. **Never** trust automatically learned traffic without poisoning controls.
5. **Never** store unmasked secrets merely for analyst convenience.
6. **Never** make AI output a direct production block rule without deterministic validation.
7. **Never** publish performance without exact test conditions.
8. **Never** add a protocol until it has conformance, fuzz, and resource-limit tests.
9. **Never** update all gateways simultaneously; use canary and rollback.
10. **Never** call the product enterprise/F5-grade until independent testing supports the claim.

---

## 👥 Role-Based Access Control (RBAC)

Ariba Shield uses a strict RBAC matrix. If you want to test different views locally, you can use these built-in test accounts:

| Role | Email | Password | Access Level |
|---|---|---|---|
| **Super Admin** | `superadmin@aribashield.local` | `admin` | Full system access. |
| **Platform Admin** | `platform@aribashield.local` | `admin` | Gateway nodes and load balancers. |
| **Security Admin** | `security@aribashield.local` | `admin` | WAF rules, policies, exceptions. |
| **App Owner** | `appowner@aribashield.local` | `admin` | Assigned apps and traffic logs. |
| **SOC Analyst** | `soc@aribashield.local` | `admin` | Events, Analytics, Webhooks. |
| **Auditor** | `auditor@aribashield.local` | `admin` | Read-only Audit logs & config state. |
| **Read Only** | `readonly@aribashield.local` | `admin` | Dashboard visualization only. |

---

## 🤝 Contributing

We welcome contributions from the community! Because of the strict security nature of a WAF, please read our [Contribution Guidelines (CONTRIBUTING.md)](CONTRIBUTING.md) before submitting a Pull Request or Architecture Decision Record (ADR).

**Development Commands:**
```sh
make lint      # lint all languages (Go, TS, Lua)
make test      # run all unit tests
make build     # build all artifacts
make gen       # regenerate SDK types from JSON schemas
```

## 📚 Documentation

For deep-dives into the architecture and usage, check out:
- [User Guide](docs/user_guide.html): End-user HTML manual for navigating the dashboard and onboarding applications.
- [API Documentation](docs/api/endpoint.md): Comprehensive list of our REST API endpoints.
- [Architecture Details](docs/architecture/): Deep dives into the WAF's isolated control-plane architecture.

## 📄 License

This project is open-sourced software licensed under the MIT license.
