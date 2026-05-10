---
name: 'step-07-global-setup'
description: 'Create global E2E test setup file'
nextStepFile: './step-08-example-test.md'
referenceFiles:
  - 'references/common/nestjs-setup.md'
---

# Step 7: Create Global Setup

## STEP GOAL

Create the global E2E test setup file with environment loading, timeout configuration, and global hooks.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/nestjs-setup.md` — global setup patterns

## EXECUTION

### 1. Create setup.ts

Create `test/e2e/setup.ts`:
- Environment variable loading
- Test timeout configuration
- Global beforeAll/afterAll hooks

### 2. Append to Output

Append to the output document:

```markdown
## Step 7: Global Setup

**File**: test/e2e/setup.ts
{{show code}}
```

## PRESENT FINDINGS

Show the global setup code.

Then ask: **[C] Continue to Step 8: Example Test / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-08-example-test.md`.
