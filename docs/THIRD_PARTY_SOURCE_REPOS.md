# Third-Party Source Repos

The universal stack can expose large third-party skill/tool projects without
copying them into every AI client. Keep one upstream checkout, index its skills
read-only, and install only the compact `universal-ai-skills` wrapper into local
agent roots.

## GStack

- Upstream: <https://github.com/garrytan/gstack>
- Local source: `%USERPROFILE%\.gstack\gstack`
- Runtime artifacts:
  - `%USERPROFILE%\.gstack\gstack\browse\dist\browse.exe`
  - `%USERPROFILE%\.gstack\gstack\design\dist\design.exe`
  - `%USERPROFILE%\.gstack\gstack\make-pdf\dist\pdf.exe`
- Router sources:
  - `gstack-gbrain`: `~\.gstack\gstack\.gbrain\skills`
  - `gstack-codex`: `~\.gstack\gstack\.agents\skills`
  - `gstack-openclaw`: `~\.gstack\gstack\openclaw\skills`

Use namespaced skill names such as `gstack-review`, `gstack-qa`,
`gstack-ship`, and `gstack-cso` to avoid collisions with generic skills.

## GBrain

- Upstream: <https://github.com/garrytan/gbrain>
- Local source: `%USERPROFILE%\gbrain`
- Runtime state: `%USERPROFILE%\.gbrain`
- Canonical CLI: `gbrain` from `%USERPROFILE%\.bun\bin`
- Windows compatibility shims:
  - `%USERPROFILE%\go\bin\bun.cmd`
  - `%USERPROFILE%\go\bin\gbrain.cmd`
- Router sources:
  - `gbrain-source`: `~\gbrain\skills`
  - `gbrain-user`: `~\.gbrain\skills`

GBrain is a persistent local knowledge system. It should be called through its
CLI for brain state operations and through `skill-router skill <name>` for
skillpack instructions.

The Universal AI Stack uses GBrain as the structured searchable mirror for
explicit shared-memory notes. MemPalace remains the authoritative durable memory
store. GBrain vectors use the local `qwen3-embedding-0.6b` llama.cpp service at
`http://127.0.0.1:18084/v1` with 1024 dimensions, so text embedding does not
require an OpenAI API key.

## Verification

```powershell
skill-router skills sources --refresh
skill-router skill search gstack
skill-router skill search gbrain
skill-router gstack status
skill-router gbrain status
gbrain --version
```

Warnings for missing optional API keys are acceptable for provider-backed query
expansion, subagent workers, or other paid optional features. They are not
expected for the local GBrain text embedding path when the Qwen embedding
service is healthy.

The `bun.cmd` and `gbrain.cmd` shims are thin delegates, not duplicate installs.
They exist because long-running AI desktop clients can inherit an older PATH
before `%USERPROFILE%\.bun\bin` was added. The canonical binaries stay in
`.bun\bin`; the shims only make the same tools visible from already-running
agent shells whose PATH already includes `%USERPROFILE%\go\bin`.
