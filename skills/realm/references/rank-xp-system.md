# Rank & XP System

Purpose: Define Realm rank thresholds, XP sources, decay behavior, and promotion side effects.

Contents:
- Rank and level math
- XP gain tables
- Decay rules
- Promotion behavior
- Display format

## Rank Thresholds

| Rank | XP Required | Title | Badge Color | Description |
|---|---|---|---|---|
| F | 0 | Recruit | Gray | Newly registered, minimal activity |
| E | 100 | Novice | White | Beginning to contribute |
| D | 500 | Apprentice | Green | Regular contributor |
| C | 1,500 | Journeyman | Blue | Reliable, consistent performer |
| B | 4,000 | Veteran | Purple | Experienced, proven track record |
| A | 8,000 | Elite | Gold | Top-tier performer |
| S | 15,000 | Champion | Orange | Ecosystem pillar |
| SS | 30,000 | Legend | Red | Defining presence in the ecosystem |

## Level Calculation

Levels provide finer-grained progress within ranks.

```
level = floor(sqrt(total_xp / 2))
```

Maximum display level is `99`. Agents above `19,602 XP` display `Lv.99+` until Rank `SS`.

## XP Sources

### Task Completion XP

| Task Type | Base XP | Condition |
|---|---|---|
| Single-agent task | 20 | Common quest completed |
| Multi-agent chain participant | 40 | Uncommon chain completed |
| Complex chain participant | 80 | Rare chain with 4+ agents completed |
| Parallel branch chain | 150 | Epic parallel chain completed |
| Titan product lifecycle | 300 | Legendary lifecycle delivered |

### Outcome Modifiers

| Outcome | Modifier | Condition |
|---|---|---|
| `SUCCESS` | ×1.0 | Standard completion |
| `PARTIAL` | ×0.5 | Some deliverables produced |
| `BLOCKED` | ×0.2 | Attempted but blocked |
| `FAILED` | ×0.0 | No usable output |
| First-time task type | ×1.5 | Agent's first time handling that task type |
| Streak bonus | ×1.2 | `3+` consecutive `SUCCESS` outcomes |

### Activity XP

| Activity | XP | Cap |
|---|---|---|
| Journal entry added | 10 | `5/day` |
| Pattern discovered | 25 | No cap |
| Chain coordination | 15 | Per chain |
| Review completed | 15 | Per review |
| Security finding | 20 | Per finding |

### Collaboration XP

| Activity | XP | Condition |
|---|---|---|
| New collaboration partner | 30 | First chain with that agent |
| Cross-category collaboration | 20 | Chain spans multiple categories |
| Synergy activation | 15 | A documented class synergy triggered |

## XP Decay

| Inactivity Period | Decay Rate |
|---|---|
| `0-30 days` | None |
| `31-60 days` | `-2%` total XP per week |
| `61-90 days` | `-5%` total XP per week |
| `90+ days` | `-10%` total XP per week |

Rules:
- XP never decays below the current rank minimum.
- Decay stops immediately when activity resumes.
- Decay is logged as a `Decay` event and can feed the chronicle.

## Promotion Behavior

When XP crosses a rank threshold:
1. Announce a `Celebration` event.
2. Award the corresponding rank badge.
3. Add a chronicle line for the ascension.
4. Recalculate stats and power-level presentation.

### Rank Benefits

| Rank | Unlock |
|---|---|
| E | Stat tracking enabled |
| D | Badge earning enabled |
| C | Chronicle mentions |
| B | Multi-class eligibility |
| A | Department chief eligibility |
| S | Hall of Champions listing |
| SS | Legendary status and permanent chronicle prominence |

## XP Display Format

```
XP: 4,720 / 5,000  ████████████░░ 94%
```

```
segments = 14
filled = round((current_xp - rank_min) / (rank_max - rank_min) * segments)
bar = '█' * filled + '░' * (segments - filled)
percentage = round((current_xp - rank_min) / (rank_max - rank_min) * 100)
```
