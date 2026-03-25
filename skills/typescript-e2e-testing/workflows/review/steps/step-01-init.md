---
name: 'step-01-init'
description: 'Initialize review workflow — set output path, load knowledge, detect continuation'
nextStepFile: './step-02-test-structure.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/best-practices.md'
  - 'references/common/examples.md'
---

# Step 1: Initialize E2E Test Review

## STEP GOAL

Set up the review session: identify the test files to review, set the output path for the review report, load knowledge base, and check for an existing report to resume.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` - Mandatory patterns to enforce
- `references/common/best-practices.md` - Quality standards
- `references/common/examples.md` - Reference implementations

## EXECUTION

### 1. Ask the User

Ask the user:
- **What tests to review?** (specific test file(s), or review all E2E tests)
- **Output path** for the review report (suggest a default: `./e2e-review-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous review

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/report-template.md`
2. Fill in the frontmatter:
   - `targetTests`: the test file(s)/scope provided by the user
   - `technologies`: detected technologies (Kafka, PostgreSQL, MongoDB, Redis, API)
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Load Technology-Specific References

Based on technologies detected in the test files:
- **Kafka tests**: `references/kafka/rules.md`, `references/kafka/isolation.md`
- **PostgreSQL tests**: `references/postgres/rules.md`
- **MongoDB tests**: `references/mongodb/rules.md`
- **Redis tests**: `references/redis/rules.md`
- **API tests**: `references/api/rules.md`

### 5. List Tests to Review

List all test files and test cases in scope:

```
Tests to Review:

1. {file1.e2e-spec.ts} - {number} test cases
2. {file2.e2e-spec.ts} - {number} test cases
...

Total: {n} test files, {m} test cases
```

### 6. Append to Report

Append to the output document:

```markdown
## Step 1: Review Scope

**Tests Under Review**:
- {{file list with test counts}}

**Technologies Detected**: {{technologies}}
**References Loaded**: {{list of loaded references}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show the user:
- Confirmation of tests to review and output path
- Technologies detected
- References loaded

Then ask: **[C] Continue to Step 2: Test Structure**

## NEXT STEP

After user confirms `[C]`, load `step-02-test-structure.md`.
