# Skills Library Normalization Report

**Date:** 2026-04-25
**Pipeline:** `pipeline.py`
**Source:** `skills/` (849 folders)
**Output:** `skills/` (773 normalized skills)

## Ledger

| Bucket | Count |
|---|---:|
| Input folders | 849 |
| Missing `SKILL.md` (dropped) | 3 |
| Missing required field after repair (dropped) | 0 |
| Empty body (dropped) | 0 |
| Duplicate name/body (canonical kept) | 73 |
| **Final cleaned skills** | **773** |
| YAML repaired (recovered) | 1 |
| Description synthesized from body | 2 |

## What was done

1. **Audit pass** — every `SKILL.md` parsed; missing files, broken YAML, missing required fields, and empty bodies all enumerated.
2. **YAML repair** — common breakage fixed automatically:
   - Unquoted colons in scalar values (e.g. `description: Use when: ...`)
   - `@`-prefixed and version-operator items in flow sequences (e.g. `[@motion-canvas/core>=3.0.0]`)
3. **Description synthesis** — for skills missing the `description` field but with a non-empty body, the first non-heading paragraph (≤ 250 chars) was promoted to `description`.
4. **Dedup** — grouped by normalized name (kebab) and by body content hash. Canonical winner per group selected deterministically:
   1. kebab-case folder name beats spaces/TitleCase
   2. shorter folder name beats longer
   3. longer body beats shorter
   4. alphabetical tiebreak
5. **Normalize** — every output `SKILL.md` re-emitted via `yaml.safe_dump`, guaranteeing valid YAML. Required fields (`name` kebab, `description` single-line, `license`) promoted to top-level frontmatter; all other fields preserved under a `metadata:` block. Body preserved verbatim.

## Verification

A re-audit of the cleaned output reports **zero defects**:

```
Total defects: 0
  no_skill_md: 0
  no_frontmatter: 0
  unparseable: 0
  missing_name: 0
  missing_description: 0
  empty_body: 0
  name_not_kebab: 0
  duplicate names: 0
  duplicate bodies: 0
```

## Format

Every skill conforms to the **Anthropic Agent Skill Format** (also adopted by Manus, Cursor, etc.):

```yaml
---
name: kebab-case-name
description: One-line description.
license: ...
metadata:           # optional, preserves source-specific fields
  domain: ...
  tags: [...]
  author: ...
  version: ...
---

# Title

## Overview
...
```

## Reports

- `MANIFEST.json` — structured decision log (every dropped folder + reason)
- `DROPPED.md` — human-readable drop log

## Out of scope

This audit verified `SKILL.md` correctness only. Scripts inside `scripts/` directories were copied verbatim and were **not** runtime-tested.

## Symlinks

The source `ui-ux-pro-max/` skill contained dead symlinks pointing to `../../../src/ui-ux-pro-max/` (a path that doesn't exist in the repo). The `SKILL.md` itself was preserved; the dead links were skipped. No data lost.
