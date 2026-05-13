# Universal AI Setup

This folder makes the local Universal AI Stack reproducible from the repository.

It is intentionally split into:

- `manifests/`: portable inventory and install policy.
- `runtime/`: router, supervisor, config templates, and client sync scripts.
- `scripts/install-universal-ai-stack.ps1`: materializes the runtime into `%USERPROFILE%\.universal-ai-stack`.
- `scripts/validate-universal-ai-stack.ps1`: validates model-specific config, portability, and secret hygiene.

The repo owns templates and code. The machine owns secrets, logs, generated state, OAuth sessions, and downloaded model files.

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

Local coding default is Qwen3-Coder-30B-A3B-Instruct `Q4_K_M` at 16k context. The local embedding default is Qwen3-Embedding-0.6B `Q8_0` at `http://127.0.0.1:18084/v1` with 1024 dimensions for GBrain. Heavier local model records are intentionally not registered in the default stack so Hermes/Paperclip cannot route to stale or missing endpoints or start redundant model servers.
