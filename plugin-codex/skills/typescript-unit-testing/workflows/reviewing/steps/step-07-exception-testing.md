---
name: 'step-07-exception-testing'
description: 'Review exception testing completeness'
nextStepFile: './step-08-compile-report.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/assertions.md'
---

# Step 7: Review Exception Testing

## STEP GOAL

Verify every exception thrown in the source code has a corresponding test that checks the exception type, message, and error code.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements
- `references/common/assertions.md` — Assertion patterns and matchers

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### 1. Inventory All Exceptions in Source Code

Scan the source file for every `throw` statement and map each to its condition:

| Exception | Condition | Test Exists | Verifies Type | Verifies Message | Verifies Code |
|-----------|-----------|-------------|---------------|------------------|---------------|
| NotFoundException | not found | [ ] | [ ] | [ ] | [ ] |
| BadRequestException | validation | [ ] | [ ] | [ ] | [ ] |

Include:
- Explicit `throw new XxxException(...)` statements
- Exceptions thrown by called methods that should propagate
- Guard clauses that throw

### 2. Check Exception Test Patterns

Each exception test should verify multiple aspects:

```typescript
// Should verify type
await expect(target.method()).rejects.toThrow(NotFoundException);

// Should verify message
await expect(target.method()).rejects.toThrow('Expected message');

// Should verify error code if applicable
await expect(target.method()).rejects.toMatchObject({
  errorCode: 'EXPECTED_CODE',
});
```

For each exception test, check:
- Does it verify the **exception type** (class)?
- Does it verify the **error message** (string content)?
- Does it verify the **error code** if the exception has one?
- Does it properly set up the condition that triggers the exception?
- Does it use `rejects.toThrow` for async methods?

### 3. Document Exception Testing Issues

```
**Exception Testing Issues:**
- [exception type]: [missing verification]
- [exception type]: Test exists but does not verify message
- [exception type]: No test found for this exception
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 7: Exception Testing
==========================

Exception Inventory:
  [PASS/FAIL] ExceptionType — "condition description"
    - Type verified:    [YES/NO]
    - Message verified: [YES/NO]
    - Code verified:    [YES/NO/N/A]
    - Test location:    [test name or MISSING]

Summary:
  Exceptions in source: N
  Fully tested: N
  Partially tested: N
  Untested: N

Issues:
  - [SEVERITY]: [description]
```

Then ask: **[C] Continue to Step 8: Compile Report**

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-08-compile-report.md`.
