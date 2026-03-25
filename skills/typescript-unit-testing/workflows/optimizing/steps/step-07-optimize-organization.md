---
name: 'step-07-optimize-organization'
description: 'Optimize test organization'
nextStepFile: './step-08-measure-improvement.md'
---

# Step 7: Optimize Test Organization

## STEP GOAL

Improve test organization for better performance and maintainability through logical grouping, redundancy elimination, and test categorization.

## EXECUTION

### Technique A: Group Related Tests

Organize tests into `describe` blocks by method for shared context and better readability:

```typescript
describe('UserService', () => {
  // Group by method
  describe('findById', () => {
    // All findById tests share context
  });

  describe('create', () => {
    // All create tests share context
  });
});
```

**When to apply**: Test files where tests for different methods are interleaved or not logically grouped. Grouping by method allows shared setup within each group and makes it easier to run specific method tests.

### Technique B: Skip Redundant Tests

Identify and skip tests that duplicate coverage already provided by other tests:

```typescript
// If testing same logic through multiple paths, skip redundant ones
it.skip('redundant test covered by other tests', () => {});

// Or use conditional skipping
const skipIfCI = process.env.CI ? it.skip : it;
skipIfCI('expensive test skipped in CI', () => {});
```

**When to apply**: Use sparingly. Only skip tests that are truly redundant (testing the exact same code path with equivalent inputs). Never skip tests just because they are slow — optimize them instead. Conditional skipping is appropriate for tests that require external resources unavailable in CI.

### Technique C: Use Test Tags

Use Jest's path filtering to run specific categories of tests:

```bash
# Run only fast tests
npm test -- --testPathIgnorePatterns="slow"

# Run only specific category
npm test -- --testPathPattern="unit"
```

**When to apply**: When the test suite has grown large enough to benefit from selective execution. Consider organizing test files into directories or using naming conventions (e.g., `*.unit.spec.ts`, `*.integration.spec.ts`) that enable path-based filtering.

### Apply Optimizations

For each test file:
1. Review `describe` block organization
2. Identify any truly redundant tests
3. Consider test categorization for the overall suite
4. Verify tests still pass:
   ```bash
   npm test -- [path/to/file.spec.ts] > /tmp/ut-${UT_SESSION}-org-check.log 2>&1
   tail -20 /tmp/ut-${UT_SESSION}-org-check.log
   ```

## PRESENT FINDINGS

Present findings to the user:

```
Step 7: Organization Optimizations
=====================================

Optimizations Applied:
  [file.spec.ts]
    - Technique A: Reorganized into [N] describe blocks  [APPLIED/SKIPPED]
    - Technique B: Skipped [N] redundant tests           [APPLIED/SKIPPED]

  Suite-wide:
    - Technique C: Test categorization                   [APPLIED/SKIPPED]

Verification:
  All modified test files: [PASS/FAIL]
```

Then ask: **[C] Continue to Step 8: Measure Improvement**

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`
- Append the organization optimization details to the report

## NEXT STEP

After user confirms `[C]`, load `step-08-measure-improvement.md`.
