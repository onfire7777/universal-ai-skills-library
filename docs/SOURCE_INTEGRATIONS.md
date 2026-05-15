# Source Integrations

This repo models a full local AI workspace as public-safe, reusable source
integrations. The goal is to keep the same capability shape as a mature personal
setup while making the install clean for other users.

## Design Goal

Every AI client should see the same capability layer:

- skills through `skill-router`
- durable memory through MemPalace
- structured memory/search through GBrain
- context-window protection through Context Mode
- browser retrieval through Lightpanda
- NotebookLM research through `nlm`
- X API workflows through `x-cli`
- Instagram one-turn and TUI workflows through `instagram-cli`
- LLM-ready web crawling and extraction through `crwl` / Crawl4AI
- hosted web search, scraping, crawling, interaction, and parsing through
  `firecrawl`
- fresh web search through the host AI when available
- external GSkills/GStack skills through read-only routing

The public repo provides the registry, policies, wrappers, adapter instructions,
and validation scripts. The user machine provides secrets, OAuth sessions,
downloaded tools, logs, local model files, and generated state.

## Public-Safe Contract

| Capability | Public repo owns | User machine owns | Default behavior |
| --- | --- | --- | --- |
| Lightpanda | wrapper paths, source policy, validation checks | Docker image/runtime, local wrapper install | on-demand; persistent bridge disabled |
| Context Mode | Codex hook sync policy and validation | npm package and local Codex config | CLI/hooks enabled; bridge disabled |
| MemPalace | shared memory wrapper instructions | memory database and optional MCP command | CLI baseline; bridge disabled |
| NotebookLM MCP CLI | source policy, canonical router skill, validation pointers | PyPI/uv tool install, Google auth session, generated notebook artifacts | CLI baseline; MCP disabled unless explicitly needed |
| x-cli | source policy, canonical router skill, validation pointers | Rust source checkout, built `x` executable, user-owned `.xrc` OAuth profile | CLI baseline; no background service |
| Instagram CLI | source policy, canonical router skill, validation pointers | Node source checkout, global `instagram-cli` link, user-owned `.instagram-cli` auth/config/log state | CLI baseline; no background service |
| Crawl4AI | source policy, canonical router skill, validation pointers | Python source checkout, dedicated venv, Playwright/Patchright browser assets, user-owned `.crawl4ai` cache/profiles/output | CLI baseline; no background service |
| Firecrawl | source policy, canonical router skill, validation pointers | shallow upstream checkout, global `firecrawl` CLI, user-owned Firecrawl login/API key, optional `firecrawl-mcp` | CLI baseline; no background service |
| Web search | host-native policy and fallback notes | host AI web/search tool or optional provider keys | no local search proxy, no committed key |
| GBrain | source/state pointers, embedding model policy | upstream checkout and local brain state | searchable mirror, not authoritative memory |
| GSkills/GStack | read-only external source policy | upstream checkout | namespaced on-demand skills |

## Registry

The portable source registry lives at:

```text
ai-setup/runtime/config/source-integrations.json
```

After install, it is materialized to:

```text
%USERPROFILE%\.universal-ai-stack\config\source-integrations.json
```

This file is intentionally a pointer and policy registry. It should not contain
real secrets, OAuth data, browser cookies, personal memory content, logs, or a
vendored copy of external repos.

## Install Behavior

Windows install:

```powershell
.\install.ps1
```

Optional local stack and startup:

```powershell
.\install.ps1 -InstallStartup -StartNow
```

The installer copies repo-owned templates into `%USERPROFILE%\.universal-ai-stack`,
expands path placeholders, creates local secrets if needed, and syncs compact AI
adapter instructions. It does not publish or commit the user's machine state.

## Non-Redundancy Rules

- Keep the full canonical skill corpus only in `skills/`.
- Keep AI-client roots to compact wrappers and instruction files.
- Keep external source repos external: GBrain stays in `%USERPROFILE%\gbrain`;
  GStack/GSkills stay in `%USERPROFILE%\.gstack\gstack`.
- Keep MemPalace data in `%USERPROFILE%\.mempalace\palace`.
- Keep NotebookLM auth, browser profiles, and generated artifacts in the user's
  local NotebookLM/Google account state; the repo only stores the router skill
  and source pointer.
- Keep x-cli auth in `%USERPROFILE%\.xrc` or an explicitly selected profile
  file; the repo stores only the router skill and source pointer.
- Keep Instagram CLI auth, config, logs, downloaded media, and private messages
  in `%USERPROFILE%\.instagram-cli`; the repo stores only the router skill and
  source pointer.
- Keep Crawl4AI runtime state, venv, browser profiles, cookies, crawl cache,
  screenshots, and extracted private content in `%USERPROFILE%\.crawl4ai`;
  the repo stores only the router skill and source pointer.
- Keep Firecrawl authentication, API keys, account state, generated private
  scrape/crawl output, screenshots, cookies, and browser/session data local;
  the repo stores only the router skill and source pointer. Do not run
  `firecrawl-cli init --all` in this stack because it duplicates skills into
  individual AI roots.
- Keep Context Mode and Lightpanda as wrappers or host tools, not vendored repo
  copies.
- Keep web search host-native by default. Do not commit web-search provider keys
  or run a default local search proxy.

## Adapter Behavior

Each AI client receives compact instructions that say:

1. Run `skill-router preflight --json "<latest user prompt>"` only for real user
   prompts.
2. Load exactly one clearly relevant skill with `skill-router skill <name>`.
3. Use `Search-UniversalAIMemory.ps1` for prior durable context.
4. Use `Save-UniversalAIMemory.ps1` only for explicit durable saves or confirmed
   project facts.
5. Treat Context Mode as scratch/context protection, not memory.
6. Use Lightpanda for controlled page fetch/extraction and CDP workflows.
7. Use `skill-router skill notebooklm-mcp-cli` before NotebookLM notebook,
   source, query, Studio artifact, or NotebookLM MCP work.
8. Use `skill-router skill x-cli` before X API account, post, search, stream,
   list, direct message, or social graph work.
9. Use `skill-router skill instagram-cli` before Instagram inbox, direct
   message, read, reply, unsend, media download, TUI, or config work.
10. Use `skill-router skill crawl4ai` before Crawl4AI install/update/setup,
    `crwl`, LLM-ready web crawling, Markdown/JSON extraction, bounded deep
    crawl, profile, or CDP workflows.
11. Use `skill-router skill firecrawl` before Firecrawl install/update/login,
    hosted Firecrawl search/scrape/crawl/map/parse/interact/agent workflows,
    SDK/API examples, or optional Firecrawl MCP configuration.
12. Use host-native web search for fresh search when available.
13. Load GBrain and GSkills/GStack skills through `skill-router`, not by copying
    their source trees into the client.

## Validation

Repo validation:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1
skill-router skills validate-manifest
```

Installed-stack validation:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIStack.ps1"
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIAdapters.ps1"
powershell -NoProfile -ExecutionPolicy Bypass -File "$env:USERPROFILE\.universal-ai-stack\scripts\Test-UniversalAIContextTools.ps1"
```

Expected clean state:

- `source-integrations.json` includes Lightpanda, Context Mode, MemPalace,
  NotebookLM MCP CLI, x-cli, Instagram CLI, Crawl4AI, Firecrawl, web search,
  GBrain, and GSkills/GStack.
- Every supported AI adapter has compact universal skill, memory, context, and
  source-integration instructions.
- Persistent MCP bridge ports may be down in the low-resource profile.
- No external source tree is copied into every AI root.
- No secrets or private local state are present in the repo.
