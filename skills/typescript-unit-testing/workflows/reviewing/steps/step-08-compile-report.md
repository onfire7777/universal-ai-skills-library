---
name: 'step-08-compile-report'
description: 'Compile final review report'
---

# Step 8: Compile Review Report

## STEP GOAL

Compile all findings from steps 2-7 into a structured, prioritized review report with actionable recommendations.

## COMPILATION PROCESS

### 1. Gather All Findings

Read through the output document and collect all issues found in steps 2-7.

### 2. Build Summary Table

```markdown
## Unit Test Review Report

### File: [path/to/file.spec.ts]
### Source: [path/to/source.ts]
### Date: [current date]

---

### Summary

| Metric | Status |
|--------|--------|
| Overall Quality | [PASS/NEEDS WORK/FAIL] |
| Coverage | [X]% |
| Critical Issues | [count] |
| Major Issues | [count] |
| Minor Issues | [count] |
```

### 3. Categorize Issues by Severity

Use these severity definitions:

| Severity | Definition | Examples |
|----------|------------|----------|
| CRITICAL | Test provides false confidence or is fundamentally broken | Only checks `toBeDefined()`, testing wrong behavior |
| MAJOR | Test has significant gaps or doesn't follow conventions | Missing mock verification, no exception testing |
| MINOR | Test works but could be improved | Missing comments, inconsistent naming |

### 4. Critical Issues (Must Fix)

```markdown
### Critical Issues (Must Fix)

1. **[Issue Title]**
   - Location: [test name / line number]
   - Problem: [description]
   - Fix: [specific action to take]
```

### 5. Major Issues (Should Fix)

```markdown
### Major Issues (Should Fix)

1. **[Issue Title]**
   - Location: [test name / line number]
   - Problem: [description]
   - Fix: [specific action to take]
```

### 6. Minor Issues (Consider Fixing)

```markdown
### Minor Issues (Consider Fixing)

1. **[Issue Title]**
   - Location: [test name / line number]
   - Problem: [description]
   - Fix: [specific action to take]
```

### 7. Missing Test Cases

```markdown
### Missing Test Cases

| Method | Missing Category | Priority |
|--------|-----------------|----------|
| [method] | [category] | [HIGH/MEDIUM/LOW] |
```

### 8. Recommendations

```markdown
### Recommendations

1. [Specific recommendation]
2. [Specific recommendation]
```

Prioritize recommendations by impact:
1. Fix critical issues first (false confidence, broken tests)
2. Add missing exception tests
3. Improve assertion quality
4. Fill coverage gaps
5. Fix structural/naming conventions

### 9. Quick Review Checklist

Append the quick review checklist to the report:

```
**Review Checklist:**

Structure:
- [ ] File named *.spec.ts
- [ ] Co-located with source
- [ ] Using `target` for SUT
- [ ] Using `mock` prefix for mocks
- [ ] Using `DeepMocked<T>` type

Setup:
- [ ] Fresh mocks in beforeEach
- [ ] `.setLogger(new MockLoggerService())`
- [ ] `jest.clearAllMocks()` in afterEach

AAA Pattern:
- [ ] Arrange-Act-Assert structure
- [ ] Comments for each phase
- [ ] Single action per test

Assertions:
- [ ] Specific value assertions
- [ ] Mock call verification
- [ ] Mock call count verification
- [ ] No conditional assertions

Coverage:
- [ ] Happy path tests
- [ ] Edge case tests
- [ ] Error case tests
- [ ] Exception tests with type/message/code
- [ ] 80%+ coverage achieved
```

## PRESENT FINDINGS

Present the compiled report to the user. Highlight critical issues first.

Offer next steps:
- **Help fix critical issues** — offer to start with the most impactful fixes
- **Re-read `references/common/rules.md`** to confirm all issues were identified
- **Prioritize fixes** with the user

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `8` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The unit test review workflow is complete. The full report is saved at the output path.
