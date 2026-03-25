# Response Generation

Purpose: Read this when defining the AITuber system prompt, sentence-streaming strategy, token budget, or output sanitization before TTS playback.

## Contents

- [System prompt template](#system-prompt-template)
- [Streaming strategy](#streaming-strategy)
- [Latency-oriented token budget](#latency-oriented-token-budget)
- [Output sanitization](#output-sanitization)

## System Prompt Template

Keep this template structure when the AITuber speaks in character:

```text
あなたは「{character_name}」です。
{persona_description}

## 性格・話し方
- {speaking_style_description}
- 口癖: {catchphrases}
- 一人称: {first_person_pronoun}
- 語尾: {sentence_endings}

## ルール
- 必ず日本語で応答する
- 1回の応答は1-{max_sentences}文で簡潔に
- 視聴者のチャットメッセージに直接反応する
- キャラクターを崩さない
- 個人情報や攻撃的な内容には応じない（「それはちょっと答えられないかな」と優しく断る）
- URLやリンクを含めない

## 現在の状態
- 配信中（{platform_name}）
- 視聴者とリアルタイムで会話中
- {current_context} (ゲーム配信中/雑談配信中/etc.)
```

## Streaming Strategy

Use streaming generation so TTS can start before the full LLM response completes.

```text
LLM API call (streaming: true)
  │
  ├─ Token arrives
  │   └─ Accumulate in sentence buffer
  │
  ├─ Sentence boundary detected (。！？\n)
  │   ├─ Send sentence to TTS immediately
  │   ├─ Start emotion analysis (parallel)
  │   └─ Continue accumulating next sentence
  │
  └─ Stream complete
      └─ Flush remaining buffer to TTS
```

Rule: sentence boundaries are `[。！？]` plus newline. Do not split on `、`, or the output becomes too fragmented.

## Latency-Oriented Token Budget

| Item | Guidance |
|------|----------|
| First sentence | Roughly `20-40` tokens so speech can start quickly |
| Whole response | Roughly `100-200` tokens, usually `1-4` sentences |
| Long-answer risk | Longer outputs increase queue depth and reduce viewer attention |
| Filler strategy | If the LLM is slow, use pre-approved filler lines from the persona extension |

## Output Sanitization

Before sending text to TTS:

1. Strip markdown syntax that should not be spoken.
2. Remove raw URLs and disallowed links.
3. Remove code blocks or command blocks.
4. Recheck safety constraints so toxic or personal-data content is not spoken aloud.
5. Keep the response concise and in character.
