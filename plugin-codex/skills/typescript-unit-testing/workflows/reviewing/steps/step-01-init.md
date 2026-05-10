---
name: 'step-01-init'
description: 'Initialize review workflow and detect continuation'
nextStepFile: './step-02-test-structure.md'
---

# Step 1: Initialize Unit Test Review

## STEP GOAL

Set up the review session: identify test files to review, set the output path for the review report, and check for an existing report to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **What test file(s) to review?** (file path or directory, e.g., `path/to/component.spec.ts` or `path/to/module/`)
- **Output path** for the review report (suggest a default: `./unit-test-review-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous review

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `targetTests`: the test file path(s) provided by the user
   - `targetSource`: the corresponding source file path(s)
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Understand Context

**Goal**: Understand what the code is supposed to do before judging how it is tested.

- Read the test file(s) to understand what is currently being tested
- Read the corresponding source file(s) to understand what **should** be tested
- Identify the class/module under test, its public methods, and its dependencies

List files to review:
```
**Files to Review:**
- Test: [path/to/file.spec.ts]
- Source: [path/to/file.ts]
```

### 5. Append Initial Context to Report

Append to the output document:

```markdown
## Step 1: Context

**Test File**: {{targetTests}}
**Source File**: {{targetSource}}
**Class/Module Under Test**: {{className}}
**Public Methods**: {{list of public methods}}
**Dependencies**: {{list of dependencies}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Fill `targetTests` and `targetSource`

## PRESENT TO USER

Show the user:
- Files identified for review
- Source code understanding (class, public methods, dependencies)
- Summary of what the test file currently covers

Then ask: **[C] Continue to Step 2: Test Structure**

## NEXT STEP

After user confirms `[C]`, load `step-02-test-structure.md`.
