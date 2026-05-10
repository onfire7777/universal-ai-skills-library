---
name: strategic-planning-engine
description: 'A comprehensive planning skill covering project planning, strategic planning, sprint planning, scenario planning, contingency planning, roadmapping, resource allocation, milestone tracking, risk management, and execution frameworks. Use when the user needs to plan anything — from a single feature to a multi-year strategy. Combines the best planning frameworks: OKRs, WBS, Gantt, critical path, agile ceremonies, scenario planning, wardley mapping, and more.'
license: Unspecified
---
# Strategic Planning Engine

## Purpose

This skill provides a complete toolkit for planning at every level of abstraction — from 5-minute task planning to multi-year strategic roadmaps. It ensures plans are realistic, risk-aware, stakeholder-aligned, and executable.

## Planning Hierarchy

Planning operates at five distinct levels. Each level requires different tools, timeframes, and thinking modes. The key insight is that higher-level plans constrain lower-level plans, but lower-level execution informs higher-level strategy.

### Level 1: Vision and Mission (3-10 year horizon)

This level answers "Why do we exist and where are we going?" The key frameworks at this level are North Star Metric (the single metric that best captures the value you deliver), Vision Statement (a vivid description of the future state you're creating), and Mission Statement (how you will achieve the vision and for whom).

The output at this level is a one-page document that any team member can recite from memory.

### Level 2: Strategy (1-3 year horizon)

This level answers "How will we win?" The key frameworks include Wardley Mapping (map the value chain, identify where to invest based on evolution stage), Porter's Five Forces (competitive intensity, supplier power, buyer power, substitutes, new entrants), Blue Ocean Strategy (ERRC grid — Eliminate, Reduce, Raise, Create), and OKRs (Objectives and Key Results — 3-5 objectives with 2-4 measurable key results each).

The output at this level is a strategy document with clear choices about what you will and will NOT do.

### Level 3: Roadmap (Quarter to year horizon)

This level answers "What will we build and when?" The key frameworks are Now/Next/Later Roadmap (avoids false precision of date-based roadmaps), Opportunity Solution Tree (connect business outcomes to opportunities to solutions to experiments), Theme-Based Roadmap (organize by strategic themes, not features), and Dependency Mapping (identify cross-team dependencies and sequence accordingly).

The output at this level is a visual roadmap that stakeholders can understand and teams can execute against.

### Level 4: Sprint/Iteration Planning (1-4 week horizon)

This level answers "What are we doing this cycle?" The key frameworks are Sprint Planning (select stories from backlog, estimate, commit), Kanban Flow (limit WIP, optimize throughput, visualize bottlenecks), and Capacity Planning (available hours minus meetings, minus overhead, equals actual capacity — typically 60-70% of calendar time).

The output at this level is a sprint backlog with clear acceptance criteria and capacity-matched commitments.

### Level 5: Task Planning (Hours to days)

This level answers "What do I do right now?" The key frameworks are Eisenhower Matrix (urgent/important 2x2), Pomodoro Technique (25-min focused blocks), Getting Things Done (capture, clarify, organize, reflect, engage), and Time-Boxing (allocate fixed time to tasks, move on when time expires).

The output at this level is a prioritized task list with time estimates and clear next actions.

## Core Planning Frameworks

### Work Breakdown Structure (WBS)

Decompose any project into a hierarchy of deliverables. The process begins by identifying the final deliverable, then breaking it into 3-7 major components. Each component is decomposed until work packages are small enough to estimate (typically 1-5 days of effort). Each work package must have a clear owner, definition of done, and estimated effort.

### Critical Path Method (CPM)

Identify the longest sequence of dependent tasks that determines the minimum project duration. List all tasks with durations and dependencies, then calculate the earliest start and finish for each task (forward pass). Next calculate the latest start and finish (backward pass). Tasks with zero float are on the critical path — any delay here delays the entire project.

### PERT Estimation

For each task, estimate three durations: Optimistic (O) — everything goes perfectly; Most Likely (M) — the realistic estimate; and Pessimistic (P) — everything that can go wrong does. The expected duration equals (O + 4M + P) / 6. Standard deviation equals (P - O) / 6. This gives you both an estimate and a confidence interval.

### Scenario Planning

Develop 3-4 plausible future scenarios and test your plan against each. Identify the 2-3 most impactful uncertainties, then create a 2x2 matrix of the top two uncertainties. Name each quadrant as a scenario and develop a narrative for each. Test your current plan against all four scenarios. Build flexibility into the plan to handle multiple scenarios.

### Contingency Planning

For each identified risk, develop a contingency plan. The structure includes a Trigger (what observable event tells you the risk has materialized), a Response (what specific actions to take), a Owner (who is responsible for executing the response), and Resources (what budget, time, or people are pre-allocated for this contingency).

### Wardley Mapping

Map the value chain from user need to underlying components. Position each component on the evolution axis (Genesis, Custom, Product, Commodity). Strategy emerges from the movement patterns — invest in components moving from Custom to Product, commoditize what's already Product, and explore what's in Genesis.

## Planning Anti-Patterns

**The Waterfall Trap** means planning everything upfront with false precision. Instead, plan in waves — detailed for the near term, directional for the medium term, aspirational for the long term.

**The Estimate Fallacy** is treating estimates as commitments. Instead, communicate estimates as ranges with confidence levels. "2-4 weeks, 80% confidence" is more honest than "3 weeks."

**The Happy Path Plan** only accounts for the scenario where everything goes right. Instead, explicitly plan for the top 3-5 risks and build buffer for unknown unknowns (typically 20-30% of total effort).

**The Kitchen Sink Plan** tries to do everything. Instead, ruthlessly prioritize. A plan that tries to do 10 things will accomplish none well. A plan that focuses on 3 things has a chance.

**The Orphan Plan** is created and then ignored. Instead, build review cadences into the plan itself — weekly check-ins, monthly reviews, quarterly strategy sessions.

## Planning Session Protocol

When asked to help plan anything, follow this protocol.

**Phase 1 — Understand.** Clarify the objective, constraints, stakeholders, and success criteria. Ask: What does success look like? What are the hard constraints? Who needs to be aligned?

**Phase 2 — Assess.** Evaluate the current state, available resources, known risks, and dependencies. Ask: Where are we starting from? What do we have to work with? What could go wrong?

**Phase 3 — Design.** Create the plan structure using the appropriate frameworks for the planning level. Include milestones, owners, estimates (as ranges), dependencies, and risk mitigations.

**Phase 4 — Stress-Test.** Run the plan through scenario analysis, pre-mortem, and dependency analysis. Ask: What's the weakest link? What happens if our biggest assumption is wrong?

**Phase 5 — Communicate.** Present the plan at the right level of abstraction for the audience. Executives get the one-page summary. Teams get the detailed breakdown. Include decision points and review cadences.

## Integration with Other Skills

This skill works best when combined with:

- `meta-thinking-engine` — for checking planning assumptions and biases
- `lyn-project-risk-register` — for structured risk tracking
- `lyn-environmental-foresight` — for PESTLE scanning before strategic planning
- `pm-roadmap-planning` — for detailed roadmap construction
- `pm-prioritization` — for framework-based prioritization
- `lyn-stakeholders-org` — for stakeholder mapping and alignment
- `thinking-pre-mortem` — for failure mode identification
- `thinking-second-order` — for anticipating consequences of plan decisions
