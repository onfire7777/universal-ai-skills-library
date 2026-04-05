---
name: model-selector
description: >
  Set and manage the default AI model Manus uses as its backbone LLM. Supports manual model
  selection, automatic best-model selection based on task type, and a cached/refreshable leaderboard
  of the latest frontier models from OpenRouter. Toggle on/off via chat commands. Use when the user
  says "model-selector", "set model", "change model", "best model", "auto model", "switch model",
  "model-selector on", "model-selector off", "list models", "show models", "model status",
  "recommend model", "refresh models", or asks which AI model to use for a task.
---

# Model Selector

Control which AI model Manus uses as its backbone LLM. Two modes: **manual** (user picks a specific model) or **auto** (dynamically selects the best model based on task category via keyword detection).

## Commands

| Command | Action |
|---|---|
| `/model-selector set <model_id>` | Set a specific model (e.g., `anthropic/claude-opus-4.6`). Validates against cache; warns if unknown but accepts (supports private models). |
| `/model-selector auto [task]` | Auto-select best model. If task provided, detects category via keywords and selects top-ranked model. |
| `/model-selector on` | Enable model selection. |
| `/model-selector off` | Disable (Manus uses its default backbone model). |
| `/model-selector status` | Show current config: enabled state, mode, model, reason, cache freshness. |
| `/model-selector refresh` | Force-refresh cached model leaderboard from OpenRouter API. Rate limit: max 1 per minute. |
| `/model-selector list [category]` | List top models by category (coding/reasoning/research/creative/math/general/all) with pricing. |
| `/model-selector recommend <task>` | Recommend best model for a task without changing config. Shows rationale and top alternatives. |

## Execution

All commands are executed by running the CLI script:
```bash
python3 scripts/model_selector.py <command> [args]
```

The script path is relative to this skill's directory. The absolute path in the Manus sandbox is:
```bash
python3 /home/ubuntu/skills/model-selector/scripts/model_selector.py <command> [args]
```

When a user sends a command (e.g., `/model-selector auto build a website`), Manus should:
1. Parse the command and arguments from the user message.
2. Invoke the CLI script via shell with the appropriate command and arguments.
3. Read the script output and report the result to the user.
4. For errors, display the error message and suggest remediation.

## Architecture

### Persistent Files

**Config** (`~/.model_selector_config.json`):
```json
{
  "enabled": true,
  "mode": "manual|auto",
  "manual_model": "model_id|null",
  "auto_category": "category|null",
  "last_selected_model": "model_id|null",
  "last_selected_reason": "human-readable explanation",
  "updated_at": "ISO-8601 timestamp"
}
```

**Cache** (`~/.model_selector_cache.json`):
```json
{
  "cached_at": 1712345678,
  "cached_at_human": "ISO-8601",
  "total_models": 349,
  "frontier_models": [{"id": "...", "provider": "...", "context_length": 0, "price_per_m_output": 0}],
  "model_lookup": {"model_id": {"id": "...", "name": "...", "context_length": 0, "pricing_prompt": "0", "pricing_completion": "0"}}
}
```

### Task Category Detection

Auto mode classifies tasks into six categories by word-boundary regex matching: **coding**, **reasoning**, **research**, **creative**, **math**, or **general** (fallback). Each category has a curated ranking of the best models based on benchmark data. No AI classification is used — keyword detection is deterministic, fast, and credit-free. Specialized categories always win ties over "general".

### Manus Integration

At the start of each task, Manus should:
1. Read `~/.model_selector_config.json`.
2. If `enabled` is `true`, note the `last_selected_model` and `mode`.
3. Use this model preference when making decisions about which LLM to invoke as the backbone.
4. If `enabled` is `false`, config is missing, or config is corrupted, use the default backbone model.

Other skills' model preferences remain unaffected.

## Workflow

1. **Parse command** from user message.
2. **Execute** via the CLI script using shell.
3. **Report result** — which model was selected, why, detected category (if auto), and alternatives.
4. **For errors**: Respond with clear error message, usage help, and suggested remediation.

## Error Handling

- **Cache missing/corrupted**: Auto-reinitialize on next use or prompt user to run `refresh`.
- **Refresh fails**: Keep stale cache, warn user, suggest retry.
- **Unknown model set**: Warn but accept (may be private/custom model).
- **Invalid command**: Show usage help with all supported commands.
- **File permission errors**: Suggest remediation steps.
- **Config corrupted**: Fall back to defaults, notify user.
- **Malformed API data**: Skip entries with missing IDs or non-numeric pricing gracefully.

## Performance Targets

- Read-only commands (status, list, recommend): under 1 second (no network calls).
- Write commands (set, auto): under 5 seconds (may refresh cache if stale).
- Cache refresh: under 30 seconds.
- Cache auto-expires after 24 hours; refreshes on next write command.

## Security

- Never expose `OPENROUTER_API_KEY` in outputs, logs, or error messages.
- Atomic file writes to prevent config/cache corruption.
- Preserve manual model settings unless explicitly changed by user.

## Requirements

- `OPENROUTER_API_KEY` environment variable (for cache refresh and model validation).
- `requests` Python package (pre-installed in Manus sandbox).

## Benchmark Reference

For detailed benchmark data, tier rankings, and cost analysis, read:
`references/model_benchmarks.md`

Only read this file when the user asks about specific benchmark data, model comparisons, or pricing details.

## Constraints

- Controls **backbone LLM preference only** — does not alter models used by other skills.
- Uses **word-boundary regex task detection** exclusively — no AI classification to avoid extra costs.
- Cache refresh rate-limited: max 1 per minute, auto-refresh every 24 hours.
- Manual model settings preserved unless explicitly changed by user.
- Graceful degradation: if disabled or config/cache unavailable, Manus uses its default model.
