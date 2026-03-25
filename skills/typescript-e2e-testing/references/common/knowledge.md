# E2E Testing Core Knowledge

## What is E2E Testing?

End-to-End (E2E) testing validates the complete flow of an application from the user's perspective, testing the integration of all components including databases, message brokers, and external services.

## E2E vs Unit vs Integration Testing

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
   /  \     - Full Kafka flow, self-healing scenarios
  /----\    Integration Tests (30%)
 /      \   - RdKafka wrapper validation, Schema Registry
/--------\  Unit Tests (60%)
            - Core logic, utilities, error handling
```

## When to Write E2E Tests

- **MUST Test:** Critical user flows, async message processing, multi-component interactions
- **SHOULD Test:** API endpoints, error handling flows, authentication/authorization
- **NICE TO HAVE:** Edge cases, performance scenarios, recovery mechanisms

## Test File Organization

```
project/
├── src/
│   └── *.spec.ts          # Unit tests (co-located)
└── test/
    ├── e2e/               # E2E tests
    │   ├── *.e2e-spec.ts  # Test files
    │   ├── setup.ts       # Global setup
    │   └── helpers/       # Test utilities
    ├── fixtures/          # Reusable test data
    └── jest-e2e.config.ts # E2E Jest config
```

## E2E Test Lifecycle

```
beforeAll (90s timeout)
├── Start Docker infrastructure
├── Wait for services healthy
├── Initialize NestJS app
├── Connect Kafka microservice
└── Pre-subscribe to topics (if Kafka)

beforeEach
├── Wait for in-flight operations (1-2s)
└── Clear data (DB collections, message buffers)

test
├── GIVEN: Setup test data
├── WHEN: Execute single action
└── THEN: Verify outcomes

afterEach
└── Clean up test data

afterAll (30s timeout)
├── Close app
├── Disconnect clients
└── Stop infrastructure (if applicable)
```

## Key Principles

1. **Test Isolation:** Each test must be independent and not affect others
2. **Real Infrastructure:** Use real databases, message brokers, not mocks
3. **Deterministic Results:** Tests must produce same results on every run
4. **Clear Intent:** Test names and structure must clearly describe behavior
5. **Fast Feedback:** Optimize for quick failure detection while maintaining reliability
