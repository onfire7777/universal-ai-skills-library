# FAQ

Common questions about the Universal AI Skills Library.

---

### What is this, in one sentence?
A **router-first skill system**: one shared corpus of ~1,812 skills, one `skill-router` CLI, and compact per-client wrappers — so any AI agent can search, route, and load skills on demand without copying the corpus. See [Home](README.md).

### Why not just paste the skills into my agent?
That's the model UASL deliberately avoids. Copying the corpus into every client bloats context and causes drift between copies. UASL keeps the corpus in one place and loads **exactly one** relevant skill when a prompt needs it. See [Architecture](Architecture.md).

### How does the router decide which skill to use?
A deterministic [preflight](Architecture.md#stage-2--route--preflight-deterministic-decision) scores the prompt against the [manifest](Manifest-and-Registry.md) and returns `route`, `ambiguous`, or `no_route` — **with no extra LLM call**. The host AI makes the final sanity check. See [Skill Router CLI](Skill-Router-CLI.md#the-preflight-contract).

### Does routing call out to the network?
No. Preflight routing is **hermetic** — it makes its decision with no network calls, which is why it's fast and reproducible. See [Performance](Performance-and-Benchmarks.md).

### How many skills are there?
~1,812 per `manifest.json` `v2.2.8` (18 core + 1,794 library). The manifest is the authority for the live count. See [Skills Corpus](Skills-Corpus.md).

### Which AI agents are supported?
30+, across three layers (skill-root wrapper, repo-instruction, hosted/MCP) — Claude Code, Codex, Cursor, Hermes, Paperclip, OpenCode, Gemini CLI, Kiro, Continue, Windsurf, Cline, Aider, OpenHands, Copilot, and more. Run `skill-router sync matrix` for your machine. See [Agent Support Matrix](Agent-Support-Matrix.md).

### Do I have to copy 1,812 skills onto disk per agent?
No — the default is a **wrapper install** (a tiny skill that calls the router). Selected-skill and full-copy modes exist but are opt-in. See [Installation & Setup](Installation-and-Setup.md#install-modes).

### What language is the router written in?
Go (`go 1.25.0`). It's a single binary (`skill-router`) built with Cobra. It also speaks [MCP](MCP-Server.md). See [Skill Router CLI](Skill-Router-CLI.md).

### Why is there both Node and Go code in `scripts/registry/`?
The registry generator is mid-migration from Node to Go. Go is now authoritative; the Node generator remains as a **byte-parity oracle** until it's removed. See [Node → Go Migration](Node-to-Go-Migration.md).

### How do I add my own skill?
Create `skills/<id>/SKILL.md` with valid frontmatter, run `skill-router skills validate-manifest`, test with `skill-router preflight`, and open a PR. See [Contributing](Contributing.md#adding-a-skill).

### Is it safe to clone publicly? Where do my secrets go?
Yes — the repo is public-safe by design. Real secrets are written only to `%USERPROFILE%\.universal-ai-stack\secrets\.env` (machine-local, never committed), and a gitleaks CI gate blocks accidental secret commits. See [Security](Security.md).

### How do I run it as an MCP server?
`skill-router serve` (stdio, JSON-RPC 2.0). It exposes `route`, `search_skills`, `load_skill`, and `compose`. See [MCP Server](MCP-Server.md).

### What license is it under?
MIT. See [`LICENSE`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/LICENSE).

---

Didn't find your answer? Check the [Glossary](Glossary.md), the in-repo [`docs/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/docs), or open an issue.
