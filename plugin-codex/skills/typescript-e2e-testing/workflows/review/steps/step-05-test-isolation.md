---
name: 'step-05-test-isolation'
description: 'Verify tests do not share state and are truly independent'
nextStepFile: './step-06-tech-patterns.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/best-practices.md'
---

# Step 5: Review Test Isolation

## STEP GOAL

Verify tests don't share state and are truly independent.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — isolation rules
- `references/common/best-practices.md` — isolation best practices

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### Isolation Checks
- [ ] Unique identifiers per test (not hardcoded)
- [ ] Database cleaned before each test
- [ ] No global variables modified
- [ ] Message buffers cleared (Kafka)
- [ ] Mocks reset between tests

### Isolation Issues to Flag

```typescript
// ❌ Hardcoded ID (causes conflicts)
const userId = 'user-123';
// ✅ Unique ID
const userId = `user-${Date.now()}-${Math.random().toString(36).slice(2)}`;

// ❌ Missing cleanup wait
beforeEach(async () => {
  await repository.clear();
});
// ✅ With wait for in-flight
beforeEach(async () => {
  await new Promise(r => setTimeout(r, 500));
  await repository.clear();
});

// ❌ Kafka: Using fromBeginning
await consumer.subscribe({ topic, fromBeginning: true });
// ✅ Use pre-subscription + buffer clearing
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 5: Isolation Review
========================

[File: {filename}]

Issues Found:
- Hardcoded IDs: {lines}
- Missing cleanup waits: {lines}
- Shared state risk: {details}

Impact: Tests may fail intermittently when run together

Fixes Required:
{specific changes}

Summary: {n} isolation issues found
```

Then ask: **[C] Continue to Step 6: Technology Patterns / [F] Fix issues first**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-tech-patterns.md`.
