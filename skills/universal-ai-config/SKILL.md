---
name: universal-ai-config
description: Configure and audit the local universal AI skills setup across Codex, ChatGPT, Claude, Claude Code, Claude Cowork, Cursor, Gemini, OpenCode, OpenSkills, Hermes Agent, OpenClaw, Cline, Continue, GitHub Copilot, Kiro, OpenHands, hosted provider compatibility surfaces, plugin manifests, optional MCP bridges, and global instruction files. Use when the user asks to organize, consolidate, optimize, rename, or repair the cross-AI skill router configuration.
---

# Universal AI Config

Use this skill to keep the user's AI stack clean, router-first, universal, and context-efficient.

## Scope

- Universal AI Skills Router source and binaries
- `skills/` source-of-truth layout
- universal plugin manifests
- Codex, ChatGPT, Claude, Claude Code, Claude Cowork, Cursor, Gemini, OpenCode, OpenSkills, Hermes Agent, OpenClaw, Cline, Continue, GitHub Copilot, Kiro, OpenHands, and hosted provider compatibility instruction surfaces
- optional MCP bridge policy
- local persistent services and scheduled tasks

## Operating Model

Prefer this order:

1. CLI direct calls for deterministic local work.
2. MCP only for persistent endpoint workflows.
3. Full skill-body loading only through `skill-router skill <name>` when a task needs that skill.

Do not put the full skills corpus into always-loaded files.

Compatibility adapter types:

- `skill-root`: `SKILL.md` roots such as OpenSkills, Claude Code, Codex, OpenCode, Cline, OpenHands, Hermes Agent, and OpenClaw.
- `repo-instruction`: compact guidance files such as `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.continue/rules`, `.kiro/steering`, `.junie/guidelines.md`, and `CONVENTIONS.md`.
- `hosted`: ChatGPT, Claude Cowork, Devin, Amazon Q Developer, Sourcegraph Cody, Augment, and similar platforms that need API, MCP, Actions, Apps SDK, or uploaded-instruction adapters.

## Audit Checklist

1. Run `skill-router --version`.
2. Run `skill-router skill search <topic>` and `skill-router skill <name>` for the current task.
3. Verify `%USERPROFILE%\universal-ai-skills-library\skills` is the preferred corpus source.
4. Verify `%USERPROFILE%\go\bin\skill-router.exe` is on PATH.
5. Check `skill-router sync matrix` for supported platform adapters.
6. Check `skill-router doctor`.
7. Check `skill-router mcp status`.
8. Confirm global AI instruction files point to the CLI and do not contain generated skills tables.
9. Confirm plugin manifests describe the CLI router rather than embedding the whole corpus.

## Client Surfaces

Keep these compact:

- `%USERPROFILE%\.agent\AGENTS.md`
- `%USERPROFILE%\.codex\AGENTS.md`
- `%USERPROFILE%\.claude\CLAUDE.md`
- `%USERPROFILE%\.cursor\rules\openskills.md`
- `%USERPROFILE%\.gemini\GEMINI.md`
- `%USERPROFILE%\.config\opencode\AGENTS.md`

Minimum rule:

```text
Use `skill-router skill <name>` to load skills on demand.
Use `skill-router skill search <query>` when the skill name is unknown.
Keep always-loaded instructions compact.
MCP bridges are optional and only needed for persistent endpoint workflows.
```

## MCP Policy

Keep these running only when useful:

- MemPalace: durable memory
- Context Mode: long-output routing and indexed session continuity
- Skill Seekers: skill generation, packaging, upload, install workflows
- Lightpanda: persistent browser/CDP automation

If Docker is not running, Lightpanda may be down. That should be a warning, not a global failure, unless the active task needs Lightpanda browser automation.

## Done Criteria

- `skill-router --version` works.
- `skill-router skill <name>` prints the requested skill.
- `manifest.json` matches the actual `skills/` tree.
- agent instruction files are compact.
- plugin metadata points to CLI-first usage.
- optional MCP bridge status is understood and documented.
