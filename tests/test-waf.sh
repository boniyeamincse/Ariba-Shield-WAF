#!/bin/bash
set -e

echo "Logging in..."
curl -s -c cookie.txt -X POST http://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superadmin@aribashield.local","password":"admin"}' > /dev/null

echo "Creating policy version..."
VERSION_RESP=$(curl -s -b cookie.txt -X POST http://localhost:8443/api/v1/security-policies/pol-demo-003/versions \
  -H "Content-Type: application/json" \
  -d '{
    "policy_id": "pol-demo-003",
    "bundle_hash": "hash123",
    "document": {
      "schema_version": "v0",
      "config_id": "config-003",
      "virtual_servers": [{"name": "default"}],
      "backend_pools": [{"name": "pool1"}],
      "waf": {
        "enabled": true,
        "paranoia_level": 3,
        "anomaly_threshold": 10,
        "managed_rules": [{"enabled": true}]
      }
    }
  }')
echo "Version response: $VERSION_RESP"
VER_ID=$(echo $VERSION_RESP | jq -r '.id')

echo "Promoting to approved..."
curl -s -b cookie.txt -X POST "http://localhost:8443/api/v1/policy-versions/$VER_ID/promote?to=approved" > /dev/null

echo "Activating policy (triggers Push)..."
curl -s -b cookie.txt -X POST "http://localhost:8443/api/v1/security-policies/pol-demo-003/activate" > /dev/null

sleep 1

echo "Checking /rules/crs.conf in waf-engine..."
docker exec ariba-shield-dev-waf-engine-1 cat /rules/crs.conf
