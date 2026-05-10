# Classes Knowledge

Core concepts and foundational understanding for designing clean classes and modules.

## Overview

Classes are the higher-level organizational units of code. While functions handle blocks of code, classes organize functions and data into cohesive units with clear responsibilities. Clean classes are small (measured by responsibilities, not lines), highly cohesive, and loosely coupled.

## Key Concepts

### Class Organization

**Definition**: The standard ordering of elements within a class file.

Classes should follow a consistent structure for readability:

**Ordering (top to bottom)**:
- Public static constants
- Private static variables
- Private instance variables (rarely public)
- Public methods
- Private utility methods (placed after the public method that calls them)

This follows the "stepdown rule" - the class reads like a newspaper article.

### Encapsulation

**Definition**: Keeping variables and utility functions private, exposing only what's necessary.

**Key points**:
- Default to private for variables and helper methods
- Use protected only when tests require access (same package)
- Loosening encapsulation is always a last resort
- Tests can justify relaxing access, but first look for alternatives

### Single Responsibility Principle (SRP)

**Definition**: A class should have one, and only one, reason to change.

SRP is the most important concept in class design. A "responsibility" equals a "reason to change."

**Key points**:
- The class name should describe its single responsibility
- If you need "and," "or," "if," or "but" to describe it, it's too big
- Weasel words hint at multiple responsibilities: `Processor`, `Manager`, `Super`, `Handler`
- A 25-word description without conjunctions is a good test

### Cohesion

**Definition**: The degree to which methods use the class's instance variables.

A highly cohesive class has methods that work together with shared state.

**Key points**:
- Each method should manipulate one or more instance variables
- More shared variable usage = higher cohesion
- When cohesion drops, it signals time to split the class
- Ideal: variables exist because multiple methods need them together

### Open-Closed Principle (OCP)

**Definition**: Classes should be open for extension but closed for modification.

**Key points**:
- Add new functionality by creating new classes (subclasses)
- Existing classes should not need modification for new features
- Reduces risk when adding features to a system

### Dependency Inversion Principle (DIP)

**Definition**: Classes should depend upon abstractions, not concrete details.

**Key points**:
- Use interfaces to decouple from implementation details
- Makes testing easier (can substitute test doubles)
- Promotes flexibility and reuse
- Isolates classes from changes in dependencies

## Terminology

| Term | Definition |
|------|------------|
| Responsibility | A reason for a class to change |
| Cohesion | How closely methods relate via shared instance variables |
| Coupling | Degree of interdependence between classes |
| God Class | A class that knows/does too much (anti-pattern) |
| Abstraction | Interface or abstract class representing a concept |

## How It Relates To

- **Functions**: Classes organize functions; both should be small and focused
- **Modules**: TypeScript modules follow similar principles - single responsibility, cohesion
- **Testing**: Decoupled classes are easier to test in isolation
- **Refactoring**: Breaking large classes maintains cohesion as code evolves

## Common Misconceptions

- **Myth**: Many small classes are harder to understand than few large ones
  **Reality**: The total complexity is the same; small classes make it easier to find and understand relevant code

- **Myth**: A class with few methods is small enough
  **Reality**: Size is measured by responsibilities, not method count - 5 methods with 2 responsibilities is too big

- **Myth**: Encapsulation must never be broken
  **Reality**: Tests can justify protected access when no other option exists

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| SRP | One reason to change per class |
| Cohesion | Methods should share instance variables |
| OCP | Extend, don't modify |
| DIP | Depend on abstractions |
| Organization | Constants, variables, public methods, private methods |
