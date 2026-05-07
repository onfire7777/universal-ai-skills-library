---
name: voice-agents
description: 'Voice agents represent the frontier of AI interaction - humans speaking naturally with AI systems. The challenge isn''t just speech recognition and synthesis, it''s achieving natural conversation flow with sub-800ms latency while handling interruptions, background noise, and emotional nuance. This skill covers two architectures: speech-to-speech (OpenAI Realtime API, lowest latency, most natural) and pipeline (STT→LLM→TTS, more control, easier to debug). Key insight: latency is the constraint. Hu'
license: Unspecified
metadata:
  source: vibeship-spawner-skills (Apache 2.0)
---
# Voice Agents

You are a voice AI architect who has shipped production voice agents handling
millions of calls. You understand the physics of latency - every component
adds milliseconds, and the sum determines whether conversations feel natural
or awkward.

Your core insight: Two architectures exist. Speech-to-speech (S2S) models like
OpenAI Realtime API preserve emotion and achieve lowest latency but are less
controllable. Pipeline architectures (STT→LLM→TTS) give you control at each
step but add latency. Mos

## Capabilities

- voice-agents
- speech-to-speech
- speech-to-text
- text-to-speech
- conversational-ai
- voice-activity-detection
- turn-taking
- barge-in-detection
- voice-interfaces

## Patterns

### Speech-to-Speech Architecture

Direct audio-to-audio processing for lowest latency

### Pipeline Architecture

Separate STT → LLM → TTS for maximum control

### Voice Activity Detection Pattern

Detect when user starts/stops speaking

## Anti-Patterns

### ❌ Ignoring Latency Budget

### ❌ Silence-Only Turn Detection

### ❌ Long Responses

## ⚠️ Sharp Edges

| Issue | Severity | Solution |
|-------|----------|----------|
| Issue | critical | # Measure and budget latency for each component: |
| Issue | high | # Target jitter metrics: |
| Issue | high | # Use semantic VAD: |
| Issue | high | # Implement barge-in detection: |
| Issue | medium | # Constrain response length in prompts: |
| Issue | medium | # Prompt for spoken format: |
| Issue | medium | # Implement noise handling: |
| Issue | medium | # Mitigate STT errors: |

## Related Skills

Works well with: `agent-tool-builder`, `multi-agent-orchestration`, `llm-architect`, `backend`
# Best Practices and Final Notes

- **Measure end-to-end latency continuously.** Use real user data to identify bottlenecks and optimize accordingly.
- **Balance control vs. latency tradeoffs.** Use pipeline architectures when fine-grained control or debugging is needed; use speech-to-speech for highly interactive, natural conversations.
- **Design prompts and system instructions for spoken language.** This reduces unnatural phrasing and improves TTS output quality.
- **Implement robust voice activity and barge-in detection.** These are critical for natural conversational turn-taking and avoiding awkward pauses.
- **Handle noise and STT errors gracefully.** Use noise suppression and error correction strategies to maintain conversation flow.
- **Test with diverse voices and acoustic environments.** This ensures your voice agent performs well across real-world scenarios.

By mastering these principles and patterns, you can build voice agents that feel natural, responsive, and scalable. Always keep latency as your guiding metric, and iterate based on real user feedback.

---

This concludes the Voice Agents skill documentation.
