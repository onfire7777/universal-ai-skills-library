# Naming Knowledge

Core concepts and foundational understanding for meaningful names in code.

## Overview

Names are everywhere in software: variables, functions, arguments, classes, modules, files, and directories. Good names make code readable without comments. The investment in choosing good names pays dividends in maintainability and comprehension.

## Key Concepts

### Intention-Revealing Names

**Definition**: Names that answer why something exists, what it does, and how it's used.

A name should make comments unnecessary. If you need a comment to explain what a variable is, the name doesn't reveal its intent.

**Key points**:
- Names should answer the "big questions" about the thing being named
- Good names eliminate the need for explanatory comments
- Take time choosing names; change them when you find better ones

### Meaningful Distinctions

**Definition**: Names that differ in ways that convey actual semantic differences.

When names must be different, they should mean something different. Number-series naming (a1, a2) and noise words (Info, Data) create distinctions without meaning.

**Key points**:
- Avoid arbitrary differences just to satisfy the compiler
- Noise words like Info, Data, Object add no meaning
- If you can't tell the difference between two names, readers can't either

### Searchability

**Definition**: Names that can be easily located across a codebase.

Single-letter names and magic numbers are nearly impossible to search for. Longer, descriptive names enable effective code navigation.

**Key points**:
- Name length should correspond to scope size
- Constants should have searchable names
- Single letters acceptable only in tiny local scopes (loop counters)

### Mental Mapping

**Definition**: The cognitive load required to translate a name into its actual meaning.

Readers shouldn't have to mentally map your names to concepts they already know. Professional programmers prioritize clarity over cleverness.

**Key points**:
- Clarity is king
- Avoid forcing readers to remember what abbreviations mean
- Use domain-appropriate terminology directly

## Terminology

| Term | Definition |
|------|------------|
| Disinformation | Names that lie or mislead about what something is |
| Noise words | Meaningless additions like Info, Data, Object, String |
| Hungarian Notation | Encoding type information in name prefixes (obsolete) |
| Solution domain | Technical/CS terms programmers already know |
| Problem domain | Business/domain-specific terminology |

## How It Relates To

- **Readability**: Names are the primary vehicle for code comprehension
- **Comments**: Good names reduce or eliminate need for comments
- **Refactoring**: Renaming is a fundamental refactoring operation
- **Communication**: Code is read more than written; names enable team communication

## Common Misconceptions

- **Myth**: Short names are always better for brevity
  **Reality**: Clarity trumps brevity; name length should match scope size

- **Myth**: Adding prefixes like `m_` or `I` adds useful information
  **Reality**: Modern IDEs make these encodings unnecessary clutter

- **Myth**: Consistent word choice means using the same word everywhere
  **Reality**: Different concepts deserve different words even if similar

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Intention-revealing | Name answers why, what, and how |
| Meaningful distinctions | Different names = different meanings |
| Searchability | Longer names for larger scopes |
| No mental mapping | Reader shouldn't need to translate |
| Solution domain | Use CS terms for technical concepts |
| Problem domain | Use business terms for domain concepts |
