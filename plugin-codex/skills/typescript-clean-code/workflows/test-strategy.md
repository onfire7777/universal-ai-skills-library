# Test Strategy Workflow

Planning test coverage for a feature or system.

## When to Use

- Starting a new project/feature
- Adding tests to legacy code
- Reviewing test coverage gaps
- Planning QA strategy

## Prerequisites

- Understanding of the feature/system
- Knowledge of available testing tools
- Access to requirements/acceptance criteria

**Reference**: `testing-strategies/rules.md`, `unit-tests/rules.md`, `acceptance-testing/rules.md`

---

## The Test Pyramid

```
        /\
       /  \        Manual Exploratory (~5%)
      /----\       System/E2E Tests (~10%)
     /      \      Integration Tests (~20%)
    /--------\     Component Tests (~25%)
   /          \    Unit Tests (~40%)
  /______________\
```

**Principle**: More tests at the bottom, fewer at the top.

---

## Workflow Steps

### Step 1: Understand What to Test

**Goal**: Identify the scope and boundaries.

- [ ] What is the feature/system?
- [ ] What are the critical paths?
- [ ] What are the edge cases?
- [ ] What could go wrong?

**Map the components**:
```
User Export Feature
├── API Layer (Express routes)
├── Service Layer (UserExporter)
├── Data Layer (UserRepository)
├── External (File storage)
└── Formatters (JSON, CSV)
```

---

### Step 2: Define Unit Test Coverage

**Goal**: Test individual units in isolation.

**Reference**: `unit-tests/rules.md`

**What to unit test**:
- Pure functions
- Business logic
- Data transformations
- Validation rules
- Edge cases

**Coverage target**: ~90%+ of business logic

**Example unit tests**:
```typescript
describe('JsonFormatter', () => {
  it('should format user as JSON', () => {
    const user = { name: 'Alice', email: 'alice@test.com' };
    const result = formatter.format(user);
    expect(JSON.parse(result)).toEqual(user);
  });

  it('should escape special characters', () => {
    const user = { name: "O'Brien", email: 'test@test.com' };
    const result = formatter.format(user);
    expect(() => JSON.parse(result)).not.toThrow();
  });

  it('should handle empty fields', () => {
    const user = { name: '', email: '' };
    const result = formatter.format(user);
    expect(result).toContain('""');
  });
});
```

**Checklist**:
- [ ] All pure functions tested
- [ ] All business rules tested
- [ ] Edge cases covered
- [ ] Error cases tested

---

### Step 3: Define Component Test Coverage

**Goal**: Test components with mocked dependencies.

**Reference**: `testing-strategies/rules.md`

**What to component test**:
- Service classes with mocked repos
- API routes with mocked services
- Happy paths and error paths

**Coverage target**: ~50% of components

**Example component test**:
```typescript
describe('UserExporter', () => {
  let exporter: UserExporter;
  let mockRepo: jest.Mocked<UserRepository>;

  beforeEach(() => {
    mockRepo = {
      findById: jest.fn(),
    };
    exporter = new UserExporter(mockRepo);
  });

  it('should export user as JSON', async () => {
    mockRepo.findById.mockResolvedValue({ 
      id: '123', 
      name: 'Alice',
      email: 'alice@test.com' 
    });

    const result = await exporter.export('123', 'json');

    expect(result.success).toBe(true);
    expect(result.data).toContain('Alice');
  });

  it('should throw when user not found', async () => {
    mockRepo.findById.mockResolvedValue(null);

    await expect(exporter.export('999', 'json'))
      .rejects.toThrow('User not found');
  });
});
```

**Checklist**:
- [ ] Happy paths tested
- [ ] Error handling tested
- [ ] Dependencies properly mocked

---

### Step 4: Define Integration Test Coverage

**Goal**: Test components working together.

**What to integration test**:
- Database interactions
- External API calls
- File system operations
- Multiple services together

**Coverage target**: Key integration points

