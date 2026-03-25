---
name: bundled-compression-runtime
description: Compress or convert generated images for delivery, preview, and social upload. Use when the user asks to compress an image, optimize image size, convert to webp, or reduce file size after generation.
metadata: { "pattern": ["tool-wrapper"], "openclaw": { "emoji": "📦", "requires": { "anyBins": ["sips", "cwebp", "magick"] } } }
---

# Bundled Compression Runtime

Compress or convert generated images for delivery, preview, and social upload. Optimize image size, convert to webp, or reduce file size after generation.

This skill uses the best available local toolchain and stays self-contained.

Script:

- `scripts/main.ts`

## Use Cases

- files are too large to share easily
- output needs to be converted to `webp`
- a lighter preview image is needed
- images should be optimized before social upload

## Example Prompts

- `Compress my cover image to webp for social upload`
- `Optimize all images in the output directory to reduce file size`
- `Convert this PNG to a lighter webp format at quality 80`

**Important**: at least one of `cwebp`, `magick`, or `sips` must be available locally.

## Backend Preference

Preferred order:

1. `cwebp` for WebP output
2. `magick` or `convert` for conversion and flexible processing
3. `sips` as a macOS-native fallback

## CLI

```bash
${BUN_X} {baseDir}/scripts/main.ts <input> [options]
```

## Options

| Option | Description |
| --- | --- |
| `<input>` | input file or directory |
| `--output`, `-o` | output file or directory |
| `--format`, `-f` | `webp`, `png`, or `jpeg` (default: `webp`) |
| `--quality`, `-q` | quality `0-100` (default: `80`) |
| `--keep`, `-k` | keep the original when the output extension matches the input |
| `--recursive`, `-r` | process directories recursively |
| `--json` | return JSON output |

## Common Examples

Single image to WebP:

```bash
${BUN_X} {baseDir}/scripts/main.ts input.png -f webp -q 80
```

Single image to JPEG:

```bash
${BUN_X} {baseDir}/scripts/main.ts input.png -f jpeg -q 82
```

Recursive directory compression:

```bash
${BUN_X} {baseDir}/scripts/main.ts ./images -r -f webp
```

JSON output:

```bash
${BUN_X} {baseDir}/scripts/main.ts input.png --json
```

## Workflow

1. Identify the input: a single image file or a directory of images.
2. Choose the target format (`webp`, `png`, or `jpeg`) and quality level.
3. Run the CLI with the chosen options.
4. Verify the output exists and is smaller than the input.
5. If processing a directory, use `--recursive` to handle nested folders.

## Output Convention

Recommended behavior:

- keep the original image untouched
- write the compressed result under a new name
- store results in the same directory or in a `compressed/` subdirectory

Example:

```text
cover-image/topic-slug/
├── cover.png
└── compressed/
    └── cover.webp
```

## Re-run Behavior

- Re-compressing the same input produces a new output file (or overwrites an existing one at the same path).
- The original input is never modified unless `--keep` is used and the output extension matches the input.
- `--recursive` re-processes all matching files in the directory, including previously compressed ones.
- To avoid double-compression, point `--output` to a separate `compressed/` subdirectory.

## Definition of Done

- The compressed output file exists and is smaller than the input.
- The original file is preserved (not overwritten) unless `--keep` is used with matching extensions.
- Output format matches the requested `--format`.
- `--json` mode produces a valid JSON summary with input size, output size, and compression ratio.

## Current Behavior

- prefers `cwebp` -> `magick` / `convert` -> `sips`
- supports single files and directories
- supports recursive scanning
- supports JSON summaries
- does not delete the original just because the format changes; use `--keep` when the output extension matches the input
- `webp` output requires `cwebp` or ImageMagick; if only `sips` is available, prefer `png` or `jpeg`
