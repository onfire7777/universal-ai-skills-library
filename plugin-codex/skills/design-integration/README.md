# Design Integration

Design a Uniswap integration: choose between Trading API, SDKs, or direct contracts; define architecture, dependencies, and security posture for projects integrating Uniswap.

→ **[SKILL.md](SKILL.md)** — Full skill specification and workflow.

## Installation

Install into Claude Code or Cursor with:

```bash
npx skills add https://github.com/wpank/Agentic-Uniswap/tree/main/.ai/skills/design-integration
```

Or via Clawhub:

```bash
npx clawhub@latest install design-integration
```

## When to use

Use this skill when:

- You're building a **new app, bot, or agent** that needs to integrate Uniswap.
- You need a **blueprint** for how to connect (API vs SDK vs contracts), where to run, and how to secure it.
- You want to surface **rate limits, failure modes, and security best practices** up front.

## Example prompts

- "Design a Uniswap integration architecture for a cross-chain trading bot using the Trading API."
- "Propose how my web app on Base should integrate Uniswap for swaps and LP, including security considerations."
- "Outline an architecture for an agent that rebalances a treasury via Uniswap across Ethereum and Arbitrum."
