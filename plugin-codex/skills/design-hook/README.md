# Design Hook

Design a Uniswap V4 hook architecture (callbacks, behavior, invariants). No code generation; produces a detailed design document only.

→ **[SKILL.md](SKILL.md)** — Full skill specification and workflow.

## Installation

Install into Claude Code or Cursor with:

```bash
npx skills add https://github.com/wpank/Agentic-Uniswap/tree/main/.ai/skills/design-hook
```

Or via Clawhub:

```bash
npx clawhub@latest install design-hook
```

## When to use

Use this skill when:

- You want to **design a V4 hook architecture** before writing any Solidity.
- You need to clarify **which callbacks to use**, state layout, and safety considerations.
- You want a **spec document** that can be handed to implementers or auditors.

## Example prompts

- "Design a Uniswap V4 hook that charges higher fees during volatile periods."
- "Draft an architecture for a V4 hook that streams a share of fees to a DAO treasury."
- "Design an anti-sniping hook for a new token launch on V4, including which callbacks to implement."
