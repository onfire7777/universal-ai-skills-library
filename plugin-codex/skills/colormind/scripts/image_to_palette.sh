#!/usr/bin/env bash
set -euo pipefail

# Sample a small palette from an image using ImageMagick, choose a base color,
# then generate a Colormind palette from that base color.
#
# Usage:
#   image_to_palette.sh <imagePath> [--model ui|default]
#
# Output:
#   JSON with sampled colors and generated Colormind palette

IMG="${1:-}"
if [[ -z "$IMG" || "$IMG" == "-h" || "$IMG" == "--help" ]]; then
  echo "Usage: $(basename "$0") <imagePath> [--model ui|default]" >&2
  exit 2
fi
shift || true

MODEL="ui"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --model)
      MODEL="${2:-ui}"; shift 2;;
    *)
      echo "Unknown arg: $1" >&2; exit 2;;
  esac
done

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Create temp files for intermediate data
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

HIST_FILE="$TEMP_DIR/histogram.txt"
SAMPLED_FILE="$TEMP_DIR/sampled.json"
COLORMIND_FILE="$TEMP_DIR/colormind.json"

# Build histogram using ImageMagick
convert "$IMG" -alpha off -strip -resize 256x256\> -colors 8 -unique-colors \
  -format "%c\n" histogram:info:- \
  | sed -E 's/^[[:space:]]+//' > "$HIST_FILE"

# Parse histogram to JSON
python3 "$SCRIPT_DIR/parse_histogram.py" < "$HIST_FILE" > "$SAMPLED_FILE"

# Extract base RGB color
BASE_RGB=$(python3 "$SCRIPT_DIR/get_base_rgb.py" "$SAMPLED_FILE")

# Generate Colormind palette
node "$SCRIPT_DIR/generate_palette.mjs" --model "$MODEL" \
  --input "$BASE_RGB" N N N N > "$COLORMIND_FILE"

# Combine results
python3 "$SCRIPT_DIR/combine_results.py" "$SAMPLED_FILE" "$COLORMIND_FILE"
