# Unit Testing Rules

## Context Efficiency Rule

**ALWAYS redirect unit test output to temp files**. Test output can be verbose and bloats agent context.

**IMPORTANT**: Use unique session ID in filenames to prevent conflicts when multiple agents run.

```bash
# Initialize session (once at start of testing session)
export UT_SESSION=$(date +%s)-$$

# Standard pattern - redirect to temp file only (no console output)
npm test > /tmp/ut-${UT_SESSION}-output.log 2>&1

# Read summary only (last 50 lines)
tail -50 /tmp/ut-${UT_SESSION}-output.log

# Get failure details
grep -B 2 -A 15 "FAIL\|✕" /tmp/ut-${UT_SESSION}-output.log

# Cleanup when done
rm -f /tmp/ut-${UT_SESSION}-*.log /tmp/ut-${UT_SESSION}-*.md
```

**Temp File Locations** (with `${UT_SESSION}` unique per agent):
- `/tmp/ut-${UT_SESSION}-output.log` - Full test output
- `/tmp/ut-${UT_SESSION}-failures.md` - Tracking file for one-by-one fixing
- `/tmp/ut-${UT_SESSION}-debug.log` - Debug runs
- `/tmp/ut-${UT_SESSION}-verify.log` - Verification runs
- `/tmp/ut-${UT_SESSION}-coverage.log` - Coverage output

---

## One-by-One Fixing Rule

**CRITICAL: When tests fail, NEVER keep running the full suite.**

```
❌ WRONG: Run full suite → See 5 failures → Run full suite again → Still 5 failures → ...
✅ RIGHT: Run full suite → See 5 failures → Fix test 1 only → Verify → Fix test 2 only → ... → Run full suite ONCE
```

**WHY**: Full suite runs waste time and pollute output. Each failing test adds noise, making debugging harder.

### Protocol:
1. Run full suite ONCE to identify all failures
2. Create tracking file `/tmp/ut-${UT_SESSION}-failures.md` listing all failing tests
3. Fix ONE test at a time using `-t "test name"`
4. Verify each fix with 3-5 runs of that specific test
5. Mark as FIXED in tracking file
6. Move to next failing test
7. Run full suite ONLY ONCE after ALL individual tests pass

---

## Core Rules

| Rule | Requirement |
|------|-------------|
| Variable for SUT | Use `target` (MANDATORY) |
| Mock prefix | `mockServiceName` |
| Mock type | `DeepMocked<T>` |
| Pattern | Arrange-Act-Assert with comments |
| Coverage | 80%+ for new code |
| Logger | `.setLogger(new MockLoggerService())` |

## What NOT to Unit Test

Do NOT create unit tests for:

| Type | Reason |
|------|--------|
| **Interfaces** | Type definitions only, no runtime behavior |
| **Enums** | Static value mappings, no logic to test |
| **Constants** | Static values, no behavior |
| **Type aliases** | Type definitions only |
| **DTOs** (plain) | Data structures without logic |

Only test files containing **executable logic** (classes with methods, functions with behavior).

## AAA Pattern (Mandatory)

ALL unit tests MUST follow Arrange-Act-Assert:

```typescript
it('should return user when found', async () => {
  // Arrange
  const userId = 'user-123';
  mockRepository.findById.mockResolvedValue({
    id: userId,
    email: 'test@example.com',
  });

  // Act
  const result = await target.getUser(userId);

  // Assert
  expect(result).toEqual({ id: userId, email: 'test@example.com' });
  expect(mockRepository.findById).toHaveBeenCalledWith(userId);
});
```

### Formatting Guidelines

| Test Length | Format |
|-------------|--------|
| <= 3 lines | No blank lines needed |
| > 3 lines | Blank line between phases |
| Complex phases | Comment labels for each phase |

## Naming Conventions

### Test Files
- Pattern: `*.spec.ts`
- Location: Co-located with source file

### Test Structure
```typescript
describe('ClassName', () => {
  describe('methodName', () => {
    it('should [expected behavior] when [condition]', () => {});
  });
});
```

### Test Naming Pattern
`should [verb] [expected] when [condition]`

```typescript
describe('UserService', () => {
  describe('findById', () => {
    it('should return user when user exists', async () => {...});
    it('should return null when user not found', async () => {...});
    it('should throw NotFoundException when id is invalid', async () => {...});
  });
});
```

## Variable Naming

| Variable | Convention | Example |
|----------|------------|---------|
| SUT | `target` | `let target: UserService;` |
| Mocks | `mock` prefix | `let mockRepository: DeepMocked<UserRepository>;` |
| Mock Type | `DeepMocked<T>` | `DeepMocked<UserRepository>` |
| Test Module | `module` | `let module: TestingModule;` |

## Setup Requirements

```typescript
beforeEach(async () => {
  // 1. Create mocks
  mockRepository = createMock<UserRepository>();

  // 2. Create testing module
  const module: TestingModule = await Test.createTestingModule({
    providers: [
      UserService,
      { provide: UserRepository, useValue: mockRepository },
    ],
  })
    .setLogger(new MockLoggerService())  // REQUIRED
    .compile();

  // 3. Get target
  target = module.get<UserService>(UserService);
});

afterEach(() => {
  jest.clearAllMocks();  // REQUIRED
});
```

## Test Isolation Rules

1. **Fresh mocks per test**: Create mocks in `beforeEach`
2. **Clear mocks after**: Use `jest.clearAllMocks()` in `afterEach`
3. **No shared state**: Never share mutable state between tests
4. **Deterministic results**: Tests must produce same results on every run

## Coverage Requirements

| Category | Priority | Must Test |
|----------|----------|-----------|
| Happy path | MANDATORY | All public methods |
| Edge cases | MANDATORY | Empty, null, boundaries |
| Error cases | MANDATORY | Not found, validation |
| Exceptions | MANDATORY | Type, code, message |
| Business rules | MANDATORY | Domain logic |

**Coverage Target:** 80%+ for new code

## Mock Rules

### Always verify mock calls:
```typescript
expect(mockRepository.save).toHaveBeenCalledWith({
  email: 'test@example.com',
  name: 'John',
});
expect(mockRepository.save).toHaveBeenCalledTimes(1);
```

### Never mock internal implementation:
```typescript
// BAD - Mocking internal mapper
expect(mockMapper.toEntity).toHaveBeenCalled();

// GOOD - Verify behavior
expect(result).toEqual(expectedUser);
expect(mockRepository.save).toHaveBeenCalledWith(
  expect.objectContaining({ email: input.email })
);
```

## Checklist

**Setup:**
- [ ] Use `target` for system under test
- [ ] Use `mock` prefix for all mocks
- [ ] Use `DeepMocked<T>` type
- [ ] Include `.setLogger(new MockLoggerService())`
- [ ] Follow AAA pattern with comments
- [ ] Direct mock methods, not `jest.spyOn()`
- [ ] Reset mocks in `afterEach`

**Test Coverage (MANDATORY):**
- [ ] Happy path tests for all public methods
- [ ] Edge case tests (empty inputs, null values, boundaries)
- [ ] Error case tests (not found, validation failures)
- [ ] Exception behavior tests (correct type, error code, message)
- [ ] Business rule tests (domain logic, calculations)

**Assertions:**
- [ ] Assert specific values, not just types or existence
- [ ] Validate ALL result properties with expected values
- [ ] Verify mock functions called with correct parameters
- [ ] Check mock call counts (`.toHaveBeenCalledTimes()`)
- [ ] No conditional assertions (no `if` statements)
- [ ] Tests fail when any field differs from expectations
