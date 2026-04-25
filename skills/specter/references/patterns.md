# Detection Patterns

Purpose: canonical detection catalog for Specter. Use these IDs, regex cues, context checks, and confidence rules during `SCAN` and `ANALYZE`.

## Contents

1. Memory leak patterns
2. Race condition patterns
3. Resource leak patterns
4. Async anti-patterns
5. Deadlock patterns
6. Scan priority and confidence

## Memory Leak Patterns

| ID | Signal | Regex Cue | Verify |
|----|--------|-----------|--------|
| `ML-001` | `addEventListener` without cleanup | `addEventListener...` without nearby `removeEventListener` | confirm `useEffect` cleanup or equivalent |
| `ML-002` | anonymous event handler | `addEventListener(..., () => ...)` | confirm handler cannot be removed later |
| `ML-003` | `setInterval` without `clearInterval` | `setInterval...` without nearby `clearInterval` | confirm lifecycle cleanup |
| `ML-004` | `setTimeout` without `clearTimeout` | `setTimeout...` without nearby `clearTimeout` | lower risk than intervals; check unmount behavior |
| `ML-005` | subscription without unsubscribe | `.subscribe(...)` without nearby `.unsubscribe` | confirm lifecycle-managed cleanup |
| `ML-006` | large object retained in closure | `useEffect` or callback captures `largeData`-like objects | manual check of retained references |

## Race Condition Patterns

| ID | Signal | Regex Cue | Verify |
|----|--------|-----------|--------|
| `RC-001` | state mutation after async boundary | `await` or `.then()` near `setState` / `setX` | confirm cancellation or staleness guard |
| `RC-002` | non-atomic read-modify-write | `x = x +/- ...` or separate read then write | confirm atomic update or compare-and-swap |
| `RC-003` | optional chaining after `await` | `await ...; ...?.` | may indicate timing uncertainty |
| `RC-004` | multiple awaits writing same state | multiple `await` then same mutator | confirm ordering and stale-response handling |
| `RC-005` | `useEffect` async work with no cleanup | async fetch/effect without return cleanup | confirm abort/cancel guard |

## Resource Leak Patterns

| ID | Signal | Regex Cue | Verify |
|----|--------|-----------|--------|
| `RL-001` | connection acquired without release | `getConnection` / `connect` without nearby `release` / `close` | confirm `try-finally`, `with`, or `defer` |
| `RL-002` | file or stream opened without close | `openSync` / `createReadStream` without close/destroy | confirm end/error cleanup |
| `RL-003` | WebSocket without close handler | `new WebSocket(...)` without nearby `.close()` | confirm component or service cleanup |

## Async Anti-Patterns

| ID | Signal | Regex Cue | Verify |
|----|--------|-----------|--------|
| `AA-001` | async call without `await` or return | `fetch|axios|api.` call without `await` | allow intentional fire-and-forget only when documented |
| `AA-002` | promise chain without `.catch()` | `.then(...)` without `.catch()` or `.finally()` | low false-positive risk |
| `AA-003` | async function without `try-catch` | async function with no nearby `try` | context-dependent; may be handled upstream |
| `AA-004` | async work inside `useEffect` with no cleanup | effect contains async/fetch/timer but no return | confirm lifecycle or framework-managed cleanup |

## Deadlock Patterns

| ID | Signal | Regex Cue | Verify |
|----|--------|-----------|--------|
| `DL-001` | `await` inside promise executor | `new Promise(... await ...)` | often signals confusion and potential lock or error-flow issues |
| `DL-002` | nested lock acquisition | `acquire|lock` followed by another `acquire|lock` | confirm global lock ordering and timeout strategy |

## Scan Priority And Confidence

### Scan Priority

High priority:
- `ML-001`, `ML-003`
- `RC-001`, `RC-005`
- `RL-001`
- `AA-002`

Medium priority:
- `ML-005`
- `RC-002`
- `AA-001`

Context-dependent:
- `ML-002`
- `AA-003`
- `DL-001`, `DL-002`

### Confidence Levels

| Confidence | Meaning |
|------------|---------|
| `HIGH` | pattern match and context confirm |
| `MEDIUM` | pattern match but context is incomplete |
| `LOW` | pattern match with strong false-positive possibility |

### False-Positive Guidance

| Pattern | False Positive Risk | Mitigation |
|---------|---------------------|------------|
| `ML-001` | Medium | check same-file cleanup and framework lifecycle |
| `RC-001` | High | verify async context and shared-state risk |
| `RL-001` | Medium | verify `try-finally`, `with`, or `defer` |
| `AA-002` | Low | direct pattern match is usually reliable |
