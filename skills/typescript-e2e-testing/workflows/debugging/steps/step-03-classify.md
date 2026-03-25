---
name: 'step-03-classify'
description: 'Classify the failure type to guide debugging approach'
nextStepFile: './step-04-isolate.md'
---

# Step 3: Classify the Failure

## STEP GOAL

Determine the category of failure to guide debugging approach.

## EXECUTION

### Failure Categories

1. **Timeout Failures**
   - Test didn't complete in time
   - Async operation never resolved
   - Infrastructure slow/unavailable

2. **Assertion Failures**
   - Expected value doesn't match actual
   - Data not in expected state
   - Wrong error code returned

3. **Connection Failures**
   - Cannot connect to database
   - Kafka broker unavailable
   - External API unreachable

4. **Race Conditions**
   - Passes alone, fails with others
   - Intermittent failures
   - Order-dependent results

5. **State Leakage**
   - Data from previous test interferes
   - Shared resources not cleaned
   - Global state modified

### Append to Output

Append classification to the output document.

## PRESENT FINDINGS

```
Failure Classification

Based on the error, this appears to be a:

**{Category}** failure

Characteristics:
- {characteristic 1}
- {characteristic 2}

Debugging approach for this type:
1. {step 1}
2. {step 2}
3. {step 3}

[C] Continue with this approach / [R] Reclassify
```

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`

## NEXT STEP

After user confirms `[C]`, load `step-04-isolate.md`.
