---
name: 'step-05-optimize-async'
description: 'Optimize async operations in tests'
nextStepFile: './step-06-optimize-data.md'
referenceFiles:
  - 'references/common/performance-optimization.md'
---

# Step 5: Optimize Async Operations

## STEP GOAL

Optimize async patterns in tests to reduce unnecessary overhead from awaiting synchronous operations, real timeouts, and suboptimal parallelization.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/performance-optimization.md` — Worker config, caching, CI optimization

## EXECUTION

### Technique A: Remove Unnecessary Async

Do not use `async`/`await` when the method under test is synchronous:

```typescript
// BEFORE
it('should return value', async () => {
  const result = await target.syncMethod(); // Method is not async
  expect(result).toBe(expected);
});

// AFTER
it('should return value', () => {
  const result = target.syncMethod();
  expect(result).toBe(expected);
});
```

**When to apply**: Any test that uses `async`/`await` but the method under test returns a plain value, not a Promise. Check the source code to confirm whether the method is actually async.

### Technique B: Mock Time Instead of Waiting

Replace real timeouts with fake timers to eliminate wait time:

```typescript
// BEFORE (slow - waits 1 second)
it('should retry after delay', async () => {
  mockService.call
    .mockRejectedValueOnce(new Error('fail'))
    .mockResolvedValue('success');

  const result = await target.callWithRetry();
  expect(result).toBe('success');
}, 5000);

// AFTER (fast - mocks time)
it('should retry after delay', async () => {
  jest.useFakeTimers();

  mockService.call
    .mockRejectedValueOnce(new Error('fail'))
    .mockResolvedValue('success');

  const promise = target.callWithRetry();
  jest.advanceTimersByTime(1000);
  const result = await promise;

  expect(result).toBe('success');
  jest.useRealTimers();
});
```

**When to apply**: Any test with explicit timeout values (the second argument to `it()`), or tests that call code with `setTimeout`/`setInterval`/retry logic. This can eliminate seconds of wait time per test.

### Technique C: Parallel Test Execution

Ensure tests are independent for parallel execution and configure optimal worker count:

```typescript
// jest.config.ts
export default {
  maxWorkers: '50%', // Use half of CPU cores
  // OR
  maxWorkers: 4, // Fixed number
};
```

**When to apply**: When tests are independent (no shared mutable state between test files). Check that no test file depends on global state, shared database records, or file system state from another test file.

### Apply Optimizations

For each test file:
1. Scan for unnecessary `async`/`await` usage
2. Identify tests with real timeouts that can use fake timers
3. Verify test independence for parallel execution
4. Apply changes and verify:
   ```bash
   npm test -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-async-check.log 2>&1
   tail -20 /tmp/ut-${UT_SESSION}-async-check.log
   ```

## PRESENT FINDINGS

Present findings to the user:

```
Step 5: Async Optimizations
=============================

Optimizations Applied:
  [file.spec.ts]
    - Technique A: Removed unnecessary async    [APPLIED/SKIPPED]
    - Technique B: Mocked time for [N] tests    [APPLIED/SKIPPED]

  Jest Config:
    - Technique C: maxWorkers set to [value]    [APPLIED/SKIPPED]

Verification:
  All modified test files: [PASS/FAIL]
```

Then ask: **[C] Continue to Step 6: Optimize Test Data**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`
- Append the async optimization details to the report

## NEXT STEP

After user confirms `[C]`, load `step-06-optimize-data.md`.
