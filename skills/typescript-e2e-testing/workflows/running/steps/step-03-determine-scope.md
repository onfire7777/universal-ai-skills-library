---
name: 'step-03-determine-scope'
description: 'Determine which tests to run — full suite, file, test name, pattern'
nextStepFile: './step-04-execute.md'
---

# Step 3: Determine Run Scope

## STEP GOAL

Identify which tests to run.

## EXECUTION

### 1. Options

1. **Full Suite**: Run all E2E tests
2. **Specific File**: Run tests in one file
3. **Specific Test**: Run a single test by name
4. **By Pattern**: Run tests matching a pattern

### 2. Commands

```bash
# Full suite
npm run test:e2e

# Specific file
npm run test:e2e -- test/e2e/user.e2e-spec.ts

# Specific test by name
npm run test:e2e -- -t "should create user"

# Pattern matching
npm run test:e2e -- -t "User API"

# Multiple files
npm run test:e2e -- test/e2e/user.e2e-spec.ts test/e2e/order.e2e-spec.ts
```

### 3. Append to Output

Append to the output document:

```markdown
## Step 3: Test Scope

**Selected Scope**: {{scope description}}
**Command**: {{command to run}}
```

## PRESENT FINDINGS

```
Test Run Scope

Available Options:
1. [A] All E2E tests ({n} files, {m} tests)
2. [F] Specific file
3. [T] Specific test by name
4. [P] Tests matching pattern

Current recommendation: {based on context}

Select option:
```

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms scope, load `step-04-execute.md`.
