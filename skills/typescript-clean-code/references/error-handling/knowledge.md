# Error Handling Knowledge

Core concepts and foundational understanding for error handling in clean code.

## Overview

Error handling is essential but should not obscure business logic. Clean error handling separates error concerns from main algorithms, making code both readable and robust. The goal is graceful, stylish error management that enhances rather than clutters code.

## Key Concepts

### Separation of Concerns

**Definition**: Error handling logic should be separate from business logic.

When error handling dominates code, it becomes impossible to see what the code actually does. Error handling is important, but if it obscures logic, it's wrong.

**Key points**:
- Two concerns should not be tangled together
- You should be able to understand each concern independently
- Business logic should read like a clean, unadorned algorithm

### Exception-Based Error Handling

**Definition**: Using exceptions instead of return codes or error flags to signal errors.

Exceptions allow error handling at a distance, separating where errors occur from where they're handled. This keeps the happy path clean.

**Key points**:
- Return codes clutter the caller with immediate checks
- Exceptions let calling code focus on its primary purpose
- Missing a return code check is easy; exceptions can't be ignored

### Normal Flow Pattern

**Definition**: Designing code so that special cases don't require exception handling in business logic.

Instead of using try/catch for expected variations, design objects that handle special cases internally. This keeps business logic clean and linear.

**Key points**:
- Push error detection to the edges of your program
- Use Special Case objects for expected variations
- Business logic shouldn't deal with exceptional behavior directly

### Transaction Scope

**Definition**: A try block defines a scope where execution can abort and resume at catch.

Think of try blocks as transactions - the catch must leave the program in a consistent state regardless of what happens in the try block.

**Key points**:
- Start with try-catch-finally when writing code that could throw
- Define what users should expect regardless of failures
- Build up logic inside the try block after establishing the scope

## Terminology

| Term | Definition |
|------|------------|
| Exception | An object thrown to signal an error condition |
| Try-Catch-Finally | Structure that defines error handling scope |
| Special Case Pattern | Object that handles edge cases internally |
| Wrapper | Class that translates third-party exceptions to your own |
| Context | Information provided with exceptions for debugging |

## How It Relates To

- **Functions**: Error handling should be extracted into separate functions
- **Testing**: Write tests that force exceptions first (TDD approach)
- **Third-Party APIs**: Wrap external APIs to control exception types
- **Null Safety**: Avoiding null prevents many error conditions

## Common Misconceptions

- **Myth**: More null checks make code safer
  **Reality**: Too many null checks indicate a design problem; return empty collections or Special Case objects instead

- **Myth**: Exceptions should match their technical source
  **Reality**: Define exceptions based on how callers will catch them, not where they originate

- **Myth**: Each error type needs its own exception class
  **Reality**: Often a single exception class per area is sufficient; use different classes only when callers need to handle them differently

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Separation | Error handling separate from business logic |
| Exceptions | Prefer throwing over return codes |
| Try-First | Start with try-catch-finally structure |
| Normal Flow | Design to minimize exception handling in business code |
| Context | Exceptions should explain what failed and why |
| Wrapping | Translate third-party exceptions to your types |
| No Null | Don't return null; use empty collections or Special Case |
