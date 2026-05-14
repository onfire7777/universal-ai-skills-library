# Design And Messaging

Use this page to keep the repository's public presentation consistent across
GitHub About, README, plugin manifests, marketplaces, and generated previews.

## Positioning

Universal AI Skills Library is not a bundle to paste into every agent. It is a
router-first skill system:

1. Keep one canonical 1,811-skill corpus in `skills/`.
2. Give each AI client compact instructions or a wrapper skill.
3. Run `skill-router preflight` for real user prompts.
4. Load exactly one matching skill only when the route clearly fits.
5. Keep local model routing, memory, embeddings, and tool adapters centralized.

## Recommended GitHub About

```text
Router-first AI skill system for Codex, Claude, Cursor, Hermes, Paperclip, OpenCode, and local AI stacks: search, preflight-route, and load 1,811 skills on demand without duplicating the corpus.
```

## Short Taglines

- One shared skill corpus. One router. Every AI agent.
- Load the skill you need, when you need it.
- Universal skill access without universal context bloat.
- Router-first skills for local and CLI AI agents.

## Visual Direction

The public front page should feel:

- modern
- minimal
- technical
- calm
- trustworthy
- public-safe

Use a dark neutral base, restrained cyan/blue/green accents, rounded geometry,
and network-like routing diagrams. Avoid cluttered screenshots, fake terminal
spam, token/key-looking text, and dense skill tables in the README hero.

## README Hero Asset

Primary asset:

```text
docs/assets/universal-ai-skills-hero.png
```

The PNG is the reviewed public README hero. Keep it public-safe, readable at
README width, and aligned with the actual architecture:

- center: `skill-router`
- agent adapters: Codex, Claude, Cursor, Hermes, Paperclip, OpenCode, Kimi,
  Gemini, and Qwen
- corpus: 1,811 canonical skills loaded on demand from one shared source
- source layer: Lightpanda, Context Mode, MemPalace, Web Search
- memory and index layer: GBrain and GSkills

## GPT Image 2 / Raster Preview Prompt

Use this prompt when refreshing the raster hero or creating a matching social
preview. Do not commit a generated image until it has been reviewed for text
accuracy, licensing, visual artifacts, and secret-like strings.

```text
Create a modern, sleek, minimal GitHub repository social-preview image for "Universal AI Skills Library".

Concept: one central AI skill router connects many AI agents to one shared skill corpus. Show a clean dark interface-inspired composition with a central node labeled "skill-router", surrounding nodes labeled Codex, Claude, Cursor, Hermes, Paperclip, and OpenCode, and a subtle line "1,811 skills loaded on demand".

Visual style: premium developer tooling, dark neutral background, refined cyan/blue/green accents, precise spacing, soft depth, no busy code screenshots, no fake secrets, no API keys, no mascot, no cartoon style.

Composition: 16:9, strong center focus, ample margins, readable type, abstract routing lines, polished SaaS/devtool aesthetic, suitable as a GitHub social preview and README hero.
```

## Metadata Consistency Checklist

Keep these surfaces aligned when the product description changes:

- GitHub About description
- `README.md`
- `marketplace.json`
- `plugin/plugin.json`
- `plugin/marketplace.json`
- `plugin/.codex-plugin/plugin.json`
- `plugin-codex/.codex-plugin/plugin.json`
- `manifest.json`
- `docs/PUBLIC_RELEASE_CHECKLIST.md`
