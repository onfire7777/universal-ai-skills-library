---
name: 'step-01-init'
description: 'Initialize running workflow — set output path, load knowledge, detect continuation'
nextStepFile: './step-02-pre-run.md'
referenceFiles:
  - 'references/common/knowledge.md'
  - 'references/common/debugging.md'
  - 'references/common/nestjs-setup.md'
---

# Step 1: Initialize E2E Test Run

## STEP GOAL

Set up the test run session: set the output path for the run log, load knowledge base, and check for existing output to resume.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/knowledge.md` - Understand test lifecycle
- `references/common/debugging.md` - Command line debugging and log analysis
- `references/common/nestjs-setup.md` - Configuration verification

## EXECUTION

### 1. Ask the User

Ask the user:
- **What tests to run?** (all, specific file, specific test, pattern)
- **Output path** for the run log (suggest a default: `./e2e-run-log-{{date}}.md`)
- **Or provide path to an existing run log** to resume

### 2. Check for Existing Output

If the user provides a path to an existing run log:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `testScope`: the test scope
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized log to the output path

### 4. Initialize Session

```bash
export E2E_SESSION=$(date +%s)-$$
```

### 5. Append to Output

Append to the output document:

```markdown
## Step 1: Session Setup

**Test Scope**: {{testScope}}
**Session ID**: {{E2E_SESSION}}
**References Loaded**: {{list of loaded references}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show confirmation of test scope and session setup.

Then ask: **[C] Continue to Step 2: Pre-Run Verification**

## NEXT STEP

After user confirms `[C]`, load `step-02-pre-run.md`.
