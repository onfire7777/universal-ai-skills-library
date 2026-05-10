---
name: 'step-04-test-structure'
description: 'Design test file structure — describe blocks, test cases, setup/cleanup plan'
nextStepFile: './step-05-test-setup.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/examples.md'
---

# Step 4: Design Test Structure

## STEP GOAL

Plan the test file structure following patterns from references.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — structural rules
- `references/common/examples.md` — reference implementations

## EXECUTION

### 1. Design Structure

1. Design describe blocks based on feature organization
2. Plan test cases with descriptive names
3. Identify shared setup vs per-test setup
4. Determine cleanup requirements

### 2. Test File Template

```typescript
describe('{Feature} E2E', () => {
  let app: INestApplication;
  let httpServer: any;
  // Technology-specific helpers

  beforeAll(async () => {
    // App initialization
    // Helper initialization (pre-subscription for Kafka)
  }, 90000);

  afterAll(async () => {
    // Cleanup connections
    // Close app
  }, 30000);

  beforeEach(async () => {
    // Wait for in-flight operations
    // Clear database/message buffers
  });

  describe('{Sub-feature}', () => {
    it('should {expected behavior}', async () => {
      // GIVEN: {setup}
      // WHEN: {action}
      // THEN: {assertions}
    });
  });
});
```

### 3. Append to Output

Append to the output document:

```markdown
## Step 4: Test Structure

**Test File**: {{testFilePath}}

**Structure**:
{{show describe block hierarchy}}

**Test Cases**:
1. should {{case 1}} - Tests {{scenario}}
2. should {{case 2}} - Tests {{scenario}}
...
```

## PRESENT FINDINGS

Present the proposed structure to the user.

Then ask: **[C] Continue to Step 5: Test Setup / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-05-test-setup.md`.
