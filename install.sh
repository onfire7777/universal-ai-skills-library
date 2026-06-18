#!/usr/bin/env bash
#
# Universal AI Skills Library installer for Linux, macOS, and WSL.
#
# This installer is intentionally router-first: it builds the `skill-router`
# CLI and leaves the skill corpus in this repository. Do not copy all skills
# into every AI client unless you are intentionally building an offline bundle.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="${GOBIN:-"$HOME/go/bin"}"
COPY_SKILLS_DIR=""
SKIP_VALIDATE=0
SYNC_CODEX=0
SYNC_CLAUDE=0

usage() {
  cat <<'USAGE'
Universal AI Skills Library installer

Usage:
  bash install.sh [options]

Options:
  --bin-dir DIR        Install skill-router into DIR (default: $GOBIN or ~/go/bin)
  --copy-skills DIR    Optional offline/full-copy export of skills/ into DIR
  --sync-codex         Install compact wrapper into ~/.codex/skills
  --sync-claude        Install compact wrapper into ~/.claude/skills
  --sync-cli-clients   Install compact wrappers for both Codex and Claude
  --skip-validate      Skip manifest validation after build
  -h, --help           Show this help

Default install:
  git clone https://github.com/onfire7777/universal-ai-skills-library.git
  cd universal-ai-skills-library
  bash install.sh

After install:
  skill-router skill search debugging
  skill-router skill universal-ai-skills
  skill-router skills validate-manifest

Notes:
  - The normal install does not duplicate the 1,812-skill corpus.
  - AI clients should load skills on demand through skill-router.
  - Windows users should prefer ./install.ps1 for the full local AI stack.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bin-dir)
      BIN_DIR="${2:?--bin-dir requires a directory}"
      shift 2
      ;;
    --copy-skills)
      COPY_SKILLS_DIR="${2:?--copy-skills requires a directory}"
      shift 2
      ;;
    --sync-codex)
      SYNC_CODEX=1
      shift
      ;;
    --sync-claude)
      SYNC_CLAUDE=1
      shift
      ;;
    --sync-cli-clients)
      SYNC_CODEX=1
      SYNC_CLAUDE=1
      shift
      ;;
    --skip-validate)
      SKIP_VALIDATE=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: Go is required to build skill-router. Install Go, then rerun this script." >&2
  exit 1
fi

mkdir -p "$BIN_DIR"

echo "==> Building skill-router"
(
  cd "$SCRIPT_DIR/skill-router-cli"
  go build -o "$BIN_DIR/skill-router" .
)

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
  echo "NOTE: $BIN_DIR is not on PATH for this shell."
  echo "      Add this to your shell profile: export PATH=\"$BIN_DIR:\$PATH\""
fi

if [[ -n "$COPY_SKILLS_DIR" ]]; then
  echo "==> Exporting full skill corpus to $COPY_SKILLS_DIR"
  mkdir -p "$COPY_SKILLS_DIR"
  cp -R "$SCRIPT_DIR/skills/." "$COPY_SKILLS_DIR/"
fi

if [[ "$SKIP_VALIDATE" -eq 0 ]]; then
  echo "==> Validating manifest"
  (
    cd "$SCRIPT_DIR"
    "$BIN_DIR/skill-router" skills validate-manifest
  )
fi

if [[ "$SYNC_CODEX" -eq 1 ]]; then
  echo "==> Installing Codex CLI wrapper"
  (
    cd "$SCRIPT_DIR"
    "$BIN_DIR/skill-router" sync codex
  )
fi

if [[ "$SYNC_CLAUDE" -eq 1 ]]; then
  echo "==> Installing Claude CLI wrapper"
  (
    cd "$SCRIPT_DIR"
    "$BIN_DIR/skill-router" sync claude
  )
fi

cat <<EOF

Universal AI Skills Library installed.

CLI:
  $BIN_DIR/skill-router

Try:
  skill-router skill search debugging
  skill-router skill universal-ai-skills
  skill-router doctor

EOF
