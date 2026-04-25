# Stat Calculation

Purpose: Define the six Realm stats, their formulas, rendering rules, and cold-start behavior.

Contents:
- Stat formulas and data sources
- Aggregation and profile types
- Update frequency
- Cold-start handling

## Overview

Each agent has six stats scored `0-100`. Realm derives them from observable ecosystem data only.

## Stat Formulas

| Stat | Meaning | Data sources | Formula | Key adjustments |
|---|---|---|---|---|
| `STR` | Output power | `PROJECT.md`, commit contributions, artifact volume | `activity_count = entries in last 90d` → `normalized = min(activity_count / median_activity_all_agents, 2.0)` → `STR = round(normalized * 50)` | Non-code categories use report or analysis counts. Agents younger than 30 days get `×1.2` to reduce cold-start penalty. |
| `DEX` | Versatility | Distinct task types, chain patterns, cross-category work | `DEX = round((task_types / max_types) * 100)` | Specialists normalize against category peers. Minimum `DEX = 10` for any active agent. |
| `INT` | Complexity handling | Sherpa complexity, chain depth, Nexus contribution | `avg_complexity = avg complexity (1-5)`; `chain_depth_bonus = avg chain position * 5`; `INT = round(min((avg_complexity / 5) * 75 + chain_depth_bonus, 100))` | If Sherpa data is missing, estimate complexity from chain length. Meta agents (`Darwin`, `Realm`, `Sigil`) gain `+15 INT`. |
| `WIS` | Learning rate | Journal growth, pattern reuse, Lore references | `growth_rate = recent_30d_entries / max(previous_30d_entries, 1)`; `base_wis = min(growth_rate * 30, 60)`; `lore_bonus = 5 * metapattern_refs` capped by total `100` | Agents younger than 60 days with more than 5 entries get `+20`. Stale journals (`>60 days`) cap `WIS` at `30`. |
| `CHA` | Collaboration | Chain participation, unique partners, completion with others | `partner_ratio = unique_partners / total_agents`; `CHA = round(min(collab_count * 2 + partner_ratio * 60, 100))` | Commander-class agents gain `+15 CHA`. Solo specialists have a floor of `20`. |
| `CON` | Reliability | AUTORUN outcomes, recovery count, Darwin RS | `success_rate = SUCCESS / total_tasks`; `recovery_bonus = recovery_count * 5` up to `20`; `rs_factor = Darwin_RS / 100 * 20`; `CON = round(min(success_rate * 60 + recovery_bonus + rs_factor, 100))` | Agents with fewer than 5 tasks use the ecosystem average as a baseline. Perfect success over 20+ tasks gains `+10`. |

## Bar Rendering

### Stat Bars

```
segments = round(stat_value / 10)
filled = '█' * segments
empty = '░' * (10 - segments)
bar = '[' + filled + empty + ']'
```

Example: `82 -> [████████░░]`

### Power Level

```
Power = (STR*20 + DEX*15 + INT*20 + WIS*15 + CHA*15 + CON*15) / 100
```

## Profile Types

| Profile | Condition | Label |
|---|---|---|
| Balanced | All stats within 20 points of mean | Well-Rounded |
| Physical | `STR + CON > 160` | Powerhouse |
| Mental | `INT + WIS > 160` | Intellectual |
| Social | `CHA + DEX > 160` | Versatile |
| Tank | `CON > 85` | Ironclad |
| Glass Cannon | `STR > 85` and `CON < 40` | Glass Cannon |

## Update Frequency

- Recalculate stats on each `/Realm` or `/Realm agent` invocation.
- Persist historical stat snapshots in `.agents/realm-state.md`.
- Show deltas as `↑`, `↓`, or `→` compared with the previous stored state.

## Cold-Start Handling

| Situation | Behavior |
|---|---|
| No activity data | Show all stats as `??` and bars as `[??????????]`. |
| Partial data (`<5` tasks) | Calculate only supported stats; leave missing stats as `??`. |
| New agent (`<14 days`) | Mark stats as provisional with a `*` suffix. |
