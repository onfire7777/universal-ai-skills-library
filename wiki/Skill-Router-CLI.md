# Skill Router CLI

`skill-router` is the single Go binary at the heart of UASL. It scores prompts against the [registry](Manifest-and-Registry.md), makes the deterministic [preflight](Architecture.md#stage-2--route--preflight-deterministic-decision) decision, loads skills, regenerates the manifest, and can run as an [MCP server](MCP-Server.md). Every client talks to the corpus *through* this binary.

Source: [`skill-router-cli/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/skill-router-cli)

---

## Identity

| | |
|---|---|
| **Binary** | `skill-router` |
| **Module** | `github.com/onfire7777/universal-ai-skills-library/skill-router-cli` |
| **Language** | Go (`go 1.25.0`, toolchain `1.26.x`) |
| **Version** | `2.2.8` (matches the manifest) |
| **Dependencies** | [`spf13/cobra`](https://github.com/spf13/cobra) (command framework) · [`fatih/color`](https://github.com/fatih/color) (terminal output) — intentionally minimal |
| **Build** | `make build` / `go build ./...` |

---

## Core routing commands

These are the commands you'll use day to day.

| Command | Purpose |
|---|---|
| `skill-router skill search <query>` | Rank skills against a free-text query. |
| `skill-router skill <name>` | Load (hydrate) one skill — prints its `SKILL.md`. (`skill` is an alias of `skills`.) |
| `skill-router skills list` | List all skills in the corpus. |
| `skill-router skills validate-manifest` | Validate that the manifest matches the corpus. |
| `skill-router route "<prompt>"` | Decide the single best skill for a prompt. |
| `skill-router route --explain "<prompt>"` | Same, but show the scoring rationale. |
| `skill-router preflight --json "<prompt>"` | Deterministic route / ambiguous / no_route decision as JSON. |
| `skill-router preflight --hook-event UserPromptSubmit --json "<prompt>"` | Preflight wired for a user-prompt hook. |
| `skill-router compose "<goal>"` | Build a working set / multi-step plan of skills. |
| `skill-router compose --pipeline "<goal>"` | Emit an ordered multi-step **DAG** pipeline (Phase 4). |

### The preflight contract

`preflight` returns one of three decisions, and that decision is computed — **there is no secondary LLM call**:

| Decision | Meaning | Host action |
|---|---|---|
| `route` | One confident match | Load it with `skill-router skill <name>` |
| `ambiguous` | Several candidates | Host AI picks from the list |
| `no_route` | Nothing relevant | Continue normally |

This is the mechanism behind UASL's "load exactly one relevant skill when needed" promise. See [Architecture](Architecture.md#the-router-first-pipeline-search--route--load).

---

## Registry commands

The manifest is generated, never hand-edited. `registry build` is the generator (see [Manifest & Registry](Manifest-and-Registry.md)).

| Command | Purpose |
|---|---|
| `skill-router registry build --write` | Regenerate `manifest.json` + `docs/build_manifest.json`. |
| `skill-router registry build --check` | Verify the committed artifacts are in sync (CI gate). |
| `skill-router registry build --print` | Print the generated manifest without writing. |
| `skill-router registry build --only <ids>` | Rebuild a subset. |

Notable flags include `--faithful`, `--optimize`, and `--repo` for pointing at an alternate checkout.

---

## MCP server mode

```bash
skill-router serve
```

Runs `skill-router` as a **Model Context Protocol** server over stdio (JSON-RPC 2.0). It exposes four tools — `route`, `search_skills`, `load_skill`, `compose` — backed by the same engine as the CLI. Full details and client wiring on the [MCP Server](MCP-Server.md) page.

On Windows there is an additional `skill-router mcp` command group (`status`, `start`, `stop`, `restart`, `logs`, `watchdog`, `setup`) that manages the bridge as a managed process.

---

## Setup, sync & system commands

| Command | Purpose |
|---|---|
| `skill-router sync matrix` | Show the [agent support matrix](Agent-Support-Matrix.md) (read-only, safe). |
| `skill-router sync installed` | Install/update the compact wrapper into every detected client root. |
| `skill-router sync codex` / `claude` / `paperclip` | Sync the wrapper for a specific client. |
| `skill-router doctor` / `doctor --json` | Diagnose the installation and environment. |
| `skill-router config` | Inspect/modify router configuration. |
| `skill-router update` | Update the router. |
| `skill-router --version` | Print the version. |

> [!TIP]
> `sync matrix` and `doctor` are the safest first commands after install — neither mutates client roots. See [Installation & Setup](Installation-and-Setup.md).

There are additional research and productivity subcommand groups (e.g. `find`, `web`, `hf`, `gstack`, `gbrain`, `files`, `schedule`, `oracle`, `models`, `audit`) that wrap optional [source integrations](Sources-and-Integrations.md). Run `skill-router --help` for the complete, version-accurate list.

---

## Internal package layout

The CLI is ~50 Go files across ten internal packages. The key one is `skillservice` — the single engine both the CLI and MCP server call.

| Package | Purpose |
|---|---|
| `skillservice/` | Core engine: routing, scoring, semantic embedding, reranking, composition |
| `registry/` | Manifest builder, frontmatter parser, artifact generation |
| `reranker/` | Feedback-driven re-ranking model |
| `eval/` | Routing-quality metrics (P@1, MRR, Recall@5) |
| `telemetry/` | Decision records, analytics, version tracking |
| `skillsync/` | Wrapper propagation into client roots |
| `platform/` | Path resolution, environment detection, asset discovery |
| `runner/` | Subprocess execution (Python, PowerShell) |
| `mcpcli/` | MCP bridge management helpers |

See [Architecture](Architecture.md#cli-and-mcp-are-two-faces-of-one-engine) for why everything funnels through one engine.

---

## How the CLI finds the corpus

The router resolves the manifest and skills directory through configuration and environment variables rather than hard-coded paths — this is the [decoupling](Architecture.md#decoupling-what-was-separated-from-what) that lets one binary serve a dev checkout, an installed stack, or a packaged corpus.

| Variable | Meaning |
|---|---|
| `SKILL_ROUTER_SKILLS` | Path to the canonical `skills/` directory. |
| `SKILL_ROUTER_SEMANTIC` | Enable semantic routing. |
| `SKILL_ROUTER_VECTORS` | Path to the int8 vector store (semantic mode). |
| `SKILL_ROUTER_HOOK_EVENT` | Hook event name passed through to routing. |
| `MCP_DIR` | Windows MCP bridge config directory. |

If unset, the router falls back to the repo-local `skills/` and `manifest.json`.

---

## Router intelligence

Beyond lexical matching, the engine includes:

- **Semantic routing** — optional offline int8 embeddings (Ollama-backed or built-in), enabled via `SKILL_ROUTER_SEMANTIC`.
- **Reranker** — a feedback-trained model that improves ranking over time; it plugs into the single engine seam rather than forking a second routing path.
- **Telemetry & eval** — decision records feed the routing-eval harness (P@1 / MRR / Recall@5). See [Testing & CI](Testing-and-CI.md).
- **Compose / DAG** — `compose --pipeline` decomposes a goal into an ordered capability pipeline (the Phase-4 multi-step feature).

---

## Build, install & test

```bash
make build           # go build ./...
make vet             # go vet ./...
make test            # go test ./...
make install-local   # install to ~/.local/bin/skill-router
make baseline        # record a routing baseline
make eval            # compute P@1 / MRR / Recall@5
make eval-check      # regression gate — fail if metrics decline
make registry-check  # verify manifest matches committed artifacts
make parity          # prove byte-identical output vs the legacy Node generator
```

`make parity` is the safety net for the [Node → Go migration](Node-to-Go-Migration.md): it proves the Go generator emits byte-for-byte identical artifacts.

---

## Related pages

- [MCP Server](MCP-Server.md) — tool schemas and client registration
- [Manifest & Registry](Manifest-and-Registry.md) — the contract the CLI reads/writes
- [Installation & Setup](Installation-and-Setup.md) — getting the binary on your PATH
- [Testing & CI](Testing-and-CI.md) — eval harness and parity gate
