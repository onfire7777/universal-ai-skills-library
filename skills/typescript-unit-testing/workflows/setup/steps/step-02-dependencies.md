---
name: 'step-02-dependencies'
description: 'Install missing testing dependencies'
nextStepFile: './step-03-configure-jest.md'
referenceFiles:
  - 'references/common/knowledge.md'
  - 'references/common/rules.md'
---

# Step 2: Install Missing Dependencies

## STEP GOAL

Identify all missing testing dependencies and install them so the project has everything needed for NestJS unit testing.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/knowledge.md` — testing philosophy and pyramid
- `references/common/rules.md` — naming conventions and setup requirements

## EXECUTION

### 1. Check Required Dependencies

Review the project's `package.json` and verify the following devDependencies are present:

```json
{
  "devDependencies": {
    "@nestjs/testing": "^11.0.12",
    "jest": "^29.7.0",
    "ts-jest": "^29.2.6",
    "@golevelup/ts-jest": "^0.4.0",
    "@types/jest": "^29.5.14"
  }
}
```

### 2. Identify Missing Dependencies

Compare installed dependencies against the required list. Note:
- Which dependencies are already installed and at what version
- Which dependencies are missing entirely
- Which dependencies have incompatible versions

### 3. Present Installation Command

If there are missing dependencies, present the install command:

```bash
npm install --save-dev [missing-packages]
```

Wait for user confirmation before running the installation.

### 4. Verify Installation

After installation completes:
- Confirm all required dependencies appear in `package.json`
- Verify no installation errors occurred

## PRESENT FINDINGS

Present findings to the user:

```
Step 2: Dependencies
====================

Required Dependencies:
  [@nestjs/testing]  : INSTALLED (version) / MISSING
  [jest]             : INSTALLED (version) / MISSING
  [ts-jest]          : INSTALLED (version) / MISSING
  [@golevelup/ts-jest]: INSTALLED (version) / MISSING
  [@types/jest]      : INSTALLED (version) / MISSING

Action Taken: [Installed X missing packages / All dependencies already present]
```

Then ask: **[C] Continue to Step 3: Configure Jest**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Append the dependency status section to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-configure-jest.md`.
