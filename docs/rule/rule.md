হ্যাঁ। **Rules menu শেষ করার পর এখন আর নতুন menu ধরব না।** আমরা একটার পর একটা পুরো product build করব—**UI → Frontend logic → Backend → Database → API → WAF engine integration → Testing**।

তোমার Ariba Shield-এর জন্য আমি এই implementation order recommend করছি।

## Phase 1 — Rules Module Complete

প্রথমে Rules module পুরোপুরি production-ready করি।

### Rules-এর Rule Types

```text
Rules
│
├── Managed Rules
│   ├── OWASP CRS
│   ├── Ariba Core Security Rules
│   └── Threat Intelligence Rules
│
└── Custom Rules
    ├── SQL Injection
    ├── XSS
    ├── Path Traversal
    ├── Command Injection
    ├── File Inclusion
    ├── SSRF
    ├── XXE
    ├── RCE
    ├── HTTP Protocol
    ├── Scanner Detection
    ├── Bot Detection
    ├── Rate Limiting
    ├── IP Access Control
    ├── Geo Blocking
    ├── Header Validation
    └── Custom
```

---

# 1. Rules List UI

Route:

```text
/bn/rules
```

### Header

```text
Rules

Manage WAF detection and protection rules.

                              [+ New Rule]
```

### Filter

```text
Search rules...

Category     [All ▼]
Type         [All ▼]
Severity     [All ▼]
Action       [All ▼]
Status       [All ▼]
```

### Table

| Rule ID     | Rule Name               | Category   | Severity | Action | Status | Priority | Actions |
| ----------- | ----------------------- | ---------- | -------- | ------ | ------ | -------: | ------- |
| ARB-SQL-001 | SQL Injection Detection | SQLi       | Critical | Block  | Active |       10 | ⋮       |
| ARB-XSS-001 | XSS Detection           | XSS        | High     | Block  | Active |       20 | ⋮       |
| ARB-LFI-001 | Path Traversal          | LFI        | High     | Block  | Active |       30 | ⋮       |
| ARB-RL-001  | Login Rate Limit        | Rate Limit | Medium   | Block  | Active |       40 | ⋮       |

Actions:

```text
View
Edit
Duplicate
Disable
Delete
Test Rule
View Logs
```

---

# 2. New Rule UI

Route:

```text
/bn/rules/create
```

আমি এটাকে **wizard** করলে UI অনেক cleaner হবে।

```text
1 Basic
   ↓
2 Match
   ↓
3 Conditions
   ↓
4 Action
   ↓
5 Scope
   ↓
6 Review
```

---

# Step 1 — Basic

```text
Rule Name *
[ SQL Injection Detection ]

Rule ID *
[ ARB-SQL-001 ]

Description
[ Detect common SQL injection patterns ]

Rule Type *
[ Custom Rule ▼ ]

Category *
[ SQL Injection ▼ ]

Severity *
[ Critical ▼ ]

Priority *
[ 10 ]

Status
[● Enabled]
```

---

# Step 2 — Match

এখানে request-এর কোন অংশ inspect করবে।

```text
Match Target *

[ URL ▼ ]

Operator *

[ Contains ▼ ]

Value *

[ union select ]

Transformation

[ URL Decode ▼ ]

Case Sensitive

[ No ▼ ]
```

Target:

```text
URL
Query Parameter
Request Header
Cookie
Request Body
HTTP Method
Source IP
User-Agent
Host
```

Operator:

```text
Equals
Not Equals
Contains
Not Contains
Starts With
Ends With
Regex
IP Match
CIDR Match
Greater Than
Less Than
```

---

# Step 3 — Conditions

Multiple conditions support করতে হবে।

```text
Conditions

┌─────────────────────────────────────────────┐
│ [Request Body] [contains] [union select]   │
│                                             │
│                 AND                         │
│                                             │
│ [Method] [equals] [POST]                   │
└─────────────────────────────────────────────┘

[ + Add Condition ]
```

Logical:

```text
AND
OR
```

Nested condition ভবিষ্যতে:

```text
(
   Condition A
   AND
   Condition B
)
OR
(
   Condition C
   AND
   Condition D
)
```

---

