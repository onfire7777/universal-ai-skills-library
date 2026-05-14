# Universal Skills Corpus Report

**Date:** 2026-05-09
**Pipeline:** universal router architecture pass
**Source of truth:** `skills/`
**Current corpus:** 1,811 canonical skills, all with `SKILL.md`
**Access model:** on-demand loading through `skill-router skill <name>` with `manus skill <name>` kept only as a legacy alias.

## Ledger

| Bucket | Count |
|---|---:|
| Current skill directories | 1,811 |
| Directories with `SKILL.md` | 1,811 |
| Missing `SKILL.md` | 0 |
| Legacy/display aliases preserved | 125 |
| Underscore duplicate directories merged | 26 |
| Current source-of-truth location | `skills/` |
| Router source | `skill-router-cli/` |
| Universal setup skill | `universal-ai-config` |

## What was done

1. **Universal naming pass** - default surfaces now use `Universal AI Skills`, `skill-router`, and `universal-ai-config`.
2. **Compatibility boundary pass** - product names such as Codex, Claude Code, OpenAI, Anthropic, Manus, and Gemini remain only where they identify a real provider, API, model, client, or compatibility adapter.
3. **Context efficiency pass** - global instructions and plugin manifests are index-only; full skill bodies are loaded on demand.
4. **Router pass** - `skill-router-cli/` is the canonical CLI source; `manus` remains a compatibility executable and lookup alias.
5. **Full-corpus import pass** - archived local skill roots were merged into `skills/` so the universal repo, not individual AI roots, is the complete source of truth.
6. **Alias pass** - legacy display names and old source names were retained as manifest aliases when they do not collide with canonical skill ids.

## Verification

A current corpus check reports all skill directories have a `SKILL.md`:

```
skill directories: 1,811
directories with SKILL.md: 1,811
missing SKILL.md: 0
```

## Format

Every skill uses the portable `SKILL.md` package shape:

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

## Reports and Manifests

- `manifest.json` - current skill index consumed by the router/plugin.
- `docs/DROPPED.md` - historical source normalization ledger, retained for provenance.
- `docs/ARCHITECTURE_V2.md` - current universal architecture.

## Out of scope

This report verifies corpus presence, naming boundaries, and router architecture. Individual third-party provider APIs, model pricing, and optional runtime adapters still require live checks when a task depends on current external behavior.
