---
name: 'step-03-optimize-setup'
description: 'Optimize test setup patterns'
nextStepFile: './step-04-fix-open-handles.md'
referenceFiles:
  - 'references/common/performance-optimization.md'
---

# Step 3: Optimize Test Setup

## STEP GOAL

Optimize test module creation and mock setup to reduce per-test overhead. Apply shared setup, lazy mock creation, and minimal module patterns.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/performance-optimization.md` — Worker config, caching, CI optimization

## EXECUTION

### Technique A: Share Expensive Setup

Move one-time setup to `beforeAll` and keep only mock resets in `beforeEach`:

```typescript
// BEFORE (slow)
describe('Service', () => {
  let target: Service;

  beforeEach(async () => {
    // Expensive: compiles module for every test
    const module = await Test.createTestingModule({
      providers: [Service, ...manyProviders],
    }).compile();
    target = module.get(Service);
  });
});

// AFTER (fast)
describe('Service', () => {
  let module: TestingModule;
  let target: Service;

  beforeAll(async () => {
    // One-time: compiles module once
    module = await Test.createTestingModule({
      providers: [Service, ...manyProviders],
    }).compile();
  });

  beforeEach(() => {
    target = module.get(Service);
    jest.clearAllMocks(); // Reset mock state
  });
});
```

**When to apply**: Any test file where `Test.createTestingModule` is inside `beforeEach`. This is the single highest-impact optimization for NestJS tests.

### Technique B: Lazy Mock Creation

Only create mocks when the specific test needs them:

```typescript
// BEFORE
beforeEach(() => {
  mockServiceA = createMock<ServiceA>();
  mockServiceB = createMock<ServiceB>();
  mockServiceC = createMock<ServiceC>();
  // Creates all mocks even if test only uses one
});

// AFTER
const getMockServiceA = () => createMock<ServiceA>();
const getMockServiceB = () => createMock<ServiceB>();

it('test using only service A', () => {
  const mockA = getMockServiceA();
  // Only creates what's needed
});
```

**When to apply**: Test files with many mocks where individual tests only use a subset.

### Technique C: Reduce Module Complexity

Only include required providers, not the full AppModule:

```typescript
// BEFORE (slow)
const module = await Test.createTestingModule({
  imports: [AppModule], // Imports everything
}).compile();

// AFTER (fast)
const module = await Test.createTestingModule({
  providers: [
    TargetService,
    { provide: DependencyA, useValue: mockA },
    { provide: DependencyB, useValue: mockB },
  ],
}).compile();
```

**When to apply**: Any test importing `AppModule` or large feature modules. Unit tests should only include the target and its direct mocked dependencies.

### Apply Optimizations

For each test file identified in Step 2:
1. Check which techniques apply
2. Apply changes one technique at a time
3. Run the specific test file after each change to verify tests still pass:
   ```bash
   npm test -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-setup-check.log 2>&1
   tail -20 /tmp/ut-${UT_SESSION}-setup-check.log
   ```

## PRESENT FINDINGS

Present findings to the user:

```
Step 3: Setup Optimizations
=============================

Optimizations Applied:
  [file.spec.ts]
    - Technique A: Moved TestingModule to beforeAll  [APPLIED/SKIPPED]
    - Technique B: Lazy mock creation                [APPLIED/SKIPPED]
    - Technique C: Reduced module complexity          [APPLIED/SKIPPED]

Verification:
  All modified test files: [PASS/FAIL]
```

Then ask: **[C] Continue to Step 4: Fix Open Handles (CRITICAL)**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the setup optimization details to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-fix-open-handles.md`.
