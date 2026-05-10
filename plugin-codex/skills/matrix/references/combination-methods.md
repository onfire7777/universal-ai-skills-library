# Combination Methods

Purpose: Use this file when you need the method definitions, reduction expectations, and constraint thresholds behind Matrix planning.

## Contents

- Full enumeration
- Pairwise
- Orthogonal arrays
- Higher-strength CIT
- Constraint handling

## Full Enumeration

Use the full Cartesian set when:

- axes `<= 2`
- the user explicitly requires exhaustive coverage
- the domain is small enough that optimization would hide important cases

Formula:

```text
total_combinations = ∏(value_count_per_axis)
```

Do not present full enumeration as the default once the matrix becomes expensive.

## Pairwise

Definition:

- Guarantee that every `2-way` value pair appears at least once.

Use pairwise when:

- axes `>= 3`
- invalid pairs stay below roughly `30%`
- the domain is not safety-critical

Typical reduction:

| Shape | Full | Pairwise | Reduction |
|---|---:|---:|---:|
| `3 axes x 3 values` | `27` | `9` | `67%` |
| `4 axes x 3 values` | `81` | `9-12` | `85-89%` |
| `10 axes x 3 values` | `59,049` | `18-27` | `99%+` |

Default evidence baseline:

- Pairwise usually detects `70-95%` of interaction faults in non-safety-critical software.

## Orthogonal Arrays

Use an orthogonal array when:

- each axis has the same number of values
- balanced representation matters as much as raw reduction
- a fixed, well-known array is easier than greedy generation

Common arrays:

| Array | Rows | Max factors | Levels |
|---|---:|---:|---:|
| `L4` | `4` | `3` | `2` |
| `L8` | `8` | `7` | `2` |
| `L9` | `9` | `4` | `3` |
| `L16` | `16` | `15` | `2` or `5 factors x 4 levels` |
| `L27` | `27` | `13` | `3` |

Choose OA over pairwise when value-balance is more important than the smallest possible row count.

## Higher-Strength CIT

Use `3-way+` coverage when:

- the system is safety-critical or regulated
- historical defects show higher-order interactions
- the user explicitly requests stronger assurance

Recommended baseline:

| Context | Minimum | Preferred |
|---|---|---|
| Normal application | `2-way` | `2-way` |
| High-quality or complex interaction domain | `2-way` | `3-way` |
| Safety-critical or regulated | `3-way` | `4-way+` |

Size guidance:

| Strength | Relative size vs pairwise | Typical use |
|---|---|---|
| `2-way` | baseline | general planning |
| `3-way` | `2-3x` | high-risk interaction zones |
| `4-way` | `3-5x` | regulated or safety-critical systems |

## Constraint Handling

Constraint types:

- `exclude`: impossible or invalid combinations
- `conditional`: one value forces another
- `require`: combinations that must appear in the final set

Constraint health:

| Exclusion rate | Interpretation | Action |
|---|---|---|
| `< 30%` | healthy | proceed |
| `30-40%` | suspicious | review model and business rules |
| `> 40%` | over-constrained | recommend redesign |

Escalate when:

- constraints remove every valid combination -> `ON_CONSTRAINT_UNKNOWN`
- constraints are too complex for simple pairwise reasoning -> mention SAT-solver style tooling

## Method Selection Summary

| Situation | Default choice |
|---|---|
| `<= 2` axes | Full |
| `3+` axes, normal risk | Pairwise |
| Uniform value counts, balanced representation needed | OA |
| Safety-critical or higher-order evidence | `3-way+` CIT |
| Heavy constraints plus budget | Constrained, budgeted optimization |
