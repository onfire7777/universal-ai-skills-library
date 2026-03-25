# Comments Rules

Guidelines for when comments are acceptable and when to avoid them.

## Good Comments

### 1. Legal Comments

Copyright and license statements at file start are acceptable.

- Keep brief - reference external license files
- Let IDE collapse them to reduce clutter

```typescript
// Copyright (C) 2024 Company Inc. All rights reserved.
// Released under the terms of the MIT License.
```

### 2. Explanation of Intent

Explain WHY a decision was made, not WHAT the code does.

```typescript
// We sort our type higher than others to ensure consistent ordering
// when mixed with external types in collections
return 1;
```

### 3. Clarification

Translate obscure library returns or arguments when you can't modify them.

```typescript
expect(a.compareTo(a)).toBe(0);    // a == a
expect(a.compareTo(b)).toBe(-1);   // a < b
expect(b.compareTo(a)).toBe(1);    // b > a
```

**Warning**: Verify accuracy - clarifying comments are high-risk for errors.

### 4. Warning of Consequences

Alert other developers to non-obvious risks.

```typescript
// DateTimeFormat is not thread-safe, create new instance each call
function makeStandardDateFormat(): Intl.DateTimeFormat {
  return new Intl.DateTimeFormat('en-US', { /* options */ });
}
```

### 5. TODO Comments

Mark planned work with context, but clean them up regularly.

```typescript
// TODO: Remove after checkout model migration (Q2 2024)
function makeVersion(): VersionInfo | null {
  return null;
}
```

### 6. Amplification

Highlight importance of something that appears trivial.

```typescript
const listItemContent = match[3].trim();
// trim() is critical - leading spaces cause incorrect list detection
```

### 7. JSDoc for Public APIs

Document public APIs that others will consume.

- Keep concise and accurate
- Update when API changes
- Skip for internal/private code

## Bad Comments

### 1. Mumbling

Unclear comments that require reading other code to understand.

- If the comment needs explanation, rewrite it or delete it

### 2. Redundant Comments

Comments that restate what the code already says.

```typescript
// Bad - says nothing the code doesn't
/** Returns the day of the month */
getDayOfMonth(): number {
  return this.dayOfMonth;
}
```

### 3. Misleading Comments

Comments that are subtly inaccurate.

- Causes debugging nightmares when developers trust the comment over code

### 4. Mandated Comments

Requiring comments on every function/variable creates noise.

```typescript
// Bad - adds no value, potential for lies
/**
 * @param title The title of the CD
 * @param author The author of the CD
 */
function addCD(title: string, author: string): void { }
```

### 5. Journal Comments

Change logs at file top - use version control instead.

### 6. Noise Comments

Comments restating the obvious.

```typescript
// Bad
/** Default constructor */
constructor() { }

/** The day of the month */
private dayOfMonth: number;
```

### 7. Position Markers / Banners

```typescript
// Bad - use sparingly if at all
// //////////////// Actions ////////////////
```

### 8. Closing Brace Comments

If you need these, your function is too long.

```typescript
// Bad
} // end while
} // end try
```

### 9. Attributions

Use version control, not comments like `// Added by Rick`.

### 10. Commented-Out Code

Delete it. Version control remembers everything.

### 11. Nonlocal Information

Don't describe system-wide behavior in local comments.

### 12. Too Much Information

Skip historical context, RFCs, or algorithms - link to external docs instead.

## Guidelines

- Prefer extracting a well-named function over writing a comment
- If comment explains WHAT, refactor the code instead
- If comment explains WHY, it might be valuable
- Scan TODOs regularly and resolve them

## Exceptions

- **Third-party code**: Clarifying comments acceptable when you can't modify the source
- **Complex algorithms**: Mathematical or algorithmic intent may need explanation
- **Regulatory requirements**: Some domains require specific documentation

## Quick Reference

| Comment Type | Verdict |
|--------------|---------|
| Legal headers | OK |
| Intent explanation | OK |
| Warning of consequences | OK |
| TODO with context | OK (clean up regularly) |
| Amplification | OK |
| Public API docs | OK |
| Redundant/obvious | Bad |
| Commented-out code | Bad |
| Journal/changelog | Bad |
| Attributions | Bad |
| Closing brace markers | Bad |
