# Event System

Purpose: Classify Realm events, define their triggers and priorities, and preserve the narrative templates used in dashboards, maps, and chronicles.

Contents:
- Event types and triggers
- Narrative template rules
- Promotion flavor text
- Priority and retention rules

## Event Types

| Event Type | Icon | Trigger | Severity |
|---|---|---|---|
| Battle | `⚔` | Incident response chain activated | High |
| Boss Fight | `🐉` | Quest difficulty `★★★★★` | High |
| Victory | `✅` | Quest completed successfully | Normal |
| Celebration | `🎉` | Rank promotion or milestone | Normal |
| Achievement | `🏅` | Badge earned | Normal |
| Discovery | `🔍` | Lore identifies a new pattern | Normal |
| Arrival | `🌟` | New agent added to the ecosystem | Normal |
| Departure | `🌅` | Darwin recommends sunset (`RS < 20`) | Low |
| Dark Age | `🌑` | EFS drops below `40` | Critical |
| Renaissance | `🌅` | EFS recovers from grade D to B+ within `30 days` | High |
| Rotation | `🔄` | Department chief changes | Low |
| Expedition | `🗺️` | Legendary quest begins | High |
| Decay | `📉` | XP decays from inactivity | Low |

## Narrative Templates

| Event | Required fields | Narrative rule |
|---|---|---|
| Battle | Incident name/ID, threat type, affected module, participating agents, outcome, duration | Frame as an attack or defense of the realm. |
| Boss Fight | Quest name, challenge summary, party, current phase or completion state | Use when Sherpa complexity is 5/5. |
| Victory | Quest name, accomplishment summary, XP amount, party size, earned badges | Keep accomplishment clear and concise. |
| Celebration | Agent, new rank, rank title, power level | Use class-appropriate flavor text from the table below. |
| Achievement | Agent, badge name, rarity, badge meaning | Add special praise for Epic or Legendary badges. |
| Discovery | Pattern name, discovering agent, description, affected agents | State how the discovery propagates across the ecosystem. |
| Arrival | Agent name, class, department, motto/signature line | Introduce the new hero and updated roster size. |
| Departure | Agent name, service length, last active date, final rank | Use respectful sunset language. |
| Dark Age | EFS score, previous score, contributing factors | Always show first and highlight as critical. |
| Renaissance | Prior low point, current EFS, key contributors, duration | Frame as recovery after crisis. |
| Rotation | Department, new chief, previous chief, new chief rank and XP | Explain the succession. |
| Expedition | Quest name, expected party, lifecycle scope | Reserve for Legendary quests. |
| Decay | Agent, inactivity duration, XP delta | Mention if a lower-rank warning is approaching. |

## Rank Promotion Flavor

| Promotion | Flavor |
|---|---|
| `F -> E` | awakens to their calling |
| `E -> D` | begins their training in earnest |
| `D -> C` | earns their place among the skilled |
| `C -> B` | proves their mettle in countless battles |
| `B -> A` | joins the company's top performers |
| `A -> S` | is crowned champion of the realm |
| `S -> SS` | transcends mortal limits and becomes a legend |

## Event Priority And Display

| Priority | Events | Display rule |
|---|---|---|
| Critical | Dark Age | Always show first and highlight |
| High | Battle, Boss Fight, Renaissance, Expedition | Place at the top of the relevant daily section |
| Normal | Victory, Celebration, Achievement, Discovery, Arrival | Standard reverse-chronological order |
| Low | Departure, Rotation, Decay | Collapse by default when space is limited |

## Persistence

- Persist events to the `Events` section of `.agents/realm-state.md`.
- Keep only the latest 100 events in state.
- Archive older events into the chronicle instead of dropping them entirely.
