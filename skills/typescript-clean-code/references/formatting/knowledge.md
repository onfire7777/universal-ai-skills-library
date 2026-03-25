# Formatting Knowledge

Core concepts and foundational understanding for code formatting.

## Overview

Code formatting is about communication, not aesthetics. The readability of your code affects all future changes, long after the original functionality has been modified. Good formatting makes code feel professional and establishes trust with readers.

## Key Concepts

### The Communication Principle

**Definition**: Code formatting is primarily about communicating intent to other developers.

Functionality changes frequently, but readability affects maintainability forever. Your coding style and discipline survive even when your original code does not.

**Key points**:
- Formatting is too important to ignore
- Formatting is too important to treat religiously
- Professional developers prioritize communication

### Vertical Formatting

**Definition**: How code is organized from top to bottom within a file.

Controls file size, concept separation, and the reading flow through your code.

**Key points**:
- Files should typically be 200-500 lines
- High-level concepts at top, details at bottom
- Related concepts should be vertically close

### Horizontal Formatting

**Definition**: How code is organized from left to right within a line.

Controls line length, whitespace usage, and visual grouping of related elements.

**Key points**:
- Keep lines under 100-120 characters
- Use whitespace to show relationships
- Indentation reveals code hierarchy

### The Newspaper Metaphor

**Definition**: Source files should read like newspaper articles - headline first, then synopsis, then increasing detail.

The file name is the headline. Top-level code provides the overview. Details and low-level functions come last.

## Terminology

| Term | Definition |
|------|------------|
| Vertical openness | Blank lines that separate distinct concepts |
| Vertical density | Tightly related code appearing close together |
| Vertical distance | Physical separation between related concepts |
| Conceptual affinity | Degree to which code elements belong together |
| Horizontal density | Elements without spaces indicating close relationship |

## How It Relates To

- **Functions**: Small functions enable proper vertical formatting
- **Naming**: Good names support the newspaper metaphor
- **Classes**: Class structure determines file organization
- **Comments**: Excessive comments disrupt vertical density

## Common Misconceptions

- **Myth**: Formatting is just about aesthetics
  **Reality**: Formatting directly impacts code comprehension and maintenance

- **Myth**: Personal style preferences should prevail
  **Reality**: Team consistency trumps individual preferences

- **Myth**: Modern IDEs eliminate formatting concerns
  **Reality**: Tools help enforce rules, but humans must choose good rules

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Purpose | Formatting is communication, not decoration |
| File size | Target 200 lines, rarely exceed 500 |
| Line width | Stay under 100-120 characters |
| Newspaper metaphor | High-level first, details last |
| Team rules | Consistency matters more than personal preference |
