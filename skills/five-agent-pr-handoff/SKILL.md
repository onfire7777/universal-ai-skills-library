---
name: five-agent-pr-handoff
description: Use when a five-agent-dev-team automation opens, updates, labels, or hands off a pull request.
metadata:
  short-description: PR handoff packet for the five-agent dev team
---

# Five Agent PR Handoff

Every material PR update needs a compact handoff.

Include:
- criterion IDs or blocker IDs touched
- files changed
- lane owner and out-of-lane touches
- verification commands and pass/fail result
- risk and rollback notes
- Docker impact
- acceptance evidence table when behavior changed
- next owner stage

Use labels when available:
- `codex-status:blocked`
- `codex-status:handoff-ready`
- `codex-status:needs-review`
- `codex-status:ready-to-merge`
- relevant `codex-stage:*` and `codex-blocker:*`

Do not mark `ready-to-merge` while draft, red, unreviewed, missing evidence, or blocked.
