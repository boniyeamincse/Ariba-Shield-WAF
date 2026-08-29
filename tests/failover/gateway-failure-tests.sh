#!/bin/sh
# gateway-failure-tests.sh — Sprint 4 failure drills (ADR-003)
# Tests: invalid-config rejection, restart recovery, last-known-good fallback,
# safe-mode. Requires Docker and a built image:
#   docker build -t ariba-shield/openresty-gateway:0.1 gateways/openresty-gateway
set -eu

IMAGE="${IMAGE:-ariba-shield/openresty-gateway:0.1}"
STORE="/tmp/shield-gw-test/store"
CONFIG_STORE="/var/lib/shield-waf/config"

GOOD_HASH="0000000000000000000000000000000000000000000000000000000000000001"
BAD_HASH="0000000000000000000000000000000000000000000000000000000000000002"

PASS=0
FAIL=0

note() { printf '[TEST] %s\n' "$1"; }
ok()   { PASS=$((PASS+1)); printf '  PASS: %s\n' "$1"; }
bad()  { FAIL=$((FAIL+1)); printf '  FAIL: %s\n' "$1"; }

rm -rf /tmp/shield-gw-test
mkdir -p "$STORE/$GOOD_HASH"

# A minimal valid shield.conf (http-level, included by entrypoint main conf).
cat > "$STORE/$GOOD_HASH/generated.conf" <<'EOF'
    server {
        listen 80 default_server;
        location / { return 200 "ok\n"; add_header Content-Type text/plain; }
        access_log off;
    }
EOF

# An intentionally invalid config (unterminated block).
mkdir -p "$STORE/$BAD_HASH"
cat > "$STORE/$BAD_HASH/generated.conf" <<'EOF'
    server {
        listen 80 default_server;
        location / { return 200
EOF

echo "$GOOD_HASH" > "$STORE/active"
echo "$GOOD_HASH" > "$STORE/last-known-good"

run_gw() {
  docker run -d --rm --name shield-failtest \
    -e CONFIG_STORE="$CONFIG_STORE" \
    -v "$STORE:$CONFIG_STORE" \
    -p 18099:80 \
    "$IMAGE" >/dev/null 2>&1
  # wait for any HTTP response (200 or 503 both mean the gateway is up)
  for _ in $(seq 1 20); do
    if curl -s -o /dev/null http://127.0.0.1:18099/ 2>/dev/null; then return 0; fi
    sleep 0.5
  done
  return 1
}

stop_gw() { docker rm -f shield-failtest >/dev/null 2>&1 || true; }

# --- Test 1: valid config serves ---
stop_gw
note "Test 1: gateway starts with valid config and serves"
if run_gw && curl -sf http://127.0.0.1:18099/ | grep -q ok; then
  ok "served 200 ok"
else
  bad "gateway did not serve"
fi
stop_gw

# --- Test 2: invalid active config falls back to last-known-good ---
note "Test 2: invalid ACTIVE config falls back to last-known-good"
echo "$BAD_HASH" > "$STORE/active"
echo "$GOOD_HASH" > "$STORE/last-known-good"
if run_gw && curl -sf http://127.0.0.1:18099/ | grep -q ok; then
  ok "fell back to last-known-good and served"
else
  bad "did not serve after fallback"
fi
stop_gw

# --- Test 3: both active and LKG invalid -> safe-mode 503 ---
note "Test 3: no valid config -> safe-mode"
echo "$BAD_HASH" > "$STORE/active"
echo "$BAD_HASH" > "$STORE/last-known-good"
if run_gw && curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:18099/ | grep -q 503; then
  ok "safe-mode returned 503"
else
  bad "safe-mode did not return 503"
fi
stop_gw

# --- Test 4: missing store -> safe-mode (never crash) ---
note "Test 4: missing config store -> safe-mode, no crash"
rm -rf "$STORE"
mkdir -p "$STORE"
if run_gw && curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:18099/ | grep -q 503; then
  ok "safe-mode 503 with empty store"
else
  bad "gateway crashed or wrong status with empty store"
fi
stop_gw

echo
echo "=========================================="
echo "  gateway failure tests: $PASS passed, $FAIL failed"
echo "=========================================="
[ "$FAIL" -eq 0 ] || exit 1
