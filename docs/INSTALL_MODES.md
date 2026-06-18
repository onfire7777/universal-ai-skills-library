# Install Modes

The Universal AI Skills Router supports three conceptual install modes. Keep them distinct.

The concrete connection map for these install modes lives in
`UNIVERSAL_AI_CONNECTION_CONFIGS.md`. Use that file when verifying which
repo-owned config installs or updates a local AI client.

## 1. Wrapper install — default/recommended

A small `universal-ai-skills` skill is installed into an agent root. The wrapper tells the agent to call:

```bash
skill-router skill search <query>
skill-router skill <name>
```

Use this for most agents. It keeps always-loaded context small and avoids copying 1,812 canonical skills into every client.

## 2. Selected-skill install — explicit

Only a chosen skill or small set of skills is installed physically into a target agent root.

Use this when an agent cannot call the router directly or when a workflow needs a self-contained skill bundle.

## 3. Full-copy install — dangerous / explicit

All skills are copied or linked into a target root.

Use only when a client requires physical skills and the user explicitly accepts the clutter/cost. Full-copy mode can flood agent roots and should never be the accidental default.

## Current CLI behavior

- `skill-router sync matrix` is read-only and safe.
- `skill-router sync installed` updates every detected local skill-root adapter with the compact wrapper only.
- `skill-router sync paperclip` updates the Paperclip wrapper root and compact Paperclip agent instructions only.
- `skill-router sync propagate`, `skill-router sync all`, `skill-router update`, `skill-router skills install`, and `skill-router skills propagate` are wrapper-only by default.
- Full-copy installation requires the explicit `--full-copy` flag.
- Default writes are limited to conservative roots returned by `platform.AgentRoots()`.
- Workspace/source trees and hosted/repo-instruction adapters remain report-only unless a client-specific adapter is added.

## Explicit command shape

```bash
skill-router sync matrix
skill-router sync installed
skill-router sync paperclip
skill-router sync propagate
skill-router sync propagate --full-copy
skill-router skills install --target ~/.codex/skills
skill-router skills install --target ~/.codex/skills --full-copy
```

Do not use `--full-copy` unless the user explicitly accepts redundant physical skill copies.
