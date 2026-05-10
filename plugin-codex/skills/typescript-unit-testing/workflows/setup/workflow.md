---
name: 'setup'
description: 'Initialize unit testing infrastructure for NestJS/TypeScript projects'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Unit Test Setup Workflow

Initialize unit testing infrastructure for a NestJS/TypeScript project. This workflow ensures consistent test setup across all projects by walking through dependency installation, Jest configuration, test helper creation, and verification.

## When to Use

- Setting up unit testing in a new NestJS/TypeScript project
- Verifying and fixing an existing test infrastructure
- Standardizing test setup across multiple projects

## Prerequisites

Have access to the target NestJS/TypeScript project with a valid `package.json`.

## Step-File Architecture

This workflow uses a **step-file architecture** for context-safe execution:

- Each step is a separate file loaded sequentially
- Progress is tracked via `stepsCompleted` in the output document's YAML frontmatter
- If context is compacted mid-workflow, step-01 detects existing output and resumes from the last completed step via step-01b

### Steps

| Step | File | Description |
|------|------|-------------|
| 1 | `step-01-init.md` | Initialize workflow, analyze existing infrastructure, detect continuation |
| 1b | `step-01b-continue.md` | Resume from last completed step |
| 2 | `step-02-dependencies.md` | Install missing testing dependencies |
| 3 | `step-03-configure-jest.md` | Configure Jest for the project |
| 4 | `step-04-test-helpers.md` | Create test helper utilities |
| 5 | `step-05-verify-setup.md` | Verify setup with sample test |
| 6 | `step-06-npm-scripts.md` | Configure npm test scripts |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before execution

## Begin

Load `steps/step-01-init.md` to start.
