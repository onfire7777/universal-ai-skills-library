# Multi-step Composition (Phase 4, plan §3.6)

> Additive to the working-set `compose` engine. `compose` (default) assembles a
> flat top-N working set for one task; **`compose --pipeline`** decomposes a
> *multi-step* prompt into an ordered, capability-typed **DAG of skills**.

This is the plan §3.6 deliverable: `"scrape → summarize → post"` becomes an
ordered pipeline of routed skills. It is built **on the canonical engine seam** —
`ComposePlan` calls `skillservice.Route` and `skillservice.Load`, adding no
parallel routing logic — so it inherits the semantic+lexical+guardrail pipeline
and the offline/public-safe guarantee.

## Invariants (test-covered)

- **Deterministic** — same prompt + manifest ⇒ identical plan. Segmentation is
  pure string processing; per-step routing reuses `Route`; output ordering comes
  from ordered slices only. Verified byte-identical across repeated runs.
- **Public-safe** — no network/remote-LLM; third-party external skills are
  **never inlined by default** (pointer-only unless `--allow-external`).
- **Context-light** — steps carry skill names + one-line load pointers, not
  bodies. Bodies inline lazily, only with `--load`, only within `--budget`.

## CLI

```sh
# Plan a pipeline (pointers only — context-light default):
skill-router skills compose --pipeline "scrape a site with crawl4ai, then summarize this chat, then make a card"

# Structured DAG:
skill-router skills compose --pipeline --json "<prompt>"

# Inline bodies within a token budget (over-budget steps stay pointers):
skill-router skills compose --pipeline --load --budget 4000 "<prompt>"
```

| Flag | Default | Purpose |
|------|---------|---------|
| `--pipeline` | `false` | Switch from working-set to multi-step DAG mode. |
| `--max-steps` | `8` | Pipeline-length guardrail; excess sub-tasks are dropped (`truncated`). |
| `--budget` | `6000` | Token ceiling for inlined bodies. |
| `--load` | `false` | Inline bodies within budget. Off ⇒ pointers only (0 tokens). |
| `--allow-external` | `false` | Permit inlining third-party external skill bodies. |

## MCP

The `compose` MCP tool (`skill-router serve`) accepts `pipeline: true`, returning
the `ComposePlanResult` DAG instead of the flat working set.

## Segmentation & capability typing

Sub-tasks split on arrows (`->`, `=>`, `→`, `➜`, `⇒`), sequence words
(`then`, `and then`, `after that`, `followed by`, `finally`, `lastly`), numbered/
bullet list markers, newlines, and `;`. A bare `and` never splits (it joins noun
phrases). Each node carries a typed capability — inferred deterministically
(`scrape → web.fetch`, `summarize → text.summarize`, …, else `general`) until the
corpus is backfilled with frontmatter `capabilities` (§3.2). Unroutable sub-tasks
become honest **gap** nodes rather than forced low-confidence matches.
