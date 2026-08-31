# Architecture Decision Record (ADR-009): System Settings Management

## 1. Context & Motivation

Ariba Shield WAF requires a flexible mechanism to store, manage, and propagate global configurations across the platform. Hardcoding settings in environment variables or configuration files requires restarts and doesn't scale for multi-tenant environments. 

We need a structured way to handle settings like **Security Policies (MFA requirements)**, **Data Retention Rules**, and **Platform Localization**, ensuring these settings are dynamically enforced by backend services without downtime.

## 2. Objective
Design and document the **System Settings** module—defining the schema, CRUD lifecycle, caching mechanisms, and how the values are consumed by background jobs (e.g., cron jobs for data retention).

## 3. Database Architecture

Settings are stored in the PostgreSQL `system_settings` table in a flat, key-value structure categorized by domains.

### Schema (`system_settings`)
* `id` (ULID)
* `organization_id` (String - isolates settings in multi-tenant mode)
* `category` (String - e.g., `general`, `security`, `localization`, `retention`)
* `key` (String - unique within a category/org)
* `value` (JSONB - holds strings, numbers, or booleans)
* `updated_at` (Timestamp)

## 4. Setting Categories & Core Keys

### A. General Settings (`general`)
* `site_name`: Display name of the platform.
* `admin_email`: Fallback email for critical system alerts.
* `maintenance_mode`: If `true`, the `control-api` returns HTTP 503 for all non-admin write requests.

### B. Security Settings (`security`)
* `max_login_attempts`: Brute-force threshold (default 5).
* `session_timeout_minutes`: JWT/Cookie expiration for the console web.
* `mfa_required`: If `true`, any user without MFA will be forced to configure it before accessing the dashboard.

### C. Localization (`localization`)
* `default_language`: `en` or `bn`. Dictates the fallback language for system emails and new user accounts.
* `timezone`: Standard timezone for log timestamps (default UTC).

### D. Data Retention (`retention`)
* `audit_log_days`: Number of days to retain `audit_logs` before cleanup.
* `security_event_days`: Number of days to retain raw `security_events` from the data plane.

## 5. API & UI Lifecycle

### Backend (`control-api`)
The API exposes `GET` and `PATCH` endpoints organized by category:
* `GET /api/v1/settings` (Fetches all, grouped by category)
* `PATCH /api/v1/settings/{category}` (Bulk upserts keys for a specific category)

All `PATCH` operations are logged to `audit_logs` to maintain an immutable record of system changes.

### Frontend (`console-web`)
The `SettingsPage` dynamically renders form fields based on a strongly-typed `SETTINGS_SCHEMA`. This allows developers to add new settings simply by updating the schema object in React, without needing new UI components.

## 6. Propagation & Consumption (Next Steps)

Currently, settings are stored and retrieved successfully. The roadmap for enforcement:

1. **In-Memory Cache (Redis or Go maps):** Read-heavy settings like `maintenance_mode` or `mfa_required` will be cached in the API middleware to prevent database hits on every request.
2. **Cron Jobs (Data Pruning):** The `event-ingestor` will query `retention.security_event_days` daily to execute `DELETE FROM security_events WHERE timestamp < NOW() - INTERVAL 'X days'`.
3. **Audit Log Pruning:** A separate worker will query `retention.audit_log_days` to clean up the `audit_logs` table.
