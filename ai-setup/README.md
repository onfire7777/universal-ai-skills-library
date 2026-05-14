# Universal AI Setup

This folder makes the optional local Universal AI Stack reproducible from the
repository. It is the runtime side of the router-first system: the skill corpus
stays centralized, while model routing, Hermes/Paperclip integration, memory,
embeddings, Context Mode, Lightpanda, and health checks are installed from
portable templates.

It is intentionally split into:

- `manifests/`: portable inventory and install policy.
- `runtime/`: router, supervisor, config templates, and client sync scripts.
- `scripts/install-universal-ai-stack.ps1`: materializes the runtime into `%USERPROFILE%\.universal-ai-stack`.
- `scripts/validate-universal-ai-stack.ps1`: validates model-specific config, portability, and secret hygiene.

The repo owns templates and code. The machine owns secrets, logs, generated state, OAuth sessions, and downloaded model files.

For the exact config map that connects each AI client, router, model endpoint,
memory system, embedding service, Context Mode, Lightpanda, and NotebookLM, see
`../docs/UNIVERSAL_AI_CONNECTION_CONFIGS.md`.

The source-integration layer is also repo-owned and public-safe. It records
Lightpanda, Context Mode, MemPalace, NotebookLM MCP CLI, host-native web search,
GBrain, and GSkills/GStack as shared integrations in
`runtime/config/source-integrations.json`.
Those entries are pointers, policies, and wrapper paths; they do not vendor
private machine state, upstream source checkouts, browser sessions, or full
external skill trees into the public repo.

## Install

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\install-universal-ai-stack.ps1 -InstallStartup -StartNow
```

Optional Kimi key:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\install-universal-ai-stack.ps1 -KimiApiKey "<key>"
```

## Validate

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1 -CheckInstalled
```

## Model Policy

Primary host-session model:

1. `gpt-5.5`, `xhigh`, fast tier, through official CLI/session auth.

HTTP/API fallback order:

1. `kimi-k2.6-thinking`
2. `qwen3-coder-30b-a3b-q4`

Only configured local generative model:

1. `qwen3-coder-30b-a3b-q4`

Local shared-memory embedding model:

1. `qwen3-embedding-0.6b-q8`

Local coding default is Qwen3-Coder-30B-A3B-Instruct `Q4_K_M` at 16k context with `--n-gpu-layers 99`, flash attention on, batch `384`, ubatch `192`, 6 threads, parallel `1`, and q4 KV cache. The local proxy refuses to start the 30B backend below `20GB` free VRAM or `6GB` free RAM, rejects local request bodies over `8MB`, runs `llama-server` below-normal priority, and caps local request waits at `300` seconds. The local embedding default is Qwen3-Embedding-0.6B `Q8_0` at `http://127.0.0.1:18084/v1` with 1024 dimensions for GBrain, guarded at `1GB` free VRAM, `2GB` free RAM, and `2MB` request bodies. Heavier local model records are intentionally not registered in the default stack so Hermes/Paperclip cannot route to stale or missing endpoints or start redundant model servers.
