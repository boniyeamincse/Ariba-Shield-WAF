#!/bin/sh
# apply-config.sh — validate, stage, and atomically switch nginx config (ADR-003).
# Usage: apply-config.sh <config_hash> <generated_config_path>

set -euo pipefail

CONFIG_STORE="${CONFIG_STORE:-/var/lib/shield-waf/config}"
ACTIVE_POINTER="${ACTIVE_POINTER:-$CONFIG_STORE/active}"
LKG_POINTER="${LKG_POINTER:-$CONFIG_STORE/last-known-good}"
CONFIG_SOURCE="${CONFIG_SOURCE:-/etc/nginx/conf.d/shield.conf}"

HASH="$1"
STAGED_CONFIG="$2"

if [ ! -f "$STAGED_CONFIG" ]; then
  echo "[apply-config] ERROR: staged config not found: $STAGED_CONFIG" >&2
  exit 1
fi

# Step 1: validate with nginx -t
staging_dir="$CONFIG_STORE/$HASH"
mkdir -p "$staging_dir"

# Write the staged config to the version directory.
cp "$STAGED_CONFIG" "$staging_dir/generated.conf"

# Validate with nginx -t using the staged config.
if ! nginx -t -c "$staging_dir/generated.conf" 2>/dev/null; then
  echo "[apply-config] REJECTED: config $HASH failed nginx -t" >&2
  echo "rejected" > "$staging_dir/meta.json"
  rm -f "$staging_dir/generated.conf"
  exit 1
fi

# Step 2: atomic switch — write active pointer, then send SIGHUP.
echo "$HASH" > "$ACTIVE_POINTER.tmp"
mv "$ACTIVE_POINTER.tmp" "$ACTIVE_POINTER"  # atomic rename

# Symlink the active config source.
ln -sf "$staging_dir/generated.conf" "$CONFIG_SOURCE"

# Reload nginx gracefully.
if nginx -s reload 2>/dev/null; then
  echo "[apply-config] ACTIVATED: config $HASH" >&2
  echo "active" > "$staging_dir/meta.json"

  # Promote to last-known-good after a successful reload.
  echo "$HASH" > "$LKG_POINTER.tmp"
  mv "$LKG_POINTER.tmp" "$LKG_POINTER"  # atomic rename
  exit 0
else
  echo "[apply-config] RELOAD FAILED; rolling back to last-known-good" >&2
  echo "rolled_back" > "$staging_dir/meta.json"
  # Rollback handled by the next resolve_config cycle on restart.
  exit 1
fi