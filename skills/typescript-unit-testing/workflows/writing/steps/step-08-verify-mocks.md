---
name: 'step-08-verify-mocks'
description: 'Verify all mock interactions and finalize test suite'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/mocking/deep-mocked.md'
---

# Step 8: Verify Mock Interactions

## STEP GOAL

Verify all mock interactions are properly asserted, run the full test suite, and check coverage meets the 80%+ threshold.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements
- `references/mocking/deep-mocked.md` — DeepMocked patterns and verification

## EXECUTION

### 1. Verify Mock Interactions Across All Tests

Review every test written in steps 4-7 and ensure mock interactions are verified:

#### Verify Exact Parameters

```typescript
expect(mockService.method).toHaveBeenCalledWith(
  expectedArg1,
  expectedArg2
);
```

#### Verify Call Count

```typescript
expect(mockService.method).toHaveBeenCalledTimes(1);
```

#### Verify NOT Called When Appropriate

```typescript
expect(mockService.method).not.toHaveBeenCalled();
```

#### Verify Call Order When Relevant

```typescript
expect(mockService.methodA).toHaveBeenCalledBefore(mockService.methodB);
```

### 2. Run All Tests

Execute the test suite:

```bash
npm test -- [path/to/component.spec.ts]
```

Fix any failing tests before proceeding.

### 3. Check Coverage

Run coverage analysis:

```bash
npm run test:cov -- [path/to/component.spec.ts]
```

Verify 80%+ coverage on the component. If coverage is below 80%, identify untested paths and add tests.

### 4. Present Completion Checklist

Present the final checklist to the user:

```
Unit Tests Written for [ComponentName]:

Setup:
- [ ] Using `target` for SUT
- [ ] Using `mock` prefix for all mocks
- [ ] Using `DeepMocked<T>` type
- [ ] Including `.setLogger(new MockLoggerService())`
- [ ] Following AAA pattern with comments
- [ ] Resetting mocks in `afterEach`

Coverage:
- [ ] Happy path tests for all public methods
- [ ] Edge case tests (empty, null, boundaries)
- [ ] Error case tests (not found, validation)
- [ ] Exception type and error code verification
- [ ] Business rule tests (domain logic)
- [ ] Mock call verification (parameters + count)

Quality:
- [ ] 80%+ coverage achieved
- [ ] Specific assertions (not just existence)
- [ ] No conditional assertions
- [ ] Tests fail when any field differs
```

### 5. Anti-Patterns Final Check

Verify none of these anti-patterns exist in the test suite:

| Don't | Why | Do Instead |
|-------|-----|------------|
| `expect(result).toBeDefined()` | Doesn't catch wrong values | Assert specific values |
| `if (result) expect(...)` | Non-deterministic | Separate test cases |
| Test private methods | Couples to implementation | Test via public interface |
| Share state between tests | Causes flaky tests | Fresh setup in beforeEach |
| Skip mock verification | Doesn't validate behavior | Verify all mock calls |

### 6. Append to Report

Append to the output document:

```markdown
## Step 8: Mock Verification & Final Report

**Mock Verification:**
- All mock calls verified with exact parameters: [yes/no]
- All mock call counts verified: [yes/no]
- Negative cases (not.toHaveBeenCalled) verified: [yes/no]

**Test Execution:**
- All tests passing: [yes/no]
- Coverage: [percentage]%

**Anti-Pattern Check:**
- No existence-only assertions: [pass/fail]
- No conditional assertions: [pass/fail]
- No private method testing: [pass/fail]
- No shared state between tests: [pass/fail]
- All mocks verified: [pass/fail]

**Completion Checklist:** [All items checked / N items remaining]
```

## PRESENT FINDINGS

Show the user:
- Mock verification summary
- Test execution results
- Coverage percentage
- Completion checklist with all items checked or flagged
- Anti-pattern check results

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `8` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The unit test writing workflow is complete. The full report is saved at the output path. The spec file is ready at the `specFilePath` location.
