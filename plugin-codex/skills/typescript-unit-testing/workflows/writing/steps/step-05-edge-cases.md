---
name: 'step-05-edge-cases'
description: 'Implement edge case tests for boundary conditions and special values'
nextStepFile: './step-06-error-cases.md'
referenceFiles:
  - 'references/common/rules.md'
---

# Step 5: Implement Edge Case Tests

## STEP GOAL

Write edge case tests covering empty values, null/undefined, boundary conditions, and other special inputs.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements

## EXECUTION

### 1. Write Edge Case Tests

Cover the following categories for each relevant public method:

#### Empty Arrays/Strings

```typescript
it('should handle empty array when no items found', async () => {
  // Arrange
  mockRepository.findAll.mockResolvedValue([]);

  // Act
  const result = await target.getAllItems();

  // Assert
  expect(result).toEqual([]);
  expect(mockRepository.findAll).toHaveBeenCalledTimes(1);
});
```

#### Null/Undefined Values

```typescript
it('should handle null value when optional field missing', async () => {
  // Arrange
  const input = { requiredField: 'value', optionalField: null };
  mockService.process.mockResolvedValue({ processed: true });

  // Act
  const result = await target.process(input);

  // Assert
  expect(result.optionalField).toBeNull();
});
```

#### Maximum/Minimum Values

```typescript
it('should handle maximum value at upper boundary', async () => {
  // Arrange
  const input = { count: Number.MAX_SAFE_INTEGER };

  // Act
  const result = await target.process(input);

  // Assert
  expect(result).toBeDefined(); // verify correct behavior at boundary
});
```

#### Zero Values

```typescript
it('should handle zero value correctly', async () => {
  // Arrange
  const input = { amount: 0 };

  // Act
  const result = await target.calculate(input);

  // Assert
  expect(result.total).toBe(0);
});
```

### 2. Boundary Conditions Checklist

Ensure coverage of:
- [ ] Empty arrays/strings
- [ ] Null/undefined values
- [ ] Maximum/minimum values
- [ ] Zero values
- [ ] Single-element arrays
- [ ] Whitespace-only strings (if relevant)
- [ ] Boundary values for numeric ranges

### 3. Append to Report

Append to the output document:

```markdown
## Step 5: Edge Case Tests

**Tests Written:**
- [methodName]: should handle [edge case] - DONE

**Edge Cases Covered:**
- Empty arrays/strings: [yes/no/N/A]
- Null/undefined: [yes/no/N/A]
- Max/min values: [yes/no/N/A]
- Zero values: [yes/no/N/A]
- Boundary conditions: [yes/no/N/A]
```

## PRESENT FINDINGS

Show the user:
- All edge case tests written
- Coverage of each boundary condition category
- Any edge cases that were not applicable to this component

Then ask: **[C] Continue to Step 6: Implement Error Case Tests**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `5` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-06-error-cases.md`.
