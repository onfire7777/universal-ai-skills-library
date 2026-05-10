---
name: 'step-02-test-structure'
description: 'Review test file organization and lifecycle hooks'
nextStepFile: './step-03-gwt-compliance.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/nestjs-setup.md'
---

# Step 2: Review Test Structure

## STEP GOAL

Verify test file organization, lifecycle hooks (beforeAll/beforeEach/afterAll), and proper configuration.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — mandatory patterns
- `references/common/nestjs-setup.md` — NestJS test setup patterns

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

For each test file, check:

### Checklist
- [ ] Test file uses `.e2e-spec.ts` suffix
- [ ] Proper describe block organization
- [ ] beforeAll: App initialization, helper setup
- [ ] beforeEach: Cleanup and isolation
- [ ] afterAll: Proper teardown
- [ ] Appropriate timeouts set

### Issues to Flag

```typescript
// ❌ Missing sequential execution
// Fix: maxWorkers: 1 in jest config

// ❌ Missing cleanup in beforeEach
beforeEach(async () => {
  // Must clean database
  await repository.clear();
});

// ❌ Missing timeout on long operations
beforeAll(async () => {
  // Add timeout
}, 90000);

// ❌ Not closing app properly
afterAll(async () => {
  await app?.close(); // Must close
}, 30000);
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 2: Structure Review
========================

[File: {filename}]

✅ Test file naming correct
❌ Missing beforeEach cleanup
⚠️ No explicit timeout on beforeAll

Issues Found: {n}

{detailed findings with line numbers}

Summary: {n} files reviewed, {m} structural issues found
```

Then ask: **[C] Continue to Step 3: GWT Compliance / [F] Fix issues first**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-gwt-compliance.md`.
