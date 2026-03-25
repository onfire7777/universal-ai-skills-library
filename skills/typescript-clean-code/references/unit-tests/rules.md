# Unit Tests Rules

Rules for writing clean, maintainable unit tests that enable confident refactoring.

## Core Rules

### 1. Three Laws of TDD

Follow these laws for test-driven development:

- **First Law**: Never write production code until you have a failing unit test
- **Second Law**: Write only enough test to fail (not compiling counts as failing)
- **Third Law**: Write only enough production code to pass the failing test

**Cycle time**: ~30 seconds per iteration

### 2. Keep Tests Clean

Test code is just as important as production code.

- Tests require thought, design, and care
- Dirty tests become harder to change as production code evolves
- Dirty tests eventually get abandoned, leading to rotting code
- Maintain tests to the same quality standards as production code

### 3. Prioritize Readability

What makes tests readable: clarity, simplicity, density of expression.

- Say a lot with few expressions
- Hide noisy details behind well-named helper functions
- Use domain-specific testing language
- Make test intent immediately clear

**Example**:
```typescript
// Bad - too much noise
it('should get page hierarchy as XML', async () => {
  await crawler.addPage(root, PathParser.parse('PageOne'));
  await crawler.addPage(root, PathParser.parse('PageOne.ChildOne'));
  request.setResource('root');
  request.addInput('type', 'pages');
  const responder = new SerializedPageResponder();
  const response = await responder.makeResponse(new Context(root), request);
  expect(response.getContentType()).toBe('text/xml');
  expect(response.getContent()).toContain('<name>PageOne</name>');
});

// Good - clear intent with helpers
it('should get page hierarchy as XML', async () => {
  await makePages('PageOne', 'PageOne.ChildOne', 'PageTwo');
  await submitRequest('root', 'type:pages');
  assertResponseIsXML();
  assertResponseContains('<name>PageOne</name>', '<name>PageTwo</name>');
});
```

### 4. Follow BUILD-OPERATE-CHECK Pattern

Structure each test in three clear parts:

- **Build**: Set up the test data
- **Operate**: Perform the action being tested
- **Check**: Verify the expected results

Also known as Arrange-Act-Assert (AAA) or Given-When-Then.

### 5. Single Concept per Test

Test one concept per test function.

- Don't test multiple unrelated things in one test
- Each test should have a single reason to fail
- Split tests that verify multiple independent behaviors
- Helps identify what broke when a test fails

### 6. Minimize Asserts per Concept

Keep assertions focused on the concept being tested.

- One assert per test is a good guideline (not strict rule)
- Multiple asserts are acceptable if testing same concept
- Avoid unrelated assertions in same test
- Create domain-specific assertion helpers

### 7. F.I.R.S.T. Principles

Clean tests follow these five principles:

| Principle | Description |
|-----------|-------------|
| **Fast** | Tests run quickly; slow tests won't be run frequently |
| **Independent** | Tests don't depend on each other; run in any order |
| **Repeatable** | Run in any environment without external dependencies |
| **Self-Validating** | Boolean output - pass or fail, no manual checking |
| **Timely** | Written just before the production code |

## Guidelines

Less strict recommendations:

- Build a domain-specific testing language over time
- Refactor tests to be more expressive as patterns emerge
- Efficiency trade-offs are acceptable in test code
- Use Given-When-Then naming convention for clarity
- Name tests to describe the behavior being verified

## Exceptions

When these rules may be relaxed:

- **Efficiency in tests**: Test code can be less efficient than production code (but never less clean)
- **Multiple asserts**: Acceptable when testing a single logical concept
- **Integration tests**: May require more setup and multiple verifications

## Quick Reference

| Rule | Summary |
|------|---------|
| TDD Laws | Failing test → minimal code → refactor |
| Keep Clean | Test code = production code quality |
| Readability | Clarity, simplicity, density of expression |
| BUILD-OPERATE-CHECK | Setup → Action → Verify |
| Single Concept | One behavior per test |
| Minimize Asserts | Few assertions per concept |
| F.I.R.S.T. | Fast, Independent, Repeatable, Self-Validating, Timely |
