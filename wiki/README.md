# Universal AI Skills Library

**One shared skill corpus. One router. Every AI agent.**

The Universal AI Skills Library (UASL) is a **router-first skill system** that lets Codex, Claude, Cursor, Hermes, Paperclip, OpenCode, and many other AI agents **search, preflight-route, and load ~1,812 skills on demand** from a single shared corpus — without copying the skills into every client.

The full corpus lives once in [`skills/`](Skills-Corpus.md). Each AI client gets only a compact instruction file or a tiny wrapper that calls the `skill-router` CLI when — and only when — a real user prompt needs a matching skill. No corpus duplication, no context bloat.

> [!NOTE]
> This wiki documents the repository at its current state on `main`. The skill catalog is generated from [`manifest.json`](Manifest-and-Registry.md) (schema `v2.2.8`). Counts and numbers throughout cite that manifest and the in-repo docs.

---

## Wiki contents

**Getting started** — [Quick Start](Quick-Start.md) · [Installation & Setup](Installation-and-Setup.md) · [Agent Support Matrix](Agent-Support-Matrix.md)

**Concepts** — [Architecture](Architecture.md) · [Skills Corpus](Skills-Corpus.md) · [Sources & Integrations](Sources-and-Integrations.md)

**Reference** — [Skill Router CLI](Skill-Router-CLI.md) · [MCP Server](MCP-Server.md) · [Manifest & Registry](Manifest-and-Registry.md)

**Project** — [Roadmap & Phases](Roadmap-and-Phases.md) · [Node → Go Migration](Node-to-Go-Migration.md) · [Performance & Benchmarks](Performance-and-Benchmarks.md) · [Testing & CI](Testing-and-CI.md) · [Contributing](Contributing.md) · [Security](Security.md) · [Known Issues](Known-Issues.md)

**Help** — [FAQ](FAQ.md) · [Glossary](Glossary.md)

---

## Why router-first?

| Problem | How UASL solves it |
|---|---|
| Every AI client wants skills in a different directory. | Compact per-client adapters all point back to one canonical router + corpus. |
| Large skill libraries bloat the model's context window. | Deterministic **preflight** decides if any skill is needed, then loads *exactly one*. |
| Local AI stacks drift across Hermes, Paperclip, Codex, Claude, IDE tools. | Portable model/memory/routing/adapter config lives in repo-owned templates. |
| Public repos can leak machine state or secrets. | Public-safe defaults, ignored local state, validation scripts, release audit. |

The design intent, in the project's own words: this is **"not a bundle to paste; a router-first skill system."**

---

## How it works (30-second tour)

```
User prompt
   │
   ▼
skill-router preflight ──►  route?  ──► load exactly one skill  ──►  host AI uses it
   (deterministic,          ├─ ambiguous ──► host AI picks from candidates
    no extra LLM call)       └─ no_route  ──► continue normally
```

1. A user-prompt hook runs `skill-router preflight --hook-event UserPromptSubmit --json "<prompt>"`.
2. If the decision is **`route`**, the host loads one skill via `skill-router skill <name>`.
3. If **`ambiguous`**, the host AI chooses from the listed candidates.
4. If **`no_route`**, nothing is loaded and the agent proceeds normally.

See [Architecture](Architecture.md) for the full search → route → load pipeline.

---

## Quick links

- **New here?** → [Quick Start](Quick-Start.md)
- **Installing?** → [Installation & Setup](Installation-and-Setup.md) · [Agent Support Matrix](Agent-Support-Matrix.md)
- **Understanding the design?** → [Architecture](Architecture.md) · [Skills Corpus](Skills-Corpus.md)
- **Using the CLI?** → [Skill Router CLI](Skill-Router-CLI.md) · [MCP Server](MCP-Server.md)
- **Contributing?** → [Contributing](Contributing.md) · [Testing & CI](Testing-and-CI.md)
- **Where's the project headed?** → [Roadmap & Phases](Roadmap-and-Phases.md) · [Node → Go Migration](Node-to-Go-Migration.md)

---

## At a glance

| | |
|---|---|
| **Skills** | ~1,812 (18 core + 1,794 library) |
| **Router** | `skill-router` — a single Go CLI (also speaks MCP) |
| **Contract surface** | [`manifest.json`](Manifest-and-Registry.md) (`v2.2.8`) |
| **Architecture** | Router-first, three decoupled layers: Router · Corpus · Registry |
| **Supported clients** | 30+ AI agents & IDE tools — [see the matrix](Agent-Support-Matrix.md) |
| **Install** | `bash install.sh` (macOS/Linux/WSL) · `.\install.ps1` (Windows) |
| **License** | MIT |
| **Repo** | [onfire7777/universal-ai-skills-library](https://github.com/onfire7777/universal-ai-skills-library) |

---

## The three layers

UASL keeps three concerns strictly separated (see [Architecture](Architecture.md)):

1. **Router** — the `skill-router` Go CLI. Owns *no* skill content. Resolves prompts to skills through the manifest.
2. **Corpus** — [`skills/`](Skills-Corpus.md). Knows nothing about the router; it is just normalized skill directories.
3. **Registry** — [`manifest.json`](Manifest-and-Registry.md). The only contract between the two, regenerated deterministically by `scripts/registry/`.

This decoupling is what makes the same corpus usable from any client and what enables the in-progress [Node → Go migration](Node-to-Go-Migration.md) to swap the router implementation with byte-for-byte identical output.

---

*Looking for the source docs? The repo's [`docs/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/docs) folder is the upstream source for everything in this wiki.*
