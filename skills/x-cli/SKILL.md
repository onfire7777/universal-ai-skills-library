---
name: x-cli
description: Use when working with sferik/x-cli or the X API from the command line: authorizing X accounts, listing accounts, posting or replying, reading timelines, searching posts or users, streaming X events, managing follows, mutes, blocks, lists, direct messages, trends, places, or diagnosing the local `x` CLI install. Do not use for generic algebra variable "x" or unrelated command-line work.
aliases:
  - sferik-x-cli
  - sferik x-cli
  - x api cli
  - x command line
  - twitter cli
  - twitter command line
  - xrc
  - .xrc
  - x-api
---

# x-cli Universal Adapter

## Source

- Upstream repository: `https://github.com/sferik/x-cli`
- Local checkout: `%USERPROFILE%\.x-cli\x-cli`
- Primary command: `%USERPROFILE%\.local\bin\x.exe`
- Compatibility alias: `%USERPROFILE%\.local\bin\x-cli.exe`
- Router delegator: `skill-router xcli`
- X profile config: `%USERPROFILE%\.xrc`
- Legacy read fallback: `%USERPROFILE%\.trc`
- Universal stack source registry: `%USERPROFILE%\.universal-ai-stack\config\source-integrations.json`

## Operating Model

Use this as a shared CLI source integration, not as copied skill bodies or a
background service. The universal install uses:

- one external source checkout for README-grounded reference,
- one built Rust release on PATH as `x`,
- a namespaced `x-cli` skill for routing and workflow instructions,
- compact `skill-router` wrappers in each AI client.

No MCP bridge is required for normal X API work.

## Authentication And Secrets

`x-cli` supports OAuth 1.0a and OAuth 2.0 and stores account/profile data in
`~/.xrc`. If `~/.xrc` is missing, upstream can read legacy `~/.trc` and migrate
on write.

Never print, paste, store, or commit:

- `~/.xrc` or `~/.trc`,
- consumer keys or secrets,
- OAuth access tokens or refresh tokens,
- direct messages, private account data, or private API output.

Start or inspect auth with:

```powershell
x authorize
x accounts
x set active <account>
```

Use `--profile <FILE>` when the user explicitly wants an isolated profile file.

## Common Workflows

Check the install:

```powershell
x version
skill-router xcli status
```

Post or reply:

```powershell
x update "post text"
x reply <post-id> "reply text"
```

Read or search:

```powershell
x timeline
x mentions
x status <post-id>
x search all "query"
x search users "name or handle"
```

Account and social graph actions:

```powershell
x accounts
x follow <screen-name>
x unfollow <screen-name>
x favorite <post-id>
x retweet <post-id>
x mute <screen-name>
x block <screen-name>
```

Lists, direct messages, trends, and places are available through the upstream
command tree. Run command-specific help before using a less common subcommand:

```powershell
x --help
x search --help
x stream --help
x list --help
x dm --help
```

## Streaming Safety

Streaming commands keep HTTP connections open. For automation and AI-agent
sessions, set `X_STREAM_MAX_EVENTS` unless the user explicitly wants an
open-ended stream:

```powershell
$env:X_STREAM_MAX_EVENTS = "25"
x stream all
```

Do not start persistent stream loops in the background without an explicit user
request and a stop plan.

## Router Selection Rules

Prefer this skill only when the task is specifically about `x-cli`, X API
automation, X/Twitter account operations, timelines, posts, search, streams,
direct messages, follows, blocks, mutes, lists, trends, places, or local
diagnostics for the `x` command.

Reject this skill for generic "x" wording, algebra, unrelated CLI tooling, or
generic social-media writing tasks that do not need the X API.
