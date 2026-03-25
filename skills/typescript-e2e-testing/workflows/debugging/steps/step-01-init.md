---
name: 'step-01-init'
description: 'Initialize debugging workflow — set output path, load knowledge, detect continuation'
nextStepFile: './step-02-capture-failure.md'
referenceFiles:
  - 'references/common/debugging.md'
  - 'references/common/best-practices.md'
  - 'references/common/rules.md'
---

# Step 1: Initialize E2E Debugging

## STEP GOAL

Set up the debugging session: identify failing tests, set the output path for the tracking document, load knowledge base, and check for existing output to resume.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/debugging.md` - VS Code config, log analysis, systematic resolution
- `references/common/best-practices.md` - Common issues and solutions
- `references/common/rules.md` - Verify test follows required patterns

## EXECUTION

### 1. Ask the User

Ask the user:
- **Which tests are failing?** (test names, file, or "run suite to find out")
- **Output path** for the tracking document (suggest a default: `./e2e-debug-tracker-{{date}}.md`)
- **Or provide path to an existing tracker** to resume

### 2. Check for Existing Output

If the user provides a path to an existing tracker:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/tracking-template.md`
2. Fill in the frontmatter:
   - `failingTests`: list of failing test names
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized tracker to the output path

### 4. Initialize Session

```bash
export E2E_SESSION=$(date +%s)-$$
```

### 5. Load Technology-Specific References

Based on failing tests:
- **Kafka tests**: `references/kafka/isolation.md`, `references/kafka/performance.md`
- **Database tests**: `references/{postgres|mongodb}/rules.md`
- **API tests**: `references/api/mocking.md`

### 6. Append to Output

Append to the output document:

```markdown
## Step 1: Session Setup

**Failing Tests**: {{list of failing tests}}
**Session ID**: {{E2E_SESSION}}
**References Loaded**: {{list of loaded references}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Set `currentTest` to first failing test

## PRESENT TO USER

Show confirmation of failing tests and session setup.

Then ask: **[C] Continue to Step 2: Capture Failure**

## NEXT STEP

After user confirms `[C]`, load `step-02-capture-failure.md`.
