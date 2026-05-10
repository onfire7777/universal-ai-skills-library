# Formatting Examples

Code examples demonstrating formatting principles in TypeScript.

## Bad Examples

### Missing Vertical Openness

```typescript
import { regex } from './utils';
import { ParentWidget } from './widgets';
export class BoldWidget extends ParentWidget {
  private static readonly REGEXP = /'''.+?'''/;
  private pattern = /'''(.+?)'''/gs;
  constructor(parent: ParentWidget, text: string) {
    super(parent);
    const match = this.pattern.exec(text);
    if (match) this.addChildWidgets(match[1]);
  }
  render(): string {
    return `<b>${this.childHtml()}</b>`;
  }
}
```

**Problems**:
- No blank lines between imports and class
- No separation between static fields and constructor
- No separation between constructor and methods
- Code appears as an unstructured block

### Excessive Comments Breaking Density

```typescript
class ReporterConfig {
  /**
   * The class name of the reporter listener
   */
  private className: string;

  /**
   * The properties of the reporter listener
   */
  private properties: Property[] = [];

  public addProperty(property: Property): void {
    this.properties.push(property);
  }
}
```

**Problems**:
- Trivial comments separate related declarations
- Forces reader to scan more vertical space
- Comments add no information beyond the names

### Instance Variables Hidden in Class

```typescript
class TestSuite implements Test {
  static createTest(testClass: TestClass, name: string): Test {
    // ... implementation
  }

  static getTestConstructor(testClass: TestClass): Constructor {
    // ... implementation
  }

  // Instance variables buried after static methods!
  private name: string;
  private tests: Test[] = [];

  constructor() {}
}
```

**Problems**:
- Instance variables declared in middle of class
- Reader must hunt to find class state
- Violates convention of declarations at top

### No Indentation

```typescript
class FitNesseServer implements SocketServer { private context: FitNesseContext; constructor(context: FitNesseContext) { this.context = context; } serve(socket: Socket): void { this.serveWithTimeout(socket, 10000); } serveWithTimeout(socket: Socket, timeout: number): void { try { const sender = new FitNesseExpediter(socket, this.context); sender.setRequestParsingTimeLimit(timeout); sender.start(); } catch (e) { console.error(e); } } }
```

**Problems**:
- Impossible to discern structure
- Cannot identify method boundaries
- Cannot see nesting levels

### Collapsed One-Liners

```typescript
class CommentWidget extends TextWidget {
  static readonly REGEXP = /^#[^\r\n]*(?:(?:\r\n)|\n|\r)?/;

  constructor(parent: ParentWidget, text: string) { super(parent, text); }
  render(): string { return ''; }
}
```

**Problems**:
- Methods are hard to scan
- Inconsistent visual rhythm
- Easy to miss method boundaries

## Good Examples

### Proper Vertical Openness

```typescript
import { regex } from './utils';

import { ParentWidget } from './widgets';

export class BoldWidget extends ParentWidget {
  private static readonly REGEXP = /'''.+?'''/;
  private pattern = /'''(.+?)'''/gs;

  constructor(parent: ParentWidget, text: string) {
    super(parent);
    const match = this.pattern.exec(text);
    if (match) {
      this.addChildWidgets(match[1]);
    }
  }

  render(): string {
    return `<b>${this.childHtml()}</b>`;
  }
}
```

**Why it works**:
- Blank lines separate logical groups
- Easy to scan and identify structure
- Each concept is visually distinct

### Related Code Kept Dense

```typescript
class ReporterConfig {
  private className: string;
  private properties: Property[] = [];

  addProperty(property: Property): void {
    this.properties.push(property);
  }
}
```

**Why it works**:
- Related properties grouped together
- No unnecessary vertical separation
- Fits in one "eye-full"

### Proper Function Ordering (Caller Above Callee)

```typescript
class WikiPageResponder implements SecureResponder {
  private page: WikiPage | null = null;
  private pageData: PageData | null = null;
  private crawler: PageCrawler;

  makeResponse(context: FitNesseContext, request: Request): Response {
    const pageName = this.getPageNameOrDefault(request, 'FrontPage');
    this.loadPage(pageName, context);

    if (!this.page) {
      return this.notFoundResponse(context, request);
    }
    return this.makePageResponse(context);
  }

  private getPageNameOrDefault(request: Request, defaultName: string): string {
    const pageName = request.getResource();
    return pageName || defaultName;
  }

  private loadPage(resource: string, context: FitNesseContext): void {
    const path = PathParser.parse(resource);
    this.crawler = context.root.getPageCrawler();
    this.page = this.crawler.getPage(context.root, path);

    if (this.page) {
      this.pageData = this.page.getData();
    }
  }

  private notFoundResponse(context: FitNesseContext, request: Request): Response {
    return new NotFoundResponder().makeResponse(context, request);
  }

  private makePageResponse(context: FitNesseContext): SimpleResponse {
    const html = this.makeHtml(context);
    const response = new SimpleResponse();
    response.setMaxAge(0);
    response.setContent(html);
    return response;
  }
}
```

**Why it works**:
- Public entry point at top
- Called functions follow their callers
- Creates natural top-down reading flow
- Easy to trace execution path

### Proper Indentation

```typescript
class FitNesseServer implements SocketServer {
  private context: FitNesseContext;

  constructor(context: FitNesseContext) {
    this.context = context;
  }

  serve(socket: Socket): void {
    this.serveWithTimeout(socket, 10000);
  }

  serveWithTimeout(socket: Socket, timeout: number): void {
    try {
      const sender = new FitNesseExpediter(socket, this.context);
      sender.setRequestParsingTimeLimit(timeout);
      sender.start();
    } catch (e) {
      console.error(e);
    }
  }
}
```

**Why it works**:
- Structure immediately visible
- Scopes clearly delineated
- Easy to navigate mentally

## Refactoring Walkthrough

### Before

```typescript
class Assert {
static assertTrue(message: string, condition: boolean): void { if (!condition) fail(message); }
static assertTrue(condition: boolean): void { assertTrue(null, condition); }
static assertFalse(message: string, condition: boolean): void { assertTrue(message, !condition); }
static assertFalse(condition: boolean): void { assertFalse(null, condition); }
}
```

### After

```typescript
class Assert {
  static assertTrue(condition: boolean, message?: string): void {
    if (!condition) {
      fail(message);
    }
  }

  static assertFalse(condition: boolean, message?: string): void {
    Assert.assertTrue(!condition, message);
  }
}
```

### Changes Made

1. Added proper indentation to reveal class structure
2. Expanded method bodies with braces
3. Combined overloads using optional parameters (TypeScript idiom)
4. Added blank line between methods for visual separation
5. Used proper braces for if statement
