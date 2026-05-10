---
name: 'step-03-naming'
description: 'Review naming conventions — intent, consistency, searchability'
nextStepFile: './step-04-class-design.md'
referenceFiles:
  - 'references/naming/rules.md'
---

# Step 3: Check Naming

## STEP GOAL

Review all names (variables, functions, classes) in the target code for intent-revealing quality, consistency, and searchability.

## REFERENCE LOADING

Before starting analysis, load and read:
- `references/naming/rules.md` — naming rules and conventions

Cite specific rules when reporting findings.

## ANALYSIS PROCESS

Check all names in the target code:

1. **Variables**: Do they reveal intent without comments?
2. **Functions**: Are they verb phrases that describe the action?
3. **Classes**: Are they noun phrases that describe responsibility?
4. **No Encodings**: No Hungarian notation or prefixes?
5. **Searchable**: Can you grep for important names?
6. **Consistent**: Same concept = same word throughout?

### Red Flags

Watch for and report:
- Single-letter variables (except loop counters)
- Abbreviations that aren't universal
- Names that require comments to explain
- Different words for the same concept

## PRESENT FINDINGS

Present findings to the user in this format:

```
Step 3: Naming
==============

[PASS/ISSUE] name (file:line)
  - Problem: description
  - Suggestion: better_name
  Rule: naming/rules.md — Rule N

Summary: N names reviewed, N issues found
```

Then ask: **[C] Continue to Step 4: Class/Module Design**

## FRONTMATTER UPDATE

Update the output document:
- Add `3` to `stepsCompleted`
- Append the findings section to the report

## NEXT STEP

After user confirms `[C]`, load `step-04-class-design.md`.
