---
name: challenge-that
description: "Force critical evaluation of any proposal, decision, or approach from 5 adversarial perspectives. Use when you've just agreed to something, accepted a plan, or want to stress-test an idea before committing. Triggers on: 'challenge this', 'stress test', 'devil's advocate', 'what could go wrong', 'poke holes', 'is this a good idea', 'evaluate this critically'."
---

# Challenge That

Force critical evaluation by analyzing from multiple adversarial perspectives.

## When Invoked

1. Identify what was just accepted/proposed in this conversation
2. State it clearly: "**Challenging:** [the proposal/decision]"
3. Analyze from each of the five perspectives below
4. Synthesize into actionable concerns
5. Ask: "Which concerns should we address before proceeding?"

## The Five Hats

| Hat | Focus | Key Questions |
|-----|-------|---------------|
| 🔴 **Skeptic** | Evidence & proof | "What evidence supports this? Has this been validated? Are we guessing?" |
| 🟡 **Pragmatist** | Cost/benefit | "Is this the simplest fix? What's the maintenance burden? Is it worth it?" |
| 🟢 **Edge Case Hunter** | Failure modes | "What breaks this? What's the worst case? What did we miss?" |
| 🔵 **Structural Critic** | Architecture | "Is this the right location? Does it fit the existing design? Will it cause problems elsewhere?" |
| 🟣 **Root Cause Analyst** | Problem diagnosis | "Is the problem correctly identified? Are we treating symptoms? What's the actual cause?" |

## Output Format

```
## Challenging: [proposal/decision being challenged]

### 🔴 Skeptic
[Challenge from evidence/proof perspective]

### 🟡 Pragmatist
[Challenge from cost/benefit perspective]

### 🟢 Edge Case Hunter
[Challenge from failure modes perspective]

### 🔵 Structural Critic
[Challenge from architecture perspective]

### 🟣 Root Cause Analyst
[Challenge from problem diagnosis perspective]

---

**Key Concerns:**
1. [Most significant concern]
2. [Second concern]
3. [Third concern]

Which of these should we address before proceeding?
```

## Rules

🚨 **Be genuinely adversarial.** Don't softball the challenges. If a perspective finds nothing wrong, say so—but look hard first.
🚨 **Challenge the proposal, not the person.** Focus on the idea's weaknesses, not who suggested it.
🚨 **Provide actionable output.** Each challenge should point to something that could be investigated or changed.