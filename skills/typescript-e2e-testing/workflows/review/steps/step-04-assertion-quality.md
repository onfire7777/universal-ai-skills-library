---
name: 'step-04-assertion-quality'
description: 'Verify assertions are meaningful and specific'
nextStepFile: './step-05-test-isolation.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/best-practices.md'
---

# Step 4: Review Assertion Quality

## STEP GOAL

Verify assertions are meaningful and specific, not just checking existence or types.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — assertion rules
- `references/common/best-practices.md` — assertion best practices

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### Quality Checks
- [ ] Uses `toMatchObject` for partial matching
- [ ] Checks specific values, not just types
- [ ] Verifies database state changes
- [ ] Verifies async side effects
- [ ] Error cases verify NO side effects

### Common Issues to Flag

```typescript
// ❌ Type-only assertion
expect(typeof response.body.id).toBe('string');
// ✅ Value assertion
expect(response.body.id).toMatch(/^user-/);

// ❌ Existence-only assertion
expect(response.body.data).toBeDefined();
// ✅ Content assertion
expect(response.body.data).toMatchObject({
  email: 'test@example.com',
  status: 'active',
});

// ❌ Missing side effect verification in error case
it('should reject invalid input', async () => {
  await request(httpServer).post('/users').send(invalid).expect(400);
  // Missing: verify no database record created
});
// ✅ With side effect check
it('should reject invalid input', async () => {
  await request(httpServer).post('/users').send(invalid).expect(400);
  const count = await userRepository.count();
  expect(count).toBe(0);
});
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 4: Assertion Quality
=========================

[File: {filename}]

Weak Assertions Found: {n}
- Line {x}: toBeDefined → should use toMatchObject
- Line {y}: Missing database verification
- Line {z}: Error case missing side effect check

Recommendations:
{specific improvements}

Summary: {n} assertions reviewed, {m} quality issues found
```

Then ask: **[C] Continue to Step 5: Test Isolation / [F] Fix issues first**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-test-isolation.md`.
