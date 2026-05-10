# Formatting Rules

Specific guidelines for vertical and horizontal code formatting.

## Vertical Formatting Rules

### 1. Keep Files Small

Target 200 lines, with an upper limit of 500 lines. Significant systems can be built with small files.

- Small files are easier to understand
- If a file grows too large, split it

### 2. Follow the Newspaper Metaphor

Structure files with high-level concepts first, details last.

- File name should tell you if you're in the right module
- Top of file: high-level concepts and algorithms
- Bottom of file: lowest-level functions and details

### 3. Separate Concepts with Blank Lines

Use blank lines between package/import statements, classes, and functions.

```typescript
// Good - concepts separated
import { Logger } from './logger';

import { Config } from './config';

export class BoldWidget extends ParentWidget {
  private static readonly REGEXP = /'''.+?'''/;

  constructor(parent: ParentWidget, text: string) {
    super(parent);
  }

  render(): string {
    return `<b>${this.childHtml()}</b>`;
  }
}
```

### 4. Keep Related Code Dense

Lines of tightly related code should appear together without blank lines.

```typescript
// Bad - unnecessary separation
class ReporterConfig {
  /**
   * The class name of the reporter listener
   */
  private className: string;

  /**
   * The properties of the reporter listener
   */
  private properties: Property[] = [];
}

// Good - related items together
class ReporterConfig {
  private className: string;
  private properties: Property[] = [];

  addProperty(property: Property): void {
    this.properties.push(property);
  }
}
```

### 5. Minimize Vertical Distance

Related concepts should be vertically close to each other.

- **Variables**: Declare close to usage, typically at function top
- **Instance variables**: Declare at top of class
- **Dependent functions**: Caller above callee
- **Conceptually related functions**: Keep near each other

### 6. Order Functions Top-Down

Functions that are called should appear below their callers.

```typescript
// Good - caller above callee
class WikiPageResponder {
  makeResponse(context: Context, request: Request): Response {
    const pageName = this.getPageNameOrDefault(request, 'FrontPage');
    this.loadPage(pageName, context);
    return this.page ? this.makePageResponse(context) : this.notFoundResponse();
  }

  private getPageNameOrDefault(request: Request, defaultName: string): string {
    return request.getResource() || defaultName;
  }

  private loadPage(resource: string, context: Context): void {
    // implementation
  }
}
```

## Horizontal Formatting Rules

### 7. Keep Lines Short

Limit lines to 100-120 characters maximum.

- Programmers prefer short lines
- Don't rely on horizontal scrolling

### 8. Use Whitespace to Show Relationships

Surround operators with spaces; keep function names attached to parentheses.

```typescript
// Good - whitespace shows relationships
function measureLine(line: string): void {
  lineCount++;
  const lineSize = line.length;
  totalChars += lineSize;
  lineWidthHistogram.addLine(lineSize, lineCount);
  recordWidestLine(lineSize);
}
```

### 9. Avoid Column Alignment

Don't align variable names or values in columns - it emphasizes wrong things.

```typescript
// Bad - aligned columns
private socket:                   Socket;
private input:                    InputStream;
private requestParsingTimeLimit:  number;

// Good - natural formatting
private socket: Socket;
private input: InputStream;
private requestParsingTimeLimit: number;
```

### 10. Use Proper Indentation

Indent code proportionally to its hierarchy level.

- Class-level declarations: no indent
- Methods: one level in from class
- Method bodies: one level in from method
- Nested blocks: one level per nesting

### 11. Never Collapse Short Statements

Always expand and indent scopes, even for one-liners.

```typescript
// Bad - collapsed
constructor(parent: ParentWidget, text: string) { super(parent, text); }
render(): string { return ''; }

// Good - expanded
constructor(parent: ParentWidget, text: string) {
  super(parent, text);
}

render(): string {
  return '';
}
```

## Team Rules

### 12. Agree on Team Standards

Every team should agree on a single formatting style.

- Decide on braces, indentation, naming conventions
- Encode rules in automated formatters (Prettier, ESLint)
- All members follow team rules, not personal preferences
- Consistency across the codebase is the goal

## Guidelines

- Use automated formatters (Prettier) to enforce consistency
- Configure ESLint for style enforcement
- Keep format configuration in version control
- Run formatters on save or pre-commit

## Exceptions

- **Legacy code**: May need gradual reformatting to avoid massive diffs
- **Generated code**: May have different formatting requirements
- **External dependencies**: Don't reformat third-party code

## Quick Reference

| Rule | Summary |
|------|---------|
| File size | 200-500 lines max |
| Line width | 100-120 characters max |
| Blank lines | Separate concepts, not related items |
| Variable declarations | Close to usage |
| Instance properties | Top of class |
| Function order | Caller above callee |
| Indentation | Always use, never collapse |
| Team rules | Follow team standards, use formatters |
