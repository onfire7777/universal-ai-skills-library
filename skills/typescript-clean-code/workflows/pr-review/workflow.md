---
name: 'pr-review'
description: 'Comprehensive process for reviewing pull requests'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/report-template.md'
---

# PR Review Workflow

Comprehensive process for reviewing pull requests.

## When to Use

- Reviewing a teammate's PR
- Doing a final check before merge
- Conducting formal code review

## Prerequisites

- Access to the PR and codebase
- Understanding of the feature/fix context
- Time to do a thorough review

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
| 2 | `step-02-tests-first.md` | Review tests first |
| 3 | `step-03-high-level.md` | Review code at high level |
| 4 | `step-04-code-details.md` | Review code details |
| 5 | `step-05-security.md` | Check for security issues |
| 6 | `step-06-performance.md` | Check for performance issues |
| 7 | `step-07-run-code.md` | Run code if needed |
| 8 | `step-08-feedback.md` | Provide constructive feedback |
| 9 | `step-09-decision.md` | Make approve/reject decision, mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before analysis
5. **Cite specific rules** - When reporting findings, cite the specific rule from the reference file

## Begin

Load `steps/step-01-init.md` to start.
