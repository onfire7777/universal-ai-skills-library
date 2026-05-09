---
name: prompt-engineer
description: Create the absolute best, highly optimized prompt for Manus AI based on user input. Use when the user asks to improve, rewrite, or optimize a prompt, or when they want to learn how to write better prompts for Manus. Requires OPENROUTER_API_KEY and/or OPENAI_API_KEY environment variables.
---

# Manus Prompt Engineer

Transform any raw user idea or prompt into an elite, Manus-optimized prompt. Uses multi-model AI analysis (Claude Opus 4, GPT-4.1, gpt-4.1-mini) and Manus-specific context engineering principles to produce prompts that save credits, prevent errors, and deliver exceptional results.

## Prerequisites

At least one AI API key must be set as an environment variable:

- `OPENROUTER_API_KEY` — Required for multi-model optimization (standard/deep depth). Powers Claude Opus 4 and GPT-4.1 via OpenRouter.
- `OPENAI_API_KEY` — Required for quick mode and as one of the parallel models (gpt-4.1-mini).

If only one key is available, the script automatically falls back to quick mode (single model). If no keys are set, the script exits with an error. Python packages `requests` and `openai` must be installed.

## Quick Start

```bash
python3 /home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py "your raw prompt here"
```

With full analysis and comparison metrics:

```bash
python3 /home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py "your raw prompt here" --show-analysis --show-comparison
```

## Workflow

1. **Understand the User's Goal**: If the request is extremely vague (fewer than 10 words or no clear task), ask ONE clarifying question. Otherwise, proceed with reasonable assumptions — the optimizer handles ambiguity well.
2. **Run the Optimizer**: Execute the script with the user's raw prompt. For long prompts, save to a file first and use `--file`.
3. **Review the Output**: Verify the generated prompt applies Manus-specific principles (One-Shot Density, Agent-First Framing, explicit tool usage).
4. **Deliver to User**: Present the optimized prompt in a markdown code block for easy copy-paste. Briefly explain *why* it is better.

## Advanced Usage

**Modes**: `--mode agent` (autonomous tasks), `--mode chat` (conversational), `--mode project` (reusable workflows). Default: auto-detect.

**Depths**: `--depth quick` (single model, fast), `--depth standard` (parallel models + merge, default), `--depth deep` (intensive multi-model with advanced patterns).

**File Input**: For long prompts, save to a file and use `--file /path/to/prompt.txt`.

**Stdin Input**: Pipe input with `echo "prompt" | python3 scripts/optimize_prompt.py --stdin`. Note: `--stdin` requires piped input; it will hang if run interactively without piping.

**Individual Results**: Use `--show-individual` to see each model's optimization before the merge.

**Save Output**: Use `-o result.md` to save the optimized prompt to a file.

## Prompt Engineering Knowledge Base

For manual prompt construction or teaching users how to write better prompts, consult these bundled reference files (confirmed present in the skill directory):

- **`references/manus_techniques.md`**: Core Manus context engineering principles (One-Shot Density, Agent-First Framing, File System as Memory, Attention Manipulation, Error Prevention) and structured prompt templates for agent tasks, research, and code development.
- **`references/prompt_patterns.md`**: Quick reference for 26+ advanced patterns (Tree of Thoughts, Self-Consistency, Red Team/Blue Team, etc.), intent-specific optimization strategies, and a 10-point optimization checklist.

## Key Principles

1. **One-Shot Density**: Pack everything into one prompt. Follow-ups cost almost as much as the initial prompt.
2. **Agent-First Framing**: Treat Manus as an autonomous agent, not a chatbot. Use phrases like "Take ownership" and "Guide step by step."
3. **Explicit Deliverables**: Always specify exact output formats, file types, and directory structures.
4. **Leverage Tools**: Explicitly mention tools like the shell, browser, or Python when relevant.
