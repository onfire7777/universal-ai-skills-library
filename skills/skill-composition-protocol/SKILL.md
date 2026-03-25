---
name: skill-composition-protocol
description: >
  The definitive guide to chaining, composing, and layering multiple skills together for complex
  multi-phase tasks. Use when a task requires more than 2-3 skills, when you need to determine
  the optimal sequence of skill application, when skills need to pass outputs to each other,
  or when you need to build a custom workflow from existing skills. This is the "programming
  language" for skill orchestration — it turns individual skills into composable pipelines.
---

# Skill Composition Protocol

## Purpose

Individual skills are functions. This skill teaches you to compose those functions into programs — multi-step workflows where the output of one skill becomes the input of the next, where skills run in parallel when independent, and where the whole pipeline produces results no single skill could achieve alone.

---

## Composition Patterns

### Pattern 1: Sequential Pipeline

Skills execute one after another, each building on the previous output.

**Structure:** Skill A → Output A → Skill B → Output B → Skill C → Final Output

**When to use:** When each skill needs the output of the previous one. Most common pattern.

**Example — Building a Feature:**

```
pm-problem-statement     → Clear problem definition
  ↓
pm-prd-development       → Requirements document
  ↓
architecture-patterns    → Technical design
  ↓
frontend-design          → UI/UX approach
  ↓
react-best-practices     → Implementation
  ↓
test-driven-development  → Tests + verified code
  ↓
code-review-excellence   → Reviewed, polished code
```

**Rules for Sequential Pipelines:**
- Each skill's output must be concrete enough to serve as the next skill's input
- If a skill's output is ambiguous, add a checkpoint to clarify before proceeding
- Document what passes between each stage so the chain is traceable

### Pattern 2: Parallel Fan-Out

Multiple skills execute simultaneously on the same input, then results are merged.

**Structure:**
```
                ┌→ Skill A → Output A ─┐
Input → Fork ──┼→ Skill B → Output B ──┼→ Merge → Combined Output
                └→ Skill C → Output C ─┘
```

**When to use:** When multiple independent analyses are needed on the same artifact.

**Example — Security + Performance + Accessibility Review:**

```
                        ┌→ owasp-security        → Security findings    ─┐
Completed feature → ────┼→ tlc-perf-lighthouse    → Performance findings ─┼→ Unified review report
                        ├→ accessibility-testing  → A11y findings        ─┤
                        └→ code-review-excellence → Code quality findings─┘
```

**Rules for Parallel Fan-Out:**
- Only parallelize skills that are truly independent (no output dependencies)
- Define a merge strategy before fanning out (union, intersection, priority-weighted)
- Resolve conflicts between parallel outputs using the Priority Stack from `master-skill-orchestrator`

### Pattern 3: Iterative Refinement Loop

A skill is applied repeatedly, with each iteration improving the output.

**Structure:**
```
Input → Skill A → Output → Evaluate → Good enough? 
                                ↓ No
                          Skill A (again, with feedback) → Output → Evaluate → ...
```

**When to use:** When quality emerges through iteration, not single-pass execution.

**Example — Design Iteration:**

```
Requirements → ui-ux-pro-max → Initial design
  ↓
ux-audit → Findings
  ↓
ui-ux-pro-max (with findings) → Improved design
  ↓
ux-audit → Fewer findings
  ↓
(repeat until audit passes)
```

**Rules for Iterative Refinement:**
- Set a maximum iteration count (typically 3) to prevent infinite loops
- Each iteration must address specific, named issues from the previous evaluation
- If quality plateaus after 2 iterations, the problem is likely upstream — go back to an earlier pipeline stage

### Pattern 4: Conditional Branching

Different skills activate based on the characteristics of the input or intermediate results.

**Structure:**
```
Input → Classify → If Type A → Skill A → Output
                   If Type B → Skill B → Output
                   If Type C → Skill C → Output
```

**When to use:** When the right approach depends on what you discover during analysis.

**Example — Bug Triage:**

```
Bug report → systematic-debugging → Classify bug type
  ↓
  ├─ If race condition    → simota-specter → Fix
  ├─ If regression        → simota-rewind → Find introducing commit → Fix
  ├─ If performance       → simota-bolt → Profile → Optimize
  ├─ If security vuln     → owasp-security → Assess severity → Patch
  └─ If unknown root cause → simota-scout → Deep investigation → Fix
```

**Rules for Conditional Branching:**
- Classification must be explicit and documented
- Always have a default/fallback branch for unexpected cases
- Re-evaluate the classification if the chosen branch doesn't produce results

### Pattern 5: Wrapper/Decorator

One skill wraps around another, adding a cross-cutting concern without changing the core skill's behavior.

