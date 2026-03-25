# Comments Knowledge

Core concepts and foundational understanding for code comments.

## Overview

Comments are, at best, a necessary evil used to compensate for our failure to express intent in code. The proper use of comments is to explain what cannot be made clear through code alone. Truth can only be found in the code itself.

## Key Concepts

### Comments as Failure

**Definition**: Every comment represents a failure to express yourself clearly in code.

When you write a comment, it means the code wasn't expressive enough on its own. The goal should be to minimize comments by writing clearer, more self-documenting code.

**Key points**:
- Comments should trigger reflection: "Can I express this in code instead?"
- Energy spent on comments is better spent improving code clarity
- Success is measured by how few comments are needed

### Comment Decay

**Definition**: Comments become increasingly inaccurate over time as code evolves.

Code changes and moves, but comments often don't follow. Comments become "orphaned blurbs of ever-decreasing accuracy."

**Key points**:
- Comments lie - not intentionally, but inevitably
- The older and farther from code, the more likely to be wrong
- Inaccurate comments are worse than no comments

### Code as Truth

**Definition**: Only the code truly tells you what the system does.

The code is the only source of accurate information. Comments can mislead, but executable code cannot lie about what it actually does.

**Key points**:
- Trust the code over the comment when they conflict
- Make the code tell the truth clearly
- Use comments only when code cannot express intent

## Terminology

| Term | Definition |
|------|------------|
| Self-documenting code | Code that explains itself through clear naming and structure |
| Comment decay | The process by which comments become inaccurate over time |
| Noise comment | A comment that adds no information beyond what code already shows |
| Orphaned comment | A comment separated from the code it was meant to describe |

## How It Relates To

- **Naming**: Good names eliminate need for explanatory comments
- **Functions**: Small, well-named functions replace comment blocks
- **Refactoring**: Often the solution to "needing" a comment is refactoring

## Common Misconceptions

- **Myth**: More comments mean better documented code
  **Reality**: Comments often indicate code that needs to be clearer

- **Myth**: Every function needs a doc comment
  **Reality**: Good function names make most doc comments redundant

- **Myth**: Comments help future maintainers
  **Reality**: Outdated comments actively mislead maintainers

## The Comment Test

Before writing a comment, ask:

1. Can I rename a variable to make this clear?
2. Can I extract a well-named function?
3. Can I restructure the code to be self-explanatory?
4. Is this comment explaining WHAT (bad) or WHY (potentially good)?

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Comments as failure | Every comment means code wasn't clear enough |
| Comment decay | Comments become lies over time |
| Code as truth | Only executable code is reliable documentation |
| Self-documenting | The goal is code that needs no comments |
