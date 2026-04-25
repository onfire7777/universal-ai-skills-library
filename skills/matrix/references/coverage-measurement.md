# Coverage Measurement

Purpose: Use this file when validating coverage claims or mapping execution results back into missing tuples and follow-up work.

## Contents

- Coverage layers
- Warning signs
- Distribution checks
- Post-execution analysis
- Recovery targets
- Escape-rate guidance

## Coverage Layers

| Layer | Formula | Use |
|---|---|---|
| Pair coverage | covered `t-way` tuples / total `t-way` tuples | optimization guarantee |
| Configuration coverage | represented space / full Cartesian space | space compression view |
| Fault-detection effectiveness | detected defects / total relevant defects | post-execution learning |

## Warning Signs

| ID | Signal | Interpretation | Action |
|---|---|---|---|
| `CM-01` | production bug in a supposedly covered tuple | value granularity is too coarse | split values or add axes |
| `CM-02` | production bug in an uncovered tuple | a parameter or constraint is wrong | expand the model |
| `CM-03` | some values appear too rarely | algorithm or constraints are biased | inspect per-value distribution |
| `CM-04` | many extra defects appear in `3-way` | `2-way` was too weak | escalate strength |
| `CM-05` | value distribution is heavily skewed | weak balance | warn when deviation exceeds `30%` |
| `CM-06` | exclusion count rises over time | model drift | review constraints |

## Distribution Checks

Per-value appearance check:

- compute expected appearances = total cases / value count
- warn when deviation exceeds `30%`
- consider `require` or OA if one value is starved

## Post-Execution Analysis

Always do these steps when results exist:

1. collect `Pass / Fail / Skip`
2. identify tuples invalidated by failures or skips
3. propose the smallest supplemental set that restores the target coverage
4. update the model if higher-order interactions are discovered

## Recovery Targets

| Domain risk | Target |
|---|---|
| Critical | `100%` recovery |
| Standard | `80%` recovery |
| Low-risk | `50%` recovery focused on critical tuples |

## Escape Rate

Formula:

```text
escape_rate = production defects that should have been caught / total relevant defects
```

Interpretation:

| Rate | Meaning | Action |
|---:|---|---|
| `< 5%` | healthy | keep the model |
| `5-15%` | needs work | refine parameters or increase strength |
| `> 15%` | poor | redesign the model from parameter identification upward |
