---
name: design-hook
description: >-
  Design a Uniswap V4 hook architecture without code generation. Use when user
  wants to plan a hook, understand which callbacks to use, or review an
  architecture before building. Returns a design document, not code.
allowed-tools: >-
  Read, Glob, Grep,
  Task(subagent_type:hook-builder),
  mcp__uniswap__get_supported_chains
model: opus
---

# Design Hook

## Overview

Designs a Uniswap V4 hook architecture without generating code. Delegates to the `hook-builder` agent in design-only mode to produce a comprehensive design document covering: which callbacks are needed, hook flag requirements, state management approach, gas estimates, security considerations, and architecture decisions. Use this to plan before building, or to evaluate feasibility.

## When to Use

Activate when the user asks:

- "Design a hook for..."
- "What callbacks do I need for..."
- "Hook architecture for..."
- "Plan a V4 hook"
- "Is it possible to build a hook that..."
- "What would a dynamic fee hook look like?"
- "Help me think through a hook design"
- "Which flags do I need for a TWAMM?"

## Parameters

| Parameter | Required | Default | Description |
| --- | --- | --- | --- |
| behavior | Yes | -- | Hook behavior description (e.g., "limit orders", "dynamic fees", "oracle pricing") |
| constraints | No | -- | Gas budget, security requirements, or specific design constraints |
| integrations | No | -- | External systems the hook needs to interact with (oracles, governance, staking) |

## Workflow

1. **Extract parameters** from the user's request: identify the hook behavior, constraints, and any external integrations.

2. **Delegate to hook-builder in design-only mode**: Invoke `Task(subagent_type:hook-builder)` with explicit instruction to produce a design document only -- no code generation, no file writes. The hook-builder will:
   - Analyze the requirements and determine which V4 callbacks are needed
   - Map callbacks to hook flags and validate the combination
   - Design the state management approach (what storage, what data structures)
   - Estimate gas overhead per callback
   - Identify security considerations specific to this hook design
   - Evaluate feasibility and flag any concerns

3. **Present the design document** to the user covering:
   - Callbacks needed and why each is required
   - Hook flags and bitmask
   - State management design (storage variables, data structures, access patterns)
   - Gas estimates and performance implications
   - Security considerations and mitigations
   - Architecture decisions with rationale
   - Comparison with alternative approaches if applicable

## Output Format

Present a structured design document:

```text
V4 Hook Design: Dynamic Fee Hook

  Callbacks Required:
    - beforeSwap: Read volatility oracle, calculate dynamic fee
    - beforeInitialize: Set initial fee parameters and oracle address

  Hook Flags: BEFORE_SWAP_FLAG | BEFORE_INITIALIZE_FLAG
  Bitmask: 0x2080

  State Management:
    - volatilityOracle: IVolatilityOracle (immutable, set in constructor)
    - baseFee: uint24 (configurable by owner)
    - maxFee: uint24 (cap to prevent excessive fees)
    - feeMultiplier: uint24 (scales with volatility)

  Gas Estimates:
    beforeSwap: ~30,000 gas (oracle read + fee calculation)
    beforeInitialize: ~25,000 gas (one-time setup)

  Security Considerations:
    - Oracle manipulation: Use TWAP, not spot price
    - Fee cap: Enforce maxFee to protect traders
    - Owner control: Fee parameters updatable by owner only

  Architecture Decisions:
    - Using beforeSwap (not afterSwap) to set fee before execution
    - External oracle for volatility data rather than on-chain calculation
    - Fee bounded between baseFee and maxFee for predictability

  Alternative Approaches:
    - On-chain volatility calculation (higher gas, no oracle dependency)
    - Fixed fee tiers with governance voting (simpler, less responsive)
```

## Important Notes

- This skill produces a design document only -- no code is generated and no files are written.
- The design document provides enough detail to proceed with `build-hook` when the user is ready.
- If the hook design is infeasible (e.g., requires callbacks that V4 doesn't support), this will be clearly communicated.
- Gas estimates are approximations based on typical implementations -- actual gas depends on implementation details.

## Error Handling

| Error | User-Facing Message | Suggested Action |
| --- | --- | --- |
| `VAGUE_REQUIREMENTS` | "Need more detail about the desired hook behavior." | Describe specific behavior (e.g., "limit orders that execute at tick boundaries") |
| `UNSUPPORTED_CALLBACK` | "V4 does not support the requested callback." | Review available V4 callbacks and adjust requirements |
| `INFEASIBLE_DESIGN` | "This hook design is not feasible with current V4 capabilities." | Simplify requirements or consider alternative approaches |
