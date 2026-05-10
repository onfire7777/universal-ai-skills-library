# Quest Mapping

Purpose: Convert Nexus chains and task executions into Realm quests, rewards, and board layouts.

Contents:
- Quest rarity and reward multipliers
- Quest fields and naming rules
- Party bonuses
- Outcome mapping and XP distribution
- Board layout requirements

## Quest Rarity Classification

| Chain Complexity | Rarity | Color | Theme | XP Multiplier |
|---|---|---|---|---|
| Single agent | Common | `⚪` | Daily quest | ×1.0 |
| 2-3 agents | Uncommon | `🟢` | Investigation quest | ×2.0 |
| 4+ agents | Rare | `🔵` | Adventure quest | ×4.0 |
| Parallel branches | Epic | `🟣` | Expedition quest | ×7.5 |
| Titan full product | Legendary | `🟠` | Grand expedition | ×15.0 |

## Quest Fields

| Field | Source | Notes |
|---|---|---|
| Name | Auto-generated from task description | Use the templates below. |
| Rarity | Agent count and branching | Follow the rarity table above. |
| Objective | Original user request or Nexus task | Preserve intent without embellishment. |
| Party | Ordered chain agents | First = leader, middle = support, final = finisher, parallel = flankers. |
| Difficulty | Sherpa complexity `1-5` | Render as stars. |
| Status | AUTORUN markers | Active, Complete, Partial, Blocked, Failed, or Abandoned. |
| Reward | Base XP × modifiers | Split by outcome rules. |
| Duration | Activity timestamps | Estimate when exact timing is missing. |
| Badges | Badge catalog matches | Include only earned badges. |

## Quest Name Generation

| Task Pattern | Quest Name Template |
|---|---|
| Bug investigation | The `[Component]` Mystery |
| Bug fix | Vanquishing the `[Component]` Bug |
| New feature | Forging the `[Feature]` |
| Refactoring | Purifying the `[Module]` |
| Security audit | The Sentinel's Vigil: `[Target]` |
| Performance optimization | Accelerating the `[Component]` |
| Documentation | Chronicling the `[Subject]` |
| Architecture design | Blueprint of `[System]` |
| Testing | Fortifying `[Component]` Defenses |
| Deployment | The Great `[Version]` Launch |
| Migration | The `[Technology]` Exodus |
| Review | The Council Reviews `[Subject]` |

## Party Bonuses

### Party Size

| Party Size | Bonus | Narrative |
|---|---|---|
| 1 | None | Lone wolf mission |
| 2-3 | +10% XP | Small party synergy |
| 4-5 | +20% XP | Full party advantage |
| 6+ | +30% XP | Raid bonus |

### Composition

| Composition | Bonus | Condition |
|---|---|---|
| Balanced Party | +15% XP | `3+` different classes |
| Specialist Team | +10% XP | All party members share one class |
| Full Stack | +25% XP | Commander + Artisan + Guardian + Sage present |

## Quest Outcome Mapping

| AUTORUN Status | Quest Status | XP Award | Distribution |
|---|---|---|---|
| All `SUCCESS` | `✅ Complete` | 100% | Equal split among all party members |
| Mix of `SUCCESS` and `PARTIAL` | `⚠️ Partial` | 50% | Only agents with `SUCCESS` receive XP |
| Any `BLOCKED` | `🔒 Blocked` | 20% | Participation credit only |
| Any `FAILED` | `❌ Failed` | 0% | No XP, keep the record |
| No completion marker | `💀 Abandoned` | 0% | Record as a negative chronicle event |

## Quest Board Requirements

- Show active quests first.
- For active quests, display rarity, party, progress bar, and current phase/owner.
- Show completed quests from the last 7 days beneath active quests.
- Maintain lifetime quest stats by status and rarity.
- Keep completed quest history capped at the latest 50 entries in persisted state.
