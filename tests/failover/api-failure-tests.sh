#!/bin/sh
# api-failure-tests.sh — Sprint 4: gateway registration + heartbeat API drills.
# Requires Docker Compose dev stack up (postgres + control-api).
set -eu

BASE="${API_BASE:-http://127.0.0.1:8443}"
GW_ID="01JARZ3NDEKTSV4RRFFQ69G5FBE"

PASS=0; FAIL=0
note() { printf '[TEST] %s\n' "$1"; }
ok()   { PASS=$((PASS+1)); printf '  PASS: %s\n' "$1"; }
bad()  { FAIL=$((FAIL+1)); printf '  FAIL: %s\n' "$1"; }

# Wait for API health
for _ in $(seq 1 30); do
  if curl -sf "$BASE/api/v1/health" >/dev/null 2>&1; then break; fi
  sleep 1
done

note "Test 1: register gateway (idempotent)"
if curl -sf -X POST "$BASE/api/v1/gateways/register" \
  -H 'Content-Type: application/json' \
  -d "{\"gateway_id\":\"$GW_ID\",\"hostname\":\"gw-test-01\",\"ip\":\"10.0.0.10\",\"version\":\"0.1.0\",\"capabilities\":[\"http/1.1\"]}" \
  | grep -q '"status":"active"'; then
  ok "registered"
else
  bad "registration failed"
fi

# Re-register: must be idempotent (200, not 409/500)
note "Test 2: re-registration is idempotent"
if curl -sf -X POST "$BASE/api/v1/gateways/register" \
  -H 'Content-Type: application/json' \
  -d "{\"gateway_id\":\"$GW_ID\",\"hostname\":\"gw-test-01\"}" >/dev/null; then
  ok "idempotent"
else
  bad "re-registration not idempotent"
fi

note "Test 3: heartbeat updates status and applied_hash"
if curl -sf -X POST "$BASE/api/v1/gateways/$GW_ID/heartbeat" \
  -H 'Content-Type: application/json' \
  -d '{"status":"active","applied_hash":"abc123","version":"0.1.0"}' >/dev/null; then
  ok "heartbeat accepted"
else
  bad "heartbeat rejected"
fi

note "Test 4: fleet list shows gateway + applied hash"
if curl -sf "$BASE/api/v1/gateways" | grep -q "abc123"; then
  ok "fleet list shows applied hash"
else
  bad "fleet list missing applied hash"
fi

note "Test 5: metrics endpoint exposes prometheus format"
if curl -sf "$BASE/api/v1/metrics" | grep -q "shield_gateways_registered"; then
  ok "metrics present"
else
  bad "metrics missing"
fi

note "Test 6: create + list application"
APP_ID=$(curl -sf -X POST "$BASE/api/v1/applications" \
  -H 'Content-Type: application/json' \
  -d '{"name":"test-app","description":"integration test"}' | sed -E 's/.*"id":"([^"]+)".*/\1/')
if [ -n "$APP_ID" ] && curl -sf "$BASE/api/v1/applications" | grep -q "test-app"; then
  ok "application created + listed"
else
  bad "application CRUD failed"
fi

note "Test 7: add domain + origin to application"
if curl -sf -X POST "$BASE/api/v1/applications/$APP_ID/domains" \
  -H 'Content-Type: application/json' -d '{"hostname":"app.test.local"}' >/dev/null \
  && curl -sf -X POST "$BASE/api/v1/applications/$APP_ID/origins" \
  -H 'Content-Type: application/json' -d '{"name":"origin-1","host":"10.0.0.20","port":8080}' >/dev/null \
  && curl -sf "$BASE/api/v1/applications/$APP_ID/origins" | grep -q "origin-1"; then
  ok "domain + origin created"
else
  bad "domain/origin CRUD failed"
fi

note "Test 8: create security policy (transparent default)"
if curl -sf -X POST "$BASE/api/v1/security-policies" \
  -H 'Content-Type: application/json' \
  -d '{"name":"default","enforcement_mode":"transparent"}' >/dev/null \
  && curl -sf "$BASE/api/v1/security-policies" | grep -q "default"; then
  ok "policy created"
else
  bad "policy creation failed"
fi

note "Test 9: invalid policy enforcement_mode rejected (400)"
CODE=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/api/v1/security-policies" \
  -H 'Content-Type: application/json' \
  -d '{"name":"bad","enforcement_mode":"nonsense"}')
if [ "$CODE" = "400" ]; then
  ok "invalid mode rejected"
else
  bad "invalid mode accepted ($CODE)"
fi

echo
echo "=========================================="
echo "  api failure tests: $PASS passed, $FAIL failed"
echo "=========================================="
[ "$FAIL" -eq 0 ] || exit 1
