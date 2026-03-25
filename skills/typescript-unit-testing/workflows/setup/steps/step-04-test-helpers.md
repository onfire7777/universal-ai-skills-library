---
name: 'step-04-test-helpers'
description: 'Create test helper utilities'
nextStepFile: './step-05-verify-setup.md'
---

# Step 4: Create Test Helpers

## STEP GOAL

Create the test helpers directory structure and the `MockLoggerService` utility so tests can suppress logger output and have a consistent mock available.

## EXECUTION

### 1. Create Directory Structure

Ensure the following directory structure exists in the project:

```
test/
└── helpers/
    └── mock-logger.service.ts
```

### 2. Check for Existing Logger

Before creating the mock:
- Check if the project has a custom logger service (e.g., extending `LoggerService`)
- If a custom logger exists, adapt the `MockLoggerService` to match its interface

### 3. Create MockLoggerService

Create `test/helpers/mock-logger.service.ts`:

```typescript
// test/helpers/mock-logger.service.ts
import { LoggerService } from '@nestjs/common';

export class MockLoggerService implements LoggerService {
  log(): void {}
  error(): void {}
  warn(): void {}
  debug(): void {}
  verbose(): void {}
}
```

If the project has a custom logger with additional methods, extend the mock accordingly.

### 4. Append to Report

Append to the output document:

```markdown
## Step 4: Test Helpers

**Directory Created**: test/helpers/
**Files Created:**
- `test/helpers/mock-logger.service.ts` — MockLoggerService

**Custom Logger**: [Found at path — adapted mock / Not found — using standard NestJS LoggerService interface]
```

## PRESENT FINDINGS

Present the test helper status:

```
Step 4: Test Helpers
====================

Directory: test/helpers/ [created/already existed]
Files:
  mock-logger.service.ts  [created/already existed]

Custom Logger: [Adapted for custom logger at path / Standard NestJS interface used]
```

Then ask: **[C] Continue to Step 5: Verify Setup**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Append the test helpers section to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-verify-setup.md`.
