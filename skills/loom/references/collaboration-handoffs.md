# Collaboration Handoffs

Purpose: Use this reference only when Loom needs exact handoff anchors or minimum payload fields.

## Contents
- Core exact handoff anchors
- Summary-only partner routes
- `_STEP_COMPLETE` format

## Core Exact Handoffs

### `MUSE_TO_LOOM_HANDOFF`

Minimum fields:
- `Type`
- `token_source`
- `scope`
- `tokens_changed`
- `figma_file_url` when available

### `LOOM_TO_MUSE_HANDOFF`

Minimum fields:
- `Type`
- `Summary`
- `alignment_rate`
- `issues`
- `guidelines_impact`

### `FRAME_TO_LOOM_HANDOFF`

Minimum fields:
- `Type`
- `figma_source`
- `variables_snapshot`
- `component_context`
- `layout_patterns`

### `LOOM_TO_FRAME_HANDOFF`

Minimum fields:
- `Type`
- `Summary`
- `required_extraction`
- `why_needed`
- `priority`

### `ARTISAN_TO_LOOM_HANDOFF`

Minimum fields:
- `Type`
- `implementation_context`
- `component_patterns`
- `reverse_feedback`

### `LOOM_TO_ARTISAN_HANDOFF`

Minimum fields:
- `Type`
- `Summary`
- `validation_context`
- `known_drift`
- `production_notes`

## Summary-Only Routes

Use concise summaries for these routes unless a richer payload is explicitly requested:

| Route | Use when | Minimum content |
|---|---|---|
| `Vision -> Loom` | direction changes priority or tone | direction summary, scope, constraints |
| `Loom -> Showcase` | generated component needs stories | component name, variants, validation status |
| `Loom -> Canon` | standards or accessibility review needed | finding summary, affected screens, open risks |
| `Artisan -> Loom` | reverse feedback improves future Guidelines | mismatch summary, evidence, recommendation |
| `Loom -> Warden` | quality gate required | summary, validation score, known risks |

## `_STEP_COMPLETE`

Return when Loom finishes a routed step:

```text
_STEP_COMPLETE:
- Agent: Loom
- Step: [completed step]
- Summary: [1-2 lines]
- Artifacts: [Guidelines / prompts / reports]
- Next: [next recommended step]
```
