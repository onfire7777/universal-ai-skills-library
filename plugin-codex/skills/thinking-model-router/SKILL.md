---
name: thinking-model-router
description: Route to the right mental model based on your domain and problem type. The single entry point for all thinking skills.
license: Unspecified
---
# Model Router

## Overview

This is the **master routing skill** for all mental models. Instead of knowing 38 frameworks, start here. Identify your domain and problem type, and this skill points you to the right model(s). Think of it as the "which tool do I use?" guide.

**Core Principle:** Don't memorize models—memorize how to find the right one. Domain + Problem Type → Model.

## Quick Router

### Step 1: What's Your Domain?

| Domain | You're working on... |
|--------|---------------------|
| **Coding/Debugging** | Bugs, errors, performance issues, root cause |
| **Architecture** | System design, technical decisions, scalability |
| **Product** | Features, user needs, prioritization, roadmap |
| **Business Strategy** | Competition, growth, market, organization |
| **Personal Decisions** | Career, life choices, major commitments |
| **Abstract/Analytical** | Arguments, ideas, theories, pure reasoning |
| **Risk/Safety** | What could go wrong, preparation, resilience |
| **Innovation** | New ideas, breakthroughs, creative solutions |

### Step 2: What's Your Problem Type?

| Type | You need to... |
|------|----------------|
| **Diagnose** | Find root cause, understand why |
| **Decide** | Choose between options |
| **Understand** | Grasp how something works |
| **Create** | Generate new solutions |
| **Evaluate** | Judge quality or validity |
| **Predict** | Forecast outcomes |
| **Optimize** | Improve performance |

---

## Domain → Model Maps

### 🖥️ Coding & Debugging

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
Bug with unknown cause           → Scientific Method, 5 Whys Plus
Performance degradation          → Theory of Constraints, Systems Thinking
Spans multiple services          → Systems Thinking, Feedback Loops
Incident postmortem              → 5 Whys Plus, Systems Thinking
Flaky/intermittent behavior      → Scientific Method (hypothesis testing)
"It works on my machine"         → Map-Territory (model vs reality gap)
```

**Default for debugging:** Start with **5 Whys Plus**, escalate to **Systems Thinking** if it spans components.

---

### 🏗️ Architecture & Technical Decisions

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
Technology choice                → Lindy Effect, Reversibility
Build vs buy                     → Opportunity Cost, First Principles
Scalability design               → Systems Thinking, Leverage Points
Microservices vs monolith        → Cynefin, Reversibility
Database selection               → Lindy Effect, Theory of Constraints
API design tradeoffs             → TRIZ (resolve contradictions)
Should we rewrite?               → Second-Order, Opportunity Cost
```

**Default for architecture:** Start with **Reversibility** (is this Type 1 or Type 2?), then **Systems Thinking** for interconnections.

---

### 📦 Product & Feature Development

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
What should we build?            → Jobs to be Done
Feature prioritization           → Opportunity Cost, Theory of Constraints
Why aren't users engaging?       → Jobs to be Done, 5 Whys Plus
New product exploration          → Effectuation, First Principles
Should we pivot?                 → Regret Minimization, Reversibility
Product-market fit               → Jobs to be Done, Scientific Method
Roadmap planning                 → Theory of Constraints, Opportunity Cost
A/B test design                  → Scientific Method, Bayesian
```

**Default for product:** Start with **Jobs to be Done** (what job is the user hiring this for?).

---

### 📈 Business Strategy

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
Competitive analysis             → Red Team, Second-Order
Market entry                     → Cynefin, Effectuation
Growth strategy                  → Feedback Loops, Leverage Points
Organizational dysfunction       → Archetypes, Systems Thinking
Resource allocation              → Theory of Constraints, Opportunity Cost
Startup strategy                 → Effectuation, Margin of Safety
M&A evaluation                   → Pre-mortem, Steel-manning
Pricing decisions                → First Principles, Fermi Estimation
```

