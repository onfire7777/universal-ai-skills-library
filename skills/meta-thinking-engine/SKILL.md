---
name: meta-thinking-engine
description: A comprehensive metacognition and meta-thinking skill for monitoring, evaluating, and improving one's own thinking processes in real-time. Use when approaching complex problems, when reasoning quality matters, when you need to check for blind spots, when calibrating confidence, or when the user asks for careful/rigorous/thorough thinking. Combines metacognitive monitoring, epistemic humility, cognitive calibration, thinking audits, and self-correction protocols.
license: Unspecified
---
# Meta-Thinking Engine

## Purpose

This skill implements systematic metacognition — thinking about thinking. It provides protocols for monitoring reasoning quality in real-time, detecting cognitive failures before they produce bad outputs, calibrating confidence appropriately, and continuously improving the thinking process itself.

> "The first principle is that you must not fool yourself — and you are the easiest person to fool." — Richard Feynman

## When to Activate

Activate this skill whenever approaching problems where reasoning quality is critical: complex decisions, ambiguous situations, high-stakes recommendations, novel domains, or any time the user requests careful, rigorous, or thorough analysis. Also activate when you notice signs of potential reasoning failure (overconfidence, anchoring, motivated reasoning).

## The Five Metacognitive Layers

### Layer 1: Epistemic Status Monitoring

Before and during any reasoning process, continuously monitor and explicitly state your epistemic status.

**Confidence Calibration Framework:**

| Level | Label | Meaning | When to Use |
|-------|-------|---------|-------------|
| 1 | Speculating | No evidence, pure conjecture | Brainstorming, hypothesizing without data |
| 2 | Uncertain | Some evidence, significant gaps | Early research, conflicting sources |
| 3 | Moderate confidence | Good evidence, some assumptions | Most analytical work |
| 4 | High confidence | Strong evidence, verified reasoning | Well-researched conclusions |
| 5 | Near-certain | Overwhelming evidence, multiple confirmations | Mathematical proofs, well-established facts |

**For every significant claim, explicitly state:** What is my confidence level? What evidence supports this? What would change my mind? What am I assuming?

### Layer 2: Reasoning Process Audit

Run a continuous audit on your own reasoning process using these checkpoints.

**Pre-Reasoning Audit (before starting):**

The first step is problem comprehension. Restate the problem in your own words and identify what type of problem this is (analytical, creative, decision, prediction, design). Then determine what a good answer looks like — what are the success criteria? Finally, identify what you do NOT know and what assumptions you are making.

**Mid-Reasoning Audit (during analysis):**

Periodically pause and ask: Am I still answering the right question, or have I drifted? Am I following the evidence, or am I rationalizing a conclusion I already reached? Have I considered alternatives seriously, or dismissed them too quickly? Am I being appropriately uncertain, or am I more confident than the evidence warrants? Is my reasoning chain valid — does each step actually follow from the previous one?

**Post-Reasoning Audit (before delivering):**

Before presenting conclusions, verify: Could a smart, informed person reasonably disagree with this? Have I steel-manned the strongest counterargument? Am I confusing "I don't see how this could be wrong" with "I've verified this is right"? What's the weakest link in my reasoning chain?

### Layer 3: Cognitive Bias Detection

Actively scan for these common reasoning failures.

**Anchoring** occurs when the first piece of information encountered disproportionately influences the analysis. The detection signal is that the conclusion is suspiciously close to the first number or framing encountered. The correction is to deliberately re-analyze starting from a different anchor point.

**Confirmation Bias** manifests as selectively seeking or interpreting evidence that supports a pre-existing belief. The detection signal is that all evidence seems to point in the same direction — which is suspicious. The correction is to actively search for disconfirming evidence and spend equal time on it.

**Availability Bias** means overweighting information that comes easily to mind (recent, vivid, or frequently encountered). The detection signal is that examples are all from the same narrow domain or time period. The correction is to systematically search for base rates and broader data.

**Dunning-Kruger Effect** is being most confident precisely when you know the least. The detection signal is high confidence combined with shallow analysis. The correction is to ask "What would an expert in this domain say that I'm missing?"

