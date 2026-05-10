---
name: 'step-05-assertions'
description: 'Review assertion quality'
nextStepFile: './step-06-coverage.md'
referenceFiles:
  - 'references/common/assertions.md'
  - 'references/common/rules.md'
---

# Step 5: Review Assertion Quality

## STEP GOAL

Check assertion quality — verify tests use specific value assertions, avoid existence-only checks, and properly verify mock interactions.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/assertions.md` — Assertion patterns and matchers
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### 1. Check for Assertion Anti-Patterns

| Anti-Pattern | Found | Location |
|--------------|-------|----------|
| `toBeDefined()` only | [ ] | |
| `toBeTruthy()` only | [ ] | |
| Missing property assertions | [ ] | |
| Conditional assertions (if/else) | [ ] | |
| No mock call verification | [ ] | |

### 2. Verify Assertions Are Specific

Assertions should verify exact values, not just existence:

```typescript
// BAD
expect(result).toBeDefined();

// GOOD
expect(result).toEqual({ id: 'user-123', email: 'test@example.com' });
```

For each test, check:
- Does it verify the **actual values** returned, not just that something was returned?
- Does it check **all relevant properties** of the result?
- Are assertions **deterministic** (no random data, no Date.now())?

### 3. Verify Mock Interactions

Every test that arranges mocks should verify those mocks were called correctly:

```typescript
expect(mockService.method).toHaveBeenCalledWith(expectedArgs);
expect(mockService.method).toHaveBeenCalledTimes(1);
```

Check:
- Are mock calls verified with `toHaveBeenCalledWith` (specific arguments)?
- Are mock call counts verified with `toHaveBeenCalledTimes`?
- Are **all** arranged mocks actually asserted against?

### 4. Document Assertion Issues

```
**Assertion Issues:**
- [CRITICAL]: [test name] - Only checks existence, not values
- [MAJOR]: [test name] - Missing mock call verification
- [MINOR]: [test name] - Could use more specific matcher
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 5: Assertion Quality
==========================

Anti-Pattern Scan:
  toBeDefined() only:            [FOUND N / NONE]
  toBeTruthy() only:             [FOUND N / NONE]
  Missing property assertions:   [FOUND N / NONE]
  Conditional assertions:        [FOUND N / NONE]
  No mock call verification:     [FOUND N / NONE]

Test-by-Test Assertion Review:
  [PASS/FAIL] "test name"
    - Value specificity: [OK / WEAK / MISSING]
    - Mock verification: [OK / INCOMPLETE / MISSING]

Summary: N assertions reviewed, N issues found
  Critical: N  |  Major: N  |  Minor: N
```

Then ask: **[C] Continue to Step 6: Test Coverage**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-coverage.md`.
