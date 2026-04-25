# Common Pitfalls & Solutions

## 1. Peeking Problem

**Problem:** Checking results repeatedly increases false positive rate.

**Solution:** Use sequential testing with alpha spending.

```typescript
function shouldContinue(
  currentPValue: number,
  currentN: number,
  targetN: number,
  alphaSpent: number = 0
): { continue: boolean; decision: 'significant' | 'not_significant' | 'continue' } {
  const progress = currentN / targetN;
  const alphaAvailable = 0.05 * progress - alphaSpent;

  if (currentPValue < alphaAvailable) {
    return { continue: false, decision: 'significant' };
  }

  if (progress >= 1) {
    return { continue: false, decision: currentPValue < 0.05 ? 'significant' : 'not_significant' };
  }

  return { continue: true, decision: 'continue' };
}
```

## 2. Multiple Comparisons

**Problem:** Testing many metrics inflates false positive rate.

**Solution:** Bonferroni correction or pre-register ONE primary metric.

```typescript
function adjustedSignificance(
  numComparisons: number,
  baseAlpha: number = 0.05
): number {
  return baseAlpha / numComparisons;
}

// 5 metrics → significance level = 0.01 per test
```

## 3. Selection Bias

**Problem:** Non-random assignment leads to confounded results.

**Solution:** Use deterministic hashing for assignment.

```typescript
// Good: Deterministic based on user ID
const variant = getVariant('experiment', userId);

// Bad: Random each time
const variant = Math.random() > 0.5 ? 'treatment' : 'control';
```
