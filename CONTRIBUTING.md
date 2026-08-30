# Contributing to Ariba Shield WAF

First off, thank you for considering contributing to Ariba Shield WAF! It's people like you that make open-source software such a great community. 

This document provides guidelines and instructions for contributing to this project. Please read it carefully to ensure a smooth collaboration process.

---

## 1. Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. We are committed to providing a welcoming and inspiring community for all. Harassment and unacceptable behavior will not be tolerated.

## 2. Core Technologies

Before diving in, you should be familiar with the core stack:
- **Control Plane API:** Go (Golang), PostgreSQL, Redis
- **Data Plane (Gateway):** OpenResty (NGINX), Lua, Coraza WAF Engine
- **Web Console (UI):** Next.js (TypeScript), Tailwind CSS / Vanilla CSS, React

### Project Directory Structure
To help you navigate the codebase, here is a high-level overview of our architecture:
- `/apps/control-api/`: The Go-based backend management API. Out-of-path configuration engine.
- `/apps/console-web/`: The Next.js frontend dashboard.
- `/gateways/openresty-gateway/`: The OpenResty container configuration and Lua scripts for live traffic routing.
- `/docs/`: Project documentation, Architecture Decision Records (ADRs), API schemas, and the User Guide.
- `/packages/`: Shared libraries, generated API clients, and JSON Schemas that act as the single source of truth.
- `/deployments/`: Docker Compose and Kubernetes manifests for standing up the stack.

## 3. Getting Started

### Prerequisites
- [Docker](https://www.docker.com/) and Docker Compose
- [Go 1.22+](https://golang.org/dl/)
- [Node.js 22+](https://nodejs.org/)
- Make

### Local Development Environment
1. **Fork and Clone the Repository:**
   ```bash
   git clone https://github.com/your-org/ariba-shield-waf.git
   cd ariba-shield-waf
   ```

2. **Run the Development Stack:**
   We use Docker Compose to spin up the entire infrastructure locally.
   ```bash
   docker compose -f deployments/compose/docker-compose.yml up -d --build
   ```

3. **Access the Application:**
   - **Console UI:** `http://localhost:3000`
   - **Control API:** `http://localhost:8443`
   - **Default Credentials:** `superadmin@aribashield.local` / `admin`

## 4. Ground Rules & Architecture Policies

To maintain enterprise-grade stability, all contributors **must** follow these rules:

1. **Architecture Changes Require an ADR:**
   Anything touching a trust boundary, data flow, data store, or network protocol requires a new Architecture Decision Record (ADR) in `docs/architecture/`. Do not submit PRs for architectural changes without an approved ADR.
2. **Contracts are Versioned:**
   The policy/event JSON Schemas in `packages/*/schema` are the source of truth. SDK types and database models are generated—never hand-edit them.
3. **No Secrets in Code:**
   No default passwords, no keys in source code, and absolutely no bodies/cookies/credentials in logs. Never commit `.env` or certificate materials.
4. **Traffic Path Invariants:**
   - Never block or crash the live traffic path during a log write.
   - Never deploy partially compiled WAF policies.
   - Never add a network protocol without rigorous conformance, fuzzing, and resource-limit tests.
5. **Claims Must Be Tested:**
   Do not claim production or enterprise-class capabilities in code, docs, or commits without evidence from our QA program.

## 5. Development Workflow

### Commands and Makefile
We use `make` for standardization. Before submitting any changes, ensure you run the appropriate targets:
- `make lint` - Runs linters for Go, TypeScript, and Lua.
- `make test` - Runs unit tests.
- `make build` - Builds all binaries and frontend assets.
- `make check-i18n` - Verifies English (`en`) and Bengali (`bn`) message catalogs match in the frontend.

### Branching Strategy
1. Create a branch from `main`. Name it descriptively: `feature/your-feature`, `bugfix/issue-description`, or `docs/update-guide`.
2. Keep your changes focused. Do not mix refactoring with new feature additions in the same PR.

### Committing Changes
Every finished task must be committed to git with a clear, descriptive commit message relating to the task. 
- Use imperative mood: "Fix port collision in gateway" instead of "Fixed port collision".
- Reference any open issues.

## 6. Submitting a Pull Request

1. Ensure your code strictly adheres to the linters (`make lint`).
2. Verify all tests pass (`make test`).
3. Update relevant documentation (see section 7).
4. Push your branch to your fork and submit a PR against the `main` branch.
5. Provide a clear and detailed PR description explaining the *why* and *how* of your changes.

## 7. Documentation Requirements

If you add a feature or change existing behavior, you **must** update the relevant documentation:
- `docs/phase0_srs.md`: The day-to-day contract for features and use cases.
- `docs/api/endpoint.md`: The API master plan.
- `docs/api/openapi-v0.yaml`: Must be kept in sync with the Go router (`apps/control-api/internal/api/router.go`).
- `docs/user_guide.html`: End-user facing documentation.
- `docs/architecture/`: Update diagrams with any boundary/data-flow changes.

## 8. Getting Help

If you have any questions about the architecture, or need help getting the project running, feel free to open a Discussion on GitHub or reach out to the core maintainers. 

We look forward to your contributions!
