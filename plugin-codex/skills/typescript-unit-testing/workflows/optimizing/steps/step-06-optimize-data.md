---
name: 'step-06-optimize-data'
description: 'Optimize test data patterns'
nextStepFile: './step-07-optimize-organization.md'
referenceFiles:
  - 'references/mocking/factories.md'
---

# Step 6: Optimize Test Data

## STEP GOAL

Optimize test data creation and mock data to reduce verbosity, improve maintainability, and eliminate unnecessary data setup overhead.

## REFERENCE LOADING

Before starting, load and read:
- `references/mocking/factories.md` — Factory patterns for test data creation

## EXECUTION

### Technique A: Use Factories

Replace repeated inline object creation with reusable factory functions:

```typescript
// BEFORE (verbose, repeated)
it('test 1', () => {
  const user = {
    id: 'user-1',
    email: 'test@example.com',
    name: 'Test User',
    role: 'admin',
    createdAt: new Date(),
  };
});

it('test 2', () => {
  const user = {
    id: 'user-2',
    email: 'test2@example.com',
    name: 'Test User 2',
    role: 'user',
    createdAt: new Date(),
  };
});

// AFTER (concise, reusable)
const createUser = (overrides = {}): User => ({
  id: 'user-1',
  email: 'test@example.com',
  name: 'Test User',
  role: 'user',
  createdAt: new Date(),
  ...overrides,
});

it('test 1', () => {
  const user = createUser({ role: 'admin' });
});

it('test 2', () => {
  const user = createUser({ id: 'user-2', email: 'test2@example.com' });
});
```

**When to apply**: Any test file where the same object shape is created more than twice with minor variations. Factories reduce duplication and make tests more readable by highlighting only the relevant differences.

### Technique B: Minimize Mock Data

Only include fields that are actually used by the code under test:

```typescript
// BEFORE (excessive data)
mockRepository.findById.mockResolvedValue({
  id: 'user-1',
  email: 'test@example.com',
  firstName: 'John',
  lastName: 'Doe',
  address: { street: '123 Main', city: 'NYC', zip: '10001' },
  preferences: { theme: 'dark', notifications: true },
  // ... 20 more fields
});

// AFTER (minimal required data)
mockRepository.findById.mockResolvedValue({
  id: 'user-1',
  email: 'test@example.com',
  // Only fields actually used by the test
});
```

**When to apply**: Any mock return value that includes fields not referenced by the code path under test. Trace the code to determine which fields are actually accessed.

### Apply Optimizations

For each test file:
1. Identify repeated object creation patterns
2. Extract factory functions for common entities
3. Trim mock return values to only required fields
4. Verify tests still pass:
   ```bash
   npm test -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-data-check.log 2>&1
   tail -20 /tmp/ut-${UT_SESSION}-data-check.log
   ```

## PRESENT FINDINGS

Present findings to the user:

```
Step 6: Data Optimizations
============================

Optimizations Applied:
  [file.spec.ts]
    - Technique A: Created [N] factory functions   [APPLIED/SKIPPED]
    - Technique B: Minimized [N] mock data objects  [APPLIED/SKIPPED]

Verification:
  All modified test files: [PASS/FAIL]
```

Then ask: **[C] Continue to Step 7: Optimize Test Organization**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`
- Append the data optimization details to the report

## NEXT STEP

After user confirms `[C]`, load `step-07-optimize-organization.md`.
