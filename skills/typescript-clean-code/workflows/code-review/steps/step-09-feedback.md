---
name: 'step-09-feedback'
description: 'Compile final feedback — prioritize, summarize, mark complete'
referenceFiles:
  - 'references/collaboration/rules.md'
---

# Step 9: Provide Feedback

## STEP GOAL

Compile all findings from steps 2-8 into a prioritized, actionable feedback summary. Mark the review as complete.

## REFERENCE LOADING

Before compiling feedback, load and read:
- `references/collaboration/rules.md` — feedback and communication rules

## COMPILATION PROCESS

### 1. Gather All Findings

Read through the output document and collect all issues found in steps 2-8.

### 2. Prioritize Issues

Categorize each issue:

| Priority | Category | Meaning |
|----------|----------|---------|
| 1 | `[BLOCKING]` | Must fix — correctness, security, major design flaw |
| 2 | `[SUGGESTION]` | Should fix — improves quality significantly |
| 3 | `[NIT]` | Nice to have — minor style or preference |

### 3. Format Feedback

For each issue, use this template:

```
[BLOCKING/SUGGESTION/NIT] [File:Line]
Issue: What's wrong
Why: Reference to specific rule
Suggestion: How to fix
```

### 4. Write Summary

Append to the output document:

```markdown
## Final Review Summary

### Blocking Issues (Must Fix)
- [list or "None found"]

### Suggestions (Should Fix)
- [list or "None"]

### Nits (Nice to Have)
- [list or "None"]

### Quick Checklist
- [ ] Functions: small, single purpose, few arguments
- [ ] Names: reveal intent, consistent
- [ ] Classes: SRP, cohesive
- [ ] Errors: exceptions, no null, context
- [ ] Tests: exist, readable, F.I.R.S.T.
- [ ] Comments: necessary, accurate
- [ ] Smells: none detected

### Verdict
[Overall assessment — approve, needs work, or major concerns]
```

## PRESENT TO USER

Show the final summary to the user. Highlight blocking issues first.

Feedback principles (from `collaboration/rules.md`):
1. **Be specific**: Point to exact lines/functions
2. **Explain why**: Reference principles, not just preferences
3. **Suggest alternatives**: Don't just criticize, propose solutions
4. **Prioritize**: Distinguish blocking issues from nice-to-haves
5. **Be respectful**: Critique code, not the person

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `9` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The code review workflow is complete. The full report is saved at the output path.
