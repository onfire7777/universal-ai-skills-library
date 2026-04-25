# Chronicle Format

Purpose: Define chapter structure, era detection, narrative voice, arc detection, and retention rules for `/Realm chronicle`.

Contents:
- Chapter structure
- Era and arc detection
- Voice and vocabulary rules
- Output limits
- Persistence

## Chapter Structure

Each chapter contains:
- Chapter title and era name
- Period (`YYYY-MM ~ YYYY-MM`)
- EFS range
- Lifecycle phase
- `2-4` narrative paragraphs
- Key events
- Heroes of the era
- Statistics

## Era Detection

Start a new chapter when any of these conditions is met:

| Trigger | Era title pattern |
|---|---|
| Lifecycle phase transition | The Age of `[Phase Name]` |
| EFS grade change by 2+ grades | The `[Rise/Fall]` of `[Grade]` |
| 10+ new agents in 30 days | The Great Expansion |
| 3+ sunsets in 30 days | The Consolidation |
| First Legendary quest completed | The `[Quest Name]` Saga |
| 90 days of stable EFS A+ | The Golden Age |
| EFS drops to D or below | The Dark Times |
| Manual user trigger | The `[User-specified]` Era |

## Voice And Vocabulary

- Use third-person omniscient narration.
- Use past tense for completed events and present tense for ongoing state.
- Keep the tone formal but warm.
- Fold metrics into narrative instead of dumping raw numbers.

### Preferred Vocabulary

| Technical term | Chronicle equivalent |
|---|---|
| Agent | Hero, champion, warrior, scholar |
| Task or chain | Quest, mission, expedition |
| Bug fix | Battle, defense, vanquishing |
| Feature | Construction, creation, forging |
| Refactoring | Purification, renewal, restoration |
| Deployment | Launch, expedition departure |
| EFS increase | Prosperity, growth, strengthening |
| EFS decrease | Decline, shadow, weakening |
| New agent | Arrival, emergence, awakening |
| Rank promotion | Ascension, crowning, advancement |
| Incident | Invasion, assault, crisis |
| Recovery | Healing, restoration, renaissance |

### Class Titles

| Class | Preferred titles |
|---|---|
| Commander | General, Marshal, Grand Commander |
| Ranger | Scout Master, Chief Tracker |
| Artisan | Master Craftsman, Grand Artisan |
| Guardian | Shield Captain, Warden |
| Paladin | Holy Knight, Grand Paladin |
| Sage | Elder, Chief Sage |
| Alchemist | Grand Alchemist, Transmuter |
| Scribe | Grand Chronicler, Lorekeeper |
| Architect | Grand Architect, Master Planner |
| Enchanter | Court Enchanter, Spell Weaver |
| Engineer | Chief Engineer, Grand Mechanic |
| Merchant | Guild Master, Trade Prince |
| Oracle | High Oracle, Chief Seer |
| Herald | Royal Herald, Grand Messenger |
| Demiurge | World Shaper, Grand Demiurge |
| Strategist | Grand Strategist, War Sage |

## Story Arc Detection

| Arc Type | Detection rule |
|---|---|
| Rise | Agent XP increases by `>200%` in `60 days` |
| Fall | Agent RS drops from Active to Dormant |
| Redemption | A fall followed by a rise within `90 days` |
| Rivalry | Two agents alternate between rank #1 and #2 |
| Journey | Agent crosses `3+` rank thresholds |
| Alliance | Collaboration frequency between two departments doubles |
| Crisis | EFS drops below `55` and later recovers above `70` |

## Output Limits

| Section | Limit |
|---|---|
| Chapter narrative | 200 words |
| Key events | 10 items |
| Heroes of the era | 5 agents |
| Arc summary | 50 words |

`/Realm chronicle` shows the latest 3 chapters in full and older chapters as a table of contents.

## Persistence

- Store chronicles under `## Chronicle` in `.agents/realm-state.md`.
- Retain all chapters permanently.
