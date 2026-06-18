---
name: provider-api
description: Universal hosted-provider API adapter guide. Use the skill-router provider-api CLI for compatible project/task/file/webhook/agent APIs, with explicit endpoint and key configuration instead of provider-specific defaults.
---

# Provider API

Use this skill when a workflow needs to call a hosted project provider through the universal `skill-router provider-api` CLI. The adapter is intentionally neutral: it does not assume any provider endpoint, API key name, or hosted UI.

## CLI Contract

Set an endpoint and key explicitly before making API calls:

```bash
export SKILL_ROUTER_API_BASE="https://provider.example/v2"
export SKILL_ROUTER_API_KEY="..."
skill-router provider-api tasks list
```

You can also pass the endpoint per command:

```bash
skill-router provider-api --api-base "https://provider.example/v2" projects list
```

## Supported Surfaces

The CLI exposes common hosted-provider surfaces:

- `tasks` for task lifecycle calls
- `projects` for reusable project context
- `files` for provider-managed file references
- `webhooks` for callback configuration
- `agents` for provider-hosted agent records
- `connectors`, `skills`, `usage`, and `websites` when the provider supports them

## Provider-Specific Docs

For exact request schemas, authentication headers, rate limits, or provider-only behavior, load the explicit provider skill or documentation for that platform. Do not hard-code provider endpoints into universal setup, install, doctor, or router defaults.
