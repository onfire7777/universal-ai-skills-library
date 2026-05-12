---
name: gbrain
description: Use when the user asks for GBrain, personal knowledge brain setup, brain-first retrieval, local PGLite brain, hybrid RAG search, gbrain import/query/sync/embed, GBrain skills, Minions, durable agent jobs, soul-audit, brain maintenance, or integrating GBrain with gstack, Hermes, OpenClaw, Claude Code, or Codex.
---

# GBrain Universal Adapter

Canonical upstream checkout: `C:\Users\burni\gbrain`

GBrain is a real CLI-backed personal knowledge brain. It is not just a skill
pack. The universal stack keeps the upstream checkout and local brain state in
one place, then exposes GBrain skills and commands through `skill-router`.

## Core Commands

```bash
gbrain --version
gbrain doctor --json
gbrain stats
gbrain import <dir> --no-embed
gbrain embed --stale
gbrain search "<query>"
gbrain query "<question>"
gbrain sync --repo <brain-repo>
gbrain jobs smoke
```

Use keyword search without API keys. Vector embeddings require
`OPENAI_API_KEY`. Anthropic query expansion and durable subagent workers require
the relevant provider key only for those optional features.

## Load GBrain Skills On Demand

The upstream GBrain skills are indexed from:

```text
C:\Users\burni\gbrain\skills
```

Examples:

```bash
skill-router skill brain-ops
skill-router skill signal-detector
skill-router skill query
skill-router skill ingest
skill-router skill media-ingest
skill-router skill maintain
skill-router skill minion-orchestrator
skill-router skill soul-audit
skill-router skill skillify
```

Before changing a GBrain workflow, read:

```bash
skill-router skill gbrain
skill-router skill brain-ops
skill-router skill query
```

For setup or repair, read the upstream installer flow first:

```bash
Get-Content C:\Users\burni\gbrain\INSTALL_FOR_AGENTS.md
Get-Content C:\Users\burni\gbrain\AGENTS.md
Get-Content C:\Users\burni\gbrain\docs\GBRAIN_VERIFY.md
```

## Integration Policy

- Keep `C:\Users\burni\gbrain` as the single upstream source checkout.
- Keep `C:\Users\burni\.gbrain` for GBrain runtime state such as PGLite data and
  user-level GBrain skills.
- Do not vendor GBrain's full skill tree into every AI root.
- Use `skill-router` for on-demand loading and `gbrain` for persistent brain
  operations.
- GStack and GBrain can cooperate: GStack supplies engineering process skills;
  GBrain supplies brain-first retrieval, memory, jobs, and durable agent context.

## Maintenance

Follow the upstream upgrade path:

```bash
git -C C:\Users\burni\gbrain pull --ff-only
cd C:\Users\burni\gbrain
bun install
gbrain init
gbrain post-upgrade
gbrain doctor --json
```

Do not run `bun install -g github:garrytan/gbrain`; upstream documents that this
skips required install behavior.
