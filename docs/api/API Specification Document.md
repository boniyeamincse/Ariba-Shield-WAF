# Ariba Shield WAF - Enterprise API Specification Document

This document provides a detailed specification for the proposed Enterprise API endpoints. It defines the payload structures, expected responses, and authentication requirements for managing advanced security features.

## Base URL
`https://control-api.ariba-shield.local/api/v1`

## Authentication
All endpoints require a valid Bearer token (JWT) or session cookie obtained via `/api/v1/auth/login`.

---

## 1. Security Policy & WAF Rules

### 1.1 Create Custom WAF Rule
**Endpoint:** `POST /security-policies/{policy_id}/rules/custom`
**Description:** Adds a custom regex or condition-based firewall rule.

**Request Body (JSON):**
```json
{
  "name": "Block SQLi in Header",
  "description": "Blocks SQL injection attempts in the User-Agent header",
  "action": "BLOCK",
  "match_conditions": [
    {
      "variable": "HEADER:User-Agent",
      "operator": "REGEX_MATCH",
      "value": "(?i)(union.*select|select.*from)"
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "id": "rule_9a8b7c6d",
  "status": "active",
  "created_at": "2026-08-29T10:00:00Z"
}
```

---

### 1.2 Enable Managed Rules (OWASP CRS)
**Endpoint:** `POST /security-policies/{policy_id}/rules/managed`
**Description:** Enable industry-standard managed rule sets (e.g., OWASP Core Rule Set) with configurable paranoia levels.

**Request Body (JSON):**
```json
{
  "rule_set": "OWASP_CRS_V3",
  "paranoia_level": 1,
  "enforcement_mode": "DETECTION",
  "anomaly_score_threshold": 5
}
```
*Note: Best practice is to run at Paranoia Level 1 or 2 in DETECTION mode for the first 3 months to monitor false positives before switching to BLOCKING.*

### 1.3 Create Exception Rule (False Positive Mitigation)
**Endpoint:** `POST /exceptions`
**Description:** Create highly specific exclusion rules to prevent valid traffic from being blocked by strict security rules.

**Request Body (JSON):**
```json
{
  "target_rule_id": "ARIBA-SQLI-001",
  "application_id": "app_internal_hr",
  "match_conditions": [
    {
      "variable": "URI_PATH",
      "operator": "EQUALS",
      "value": "/api/reports/query"
    }
  ],
  "reason": "Internal HR reporting payload contains SQL-like syntax",
  "expiry_days": 90
}
```

### 1.4 Automated Policy Testing / Dry-Run Pipeline
**Endpoint:** `POST /rule-testing/simulate`
**Description:** Simulate a raw HTTP request against a draft policy to ensure it catches attacks without breaking legitimate traffic (Used heavily in CI/CD pipelines).

**Request Body (JSON):**
```json
{
  "policy_id": "policy_draft_002",
  "request": {
    "method": "POST",
    "path": "/login",
    "headers": {
      "User-Agent": "Mozilla/5.0",
      "Content-Type": "application/x-www-form-urlencoded"
    },
    "body": "username=admin' OR 1=1--"
  }
}
```

**Response (200 OK):**
```json
{
  "action": "BLOCK",
  "matched_rules": ["SQLI-001", "OWASP-942100"],
  "anomaly_score": 10
}
```

---

## 2. Traffic Control & Bot Management

### 2.1 Configure Rate Limiting
**Endpoint:** `POST /traffic-control/rate-limits`
**Description:** Apply rate limiting to prevent DDoS or brute-force attacks.

**Request Body (JSON):**
```json
{
  "name": "Global API Rate Limit",
  "target_application_id": "app_12345",
  "limit": 1000,
  "window_seconds": 60,
  "match_by": ["CLIENT_IP"],
  "action": "THROTTLE",
  "response_code": 429
}
```

### 2.2 Configure Geo-Blocking
**Endpoint:** `POST /traffic-control/geo-blocking`
**Description:** Restrict traffic from specific countries.

**Request Body (JSON):**
```json
{
  "policy_id": "pol_9999",
  "blocked_countries": ["XX", "YY"],
  "allowed_countries": [],
  "action": "DROP"
}
```

---

## 3. Certificate Management (mTLS & SSL)

### 3.1 Provision Let's Encrypt Certificate
**Endpoint:** `POST /certificates/acme`
**Description:** Trigger automatic certificate provisioning via ACME.

**Request Body (JSON):**
```json
{
  "domain": "secure.ariba-shield.local",
  "provider": "letsencrypt",
  "email": "admin@ariba-shield.local"
}
```

**Response (202 Accepted):**
```json
{
  "task_id": "task_445566",
  "status": "provisioning"
}
```

---

## 4. SIEM Integration

### 4.1 Configure Splunk Forwarder
**Endpoint:** `POST /integrations/log-forwarders`
**Description:** Forward access and WAF event logs to Splunk HEC.

**Request Body (JSON):**
```json
{
  "type": "SPLUNK_HEC",
  "name": "Enterprise Splunk SOC",
  "endpoint": "https://splunk.internal:8088/services/collector",
  "token": "secret-token-value",
  "log_types": ["WAF_ALERTS", "ACCESS_LOGS", "AUDIT_LOGS"],
  "enabled": true
}
```

---

## 5. Enterprise IAM & SSO

### 5.1 Configure SAML (Okta Integration)
**Endpoint:** `POST /iam/sso/saml`
**Description:** Configure SAML SSO for admin console access.

**Request Body (JSON):**
```json
{
  "provider_name": "Okta",
  "idp_entity_id": "http://www.okta.com/exk123456",
  "sso_url": "https://company.okta.com/app/ariba/sso/saml",
  "x509_certificate": "-----BEGIN CERTIFICATE-----\nMIIB...-----END CERTIFICATE-----"
}
```

