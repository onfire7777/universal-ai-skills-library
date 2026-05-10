Purpose: Use this file when you are selecting models, designing caching and batching, or setting budget and token-efficiency rules.

## Contents
- Token economics
- Cost reduction stages
- Effort tuning
- Model registry and dashboard
- Oracle gates

# Cost Optimization

## Token Economics

| Model | Input / 1M | Output / 1M | Speed | Quality |
|-------|------------|-------------|-------|---------|
| Claude Opus 4.6 | `$15.00` | `$75.00` | Slow | Highest |
| Claude Sonnet 4.6 | `$3.00` | `$15.00` | Medium | High |
| Claude Haiku 4.5 | `$0.80` | `$4.00` | Fast | Good |
| GPT-4o | `$2.50` | `$10.00` | Medium | High |
| GPT-4o-mini | `$0.15` | `$0.60` | Fast | Good |

Formula:
- monthly cost = `(input cost + output cost) * requests per day * 30`

## 5-Stage Cost Reduction Framework

| Stage | Typical savings | Core moves |
|-------|-----------------|------------|
| prompt optimization | `-10%` to `-30%` | remove redundant context, trim few-shot, right-size `max_tokens` |
| model routing | `-30%` to `-60%` | start cheap, escalate on validation failure |
| caching | `-10%` to `-90%` | prompt cache, exact cache, semantic cache |
| batching | `-20%` to `-40%` | batch or continuously batch async work |
| error cost elimination | `-5%` to `-10%` | kill retry waste and failed re-generation |

Rules:
- start with the cheapest viable model;
- Haiku-first escalation can cut cost materially when validation backs it;
- stable prompt shapes are required for strong prompt-cache hit rates.

### Cache Expectations

| Strategy | Expected hit rate / effect |
|----------|----------------------------|
| prompt cache | up to `90%` input-cost reduction |
| exact cache | `40-70%` hit rate for repeated FAQ/classification |
| semantic cache | `10-30%` hit rate for similar chat queries |

## Effort Parameter

| Effort | Default use |
|--------|-------------|
| `low` | simple extraction and classification |
| `medium` | default for most production tasks |
| `high` | complex reasoning and agentic work |
| `max` | deep research only |

Rule: do not default to `high`; `medium` is the standard starting point.

## Registry And Monitoring

Keep per-model records for:
- provider
- exact model name
- prompt version
- status
- quality metrics
- latency
- cost per 1k queries
- rollback target

Dashboard thresholds:
- daily spend `> 120%` of budget
- cost per query `> 2x` baseline
- cache hit rate `< 50%` of expected
- wasted-token cost `> 5%` of total
- unexpected thinking-token spikes

## Checklist

- cheapest viable model
- prompt caching enabled for stable prompts
- app-level caching where repeat traffic exists
- realistic `max_tokens`
- explicit effort choice
- token waste monitored
- batching used when latency budget allows
- hard limits and alerts per endpoint
- cost tracked per feature, not only globally

## Oracle Gates

- no cost estimate -> require budget projection
- Opus for simple tasks -> require routing justification
- no caching strategy -> require ROI analysis
- default `max_tokens` without need analysis -> require right-sizing