# Step 4 — Action

```text
Action *

○ Allow
○ Log
● Block
○ Challenge
○ Rate Limit
```

Block হলে:

```text
HTTP Status
[403]

Response Type
[JSON ▼]

Message
[Request blocked by security policy]
```

Log:

```text
Log Event
[✓]

Generate Security Event
[✓]
```

---

# Step 5 — Scope

Rule কোথায় apply করবে?

```text
Scope

○ All Applications

● Selected Applications

Applications:

☑ ERP
☑ HR Portal
☐ CRM
☐ Website
```

Path:

```text
Apply Path

/api/*
/login
/admin/*
```

Methods:

```text
☑ GET
☑ POST
☑ PUT
☑ PATCH
☑ DELETE
☐ OPTIONS
☐ HEAD
```

---

# Step 6 — Review

শেষে পুরো rule summary:

```text
Review Rule
────────────────────────────────

Name:
SQL Injection Detection

ID:
ARB-SQL-001

Category:
SQL Injection

Severity:
Critical

Priority:
10

Match:
Request Body contains "union select"

Condition:
AND Method = POST

Action:
BLOCK

Status:
Enabled

Scope:
ERP, HR Portal
```

Buttons:

```text
[ Back ]       [ Save Rule ]
```

---

# 3. Rule Details UI

Route:

```text
/bn/rules/{id}
```

এখানে dashboard-style detail page।

```text
SQL Injection Detection
ARB-SQL-001

● Active                         [Edit Rule]
```

### Overview

```text
Category       SQL Injection
Severity       Critical
Action         Block
Priority       10
Created        31 Aug 2026
Updated        31 Aug 2026
```

### Rule Logic

```text
IF

Request Body
    contains
"union select"

AND

Method
    equals
POST

THEN

BLOCK
403
```

### Statistics

```text
Matches Today       1,248
Blocked             1,241
Allowed                 7
Applications             3
```

### Tabs

```text
Overview
Conditions
Applications
Events
History
```

---

# 4. Backend Database Design

এখন আসল backend।

আমি minimum এই tables রাখব।

### `waf_rules`

```text
id
uuid
name
rule_id
description
type
category
severity
priority
action
status
created_by
created_at
updated_at
```

### `waf_rule_conditions`

```text
id
rule_id
group_id
field
operator
value
transformation
case_sensitive
created_at
```

Example:

```text
rule_id = 1

field = request_body
operator = contains
value = union select
```

---

### `waf_rule_scopes`

```text
id
rule_id
application_id
path_pattern
methods
created_at
```

---

### `waf_rule_actions`

```text
id
rule_id
action
status_code
response_type
response_message
log_enabled
created_at
```

---

### `waf_rule_versions`

Very important.

```text
id
rule_id
version
configuration
changed_by
change_type
created_at
```

এটা রাখলে পরে:

```text
Rule History

v1
v2
v3
v4
```

এবং rollback করতে পারবে।

---

# 5. Backend API

REST API structure:

```text
GET     /api/rules
POST    /api/rules

GET     /api/rules/{id}
PUT     /api/rules/{id}
DELETE  /api/rules/{id}

POST    /api/rules/{id}/enable
POST    /api/rules/{id}/disable

POST    /api/rules/{id}/duplicate

POST    /api/rules/{id}/test

GET     /api/rules/{id}/events
GET     /api/rules/{id}/history
```

---

# 6. Rule Validation

Backend অবশ্যই validation করবে।

Example:

```text
Rule Name
required
max:255

Rule ID
required
unique

Category
required

Severity
required

Priority
required|integer

Action
required

Conditions
required|array
```

আর security validation:

```text
Regex validation
CIDR validation
IP validation
HTTP method validation
JSON validation
Path validation
```

---

# 7. Rule Engine

এটাই সবচেয়ে important backend component।

Request আসবে:

```text
Client
   ↓
Ariba Shield
   ↓
Request Parser
   ↓
Rule Engine
   ↓
Rule Matching
   ↓
Action
```

Example:

```text
POST /login

Body:
username=admin' OR 1=1--
```

Engine:

