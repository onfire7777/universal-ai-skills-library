---
name: five-agent-security-review
description: Use for five-agent-dev-team secret handling, workflow permissions, Docker safety, dependency audit, local binding, and supply-chain review.
metadata:
  short-description: Security review for the five-agent dev team
---

# Five Agent Security Review

Check:
- least-privilege workflow permissions
- no secrets in code, config, logs, PR bodies, state, or docs
- no broad OAuth/GitHub scopes without a spec reason
- Docker compose output is not resolved into state or logs
- services bind to `127.0.0.1` unless explicitly required
- `npm audit fix --force` is never used
- moderate advisories are handled by documented policy, not hidden

Evidence should be compact: command names, pass/fail, and redacted risk notes.
