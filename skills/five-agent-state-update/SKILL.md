---
name: five-agent-state-update
description: Use when updating five-agent-dev-team swarm state, ledger entries, blockers, or next-stage handoffs.
metadata:
  short-description: Compact state and ledger updates
---

# Five Agent State Update

Update two files only:
- `~/.codex/state/five-agent-dev-team-swarm.md`
- `~/.codex/state/five-agent-dev-team-ledger.ndjson`

Rolling state must stay under 200 lines and keep only:
- current cycle
- active PR queue
- do-not-duplicate scope
- blockers
- last stage summaries
- next handoff

Ledger entries are one NDJSON object per run with timestamp, stage, cycleId, baseSha, headSha, PR, files touched, checks, blockers, and next handoff.

Do not paste full logs or secret-like values.
