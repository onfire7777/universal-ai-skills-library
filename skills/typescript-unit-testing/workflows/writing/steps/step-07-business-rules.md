---
name: 'step-07-business-rules'
description: 'Implement business rule tests for domain logic and conditional branches'
nextStepFile: './step-08-verify-mocks.md'
referenceFiles:
  - 'references/common/rules.md'
---

# Step 7: Implement Business Rule Tests

## STEP GOAL

Write tests for domain logic, business rules, conditional branches, and calculations.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements

## EXECUTION

### 1. Write Business Rule Tests

Test all business logic branches, calculations, and conditions identified in Step 1.

#### Conditional Logic — Both Branches

```typescript
it('should apply discount when user is premium member', async () => {
  // Arrange
  const user = { id: 'user-1', isPremium: true };
  const order = { total: 100 };
  mockUserService.getUser.mockResolvedValue(user);

  // Act
  const result = await target.calculateTotal(order, user.id);

  // Assert
  expect(result.discount).toBe(10); // 10% premium discount
  expect(result.finalTotal).toBe(90);
});

it('should not apply discount when user is not premium member', async () => {
  // Arrange
  const user = { id: 'user-1', isPremium: false };
  const order = { total: 100 };
  mockUserService.getUser.mockResolvedValue(user);

  // Act
  const result = await target.calculateTotal(order, user.id);

  // Assert
  expect(result.discount).toBe(0);
  expect(result.finalTotal).toBe(100);
});
```

### 2. Business Rule Categories

For each business rule, ensure tests cover:
- [ ] All conditional branches (if/else, switch cases)
- [ ] Calculations and formulas (verify exact numeric results)
- [ ] Validation rules (what passes, what fails)
- [ ] State transitions (before and after)
- [ ] Permission checks (authorized vs unauthorized)
- [ ] Data transformations (input shape vs output shape)

### 3. Test All Branches

For every `if/else`, `switch`, or ternary in the component:
- Write a test for the **true/matching** branch
- Write a test for the **false/default** branch
- Verify the correct side effects for each branch

### 4. Append to Report

Append to the output document:

```markdown
## Step 7: Business Rule Tests

**Tests Written:**
- [methodName]: should [business behavior] when [condition] - DONE

**Branch Coverage:**
- [methodName]: [N] branches, [N] tested
```

## PRESENT FINDINGS

Show the user:
- All business rule tests written
- Branch coverage for each method with conditional logic
- Any business rules that were ambiguous (ask for clarification)

Then ask: **[C] Continue to Step 8: Verify Mock Interactions**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `7` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-08-verify-mocks.md`.
