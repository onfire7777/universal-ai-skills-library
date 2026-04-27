# Manus Skills Library

Canonical flat store of **770 production-quality skills** in the Anthropic / OpenAI Agent Skills standard format.

## Layout

```
skills/
├── aes-encryption/SKILL.md
├── motion-canvas/SKILL.md
├── thinking-bayesian/SKILL.md
└── ... (770 skills total)
```

Every skill folder contains a `SKILL.md` with YAML frontmatter (`name`, `description`, optional `license` and `metadata`) and markdown body. Some include `scripts/`, `references/`, `assets/` subfolders.

## Direct flat install

Drop into any agent's flat skills folder:

```bash
# Claude Code / Cowork
cp -r skills/* ~/.claude/skills/

# OpenAI Codex CLI
cp -r skills/* ~/.codex/skills/
```

Both auto-detect on next session start. The format is identical across both ecosystems (Anthropic donated Skills as an open spec on 2025-12-18; OpenAI adopted it for Codex and ChatGPT in December 2025).

## For plugin-style install

See [`manus-skills-marketplace`](https://github.com/onfire7777/manus-skills-marketplace).

## For browse / discovery

See [`manus-skills-organized`](https://github.com/onfire7777/manus-skills-organized).

## Categories

| # | Browse folder | Plugin name | Skills | Purpose |
|--:|---|---|--:|---|
| 01 | `01-security` | `manus-security` | 134 | Security: pentesting, exploits, hardening, crypto, AppSec/NetSec |
| 02 | `02-privacy-compliance` | `manus-privacy-compliance` | 55 | Privacy and compliance: GDPR/CCPA/HIPAA, DPIAs, AI Act, ISO 27701 |
| 03 | `03-frontend-design` | `manus-frontend-design` | 113 | Frontend and design: UI/UX, web frontend, design systems, accessibility |
| 04 | `04-backend-architecture` | `manus-backend-architecture` | 55 | Backend architecture: APIs, system design, distributed systems, databases |
| 05 | `05-devops-infra` | `manus-devops-infra` | 21 | DevOps and infrastructure: CI/CD, containers, cloud, IaC, SRE |
| 06 | `06-ai-prompt-engineering` | `manus-ai-prompt-engineering` | 32 | AI and prompt engineering: LLM apps, agents, RAG, evals |
| 07 | `07-data-analysis` | `manus-data-analysis` | 16 | Data analysis: analytics, viz, statistics, ML/data pipelines |
| 08 | `08-reasoning-thinking` | `manus-reasoning-thinking` | 72 | Reasoning and thinking: mental models, critical thinking, epistemics |
| 09 | `09-creativity-innovation` | `manus-creativity-innovation` | 23 | Creativity and innovation: creative writing, art, ideation |
| 10 | `10-planning-strategy` | `manus-planning-strategy` | 68 | Planning and strategy: PM, strategic planning, roadmaps, OKRs |
| 11 | `11-communication` | `manus-communication` | 25 | Communication: writing, presentations, facilitation, leadership |
| 12 | `12-quality-engineering` | `manus-quality-engineering` | 46 | Quality engineering: code review, testing, debugging, refactoring |
| 13 | `13-meta-skills` | `manus-meta-skills` | 109 | Meta-skills: skill creation, learning, productivity |
| 14 | `14-domain-specific` | `manus-domain-specific` | 1 | Narrow verticals: medical, legal, finance, scientific domains |

## Companion repos

These three repos hold the **same 770 skills** in different layouts. Use whichever fits your workflow:

| Repo | Role | Layout |
|---|---|---|
| [`manus-skills-library`](https://github.com/onfire7777/manus-skills-library) | **Canonical store** | Flat `skills/<slug>/SKILL.md` |
| [`manus-skills-organized`](https://github.com/onfire7777/manus-skills-organized) | **Browse view** | Numbered category folders `skills/<NN-category>/<slug>/` |
| [`manus-skills-marketplace`](https://github.com/onfire7777/manus-skills-marketplace) | **Install distribution** | Plugin format for Claude Code (`plugins/`) and OpenAI Codex CLI (`ai-codex/plugins/`) |

## Naming convention (consistent across all three repos)

Same 14 category slugs everywhere. Prefix differs only by purpose:

- **Skill slug:** `aes-encryption`, `motion-canvas`, `thinking-bayesian` (kebab-case, no prefix)
- **Browse folder:** `01-security`, `02-privacy-compliance`, ... (numeric prefix sorts naturally for navigation)
- **Plugin identifier:** `manus-security`, `manus-privacy-compliance`, ... (`manus-` namespace prefix for plugin ecosystems)


## Format

```yaml
---
name: kebab-case-slug
description: One-line description of what the skill does and when it triggers.
license: ...
metadata:
  domain: ...
  tags: [...]
---

# Title

## Overview
...
```

## Provenance

Built by deduplicating, normalizing, and re-categorizing skills from the original Manus skills ecosystem. 849 input → 773 cleaned (3 unsalvageable, 73 duplicates collapsed) → 770 final after additional plugin-canonical dedup. Zero defects on the final output (verified): no broken YAML, no missing required fields, no name collisions, no duplicate bodies.
