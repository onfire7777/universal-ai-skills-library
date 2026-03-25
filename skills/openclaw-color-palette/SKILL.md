---
name: openclaw-color-palette
version: 2.0.0
description: "Design-grade color processing suite — generate harmonious palettes from any seed color, check WCAG 2.1 AA/AAA contrast accessibility, convert between hex/RGB/HSL/HSV/CMYK, simulate 8 color blindness types, and extract dominant colors from images. Use when: (1) user says 'generate a color palette' or 'give me a color scheme for my website', (2) user asks 'check contrast' or 'is this readable' or 'WCAG accessibility', (3) user needs to 'convert hex to RGB' or 'what is #FF5733 in HSL', (4) user wants to 'simulate color blindness' or 'how do colorblind users see this', (5) user requests 'extract colors from image' or 'pull brand colors from this logo'. 6 harmony algorithms, 140 CSS named colors, k-means clustering. Zero install, sub-100ms on Cloudflare Workers. Free + Pro $9/mo."
read_when:
  - "User asks to generate a color palette, color scheme, or harmonious colors for a project"
  - "User needs to check text/background contrast ratio or WCAG accessibility compliance"
  - "User wants to convert a color between formats: hex, RGB, HSL, HSV, CMYK, or CSS named colors"
  - "User asks how colors appear to colorblind users or wants a color blindness simulation"
  - "User wants to extract dominant colors from an image, logo, or screenshot"
commands:
  - generate_palette
  - check_contrast
  - convert_color
  - color_blindness_sim
  - extract_palette
metadata:
  openclaw:
    emoji: "\U0001F3A8"
    homepage: https://github.com/yedanyagamiai-cmd/openclaw-mcp-servers
    always: false
---

# OpenClaw Color Palette v2.0

**Design-grade color tools -- 5 tools for palettes, accessibility, conversion, simulation, and extraction.**

## Quick Reference

| Situation | Action | Tool |
|-----------|--------|------|
| Starting a new website or app design | Generate a harmonious palette from a brand color | `generate_palette` |
| Text looks hard to read on a background | Check WCAG 2.1 AA/AAA contrast ratio | `check_contrast` |
| Designer gave hex, you need RGB for code | Convert between any color format | `convert_color` |
| Need to ensure inclusive design | Simulate 8 types of color vision deficiency | `color_blindness_sim` |
| Client sent a logo, need brand colors | Extract dominant colors with frequency data | `extract_palette` |
| Building a design system from scratch | Generate palette + check contrast + simulate CVD | All 3 combined |

## What's New in v2.0

- :art: **HarmonyEngine** -- 6 color theory algorithms: complementary, analogous, triadic, split-complementary, tetradic, monochromatic. Mathematically precise palette generation from any seed color.
- :shield: **AccessGuard Protocol** -- Full WCAG 2.1 AA and AAA checking for normal and large text. When contrast fails, auto-suggests the nearest accessible color pair with minimal visual change.
- :eye: **VisionSim 8-Type** -- Simulate protanopia, deuteranopia, tritanopia, achromatopsia, protanomaly, deuteranomaly, tritanomaly, and achromatomaly. Warns when palette colors become indistinguishable.
- :dropper: **ChromaExtract** -- K-means clustering extracts 5-10 dominant colors from any image URL or base64 data, sorted by area percentage.

## MCP Quick Start

```json
{
  "openclaw-colors": {
    "type": "streamable-http",
    "url": "https://color-palette-mcp.yagami8095.workers.dev/mcp"
  }
}
```

Add to Claude Desktop, Cursor, Windsurf, VS Code, or any MCP-compatible client. No API key needed for Free tier. Works immediately.

## Detection Triggers

This skill activates when you say:

- "generate a color palette" / "color scheme for my website" / "give me colors for"
- "check contrast" / "is this readable" / "WCAG accessibility" / "contrast ratio"
- "convert hex to RGB" / "convert color" / "what is #FF5733 in HSL" / "hex to CMYK"
- "simulate color blindness" / "how do colorblind users see this" / "protanopia simulation"
- "extract colors from image" / "pull palette from logo" / "brand colors from screenshot"
- "complementary colors" / "analogous palette" / "triadic" / "split-complementary" / "monochromatic"

## Tools (5)

### `generate_palette` -- Harmony-Based Palettes (HarmonyEngine)

Generate mathematically harmonious color palettes from any seed color using 6 color theory algorithms.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `color` | string | required | Seed color in any format: hex, RGB, HSL, or CSS name |
| `harmony` | string | "complementary" | Algorithm: `complementary`, `analogous`, `triadic`, `split-complementary`, `tetradic`, `monochromatic` |
| `count` | number | 5 | Number of colors to generate (3-8) |

**Output**: Array of colors, each with `hex`, `rgb`, `hsl`, `name` (nearest CSS name), plus ready-to-paste CSS custom properties (`--color-primary`, `--color-secondary`, etc.).

#### WRONG vs RIGHT

**WRONG** -- Picking colors by gut feeling in a color picker:
```
User: I need a color palette for my SaaS dashboard
Agent: *suggests random colors* "#3498db, #e74c3c, #2ecc71" — no harmony, no theory, may clash
```

**RIGHT** -- Use `generate_palette`:
```
User: I need a color palette based on #3498db
Agent: calls generate_palette with triadic harmony → returns 5 mathematically balanced colors with hex/RGB/HSL, CSS variables, and named descriptors
```

---

### `check_contrast` -- WCAG Accessibility Checker (AccessGuard Protocol)

Check foreground/background color contrast ratio against WCAG 2.1 standards for both AA and AAA levels, normal and large text.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `foreground` | string | required | Text color (any format) |
| `background` | string | required | Background color (any format) |

