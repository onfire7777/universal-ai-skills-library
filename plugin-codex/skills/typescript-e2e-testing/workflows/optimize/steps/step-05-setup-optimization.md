---
name: 'step-05-setup-optimization'
description: 'Reduce setup and teardown overhead'
nextStepFile: './step-06-execution-optimization.md'
---

# Step 5: Test Setup Optimization

## STEP GOAL

Reduce setup and teardown overhead.

## EXECUTION

### Optimizations

#### 1. Reuse App Instance

```typescript
// ❌ App per file (slow)
describe('Feature', () => {
  beforeAll(async () => {
    app = await createTestApp();
  });
  afterAll(async () => {
    await app.close();
  });
});

// ✅ Global app (if tests are isolated)
let app: INestApplication;
export const getTestApp = async () => {
  if (!app) {
    app = await createTestApp();
  }
  return app;
};
```

#### 2. Optimize beforeEach

```typescript
// ❌ Over-aggressive cleanup
beforeEach(async () => {
  await new Promise(r => setTimeout(r, 2000)); // Too long
  await userRepository.clear();
  await orderRepository.clear();
  await productRepository.clear();
  await auditRepository.clear(); // Unnecessary if not used
});

// ✅ Targeted cleanup
beforeEach(async () => {
  await new Promise(r => setTimeout(r, 500)); // Shorter wait
  // Only clean tables used in this describe block
  await userRepository.clear();
});
```

#### 3. Lazy Initialization

```typescript
// ❌ Always initialize everything
beforeAll(async () => {
  kafkaHelper = new KafkaTestHelper();
  await kafkaHelper.connect();
  await kafkaHelper.subscribeToAllTopics(); // Expensive
});

// ✅ Initialize on demand
let kafkaHelper: KafkaTestHelper;
const getKafkaHelper = async () => {
  if (!kafkaHelper) {
    kafkaHelper = new KafkaTestHelper();
    await kafkaHelper.connect();
  }
  return kafkaHelper;
};
```

### Append to Output

Append setup optimization details to the output document.

## PRESENT FINDINGS

```
Setup Optimization

Current beforeAll time: {X}s
Current beforeEach time: {Y}s per test

Proposed Changes:
1. {optimization 1}
   - Savings: {time}

2. {optimization 2}
   - Savings: {time}

Estimated Total Savings: {time}

[A] Apply / [M] Modify / [S] Skip
```

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`

## NEXT STEP

After user confirms, load `step-06-execution-optimization.md`.
