---
name: 'step-01-init'
description: 'Initialize unit test setup workflow — set output path and detect continuation'
nextStepFile: './step-02-dependencies.md'
---

# Step 1: Initialize Unit Test Setup

## STEP GOAL

Set up the unit test setup session: identify the target project, set the output path for the setup report, check for an existing report to resume, and perform initial infrastructure analysis.

## EXECUTION

### 1. Ask the User

Ask the user:
- **Target project path** — the root directory of the NestJS/TypeScript project
- **Output path** for the setup report (suggest a default: `./unit-test-setup-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous setup

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty — **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `targetProject`: the project path provided by the user
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Initial Infrastructure Analysis

Analyze the existing test infrastructure in the target project:

1. **Check for existing test configuration:**
   - Look for `jest.config.ts` or `jest.config.js` in project root
   - Check for existing `test/` directory structure
   - Look for `*.spec.ts` files to understand existing patterns

2. **Check package.json for testing dependencies:**
   - `@nestjs/testing`
   - `jest`
   - `ts-jest`
   - `@golevelup/ts-jest`
   - `@types/jest`

3. **Check for existing spec files:**
   - Count `*.spec.ts` files in the project
   - Note any existing test patterns

### 5. Append Findings to Report

Append to the output document:

```markdown
## Step 1: Infrastructure Analysis

**Test Infrastructure Analysis:**
- Jest config: [Found/Not found] at [path]
- Test directory: [Found/Not found] at [path]
- Existing spec files: [count] files found
- Dependencies status: [list installed and missing dependencies]

**Existing Patterns:**
- [describe any existing test patterns found, or "No existing tests found"]
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show the user:
- Confirmation of the target project and output path
- Infrastructure analysis summary
- List of installed vs missing dependencies

Then ask: **[C] Continue to Step 2: Install Missing Dependencies**

## NEXT STEP

After user confirms `[C]`, load `step-02-dependencies.md`.
