---
name: five-agent-captain-merge-gate
description: Use before the five-agent-dev-team Captain Integrator merges, rebases, labels ready-to-merge, pushes main, or advances main build state.
metadata:
  short-description: Captain merge gate for the five-agent dev team
---

# Five Agent Captain Merge Gate

Captain is integration-only.

Refuse merge unless all are true:
- Captain lock is acquired or no integration mutation is needed
- local `main` is clean and equals `origin/main`
- target PR is not draft
- target PR is up to date with the base
- latest CI on target head SHA is success
- review has run on the latest head SHA
- no unresolved P0/P1 findings
- P2 security/test/reliability findings are fixed or explicitly accepted
- no blocked, CI, security, or merge-conflict labels
- validation packet exists with commands, evidence, risk, rollback, and Docker impact

If any condition fails, no-op and write a precise handoff.
