---
name: 'step-07-error-cases'
description: 'Write error case tests — validation, not found, authorization, no side effects'
nextStepFile: './step-08-async-tests.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/test-case-creation-guide.md'
---

# Step 7: Write Error Case Tests

## STEP GOAL

Test error handling paths with proper verification that no side effects occur.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — error test rules
- `references/common/test-case-creation-guide.md` — error test templates

## EXECUTION

### 1. Error Cases to Cover

1. Test validation errors (400)
2. Test not found errors (404)
3. Test authorization errors (401/403)
4. Test business logic errors
5. Verify NO side effects occur on error

### 2. Error Test Template

```typescript
it('should return 400 when {invalid condition}', async () => {
  // GIVEN: Invalid input
  const invalidData = { field: 'invalid-value' };

  // WHEN: Attempting action
  const response = await request(httpServer)
    .post('/endpoint')
    .send(invalidData)
    .expect(400);

  // THEN: Validation error returned
  expect(response.body).toMatchObject({
    code: 'VALIDATION_ERROR',
    errors: expect.arrayContaining([
      expect.objectContaining({ field: 'field' }),
    ]),
  });

  // AND: No side effects
  const count = await repository.count();
  expect(count).toBe(0);
});
```

### 3. Append to Output

Append to the output document:

```markdown
## Step 7: Error Case Tests

### Test: should return {{status}} when {{condition}}
{{show code}}

Verifies:
- Correct error response
- No database changes
- No events published
```

## PRESENT FINDINGS

For each error case, present the code and verifications.

Then ask: **[C] Continue to Step 8: Async Tests / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `7` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-08-async-tests.md`.
