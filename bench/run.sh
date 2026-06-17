#!/usr/bin/env bash
# Capture a reproducible performance/bloat snapshot of the skill-router + library.
#
#   bench/run.sh before     # baseline on canonical main (run now)
#   bench/run.sh after      # re-run after B2/B4/B5 land; compare deltas
#
# Writes bench/results/<label>.json. Builds the router fresh (default `go build`,
# matching CI) into a temp dir so no binary artifact is committed.
set -euo pipefail

LABEL="${1:-before}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="${REPO:-$(cd "$SCRIPT_DIR/.." && pwd)}"

TMPDIR_BIN="$(mktemp -d)"
trap 'rm -rf "$TMPDIR_BIN"' EXIT
BIN="$TMPDIR_BIN/skill-router-bench"

echo "[bench] building router (go build) -> $BIN"
( cd "$REPO/skill-router-cli" && go build -o "$BIN" . )

OUT="$REPO/bench/results/$LABEL.json"
echo "[bench] measuring (label=$LABEL) ..."
python3 "$REPO/bench/measure.py" --bin "$BIN" --repo "$REPO" --label "$LABEL" --out "$OUT" >/dev/null
echo "[bench] wrote $OUT"
