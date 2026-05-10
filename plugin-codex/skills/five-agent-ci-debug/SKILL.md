---
name: five-agent-ci-debug
description: Use when five-agent-dev-team GitHub Actions, Gitleaks, Docker, audit, lint, test, or formatting checks fail.
metadata:
  short-description: CI failure triage for the five-agent dev team
---

# Five Agent CI Debug

Triage in this order.

1. Identify the current failing job and first failing step.
2. Compare with the previous run to avoid fixing stale failures.
3. Extract only relevant error lines; do not paste raw logs containing secrets.
4. Classify owner lane: CI/security, quality/test, frontend, backend, docs, or dependency.
5. Reproduce locally if feasible with the PR branch's own tooling and config.
6. Apply the smallest fix in the owning lane.
7. Run the same local command that failed remotely.
8. Push only after local verification passes.
9. Update labels and rolling state.

For Docker compose validation in reportable contexts, use `docker compose config --no-interpolate`.
