# AI and LLM Persona Generation Risks

Purpose: Define AI-specific bias, ethics, validation, and confidence guardrails for persona generation.

## Contents

1. Risk categories
2. Bias effects
3. Hybrid workflow
4. Guardrails
5. Ethics
6. Cast integration

## Risk Categories

| Category | Risk |
|---|---|
| Data quality | biased training data, weak representativeness, stale snapshot effects |
| Bias | gender, age, cultural, socioeconomic stereotyping |
| Ethics | weak consent, privacy leakage, deepfake personas, weak explainability |
| Methodology | over-homogenization, weak segment discovery, amplified confirmation bias, weak validation |
| Organization | over-automation, skill erosion, cost illusion, unclear accountability |

## Bias Effects

Typical failure modes:

- stereotyped gender-role assumptions
- oversimplified older-user representations
- English- or majority-culture bias in global products
- persuasive but false personas produced by prompt injection or weak review

## Recommended Hybrid Approach

| Phase | Owner | Output |
|---|---|---|
| AI hypothesis generation | AI | draft proto-persona |
| Human deepening | Human reviewer | corrected and contextualized persona |
| Data validation | Research / analytics | validated persona |
| Continuous monitoring | AI + fresh data | drift alerts and update suggestions |

## Mandatory Guardrails

| Guardrail | Rule |
|---|---|
| Bias check | Review for gender, age, and cultural stereotyping |
| Source transparency | Record where each important attribute came from |
| Human review | AI-generated personas always require review |
| Confidence cap | AI-only generation is capped at `0.50` |
| Diversity check | Evaluate persona-set diversity, not just individual quality |
| Version tracking | Record AI-generated vs human-updated history |

## Ethics Principles

- Transparency: record model, prompt, and sources when AI materially shaped output.
- Fairness: audit for bias regularly.
- Human-centricity: real user evidence remains the authority.
- Minimality: use only the minimum sensitive data needed.
- Accountability: humans own the decision, not the generated persona.

## Cast Integration

Recommended generation-method labels:

- `manual`
- `ai_assisted`
- `ai_generated`
- `data_driven`

Promotion path:

- `ai_generated` or `proto`
- `ai_assisted` after human review
- `validated` after data-backed validation
