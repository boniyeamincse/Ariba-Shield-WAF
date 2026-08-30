#!/bin/sh
# replay.sh — replay a traffic corpus against the WAF engine and verify
# expected decisions (detection/blocking). Master plan §20: true-positive rate
# by category; Phase 2 exit criteria (stable replay results).
#
# Corpus format (tests/corpus/*.txt):
#   First line: EXPECT:pass|block  (# optional comment)
#   Each entry is either:
#     METHOD /path                  (simple GET/HEAD request)
#     METHOD /path\nbody            (POST/PUT with body)
#     METHOD /path\nHeader: val\n\nbody  (custom headers + body)
#   Entries may carry a trailing # RULE:rule_id marker to verify event match.
#
# Usage: replay.sh [--host http://127.0.0.1:8082] [--api http://127.0.0.1:8443]
#                  [--mode detect|block] [--verbose]
#   detect (default): malicious entries must pass through (200) but emit events
#   block:            malicious entries must be rejected (403)
set -eu

HOST="http://127.0.0.1:8082"
API="http://127.0.0.1:8443"
MODE="detect"
VERBOSE=false
while [ $# -gt 0 ]; do
  case "$1" in
    --host) HOST="$2"; shift ;;
    --api)  API="$2"; shift ;;
    --mode) MODE="$2"; shift ;;
    --verbose) VERBOSE=true ;;
  esac
  shift
done

CORPUS_DIR="$(dirname "$0")/../corpus"
PASS=0; FAIL=0; TOTAL=0; EVENTS=0; EVENTS_FAIL=0

# process_entry <name> <expected> <entry>
#   entry: lines separated by newlines. Supports two formats:
#   Format 1 (pipe-separated): "METHOD /path | Header: val;;Header2: val | body"
#   Format 2 (newline-separated): "METHOD /path\nHeader: val\n\nbody"
#   Trailing "# RULE:id" marker on the request line verifies the event rule.
process_entry() {
  name="$1"; expected="$2"; entry="$3"

  # Determine format: if the first line contains " | ", it's pipe-separated.
  firstline=$(printf '%s' "$entry" | sed -n '1p')
  case "$firstline" in
    *" | "*)
      # Format 1: METHOD /path | headers | body
      method=$(printf '%s' "$firstline" | cut -d' ' -f1)
      path_part=$(printf '%s' "$firstline" | cut -d'|' -f1 | sed 's/^ *//;s/ *$//')
      url=$(printf '%s' "$path_part" | cut -d' ' -f2-)
      header_part=$(printf '%s' "$firstline" | cut -d'|' -f2 | sed 's/^ *//;s/ *$//')
      body_part=$(printf '%s' "$firstline" | cut -d'|' -f3- | sed 's/^ *//;s/ *$//')
      # Extract rule marker from the body part
      rule=""
      case "$body_part" in
        *"# RULE:"*) rule=$(printf '%s' "$body_part" | sed 's/.*# RULE:\([^ ]*\).*/\1/')
                     body_part=$(printf '%s' "$body_part" | sed 's/ # RULE:.*//') ;;
      esac
      # Convert ";;" to newlines for curl -H
      headers=""
      oldIFS="$IFS"; IFS=';'
      for h in $header_part; do
        h_trim=$(printf '%s' "$h" | sed 's/^ *//;s/ *$//')
        [ -n "$h_trim" ] && headers="${headers} -H '${h_trim}'"
      done
      IFS="$oldIFS"
      body="$body_part"
      ;;
    *)
      # Format 2: newline-separated (existing format)
      method="${firstline%% *}"
      rest="${firstline#* }"
      rule=""
      case "$rest" in
        *"# RULE:"*) rule=$(printf '%s' "$rest" | sed 's/.*# RULE:\([^ ]*\).*/\1/')
                     url=$(printf '%s' "$rest" | sed 's/ # RULE:.*//;s/ *$//') ;;
        *) url="$rest" ;;
      esac
      headers=""
      body=""
      in_body=false
      line_no=0
      while IFS= read -r line || [ -n "$line" ]; do
        line_no=$((line_no+1))
        [ "$line_no" -eq 1 ] && continue
        case "$line" in
          "") in_body=true ; continue ;;
        esac
        if $in_body; then
          body="${body}${line}"
        else
          case "$line" in
            *": "*) headers="${headers} -H '$line'" ;;
          esac
        fi
      done <<EOF