---

## 6. Multi-Tenancy

### 6.1 Create New Organization (Tenant)
**Endpoint:** `POST /organizations`
**Description:** Create a new isolated workspace/tenant for multi-tenant deployments.

**Request Body (JSON):**
```json
{
  "name": "Acme Corp",
  "contact_email": "security@acme.corp",
  "plan": "ENTERPRISE",
  "quotas": {
    "max_applications": 50,
    "max_gateways": 10
  }
}
```

---

## 7. Data Privacy & DLP (Data Loss Prevention)

### 7.1 Create DLP Profile
**Endpoint:** `POST /dlp/profiles`
**Description:** Create a profile to mask sensitive data (like credit cards) in outbound HTTP responses.

**Request Body (JSON):**
```json
{
  "name": "PCI-DSS Data Masking",
  "scan_targets": ["RESPONSE_BODY"],
  "rules": [
    {
      "pattern_type": "PREDEFINED",
      "pattern_name": "CREDIT_CARD",
      "action": "MASK",
      "mask_character": "*"
    }
  ]
}
```

---

## 8. Incident Response & Webhooks

### 8.1 Register Alert Webhook
**Endpoint:** `POST /incident-response/webhooks`
**Description:** Register a webhook URL to receive real-time JSON payloads when critical security events occur.

**Request Body (JSON):**
```json
{
  "name": "Security Operations Center (SOC) Alert",
  "url": "https://soc.acme.corp/api/alerts/webhook",
  "secret_token": "whsec_12345",
  "trigger_conditions": {
    "severity": ["HIGH", "CRITICAL"],
    "event_types": ["SQL_INJECTION", "RATE_LIMIT_EXCEEDED"]
  }
}
```

---

## 9. Dynamic Threat Intelligence
### 9.1 Configure Threat Feeds
**Endpoint:** `POST /threat-feeds`
**Description:** Integrate real-time threat intelligence feeds to block malicious IPs, ASNs, and domains dynamically.

**Request Body (JSON):**
```json
{
  "feed_name": "Emerging Threats Pro",
  "feed_type": "IP | DOMAIN | URL",
  "update_interval": "3600s",
  "confidence_threshold": 70,
  "action": "BLOCK | ALERT | SCORE",
  "source_provenance": "proofpoint.com/et"
}
```

---

## 10. Upstream mTLS & Certificate Pinning
### 10.1 Configure Backend TLS
**Endpoint:** `POST /applications/{id}/backend-tls`
**Description:** Secure communication between the WAF and upstream microservices using mTLS.

**Request Body (JSON):**
```json
{
  "enable_mtls": true,
  "client_certificate_id": "cert-12345",
  "ca_bundle_id": "ca-prod-01",
  "verify_hostname": true,
  "ssl_ciphers": ["TLS_AES_256_GCM_SHA384"]
}
```

---

## 11. AI Learning Mode & Anomaly Baselines
### 11.1 Configure Traffic Baselines
**Endpoint:** `POST /learning/baseline`
**Description:** Enable behavioral analysis to detect seasonal anomalies and non-signature-based zero-day attacks.

**Request Body (JSON):**
```json
{
  "application_id": "app-001",
  "learning_window": "7d",
  "sensitivity": "MEDIUM",
  "metrics": ["REQUEST_RATE", "ERROR_RATE", "AVG_PAYLOAD_SIZE"],
  "auto_mitigation": false,
  "deviations": {
    "threshold": 3.5,
    "action": "INCREASE_ANOMALY_SCORE"
  }
}
```

---

## 12. Response Modification & Data Masking
### 12.1 Transform Upstream Responses
**Endpoint:** `POST /policies/response-modification`
**Description:** Modify or mask sensitive data (like PANs or stack traces) leaking from upstream servers before sending to the client.

**Request Body (JSON):**
```json
{
  "rule_id": "resp-001",
  "match_pattern": "\\b\\d{4}-\\d{4}-\\d{4}-\\d{4}\\b",
  "replacement": "****-****-****-****",
  "content_types": ["application/json", "text/html"],
  "status_code_range": [200, 299]
}
```

---

## 13. Advanced Geo-Fencing
### 13.1 Regional Access Control
**Endpoint:** `POST /policies/geo-access`
**Description:** Drop traffic from unauthorized countries or ASNs at the edge.

**Request Body (JSON):**
```json
{
  "application_id": "banking-api",
  "allowed_countries": ["BD", "IN"],
  "blocked_countries": ["RU", "CN", "KP"],
  "action": "BLOCK",
  "custom_block_message": "Access restricted to Bangladesh and India."
}
```

---

## 14. SOAR Integration & Advanced Webhooks
### 14.1 Configure Incident Playbook Webhooks
**Endpoint:** `POST /integrations/webhook`
**Description:** Trigger external SOAR playbooks (like blocking an IP automatically) when critical WAF events occur.

**Request Body (JSON):**
```json
{
  "name": "SOC-Playbook",
  "trigger_events": ["SQL_INJECTION", "RCE"],
  "severity": "HIGH",
  "webhook_url": "https://soc.company.com/incident/create",
  "headers": {"X-API-Key": "encrypted_value"},
  "retry_policy": {"max_attempts": 3, "backoff": "exponential"}
}
```

---

## 15. Policy as Code (PaC) & GitOps
### 15.1 Sync Declarative Configuration
**Endpoint:** `POST /deployments/sync`
**Description:** API endpoint for CI/CD tools (Terraform/Ansible) or CLI (`wafctl apply -f policy.yaml`) to push declarative configurations directly to the WAF control plane.
