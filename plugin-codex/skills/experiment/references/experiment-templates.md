# Experiment Templates

## Hypothesis Document Template

```markdown
## Experiment: [Experiment Name]

### Hypothesis
**If** [we make this change]
**Then** [this metric will improve]
**Because** [this is the underlying mechanism]

### Background
- **Problem Statement:** [What problem are we solving?]
- **Current State:** [Current metric value and user behavior]
- **Evidence:** [What data/research supports this hypothesis?]

### Variants
| Variant | Description | Traffic Allocation |
|---------|-------------|--------------------|
| Control | Current experience | 50% |
| Treatment | [Describe change] | 50% |

### Metrics
**Primary Metric (決定指標):**
- Metric: [Name]
- Definition: [Exact calculation]
- Current Baseline: [X%]
- MDE (Minimum Detectable Effect): [Y%]
- Expected Lift: [Z%]

**Secondary Metrics (参考指標):**
1. [Metric name] - [Definition]
2. [Metric name] - [Definition]

**Guardrail Metrics (ガードレール指標):**
1. [Metric name] - [Threshold that should not be crossed]
2. [Metric name] - [Threshold]

### Sample Size & Duration
- Required Sample Size: [N per variant]
- Current Daily Traffic: [N users]
- Expected Duration: [X days/weeks]
- Statistical Power: 80%
- Significance Level: 5%

### Success Criteria
- [ ] Primary metric shows statistically significant improvement
- [ ] No guardrail metrics violated
- [ ] Lift >= MDE

### Rollout Plan
- **If wins:** Roll out to 100% on [date]
- **If loses:** Revert and [next action]
- **If inconclusive:** [Extend / iterate / abandon]
```

---

## Experiment Report Template

```markdown
## Experiment Report: [Experiment Name]

### Summary
| Metric | Control | Treatment | Lift | P-Value | Significant |
|--------|---------|-----------|------|---------|-------------|
| Primary: [Name] | X% | Y% | +Z% | 0.0XX | Yes/No |
| Secondary: [Name] | X | Y | +Z% | 0.0XX | Yes/No |
| Guardrail: [Name] | X | Y | -Z% | 0.0XX | No violation |

### Recommendation
**[SHIP / ITERATE / ABANDON]**

[1-2 sentences explaining the recommendation]

### Key Findings
1. [Finding 1 with data support]
2. [Finding 2 with data support]
3. [Finding 3 with data support]

### Detailed Results

#### Primary Metric: [Name]
- Control: [X%] (n=[N])
- Treatment: [Y%] (n=[N])
- Relative Lift: [+Z%]
- 95% CI: [[L%, U%]]
- P-Value: [0.0XX]
- Statistical Power Achieved: [X%]

#### Segment Analysis
| Segment | Control | Treatment | Lift | Significant |
|---------|---------|-----------|------|-------------|
| Mobile | X% | Y% | +Z% | Yes/No |
| Desktop | X% | Y% | +Z% | Yes/No |
| New Users | X% | Y% | +Z% | Yes/No |
| Returning Users | X% | Y% | +Z% | Yes/No |

### Timeline
- Started: [Date]
- Ended: [Date]
- Duration: [X days]
- Total Participants: [N]

### Learnings & Next Steps
1. [Learning 1] → [Next step]
2. [Learning 2] → [Next step]

### Appendix
- [Link to hypothesis document]
- [Link to raw data]
- [Link to dashboard]
```
