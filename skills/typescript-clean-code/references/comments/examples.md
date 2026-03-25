# Comments Examples

Code examples demonstrating comment principles.

## Bad Examples

### Explaining What Instead of Refactoring

```typescript
// Check to see if the employee is eligible for full benefits
if ((employee.flags & HOURLY_FLAG) && (employee.age > 65))
```

**Problems**:
- Comment explains WHAT, not WHY
- Code could be made self-explanatory

### Redundant Comment

```typescript
// Utility method that returns when this.closed is true.
// Throws an exception if the timeout is reached.
async waitForClose(timeoutMillis: number): Promise<void> {
  if (!this.closed) {
    await this.wait(timeoutMillis);
    if (!this.closed) {
      throw new Error('Could not be closed');
    }
  }
}
```

**Problems**:
- Comment is less precise than code
- Takes longer to read than the code itself
- Actually misleading (returns IF closed, not WHEN)

### Noise Comments

```typescript
/** The processor delay for this component */
protected backgroundProcessorDelay = -1;

/** The lifecycle event support for this component */
protected lifecycle = new LifecycleSupport(this);

/** The container event listeners for this Container */
protected listeners: EventListener[] = [];
```

**Problems**:
- Comments add zero information
- Variable names already describe purpose
- Creates clutter that obscures code

### Commented-Out Code

```typescript
const response = new InputStreamResponse();
response.setBody(formatter.getResultStream(), formatter.getByteCount());
// const resultsStream = formatter.getResultStream();
// const reader = new StreamReader(resultsStream);
// response.setContent(reader.read(formatter.getByteCount()));
```

**Problems**:
- No one will delete it (assumes it's important)
- Gathers like sediment over time
- Version control already preserves history

### Frustration Venting

```typescript
catch (e) {
  // Give me a break!
}
```

**Problems**:
- Provides no useful information
- Signals need for refactoring, not commenting

## Good Examples

### Explanation of Intent

```typescript
compareTo(other: unknown): number {
  if (other instanceof WikiPagePath) {
    const compressedName = this.names.join('');
    const compressedOther = other.names.join('');
    return compressedName.localeCompare(compressedOther);
  }
  return 1; // We are greater because we are the right type
}
```

**Why it works**:
- Explains the decision, not the mechanics
- Documents non-obvious business logic

### Warning of Consequences

```typescript
// Intl.DateTimeFormat caches results per locale, but options vary
// Create new instance to avoid stale format conflicts
function makeStandardHttpDateFormat(): Intl.DateTimeFormat {
  return new Intl.DateTimeFormat('en-US', {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZone: 'GMT'
  });
}
```

**Why it works**:
- Prevents future "optimization" that would break code
- Explains WHY, not WHAT

### Clarifying External Library Returns

```typescript
test('compareTo ordering', () => {
  const a = PathParser.parse('PageA');
  const b = PathParser.parse('PageB');
  const ab = PathParser.parse('PageA.PageB');

  expect(a.compareTo(a)).toBe(0);    // a == a
  expect(a.compareTo(b)).toBe(-1);   // a < b
  expect(b.compareTo(a)).toBe(1);    // b > a
  expect(ab.compareTo(a)).toBe(1);   // ab > a (child > parent)
});
```

**Why it works**:
- Makes test assertions readable
- Clarifies comparison semantics

## Refactoring Walkthrough

### Before

```typescript
// Check to see if the employee is eligible for full benefits
if ((employee.flags & HOURLY_FLAG) && (employee.age > 65))
```

### After

```typescript
if (employee.isEligibleForFullBenefits())
```

### Changes Made

1. Extracted condition to descriptive method name
2. Comment became unnecessary - code speaks for itself
3. Logic now reusable across codebase

---

### Before

```typescript
// does the module from the global list <mod> depend on the
// subsystem we are part of?
if (smodule.getDependSubsystems().contains(subSysMod.getSubSystem()))
```

### After

```typescript
const moduleDependencies = smodule.getDependSubsystems();
const ourSubSystem = subSysMod.getSubSystem();
if (moduleDependencies.includes(ourSubSystem))
```

### Changes Made

1. Extracted intermediate variables with descriptive names
2. Comment became redundant
3. Code now reads as English

---

### Before (with frustration comment)

```typescript
private startSending(): void {
  try {
    this.doSending();
  } catch (e) {
    if (e instanceof SocketException) {
      // normal. someone stopped the request.
    } else {
      try {
        this.response.add(ErrorResponder.makeExceptionString(e));
        this.response.closeAll();
      } catch {
        // Give me a break!
      }
    }
  }
}
```

### After

```typescript
private startSending(): void {
  try {
    this.doSending();
  } catch (e) {
    if (e instanceof SocketException) {
      // Expected when client disconnects mid-request
      return;
    }
    this.addExceptionAndCloseResponse(e);
  }
}

private addExceptionAndCloseResponse(e: Error): void {
  try {
    this.response.add(ErrorResponder.makeExceptionString(e));
    this.response.closeAll();
  } catch {
    // Response already closed, safe to ignore
  }
}
```

### Changes Made

1. Extracted nested try/catch to separate method
2. Replaced frustrated comment with meaningful one
3. Reduced complexity, improved readability
