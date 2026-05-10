---
name: 'writing'
description: 'Write high-quality E2E tests following GWT pattern with proper isolation, assertions, and cleanup'
firstStepFile: './steps/step-01-init.md'
templateFile: './templates/output-template.md'
---

# Writing E2E Test Workflow

Write high-quality E2E tests following GWT pattern with proper isolation, assertions, and cleanup for TypeScript/NestJS applications.

## When to Use

- Creating new E2E test cases with proper GWT pattern
- Testing API endpoints, workflows, or complete features end-to-end
- Testing with real databases, message brokers, or external services
- Testing Kafka consumers/producers, event-driven workflows

## Prerequisites

Have the source code to test accessible and understand what behavior needs to be tested.

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
| 2 | `step-02-requirements.md` | Understand test requirements |
| 3 | `step-03-code-analysis.md` | Analyze existing code to be tested |
| 4 | `step-04-test-structure.md` | Design test file structure |
| 5 | `step-05-test-setup.md` | Write test lifecycle hooks |
| 6 | `step-06-happy-path.md` | Write happy path test cases |
| 7 | `step-07-error-cases.md` | Write error case tests |
| 8 | `step-08-async-tests.md` | Write async tests (if applicable) |
| 9 | `step-09-final-review.md` | Final review against rules, mark complete |

### Rules

1. **Load one step at a time** - Read the step file, execute it, then load the next
2. **Update frontmatter after each step** - Add the step number to `stepsCompleted`
3. **Wait for user confirmation** - Present findings and wait for `[C]` before proceeding
4. **Load reference files** - Each step specifies which reference files to load before analysis
5. **Cite specific rules** - When writing tests, follow patterns from reference files

### Mandatory Execution Rules

- **ALWAYS load knowledge references before writing any test**
- **ALWAYS follow GWT (Given-When-Then) pattern**
- **ALWAYS ensure test isolation (no shared state between tests)**
- **NEVER write tests with conditional assertions**
- **NEVER mock databases or message brokers in E2E tests**

### Anti-Patterns to Avoid

1. Multiple WHEN actions in one test
2. Conditional assertions (`if/else` in THEN)
3. Hardcoded test IDs
4. Mocking databases
5. Generic assertions (`toBeDefined`, `toBeTruthy`)
6. Fixed waits for async operations
7. Shared state between tests
8. fromBeginning: true for Kafka consumers

## Begin

Load `steps/step-01-init.md` to start.
