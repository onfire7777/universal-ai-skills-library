# Architecture v2

## Goal

Make the AI setup fully available without loading every skill into every agent context.

## Decisions

1. `skills/` is the only skill source of truth.
2. `manus-cli/` is the only normal runtime router.
3. `plugin/` is metadata and compact instructions, not another skill copy.
4. MCP bridges are optional and reserved for persistent endpoints.
5. Local AI clients should carry small operating rules that point to `manus`.

## Runtime Paths

```text
C:\Users\burni\manus-skills-library\skills       source corpus
C:\Users\burni\manus-skills-library\manus-cli   CLI source
C:\Users\burni\go\bin\manus.exe                 installed CLI
C:\ProgramData\manus-mcps                       optional bridge logs/scripts
```

## Skill Loading

Preferred:

```bash
manus skill <name>
manus skill search <query>
```

Fallback:

```bash
npx openskills read <name>
```

Do not copy full skill bodies into `AGENTS.md`, `CLAUDE.md`, Cursor rules, or other always-loaded instruction files.

## MCP Policy

Use direct CLI for:

- skill loading
- skill search
- audits
- local report generation
- file organization
- command orchestration

Use MCP for:

- MemPalace durable memory
- Context Mode indexed long-output routing
- Skill Seekers dynamic skill generation workflows
- Lightpanda persistent browser/CDP workflows

If an MCP endpoint is down but the equivalent CLI works, the system remains usable.

## Local Client Surfaces

Each AI client should get the same compact rule:

```text
Use `manus skill <name>` to load skills on demand.
Use `manus skill search <query>` when the skill name is unknown.
Keep always-loaded instructions compact.
MCP bridges are optional and only needed for persistent endpoint workflows.
```

## Verification

```bash
manus --version
manus skill persistent-computing
manus skill search debugger
manus skill list
manus mcp status
go test ./...
```
