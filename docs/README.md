# Documentation Hub

This folder documents the public, portable Universal AI Skills setup. Start with
the quickstart if you only want to install the router. Use the architecture,
compatibility, and connection docs when integrating new AI clients or changing
the local Universal AI Stack.

## Start Here

| Page | Use it for |
| --- | --- |
| [Quickstart](QUICKSTART.md) | Install from a clean clone and run the first router commands. |
| [Install Modes](INSTALL_MODES.md) | Choose router-only, stack-enabled, startup, or validation-focused installs. |
| [Universal AI Setup](UNIVERSAL_AI_SETUP.md) | Understand model routing, local Qwen fallback, embeddings, Hermes, Paperclip, memory, Context Mode, Lightpanda, NotebookLM MCP CLI, x-cli, Instagram CLI, Crawl4AI, web search, GBrain, and GSkills/GStack. |
| [Universal AI Connection Configs](UNIVERSAL_AI_CONNECTION_CONFIGS.md) | See the exact repo-owned and machine-generated files that connect each AI client. |
| [Source Integrations](SOURCE_INTEGRATIONS.md) | Understand the public-safe source layer for Lightpanda, Context Mode, MemPalace, NotebookLM MCP CLI, x-cli, Instagram CLI, Crawl4AI, web search, GBrain, and GSkills/GStack. |

## Architecture And Compatibility

| Page | Use it for |
| --- | --- |
| [Architecture](ARCHITECTURE_V2.md) | Understand the router-first design, source-of-truth layout, and naming contract. |
| [Universal Compatibility](UNIVERSAL_COMPATIBILITY.md) | Decide whether a client should use a wrapper skill, compact repo instructions, hosted/API/MCP integration, or report-only support. |
| [Agent Support Matrix](AGENT_SUPPORT_MATRIX.md) | Check current adapter policy for Codex, Claude, Cursor, Hermes, Paperclip, OpenCode, Kimi, Qwen, OpenHands, and other clients. |
| [AI Repo Tools Summary](AI_REPO_TOOLS_SUMMARY.md) | Review the local AI/repo tooling surface around `skill-router`. |

## Public Release And Maintenance

| Page | Use it for |
| --- | --- |
| [Public Release Checklist](PUBLIC_RELEASE_CHECKLIST.md) | Validate that the repository is safe to publish or tag. |
| [Third-Party Source Repos](THIRD_PARTY_SOURCE_REPOS.md) | Track indexed external skill sources and read-only source policy. |
| [Normalization Report](NORMALIZATION_REPORT.md) | Review corpus normalization and alias decisions. |
| [Dropped Items](DROPPED.md) | See removed or intentionally excluded material. |
| [Design And Messaging](DESIGN_AND_MESSAGING.md) | Keep the GitHub About text, README hero, plugin metadata, and social-preview copy consistent. |

## Maintenance Rule

When a connection path changes, update both the repo-owned config or script and
the matching documentation page. Do not rely on hand-edited machine-local files
as the long-term source of truth.
