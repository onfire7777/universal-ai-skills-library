---
name: skill-debugger
description: Deep dual-model debugging of Manus skills using Claude Opus 4.6 and Manus gpt-4.1-mini. Use when asked to debug a skill, find bugs in a skill, review a skill for issues, fix a broken skill, or audit a skill's quality. Also use when a skill is not working correctly or producing unexpected results.
---

# Skill Debugger

Debug Manus skills using two complementary AI models in parallel: **Claude Opus 4.6** (deep code reasoning, security, logic bugs) and **Manus gpt-4.1-mini** (structural integrity, integration quality, trigger accuracy). Findings are merged by consensus — issues confirmed by both models get elevated confidence.

## When to Use

- A skill is broken, crashing, or producing wrong results
- Before deploying a new or modified skill
- To audit an existing skill for hidden bugs
- When a skill triggers at the wrong time or fails to trigger
- After significant changes to a skill's scripts or instructions

## Quick Start

```bash
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py <skill-name>
```

## Workflow

### Step 1: Run the Debugger

Standard analysis (fast, covers most issues):
```bash
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py <skill-name>
```

Deep analysis (extended checks for race conditions, resource leaks, edge cases):
```bash
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py <skill-name> --deep
```

Single-model mode (when one API is unavailable):
```bash
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py <skill-name> --model claude
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py <skill-name> --model manus
```

Debug by path (for skills not in the standard directory):
```bash
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py /path/to/skill-dir
```

### Step 2: Review the Report

The script generates `DEBUG_REPORT.md` inside the skill directory with:
- Model status (which models responded successfully)
- Overall health assessment (healthy / degraded / broken)
- Findings summary table sorted by severity
- Detailed findings with problematic code, explanation, and exact fix
- Consensus badges showing which findings both models agree on

Raw JSON data is saved to `.debug_raw.json` for programmatic access.

### Step 3: Apply Fixes

Read the `DEBUG_REPORT.md` and apply fixes in order of severity (critical first). Each finding includes:
- The exact problematic code to find
- Why it's a problem with a concrete failure scenario
- The exact replacement code or instruction

After applying fixes, re-run the debugger to verify:
```bash
python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py <skill-name>
```

## Analysis Dimensions

The debugger examines five dimensions, with each model contributing its strengths:

| Dimension | Claude Focus | Manus Focus |
|---|---|---|
| Structural | File existence, path correctness | Frontmatter quality, trigger accuracy |
| Scripts | Logic bugs, security, edge cases | Import errors, argument parsing, output format |
| Robustness | Race conditions, resource leaks, memory | Missing error handling, hardcoded paths |
| Security | Command injection, path traversal, key exposure | Input validation, unsafe deserialization |
| Integration | API contract violations | Manus-specific conventions, instruction clarity |

## Prompt Engineering

The debugger uses elite prompt engineering techniques documented in `references/prompt_engineering.md`. Key techniques: role priming, chain-of-thought enforcement, metacognitive verification ("Would I bet $100?"), negative prompting (explicit exclusion list), and severity calibration with concrete criteria.

## Requirements

- `OPENROUTER_API_KEY` environment variable (for Claude Opus 4.6)
- `OPENAI_API_KEY` environment variable (for Manus gpt-4.1-mini)
- Python packages: `requests`, `openai` (auto-installed if missing)
