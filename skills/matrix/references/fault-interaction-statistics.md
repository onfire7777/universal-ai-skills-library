# Fault Interaction Statistics

Purpose: Use this file when deciding whether `2-way`, `3-way`, `4-way`, or mixed-strength coverage is justified.

## Contents

- Interaction distributions
- Strength selection
- Escalation strategy
- Mixed strength
- Coverage targets

## Interaction Distributions

NIST-style evidence summary:

| Domain | `<=2-way` | `<=3-way` | `<=4-way` | `<=5-way` | `<=6-way` |
|---|---:|---:|---:|---:|---:|
| Medical devices | `97%` | `99%` | `100%` | - | - |
| Browsers | `70%` | `90%` | `95%` | `99%` | `100%` |
| Servers | `76%` | `95%` | `99%` | `100%` | - |
| General software | `70-95%` | `90-99%` | `97-100%` | - | - |

Operational takeaway:

- `2-way` is enough for many normal systems.
- `3-way+` is justified once history, regulation, or criticality says pairwise is insufficient.

## Strength Selection

| Strength | Typical detection | Relative cost | Use when |
|---|---:|---:|---|
| `2-way` | `70-95%` | baseline | normal applications |
| `3-way` | `90-99%` | `2-3x` | high-quality or interaction-heavy areas |
| `4-way` | `97-100%` | `3-5x` | safety-critical or regulated systems |
| `5-way+` | `99-100%` | very high | exceptional regulatory cases |

## Escalation Strategy

Use this escalation path:

1. start with `2-way`
2. if additional higher-order defects appear, move to `3-way`
3. if `3-way` still exposes meaningful new defects, move to `4-way`
4. stop escalating when the next step yields no meaningful new defects

## Mixed Strength

Use mixed strength when only part of the model is high risk.

Example:

- authentication x privilege x data sensitivity -> `3-way`
- browser x OS -> `2-way`
- locale or theme -> `1-way` sampling

This often delivers better cost efficiency than applying `3-way` to the entire matrix.

## Coverage Targets

| Domain | Minimum | Preferred |
|---|---|---|
| General web application | `2-way 100%` | `2-way 100%` |
| Finance | `2-way 100%` | `3-way 100%` |
| Medical or safety-critical | `3-way 100%` | `4-way 100%` |
| IoT / embedded | `2-way 100%` | `3-way 100%` |
| Security testing | `2-way 100%` | mixed strength with `3-way` for high-risk areas |

Quality gates:

- `2-way` coverage below `100%` -> warning
- safety-critical domain + `2-way` only -> escalate recommendation
