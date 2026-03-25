---
name: 'step-01-init'
description: 'Initialize writing workflow — identify target component and detect continuation'
nextStepFile: './step-02-test-file-structure.md'
---

# Step 1: Initialize Unit Test Writing

## STEP GOAL

Set up the unit test writing session: identify the target component, set the output path for the writing report, and check for an existing report to resume.

## EXECUTION

### 1. Ask the User

Ask the user:
- **Source file to test** — the path to the component to write tests for
- **Output path** for the writing report (suggest a default: `./unit-test-writing-report-{{date}}.md`)
- **Or provide path to an existing report** to resume a previous session

### 2. Check for Existing Report

If the user provides a path to an existing report file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty — **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `targetComponent`: the source file path provided by the user
   - `componentType`: detected component type (Service | Controller | Guard | Interceptor | Pipe | Filter | Repository)
   - `specFilePath`: the spec file path (co-located with source, e.g., `src/path/to/component.spec.ts`)
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized report to the output path

### 4. Analyze Target Component

Read the source file completely to understand:

1. **All public methods** — list every method that needs testing
2. **Constructor dependencies** — identify all injected services/repositories
3. **Exception handling** — find all thrown exceptions (NotFoundException, BadRequestException, etc.)
4. **Business logic branches** — identify conditional logic, calculations, validations

### 5. List Dependencies for Mocking

List all dependencies that need mocking:

```
**Dependencies to Mock:**
- [DependencyName]: [DependencyType]
- [DependencyName]: [DependencyType]
```

### 6. Append Component Analysis to Report

Append to the output document:

```markdown
## Step 1: Component Analysis

**Source File**: {{targetComponent}}
**Component Type**: {{componentType}}
**Spec File**: {{specFilePath}}

**Public Methods:**
- [methodName]: [brief description]

**Dependencies to Mock:**
- [DependencyName]: [DependencyType]

**Exception Handling:**
- [ExceptionType] thrown when [condition]

**Business Logic:**
- [description of business rules/branches]
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`
- Fill in `targetComponent`, `componentType`, `specFilePath`

## PRESENT TO USER

Show the user:
- Component summary (type, public methods, dependencies)
- List of all dependencies needing mocks
- Detected exception handling and business logic

Then ask: **[C] Continue to Step 2: Create Test File Structure**

## NEXT STEP

After user confirms `[C]`, load `step-02-test-file-structure.md`.
