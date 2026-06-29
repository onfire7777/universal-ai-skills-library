# Glossary

Key terms used throughout the Universal AI Skills Library and this wiki.

| Term | Definition |
|---|---|
| **UASL** | Universal AI Skills Library — this project. |
| **Router-first** | The architecture principle that clients reach skills *through* a router, never by holding their own copy. See [Architecture](Architecture.md). |
| **`skill-router`** | The single Go CLI binary that routes prompts to skills, loads them, generates the registry, and serves MCP. See [Skill Router CLI](Skill-Router-CLI.md). |
| **Corpus** | The set of all skills, stored once under `skills/`. See [Skills Corpus](Skills-Corpus.md). |
| **Skill** | A directory `skills/<id>/` containing a `SKILL.md` (+ optional files) that teaches an agent how to do something. |
| **`SKILL.md`** | The required file in every skill: YAML frontmatter + a markdown body. |
| **Registry / Manifest** | `manifest.json` — the generated, machine-readable index of the corpus and the only contract between corpus and router. See [Manifest & Registry](Manifest-and-Registry.md). |
| **Canonical ID** | A skill's kebab-case top-level directory name; the authoritative identifier. |
| **Alias** | A legacy or display name that resolves to a canonical ID via `compatibility_aliases.json`. |
| **Preflight** | The deterministic decision step that returns `route` / `ambiguous` / `no_route` for a prompt, with no extra LLM call. |
| **Route** | A preflight decision that one skill is a confident match (also the name of the `route` verb/command). |
| **Ambiguous** | A preflight decision with several candidate skills; the host AI chooses. |
| **No-route** | A preflight decision that no skill is relevant; the agent proceeds normally. |
| **Compose / Pipeline** | Building a multi-step plan (DAG) that chains skills — `skill-router compose --pipeline`. |
| **`skillservice`** | The internal Go engine both the CLI and the MCP server call — the single routing seam. |
| **MCP** | Model Context Protocol — the tool-server interface exposed by `skill-router serve`. See [MCP Server](MCP-Server.md). |
| **Wrapper / Adapter** | The compact, per-client instruction or skill that calls the router instead of embedding the corpus. |
| **Wrapper install** | The default install mode: a tiny wrapper, not the whole corpus. See [Installation & Setup](Installation-and-Setup.md#install-modes). |
| **Universal AI Stack** | The optional local runtime (`ai-setup/`) that routes model calls across providers with health checks and a guarded local fallback. |
| **Byte-parity** | The guarantee that the Go registry generator emits output byte-for-byte identical to the legacy Node generator (`make parity`). See [Node → Go Migration](Node-to-Go-Migration.md). |
| **Routing-eval** | The harness (`tests/routing-eval/`) that measures routing quality with P@1, MRR, Recall@5, and abstention accuracy. See [Testing & CI](Testing-and-CI.md#routing-eval-harness). |
| **P@1** | Precision at rank 1 — how often the top-ranked skill is correct. |
| **MRR** | Mean Reciprocal Rank — average of 1/(rank of the first correct skill). |
| **Recall@5** | Whether a correct skill appears within the top 5 results. |
| **Abstention accuracy** | How reliably the router *refuses* to route on off-domain prompts. |
| **Characterization test** | A test that pins existing behavior so refactors don't change it (`tests/characterization/`). |
| **Public-safe** | Defaults that keep machine state and secrets out of the repo, enforced by audits and gitleaks. |

---

See also: [FAQ](FAQ.md) · [Architecture](Architecture.md) · [Home](README.md).
