---
name: 'step-07-summary'
description: 'Generate comprehensive review summary with prioritized fixes'
referenceFiles:
  - 'references/common/rules.md'
---

# Step 7: Generate Review Summary

## STEP GOAL

Provide comprehensive review summary with prioritized fixes, compiling all findings from steps 2-6.

## REFERENCE LOADING

Before compiling summary, load and read:
- `references/common/rules.md` — to verify all findings are properly categorized

## COMPILATION PROCESS

### 1. Gather All Findings

Read through the output document and collect all issues found in steps 2-6.

### 2. Categorize by Severity

Categorize each issue:

| Severity | Meaning |
|----------|---------|
| **Critical (Must Fix)** | Tests with multiple WHEN actions, conditional assertions, missing cleanup, hardcoded Kafka IDs, fromBeginning: true |
| **Important (Should Fix)** | Generic assertions, missing side effect verification, hardcoded test IDs, missing GWT comments, inappropriate timeouts |
| **Minor (Nice to Fix)** | Suboptimal organization, missing factory functions, verbose setup, missing error case tests |

### 3. Write Summary

Append to the output document:

```markdown
## Review Summary

### Files Reviewed
- {file1}: {n} tests
- {file2}: {m} tests

### Critical Issues (Must Fix)
1. {issue} - {file}:{line}
2. {issue} - {file}:{line}

### Important Issues (Should Fix)
1. {issue} - {file}:{line}

### Minor Issues (Nice to Fix)
1. {issue} - {file}:{line}

### Good Practices Found
- {positive finding}
- {positive finding}

### Recommended Reading
- references/common/{relevant}.md
- references/{technology}/{relevant}.md

### Next Steps
1. Fix critical issues first
2. Run tests to verify fixes
3. Address important issues
4. Consider refactoring for minor issues
```

## PRESENT TO USER

Show the final summary. Then offer:

```
Would you like me to:
[F] Fix critical issues now
[D] Generate detailed fix instructions
[E] Export review as markdown file
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `7` to `stepsCompleted`
- Set `status` to `'complete'`

### Post-Workflow Recommendations

After review, recommend:
1. Re-read `references/common/rules.md` if many GWT violations
2. Re-read technology-specific guides if patterns incorrect
3. Review `references/common/examples.md` for reference implementations

## WORKFLOW COMPLETE

The review workflow is complete. The full report is saved at the output path.
