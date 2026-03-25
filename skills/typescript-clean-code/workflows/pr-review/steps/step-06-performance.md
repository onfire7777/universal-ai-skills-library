---
name: 'step-06-performance'
description: 'Check for performance issues — N+1, memory, blocking'
nextStepFile: './step-07-run-code.md'
---

# Step 6: Check for Performance Issues

## STEP GOAL

Identify obvious performance problems in the PR changes.

## ANALYSIS PROCESS

Check all changed code for:

1. **N+1 query problems**?
2. **Large data sets in memory**?
3. **Unnecessary database calls**?
4. **Missing indexes for queries**?
5. **Blocking operations in async code**?

### Common Issues

```typescript
// BAD: N+1 query
for (const user of users) {
  const orders = await db.getOrdersForUser(user.id);
}

// GOOD: Single query
const orders = await db.getOrdersForUsers(userIds);
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 6: Performance
===================

[PASS/ISSUE] issue_type (file:line)
  - Impact: HIGH/MEDIUM/LOW
  - Description: what the performance issue is
  - Suggestion: how to optimize

Summary: N performance issues found
```

Then ask: **[C] Continue to Step 7: Run Code**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-07-run-code.md`.
