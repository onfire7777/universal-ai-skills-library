---
name: five-agent-preflight-check
description: Use before any five-agent-dev-team automation mutates files, opens PRs, updates branches, or performs Captain integration.
metadata:
  short-description: Preflight checks for the five-agent dev team repo
---

# Five Agent Preflight Check

Run this before mutation.

1. Read `~/.codex/specs/five-agent-dev-team-automation-team.md` and the rolling swarm state.
2. Check repo status: `git status --porcelain`.
3. Fetch safely: `git fetch --prune origin`.
4. Check sync against the intended base with `git rev-list --left-right --count`.
5. Inspect active PRs, labels, draft state, latest checks, and changed-file overlap.
6. Stop and write a handoff if the repo is dirty, behind, diverged, locked, duplicating PR #7, or unable to validate locally.
7. Never paste secrets, resolved environment values, or raw Docker compose output.
