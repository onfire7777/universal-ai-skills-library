# Skills Corpus

The corpus is the actual library — every skill lives once under [`skills/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/skills) as a directory containing a `SKILL.md`. This page covers a skill's anatomy, the frontmatter schema, normalization, and dropped skills.

Source docs: [`docs/NORMALIZATION_REPORT.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/NORMALIZATION_REPORT.md) · [`docs/DROPPED.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/DROPPED.md) · [`schemas/skill.schema.json`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/schemas/skill.schema.json)

---

## Size & shape

| | |
|---|---|
| **Total skills** | ~1,812 (18 core + 1,794 library) per `manifest.json` `v2.2.8` |
| **Storage** | `skills/<kebab-name>/` — one directory per skill |
| **Canonical ID** | the kebab-case directory name (see [policy](Manifest-and-Registry.md#canonical-id-policy)) |
| **Aliases** | ~214 legacy/display names mapped in `compatibility_aliases.json` |

**Core vs library:** 18 *core* skills are always available (e.g. `skill-creator`, `prompt-engineer`, `model-selector`, `universal-ai-skills`, `multi-model-oracle`). The remaining ~1,794 *library* skills are indexed and loaded on demand.

---

## Anatomy of a skill

Minimum — a single file:

```
my-skill/
└── SKILL.md          # required: YAML frontmatter + markdown body
```

A richer skill can carry supporting files:

```
acceptance-testing/
├── SKILL.md          # required
├── README.md         # optional extended docs
├── references/       # optional linked docs / citations
├── scripts/          # optional utility scripts (.py, .js, …)
├── resources/        # optional templates, snippets
├── assets/           # optional media
├── examples/         # optional working examples
├── rules/            # optional domain rules
├── LICENSE           # optional
└── .openskills.json  # optional OpenSkills compatibility metadata
```

From a 100-skill sample, `SKILL.md` is universal; `.openskills.json` (~54%), `references/` (~20%), and `scripts/` (~13%) are the most common extras.

---

## `SKILL.md` frontmatter schema

Validated against [`schemas/skill.schema.json`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/schemas/skill.schema.json). Two fields are **required**; the rest are optional and many are Phase-0 routing-intelligence fields (currently validated in warn mode).

```yaml
---
name: acceptance-testing          # REQUIRED — kebab-case, MUST equal directory name
description: Write executable...   # REQUIRED — single-line, whitespace-collapsed
license: MIT
version: 1.2.3
aliases: [bdd-testing, cucumber-tests]
domain: engineering
tags: [python, testing, bdd]
triggers: ["write BDD tests", "cucumber scenarios"]
capabilities: [test.bdd, code.cucumber]
requires: { language: python, framework: pytest }
quality_tier: stable
supersedes: [old-bdd-skill]
conflicts: [competing-bdd-tool]
metadata:
  displayName: Acceptance Testing
  author: ...
  compatibility: [claude, copilot, cursor]
  references: [{ title: ..., url: ... }]
---
```

| Field | Req | Purpose |
|---|---|---|
| `name` | ✅ | Canonical kebab-case ID; must match the directory name. |
| `description` | ✅ | Primary single-line description used in routing. |
| `license` | | License identifier (MIT, Apache-2.0, …). |
| `version` | | Per-skill semantic version (Phase 0). |
| `aliases` | | Alternate exact-match lookup names. |
| `domain` | | Single taxonomy bucket (engineering, design, data, …) (Phase 0). |
| `tags` | | Free-form facets — language, framework, capability (Phase 0). |
| `triggers` | | Example prompts that should route here (Phase 0). |
| `capabilities` | | Typed dotted verbs, e.g. `test.bdd` (Phase 0). |
| `requires` | | Declared dependencies (806+ corpus skills already use this). |
| `quality_tier` | | Lifecycle/trust tier; a reranker feature (Phase 0). |
| `supersedes` / `conflicts` | | Lifecycle edges between skills (Phase 0). |
| `metadata` | | Display name, author, compatibility, references. |

---

## Normalization

Incoming skills from many sources are normalized into one consistent corpus. The process:

1. **Deduplicate** — drop exact duplicates, keep the canonical (cleanest-ID) version.
2. **Standardize IDs** — convert to kebab-case top-level directory names.
3. **Repair YAML** — fix malformed frontmatter (unquoted colons, broken structure).
4. **Detect missing metadata** — flag skills lacking a valid `SKILL.md`.
5. **Map aliases** — preserve old/display names as lookup aliases.

A documented normalization run took **849 input folders → 773 skills** (76 dropped); the corpus has since grown to ~1,812 through further integration. `manifest.json` is the authority for the current count.

---

## Dropped skills

The historical ledger in [`docs/DROPPED.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/DROPPED.md) records **76 drops** from that run, retained for auditability:

| Reason | Count | Examples |
|---|---|---|
| Duplicate of a canonical skill | 73 | `commandcode-brainstorming` → kept `brainstorming`; `semgrep-official` → kept `tob-semgrep` |
| Missing `SKILL.md` | 3 | `excalidraw-diagram`, `mlevison-systems-thinking`, `socratic-architect` |
| YAML repaired (recovered, not dropped) | 1 | `davila-motion-canvas` |

---

## How the corpus connects to routing

The corpus itself knows nothing about the router. The [registry](Manifest-and-Registry.md) (`manifest.json`) indexes it, and the [router](Skill-Router-CLI.md) scores against that index — never by scanning 1,800 files at query time. That separation is the heart of the [architecture](Architecture.md).

---

## Related pages

- [Manifest & Registry](Manifest-and-Registry.md) — the generated index and ID policy
- [Sources & Integrations](Sources-and-Integrations.md) — where skills come from
- [Contributing](Contributing.md) — how to add a skill
