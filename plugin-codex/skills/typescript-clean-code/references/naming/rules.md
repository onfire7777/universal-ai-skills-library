# Naming Rules

Specific guidelines for choosing meaningful names in code.

## Core Rules

### 1. Use Intention-Revealing Names

Names should tell you why something exists, what it does, and how it's used.

- A name requiring a comment doesn't reveal its intent
- Include units of measurement when relevant

```typescript
// Bad
const d = 0; // elapsed time in days

// Good
const elapsedTimeInDays = 0;
```

### 2. Avoid Disinformation

Don't use names that lie about what something is.

- Don't use `list` unless it's actually a List type
- Avoid names that look too similar to each other
- Never use `l` (lowercase L) or `O` as variable names

### 3. Make Meaningful Distinctions

If names must be different, they should mean something different.

- Avoid number-series naming (a1, a2, a3)
- Avoid noise words (Info, Data, Object, String)

```typescript
// Bad: a1, a2 mean nothing
function copyChars(a1: string[], a2: string[]): void { }

// Good: source, destination are meaningful
function copyChars(source: string[], destination: string[]): void { }
```

### 4. Use Pronounceable Names

If you can't say it, you can't discuss it.

- Programming is a social activity
- New developers shouldn't need names explained

### 5. Use Searchable Names

Single letters and magic numbers can't be found in a codebase.

- Name length should correspond to scope size
- Single letters only in small local scopes
- Constants must have descriptive names

```typescript
// Bad: can't search for 5
s += (t[j] * 4) / 5;

// Good: WORK_DAYS_PER_WEEK is searchable
const realTaskWeeks = realTaskDays / WORK_DAYS_PER_WEEK;
```

### 6. Avoid Encodings

Don't encode type or scope information into names.

- No Hungarian Notation (strName, intCount)
- No member prefixes (m_description, _name)
- Leave interfaces unadorned; encode implementations if needed (ShapeFactoryImpl)

### 7. Class Names Should Be Nouns

- Good: Customer, WikiPage, Account, AddressParser
- Avoid: Manager, Processor, Data, Info
- Never use verbs for class names

### 8. Method Names Should Be Verbs

- Accessors: `get` prefix (getName)
- Mutators: `set` prefix (setName)
- Predicates: `is`/`has` prefix (isPosted, hasPermission)
- Use static factory methods over overloaded constructors

### 9. Pick One Word Per Concept

- Don't mix fetch, retrieve, and get for the same concept
- Don't mix controller, manager, and driver arbitrarily
- Consistent lexicon helps programmers navigate code

### 10. Don't Pun

Don't use the same word for different concepts.

- If `add` concatenates values, don't use it for collection insertion
- Use `insert` or `append` when semantics differ

### 11. Add Meaningful Context

Place names in context through enclosing structures.

- Variables alone may be ambiguous (`state` could be anything)
- Group related variables in classes
- Use prefixes only as a last resort

### 12. Don't Add Gratuitous Context

Shorter names are better if they're clear.

- Don't prefix every class with application name
- Add only necessary context

## Guidelines

- Use solution domain names (CS terms) for technical concepts
- Use problem domain names for business concepts
- Rename fearlessly; modern tools make it safe

## Exceptions

- **Loop counters**: `i`, `j`, `k` acceptable in small loops
- **Lambda parameters**: Short names OK in tiny scopes
- **Well-known abbreviations**: `id`, `url`, `html` are acceptable

## Quick Reference

| Rule | Summary |
|------|---------|
| Intention-revealing | Names answer why, what, how |
| No disinformation | Don't lie about types or purpose |
| Meaningful distinctions | Different names = different meanings |
| Pronounceable | You can discuss it out loud |
| Searchable | Can be found with grep/search |
| No encodings | No type prefixes or Hungarian Notation |
| Nouns for classes | Customer, not Manage |
| Verbs for methods | calculatePay, not payment |
| One word per concept | Consistent vocabulary |
| No puns | Same word = same concept |
