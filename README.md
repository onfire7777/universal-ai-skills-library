# Universal AI Skills Library

CLI-first skill routing for Codex, Claude Code, Cursor, Gemini, OpenCode, Manus compatibility workflows, and local desktop AI setups.

The repository keeps a large skill corpus available without loading it all into every agent context. The source of truth is `skills/`; the router is `skill-router-cli/`; plugin metadata lives in `plugin/`.

## Layout

```text
universal-ai-skills-library/
├── README.md
├── manifest.json
├── skill-router-cli/       # Go CLI source for the universal router
├── skills/                 # Source-of-truth skill corpus: 1,804 canonical skills
├── plugin/                 # Universal plugin metadata plus Codex/Claude adapters
├── infrastructure/         # Optional persistent MCP bridge scripts
└── docs/                   # Architecture, compatibility, and migration notes
```

## Primary Command

```bash
skill-router skill <name>
skill-router skill search <query>
skill-router skill list
```

`manus` remains a legacy compatibility executable for existing local rules and scripts. New docs and local setup should prefer `skill-router`.

## Install Or Build

```bash
git clone https://github.com/onfire7777/universal-ai-skills-library.git
cd universal-ai-skills-library/skill-router-cli
go build -o "$HOME/go/bin/skill-router" .
```

On Windows the installed binaries are:

```text
%USERPROFILE%\go\bin\skill-router.exe
%USERPROFILE%\go\bin\manus.exe       # compatibility alias
```

## Common Commands

```bash
skill-router skill <name>          # Print one SKILL.md
skill-router skill search <query>  # Search manifest names and descriptions
skill-router skill list            # List all skills from manifest.json
skill-router doctor                # Check local AI stack health
skill-router mcp status            # Check optional MCP bridge endpoints
skill-router audit <path>          # Run audit workflows
skill-router oracle <question>     # Query multi-model oracle flow
skill-router print <api>           # Generate API-specific CLI scaffolds
```

## Core Skill Groups

Core universal workflow skills include:

- `universal-ai-config` - audit and repair the cross-AI setup
- `skill-creator` - create portable AI skills
- `skill-sync` - sync skills into local or global roots
- `skill-debugger` - debug one skill with configured model roles
- `prompt-engineer` - optimize prompts for agentic workflows
- `context-anchor` - keep long sessions focused
- `chat-summarizer` - create AI-to-AI handoff summaries
- `persistent-computing` - choose local, managed, or third-party persistent compute

Compatibility-specific skills keep their platform names, for example `manus-api` because it targets the real Manus API.

## Architecture Rule

Use neutral names for universal capabilities and platform names only for adapters:

- Universal core: `skill-router`, `skills/`, `universal-ai-config`, `plugin/`.
- Compatibility adapters: Manus API, `.manus` roots, legacy `manus.exe`, existing Windows scheduled task names.
- Optional persistent endpoints: MemPalace, Context Mode, Skill Seekers, Lightpanda.

Detailed policy: `docs/UNIVERSAL_COMPATIBILITY.md`.

## MCP Policy

The router should be the default path. MCP bridges are optional adapters for services that need persistent endpoints or shared state:

- MemPalace: durable memory
- Context Mode: large-output routing and session continuity
- Skill Seekers: skill creation and packaging workflows
- Lightpanda: browser/CDP automation when Docker is available

If a workflow can run through `skill-router`, use the CLI. If it needs a long-running endpoint, shared memory service, credential bridge, or browser/CDP runtime, use MCP.

## Agent Rule Snippet

```text
Use Universal AI Skills Router for skills:
- skill-router skill <name>
- skill-router skill search <query>
- skill-router doctor

Keep always-loaded instructions compact. Do not paste the full skills corpus into global rules. MCP bridges are optional adapters for persistent endpoint workflows.
```

## Verification

```bash
skill-router --version
skill-router skill universal-ai-config
skill-router skill search debugging
skill-router doctor
```

Expected current corpus size: 1,804 canonical skills.
