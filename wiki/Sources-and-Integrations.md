# Sources & Integrations

UASL integrates external skill sources and optional tools **by reference, not by vendoring**. The canonical corpus stays in [`skills/`](Skills-Corpus.md); external roots are indexed read-only, and optional services are reached through router pointers that keep their auth/state local.

Source docs: [`docs/SOURCE_INTEGRATIONS.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/SOURCE_INTEGRATIONS.md) · [`docs/THIRD_PARTY_SOURCE_REPOS.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/THIRD_PARTY_SOURCE_REPOS.md)

---

## Non-redundancy rule

> The canonical corpus is single-sourced in `skills/`. External sources are **indexed read-only**; nothing is copied into the repo.

This means the router can *see* skills that live in other tools' directories without the repo carrying a second copy that could drift. Secrets and runtime state always stay in the user's own home directory, never in the repo.

---

## Third-party skill sources (indexed read-only)

| Source | Local path (typical) | Notes |
|---|---|---|
| **GStack** | `~/.gstack/gstack` | Engineering skill stack, maintained separately. Router sub-sources: `gstack-gbrain`, `gstack-codex`, `gstack-openclaw`. |
| **GBrain** | `~/gbrain` | Personal knowledge brain; indexed read-only as a memory/source pointer. |

GStack's router sources map to:
- `gstack-gbrain` → `~/.gstack/gstack/.gbrain/skills`
- `gstack-codex` → `~/.gstack/gstack/.agents/skills`
- `gstack-openclaw` → `~/.gstack/gstack/openclaw/skills`

---

## Optional services & tool wrappers (router pointers only)

These are reached through a router skill/pointer; their credentials and artifacts stay in the user's account or home directory.

| Integration | What it adds | Local state |
|---|---|---|
| **MemPalace** | Durable memory store | `~/.mempalace/palace` |
| **NotebookLM** | Notebook research source | Auth/artifacts in the user's Google account |
| **x-cli** | X (Twitter) API workflows | `~/.xrc` |
| **Instagram CLI** | Instagram workflows | `~/.instagram-cli` |
| **Crawl4AI** | Local, LLM-ready crawling | `~/.crawl4ai` |
| **Firecrawl** | Hosted web-data API/CLI | Keys kept local |
| **Context Mode** | Context-window protection (not durable memory) | Host tool/wrapper |
| **Lightpanda** | On-demand browser automation | Host tool/wrapper |
| **OpenRouter** | Optional OpenAI-compatible model fallback | `OPENROUTER_API_KEY` env var |

> [!NOTE]
> Web search is **host-native by default** — no provider keys live in the repo. Lightpanda is used only for controlled, on-demand retrieval.

---

## The local runtime (Universal AI Stack)

On a full install, `ai-setup/` materializes a local runtime that routes model calls across providers with health checks and a guarded local fallback:

- **Router endpoint:** `http://127.0.0.1:18100/v1`, model `auto-coding`.
- **Default model failover order:** GPT → Kimi (Moonshot, HTTP fallback) → Claude Opus → local Qwen (llama.cpp).
- **Config (repo-owned source of truth):** `ai-setup/runtime/config/model-registry.json`, `routing-policy.json` (240s global timeout, ~10-min circuit breaker), `source-integrations.json`.

See [Installation & Setup](Installation-and-Setup.md#the-packaging-directories) for how the stack is laid out and which files are repo-owned vs machine-owned.

---

## How a source becomes routable

1. A pointer/config entry registers the external root in `source-integrations.json`.
2. The router indexes that root read-only alongside the canonical corpus.
3. Prompts route to skills from any indexed root — without duplicating them into the repo.

This is the same decoupling principle as the rest of the [architecture](Architecture.md): the router resolves *locations*; it never owns the content.
