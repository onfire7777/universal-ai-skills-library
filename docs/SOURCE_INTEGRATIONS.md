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
7. Use host-native web search for fresh search when available.
8. Load GBrain and GSkills/GStack skills through `skill-router`, not by copying
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

- `source-integrations.json` includes Lightpanda, Context Mode, MemPalace, web
  search, GBrain, and GSkills/GStack.
- Every supported AI adapter has compact universal skill, memory, context, and
  source-integration instructions.
- Persistent MCP bridge ports may be down in the low-resource profile.
- No external source tree is copied into every AI root.
- No secrets or private local state are present in the repo.
