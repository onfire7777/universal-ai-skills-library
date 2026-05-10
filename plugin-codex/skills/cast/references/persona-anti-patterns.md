# Persona Anti-Patterns

Purpose: Define the common failure modes Cast should detect or avoid during creation, maintenance, and organizational rollout.

## Contents

1. Core anti-patterns
2. Persona fatigue
3. Anti-persona usage
4. Common failure factors

## Core Anti-Patterns

| ID | Name | What goes wrong | Mitigation |
|---|---|---|---|
| `PA-01` | Demographics Fixation | Persona is mostly age/gender/job labels | Anchor on goals, pain points, and behaviors |
| `PA-02` | Single Monolithic Persona | One persona tries to represent everyone | Keep at least `P0/P1/P2` by default |
| `PA-03` | Happy Path Persona | Only ideal users are represented | Include friction-heavy or underserved users |
| `PA-04` | Proto-Persona Ossification | Hypotheses are treated as stable truth | Keep validation status explicit |
| `PA-05` | User-Buyer Conflation | Buyer and end user are merged | Split if goals or behaviors differ materially |
| `PA-06` | One-Shot Creation | Persona is created once and never updated | Use `AUDIT` and `EVOLVE` regularly |
| `PA-07` | Over-Designed Artifact | Persona looks polished but is weakly evidenced | Favor evidence density over visual polish |
| `PA-08` | Specificity Imbalance | Too vague or too fictional | Keep roughly 80% evidence / 20% inference |
| `PA-09` | Silo Creation | Persona is not shared or reusable | Register and distribute systematically |
| `PA-10` | Gallery Display | Persona exists as decoration only | Tie personas to downstream agent tasks |

## Persona Fatigue

### What causes it

- too many personas
- stale personas
- personas not used in real decisions
- overly repetitive or decorative artifacts

### Mitigation

- keep persona count manageable
- deprecate or archive stale personas
- track actual downstream use
- distribute task-specific versions, not full archives

## Anti-Persona

Use anti-personas to define who the product should not optimize for.

Recommended steps:

1. Identify clearly mismatched segments.
2. Document why they are out-of-scope.
3. Record the associated cost/risk.
4. Keep anti-personas separate from primary personas.
5. Revisit during major strategy shifts.

## Common Failure Factors

Common failure patterns to watch:

- weak stakeholder buy-in
- poor evidence quality
- no update cadence
- persona sprawl
- lack of integration with decision workflows