```text
Request
   ↓
Load active rules
   ↓
Priority sort
   ↓
Evaluate condition
   ↓
SQLi rule matched
   ↓
Action = BLOCK
   ↓
403
   ↓
Security Event
```

---

# 8. Rule Engine-এর জন্য Internal Structure

Frontend থেকে rule এলে backend normalized configuration বানাবে:

```json
{
  "id": "ARB-SQL-001",
  "priority": 10,
  "status": "enabled",
  "conditions": [
    {
      "field": "request_body",
      "operator": "contains",
      "value": "union select"
    },
    {
      "field": "method",
      "operator": "equals",
      "value": "POST"
    }
  ],
  "logic": "AND",
  "action": {
    "type": "block",
    "status_code": 403
  }
}
```

তারপর WAF engine এই configuration consume করবে।

---

# 9. Rule Test Engine

UI:

```text
Test Rule

Method
[POST]

URL
[/login]

Headers
[Content-Type: application/x-www-form-urlencoded]

Body

username=admin' OR 1=1--
password=test

             [ Test Rule ]
```

Result:

```text
✓ RULE MATCHED

Rule:
ARB-SQL-001

Matched:
request_body

Action:
BLOCK

Status:
403
```

আর clean request:

```text
✓ RULE NOT MATCHED

Action:
ALLOW
```

---

# 10. Security Events Integration

Rule match করলে:

```text
Rule Engine
     ↓
Security Event
     ↓
Traffic / Security Events
```

Event:

```text
Event ID
EVT-98231

Rule
ARB-SQL-001

Source IP
10.x.x.x

Destination
ERP

URL
/login

Method
POST

Category
SQL Injection

Severity
Critical

Action
Blocked

Timestamp
...
```

অর্থাৎ Rules module আলাদা হলেও **Traffic/Security Events-এর সাথে connected থাকবে।**

---

# 11. প্রথমে যেসব Rules বানাবে

আমি development-এর জন্য প্রথম batch হিসেবে এগুলো করতাম:

### Critical

```text
ARB-SQL-001  SQL Injection
ARB-RCE-001  Remote Code Execution
ARB-CMD-001  Command Injection
ARB-SSRF-001 SSRF
```

### High

```text
ARB-XSS-001  Cross-Site Scripting
ARB-LFI-001  Local File Inclusion
ARB-PT-001   Path Traversal
ARB-XXE-001  XXE
ARB-FI-001   File Inclusion
ARB-HTTP-001 HTTP Protocol Attack
```

### Medium

```text
ARB-SCAN-001 Scanner Detection
ARB-BOT-001  Malicious Bot
ARB-RL-001   Rate Limit
ARB-HEAD-001 Header Validation
```

### Access Control

```text
ARB-IP-001   IP Block
ARB-IP-002   IP Allow
ARB-GEO-001  Geo Block
```

---

# 12. Rule Development Priority

**সব rule একসাথে code করবে না।**

এই order follow করো:

```text
PHASE 1
│
├── Rules List UI
├── New Rule UI
├── Edit Rule
├── Delete
├── Enable/Disable
└── Duplicate
        ↓
PHASE 2
│
├── Database
├── Models
├── Migration
├── CRUD API
└── Validation
        ↓
PHASE 3
│
├── Rule Engine
├── Condition Engine
├── Operator Engine
├── Action Engine
└── Priority Engine
        ↓
PHASE 4
│
├── Test Rule
├── Security Events
├── Rule Statistics
└── Rule History
        ↓
PHASE 5
│
├── OWASP CRS
├── Managed Rules
├── Import/Export
└── Rule Version/Rollback
```

## এরপর কী করব

তুমি যেহেতু **Rules menu already complete** করেছ, এখন আমার মতে **অন্য কোনো menu-তে যাওয়া যাবে না**।

**Next task = `Rules → New Rule` পুরো UI তৈরি করা।**

তারপর একই form-এর জন্য:

```text
Frontend
   ↓
Form State
   ↓
Validation
   ↓
API
   ↓
go Backend
   ↓
Database
   ↓
Rule Engine
```

এই flow শেষ করব।

**প্রথম implementation হিসেবে `New Rule`-এর UI + exact fields + frontend data structure + structure** তৈরি করাই সবচেয়ে logical next step।
