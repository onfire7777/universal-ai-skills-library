---
name: arc-design-review
description: "Three-phase architecture design review pipeline: Architect creates initial design, Refiner improves it, Critic challenges it. Then facilitated discussion produces an Architecture Decision Record. Use for system design, major refactors, new features, or any significant architectural decision. Triggers on: 'architecture review', 'design review', 'ADR', 'architecture decision', 'system design review', 'review this architecture', 'should we build it this way'."
---

# Architect-Refine-Critique Design Review

## Overview

A three-phase design review pipeline that produces an Architecture Decision Record.

## Phase 1: Architect
- Research the codebase and project conventions
- Identify constraints, requirements, integration points
- Create initial design with component diagrams
- Document structural decisions with reasoning

## Phase 2: Refine
- Review the initial design for gaps
- Improve based on best practices and current trends
- Challenge assumptions, verify feasibility
- Produce refined design document

## Phase 3: Critique
- Adversarial evaluation of the refined design
- Rate findings by severity (Critical, High, Medium, Low)
- Identify risks, failure modes, scaling concerns
- Propose alternatives where design is weak

## Phase 4: Review & ADR
Facilitated discussion with the user:
1. Present initial design summary — discuss strengths/weaknesses
2. Present refinements — discuss whether they addressed right issues
3. Present critique findings one by one — user decides Accept/Reject/Defer
4. Synthesize into ADR:

```markdown
# ADR: [Title]
**Status:** Proposed | **Date:** [today]
## Context: [problem + constraints]
## Decision: [chosen architecture + reasoning]
## Consequences: [positive, negative, mitigations]
## Alternatives Considered: [rejected options + why]
## Process Improvements: [patterns identified during review]
```

## Rules
🚨 Ground every decision in a named principle (SRP, DIP, separation-of-concerns, etc.)
🚨 Discuss ALL phases in order — don't skip.
🚨 Don't decide for the user. Present, discuss, let them choose.