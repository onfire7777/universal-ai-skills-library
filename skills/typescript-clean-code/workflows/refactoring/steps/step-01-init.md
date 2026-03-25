---
name: 'step-01-init'
description: 'Initialize refactoring workflow — verify tests, set output path, detect continuation'
nextStepFile: './step-02-identify-smell.md'
---

# Step 1: Initialize Refactoring

## STEP GOAL

Set up the refactoring session: identify the target code, verify tests pass, set the output path for the refactoring log, and check for an existing log to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **What code to refactor?** (file path, module, or class)
- **What's the improvement goal?** (what smell to fix, what to clean up)
- **Output path** for the refactoring log (suggest a default: `./refactoring-log-{{date}}.md`)
- **Or provide path to an existing log** to resume a previous session

### 2. Check for Existing Log

If the user provides a path to an existing log file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/log-template.md`
2. Fill in the frontmatter:
   - `targetCode`: the code path provided by the user
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized log to the output path

### 4. Verify Tests Pass

**Goal**: Establish a green baseline.

```bash
npm test  # All tests must pass
```

- All existing tests pass
- Coverage is adequate for the code you'll change
- You understand what the tests verify

**If tests fail**: Fix them first. Don't refactor broken code.

**If no tests exist**: Write characterization tests first. See `test-strategy.md` workflow.

### 5. Append to Log

Append to the output document:

```markdown
## Step 1: Initialization

**Target**: {{targetCode}}
**Goal**: {{improvement goal}}
**Tests**: {{PASS / FAIL — details}}
**Coverage**: {{adequate / needs improvement}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Set `testsGreen` to `true` (if tests passed)

## PRESENT TO USER

Show the user:
- Confirmation of the refactoring target and goal
- Test results
- Output path

Then ask: **[C] Continue to Step 2: Identify the Smell**

## NEXT STEP

After user confirms `[C]`, load `step-02-identify-smell.md`.
