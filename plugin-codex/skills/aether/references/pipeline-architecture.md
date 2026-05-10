# Pipeline Architecture

Purpose: Read this when designing the end-to-end AITuber pipeline, IPC choices, latency budget, queueing, fallback behavior, or resource planning.

## Contents

- [Pipeline overview](#pipeline-overview)
- [Component communication](#component-communication)
- [Latency budget](#latency-budget)
- [Streaming overlap](#streaming-overlap)
- [Orchestrator contract](#orchestrator-contract)
- [Fallback and recovery](#fallback-and-recovery)
- [Resource guide](#resource-guide)

## Pipeline Overview

```text
Chat Listener -> Message Queue -> LLM Engine -> TTS Engine -> Lip Sync -> Avatar Renderer -> OBS Output
Target: Chat message -> audible speech < 3000ms
```

Keep the path asynchronous or streaming. A synchronous end-to-end chain is for debugging only.

## Component Communication

| From -> To | Protocol | Pattern | Use when |
|------------|----------|---------|----------|
| Chat Listener -> Message Queue | EventEmitter / in-process | Pub/Sub | Single-process v1 |
| Message Queue -> LLM Engine | Async call | Pull consumer | Standard request flow |
| LLM Engine -> TTS Engine | EventEmitter / stream pipe | Sentence-level streaming | Production path |
| TTS Engine -> Lip Sync | Direct data pass | Synchronous timing handoff | TTS exposes phoneme timing |
| Lip Sync -> Avatar Renderer | WebSocket / IPC | Push | Avatar runs in another process or runtime |
| Avatar Renderer -> OBS | Browser source / virtual camera / NDI | Video stream | Stage output |
| TTS Engine -> OBS | Audio pipe / browser audio / virtual device | Audio stream | Spoken output |

### IPC Options

| Method | Latency | Complexity | Recommended use |
|--------|---------|------------|-----------------|
| In-process EventEmitter | `< 1ms` | Low | Single Node.js process |
| Child-process IPC | `1-5ms` | Medium | Local multi-process |
| Redis Pub/Sub | `1-10ms` | Medium | Multi-process or multi-machine |
| Local WebSocket | `5-20ms` | Medium | Cross-language components |
| gRPC | `1-5ms` | High | Typed service boundaries |
| Durable queue | `10-50ms` | High | Persistence and retry heavy workflows |

## Latency Budget

### Required Budget

| Stage | Target | Max | Optimization |
|-------|--------|-----|--------------|
| Chat listener | `200ms` | `500ms` | Prefer WebSocket over polling where available |
| Message queue | `50ms` | `100ms` | In-memory priority queue |
| LLM first token | `500ms` | `800ms` | Streaming API and prompt cache |
| LLM sentence completion | `1000ms` | `1500ms` | Sentence-boundary detection |
| TTS synthesis | `800ms` | `1200ms` | Warm engine, chunked synthesis if possible |
| Lip sync computation | `50ms` | `100ms` | Use TTS timing data directly |
| Avatar frame update | `100ms` | `200ms` | Keep renderer light |
| OBS encoding and output | `250ms` | `400ms` | Prefer hardware encoding |
| Pipeline total | `~2950ms` | `~4800ms` | Launch gate must stay below `3000ms` |

## Streaming Overlap

Key rule: start TTS on sentence boundaries while the LLM continues generating later sentences. This is the main path to perceived latency below the hard budget.

## Orchestrator Contract

```typescript
interface PipelineConfig {
  chatListener: ChatListenerAdapter;
  messageQueue: MessageQueue;
  llmEngine: LLMAdapter;
  ttsEngine: TTSAdapter;
  lipSync: LipSyncEngine;
  avatarController: AvatarController;
  obsController: OBSController;
}

interface PipelineMetrics {
  chatToSpeechLatency: number;
  ttsQueueDepth: number;
  droppedFrames: number;
  messagesProcessed: number;
  messagesSkipped: number;
}
```

### Message Priority

| Priority | Message type | Queue rule |
|----------|--------------|------------|
| `1` | Superchat / donation | Always process |
| `2` | Command | High priority |
| `3` | Mention | Prefer over regular chat |
| `4` | Regular chat | Keep newest first when overloaded |

## Fallback and Recovery

| Component | Failure | Fallback | Recovery |
|-----------|---------|----------|----------|
| Chat listener | Rate limit / auth error | Slow polling or pause input | Backoff or token refresh |
| Message queue | Queue full | Drop lowest priority or oldest low-value item | Auto-drain |
| LLM | Timeout / API error | Cached or filler response | Retry with shorter prompt |
| TTS | Synthesis failure | Fallback TTS, then text overlay | Restart or cool down engine |
| Lip sync | Desync | Neutral mouth reset | Re-sync on next sentence |
| Avatar | Render crash / low FPS | Static image or simpler mode | Restart renderer or reduce load |
| OBS | WebSocket loss / overload | Preserve state and reconnect | Exponential backoff |

### Circuit-Breaker Defaults

- `LLM`: open after `3` consecutive failures, cool down for `30s`
- `TTS`: open after `2` consecutive failures, cool down for `15s`
- `Avatar`: restart immediately after `1` crash, then cool down for `10s`

## Resource Guide

| Component | CPU | Memory | GPU | Notes |
|-----------|-----|--------|-----|-------|
| Chat listener | Low | ~100MB | No | I/O bound |
| Queue / orchestrator | Low | ~200MB | No | In-memory state |
| API LLM client | Low | ~200MB | No | External call path |
| Local TTS | Medium to high | ~500MB to 1GB | Often helpful | Depends on engine |
| Avatar renderer | Medium | ~500MB to 800MB | Recommended | WebGL-heavy |
| OBS | Medium | ~500MB | Recommended | Hardware encode preferred |
| Single-machine v1 total | Medium | ~3GB | Optional to recommended | Feasible on one machine |
