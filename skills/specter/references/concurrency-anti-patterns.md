# Concurrency & Async Anti-Patterns

Purpose: async/concurrency guardrails, promise anti-pattern IDs, race-prevention strategy matrix, and deadlock rules.

## Contents

1. Promise anti-patterns
2. Race-prevention strategies
3. Deadlock prevention
4. Concurrency-safe patterns

## Promise Anti-Patterns

| ID | Anti-Pattern | Risk | Preferred Response |
|----|--------------|------|--------------------|
| `AP-001` | Explicit Promise construction around async work | swallowed errors and confusion | return the original promise or use `async` function directly |
| `AP-002` | `forEach` with async | loop does not await completion | use `for...of` for sequential or `Promise.all` for parallel |
| `AP-003` | fragile `Promise.all()` error handling | one reject hides partial success | use `Promise.allSettled()` when partial results matter |
| `AP-004` | `return` without `await` inside `try-catch` | rejection bypasses `catch` | `return await ...` when local catch is required |
| `AP-005` | `await` in loop without intent | accidental serialization | parallelize unless ordering, rate limit, or dependency requires serial execution |

## Race-Prevention Strategy Matrix

| Strategy | Use when | Implementation hint |
|----------|----------|---------------------|
| AbortController | cancellable fetch or React effect | abort stale requests |
| Request ID tracking | latest result should win | ignore stale responses |
| Debouncing | user-input-triggered async work | delay until input stabilizes |
| Atomic update | DB or state mutation | `$inc`, conditional update, compare-and-swap |
| Optimistic locking | conflict-prone writes | version field |
| Distributed lock | cross-process coordination | Redis or mutex service |
| Cancellation token | long-running background work | cooperative cancellation |

## Deadlock Prevention

### Coffman Conditions

Deadlock requires all four:
1. mutual exclusion
2. hold and wait
3. no preemption
4. circular wait

Break any one of them to prevent deadlock.

### JavaScript / Node.js Deadlock Traps

| Trap | Guardrail |
|------|-----------|
| nested async locks | enforce a single lock acquisition order |
| promise dependency cycles | inspect dependency graph and break cycles |
| connection-pool exhaustion | set pool size and acquire timeout |
| worker-thread synchronous waiting | keep communication asynchronous |

### Deadlock Checklist

- lock acquisition order is globally consistent
- locks have timeouts
- mutex/semaphore libraries are used safely
- connection pools have `maxConnections` and `acquireTimeout`
- workers and main thread do not synchronously wait on each other

## Concurrency-Safe Patterns

- immutability for shared state
- actor-style message passing
- semaphore or mutex when bounded concurrency is required
