---
name: 'step-06-apply-fix'
description: 'Apply targeted fix — minimum change needed'
nextStepFile: './step-07-verify-fix.md'
---

# Step 6: Apply Targeted Fix

## STEP GOAL

Fix the specific root cause without broad changes.

## EXECUTION

### Fix Principles

1. Make the MINIMUM change needed
2. Don't refactor unrelated code
3. Don't "clean up" while fixing
4. Document what was changed

### Common Fixes by Category

#### Timeout Fixes
```typescript
// Increase specific timeout
it('should process event', async () => {
  // ...
}, 30000);

// Add polling instead of fixed wait
const result = await waitFor(
  () => repository.findOne({ id }),
  15000,  // timeout
  100     // poll interval
);

// Ensure Kafka subscription before test
beforeAll(async () => {
  await kafkaHelper.subscribeToTopic(outputTopic);
  await new Promise(r => setTimeout(r, 2000));
}, 90000);
```

#### Assertion Fixes
```typescript
// Fix expected value
expect(response.body.code).toBe('SUCCESS');

// Use partial matching for flexible assertions
expect(response.body).toMatchObject({
  code: 'SUCCESS',
});

// Verify precondition
const setupCount = await repository.count();
expect(setupCount).toBe(1);
```

#### Connection Fixes
```typescript
// Add retry logic
const connectWithRetry = async (retries = 3) => {
  for (let i = 0; i < retries; i++) {
    try {
      await dataSource.initialize();
      return;
    } catch (e) {
      if (i === retries - 1) throw e;
      await new Promise(r => setTimeout(r, 1000));
    }
  }
};
```

#### Race Condition Fixes
```typescript
// Ensure cleanup complete
beforeEach(async () => {
  await new Promise(r => setTimeout(r, 1000));
  await repository.clear();
  kafkaHelper.clearMessages(topic);
  jest.clearAllMocks();
});

// Use unique identifiers
const id = `test-${Date.now()}-${Math.random().toString(36).slice(2)}`;
```

### Append to Output

Append fix details to the output document.

## PRESENT FINDINGS

```
Proposed Fix

File: {file}
Line: {line}

Change:
- {old code}
+ {new code}

Rationale: {why this fixes the issue}

[A] Apply fix / [M] Modify fix / [S] Suggest alternative
```

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`

## NEXT STEP

After fix is applied, load `step-07-verify-fix.md`.
