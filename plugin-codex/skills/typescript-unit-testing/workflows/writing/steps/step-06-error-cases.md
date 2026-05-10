---
name: 'step-06-error-cases'
description: 'Implement error case tests for exceptions and error handling'
nextStepFile: './step-07-business-rules.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/assertions.md'
---

# Step 6: Implement Error Case Tests

## STEP GOAL

Write exception and error handling tests that verify exception types, error messages, and error codes.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements
- `references/common/assertions.md` — Assertion patterns and matchers

## EXECUTION

### 1. Write Exception Tests

For each exception identified in Step 1, write tests covering the exception type, message, and error code.

#### NotFoundException

```typescript
it('should throw NotFoundException when resource not found', async () => {
  // Arrange
  mockRepository.findById.mockResolvedValue(null);

  // Act & Assert
  await expect(target.getById('invalid-id')).rejects.toThrow(NotFoundException);
});
```

#### BadRequestException with Message

```typescript
it('should throw BadRequestException with correct message when validation fails', async () => {
  // Arrange
  const invalidInput = { email: 'invalid-email' };

  // Act & Assert
  await expect(target.create(invalidInput)).rejects.toThrow(
    expect.objectContaining({
      message: expect.stringContaining('Invalid email'),
    })
  );
});
```

#### Exception with Error Code

```typescript
it('should throw exception with correct error code', async () => {
  // Arrange
  mockService.process.mockRejectedValue(new Error('External failure'));

  // Act & Assert
  await expect(target.execute()).rejects.toMatchObject({
    errorCode: 'EXTERNAL_SERVICE_ERROR',
  });
});
```

### 2. Error Verification Checklist

For each exception thrown by the component, verify:
- [ ] Exception type (NotFoundException, BadRequestException, etc.)
- [ ] Error message content
- [ ] Error code (if applicable)
- [ ] Error is propagated correctly from dependencies

### 3. Common Error Patterns

Use these assertion patterns for error testing:

| Pattern | Use When |
|---------|----------|
| `rejects.toThrow(ExceptionType)` | Verifying exception class |
| `rejects.toThrow(expect.objectContaining({...}))` | Verifying exception properties |
| `rejects.toMatchObject({...})` | Verifying error shape |
| `rejects.toThrow('message')` | Verifying error message string |

### 4. Append to Report

Append to the output document:

```markdown
## Step 6: Error Case Tests

**Tests Written:**
- [methodName]: should throw [ExceptionType] when [condition] - DONE

**Error Verification:**
- Exception types: [verified/not applicable]
- Error messages: [verified/not applicable]
- Error codes: [verified/not applicable]
```

## PRESENT FINDINGS

Show the user:
- All error case tests written
- Verification coverage for exception types, messages, and codes
- Any error paths that could not be tested (explain why)

Then ask: **[C] Continue to Step 7: Implement Business Rule Tests**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `6` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-07-business-rules.md`.
