---
name: 'step-05-test-setup'
description: 'Write test lifecycle hooks — beforeAll, beforeEach, afterAll'
nextStepFile: './step-06-happy-path.md'
referenceFiles:
  - 'references/common/nestjs-setup.md'
---

# Step 5: Write Test Setup

## STEP GOAL

Write proper test lifecycle hooks (beforeAll/beforeEach/afterAll).

## REFERENCE LOADING

Before starting, load and read:
- `references/common/nestjs-setup.md` — NestJS test app bootstrap

Also load technology-specific test-helper files as needed:
- `references/kafka/test-helper.md`
- `references/postgres/test-helper.md`
- `references/mongodb/test-helper.md`
- `references/redis/test-helper.md`
- `references/api/test-helper.md`

## EXECUTION

### 1. Write beforeAll

- Create NestJS test module
- Initialize app with proper config
- Set up technology-specific helpers
- For Kafka: Pre-subscribe to output topics

### 2. Write beforeEach

- Wait for in-flight operations (500-1000ms)
- Clear database tables/collections
- Clear message buffers

### 3. Write afterAll

- Disconnect helpers
- Close app

### 4. Append to Output

Append to the output document:

```markdown
## Step 5: Test Setup

### beforeAll
{{show beforeAll code}}

### beforeEach
{{show beforeEach code}}

### afterAll
{{show afterAll code}}

**Patterns from**: references/common/nestjs-setup.md, references/{tech}/test-helper.md
```

## PRESENT FINDINGS

Present the setup code to the user.

Then ask: **[C] Continue to Step 6: Happy Path Tests / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `5` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-06-happy-path.md`.
