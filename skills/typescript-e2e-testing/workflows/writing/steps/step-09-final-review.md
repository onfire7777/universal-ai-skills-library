---
name: 'step-09-final-review'
description: 'Final review against rules — verify GWT, isolation, assertions, mark complete'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/examples.md'
---

# Step 9: Final Review

## STEP GOAL

Review the complete test file against all rules and patterns before saving.

## REFERENCE LOADING

Before starting review, load and read:
- `references/common/rules.md` — verify compliance
- `references/common/examples.md` — compare with reference implementations

## EXECUTION

### 1. Review Checklist

- [ ] All tests follow GWT pattern
- [ ] Each test has exactly ONE action in WHEN
- [ ] No conditional assertions
- [ ] Unique identifiers used (no hardcoded IDs)
- [ ] Database cleaned in beforeEach
- [ ] Specific assertions (not just `toBeDefined`)
- [ ] Error cases verify no side effects
- [ ] Timeouts appropriate for operations
- [ ] Async tests use polling, not fixed waits

### 2. Present Complete Test File

Show the complete test file for final confirmation.

### 3. Append to Output

Append to the output document:

```markdown
## Step 9: Final Review

### Compliance Check
{{checklist results}}

### Complete Test File
{{show final code}}
```

## PRESENT TO USER

```
Test File Review:

✅ {check 1}
✅ {check 2}
⚠️ {warning if any}

Complete test file:
{show final code}

[C] Confirm and save / [M] Modify
```

### Post-Workflow Recommendations

After writing tests, recommend reviewing:
1. `references/common/examples.md` - Verify patterns match examples
2. Technology-specific `examples.md` - Compare with reference implementations
3. `references/common/debugging.md` - For when tests fail

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `9` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The test writing workflow is complete. The test plan and code are saved at the output path.
