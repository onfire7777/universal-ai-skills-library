---
name: 'step-02-test-structure'
description: 'Review test file structure and naming conventions'
nextStepFile: './step-03-setup-teardown.md'
referenceFiles:
  - 'references/common/rules.md'
---

# Step 2: Review Test Structure

## STEP GOAL

Check test file structure against conventions — file naming, co-location, variable naming, mock patterns, and describe/it hierarchy.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

### 1. Check Test File Structure Against Conventions

| Criterion | Expected | Status |
|-----------|----------|--------|
| File naming | `*.spec.ts` | [ ] |
| File location | Co-located with source | [ ] |
| SUT variable | Named `target` | [ ] |
| Mock variables | Prefixed with `mock` | [ ] |
| Mock types | Using `DeepMocked<T>` | [ ] |
| Logger | Using `MockLoggerService` | [ ] |

### 2. Check describe/it Structure

Verify the test file follows the nested describe/it convention:

```typescript
describe('ClassName', () => {
  describe('methodName', () => {
    it('should [expected] when [condition]', () => {});
  });
});
```

Check for:
- Top-level `describe` matches the class/module name
- Second-level `describe` blocks for each method
- `it` blocks follow the `should [expected] when [condition]` pattern
- No deeply nested describes beyond method level (unless grouping edge cases)

### 3. Document Structural Issues

For each issue found, classify severity:

```
**Structural Issues:**
- [CRITICAL/MAJOR/MINOR]: [Description]
```

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 2: Test Structure
======================

Convention Checklist:
  File naming (*.spec.ts):        [PASS/FAIL]
  Co-located with source:         [PASS/FAIL]
  SUT named `target`:             [PASS/FAIL]
  Mock prefix `mock`:             [PASS/FAIL]
  DeepMocked<T> types:            [PASS/FAIL]
  MockLoggerService:              [PASS/FAIL]

Describe/It Structure:
  Top-level matches class name:   [PASS/FAIL]
  Method-level describes:         [PASS/FAIL]
  It naming convention:           [PASS/FAIL]

Issues Found: N
  - [SEVERITY]: description
```

Then ask: **[C] Continue to Step 3: Setup & Teardown**

## FRONTMATTER UPDATE

Update the output document:
- Add `2` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-03-setup-teardown.md`.
