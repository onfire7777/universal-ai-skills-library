# Estimation Workflow

PERT-based task estimation process for accurate estimates.

## When to Use

- Sprint planning
- Project scoping
- When asked "how long will this take?"
- Committing to deadlines

## Prerequisites

- Understanding of the task/feature
- Knowledge of the codebase
- Ability to break down work

**Reference**: `estimation/rules.md`, `estimation/examples.md`, `commitment/rules.md`

---

## The Golden Rules

1. **Estimates are NOT commitments** - They're probability distributions
2. **Never give a single number** - Always provide a range
3. **Break it down** - Large tasks have large uncertainty
4. **Communicate uncertainty** - Be explicit about what you don't know

---

## Workflow Steps

### Step 1: Clarify the Request

**Goal**: Understand exactly what's being asked.

- [ ] What is the scope of the work?
- [ ] What are the acceptance criteria?
- [ ] What's NOT included?
- [ ] Are there dependencies on others?

**Ask**:
- "Can you clarify what 'done' means for this?"
- "Does this include testing/documentation/deployment?"
- "Are there any constraints I should know about?"

**Reference**: `acceptance-testing/knowledge.md`

---

### Step 2: Break Down the Task

**Goal**: Decompose into estimable pieces.

Break the task into subtasks that are:
- Independent (can be done separately)
- Small (< 1 day ideal, < 3 days max)
- Clear (you know what "done" means)

**Example**:
```
Feature: User Export

Subtasks:
1. Create UserExporter interface
2. Implement JSON formatter
3. Implement CSV formatter  
4. Add error handling
5. Write unit tests
6. Write integration tests
7. Update API documentation
```

- [ ] Each subtask is small enough to estimate
- [ ] No subtask > 3 days
- [ ] Dependencies identified

---

### Step 3: Apply PERT to Each Subtask

**Goal**: Get trivariate estimates for each subtask.

For each subtask, estimate three values:

| Value | Symbol | Meaning |
|-------|--------|---------|
| Optimistic | O | Best case (everything goes right) |
| Nominal | N | Most likely (normal conditions) |
| Pessimistic | P | Worst case (things go wrong) |

**Example**:
```
Subtask: Implement JSON formatter
- Optimistic (O): 2 hours (straightforward, no surprises)
- Nominal (N): 4 hours (normal development)
- Pessimistic (P): 12 hours (unexpected complexity, bugs)
```

**Tips for each value**:
- **Optimistic**: "If everything goes perfectly"
- **Nominal**: "Most likely scenario"
- **Pessimistic**: "If I hit problems but still finish"

---

### Step 4: Calculate PERT Values

**Goal**: Compute expected time and uncertainty.

**PERT Formulas**:
```
Expected Time (μ) = (O + 4N + P) / 6
Standard Deviation (σ) = (P - O) / 6
```

**Example Calculation**:
```
O = 2, N = 4, P = 12

μ = (2 + 4×4 + 12) / 6 = (2 + 16 + 12) / 6 = 30/6 = 5 hours
σ = (12 - 2) / 6 = 10/6 = 1.67 hours
```

**Calculate for each subtask**:
```
| Subtask              | O  | N  | P   | μ     | σ    |
|---------------------|----|----|-----|-------|------|
| UserExporter interface| 1  | 2  | 4   | 2.2   | 0.5  |
| JSON formatter       | 2  | 4  | 12  | 5.0   | 1.7  |
| CSV formatter        | 2  | 4  | 10  | 4.7   | 1.3  |
| Error handling       | 1  | 2  | 6   | 2.5   | 0.8  |
| Unit tests           | 2  | 4  | 8   | 4.3   | 1.0  |
| Integration tests    | 2  | 4  | 10  | 4.7   | 1.3  |
| Documentation        | 1  | 2  | 4   | 2.2   | 0.5  |
```

---

### Step 5: Sum and Calculate Total

**Goal**: Combine subtask estimates.

**Total Expected Time**:
```
μ_total = Σ μ_i  (sum of all expected times)
```

**Total Standard Deviation**:
```
σ_total = √(Σ σ_i²)  (square root of sum of squares)
```

**Example**:
```
μ_total = 2.2 + 5.0 + 4.7 + 2.5 + 4.3 + 4.7 + 2.2 = 25.6 hours

σ_total = √(0.5² + 1.7² + 1.3² + 0.8² + 1.0² + 1.3² + 0.5²)
        = √(0.25 + 2.89 + 1.69 + 0.64 + 1.0 + 1.69 + 0.25)
        = √8.41
        = 2.9 hours
```

---

### Step 6: Communicate the Estimate

**Goal**: Present estimate with appropriate uncertainty.

**Never say**: "It will take 26 hours"

**Instead say**:
```
"I estimate this will take about 26 hours, give or take 3 hours.
There's about a 95% chance it will be done within 32 hours."
```

**Confidence Intervals**:
| Confidence | Formula | Example |
|------------|---------|---------|
| 68% | μ ± 1σ | 23-29 hours |
| 95% | μ ± 2σ | 20-32 hours |
| 99.7% | μ ± 3σ | 17-35 hours |

**Communicate**:
- [ ] Expected time (μ)
- [ ] Range (μ ± σ or μ ± 2σ)
- [ ] Key assumptions
- [ ] What could change the estimate

---

### Step 7: Document Assumptions

**Goal**: Make uncertainty explicit.

List what you assumed:
```
Assumptions:
- API contract is already defined
- No major refactoring needed
- Test environment is available
- No dependencies on other teams

Risks that could increase time:
- If we need to refactor the data layer: +8-16 hours
- If API contract changes: +4-8 hours
- If test environment issues: +2-4 hours
```

- [ ] Assumptions documented
- [ ] Risks identified
- [ ] Impact of risks quantified

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────┐
│                PERT ESTIMATION                       │
│                                                      │
│  For each task, estimate:                           │
│    O = Optimistic (best case)                       │
│    N = Nominal (most likely)                        │
│    P = Pessimistic (worst case)                     │
│                                                      │
│  Calculate:                                          │
│    Expected: μ = (O + 4N + P) / 6                   │
│    Std Dev:  σ = (P - O) / 6                        │
│                                                      │
│  For multiple tasks:                                │
│    μ_total = sum of all μ                           │
│    σ_total = √(sum of all σ²)                       │
│                                                      │
│  Report as range:                                   │
│    "About X hours, give or take Y hours"            │
│    or "Between X and Z hours (95% confidence)"      │
└─────────────────────────────────────────────────────┘
```

---

## Common Mistakes to Avoid

| Mistake | Why It's Bad | Do Instead |
|---------|--------------|------------|
| Single number | Creates false precision | Give a range |
| Optimistic only | You'll always be "late" | Include pessimistic |
| Padding | Dishonest, still often wrong | Use PERT honestly |
| Commitment language | "I'll have it done by..." | "I estimate..." |
| Skipping breakdown | Large tasks = large error | Break it down |

---

## Exit Criteria

Estimation is complete when:
- [ ] Task broken into subtasks (< 3 days each)
- [ ] Each subtask has O, N, P values
- [ ] Total μ and σ calculated
- [ ] Estimate communicated as a range
- [ ] Assumptions documented
- [ ] Risks identified
