---
name: 'step-03-configure-jest'
description: 'Configure Jest for the project'
nextStepFile: './step-04-test-helpers.md'
---

# Step 3: Configure Jest

## STEP GOAL

Create or verify the Jest configuration file (`jest.config.ts`) so the project has a correct and complete test runner setup.

## EXECUTION

### 1. Check for Existing Configuration

Look for `jest.config.ts` or `jest.config.js` in the project root.

### 2a. If No Config Exists — Create It

Create `jest.config.ts` with the following standard configuration:

```typescript
import type { Config } from 'jest';

const config: Config = {
  moduleFileExtensions: ['js', 'json', 'ts'],
  rootDir: '.',
  testRegex: '.*\\.spec\\.ts$',
  transform: {
    '^.+\\.(t|j)s$': 'ts-jest',
  },
  collectCoverageFrom: ['**/*.(t|j)s'],
  coverageDirectory: './coverage',
  testEnvironment: 'node',
  roots: ['<rootDir>/src/', '<rootDir>/test/'],
  moduleNameMapper: {
    '^src/(.*)$': '<rootDir>/src/$1',
  },
};

export default config;
```

### 2b. If Config Exists — Verify It

Check that the existing configuration includes:
- Correct `testRegex` pattern for `.spec.ts` files
- `ts-jest` transform for TypeScript compilation
- Coverage configuration (`collectCoverageFrom`, `coverageDirectory`)
- `roots` including both `src/` and `test/` directories
- Path alias mapping via `moduleNameMapper` (if project uses path aliases in `tsconfig.json`)

Report any missing or incorrect settings and propose corrections.

### 3. Append to Report

Append to the output document:

```markdown
## Step 3: Jest Configuration

**Config File**: [created/verified] at [path]
**Settings Verified:**
- testRegex: [OK/Fixed]
- ts-jest transform: [OK/Fixed]
- Coverage config: [OK/Fixed]
- Module name mapper: [OK/Fixed]
- Roots: [OK/Fixed]

**Changes Made**: [list changes or "No changes needed"]
```

## PRESENT FINDINGS

Present the Jest configuration status:

```
Step 3: Jest Configuration
==========================

Config: [Created new / Verified existing] jest.config.ts
  - testRegex:         OK
  - ts-jest transform: OK
  - Coverage:          OK
  - Module mapper:     OK
  - Roots:             OK

Changes: [list or "None needed"]
```

Then ask: **[C] Continue to Step 4: Create Test Helpers**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the Jest configuration section to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-test-helpers.md`.
