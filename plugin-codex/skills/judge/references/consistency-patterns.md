# Consistency Pattern Detection

Detection heuristics for cross-file pattern inconsistencies. Judge identifies violations; Zen performs the actual unification.

---

## Overview

Consistency detection identifies places where the same concept is expressed differently across files. These are not bugs — they're maintenance liabilities that increase cognitive load and error risk.

**Judge's role:** Detect and report consistency violations with severity.
**Zen's role:** Perform the actual pattern unification (via consistency audit).

---

## 6 Consistency Categories

### 1. Error Handling Style

**What to detect:** Mixed error handling patterns within the same codebase layer.

**Detection heuristics:**
- Mixed `try/catch` + `.catch()` in the same service layer
- Inconsistent error class usage (`new Error` vs custom error classes)
- Some functions throw, others return error codes/null
- Inconsistent error message format (template literals vs concatenation)

**BAD (mixed):**
```typescript
// Service A: throws custom error
throw new AppError('NOT_FOUND', `User ${id} not found`);

// Service B: throws generic Error
throw new Error('User not found');

// Service C: returns null on error
return null;
```

**GOOD (consistent):**
```typescript
// All services: custom error with code
throw new AppError('NOT_FOUND', `User ${id} not found`);
```

**Severity:** HIGH — Mixed error handling makes catch blocks unreliable.

### 2. Null Safety

**What to detect:** Inconsistent null/undefined handling patterns.

**Detection heuristics:**
- Mix of optional chaining (`?.`) and manual null checks
- Some functions check for null, similar functions don't
- Inconsistent use of default values (`??` vs `||` vs ternary)

**BAD (mixed):**
```typescript
// File A: optional chaining
const name = user?.profile?.name ?? 'Unknown';

// File B: manual null check
const name = user && user.profile ? user.profile.name : 'Unknown';

// File C: no null check
const name = user.profile.name; // potential crash
```

**GOOD (consistent):**
```typescript
// All files: optional chaining with nullish coalescing
const name = user?.profile?.name ?? 'Unknown';
```

**Severity:** HIGH — Inconsistent null handling leads to runtime crashes.

### 3. Async Pattern

**What to detect:** Mixed asynchronous code styles.

**Detection heuristics:**
- Mix of `async/await` and `.then()/.catch()` chains
- Inconsistent error handling in async code (try/catch vs .catch)
- Mixed `Promise.all` vs sequential awaits for independent operations
- Callback-style mixed with Promise-style

**BAD (mixed):**
```typescript
// File A: async/await
async function fetchUser(id) {
  try {
    const user = await db.findUser(id);
    return user;
  } catch (err) {
    throw new AppError('FETCH_FAILED', err.message);
  }
}

// File B: .then() chain
function fetchOrder(id) {
  return db.findOrder(id)
    .then(order => order)
    .catch(err => { throw new AppError('FETCH_FAILED', err.message); });
}
```

**GOOD (consistent):**
```typescript
// All files: async/await with try/catch
async function fetchUser(id) {
  try {
    return await db.findUser(id);
  } catch (err) {
    throw new AppError('FETCH_FAILED', err.message);
  }
}
```

**Severity:** MEDIUM — Mixed async patterns increase cognitive load and error-prone patterns.

### 4. Naming Convention

**What to detect:** Inconsistent naming across same-type identifiers.

**Detection heuristics:**
- Mixed casing in same context (camelCase vs snake_case for same type)
- Inconsistent prefixes/suffixes (`IUser` vs `UserInterface` vs `User`)
- Mixed terminology for same concept (`user`/`account`/`member`)
- Boolean naming inconsistency (`isActive` vs `active` vs `enabled`)

**BAD (mixed):**
```typescript
// File A
interface IUserService { ... }
const is_active = user.isActive;

// File B
interface OrderServiceInterface { ... }
const isActive = order.active;

// File C
type AccountData = { ... }; // Same concept as "User"
```

**GOOD (consistent):**
```typescript
// All files: no prefix, consistent casing, same terminology
interface UserService { ... }
interface OrderService { ... }
const isActive = user.isActive;
```

**Severity:** LOW — Naming inconsistency reduces readability but doesn't cause bugs.

### 5. Import/Export Style

**What to detect:** Mixed module import/export conventions.

**Detection heuristics:**
- Mix of default and named exports for similar constructs
- Inconsistent import path style (relative vs alias)
- Some files use barrel re-exports, others don't
- Mixed CommonJS (`require`) and ES modules (`import`)