**Default for strategy:** Start with **Cynefin** (what domain is this problem in?), then match approach.

---

### 🧑 Personal & Career Decisions

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
Should I take this job?          → Regret Minimization, Reversibility
Career direction                 → Circle of Competence, Regret Minimization
Major life decision              → Regret Minimization, Pre-mortem
Learning what to learn           → Circle of Competence, Lindy Effect
Negotiation prep                 → Steel-manning, Red Team
Should I start a company?        → Effectuation, Margin of Safety, Pre-mortem
Time allocation                  → Opportunity Cost, Theory of Constraints
```

**Default for personal:** Start with **Regret Minimization** (what will 80-year-old you think?).

---

### 🧠 Abstract & Analytical Thinking

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
Evaluating an argument           → Steel-manning, Bayesian
Challenging assumptions          → First Principles, Socratic
Estimating unknowns              → Fermi Estimation, Probabilistic
Updating beliefs                 → Bayesian, Probabilistic
Exploring edge cases             → Thought Experiment, Inversion
Finding logical flaws            → Inversion, Steel-manning
Complex causation                → Systems Thinking, Feedback Loops
Philosophical questions          → Thought Experiment, First Principles
```

**Default for abstract:** Start with **Steel-manning** (argue the strongest opposing view first).

---

### ⚠️ Risk & Safety

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
What could go wrong?             → Pre-mortem, Red Team
Security review                  → Red Team, Inversion
Disaster preparation             → Pre-mortem, Margin of Safety
Avoiding catastrophic failure    → Margin of Safety, Via Negativa
Stress-testing plans             → Red Team, Pre-mortem
Probability of failure           → Probabilistic, Bayesian
Building resilience              → Via Negativa, Margin of Safety
```

**Default for risk:** Start with **Pre-mortem** (assume failure, explain why).

---

### 💡 Innovation & Creativity

```
PROBLEM                          → MODEL(S)
─────────────────────────────────────────────────────
Breakthrough needed              → First Principles, TRIZ
Stuck on contradictions          → TRIZ
Limited resources                → Effectuation, Via Negativa
Simplification                   → Via Negativa, Occam's Razor
Challenging "impossible"         → First Principles, TRIZ
New market creation              → Effectuation, Jobs to be Done
Removing complexity              → Via Negativa, Occam's Razor
```

**Default for innovation:** Start with **First Principles** (strip to fundamentals, rebuild).

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────┐
│                    MENTAL MODEL QUICK ROUTER                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  "Why is this broken?"        → 5 Whys Plus, Scientific Method  │
│  "How does this system work?" → Systems Thinking, Feedback Loops│
│  "What should we build?"      → Jobs to be Done                 │
│  "Should I do this?"          → Reversibility, Regret Min       │
│  "What could go wrong?"       → Pre-mortem, Red Team            │
│  "How do I innovate?"         → First Principles, TRIZ          │
│  "What's the probability?"    → Bayesian, Probabilistic         │
│  "Where's the bottleneck?"    → Theory of Constraints           │
│  "What am I giving up?"       → Opportunity Cost                │
│  "Is this argument valid?"    → Steel-manning                   │
│  "Will this technology last?" → Lindy Effect                    │
│  "How complex is this?"       → Cynefin                         │
│  "What to remove?"            → Via Negativa                    │
│  "Is this safe enough?"       → Margin of Safety                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Decision Flow

```
START HERE
    │
    ▼
┌─────────────────────┐
│ What's your domain? │
└──────────┬──────────┘
           │
     ┌─────┴─────┬──────────┬──────────┬──────────┐
     ▼           ▼          ▼          ▼          ▼
  Coding    Architecture  Product   Strategy   Personal
     │           │          │          │          │
     ▼           ▼          ▼          ▼          ▼
