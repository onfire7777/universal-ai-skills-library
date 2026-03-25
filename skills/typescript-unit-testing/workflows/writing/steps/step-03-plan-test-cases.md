---
name: 'step-03-plan-test-cases'
description: 'Plan test cases for all public methods before implementation'
nextStepFile: './step-04-happy-path.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/assertions.md'
---

# Step 3: Plan Test Cases

## STEP GOAL

Plan all test cases before implementation. Create a comprehensive test plan covering happy paths, edge cases, error cases, and business rules for every public method.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements
- `references/common/assertions.md` — Assertion patterns and matchers

## EXECUTION

### 1. Plan Test Cases for Each Public Method

For each public method identified in Step 1, create a test plan table:

**Method: `[methodName]`**

| Category | Test Case | Priority |
|----------|-----------|----------|
| Happy path | should [expected behavior] when [valid input] | MANDATORY |
| Edge case | should [expected behavior] when [empty/null/boundary] | MANDATORY |
| Error case | should throw [Exception] when [condition] | MANDATORY |
| Business rule | should [expected behavior] when [business condition] | MANDATORY |

### 2. Ensure Complete Coverage

For each method, verify the test plan covers:
- **Happy path**: At least one test with valid, typical input
- **Edge cases**: Empty arrays/strings, null/undefined, boundary values
- **Error cases**: All exceptions thrown by the method
- **Business rules**: All conditional logic branches, calculations, validations

### 3. Present Test Plan to User

Present the complete test plan and ask the user to:
- Confirm the plan is complete
- Add any missing test cases
- Adjust priorities if needed

### 4. Append to Report

Append to the output document:

```markdown
## Step 3: Test Plan

### Method: [methodName]
| Category | Test Case | Priority |
|----------|-----------|----------|
| ... | ... | ... |

**Total Test Cases**: [count]
```

## PRESENT FINDINGS

Show the user:
- Complete test plan table for each public method
- Total number of test cases planned
- Any coverage gaps identified

Then ask: **[C] Continue to Step 4: Implement Happy Path Tests**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-04-happy-path.md`.
