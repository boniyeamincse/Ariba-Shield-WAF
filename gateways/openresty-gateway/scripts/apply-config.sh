#!/bin/sh
# apply-config.sh — validate, stage, and atomically switch nginx config (ADR-003).
# Usage: apply-config.sh <config_hash> <generated_config_path>
#
# The generated config is an http-level fragment intended to be `include`d from
# the main nginx.conf. Validation must therefore run against a wrapper main
# config, exactly as entrypoint.sh does at boot (P0.1).

set -euo pipefail

CONFIG_STORE="${CONFIG_STORE:-/var/lib/shield-waf/config}"
ACTIVE_POINTER="${ACTIVE_POINTER:-$CONFIG_STORE/active}"
LKG_POINTER="${LKG_POINTER:-$CONFIG_STORE/last-known-good}"
CONFIG_SOURCE="${CONFIG_SOURCE:-/etc/nginx/conf.d/shield.conf}"
NGINX_CONF="${NGINX_CONF:-/usr/local/openresty/nginx/conf/nginx.conf}"
SOAK_SECONDS="${SOAK_SECONDS:-60}"
HEALTH_URL="${HEALTH_URL:-http://127.0.0.1/shield-health}"

HASH="$1"
STAGED_CONFIG="$2"

if [ ! -f "$STAGED_CONFIG" ]; then
  echo "[apply-config] ERROR: staged config not found: $STAGED_CONFIG" >&2
  exit 1
fi

# write_main_conf writes the wrapper main config that includes the generated
# fragment (mirrors entrypoint.sh; keep in sync).
write_main_conf() {
  cat > "$NGINX_CONF" <<'NEOF'
worker_processes auto;
error_log /var/log/shield-waf/error.log info;
pid /usr/local/openresty/nginx/logs/nginx.pid;

events {
    worker_connections 4096;
    multi_accept on;
    use epoll;
}

http {
    include /usr/local/openresty/nginx/conf/mime.types;
    default_type application/octet-stream;
    include /etc/nginx/conf.d/*.conf;
}
NEOF
}

# Step 1: stage the config into the immutable version store.
staging_dir="$CONFIG_STORE/$HASH"
mkdir -p "$staging_dir"
cp "$STAGED_CONFIG" "$staging_dir/generated.conf"

# Symlink it as the active config source so the wrapper main config picks it up.
mkdir -p "$(dirname "$CONFIG_SOURCE")"
ln -sf "$staging_dir/generated.conf" "$CONFIG_SOURCE"

# Ensure the wrapper main config exists, then validate the WHOLE config.
write_main_conf
if ! nginx -t -c "$NGINX_CONF" >/dev/null 2>&1; then
  echo "[apply-config] REJECTED: config $HASH failed nginx -t" >&2
  echo "rejected" > "$staging_dir/meta.json"
  rm -f "$staging_dir/generated.conf" "$CONFIG_SOURCE"
  exit 1
fi

# Step 2: atomic switch — update active pointer, then graceful reload.
echo "$HASH" > "$ACTIVE_POINTER.tmp"
mv "$ACTIVE_POINTER.tmp" "$ACTIVE_POINTER"  # atomic rename

if ! nginx -s reload >/dev/null 2>&1; then
  echo "[apply-config] RELOAD FAILED; rolling back to last-known-good" >&2
  echo "rolled_back" > "$staging_dir/meta.json"
  rollback_active
  exit 1
fi

# Step 3: soak — only promote to last-known-good after the config serves
# healthily for the soak window (ADR-003 D2). During this window the new config
# is live but not yet the recovery point.
soak_ok=1
deadline=$(( $(date +%s) + SOAK_SECONDS ))
while [ "$(date +%s)" -lt "$deadline" ]; do
  if ! curl -sf "$HEALTH_URL" >/dev/null 2>&1; then
    soak_ok=0
    break
  fi
  sleep 5
done

if [ "$soak_ok" -eq 0 ]; then
  echo "[apply-config] SOAK FAILED; rolling back to last-known-good" >&2
  echo "rolled_back" > "$staging_dir/meta.json"
  rollback_active
  exit 1
fi

echo "[apply-config] ACTIVATED: config $HASH" >&2
echo "active" > "$staging_dir/meta.json"
echo "$HASH" > "$LKG_POINTER.tmp"
mv "$LKG_POINTER.tmp" "$LKG_POINTER"  # atomic rename
exit 0

# rollback_active re-points the active pointer to the last-known-good config
# and reloads, restoring service immediately (ADR-003 D4).
rollback_active() {
  if [ -f "$LKG_POINTER" ]; then
    LKG=$(cat "$LKG_POINTER" 2>/dev/null || true)
    if [ -n "$LKG" ]; then
      echo "$LKG" > "$ACTIVE_POINTER.tmp"
      mv "$ACTIVE_POINTER.tmp" "$ACTIVE_POINTER"
      ln -sf "$CONFIG_STORE/$LKG/generated.conf" "$CONFIG_SOURCE"
      nginx -s reload >/dev/null 2>&1 || true
      echo "[apply-config] rolled back active to $LKG" >&2
    fi
  fi
}
