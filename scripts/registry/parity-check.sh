#!/usr/bin/env bash
#
# parity-check.sh — the registry migration cut-over gate.
#
# Proves the Go `skill-router registry build` emits BYTE-IDENTICAL artifacts to
# the legacy Node generator (generate-registry.mjs), for every artifact in both
# the optimize (committed) and faithful (legacy byte-for-byte) modes. Exit 0 only
# when every comparison matches. Run from anywhere; it locates the repo root.
#
#   make parity            # convenience wrapper
#   SKILL_ROUTER_BIN=/path/to/skill-router scripts/registry/parity-check.sh
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

NODE_GEN="scripts/registry/generate-registry.mjs"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

# Resolve the Go binary: use $SKILL_ROUTER_BIN if set, else build a throwaway.
GO_BIN="${SKILL_ROUTER_BIN:-}"
if [ -z "$GO_BIN" ]; then
  GO_BIN="$TMP/skill-router"
  echo "building skill-router (set SKILL_ROUTER_BIN to skip)…"
  ( cd skill-router-cli && go build -o "$GO_BIN" . )
fi

fail=0
compare() { # <label> <mode-flag> <artifact...>
  local label="$1"; shift
  local flag="$1"; shift
  local art
  for art in "$@"; do
    if ! node "$NODE_GEN" "$flag" --print "$art" > "$TMP/node.$art" 2>"$TMP/node.$art.err"; then
      echo "FAIL [$label] node --print $art errored:"; sed 's/^/    /' "$TMP/node.$art.err"; fail=1; continue
    fi
    if ! "$GO_BIN" registry build "$flag" --print "$art" > "$TMP/go.$art" 2>"$TMP/go.$art.err"; then
      echo "FAIL [$label] go registry build --print $art errored:"; sed 's/^/    /' "$TMP/go.$art.err"; fail=1; continue
    fi
    if diff -u "$TMP/node.$art" "$TMP/go.$art" > "$TMP/diff.$art"; then
      printf 'OK   [%-8s] %-18s %8s bytes\n' "$label" "$art" "$(wc -c < "$TMP/node.$art" | tr -d ' ')"
    else
      echo "FAIL [$label] $art differs (Node left / Go right):"; head -40 "$TMP/diff.$art"; fail=1
    fi
  done
}

echo "== optimize (the committed form) =="
compare optimize --optimize manifest marketplace codex-marketplace build-manifest
echo "== faithful (legacy byte-for-byte) =="
compare faithful --faithful manifest marketplace

echo
if [ "$fail" -ne 0 ]; then
  echo "PARITY GATE: FAIL — Go output diverges from the Node generator (see diffs above)."
  exit 1
fi
echo "PARITY GATE: PASS — Go output is byte-identical to the Node generator."
