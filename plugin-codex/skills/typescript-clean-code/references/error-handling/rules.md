# Error Handling Rules

Rules for writing clean, robust error handling code.

## Core Rules

### 1. Use Exceptions, Not Return Codes

Throw exceptions when you encounter errors instead of returning error codes or flags.

- Return codes clutter caller with immediate checks
- Easy to forget to check return codes
- Exceptions separate error handling from main logic

**Example**:
```typescript
// Bad - return codes
function getUser(id: string): User | null {
  const user = db.find(id);
  if (!user) return null;  // Caller must check
  return user;
}

// Good - exceptions
function getUser(id: string): User {
  const user = db.find(id);
  if (!user) throw new UserNotFoundError(`User ${id} not found`);
  return user;
}
```

### 2. Write Try-Catch-Finally First

Start with the error handling structure before writing the logic.

- Try blocks define a transaction scope
- Catch must leave program in consistent state
- Establishes expectations for callers upfront

**Example**:
```typescript
// Start with the structure
async function readConfig(path: string): Promise<Config> {
  try {
    // Add logic here after establishing scope
    const content = await fs.readFile(path, 'utf-8');
    return JSON.parse(content);
  } catch (error) {
    throw new ConfigError(`Failed to read config: ${path}`, { cause: error });
  }
}
```

### 3. Provide Context with Exceptions

Include enough information to determine source and location of errors.

- Mention the operation that failed
- Include the type of failure
- Provide enough info for logging

**Example**:
```typescript
// Bad - no context
throw new Error('Failed');

// Good - rich context
throw new StorageError(
  `Failed to save user ${userId} to database: connection timeout after ${timeout}ms`
);
```

### 4. Define Exceptions by Caller's Needs

Classify exceptions based on how they will be caught, not their source.

- Wrap third-party exceptions into your own types
- One exception class per area is often sufficient
- Use different classes only when callers handle them differently

**Example**:
```typescript
// Bad - exposing third-party exceptions
try {
  await externalApi.call();
} catch (e) {
  if (e instanceof AxiosError) { /* ... */ }
  if (e instanceof TimeoutError) { /* ... */ }
  if (e instanceof NetworkError) { /* ... */ }
}

// Good - wrapped with common type
try {
  await apiClient.call();
} catch (e) {
  if (e instanceof ApiError) {
    logger.error(e.message);
  }
}
```

### 5. Define the Normal Flow

Use Special Case pattern to avoid exceptions in business logic.

- Return default objects instead of throwing for expected cases
- Encapsulate special behavior in the object itself
- Keep business logic clean and linear

**Example**:
```typescript
// Bad - exception for normal case
try {
  const expenses = await getMealExpenses(employeeId);
  total += expenses.getTotal();
} catch (e) {
  if (e instanceof NoExpensesError) {
    total += getPerDiem();
  }
}

// Good - Special Case pattern
const expenses = await getMealExpenses(employeeId); // Returns PerDiemExpenses if none
total += expenses.getTotal();
```

### 6. Don't Return Null

Return empty collections, throw exceptions, or use Special Case objects.

- Null forces callers to add null checks everywhere
- One missing check causes runtime errors
- Empty collections are safe to iterate

**Example**:
```typescript
// Bad - returning null
function getEmployees(): Employee[] | null {
  if (noEmployees) return null;
  return employees;
}

// Good - return empty array
function getEmployees(): Employee[] {
  if (noEmployees) return [];
  return employees;
}
```

### 7. Don't Pass Null

Avoid passing null as function arguments.

- Forces defensive null checks in every function
- No good way to handle accidentally passed null
- Forbid null by default in your codebase

**Example**:
```typescript
// Bad - allows null
function calculateDistance(p1: Point | null, p2: Point | null): number {
  if (!p1 || !p2) throw new Error('Invalid arguments');
  return Math.sqrt((p2.x - p1.x) ** 2 + (p2.y - p1.y) ** 2);
}

// Good - require valid arguments, use TypeScript strict mode
function calculateDistance(p1: Point, p2: Point): number {
  return Math.sqrt((p2.x - p1.x) ** 2 + (p2.y - p1.y) ** 2);
}
```

## Guidelines

- Wrap third-party APIs to minimize dependencies and enable mocking
- Use a single exception class per module/area when possible
- Push error detection to the edges of your program
- Use TypeScript strict null checks to catch null issues at compile time
- Use optional chaining (`?.`) and nullish coalescing (`??`) for safe property access

## Exceptions

- **External APIs expecting null**: When calling APIs that require null, pass it
- **Performance-critical paths**: Return codes may be acceptable in hot loops
- **Optional parameters**: Use TypeScript optional parameters (`param?: Type`) instead of null

## Quick Reference

| Rule | Summary |
|------|---------|
| Use Exceptions | Throw instead of return error codes |
| Try-First | Start with try-catch-finally structure |
| Context | Include operation and failure type in errors |
| Caller's Needs | Define exceptions by how they're caught |
| Normal Flow | Use Special Case pattern for expected variations |
| No Return Null | Return empty collections or throw |
| No Pass Null | Forbid null arguments by default |
