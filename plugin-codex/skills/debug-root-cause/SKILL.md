---
name: debug-root-cause
description: "Systematic root-cause debugging — reproduce, isolate, hypothesize, fix the actual cause not symptoms. Use when something is broken, behaving unexpectedly, or when a bug needs methodical investigation. Triggers on: 'debug this', 'why is this broken', 'this isn't working', 'find the bug', 'root cause', 'troubleshoot', 'figure out why', 'diagnose this error'."
---

# Debug — Root Cause Analysis

## Thinking Protocol

Before touching code, answer silently:
1. What **exactly** is the symptom? (Be precise)
2. When did it **last work**? What changed since then?
3. Am I fixing a symptom or the **root cause**?

## Execution

### Phase 1: Reproduce
- Create shortest path to trigger the bug
- Record exact steps, input data, environment
- If not reproducible, focus on gathering conditions

### Phase 2: Isolate
- Binary search the cause space:
  - `git bisect` for regressions
  - Comment out / disable components to narrow scope
  - Add targeted logging at system boundaries
- Goal: identify the single component where behavior diverges

### Phase 3: Hypothesize & Test
- Form exactly 2-3 hypotheses for root cause
- For each: design a test that would DISPROVE it (seek disconfirmation)
- Run tests, record results, eliminate until one remains

### Phase 4: Fix Root Cause
- Apply minimum change that fixes root cause
- If tempted to add a workaround, stop and dig deeper
- Verify fix resolves original reproduction case

### Phase 5: Prevent Recurrence
- Write a test that fails without the fix, passes with it
- If bug class is systemic, add lint rule or validation
- Document in commit message

## Rules
🚨 Fix-it-never-work-around-it. Band-aid fixes ARE band-aid fixes.
🚨 Don't guess. Hypothesize, then test. Can't test? Gather more data.
🚨 The bug is in your code until proven otherwise.