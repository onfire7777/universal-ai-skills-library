# Gemini Delegation Reference

Purpose: load this when the user explicitly asks for Gemini delegation or when exploratory single-sprite SVG generation is faster than hand-authored pixel placement.

## Contents

- [When to use Gemini](#when-to-use-gemini)
- [Execution flow](#execution-flow)
- [Capability limits](#capability-limits)
- [Prompt template](#prompt-template)
- [CLI command](#cli-command)
- [Sanitize step](#sanitize-step)
- [Naming convention](#naming-convention)
- [Limitations](#limitations)
- [Direct vs delegation](#direct-vs-delegation)

## When to Use Gemini

| Situation | Use Gemini CLI | Use Dot Direct |
|-----------|:--------------:|:--------------:|
| quick prototype or variation | Yes | - |
| strict pixel placement is required | - | Yes |
| the user says `gemini`, `delegate`, or external generation | Yes | - |
| spritesheet or animation work | - | Yes |
| exploratory single-sprite generation | Yes | - |

## Execution Flow

`PLAN -> PALETTE -> PROMPT -> GEMINI CLI -> SANITIZE -> PREVIEW`

1. Define asset requirements and palette as usual.
2. Embed Dot constraints into the prompt.
3. Run `gemini '<prompt>' > output.svg`.
4. Extract the raw SVG only.
5. Run the usual pixel-perfect preview checks.

## Capability Limits

| Grid | Rect count (estimated) | Gemini feasibility | Notes |
|------|------------------------|:------------------:|-------|
| `8x8` | `~30-50` | Excellent | highest reliability |
| `16x16` | `~100-180` | Good | standard use case |
| `32x32` | `~400-700` | Fair | use run-length compression prompts |
| `64x64` | `~1500-3000` | Poor | close to token limits; prefer splitting or Pillow |
| `128x128` | `~8000+` | Not recommended | use Dot direct or Pillow |

Rules:

- For `32x32+`, explicitly ask for run-length `<rect width="N">` compression.
- For `64x64+`, recommend switching to Pillow script generation.
- Do not use Gemini for spritesheet consistency or frame-by-frame animation.

## Prompt Template

The prompt must include all of the following constraints:

```text
- Canvas: exactly {W}x{H} grid
- Palette: exactly {N} colors (list hex values)
- Every pixel is a <rect x="N" y="N" width="1" height="1" fill="#hex"/>
- Merge consecutive same-color pixels into wider rects (e.g., width="5") to reduce element count
- Use -1 (transparent) for empty pixels; do NOT emit a rect for them
- shape-rendering="crispEdges" on the <svg> element
- style="image-rendering: pixelated;" on the <svg> element
- viewBox="0 0 {W} {H}", display size width="{DISPLAY}" height="{DISPLAY}"
- NO anti-aliasing, NO gradients, NO filters, NO rounded corners
- Return ONLY the complete SVG code. No markdown, no explanation, no code fences.
```

Display size must use integer scaling. Example: `32x32` -> `width="256" height="256"` at `8x`.

## CLI Command

```bash
gemini '<prompt with Dot constraints>' > {name}-{W}x{H}-gemini.svg
```

## Sanitize Step

Gemini output may include logs or markdown fences. Extract the SVG body only:

```bash
# Extract only the SVG content
grep -o '<svg.*</svg>' output.svg > clean.svg && mv clean.svg output.svg
```

## Naming Convention

Use the `-gemini` suffix for delegated output:

- Dot direct: `goblin-16x16.svg`
- Gemini delegation: `goblin-16x16-gemini.svg`

## Limitations

- Gemini may not fully honor prompt constraints.
- The CLI may try a policy-rejected image path before falling back to text generation.
- Spritesheet and animation consistency are not guaranteed.
- Use Dot direct when exact pixel placement matters.

## Direct vs Delegation

| Aspect | Dot Direct | Gemini CLI |
|--------|------------|------------|
| pixel precision | full control | best effort |
| speed | slower | faster |
| consistency | high | medium to low |
| variation | manual | easy via prompt changes |
| spritesheet support | supported | not recommended |
| dependency | none | `gemini` CLI |
