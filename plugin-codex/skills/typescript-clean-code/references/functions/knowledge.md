# Functions Knowledge

Core concepts and foundational understanding for writing clean functions.

## Overview

Functions are the first line of organization in any program. Well-written functions tell a story, are easy to read, and communicate intent clearly. The goal is to make code read like a top-down narrative, with each function leading naturally to the next.

## Key Concepts

### Small Functions

**Definition**: Functions should be very small, ideally 2-5 lines, rarely exceeding 20 lines.

Small functions are easier to read, understand, and test. Each function should be "transparently obvious" and tell a story that leads to the next function.

**Key points**:
- Blocks in `if`/`else`/`while` should be one line (a function call)
- Indent level should not exceed one or two
- If it's hard to shrink further, you've reached the right size

### Do One Thing

**Definition**: A function should do one thing, do it well, and do it only.

A function does "one thing" if you can describe it as a brief TO paragraph and all steps are one level of abstraction below the function name.

**Key points**:
- If you can extract another function with a meaningful (non-restating) name, it does more than one thing
- Functions that can be divided into sections are doing more than one thing

### Abstraction Levels

**Definition**: All statements within a function should be at the same level of abstraction.

Mixing high-level concepts (`getHtml()`) with low-level details (`.append("\n")`) is confusing. Readers can't tell what's essential vs. detail.

**Key points**:
- High level: business logic, domain operations
- Intermediate level: utility operations
- Low level: string manipulation, data formatting

### The Stepdown Rule

**Definition**: Code should read like a top-down narrative, with each function followed by those at the next level of abstraction.

Write code as a set of TO paragraphs, each describing the current level and referencing the next level down.

## Terminology

| Term | Definition |
|------|------------|
| Niladic | Function with zero arguments (ideal) |
| Monadic | Function with one argument |
| Dyadic | Function with two arguments |
| Triadic | Function with three arguments (avoid) |
| Polyadic | Function with more than three arguments (never) |
| Side Effect | Hidden action beyond the function's stated purpose |
| Temporal Coupling | When a function can only be called at certain times |
| Command | Function that changes state |
| Query | Function that returns information |

## How It Relates To

- **Naming**: Small, focused functions are easier to name descriptively
- **Testing**: Fewer arguments and single responsibility make testing simpler
- **Error Handling**: Separate error handling from business logic using exceptions
- **DRY Principle**: Extract duplicated code into well-named functions

## Common Misconceptions

- **Myth**: Functions need comments to explain what they do
  **Reality**: A well-named small function is self-documenting

- **Myth**: Extracting small functions hurts performance
  **Reality**: Readability matters more; compilers optimize well

- **Myth**: Functions should have one return statement
  **Reality**: Multiple returns are fine in small functions

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Small | 2-5 lines ideal, max 20 lines |
| Do One Thing | One level of abstraction below the function name |
| Abstraction | All statements at same level |
| Stepdown | Code reads top-down like TO paragraphs |
| Arguments | Zero is best, three is maximum |
| Side Effects | Functions should not have hidden behaviors |
| Command/Query | Do something OR answer something, not both |
