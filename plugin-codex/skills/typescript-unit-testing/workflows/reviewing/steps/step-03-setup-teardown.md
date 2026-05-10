---
name: 'step-03-setup-teardown'
description: 'Review setup and teardown patterns'
nextStepFile: './step-04-aaa-pattern.md'
referenceFiles:
  - 'references/common/rules.md'
---

# Step 3: Review Setup and Teardown

## STEP GOAL

Verify beforeEach/afterEach patterns follow conventions — fresh mock creation, TestingModule setup, logger configuration, and proper cleanup.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### 1. Verify beforeEach Setup

| Criterion | Status | Issue |
|-----------|--------|-------|
| Mocks created fresh | [ ] | |
| TestingModule built | [ ] | |
| `.setLogger()` included | [ ] | |
| Target retrieved from module | [ ] | |

Check that:
- All mocks are created fresh in each `beforeEach` (not shared across tests)
- `TestingModule` is built with proper providers
- `.setLogger(new MockLoggerService())` is called on the target
- The SUT (`target`) is retrieved from the module via `module.get()`

### 2. Verify afterEach Cleanup

| Criterion | Status | Issue |
|-----------|--------|-------|
| `jest.clearAllMocks()` called | [ ] | |
| No shared mutable state | [ ] | |

### 3. Check for Setup Anti-Patterns

Watch for and report:
- **Shared state between tests**: Variables mutated in one test affecting another
- **Missing mock reset**: Mocks not cleared between tests, leading to leaky state
- **Expensive setup in beforeEach**: Operations that should be in `beforeAll` instead (e.g., reading config files, static data)
- **Setup duplication**: Same arrangement repeated in every test instead of using `beforeEach`
- **Over-mocking in setup**: Mocking things not needed for most tests in the suite

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 3: Setup and Teardown
===========================

beforeEach:
  Fresh mock creation:          [PASS/FAIL]
  TestingModule built:          [PASS/FAIL]
  .setLogger() included:        [PASS/FAIL]
  Target retrieved from module: [PASS/FAIL]

afterEach:
  jest.clearAllMocks():         [PASS/FAIL]
  No shared mutable state:      [PASS/FAIL]

Anti-Patterns Found: N
  - [SEVERITY]: description
```

Then ask: **[C] Continue to Step 4: AAA Pattern**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-aaa-pattern.md`.
