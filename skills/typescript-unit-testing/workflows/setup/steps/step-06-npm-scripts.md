---
name: 'step-06-npm-scripts'
description: 'Configure npm test scripts'
---

# Step 6: Configure npm Scripts

## STEP GOAL

Verify and add the required npm test scripts to `package.json` so the project has consistent commands for running, watching, debugging, and generating coverage for tests.

## EXECUTION

### 1. Check Existing Scripts

Read the `scripts` section of `package.json` and check for the following required scripts:

```json
{
  "scripts": {
    "test": "jest",
    "test:watch": "jest --watch",
    "test:cov": "jest --coverage",
    "test:debug": "node --inspect-brk -r tsconfig-paths/register -r ts-node/register node_modules/.bin/jest --runInBand"
  }
}
```

### 2. Add Missing Scripts

For each missing script:
- Present the script to the user
- Add it to `package.json`

Do not overwrite existing scripts that already serve the same purpose (e.g., if `test` already runs Jest with custom flags, confirm with the user before replacing).

### 3. Verify Scripts Work

Run a quick verification:

```bash
npm test -- --help
```

Confirm the test command is wired up correctly.

### 4. Append to Report

Append to the output document:

```markdown
## Step 6: npm Scripts

**Scripts Verified:**
- `test`: [present/added] — `jest`
- `test:watch`: [present/added] — `jest --watch`
- `test:cov`: [present/added] — `jest --coverage`
- `test:debug`: [present/added] — `node --inspect-brk ...`

**Changes Made**: [list changes or "All scripts already present"]
```

## PRESENT FINDINGS

Present the npm scripts status and full completion summary:

```
Step 6: npm Scripts
===================

Scripts:
  test:        [present/added]
  test:watch:  [present/added]
  test:cov:    [present/added]
  test:debug:  [present/added]

Changes: [list or "None needed"]

========================================
Unit Test Setup Complete
========================================

  [1] Infrastructure Analysis    DONE
  [2] Dependencies Installed     DONE
  [3] Jest Configured            DONE
  [4] Test Helpers Created       DONE
  [5] Setup Verified             DONE
  [6] npm Scripts Configured     DONE

Next Steps:
- Use the `writing-unit-test.md` workflow to write tests
- Follow AAA pattern and naming conventions from references/common/rules.md
- Target 80%+ coverage for new code
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `6` to `stepsCompleted`
- Set `status` to `'complete'`
- Append the npm scripts section and completion summary to the report

## WORKFLOW COMPLETE

The unit test setup workflow is complete. The full report is saved at the output path.
