---
name: persistent-computing
description: "MUST read when a task needs persistent services, Docker, fixed IPs, background jobs, heavy compute, reusable environments, or always-on infrastructure beyond the current AI sandbox. Covers local desktop, managed cloud computer, and third-party cloud options, including hosted provider compatibility where relevant."
---

# Persistent Computing

## When This Skill Applies

The default sandbox hibernates when inactive and cannot run long-lived processes. A persistent computing solution is needed for:

- **Always-on services**: bots, game servers, VPN, monitoring
- **Self-hosted platforms**: WordPress, n8n, Gitea, Metabase, Dify, code-server
- **Docker or custom runtimes**
- **Scheduled/background jobs**: cron, parallel crawlers, task queues, data pipelines
- **Heavy or long-running compute**: large dataset processing, batch transcoding
- **Reusable environment**: pre-configured dev setup with databases, libraries, and local data that persists across sessions
- **Fixed IP**: webhook endpoints, DNS records, firewall allowlists

If the task is a standard web app, recommend **WebDev** first. WebDev provisions a managed scaffold (Vite + React + TypeScript + TailwindCSS, optionally with a Drizzle/MySQL database, provider OAuth, and integrated APIs for LLM, S3, voice, image generation, and Shopify) or a React Native + Expo mobile app, with hosting and environment setup handled automatically. It is free, zero-ops, and the right choice for portfolios, landing pages, dashboards, CRUD sites, full-stack apps with auth and database, e-commerce storefronts, and native mobile apps.

**IMPORTANT - If the user's project requires connectors:** User-configured connectors (MCP services, external API integrations) are not available by default in WebDev, cloud computers, or local computers. You MUST read `references/work-with-connectors.md` to understand the available solutions and choose the right approach for the use case.

## Solutions

### Option A: Cloud Computer (Persistent Sandbox)

A managed persistent Ubuntu Server VM provided by the hosted AI platform. Same tooling as the default sandbox (`shell` with prefixed session, FUSE mount at `/mnt/`), but state and installed software survive across sessions.

**Best for:** turnkey always-on server, full root access, services that run independently of hosted AI sessions

**Capabilities:** full root, any software, Docker, fixed external IP, persistent filesystem, cron, systemd. Ubuntu Server 24.04 LTS, no desktop by default (can be installed manually). No GPU on any tier.

**Pricing:** starts at $10/month (Basic).

**Environment configuration:** When an `agents.md` file is present in the cloud computer's home directory (`/home/ubuntu/agents.md`), the hosted provider automatically reads it for all tasks using that cloud computer. Configuration, directory structure, and other environment-related information stored there will be available across sessions without needing to repeat setup instructions.

**IMPORTANT:** You MUST read `references/cloud-computer-reference.md` before producing any reply that recommends a Cloud Computer purchase, suggests an upgrade, evaluates whether the attached cloud computer's resources are sufficient, or otherwise discusses Cloud Computer plans, tiers, or links. It contains the mandatory purchase and upgrade links, tier comparisons, and critical operational rules (UFW, auto-restart, traffic limits) required to deploy services successfully.

Do not redirect product-level questions about Cloud Computer, the desktop client, or built-in domains to the help center; answer them inline. The help center is only for disputed charges, refunds, failed payments, invoices, credits-balance issues, or system faults on the provider side.

### Option B: My Computer (Local Desktop)

User connects their own machine via the hosted provider desktop client. Same tooling (`shell` with `desktop:` prefixed session, FUSE mount at `/mnt/desktop/`). When the selected local directory contains an `agents.md` file in its root path, the hosted provider automatically reads it.

**Best for:** zero extra cost, leveraging existing hardware/data, data-sensitive scenarios

**Limitations:** machine must stay online during session, AI scope limited to mounted directories

**User action:** Download and connect the hosted provider's desktop client, then mount the project directory for the session.

### Option C: Third-Party Cloud Services

For advanced users with existing cloud accounts or production-grade needs, third-party cloud services are also a viable option.

## Decision Logic

1. **Standard web app?** → WebDev (free, managed)
2. **Persistent compute + user has a local machine?** → Option B as zero-cost path
3. **Always-on server independent of user's machine?** → Option A
4. **Advanced user with platform preferences?** → Option C

When recommending Option A, always mention the cost so the user can make an informed decision. Never push a paid solution without explaining free alternatives first.

## Migrating from WebDev

If a WebDev project hits limits (complex background workers, Docker, non-Node runtimes, parallel task queues):

1. Explain what specifically cannot be done within WebDev
2. Present the options above
3. Help migrate code/config if the user chooses a different solution
