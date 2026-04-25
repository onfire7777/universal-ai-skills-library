---
name: skill-connection-map
description: Maps the hidden connections between all 630+ installed skills, enabling cross-pollination of ideas across domains. Use when you want to find unexpected connections, when a problem in one domain might benefit from thinking in another domain, when you want to transfer patterns between fields, or when you want to see the bigger picture of how skills relate. This is the "neural network" that connects all skills into a living, interconnected knowledge system.
license: Unspecified
---
# Skill Connection Map

## Purpose

Individual skills are powerful. But the real magic happens at the **intersections** — when a security insight improves a design decision, when a creativity technique solves a debugging problem, when a planning framework reveals an architecture flaw. This skill maps those connections and teaches you to find new ones.

> "The most exciting phrase in science is not 'Eureka!' but 'That's funny...'" — Isaac Asimov

The best ideas come from connecting things that don't obviously belong together.

---

## The Connection Taxonomy

Skills connect across seven types of relationships. When exploring connections, systematically check each type.

### Type 1: Structural Analogy

Two skills from different domains share the same underlying structure or pattern.

| Skill A | Skill B | Shared Structure |
|---------|---------|-----------------|
| `thinking-systems` (feedback loops) | `simota-bolt` (caching invalidation) | Both are about managing circular dependencies in complex systems |
| `thinking-first-principles` (decompose to fundamentals) | `pm-epic-breakdown` (decompose epics to stories) | Both use recursive decomposition to make the complex manageable |
| `thinking-red-team` (adversarial thinking) | `owasp-security` (threat modeling) | Both systematically think like an attacker to find weaknesses |
| `color-palette` (harmony from constraints) | `architecture-patterns` (elegance from constraints) | Both create beauty/quality by working within deliberate constraints |
| `wondelai-refactoring-ui` (visual hierarchy) | `sql-optimization` (query execution plan) | Both are about ordering operations for maximum efficiency of processing |
| `lyn-bayesian-calibration` (update beliefs with evidence) | `test-driven-development` (update code with failing tests) | Both use a tight feedback loop of hypothesis → evidence → correction |

**How to use:** When stuck in one domain, ask "What skill from a completely different domain has the same underlying structure?" Then apply that skill's techniques with domain-appropriate translation.

### Type 2: Inverse/Complement

Two skills that are mirror images — one creates, the other validates; one builds, the other breaks.

| Creator Skill | Validator Skill | Relationship |
|--------------|----------------|-------------|
| `frontend-design` | `ux-audit` | Design creates, audit validates |
| `architecture-patterns` | `jeffallan-the-fool` | Architecture proposes, The Fool stress-tests |
| `pm-prd-development` | `thinking-pre-mortem` | PRD plans for success, pre-mortem plans for failure |
| `creative-connections-engine` | `meta-thinking-engine` | Creativity generates, metacognition evaluates |
| `react-best-practices` | `code-review-excellence` | Implementation follows patterns, review checks adherence |
| `strategic-planning-engine` | `lyn-project-risk-register` | Planning charts the course, risk register maps the dangers |
| `auth-patterns` | `tob-codeql-analysis` | Auth implements security, CodeQL finds the holes |

**How to use:** Every creative/constructive skill should be paired with its inverse. Never build without validating. Never plan without stress-testing. This is the fundamental rhythm of quality work.

### Type 3: Prerequisite Chain

Skill A must be applied before Skill B can be effective.

| First | Then | Why |
|-------|------|-----|
| `pm-problem-statement` | `pm-prd-development` | You can't write requirements for a problem you haven't defined |
| `meta-thinking-engine` | `thinking-model-router` | Understand your thinking mode before selecting a mental model |
| `ui-ux-pro-max` | `frontend-design` | Design the experience before implementing the interface |
| `color-palette` | `design-system-creation` | Generate your palette before building the system that uses it |
| `owasp-security` | `semgrep-code-security` | Understand the threats before scanning for them |
| `pm-discovery-process` | `pm-roadmap-planning` | Validate the problem before planning the solution |
| `systematic-debugging` | `interactive-debugger` | Form a hypothesis before setting breakpoints |
| `thinking-first-principles` | `thinking-second-order` | Decompose to fundamentals, then trace the consequences |

**How to use:** Before applying any skill, ask "Is there a prerequisite skill that should fire first?" Working out of order wastes effort and produces inferior results.

### Type 4: Cross-Domain Transfer

A technique from one domain that directly applies in another, often with surprising effectiveness.

| Source Skill | Target Domain | Transfer |
|-------------|--------------|----------|
| `thinking-scientific-method` (hypothesis → experiment → conclusion) | Debugging | Every bug fix should be a hypothesis tested against evidence, not a random change |
| `wondelai-never-split-difference` (negotiation) | Code Review | "Mirroring" and "labeling" techniques make code review feedback more effective |
| `thinking-inversion` (what would make this fail?) | UI Design | Design for the error state first, then the happy path |
| `thinking-theory-of-constraints` (find the bottleneck) | Performance Optimization | Don't optimize everything — find and fix the single biggest bottleneck |
| `wondelai-design-everyday-things` (affordances) | API Design | APIs should have "affordances" — their interface should suggest correct usage |
| `thinking-ooda-loop` (Observe-Orient-Decide-Act) | Incident Response | Military decision cycle maps perfectly to production incident management |
| `thinking-margin-of-safety` (engineering buffers) | Sprint Planning | Always plan with 20-30% buffer — the margin of safety principle |
| `graphic-design` (Gestalt principles) | Information Architecture | Proximity, similarity, and closure apply to data organization, not just visual layout |
| `wondelai-hooked` (habit loops) | Developer Experience | Make your tools habit-forming: trigger → action → reward → investment |
| `thinking-cynefin` (complexity domains) | Project Management | Different project types need different management approaches — don't use agile for everything |