**Output**: `{ ratio: "4.5:1", AA_normal: "pass", AA_large: "pass", AAA_normal: "fail", AAA_large: "pass", suggestion: { foreground: "#...", background: "#..." } }`

When contrast fails, `suggestion` provides the nearest accessible color pair that passes with minimal visual change from the original colors.

#### WRONG vs RIGHT

**WRONG** -- Eyeballing if light gray text on white is readable:
```
User: Is #999999 text readable on white?
Agent: *guesses* "It looks okay to me" — fails WCAG AA at 2.85:1 (needs 4.5:1)
```

**RIGHT** -- Use `check_contrast`:
```
User: Is #999999 text readable on white?
Agent: calls check_contrast → "Ratio 2.85:1 — FAILS AA normal text (needs 4.5:1). Suggestion: darken to #767676 for AA pass at 4.54:1"
```

---

### `convert_color` -- Universal Format Conversion (ChromaBridge)

Convert any color between hex, RGB, HSL, HSV, CMYK, and all 140 CSS named colors. Auto-detects input format.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `color` | string | required | Color in any recognized format |
| `format` | string | "all" | Target format: `hex`, `rgb`, `hsl`, `hsv`, `cmyk`, `all` |

**Auto-detected inputs**: `#FF5733`, `#F53`, `rgb(255,87,51)`, `rgba(255,87,51,0.8)`, `hsl(11,100%,60%)`, `hsv(11,80%,100%)`, `cmyk(0,66,80,0)`, `tomato`, `rebeccapurple`.

**Output**: All formats simultaneously plus `nearestCssName`, `luminance`, `isLight` (boolean for text color selection).

---

### `color_blindness_sim` -- Vision Deficiency Simulation (VisionSim 8-Type)

Simulate how any color or full palette appears to users with 8 different types of color vision deficiency.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `colors` | string[] | required | Array of colors to simulate (any format) |
| `types` | string[] | ["all"] | CVD types to simulate, or `"all"` for all 8 |

**8 types**:
- **Protanopia** -- no red cones (1% of males)
- **Deuteranopia** -- no green cones (1% of males)
- **Tritanopia** -- no blue cones (0.003% of population)
- **Achromatopsia** -- complete color blindness (rare)
- **Protanomaly** -- weak red cones (1% of males)
- **Deuteranomaly** -- weak green cones (5% of males, most common CVD)
- **Tritanomaly** -- weak blue cones (very rare)
- **Achromatomaly** -- weak all cones (rare)

**Output**: For each color, simulated hex values for each CVD type plus `warnings[]` where palette colors become indistinguishable (Delta E < 5).

---

### `extract_palette` -- Image Color Extraction (ChromaExtract)

Extract dominant colors from an image URL or base64 data using k-means clustering algorithm.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `image` | string | required | Image URL or base64 data string |
| `count` | number | 5 | Number of dominant colors to extract (3-10) |

**Output**: Array of colors sorted by area percentage, each with `hex`, `rgb`, `percentage`, `name` (nearest CSS name). Total coverage percentage included.

**Use for**: Brand color analysis, design inspiration, mood board creation, logo color extraction, competitive design analysis.

## What NOT to Do

| Don't | Why | Do Instead |
|-------|-----|------------|
| Expect print-accurate CMYK values | Mathematical conversion only, no ICC profiles | Use a color management system for print |
| Send images over 5MB | Edge workers have memory limits | Resize or compress images first |
| Use CVD simulation for medical diagnosis | Illustrative only, not clinically calibrated | Refer to ophthalmology resources |
| Rely solely on contrast ratio for accessibility | Color is one aspect of accessibility | Also consider font size, weight, and spacing |
| Input colors with alpha in CMYK conversion | CMYK has no alpha channel concept | Remove alpha or use RGB/HSL output |

## Security & Privacy

| What | Status |
|------|--------|
| Reads your color values and images | Yes -- for processing only |
| Stores your data or images | **No** -- zero persistence, stateless edge processing |
| Sends data to third parties | **No** -- processed entirely on Cloudflare Workers |
| Logs request content | **No** -- only anonymous usage counters |
| Downloads images to server | Fetched transiently for extraction, never stored |
| Requires authentication | Free tier: no. Pro: API key header only |

## Pricing

| Tier | Calls/Day | Price | What You Get |
|------|-----------|-------|--------------|
| Free | 25 | $0 | All 5 color tools, no signup |
| Pro | 1,000 | $9/mo | All 9 OpenClaw servers (49 tools), priority routing |
| x402 | Pay-per-call | $0.05 USDC | No account needed, crypto-native |

**Get Pro Key**: https://buy.stripe.com/4gw5na5U19SP9TW288

## The OpenClaw Intelligence Stack

| Server | Tools | Best For |
|--------|-------|----------|
| **JSON Toolkit** | 6 | Format, validate, diff, query, transform JSON |
| **Regex Engine** | 5 | Test, extract, replace, explain regex patterns |
| **Color Palette** | 5 | Generate, convert, harmonize, accessibility-check colors |
| **Timestamp Converter** | 5 | Parse, format, diff, timezone-convert timestamps |
| **Prompt Enhancer** | 6 | Optimize, rewrite, score, A/B test AI prompts |
| **Market Intelligence** | 6 | AI market trends, GitHub stats, competitor analysis |
| **Fortune & Tarot** | 3 | Daily fortune, tarot readings, I Ching |
| **Content Publisher** | 8 | MoltBook posts, social content, newsletter |
| **AgentForge Compare** | 5 | Compare AI tools, frameworks, MCP servers |

All 9 servers share one Pro key. **$9/mo = 49 tools.**
