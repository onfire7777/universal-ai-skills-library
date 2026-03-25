# Functions Checklist

Use when writing new functions or reviewing existing code.

## Size and Structure

- [ ] Function is 20 lines or less
- [ ] Indent level is 2 or less
- [ ] Blocks in `if`/`else`/`while` are one line (function calls)
- [ ] No sections or comments dividing the function

## Single Responsibility

- [ ] Can describe what function does in one sentence without "and"
- [ ] Cannot extract another function with a meaningful name
- [ ] All statements are at the same abstraction level
- [ ] Function does not have hidden side effects

## Arguments

- [ ] Has 2 or fewer arguments (3 max in rare cases)
- [ ] No boolean/flag arguments
- [ ] No output arguments (use return values instead)
- [ ] Related arguments grouped into objects

## Naming

- [ ] Name describes what the function does
- [ ] Verb/noun pair for monadic functions
- [ ] Consistent with similar functions in the codebase
- [ ] Long enough to be clear (don't abbreviate)

## Error Handling

- [ ] Uses exceptions, not error codes
- [ ] Try/catch bodies extracted to separate functions
- [ ] Error handling function does nothing else
- [ ] No error code enums

## Command Query Separation

- [ ] Function either changes state OR returns information, not both
- [ ] If returning boolean, it's asking a question, not performing action
- [ ] State changes don't return status codes

## Red Flags

Stop and refactor if you find:

- Function longer than 20 lines
- More than 3 arguments
- Boolean parameters controlling behavior
- Mixed abstraction levels in one function
- Same switch statement in multiple places
- Function name includes "And" or "Or"
- Side effects hidden in function body
- Deep nesting (3+ levels)

## Quick Reference

| Aspect | Ideal | Acceptable | Red Flag |
|--------|-------|------------|----------|
| Lines | 2-5 | 6-20 | >20 |
| Arguments | 0-1 | 2 | 3+ |
| Indent depth | 1 | 2 | 3+ |
| Abstraction levels | 1 | 1 | 2+ mixed |
| Side effects | 0 | Named | Hidden |
| Returns per function | 1 | 2-3 small fn | Many in large fn |
