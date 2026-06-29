# Architecture

UASL is built on one organizing idea: **the router, the corpus, and the registry are three independent things, connected only through a generated contract.** This page explains the V2 router-first model, the decoupling that makes it work, and the runtime data flow.

Source docs: [`docs/ARCHITECTURE_V2.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/ARCHITECTURE_V2.md) · [`docs/ARCHITECTURE-decoupling.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/ARCHITECTURE-decoupling.md) · [`docs/DESIGN_AND_MESSAGING.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/DESIGN_AND_MESSAGING.md)

---

## The three layers

```
┌──────────────────────────────────────────────────────────────┐
│  AI CLIENTS  (Codex · Claude · Cursor · Hermes · Paperclip …) │
│  each holds only a compact wrapper / instruction file         │
└───────────────────────────┬──────────────────────────────────┘
                            │ calls
                            ▼
┌──────────────────────────────────────────────────────────────┐
│  ROUTER  —  skill-router (Go CLI + MCP server)                │
│  owns NO skill content; resolves prompts → skills             │
└───────────────────────────┬──────────────────────────────────┘
                            │ reads the only contract
                            ▼
┌──────────────────────────────────────────────────────────────┐
│  REGISTRY  —  manifest.json   (generated, deterministic)     │
└───────────────────────────┬──────────────────────────────────┘
                            │ indexes
                            ▼
┌──────────────────────────────────────────────────────────────┐
│  CORPUS  —  skills/<kebab-name>/SKILL.md (+ optional files)  │
│  knows nothing about the router                              │
└──────────────────────────────────────────────────────────────┘
```

| Layer | Lives in | Responsibility | Knows about |
|---|---|---|---|
| **Router** | `skill-router-cli/` | Score prompts, preflight, load skills, serve MCP | The registry only |
| **Registry** | `manifest.json` (+ `docs/build_manifest.json`) | Machine-readable index of every skill | Nothing — it is data |
| **Corpus** | `skills/` | Hold the actual skill content | Nothing |

Because the corpus knows nothing about the router and the router only reads the registry, you can: regenerate the registry without touching skills, swap the router implementation (see [Node → Go](Node-to-Go-Migration.md)) without touching the corpus, and point any client at the same canonical corpus instead of duplicating it.

---

## The router-first pipeline: search → route → load

The V2 engine exposes **four canonical verbs**, all served by one internal engine (`skillservice`) regardless of whether they're reached via the CLI or the MCP server:

| Verb | What it does |
|---|---|
| **`route`** | Given a prompt, decide whether a skill is needed and which one (`route` / `ambiguous` / `no_route`). |
| **`search_skills`** | Rank skills against a query using manifest metadata (names, descriptions, aliases, tags). |
| **`load_skill`** | Hydrate one skill — return its `SKILL.md` body and referenced files. |
| **`compose`** | Build a multi-step plan that chains skills (the Phase-4 DAG / pipeline feature). |

### Stage 1 — Search (metadata only)
Ranking happens against the **registry**, not the raw corpus. The router never has to read 1,800 `SKILL.md` files to answer "what's relevant?" — it scores against the compact manifest. This is what keeps routing fast (see [Performance](Performance-and-Benchmarks.md)).

### Stage 2 — Route / preflight (deterministic decision)
`preflight` makes a **deterministic decision with no secondary LLM call**. It returns one of:
- **`route`** — a single confident match; the host loads it.
- **`ambiguous`** — several candidates; the host AI picks.
- **`no_route`** — nothing relevant; continue normally.

The host AI provides the final sanity check, but the *decision* is computed, not delegated to another model. That determinism is the core of the "preflight-route" promise.

### Stage 3 — Load (hydration)
Only after a route is confirmed does `load_skill` read the chosen skill's `SKILL.md` (and any `references/`, `scripts/`, `resources/`) into context. Exactly one skill, exactly when needed.

---

## Decoupling: what was separated from what

The decoupling work (`ARCHITECTURE-decoupling.md`) made each layer resolvable independently:

- **Corpus location is configurable.** The router discovers the manifest and `skills/` through an environment-variable / config-driven resolution hierarchy rather than hard-coded paths — so the same binary works against a dev checkout, an installed stack, or a packaged corpus.
- **The registry is the single contract surface.** Nothing crosses between router and corpus except `manifest.json`. The generator in `scripts/registry/` emits both `manifest.json` and the slim `docs/build_manifest.json` **in lockstep**, and CI enforces they stay in sync with a `--check` mode.
- **Adapters are compact and uniform.** Every client receives the *same* policy contract (preflight → load one skill → don't copy the corpus), differing only in file location and format.

---

## CLI and MCP are two faces of one engine

Both the command-line surface and the [MCP server](MCP-Server.md) call the **same `skillservice` engine**. There is no parallel routing path — a request to `skill-router route "..."` and the equivalent MCP tool call run identical code. This is deliberate: it guarantees the router behaves the same no matter how an agent reaches it, and it's why the router-intelligence features (reranker, telemetry, compose) only need to be wired into one seam.

---

## Adapter model: from copy-propagation to one canonical engine

The **old** model physically copied skills into each downstream agent's directory — every client carried its own drifting copy of the corpus. The **new** model ships a single canonical engine accessed via the CLI or MCP; clients hold only a wrapper. See [Adapter deprecation](Roadmap-and-Phases.md#adapter-deprecation) for the migration path and what replaced the old adapters.

---

## Design principles (from `DESIGN_AND_MESSAGING.md`)

- **Positioning:** "Not a bundle to paste; a router-first skill system."
- **GitHub About line:** *"Router-first AI skill system … search, preflight-route, and load skills on demand."*
- **Visual identity:** modern, minimal, technical, calm, trustworthy — dark base with cyan / blue / green accents (the README badges and hero image follow this).
- **Public-safe by default:** machine state and secrets are git-ignored; a release audit script gates anything that could leak.

---

## Related pages

- [Skills Corpus](Skills-Corpus.md) — how a skill is structured and normalized
- [Manifest & Registry](Manifest-and-Registry.md) — the contract format and generator
- [Skill Router CLI](Skill-Router-CLI.md) — every command and flag
- [Node → Go Migration](Node-to-Go-Migration.md) — swapping the router with byte-parity
- [Roadmap & Phases](Roadmap-and-Phases.md) — where each capability sits
