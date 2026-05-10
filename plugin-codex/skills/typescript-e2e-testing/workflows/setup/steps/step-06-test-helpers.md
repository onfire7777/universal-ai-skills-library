---
name: 'step-06-test-helpers'
description: 'Create reusable test helper classes for each technology'
nextStepFile: './step-07-global-setup.md'
referenceFiles:
  - 'references/kafka/test-helper.md'
  - 'references/postgres/test-helper.md'
  - 'references/mongodb/test-helper.md'
  - 'references/redis/test-helper.md'
  - 'references/api/test-helper.md'
---

# Step 6: Create Test Helpers

## STEP GOAL

Create reusable test helper classes for each detected technology.

## REFERENCE LOADING

Before starting, load and read the test-helper references for detected technologies:
- `references/kafka/test-helper.md` — KafkaTestHelper implementation
- `references/postgres/test-helper.md` — PostgresTestHelper implementation
- `references/mongodb/test-helper.md` — MongoDbTestHelper implementation
- `references/redis/test-helper.md` — RedisTestHelper implementation
- `references/api/test-helper.md` — Auth and Supertest helpers

## EXECUTION

### 1. Create test-app.helper.ts

Create `test/e2e/helpers/test-app.helper.ts`:
- NestJS app bootstrap
- HTTP server access
- Cleanup utilities

### 2. Create Technology-Specific Helpers

For each detected technology, create `{technology}.helper.ts` based on the reference file.

### 3. Append to Output

Append to the output document:

```markdown
## Step 6: Test Helpers

### test-app.helper.ts
{{show code preview}}

### {technology}.helper.ts
{{show code preview}}
```

## PRESENT FINDINGS

Show the test helpers to the user.

Then ask: **[C] Continue to Step 7: Global Setup / [M] Modify**

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-07-global-setup.md`.