**BAD (mixed):**
```typescript
// File A: default export
export default class UserService { ... }

// File B: named export
export class OrderService { ... }

// File C: CommonJS
module.exports = { PaymentService };
```

**GOOD (consistent):**
```typescript
// All files: named exports
export class UserService { ... }
export class OrderService { ... }
export class PaymentService { ... }
```

**Severity:** LOW — Import inconsistency mainly affects maintainability and tree-shaking.

### 6. Error Type

**What to detect:** Mixed error type definitions and hierarchies.

**Detection heuristics:**
- Multiple custom error classes with no common base
- Some modules use error codes, others use error classes
- Inconsistent error properties (`code`, `statusCode`, `status`, `errorCode`)
- Mixed HTTP status code usage for same error type

**BAD (mixed):**
```typescript
// Module A
class NotFoundError extends Error { statusCode = 404; }

// Module B
class ResourceNotFound extends AppError { code = 'NOT_FOUND'; status = 404; }

// Module C
throw { error: 'not_found', httpStatus: 404 };
```

**GOOD (consistent):**
```typescript
// All modules: single error hierarchy
class AppError extends Error {
  constructor(public code: string, message: string, public statusCode: number) {
    super(message);
  }
}
class NotFoundError extends AppError {
  constructor(resource: string, id: string) {
    super('NOT_FOUND', `${resource} ${id} not found`, 404);
  }
}
```

**Severity:** HIGH — Mixed error types break centralized error handling middleware.

---

## False Positive Filtering

### When NOT to Flag

| Situation | Reason | Example |
|-----------|--------|---------|
| **Framework-required pattern** | Framework mandates specific style | React hooks must use specific patterns |
| **Legacy migration in progress** | File marked with migration TODO | `// TODO(migrate): convert to async/await` |
| **Adapter/wrapper layers** | Boundary code adapting external APIs | GraphQL resolver wrapping REST |
| **Test-specific patterns** | Tests may intentionally use different patterns | Mock factories vs production factories |
| **Generated code** | Auto-generated files have their own style | Prisma client, protobuf |
| **Performance-critical path** | Justified deviation for performance | Using callbacks for hot path |

### Confidence Threshold

Only report as consistency violation when:
- Pattern deviation appears in ≥3 files (not a one-off)
- Dominant pattern has ≥70% usage
- The deviation is not in the exclusion list above

---

## Reporting Format

### Finding Format

```markdown
#### [CONSISTENCY-NNN] [Category]: [Title]

- **Category**: [Error Handling / Null Safety / Async Pattern / Naming / Import-Export / Error Type]
- **Severity**: [HIGH / MEDIUM / LOW]
- **Dominant Pattern**: [Description of canonical pattern]
- **Deviation**: [Description of non-canonical pattern]
- **Files Affected**:
  - `path/to/file1.ts:42` — [specific deviation]
  - `path/to/file2.ts:18` — [specific deviation]
- **Canonical Example**:
  ```typescript
  // Expected pattern
  ```
- **Remediation Agent**: Zen (consistency audit)
```

### Summary in Review Report

Add to the standard Judge Review Report under a new section:

```markdown
### Consistency Findings

| ID | Category | Severity | Files | Dominant Pattern |
|----|----------|----------|-------|-----------------|
| CONSISTENCY-001 | Error Handling | HIGH | 3/12 deviate | AppError class |
| CONSISTENCY-002 | Async Pattern | MEDIUM | 5/20 deviate | async/await |

**Recommendation**: Route to Zen for consistency audit on Error Handling (HIGH priority).
```

---

## Handoff to Zen

When consistency violations are detected, include in `JUDGE_TO_ZEN_HANDOFF`:

```markdown
## JUDGE_TO_ZEN_HANDOFF (Consistency)

**Review ID**: [PR# or audit scope]
**Type**: Consistency Violation

**Consistency Findings**:

### [CONSISTENCY-001] Error Handling: Mixed error types
| Aspect | Detail |
|--------|--------|
| Category | Error Handling |
| Severity | HIGH |
| Dominant Pattern | `throw new AppError(code, message)` (9/12 files) |
| Deviations | 3 files using `throw new Error(message)` |
| Files | `src/legacy/handler.ts:42`, `src/api/v1/client.ts:89`, `src/utils/parser.ts:15` |

**Request**: Apply consistency audit — unify error handling to AppError pattern.
**Scope Tier**: Module (3 files, mechanical replacement)
```
