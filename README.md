# Manus Skills Library v2.0

CLI-first universal skill routing for Codex, Claude Code, Manus, Cursor, Gemini, OpenCode, and local desktop AI workflows.

This repository keeps the complete skill corpus available without loading it all into agent context. The source of truth is `skills/`; the router is `manus-cli/`; the universal plugin metadata is `plugin/`.

## Architecture

```text
manus-skills-library/
├── README.md
├── manifest.json
├── manus-cli/              # Go CLI: one entry point for skills and tools
├── infrastructure/         # Optional persistent MCP bridge setup
│   └── mcp-bridges/
├── skills/                 # All 785 skills, including custom core skills
├── plugin/                 # Universal plugin metadata and AI instructions
└── docs/                   # Supporting docs and inventory notes
```

Primary rule:

```bash
manus skill <name>
```

Agents should load only the skill needed for the current task. Do not paste the full skills list or full skill bodies into always-loaded instruction files.

## Quick Start

```bash
git clone https://github.com/onfire7777/manus-skills-library.git
cd manus-skills-library/manus-cli
go install .
manus --version
manus skill search debugging
manus skill skill-debugger
```

Useful commands:

```bash
manus skill <name>          # Print one SKILL.md
manus skill search <query>  # Search manifest names and descriptions
manus skill list            # List all skills from manifest.json
manus doctor                # Check local AI stack health
manus mcp status            # Check optional MCP bridge endpoints
manus audit <path>          # Run audit workflows
manus oracle <question>     # Query multi-model oracle flow
manus print <api>           # Generate API-specific CLI scaffolds
```

## Skills

All skills live under `skills/`.

The previous split between top-level custom skills and library skills is removed. The 14 custom core skills now sit beside the 770 library skills:

- `skills/chat-summarizer`
- `skills/context-anchor`
- `skills/file-organizer`
- `skills/manus-api`
- `skills/model-selector`
- `skills/multi-model-code-auditor`
- `skills/multi-model-oracle`
- `skills/music-prompter`
- `skills/persistent-computing`
- `skills/prompt-engineer`
- `skills/skill-creator`
- `skills/skill-debugger`
- `skills/skill-sync`
- `skills/ultimate-skill-creator`

`manifest.json` is the machine-readable catalog. Each skill entry points to its `skills/<name>` directory.

## Plugin Model

`plugin/plugin.json` describes one universal plugin surface:

- no full skill bundle injected into context
- no duplicate per-client skill copies required for normal operation
- CLI-first access through `manus skill <name>`
- optional client-specific instruction payloads in `plugin/codex/` and `plugin/claude/`

The plugin exists to route agents to the CLI and preserve context. It is not a second copy of the skill corpus.

## MCP Policy

MCP servers are optional. Prefer CLI calls when the task is local, deterministic, or skill-loading related.

Keep MCP bridges only when a persistent endpoint is genuinely useful:

- Skill Seekers: dynamic skill generation and packaging
- MemPalace: durable cross-session memory
- Context Mode: context-window protection and indexed long-output routing
- Lightpanda: browser automation when a persistent CDP endpoint is useful

Current local bridge convention:

```text
Skill Seekers  http://127.0.0.1:8875/mcp
MemPalace      http://127.0.0.1:8876/mcp
Context Mode   http://127.0.0.1:8877/mcp
Lightpanda     http://127.0.0.1:8878/mcp
```

If a workflow can be done through `manus`, use `manus`. If it needs a long-running server, credentials bridge, browser endpoint, or shared memory service, use MCP.

## Local AI Setup

Recommended always-loaded instruction:

```text
Use Manus CLI for skills:
- manus skill <name>
- manus skill search <query>
- manus doctor
Do not load the full skill corpus into context.
```

Client targets:

- Codex: global `AGENTS.md` plus optional plugin cache entry
- Claude Code: global `CLAUDE.md` plus optional plugin entry
- Cursor: rules file pointing to the CLI
- Gemini/OpenCode/Manus: compact instruction rule plus CLI availability

The local setup should be an index, not the library itself.

## Install Script

For environments that still need physical skill directories:

```bash
bash install.sh --target /path/to/skills
```

That copies `skills/*` only. It no longer copies separate top-level custom skill directories.

## Development

Build and test the CLI:

```bash
cd manus-cli
go test ./...
go build .
```

Verify repository structure:

```bash
manus skill list
manus skill persistent-computing
manus mcp status
git status --short
```

## Design Principles

- One source of truth: `skills/`
- One router: `manus`
- One plugin concept: CLI-first universal plugin
- Context stays small: load skills on demand
- MCP stays purposeful: persistent endpoints only
- Local AI configs stay compact and point to the router
