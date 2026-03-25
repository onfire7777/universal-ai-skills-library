---
name: infographic
description: Generate high-density infographics with structured layout and style choices and bundled generation tooling. Use when the user asks to create an infographic, a visual summary, or a dense single-page information graphic.
metadata: { "pattern": ["generator", "pipeline"], "openclaw": { "emoji": "📊", "primaryEnv": "IMAGE_GEN_API_KEY", "requires": { "env": ["IMAGE_GEN_API_KEY"], "anyBins": ["bun", "npx"], "bins": ["node", "npm"] } } }
---

# Infographic Generation (`infographic`)

## Reference Images (Important)

If you use reference images (image-to-image / series reference / consistency refs):

- Reference images must be public URLs.
- **HTTPS is strongly recommended.**
- `http://` may work but is insecure and can be blocked by some networks.
- Local file paths and `data:` URLs are not supported by the WeryAI gateway.


This skill turns complex source material into a single high-density infographic.

Maintain the layout and style mapping when the bundled generation runtime updates, and re-check recent runtime behavior if a newer version changes prompt handling.

Script:

- `scripts/scaffold.ts`
- `scripts/build-prompt.ts`
- `scripts/build-batch.ts`

## Safety & Scope

- **Network**: This skill calls the WeryAI gateway over HTTPS (`https://api.weryai.com`).
- **Auth**: Uses `IMAGE_GEN_API_KEY`. The key is never printed. It may be persisted **only** when you explicitly run `npm run setup -- --persist-api-key`.
- **Reference images**: Must be public URLs (`https://` recommended). `http://` may work but is insecure. Local file paths and `data:` URLs are rejected.
- **No arbitrary shell**: The generation runtime does not execute arbitrary shell commands.
- **Files written**: Output images and optional local config under `.image-skills/infographic/` (project) and/or `~/.image-skills/infographic/` (home).


## Use Cases

- high-density information graphics
- visual summaries
- data relationship graphics
- process posters
- content that needs to be explained in one image

Not a good fit for:

- a single article cover image
- multi-card RedNote content
- a sequential comic narrative

## Core Dimensions

1. `layout`: the information structure
2. `style`: the visual language

See:

- [references/analysis-framework.md](references/analysis-framework.md)
- [references/structured-content-template.md](references/structured-content-template.md)
- [references/layouts-and-styles.md](references/layouts-and-styles.md)
- [references/prompt-template.md](references/prompt-template.md)

## Commands

| Script | Purpose |
| --- | --- |
| `scripts/scaffold.ts` | Initialize `analysis.md`, `structured-content.md`, and prompts |
| `scripts/build-batch.ts` | Generate `batch.json` from prompt variants |
| `scripts/build-prompt.ts` | Build final prompt from `structured-content.md` |
| `npm run generate` | Generate the infographic |
| `./scripts/vendor/compression-runtime/scripts/main.ts` | Compress output for delivery |

## Workflow

### Step 1: Initialize Working Files

Create the working directory and starter files:

```bash
${BUN_X} {baseDir}/scripts/scaffold.ts \
  --output-dir infographic/topic-slug \
  --topic "Topic name" \
  --lang en
```

This creates:

- `analysis.md`
- `structured-content.md`
- `prompts/infographic.md`

### Step 2: Understand the Source

Extract:

- the main topic
- the target audience
- the density of information
- whether the content is driven by data, lists, hierarchy, process, or comparison
- the user's language, especially for titles, labels, and annotations

Save the result to `analysis.md`, using:

- [references/analysis-framework.md](references/analysis-framework.md)

### Step 3: Organize Structured Content

Before selecting the final layout and style, rewrite the material into `structured-content.md`:

- [references/structured-content-template.md](references/structured-content-template.md)

This step should make titles, sections, labels, and data points explicit.

### Step 4: Choose `layout` and `style`

Default priorities:

- `layout`: `bento-grid`
- `style`: `infographic`

Recommended rules:

- step-by-step tutorial -> `linear-progression`
- comparison -> `binary-comparison`
- hierarchy -> `hierarchical-layers`
- broad knowledge overview -> `bento-grid`
- dense modular content -> `dense-modules`

If the user explicitly asks for a specific visual direction, follow that preference.

### Step 5: Build `batch.json` for Variants

If you want to explore multiple infographic directions, place multiple prompt files in `prompts/`, for example:

- `01-technical.md`
- `02-minimal.md`
- `03-bold.md`

Then generate a batch file:

```bash
${BUN_X} {baseDir}/scripts/build-batch.ts \
  --prompts infographic/topic-slug/prompts \
  --output infographic/topic-slug/batch.json \
  --images-dir infographic/topic-slug \
  --model "$M" \
  --jobs 3
```

The script reads `Style direction:` and `Aspect ratio:` from each prompt file when possible, then maps them into generation task fields.

### Step 6: Build the Final Prompt

Convert `structured-content.md` into the final prompt:

