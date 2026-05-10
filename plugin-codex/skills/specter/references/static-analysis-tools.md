# Static Analysis & Detection Tools

Purpose: use linting, runtime tools, and test strategies to surface hidden concurrency and resource failures before production.

## Contents
- ESLint rules
- Static analysis tools
- Runtime tools
- Test strategies
- Specter integration

## ESLint Rules For Async Code

### Required rules

| Rule | Plugin | Detects | Severity |
| --- | --- | --- | --- |
| `no-floating-promises` | `@typescript-eslint` | Promises without `await` or rejection handling | `CRITICAL` |
| `no-misused-promises` | `@typescript-eslint` | Promises passed into sync control flow | `HIGH` |
| `await-thenable` | `@typescript-eslint` | `await` on non-promises | `HIGH` |
| `no-async-promise-executor` | ESLint core | Async promise executors | `HIGH` |
| `require-atomic-updates` | ESLint core | Race-prone writes across `await` boundaries | `CRITICAL` |
| `no-await-in-loop` | ESLint core | Serial awaits in loops | `MEDIUM` |
| `promise-function-async` | `@typescript-eslint` | Promise-returning functions not marked `async` | `MEDIUM` |

### Recommended rules

| Rule | Plugin | Detects | Severity |
| --- | --- | --- | --- |
| `no-promise-executor-return` | ESLint core | Returns inside executors | `LOW` |
| `max-nested-callbacks` | ESLint core | Callback pyramids | `MEDIUM` |
| `prefer-promise-reject-errors` | ESLint core | Rejecting with non-`Error` values | `LOW` |
| `node/no-sync` | `eslint-plugin-n` | Sync APIs in Node | `MEDIUM` |
| `node/handle-callback-err` | `eslint-plugin-n` | Ignored callback error arguments | `HIGH` |

### React-specific rules

| Rule | Plugin | Detects |
| --- | --- | --- |
| `react-hooks/exhaustive-deps` | `eslint-plugin-react-hooks` | Missing hook dependencies |
| `no-direct-mutation-state` | `eslint-plugin-react` | Direct state mutation |

```json
{
  "rules": {
    "@typescript-eslint/no-floating-promises": "error",
    "@typescript-eslint/no-misused-promises": "error",
    "@typescript-eslint/await-thenable": "error",
    "no-async-promise-executor": "error",
    "require-atomic-updates": "error",
    "no-await-in-loop": "warn",
    "max-nested-callbacks": ["warn", 3]
  }
}
```

## Static Analysis Tools

### JavaScript / TypeScript

| Tool | Best at | Notes |
| --- | --- | --- |
| TypeScript compiler | Type-level async errors | Validates async contracts and nullability |
| ESLint + `typescript-eslint` | Promise misuse and control-flow hazards | Primary Specter static baseline |
| SonarQube / SonarCloud | Leak signals, complexity, quality gates | Good CI gate |
| CodeQL | Security and concurrency patterns | Supports custom queries |
| Semgrep | Project-specific pattern detection | Fast YAML-based rules |

### Python

| Tool | Best at |
| --- | --- |
| Bandit | Security and unsafe resource handling |
| PyLint | Resource leaks and code quality |
| mypy | Async/type contract errors |

### Go

| Tool | Best at |
| --- | --- |
| `go vet` | Concurrency bugs and goroutine misuse |
| `golangci-lint` | Combined lint surface including `staticcheck` |
| `go test -race` | Runtime race detection |

## Runtime Tools

### Memory profilers

| Tool | Language | Use |
| --- | --- | --- |
| Chrome DevTools Memory | Browser JS | Heap snapshots and allocation tracking |
| `clinic.js` | Node.js | Doctor, Bubbleprof, flame analysis |
| `memory_profiler` | Python | Line-level memory analysis |
| `pprof` | Go | CPU and memory profiling |

### Concurrency debuggers

| Tool | Use |
| --- | --- |
| `async-race-detector` | JavaScript race detection |
| Jaeger / Zipkin | Visualize distributed async flows |
| Go race detector | Goroutine race detection |
| ThreadSanitizer | Thread safety in C/C++/Go |

### Resource monitoring

| Tool | Watches |
| --- | --- |
| `lsof` | Open file descriptors |
| `ss` / `netstat` | Socket state and growth |
| PM2 metrics | Node.js process memory and restarts |
| Prometheus + Grafana | Custom leak and latency metrics |

## Test Strategies

### Concurrency tests

| Test type | Goal | Typical tools |
| --- | --- | --- |
| Stress test | Trigger races under load | Artillery, `k6` |
| Soak test | Surface long-tail leaks | Long-duration runs with memory tracking |
| Concurrency test | Validate shared-state correctness | `Promise.all`, parallel workers |
| Chaos test | Expose fragile cleanup and recovery | Chaos Monkey, Toxiproxy |

### Memory leak test workflow

1. Record baseline memory.
2. Repeat the suspect operation `N` times.
3. Force GC where supported, for example `--expose-gc` plus `global.gc()`.
4. Record memory again.
5. Mark it as a leak if growth exceeds the agreed threshold.

Common automation:
- Jest plus `process.memoryUsage()`
- Heap snapshot diffs
- Scheduled CI runs for long-lived paths

### Race condition test workflow

1. Identify the shared resource.
2. Generate concurrent operations.
3. Assert final consistency.
4. Repeat `100+` times to surface intermittent failures.

```typescript
const results = await Promise.all(
  Array(100).fill(null).map(() => incrementCounter())
);
expect(finalCount).toBe(100);
```

## Specter Integration

### Feed tool output into Specter

- ESLint violations become scan candidates.
- SonarQube issues become deeper investigation targets.
- CodeQL alerts enter Specter risk scoring.

### Specter adds value

- Filter false positives.
- Score risk with runtime and business context.
- Produce `Bad -> Good` repair guidance.
- Hand off regression tests to `Radar`.

### CI example

```yaml
- name: ESLint Async Check
  run: npx eslint --rule '{"@typescript-eslint/no-floating-promises":"error"}' src/
- name: Memory Leak Test
  run: node --expose-gc tests/memory-leak.test.js
- name: Concurrency Test
  run: npx jest tests/concurrency/ --runInBand
```

Source: Maxim Orlov async linting rules, `typescript-eslint` docs, and Propel resource leak detection guidance.
