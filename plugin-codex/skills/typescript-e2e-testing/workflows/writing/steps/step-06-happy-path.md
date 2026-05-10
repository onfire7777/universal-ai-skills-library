---
name: 'step-06-happy-path'
description: 'Write happy path test cases following GWT pattern'
nextStepFile: './step-07-error-cases.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/test-case-creation-guide.md'
---

# Step 6: Write Happy Path Tests

## STEP GOAL

Write each happy path test case following GWT pattern.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — GWT rules
- `references/common/test-case-creation-guide.md` — test case templates

## EXECUTION

For each happy path test case:

### 1. GIVEN Section (Setup)
- Create test data using factories
- Seed database if needed
- Configure mocks for external APIs
- Use unique identifiers (timestamp + random)

### 2. WHEN Section (Action)
- Execute EXACTLY ONE action
- Capture response/result
- For async: Trigger event and wait for result

### 3. THEN Section (Assertions)
- Assert specific values, not just existence
- Use `toMatchObject` for partial matching
- Verify database state changes
- Verify async side effects (Kafka messages)

### GWT Template

```typescript
it('should {expected behavior}', async () => {
  // GIVEN: {describe setup}
  const testData = createTestData({
    id: `test-${Date.now()}-${Math.random().toString(36).slice(2)}`,
  });
  await repository.save(testData);

  // WHEN: {describe action}
  const response = await request(httpServer)
    .post('/endpoint')
    .send(requestBody)
    .expect(201);

  // THEN: {describe expectations}
  expect(response.body).toMatchObject({
    code: 'SUCCESS',
    data: {
      id: expect.any(String),
      field: expectedValue,
    },
  });

  // Verify database state
  const saved = await repository.findOne({ where: { id: response.body.data.id } });
  expect(saved).toMatchObject({ field: expectedValue });
});
```

### Append to Output

Append to the output document:

```markdown
## Step 6: Happy Path Tests

### Test: should {{description}}
{{show code}}

Verifies:
- {{verification 1}}
- {{verification 2}}

Pattern Reference: references/common/test-case-creation-guide.md
```

## PRESENT FINDINGS

For each test case, present the code and what it verifies.

Then ask: **[C] Continue to Step 7: Error Cases / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-07-error-cases.md`.