**Example integration test**:
```typescript
describe('UserExporter Integration', () => {
  let db: TestDatabase;
  let exporter: UserExporter;

  beforeAll(async () => {
    db = await TestDatabase.create();
    const repo = new UserRepository(db);
    exporter = new UserExporter(repo);
  });

  afterAll(async () => {
    await db.destroy();
  });

  it('should export user from real database', async () => {
    // Arrange
    await db.insertUser({ id: '123', name: 'Alice', email: 'alice@test.com' });

    // Act
    const result = await exporter.export('123', 'json');

    // Assert
    expect(result.success).toBe(true);
  });
});
```

**Checklist**:
- [ ] Database operations tested
- [ ] External integrations tested
- [ ] Uses test database/sandbox

---

### Step 5: Define E2E/System Test Coverage

**Goal**: Test the system as a whole.

**What to E2E test**:
- Critical user journeys
- API contracts
- Authentication flows
- Key business workflows

**Coverage target**: ~10% (critical paths only)

**Example E2E test**:
```typescript
describe('User Export API', () => {
  it('should export user via API', async () => {
    // Create user via API
    const createResponse = await request(app)
      .post('/api/users')
      .send({ name: 'Alice', email: 'alice@test.com' });
    
    const userId = createResponse.body.id;

    // Export via API
    const exportResponse = await request(app)
      .get(`/api/users/${userId}/export?format=json`)
      .expect(200);

    expect(exportResponse.body.success).toBe(true);
    expect(exportResponse.body.data).toContain('Alice');
  });
});
```

**Checklist**:
- [ ] Critical paths covered
- [ ] End-to-end flows work
- [ ] API contracts verified

---

### Step 6: Plan Manual/Exploratory Testing

**Goal**: Catch issues automation misses.

**What to manually test**:
- UI/UX flows
- Edge cases not worth automating
- Exploratory "what if" scenarios
- Accessibility

**Create test scenarios**:
```
Manual Test Scenarios:
1. Export user with very long name (1000 chars)
2. Export while database is slow
3. Export from mobile browser
4. Export with expired session
5. Try to export another user's data
```

**Checklist**:
- [ ] Exploratory scenarios documented
- [ ] Not automatable cases listed
- [ ] Security edge cases included

---

### Step 7: Document the Test Strategy

**Goal**: Create a reference for the team.

**Test Strategy Document**:
```markdown
# User Export - Test Strategy

## Test Pyramid Distribution
- Unit Tests: 15 tests (~40%)
- Component Tests: 10 tests (~25%)  
- Integration Tests: 5 tests (~20%)
- E2E Tests: 3 tests (~10%)
- Manual: 5 scenarios (~5%)

## Coverage by Component
| Component | Unit | Component | Integration | E2E |
|-----------|------|-----------|-------------|-----|
| JsonFormatter | 5 | - | - | - |
| CsvFormatter | 5 | - | - | - |
| UserExporter | 3 | 6 | 2 | - |
| API Routes | 2 | 4 | 3 | 3 |

## Critical Paths (must have E2E)
1. Export user as JSON
2. Export user as CSV
3. Handle user not found

## Test Data Strategy
- Unit: Hardcoded fixtures
- Component: Factory functions
- Integration: Test database with seeding
- E2E: Isolated test environment

## Not Covered (accepted risk)
- Performance under load (separate perf tests)
- Export of 10M+ users (accepted limitation)
```

---

## Test Type Quick Reference

| Type | Speed | Isolation | Scope | When to Use |
|------|-------|-----------|-------|-------------|
| Unit | Fast | Complete | Single function/class | Business logic |
| Component | Fast | Partial (mocks) | One component | Service behavior |
| Integration | Slow | None | Multiple components | DB, APIs |
| E2E | Slowest | None | Full system | Critical paths |
| Manual | N/A | None | Everything | Edge cases, UX |

---

## Exit Criteria

Test strategy is complete when:
- [ ] All test types defined
- [ ] Coverage targets set
- [ ] Critical paths identified
- [ ] Test data strategy defined
- [ ] Manual test scenarios documented
- [ ] Strategy documented for team
