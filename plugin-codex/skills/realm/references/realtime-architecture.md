# Real-Time Update Architecture

Purpose: Guide future evolution of Realm live mode beyond the current polling implementation while preserving the data limits and diff semantics Realm already relies on.

Contents:
- Current runtime baseline
- Protocol recommendation
- Diff protocol
- Layered architecture
- Browser limits
- Upgrade priorities

## Current Runtime Baseline

Current `serve.py` live mode uses HTTP polling:
- `/api/hash` every `3s`
- `/api/activity` every `5s`

Use this reference when you need a more scalable real-time architecture, not for the baseline static/live contract already implemented in `serve.py`.

## Protocol Recommendation

Prefer `SSE` for future live-mode evolution.

| Criterion | WebSocket | SSE | Recommendation |
|---|---|---|---|
| Direction | Bi-directional | Server to client | `SSE` |
| Complexity | Higher | Lower | `SSE` |
| Auto reconnect | Manual | Browser-native | `SSE` |
| HTTP/2 compatibility | Separate connection | Native multiplexing | `SSE` |
| Infra footprint | Dedicated stateful infra | Existing HTTP infra | `SSE` |

## Three-Layer Buffering

| Layer | Cadence | Responsibility |
|---|---|---|
| Data collection | `1s` | Detect changes from git, `PROJECT.md`, `ECOSYSTEM.md`, and journals |
| Batch aggregation | `3s` | Send only diffs after the initial full state |
| Client render | `requestAnimationFrame` | Queue updates and animate transitions smoothly |

## Diff Protocol

### Initial payload

```json
{ "type": "full_state", "agents": [...], "quests": [...], "departments": [...] }
```

### Delta payload

```json
{
  "type": "delta",
  "timestamp": 1709550000,
  "changes": [
    { "path": "agents.builder.xp", "value": 4250 },
    { "path": "quests.Q-42.status", "value": "completed" },
    { "path": "events", "append": { "type": "Victory", "summary": "..." } }
  ]
}
```

## Layered Architecture

| Layer | Responsibility |
|---|---|
| Client layer | Render HTML/Phaser and manage interaction state |
| Data bridge layer | Transport updates and apply diffs |
| State manager layer | Read and write `.agents/realm-state.md`, compute change sets |
| Data collector layer | Watch git, journals, project logs, and ecosystem state |

## Browser Data Limits

- Keep only the latest 100 events on the client.
- Keep only the latest agent stats on the client; historical trends stay server-side.
- Keep only the latest 50 completed quests on the client.

## Upgrade Priorities

| Priority | Upgrade | Cost | Benefit |
|---|---|---|---|
| P0 | Polling to SSE migration | Medium | Lower complexity and better reliability |
| P0 | Diff-only update protocol | Medium | Lower bandwidth and faster refresh |
| P1 | Three-layer buffering | Small | Smoother rendering |
| P2 | Fully layered architecture | Large | Better maintainability |
