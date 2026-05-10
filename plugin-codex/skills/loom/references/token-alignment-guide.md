# Token Alignment Guide

Purpose: Use this reference when Loom must compare code tokens and Figma Variables, prioritize drift, or produce a remediation report.

## Contents
- Audit axes
- Source inventory
- Diff categories
- Priority model
- Reporting format

## Audit Axes

Run the audit across four axes:

| Axis | Meaning | Priority |
|---|---|---|
| `Name` | token name alignment | High |
| `Value` | actual value alignment | High |
| `Semantics` | intended usage alignment | Medium |
| `Hierarchy` | group, collection, or mode alignment | Low |

## Source Inventory

Inventory both sides before comparing:
- code tokens
- generated token artifacts
- Figma Variables
- Figma collections and modes
- package-level aliases or transformed names

## Diff Categories

| Category | Code | Meaning | Priority |
|---|---|---|---|
| `MISSING_IN_FIGMA` | `MIF` | exists in code, missing in Figma | High |
| `MISSING_IN_CODE` | `MIC` | exists in Figma, missing in code | Medium |
| `VALUE_MISMATCH` | `VM` | same token, different value | High |
| `NAME_DRIFT` | `ND` | different name, same intent | Medium |
| `SEMANTIC_GAP` | `SG` | same primitive, different semantic usage | Medium |
| `ORPHANED` | `OR` | exists on one side with no likely match | Low |
| `STRUCTURE_DIFF` | `SD` | collection, mode, or hierarchy mismatch | Low |

## Priority Model

Score each issue by:
- `Impact` (`1-3`)
- `Frequency` (`1-3`)
- `Visibility` (`1-3`)

Priority bands:

| Score | Priority | Action |
|---|---|---|
| `18-27` | `P0` | block delivery or escalate immediately |
| `9-17` | `P1` | deliver with an explicit remediation plan |
| `4-8` | `P2` | fix in next update cycle |
| `1-3` | `P3` | record only |

## Reporting Format

```md
# Token Alignment Report

## Scope
- Sources compared
- Collections or files covered

## Summary
- Alignment rate: [XX%]
- Critical issues: [count]

## Findings
| Token | Figma variable | Type | Priority | Recommendation |
|---|---|---|---|---|
| [token] | [variable] | [VM/MIF/ND/...] | [P0-P3] | [action] |

## Guidelines Impact
- Rules to add or change
- Rules to delete

## Routing
- Muse follow-up
- Frame follow-up
- Loom-only fixes
```

## Blocking Conditions

Escalate before final delivery when:
- a `P0` mismatch changes brand, accessibility, or system behavior
- mode or collection structure makes the Guidelines misleading
- token ownership is unclear between code and Figma
