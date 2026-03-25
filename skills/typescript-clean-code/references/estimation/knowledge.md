# Estimation Knowledge

Core concepts and foundational understanding for software estimation.

## Overview

Estimation is one of the most misunderstood activities in software development. The fundamental problem is that business and developers view estimates differently - business sees commitments, developers see guesses. Understanding this distinction is critical for professional software development.

## Key Concepts

### Estimate vs Commitment

**Definition**: An estimate is a guess with no commitment implied; a commitment is a promise that must be achieved.

This is the most important distinction in professional estimation. Confusing these two concepts leads to broken promises, damaged reputations, and failed projects.

**Key points**:
- Estimates have no promises attached - missing an estimate is not dishonorable
- Commitments require certainty - professionals only commit when they know they can deliver
- Business needs both, but they serve different purposes

### An Estimate is a Distribution

**Definition**: An estimate is not a single number but a probability distribution describing the range of possible completion times.

When someone says "three days," they're giving you the most likely duration, not a guarantee. The actual completion time could fall anywhere on a probability curve, with a tail extending toward longer durations.

**Key points**:
- The "estimate" typically represents the peak of probability (most likely)
- There's always a tail extending toward longer times
- Professional estimates communicate this uncertainty explicitly

### Implied Commitments

**Definition**: Implicit promises created through language that sounds like agreement, especially the word "try."

When you agree to "try" to meet a deadline, you've made a commitment. There's no other interpretation - agreeing to try is agreeing to succeed.

**Key points**:
- "Can you try?" = "Will you commit?"
- Saying yes to "try" implies working extra hours, weekends, skipping vacations
- Professionals carefully avoid implied commitments

### PERT (Program Evaluation and Review Technique)

**Definition**: A technique using three estimates (optimistic, nominal, pessimistic) to create probability distributions suitable for planning.

Created in 1957 for the U.S. Navy's Polaris submarine project, PERT converts estimates into realistic probability distributions.

**Key points**:
- Uses trivariate analysis: O (optimistic), N (nominal), P (pessimistic)
- Expected duration: (O + 4N + P) / 6
- Standard deviation: (P - O) / 6

### Law of Large Numbers

**Definition**: Breaking large tasks into smaller ones and estimating independently produces more accurate total estimates because errors tend to integrate out.

**Key points**:
- Small task errors tend to cancel each other
- Underestimation bias means integration isn't perfect
- Breaking tasks up also helps discover hidden complexity

## Terminology

| Term | Definition |
|------|------------|
| Commitment | A promise you must achieve, regardless of required effort |
| Estimate | A guess about duration with no promise implied |
| Trivariate Analysis | Using three estimates (O, N, P) to model uncertainty |
| Expected Duration (mu) | The statistically expected completion time |
| Standard Deviation (sigma) | A measure of uncertainty in the estimate |

## How It Relates To

- **Saying No**: Declining impossible commitments while providing honest estimates
- **Time Management**: Realistic estimates enable better planning
- **Professionalism**: Clear communication about certainty builds trust

## Common Misconceptions

- **Myth**: A good developer should be able to estimate accurately
  **Reality**: Estimates are inherently uncertain - that's why they're called estimates

- **Myth**: Agreeing to "try" is different from committing
  **Reality**: Agreeing to try IS committing - there's no other interpretation

- **Myth**: The nominal estimate is when you'll be done
  **Reality**: The nominal is just the most likely point on a distribution curve

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Estimate | A probability distribution, not a single number |
| Commitment | A promise requiring certainty before making |
| PERT Formula | (O + 4N + P) / 6 for expected duration |
| Standard Deviation | (P - O) / 6 measures uncertainty |
| "Try" | Agreeing to try = agreeing to succeed |
