# Classes Rules

Guidelines for designing clean, maintainable classes and modules in TypeScript.

## Core Rules

### 1. Classes Should Be Small (Measured by Responsibilities)

Size is measured by counting responsibilities, not lines of code.

- A class should have ONE reason to change
- If the class name needs weasel words (`Manager`, `Processor`, `Handler`, `Super`), it's too big
- You should describe the class in ~25 words without "if," "and," "or," or "but"

**Example**:
```typescript
// Bad - multiple responsibilities
class SuperDashboard {
  getLastFocusedComponent(): Component { }
  setLastFocused(component: Component): void { }
  getMajorVersionNumber(): number { }
  getMinorVersionNumber(): number { }
  getBuildNumber(): number { }
}

// Good - single responsibility extracted
class Version {
  getMajorVersionNumber(): number { }
  getMinorVersionNumber(): number { }
  getBuildNumber(): number { }
}
```

### 2. Follow the Single Responsibility Principle (SRP)

Every class should have one, and only one, reason to change.

- Identify responsibilities by asking "what would cause this class to change?"
- Different reasons to change = different classes
- Extract responsibilities into their own classes

### 3. Maintain High Cohesion

Methods should use the class's instance variables.

- Each method should manipulate one or more instance variables
- If a subset of variables is only used by a subset of methods, consider splitting
- When cohesion drops after extracting functions, split the class

**Example**:
```typescript
// Good - highly cohesive
class Stack<T> {
  private topOfStack = 0;
  private elements: T[] = [];

  size(): number {
    return this.topOfStack;
  }

  push(element: T): void {
    this.topOfStack++;
    this.elements.push(element);
  }

  pop(): T {
    if (this.topOfStack === 0) throw new Error('Stack is empty');
    this.topOfStack--;
    return this.elements.pop()!;
  }
}
```

### 4. Organize for Change (Open-Closed Principle)

Structure classes to minimize modification when adding features.

- Prefer extension over modification
- New features should be new classes, not changes to existing ones
- Private methods used by only one public method may indicate extraction opportunity

### 5. Isolate from Change (Dependency Inversion)

Depend on abstractions, not concrete implementations.

- Use interfaces to decouple from external dependencies
- Inject dependencies through constructors
- Enables testing with mocks/stubs

**Example**:
```typescript
// Bad - depends on concrete class
class Portfolio {
  private exchange = new TokyoStockExchange();
  
  getValue(): Money { }
}

// Good - depends on abstraction
interface StockExchange {
  currentPrice(symbol: string): Money;
}

class Portfolio {
  constructor(private exchange: StockExchange) { }
  
  getValue(): Money { }
}
```

### 6. Follow Standard Class Organization

Order elements consistently within a class.

- Public static constants first
- Private static variables
- Private instance variables
- Public methods
- Private methods (after the public method that calls them)

## Guidelines

Less strict recommendations:

- Prefer many small, single-purpose classes over few large ones
- When breaking up functions creates shared variables, consider if those variables define a new class
- Tests can justify relaxing encapsulation (protected access), but try other approaches first
- A class that is "logically complete" with no anticipated changes can stay as-is

## Exceptions

When these rules may be relaxed:

- **Testing**: Make methods/properties protected when tests need access (last resort)
- **Stable code**: Don't preemptively split a class if no changes are anticipated
- **Simple utilities**: Pure utility functions may stay together if truly cohesive

## Quick Reference

| Rule | Summary |
|------|---------|
| Small classes | Measure by responsibilities, not lines |
| SRP | One reason to change |
| Cohesion | Methods share instance variables |
| OCP | Extend, don't modify |
| DIP | Depend on abstractions |
| Organization | Constants > variables > public > private |
