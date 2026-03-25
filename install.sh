#!/bin/bash
# =============================================================================
# Manus Skills Library Installer
# Installs all 791 curated skills into /home/ubuntu/skills/
# Usage: curl -sSL <raw-url>/install.sh | bash
#   OR:  git clone <repo> /tmp/manus-skills-library && bash /tmp/manus-skills-library/install.sh
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILLS_DIR="/home/ubuntu/skills"

echo "============================================="
echo "  Manus Skills Library Installer"
echo "  791 curated skills for software development"
echo "============================================="
echo ""

# Check if we're running from the cloned repo
if [ -d "$SCRIPT_DIR/skills" ]; then
    echo "[1/3] Installing skills from local repo..."
    cp -r "$SCRIPT_DIR/skills/"* "$SKILLS_DIR/" 2>/dev/null || true
else
    echo "[1/3] Cloning skills repository..."
    REPO_DIR="/tmp/manus-skills-library-install"
    rm -rf "$REPO_DIR"
    gh repo clone manus-skills-library "$REPO_DIR" -- --depth 1 2>/dev/null || \
        git clone --depth 1 "$(git remote get-url origin 2>/dev/null || echo 'https://github.com/USER/manus-skills-library')" "$REPO_DIR"
    echo "[2/3] Copying skills..."
    cp -r "$REPO_DIR/skills/"* "$SKILLS_DIR/" 2>/dev/null || true
    rm -rf "$REPO_DIR"
fi

echo "[2/3] Verifying installation..."
TOTAL=$(find "$SKILLS_DIR" -maxdepth 2 -name "SKILL.md" -type f | wc -l)
EMPTY=$(find "$SKILLS_DIR" -maxdepth 2 -name "SKILL.md" -empty | wc -l)

echo ""
echo "============================================="
echo "  Installation Complete!"
echo "  Total skills: $TOTAL"
echo "  Empty/broken: $EMPTY"
echo "============================================="
echo ""
echo "[3/3] Done! All skills are now active in /home/ubuntu/skills/"
