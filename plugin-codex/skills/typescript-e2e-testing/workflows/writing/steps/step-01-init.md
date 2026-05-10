---
name: 'step-01-init'
description: 'Initialize writing workflow — set output path, load knowledge, detect continuation'
nextStepFile: './step-02-requirements.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/test-case-creation-guide.md'
  - 'references/common/best-practices.md'
---

# Step 1: Initialize E2E Test Writing

## STEP GOAL

Set up the test writing session: identify the feature to test, set the output path, load knowledge base, and check for existing output to resume.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` - GWT pattern, timeout rules, isolation requirements
- `references/common/test-case-creation-guide.md` - Templates for all scenarios
- `references/common/best-practices.md` - Test design patterns

## EXECUTION

### 1. Ask the User

Ask the user:
- **What feature/flow to test?** (endpoint, workflow, event handler)
- **Output path** for the test plan document (suggest a default: `./e2e-test-plan-{{date}}.md`)
- **Or provide path to an existing plan** to resume

### 2. Check for Existing Output

If the user provides a path to an existing plan file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `targetFeature`: the feature/flow to test
   - `technologies`: detected technologies
   - `testFilePath`: planned test file path
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized plan to the output path

### 4. Load Technology-Specific References

Based on technologies involved:
- **Kafka tests**: `references/kafka/rules.md`, `references/kafka/isolation.md`, `references/kafka/examples.md`
- **PostgreSQL tests**: `references/postgres/rules.md`, `references/postgres/examples.md`
- **MongoDB tests**: `references/mongodb/rules.md`, `references/mongodb/examples.md`
- **Redis tests**: `references/redis/rules.md`, `references/redis/examples.md`
- **API tests**: `references/api/rules.md`, `references/api/examples.md`, `references/api/mocking.md`

### 5. Append to Output

Append to the output document:

```markdown
## Step 1: Session Setup

**Feature to Test**: {{targetFeature}}
**Technologies**: {{technologies}}
**Test File**: {{testFilePath}}
**References Loaded**: {{list of loaded references}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show the user:
- Confirmation of feature, technologies, and test file path
- References loaded

Then ask: **[C] Continue to Step 2: Requirements**

## NEXT STEP

After user confirms `[C]`, load `step-02-requirements.md`.
