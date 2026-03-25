---
name: 'step-09-document'
description: 'Document optimizations and mark complete'
---

# Step 9: Document Optimizations

## STEP GOAL

Record what was changed for future reference and mark the workflow complete.

## EXECUTION

### 1. Compile Documentation

Append to the output document:

```markdown
## Optimization Summary

### Changes Applied
1. **Infrastructure**: {{description}}
   - Impact: -{{X}}s setup time

2. **Setup**: {{description}}
   - Impact: -{{Y}}s per test

3. **Execution**: {{description}}
   - Impact: -{{Z}}s average per affected test

4. **Code Organization**: {{description}}
   - Impact: {{description}}

### Performance Results
- Total time: {{before}} → {{after}} ({{improvement}}%)
- All tests passing
- Stability: {{X}}/{{X}} runs passed

### Notes
- {{any caveats or limitations}}
```

### 2. Clean Up Temp Files

```bash
rm -f /tmp/e2e-${E2E_SESSION}-*.log
```

## PRESENT TO USER

```
Optimization Complete

Summary:
- Total improvement: {X}%
- Time saved: {Y} minutes per run

Changes documented in: {output path}

Recommendations for ongoing performance:
1. Monitor for new fixed waits in PR reviews
2. Re-run baseline monthly
3. Consider parallel test execution if tests are fully isolated

[D] Done
```

### Post-Workflow Recommendations

After optimization, review and update:
1. `references/common/best-practices.md` - Add new patterns discovered
2. `references/kafka/performance.md` - Update Kafka-specific optimizations
3. Test helpers - Incorporate optimized patterns

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `9` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The optimization workflow is complete. The report is saved at the output path.