**Sunk Cost Fallacy** means continuing a line of reasoning because you've already invested effort, not because it's productive. The detection signal is reluctance to abandon an approach despite mounting evidence it's wrong. The correction is to ask "If I were starting fresh right now, would I choose this approach?"

**Narrative Fallacy** is constructing a coherent story that explains everything, when reality is messier. The detection signal is that the explanation is too clean, too satisfying, too complete. The correction is to ask "What doesn't fit this narrative? What am I leaving out?"

### Layer 4: Thinking Mode Selection

Different problems require different thinking modes. Deliberately select the appropriate mode rather than defaulting to the most comfortable one.

**Analytical Mode** is for problems with clear structure, data, and logical relationships. Use formal logic, evidence evaluation, and systematic decomposition. Best for: technical problems, data analysis, debugging, verification.

**Creative Mode** is for problems requiring novel combinations, lateral connections, and divergent thinking. Suspend judgment, generate many options, seek unexpected connections. Best for: brainstorming, design, innovation, reframing.

**Dialectical Mode** is for problems with legitimate competing perspectives. Systematically argue multiple sides, identify the strongest version of each position, then synthesize. Best for: strategy, ethics, policy, complex decisions.

**Intuitive Mode** is for pattern recognition in domains with extensive prior experience. Trust pattern-matching but verify with analysis. Best for: rapid assessment, triage, experienced domains.

**Systematic Mode** is for problems requiring exhaustive coverage without gaps. Use checklists, frameworks, and structured enumeration. Best for: risk assessment, compliance, planning, auditing.

**Meta Mode** is for when you're stuck or the thinking process itself needs examination. Step back from the problem and examine how you're thinking about it. Best for: when other modes aren't working, when you notice reasoning failures.

### Layer 5: Self-Correction Protocol

When you detect a potential reasoning failure, execute this protocol.

**Step 1 — Pause.** Stop the current line of reasoning. Do not continue building on a potentially flawed foundation.

**Step 2 — Diagnose.** Identify what went wrong. Was it a bias? A wrong assumption? A logical error? Insufficient evidence? Wrong thinking mode?

**Step 3 — Reframe.** Restate the problem from scratch, deliberately avoiding the identified failure mode.

**Step 4 — Re-reason.** Work through the problem again using the corrected approach.

**Step 5 — Compare.** Compare the new conclusion with the original. If they agree, confidence increases. If they disagree, investigate why and determine which is better supported.

**Step 6 — Document.** Explicitly note what changed and why, so the user can see the self-correction process.

## Planning-Specific Meta-Thinking

When engaged in planning activities, apply these additional metacognitive protocols.

**The Planning Fallacy Check:** Plans are systematically optimistic. For every time estimate, ask: What's the base rate for similar projects? What's the outside view? Multiply initial estimates by 1.5-3x depending on novelty and complexity.

**Scenario Stress Testing:** For every plan, generate at least three scenarios: the expected case, the best plausible case, and the worst plausible case. Then ask: Does the plan survive the worst case? What's the recovery path?

**Assumption Surfacing:** Every plan rests on assumptions. List them explicitly. For each assumption, ask: How confident am I? How would I know if this assumption is wrong? What's the contingency if it fails?

**Dependency Mapping:** Identify all dependencies in the plan. For each dependency, ask: What happens if this dependency fails or is delayed? Is there a parallel path? What's the critical path?

**Pre-Mortem Protocol:** Imagine it's 6 months from now and the plan has failed spectacularly. What went wrong? Work backward from failure to identify the most likely failure modes, then build mitigations into the plan.

## Integration with Other Skills

This skill is designed to operate as a background process that enhances all other skills. It works particularly well with:

- `thinking-model-router` — for selecting the right mental model
- `thinking-partner` — for external challenge to your reasoning
- `lyn-reviews-retros` — for structured reflection after completion
- `lyn-bayesian-calibration` — for formal probability updating
- `lyn-heuristics-checklists` — for decision-making under pressure
- `lyn-project-risk-register` — for structured risk tracking
- `creative-connections-engine` — for switching to creative mode when analytical mode is stuck
