#!/bin/sh
# entrypoint.sh — Ariba Shield WAF gateway entrypoint
# Starts nginx with last-known-good config recovery (ADR-003).

set -e

CONFIG_STORE="${CONFIG_STORE:-/var/lib/shield-waf/config}"
ACTIVE_POINTER="${ACTIVE_POINTER:-$CONFIG_STORE/active}"
LKG_POINTER="${LKG_POINTER:-$CONFIG_STORE/last-known-good}"
NGINX_CONF="${NGINX_CONF:-/usr/local/openresty/nginx/conf/nginx.conf}"
SHIELD_CONF_SRC="${SHIELD_CONF_SRC:-/etc/nginx/conf.d/shield.conf}"
OPENRESTY_PREFIX="${OPENRESTY_PREFIX:-/usr/local/openresty}"

# Write a minimal main nginx.conf that includes the generated shield config
# and references OpenResty's standard paths.
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

# Phase 1: determine which config to start with.
resolve_config() {
  if [ -f "$ACTIVE_POINTER" ]; then
    HASH=$(cat "$ACTIVE_POINTER" 2>/dev/null)
    CANDIDATE="$CONFIG_STORE/$HASH/generated.conf"
    if [ -f "$CANDIDATE" ]; then
      ln -sf "$CANDIDATE" "$SHIELD_CONF_SRC" 2>/dev/null || cp "$CANDIDATE" "$SHIELD_CONF_SRC"
      write_main_conf
      if nginx -t -c "$NGINX_CONF" 2>/dev/null; then
        echo "$CANDIDATE"
        return
      fi
    fi
    echo "[entrypoint] Active config $HASH failed validation; falling back to last-known-good." >&2
  fi

  if [ -f "$LKG_POINTER" ]; then
    HASH=$(cat "$LKG_POINTER" 2>/dev/null)
    CANDIDATE="$CONFIG_STORE/$HASH/generated.conf"
    if [ -f "$CANDIDATE" ]; then
      ln -sf "$CANDIDATE" "$SHIELD_CONF_SRC" 2>/dev/null || cp "$CANDIDATE" "$SHIELD_CONF_SRC"
      write_main_conf
      if nginx -t -c "$NGINX_CONF" 2>/dev/null; then
        echo "$CANDIDATE"
        return
      fi
    fi
    echo "[entrypoint] Last-known-good config $HASH also failed validation." >&2
  fi

  echo "[entrypoint] No valid config found; starting safe-mode." >&2
  echo "safe-mode"
}

CONFIG_PATH=$(resolve_config)

if [ "$CONFIG_PATH" = "safe-mode" ]; then
  mkdir -p /etc/nginx/conf.d
  cat > "$SHIELD_CONF_SRC" <<'SAFE'
server {
    listen 80 default_server;
    location / { return 503 "Service Unavailable\n"; add_header Content-Type text/plain; }
    access_log off;
}
SAFE
  write_main_conf
  echo "[entrypoint] Safe-mode config written."
else
  echo "[entrypoint] Using config: $CONFIG_PATH"
fi

# Phase 2: start nginx
exec "$@"