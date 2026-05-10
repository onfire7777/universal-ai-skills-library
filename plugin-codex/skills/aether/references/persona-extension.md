# Persona Extension

Purpose: Read this when extending a Cast persona into an AITuber persona with streaming behavior, voice mapping, expression rules, and viewer-interaction settings.

## Contents

- [Extension workflow](#extension-workflow)
- [AITuber extension schema](#aituber-extension-schema)
- [Design notes](#design-notes)

## Extension Workflow

1. Receive a Cast persona or raw character concept.
2. Add streaming-specific traits, voice parameters, expression mapping, and interaction rules.
3. Keep Cast as the canonical persona source of truth.
4. Send behavior updates back through `AETHER_TO_CAST_EVOLVE`; never edit Cast files directly.

## AITuber Extension Schema

Append the following block to the Cast persona when the system needs stream-specific behavior:

```yaml
# AITuber Extension (appended to Cast persona)
aituber:
  streaming_personality:
    reaction_speed: fast | normal | slow     # Chat reaction speed
    humor_style: witty | warm | deadpan      # Humor style
    catchphrases:                             # Signature phrases
      greeting: "はいはーい！"
      farewell: "またねー！"
      thinking: "うーんとね..."
      superchat: "わぁ！ありがとう！"
    filler_phrases:                           # Filler lines for silence avoidance
      - "えっとね..."
      - "ちょっと待ってね"
      - "それはね..."

  voice_mapping:
    tts_engine: voicevox                     # Selected TTS
    speaker_id: 3                            # VOICEVOX speaker ID
    base_params:
      speed: 1.1
      pitch: 0.02
      intonation: 1.2
      volume: 1.0
    emotion_overrides:                        # Per-emotion parameter overrides
      joy: { speed: 1.2, pitch: 0.05, intonation: 1.4 }
      sad: { speed: 0.9, pitch: -0.03, intonation: 0.8 }
      angry: { speed: 1.15, pitch: 0.03, intonation: 1.5 }
      surprised: { speed: 1.25, pitch: 0.08, intonation: 1.6 }

  expression_map:
    framework: live2d                        # live2d | vrm
    # See references/lip-sync-expression.md for full parameter tables

  interaction_rules:
    superchat_always_respond: true
    command_prefix: "!"
    mention_priority: high
    max_response_sentences: 5
    greeting_on_first_message: true
    farewell_on_leave: false                 # Usually not available
```

## Design Notes

| Area | Keep explicit |
|------|---------------|
| Streaming personality | Reaction speed, humor style, greeting/filler behavior, and superchat treatment |
| Voice mapping | Chosen TTS engine, speaker ID, base parameters, and emotion-specific overrides |
| Expression map | Framework choice (`live2d` or `vrm`) and the reference used for parameter tables |
| Interaction rules | Command prefix, mention priority, response-length cap, and greeting behavior |

> Keep user-facing examples in Japanese inside the schema because they affect character behavior and stream output tone.
