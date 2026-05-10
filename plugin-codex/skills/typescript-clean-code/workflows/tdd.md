# TDD Workflow

Red-Green-Refactor cycle for test-driven development.

## When to Use

- Writing new functionality
- Adding features to existing code
- Fixing bugs (write failing test first)

## Prerequisites

- Test framework set up (Jest, Vitest, etc.)
- Understanding of the feature to implement

**Reference**: `tdd/rules.md`, `tdd/examples.md`

---

## The Three Laws of TDD

1. **Write no production code** until you have a failing test
2. **Write only enough test** to demonstrate a failure
3. **Write only enough production code** to pass the test

**Cycle time**: ~30 seconds per iteration

---

## Workflow Steps

### Step 1: Think

**Goal**: Understand what you're building before writing any code.

- [ ] What is the next smallest piece of functionality?
- [ ] What behavior should it have?
- [ ] What is the simplest test that would force this behavior?

**Ask**: "What's the smallest step I can take?"

---

### Step 2: RED - Write a Failing Test

**Goal**: Write a test that fails for the right reason.

```typescript
// Example: Building a Stack
describe('Stack', () => {
  it('should be empty when created', () => {
    const stack = new Stack();
    expect(stack.isEmpty()).toBe(true);
  });
});
```

**Checklist**:
- [ ] Test is small and focused on ONE behavior
- [ ] Test name describes the expected behavior
- [ ] Test follows Arrange-Act-Assert pattern
- [ ] Run the test - it MUST fail
- [ ] Failure is for the right reason (not syntax error)

**Common Mistakes**:
- Writing too much test at once
- Test passes immediately (you wrote too much production code)
- Test fails for wrong reason (compilation error, not assertion)

---

### Step 3: GREEN - Make It Pass

**Goal**: Write the minimum code to make the test pass.

```typescript
// Minimum code to pass
class Stack {
  isEmpty(): boolean {
    return true;  // Simplest thing that works
  }
}
```

**Checklist**:
- [ ] Write the simplest code that passes
- [ ] Don't write code for future tests
- [ ] Don't refactor yet
- [ ] Run the test - it MUST pass
- [ ] All previous tests still pass

**Common Mistakes**:
- Writing more code than needed
- Generalizing before you have multiple examples
- Refactoring during this phase

---

### Step 4: REFACTOR - Clean Up

**Goal**: Improve code quality while keeping tests green.

**Checklist**:
- [ ] Remove duplication
- [ ] Improve names
- [ ] Extract methods if needed
- [ ] Improve readability
- [ ] Run tests after EVERY change
- [ ] All tests still pass

**Reference**: `refactoring.md` workflow

**Common Mistakes**:
- Refactoring without running tests
- Adding new functionality (that's Step 2)
- Skipping this step

---

### Step 5: Repeat

**Goal**: Continue the cycle until feature is complete.

```
Think → RED → GREEN → REFACTOR → Think → RED → GREEN → REFACTOR → ...
```

**Next Iteration Example**:
```typescript
// Next test
it('should not be empty after push', () => {
  const stack = new Stack();
  stack.push(1);
  expect(stack.isEmpty()).toBe(false);
});

// Make it pass
class Stack {
  private items: number[] = [];
  
  isEmpty(): boolean {
    return this.items.length === 0;
  }
  
  push(item: number): void {
    this.items.push(item);
  }
}
```

---

## TDD Cycle Timing

| Phase | Time | Activity |
|-------|------|----------|
| RED | ~1 min | Write failing test |
| GREEN | ~2 min | Make it pass |
| REFACTOR | ~2 min | Clean up |
| **Total** | ~5 min | One cycle |

If a cycle takes longer than 10 minutes, you're taking too big a step.

---

## Test Progression Strategy

Build up complexity gradually:

1. **Degenerate cases**: Empty, null, zero
2. **Single item**: One element scenarios
3. **Multiple items**: Normal operation
4. **Edge cases**: Boundaries, errors
5. **Error cases**: Invalid inputs, exceptions

**Example for Stack**:
```
1. isEmpty() returns true for new stack
2. isEmpty() returns false after push
3. pop() returns pushed item
4. pop() returns items in LIFO order
5. pop() on empty stack throws
```

---

## Quick Reference

```
┌─────────────────────────────────────────────────┐
│                  TDD CYCLE                       │
│                                                  │
│    ┌───────┐                                    │
│    │ THINK │ ← What's the smallest step?        │
│    └───┬───┘                                    │
│        ▼                                        │
│    ┌───────┐                                    │
│    │  RED  │ ← Write failing test               │
│    └───┬───┘                                    │
│        ▼                                        │
│    ┌───────┐                                    │
│    │ GREEN │ ← Minimum code to pass             │
│    └───┬───┘                                    │
│        ▼                                        │
│    ┌──────────┐                                 │
│    │ REFACTOR │ ← Clean up, tests stay green    │
│    └────┬─────┘                                 │
│         │                                       │
│         └──────────► Repeat                     │
│                                                  │
└─────────────────────────────────────────────────┘
```

---

## Exit Criteria

Feature is complete when:
- [ ] All acceptance criteria have tests
- [ ] All tests pass
- [ ] Code is clean (no obvious refactoring needed)
- [ ] You can't think of any more tests that would fail
