# TDD Examples

Code examples demonstrating the TDD red-green-refactor cycle in TypeScript.

## The TDD Cycle in Practice

### Example: Building a Stack

We'll build a simple Stack class following the three laws of TDD.

#### Cycle 1: Empty Stack

**Red** - Write a failing test:

```typescript
// stack.test.ts
import { Stack } from './stack';

describe('Stack', () => {
  it('should be empty when created', () => {
    const stack = new Stack<number>();
    expect(stack.isEmpty()).toBe(true);
  });
});
```

**Green** - Write minimal code to pass:

```typescript
// stack.ts
export class Stack<T> {
  isEmpty(): boolean {
    return true;
  }
}
```

**Why it works**:
- Test fails first (Stack doesn't exist)
- We write just enough to compile and pass
- No unnecessary code

#### Cycle 2: Push Makes Non-Empty

**Red** - Write the next failing test:

```typescript
it('should not be empty after push', () => {
  const stack = new Stack<number>();
  stack.push(1);
  expect(stack.isEmpty()).toBe(false);
});
```

**Green** - Make it pass:

```typescript
export class Stack<T> {
  private items: T[] = [];

  isEmpty(): boolean {
    return this.items.length === 0;
  }

  push(item: T): void {
    this.items.push(item);
  }
}
```

#### Cycle 3: Pop Returns Pushed Item

**Red**:

```typescript
it('should pop the pushed item', () => {
  const stack = new Stack<number>();
  stack.push(42);
  expect(stack.pop()).toBe(42);
});
```

**Green**:

```typescript
export class Stack<T> {
  private items: T[] = [];

  isEmpty(): boolean {
    return this.items.length === 0;
  }

  push(item: T): void {
    this.items.push(item);
  }

  pop(): T | undefined {
    return this.items.pop();
  }
}
```

#### Cycle 4: LIFO Order

**Red**:

```typescript
it('should pop in LIFO order', () => {
  const stack = new Stack<number>();
  stack.push(1);
  stack.push(2);
  expect(stack.pop()).toBe(2);
  expect(stack.pop()).toBe(1);
});
```

**Green**: Already passes! Our implementation naturally handles this.

**Refactor**: Tests are clean, code is clean. Nothing to refactor yet.

---

## Refactoring Example

### Before Refactoring

After several TDD cycles, we might have:

```typescript
export class UserService {
  createUser(name: string, email: string): User {
    if (!name || name.trim() === '') {
      throw new Error('Name required');
    }
    if (!email || !email.includes('@')) {
      throw new Error('Valid email required');
    }
    const user = { id: generateId(), name: name.trim(), email: email.toLowerCase() };
    saveToDatabase(user);
    return user;
  }
}
```

### After Refactoring

With tests protecting us, we safely extract validation methods:

```typescript
export class UserService {
  createUser(name: string, email: string): User {
    const validatedName = this.validateName(name);
    const validatedEmail = this.validateEmail(email);
    return this.buildUser(validatedName, validatedEmail);
  }

  private validateName(name: string): string { /* validation logic */ }
  private validateEmail(email: string): string { /* validation logic */ }
  private buildUser(name: string, email: string): User { /* create user */ }
}
```

**Key point**: Tests still pass - refactoring is safe

---

## Quick Reference

| Phase | Action | Duration |
|-------|--------|----------|
| Red | Write failing test | ~10 sec |
| Green | Write minimal code | ~15 sec |
| Refactor | Clean up | ~5 sec |
| **Total** | One cycle | ~30 sec |
