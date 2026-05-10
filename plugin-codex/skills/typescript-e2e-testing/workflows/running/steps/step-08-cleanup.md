---
name: 'step-08-cleanup'
description: 'Post-run cleanup — stop Docker, delete tracking files, mark complete'
---

# Step 8: Post-Run Cleanup

## STEP GOAL

Clean up after test run.

## EXECUTION

### 1. Stop Docker if Not Needed

```bash
docker-compose -f docker-compose.e2e.yml down
```

### 2. Delete Tracking Files

```bash
rm -f /tmp/e2e-${E2E_SESSION}-*.log /tmp/e2e-${E2E_SESSION}-*.md
```

### 3. Review and Commit Fixes

If test fixes were made, present for review.

### 4. Append to Output

Append to the output document:

```markdown
## Step 8: Cleanup

**Docker**: {{stopped/kept running}}
**Temp Files**: {{cleaned/kept}}
**Results Summary**:
- Passed: {{x}}
- Failed: {{y}}
- Fixed: {{z}}
```

## PRESENT TO USER

```
Test Run Complete

Results Summary:
- Passed: {x}
- Failed: {y}
- Fixed: {z}

Cleanup Options:
[D] Stop Docker services
[K] Keep Docker running
[C] Commit test fixes
```

### Post-Workflow Recommendations

After running tests, if failures occurred:
1. Review `references/common/debugging.md` for systematic debugging
2. Review technology-specific guides for common issues
3. Check `references/common/examples.md` for correct patterns

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `8` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The running workflow is complete. The run log is saved at the output path.
