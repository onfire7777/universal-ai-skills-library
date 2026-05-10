---
name: 'step-09-verification'
description: 'Verify setup works correctly — start Docker, run example test, mark complete'
referenceFiles:
  - 'references/common/debugging.md'
---

# Step 9: Verification

## STEP GOAL

Verify the setup works correctly by starting Docker infrastructure and running the example test.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/debugging.md` — for troubleshooting if verification fails

## EXECUTION

### 1. Start Docker Infrastructure

```bash
npm run docker:e2e
```

### 2. Wait for Services to be Healthy

Verify each service is responsive.

### 3. Run Example Test

```bash
npm run test:e2e -- test/e2e/example.e2e-spec.ts
```

### 4. Report Results

### 5. Append to Output

Append to the output document:

```markdown
## Step 9: Verification

**Docker Services**: {{status of each service}}
**Example Test**: {{passed/failed}}

### Next Steps
1. Review references/common/rules.md for mandatory patterns
2. Use references/common/test-case-creation-guide.md for writing tests
3. Check technology-specific examples in references/{tech}/examples.md
```

## PRESENT TO USER

```
Setup Verification:

✅ Docker services started
✅ Database connection successful
✅ Example test passed

Your E2E testing infrastructure is ready!

Next steps:
1. Review references/common/rules.md for mandatory patterns
2. Use references/common/test-case-creation-guide.md for writing tests
3. Check technology-specific examples in references/{tech}/examples.md
```

### Post-Workflow Recommendations

After completing setup, recommend the user review:
1. `references/common/rules.md` - Mandatory patterns they must follow
2. `references/common/best-practices.md` - Test design guidelines
3. Technology-specific `rules.md` files for their stack

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `9` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The setup workflow is complete. The setup checklist is saved at the output path.
