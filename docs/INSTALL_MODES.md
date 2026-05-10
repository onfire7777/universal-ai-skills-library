# Install Modes

The Universal AI Skills Router supports three conceptual install modes. Keep them distinct.

## 1. Wrapper install — default/recommended

A small `universal-ai-skills` skill is installed into an agent root. The wrapper tells the agent to call:

```bash
skill-router skill search <query>
skill-router skill <name>
```

Use this for most agents. It keeps always-loaded context small and avoids copying 1,805 canonical skills into every client.

## 2. Selected-skill install — explicit

Only a chosen skill or small set of skills is installed physically into a target agent root.

Use this when an agent cannot call the router directly or when a workflow needs a self-contained skill bundle.

## 3. Full-copy install — dangerous / explicit

All skills are copied or linked into a target root.

Use only when a client requires physical skills and the user explicitly accepts the clutter/cost. Full-copy mode can flood agent roots and should never be the accidental default.

## Current CLI behavior

- `skill-router sync matrix` is read-only and safe.
- Existing `sync propagate` behavior remains conservative: it writes only to legacy default roots returned by `platform.AgentRoots()`.
- Newly detected roots such as Windsurf, Roo, Continue, Qwen, and Kimi/OpenClaw are report-only until adapter semantics are confirmed.

## Desired future commands

```bash
skill-router sync plan
skill-router sync install-wrapper --all
skill-router sync install-wrapper --agent claude
skill-router sync install-selected <skill> --agent codex
skill-router sync propagate-full --all --confirm-full-copy
```

These commands should make install intent explicit before changing any agent root.
