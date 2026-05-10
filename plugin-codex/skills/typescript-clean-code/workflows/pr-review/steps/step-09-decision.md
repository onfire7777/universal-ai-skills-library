---
name: 'step-09-decision'
description: 'Make the final PR decision — approve, request changes, or comment'
---

# Step 9: Make the Decision

## STEP GOAL

Make the final review decision: approve, request changes, or comment. Complete the review report.

## DECISION CRITERIA

### Approve when:
- All blocking issues resolved
- Tests are adequate
- Code quality is acceptable
- You'd be comfortable maintaining this code

### Request Changes when:
- Blocking issues exist
- Tests are missing for critical paths
- Security vulnerabilities found
- Major design problems

### Comment when:
- You have questions but no blockers
- You want to discuss approaches
- You're not the final approver

## COMPILATION

### 1. Write Final Decision

Append to the output document:

```markdown
## Final Decision

### PR Review Checklist

Context:
- [ ] PR description read
- [ ] Ticket/issue understood
- [ ] Scope is appropriate

Tests:
- [ ] New code has tests
- [ ] Tests are meaningful
- [ ] Edge cases covered

Code Quality:
- [ ] Functions are small
- [ ] Names are clear
- [ ] No code smells
- [ ] Error handling is proper

Security:
- [ ] No vulnerabilities
- [ ] Input validated
- [ ] No secrets in code

Performance:
- [ ] No obvious issues
- [ ] No N+1 queries

Feedback:
- [ ] Constructive and specific
- [ ] Clear blocking vs suggestions
- [ ] Actionable recommendations

### Decision: [APPROVE / REQUEST CHANGES / COMMENT]

**Reason**: [summary]

### Common Issues Found

| Issue | Location | Category |
|-------|----------|----------|
| ... | file:line | BLOCKING/SUGGESTION/NIT |
```

## PRESENT TO USER

Show the final decision and checklist to the user.

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `9` to `stepsCompleted`
- Set `status` to `'complete'`
- Set `decision` to the chosen decision

## WORKFLOW COMPLETE

The PR review workflow is complete. The full report is saved at the output path.
