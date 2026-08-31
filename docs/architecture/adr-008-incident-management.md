# Architecture Decision Record (ADR-008): Incident Management

## 1. Context & Motivation

As a web application firewall (WAF) platform, Ariba Shield generates a substantial amount of security events. To prevent alert fatigue and enable SOC analysts to triage critical threats effectively, raw `security_events` must be correlated into actionable **Incidents**. 

Currently, basic CRUD operations for incidents exist in the Control API, and a frontend dashboard page has been built. However, a formalized master plan is required to dictate how incidents are managed, automatically escalated, and integrated with external alerting channels.

## 2. Objective
Establish a structured Incident Management Lifecycle (Creation, Triage, Escalation, Resolution, and Audit) to govern how security anomalies are tracked and mitigated.

## 3. Incident Lifecycle Workflow

An incident transitions through the following states:
1. **Open:** A new incident is created manually by a SOC Analyst or automatically via anomaly detection rules.
2. **Investigating:** An analyst is assigned to the incident and is actively reviewing related security events and traffic logs.
3. **Resolved:** The root cause is identified, false positives are tuned, or mitigation (like blocking the IP) has been applied.
4. **Closed:** The incident is archived.

### Transitions & Triggers
* **Auto-Escalation:** If an incident has the severity `critical` and remains `open` for more than 15 minutes, it will automatically escalate and notify on-call engineers.
* **Reopen:** If the same attack vector is detected within 24 hours of an incident being `resolved`, it will automatically transition back to `open` rather than creating a duplicate incident.

## 4. Core Features & Roadmap

### Phase 1: Basic Management (Current State)
* **CRUD Endpoints:** `control-api` has endpoints to list, create, edit, delete, and view incident timelines.
* **SOC UI:** The React dashboard (`/incidents`) supports viewing incidents, filtering by status/severity, manual creation, and lifecycle actions (escalate, close, reopen).

### Phase 2: Automatic Correlation
* **Event Grouping:** The `event-ingestor` service will analyze incoming Coraza WAF logs. If multiple blocks (e.g., SQLi) occur against the same `application_id` from the same IP within a 5-minute window, an incident is automatically generated.
* **Timeline Association:** New related events will be appended to the `related_events` array of the active incident.

### Phase 3: External Integrations & Webhooks
* **Slack / Teams Alerts:** Sending an immediate webhook notification for any `high` or `critical` incident.
* **PagerDuty Integration:** Escalating an incident via API calls to on-call incident management tools.

### Phase 4: AI-Assisted Triage
* **Sherlog AI Integration:** Introduce an LLM-based summarization of the `related_events`. The AI will provide a natural language summary of the attack vector, suggesting potential WAF tuning or False Positive exceptions.

## 5. Security & Access Control (RBAC)

Access to incidents is governed by strict RBAC constraints:
* **SOC Analysts & Security Admins:** Full read/write access (create, assign, escalate, close).
* **App Owners:** Read-only access restricted strictly to incidents that impact their assigned applications.
* **Auditors:** Read-only global access.

## 6. Database Schema Summary

The `incidents` table holds the aggregate context:
* `id` (ULID)
* `organization_id` (String)
* `title` (String)
* `severity` (Enum: low, medium, high, critical)
* `status` (Enum: open, investigating, resolved, closed)
* `owner_user_id` (String - Assigned analyst)
* `related_events` (JSONB Array of Event IDs)
* `notes` (String)

*Note: All actions (close, reopen, escalate) also log to the immutable `audit_logs` table for compliance tracking.*
