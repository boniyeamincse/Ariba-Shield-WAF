#!/bin/bash
set -e

echo "=== 0. Login ==="
curl -s -c cookie.txt -X POST http://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superadmin@aribashield.local","password":"admin"}' > /dev/null

echo "Logged in."

echo "=== 1. Create a New Rule (SQLi Detection) ==="
RULE_RESPONSE=$(curl -s -b cookie.txt -X POST http://localhost:8443/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "rule_id": "ARB-SQL-'"$RANDOM"'",
    "name": "SQL Injection Detection",
    "description": "Detect common SQL injection patterns",
    "type": "custom",
    "category": "sqli",
    "severity": "critical",
    "priority": 10,
    "action": "block",
    "status": "active",
    "logic": "AND",
    "conditions": [
      {
        "group_id": 0,
        "field": "request_body",
        "operator": "contains",
        "value": "union select",
        "transformation": "lowercase",
        "case_sensitive": false
      },
      {
        "group_id": 0,
        "field": "method",
        "operator": "equals",
        "value": "POST",
        "transformation": "",
        "case_sensitive": false
      }
    ],
    "scopes": [
      {
        "path_pattern": "/*",
        "methods": ["GET", "POST", "PUT", "PATCH", "DELETE"]
      }
    ]
  }')

echo "Response: $RULE_RESPONSE"
RULE_ID=$(echo $RULE_RESPONSE | jq -r '.id')

if [ "$RULE_ID" == "null" ] || [ -z "$RULE_ID" ]; then
    echo "Failed to create rule."
    exit 1
fi
echo "Rule created successfully with ID: $RULE_ID"
echo ""

echo "=== 2. List Rules ==="
curl -s -b cookie.txt "http://localhost:8443/api/v1/rules?category=sqli" | jq .
echo ""

echo "=== 3. Test Rule Match (Should Match) ==="
curl -s -b cookie.txt -X POST "http://localhost:8443/api/v1/rules/$RULE_ID/test" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "POST",
    "url": "/login",
    "body": "username=admin&password=union SELECT * FROM users"
  }' | jq .
echo ""

echo "=== 4. Test Rule Match (Should NOT Match) ==="
curl -s -b cookie.txt -X POST "http://localhost:8443/api/v1/rules/$RULE_ID/test" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "GET",
    "url": "/login",
    "body": "username=admin&password=password123"
  }' | jq .
echo ""

rm cookie.txt
