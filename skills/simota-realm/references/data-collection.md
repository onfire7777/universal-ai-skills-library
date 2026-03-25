# Data Collection

Purpose: Define Realm's source inventory, collection order, persistent state schema, cold-start rules, and retention limits.

Contents:
- Source inventory
- Collection flow and freshness thresholds
- `.agents/realm-state.md` schema
- Cold-start behavior
- Integrity and performance rules

## Source Inventory

### Primary Sources

| Source | Path | Used for |
|---|---|---|
| Activity log | `.agents/PROJECT.md` | STR, quest history, XP, activity counts |
| Ecosystem state | `.agents/ECOSYSTEM.md` | EFS, RS, lifecycle phase, department health |
| Agent journals | `.agents/{name}.md` | WIS, learning rate, insight volume |
| Git history | `git log` | Activity trends, era shifts, live activity feed |
| Agent catalog | `architect/references/agent-categories.md` | Roster and class assignment |

### Secondary Sources

| Source | Path | Used for |
|---|---|---|
| Nexus chains | AUTORUN markers and logs | Quest composition, CHA, party data |
| Sherpa tasks | decomposition data | INT, quest difficulty |
| Lore patterns | `.agents/METAPATTERNS.md` | WIS bonuses, discovery events |
| Darwin scores | `.agents/ECOSYSTEM.md` | CON and health calculations |
| Retain patterns | `retain/references/gamification.md` | Optional gamification overlays |

## Collection Flow

1. Read `.agents/PROJECT.md` and count activity entries by agent.
2. Read `.agents/ECOSYSTEM.md` and extract EFS, RS, and lifecycle phase.
3. Scan `.agents/*.md` journals and count recent entries.
4. Read recent git history, capped at the last 100 entries.
5. Read the agent catalog and refresh roster/category mappings.
6. Parse AUTORUN markers when quest or chain details are needed.

## Freshness Thresholds

| Source | Refresh | Warning threshold |
|---|---|---|
| `PROJECT.md` | Every invocation | Stale after `>7 days` |
| `ECOSYSTEM.md` | Every invocation | Stale after `>14 days` |
| Agent journals | Every invocation | No staleness warning |
| Git history | Every invocation | Always fresh |
| Agent catalog | `/Realm` and `/Realm map` | Roster drift warning after `>30 days` |

## Persistent State Schema

Persist all Realm state to `.agents/realm-state.md`.

| Section | Required fields | Retention |
|---|---|---|
| `Agents` | Agent, Class, Rank, XP, Level, STR/DEX/INT/WIS/CHA/CON, Badge count | Full roster |
| `Active Quests` | ID, Name, Rarity, Party, Progress, Started | Current only |
| `Completed Quests` | ID, Name, Rarity, Party, XP, Completed | Keep last 50 |
| `Badges Earned` | Agent, Badge ID, Badge Name, Earned Date | Permanent |
| `Events` | Date, Type, Summary, Agents | Keep last 100 |
| `Department Health` | Department, Chief, Health, Status | Latest snapshot |
| `Ecosystem Badges` | Badge ID, Badge Name, Earned Date | Permanent |
| `Chronicle` | Chapters with era title and narrative | Permanent |

## Cold Start

When `.agents/realm-state.md` does not exist:
1. Create the file with empty sections.
2. Initialize all known agents with their mapped classes.
3. Set rank to `F`, XP to `0`, and unknown stats to `??`.
4. Backfill from historical activity where data exists.
5. Generate the first chronicle chapter: `Chapter I: The Beginning`.

## Integrity Rules

1. Never fabricate missing data. Use `??` or `N/A`.
2. Never delete earned badges or chronicle chapters.
3. Timestamp every persisted update.
4. Make updates idempotent within one session.
5. Use cached values with a warning when a source is stale.
6. Every displayed metric must trace back to a source or a documented Realm formula.

## Performance Limits

- Cache file reads within a session.
- Limit git history to the most recent 100 entries.
- Count journal entries; do not parse full journals unless needed.
- Prune completed quests beyond 50 and events beyond 100.
- Delegate graph-heavy rendering to Canvas instead of expanding Realm logic.
