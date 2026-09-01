<div align="center">
  <img src="https://img.icons8.com/color/144/000000/shield.png" alt="Ariba Shield WAF Logo" width="120" />
  <h1>🛡️ Ariba Shield WAF</h1>
  <p><b>Centralized, high-performance, enterprise-grade Web Application Firewall (WAF)</b></p>

  <a href="https://github.com/your-org/ariba-shield-waf/releases"><img src="https://img.shields.io/github/v/release/your-org/ariba-shield-waf?style=for-the-badge&color=blue" alt="Release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=for-the-badge" alt="License"></a>
  <a href="CONTRIBUTING.md"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge" alt="PRs Welcome"></a>
  <a href="SECURITY.md"><img src="https://img.shields.io/badge/Security-Policy-red.svg?style=for-the-badge" alt="Security Policy"></a>
</div>

<br>

**Ariba Shield** is an open-source, next-generation Web Application Firewall designed for modern enterprise infrastructure. It fundamentally isolates the **Control Plane** from the **Data Plane**, ensuring that your live web traffic is never bottlenecked by management APIs, log aggregation, or UI tasks.

Powered by a robust **Go** control plane, an **OpenResty (Coraza)** data plane, and a sleek **Next.js** (Glassmorphism) dashboard, Ariba Shield brings enterprise-level security visibility to everyone.

---

## ✨ Key Features

- **Advanced WAF Engine:** Powered by Coraza and the OWASP Core Rule Set (CRS) to actively mitigate SQLi, XSS, RCE, and LFI.
- **Out-of-Path Control Plane:** Traffic routing is 100% physically isolated from the management APIs. If the control API goes down, the WAF keeps protecting traffic.
- **Dynamic Rule Engine:** Create, test, and deploy modular security policies using our intuitive visual builder.
- **Zero-Downtime Deployments:** Atomic configuration updates with bundle hashing and automatic rollback on failure.
- **Granular RBAC:** 7-tier role system natively built-in (from Super Admin to Read-Only Auditor).
- **Premium Dashboard:** Bilingual (English/Bengali) Next.js interface with real-time metric visualization, incident auto-correlation, and rich analytics.

---

## 🚀 Quick Start

The fastest way to test Ariba Shield is by running the local development stack via Docker Compose.

### 1. Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Node.js](https://nodejs.org/en/) v18+ (for frontend development)
- [Go](https://golang.org/doc/install) 1.22+ (for backend development)

### 2. Spin up the infrastructure
```bash
# Clone the repository
git clone https://github.com/your-org/ariba-shield-waf.git
cd ariba-shield-waf

# Spin up the Data Plane and Control Plane via Docker Compose
docker compose -f infra/compose/docker-compose.yml up -d --build
```

### 3. Accessing the Dashboard

Once the containers are running and healthy, open your browser:
- **URL:** `http://localhost:3000` 
- **Login Email:** `superadmin@aribashield.local`
- **Password:** `admin`

*Note: In production environments, these mock passwords MUST be changed and MFA enabled.*

---

## 📁 Repository Architecture

This is a monorepo containing both the control plane, data plane, and shared schemas.

| Directory | Purpose | Tech Stack |
|---|---|---|
| `apps/console-web/` | The management dashboard (UI) | Next.js, React, TypeScript |
| `apps/control-api/` | The centralized management API | Go, PostgreSQL, Redis |
| `services/waf-engine/`| The actual Data Plane processing HTTP requests | Go, Coraza WAF |
| `services/event-ingestor/`| Asynchronous log processor for security analytics | Go, Kafka/Redis |
| `packages/` | Shared schemas, API contracts, and localization | JSON Schema |
| `docs/` | Architecture ADRs, SRS, API schemas, and Master Plans | Markdown |
| `.github/` | CI/CD pipelines, Issue Templates, and PR templates | YAML |

---

## 🔒 10 Core Architectural Invariants

To maintain enterprise-grade resilience, this project strictly adheres to 10 golden rules:
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

We love contributions! Because of the strict security nature of a WAF, please ensure you review our guidelines before submitting a Pull Request.

1. **[Contribution Guidelines](CONTRIBUTING.md)**: Start here to understand our workflow.
2. **[Code of Conduct](CODE_OF_CONDUCT.md)**: Please adhere to our community standards.
3. **[Security Policy](SECURITY.md)**: Learn how to safely report vulnerabilities.

### Development Commands
```sh
make lint      # lint all languages (Go, TS)
make test      # run all unit tests
make build     # build all artifacts
make gen       # regenerate SDK types from JSON schemas
```

---

## 📚 Documentation

For deep-dives into the architecture and usage, check out:
- **[Architecture (ADRs)](docs/architecture/)**: Read our Architecture Decision Records to understand *why* we built things this way.
- **[API Specs](docs/api/)**: Comprehensive list of our REST API endpoints and schemas.

---

## 🌟 Contributors

A massive thank you to everyone who has contributed to Ariba Shield! 

<a href="https://github.com/your-org/ariba-shield-waf/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=your-org/ariba-shield-waf" alt="Contributors" />
</a>

*Want to see your face here? Check out our [Contributing Guide](CONTRIBUTING.md) to get started!*

---

## 📄 License

This project is licensed under the **Apache 2.0 License**. See the [LICENSE](LICENSE) file for details.
