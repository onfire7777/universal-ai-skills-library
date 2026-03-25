# New Feature Workflow

Writing new code following clean code principles from the start.

## When to Use

- Implementing a new feature or user story
- Adding new functionality to an existing system
- Creating a new module or service

## Prerequisites

- Clear understanding of requirements
- Access to the codebase
- Test framework set up

**Reference**: `functions/rules.md`, `naming/rules.md`, `classes/rules.md`, `tdd/rules.md`

---

## Workflow Steps

### Step 1: Understand Requirements

**Goal**: Know exactly what you're building before writing code.

- [ ] Read the ticket/story/requirement completely
- [ ] Identify acceptance criteria
- [ ] Clarify any ambiguities with stakeholders
- [ ] Understand the "why" behind the feature

**Ask**:
- What problem does this solve?
- What are the edge cases?
- What should happen on errors?
- How will this be tested?

**Reference**: `acceptance-testing/knowledge.md`

---

### Step 2: Plan the Approach

**Goal**: Design before coding.

- [ ] Identify the components/classes needed
- [ ] Define the interfaces/contracts
- [ ] Consider dependencies
- [ ] Plan for testability

**Design Questions**:
- What are the nouns? (potential classes)
- What are the verbs? (potential methods)
- What are the dependencies?
- How will I test this?

**Sketch the structure**:
```typescript
// Example: Feature to export user data

// Nouns: UserExporter, ExportFormat, ExportResult
// Verbs: export, format, validate

interface UserExporter {
  export(userId: string, format: ExportFormat): Promise<ExportResult>;
}
```

---

### Step 3: Write Acceptance Test(s)

**Goal**: Define "done" with executable tests.

**Reference**: `acceptance-testing/rules.md`

```typescript
describe('User Export Feature', () => {
  it('should export user data as JSON', async () => {
    // Arrange
    const userId = 'user-123';
    await createTestUser(userId, { name: 'Alice', email: 'alice@example.com' });
    
    // Act
    const result = await userExporter.export(userId, 'json');
    
    // Assert
    expect(result.success).toBe(true);
    expect(result.data).toContain('"name":"Alice"');
  });
});
```

- [ ] Tests describe the acceptance criteria
- [ ] Tests are readable by non-developers
- [ ] Tests fail (feature doesn't exist yet)

---

### Step 4: Build with TDD

**Goal**: Implement using red-green-refactor cycles.

**Reference**: `tdd.md` workflow

For each piece of functionality:

1. **RED**: Write a small failing unit test
2. **GREEN**: Write minimum code to pass
3. **REFACTOR**: Clean up while green

**Start with the simplest case**:
```typescript
describe('UserExporter', () => {
  it('should throw if user not found', async () => {
    const exporter = new UserExporter(mockRepo);
    await expect(exporter.export('nonexistent', 'json'))
      .rejects.toThrow('User not found');
  });
});
```

**Build up complexity**:
```
1. Error case: user not found
2. Simple case: export one field
3. Full case: export all fields
4. Format: JSON format
5. Format: CSV format
6. Edge: special characters
```

---

### Step 5: Apply Clean Code Principles

**Goal**: Write clean code from the start.

While building, continuously apply:

**Functions** (`functions/rules.md`):
- [ ] Small (5-20 lines)
- [ ] Do one thing
- [ ] Descriptive names
- [ ] Few arguments (≤3)

**Naming** (`naming/rules.md`):
- [ ] Intention-revealing names
- [ ] No abbreviations
- [ ] Consistent terminology

**Error Handling** (`error-handling/rules.md`):
- [ ] Use exceptions, not error codes
- [ ] Don't return null
- [ ] Provide context in errors

**Example of clean implementation**:
```typescript
class UserExporter {
  constructor(private readonly userRepository: UserRepository) {}

  async export(userId: string, format: ExportFormat): Promise<ExportResult> {
    const user = await this.findUserOrThrow(userId);
    const formatter = this.getFormatter(format);
    const data = formatter.format(user);
    
    return { success: true, data };
  }

  private async findUserOrThrow(userId: string): Promise<User> {
    const user = await this.userRepository.findById(userId);
    if (!user) {
      throw new UserNotFoundError(userId);
    }
    return user;
  }

  private getFormatter(format: ExportFormat): Formatter {
    const formatters: Record<ExportFormat, Formatter> = {
      json: new JsonFormatter(),
      csv: new CsvFormatter(),
    };
    return formatters[format] ?? throw new UnsupportedFormatError(format);
  }
}
```

---

### Step 6: Review Against Checklist

**Goal**: Self-review before considering "done".

**Reference**: `functions/checklist.md`, `code-review.md` workflow

- [ ] All acceptance tests pass
- [ ] Unit test coverage is adequate
- [ ] Functions are small and focused
- [ ] Names reveal intent
- [ ] No code smells
- [ ] Error handling is clean
- [ ] No commented-out code
- [ ] Code is formatted consistently

---

### Step 7: Refactor if Needed

**Goal**: Polish the implementation.

**Reference**: `refactoring.md` workflow

- [ ] Remove any duplication
- [ ] Improve any unclear names
- [ ] Extract any long methods
- [ ] Simplify any complex logic
- [ ] All tests still pass

---

### Step 8: Document if Necessary

**Goal**: Add only necessary documentation.

**Reference**: `comments/rules.md`

- [ ] Public APIs have doc comments
- [ ] Complex "why" decisions are explained
- [ ] No redundant comments
- [ ] README updated if needed

**Good documentation**:
```typescript
/**
 * Exports user data in the specified format.
 * 
 * @throws UserNotFoundError if user doesn't exist
 * @throws UnsupportedFormatError if format is not supported
 */
export(userId: string, format: ExportFormat): Promise<ExportResult>
```

---

## Quick Checklist

```
[ ] Requirements understood
[ ] Approach planned
[ ] Acceptance tests written
[ ] Built with TDD (red-green-refactor)
[ ] Clean code principles applied
[ ] Self-reviewed against checklist
[ ] Refactored if needed
[ ] Documented if necessary
[ ] All tests pass
[ ] Ready for code review
```

---

## Exit Criteria

Feature is complete when:
- [ ] All acceptance criteria are met
- [ ] All tests pass
- [ ] Code follows clean code principles
- [ ] Self-review checklist passes
- [ ] Ready for peer code review
