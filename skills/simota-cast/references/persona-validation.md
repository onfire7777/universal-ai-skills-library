# Persona Validation Methods

Purpose: Define how Cast validates personas with triangulation, survey evidence, clustering, and staged confidence upgrades.

## Contents

1. Why validation matters
2. Triangulation patterns
3. Survey thresholds
4. Clustering guidance
5. Synthetic persona rules
6. Validation statuses

## Why Validation Matters

Validation exists to avoid:

- proto-personas treated as facts
- confirmation bias from creators
- over-generalizing from too little data
- stale personas surviving market or product change

## Triangulation Patterns

| Pattern | Methods | Strength |
|---|---|---|
| Basic | interview `5-10` + survey `350+` | Cost-efficient |
| Behavioral | interview + behavior logs + experiment / test evidence | Verifies saying vs doing |
| Full | interview + survey + behavior logs + usability evidence | Highest confidence |

## Quantitative Survey Thresholds

- `350+` respondents per segment for `95%` confidence
- `1000+` total when comparing multiple segments
- Prefer behavior-based questions over preference-only questions
- Likert + free text is the default hybrid

## ML Clustering Guidance

### Algorithm Fit

| Algorithm | Use when | Caveat |
|---|---|---|
| K-means | clearly separated segments | requires preselected cluster count |
| DBSCAN | irregular clusters and outliers | parameter-sensitive |
| Hierarchical clustering | exploratory structure analysis | weak for very large datasets |
| Gaussian mixture | overlapping segments | higher computational cost |

### Cluster Count Rules

- Use `Elbow + Silhouette + Gap` together.
- `Silhouette > 0.5` is a good signal.
- Recommended persona count is `3-7`.
  - early product: `3-4`
  - mature product: `5-7`

## Validation Workflow

1. Collect behavior, survey, support, and satisfaction data.
2. Preprocess and normalize.
3. Cluster with more than one method when possible.
4. Match clusters against current personas.
5. Treat uncovered clusters as new persona candidates.
6. Raise confidence only after evidence-backed mapping.

## Synthetic Persona Rules

- Synthetic personas are hypothesis tools, not production truth.
- Use them to improve guides, expose gaps, or explore edge cases.
- Never treat them as substitutes for real user validation.
- Keep synthetic and real-data-backed personas explicitly separated.

## Validation Statuses

| Status | Meaning |
|---|---|
| `proto` | hypothesis only |
| `partial` | validated by one stream only |
| `validated` | triangulated |
| `ml_validated` | supported by clustering evidence |

### Confidence Contributions

| Validation state | Contribution |
|---|---|
| Proto baseline | `0.30` |
| Interview validation | `+0.20` |
| Survey validation | `+0.15` |
| ML validation | `+0.20` |
| Triangulation complete | `+0.10` |
