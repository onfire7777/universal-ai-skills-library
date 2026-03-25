---
name: 'step-03-gwt-compliance'
description: 'Verify all tests follow Given-When-Then pattern correctly'
nextStepFile: './step-04-assertion-quality.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/test-case-creation-guide.md'
---

# Step 3: Review GWT Pattern Compliance

## STEP GOAL

Verify all tests follow the Given-When-Then pattern correctly with proper comments, single actions, and specific assertions.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — GWT pattern rules
- `references/common/test-case-creation-guide.md` — GWT templates for all scenarios

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

For each test case, check:

### 1. GIVEN Section
- [ ] Has explicit GIVEN comment
- [ ] Setup is clear and complete
- [ ] Uses unique identifiers
- [ ] Test data created via factories

### 2. WHEN Section
- [ ] Has explicit WHEN comment
- [ ] EXACTLY ONE action
- [ ] Action clearly matches test name

### 3. THEN Section
- [ ] Has explicit THEN comment
- [ ] Specific assertions (not generic)
- [ ] No conditional logic
- [ ] Verifies all expected outcomes

### Pattern Violations to Flag

```typescript
// ❌ Multiple WHEN actions
it('should create and update user', async () => {
  const createRes = await request(httpServer).post('/users');
  const updateRes = await request(httpServer).patch(`/users/${id}`);
  // Split into separate tests!
});

// ❌ Generic assertion
expect(response.body.data).toBeDefined();
// Fix: expect(response.body.data.email).toBe('expected@email.com');

// ❌ Conditional assertion
if (response.status === 200) {
  expect(response.body.data).toBeDefined();
}
// Fix: Create separate tests for each scenario
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 3: GWT Pattern Review
===========================

[File: {filename}]

Test: "should {name}"
- GIVEN: ✅ Present / ❌ Missing
- WHEN: ✅ Single action / ❌ Multiple actions
- THEN: ✅ Specific / ❌ Generic

Issues:
{list issues with line numbers}

Recommendations:
{specific fixes}

Summary: {n} tests reviewed, {m} GWT issues found
```

Then ask: **[C] Continue to Step 4: Assertion Quality / [F] Fix issues first**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-assertion-quality.md`.