**How to use:** This is the most powerful connection type. Actively ask "What would an expert from a completely different field do with this problem?" Then find the skill from that field and apply its principles.

### Type 5: Emergent Combination

Two or more skills that, when combined, create a capability that neither has alone.

| Skills Combined | Emergent Capability |
|----------------|-------------------|
| `creative-connections-engine` + `thinking-model-router` + `meta-thinking-engine` | **Adaptive Creative Reasoning** — Generate creative ideas, route them through the right analytical framework, then metacognitively evaluate the quality of both the ideas and the reasoning |
| `ui-ux-pro-max` + `color-palette` + `wondelai-web-typography` + `interaction-design` | **Complete Design Language** — Not just a UI, but a coherent visual language with color theory, typographic rigor, and interaction patterns that feel inevitable |
| `owasp-security` + `privacy-by-design` + `auth-patterns` + `encryption-aes` | **Security-First Architecture** — Not security bolted on, but security as the foundation that shapes every architectural decision |
| `pm-discovery-process` + `pm-opportunity-tree` + `thinking-jobs-to-be-done` + `wondelai-mom-test` | **Evidence-Based Product Discovery** — Discover real problems, map the solution space, validate with users, all grounded in Jobs-to-be-Done theory |
| `strategic-planning-engine` + `thinking-systems` + `lyn-environmental-foresight` + `thinking-second-order` | **Systems-Aware Strategic Planning** — Plans that account for feedback loops, second-order effects, and environmental forces, not just linear task lists |
| `systematic-debugging` + `simota-scout` + `simota-specter` + `tob-codeql-analysis` + `simota-rewind` | **Deep Bug Forensics** — Not just finding bugs, but understanding their root cause, their history, their systemic patterns, and preventing their recurrence |

**How to use:** For complex tasks, don't just pick one skill — compose multiple skills into an emergent capability that's greater than the sum of its parts.

### Type 6: Tension/Dialectic

Two skills that pull in opposite directions, creating productive tension that leads to better outcomes.

| Skill A (Thesis) | Skill B (Antithesis) | Synthesis |
|------------------|---------------------|-----------|
| `thinking-first-principles` (simplify to essentials) | `thinking-systems` (everything is connected) | Know when to decompose and when to see the whole |
| `lyn-focus-timeboxing` (ship fast, 80/20) | `code-review-excellence` (thoroughness, quality) | Ship fast on the right things, be thorough on the critical things |
| `creative-connections-engine` (diverge, explore) | `meta-thinking-engine` (converge, evaluate) | Alternate between creative expansion and critical contraction |
| `thinking-via-negativa` (remove, simplify) | `pm-feature-investment` (add features for ROI) | Remove what doesn't matter so you can invest in what does |
| `simota-bolt` (optimize for speed) | `wondelai-clean-code` (optimize for readability) | Fast code that's also readable — the real craft |
| `thinking-effectuation` (work with what you have) | `thinking-first-principles` (reimagine from scratch) | Know when to be pragmatic and when to be revolutionary |

**How to use:** When two skills seem to conflict, don't pick one — hold both in tension. The synthesis is usually better than either extreme.

### Type 7: Recursive/Fractal

Skills that apply at multiple scales — the same skill works at the code level, the architecture level, the organization level, and the strategy level.

| Skill | Code Scale | Architecture Scale | Organization Scale | Strategy Scale |
|-------|-----------|-------------------|-------------------|---------------|
| `thinking-systems` | Function dependencies | Service dependencies | Team dependencies | Market dynamics |
| `thinking-theory-of-constraints` | Slow function | Bottleneck service | Overloaded team | Market constraint |
| `thinking-feedback-loops` | Recursive calls | Event loops | Retrospectives | Customer feedback |
| `thinking-first-principles` | Why this algorithm? | Why this architecture? | Why this org structure? | Why this market? |
| `thinking-inversion` | What breaks this function? | What breaks this system? | What breaks this team? | What breaks this strategy? |

**How to use:** When a skill works well at one scale, try applying it at a different scale. The same pattern often reveals insights at every level.

---

## The Connection Discovery Protocol

When facing any problem, use this protocol to find unexpected skill connections:

**Step 1: Name the problem domain.** What field is this problem obviously in?

**Step 2: Name three distant domains.** What fields have nothing to do with this problem?

**Step 3: For each distant domain, find one skill.** What skill from that domain might have a structural analogy to this problem?

**Step 4: Apply the distant skill's framework.** Translate the distant skill's key concepts into the language of your problem.

**Step 5: Evaluate the insight.** Did the cross-domain application reveal something you wouldn't have seen otherwise?

This protocol is based on Koestler's bisociation principle (from `creative-connections-engine`) — creative breakthroughs happen when two previously unconnected frames of reference collide.

---

## Living Map Principle

This connection map is not static. Every time you discover a new connection between skills, it strengthens the network. The more connections you find, the more powerful the entire skill library becomes. Each task is an opportunity to discover new connections that make future tasks easier.
