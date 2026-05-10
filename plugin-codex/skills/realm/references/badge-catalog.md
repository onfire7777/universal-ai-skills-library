# Badge Catalog

Purpose: Define badge rarity, earn conditions, display behavior, and awarding rules for both agents and the full ecosystem.

Contents:
- Rarity tiers
- Individual badge conditions
- Ecosystem badge conditions
- Display rules
- Awarding rules

## Badge Rarity Tiers

| Rarity | Color | Typical prevalence | Description |
|---|---|---|---|
| Common | White `⚪` | ~60% of agents | Basic milestones |
| Uncommon | Green `🟢` | ~35% of agents | Notable achievements |
| Rare | Blue `🔵` | ~15% of agents | Exceptional performance |
| Epic | Purple `🟣` | ~5% of agents | Outstanding accomplishments |
| Legendary | Orange `🟠` | ~1% of agents | Once-in-a-lifetime feats |

## Individual Badges

### Activity

| Badge ID | Name | Rarity | Condition |
|---|---|---|---|
| `first_quest` | First Steps | Common | Complete first quest |
| `quest_10` | Seasoned Adventurer | Common | Complete 10 quests |
| `quest_50` | Quest Master | Uncommon | Complete 50 quests |
| `quest_100` | Centurion | Rare | Complete 100 quests |
| `quest_500` | Legendary Hero | Legendary | Complete 500 quests |
| `solo_10` | Lone Wolf | Common | Complete 10 solo quests |
| `daily_streak_7` | Weekly Warrior | Common | Active 7 consecutive days |
| `daily_streak_30` | Monthly Champion | Uncommon | Active 30 consecutive days |
| `daily_streak_90` | Seasonal Stalwart | Rare | Active 90 consecutive days |
| `daily_streak_365` | Eternal Flame | Legendary | Active 365 consecutive days |

### Collaboration

| Badge ID | Name | Rarity | Condition |
|---|---|---|---|
| `first_chain` | Party Member | Common | Join the first multi-agent chain |
| `collab_5` | Team Player | Common | Collaborate with 5 unique agents |
| `collab_15` | Social Butterfly | Uncommon | Collaborate with 15 unique agents |
| `collab_30` | Ambassador | Rare | Collaborate with 30 unique agents |
| `collab_all` | Universal Connector | Legendary | Collaborate with every agent in the ecosystem |
| `cross_category_5` | Border Crosser | Uncommon | Join chains spanning 5+ categories |
| `synergy_3` | Synergy Seeker | Uncommon | Trigger 3 class synergy bonuses |
| `synergy_10` | Synergy Master | Rare | Trigger 10 class synergy bonuses |
| `raid_party` | Raid Leader | Rare | Join a 6+ agent chain |

### Quality

| Badge ID | Name | Rarity | Condition |
|---|---|---|---|
| `perfect_10` | Flawless Run | Uncommon | 10 consecutive `SUCCESS` completions |
| `perfect_25` | Untouchable | Rare | 25 consecutive `SUCCESS` completions |
| `perfect_50` | Perfection Incarnate | Epic | 50 consecutive `SUCCESS` completions |
| `bug_slayer_5` | Bug Slayer | Common | Join 5 bug-fix quests |
| `bug_slayer_20` | Exterminator | Uncommon | Join 20 bug-fix quests |
| `zero_defects` | Zero Defect | Rare | 30 days without a `FAILED` quest |
| `recovery_3` | Comeback Kid | Uncommon | Recover from 3 blocked/failed states |
| `recovery_10` | Phoenix | Rare | Recover from 10 blocked/failed states |

### Growth

| Badge ID | Name | Rarity | Condition |
|---|---|---|---|
| `rank_e` | Awakened | Common | Reach Rank E |
| `rank_d` | Apprentice's Mark | Common | Reach Rank D |
| `rank_c` | Journeyman's Seal | Uncommon | Reach Rank C |
| `rank_b` | Veteran's Crest | Uncommon | Reach Rank B |
| `rank_a` | Elite Insignia | Rare | Reach Rank A |
| `rank_s` | Champion's Crown | Epic | Reach Rank S |
| `rank_ss` | Legend's Mantle | Legendary | Reach Rank SS |
| `journal_10` | Reflective Mind | Common | Write 10 journal entries |
| `journal_50` | Deep Thinker | Uncommon | Write 50 journal entries |
| `multi_class` | Dual Wielder | Rare | Earn multi-class status |
| `stat_90` | Transcendent | Epic | Any stat reaches 90+ |
| `all_stats_70` | Renaissance Agent | Legendary | All 6 stats reach 70+ simultaneously |

### Special

| Badge ID | Name | Rarity | Condition |
|---|---|---|---|
| `first_legendary` | Legendary First | Epic | Complete the first Legendary quest |
| `speed_run` | Speed Demon | Uncommon | Complete a quest in under 5 minutes |
| `marathon` | Marathon Runner | Uncommon | Join a quest lasting 2+ hours |
| `night_watch` | Night Owl | Uncommon | Complete a quest between 00:00-05:00 |
| `department_chief` | Department Chief | Rare | Serve as department chief |
| `mentor` | Mentor | Rare | Join a chain with 3+ lower-ranked agents |
| `dark_age_survivor` | Dark Age Survivor | Epic | Contribute during EFS < 40 and help recovery |
| `ecosystem_founder` | Founding Agent | Legendary | Be one of the first 10 agents |

## Ecosystem Badges

| Badge ID | Name | Rarity | Condition |
|---|---|---|---|
| `eco_genesis` | Company Founded | Common | 10+ active agents |
| `eco_expansion` | Great Expansion | Uncommon | 30+ agents |
| `eco_empire` | Mighty Empire | Rare | 60+ agents |
| `eco_efs_a` | Thriving Company | Uncommon | EFS reaches 85+ |
| `eco_efs_s` | Golden Age | Epic | EFS reaches 95+ |
| `eco_quest_100` | Century of Quests | Uncommon | 100 total quests completed |
| `eco_quest_1000` | Thousand Tales | Epic | 1000 total quests completed |
| `eco_all_active` | Full Roster | Rare | All agents have RS > 60 simultaneously |
| `eco_no_sunset` | Enduring Company | Rare | No sunset in 6 months |
| `eco_collaboration` | United Realm | Uncommon | Every category joins at least one chain |
| `eco_legendary_3` | Hall of Legends | Epic | 3+ agents reach Rank SS |
| `eco_recovery` | Phoenix Company | Rare | EFS recovers from D to B+ within 30 days |
| `eco_diversity` | Rainbow Coalition | Uncommon | 15+ categories have active agents |

## Display Rules

- Character sheets show earned badges inline and can surface the top earned set first.
- Badge catalog views show:
  - Earned badges
  - Next three closest badges
  - Progress toward each closest badge
- Ecosystem badges appear on the company dashboard.

## Awarding Rules

1. Re-evaluate badge conditions on every `/Realm` invocation.
2. Announce new badges through the event system.
3. Persist awarded badges with their earn dates in `.agents/realm-state.md`.
4. Never revoke badges once earned.
5. Track the top 3 closest badges for each agent.
