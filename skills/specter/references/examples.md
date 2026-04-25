# Usage Examples & Report Samples

Purpose: preserve must-keep invocation examples, canonical report structure, confidence wording, and AUTORUN output shape.

## Contents

1. Invocation examples
2. Report templates
3. Confidence wording
4. False-positive notes
5. AUTORUN output

## Invocation Examples

### Example 1: Vague Symptom

**User Input**
```text
アプリを使っていると段々重くなってきます
```

**Interpretation**
- primary hypothesis: memory leak
- secondary hypothesis: resource leak
- scan priority: `useEffect`, listeners, timers, subscriptions

### Example 2: Concurrency Data Corruption

**User Input**
```text
チェックアウト処理で時々在庫数がおかしくなります。同時に複数のユーザーが注文すると起きるようです。
```

**Interpretation**
- primary hypothesis: race condition
- secondary hypothesis: non-atomic read-modify-write
- scan priority: checkout flow, inventory service, concurrent updates

### Example 3: Full Scan

**User Input**
```text
/specter scan全体
```

Use a full-category scan and summarize findings by ghost type and severity.

## Report Templates

### Standard Report

```markdown
## Specter Detection Report

### Summary
**Ghost Category:** [Category]
**Issues Found:** X CRITICAL, Y HIGH, Z MEDIUM, W LOW
**Confidence:** [HIGH/MEDIUM/LOW]
**Scan Scope:** [Files or subsystems]

### Critical Issues
#### SPECTER-001: [Title]
**Location:** `path/to/file.ts:123`
**Risk Score:** 8.7/10 (CRITICAL)
**Category:** [Ghost type]
**Detection Pattern:** [Pattern ID]
**Evidence:** [Why this is real]

**Bad:**
```code
// current problematic code
```

**Good:**
```code
// remediation example
```

**Risk Breakdown:**
| Dimension | Score | Rationale |
|-----------|-------|-----------|
| Detectability | X | ... |
| Impact | X | ... |
| Frequency | X | ... |
| Recovery | X | ... |
| Data Risk | X | ... |

**Suggested Tests:**
- [test 1]
- [test 2]

### Recommendations
1. [Priority fix order]

### False Positive Notes
- [If any]
```

### Full Scan Summary

```markdown
## Specter Full Scan Report

### Statistics
| Category | CRITICAL | HIGH | MEDIUM | LOW |
|----------|----------|------|--------|-----|
| Memory Leaks | X | X | X | X |
| Race Conditions | X | X | X | X |
| Resource Leaks | X | X | X | X |
| Async Issues | X | X | X | X |
```

## Confidence Wording

| Level | Use when |
|-------|----------|
| `HIGH` | pattern and context both confirm |
| `MEDIUM` | pattern matches but architecture or call-site handling is unclear |
| `LOW` | likely intentional or framework-managed pattern |

## False-Positive Notes

Use explicit notes for:
- intentional fire-and-forget
- framework-managed cleanup
- global error-boundary handling

## AUTORUN Output

```yaml
_STEP_COMPLETE:
  Agent: Specter
  Status: SUCCESS
  Output:
    ghost_category: Memory Leak
    issues_found:
      critical: 2
      high: 3
      medium: 5
      low: 1
    confidence: HIGH
    top_issue:
      location: src/components/Modal.tsx:45
      risk_score: 8.7
      category: Event Listener Leak
  Next: Builder
```
