---
name: 'step-08-example-test'
description: 'Create a sample E2E test file demonstrating patterns'
nextStepFile: './step-09-verification.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/examples.md'
---

# Step 8: Create Example Test

## STEP GOAL

Create a sample E2E test file demonstrating GWT pattern, lifecycle hooks, helper usage, and cleanup patterns.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — GWT and test rules
- `references/common/examples.md` — reference implementations

## EXECUTION

### 1. Create Example Test

Create `test/e2e/example.e2e-spec.ts`:
- Proper GWT pattern
- beforeAll/beforeEach/afterAll lifecycle
- Helper usage examples
- Cleanup patterns

### 2. Append to Output

Append to the output document:

```markdown
## Step 8: Example Test

**File**: test/e2e/example.e2e-spec.ts
{{show code}}

This demonstrates:
- GWT pattern usage
- Test isolation
- Helper integration
- Proper cleanup
```

## PRESENT FINDINGS

Show the example test code and what it demonstrates.

Then ask: **[C] Continue to Step 9: Verification / [S] Skip**

## FRONTMATTER UPDATE

Update the output document:
- Add `8` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-09-verification.md`.
