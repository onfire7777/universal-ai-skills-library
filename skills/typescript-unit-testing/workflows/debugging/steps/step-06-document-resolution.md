---
name: 'step-06-document-resolution'
description: 'Document the complete debugging resolution'
---

# Step 6: Document Resolution

## STEP GOAL

Record the complete debugging session with root cause, fix, lessons learned, and a reference of debug techniques for future use.

## EXECUTION

### 1. Compile Resolution Report

Append the full resolution report to the output document:

```markdown
## Debug Resolution Report

### Original Failure
- Test: [test name]
- Error: [error message]
- Date: [date]

### Root Cause
[Description of what caused the failure]

### Fix Applied
- File: [path]
- Change: [specific change made]

### Lessons Learned
[Any insights for preventing similar issues]

### Verification
- Fixed test: PASS
- Related tests: PASS
- Regression check: PASS
```

### 2. Include Debug Techniques Reference

Append the following reference section to the report for future use:

```markdown
## Debug Techniques Reference

### Console Logging

```typescript
it('test with logging', async () => {
  console.log('Input:', input);
  console.log('Mock calls before:', mockService.method.mock.calls);

  const result = await target.methodName(input);

  console.log('Result:', result);
  console.log('Mock calls after:', mockService.method.mock.calls);
});
```

### Check Mock State

```typescript
// Check if mock was called
console.log('Called:', mockService.method.mock.calls.length);

// Check call arguments
console.log('Args:', mockService.method.mock.calls[0]);

// Check return value
console.log('Returns:', mockService.method.mock.results);
```

### Step Through with Debugger

```bash
# Run Jest with debugger (this requires console for interactive debugging)
node --inspect-brk node_modules/.bin/jest --runInBand -t "[test name]"
```

Then attach VS Code debugger or Chrome DevTools.

### Isolate Setup Issues

```typescript
describe('Target', () => {
  it('should create module', async () => {
    const module = await Test.createTestingModule({
      providers: [Target],
    }).compile();

    expect(module).toBeDefined();
  });

  it('should get target', async () => {
    const module = await Test.createTestingModule({
      providers: [Target],
    }).compile();

    const target = module.get<Target>(Target);
    expect(target).toBeDefined();
  });
});
```
```

### 3. Include Common Issues and Solutions

Append the quick-reference table:

```markdown
## Common Issues and Solutions

| Issue | Symptom | Solution |
|-------|---------|----------|
| Mock not returning | `undefined` result | Add `mockResolvedValue()` |
| Mock not called | Call count is 0 | Verify mock injection |
| Wrong mock called | Different mock has calls | Check dependency injection |
| Async not awaited | Promise object in result | Add `await` to call |
| Mock reset missing | State from other tests | Add `jest.clearAllMocks()` |
| Private method tested | Cannot access | Test via public interface |
| Logger errors | Console noise | Add `MockLoggerService` |
```

### 4. Post-Workflow Verification Checklist

```markdown
## Post-Workflow Verification

1. Re-read `references/common/rules.md` to ensure fix follows conventions
2. Confirm all tests pass
3. Run coverage to ensure fix did not reduce coverage
4. Consider if similar issues exist elsewhere
```

### 5. Cleanup Temp Files

```bash
rm -f /tmp/ut-${UT_SESSION}-*.log /tmp/ut-${UT_SESSION}-*.md
```

## PRESENT FINDINGS

Present to the user the compiled resolution report including:
- Original failure summary
- Root cause
- Fix applied
- Lessons learned
- Verification status

## FRONTMATTER UPDATE

Update the output document:
- Add `6` to `stepsCompleted`
- Set `status` to `'complete'`

## WORKFLOW COMPLETE

The debugging workflow is finished. The full debug resolution report is saved at the output path. Use this report as a reference if similar test failures occur in the future.
