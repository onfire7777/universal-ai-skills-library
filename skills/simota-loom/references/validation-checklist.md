# Validation Checklist

Purpose: Use this reference when Loom must score Figma Make output, decide a verdict, or write a validation report.

## Contents
- Validation model
- Weighted scorecard
- Required pass conditions
- Verdict bands
- Reporting format
- Guidelines quality tracking

## Validation Model

Validate against representative prompts, not a synthetic best case.

Always verify:
- token usage
- naming quality
- Auto Layout quality
- accessibility basics
- responsive behavior

## Weighted Scorecard

| Category | Weight | What to check |
|---|---:|---|
| `Token Usage` | `30%` | token adherence, no hardcoded values, semantic token use |
| `Naming` | `20%` | semantic layer names, no defaults, component/property naming |
| `Auto Layout` | `20%` | Auto Layout coverage, direction, gap, padding, sizing |
| `Accessibility` | `15%` | contrast, hierarchy, touch target basics, state clarity |
| `Responsive` | `15%` | resizing rules, width constraints, layout resilience |

## Required Pass Conditions

Fail immediately if any of the following is true:
- token usage is mostly hardcoded
- Auto Layout is missing across key frames
- default layer names remain widespread
- responsive behavior is absent on a responsive target
- accessibility issues block basic use

## Verdict Bands

| Score | Verdict | Meaning |
|---|---|---|
| `90-100` | `PASS` | ready to ship as Guidelines package |
| `70-89` | `CONDITIONAL` | usable after targeted fixes |
| `50-69` | `REVISE` | prompt or Guidelines logic must change |
| `<50` | `REBUILD` | change the strategy, not just the wording |

## Validation Report Format

```md
# Validation Report

## Scope
- Target file or screen
- Prompt set used

## Score
- Total: [0-100]
- Verdict: [PASS | CONDITIONAL | REVISE | REBUILD]

## Category Results
| Category | Score | Status | Key findings |
|---|---:|---|---|
| Token Usage | [score] | [PASS/FAIL] | [...] |
| Naming | [score] | [PASS/FAIL] | [...] |
| Auto Layout | [score] | [PASS/FAIL] | [...] |
| Accessibility | [score] | [PASS/FAIL] | [...] |
| Responsive | [score] | [PASS/FAIL] | [...] |

## Top Issues
1. [issue]
2. [issue]
3. [issue]

## Fix Plan
- Loom fixes
- Muse fixes
- Frame fixes
- Artisan fixes
```

## Issue Ownership

| Problem | Primary owner |
|---|---|
| Token misuse | Loom |
| Naming violation | Loom |
| Layout rule gap | Loom |
| Token drift | Muse |
| Structure mismatch | Frame |
| Production-pattern gap | Artisan |

## Refinement Rules

- Maximum `3` refinement cycles per Guidelines version.
- If the same class of failure survives `3` cycles, escalate with root-cause analysis.

## Guidelines Quality Tracking

Track these downstream metrics:

| Metric | Weight | Meaning |
|---|---:|---|
| `Implementation Fidelity` | `40%` | Make output vs production translation gap |
| `Validation Pass Rate` | `30%` | share of first-pass results at `CONDITIONAL` or better |
| `Token Coverage` | `20%` | percentage of intended token usage aligned with Figma Variables |
| `Refinement Cycles` | `10%` | fewer cycles is better |

Grade:
- `A`: `90+`
- `B`: `80+`
- `C`: `70+`
- `D`: `<70` and redesign recommended
