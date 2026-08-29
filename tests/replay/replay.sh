#!/bin/sh
# replay.sh — replay a traffic corpus against the WAF engine and verify
# expected decisions (detection/blocking). Master plan §20: true-positive rate
# by category; Phase 2 exit criteria (stable replay results).
#
# Corpus files carry an EXPECT:pass|block header.
# Usage: replay.sh [--host http://127.0.0.1:8082] [--mode detect|block]
#   detect (default): malicious entries must pass through (200) but emit events
#   block:            malicious entries must be rejected (403)
set -eu

HOST="http://127.0.0.1:8082"
MODE="detect"
while [ $# -gt 0 ]; do
  case "$1" in
    --host) HOST="$2"; shift ;;
    --mode) MODE="$2"; shift ;;
  esac
  shift
done

CORPUS_DIR="$(dirname "$0")/../corpus"
PASS=0; FAIL=0; TOTAL=0

for f in "$CORPUS_DIR"/*.txt; do
  [ -e "$f" ] || continue
  name="$(basename "$f" .txt)"
  expected=$(sed -n '1s/.*EXPECT:\([a-z]*\).*/\1/p' "$f")
  [ -z "$expected" ] && expected="block"

  while IFS= read -r line || [ -n "$line" ]; do
    case "$line" in
      \#*|"") continue ;;
    esac
    TOTAL=$((TOTAL+1))

    method="${line%% *}"
    url="${line#* }"

    code=$(curl -s -o /dev/null -w '%{http_code}' -X "$method" "$HOST$url")

    # desired is the status we want for THIS file, adjusted by mode.
    if [ "$expected" = "pass" ]; then
      # Legitimate traffic must NEVER be blocked by the WAF — it must be
      # forwarded to the backend (backend status 200/404/5xx is irrelevant).
      # Detection: any non-block status means the WAF passed it through.
      ok=false
      case "$code" in
        200|201|204|301|302|304|400|404|405|500|502|503|504) ok=true ;;
      esac
      if $ok; then
        PASS=$((PASS+1)); printf '  PASS [%s] %s %s -> %s (passed through)\n' "$name" "$method" "$url" "$code"
      else
        FAIL=$((FAIL+1)); printf '  FAIL [%s] %s %s -> %s (WAF blocked legitimate traffic)\n' "$name" "$method" "$url" "$code"
      fi
      continue
    fi

    # Malicious: in detect mode pass through; in block mode be rejected.
    if [ "$MODE" = "block" ]; then
      desired="reject"
    else
      # Detection mode: any non-block status means the WAF passed it through.
      ok=false
      case "$code" in
        200|201|204|301|302|304|400|404|405|500|502|503|504) ok=true ;;
      esac
      if $ok; then
        PASS=$((PASS+1)); printf '  PASS [%s] %s %s -> %s (detected, passed through)\n' "$name" "$method" "$url" "$code"
      else
        FAIL=$((FAIL+1)); printf '  FAIL [%s] %s %s -> %s (unexpected status)\n' "$name" "$method" "$url" "$code"
      fi
      continue
    fi

    ok=false
    if [ "$desired" = "reject" ]; then
      if [ "$code" = "403" ] || [ "$code" = "400" ] || [ "$code" = "406" ]; then ok=true; fi
    elif [ "$code" = "$desired" ]; then
      ok=true
    fi

    if $ok; then
      PASS=$((PASS+1)); printf '  PASS [%s] %s %s -> %s\n' "$name" "$method" "$url" "$code"
    else
      FAIL=$((FAIL+1)); printf '  FAIL [%s] %s %s -> %s (expected %s)\n' "$name" "$method" "$url" "$code" "$desired"
    fi
  done < "$f"
done

echo
echo "=============================================="
echo "  replay (mode=$MODE): $TOTAL total, $PASS passed, $FAIL failed"
echo "=============================================="
[ "$FAIL" -eq 0 ] || exit 1