┌─────────────────────┐
│ What problem type?  │
│ Diagnose/Decide/    │
│ Understand/Create/  │
│ Evaluate/Predict    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Look up in domain   │
│ table above         │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Single model enough?│
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    ▼             ▼
   YES           NO
    │             │
    ▼             ▼
 Apply it    Use Model Combination
             (Sequential/Parallel/Nested)
```

## Model Inventory by Category

### Diagnostic Models (Find root cause)
- `5 Whys Plus` - Iterative "why" with bias guards
- `Scientific Method` - Hypothesis → Test → Learn
- `Kepner-Tregoe` - Systematic problem/decision analysis

### Decision Models (Choose wisely)
- `Reversibility` - Type 1 vs Type 2 decisions
- `Regret Minimization` - What will future-you think?
- `Opportunity Cost` - What are you giving up?
- `Bayesian` - Update beliefs with evidence
- `Probabilistic` - Calibrated probability estimates

### Systems Models (Understand interconnections)
- `Systems Thinking` - Feedback, emergence, non-linearity
- `Feedback Loops` - Reinforcing and balancing loops
- `Theory of Constraints` - Find and exploit the bottleneck
- `Leverage Points` - Where small changes have big effects
- `Archetypes` - Recurring system patterns

### Risk Models (Prepare for failure)
- `Pre-mortem` - Assume failure, explain why
- `Red Team` - Attack your own plan
- `Margin of Safety` - Build in buffers
- `Inversion` - Identify paths to failure, avoid them

### Innovation Models (Create breakthroughs)
- `First Principles` - Strip to fundamentals, rebuild
- `TRIZ` - Resolve technical contradictions
- `Effectuation` - Start with means, not goals
- `Via Negativa` - Improve by removing

### Evaluation Models (Judge quality)
- `Steel-manning` - Argue strongest opposing view
- `Lindy Effect` - Older = likely to last longer
- `Circle of Competence` - Know your expertise boundaries
- `Occam's Razor` - Prefer simpler explanations

### Context Models (Match approach to situation)
- `Cynefin` - Clear/Complicated/Complex/Chaotic domains
- `Model Selection` - Choose the right model
- `Model Combination` - Use multiple models together

### Product Models (Build the right thing)
- `Jobs to be Done` - What job is user hiring this for?
- `Thought Experiment` - Structured imagination

### Estimation Models (Size unknowns)
- `Fermi Estimation` - Order-of-magnitude calculations

## When to Combine Models

Use **Model Combination** when:

| Situation | Combination Pattern | Example |
|-----------|--------------------| --------|
| High-stakes decision | Sequential | Reversibility → Pre-mortem → Opportunity Cost |
| System diagnosis | Nested | Cynefin (macro) → ToC (meso) → OODA (micro) |
| Validating strategy | Parallel | Red Team + Steel-manning + Second-Order |
| Innovation under constraints | Sequential | First Principles → TRIZ → Effectuation |
| Career decision | Temporal | 5 Whys (past) → Circle of Competence (present) → Regret Min (future) |

## Template

```markdown
# Model Router Analysis

## Context
Domain: [Coding/Architecture/Product/Strategy/Personal/Abstract/Risk/Innovation]
Problem: [Brief description]
Problem Type: [Diagnose/Decide/Understand/Create/Evaluate/Predict/Optimize]

## Routed Models
Primary: [Main model to use]
Secondary: [If needed]
Combination pattern: [Sequential/Parallel/Nested/None]

## Application
[Apply the selected model(s) here]

## Verification
- [ ] Domain correctly identified
- [ ] Problem type matches
- [ ] Model fits the situation
- [ ] Considered if combination needed
```

## Key Questions

- "What domain am I operating in?"
- "What type of problem is this—diagnose, decide, understand, create, or evaluate?"
- "What's the default model for this domain + type?"
- "Does one model cover it, or do I need to combine?"
- "Am I using a model because it fits, or because it's familiar?"

---

**Remember:** You don't need to know all 38 models. You need to know how to find the right one. Start with domain, identify problem type, look it up, apply. That's it.
