# Model Benchmarks Reference

This file contains curated benchmark data for frontier AI models. It serves as a fallback when the live cache is unavailable and as context for understanding model strengths.

**Last updated**: April 2026

## Frontier Model Tiers (April 2026)

| Model | Provider | Coding | Reasoning | Research | Creative | Math | Context | Cost (in/out per 1M) |
|---|---|---|---|---|---|---|---|---|
| Claude Opus 4.6 | Anthropic | S | S | A | S | A | 1M | $5/$25 |
| GPT-5.4 | OpenAI | A | S | A | A | S | 1M | $2.50/$15 |
| Gemini 3.1 Pro | Google | A | A | S | B | A | 1M | $2/$12 |
| Grok 4.20 | xAI | A | A | A | A | A | 2M | $2/$6 |
| DeepSeek V3.2 | DeepSeek | A | B | B | B | A | 164K | $0.26/$0.38 |
| Claude Sonnet 4.6 | Anthropic | A | A | A | A | B | 1M | $3/$15 |
| GPT-5.4 Mini | OpenAI | B | B | B | B | B | 400K | low |
| GPT-5.4 Nano | OpenAI | C | C | C | C | C | 400K | lowest |

**Tier key**: S = Best-in-class, A = Excellent, B = Good, C = Adequate

## Category-Specific Leaders

### Coding
1. **Claude Opus 4.6** — 80.8% SWE-bench, best for complex multi-file codebases and refactoring
2. **GPT-5.4** — 57.7% SWE-bench Pro, 75.1% Terminal-Bench, strong tool use
3. **Grok 4.20** — 2M context enables full-repo analysis

### Reasoning & Analysis
1. **Claude Opus 4.6** — Strongest chain-of-thought and structured analysis
2. **GPT-5.4** — Excellent structured reasoning with tool integration
3. **DeepSeek R1-0528** — Dedicated reasoning model, strong on complex logic

### Research & Synthesis
1. **Gemini 3.1 Pro** — 1M context, best for large document synthesis
2. **Grok 4.20** — 2M context for massive corpus analysis
3. **Claude Opus 4.6** — Best synthesis quality

### Creative Writing
1. **Claude Opus 4.6** — Best prose quality and creative expression
2. **Claude Sonnet 4.6** — Excellent creative at lower cost
3. **GPT-5.4** — Strong narrative control and style adaptation

### Mathematics
1. **GPT-5.4** — Strongest math benchmarks overall
2. **DeepSeek R1-0528** — Dedicated reasoning, excellent on proofs
3. **Claude Opus 4.6** — Strong mathematical reasoning

## Cost-Efficiency Tiers

### Premium (best quality, higher cost)
- Claude Opus 4.6: $5/$25 per 1M tokens
- GPT-5.4: $2.50/$15 per 1M tokens

### Balanced (great quality, moderate cost)
- Gemini 3.1 Pro: $2/$12 per 1M tokens
- Grok 4.20: $2/$6 per 1M tokens
- Claude Sonnet 4.6: $3/$15 per 1M tokens

### Value (good quality, low cost)
- DeepSeek V3.2: $0.26/$0.38 per 1M tokens
- Meta Llama 4 Maverick: $0.15/$0.60 per 1M tokens

## Data Sources

- SWE-bench / SWE-bench Pro (coding)
- Terminal-Bench (tool use)
- MMLU-Pro / GPQA Diamond (reasoning)
- Chatbot Arena / LMSys (overall quality)
- OpenRouter API (pricing, availability)
- Vellum AI Leaderboard, Onyx LLM Leaderboard, LM Council Benchmarks
