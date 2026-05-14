---
name: instagram-cli
description: Use when working with supreme-gg-gg/instagram-cli or Instagram from the terminal: logging in or checking auth, listing or searching inbox threads, sending Instagram direct messages, reading messages, replying, unsending, downloading message media, viewing feed/stories/notifications/profile TUI screens, configuring the Instagram CLI, or diagnosing the local `instagram-cli` install. Prefer this for agent one-turn Instagram workflows that need JSON output. Do not use for generic social-media copywriting or unrelated Instagram marketing strategy.
aliases:
  - instagram cli
  - instagram-cli
  - instagram-skill
  - i7m
  - i7m instagram
  - "@i7m/instagram-cli"
  - instagram tui
  - instagram terminal client
  - instagram direct messages
  - instagram messages cli
  - instagram dm cli
  - instagram inbox
  - instagram auth
  - instagram one-turn commands
---

# Instagram CLI Universal Adapter

## Source

- Upstream repository: `https://github.com/supreme-gg-gg/instagram-cli`
- Local source checkout: `%USERPROFILE%\.instagram-cli-source\instagram-cli`
- Installed command: `instagram-cli`
- Router delegator: `skill-router instagram`
- Installed npm package link: `%APPDATA%\npm\node_modules\@i7m\instagram-cli`
- User config/state root: `%USERPROFILE%\.instagram-cli`
- Config file: `%USERPROFILE%\.instagram-cli\config.ts.yaml`
- Logs: `%USERPROFILE%\.instagram-cli\logs`
- Universal stack source registry: `%USERPROFILE%\.universal-ai-stack\config\source-integrations.json`

## Operating Model

Use this as a shared CLI source integration, not as copied skill bodies or a
background service. The universal install uses:

- one external source checkout for README-grounded reference,
- one globally linked `instagram-cli` command built from that checkout,
- this canonical skill for routing and workflow instructions,
- compact `skill-router` wrappers in each AI client.

No MCP bridge is required for normal Instagram CLI work.

## Safety And Account Boundaries

This is an unofficial Instagram client and may violate Meta/Instagram terms of
service. Use it only when the user explicitly asks for Instagram CLI work and
accepts that risk.

Never print, paste, store, or commit:

- `%USERPROFILE%\.instagram-cli` session files,
- `%USERPROFILE%\.instagram-cli\logs`,
- usernames/passwords,
- two-factor codes,
- cookies or device/session identifiers,
- private direct messages,
- downloaded private media.

Do not attempt to log in for the user. If auth is missing, tell the user to run:

```powershell
instagram-cli auth login
```

Check the active account without exposing private session files:

```powershell
instagram-cli auth whoami
```

Switch accounts only when the user names the intended account:

```powershell
instagram-cli auth switch <username>
```

## Agent One-Turn Commands

Prefer non-interactive one-turn commands for AI agents. They print to stdout and
exit. Use `--output json` or `-o json` whenever the result will be parsed or used
by an agent.

List inbox threads:

```powershell
instagram-cli inbox --limit 10 --output json
```

Search threads:

```powershell
instagram-cli inbox --search "alice" --output json
```

Send a text message:

```powershell
instagram-cli send <thread> --text "Message text" --output json
```

Send media:

```powershell
instagram-cli send <thread> --file "C:\path\photo.jpg" --type photo --output json
```

Read messages:

```powershell
instagram-cli read <thread> --limit 20 --output json
```

Read and mark seen only when the user explicitly wants that side effect:

```powershell
instagram-cli read <thread> --limit 20 --mark-seen --output json
```

Download media from a message:

```powershell
instagram-cli read <thread> --message-id <id> --download "C:\path\media" --output json
```

Reply to a specific message:

```powershell
instagram-cli reply <thread> --message-id <id> --text "Reply text" --output json
```

Unsend a message only after explicit user confirmation:

```powershell
instagram-cli unsend <thread> --message-id <id> --output json
```

`<thread>` can be a thread ID, username, or fuzzy thread title. Prefer a thread
ID from `instagram-cli inbox --output json` to avoid accidental misrouting.

All one-turn commands accept:

```powershell
--username <username>
```

Use it for multi-account work instead of switching the default account when the
user wants a single command to target a specific login.

## TUI Commands

Use TUI commands only when the user wants an interactive terminal experience:

```powershell
instagram-cli chat -u <username> -t <title>
instagram-cli feed
instagram-cli stories
instagram-cli notify
instagram-cli profile <username>
```

Do not launch long-lived TUI sessions from background jobs or automatic hooks.

## Configuration

Inspect configuration:

```powershell
instagram-cli config
```

Set a config key:

```powershell
instagram-cli config image.protocol ascii
instagram-cli config feed.feedType list
```

Open config for manual edits only when the user asks:

```powershell
instagram-cli config edit
```

## Diagnostics

Check local install:

```powershell
instagram-cli --version
instagram-cli --help
skill-router instagram status
```

If commands fail after Instagram changes its API, capture only redacted error
summaries. Do not paste raw logs unless secrets and private message content have
been removed.

## Router Selection Rules

Prefer this skill only when the task is specifically about `instagram-cli`,
Instagram direct messages, Instagram inbox/search/read/reply/unsend/media
commands, Instagram terminal/TUI workflows, multi-account CLI usage, or local
diagnostics for the `instagram-cli` command.

Reject this skill for generic Instagram marketing, content strategy, image
generation, social caption writing, unrelated social media analysis, or any
request that does not need the local Instagram CLI.
