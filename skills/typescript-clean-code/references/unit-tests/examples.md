# Unit Tests Examples

Code examples demonstrating clean test principles in TypeScript (Jest/Vitest).

## Bad Examples

### Noisy Test with Too Much Detail

```typescript
it('should get page hierarchy as XML', async () => {
  await crawler.addPage(root, PathParser.parse('PageOne'));
  await crawler.addPage(root, PathParser.parse('PageOne.ChildOne'));
  request.setResource('root');
  request.addInput('type', 'pages');
  const responder = new SerializedPageResponder();
  const response = await responder.makeResponse(new Context(root), request);
  expect(response.getContentType()).toBe('text/xml');
  expect(response.getContent()).toContain('<name>PageOne</name>');
});
```

**Problems**: PathParser calls are noise, setup details obscure intent, reader must understand implementation.

### Multiple Concepts in One Test

```typescript
it('should add months correctly', () => {
  const d1 = SerialDate.createInstance(31, 5, 2004);
  const d2 = SerialDate.addMonths(1, d1);
  expect(d2.getDayOfMonth()).toBe(30);  // Concept 1: cap to 30-day month
  const d3 = SerialDate.addMonths(2, d1);
  expect(d3.getDayOfMonth()).toBe(31);  // Concept 2: preserve in 31-day month
  const d4 = SerialDate.addMonths(1, SerialDate.addMonths(1, d1));
  expect(d4.getDayOfMonth()).toBe(30);  // Concept 3: chained addition
});
```

**Problems**: Three concepts in one test, hard to identify which failed, missing edge cases.

### Hard-to-Read State Assertions

```typescript
it('should turn on low temp alarm at threshold', () => {
  hw.setTemp(WAY_TOO_COLD);
  controller.tic();
  expect(hw.heaterState()).toBe(true);
  expect(hw.blowerState()).toBe(true);
  expect(hw.coolerState()).toBe(false);
  expect(hw.hiTempAlarm()).toBe(false);
  expect(hw.loTempAlarm()).toBe(true);
});
```

**Problems**: Eyes bounce between state names and values, tedious to read.

## Good Examples

### Clean Test with Domain Language

```typescript
it('should get page hierarchy as XML', async () => {
  await makePages('PageOne', 'PageOne.ChildOne', 'PageTwo');
  await submitRequest('root', 'type:pages');
  assertResponseIsXML();
  assertResponseContains('<name>PageOne</name>', '<name>PageTwo</name>');
});

it('should exclude symbolic links from hierarchy', async () => {
  const page = await makePage('PageOne');
  await addLinkTo(page, 'PageTwo', 'SymPage');
  await submitRequest('root', 'type:pages');
  assertResponseDoesNotContain('SymPage');
});
```

**Why it works**: Clear BUILD-OPERATE-CHECK structure, details hidden in helpers.

### Compact State Representation

```typescript
// State: Heater, Blower, Cooler, HiAlarm, LoAlarm (uppercase=ON)
it('should turn on cooler and blower if too hot', () => {
  tooHot();
  expect(hw.getState()).toBe('hBChl');
});

it('should turn on heater and blower if too cold', () => {
  tooCold();
  expect(hw.getState()).toBe('HBchl');
});

it('should turn on lo-temp alarm at threshold', () => {
  wayTooCold();
  expect(hw.getState()).toBe('HBchL');
});
```

**Why it works**: Compact, consistent format; once learned, easy to scan.

### Single Concept per Test

```typescript
describe('SerialDate.addMonths', () => {
  describe('given last day of 31-day month', () => {
    const may31 = SerialDate.createInstance(31, 5, 2004);

    it('should cap to 30th when target month has 30 days', () => {
      const june = SerialDate.addMonths(1, may31);
      expect(june.getDayOfMonth()).toBe(30);
    });

    it('should preserve 31st when target month has 31 days', () => {
      const july = SerialDate.addMonths(2, may31);
      expect(july.getDayOfMonth()).toBe(31);
    });
  });
});
```

**Why it works**: One concept per test, easy to add edge cases, clear what failed.

## Refactoring Walkthrough

### Before

```typescript
it('should handle page requests', async () => {
  await crawler.addPage(root, PathParser.parse('PageOne'));
  request.setResource('root');
  request.addInput('type', 'pages');
  const responder = new SerializedPageResponder();
  const response = await responder.makeResponse(ctx, request);
  expect(response.getContentType()).toBe('text/xml');
  expect(response.getContent()).toContain('<name>PageOne</name>');
});
```

### After

```typescript
it('should return XML containing requested pages', async () => {
  await makePages('PageOne');
  await submitRequest('root', 'type:pages');
  assertResponseIsXML();
  assertResponseContains('<name>PageOne</name>');
});
```

### Test Helpers (Domain-Specific Testing Language)

```typescript
async function makePages(...names: string[]): Promise<void> {
  for (const name of names) {
    await crawler.addPage(root, PathParser.parse(name));
  }
}

async function submitRequest(resource: string, params: string): Promise<void> {
  request.setResource(resource);
  const [key, value] = params.split(':');
  request.addInput(key, value);
  response = await new SerializedPageResponder().makeResponse(ctx, request);
}

function assertResponseIsXML(): void {
  expect(response.getContentType()).toBe('text/xml');
}

function assertResponseContains(...substrings: string[]): void {
  for (const s of substrings) {
    expect(response.getContent()).toContain(s);
  }
}
```

### Changes Made

1. **Extracted helpers** - Hide implementation noise (PathParser, request setup)
2. **Created assertion DSL** - Domain-specific language for responses
3. **Improved test name** - Describes behavior, not implementation
4. **Clear structure** - BUILD-OPERATE-CHECK pattern is visible
