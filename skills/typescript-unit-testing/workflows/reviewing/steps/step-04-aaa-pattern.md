---
name: 'step-04-aaa-pattern'
description: 'Review AAA pattern compliance'
nextStepFile: './step-05-assertions.md'
referenceFiles:
  - 'references/common/rules.md'
---

# Step 4: Review AAA Pattern Compliance

## STEP GOAL

Verify every test follows the Arrange-Act-Assert pattern with proper separation, comments, and single-action discipline.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### 1. Verify AAA Pattern in Each Test

For each test, check all three phases:

| Test | Arrange | Act | Assert | Comments |
|------|---------|-----|--------|----------|
| test name | [ ] | [ ] | [ ] | [ ] |

- **Arrange**: Sets up test data, configures mocks
- **Act**: Calls the method under test (single action)
- **Assert**: Verifies the result and/or mock interactions

### 2. Check for AAA Violations

For each test, look for:
- **Missing phase comments**: Each phase should have a `// Arrange`, `// Act`, `// Assert` comment
- **Multiple actions in single test**: Only ONE call to the method under test per `it` block
- **Assertions in Arrange phase**: Assertions should only appear in the Assert phase
- **Missing assertions**: Every test must have at least one assertion
- **Blurred boundaries**: Act and Assert combined (e.g., `expect(target.method()).toEqual(...)` — acceptable for simple cases but note it)

### 3. Document AAA Issues

```
**AAA Pattern Issues:**
- [test name]: [issue description]
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 4: AAA Pattern Compliance
===============================

Test-by-Test Review:
  [PASS/FAIL] "should [expected] when [condition]"
    - Arrange: [OK / MISSING / ISSUE]
    - Act: [OK / MULTIPLE ACTIONS / MISSING]
    - Assert: [OK / MISSING / WEAK]
    - Comments: [OK / MISSING]

Summary: N tests reviewed, N with AAA violations

Issues:
  - [test name]: [description]
```

Then ask: **[C] Continue to Step 5: Assertions Quality**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-assertions.md`.
