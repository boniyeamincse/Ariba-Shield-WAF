# Implementation Plan: Managed Rulesets (OWASP CRS)

This document outlines the step-by-step implementation plan for integrating Managed Rulesets (OWASP Core Rule Set - CRS) into the Ariba Shield WAF.

## 1. Database & Backend API Updates
- **Current State:** A basic `managed_rules` table and `managed_rules.go` handler exists but is currently mocked.
- **Action Items:**
  1. Update `managed_rules` schema to include `paranoia_level` (Sensitivity: Low, Medium, High, Strict) and `action` (Block, Log).
  2. Create a database seeder script to populate default OWASP CRS categories (e.g., SQLi, XSS, LFI, RFI, RCE, Scanner Detection, Protocol Validation).
  3. Ensure `ListManagedRules` and `ConfigureManagedRules` APIs are registered in `apps/control-api/internal/api/router.go`.
  4. Implement `GET /api/v1/rules/managed` and `PATCH /api/v1/rules/managed/:id`.

## 2. Frontend UI Development
- **New Page (`/rules/managed`):** Create a dedicated page for configuring Managed Rules.
  - Display a toggle to enable/disable the entire OWASP Core Rule Set.
  - Display a list of CRS rule groups (SQLi, XSS, etc.) with individual toggle switches.
  - Add a dropdown for Global Paranoia Level (Sensitivity) tuning.
  - Add an Anomaly Threshold configuration (Score-based blocking).
- **Navigation:** Update the Sidebar or Rules listing page to add a tab/button linking to "Managed Rules (CRS)".

## 3. Data-Plane & WAF Engine Integration
- **Concept:** Managed rules typically rely on Coraza WAF's built-in CRS support.
- **Action Items:**
  1. The `policy-compiler` service needs to read the `managed_rules` state from the database.
  2. If OWASP CRS is enabled, the compiler should inject the `Include @owasp_crs/*.conf` directive into the generated `shield.conf` or Coraza policy.
  3. Map the Sensitivity (Paranoia Level) from the DB to Coraza `SecAction "id:900000,phase:1,nolog,pass,t:none,setvar:tx.paranoia_level=X"`.
  4. Map the Anomaly Threshold to `tx.inbound_anomaly_score_threshold`.

## 4. Third-Party Heuristics Integration (Future / AI Note)
- **Goal:** Expand Managed Rules beyond standard OWASP CRS.
- **Action:** Integrate ModSecurity OSINT rules and heuristics from `https://github.com/w8mej/WAFRulesHeuristics`.
- **Implementation Strategy:**
  1. Add a new category in `managed_rules` table (e.g., `osint-heuristics`, `comodo-rules`).
  2. Map these `.conf` files to the `policy-compiler` so that they can be activated via the UI.
  3. Ensure Coraza engine compatibility for any specific SecRule variations found in the heuristics repo.

## 5. Testing & Validation
- Write an API test script for configuring Managed Rules.
- Test the WAF Engine against standard OWASP payloads to ensure CRS rules are triggering properly when enabled.
- Verify that False Positives decrease when Paranoia Level is lowered.

---
**Status:** In Progress
