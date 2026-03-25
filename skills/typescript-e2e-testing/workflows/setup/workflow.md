---
name: 'setup'
description: 'Set up a complete E2E testing infrastructure for a TypeScript/NestJS project'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Setup E2E Test Workflow

Set up a complete E2E testing infrastructure for a TypeScript/NestJS project with all necessary configuration, helpers, and Docker infrastructure.

## When to Use

- Setting up E2E testing for a new project
- Configuring docker-compose for testing (Kafka, PostgreSQL, MongoDB, Redis)
- Creating jest-e2e.config.ts or E2E Jest configuration
- Setting up test helpers for database, Kafka, or Redis
- Configuring .env.e2e environment variables
- Creating test/e2e directory structure

## Prerequisites

Have access to the project's source code and understand its infrastructure dependencies.

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
| 2 | `step-02-project-analysis.md` | Analyze project structure and dependencies |
| 3 | `step-03-directory-structure.md` | Create E2E test directory structure |
| 4 | `step-04-jest-config.md` | Configure Jest for E2E testing |
| 5 | `step-05-docker-compose.md` | Create Docker Compose for test infrastructure |
| 6 | `step-06-test-helpers.md` | Create reusable test helper classes |
| 7 | `step-07-global-setup.md` | Create global E2E test setup file |
| 8 | `step-08-example-test.md` | Create example E2E test demonstrating patterns |
| 9 | `step-09-verification.md` | Verify setup works correctly, mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before action
5. **Ask before creating files** - Always confirm with user before creating/modifying files

### Mandatory Execution Rules

- **NEVER proceed without understanding project context first**
- **ALWAYS load and review knowledge references before taking action**
- **ALWAYS ask user for confirmation before creating files**
- **ALWAYS verify Docker and dependencies are available**

### Failure Recovery

If any step fails:
1. Report the specific error
2. Reference relevant documentation from `references/`
3. Offer manual troubleshooting steps
4. Ask user if they want to retry or skip

## Begin

Load `steps/step-01-init.md` to start.
