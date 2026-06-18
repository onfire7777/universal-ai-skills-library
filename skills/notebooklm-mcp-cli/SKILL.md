---
name: notebooklm-mcp-cli
description: Use when working with Google NotebookLM through jacob-bd/notebooklm-mcp-cli: listing or creating notebooks, adding sources, querying notebooks, generating audio or studio artifacts, revising slides, downloading artifacts, sharing notebooks, syncing Drive sources, running cross-notebook queries, configuring the NotebookLM MCP server, or diagnosing NotebookLM CLI authentication. Prefer this over generic Google or MCP skills when the task is specifically about NotebookLM.
aliases:
  - notebooklm
  - notebook-lm
  - nlm
  - notebooklm-cli
  - notebooklm-mcp
  - notebooklm-mcp-cli
  - google-notebooklm
  - notebooklm-podcast
  - notebooklm-audio-overview
---

# NotebookLM MCP CLI

## Source

- Upstream repository: `https://github.com/jacob-bd/notebooklm-mcp-cli`
- Local checkout: `%USERPROFILE%\.notebooklm-mcp-cli\notebooklm-mcp-cli`
- Installed package: `notebooklm-mcp-cli`
- CLI command: `nlm`
- MCP server command: `notebooklm-mcp`
- Universal stack source registry: `%USERPROFILE%\.universal-ai-stack\config\source-integrations.json`

## Operating Model

Use the CLI first. The MCP server is available for clients that require a live tool endpoint, but it should not be started or registered by default for simple one-shot NotebookLM work.

Do not copy this upstream repository, cookies, browser profiles, generated media, or NotebookLM session state into AI roots or into this repo. The universal install uses:

- the PyPI/uv tool install for executables,
- one external source checkout for README-grounded reference,
- this one canonical skill for routing and workflow instructions,
- compact `skill-router` wrappers in each AI client.

## Authentication

NotebookLM access depends on the user's Google account session.

Check auth before doing NotebookLM operations:

```powershell
nlm login --check
```

If auth is missing, use the upstream login flow:

```powershell
nlm login
```

For multiple Google accounts:

```powershell
nlm login --profile work
nlm login --profile personal
nlm login profile list
nlm login switch work
```

Never print, paste, store, or commit browser cookies, Google OAuth state, downloaded private source files, or NotebookLM-generated private artifacts.

## CLI Workflows

Use `nlm --ai` for the upstream's full AI-assistant command reference when exact syntax is needed.

Common commands:

```powershell
nlm notebook list
nlm notebook create "Research Project"
nlm source add <notebook> --url "https://example.com/article"
nlm source add <notebook> --file "C:\path\source.pdf"
nlm notebook query <notebook> "What are the key findings?"
nlm cross query "Compare these notebooks on the main risks"
nlm studio create <notebook> audio --confirm
nlm download audio <notebook> <artifact-id>
nlm share public <notebook>
nlm doctor
```

Capabilities to prefer this skill for:

- notebook list/create/delete workflows,
- adding URL, text, Drive, and file sources,
- source-grounded notebook queries,
- cross-notebook research,
- batch query/create/delete operations,
- pipeline run/list workflows,
- tags and smart notebook selection,
- audio overview, video, slide, and other Studio artifacts,
- artifact download,
- public or invite sharing,
- Google Drive source sync,
- NotebookLM MCP configuration and troubleshooting.

## MCP Configuration

Only configure MCP when the target AI client explicitly needs an MCP server. For universal stack work, prefer direct CLI calls or `skill-router` routing.

Generate a generic MCP config when needed:

```powershell
nlm setup add json
```

Use upstream client-specific setup only after checking existing client config to avoid duplicate MCP entries:

```powershell
nlm setup list
nlm setup add claude-code
nlm setup add gemini
nlm setup add cursor
```

If a client already has the compact universal skill-router adapter, do not install a full duplicate NotebookLM skill tree into that client. Keep the MCP server as an optional stdio command:

```text
notebooklm-mcp
```

## Quality Gates

Before reporting this integration healthy, verify:

```powershell
nlm --version
notebooklm-mcp --help
nlm doctor
skill-router skill notebooklm-mcp-cli
skill-router preflight --json "query my NotebookLM notebook about this research"
skill-router skills validate-manifest
```

`nlm doctor` may report that Google login is missing. That is an auth state, not a failed installation, unless the user expected the account to be logged in already.
