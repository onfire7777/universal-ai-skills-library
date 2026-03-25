# AI-Assisted Prototyping Guide

> Purpose: use AI to accelerate Forge safely without handing it architecture, domain truth, or quality judgment.

## Contents

- Tool map
- Phase-by-phase use
- Prompt strategy
- Safe use boundaries
- Quality checks

## Tool Map

| Category | Examples | Best use | Common weakness |
|---|---|---|---|
| General AI chat | Claude, ChatGPT | concept shaping, boilerplate, quick variants | weak long-range state management |
| Design-to-code tools | v0, Lovable, Make-style builders | fast visual scaffolds | low reliability for nuanced behavior |
| IDE assistants | Claude Code, Cursor, Copilot | local integration and repetitive edits | quality depends on existing project context |

## Best Use by Forge Phase

| Forge phase | Use AI for | Do not delegate |
|---|---|---|
| `SCAFFOLD` | option generation, stack choice, slice framing | deciding product truth without evidence |
| `STRIKE` | boilerplate, starter TSX, mock handlers, sample data | complex domain rules or final architecture |
| `COOL` | quick bug triage, edge-case prompts, checklist generation | final quality judgment |
| `PRESENT` | demo steps, summary bullets, validation prompts | stakeholder decision-making |

## Prompt Strategy

Use staged prompts instead of one giant ask:

1. Ask for the minimum skeleton.
2. Ask for one interaction.
3. Ask for one mock strategy.
4. Ask for review against the current hypothesis.

Good prompt qualities:
- explicit scope
- explicit time-box
- declared prototype type: `Throwaway` or `Evolutionary`
- declared framework and output files
- one request per pass

## Safe AI Boundaries

AI is good at:
- CRUD UI skeletons
- form and table scaffolds
- starter `MSW` handlers
- realistic sample data shapes
- repetitive boilerplate

AI is weak at:
- hidden business rules
- multi-service orchestration
- production-ready state architecture
- long-lived maintainability decisions

Safe Vibe Coding scope:
- throwaway demos
- internal exploration
- quick local-only prototypes

Ask first before relying on AI-generated structure when:
- the prototype is intended to evolve into production
- the domain has non-obvious business constraints
- the prototype spans multiple services

## One-Day Prototype Pattern

Morning:
- frame the hypothesis
- generate the base UI
- connect fixtures or `MSW`

Afternoon:
- add interaction and failure states
- run `COOL` checks
- prepare the demo and decision notes

## Quality Checks for AI-Generated Code

Check these before accepting AI output:

- build passes
- component actually renders
- mock contract matches expected field names
- no unexplained `any`
- no silent missing states
- no invented API behavior
- no hidden third-party dependency creep

Common failure signals:
- wrong framework API usage
- stale library syntax
- inconsistent data shape
- missing loading / error state
- unexplained business assumptions

## Forge Integration Rules

- Use AI to accelerate scaffolding, not to define the prototype’s meaning.
- Require a human review during `COOL`.
- Resolve TypeScript or build errors before `PRESENT`.
- Document why a given AI tool was chosen when the prototype survives beyond a quick check.
