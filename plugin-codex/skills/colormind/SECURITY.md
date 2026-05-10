# Security Notes

## What This Skill Does

The Colormind skill generates color palettes using the Colormind.io API. It has two main functions:

1. **Generate random palettes** - Makes HTTP requests to colormind.io API
2. **Sample colors from images** - Uses ImageMagick to extract dominant colors from an image, then generates a harmonious palette

## Why Security Scanners May Flag This

The `image_to_palette.sh` script may trigger false positives because it:

- Uses **ImageMagick `convert`** command (historically had CVEs, but we use safe operations only)
- Creates **temporary files** for data passing between scripts
- Uses **Python scripts** to parse color data

## What It Does NOT Do

- ❌ Download or execute external code
- ❌ Access credentials or environment variables
- ❌ Write outside the workspace
- ❌ Make network calls except to colormind.io API
- ❌ Execute user-provided code
- ❌ Access sensitive files

## Architecture (Refactored 2026-02-18)

Previously used inline Python with heredocs, which triggered security scanners. Now uses:

```
image_to_palette.sh
  ├─> ImageMagick (extract histogram)
  ├─> parse_histogram.py (parse colors)
  ├─> get_base_rgb.py (extract base color)
  ├─> generate_palette.mjs (call Colormind API)
  └─> combine_results.py (merge results)
```

All data is passed via temporary JSON files in a secure temp directory that's cleaned up on exit.

## Verification

You can audit the code yourself:
- All scripts are in `scripts/`
- No obfuscation, no eval, no remote code execution
- Simple, readable Python and bash

## Dependencies

- **Node.js** (18+) - for HTTP requests to Colormind API
- **ImageMagick** (`convert`) - for image color sampling (use recent patched version)
- **Python 3** - for JSON parsing

All are standard development tools.

## Known Limitations & Privacy Considerations

### 1. Unencrypted Transport
The Colormind.io API only supports **plain HTTP** (their HTTPS endpoint has a self-signed certificate). This means:
- Color palette data is sent unencrypted
- Network observers can see what colors you're generating
- **Mitigation:** Don't use with sensitive color schemes or proprietary brand colors

### 2. External Data Sharing
When using `image_to_palette.sh`:
- The script extracts dominant colors from your local images
- These RGB values are sent to colormind.io to generate a harmonious palette
- **Do not use with:**
  - Private photos
  - Sensitive documents
  - Proprietary designs
  - Any image you wouldn't want color data extracted from

### 3. ImageMagick Safety
The script uses ImageMagick's `convert` command, which has had security vulnerabilities:
- **Use a recent, patched ImageMagick version**
- Consider running in a sandboxed environment
- Avoid processing untrusted images from external sources

### 4. Automated Use
If running this in an automated agent:
- Add explicit user consent prompts before processing images
- Restrict which images can be processed (e.g., only from specific directories)
- Log what images were processed and when

## Recommendations

**For public/non-sensitive use:** This skill works as intended for general palette generation.

**For production/sensitive use:**
- Host your own Colormind-compatible API with HTTPS
- Or use a different palette generation service with proper TLS
- Consider client-side palette generation libraries that don't require external APIs
