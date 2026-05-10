---
name: 'step-05-security'
description: 'Check for security vulnerabilities — injection, XSS, auth, secrets'
nextStepFile: './step-06-performance.md'
---

# Step 5: Check for Security Issues

## STEP GOAL

Identify potential security vulnerabilities in the PR changes.

## ANALYSIS PROCESS

Check all changed code for:

1. **Input validation** present?
2. **SQL injection** possible?
3. **XSS vulnerabilities**?
4. **Sensitive data** exposed?
5. **Authentication/authorization** checked?
6. **Secrets in code**?

### Common Issues

```typescript
// BAD: SQL injection
const query = `SELECT * FROM users WHERE id = ${userId}`;

// GOOD: Parameterized query
const query = 'SELECT * FROM users WHERE id = ?';
db.query(query, [userId]);
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 5: Security
================

[PASS/ISSUE] vulnerability_type (file:line)
  - Risk: HIGH/MEDIUM/LOW
  - Description: what the vulnerability is
  - Attack vector: how it could be exploited
  - Fix: how to remediate

Summary: N security issues found
  - HIGH: N
  - MEDIUM: N
  - LOW: N
```

Then ask: **[C] Continue to Step 6: Performance**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-performance.md`.
