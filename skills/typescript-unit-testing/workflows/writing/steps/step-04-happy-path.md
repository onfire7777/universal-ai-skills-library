---
name: 'step-04-happy-path'
description: 'Implement happy path tests for all public methods'
nextStepFile: './step-05-edge-cases.md'
referenceFiles:
  - 'references/common/rules.md'
  - 'references/common/assertions.md'
---

# Step 4: Implement Happy Path Tests

## STEP GOAL

Write happy path tests for all public methods following the AAA pattern with specific assertions.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/rules.md` — AAA pattern, naming conventions, coverage requirements
- `references/common/assertions.md` — Assertion patterns and matchers

## EXECUTION

### 1. Write Happy Path Tests

For each public method, write tests with valid, typical input:

```typescript
describe('methodName', () => {
  it('should return expected result when valid input provided', async () => {
    // Arrange
    const input = { /* valid input data */ };
    mockDependency.method.mockResolvedValue({ /* expected return */ });

    // Act
    const result = await target.methodName(input);

    // Assert
    expect(result).toEqual({ /* expected output */ });
    expect(mockDependency.method).toHaveBeenCalledWith(/* expected args */);
    expect(mockDependency.method).toHaveBeenCalledTimes(1);
  });
});
```

### 2. Follow AAA Pattern

Every test must have clearly separated sections with comments:
- **Arrange** — set up input data and mock return values
- **Act** — call the method under test
- **Assert** — verify the result and mock interactions

### 3. Verify Assertions Are Specific

Assertions must validate specific values, not just existence:

| Don't | Do Instead |
|-------|------------|
| `expect(result).toBeDefined()` | `expect(result).toEqual({ specific: 'value' })` |
| `expect(result).toBeTruthy()` | `expect(result).toBe(true)` |
| `expect(result.length).toBeGreaterThan(0)` | `expect(result).toHaveLength(3)` |

### 4. Append to Report

Append to the output document:

```markdown
## Step 4: Happy Path Tests

**Tests Written:**
- [methodName]: should [description] - DONE
- [methodName]: should [description] - DONE

**Assertions Verified**: All use specific value checks
```

## PRESENT FINDINGS

Show the user:
- All happy path tests written
- Confirmation that assertions are specific
- Any methods where happy path was unclear (ask for clarification)

Then ask: **[C] Continue to Step 5: Implement Edge Case Tests**

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `4` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-05-edge-cases.md`.
