---
name: multi-model-oracle
description: Get a merged answer from a configured multi-model pool. Automatically prompt-engineers the query first, queries available models in parallel, and merges results into one best answer. Use when the user wants a best-effort answer, model comparison, multi-model response, oracle answer, or consensus synthesis.
---

# Multi-Model Oracle

Query a configured frontier-model pool in parallel, with automatic prompt engineering and intelligent response merging. The pipeline: optimize the prompt, query available models, merge into one best-effort answer.

## Models

| Model | Provider | Strengths |
|---|---|---|
| Reasoning model | OpenRouter or configured provider | Nuanced reasoning, creative depth, multi-perspective analysis |
| Coding model | OpenRouter or configured provider | Technical precision, breadth, systematic coverage, code |
| Fast synthesis model | OpenAI-compatible provider | Efficient prompt engineering, normalization, and merge synthesis |

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

**Stage 1 - Prompt Engineering**: Detects query intent (code, creative, analysis, research, reasoning, practical, general), then uses the configured fast synthesis model to create model-specific prompt variants optimized for each model role.

**Stage 1.5 - Parallel Prompt Engineering**: The model-specific prompt variants are generated in parallel via ThreadPoolExecutor for faster startup.

**Stage 2 - Parallel Query**: Sends prompts simultaneously via ThreadPoolExecutor. Each model has retry logic (2 retries with exponential backoff) and 180s timeout. If a model fails, the pipeline continues with remaining models.

**Stage 3 - Intelligent Merge**: Feeds all successful responses to the configured fast synthesis model with a merge prompt that extracts the strongest elements, resolves contradictions, eliminates redundancy, and produces a unified answer that reads as if from a single expert. Has its own retry logic (2 retries) with concatenation fallback if merge fails.

**Stage 4 - Output**: Formats the final result with metadata (models used, timing, token counts).

## Requirements

- `OPENROUTER_API_KEY` environment variable for OpenRouter-backed model roles
- `OPENAI_API_KEY` environment variable for OpenAI-compatible synthesis models
- Python packages: `requests`, `openai`

## Prompt Engineering Reference

For details on the intent detection patterns, model-specific tailoring, and merge strategy, read `references/prompt_engineering.md`.
