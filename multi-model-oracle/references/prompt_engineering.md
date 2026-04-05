# Prompt Engineering Reference — Multi-Model Oracle

## Core Principles

The oracle's prompt engineering pipeline transforms raw user queries into optimized prompts tailored to each model's strengths. This is not generic "make it better" — it's model-aware, intent-aware optimization.

## Intent Detection

The system classifies queries into 7 categories, each triggering different optimization strategies:

| Intent | Optimization Strategy |
|---|---|
| code | Add language, framework, error handling, testing requirements |
| creative | Add tone, audience, style constraints, emotional targets |
| analysis | Add evaluation criteria, structure (pros/cons, matrix), depth |
| research | Add scope, level of detail, examples, analogies |
| reasoning | Add step-by-step requirement, verification, confidence levels |
| practical | Add specificity (OS, version, context), actionable format |
| general | Add structure, depth calibration, format preference |

## Model-Specific Tailoring

Each model receives a prompt variant optimized for its architecture:

**Anthropic Claude Opus 4.6**: Leverage its strength in nuanced multi-perspective analysis. Prompts emphasize considering edge cases, acknowledging uncertainty, and providing balanced viewpoints. Works best with explicit thinking structure.

**OpenAI GPT-5.4**: Leverage its strength in systematic coverage and technical precision. Prompts emphasize completeness, structured output, and concrete examples. Works best with clear format specifications.

**Manus gpt-4.1-mini**: Leverage its speed and efficiency. Prompts are more concise, emphasizing directness and actionability. Used as the prompt engineering engine itself (meta-level) and as the merge engine.

## Techniques Applied

| Technique | When Used | Effect |
|---|---|---|
| Role Priming | Always | Sets expert persona matching the domain |
| Specificity Injection | Vague queries | Adds concrete parameters (format, length, depth) |
| Structure Request | Analysis/research | Requests organized output sections |
| Chain of Thought | Reasoning tasks | Forces step-by-step verification |
| Negative Constraints | All queries | Eliminates filler, hedging, generic advice |
| Context Inference | Implicit requirements | Surfaces unstated but likely needs |

## Merge Strategy

The merge engine (Stage 3) follows a strict protocol:

1. Extract unique insights from each response
2. Resolve contradictions by majority or dual-perspective presentation
3. Eliminate redundancy without losing nuance
4. Structure the output as if from a single expert source
5. Never reference the individual models in the final output
