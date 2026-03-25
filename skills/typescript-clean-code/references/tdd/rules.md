# TDD Rules

The three laws of TDD and guidelines for professional practice.

## Core Rules

### 1. First Law: No Production Code Without a Failing Test

You are not allowed to write any production code until you have first written a failing unit test.

- Start every feature or fix with a test
- The test defines what you're about to build
- Even a compile error counts as a failing test

### 2. Second Law: Minimal Failing Test

You are not allowed to write more of a unit test than is sufficient to fail - and not compiling is failing.

- Write just enough test to fail
- A compilation failure is a valid failure
- Stop writing test code as soon as it fails to compile or run

### 3. Third Law: Minimal Production Code

You are not allowed to write more production code than is sufficient to pass the currently failing unit test.

- Write only enough code to make the test pass
- Don't anticipate future needs
- Stop as soon as the test goes green

## The TDD Cycle

The three laws lock you into a ~30 second cycle:

1. **Red**: Write a small failing test
2. **Green**: Write minimal code to pass
3. **Refactor**: Clean up while tests stay green
4. Repeat

**Key insight**: The test code and production code grow simultaneously into complementary components - "like an antibody fits an antigen."

## Guidelines

### The Professional Option

TDD enhances:
- **Certainty** - Know your code works
- **Courage** - Fearlessly improve code
- **Defect reduction** - 2X-10X fewer bugs
- **Documentation** - Tests describe the system
- **Design** - Forces decoupling

> It could be considered *unprofessional* not to use TDD.

### Rapid Cycle Time

- Aim for ~30 second cycles
- Run tests after every small change
- Don't go more than a few minutes without green tests

### Test Coverage Through TDD

- Following the three laws naturally produces high coverage
- Every object creation method gets tested
- Every meaningful function call gets tested
- No production code exists without a corresponding test

## Exceptions

When TDD may not apply:

- **Impractical situations**: Rare cases where the discipline does more harm than good
- **Exploratory spikes**: Sometimes you need to explore before you test
- **Legacy code**: May need different strategies for untested codebases

> No professional developer should ever follow a discipline when that discipline does more harm than good.

## What TDD Is Not

- **Not a religion**: It's a practical discipline, not dogma
- **Not a magic formula**: Following the laws doesn't guarantee good code
- **Not a guarantee**: You can still write bad code and bad tests
- **Not always applicable**: Use judgment about when it applies

## Quick Reference

| Rule | Summary |
|------|---------|
| Law 1 | No production code without a failing test first |
| Law 2 | Write only enough test to fail (compile failure counts) |
| Law 3 | Write only enough code to pass the failing test |
| Cycle Time | ~30 seconds between test runs |
| Coverage | Emerges naturally from following the three laws |
| Professional | TDD is the professional option, not an optional extra |
