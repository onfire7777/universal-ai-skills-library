---
name: 'step-04-jest-config'
description: 'Create or update Jest E2E configuration'
nextStepFile: './step-05-docker-compose.md'
referenceFiles:
  - 'references/common/nestjs-setup.md'
  - 'references/common/rules.md'
---

# Step 4: Configure Jest for E2E

## STEP GOAL

Create or update Jest E2E configuration with proper settings for sequential execution, timeouts, and test matching.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/nestjs-setup.md` — Jest configuration patterns
- `references/common/rules.md` — mandatory configuration rules

## EXECUTION

### 1. Create jest-e2e.config.ts

Create `test/jest-e2e.config.ts` with:
- `preset: 'ts-jest'`
- `testEnvironment: 'node'`
- `testMatch: ['**/*.e2e-spec.ts']`
- `testTimeout: 25000`
- `maxWorkers: 1` (CRITICAL for E2E isolation)
- `setupFilesAfterEnv` if needed
- `forceExit: true`
- `detectOpenHandles: true`

### 2. Update package.json

Add/update scripts:
- `"test:e2e": "jest --config test/jest-e2e.config.ts"`

### 3. Append to Output

Append to the output document:

```markdown
## Step 4: Jest Configuration

**File**: test/jest-e2e.config.ts
{{show configuration}}

**package.json scripts updated**:
- test:e2e: jest --config test/jest-e2e.config.ts
```

## PRESENT FINDINGS

Show the configuration to the user.

Then ask: **[C] Continue to Step 5: Docker Compose / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-05-docker-compose.md`.