**Structure:**
```
meta-thinking-engine ─────────────────────────────────────┐
│                                                          │
│  Input → Core Skill → Output                            │
│                                                          │
│  (metacognition monitors reasoning quality throughout)   │
└──────────────────────────────────────────────────────────┘
```

**When to use:** When a concern (security, quality, metacognition) should be active throughout, not just at one stage.

**Permanent Wrappers (always active):**

| Wrapper Skill | What It Monitors | When It Intervenes |
|--------------|-----------------|-------------------|
| `meta-thinking-engine` | Reasoning quality, bias, confidence calibration | When it detects overconfidence, anchoring, or logical errors |
| `big-picture-thinking-engine` | Alignment with bigger goals, long-term implications | When tactical work drifts from strategic direction |
| `master-skill-orchestrator` | Whether the right skills are being used | When it detects a missed skill or wrong skill application |

**Situational Wrappers (active when relevant):**

| Wrapper Skill | Active When | What It Monitors |
|--------------|------------|-----------------|
| `owasp-security` | Any code that handles user data, auth, or external input | Security vulnerabilities |
| `ux-audit` | Any UI work | Usability issues |
| `code-review-excellence` | Any code production | Code quality |

### Pattern 6: Ensemble/Voting

Multiple skills independently analyze the same question, then their conclusions are compared and synthesized.

**Structure:**
```
Question → Skill A → Conclusion A ─┐
Question → Skill B → Conclusion B ──┼→ Compare → Synthesize → Best conclusion
Question → Skill C → Conclusion C ─┘
```

**When to use:** For high-stakes decisions where you want multiple independent perspectives.

**Example — Architecture Decision:**

```
"Should we use microservices or monolith?"
  ↓
  ├→ thinking-first-principles  → Conclusion based on fundamentals
  ├→ thinking-cynefin            → Conclusion based on complexity domain
  ├→ thinking-lindy-effect       → Conclusion based on what's proven
  ├→ thinking-reversibility      → Conclusion based on reversibility
  └→ thinking-second-order       → Conclusion based on long-term effects
  ↓
  Compare all 5 conclusions → Where do they agree? Where do they disagree?
  ↓
  Synthesize → Decision with high confidence (if agreement) or identified uncertainty (if disagreement)
```

**Rules for Ensemble:**
- Use at least 3 independent skills for meaningful ensemble results
- Weight conclusions by relevance (a security skill's opinion on security matters more than a creativity skill's)
- Disagreement between skills is a signal, not a problem — it reveals genuine uncertainty

---

## Composition Depth Levels

Match composition complexity to task complexity.

| Task Complexity | Composition Depth | Example |
|----------------|------------------|---------|
| **Trivial** | Single skill, no composition | Fix a typo → `code-review-excellence` |
| **Simple** | 2-skill sequential | Add a button → `frontend-design` → `react-best-practices` |
| **Moderate** | 3-5 skill pipeline | Build a form → `ui-ux-pro-max` → `frontend-design` → `react-best-practices` → `test-driven-development` |
| **Complex** | Pipeline + parallel + wrapper | Full feature → Sequential pipeline with parallel review fan-out and metacognition wrapper |
| **Wicked** | All patterns, multi-phase | System redesign → Conditional branching into multiple pipelines with iterative refinement and ensemble decisions |

---

## The Composition Checklist

Before executing a composed skill pipeline, verify:

1. **Completeness:** Are all necessary skills included? Check `master-skill-orchestrator` for commonly missed skills.
2. **Order:** Are skills in the right sequence? Check `skill-connection-map` for prerequisite chains.
3. **Independence:** Are parallel skills truly independent? If not, sequentialize them.
4. **Termination:** Do iterative loops have exit conditions? Set maximum iterations.
5. **Conflict Resolution:** If parallel skills disagree, how will you resolve it? Define the strategy upfront.
6. **Wrappers:** Are the right cross-cutting concerns active? At minimum: `meta-thinking-engine` for reasoning quality.
7. **Checkpoints:** Are there intermediate checkpoints where you can evaluate progress and re-route if needed?
8. **Traceability:** Can you trace the reasoning chain from input to output through every skill? If not, add documentation between stages.

---

## Integration

This skill is the "how" that complements the other meta-skills:

| Meta-Skill | Role |
|-----------|------|
| `master-skill-orchestrator` | **Which** skills to use (selection) |
| `skill-connection-map` | **Why** these skills connect (relationships) |
| `big-picture-thinking-engine` | **Whether** this is the right thing to be doing (alignment) |
| `skill-composition-protocol` | **How** to combine skills into effective workflows (execution) |
| `meta-thinking-engine` | **Quality** of the reasoning throughout (monitoring) |
