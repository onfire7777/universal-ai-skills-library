---
name: 'step-01-init'
description: 'Initialize optimization workflow — set output path, load knowledge, detect continuation'
nextStepFile: './step-02-baseline.md'
referenceFiles:
  - 'references/common/best-practices.md'
  - 'references/common/knowledge.md'
  - 'references/common/rules.md'
---

# Step 1: Initialize E2E Optimization

## STEP GOAL

Set up the optimization session: identify the test suite, set the output path, load knowledge base, and check for existing output to resume.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/best-practices.md` - Performance patterns
- `references/common/knowledge.md` - Test lifecycle understanding
- `references/common/rules.md` - Ensure optimizations don't violate rules

## EXECUTION

### 1. Ask the User

Ask the user:
- **Which test suite to optimize?** (all E2E tests, specific files)
- **Output path** for the optimization report (suggest a default: `./e2e-optimization-report-{{date}}.md`)
- **Or provide path to an existing report** to resume

### 2. Check for Existing Output

If the user provides a path to an existing report:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/report-template.md`
2. Fill in the frontmatter:
   - `testSuite`: the test suite scope
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Initialize Session

```bash
export E2E_SESSION=$(date +%s)-$$
```

### 5. Load Technology-Specific References

- **Kafka**: `references/kafka/performance.md` - Kafka-specific optimizations
- **Database**: `references/{postgres|mongodb}/rules.md` - Query optimization
- **API**: `references/api/rules.md` - Request optimization

### 6. Append to Output

Append to the output document:

```markdown
## Step 1: Session Setup

**Test Suite**: {{testSuite}}
**Session ID**: {{E2E_SESSION}}
**References Loaded**: {{list of loaded references}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show confirmation of test suite and session setup.

Then ask: **[C] Continue to Step 2: Baseline Assessment**

## NEXT STEP

After user confirms `[C]`, load `step-02-baseline.md`.
