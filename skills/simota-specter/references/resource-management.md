# Resource Management Patterns

Purpose: detect resource leaks, validate cleanup discipline, and size pools safely.

## Contents
- Resource classes
- Language patterns
- Pool thresholds
- Review checklist
- Anti-patterns
- Incident lessons

## Resource Classes

| Resource | Typical limit | Failure mode when exhausted | Detection difficulty |
| --- | --- | --- | --- |
| Memory | Process memory limit | OOM kill, crash | Medium |
| File handle | `1,024-65,536` per process | `EMFILE`, file ops fail | High |
| DB connection | Pool-size dependent | Timeouts, cascading outages | Medium |
| Network socket | `~65,535` per host | Connection failures, timeouts | High |
| Thread/worker | Pool-size dependent | New work stalls | Medium |

## Language Patterns

### JavaScript / TypeScript

```typescript
// TypeScript 5.2+ explicit resource management surface
class DatabaseConnection {
  [Symbol.dispose]() {
    this.close();
  }
}

// Current safe default: try-finally
async function query() {
  const conn = await pool.getConnection();
  try {
    return await conn.query('SELECT ...');
  } finally {
    conn.release();
  }
}
```

### Python

```python
with open('file.txt') as f:
    data = f.read()

async with aiohttp.ClientSession() as session:
    async with session.get(url) as response:
        data = await response.json()
```

### Go

```go
func query() error {
    conn, err := pool.Get()
    if err != nil { return err }
    defer conn.Close()
    // ...
    return nil
}
```

## Pool Thresholds

| Parameter | Recommended range | Notes |
| --- | --- | --- |
| `maxConnections` | `10-50` | Depends on workload and DB capacity |
| `minConnections` | `2-5` | Keep warm capacity without overcommitting |
| `acquireTimeout` | `5000-30000ms` | Fail fast enough to avoid request pileups |
| `idleTimeout` | `30000-600000ms` | Release idle connections predictably |
| `leakDetectionThreshold` | `2000-5000ms` | Warn on long-held pooled resources |
| `maxLifetime` | `1800000ms (30min)` | Rotate long-lived connections |

### Leak Detection

- Record acquire timestamps and emit warnings when hold time exceeds `leakDetectionThreshold`.
- Capture stack traces on acquire where the pool supports it.
- Track acquire/release symmetry and periodically inspect pool utilization such as `pool.numUsed()`.

### Exhaustion Cascade

1. Request A acquires a connection and never releases it.
2. The pool slowly depletes.
3. New requests fail on `acquireTimeout`.
4. Health checks also fail to acquire connections.
5. The load balancer marks the instance unhealthy.
6. Traffic shifts to the remaining instances.
7. The same leak repeats and the failure becomes systemic.

Real incident class: DB pool exhaustion impacted `75%` of global users.

## Review Checklist

### Acquisition and cleanup

- Every acquisition has a matching release.
- Cleanup is guaranteed with `finally`, `defer`, `with`, `async with`, or equivalent.
- Early returns do not bypass cleanup.
- Exception paths still release resources.
- Loop-scoped acquisitions are not released only after the loop finishes.

### Connections and streams

- DB clients define timeouts.
- HTTP clients define timeouts.
- WebSocket flows include reconnect and cleanup logic.
- Streams handle both piping and error paths.
- Temporary files are removed after use.

### Events and timers

- `addEventListener` has matching `removeEventListener`.
- `setInterval` has matching `clearInterval`.
- Subscriptions have matching unsubscribe/dispose.
- React `useEffect` returns cleanup when needed.
- `AbortController` is used where fetches can outlive the caller.

### Pools and caches

- Pool size matches workload and infrastructure limits.
- Caches have a size limit or TTL.
- Global collections define eviction.
- `WeakMap` or `WeakRef` is considered where ownership is weak.

## Anti-Patterns

| Anti-pattern | Symptom | Safer rule |
| --- | --- | --- |
| Happy Path Only | Cleanup exists only on success | Require `try-finally` or equivalent |
| Early Return Bypass | Return exits before cleanup | Put cleanup in a guaranteed finalizer |
| Nested Resource | Only the outer resource is released | Guard each resource independently |
| Pooled but Unreleased | Pool usage climbs over time | Enable leak detection and audit symmetry |
| Fire-and-Forget Close | Async close is not awaited | Await close/dispose completion |
| Conditional Cleanup | Cleanup runs only on some branches | Make cleanup unconditional |

## Incident Lessons

| Incident | Root cause | Operational lesson |
| --- | --- | --- |
| `CVE-2024-21626` | `runc` file descriptor leak | Resource leaks can become security issues |
| Global DB outage | Pool exhaustion cascade | Small leaks can create fleet-wide failures |
| Google Shakespeare | Leaks under load | High traffic amplifies latent leaks |
| CheckMK handle leak | Windows handle accumulation | Long-lived agents need explicit cleanup audits |

Source: Propel resource leak guide, DZone DB leak detection, DoHost connection pool exhaustion, Etleap DB connection leak prevention.
