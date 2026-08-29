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

