# Composition (Phase 4)

> Implements [`ARCHITECTURE_IMPROVEMENT_PLAN.md`](ARCHITECTURE_IMPROVEMENT_PLAN.md) §3.6 / roadmap Phase 4.
> Status: **shipped** in `skill-router-cli` (`cmd/skills/route_compose.go`).

Composition turns a multi-step prompt — `"scrape → summarize → post"` — into an
ordered, capability-typed **DAG of skills**, so an agent can plan a whole workflow
in one call instead of routing each step by hand.

It is the value-add layer on top of routing. The default path stays single-skill;
composition is **explicit opt-in**.

## Invariants

These hold by construction and are covered by tests:

- **Deterministic given inputs.** Same prompt + same manifest ⇒ byte-identical DAG.
  Segmentation is pure string processing; per-step routing reuses the existing
  deterministic lexical preflight; output ordering comes from ordered slices only
  (never map iteration).
- **Public-safe.** No network, no remote LLM — composition reuses the offline
  router. Third-party (external) skill bodies are **never inlined by default**;
  they appear as pointers only, behind an explicit `--allow-external` opt-in.
- **Context-light.** The DAG carries skill names and a one-line load command per
  step, not skill bodies. Bodies load lazily, and only within an explicit token
  budget.

## CLI

```sh
# Plan a pipeline (pointers only — the context-light default):
skill-router skills compose "scrape a site with crawl4ai, then summarize this chat, then make a printable card"

# Same thing via the route command's flag:
skill-router skills route --compose "<prompt>"

# Structured DAG for programmatic callers:
skill-router skills compose --json "<prompt>"

# Inline bodies, but only within a token budget (over-budget steps stay pointers):
skill-router skills compose --load --budget 4000 "<prompt>"

# Opt into composing third-party external skills:
skill-router skills compose --allow-external "<prompt>"
```

| Flag | Default | Purpose |
|------|---------|---------|
| `--max-steps` | `8` | Hard cap on pipeline length (context-budget guardrail). Excess sub-tasks are dropped and the plan is flagged `truncated`. |
| `--budget` | `6000` | Estimated token ceiling for lazily inlined bodies. |
| `--load` | `false` | Inline bodies within budget. Off ⇒ pointers only (0 tokens used). |
| `--allow-external` | `false` | Permit composing third-party external skill bodies. |
| `--json` | `false` | Emit the DAG as structured JSON. |

## Segmentation

A prompt is split into ordered sub-tasks on these deterministic connectors:

- arrows: `->`, `=>`, `→`, `➜`, `⇒`
- sequence words: `then`, `and then`, `after that`, `followed by`, `finally`, `lastly`
- list / line structure: numbered (`1.`, `2)`) and bullet (`-`, `*`) markers, newlines, `;`

A bare `and` is intentionally **not** a connector — it joins noun phrases
("organize **and** rename files") rather than tasks.

## Capability typing

Each node carries a typed capability. When the routed skill's frontmatter declares
`capabilities:` (§3.2), those are used verbatim. The corpus is not yet backfilled,
so until then composition falls back to deterministic per-step inference
(`scrape → web.fetch`, `summarize → text.summarize`, `post → message.publish`,
`write → content.generate`, `test → test.author`, …, else `general`). Backfilled
metadata flows through automatically with no code change.

## JSON shape

```json
{
  "prompt": "...",
  "multi_step": true,
  "steps": [
    {
      "index": 0,
      "text": "scrape a site with crawl4ai",
      "decision": "route",
      "skill": "crawl4ai",
      "source": "canonical",
      "score": 353,
      "capabilities": ["web.fetch"],
      "load": "skill-router skill crawl4ai",
      "est_tokens": 1392,
      "loaded": false
    }
  ],
  "edges": [{ "from": 0, "to": 1 }],
  "truncated": false,
  "budget_tokens": 6000,
  "est_tokens_used": 0,
  "notes": ["composition is opt-in; bodies load lazily via each step's load command"]
}
```

Sub-tasks with no confident route become honest **gap** nodes (`decision` ≠
`route`, empty `skill`) with a note, rather than a forced low-confidence match.

## Boundary with the MCP `compose` tool

`buildComposition(prompt, opts)` is a pure function with no I/O beyond the manifest
read. The Phase 2 `cmd/serve` daemon will expose the MCP `compose` tool (§3.4) by
calling it directly — no reimplementation. This file documents the engine the tool
will wrap.
