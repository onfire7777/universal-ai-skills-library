# Statistical Methods for Experiments

## Test Selection Guide

| Metric Type | Test | When to Use |
|-------------|------|-------------|
| Binary (conversion) | Z-test for proportions | Standard A/B test |
| Continuous (revenue) | Welch's t-test | Revenue per user |
| Count (clicks) | Chi-square | Multiple categories |
| Time-to-event | Log-rank test | Time to conversion |

## Interpreting Results

| p-value | Conclusion |
|---------|------------|
| < 0.01 | Strong evidence |
| 0.01-0.05 | Moderate evidence |
| 0.05-0.10 | Weak evidence |
| > 0.10 | No significant evidence |

## Common Pitfalls
1. **Peeking**: Checking results before reaching sample size inflates false positives
2. **Multiple comparisons**: Testing many metrics without correction
3. **Selection bias**: Non-random assignment to variants
4. **Novelty effect**: Short tests may capture excitement, not sustained behavior

---

## Z-Test for Proportions (Binary Metrics)

```typescript
interface ExperimentResult {
  control: { conversions: number; total: number };
  treatment: { conversions: number; total: number };
}

interface AnalysisResult {
  controlRate: number;
  treatmentRate: number;
  relativeLift: number;
  absoluteLift: number;
  zScore: number;
  pValue: number;
  isSignificant: boolean;
  confidenceInterval: [number, number];
}

function analyzeExperiment(
  result: ExperimentResult,
  significance: number = 0.05
): AnalysisResult {
  const { control, treatment } = result;

  // Conversion rates
  const p1 = control.conversions / control.total;
  const p2 = treatment.conversions / treatment.total;

  // Pooled proportion
  const pPooled = (control.conversions + treatment.conversions) /
                  (control.total + treatment.total);

  // Standard error
  const se = Math.sqrt(
    pPooled * (1 - pPooled) * (1/control.total + 1/treatment.total)
  );

  // Z-score
  const zScore = (p2 - p1) / se;

  // P-value (two-tailed)
  const pValue = 2 * (1 - normalCDF(Math.abs(zScore)));

  // 95% Confidence interval for the difference
  const zAlpha = 1.96;
  const seDiff = Math.sqrt(
    p1 * (1 - p1) / control.total +
    p2 * (1 - p2) / treatment.total
  );
  const ci: [number, number] = [
    (p2 - p1) - zAlpha * seDiff,
    (p2 - p1) + zAlpha * seDiff
  ];

  return {
    controlRate: p1,
    treatmentRate: p2,
    relativeLift: (p2 - p1) / p1,
    absoluteLift: p2 - p1,
    zScore,
    pValue,
    isSignificant: pValue < significance,
    confidenceInterval: ci
  };
}

// Normal CDF approximation
function normalCDF(x: number): number {
  const a1 =  0.254829592;
  const a2 = -0.284496736;
  const a3 =  1.421413741;
  const a4 = -1.453152027;
  const a5 =  1.061405429;
  const p  =  0.3275911;

  const sign = x < 0 ? -1 : 1;
  x = Math.abs(x) / Math.sqrt(2);

  const t = 1.0 / (1.0 + p * x);
  const y = 1.0 - (((((a5 * t + a4) * t) + a3) * t + a2) * t + a1) * t * Math.exp(-x * x);

  return 0.5 * (1.0 + sign * y);
}
```

## Example Analysis

```typescript
const result = analyzeExperiment({
  control: { conversions: 500, total: 10000 },      // 5.0%
  treatment: { conversions: 550, total: 10000 }     // 5.5%
});

console.log(`
Control Rate: ${(result.controlRate * 100).toFixed(2)}%
Treatment Rate: ${(result.treatmentRate * 100).toFixed(2)}%
Relative Lift: ${(result.relativeLift * 100).toFixed(1)}%
P-Value: ${result.pValue.toFixed(4)}
Significant: ${result.isSignificant ? 'Yes' : 'No'}
95% CI: [${(result.confidenceInterval[0] * 100).toFixed(2)}%, ${(result.confidenceInterval[1] * 100).toFixed(2)}%]
`);
```
