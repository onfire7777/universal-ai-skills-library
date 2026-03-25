# Consistency Audit Framework

Purpose: Use this file when Zen must standardize one pattern across multiple files without changing behavior.

## Contents
- [Audit Categories](#audit-categories)
- [Cross-File Pattern Detection](#cross-file-pattern-detection)
- [Pattern Replacement Recipes](#pattern-replacement-recipes)
- [Consistency Report Template](#consistency-report-template)
- [Scope Tiers](#scope-tiers)
- [False Positive Filtering](#false-positive-filtering)

Cross-file pattern consistency analysis and unification guide.

---

## Overview

Consistency Audit detects divergent patterns across files and provides a structured approach to unify them. Unlike single-file refactoring, this focuses on **cross-file pattern alignment** — ensuring the same concept is expressed the same way throughout the codebase.

**Scope tiers apply:** See [Scope Tiers](#scope-tiers) for limits on single-pass changes.

---

## Audit Categories

### 6 Core Categories

| # | Category | What to Detect | Impact |
|---|----------|---------------|--------|
| 1 | **Error Handling** | Mixed try/catch styles, inconsistent error types, varied recovery patterns | HIGH |
| 2 | **API Call** | Different HTTP clients, inconsistent retry/timeout, mixed auth patterns | HIGH |
| 3 | **State Management** | Mixed state patterns (local vs global), inconsistent update flows | MEDIUM |
| 4 | **Logging** | Mixed logger usage, inconsistent log levels, varied structured fields | MEDIUM |
| 5 | **Naming** | Inconsistent casing, mixed terminology (e.g., user/account), prefix conflicts | LOW |
| 6 | **Import/Export** | Mixed default/named exports, inconsistent barrel files, path alias inconsistency | LOW |

### Severity Classification

| Severity | Criteria | Action |
|----------|----------|--------|
| **HIGH** | Causes runtime inconsistency, hides bugs, or complicates debugging | Fix in current cycle |
| **MEDIUM** | Reduces readability, increases cognitive load | Fix if touching related code |
| **LOW** | Cosmetic inconsistency, does not affect behavior | Document for future cleanup |

---

## Cross-File Pattern Detection

### 5-Step Process

```
1. SCAN     → Collect pattern instances across target files
2. CLASSIFY → Group by category (error handling, API call, etc.)
3. IDENTIFY → Determine dominant pattern (≥70% usage = canonical)
4. DEVIATE  → List files/locations that deviate from canonical
5. PLAN     → Create migration plan respecting scope tier
```

### Step 1: Scan

Collect all instances of a pattern category. Use AST-level tools when available:

| Language | Tool | Usage |
|----------|------|-------|
| TypeScript/JS | ast-grep | `ast-grep -p 'try { $$$ } catch($E) { $$$ }' src/` |
| TypeScript/JS | ESLint | Custom rule or `--rule` flag |
| Python | ruff | `ruff check --select E src/` |
| Go | staticcheck | `staticcheck ./...` |
| Multi-lang | grep (fallback) | Pattern-specific regex |

### Step 2: Classify

Group instances by variation:

```markdown
### Error Handling Variants Found

| Variant | Count | Files | Example |
|---------|-------|-------|---------|
| try/catch with custom Error class | 12 | auth/, payment/ | `throw new AppError(...)` |
| try/catch with string throw | 3 | legacy/ | `throw "failed"` |
| .catch() on Promise | 8 | api/ | `.catch(err => ...)` |
| Result type (no throw) | 2 | utils/ | `return { ok: false, error }` |
```

### Step 3: Identify Canonical Pattern

**Rule:** The dominant pattern (≥70% of occurrences) becomes the canonical pattern. If no pattern reaches 70%, flag for manual decision.

```
Dominant: try/catch with custom Error class (12/25 = 48%) ← Below 70%
→ Flag: "No clear dominant pattern. Manual decision required."
  Options:
    A) Adopt AppError class (most frequent, 48%)
    B) Adopt Result type (most modern, 8%)
    C) Mixed approach with migration plan
```

### Step 4: Identify Deviations

List every deviation from the canonical pattern:

```markdown
| # | File | Line | Current Pattern | Canonical Pattern | Migration Effort |
|---|------|------|----------------|-------------------|-----------------|
| 1 | `src/legacy/handler.ts` | 42 | string throw | AppError | Low |
| 2 | `src/api/client.ts` | 89 | .catch() chain | try/catch + AppError | Medium |
| 3 | `src/utils/parser.ts` | 15 | Result type | AppError | High (design choice) |
```

### Step 5: Create Migration Plan

See [Scope Tiers](#scope-tiers) for how to phase the migration.

---

## Pattern Replacement Recipes

### Recipe: Unify Error Handling

**Goal:** Standardize error handling to a single pattern across the codebase.

**Before (mixed):**
```typescript
// File A: Custom Error class
throw new AppError('NOT_FOUND', 'User not found');

// File B: String throw
throw 'User not found';

// File C: Inline Error
throw new Error('User not found');
```

**After (unified):**
```typescript
// All files: Custom Error class with code
throw new AppError('NOT_FOUND', 'User not found');
```

**Migration steps:**
1. Ensure `AppError` class exists and is importable
2. Replace string throws → `new AppError(code, message)`
3. Replace `new Error(msg)` → `new AppError(code, msg)`
4. Update catch blocks to handle `AppError` properties
5. Run tests after each file change

### Recipe: Unify API Client

**Goal:** Standardize HTTP client usage and configuration.

**Before (mixed):**
```typescript
// File A: Raw fetch
const res = await fetch(url, { headers: { Authorization: token } });

// File B: Axios instance
const res = await axios.get(url);

// File C: Custom wrapper
const res = await apiClient.get(url);
```

**After (unified):**
```typescript
// All files: Shared API client with built-in auth, retry, timeout
const res = await apiClient.get(url);
```

**Migration steps:**
1. Identify or create shared API client with standard config
2. Replace raw fetch calls → `apiClient` methods
3. Replace axios direct usage → `apiClient` methods
4. Consolidate auth header injection into client interceptor
5. Verify timeout and retry behavior is consistent

### Recipe: Unify State Management

**Goal:** Standardize state update patterns within a framework.

**React example — Before (mixed):**
```typescript
// File A: useState + useEffect
const [data, setData] = useState(null);
useEffect(() => { fetchData().then(setData); }, []);

// File B: useReducer
const [state, dispatch] = useReducer(reducer, initialState);

// File C: External store (Zustand)
const data = useStore(state => state.data);
```

**After (unified per layer):**
- Local component state → `useState`
- Complex component state → `useReducer`
- Shared/global state → Zustand store
- Server state → TanStack Query / SWR

### Recipe: Unify Import Style

**Goal:** Standardize import/export conventions.

**Before (mixed):**
```typescript
// File A: Default export
export default class UserService { ... }
import UserService from './UserService';

// File B: Named export
export class OrderService { ... }
import { OrderService } from './OrderService';

// File C: Barrel re-export
export * from './services';
```

**After (unified):**
```typescript
// All files: Named exports (more explicit, better tree-shaking)
export class UserService { ... }
import { UserService } from './UserService';

// Barrel files: Explicit named re-exports
export { UserService } from './UserService';
export { OrderService } from './OrderService';
```

---

## Consistency Report Template

```markdown
## Consistency Audit Report

### Audit Scope
- **Category**: [Error Handling / API Call / State Mgmt / Logging / Naming / Import-Export]
- **Files Scanned**: [count]
- **Scope Tier**: [Focused / Module / Project-wide]

### Canonical Pattern
| Aspect | Detail |
|--------|--------|
| Pattern | [Description] |
| Usage | [X/Y files, Z%] |
| Justification | [Why this is canonical] |

### Deviation Summary
| Severity | Count | Action |
|----------|-------|--------|
| HIGH | X | Fix now |
| MEDIUM | Y | Fix if touching |
| LOW | Z | Document |

### Deviation Details
| # | File:Line | Current | Canonical | Severity | Effort |
|---|-----------|---------|-----------|----------|--------|
| 1 | `path:42` | [variant] | [canonical] | HIGH | Low |

### Migration Plan
| Phase | Files | Pattern Change | Estimated Lines |
|-------|-------|---------------|----------------|
| 1 | [1-3 files] | [change] | [~N lines] |
| 2 | [next files] | [change] | [~N lines] |

### Risks
- [False positive risk: framework-required pattern differences]
- [Migration risk: behavioral side effects]
```

---

## Scope Tiers

| Tier | Files | Max Lines Changed | Zen Approach |
|------|-------|-------------------|-------------|
| **Focused** | 1-3 | ≤50 | Direct refactoring (standard Zen scope) |
| **Module** | 4-10 | ≤100 | Mechanical pattern replacement only (same pattern, no logic changes) |
| **Project-wide** | 10+ | Plan only | Generate migration plan, execute in multiple PRs |

### Tier Rules

**Focused (1-3 files):**
- Standard Zen 50-line limit applies
- Full Before/After report
- Tests BEFORE and AFTER

**Module (4-10 files):**
- Extended to 100 lines maximum
- **Restricted to mechanical pattern replacement** (same transform applied identically)
- No logic changes, no behavioral modifications
- Must be the **same pattern category** across all files
- Tests BEFORE and AFTER each batch

**Project-wide (10+ files):**
- Zen produces a **migration plan only** (no code changes)
- Plan includes: file list, priority order, estimated effort per file
- Execution split into multiple Focused/Module-tier PRs
- Each PR independently testable

---

## Tool Reference

### AST-level Pattern Search

```bash
# ast-grep: Find all try-catch blocks
ast-grep -p 'try { $$$ } catch($ERR) { $$$ }' src/

# ast-grep: Find all fetch calls
ast-grep -p 'fetch($URL, $$$)' src/

# ast-grep: Find all console.log
ast-grep -p 'console.log($$$)' src/

# ESLint: Find inconsistent error handling
npx eslint --rule 'no-throw-literal: error' src/

# ruff: Python pattern checks
ruff check --select E,W src/

# staticcheck: Go consistency
staticcheck -checks all ./...
```

### Pattern Counting

```bash
# Count error handling variants (TypeScript)
# try-catch blocks
ast-grep -p 'try { $$$ } catch($E) { $$$ }' src/ --json | jq length

# throw new Error
ast-grep -p 'throw new Error($$$)' src/ --json | jq length

# throw string literal
ast-grep -p 'throw "$$$"' src/ --json | jq length
```

---

## False Positive Filtering

Not all pattern variations are inconsistencies. Filter these:

| Exception | Reason | Action |
|-----------|--------|--------|
| Framework-required patterns | React hooks must follow rules of hooks | Skip |
| Test-specific patterns | Test files may use different assertion styles | Skip |
| Generated code | Auto-generated files have their own patterns | Exclude from scan |
| Legacy migration in progress | Files marked with `// TODO: migrate` | Document, don't flag |
| Third-party adapter patterns | Wrappers around external APIs may need different patterns | Skip |
| Performance-critical paths | Optimized code may justify pattern deviation | Flag as LOW |

### Dominant Pattern Threshold

- **≥70%**: Clear canonical pattern → proceed with migration
- **50-69%**: Weak consensus → flag for team decision, recommend strongest candidate
- **<50%**: No consensus → escalate to Atlas for architectural guidance

---

## Zen vs Other Agents

| Task | Agent | Reason |
|------|-------|--------|
| Detect pattern inconsistencies | **Judge** (consistency detection) | Pattern violation as review finding |
| Unify patterns (refactor) | **Zen** (consistency audit) | Mechanical replacement, no behavior change |
| Architectural pattern decisions | **Atlas** | When no dominant pattern exists |
| Test pattern consistency | **Radar** | Test-specific patterns (fixtures, mocks) |
| Document canonical patterns | **Quill** | ADR or coding standards documentation |
