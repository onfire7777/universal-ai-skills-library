# Functions Rules

Guidelines for writing clean, readable, and maintainable functions.

## Core Rules

### 1. Keep Functions Small

Functions should be very small, ideally 2-5 lines, rarely exceeding 20.

- Blocks in `if`/`else`/`while` should be one line (a function call)
- Maximum indent level of one or two
- Each function should be transparently obvious

**Example**:
```typescript
// Bad - too long, multiple levels
const processOrder = (order: Order): void => {
  if (order.isValid()) {
    for (const item of order.items) {
      if (item.inStock) {
        inventory.reserve(item);
        // ... 20 more lines
      }
    }
  }
};

// Good - small, delegating
const processOrder = (order: Order): void => {
  if (order.isValid()) {
    reserveItems(order);
    calculateTotals(order);
    notifyCustomer(order);
  }
};
```

### 2. Do One Thing

Functions should do one thing, do it well, and do it only.

- All steps should be one level of abstraction below the function name
- If you can extract a meaningful function, it's doing too much
- If it has sections, it's doing too much

### 3. One Level of Abstraction

All statements in a function should be at the same abstraction level.

- Don't mix `getUser()` with `str.toLowerCase()`
- High-level calls with high-level, low-level with low-level

### 4. Minimize Arguments

Zero arguments is ideal. Three is the maximum.

- Zero (niladic): Best
- One (monadic): Good for questions or transformations
- Two (dyadic): Acceptable, but harder to understand
- Three (triadic): Avoid - very hard to understand
- More: Never

**Example**:
```typescript
// Bad - too many arguments
const createUser = (
  name: string,
  email: string,
  age: number,
  address: string,
  phone: string
): User => { /* ... */ };

// Good - use an object
interface CreateUserParams {
  name: string;
  email: string;
  age: number;
  address: string;
  phone: string;
}

const createUser = (params: CreateUserParams): User => { /* ... */ };
```

### 5. No Flag Arguments

Never pass a boolean to control function behavior.

- Flag arguments mean the function does more than one thing
- Split into two functions instead

**Example**:
```typescript
// Bad
const render = (data: Data, isSuite: boolean): string => { /* ... */ };

// Good
const renderSuite = (data: Data): string => { /* ... */ };
const renderSingleTest = (data: Data): string => { /* ... */ };
```

### 6. No Side Effects

Functions should not have hidden behaviors beyond their stated purpose.

- Don't modify global state unexpectedly
- Don't modify input parameters unexpectedly
- If side effects are necessary, make them explicit in the name

### 7. Command Query Separation

Functions should either do something OR answer something, never both.

- Commands: Change state, return nothing (or void)
- Queries: Return information, change nothing

**Example**:
```typescript
// Bad - does both
const set = (attr: string, value: string): boolean => { /* ... */ };
if (set("username", "bob")) { /* confusing! */ }

// Good - separated
const attributeExists = (attr: string): boolean => { /* ... */ };
const setAttribute = (attr: string, value: string): void => { /* ... */ };

if (attributeExists("username")) {
  setAttribute("username", "bob");
}
```

### 8. Use Exceptions, Not Error Codes

Throw exceptions instead of returning error codes.

- Error codes force immediate handling, causing deep nesting
- Exceptions separate happy path from error handling
- Error code enums become dependency magnets

### 9. Extract Try/Catch Blocks

Error handling is one thing - functions that handle errors should do nothing else.

- If `try` exists, it should be the first word in the function
- Nothing should come after `catch`/`finally` blocks
- Extract the try body and catch body into separate functions

### 10. Use Descriptive Names

Long descriptive names are better than short cryptic ones.

- Names should say what the function does
- Be consistent: use same phrases, nouns, verbs
- Spend time choosing names - it clarifies design

## Guidelines

- Use verb/noun pairs for monadic functions: `writeField(name)`
- Encode argument names in function name: `assertExpectedEqualsActual(expected, actual)`
- Avoid output arguments - use return values or `this`
- Multiple `return`/`break`/`continue` are fine in small functions

## Exceptions

- **Switch statements**: Acceptable if used once, to create polymorphic objects, hidden behind a factory
- **Multiple returns**: Fine in small functions where it improves clarity
- **Two arguments**: Acceptable for ordered pairs like `Point(x, y)` or natural pairs

## Quick Reference

| Rule | Guideline |
|------|-----------|
| Size | 2-5 lines ideal, max 20 |
| Arguments | 0-2 preferred, max 3 |
| Flag args | Never use |
| Side effects | Make explicit or eliminate |
| Command/Query | Separate state changes from returns |
| Error handling | Use exceptions, extract try/catch |
| Naming | Long and descriptive beats short |
