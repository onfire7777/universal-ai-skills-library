# Unit Testing Core Knowledge

## What is Unit Testing?

Unit testing validates individual functions, methods, and classes in complete isolation by mocking all external dependencies. Tests should be fast (milliseconds) and deterministic.

## Unit vs Integration vs E2E Testing

| Aspect | Unit Tests | Integration Tests | E2E Tests |
|--------|------------|-------------------|-----------|
| Scope | Single function/class | Multiple components | Full system |
| Dependencies | All mocked | Partial real | All real |
| Speed | Fast (ms) | Medium (s) | Slow (min) |
| Pattern | AAA (Arrange-Act-Assert) | AAA or GWT | GWT (Given-When-Then) |
| Isolation | Complete | Partial | None |
| Purpose | Verify logic | Verify integration | Verify behavior |

## Test Pyramid Distribution

```
    /\      E2E Tests (10%)
   /  \     - Critical user flows only
  /----\    Integration Tests (20%)
 /      \   - Component interactions
/--------\  Unit Tests (70%)
            - Core logic, utilities, error handling
```

## Framework & Dependencies

```json
{
  "@nestjs/testing": "^11.0.12",
  "jest": "^29.7.0",
  "@golevelup/ts-jest": "^0.4.0",
  "mongodb-memory-server": "^9.0.0",
  "pg-mem": "^1.0.0"
}
```

## Testing Philosophy

### Core Principles

- **Comprehensive Coverage**: Fully test core user flows, edge cases, and all expected behaviors
- **Behavior-Driven**: Test behavior and outcomes, not implementation details
- **Fast Execution**: Keep tests fast (milliseconds) to enable rapid feedback
- **Isolated Testing**: Mock external dependencies to ensure unit isolation

### What MUST Be Tested

| Category | Description | Priority |
|----------|-------------|----------|
| **Core User Flows** | Happy path scenarios | MANDATORY |
| **Edge Cases** | Boundary conditions, empty inputs, null values | MANDATORY |
| **Error Handling** | Exception paths, validation failures | MANDATORY |
| **Exception Behavior** | Correct type, error code, message | MANDATORY |
| **Business Rules** | Domain logic, calculations | MANDATORY |
| **Input Validation** | Invalid inputs, malformed data | MANDATORY |

### What NOT to Test

- Logger calls (implementation detail)
- Framework behavior (NestJS internals)
- Third-party library internals
- Private methods directly (test via public interface)

## Test File Organization

```
src/
├── user/
│   ├── user.service.ts
│   ├── user.service.spec.ts        # Unit tests (co-located)
│   ├── user.controller.ts
│   └── user.controller.spec.ts
test/
├── helpers/                         # Shared test utilities
│   ├── mongodb-test.helper.ts
│   ├── pgmock.helper.ts
│   └── factory/
│       └── user.factory.ts
└── jest.config.ts
```

## Test Lifecycle

```
beforeAll (if needed)
├── Expensive one-time setup

beforeEach
├── Create mocks
├── Create TestingModule
├── Get target instance
└── Reset mock state

test
├── Arrange: Setup test data and mock returns
├── Act: Execute single action
└── Assert: Verify outcomes

afterEach
└── Clear all mocks

afterAll (if needed)
└── Cleanup resources
```
