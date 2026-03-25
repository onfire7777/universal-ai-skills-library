---
name: 'step-04-implement-fix'
description: 'Apply targeted fix based on identified root cause'
nextStepFile: './step-05-verify-fix.md'
---

# Step 4: Implement Fix

## STEP GOAL

Apply a minimal, targeted fix based on the root cause identified in Step 3. Fix only what is necessary -- avoid unnecessary changes.

## EXECUTION

### 1. Select Fix Based on Root Cause

Use the root cause from Step 3 to determine the appropriate fix:

| Root Cause | Fix |
|------------|-----|
| Missing mock return | Add `mockResolvedValue()` / `mockReturnValue()` |
| Wrong mock method | Fix method name or use correct mock |
| Logic error in test | Fix test arrangement or expectations |
| Logic error in target | Fix target code (create as separate task) |
| Missing mock | Add mock to provider list |
| Import error | Fix import path |

### 2. Apply Minimal Fix

Apply the smallest possible change that addresses the root cause:
- Change only the lines directly related to the root cause
- Do not refactor or improve unrelated code
- If the fix requires a target code change (not test code), flag it as a separate task

### 3. Document the Fix

Record what was changed:

```
**Fix Applied:**
- File: [path to modified file]
- Change: [description of the specific change]
- Reason: [why this fixes the root cause]
```

## PRESENT FINDINGS

Present to the user:

```
Step 4: Fix Applied
===================

File: [path]
Change: [description]
Reason: [why this addresses the root cause]

Lines Modified:
  [before -> after summary]
```

Then ask: **[C] Continue to Step 5: Verify Fix**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted`
- Append the fix details to the report

## NEXT STEP

After user confirms `[C]`, load `step-05-verify-fix.md`.
