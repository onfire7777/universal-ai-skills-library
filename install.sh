#!/bin/bash
# =============================================================================
# Universal AI Skills Library Installer
# Installs the unified 785-skill library into the target skills directory
#
# Usage:
#   From cloned repo:  bash install.sh [--target /path/to/skills]
#   One-liner:         git clone https://github.com/onfire7777/manus-skills-library.git /tmp/msl && bash /tmp/msl/install.sh
#
# Supported platforms: Linux, macOS, WSL
# Target AI platforms: Claude Code, Codex CLI, OpenCode, Gemini CLI, Manus
# =============================================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="/home/ubuntu/skills"

# Parse --target flag
while [[ $# -gt 0 ]]; do
    case $1 in
        --target) TARGET_DIR="$2"; shift 2;;
        *) shift;;
    esac
done

echo "============================================================"
echo "  Universal AI Skills Library Installer"
echo "  785 Unified Skills"
echo "============================================================"
echo ""
echo "  Target directory: $TARGET_DIR"
echo ""

mkdir -p "$TARGET_DIR"

echo "[1/3] Installing unified skills..."
if [ -d "$SCRIPT_DIR/skills" ]; then
    cp -r "$SCRIPT_DIR/skills/"* "$TARGET_DIR/"
    SKILL_COUNT=$(find "$SCRIPT_DIR/skills" -maxdepth 1 -mindepth 1 -type d | wc -l)
    echo "      Done: $SKILL_COUNT skills installed"
else
    echo "      ERROR: skills/ directory not found. Run from repo root."
    exit 1
fi

echo "[2/3] Verifying installation..."
TOTAL=$(find "$TARGET_DIR" -maxdepth 2 -name "SKILL.md" -type f | wc -l)
EMPTY=$(find "$TARGET_DIR" -maxdepth 2 -name "SKILL.md" -empty | wc -l)

echo ""
echo "============================================================"
echo "  Installation Complete!"
echo "  Total skills installed: $TOTAL"
echo "  Empty/broken: $EMPTY"
echo "============================================================"
echo ""

if [ "$EMPTY" -gt 0 ]; then
    echo "  Warning: $EMPTY skills have empty SKILL.md files."
    echo ""
fi

echo "[3/3] Post-install notes:"
echo "  - Skills are now active in: $TARGET_DIR"
echo "  - Use CLI-first loading: manus skill <name>"
echo "  - Keep agent instruction files as indexes, not full skill bodies"
echo "  - For Windows cross-platform install: infrastructure/scripts/install_skills.ps1"
echo "  - For MCP bridge setup: infrastructure/scripts/setup_mcp_bridges.ps1"
echo ""
echo "Done. Your AI agents now have access to all skills."
