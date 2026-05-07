#!/bin/bash
# =============================================================================
# Universal AI Skills Library Installer
# Installs 14 core skills + 770 library skills into the target skills directory
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
echo "  14 Core Skills + 770 Library Skills"
echo "============================================================"
echo ""
echo "  Target directory: $TARGET_DIR"
echo ""

mkdir -p "$TARGET_DIR"

CORE_SKILLS=(
    chat-summarizer
    context-anchor
    file-organizer
    manus-api
    model-selector
    multi-model-code-auditor
    multi-model-oracle
    music-prompter
    persistent-computing
    prompt-engineer
    skill-creator
    skill-debugger
    skill-sync
    ultimate-skill-creator
)

echo "[1/4] Installing 14 core skills..."
CORE_COUNT=0
for skill in "${CORE_SKILLS[@]}"; do
    if [ -d "$SCRIPT_DIR/$skill" ]; then
        cp -r "$SCRIPT_DIR/$skill" "$TARGET_DIR/"
        CORE_COUNT=$((CORE_COUNT + 1))
    fi
done
echo "      Done: $CORE_COUNT core skills installed"

echo "[2/4] Installing 770 library skills..."
if [ -d "$SCRIPT_DIR/skills" ]; then
    cp -r "$SCRIPT_DIR/skills/"* "$TARGET_DIR/"
    LIB_COUNT=$(find "$SCRIPT_DIR/skills" -maxdepth 1 -mindepth 1 -type d | wc -l)
    echo "      Done: $LIB_COUNT library skills installed"
else
    echo "      ERROR: skills/ directory not found. Run from repo root."
    exit 1
fi

echo "[3/4] Verifying installation..."
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

echo "[4/4] Post-install notes:"
echo "  - Skills are now active in: $TARGET_DIR"
echo "  - For Windows cross-platform install: infrastructure/scripts/install_skills.ps1"
echo "  - For MCP bridge setup: infrastructure/scripts/setup_mcp_bridges.ps1"
echo ""
echo "Done. Your AI agents now have access to all skills."
