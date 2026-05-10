# Memory Leak Diagnosis

Purpose: memory-leak investigation workflow, tool selection, root-cause tracing, and production monitoring thresholds.

## Contents

1. Diagnosis workflow
2. Tool guide
3. Leak-source catalog
4. Weak references
5. Production thresholds

## Diagnosis Workflow

1. Confirm the symptom over time. If memory does not fall after GC, suspect a leak.
2. Capture a baseline heap snapshot.
3. Repeat the suspicious action and capture a second snapshot.
4. Diff by retained size and retainer chain.
5. Hand remediation to `Builder`, then re-run the same workflow.

## Tool Guide

### Node.js

| Tool | Use |
|------|-----|
| `v8.writeHeapSnapshot()` | heap snapshot |
| `--inspect` | Chrome DevTools attach |
| `process.memoryUsage()` | live memory metrics |
| `heapdump` | snapshot on signal |
| `clinic.js` | broad async and CPU diagnosis |
| `why-is-node-running` | active handle and timer inspection |

### Chrome DevTools

Use:
- Heap Snapshot
- Comparison view
- Allocation Timeline
- Allocation Sampling

## Leak-Source Catalog

### Browser / React

- listener accumulation
- detached DOM references
- large closure captures
- uncleared timers
- state updates after unmount

### Node.js / Server

- global object growth
- unconsumed streams
- dynamic module caching
- buffer accumulation
- implicit globals

## Weak References

Use `WeakRef` or `FinalizationRegistry` only for weak caching or weak observer relationships.

Do not use them as the primary cleanup path for:
- file handles
- database connections
- critical resource release

## Production Thresholds

| Metric | Warning Threshold | Action |
|--------|-------------------|--------|
| `heapUsed` | rising continuously for `30 minutes` | capture heap snapshot |
| `rss` | `80%` of process memory budget | consider scale-out and root-cause analysis |
| `external` | `2x` baseline | inspect buffers or native add-ons |
| `GC pause` | `>100ms` | inspect pressure and GC tuning |

`pm2 --max-memory-restart` is a temporary mitigation, not a fix.
