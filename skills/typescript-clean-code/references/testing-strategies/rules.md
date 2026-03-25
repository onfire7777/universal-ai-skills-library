# Testing Strategies Rules

Rules for implementing a comprehensive testing strategy across all levels of the test pyramid.

## Core Rules

### 1. Unit Tests

Tests written by programmers, for programmers, in the system's language.

- Write tests BEFORE production code (TDD)
- Execute as part of Continuous Integration
- Target ~90%+ true coverage (not false tests that execute without asserting)

**Coverage**: As close to 100% as practical (typically 90s%)

### 2. Component Tests

Acceptance tests for individual components encapsulating business rules.

- Wrap individual components with input/output verification
- Decouple from other components using mocks and test doubles
- Business should be able to read and interpret these tests
- Written by QA and Business with development assistance

**Coverage**: ~50% of system (happy paths + obvious edge cases)

### 3. Integration Tests

Tests that verify component assemblies communicate correctly.

- Test choreography, not business rules
- Verify plumbing and connections between components
- Written by system architects or lead designers
- Run periodically (nightly/weekly), not on every CI build
- Include performance and throughput tests at this level

**Coverage**: Architectural soundness verification

### 4. System Tests

Automated tests against the entire integrated system.

- Ultimate integration tests
- Verify system is wired together correctly
- Test that parts interoperate according to plan
- Include throughput and performance tests
- Written by architects and technical leads

**Coverage**: ~10% of system (construction, not behavior)

### 5. Manual Exploratory Tests

Human-driven, unscripted testing.

- NOT automated
- NOT scripted (no written test plans)
- Explore for unexpected behaviors
- Confirm expected behaviors
- Use human creativity to investigate the system

**Goal**: Find peculiarities, not prove coverage

## Guidelines

- Run tests as frequently as possible for maximum feedback
- Lower pyramid levels should run more frequently
- Business writes happy-path tests; QA writes corner/boundary/unhappy-path tests
- React to every QA-found bug with process improvement
- Keep system continuously clean through continuous testing

## Exceptions

- **Small systems**: May skip integration tests if few components exist
- **Performance-critical paths**: May need additional system-level tests beyond 10%
- **Rapidly changing UI**: May reduce component test coverage in favor of unit tests

## Quick Reference

| Test Level | Coverage | Frequency | Written By |
|------------|----------|-----------|------------|
| Unit | ~90%+ | Every CI build | Developers |
| Component | ~50% | Every CI build | QA + Business + Dev |
| Integration | Architectural | Periodic (nightly/weekly) | Architects |
| System | ~10% | Periodic | Tech Leads |
| Exploratory | N/A | As needed | Humans (anyone) |
