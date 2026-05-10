---
name: claude-native-code-review
description: "Adversarial code review — find what will break in production, not what looks wrong in theory. Covers correctness, security, failure modes, performance, and maintainability. Use when reviewing PRs, auditing code quality, or before merging. Triggers on: 'review this code', 'code review', 'check this PR', 'is this code good', 'review before merge', 'look over this code'."
---

# Code Review — Adversarial

## Thinking Protocol

Before reviewing, answer silently:
1. What does this code **actually do** vs what it **claims to do**?
2. What would an attacker target here?
3. What will break at 10x scale?

## Review Dimensions (priority order)

### 1. Correctness
- Trace happy path end-to-end, then every error path
- Identify implicit assumptions (nullability, ordering, concurrency)
- Check boundaries: empty inputs, max values, unicode, concurrent access

### 2. Security
- Injection surfaces (SQL, XSS, SSRF, command, path traversal)
- Auth/authz gaps, privilege escalation paths
- Secrets in code or logs

### 3. Failure Modes
- Network timeouts, partial failures, retry storms
- Resource leaks (connections, file handles, memory)
- Missing circuit breakers or fallbacks

### 4. Performance
- N+1 queries, unbounded loops, missing pagination
- Unnecessary allocations in hot paths
- Blocking operations in async contexts

### 5. Maintainability
- Naming that reveals intent vs implementation
- Coupling between components
- Test coverage of actual behavior

## Output Format

Per finding:
```
[CRITICAL|HIGH|MEDIUM|LOW] file:line
Problem: [one sentence]
Impact: [what actually breaks]
Fix: [specific code change]
```

## Rules
🚨 Recommend root cause fixes, not band-aids.
🚨 Don't list things that are fine. Only report actual findings.
🚨 If the code is solid, say so in one sentence and stop.