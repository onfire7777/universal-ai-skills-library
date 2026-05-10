---
name: 'step-05-verify-setup'
description: 'Verify setup with sample test'
nextStepFile: './step-06-npm-scripts.md'
---

# Step 5: Verify Setup with Sample Test

## STEP GOAL

Create and run a verification test that exercises the core testing infrastructure — `@nestjs/testing` module creation, `MockLoggerService`, and `@golevelup/ts-jest` deep mocking — to confirm the setup is working correctly.

## EXECUTION

### 1. Create Verification Test

Create `test/setup.spec.ts` with the following content (skip if the file already exists):

```typescript
// test/setup.spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { createMock, DeepMocked } from '@golevelup/ts-jest';
import { MockLoggerService } from './helpers/mock-logger.service';

describe('Test Setup Verification', () => {
  it('should create testing module successfully', async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [],
    })
      .setLogger(new MockLoggerService())
      .compile();

    expect(module).toBeDefined();
  });

  it('should create mocks with DeepMocked', () => {
    interface SampleService {
      getData(): Promise<string>;
    }

    const mockService: DeepMocked<SampleService> = createMock<SampleService>();
    mockService.getData.mockResolvedValue('test');

    expect(mockService.getData).toBeDefined();
  });
});
```

### 2. Run Verification Test

Execute the verification test:

```bash
npm test -- test/setup.spec.ts
```

### 3. Evaluate Results

- If tests pass: Setup is confirmed working
- If tests fail: Analyze the error output and identify the root cause (missing dependency, configuration issue, path resolution problem)
- Propose and apply fixes if needed, then re-run

### 4. Append to Report

Append to the output document:

```markdown
## Step 5: Setup Verification

**Verification Test**: test/setup.spec.ts [created/already existed]
**Test Results:**
- TestingModule creation: [PASS/FAIL]
- DeepMocked creation: [PASS/FAIL]

**Issues Found**: [list issues and fixes applied, or "None"]
```

## PRESENT FINDINGS

Present the verification results:

```
Step 5: Setup Verification
==========================

Verification Test: test/setup.spec.ts
  [PASS] should create testing module successfully
  [PASS] should create mocks with DeepMocked

Result: [All tests passing / Issues found and fixed / Issues require attention]
```

Then ask: **[C] Continue to Step 6: Configure npm Scripts**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the verification section to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-npm-scripts.md`.