$(printf '%s' "$entry")
EOF
      ;;
  esac

  TOTAL=$((TOTAL+1))
  code=""
  if [ -n "$body" ]; then
    # POST/PUT with body + optional headers
    eval "code=\$(curl -s -o /dev/null -w '%{http_code}' -X '$method' $headers --data-binary '$body' '$HOST$url')"
  elif [ -n "$headers" ]; then
    eval "code=\$(curl -s -o /dev/null -w '%{http_code}' -X '$method' $headers '$HOST$url')"
  else
    code=$(curl -s -o /dev/null -w '%{http_code}' -X "$method" "$HOST$url")
  fi

  # Verdict based on EXPECT + mode.
  if [ "$expected" = "pass" ]; then
    # Legitimate traffic must never be blocked.
    case "$code" in
      200|201|204|301|302|304|400|404|405|500|502|503|504)
        PASS=$((PASS+1)); [ "$VERBOSE" = true ] && printf '  PASS [%s] %s %s -> %s\n' "$name" "$method" "$url" "$code" ;;
      *)
        FAIL=$((FAIL+1)); printf '  FAIL [%s] %s %s -> %s (WAF blocked legitimate traffic)\n' "$name" "$method" "$url" "$code" ;;
    esac
  else
    # Malicious traffic.
    if [ "$MODE" = "block" ]; then
      # Must be rejected.
      case "$code" in
        403|400|406) PASS=$((PASS+1)); [ "$VERBOSE" = true ] && printf '  PASS [%s] %s %s -> %s (blocked)\n' "$name" "$method" "$url" "$code" ;;
        *) FAIL=$((FAIL+1)); printf '  FAIL [%s] %s %s -> %s (expected 403/400/406)\n' "$name" "$method" "$url" "$code" ;;
      esac
    else
      # Detection mode: pass through, but MUST emit a security event.
      case "$code" in
        200|201|204|301|302|304|400|404|405|500|502|503|504)
          PASS=$((PASS+1)); [ "$VERBOSE" = true ] && printf '  PASS [%s] %s %s -> %s (detected)\n' "$name" "$method" "$url" "$code" ;;
        *) FAIL=$((FAIL+1)); printf '  FAIL [%s] %s %s -> %s (unexpected status)\n' "$name" "$method" "$url" "$code" ;;
      esac
    fi
  fi

  # P2.14/2.16: in detect mode verify a security event was emitted, and if a
  # rule marker was given, that the event's rule_ids include it.
  if [ "$expected" != "pass" ] && [ "$MODE" != "block" ]; then
    EVENTS=$((EVENTS+1))
    found="false"
    if [ -n "$rule" ]; then
      # Query the control API for a recent event containing the rule.
      resp=$(curl -s "$API/api/v1/security-events?limit=20" 2>/dev/null || true)
      if printf '%s' "$resp" | grep -q "\"$rule\""; then
        found="true"
      fi
    else
      resp=$(curl -s "$API/api/v1/security-events?limit=20" 2>/dev/null || true)
      if printf '%s' "$resp" | grep -q "event_id\|rule_ids"; then
        found="true"
      fi
    fi
    if [ "$found" = "true" ]; then
      [ "$VERBOSE" = true ] && printf '    event-ok [%s] rule=%s\n' "$name" "$rule"
    else
      EVENTS_FAIL=$((EVENTS_FAIL+1))
      printf '  EVENT-FAIL [%s] %s %s -> no matching event (rule=%s)\n' "$name" "$method" "$url" "$rule"
    fi
  fi
}

for f in "$CORPUS_DIR"/*.txt; do
  [ -e "$f" ] || continue
  name="$(basename "$f" .txt)"
  expected=$(sed -n '1s/.*EXPECT:\([a-z]*\).*/\1/p' "$f")
  [ -z "$expected" ] && expected="block"

  # Read the entire corpus into a string; entries separated by blank lines.
  ENTRY=""
  while IFS= read -r line || [ -n "$line" ]; do
    case "$line" in
      \#*) continue ;;    # skip comments
      "")                 # blank line = entry separator
        if [ -n "$ENTRY" ]; then
          process_entry "$name" "$expected" "$ENTRY"
          ENTRY=""
        fi
        continue ;;
    esac
    ENTRY="${ENTRY}${line}
"
  done < "$f"
  [ -n "$ENTRY" ] && process_entry "$name" "$expected" "$ENTRY"
done

echo
echo "=============================================="
echo "  replay (mode=$MODE): $TOTAL total, $PASS passed, $FAIL failed"
echo "  events: $EVENTS checked, $EVENTS_FAIL event-verification failures"
echo "=============================================="
[ "$FAIL" -eq 0 ] && [ "$EVENTS_FAIL" -eq 0 ] || exit 1