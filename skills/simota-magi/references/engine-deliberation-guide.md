# Engine Deliberation Guide

Engine Mode specification for Magi's three-engine deliberation system. In Engine Mode, Claude, Codex, and Gemini each provide an independent, integrated analysis rather than simulating three internal perspectives.

---

## Engine Availability Check

Before starting Engine Mode deliberation, verify external engine availability:

```bash
which codex && echo "codex: available" || echo "codex: not found"
which gemini && echo "gemini: available" || echo "gemini: not found"
```

| Engine | Command | Role |
|--------|---------|------|
| **Claude** | (internal) | Primary deliberator, orchestrator |
| **Codex** | `codex exec --full-auto "{prompt}"` | Independent external deliberator |
| **Gemini** | `gemini -p "{prompt}" --yolo` | Independent external deliberator |

---

## Deliberation Prompt Template (Codex)

Codex is invoked via `codex exec --full-auto` with a concise, directive prompt optimized for structured output.

**Important:** This is a text analysis task, not code generation. The prompt must explicitly request analytical output.

### Prompt Structure

```text
codex exec --full-auto "You are a decision analyst. Evaluate the following decision independently.

DECISION: {subject}
TYPE: {decision_type}
CONTEXT: {context_summary}
OPTIONS: {options}
CONSTRAINTS: {constraints}

Provide your analysis as YAML in a code block:

\`\`\`yaml
position: APPROVE | REJECT | ABSTAIN
confidence: 0-100
rationale: "2-3 sentence explanation of your position"
key_evidence:
  - "Evidence point 1"
  - "Evidence point 2"
risks:
  - "Risk 1"
  - "Risk 2"
conditions:
  - "Condition for approval (if any)"
dissent_note: "Key concern if this decision goes the other way"
\`\`\`

Be direct. State your position clearly. Do not hedge."
```

### Codex Optimization Notes

- Keep prompts concise and directive (Codex performs best with clear instructions)
- Avoid verbose background sections
- Request YAML output explicitly
- Use `--full-auto` to avoid interactive prompts

---

## Deliberation Prompt Template (Gemini)

Gemini is invoked via `gemini -p` with `--yolo` flag for non-interactive execution. Gemini benefits from additional context due to its larger context window.

### Prompt Structure

```text
gemini -p "You are a decision analyst. Evaluate the following decision independently.

## Background
{extended_context}

## Decision
DECISION: {subject}
TYPE: {decision_type}
OPTIONS: {options}
CONSTRAINTS: {constraints}

## Instructions
Provide your analysis as YAML in a code block:

\`\`\`yaml
position: APPROVE | REJECT | ABSTAIN
confidence: 0-100
rationale: "2-3 sentence explanation of your position"
key_evidence:
  - "Evidence point 1"
  - "Evidence point 2"
risks:
  - "Risk 1"
  - "Risk 2"
conditions:
  - "Condition for approval (if any)"
dissent_note: "Key concern if this decision goes the other way"
\`\`\`

Be direct. State your position clearly. Do not hedge." --yolo
```

### Gemini Optimization Notes

- Include a `## Background` section with extended context (leverages larger context window)
- Same YAML output format as Codex for consistent parsing
- `--yolo` flag enables non-interactive execution

---

## Claude Internal Deliberation (Engine Mode)

In Engine Mode, Claude does **not** simulate three perspectives (that is Simple Mode). Instead, Claude provides a single, integrated analysis as one unified viewpoint.

### Claude's Engine Mode Behavior

1. Analyze the decision from a **unified perspective** (combining technical, human, and strategic considerations)
2. Output in the **same YAML format** as Codex and Gemini
3. Complete analysis **before** collecting external engine outputs (contamination prevention)

### Claude Output Format

```yaml
position: APPROVE | REJECT | ABSTAIN
confidence: 0-100
rationale: "2-3 sentence integrated analysis"
key_evidence:
  - "Evidence point 1"
  - "Evidence point 2"
risks:
  - "Risk 1"
  - "Risk 2"
conditions:
  - "Condition for approval (if any)"
dissent_note: "Key concern if this decision goes the other way"
```

---

## Output Parsing Strategy