```bash
${BUN_X} {baseDir}/scripts/build-prompt.ts \
  --structured-content infographic/topic-slug/structured-content.md \
  --output infographic/topic-slug/prompts/infographic.md \
  --audience "product teams" \
  --style technical-schematic \
  --aspect 4:3 \
  --lang en
```

This keeps the prompt synchronized with the structured content instead of rewriting it by hand.

### Step 7: Map to the Bundled Runtime

The bundled image runtime currently exposes one structured style argument, `--style`, so:

- map `style` to `--style`
- write `layout`, labels, hierarchy, and information relationships into the prompt body
- prefer `3:4`, `4:3`, or `16:9` depending on content shape

Recommended mapping:

| infographic style | runtime `--style` |
| --- | --- |
| `craft-handmade` | `watercolor` |
| `chalkboard` | `chalk` |
| `corporate-memphis` | `flat-illustration` |
| `technical-schematic` | `infographic` |
| `bold-graphic` | `poster` |
| `storybook-watercolor` | `watercolor` |
| `retro-pop-grid` | `poster` |
| `minimal` | `editorial` |

Use [references/prompt-template.md](references/prompt-template.md) and make sure to:

- define the main title, subtitle, sections, and labels explicitly
- include numbers, years, and terminology directly when possible
- emphasize hierarchy, whitespace, and readability instead of abstract style-only wording
- state the target language explicitly if the image needs on-canvas text

### Step 8: Run Generation

On first use in a new project, run `npm run ensure-ready -- --project <your-project> --workflow infographic` from this skill directory before generation. This reads the doctor report and auto-runs `bootstrap` if local script dependencies are still missing. If the report shows a missing `IMAGE_GEN_API_KEY` and the user approves, run `npm run setup -- --project <your-project> --workflow infographic --persist-api-key` when the key is already in env, or persist it to `.image-skills/infographic/.env` on the user's behalf, then continue without leaving this workflow.

When this skill is first connected, tell the user that the default generation model is **Nano Banana 2** (`GEMINI_3_1_FLASH_IMAGE`). Also tell them it can be switched later whenever another model fits the task better.

```bash
${BUN_X} {baseDir}/npm run generate \
  --promptfiles prompts/infographic.md \
  --style infographic \
  --image infographic/topic-slug/infographic.png \
  --ar 3:4 \
  -m "$M"
```

Batch example:

```bash
${BUN_X} {baseDir}/npm run generate \
  --batchfile infographic/topic-slug/batch.json \
  --jobs 3
```

If the user has not chosen a model yet, follow this skill's model-selection rules first.

## Output Convention

Suggested output directory:

```text
infographic/<topic-slug>/
```

Suggested minimum files:

- `analysis.md`
- `structured-content.md`
- `prompts/infographic.md`
- `batch.json`
- `infographic.png`

## Re-run Behavior

- `scaffold.ts` on an existing directory overwrites `analysis.md`, `structured-content.md`, and the starter prompt files.
- `build-prompt.ts` overwrites the target prompt file each time it runs.
- `build-batch.ts` overwrites `batch.json`.
- Re-running the base generator with the same `--image` path overwrites the existing infographic output.
- If multiple prompt variants are listed in `batch.json`, re-running the batch regenerates all listed variants unless you remove the ones you want to keep.

## Definition of Done

- `analysis.md`, `structured-content.md`, and `prompts/infographic.md` exist in the output directory.
- The generated infographic matches the chosen `layout / style / aspect`.
- The image is shown directly to the user with a summary of the parameters used.
- A compressed webp version is produced for delivery.

## Iteration

When the user wants changes after seeing the generated infographic:

- **Layout issues** ("too crowded", "unclear hierarchy") → switch `layout` in `structured-content.md` and rebuild the prompt. Try a different layout pattern.
- **Style mismatch** ("too busy", "not professional enough") → change `style` / `--style`, re-generate.
- **Content issues** ("missing X", "remove this section") → revise `structured-content.md`, rebuild prompt with `build-prompt.ts`, re-generate.
- **Text readability** ("text too small", "can't read it") → reduce information density in the prompt, or switch to a less text-heavy layout.
- **Want to explore alternatives** → create multiple prompt variants in `prompts/`, batch-generate to compare layouts and styles.

Infographics are single-image outputs, so each iteration only re-generates one image.

## Delivery

When the infographic is ready:

1. **Show the image directly** — do not just print a file path.
2. Briefly state: layout type, style, aspect ratio.
3. Ask if the user wants changes (layout, content, style) or is satisfied.
4. **Auto-compress**: once confirmed, run the bundled compression runtime to produce a webp version.

```bash
${BUN_X} {baseDir}/./scripts/vendor/compression-runtime/scripts/main.ts infographic/topic-slug/infographic.png -f webp -q 85
```

If multiple variants were generated via batch, show all of them and let the user pick the best one.

Internal checklist (for agent): `layout / style / aspect`, information division, model, text language, compression done.
