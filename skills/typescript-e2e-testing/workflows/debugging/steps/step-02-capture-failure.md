---
name: 'step-02-capture-failure'
description: 'Document the failure completely before investigating'
nextStepFile: './step-03-classify.md'
---

# Step 2: Capture Failure Information

## STEP GOAL

Document the failure completely before investigating.

## EXECUTION

### 1. Create or Update Tracking Entry

```markdown
## Test: "{test name}"
- **File**: {file}:{line}
- **Error Type**: {Timeout/Assertion/Connection/etc}
- **Error Message**:
  ```
  {full error message}
  ```
- **Stack Trace**:
  ```
  {relevant stack trace}
  ```
- **Status**: INVESTIGATING
- **Related Tests**: {list if part of a pattern}
```

### 2. Gather Context

- When did it start failing?
- Does it fail consistently or intermittently?
- Does it fail alone or only with other tests?
- What changed recently?

### 3. Append to Output

Append failure details to the output document.

## PRESENT FINDINGS

```
Failure Captured

Test: "{test name}"
File: {file}:{line}

Error Type: {type}
Error: {message}

Questions:
1. Is this a new failure or existing?
2. Does it fail consistently?
3. Did anything change recently (code, deps, config)?

[C] Continue investigation / [P] Provide more context
```

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Update `currentTest`

## NEXT STEP

After user confirms `[C]`, load `step-03-classify.md`.
