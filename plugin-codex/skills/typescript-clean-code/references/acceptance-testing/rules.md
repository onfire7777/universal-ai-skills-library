# Acceptance Testing Rules

Rules for writing, managing, and executing acceptance tests as professional developers.

## Core Rules

### 1. Always Automate Acceptance Tests

Acceptance tests must always be automated. Manual test plans are economically unsustainable.

- Manual testing costs grow with each execution
- Automated tests can run multiple times daily
- Manual tests get cut when budgets tighten, leaving gaps in coverage

**The cost comparison**: A manual test plan costing $1M+ per execution vs. one-time automation investment.

### 2. Define "Done" With Passing Tests

Done means done: all code written, all tests pass, QA and stakeholders accept.

- Create automated tests that define completion criteria
- When acceptance tests pass, the feature is deployable
- No ambiguous states like "done" vs "done-done"

### 3. Write Tests Late, But Before Implementation

Follow "late precision" - write acceptance tests just before implementing the feature.

- Tests should be ready by the first day of iteration/sprint
- All tests ready by iteration midpoint
- If tests aren't ready, developers help complete them

### 4. Separate Test Authors From Implementers

The developer who writes the test should not implement the tested feature.

- Business analysts write happy path tests (business value)
- QA writes unhappy path tests (edge cases, exceptions)
- Developers review for consistency and feasibility

### 5. Never Be Passive-Aggressive With Tests

If a test doesn't make sense, negotiate - don't blindly implement.

- Tests may be too complicated, awkward, or wrong
- Professionals help the team create the best software
- "That's what the test says" is not a valid excuse

**Bad approach**:
```
// Test says 2-second guarantee, so I'll implement exactly that
// even though it's statistically impossible
```

**Good approach**:
```
// Negotiate with test author for realistic criteria
// e.g., "99.5% of requests complete in 2 seconds"
```

### 6. Test Through APIs, Not GUIs

Write business rule tests through an API below the GUI.

- GUIs change frequently, making tests fragile
- GUI changes shouldn't break business rule tests
- Keep GUI-specific tests minimal
- Decouple GUI from business rules for testing

### 7. Keep CI Tests Always Passing

Broken tests in continuous integration are emergencies.

- Run all tests on every commit
- Failed tests = "stop the presses"
- Never disable failing tests to meet deadlines
- Never remove tests from the build

## Guidelines

Less strict recommendations for effective acceptance testing:

- Use unique IDs for GUI elements rather than positional selectors
- Make intermediate calculations visible in test reports for verification
- Create reusable scenarios and fixtures across tests
- Include error bars in estimates when requirements are imprecise
- Involve testers early in requirements discussions

## Exceptions

When rules may be relaxed:

- **Manual testing**: Acceptable for exploratory testing and aesthetics, but not for acceptance criteria
- **Developer-written tests**: When stakeholders lack time, developers may write tests - but not for features they implement
- **GUI testing**: Necessary when testing GUI behavior specifically, but use stubs for business rules

## Quick Reference

| Rule | Summary |
|------|---------|
| Automate | Never use manual acceptance tests |
| Define Done | Passing tests = feature complete |
| Late Precision | Write tests just before implementation |
| Separate Authors | Test writer != implementer |
| Negotiate Tests | Push back on problematic tests |
| API Over GUI | Test business rules below UI layer |
| CI Always Green | Broken build = emergency |

## Acceptance Tests vs Unit Tests

| Aspect | Acceptance Tests | Unit Tests |
|--------|-----------------|------------|
| Written by | Business/QA/BAs | Programmers |
| Written for | Business + Programmers | Programmers |
| Tests through | API or UI level | Method calls |
| Documents | System behavior | Code structure |
| Primary purpose | Specification | Design documentation |

Both test similar things through different pathways - neither is redundant.
