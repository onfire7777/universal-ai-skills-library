---
name: multi-model-oracle
description: Get the ultimate merged answer from the 3 best AI models (Anthropic Claude Opus 4.6, OpenAI GPT-5.4, Manus gpt-4.1-mini). Automatically prompt-engineers the query first, then queries all models in parallel, then merges into one best answer. Use when the user wants the best possible answer, wants to compare models, asks for a multi-model response, or says oracle, ultimate answer, or best answer.
---

# Multi-Model Oracle

Query the 3 best frontier AI models in parallel, with automatic prompt engineering and intelligent response merging. The pipeline: optimize the prompt, query all models, merge into one ultimate answer.

## Models

| Model | Provider | Strengths |
|---|---|---|
| Claude Opus 4.6 | Anthropic (OpenRouter) | Nuanced reasoning, creative depth, multi-perspective analysis |
| GPT-5.4 | OpenAI (OpenRouter) | Technical precision, breadth, systematic coverage, code |
| gpt-4.1-mini | Manus (built-in) | Fast, efficient, structured tasks, prompt engineering engine |

## Quick Start

```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py "Your question here"
```

## Workflow

### Step 1: Determine the Query

Accept the query from the user. It can be any type: code, creative writing, analysis, research, reasoning, practical advice, or general questions.

For long or complex prompts, save to a file first:
```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py --file /path/to/prompt.txt
```

### Step 2: Run the Oracle

Standard mode (prompt engineering + parallel query + merge):
```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py "query"
```

With individual model responses visible:
```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py "query" --show-individual
```

With prompt engineering details visible:
```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py "query" --show-prompts
```

Raw mode (skip prompt engineering, send query as-is):
```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py "query" --raw
```

Save output to file:
```bash
python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py "query" -o result.md
```

Read query from stdin (pipe input):
```bash
echo "Your question" | python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py --stdin
```

### Step 3: Deliver the Result

The merged answer is printed to stdout and optionally saved to a file. Present the merged answer to the user. If `--show-individual` was used, also share the individual responses for comparison.

## Pipeline Stages

**Stage 1 — Prompt Engineering**: Detects query intent (code, creative, analysis, research, reasoning, practical, general), then uses gpt-4.1-mini to create 3 model-specific prompt variants optimized for each model's strengths.

**Stage 1.5 — Parallel Prompt Engineering**: The 3 model-specific prompt variants are generated in parallel via ThreadPoolExecutor for faster startup.

**Stage 2 — Parallel Query**: Sends all 3 prompts simultaneously via ThreadPoolExecutor. Each model has retry logic (2 retries with exponential backoff) and 180s timeout. If a model fails, the pipeline continues with remaining models.

**Stage 3 — Intelligent Merge**: Feeds all successful responses to gpt-4.1-mini with a merge prompt that extracts the strongest elements, resolves contradictions, eliminates redundancy, and produces a unified answer that reads as if from a single expert. Has its own retry logic (2 retries) with concatenation fallback if merge fails.

**Stage 4 — Output**: Formats the final result with metadata (models used, timing, token counts).

## Requirements

- `OPENROUTER_API_KEY` environment variable (for Claude Opus 4.6 and GPT-5.4)
- `OPENAI_API_KEY` environment variable (for Manus gpt-4.1-mini)
- Python packages: `requests`, `openai`

## Prompt Engineering Reference

For details on the intent detection patterns, model-specific tailoring, and merge strategy, read `references/prompt_engineering.md`.
