# MCP Server

`skill-router` can run as a **Model Context Protocol** server so that hosted clients (and any MCP-capable agent) reach the exact same routing engine as the CLI.

Source: [`skill-router-cli/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/skill-router-cli) (`serve` command, `internal/mcpcli`)

---

## Starting the server

```bash
skill-router serve
```

- **Protocol:** JSON-RPC 2.0, MCP version `2024-11-05`.
- **Transport:** stdio (stdin/stdout). It is hand-rolled — no external MCP SDK dependency.
- **Engine:** the same `skillservice` engine the CLI uses. There is no separate routing path, so MCP and CLI behavior are identical by construction (see [Architecture](Architecture.md#cli-and-mcp-are-two-faces-of-one-engine)).

---

## Exposed tools

The server exposes four tools — the [four canonical verbs](Architecture.md#the-router-first-pipeline-search--route--load):

| Tool | Purpose |
|---|---|
| `route` | Return the single best-match skill for a prompt (with the route / ambiguous / no_route decision). |
| `search_skills` | Full-text search the corpus by query and return ranked candidates. |
| `load_skill` | Fetch a skill's `SKILL.md` body by name. |
| `compose` | Assemble a working set, or — with the `pipeline` argument — a multi-step DAG plan. |

---

## Registering the server with a client

Point any MCP client at the `skill-router serve` command over stdio. A typical MCP client config entry:

```json
{
  "mcpServers": {
    "skill-router": {
      "command": "skill-router",
      "args": ["serve"]
    }
  }
}
```

For Codex-style clients the equivalent goes in `~/.codex/config.toml`; for Claude clients, the MCP server block in the client's config. The [Installation & Setup](Installation-and-Setup.md#per-client-wiring) page lists per-client config locations.

---

## Windows managed bridge

On Windows there is an additional command group that runs and supervises the bridge as a managed process:

```bash
skill-router mcp status
skill-router mcp start
skill-router mcp stop
skill-router mcp restart
skill-router mcp logs
skill-router mcp watchdog
skill-router mcp setup
```

The bridge config directory is resolved via the `MCP_DIR` environment variable. `skill-router mcp status` is a safe, read-only health check.

---

## When to use MCP vs. the CLI

| Use the **CLI** when | Use **MCP** when |
|---|---|
| You drive the router from a hook/script (`preflight`, `skill`). | The client speaks MCP natively and you want tool-call integration. |
| You want one-shot commands in a terminal. | The client is hosted (no local skill directory) — see the [hosted layer](Agent-Support-Matrix.md#hosted-layer-api--mcp-only). |

Both reach the same engine, so the routing decisions are the same either way.
