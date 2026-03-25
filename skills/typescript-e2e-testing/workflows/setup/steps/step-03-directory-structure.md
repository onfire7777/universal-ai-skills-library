---
name: 'step-03-directory-structure'
description: 'Create E2E test directory structure'
nextStepFile: './step-04-jest-config.md'
referenceFiles:
  - 'references/common/nestjs-setup.md'
---

# Step 3: Create Directory Structure

## STEP GOAL

Set up the E2E test directory structure.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/nestjs-setup.md` — project structure patterns

## EXECUTION

### 1. Plan Structure

```
test/
├── e2e/
│   ├── setup.ts           # Global test setup
│   └── helpers/           # Test utilities
│       ├── test-app.helper.ts
│       └── {technology}.helper.ts
└── jest-e2e.config.ts     # E2E Jest configuration
```

### 2. Present to User for Confirmation

Show the planned structure before creating.

### 3. Create Directories

After user confirms, create the directory structure.

### 4. Append to Output

Append to the output document:

```markdown
## Step 3: Directory Structure

**Created**:
{{show tree structure}}
```

## PRESENT FINDINGS

```
I'll create this directory structure:

{show tree structure}

[C] Continue / [S] Skip (keep existing) / [M] Modify
```

Then ask: **[C] Continue to Step 4: Jest Config**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-04-jest-config.md`.