Parse engine outputs using a progressive fallback approach:

### Stage 1: YAML Block Extraction

Extract content from ` ```yaml ... ``` ` fenced code blocks in the engine output.

### Stage 2: Key Validation

Required keys (parse fails without these):
- `position` (required)
- `confidence` (required)
- `rationale` (required)

Optional keys:
- `key_evidence`
- `risks`
- `conditions`
- `dissent_note`

### Stage 3: Value Validation

| Field | Valid Values | Default |
|-------|------------|---------|
| `position` | APPROVE, REJECT, ABSTAIN | ABSTAIN |
| `confidence` | Integer 0-100 | 0 |
| `rationale` | Non-empty string | "No rationale provided" |

### Stage 4: Fallback Text Analysis

If YAML parsing fails, scan the full output text for:

| Keyword Pattern | Inferred Position |
|----------------|-------------------|
| "approve", "recommend", "proceed", "yes" | APPROVE |
| "reject", "against", "deny", "no" | REJECT |
| "abstain", "uncertain", "insufficient" | ABSTAIN |

Confidence is estimated from language strength:
- Strong language ("clearly", "strongly", "definitely") → 70
- Moderate language ("likely", "probably", "reasonable") → 50
- Weak language ("possibly", "might", "uncertain") → 30

### Stage 5: Complete Failure

If all parsing stages fail:
```yaml
position: ABSTAIN
confidence: 0
rationale: "Engine output could not be parsed"
```

---

## Fallback Strategy

Engine Mode adapts based on available engines:

| Available Engines | Mode | Behavior |
|---|---|---|
| 3 (Claude + Codex + Gemini) | **Full Engine Mode** | 3 engines deliberate independently |
| 2 (Claude + 1 external) | **2-Engine Mode** | 2 engines, consensus patterns: 2-0 / 1-1 / 0-2 |
| 1 (Claude only) | **Auto-fallback** | Automatic switch to Simple Mode, notify user |

### 2-Engine Mode Details

When only one external engine is available:
- Claude + available engine deliberate independently
- Consensus patterns reduce to: 2-0 (unanimous), 1-1 (split), 0-2 (unanimous rejection)
- Weighted confidence is 2-engine average
- Split (1-1) always escalates to user

### Auto-Fallback Notification

When falling back to Simple Mode, inform the user:

```
⚠ Engine Mode requested but external engines unavailable.
  Codex: [available/not found]
  Gemini: [available/not found]
→ Falling back to Simple Mode (Logos/Pathos/Sophia internal deliberation).
```

---

## Error Handling

| Error | Response |
|---|---|
| **Timeout** (engine does not respond within 60s) | Retry once → if retry fails, treat as ABSTAIN |
| **Parse failure** (output not parseable) | Apply fallback text analysis → if complete failure, ABSTAIN |
| **API/CLI error** (non-zero exit code) | Treat as ABSTAIN, record error in risk register |
| **Both external engines fail** | Fall back to Simple Mode with notification |

### Error Recording Format

```yaml
engine_error:
  engine: "[Codex | Gemini]"
  error_type: "[timeout | parse_failure | cli_error]"
  detail: "[Error message or description]"
  action_taken: "[ABSTAIN | retry_succeeded | fallback_to_simple]"
```

---

## Engine Mode Execution Flow

```
1. AVAILABILITY CHECK
   └─ Verify codex/gemini CLI availability
   └─ Determine mode (Full / 2-Engine / Auto-fallback)

2. CLAUDE ANALYSIS (FIRST — contamination prevention)
   └─ Claude completes integrated analysis
   └─ Output stored before external calls

3. EXTERNAL ENGINE CALLS (PARALLEL)
   └─ codex exec --full-auto "{prompt}"
   └─ gemini -p "{prompt}" --yolo

4. OUTPUT PARSING
   └─ Parse each engine's YAML output
   └─ Apply fallback parsing if needed

5. VOTE ASSEMBLY
   └─ Combine 3 (or 2) engine positions
   └─ Calculate weighted confidence
   └─ Determine consensus pattern

6. SYNTHESIS & DELIVERY
   └─ Present Engine Mode MAGI display
   └─ Include risk register and next steps
```
