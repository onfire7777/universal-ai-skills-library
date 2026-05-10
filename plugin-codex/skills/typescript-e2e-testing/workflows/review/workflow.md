---
name: 'review'
description: 'Review existing E2E tests for quality, correctness, and adherence to best practices'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/report-template.md'
---

# Review E2E Test Workflow

Review existing E2E tests for quality, correctness, and adherence to best practices. Identify issues and provide actionable improvements.

## When to Use

- Reviewing existing E2E tests for quality
- Checking test isolation and cleanup patterns
- Auditing GWT pattern compliance
- Evaluating assertion quality and specificity
- Checking for anti-patterns

## Prerequisites

Have the E2E test file(s) to review accessible.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, set output path, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-test-structure.md` | Review test file organization and lifecycle hooks |
| 3 | `step-03-gwt-compliance.md` | Verify GWT pattern compliance |
| 4 | `step-04-assertion-quality.md` | Verify assertions are meaningful and specific |
| 5 | `step-05-test-isolation.md` | Verify tests don't share state |
| 6 | `step-06-tech-patterns.md` | Verify technology-specific patterns |
| 7 | `step-07-summary.md` | Generate review summary, mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before analysis
5. **Cite specific rules** - When reporting findings, cite the specific rule from the reference file

### Mandatory Execution Rules

- **ALWAYS load knowledge references before reviewing**
- **ALWAYS check against mandatory rules from references/common/rules.md**
- **ALWAYS provide specific, actionable feedback**
- **NEVER approve tests with anti-patterns**

### Review Severity Levels

#### Critical (Must Fix)
- Tests with multiple WHEN actions
- Conditional assertions
- Missing database cleanup
- Hardcoded IDs in Kafka consumer groups
- fromBeginning: true in Kafka tests

#### Important (Should Fix)
- Generic assertions (toBeDefined)
- Missing side effect verification
- Hardcoded test IDs
- Missing GWT comments
- Inappropriate timeouts

#### Minor (Nice to Fix)
- Suboptimal test organization
- Missing factory functions
- Verbose setup code
- Missing error case tests

## Begin

Load `steps/step-01-init.md` to start.
