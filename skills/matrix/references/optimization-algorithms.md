# Optimization Algorithms

Purpose: Use this file when you must justify the chosen optimization method, budget tradeoffs, or coverage verification logic.

## Contents

- Method selection flow
- Pairwise generation
- Orthogonal array application
- Constraint-aware optimization
- Budgeted optimization
- Prioritization
- Coverage verification
- Performance guidance

## Method Selection Flow

| Condition | Default choice |
|---|---|
| `max_combinations` exists | budgeted optimization |
| invalid pairs `> 30%` | constrained pairwise |
| uniform value counts | orthogonal array |
| safety-critical requirement | `3-way+` CIT |
| axes `<= 2` | full enumeration |
| otherwise | pairwise |

## Pairwise Generation

Implementation expectation:

- expand the uncovered `2-way` tuple set
- greedily add rows that cover the most remaining tuples
- stop only when `2-way` coverage reaches `100%`

Use pairwise when:

- axes are numerous
- the domain is normal-risk
- constraints are manageable

## Orthogonal Arrays

Use OA when:

- value counts are uniform
- balanced distribution matters
- a known array such as `L9` or `L16` fits the matrix

Reference defaults:

| Levels | Factors | Array |
|---|---:|---|
| `2` | `2-3` | `L4` |
| `2` | `4-7` | `L8` |
| `3` | `2-4` | `L9` |
| `2` | `8-15` | `L16` |
| `4` | `4-5` | `L16` |

## Constraint-Aware Optimization

Apply constraints in this order:

1. remove invalid combinations
2. optimize the remaining space
3. enforce conditional rules
4. inject required cases

Escalate when:

- the valid space becomes empty
- constraints exceed the model’s explanatory power

## Budgeted Optimization

When `max_combinations` is smaller than the optimized set:

- keep the highest-coverage rows first
- report the achieved coverage rate
- list missing tuples explicitly
- do not claim the original coverage guarantee if the budget breaks it

## Prioritization

Rank after optimization, not before.

Useful strategies:

- weighted axis priority
- value-specific risk scores
- percentile bucketing for `Critical / High / Medium / Low`

Distribution guardrails:

- keep `Critical` around the top `10-20%`
- keep `Critical + High` around `<= 30%`
- warn when all rows collapse into one priority bucket

## Coverage Verification

Minimum verification fields:

- total tuples
- covered tuples
- coverage rate
- missing tuples

For pairwise:

```text
coverage_rate = covered_2_way_tuples / total_2_way_tuples
```

For higher-strength plans, verify the selected `t-way` explicitly.

## Performance Guidance

| Scale | Recommended approach | Typical runtime target |
|---|---|---|
| `2-5 axes`, `2-5 values` | Pairwise | `< 1s` |
| `6-10 axes`, `2-5 values` | Pairwise | `1-10s` |
| `10-15 axes`, `2-10 values` | Pairwise or OA | `10-60s` |
| very large or heavily constrained | specialized external tooling | context-dependent |
