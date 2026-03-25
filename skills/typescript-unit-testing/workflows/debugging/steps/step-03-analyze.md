---
name: 'step-03-analyze'
description: 'Deep analysis based on failure type to identify root cause'
nextStepFile: './step-04-implement-fix.md'
referenceFiles:
  - 'references/mocking/deep-mocked.md'
  - 'references/mocking/jest-native.md'
---

# Step 3: Analyze Based on Failure Type

## STEP GOAL

Perform deep analysis specific to the classified failure type to identify the exact root cause.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/mocking/deep-mocked.md` - DeepMocked patterns and usage
- `references/mocking/jest-native.md` - Native Jest mock patterns

These references help identify whether mocks are correctly set up and used.

## EXECUTION

Based on the failure type classified in Step 2, follow the corresponding analysis path below. Apply **all checks** in the relevant section.

### Type A: Assertion Failure

1. Compare expected vs received in detail:
   ```
   Expected: { id: 'user-123', email: 'test@example.com' }
   Received: { id: 'user-123', email: undefined }
   ```

2. Check possible causes:
   - [ ] Mock not returning expected data
   - [ ] Wrong mock method being called
   - [ ] Logic error in target code
   - [ ] Test expectation incorrect

3. Add debug logging to investigate:
   ```typescript
   it('failing test', async () => {
     // Arrange
     console.log('Mock setup:', mockService.method.getMockImplementation());

     // Act
     const result = await target.methodName(input);
     console.log('Result:', JSON.stringify(result, null, 2));

     // Assert
     expect(result).toEqual(expected);
   });
   ```

### Type B: Unexpected Exception

1. Identify where the exception is thrown:
   ```typescript
   it('failing test', async () => {
     try {
       const result = await target.methodName(input);
       console.log('Success:', result);
     } catch (error) {
       console.log('Error type:', error.constructor.name);
       console.log('Error message:', error.message);
       console.log('Stack:', error.stack);
       throw error;
     }
   });
   ```

2. Check possible causes:
   - [ ] Mock throwing instead of resolving
   - [ ] Missing required mock return value
   - [ ] Validation failing in target code
   - [ ] Dependency not properly mocked

### Type C: Timeout

1. Check for unresolved promises:
   - Mock method called but not configured
   - Async operation waiting forever
   - Event that never fires

2. Debug technique:
   ```typescript
   it('timing out test', async () => {
     console.log('Before call');
     const promise = target.methodName(input);
     console.log('After call, waiting...');
     const result = await promise;
     console.log('Resolved:', result);
   }, 30000); // Extended timeout for debugging
   ```

3. Common fixes:
   - Add `mockResolvedValue` to async mocks
   - Add `mockReturnValue` to sync mocks
   - Check if mock method name matches actual call

### Type D: Mock Error

1. Verify mock is properly created:
   ```typescript
   beforeEach(() => {
     mockService = createMock<ServiceType>();
     console.log('Mock created:', Object.keys(mockService));
   });
   ```

2. Verify mock is injected:
   ```typescript
   const module = await Test.createTestingModule({
     providers: [
       Target,
       { provide: ServiceType, useValue: mockService },
     ],
   }).compile();

   // Verify injection
   const injected = module.get(ServiceType);
   console.log('Injected same as mock:', injected === mockService);
   ```

### For Module Error or Type Error

Apply the same systematic approach:
- **Module Error**: Verify import paths, check for circular dependencies, confirm module exists
- **Type Error**: Check variable declarations, verify imports, check for typos in identifiers

## PRESENT FINDINGS

Present to the user:

```
Step 3: Root Cause Analysis
===========================

Analysis Type: [A/B/C/D]
Root Cause: [specific root cause identified]
Evidence: [what pointed to this cause]
Affected Code: [file and line]
```

Then ask: **[C] Continue to Step 4: Implement Fix**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the root cause analysis findings to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-implement-fix.md`.
